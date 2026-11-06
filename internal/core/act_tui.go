package core

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/platform"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// While a headless Mate's session is open in the harness's TUI (CON-227),
// `mate mate tui` holds `.mate/tui.lock`, the workspace's TUI flag. The
// engine reads it for every UI: the snapshot carries it, and Act refuses
// every action but a peek - a read - and the force close itself. The
// engine's loops run on: the observer, the PR sweep, the auto daemon of
// every other Project, the watch over the hosts.

// tuiCloseWait is how long a force close waits for the flag to go.
const tuiCloseWait = 10 * time.Second

// tuiPassable are the actions the TUI flag lets through.
var tuiPassable = map[coreapi.Action]bool{coreapi.ActionPeek: true, coreapi.ActionCloseTUI: true}

// tuiOpen is the workspace's TUI flag, nil when no TUI is open or the
// engine has no workspace.
func tuiOpen(ws *store.Workspace) *query.TUIStatus {
	if ws == nil {
		return nil
	}
	h, held, err := ws.TUIRunning()
	if err != nil || !held {
		return nil
	}
	st := TUIStatus(h)
	return &st
}

// TUIStatus is the flag's holder as the UIs say it: the pid is the
// harness's, the TUI the captain sees, once it runs - a kill -9 of it ends
// the TUI and the flag with it - and `mate mate tui`'s before.
func TUIStatus(h store.TUIHolder) query.TUIStatus {
	pid := h.PID
	if h.HarnessPID > 0 {
		pid = h.HarnessPID
	}
	return query.TUIStatus{Project: h.Project, PID: pid, Since: h.Started}
}

// refuseForTUI is Act's refusal of req while a TUI is open, nil when it
// may run.
func (e *Engine) refuseForTUI(req coreapi.ActionRequest) error {
	if tuiPassable[req.Action] {
		return nil
	}
	if open := tuiOpen(e.o.Workspace); open != nil {
		return &coreapi.TUIError{TUI: *open}
	}
	return nil
}

// closeTUI is ActionCloseTUI: SIGTERM to the `mate mate tui` holding the
// flag, which passes it to the TUI, then up to tuiCloseWait for the flag
// to go. A signal goes only to the pid the lock's holder wrote, and only
// while the lock is held.
func (e *Engine) closeTUI(ctx context.Context, _ coreapi.ActionRequest) (string, error) {
	if e.o.Workspace == nil {
		return "", errNoWorkspace
	}
	h, held, err := e.o.Workspace.TUIRunning()
	if err != nil {
		return "", err
	}
	if !held {
		return "no Mate's TUI is open · nothing changed", nil
	}
	open := TUIStatus(h)
	// `mate mate tui` passes SIGTERM on to the harness and cleans up after
	// it; with it gone (Orphaned) the harness itself is told.
	target := h.PID
	if h.Orphaned {
		target = h.HarnessPID
	}
	if target <= 0 {
		return "", fmt.Errorf("%s, but its record names no process; close it in its terminal", open.Line())
	}
	if err := e.signalTUI(target, syscall.SIGTERM); err != nil {
		return "", fmt.Errorf("could not signal the TUI of %s (pid %d): %w", open.Project, target, err)
	}
	wait := tuiCloseWait
	if e.tuiWait > 0 {
		wait = e.tuiWait
	}
	deadline := time.NewTimer(wait)
	defer deadline.Stop()
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for {
		if tuiOpen(e.o.Workspace) == nil {
			return fmt.Sprintf("closed the TUI of %s; its Mate runs on in its host", open.Project), nil
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-deadline.C:
			return "", fmt.Errorf("the TUI of %s (pid %d) was sent SIGTERM and is still open after %s; close it in its terminal", open.Project, target, durationWord(wait))
		case <-tick.C:
		}
	}
}

// signalTUI sends sig to pid; a test stands in for it.
func (e *Engine) signalTUI(pid int, sig syscall.Signal) error {
	if e.tuiSignal != nil {
		return e.tuiSignal(pid, sig)
	}
	return platform.Signal(pid, sig)
}

// durationWord is d the way a line says it: "10s".
func durationWord(d time.Duration) string {
	if d < time.Second {
		return d.String()
	}
	return d.Round(time.Second).String()
}
