package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/watch"
)

// TestWatchBeginsARoundOnARoundRuntime: Poll begins a round on a Runtime
// that keeps answers per round (CON-264) before it asks anything, so a
// round sees the agent Herdr has now and not the one the last round listed,
// and every round still reads the pane.
func TestWatchBeginsARoundOnARoundRuntime(t *testing.T) {
	f := newFixture(t)
	cache := runtime.NewRoundCache(f.rt)
	deps := f.deps()
	deps.Runtime = cache
	f.w = watch.New(f.ws, deps)
	for i := 0; i < 3; i++ {
		f.poll()
		f.clock.advance(5 * time.Second)
	}
	if got := len(f.rt.ReadCalls); got != 3 {
		t.Fatalf("pane reads in three rounds = %d, want 3: every round reads", got)
	}
	f.rt.DropAgent(f.handle)
	f.poll()
	f.assertIncidents("k3 runtime_lost open")
	if _, err := cache.InspectAgent(context.Background(), f.handle); !runtime.IsAgentNotFound(err) {
		t.Fatalf("InspectAgent after the round = %v, want the round's agent_not_found", err)
	}
}
