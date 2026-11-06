package matehost

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// A headless Mate's session can be opened in the harness's own TUI for a
// while (CON-227): `mate mate-agent tui <project>`, which the `tui` script in the
// Mate's directory runs. The host stays alive and runs no turn meanwhile,
// so the core's watch over it (CON-226) has nothing to restart, and only
// one process writes the session at a time: the TUI holds the turn lock,
// and the workspace's TUI flag (`.mate/tui.lock`) keeps the host from
// starting a turn and the UIs from acting. When the TUI ends, however it
// ends, the flag goes and the host runs on, on the session the TUI left.

// TUIName is the script in the Mate's directory that opens the TUI.
const TUIName = "tui"

// TUIScript is the `tui` script's text: a shell over `mate mate-agent tui`.
// argv is the command, quoted word by word for /bin/sh.
func TUIScript(argv []string) string {
	words := make([]string, len(argv))
	for i, a := range argv {
		words[i] = shellQuote(a)
	}
	return "#!/bin/sh\n# Opens the Mate's session in its harness's TUI; the host waits meanwhile (CON-227).\nexec " + strings.Join(words, " ") + "\n"
}

func shellQuote(s string) string {
	if s != "" && strings.IndexFunc(s, func(r rune) bool {
		return !(r == '/' || r == '-' || r == '_' || r == '.' || r == ':' || r == '=' || r == '+' || r == ',' ||
			(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
	}) < 0 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// TUIFlag is the workspace's TUI flag, held by this process.
type TUIFlag interface {
	// Started records the harness's process once it runs.
	Started(harnessPID int) error
	// Release lets the flag go.
	Release()
}

// TUIRun is what RunTUI opens and how it reaches the rest of mate.
type TUIRun struct {
	Project string
	// Launch is the host's launch (host.launch.json): the TUI runs its
	// program, args, directory and environment.
	Launch Launch
	// Resume turns the launch into the TUI's argv.
	Resume harness.TUIResume
	// Opened reports whether a turn has created the session; nil is
	// always.
	Opened func(dir, sessionID string) bool
	// Flag takes the workspace's TUI flag, or says who holds it.
	Flag func() (TUIFlag, error)
	// TurnLock is the Mate's turn lock (store.TryMateLock(project,
	// TurnLock)), held while the TUI runs.
	TurnLock LockFunc
	// Session is the Mate's session id as mate.meta records it now: the
	// Mate's hooks record the session the TUI runs in, as they do a pane
	// Mate's. Empty when it records none.
	Session func() (string, error)
	// SaveLaunch records the launch the host runs from now on
	// (WriteLaunch).
	SaveLaunch func(Launch) error
	// MoveHost tells the running host to resume another session
	// (Client.SetSession).
	MoveHost func(ctx context.Context, sessionID string) error
	// Stdin, Stdout and Stderr are the terminal the TUI runs on.
	Stdin         io.Reader
	Stdout, Notes io.Writer
	Stderr        io.Writer
	// Signals are the signals sent to this process: SIGTERM and SIGHUP are
	// passed on to the TUI, and while the turn in flight is waited for,
	// any of them gives up.
	Signals <-chan os.Signal
	// Poll is how often the turn lock is tried; zero is 200 ms. Grace is
	// how long a TUI passed SIGTERM has before it is killed; zero is 6 s,
	// so with the host told (moveWait) the flag goes within the 10 s a
	// force close waits.
	Poll, Grace time.Duration
}

// RunTUI opens the Mate's session in the harness's TUI on this terminal
// and returns once it has ended (CON-227), in this order:
//
//  1. it takes the workspace's TUI flag, or refuses with its holder;
//  2. it waits for the turn in flight, saying so, and holds the turn lock;
//  3. it runs the launch's program with Resume's argv on the session
//     mate.meta records;
//  4. once the program ends - an exit, ^C, a crash, a kill - it moves the
//     host to the session mate.meta records then, when that is another,
//     and lets the turn lock and the flag go.
//
// The host then starts what waited, and the outbox delivers.
func RunTUI(ctx context.Context, r TUIRun) error {
	if r.Resume == nil {
		return fmt.Errorf("the %s harness cannot open a Mate's session in its TUI", r.Launch.Harness)
	}
	flag, err := r.Flag()
	if err != nil {
		return err
	}
	defer flag.Release()
	release, err := r.waitTurn(ctx)
	if err != nil {
		return err
	}
	defer release()

	sid := r.session()
	opened := r.Opened == nil || r.Opened(r.Launch.Dir, sid)
	cmd := exec.Command(r.Launch.Program, r.Resume.TUIArgs(r.Launch.Args, sid, opened)...)
	cmd.Dir = r.Launch.Dir
	cmd.Env = launchEnv(os.Environ(), r.Launch)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = r.Stdin, r.Stdout, r.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("mate mate-agent tui: start %s: %w", r.Launch.Program, err)
	}
	_ = flag.Started(cmd.Process.Pid)
	done := make(chan struct{})
	go r.forward(cmd.Process, done)
	waitErr := cmd.Wait()
	close(done)

	after := r.session()
	if after != r.Launch.SessionID {
		r.move(ctx, after)
	}
	status := "closed"
	if waitErr != nil {
		status = "ended (" + waitErr.Error() + ")"
	}
	fmt.Fprintf(r.Notes, "mate: TUI %s; Mate %s runs on in its host, session %s\n", status, r.Project, after)
	return nil
}

// session is mate.meta's session id, or the launch's when it records none.
func (r TUIRun) session() string {
	if r.Session != nil {
		if id, err := r.Session(); err == nil && strings.TrimSpace(id) != "" {
			return strings.TrimSpace(id)
		}
	}
	return r.Launch.SessionID
}

// waitTurn takes the turn lock, waiting for a turn in flight to end.
func (r TUIRun) waitTurn(ctx context.Context) (func(), error) {
	poll := r.Poll
	if poll <= 0 {
		poll = 200 * time.Millisecond
	}
	told := false
	for {
		release, ok, err := r.TurnLock()
		if err != nil {
			return nil, fmt.Errorf("mate mate-agent tui: turn lock: %w", err)
		}
		if ok {
			return release, nil
		}
		if !told {
			fmt.Fprintf(r.Notes, "mate: Mate %s is in the middle of a turn; the TUI opens once it ends (^C gives up)\n", r.Project)
			told = true
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case sig := <-r.Signals:
			return nil, fmt.Errorf("mate mate-agent tui: %v while waiting for the turn; the TUI was not opened", sig)
		case <-time.After(poll):
		}
	}
}

// forward passes SIGTERM and SIGHUP on to the TUI, and kills it when it is
// still running Grace after. SIGINT is the TUI's own key: the terminal
// sends it to the TUI too, and this process stays.
func (r TUIRun) forward(p *os.Process, done <-chan struct{}) {
	grace := r.Grace
	if grace <= 0 {
		grace = 6 * time.Second
	}
	var kill <-chan time.Time
	for {
		select {
		case <-done:
			return
		case <-kill:
			_ = p.Kill()
			return
		case sig := <-r.Signals:
			if sig == syscall.SIGTERM || sig == syscall.SIGHUP {
				_ = p.Signal(sig)
				if kill == nil {
					kill = time.After(grace)
				}
			}
		}
	}
}

// move points the host at session id: the launch first, so a host started
// again reads it, then the running host. A host that is not running has
// nothing to be told.
func (r TUIRun) move(ctx context.Context, id string) {
	l := r.Launch
	l.SessionID = id
	if r.SaveLaunch != nil {
		if err := r.SaveLaunch(l); err != nil {
			fmt.Fprintf(r.Notes, "mate: the host's launch could not record session %s: %v\n", id, err)
		}
	}
	if r.MoveHost == nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, moveWait)
	defer cancel()
	err := r.MoveHost(ctx, id)
	var dial *net.OpError
	switch {
	case err == nil:
		fmt.Fprintf(r.Notes, "mate: the TUI left the conversation in session %s; the host resumes it\n", id)
	case errors.As(err, &dial) && dial.Op == "dial":
		// No host runs: the next one reads the launch.
	default:
		fmt.Fprintf(r.Notes, "mate: the host could not be moved to session %s (%v); restart the Mate to resume it\n", id, err)
	}
}

// moveWait bounds telling the host its new session.
const moveWait = 2 * time.Second

// launchEnv is the environment a launch runs with: base without what the
// launch unsets, the nested-session variables of whatever Claude Code
// session base came from (harness.NestedSessionEnv), and EnvHostTurn, then
// the launch's own variables.
func launchEnv(base []string, l Launch) []string {
	drop := func(key string) bool {
		if harness.IsNestedSessionEnv(key) || key == EnvHostTurn {
			return true
		}
		for _, k := range l.UnsetEnv {
			if k == key {
				return true
			}
		}
		for _, kv := range l.Env {
			if k, _, _ := strings.Cut(kv, "="); k == key {
				return true
			}
		}
		return false
	}
	out := make([]string, 0, len(base)+len(l.Env))
	for _, kv := range base {
		if k, _, _ := strings.Cut(kv, "="); !drop(k) {
			out = append(out, kv)
		}
	}
	return append(out, l.Env...)
}
