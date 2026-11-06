package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

func tuiWorkspace(t *testing.T) *store.Workspace {
	t.Helper()
	w, err := store.Init(t.TempDir(), store.Defaults{MateHarness: "claude", CrewHarness: "claude"})
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"shop", "blog"} {
		if err := w.AddProject(p, store.ProjectConfig{}); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(w.MateDir(p), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return w
}

// `mate mate tui` takes the workspace's flag, refuses a second one with
// who holds it, and runs only over a headless Mate's launch, with the
// harness's TUI argv (CON-227). The flag pauses only that Mate's host.
func TestMateTUIRunsOverAHeadlessMatesLaunch(t *testing.T) {
	t.Parallel()
	w := tuiWorkspace(t)
	if _, err := tuiRun(w, "shop"); err == nil || !strings.Contains(err.Error(), "not headless") {
		t.Fatalf("a Mate with no headless record = %v, want a refusal", err)
	}
	if err := w.WriteMateMeta("shop", map[string]string{spawn.MetaMode: spawn.ModeHeadless, spawn.MetaHarness: "claude", spawn.MetaSessionID: "s-meta"}); err != nil {
		t.Fatal(err)
	}
	if _, err := tuiRun(w, "shop"); err == nil || !strings.Contains(err.Error(), "no headless launch") {
		t.Fatalf("a Mate with no launch = %v, want a refusal", err)
	}
	launch := matehost.Launch{Project: "shop", Harness: "claude", Program: "claude", Args: []string{"--model", "opus", "--session-id", "s-launch"}, Dir: w.MateDir("shop"), SessionID: "s-launch"}
	if err := matehost.WriteLaunch(w.MateDir("shop"), launch); err != nil {
		t.Fatal(err)
	}
	run, err := tuiRun(w, "shop")
	if err != nil {
		t.Fatalf("tuiRun: %v", err)
	}
	if run.Resume == nil || run.Opened == nil {
		t.Fatalf("run = %+v, want claude's TUI argv and session check", run)
	}
	if got := strings.Join(run.Resume.TUIArgs(run.Launch.Args, "s-meta", true), " "); got != "--model opus --resume s-meta" {
		t.Fatalf("TUI argv = %q", got)
	}
	if id, err := run.Session(); err != nil || id != "s-meta" {
		t.Fatalf("session = %q, %v; want mate.meta's", id, err)
	}

	shopPaused, blogPaused := hostPaused(w, "shop"), hostPaused(w, "blog")
	flag, err := run.Flag()
	if err != nil {
		t.Fatalf("flag: %v", err)
	}
	if _, err := takeTUIFlag(w, "blog"); err == nil || !strings.HasPrefix(err.Error(), "Mate shop đang mở TUI (pid ") {
		t.Fatalf("a second TUI = %v, want the holder named", err)
	}
	if !shopPaused() || blogPaused() {
		t.Fatalf("paused: shop %v, blog %v; want only shop's host paused", shopPaused(), blogPaused())
	}
	if err := flag.Started(99); err != nil {
		t.Fatal(err)
	}
	if h, held, _ := w.TUIRunning(); !held || h.Project != "shop" || h.HarnessPID != 99 || h.PID != os.Getpid() {
		t.Fatalf("flag record = %+v, %v", h, held)
	}
	flag.Release()
	if shopPaused() {
		t.Fatal("the host is still paused after the flag went")
	}
	if _, err := os.Stat(filepath.Join(w.MateDir("shop"), matehost.TUIName)); err == nil {
		t.Fatal("tuiRun wrote the script; only a start does")
	}
}
