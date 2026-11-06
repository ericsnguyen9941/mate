package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// cmdMateTUI implements `mate mate tui <project>` (CON-227): the headless
// Mate's session in the harness's own TUI, in this terminal, until the
// captain closes it. The `tui` script in the Mate's directory runs it, and
// the shell pane `mate mate start` opens beside a headless Mate starts
// there. It needs no core: the flag and the turn lock are files.
func cmdMateTUI(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("mate tui", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate mate tui <project> [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return newUsageError("mate mate tui: want exactly 1 argument: <project>")
	}
	project := fs.Arg(0)
	w, err := resolveWorkspace(*workspaceFlag)
	if err != nil {
		return err
	}
	run, err := tuiRun(w, project)
	if err != nil {
		return err
	}
	sigs := make(chan os.Signal, 1)
	// ^C is the TUI's own key; this process waits for the TUI to end.
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	defer signal.Stop(sigs)
	run.Stdin, run.Stdout, run.Stderr, run.Notes, run.Signals = os.Stdin, os.Stdout, os.Stderr, stderr, sigs
	return matehost.RunTUI(context.Background(), run)
}

// tuiRun is what `mate mate tui` runs for project's Mate: its host's
// launch, over this workspace's flag, turn lock and mate.meta.
func tuiRun(w *store.Workspace, project string) (matehost.TUIRun, error) {
	if _, ok := w.Project(project); !ok {
		return matehost.TUIRun{}, fmt.Errorf("mate mate tui: no project %s in %s", project, w.Root())
	}
	meta, err := w.ReadMateMeta(project)
	if err != nil {
		return matehost.TUIRun{}, err
	}
	if meta[spawn.MetaMode] != spawn.ModeHeadless {
		return matehost.TUIRun{}, fmt.Errorf("mate mate tui: the Mate of %s is not headless; its TUI is its own pane", project)
	}
	mateDir := w.MateDir(project)
	launch, err := matehost.ReadLaunch(mateDir)
	if err != nil {
		return matehost.TUIRun{}, fmt.Errorf("mate mate tui: the Mate of %s has no headless launch: %w (start it with `mate mate start %s`)", project, err, project)
	}
	profile, err := harnesses.Lookup(harness.Kind(launch.Harness))
	if err != nil {
		return matehost.TUIRun{}, err
	}
	resume, _ := profile.(harness.TUIResume)
	run := matehost.TUIRun{
		Project:  project,
		Launch:   launch,
		Resume:   resume,
		Flag:     func() (matehost.TUIFlag, error) { return takeTUIFlag(w, project) },
		TurnLock: func() (func(), bool, error) { return w.TryMateLock(project, matehost.TurnLock) },
		Session: func() (string, error) {
			meta, err := w.ReadMateMeta(project)
			return meta[spawn.MetaSessionID], err
		},
		SaveLaunch: func(l matehost.Launch) error { return matehost.WriteLaunch(mateDir, l) },
		MoveHost:   matehost.Dial(mateDir).SetSession,
	}
	if turns, ok := profile.(harness.HeadlessTurns); ok {
		run.Opened = turns.SessionOpened
	}
	return run, nil
}

// tuiFlag is the workspace's TUI flag held by `mate mate tui`.
type tuiFlag struct {
	lock   *store.TUILock
	holder store.TUIHolder
}

func (f *tuiFlag) Started(pid int) error {
	f.holder.HarnessPID = pid
	return f.lock.Write(f.holder)
}

func (f *tuiFlag) Release() { f.lock.Unlock() }

// takeTUIFlag takes `.mate/tui.lock` for project, or refuses with who
// holds it: "Mate shop đang mở TUI (pid 4242, từ 09:30)".
func takeTUIFlag(w *store.Workspace, project string) (*tuiFlag, error) {
	holder := store.TUIHolder{Project: project, PID: os.Getpid(), Started: time.Now()}
	lock, err := w.LockTUI(holder)
	if errors.Is(err, store.ErrTUIOpen) {
		if h, held, _ := w.TUIRunning(); held {
			return nil, errors.New(tuiStatus(h).Line())
		}
		return nil, errors.New("another Mate's TUI is open on this workspace")
	}
	if err != nil {
		return nil, err
	}
	return &tuiFlag{lock: lock, holder: holder}, nil
}

// tuiStatus is the flag's holder as the snapshot says it.
func tuiStatus(h store.TUIHolder) query.TUIStatus {
	return query.TUIStatus{Project: h.Project, PID: h.PID, Since: h.Started}
}
