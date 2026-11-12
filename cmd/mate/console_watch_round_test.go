package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/watch"
)

// The Herdr cost of one watcher round (CON-264, epic CON-253: the core
// should cost nothing when idle). Each call the fake Herdr records is one
// `herdr` process the live adapter would launch.

// idleCrews spawns n crews on the fake Herdr, each with an idle Codex
// composer that does not change.
func idleCrews(t *testing.T, n int) (*store.Workspace, *runtime.Fake, spawn.Deps, []string, []runtime.AgentHandle) {
	t.Helper()
	w, deps := consoleFixture(t, "shop")
	rt := deps.Runtime.(*runtime.Fake)
	crews := make([]string, 0, n)
	handles := make([]runtime.AgentHandle, 0, n)
	for i := 0; i < n; i++ {
		crew := fmt.Sprintf("k%d", i)
		res := spawnFakeCrew(t, w, deps, "shop", crew)
		rt.SetReadOutput(spawnedHandle(res), codexEmptyScreen)
		crews = append(crews, crew)
		handles = append(handles, spawnedHandle(res))
	}
	return w, rt, deps, crews, handles
}

// roundCalls builds the watcher and the ingest's session ref over d, as
// consoleWatcherWithTimeline does, runs a first round (which reads every
// pane: there is no snapshot yet), calls between, and returns the Herdr
// calls of the second: the watcher's Poll, then the ingest asking each
// crew's session ref, as its locator does for a transcript it has not
// recorded.
func roundCalls(t *testing.T, w *store.Workspace, rt *runtime.Fake, d spawn.Deps, crews []string, between func()) []string {
	t.Helper()
	ctx := context.Background()
	watcher := watch.New(w, watch.Deps{
		Runtime: d.Runtime,
		Handle:  consoleCrewHandle(w, d),
		Session: consoleSession(w, d),
	})
	ref := consoleSessionRef(w, d)
	round := func() {
		if err := watcher.Poll(ctx); err != nil {
			t.Fatalf("Poll: %v", err)
		}
		for _, crew := range crews {
			if _, err := ref(ctx, "shop", crew); err != nil {
				t.Fatalf("session ref of %s: %v", crew, err)
			}
		}
		for _, crew := range crews {
			if h, ok := watcher.Health("shop", crew); !ok || !h.AgentPresent {
				t.Fatalf("health of %s = %+v (ok=%v), want a present agent", crew, h, ok)
			}
		}
	}
	round()
	if between != nil {
		between()
	}
	rt.Calls = nil
	round()
	return append([]string(nil), rt.Calls...)
}

func tally(calls []string) string {
	n := map[string]int{}
	for _, c := range calls {
		n[c]++
	}
	keys := make([]string, 0, len(n))
	for k := range n {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", k, n[k]))
	}
	return strings.Join(parts, " ")
}

// TestWatcherRoundAsksHerdrOnce: ten open crews. Before CON-264 a round
// was, per crew, a session lookup, an `agent get` and an `agent read` for
// the watcher and a session lookup and an `agent get` for the ingest, plus
// the session check: 51 `herdr` processes. Through roundDeps it is one
// `session list` and one `agent list`, shared by the watcher and the
// ingest, and still one `agent read` per crew (watch.go says why every
// pane is read every round): 12.
func TestWatcherRoundAsksHerdrOnce(t *testing.T) {
	t.Parallel()
	w, rt, deps, crews, _ := idleCrews(t, 10)

	before := roundCalls(t, w, rt, deps, crews, nil)
	if got, want := tally(before), "InspectAgent=20 LookupSession=21 ReadAgent=10"; got != want {
		t.Fatalf("calls of a round asked directly = %s, want %s", got, want)
	}

	after := roundCalls(t, w, rt, roundDeps(deps), crews, nil)
	if got, want := tally(after), "ListAgents=1 LookupSession=1 ReadAgent=10"; got != want {
		t.Fatalf("calls of a round through roundDeps = %s, want %s", got, want)
	}
	t.Logf("herdr calls per round, 10 crews: %d asked directly, %d through roundDeps", len(before), len(after))
}

// TestRoundDepsTreatsAnAgentGoneFromTheListAsGone: a crew whose agent
// Herdr no longer has opens runtime_lost through the round's list exactly
// as through `agent get`.
func TestRoundDepsTreatsAnAgentGoneFromTheListAsGone(t *testing.T) {
	t.Parallel()
	w, deps := consoleFixture(t, "shop")
	res := spawnFakeCrew(t, w, deps, "shop", "k3")
	rt := deps.Runtime.(*runtime.Fake)
	rt.SetReadOutput(spawnedHandle(res), codexEmptyScreen)

	d := roundDeps(deps)
	watcher := watch.New(w, watch.Deps{Runtime: d.Runtime, Handle: consoleCrewHandle(w, d)})
	if err := watcher.Poll(context.Background()); err != nil {
		t.Fatalf("Poll: %v", err)
	}
	rt.DropAgent(spawnedHandle(res))
	if err := watcher.Poll(context.Background()); err != nil {
		t.Fatalf("Poll: %v", err)
	}
	entries, _, err := w.ReadIncidents("shop", 0)
	if err != nil {
		t.Fatalf("ReadIncidents: %v", err)
	}
	if len(entries) != 1 || entries[0].Kind != "runtime_lost" || !entries[0].Open() {
		t.Fatalf("incidents = %+v, want one open runtime_lost", entries)
	}
	if h, ok := watcher.Health("shop", "k3"); !ok || h.AgentPresent {
		t.Fatalf("health = %+v (ok=%v), want the agent absent", h, ok)
	}
}
