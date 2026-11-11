package timeline_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// The seal (CON-259): a crew closed over an hour ago whose every conclusion
// is on the table is skipped whole by later polls of the same Ingester.

// sealLater is a clock two hours after the fixture crew closed: past
// settledFor, so the crew may be sealed.
var sealLater = fixtureNow.Add(2 * time.Hour)

// sealIngester is one long-lived ingest - the observer's - over the
// fixture, two hours after the crew closed, its git and its runtime asks
// counted.
func (f *fixture) sealIngester(now time.Time) (*timeline.Ingester, *countingGit, *refCounter) {
	counter, refs := &countingGit{}, &refCounter{}
	return timeline.New(f.ws, f.db, timeline.Deps{
		Now:        func() time.Time { return now },
		SessionRef: refs.ref,
		Harnesses:  catalog.Default(),
		TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"),
			filepath.Join(f.root, "no-codex-sessions")),
		Git: gitx.Git{Runner: counter},
	}), counter, refs
}

// refCounter counts the runtime's session-ref asks, which a Codex crew's
// locate makes before it adopts by cwd.
type refCounter struct {
	mu sync.Mutex
	n  int
}

func (r *refCounter) ref(context.Context, string, string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.n++
	return "", nil
}

func (r *refCounter) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.n
}

// tamperCrewActor overwrites a column the meta pass upserts, so a later
// pass that read the crew's `.meta` again puts it back.
func (f *fixture) tamperCrewActor(t *testing.T) {
	t.Helper()
	if _, err := f.db.SQL().Exec(`UPDATE actor SET harness = 'tampered' WHERE id = ?`, f.crewActor()); err != nil {
		t.Fatal(err)
	}
}

// crewWasRead reports whether a pass since tamperCrewActor read the crew.
func (f *fixture) crewWasRead(t *testing.T) bool {
	t.Helper()
	var harness string
	if err := f.db.SQL().QueryRow(`SELECT harness FROM actor WHERE id = ?`, f.crewActor()).Scan(&harness); err != nil {
		t.Fatal(err)
	}
	return harness != "tampered"
}

// unreadableProfile leaves the crew a harness profile no pass can read, so
// a pass that reads the crew's profiles fails.
func (f *fixture) unreadableProfile(t *testing.T) {
	t.Helper()
	dir := filepath.Join(f.ws.CrewsDir(fixtureProject), fixtureCrew)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "harness-profile.json")
	if err := os.WriteFile(path, []byte(`{}`), 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })
}

