package spawn_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/brief/brieftest"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/teardown"
)

// setTeardown gives the project's only repo a teardown command.
func setTeardown(t *testing.T, w *store.Workspace, project, cmd string) {
	t.Helper()
	cfg, err := w.LoadProject(project)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Repos[0].Teardown = cmd
	if err := w.SaveProject(project, cfg); err != nil {
		t.Fatal(err)
	}
}

// recordingCleaner remembers what it was asked to tear down and whether the
// worktree was still there to tear down from.
type recordingCleaner struct {
	targets     []teardown.Target
	sawWorktree bool
	report      teardown.Report
}

func (c *recordingCleaner) Clean(_ context.Context, tg teardown.Target) teardown.Report {
	c.targets = append(c.targets, tg)
	if _, err := os.Stat(tg.Worktree); err == nil {
		c.sawWorktree = true
	}
	return c.report
}

// The repo's teardown command runs in the crew's worktree, with the crew
// named in its environment, before the worktree is removed: `make down` and
// plantrip's `make worktree-destroy` need the checkout's files (CON-111).
func TestStopCrewRunsTheRepoTeardownInTheWorktreeBeforeRemovingIt(t *testing.T) {
	w := crewWorkspace(t, "shop")
	marker := filepath.Join(t.TempDir(), "torn-down")
	setTeardown(t, w, "shop", `test -f README.md && echo "$MATE_PROJECT/$MATE_CREW $MATE_WORKTREE_NAME" > `+marker)
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	res, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project: "shop", Crew: "k3", BriefText: brieftest.Ship("work"),
	})
	if err != nil {
		t.Fatalf("SpawnCrew: %v", err)
	}

	stopped, err := spawn.StopCrew(context.Background(), w, deps, "shop", "k3", false)
	if err != nil {
		t.Fatalf("StopCrew: %v", err)
	}
	got, err := os.ReadFile(marker)
	if err != nil {
		t.Fatalf("the teardown command did not run in the worktree: %v", err)
	}
	if want := "shop/k3 shop-k3\n"; string(got) != want {
		t.Fatalf("teardown saw %q, want %q", got, want)
	}
	if stopped.Cleanup.Hook != "ok" || !stopped.WorktreeRemoved {
		t.Fatalf("stop = %+v, want the hook ok and the worktree removed", stopped)
	}
	if _, err := os.Stat(res.Worktree); !os.IsNotExist(err) {
		t.Fatalf("the worktree survived: %v", err)
	}
	meta, err := w.ReadCrewMeta("shop", "k3")
	if err != nil {
		t.Fatal(err)
	}
	if meta[spawn.MetaCleanup] != "done (hook ok)" {
		t.Fatalf("cleanup = %q, want done (hook ok)", meta[spawn.MetaCleanup])
	}
	if _, err := os.Stat(w.CrewTeardownLog("shop", "k3")); err != nil {
		t.Fatalf("no teardown log: %v", err)
	}
}

// Whatever a teardown could not remove is recorded and reported, and the
// task still closes: a stopped docker daemon must not keep a merged task
// open.
func TestStopCrewClosesTheTaskWhenTheCleanupLeavesSomethingBehind(t *testing.T) {
	w := crewWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	cleaner := &recordingCleaner{report: teardown.Report{
		Removed: []string{"simulator Shop-shop-k3"},
		Pending: []string{"docker: docker ps exited 1: Cannot connect to the Docker daemon"},
	}}
	deps.Cleaner = cleaner
	res, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project: "shop", Crew: "k3", BriefText: brieftest.Ship("work"),
	})
	if err != nil {
		t.Fatalf("SpawnCrew: %v", err)
	}

	stopped, err := spawn.StopCrew(context.Background(), w, deps, "shop", "k3", false)
	if err != nil {
		t.Fatalf("StopCrew: %v", err)
	}
	if len(cleaner.targets) != 1 || cleaner.targets[0].Worktree != res.Worktree || !cleaner.sawWorktree {
		t.Fatalf("cleaner saw %+v (worktree there: %v), want one teardown of %s before it went", cleaner.targets, cleaner.sawWorktree, res.Worktree)
	}
	if stopped.State != spawn.CrewStateFinished || !stopped.WorktreeRemoved {
		t.Fatalf("stop = %+v, want the task finished and the worktree removed", stopped)
	}
	meta, err := w.ReadCrewMeta("shop", "k3")
	if err != nil {
		t.Fatal(err)
	}
	if got := meta[spawn.MetaCleanup]; !strings.HasPrefix(got, "pending: docker:") || !strings.Contains(got, "simulator Shop-shop-k3") {
		t.Fatalf("cleanup = %q, want the pending docker line and what was removed", got)
	}
}

// A refused stop changes nothing, so nothing is torn down: the crew is
// still working in that stack.
func TestStopCrewRefusalTearsNothingDown(t *testing.T) {
	w := crewWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	cleaner := &recordingCleaner{}
	deps.Cleaner = cleaner
	res, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project: "shop", Crew: "k3", BriefText: brieftest.Ship("work"),
	})
	if err != nil {
		t.Fatalf("SpawnCrew: %v", err)
	}
	commitInWorktree(t, res.Worktree, "new.txt", "unlanded\n")

	if _, err := spawn.StopCrew(context.Background(), w, deps, "shop", "k3", false); !errors.Is(err, spawn.ErrUnlandedWork) {
		t.Fatalf("err = %v, want ErrUnlandedWork", err)
	}
	if len(cleaner.targets) != 0 {
		t.Fatalf("a refused stop tore down %+v", cleaner.targets)
	}
	meta, err := w.ReadCrewMeta("shop", "k3")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := meta[spawn.MetaCleanup]; ok {
		t.Fatalf("a refused stop wrote cleanup=: %v", meta)
	}
}
