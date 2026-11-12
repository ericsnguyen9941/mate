package timeline_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// A pass reads only the projects of its scope, and within one skips what
// has not moved since a committed pass read it (CON-268).

var shopOnly = timeline.Scope{Projects: []string{fixtureProject}}

func ingestScope(t *testing.T, ing *timeline.Ingester, scope timeline.Scope) {
	t.Helper()
	if err := ing.IngestScope(context.Background(), scope); err != nil {
		t.Fatalf("IngestScope(%+v): %v", scope, err)
	}
}

// touch moves a file's mtime forward, as a write would.
func touch(t *testing.T, path string) {
	t.Helper()
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	at := fi.ModTime().Add(time.Second)
	if err := os.Chtimes(path, at, at); err != nil {
		t.Fatal(err)
	}
}

func TestAPassReadsOnlyTheProjectsOfItsScope(t *testing.T) {
	f := newFixture(t)
	ingestScope(t, f.ing, timeline.Scope{Projects: []string{"elsewhere"}})
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE project = ?`, fixtureProject); n != 0 {
		t.Fatalf("a pass whose scope leaves out %s recorded %d of its events", fixtureProject, n)
	}
	ingestScope(t, f.ing, shopOnly)
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE project = ?`, fixtureProject); n == 0 {
		t.Fatalf("a pass over %s recorded nothing of it", fixtureProject)
	}
}

// The rank of the Mate's sessions (rankSessions) stands for every write of
// its part of a pass: a pass that skips the Mate does not rank them again.
// mateSessions counts the ranked ones; dropMateSessions unranks them.
func mateSessions(t *testing.T, f *fixture) int {
	t.Helper()
	return f.count(t, `SELECT COUNT(*) FROM session WHERE actor_id = ? AND crew_session > 0`, f.mateActor())
}

func dropMateSessions(t *testing.T, f *fixture) {
	t.Helper()
	if _, err := f.db.SQL().Exec(`UPDATE session SET crew_session = 0 WHERE actor_id = ?`, f.mateActor()); err != nil {
		t.Fatal(err)
	}
}

// moveMate points the Mate's meta at a copy of its transcript the test may
// change.
func moveMate(t *testing.T, f *fixture) string {
	t.Helper()
	data, err := os.ReadFile(abs(t, claudeFixture))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "mate.jsonl")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	meta, err := f.ws.ReadMateMeta(fixtureProject)
	if err != nil {
		t.Fatal(err)
	}
	meta[timeline.MetaTranscript] = path
	if err := f.ws.WriteMateMeta(fixtureProject, meta); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAMateWhoseFilesDidNotMoveIsSkippedAndReadAgainWhenTheyDo(t *testing.T) {
	f := newFixture(t)
	transcript := moveMate(t, f)
	ingestScope(t, f.ing, shopOnly)
	ingestScope(t, f.ing, shopOnly)
	if mateSessions(t, f) == 0 {
		t.Fatal("the Mate has no session after two passes")
	}

	dropMateSessions(t, f)
	ingestScope(t, f.ing, shopOnly)
	if n := mateSessions(t, f); n != 0 {
		t.Fatalf("a pass over a Mate whose files did not move wrote %d session row(s): it was not skipped", n)
	}

	// Its transcript grew: read again.
	touch(t, transcript)
	ingestScope(t, f.ing, shopOnly)
	if mateSessions(t, f) == 0 {
		t.Fatal("a pass after the Mate's transcript moved did not read the Mate")
	}

	dropMateSessions(t, f)
	touch(t, f.ws.MateMeta(fixtureProject))
	ingestScope(t, f.ing, shopOnly)
	if mateSessions(t, f) == 0 {
		t.Fatal("a pass after mate.meta moved did not read the Mate")
	}

	// A sweep trusts no key.
	dropMateSessions(t, f)
	ingestScope(t, f.ing, shopOnly)
	if mateSessions(t, f) != 0 {
		t.Fatal("the Mate was read again though nothing moved")
	}
	ingestScope(t, f.ing, timeline.SweepScope())
	if mateSessions(t, f) == 0 {
		t.Fatal("a sweep skipped the Mate")
	}
}

// A Mate whose transcript the pass could not find is read every pass: the
// file the harness's rules look for may appear with mate.meta unchanged.
func TestAMateWithNoTranscriptFoundIsNeverSkipped(t *testing.T) {
	f := newFixture(t)
	meta, err := f.ws.ReadMateMeta(fixtureProject)
	if err != nil {
		t.Fatal(err)
	}
	meta[timeline.MetaTranscript] = filepath.Join(t.TempDir(), "not-yet.jsonl")
	if err := f.ws.WriteMateMeta(fixtureProject, meta); err != nil {
		t.Fatal(err)
	}
	ingestScope(t, f.ing, shopOnly)
	ingestScope(t, f.ing, shopOnly)
	dropMateSessions(t, f)
	ingestScope(t, f.ing, shopOnly)
	if mateSessions(t, f) == 0 {
		t.Fatal("a Mate whose transcript was not found was skipped")
	}
}

