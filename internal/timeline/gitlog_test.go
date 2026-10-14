package timeline_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// countingGit runs the real git and remembers every command it was asked to
// run, so a test can say how much git one ingest pass cost.
type countingGit struct {
	mu   sync.Mutex
	runs []string
}

func (c *countingGit) Run(ctx context.Context, cmd gitx.Command) (gitx.Result, error) {
	c.mu.Lock()
	c.runs = append(c.runs, strings.Join(cmd.Args, " "))
	c.mu.Unlock()
	return gitx.ExecRunner{}.Run(ctx, cmd)
}

func (c *countingGit) take() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	runs := c.runs
	c.runs = nil
	return runs
}

// countedIngester is a second ingest over the fixture's workspace and
// database whose git is counted.
func (f *fixture) countedIngester(t *testing.T) (*timeline.Ingester, *countingGit) {
	t.Helper()
	counter := &countingGit{}
	ing := timeline.New(f.ws, f.db, timeline.Deps{
		Now:       func() time.Time { return fixtureNow },
		Harnesses: catalog.Default(),
		TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"),
			filepath.Join(f.root, "no-codex-sessions")),
		Git: gitx.Git{Runner: counter},
	})
	return ing, counter
}

func ingestWith(t *testing.T, ing *timeline.Ingester) {
	t.Helper()
	if err := ing.Ingest(context.Background()); err != nil {
		t.Fatalf("Ingest: %v", err)
	}
}

// closeCrew marks the fixture's crew closed with the given state, pointed at
// a transcript.
func (f *fixture) closeCrew(t *testing.T, state, transcript string) {
	t.Helper()
	meta, err := f.ws.ReadCrewMeta(fixtureProject, fixtureCrew)
	if err != nil {
		t.Fatalf("ReadCrewMeta: %v", err)
	}
	if transcript != "" {
		meta[timeline.MetaTranscript] = transcript
	}
	meta[timeline.MetaState] = state
	meta[timeline.MetaStoppedAt] = mustTime("2026-09-19T10:47:00Z").Format(time.RFC3339)
	if err := f.ws.WriteCrewMeta(fixtureProject, fixtureCrew, meta); err != nil {
		t.Fatalf("WriteCrewMeta: %v", err)
	}
}

// commitOnBranch makes one commit on branch (creating it from main when it
// does not exist) and returns to main. It returns the commit's short and full
// sha.
func commitOnBranch(t *testing.T, repo, branch, line string) (string, string) {
	t.Helper()
	if strings.TrimSpace(gitOutput(t, repo, "branch", "--list", branch)) == "" {
		gitRun(t, repo, "checkout", "-b", branch)
	} else {
		gitRun(t, repo, "checkout", branch)
	}
	f, err := os.OpenFile(filepath.Join(repo, "README.md"), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open README: %v", err)
	}
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatalf("append README: %v", err)
	}
	_ = f.Close()
	gitRun(t, repo, "commit", "-am", "docs: add Buy link")
	short := strings.TrimSpace(gitOutput(t, repo, "rev-parse", "--short", "HEAD"))
	full := strings.TrimSpace(gitOutput(t, repo, "rev-parse", "HEAD"))
	gitRun(t, repo, "checkout", "main")
	return short, full
}

// A crew that is closed, whose branch is gone and whose commit and merge are
// already on the timeline has nothing left for git to say. Measured
// 2026-10-14 on a live workspace of 105 such crews, asking anyway cost 2.2 s
// of git on every five-second poll, inside the ingest's write transaction,
// and recorded nothing.
func TestASettledCrewCostsNoGitOnTheNextPass(t *testing.T) {
	f := newFixture(t)
	repo := filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject)
	short, _ := commitOnBranch(t, repo, fixtureBranch, "[Buy](pages/checkout-express.html)")
	gitRun(t, repo, "merge", "--ff-only", fixtureBranch)
	gitRun(t, repo, "branch", "-D", fixtureBranch)
	rollout := filepath.Join(f.root, "crew-rollout.jsonl")
	writeCommitRollout(t, rollout, repo, fixtureBranch, short)
	f.closeCrew(t, "finished", rollout)

	ing, counter := f.countedIngester(t)
	ingestWith(t, ing)
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ?`, timeline.KindMergeDone); n != 1 {
		t.Fatalf("%d merge.done event(s) after the first pass, want 1", n)
	}
	counter.take()

	ingestWith(t, ing)
	if runs := counter.take(); len(runs) != 0 {
		t.Fatalf("the second pass ran %d git command(s) for a crew whose story is complete: %q", len(runs), runs)
	}

	// A rebuild from scratch on the same ingest still records all of it.
	if err := ing.Reindex(context.Background()); err != nil {
		t.Fatalf("Reindex: %v", err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ? AND actor_id = ?`,
		timeline.KindGitCommitted, f.crewActor()); n != 1 {
		t.Fatalf("%d git.committed event(s) after a reindex, want 1", n)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ?`, timeline.KindMergeDone); n != 1 {
		t.Fatalf("%d merge.done event(s) after a reindex, want 1", n)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM task WHERE crew_actor_id = ? AND merged_event_id IS NOT NULL`,
		f.crewActor()); n != 1 {
		t.Fatal("the task lost its merge after a reindex")
	}
}

