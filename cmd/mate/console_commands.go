package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/consolecmd"
	"github.com/nguyenngocanh94/mate/internal/host"
	"github.com/nguyenngocanh94/mate/internal/panerun"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/ui/console"
)

func consoleCommandRole(i int) string { return "command-" + strconv.Itoa(i) }

// GUI launchers can leave a background process holding their output pipes.
// Bound that wait too, so such a launcher cannot wedge a console action.
func runConsoleLauncher(ctx context.Context, inv consolecmd.Invocation, in io.Reader, out, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, inv.Name, inv.Args...)
	cmd.Dir = inv.Dir
	cmd.Stdin, cmd.Stdout, cmd.Stderr = in, out, stderr
	cmd.WaitDelay = time.Second
	err := cmd.Run()
	if errors.Is(err, exec.ErrWaitDelay) {
		// The launcher exited successfully; its background application kept
		// a pipe open. The GUI owns that process now.
		return nil
	}
	return err
}

func loadConsoleCommands(ws *store.Workspace) (consolecmd.Config, error) {
	path := filepath.Join(ws.StateDir(), "console.yaml")
	cfg, err := consolecmd.Load(path)
	if err != nil {
		return cfg, err
	}
	seen := map[string]bool{}
	for _, b := range cfg.Bindings {
		if slices.Contains(consolecmd.ReservedKeys(), b.Key) {
			return cfg, fmt.Errorf("%s: key %q is reserved by the console", path, b.Key)
		}
		slot := b.Scope + "/" + b.Key
		if seen[slot] {
			return cfg, fmt.Errorf("%s: key %q is already bound on %s rows", path, b.Key, b.Scope)
		}
		seen[slot] = true
	}
	return cfg, nil
}

// crewWorktree is the crew's git worktree as its meta records it, or ""
// when none is recorded or it is no longer on disk.
func crewWorktree(ws *store.Workspace, target console.StageTarget) string {
	meta, err := ws.ReadCrewMeta(target.ProjectID, target.ID)
	if err != nil || meta[spawn.MetaWorktree] == "" {
		return ""
	}
	dir := filepath.Join(ws.Root(), filepath.FromSlash(meta[spawn.MetaWorktree]))
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return ""
	}
	return dir
}

func consoleCommandBindings(cfg consolecmd.Config) []query.CommandBinding {
	var bindings []query.CommandBinding
	for i, b := range cfg.Bindings {
		bindings = append(bindings, query.CommandBinding{Key: b.Key, Label: b.Label, Scope: b.Scope, Role: consoleCommandRole(i)})
	}
	return bindings
}

// addCommandTabs plans terminal commands only. GUI launchers need neither
// a host nor a pane, and are never owned or closed by the console.
func (c *consoleColumns) addCommandTabs(cfg consolecmd.Config) error {
	for i, b := range cfg.Bindings {
		if b.Mode != "terminal" {
			continue
		}
		role := consoleCommandRole(i)
		// A separate short directory also works when there are many bindings.
		dir, err := panerun.SocketDir("mate-cmd-", "pane.sock")
		if err != nil {
			return err
		}
		c.commandDirs = append(c.commandDirs, dir)
		socket := filepath.Join(dir, "pane.sock")
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		idle := fmt.Sprintf("mate · %s\r\n\r\n%s on a %s opens this command.", b.Label, b.Key, b.Scope)
		c.tabs[role] = commandTab{socket: socket, col: host.Column{Role: role, Argv: []string{
			exe, "pane", "serve", "--role", role, "--socket", socket,
			"--owner", strconv.Itoa(os.Getpid()), "--idle", idle,
		}}}
	}
	return nil
}

// consoleCommandView launches commands from the workspace config.
// A GUI command works even when no supported terminal host is present.
func consoleCommandView(ws *store.Workspace, c *consoleColumns, cfg consolecmd.Config, runner consolecmd.Runner) console.CommandViewFunc {
	if len(cfg.Bindings) == 0 {
		return nil
	}
	return func(ctx context.Context, key string, target console.StageTarget) error {
		scope := ""
		if target.Kind == console.StageCrew && target.ID != "" {
			scope = "crew"
		} else if target.Kind == "" && target.ID == "" {
			scope = "project"
		}
		for i, b := range cfg.Bindings {
			if b.Key != key || b.Scope != scope {
				continue
			}
			if err := requireProject(ws, target.ProjectID); err != nil {
				return err
			}
			vars := map[string]string{"workspace_dir": ws.Root(), "project_dir": ws.ProjectHome(target.ProjectID)}
			if scope == "crew" {
				dir, report, err := reviewReport(ws, target)
				if err != nil {
					return err
				}
				vars["crew_dir"], vars["report_path"] = dir, report
				vars["worktree_dir"] = crewWorktree(ws, target)
			}
			args, dir, err := b.Resolve(vars)
			if err != nil {
				return fmt.Errorf("%s: %w", b.Label, err)
			}
			binary := findTool(os.Getenv, args[0])
			if c != nil {
				binary = c.findTool(args[0])
			}
			if !filepath.IsAbs(binary) {
				return fmt.Errorf("%s: command %q not found; install it or set an absolute executable path in .mate/console.yaml", b.Label, args[0])
			}
			args[0] = binary
			if b.Mode == "external" {
				// These commands are launchers, not editors with --wait. Capture
				// their output so it cannot corrupt the console's alt screen.
				launchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				var diagnostic bytes.Buffer
				err := runner(launchCtx, consolecmd.Invocation{Name: binary, Args: args[1:], Dir: dir}, nil, &diagnostic, &diagnostic)
				if err != nil {
					return fmt.Errorf("%s: %w: %s", b.Label, err, strings.TrimSpace(diagnostic.String()))
				}
				return nil
			}
			if c == nil {
				return fmt.Errorf("%s: terminal commands need a supported host; use mode: external for a GUI launcher", b.Label)
			}
			tab, ok := c.tabs[consoleCommandRole(i)]
			if !ok {
				return fmt.Errorf("%s: command tab is unavailable", b.Label)
			}
			return c.showRoleTab(ctx, tab.col, tab.socket, panerun.Command{Argv: args, Dir: dir, Env: c.env})
		}
		return fmt.Errorf("no command bound to %s on this row", key)
	}
}
