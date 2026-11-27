package agy

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

const fixtureDB = "testdata/transcript/conversation.db"

// The fixture's figures are the ones `agy -p --output-format stream-json`
// printed for the same calls (testdata/README.md).
func TestAgyTranscriptReadsTheMeasuredUsage(t *testing.T) {
	b, err := agyTranscripts{}.Read(harness.TranscriptReadRequest{Path: fixtureDB})
	if err != nil {
		t.Fatal(err)
	}
	if b.Malformed != nil || len(b.UsageFailures) != 0 {
		t.Fatalf("malformed %+v, failures %+v", b.Malformed, b.UsageFailures)
	}
	want := []harness.TokenUsage{
		{Input: 15098, Output: 35},
		{Input: 2180, CacheRead: 13044, Output: 6},
		{Input: 15105, Output: 168, Reasoning: 161},
	}
	if len(b.UsageTurns) != len(want) {
		t.Fatalf("%d usage turns, want %d", len(b.UsageTurns), len(want))
	}
	for i, u := range b.UsageTurns {
		if u.Usage != want[i] {
			t.Errorf("turn %d usage %+v, want %+v", i, u.Usage, want[i])
		}
	}
	first := b.UsageTurns[0]
	if got := first.EndedAt.Sub(first.StartedAt); got < 5500*time.Millisecond || got > 5600*time.Millisecond {
		t.Errorf("first call took %v; stream-json said 5.54s", got)
	}
	if len(first.Calls) != 1 || first.Calls[0].ToolName != "read_file" || first.Calls[0].CommandClass != harness.CommandRead ||
		first.Calls[0].InputJSON != `{"Path":"/tmp/mate-repro/agyprobe/a.txt"}` {
		t.Fatalf("first call's tool calls = %+v", first.Calls)
	}
	if len(b.ToolResults) != 1 || b.ToolResults[0].CallID != first.Calls[0].SourceRef || b.ToolResults[0].IsError {
		t.Fatalf("tool results = %+v", b.ToolResults)
	}
	if b.UsageTurns[2].HarnessTurnRef != "4" || first.HarnessTurnRef != "0" {
		t.Errorf("turn refs %q, %q: each prompt's step idx", first.HarnessTurnRef, b.UsageTurns[2].HarnessTurnRef)
	}
	if b.ConsumedBytes != 6 || b.TotalBytes != 6 {
		t.Errorf("consumed %d of %d, want every step", b.ConsumedBytes, b.TotalBytes)
	}
}

// A step whose metadata is not protobuf costs only itself.
func TestAgyTranscriptSkipsAnUndecodableStep(t *testing.T) {
	steps := []step{
		{idx: 0, kind: stepResponse, status: statusDone, metadata: []byte{0xff}},
		{idx: 1, kind: 101, status: statusDone, metadata: []byte{0x08}},
		{idx: 2, kind: stepResponse, status: statusDone, metadata: []byte{0x4a, 0x04, 0x10, 0x05, 0x18, 0x02}},
	}
	b := stepsBatch(steps)
	if b.Malformed != nil || len(b.UsageFailures) != 1 || b.UsageFailures[0].Offset != 0 {
		t.Fatalf("malformed %+v, failures %+v", b.Malformed, b.UsageFailures)
	}
	if len(b.UsageTurns) != 1 || b.UsageTurns[0].Usage != (harness.TokenUsage{Input: 5, Output: 2}) {
		t.Fatalf("usage turns %+v", b.UsageTurns)
	}
}

// Locate finds the conversation agy's map names for the Crew's worktree,
// and only one written since the launch.
func TestAgyTranscriptLocatesThroughTheCwdMap(t *testing.T) {
	cwd := t.TempDir()
	resolved, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		t.Fatal(err)
	}
	home := writeHome(t, resolved, testConversation, time.Now())
	src := agyTranscripts{}
	loc, reason := src.Locate(harness.TranscriptLocateRequest{Cwd: cwd, LaunchedAt: time.Now().Add(-time.Minute), Root: home})
	if reason != "" || loc.Path != ConversationPath(home, testConversation) || loc.SessionID != testConversation {
		t.Fatalf("Locate = %+v, %q", loc, reason)
	}
	if _, reason := src.Locate(harness.TranscriptLocateRequest{Cwd: cwd, LaunchedAt: time.Now().Add(time.Hour), Root: home}); reason != harness.LocateNoSession {
		t.Fatalf("a conversation older than the launch was located: %q", reason)
	}
	if err := os.Remove(ConversationPath(home, testConversation)); err != nil {
		t.Fatal(err)
	}
	if _, reason := src.Locate(harness.TranscriptLocateRequest{Cwd: cwd, SessionID: testConversation, Root: home}); reason != harness.LocateNotFound {
		t.Fatalf("a missing database was located: %q", reason)
	}
}