// A sha in a transcript that git cannot resolve - the crew committed in a
// clone this workspace does not have - is asked about once, not every pass.
func TestAnUnresolvableCommitOfAClosedCrewIsAskedAboutOnce(t *testing.T) {
	f := newFixture(t)
	repo := filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject)
	rollout := filepath.Join(f.root, "crew-rollout.jsonl")
	writeCommitRollout(t, rollout, repo, fixtureBranch, "deadbee")
	f.closeCrew(t, "failed", rollout)

	ing, counter := f.countedIngester(t)
	ingestWith(t, ing)
	if runs := counter.take(); len(runs) == 0 {
		t.Fatal("the first pass asked git nothing; the test is not exercising the sighting")
	}
	ingestWith(t, ing)
	if runs := counter.take(); len(runs) != 0 {
		t.Fatalf("the second pass ran %d git command(s) for a closed crew with nothing new: %q", len(runs), runs)
	}
}

// A crew still at work keeps being read: a commit made after a pass is
// recorded by the next one.
func TestAnOpenCrewStillGetsItsNewCommitsRecorded(t *testing.T) {
	f := newFixture(t)
	repo := filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject)
	commitOnBranch(t, repo, fixtureBranch, "first")

	ing, _ := f.countedIngester(t)
	ingestWith(t, ing)
	ingestWith(t, ing)
	_, second := commitOnBranch(t, repo, fixtureBranch, "second")
	ingestWith(t, ing)

	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ? AND actor_id = ?`,
		timeline.KindGitCommitted, f.crewActor()); n != 2 {
		t.Fatalf("%d git.committed event(s), want both commits of the open crew", n)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ? AND json_extract(payload, '$.sha') = ?`,
		timeline.KindGitCommitted, second); n != 1 {
		t.Fatal("the commit made after the first passes was not recorded")
	}
}

// A finished crew whose branch was deleted before anyone merged it can still
// land later - its commit merged by hand - so its merge keeps being asked
// about, and only that.
func TestAFinishedCrewWhoseBranchIsGoneCanStillLandLater(t *testing.T) {
	f := newFixture(t)
	repo := filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject)
	_, full := commitOnBranch(t, repo, fixtureBranch, "[Buy](pages/checkout-express.html)")
	f.closeCrew(t, "finished", "")

	ing, counter := f.countedIngester(t)
	ingestWith(t, ing)
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ? AND actor_id = ?`,
		timeline.KindGitCommitted, f.crewActor()); n != 1 {
		t.Fatalf("%d git.committed event(s) while the branch was there, want 1", n)
	}
	gitRun(t, repo, "branch", "-D", fixtureBranch)
	ingestWith(t, ing)
	counter.take()

	ingestWith(t, ing)
	runs := counter.take()
	if len(runs) != 1 || !strings.HasPrefix(runs[0], "merge-base --is-ancestor") {
		t.Fatalf("a finished, unmerged crew whose branch is gone ran %q; want only the merge question", runs)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ?`, timeline.KindMergeDone); n != 0 {
		t.Fatal("a merge was recorded before the commit landed")
	}

	gitRun(t, repo, "merge", "--ff-only", full)
	ingestWith(t, ing)
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ?`, timeline.KindMergeDone); n != 1 {
		t.Fatalf("%d merge.done event(s) after the commit landed, want 1", n)
	}
}
