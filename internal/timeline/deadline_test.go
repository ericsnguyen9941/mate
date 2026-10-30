package timeline_test

import (
	"context"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// hangingGit is the real git runner with a switch that makes every call hang
// until its context ends, as a git wedged on a lock or a network remote does.
type hangingGit struct {
	hang atomic.Bool
}

func (g *hangingGit) Run(ctx context.Context, cmd gitx.Command) (gitx.Result, error) {
	if g.hang.Load() {
		<-ctx.Done()
		return gitx.Result{}, ctx.Err()
	}
	return gitx.ExecRunner{}.Run(ctx, cmd)
}

// boundedIngester is an ingest over the fixture whose git calls are held to
// bound per pass.
func (f *fixture) boundedIngester(git gitx.Runner, bound time.Duration) *timeline.Ingester {
	return timeline.New(f.ws, f.db, timeline.Deps{
		Now:             func() time.Time { return fixtureNow },
		Harnesses:       catalog.Default(),
		TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"), filepath.Join(f.root, "no-codex-sessions")),
		Git:             gitx.Git{Runner: git},
		ProjectTimeout:  bound,
	})
}

// TestAHungGitCostsOneBoundAndThePassStillCommits: the ingest runs at the end
// of every observer round, so a project whose git never answers must cost the
// round one bounded wait, not the round itself. What the files say is still
// recorded: git having said nothing is asked again by the next pass, the one
// whose git answers.
func TestAHungGitCostsOneBoundAndThePassStillCommits(t *testing.T) {
	f := newFixture(t)
	repo := filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject)
	commitOnBranch(t, repo, fixtureBranch, "[Buy](pages/checkout-express.html)")
	git := &hangingGit{}
	git.hang.Store(true)
	const bound = 100 * time.Millisecond

	done := make(chan error, 1)
	start := time.Now()
	go func() { done <- f.boundedIngester(git, bound).Ingest(context.Background()) }()
	var err error
	select {
	case err = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Ingest did not return: a hung git held the pass past its bound")
	}
	if took := time.Since(start); took > 2*time.Second {
		t.Fatalf("Ingest took %s with a %s bound on git", took, bound)
	}
	if err != nil {
		t.Fatalf("Ingest = %v, want the pass recorded without git", err)
	}
	if len(f.story(t)) == 0 {
		t.Fatal("a pass whose git hung recorded nothing from the files")
	}
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ?`, timeline.KindGitCommitted); n != 0 {
		t.Fatalf("%d git.committed event(s) from a git that never answered", n)
	}

	git.hang.Store(false)
	// The next poll gets the production bound: the 100ms one was for the
	// hung git, and a real git on a loaded machine takes longer than that
	// (CON-132).
	if err := f.boundedIngester(git, 0).Ingest(context.Background()); err != nil {
		t.Fatalf("Ingest after git answered again: %v", err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ?`, timeline.KindGitCommitted); n != 1 {
		t.Fatalf("%d git.committed event(s) once git answered, want 1", n)
	}
}

// TestAPassLongerThanItsBoundStillCommits is CON-139: the bound used to cover
// the whole pass and roll it back, so once a project's backlog took longer
// than the bound to read, every poll redid it and threw it away, and the
// project's timeline froze. A bound already spent before the pass reads a
// file must still leave the files recorded.
func TestAPassLongerThanItsBoundStillCommits(t *testing.T) {
	f := newFixture(t)
	if err := f.boundedIngester(gitx.ExecRunner{}, time.Nanosecond).Ingest(context.Background()); err != nil {
		t.Fatalf("Ingest = %v, want the pass committed", err)
	}
	if len(f.story(t)) == 0 {
		t.Fatal("a pass that outlived its bound recorded nothing")
	}
}
