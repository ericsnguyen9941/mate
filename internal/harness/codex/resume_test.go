package codex

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writeRollout writes a minimal rollout whose first record is a session_meta
// of the shape codex-cli 0.154.0 writes, under sessions/YYYY/MM/DD.
func writeRollout(t *testing.T, sessions string, day time.Time, id, cwd string, at time.Time) string {
	t.Helper()
	dir := filepath.Join(sessions, day.Format("2006"), day.Format("01"), day.Format("02"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, fmt.Sprintf("rollout-%s-%s.jsonl", day.Format("2006-01-02T15-04-05"), id))
	line := fmt.Sprintf(`{"timestamp":%q,"type":"session_meta","payload":{"session_id":%q,"id":%q,"timestamp":%q,"cwd":%q,"cli_version":"0.154.0","model_provider":"openai"}}`+"\n",
		at.UTC().Format(time.RFC3339Nano), id, id, at.UTC().Format(time.RFC3339Nano), cwd)
	if err := os.WriteFile(path, []byte(line), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCodexRolloutPathFindsTheSessionByName(t *testing.T) {
	t.Parallel()
	sessions := t.TempDir()
	now := time.Date(2026, 9, 24, 14, 45, 54, 0, time.Local)
	const id = "01a0d260-cd47-77d2-bee7-46d98aa0461a"
	want := writeRollout(t, sessions, now, id, "/w/mate", now)
	writeRollout(t, sessions, now, "01a0d264-b0af-7ca3-a044-b7fa841655b1", "/w/mate", now)

	got, ok := CodexRolloutPath(sessions, id)
	if !ok || got != want {
		t.Fatalf("CodexRolloutPath = %q, %v; want %q", got, ok, want)
	}
	if _, ok := CodexRolloutPath(sessions, "01a0d260-0000-7000-8000-000000000000"); ok {
		t.Fatal("an id with no rollout must not be found")
	}
	if _, ok := CodexRolloutPath(sessions, ""); ok {
		t.Fatal("an empty id must not match every rollout")
	}
	if _, ok := CodexRolloutPath(filepath.Join(sessions, "missing"), id); ok {
		t.Fatal("a missing sessions directory has no rollout")
	}
}

func TestCodexRolloutCandidatesSinceSkipsOlderDays(t *testing.T) {
	t.Parallel()
	sessions := t.TempDir()
	launch := time.Date(2026, 9, 24, 9, 0, 0, 0, time.Local)
	old := launch.AddDate(0, 0, -5)
	writeRollout(t, sessions, old, "01a0aaaa-0000-7000-8000-000000000001", "/w/mate", old)
	// The day before the launch is still walked: a Codex whose clock zone
	// is behind mate's files its rollout under the previous date.
	eve := launch.AddDate(0, 0, -1)
	writeRollout(t, sessions, eve, "01a0aaaa-0000-7000-8000-000000000002", "/w/mate", eve)
	today := writeRollout(t, sessions, launch, "01a0aaaa-0000-7000-8000-000000000003", "/w/mate", launch.Add(time.Second))
	// Not a rollout: never parsed, never a candidate.
	if err := os.WriteFile(filepath.Join(filepath.Dir(today), "notes.jsonl"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := CodexRolloutCandidatesSince(sessions, launch)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("candidates = %+v, want the eve and today rollouts only", got)
	}
	adopted := AdoptCodexRollout(got, "/w/mate", launch, "")
	if adopted.Status != CodexAdoptionKnown || adopted.Candidate.Path != today {
		t.Fatalf("adoption = %+v, want today's rollout", adopted)
	}
}

// Each rollout's head is read once per process (CON-255): the ingest asks
// for candidates on every poll, and a rollout's first record does not change
// once it is whole. A file that shrank is not that file any more and is read
// again; a head Codex has not finished writing is not kept.
func TestCodexRolloutCandidatesSinceReadsEachWholeHeadOnce(t *testing.T) {
	t.Parallel()
	sessions := t.TempDir()
	launch := time.Date(2026, 11, 11, 9, 0, 0, 0, time.Local)
	const id = "01a0bbbb-0000-7000-8000-000000000001"
	path := writeRollout(t, sessions, launch, id, "/w/mate", launch)
	whole, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ids := func() []string {
		t.Helper()
		got, err := CodexRolloutCandidatesSince(sessions, launch)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, c := range got {
			out = append(out, c.Meta.SessionID)
		}
		return out
	}

	if got := ids(); len(got) != 1 || got[0] != id {
		t.Fatalf("first walk = %v, want [%s]", got, id)
	}
	// Same size, head no longer parses: the meta read the first time stands.
	garbage := []byte(strings.Repeat("x", len(whole)-1) + "\n")
	if err := os.WriteFile(path, garbage, 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ids(); len(got) != 1 || got[0] != id {
		t.Fatalf("second walk = %v, want the kept meta of %s", got, id)
	}
	// Smaller than it was: another file, read again, and it does not parse.
	if err := os.WriteFile(path, []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ids(); len(got) != 0 {
		t.Fatalf("walk after the file shrank = %v, want none", got)
	}
	// A head still being written is not kept: once whole, it is a candidate.
	if err := os.WriteFile(path, whole[:len(whole)/2], 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ids(); len(got) != 0 {
		t.Fatalf("walk over a half-written head = %v, want none", got)
	}
	if err := os.WriteFile(path, whole, 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ids(); len(got) != 1 || got[0] != id {
		t.Fatalf("walk once the head is whole = %v, want [%s]", got, id)
	}
}
