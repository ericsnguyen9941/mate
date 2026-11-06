package spawn

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// A headless Mate has no pane, so a start gives the captain one beside it
// (CON-227): a shell in the Mate's directory, in the workspace's Herdr
// session, where `./tui` opens the Mate's session in the harness's own
// TUI. The pane is the captain's: a start reuses the one there, a stop
// leaves it open.

// ShellPaneLabel is the Herdr workspace and tab of project's shell pane.
func ShellPaneLabel(project string) string { return "mate-" + project + "-shell" }

// TUIArgv is what the `tui` script in project's Mate directory runs.
func TUIArgv(binary, root, project string) []string {
	return []string{binary, "mate", "tui", project, "--workspace", root}
}

// writeTUIScript writes the `tui` script into the Mate's directory, again
// on every headless start, so it names the binary that started the Mate.
func writeTUIScript(w *store.Workspace, deps Deps, project, mateDir string) error {
	binary, err := deps.binary()
	if err != nil {
		return err
	}
	path := filepath.Join(mateDir, matehost.TUIName)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(matehost.TUIScript(TUIArgv(binary, w.Root(), project))), 0o755); err != nil {
		return err
	}
	if err := os.Chmod(tmp, 0o755); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// ensureMateShell opens project's shell pane in the workspace's Herdr
// session, cwd the Mate directory, or finds the one already open: a Herdr
// workspace and tab both labelled ShellPaneLabel. It carries no MATE_*
// variable: what the captain runs there is the captain's, not the Mate's.
//
// With noServer it starts no Herdr server: a session that is not running
// is errNoHerdrSession.
func ensureMateShell(ctx context.Context, deps Deps, spec runtime.SessionSpec, project, mateDir string, noServer bool) (runtime.TabHandle, error) {
	var session runtime.SessionHandle
	var err error
	if noServer {
		var found bool
		session, found, err = deps.Runtime.LookupSession(ctx, spec)
		if err == nil && !found {
			err = errNoHerdrSession
		}
	} else {
		session, err = deps.Runtime.EnsureSession(ctx, spec)
	}
	if err != nil {
		return runtime.TabHandle{}, err
	}
	label := ShellPaneLabel(project)
	ws, err := deps.Runtime.EnsureProjectWorkspace(ctx, runtime.WorkspaceSpec{Session: session, Label: label, Cwd: mateDir})
	if err != nil {
		return runtime.TabHandle{}, err
	}
	return deps.Runtime.CreateAgentTab(ctx, runtime.TabSpec{Workspace: ws, Label: label, Cwd: mateDir})
}

// errNoHerdrSession is a shell pane a restart of mate's own did not open:
// the workspace's Herdr session is not running, and it starts none.
var errNoHerdrSession = errors.New("the workspace's Herdr session is not running")

// withMateShell gives a headless Mate that started its `tui` script and
// its shell pane, and says in res what came of them. Neither fails the
// start: the Mate runs, and `mate mate tui` works from any terminal.
func withMateShell(ctx context.Context, w *store.Workspace, deps Deps, spec runtime.SessionSpec, req StartRequest, mateDir string, res *StartResult) {
	project := req.Project
	if err := writeTUIScript(w, deps, project, mateDir); err != nil {
		res.ShellNote = fmt.Sprintf("no %s script in %s (%s); run `mate mate tui %s`", matehost.TUIName, mateDir, oneLineErr(err), project)
		return
	}
	tab, err := ensureMateShell(ctx, deps, spec, project, mateDir, req.NoHerdrServer)
	if err != nil {
		res.ShellNote = fmt.Sprintf("no shell pane for the Mate's TUI (%s); run ./%s in %s from any terminal", oneLineErr(err), matehost.TUIName, mateDir)
		return
	}
	res.ShellPane = tab.PaneID
}
