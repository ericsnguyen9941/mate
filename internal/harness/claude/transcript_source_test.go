package claude

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// claudeTailFixture is a synthetic live transcript built to stress a tail
// read: two harness turns, multi-record message groups with housekeeping and
// a tool result interleaved, a compaction, a result (toolu_1) that lands after
// the next message group has started, and a trailing group (msg_3) that a
// live read withholds.
const claudeTailFixture = `{"type":"user","uuid":"u1","timestamp":"2026-09-08T10:00:00.000Z","promptId":"p1","message":{"role":"user","content":"start"}}
{"type":"assistant","uuid":"a1","timestamp":"2026-09-08T10:00:01.000Z","message":{"id":"msg_1","model":"claude-sonnet-4","usage":{"input_tokens":2,"output_tokens":5,"cache_read_input_tokens":7,"cache_creation_input_tokens":11,"output_tokens_details":{"thinking_tokens":1}},"content":[{"type":"text","text":"one"},{"type":"tool_use","id":"toolu_1","name":"Bash","input":{"command":"sleep 9"}}]}}
{"type":"attachment","uuid":"h1","timestamp":"2026-09-08T10:00:01.050Z","attachment":{"type":"file"}}
{"type":"assistant","uuid":"a2","timestamp":"2026-09-08T10:00:01.100Z","message":{"id":"msg_1","model":"claude-sonnet-4","stop_reason":"tool_use","usage":{"input_tokens":2,"output_tokens":9,"cache_read_input_tokens":7,"cache_creation_input_tokens":11,"output_tokens_details":{"thinking_tokens":3}},"content":[{"type":"text","text":"two"}]}}
{"type":"system","uuid":"s1","timestamp":"2026-09-08T10:00:02.000Z","subtype":"compact_boundary"}
{"type":"user","uuid":"u2","timestamp":"2026-09-08T10:00:03.000Z","promptId":"p2","isCompactSummary":true,"message":{"role":"user","content":"summary of before"}}
{"type":"assistant","uuid":"b1","timestamp":"2026-09-08T10:00:04.000Z","message":{"id":"msg_2","model":"claude-sonnet-4","usage":{"input_tokens":3,"output_tokens":4,"cache_read_input_tokens":20,"cache_creation_input_tokens":0},"content":[{"type":"tool_use","id":"toolu_2","name":"Read","input":{"path":"x"}}]}}
{"type":"mode","mode":"normal","uuid":"h2","timestamp":"2026-09-08T10:00:04.050Z"}
{"type":"user","uuid":"u3","timestamp":"2026-09-08T10:00:20.000Z","promptId":"p2","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_1","content":"slept"}]}}
{"type":"assistant","uuid":"b2","timestamp":"2026-09-08T10:00:04.100Z","message":{"id":"msg_2","model":"claude-sonnet-4","stop_reason":"tool_use","usage":{"input_tokens":3,"output_tokens":8,"cache_read_input_tokens":20,"cache_creation_input_tokens":0},"content":[{"type":"text","text":"reading"}]}}

{"type":"user","uuid":"u4","timestamp":"2026-09-08T10:00:05.000Z","promptId":"p2","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_2","content":"data","is_error":true}]}}
{"type":"assistant","uuid":"c1","timestamp":"2026-09-08T10:00:06.000Z","message":{"id":"msg_3","model":"claude-sonnet-4","usage":{"input_tokens":4,"output_tokens":6,"cache_read_input_tokens":30,"cache_creation_input_tokens":1},"content":[{"type":"text","text":"done"}]}}
{"type":"assistant","uuid":"c2","timestamp":"2026-09-08T10:00:06.100Z","message":{"id":"msg_3","model":"claude-sonnet-4","stop_reason":"end_turn","usage":{"input_tokens":4,"output_tokens":7,"cache_read_input_tokens":30,"cache_creation_input_tokens":1},"content":[{"type":"text","text":"really"}]}}
`

