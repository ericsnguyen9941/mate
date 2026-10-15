package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/matehost"
)

// hostHistory is how many items a host shows again when it starts.
const hostHistory = 200

// cmdMateHost implements `mate mate host <project>` (CON-78): the program a
// headless Mate's pane runs. spawn.StartMate starts it there; run by hand
// in a terminal it works the same, without telling Herdr anything.
func cmdMateHost(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("mate host", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate mate host <project> [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return newUsageError("mate mate host: want exactly 1 argument: <project>")
	}
	project := fs.Arg(0)
	w, err := resolveWorkspace(*workspaceFlag)
	if err != nil {
		return err
	}
	release, ok, err := w.TryMateLock(project, matehost.HostLock)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("mate mate host: the Mate of %s already has a host running", project)
	}
	defer release()
	mateDir := w.MateDir(project)
	launch, err := matehost.ReadLaunch(mateDir)
	if err != nil {
		return fmt.Errorf("mate mate host: %w (start it with `mate mate start %s`, which prepares the launch)", err, project)
	}
	profile, err := harnesses.Lookup(harness.Kind(launch.Harness))
	if err != nil {
		return err
	}
	turns, ok := profile.(harness.HeadlessTurns)
	if !ok {
		return fmt.Errorf("mate mate host: the %s harness cannot run a headless Mate", launch.Harness)
	}
	runner := &matehost.Runner{
		Launch: launch,
		Turns:  turns,
		Lock:   func() (func(), bool, error) { return w.TryMateLock(project, matehost.TurnLock) },
	}
	opts := matehost.Options{History: matehost.HistoryPath(mateDir)}
	reporter, inHerdr := matehost.NewHerdrReporter(launch.Agent)
	if inHerdr {
		opts.Reporter = reporter
	}
	host := matehost.New(runner, opts)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := matehost.Program(matehost.NewModel(host, project, matehost.History(opts.History, hostHistory)))
	serveErr := make(chan error, 1)
	go func() { serveErr <- host.Serve(ctx, matehost.SocketPath(mateDir), func() { p.Quit() }) }()
	if inHerdr {
		// After the socket is up, so whoever waits for the agent can also
		// reach the host.
		go func() {
			if err := waitSocket(ctx, mateDir); err == nil {
				if err := reporter.Register(); err != nil {
					p.Send(fmt.Errorf("could not register with Herdr: %w", err))
				}
			}
		}()
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-sigs
		host.Stop(0, true)
		p.Quit()
	}()

	_, runErr := p.Run()
	host.Stop(0, true)
	cancel()
	select {
	case err := <-serveErr:
		if err != nil && runErr == nil {
			runErr = err
		}
	case <-time.After(2 * time.Second):
	}
	return runErr
}

func waitSocket(ctx context.Context, mateDir string) error {
	c := matehost.Dial(mateDir)
	for {
		if err := c.Ping(ctx); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}
