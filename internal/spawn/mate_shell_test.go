package spawn_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness/claude"
	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/spawn"
)

// onlyShellPane fails when rt was asked for anything but the captain's
// shell pane beside a headless Mate: a session, a workspace and a tab
// (CON-227). No agent is started, inspected or stopped in Herdr.
func onlyShellPane(t *testing.T, rt *runtime.Fake, from int) {
	t.Helper()
	for _, c := range rt.Calls[from:] {
		switch strings.Fields(c)[0] {
		case "EnsureSession", "EnsureProjectWorkspace", "CreateAgentTab":
		case "LookupProjectWorkspace", "RunInPane", "StartAgent", "WaitAgent", "InspectAgent", "ListAgents", "StopAgent", "RemoveTab", "LookupSession":
			t.Fatalf("Herdr was asked %s: %q", c, rt.Calls[from:])
		}
	}
}

// shellTabs are the Fake's tabs labelled as project's shell pane.
func shellTabs(rt *runtime.Fake, project string) []runtime.TabHandle {
	var out []runtime.TabHandle
	for _, tab := range rt.Tabs {
		if tab.Label == spawn.ShellPaneLabel(project) {
			out = append(out, tab)
		}
	}
	return out
}

// A headless start opens one shell pane in the workspace's Herdr session,
// cwd the Mate's directory, and writes the `tui` script there; a start
// after a stop finds the same pane and opens no second one, and the stop
// leaves it open.
func TestHeadlessStartOpensOneShellPaneInTheMateDir(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv(claude.ClaudeConfigDirEnv, "")
	w := shortWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	playHost(t, w, &deps, "shop")
	mateDir := w.MateDir("shop")

	res, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "shop", Mode: spawn.ModeHeadless})
	if err != nil {
		t.Fatalf("StartMate: %v", err)
	}
	if res.ShellNote != "" || res.ShellPane == "" {
		t.Fatalf("result = %+v, want a shell pane and no note", res)
	}
	tabs := shellTabs(rt, "shop")
	if len(tabs) != 1 || tabs[0].Cwd != mateDir || tabs[0].PaneID != res.ShellPane || len(tabs[0].Env) != 0 {
		t.Fatalf("shell tabs = %+v, want one in %s with no MATE_* variable", tabs, mateDir)
	}
	if rt.WorkspaceCreates() != 1 {
		t.Fatalf("workspace creates = %d, want 1", rt.WorkspaceCreates())
	}
	script := filepath.Join(mateDir, matehost.TUIName)
	info, err := os.Stat(script)
	if err != nil || info.Mode().Perm()&0o100 == 0 {
		t.Fatalf("tui script = %v, %v; want an executable file", info, err)
	}
	data, _ := os.ReadFile(script)
	if want := "exec " + deps.Binary + " mate tui shop --workspace " + w.Root(); !strings.Contains(string(data), want) || !strings.HasPrefix(string(data), "#!/bin/sh\n") {
		t.Fatalf("tui script = %q, want it to %s", data, want)
	}

	calls := len(rt.Calls)
	if _, err := spawn.StopMate(context.Background(), w, deps, "shop"); err != nil {
		t.Fatalf("StopMate: %v", err)
	}
	if slices.Contains(rt.Calls[calls:], "RemoveTab") || len(shellTabs(rt, "shop")) != 1 {
		t.Fatalf("the stop closed the shell pane: %q", rt.Calls[calls:])
	}

	again, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "shop", Resume: true})
	if err != nil {
		t.Fatalf("StartMate again: %v", err)
	}
	if tabs := shellTabs(rt, "shop"); len(tabs) != 1 || again.ShellPane != res.ShellPane || rt.WorkspaceCreates() != 1 {
		t.Fatalf("after a second start: shell tabs %+v, pane %q, creates %d; want the first pane again", tabs, again.ShellPane, rt.WorkspaceCreates())
	}
}

// A Herdr that cannot open the pane leaves the Mate running and says why.
func TestHeadlessStartRunsWithoutItsShellPane(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv(claude.ClaudeConfigDirEnv, "")
	w := shortWorkspace(t, "shop")
	rt := runtime.NewFake()
	rt.CreateTabErr = errors.New("herdr: no server")
	deps := fakeDeps(t, rt)
	playHost(t, w, &deps, "shop")

	res, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "shop", Mode: spawn.ModeHeadless})
	if err != nil {
		t.Fatalf("StartMate: %v", err)
	}
	if !res.Headless || res.ShellPane != "" || !strings.Contains(res.ShellNote, "no shell pane") || !strings.Contains(res.ShellNote, "herdr: no server") {
		t.Fatalf("result = %+v, want the Mate running and the note", res)
	}
	if !w.MateHeadless("shop") {
		t.Fatal("the Mate is not recorded running")
	}
}

// A restart of mate's own (recovery, the core's watch) opens the shell pane
// only in a Herdr session already running: it never starts a server.
func TestARestartOfMatesOwnStartsNoHerdrServerForTheShellPane(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv(claude.ClaudeConfigDirEnv, "")
	w := shortWorkspace(t, "shop")
	rt := runtime.NewFake()
	rt.SessionNotRunning = true
	deps := fakeDeps(t, rt)
	playHost(t, w, &deps, "shop")

	res, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "shop", Mode: spawn.ModeHeadless, NoHerdrServer: true})
	if err != nil {
		t.Fatalf("StartMate: %v", err)
	}
	if slices.Contains(rt.Calls, "EnsureSession") || res.ShellPane != "" || !strings.Contains(res.ShellNote, "not running") {
		t.Fatalf("calls %q, result %+v; want no server started and the note", rt.Calls, res)
	}
	if _, err := os.Stat(filepath.Join(w.MateDir("shop"), matehost.TUIName)); err != nil {
		t.Fatalf("the tui script is written all the same: %v", err)
	}
}

// A Mate in its pane gets no shell pane: the pane is its TUI.
func TestPaneStartOpensNoShellPane(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv(claude.ClaudeConfigDirEnv, "")
	w := shortWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	res, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "shop", Mode: spawn.ModePane})
	if err != nil {
		t.Fatalf("StartMate: %v", err)
	}
	if res.ShellPane != "" || len(shellTabs(rt, "shop")) != 0 {
		t.Fatalf("a pane Mate got a shell pane: %+v", res)
	}
	if _, err := os.Stat(filepath.Join(w.MateDir("shop"), matehost.TUIName)); err == nil {
		t.Fatal("a pane Mate got a tui script")
	}
}
