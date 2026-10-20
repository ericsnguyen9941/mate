package main

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/teardown"
)

// gc collects for closed crews whose worktree is gone and for what docker
// or Xcode still name under .worktrees/, never for a worktree that exists
// or a crew still open; without --apply it only lists (CON-114).
func TestGCTargetsOnlyWorktreesThatAreGone(t *testing.T) {
	w := wakeFixture(t)
	live := filepath.Join(w.WorktreesDir(), "shop-live")
	if err := os.MkdirAll(live, 0o755); err != nil {
		t.Fatal(err)
	}
	for crew, meta := range map[string]map[string]string{
		"gone-1": {spawn.MetaState: spawn.CrewStateFinished, spawn.MetaWorktree: ".worktrees/shop-gone-1"},
		"rev-r1": {spawn.MetaState: spawn.CrewStateReviewed, spawn.MetaWorktree: ".worktrees/shop-rev-r1"},
		"live":   {spawn.MetaState: spawn.CrewStateFinished, spawn.MetaWorktree: ".worktrees/shop-live"},
		"open":   {spawn.MetaState: spawn.CrewStateSpawned, spawn.MetaWorktree: ".worktrees/shop-open"},
	} {
		if err := w.WriteCrewMeta("shop", crew, meta); err != nil {
			t.Fatal(err)
		}
	}
	var dry []bool
	var got []string
	d := gcDeps{
		Orphans: func(_ context.Context, root string) []teardown.Target {
			if root != w.WorktreesDir() {
				t.Fatalf("orphans looked under %s", root)
			}
			return []teardown.Target{{Worktree: filepath.Join(root, "hellovietnam-old")}, {Worktree: filepath.Join(root, "shop-gone-1")}}
		},
		Collect: func(_ context.Context, ts []teardown.Target, d bool) teardown.Report {
			dry = append(dry, d)
			for _, tg := range ts {
				got = append(got, tg.Name())
			}
			return teardown.Report{Removed: []string{"docker shop-gone-1 (2 volumes)"}}
		},
	}

	targets, r := runGC(context.Background(), w, d, false)
	slices.Sort(got)
	if want := []string{"hellovietnam-old", "shop-gone-1", "shop-rev-r1"}; !slices.Equal(got, want) {
		t.Fatalf("targets = %q, want %q", got, want)
	}
	if !slices.Equal(dry, []bool{true}) {
		t.Fatalf("dry = %v, want one dry run", dry)
	}
	out := gcReport(targets, r, gcIdleReport{}, false)
	if !strings.Contains(out, "dry run, nothing removed") || !strings.Contains(out, "would remove: docker shop-gone-1 (2 volumes)") {
		t.Fatalf("report:\n%s", out)
	}

	if _, _ = runGC(context.Background(), w, d, true); !slices.Equal(dry, []bool{true, false}) {
		t.Fatalf("dry = %v, want --apply to remove", dry)
	}
}

func TestGCScheduleRunsAtFirstTickThenEveryPeriodOneAtATime(t *testing.T) {
	now := time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC)
	var runs atomic.Int32
	release := make(chan struct{})
	g := &gcSchedule{every: time.Hour, now: func() time.Time { return now }, run: func(context.Context) {
		runs.Add(1)
		<-release
	}}
	g.Tick(context.Background())
	g.Tick(context.Background()) // still running
	now = now.Add(2 * time.Hour)
	g.Tick(context.Background()) // due, but the first is still running
	close(release)
	waitFor(t, func() bool { g.mu.Lock(); defer g.mu.Unlock(); return !g.running })
	if runs.Load() != 1 {
		t.Fatalf("runs = %d, want 1 while one was running", runs.Load())
	}
	g.Tick(context.Background()) // last was set at 9:00, now 11:00: due
	waitFor(t, func() bool { return runs.Load() == 2 })
	g.Tick(context.Background()) // just ran
	time.Sleep(20 * time.Millisecond)
	if runs.Load() != 2 {
		t.Fatalf("runs = %d, want 2", runs.Load())
	}
	var nilSchedule *gcSchedule
	nilSchedule.Tick(context.Background())
}

