package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/query"
)

// While a headless Mate's session is open in the harness's TUI (CON-227),
// the console does not act: it asks the captain, in the terminal, whether
// to close that TUI - when it opens, and whenever the flag goes up while
// it is open, which ends its TUI first. `y` has the core close the TUI and
// the console goes on; anything else quits it.

// tuiGate asks the captain about an open TUI before the console draws.
type tuiGate struct {
	// load is the core's snapshot; act runs the force close.
	load coreapi.LoadFunc
	act  coreapi.ActionFunc
	// in is the captain's answer, out where the question goes.
	in  io.Reader
	out io.Writer
}

// check asks about the open TUI, if one is open, and closes it on `y`. It
// returns an error when the console must quit: the captain said no, or the
// TUI did not close.
func (g tuiGate) check(ctx context.Context) error {
	snap, err := g.load(ctx)
	if err != nil || snap.TUI == nil {
		// A snapshot that cannot be read is the console's to say.
		return nil
	}
	return g.ask(ctx, *snap.TUI)
}

func (g tuiGate) ask(ctx context.Context, open query.TUIStatus) error {
	fmt.Fprintf(g.out, "%s. Tắt TUI? [y/N] ", open.Line())
	line, _ := bufio.NewReader(g.in).ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
	default:
		return fmt.Errorf("close the TUI of %s first", open.Project)
	}
	said, err := g.act(ctx, coreapi.ActionRequest{Action: coreapi.ActionCloseTUI, Target: open.Project, TargetKind: "workspace"})
	if err != nil {
		return fmt.Errorf("mate console: %w", err)
	}
	fmt.Fprintln(g.out, said)
	return nil
}

// watchTUI is load telling seen, once, when the snapshot it read says a
// TUI is open: the console then ends its TUI and asks.
func watchTUI(load coreapi.LoadFunc, seen chan<- struct{}) coreapi.LoadFunc {
	return func(ctx context.Context) (query.Snapshot, error) {
		snap, err := load(ctx)
		if err == nil && snap.TUI != nil {
			select {
			case seen <- struct{}{}:
			default:
			}
		}
		return snap, err
	}
}