func (f *fixture) readableProfile(t *testing.T) {
	t.Helper()
	path := filepath.Join(f.ws.CrewsDir(fixtureProject), fixtureCrew, "harness-profile.json")
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAClosedCrewTheIngestHasFinishedWithIsSkippedOnLaterPolls(t *testing.T) {
	f := newFixture(t)
	f.closeCrew(t, "failed", "")
	ing, counter, _ := f.sealIngester(sealLater)
	ingestWith(t, ing)
	if n := f.count(t, `SELECT COUNT(*) FROM turn WHERE actor_id = ?`, f.crewActor()); n != codexTurns {
		t.Fatalf("%d turn(s) recorded for the crew before it was sealed, want %d", n, codexTurns)
	}

	f.unreadableProfile(t)
	f.tamperCrewActor(t)
	counter.take()
	if err := ing.Ingest(context.Background()); err != nil {
		t.Fatalf("a poll read a sealed crew's harness profile: %v", err)
	}
	if f.crewWasRead(t) {
		t.Fatal("a poll upserted a sealed crew's actor: its `.meta` was read again")
	}
	if runs := counter.take(); len(runs) != 0 {
		t.Fatalf("a poll asked git %q about a sealed crew", runs)
	}

	// The proof is the seal and not the files: an Ingester that has sealed
	// nothing reads the crew, and its profile with it.
	cold, _, _ := f.sealIngester(sealLater)
	if err := cold.Ingest(context.Background()); err == nil {
		t.Fatal("a fresh Ingester left the crew unread: the test proves nothing")
	}
}

func TestACrewClosedMomentsAgoIsNotSealed(t *testing.T) {
	f := newFixture(t)
	f.closeCrew(t, "failed", "")
	ing, _, _ := f.sealIngester(fixtureNow)
	ingestWith(t, ing)
	f.tamperCrewActor(t)
	ingestWith(t, ing)
	if !f.crewWasRead(t) {
		t.Fatal("a crew closed moments ago was sealed: its commit sightings may still need git")
	}
}

func TestAReopenedCrewIsUnsealed(t *testing.T) {
	for _, tc := range []struct {
		name   string
		reopen func(t *testing.T, f *fixture)
	}{
		{"relaunched", func(t *testing.T, f *fixture) {
			// What `mate crew relaunch` writes: a new launch, state=spawned,
			// the close keys gone.
			meta, err := f.ws.ReadCrewMeta(fixtureProject, fixtureCrew)
			if err != nil {
				t.Fatal(err)
			}
			meta[timeline.MetaState] = "spawned"
			meta[timeline.MetaLaunchedAt] = sealLater.Format(time.RFC3339)
			delete(meta, timeline.MetaStoppedAt)
			if err := f.ws.WriteCrewMeta(fixtureProject, fixtureCrew, meta); err != nil {
				t.Fatal(err)
			}
		}},
		{"a status line after the seal", func(t *testing.T, f *fixture) {
			if err := f.ws.AppendStatus(fixtureProject, fixtureCrew, "working: a line after the seal"); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			f.closeCrew(t, "failed", "")
			ing, _, _ := f.sealIngester(sealLater)
			ingestWith(t, ing)
			f.tamperCrewActor(t)
			ingestWith(t, ing)
			if f.crewWasRead(t) {
				t.Fatal("the crew was not sealed to begin with")
			}

			tc.reopen(t, f)
			ingestWith(t, ing)
			if !f.crewWasRead(t) {
				t.Fatal("a poll skipped a crew whose files changed after it was sealed")
			}
		})
	}
}

func TestARelaunchedCrewStaysUnsealedWhileOpen(t *testing.T) {
	f := newFixture(t)
	f.closeCrew(t, "failed", "")
	ing, _, _ := f.sealIngester(sealLater)
	ingestWith(t, ing)
	meta, err := f.ws.ReadCrewMeta(fixtureProject, fixtureCrew)
	if err != nil {
		t.Fatal(err)
	}
	meta[timeline.MetaState] = "spawned"
	delete(meta, timeline.MetaStoppedAt)
	if err := f.ws.WriteCrewMeta(fixtureProject, fixtureCrew, meta); err != nil {
		t.Fatal(err)
	}
	ingestWith(t, ing)
	f.tamperCrewActor(t)
	ingestWith(t, ing)
	if !f.crewWasRead(t) {
		t.Fatal("an open crew was sealed")
	}
}

func TestReindexReadsASealedCrew(t *testing.T) {
	f := newFixture(t)
	f.closeCrew(t, "failed", "")
	ing, _, _ := f.sealIngester(sealLater)
	ingestWith(t, ing)

	f.unreadableProfile(t)
	if err := ing.Reindex(context.Background()); err == nil {
		t.Fatal("reindex left a sealed crew's harness profile unread")
	}
	f.readableProfile(t)
	if err := ing.Reindex(context.Background()); err != nil {
		t.Fatalf("Reindex: %v", err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM turn WHERE actor_id = ?`, f.crewActor()); n != codexTurns {
		t.Fatalf("%d turn(s) for the sealed crew after a reindex, want %d", n, codexTurns)
	}
	// And the poll after a reindex reads it again, then seals it anew.
	f.tamperCrewActor(t)
	ingestWith(t, ing)
	if !f.crewWasRead(t) {
		t.Fatal("the first poll after a reindex trusted a seal made against the emptied tables")
	}
	f.tamperCrewActor(t)
	ingestWith(t, ing)
	if f.crewWasRead(t) {
		t.Fatal("the crew was not sealed again after the reindex")
	}
}

// A finished crew is sealed only once its merge is decided: recorded, or
// impossible to prove because its branch is gone and the timeline holds no
// commit of it. One whose branch is still there may yet be merged.
func TestAFinishedCrewIsSealedOnlyOnceItsMergeIsDecided(t *testing.T) {
	for _, tc := range []struct {
		name   string
		setup  func(t *testing.T, repo string)
		sealed bool
	}{
		{"merged", func(t *testing.T, repo string) {
			commitOnBranch(t, repo, fixtureBranch, "[Buy](pages/checkout-express.html)")
			gitRun(t, repo, "merge", "--ff-only", fixtureBranch)
			gitRun(t, repo, "branch", "-D", fixtureBranch)
		}, true},
		{"branch gone, no commit recorded", func(*testing.T, string) {}, true},
		{"branch still there", func(t *testing.T, repo string) {
			commitOnBranch(t, repo, fixtureBranch, "[Buy](pages/checkout-express.html)")
		}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			repo := filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject)
			tc.setup(t, repo)
			f.closeCrew(t, "finished", "")
			ing, _, _ := f.sealIngester(sealLater)
			ingestWith(t, ing)
			f.tamperCrewActor(t)
			ingestWith(t, ing)
			if sealed := !f.crewWasRead(t); sealed != tc.sealed {
				t.Fatalf("sealed=%v, want %v", sealed, tc.sealed)
			}
		})
	}
}

// A Codex crew that never adopted a rollout was adopted again on every poll
// (CON-255 deferred this here). Closed long ago, it is located once.
func TestAClosedCodexCrewThatNeverAdoptedIsNotAdoptedAgain(t *testing.T) {
	for _, tc := range []struct {
		name  string
		state string
		setup func(t *testing.T, repo string)
	}{
		// Sealed: nothing about the crew is read again.
		{"sealed", "failed", func(*testing.T, string) {}},
		// Not sealable, its branch may still be merged: the locate's answer
		// is remembered all the same.
		{"waiting for its merge", "finished", func(t *testing.T, repo string) {
			commitOnBranch(t, repo, fixtureBranch, "[Buy](pages/checkout-express.html)")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			tc.setup(t, filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject))
			f.closeCrew(t, tc.state, "")
			clearCrewTranscript(t, f)
			ing, _, refs := f.sealIngester(sealLater)
			ingestWith(t, ing)
			asked := refs.count()
			if asked == 0 {
				t.Fatal("the first poll never tried to locate the crew's rollout")
			}
			if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ? AND actor_id = ?`,
				timeline.KindIngestUnresolved, f.crewActor()); n == 0 {
				t.Fatal("the crew's unlocated rollout was not recorded")
			}
			ingestWith(t, ing)
			ingestWith(t, ing)
			if n := refs.count(); n != asked {
				t.Fatalf("the crew's rollout was located %d more time(s) after the first poll", n-asked)
			}
		})
	}

	// Closed moments ago, it is still looked for: Codex may not have
	// written the rollout's head yet.
	f := newFixture(t)
	f.closeCrew(t, "failed", "")
	clearCrewTranscript(t, f)
	ing, _, refs := f.sealIngester(fixtureNow)
	ingestWith(t, ing)
	asked := refs.count()
	ingestWith(t, ing)
	if refs.count() == asked {
		t.Fatal("a crew closed moments ago stopped being located")
	}
}

// A pass past its deadline for git and Herdr seals nothing: a call that ran
// out of time said nothing.
func TestAPassPastItsDeadlineSealsNothing(t *testing.T) {
	f := newFixture(t)
	f.closeCrew(t, "failed", "")
	ing := timeline.New(f.ws, f.db, timeline.Deps{
		Now:       func() time.Time { return sealLater },
		Harnesses: catalog.Default(),
		TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"),
			filepath.Join(f.root, "no-codex-sessions")),
		ProjectTimeout: time.Nanosecond,
	})
	ingestWith(t, ing)
	f.tamperCrewActor(t)
	ingestWith(t, ing)
	if !f.crewWasRead(t) {
		t.Fatal("a pass past its deadline sealed the crew")
	}
}
