package spawn_test

import (
	"context"
	"os"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/brief/brieftest"
	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// HaltCrew is the hard budget ceiling's stop for unlanded work (CON-94):
// the agent and its tab go, the worktree, the branch and the task stay
// open, and the meta names no agent - a crew `crew relaunch` can start
// again and `crew stop` can still close.
func TestHaltCrewStopsTheAgentAndKeepsTheWork(t *testing.T) {
	w := crewWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	res, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project: "shop", Crew: "k3", BriefText: brieftest.Ship("work"), Budget: 2_000_000,
	})
	if err != nil {
		t.Fatalf("SpawnCrew: %v", err)
	}
	commitInWorktree(t, res.Worktree, "new.txt", "unlanded\n")

	stopped, err := spawn.HaltCrew(context.Background(), w, deps, "shop", "k3")
	if err != nil {
		t.Fatalf("HaltCrew: %v", err)
	}
	if !stopped {
		t.Fatal("HaltCrew reported no live agent stopped")
	}
	if _, ok := rt.Tabs[res.Pane]; ok {
		t.Fatal("the crew's tab is still open")
	}
	if _, statErr := os.Stat(res.Worktree); statErr != nil {
		t.Fatalf("the worktree is gone: %v", statErr)
	}
	exists, err := gitx.New().BranchExists(context.Background(), w.RepoDir("shop/shop"), "mate/k3")
	if err != nil || !exists {
		t.Fatalf("the branch is gone: %v, %v", exists, err)
	}
	meta, err := w.ReadCrewMeta("shop", "k3")
	if err != nil {
		t.Fatal(err)
	}
	if meta[spawn.MetaAgent] != "" || meta[spawn.MetaPane] != "" {
		t.Fatalf("meta still names an agent: %v", meta)
	}
	if meta[store.MetaBudget] != "2000000" || meta[spawn.MetaBranch] == "" {
		t.Fatalf("meta lost what the task needs: %v", meta)
	}
	if crewstate.State(meta[spawn.MetaState]).Closed() {
		t.Fatalf("state = %q, want the task still open", meta[spawn.MetaState])
	}
}
