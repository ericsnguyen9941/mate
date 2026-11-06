package core

import (
	"context"
	"errors"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// While a Mate's session is open in its TUI (CON-227) the engine refuses
// every action from every UI, the chat included and on every Project,
// with the flag's line; a peek is a read and runs. The snapshot says who
// holds the flag. Once the flag goes, actions run again.
func TestActRefusesEveryActionButAPeekWhileATUIIsOpen(t *testing.T) {
	w := shopWorkspace(t)
	if err := w.AddProject("blog", store.ProjectConfig{}); err != nil {
		t.Fatal(err)
	}
	h := newHeldActs()
	e := New(Options{Workspace: w, Act: h.act})
	defer e.Stop()
	since := time.Date(2026, 11, 6, 9, 30, 0, 0, time.Local)
	lock, err := w.LockTUI(store.TUIHolder{Project: "shop", PID: 4242, Started: since})
	if err != nil {
		t.Fatal(err)
	}
	for _, req := range []coreapi.ActionRequest{
		stopMate, mergeCrew, addProject,
		{Action: coreapi.ActionSay, Target: "shop", TargetKind: "mate", Input: "hi"},
		{Action: coreapi.ActionStart, Target: "blog", TargetKind: "mate"},
		{Action: coreapi.ActionOkayWork, Target: "blog", TargetKind: "project", Input: "1"},
	} {
		_, err := e.Act(context.Background(), req)
		var open *coreapi.TUIError
		if !errors.Is(err, coreapi.ErrTUI) || !errors.As(err, &open) || open.TUI.Project != "shop" || open.TUI.PID != 4242 {
			t.Fatalf("%s on %s with shop's TUI open = %v, want a TUI refusal", req.Action, req.Target, err)
		}
		if want := "Mate shop đang mở TUI (pid 4242, từ 09:30) · nothing started"; err.Error() != want {
			t.Fatalf("refusal = %q, want %q", err, want)
		}
	}
	if h.count() != 0 {
		t.Fatalf("%d actions ran with the TUI open", h.count())
	}
	if _, err := e.Act(context.Background(), coreapi.ActionRequest{Action: coreapi.ActionPeek, Target: "shop", TargetKind: "crew", Crew: "buy"}); err != nil {
		t.Fatalf("peek with the TUI open: %v", err)
	}
	snap, err := e.FreshSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snap.TUI == nil || snap.TUI.Project != "shop" || snap.TUI.PID != 4242 || !snap.TUI.Since.Equal(since) {
		t.Fatalf("snapshot TUI = %+v, want shop's flag", snap.TUI)
	}

	lock.Unlock()
	if _, err := e.Act(context.Background(), stopMate); err != nil {
		t.Fatalf("stop once the TUI closed: %v", err)
	}
	if snap, err := e.FreshSnapshot(context.Background()); err != nil || snap.TUI != nil {
		t.Fatalf("snapshot TUI after it closed = %+v, %v", snap.TUI, err)
	}
}

// The force close sends SIGTERM to the pid the flag's holder wrote and
// waits for the flag to go; a TUI that does not close in time is said, and
// with no TUI open nothing is sent.
func TestCloseTUISignalsTheHolderAndWaitsForTheFlag(t *testing.T) {
	w := shopWorkspace(t)
	e := New(Options{Workspace: w})
	defer e.Stop()
	closeReq := coreapi.ActionRequest{Action: coreapi.ActionCloseTUI, Target: "shop", TargetKind: "workspace"}

	var sent []string
	e.tuiSignal = func(pid int, sig syscall.Signal) error {
		sent = append(sent, sig.String())
		return nil
	}
	line, err := e.Act(context.Background(), closeReq)
	if err != nil || !strings.Contains(line, "no Mate's TUI is open") || len(sent) != 0 {
		t.Fatalf("close with no TUI = %q, %v, sent %v", line, err, sent)
	}

	lock, err := w.LockTUI(store.TUIHolder{Project: "shop", PID: 4242, Started: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	e.tuiSignal = func(pid int, sig syscall.Signal) error {
		if pid != 4242 || sig != syscall.SIGTERM {
			t.Errorf("signal %v to %d, want SIGTERM to 4242", sig, pid)
		}
		go func() { time.Sleep(200 * time.Millisecond); lock.Unlock() }()
		return nil
	}
	line, err = e.Act(context.Background(), closeReq)
	if err != nil || line != "closed the TUI of shop; its Mate runs on in its host" {
		t.Fatalf("close = %q, %v", line, err)
	}
	if _, held, _ := w.TUIRunning(); held {
		t.Fatal("the flag is still up")
	}
}

func TestCloseTUISaysATUIThatStaysOpen(t *testing.T) {
	w := shopWorkspace(t)
	e := New(Options{Workspace: w})
	defer e.Stop()
	e.tuiWait = 300 * time.Millisecond
	lock, err := w.LockTUI(store.TUIHolder{Project: "shop", PID: 4242, Started: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Unlock()
	e.tuiSignal = func(int, syscall.Signal) error { return nil }
	_, err = e.Act(context.Background(), coreapi.ActionRequest{Action: coreapi.ActionCloseTUI, Target: "shop", TargetKind: "workspace"})
	if err == nil || !strings.Contains(err.Error(), "still open after 300ms") {
		t.Fatalf("close of a TUI that stays = %v", err)
	}
}