func waitFor(t *testing.T, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !ok() {
		if time.Now().After(deadline) {
			t.Fatal("timed out")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// The idle pass shuts down the simulator of a crew that has waited on the
// Mate or the captain for gcIdleAfter, never of one still working, one
// that only just went quiet, one with no worktree, or a closed one; without
// --apply it only lists (CON-119).
func TestGCIdleShutsDownSimulatorsOfCrewsThatWaited(t *testing.T) {
	w := wakeFixture(t)
	for crew, meta := range map[string]map[string]string{
		"idle-1":  {spawn.MetaState: spawn.CrewStateSpawned, spawn.MetaWorktree: ".worktrees/shop-idle-1"},
		"fresh":   {spawn.MetaState: spawn.CrewStateSpawned, spawn.MetaWorktree: ".worktrees/shop-fresh"},
		"ask-k4":  {spawn.MetaState: spawn.CrewStateSpawned, spawn.MetaWorktree: ".worktrees/shop-ask-k4"},
		"busy-k5": {spawn.MetaState: spawn.CrewStateSpawned, spawn.MetaWorktree: ".worktrees/shop-busy-k5"},
		"old-wt":  {spawn.MetaState: spawn.CrewStateFinished, spawn.MetaWorktree: ".worktrees/shop-old-wt"},
	} {
		if err := w.WriteCrewMeta("shop", crew, meta); err != nil {
			t.Fatal(err)
		}
	}
	for crew, line := range map[string]string{
		"idle-1": "wait-mate: ready in branch mate/idle-1",
		"fresh":  "wait-mate: ready in branch mate/fresh",
		"old-wt": "wait-mate: long done",
	} {
		if err := w.AppendStatus("shop", crew, line); err != nil {
			t.Fatal(err)
		}
	}
	old := time.Now().Add(-2 * time.Hour)
	for _, crew := range []string{"idle-1", "ask-k4", "busy-k5", "old-wt"} {
		if err := os.Chtimes(w.CrewStatus("shop", crew), old, old); err != nil {
			t.Fatal(err)
		}
	}
	var dry []bool
	var got []string
	d := gcDeps{
		Collect: func(context.Context, []teardown.Target, bool) teardown.Report { return teardown.Report{} },
		Shutdown: func(_ context.Context, ts []teardown.Target, isDry bool) ([]string, error) {
			dry = append(dry, isDry)
			got = got[:0]
			for _, tg := range ts {
				got = append(got, tg.Name())
			}
			return []string{"simulator Shop-shop-idle-1"}, nil
		},
	}

	idle := runGCIdle(context.Background(), w, d, false)
	slices.Sort(got)
	if want := []string{"shop-ask-k4", "shop-idle-1"}; !slices.Equal(got, want) {
		t.Fatalf("targets = %q, want %q", got, want)
	}
	if !slices.Equal(dry, []bool{true}) {
		t.Fatalf("dry = %v, want one dry run", dry)
	}
	out := gcReport(nil, teardown.Report{}, idle, false)
	for _, want := range []string{"2 idle crew(s)", "shut down 1 simulator(s)", "would shut down: simulator Shop-shop-idle-1 (crew shop/idle-1 wait-mate 2h)"} {
		if !strings.Contains(out, want) {
			t.Fatalf("report lacks %q:\n%s", want, out)
		}
	}

	if _ = runGCIdle(context.Background(), w, d, true); !slices.Equal(dry, []bool{true, false}) {
		t.Fatalf("dry = %v, want --apply to shut down", dry)
	}
	if r := runGCIdle(context.Background(), w, gcDeps{}, true); len(r.Crews) != 0 {
		t.Fatalf("no Shutdown dep still found %d crews", len(r.Crews))
	}
}

func TestCompactDuration(t *testing.T) {
	for d, want := range map[time.Duration]string{45 * time.Minute: "45m", 2 * time.Hour: "2h", 2*time.Hour + 15*time.Minute + 40*time.Second: "2h16m", 30 * time.Second: "1m"} {
		if got := compactDuration(d); got != want {
			t.Errorf("compactDuration(%v) = %q, want %q", d, got, want)
		}
	}
}