// Git is asked about an open crew only on a hint: its branch moved, it
// wrote a status line, or the sweep came.
func TestGitIsAskedAboutACrewOnlyOnAHint(t *testing.T) {
	f := newFixture(t)
	repo := filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject)
	commitOnBranch(t, repo, fixtureBranch, "first")
	ing, counter := f.countedIngester(t)

	ingestScope(t, ing, shopOnly)
	if runs := counter.take(); len(runs) == 0 {
		t.Fatal("the first pass asked git nothing: the test does not exercise git")
	}
	ingestScope(t, ing, shopOnly)
	if runs := counter.take(); len(runs) != 0 {
		t.Fatalf("a pass with no hint ran %d git command(s): %q", len(runs), runs)
	}

	// The branch moved: asked, and the new commit recorded.
	_, second := commitOnBranch(t, repo, fixtureBranch, "second")
	ingestScope(t, ing, shopOnly)
	if runs := counter.take(); len(runs) == 0 {
		t.Fatal("a pass after the branch moved asked git nothing")
	}
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ? AND json_extract(payload, '$.sha') = ?`,
		timeline.KindGitCommitted, second); n != 1 {
		t.Fatal("the commit made after the first passes was not recorded")
	}
	ingestScope(t, ing, shopOnly)
	if runs := counter.take(); len(runs) != 0 {
		t.Fatalf("a pass with no hint ran %d git command(s): %q", len(runs), runs)
	}

	// A status line.
	if err := f.ws.AppendStatus(fixtureProject, fixtureCrew, "working: one more"); err != nil {
		t.Fatal(err)
	}
	ingestScope(t, ing, shopOnly)
	if runs := counter.take(); len(runs) == 0 {
		t.Fatal("a pass after a status line asked git nothing")
	}

	// The sweep.
	ingestScope(t, ing, shopOnly)
	counter.take()
	ingestScope(t, ing, timeline.SweepScope())
	if runs := counter.take(); len(runs) == 0 {
		t.Fatal("a sweep asked git nothing")
	}
}

// Reindex forgets every memo and reads everything: what the skipping
// passes left on the table is what a rebuild gives.
func TestReindexReadsWhatThePassesSkipped(t *testing.T) {
	f := newFixture(t)
	moveMate(t, f)
	for range 3 {
		ingestScope(t, f.ing, shopOnly)
	}
	events := f.count(t, `SELECT COUNT(*) FROM event`)
	turns := f.count(t, `SELECT COUNT(*) FROM turn`)
	if err := f.ing.Reindex(context.Background()); err != nil {
		t.Fatalf("Reindex: %v", err)
	}
	if got := f.count(t, `SELECT COUNT(*) FROM event`); got != events {
		t.Fatalf("reindex recorded %d events, the passes %d", got, events)
	}
	if got := f.count(t, `SELECT COUNT(*) FROM turn`); got != turns {
		t.Fatalf("reindex recorded %d turns, the passes %d", got, turns)
	}
	if mateSessions(t, f) == 0 {
		t.Fatal("reindex did not read the Mate")
	}
	// And the next pass after it reads the Mate again rather than trust a
	// memo of tables the rebuild emptied.
	dropMateSessions(t, f)
	ingestScope(t, f.ing, shopOnly)
	if mateSessions(t, f) == 0 {
		t.Fatal("the first pass after a reindex skipped the Mate")
	}
}

// The no-turn check reads the projects of its scope, from the crews the
// round's pass read.
func TestCheckNoTurnReadsTheProjectsOfItsScope(t *testing.T) {
	f := newFixture(t)
	crewOnlyWorking(t, f)
	ctx := context.Background()
	if err := f.ing.CheckNoTurnScope(ctx, timeline.Scope{Projects: []string{"elsewhere"}}, 2*time.Minute); err != nil {
		t.Fatal(err)
	}
	if open, _ := noTurnIncidents(t, f); open != 0 {
		t.Fatalf("a check whose scope leaves the project out filed %d incident(s)", open)
	}
	if err := f.ing.CheckNoTurnScope(ctx, shopOnly, 2*time.Minute); err != nil {
		t.Fatal(err)
	}
	if open, _ := noTurnIncidents(t, f); open != 1 {
		t.Fatalf("no_turn open = %d over the project with the crew, want 1", open)
	}
}
