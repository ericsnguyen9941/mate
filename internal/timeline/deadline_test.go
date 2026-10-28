package timeline_test

import (
	"context"
	"errors"
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

// TestIngestBoundsEachProjectsPass: the ingest runs at the end of every
// observer round, so a project whose git never answers must cost the round
// one bounded wait, not the round itself. The bounded pass is rolled back
// whole, and the next pass whose git answers records the project as if the
// hung one had never run.
func TestIngestBoundsEachProjectsPass(t *testing.T) {
	f := newFixture(t)
	git := &hangingGit{}
	git.hang.Store(true)
	const bound = 100 * time.Millisecond
	ing := timeline.New(f.ws, f.db, timeline.Deps{
		Now:             func() time.Time { return fixtureNow },
		Harnesses:       catalog.Default(),
		TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"), filepath.Join(f.root, "no-codex-sessions")),
		Git:             gitx.Git{Runner: git},
		ProjectTimeout:  bound,
	})

	done := make(chan error, 1)
	start := time.Now()
	go func() { done <- ing.Ingest(context.Background()) }()
	var err error
	select {
	case err = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Ingest did not return: a hung git held the pass past its bound")
	}
	if took := time.Since(start); took > 2*time.Second {
		t.Fatalf("Ingest took %s with a %s bound per project", took, bound)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Ingest = %v, want the project's deadline reported", err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM event`); n != 0 {
		t.Fatalf("a pass that ran out its bound left %d event(s) behind", n)
	}

	git.hang.Store(false)
	// The next poll gets the production bound: the 100ms one was for the
	// hung git, and a real git on a loaded machine takes longer than that
	// (CON-132). The cursors are the database's, so a new Ingester is the
	// same next pass.
	next := timeline.New(f.ws, f.db, timeline.Deps{
		Now:             func() time.Time { return fixtureNow },
		Harnesses:       catalog.Default(),
		TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"), filepath.Join(f.root, "no-codex-sessions")),
		Git:             gitx.Git{Runner: git},
	})
	if err := next.Ingest(context.Background()); err != nil {
		t.Fatalf("Ingest after git answered again: %v", err)
	}
	if len(f.story(t)) == 0 {
		t.Fatal("the pass after the bounded one recorded no story")
	}
}
