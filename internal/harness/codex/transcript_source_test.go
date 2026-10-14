package codex

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// The acceptance rollout internal/timeline ingests (codex-cli 0.154).
const codexAcceptanceRollout = "../../timeline/testdata/codex-0.154-rollout.jsonl"

// A rollout that only grew is read from where the prior batch stopped, and
// the batch that comes back is the one a whole read gives: usage, offsets,
// prompt attribution, cursor and usage turns alike.
func TestCodexReadOfAGrownRolloutMatchesAWholeRead(t *testing.T) {
	data, err := os.ReadFile(codexAcceptanceRollout)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "rollout.jsonl")
	// Stop in the middle of a source record: the tail read must resume from
	// ConsumedBytes, keeping the fragment and the cumulative usage state.
	cut := len(data) / 2
	if err := os.WriteFile(path, data[:cut], 0o600); err != nil {
		t.Fatal(err)
	}
	src := codexTranscripts{}
	first, err := src.Read(harness.TranscriptReadRequest{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	if first.ConsumedBytes >= int64(cut) {
		t.Fatal("the cut did not leave a partial line")
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := src.Read(harness.TranscriptReadRequest{Path: path, Prior: first, Grown: true})
	if err != nil {
		t.Fatal(err)
	}
	want, err := src.Read(harness.TranscriptReadRequest{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	if len(want.UsageTurns) == 0 {
		t.Fatal("the whole read priced no turn; the fixture is not exercised")
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatal("the tail read differs from the whole read")
	}
}

// A rollout said to have grown that is in fact shorter than the prior batch
// read - stat saw it grow, then it was truncated - is read whole again.
func TestCodexReadOfARolloutTruncatedAfterGrowingRereadsIt(t *testing.T) {
	data, err := os.ReadFile(codexAcceptanceRollout)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "rollout.jsonl")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	src := codexTranscripts{}
	full, err := src.Read(harness.TranscriptReadRequest{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	cut := len(data) / 2
	if err := os.WriteFile(path, data[:cut], 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := src.Read(harness.TranscriptReadRequest{Path: path, Prior: full, Grown: true})
	if err != nil {
		t.Fatal(err)
	}
	if got.TotalBytes != int64(cut) {
		t.Fatalf("did not reread the truncated rollout: TotalBytes %d, want %d", got.TotalBytes, cut)
	}
}

// Codex reports a running total; its usage turns are per-call deltas with
// cache reads taken out of input, and the sample keeps the running total.
func TestCodexUsageTurnsAreDeltasNetOfCacheReads(t *testing.T) {
	src := codexTranscripts{}
	b, err := src.Read(harness.TranscriptReadRequest{Path: codexAcceptanceRollout})
	if err != nil {
		t.Fatal(err)
	}
	snaps := codexSnapshots(b)
	if len(b.UsageTurns) != len(snaps) || len(snaps) == 0 {
		t.Fatalf("%d usage turns from %d snapshots", len(b.UsageTurns), len(snaps))
	}
	var sum harness.TokenUsage
	for i, turn := range b.UsageTurns {
		snap := snaps[i]
		if turn.Usage.Input != snap.Delta.Input-snap.Delta.CacheRead || turn.Usage.CacheRead != snap.Delta.CacheRead {
			t.Errorf("turn %d: input %d cache-read %d, want %d net of %d", i, turn.Usage.Input, turn.Usage.CacheRead, snap.Delta.Input, snap.Delta.CacheRead)
		}
		if !turn.Sample.Cumulative || turn.Sample.Usage.Input != snap.Cumulative.Input || !turn.Sample.At.Equal(snap.OccurredAt) {
			t.Errorf("turn %d: sample %+v, want the cumulative snapshot at %s", i, turn.Sample, snap.OccurredAt)
		}
		sum.Input += turn.Usage.Input + turn.Usage.CacheRead
		sum.Output += turn.Usage.Output
	}
	last := snaps[len(snaps)-1].Cumulative
	if sum.Input != last.Input || sum.Output != last.Output {
		t.Errorf("deltas sum to input %d output %d, want the last running total %d / %d", sum.Input, sum.Output, last.Input, last.Output)
	}
}

// The stored telemetry cursor is where a Codex telemetry read starts, and a
// cursor at the end of the file reads nothing.
func TestCodexTelemetryReadsFromTheCursor(t *testing.T) {
	data, err := os.ReadFile(codexAcceptanceRollout)
	if err != nil {
		t.Fatal(err)
	}
	src := codexTranscripts{}
	size := int64(len(data))
	whole, err := src.Telemetry(harness.TelemetryRequest{Path: codexAcceptanceRollout, Size: size})
	if err != nil {
		t.Fatal(err)
	}
	if len(whole.Facts) == 0 || whole.Offset != size {
		t.Fatalf("a first read found %d facts and stopped at %d of %d", len(whole.Facts), whole.Offset, size)
	}
	again, err := src.Telemetry(harness.TelemetryRequest{Path: codexAcceptanceRollout, Size: size, Offset: whole.Offset, State: whole.State})
	if err != nil {
		t.Fatal(err)
	}
	if len(again.Facts) != 0 || again.Offset != size {
		t.Fatalf("a read at the end found %d facts and moved to %d", len(again.Facts), again.Offset)
	}
	if len(whole.Gaps) == 0 {
		t.Fatal("no telemetry gap declared")
	}
}

// Locate finds a rollout by the runtime's session first, then by the id the
// meta records, then by adoption; a request with nothing to go on waits.
func TestCodexLocateOrder(t *testing.T) {
	root := t.TempDir()
	day := filepath.Join(root, "2026", "09", "19")
	if err := os.MkdirAll(day, 0o755); err != nil {
		t.Fatal(err)
	}
	const fromRuntime, fromMeta = "01a0aaaa-0000-7000-8000-000000000001", "01a0aaaa-0000-7000-8000-000000000002"
	for _, id := range []string{fromRuntime, fromMeta} {
		if err := os.WriteFile(filepath.Join(day, "rollout-2026-09-19T10-44-09-"+id+".jsonl"), []byte("{}\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	src := codexTranscripts{}
	asked := 0
	runtime := func() string { asked++; return fromRuntime }

	loc, reason := src.Locate(harness.TranscriptLocateRequest{Root: root, SessionID: fromMeta, RuntimeSession: runtime})
	if reason != "" || loc.SessionID != fromRuntime || loc.Rule != codexRuleHerdrSession || asked != 1 {
		t.Fatalf("with a runtime session: %+v %q (asked %d)", loc, reason, asked)
	}
	loc, reason = src.Locate(harness.TranscriptLocateRequest{Root: root, SessionID: fromMeta})
	if reason != "" || filepath.Base(loc.Path) != "rollout-2026-09-19T10-44-09-"+fromMeta+".jsonl" || loc.SessionID != "" {
		t.Fatalf("with the meta's id: %+v %q", loc, reason)
	}
	if _, reason := src.Locate(harness.TranscriptLocateRequest{Root: root, Cwd: t.TempDir(), LaunchedAt: time.Now()}); reason != codexAdoptPending {
		t.Fatalf("with nothing to go on: reason %q, want %q", reason, codexAdoptPending)
	}
}

// canonicalTempDir is a fresh directory by the path Locate canonicalises
// a cwd to, so a rollout's session_meta cwd can match it.
func canonicalTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

// Adoption reads the whole first record, however large: Codex's
// session_meta carries its base instructions and was already 18-23 KiB on
// 2026-10-14, so a 64 KiB head is one prompt revision from losing it.
func TestCodexLocateAdoptsARolloutWhoseSessionMetaIsOver64KiB(t *testing.T) {
	t.Parallel()
	root, cwd := t.TempDir(), canonicalTempDir(t)
	launch := time.Date(2026, 10, 14, 9, 0, 0, 0, time.Local)
	const id = "01a0bbbb-0000-7000-8000-000000000001"
	path := writeRollout(t, root, launch, id, cwd, launch.Add(time.Second))
	at := launch.Add(time.Second).UTC().Format(time.RFC3339Nano)
	line := fmt.Sprintf(`{"timestamp":%q,"type":"session_meta","payload":{"session_id":%q,"id":%q,"timestamp":%q,"cwd":%q,"cli_version":"0.154.0","model_provider":"openai","base_instructions":{"text":%q}}}`+"\n",
		at, id, id, at, cwd, strings.Repeat("x", 100*1024))
	if err := os.WriteFile(path, []byte(line), 0o644); err != nil {
		t.Fatal(err)
	}

	loc, reason := codexTranscripts{}.Locate(harness.TranscriptLocateRequest{Root: root, Cwd: cwd, LaunchedAt: launch})
	if reason != "" || loc.Path != path || loc.SessionID != id || loc.Rule != codexRuleAdopt {
		t.Fatalf("Locate = %+v %q, want %s adopted", loc, reason, path)
	}
}

// Adoption walks only the days a rollout of this launch can be filed
// under: a rollout in an older day directory is never read, even one that
// would otherwise make the match ambiguous.
func TestCodexLocateAdoptsOnlyFromTheLaunchDays(t *testing.T) {
	t.Parallel()
	root, cwd := t.TempDir(), canonicalTempDir(t)
	launch := time.Date(2026, 10, 14, 9, 0, 0, 0, time.Local)
	writeRollout(t, root, launch.AddDate(0, 0, -3), "01a0bbbb-0000-7000-8000-000000000002", cwd, launch.Add(time.Minute))
	want := writeRollout(t, root, launch, "01a0bbbb-0000-7000-8000-000000000003", cwd, launch.Add(time.Second))

	loc, reason := codexTranscripts{}.Locate(harness.TranscriptLocateRequest{Root: root, Cwd: cwd, LaunchedAt: launch})
	if reason != "" || loc.Path != want || loc.Rule != codexRuleAdopt {
		t.Fatalf("Locate = %+v %q, want today's rollout %s", loc, reason, want)
	}
}

// Codex files a rollout under the local date it was created, which is
// never before the launch's local date. A launch just after local midnight
// therefore finds its rollout under that day, and the one extra day walked
// still covers a Codex whose clock zone is behind mate's and filed it under
// the previous date; two days back is not walked.
func TestCodexLocateAdoptsALaunchJustAfterLocalMidnight(t *testing.T) {
	t.Parallel()
	launch := time.Date(2026, 10, 14, 0, 0, 30, 0, time.Local)
	for _, tc := range []struct {
		name string
		day  time.Time
	}{
		{"filed under the launch date", launch},
		{"filed under the previous date", launch.AddDate(0, 0, -1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, cwd := t.TempDir(), canonicalTempDir(t)
			writeRollout(t, root, launch.AddDate(0, 0, -2), "01a0bbbb-0000-7000-8000-000000000004", cwd, launch.Add(time.Minute))
			want := writeRollout(t, root, tc.day, "01a0bbbb-0000-7000-8000-000000000005", cwd, launch.Add(time.Second))

			loc, reason := codexTranscripts{}.Locate(harness.TranscriptLocateRequest{Root: root, Cwd: cwd, LaunchedAt: launch})
			if reason != "" || loc.Path != want {
				t.Fatalf("Locate = %+v %q, want %s", loc, reason, want)
			}
		})
	}
}

// A session id finds only a file Codex named for it,
// `rollout-<time>-<id>.jsonl`; any other file carrying the id is not a
// rollout, and Locate falls through to adoption.
func TestCodexLocateByIDMatchesOnlyARolloutName(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	day := filepath.Join(root, "2026", "10", "14")
	if err := os.MkdirAll(day, 0o755); err != nil {
		t.Fatal(err)
	}
	const id = "01a0bbbb-0000-7000-8000-000000000006"
	if err := os.WriteFile(filepath.Join(day, "notes-"+id+".jsonl"), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	src := codexTranscripts{}
	if loc, reason := src.Locate(harness.TranscriptLocateRequest{Root: root, SessionID: id, Cwd: t.TempDir()}); reason != codexAdoptPending {
		t.Fatalf("by the meta's id: %+v %q, want %q", loc, reason, codexAdoptPending)
	}
	runtime := func() string { return id }
	if loc, reason := src.Locate(harness.TranscriptLocateRequest{Root: root, RuntimeSession: runtime, Cwd: t.TempDir()}); reason != codexAdoptPending {
		t.Fatalf("by the runtime's id: %+v %q, want %q", loc, reason, codexAdoptPending)
	}
}

// A sessions directory that does not exist yet has nothing to adopt.
func TestCodexLocateWithNoSessionsDirectoryWaits(t *testing.T) {
	t.Parallel()
	root := filepath.Join(t.TempDir(), "missing")
	if loc, reason := (codexTranscripts{}).Locate(harness.TranscriptLocateRequest{Root: root, Cwd: t.TempDir(), LaunchedAt: time.Now()}); reason != codexAdoptPending {
		t.Fatalf("Locate = %+v %q, want %q", loc, reason, codexAdoptPending)
	}
}

// An agent whose launch time is unknown is adopted from every day, as
// before the day bound.
func TestCodexLocateWithNoLaunchTimeAdoptsFromAnyDay(t *testing.T) {
	t.Parallel()
	root, cwd := t.TempDir(), canonicalTempDir(t)
	old := time.Date(2024, 1, 2, 9, 0, 0, 0, time.Local)
	want := writeRollout(t, root, old, "01a0bbbb-0000-7000-8000-000000000007", cwd, old)

	loc, reason := codexTranscripts{}.Locate(harness.TranscriptLocateRequest{Root: root, Cwd: cwd})
	if reason != "" || loc.Path != want {
		t.Fatalf("Locate = %+v %q, want %s", loc, reason, want)
	}
}
