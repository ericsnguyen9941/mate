package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/send"
	"github.com/nguyenngocanh94/mate/internal/watch"
)

// The pane read is skipped while Herdr's output revision of the pane has
// not moved (CON-264): the fixtures' fake Herdr reports a revision that
// moves exactly when the scripted screen changes, and Fake.ReadCalls counts
// the `agent read`s the real one would have run.

func (f *fixture) reads() int {
	f.t.Helper()
	return len(f.rt.ReadCalls)
}

func TestWatchReadsThePaneOnlyWhenItsRevisionMoves(t *testing.T) {
	f := newFixture(t)
	f.poll()
	if got := f.reads(); got != 1 {
		t.Fatalf("reads after the first round = %d, want 1: there is no snapshot yet", got)
	}

	for i := 0; i < 3; i++ {
		f.clock.advance(5 * time.Second)
		f.poll()
	}
	if got := f.reads(); got != 1 {
		t.Fatalf("reads after three still rounds = %d, want still 1", got)
	}
	h, ok := f.health("k3")
	if !ok || !h.AgentPresent || h.Composer != send.StateEmpty || !h.ObservedAt.Equal(f.clock.Now()) {
		t.Fatalf("health on a still pane = %+v (ok=%v), want this round's reading of the kept snapshot", h, ok)
	}
	if h.QuietFor != 15*time.Second {
		t.Fatalf("QuietFor = %s, want 15s: a skipped read is a pane that did not move", h.QuietFor)
	}

	f.setScreen(codexBusyScreen)
	f.clock.advance(5 * time.Second)
	f.poll()
	if got := f.reads(); got != 2 {
		t.Fatalf("reads after the pane moved = %d, want 2", got)
	}
	h, _ = f.health("k3")
	if h.Composer != send.StateBusy || h.QuietFor != 0 {
		t.Fatalf("health after the pane moved = %+v, want busy and the quiet clock reset", h)
	}
}

// TestWatchReadsEveryRoundWithoutARevision: a Herdr that reports no
// revision (0.8.x) cannot say a pane is unchanged, so every round reads it,
// as every round did before CON-264.
func TestWatchReadsEveryRoundWithoutARevision(t *testing.T) {
	f := newFixture(t)
	f.rt.NoRevision = true
	for i := 0; i < 3; i++ {
		f.poll()
		f.clock.advance(5 * time.Second)
	}
	if got := f.reads(); got != 3 {
		t.Fatalf("reads in three rounds without a revision = %d, want 3", got)
	}
}

// TestWatchReadsAPaneTheAgentMovedTo: a revision counts one pane's output,
// so the same number in another pane says nothing about the snapshot held.
func TestWatchReadsAPaneTheAgentMovedTo(t *testing.T) {
	f := newFixture(t)
	f.poll()
	moved := f.handle
	moved.Tab.PaneID = "pane-k3-again"
	f.rt.PutAgent(moved, runtime.AgentIdle)
	f.clock.advance(5 * time.Second)
	f.poll()
	if got := f.reads(); got != 2 {
		t.Fatalf("reads after the agent moved pane = %d, want 2", got)
	}
}

// TestWatchReadsAgainWhenTheStatusChanges: a working Codex is read through
// the visible screen rather than the scrollback (internal/runtime/herdr
// readAgent), which is other text from the same output, so a status change
// re-reads a pane whose revision did not move.
func TestWatchReadsAgainWhenTheStatusChanges(t *testing.T) {
	f := newFixture(t)
	f.poll()
	f.rt.PutAgent(f.handle, runtime.AgentWorking)
	f.clock.advance(5 * time.Second)
	f.poll()
	if got := f.reads(); got != 2 {
		t.Fatalf("reads after the status changed = %d, want 2", got)
	}
}

// TestWatchOpensStaleOnAStillPaneWithoutReadingIt: the stale rule counts
// quiet from the last time the pane moved, and a pane whose revision holds
// has not moved - so three minutes of rounds that read nothing still open
// stale on time, and not before.
func TestWatchOpensStaleOnAStillPaneWithoutReadingIt(t *testing.T) {
	f := newFixture(t)
	f.appendStatus("k3", "working: reading the tests")
	f.poll()
	for elapsed := 5 * time.Second; elapsed < 3*time.Minute; elapsed += 5 * time.Second {
		f.clock.advance(5 * time.Second)
		f.poll()
	}
	f.assertIncidents()

	f.clock.advance(5 * time.Second)
	f.poll()
	f.assertIncidents("k3 stale open")
	if got := f.reads(); got != 1 {
		t.Fatalf("reads over three still minutes = %d, want 1", got)
	}

	// And the pane moving is still seen, through the revision, and
	// resolves it.
	f.setScreen(codexBusyScreen)
	f.clock.advance(5 * time.Second)
	f.poll()
	f.assertIncidents("k3 stale open", "k3 stale resolved")
}

// TestWatchBeginsARoundOnARoundRuntime: Poll begins a round on a Runtime
// that keeps answers per round before it asks anything, so a round sees the
// agent Herdr has now and not the one the last round listed.
func TestWatchBeginsARoundOnARoundRuntime(t *testing.T) {
	f := newFixture(t)
	cache := runtime.NewRoundCache(f.rt)
	deps := f.deps()
	deps.Runtime = cache
	f.w = watch.New(f.ws, deps)
	f.poll()
	f.rt.DropAgent(f.handle)
	f.clock.advance(5 * time.Second)
	f.poll()
	f.assertIncidents("k3 runtime_lost open")
	if _, err := cache.InspectAgent(context.Background(), f.handle); !runtime.IsAgentNotFound(err) {
		t.Fatalf("InspectAgent after the round = %v, want the round's agent_not_found", err)
	}
}