func readClaude(t *testing.T, req harness.TranscriptReadRequest) harness.TranscriptBatch {
	t.Helper()
	b, err := claudeTranscripts{}.Read(req)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// A transcript that only grew is read from where the prior batch stopped,
// and the batch that comes back is the one a whole read gives: every field,
// at every place the prior read could have stopped - mid-line, inside a
// multi-record message group, on a blank line. The prior read's own bytes
// are then overwritten: a tail read never looks at them again, so it still
// matches the whole read of the original file. A read that re-parsed the
// whole file would see the garbage instead.
func TestClaudeReadOfAGrownTranscriptMatchesAWholeRead(t *testing.T) {
	data := []byte(claudeTailFixture)
	path := filepath.Join(t.TempDir(), "session.jsonl")
	write := func(b []byte) {
		t.Helper()
		if err := os.WriteFile(path, b, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write(data)
	want := readClaude(t, harness.TranscriptReadRequest{Path: path})
	if len(want.UsageTurns) != 2 || want.OpenSourceRef != "msg_3" || len(want.Compactions) != 2 {
		t.Fatalf("whole read: %d usage turns, open %q, %d compactions; the fixture is not exercised",
			len(want.UsageTurns), want.OpenSourceRef, len(want.Compactions))
	}
	tailed := 0
	for cut := 0; cut < len(data); cut++ {
		write(data[:cut])
		prior := readClaude(t, harness.TranscriptReadRequest{Path: path})
		write(data)
		got := readClaude(t, harness.TranscriptReadRequest{Path: path, Prior: prior, Grown: true})
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("cut at %d: the grown read differs from the whole read\n got %s\nwant %s", cut, describeBatch(got), describeBatch(want))
		}
		if prior.ConsumedBytes == 0 {
			continue
		}
		tailed++
		garbled := bytes.Clone(data)
		copy(garbled, bytes.Repeat([]byte("#"), int(prior.ConsumedBytes)))
		write(garbled)
		got = readClaude(t, harness.TranscriptReadRequest{Path: path, Prior: prior, Grown: true})
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("cut at %d: a grown read re-read the %d bytes the prior batch had already consumed\n got %s\nwant %s",
				cut, prior.ConsumedBytes, describeBatch(got), describeBatch(want))
		}
	}
	if tailed == 0 {
		t.Fatal("no cut left a prior batch to tail from")
	}
}

// A chain of grown reads carries its state from read to read: three passes
// over a growing file end where one whole read does.
func TestClaudeReadOfATranscriptGrowingAcrossSeveralReadsMatchesAWholeRead(t *testing.T) {
	data := []byte(claudeTailFixture)
	path := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	want := readClaude(t, harness.TranscriptReadRequest{Path: path})
	for first := 1; first < len(data); first += 61 {
		for second := first + 1; second < len(data); second += 97 {
			var b harness.TranscriptBatch
			for n, end := range []int{first, second, len(data)} {
				if err := os.WriteFile(path, data[:end], 0o600); err != nil {
					t.Fatal(err)
				}
				b = readClaude(t, harness.TranscriptReadRequest{Path: path, Prior: b, Grown: n > 0})
			}
			if !reflect.DeepEqual(b, want) {
				t.Fatalf("cuts %d, %d: the chained read differs from the whole read\n got %s\nwant %s", first, second, describeBatch(b), describeBatch(want))
			}
		}
	}
}

// A prior batch the tail cannot trust - one stopped by a malformed line, or
// one read at rest - is not carried on from: the file is read whole.
func TestClaudeReadOfAGrownTranscriptRereadsWhatItCannotCarryOn(t *testing.T) {
	data := []byte(claudeTailFixture)
	path := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	want := readClaude(t, harness.TranscriptReadRequest{Path: path})
	cut := bytes.Index(data, []byte(`{"type":"assistant","uuid":"b1"`))

	malformed := readClaude(t, harness.TranscriptReadRequest{Path: path})
	malformed.Malformed = &harness.TranscriptMalformed{Offset: 1, Reason: harness.MalformedInvalidJSON}
	malformed.ConsumedBytes = 1
	if got := readClaude(t, harness.TranscriptReadRequest{Path: path, Prior: malformed, Grown: true}); !reflect.DeepEqual(got, want) {
		t.Errorf("a malformed prior was carried on from\n got %s\nwant %s", describeBatch(got), describeBatch(want))
	}

	// A read at rest emits msg_2's first record as the whole message; a
	// tail from there would meet its second record as a new message.
	mid := bytes.Index(data, []byte(`{"type":"user","uuid":"u3"`))
	if err := os.WriteFile(path, data[:mid], 0o600); err != nil {
		t.Fatal(err)
	}
	final := readClaude(t, harness.TranscriptReadRequest{Path: path, AtRest: true})
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if got := readClaude(t, harness.TranscriptReadRequest{Path: path, Prior: final, Grown: true}); !reflect.DeepEqual(got, want) {
		t.Errorf("a prior read at rest was carried on from\n got %s\nwant %s", describeBatch(got), describeBatch(want))
	}

	foreign := harness.TranscriptBatch{ConsumedBytes: int64(cut)}
	if got := readClaude(t, harness.TranscriptReadRequest{Path: path, Prior: foreign, Grown: true}); !reflect.DeepEqual(got, want) {
		t.Errorf("a prior this source did not read was carried on from\n got %s\nwant %s", describeBatch(got), describeBatch(want))
	}

	longer := readClaude(t, harness.TranscriptReadRequest{Path: path})
	if err := os.WriteFile(path, data[:cut], 0o600); err != nil {
		t.Fatal(err)
	}
	shortWant := readClaude(t, harness.TranscriptReadRequest{Path: path})
	if got := readClaude(t, harness.TranscriptReadRequest{Path: path, Prior: longer, Grown: true}); !reflect.DeepEqual(got, shortWant) {
		t.Errorf("a prior that consumed more than the file holds was carried on from\n got %s\nwant %s", describeBatch(got), describeBatch(shortWant))
	}
}

func describeBatch(b harness.TranscriptBatch) string {
	return fmt.Sprintf("%+v", b)
}

// On real sessions, a read that carried on from a live read of the same
// file cut short gives the batch a whole read gives, at several cuts per
// session. Opt in with MATE_LIVE=1 and MATE_CLAUDE_PROJECTS_DIR.
func TestLiveClaudeTranscriptCorpusTailReadMatchesAWholeRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.jsonl")
	files, tails := 0, 0
	walkTranscriptCorpus(t, func(name string, data []byte) {
		files++
		write := func(b []byte) {
			t.Helper()
			if err := os.WriteFile(path, b, 0o600); err != nil {
				t.Fatal(err)
			}
		}
		write(data)
		want := readClaude(t, harness.TranscriptReadRequest{Path: path})
		for _, frac := range []int{3, 6, 7} {
			cut := len(data) * frac / 8
			write(data[:cut])
			prior := readClaude(t, harness.TranscriptReadRequest{Path: path})
			write(data)
			got := readClaude(t, harness.TranscriptReadRequest{Path: path, Prior: prior, Grown: true})
			if prior.Malformed == nil {
				tails++
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s: cut at %d of %d: the grown read differs from the whole read (turns %d vs %d, consumed %d vs %d, malformed %v vs %v)",
					name, cut, len(data), len(got.Turns), len(want.Turns), got.ConsumedBytes, want.ConsumedBytes, got.Malformed, want.Malformed)
			}
		}
	})
	t.Logf("%d session(s), %d tail read(s) compared with a whole read", files, tails)
}
