package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/query"
)

// gateCore is a core with a TUI flag and a force close that lowers it.
type gateCore struct {
	open   *query.TUIStatus
	acts   []coreapi.ActionRequest
	refuse error
}

func (c *gateCore) load(context.Context) (query.Snapshot, error) {
	return query.Snapshot{TUI: c.open}, nil
}

func (c *gateCore) act(_ context.Context, req coreapi.ActionRequest) (string, error) {
	c.acts = append(c.acts, req)
	if c.refuse != nil {
		return "", c.refuse
	}
	c.open = nil
	return "closed the TUI of shop; its Mate runs on in its host", nil
}

var shopTUI = query.TUIStatus{Project: "shop", PID: 4242, Since: time.Date(2026, 11, 6, 9, 30, 0, 0, time.Local)}

// A console that opens while a Mate's TUI is open asks first (CON-227):
// `y` has the core close it and the console goes on; N or Enter quits
// with the sentence the captain acts on; with no TUI it asks nothing.
func TestTheConsoleAsksAboutAnOpenTUI(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		answer  string
		wantErr string
		closes  bool
	}{
		{answer: "y\n", closes: true},
		{answer: "YES\n", closes: true},
		{answer: "n\n", wantErr: "close the TUI of shop first"},
		{answer: "\n", wantErr: "close the TUI of shop first"},
		{answer: "", wantErr: "close the TUI of shop first"},
	} {
		open := shopTUI
		c := &gateCore{open: &open}
		var out bytes.Buffer
		err := tuiGate{load: c.load, act: c.act, in: strings.NewReader(tc.answer), out: &out}.check(context.Background())
		if !strings.HasPrefix(out.String(), "Mate shop đang mở TUI (pid 4242, từ 09:30). Tắt TUI? [y/N] ") {
			t.Fatalf("answer %q: question = %q", tc.answer, out.String())
		}
		switch {
		case tc.closes && (err != nil || len(c.acts) != 1 || c.acts[0].Action != coreapi.ActionCloseTUI || c.acts[0].Target != "shop"):
			t.Fatalf("answer %q: %v, acts %+v; want the TUI closed", tc.answer, err, c.acts)
		case !tc.closes && (err == nil || err.Error() != tc.wantErr || len(c.acts) != 0):
			t.Fatalf("answer %q: %v, acts %+v; want %q and nothing sent", tc.answer, err, c.acts, tc.wantErr)
		}
	}

	c := &gateCore{}
	var out bytes.Buffer
	if err := (tuiGate{load: c.load, act: c.act, in: strings.NewReader("y\n"), out: &out}).check(context.Background()); err != nil || out.Len() != 0 {
		t.Fatalf("no TUI: %v, said %q", err, out.String())
	}

	open := shopTUI
	c = &gateCore{open: &open, refuse: errors.New("the TUI of shop (pid 4242) was sent SIGTERM and is still open after 10s; close it in its terminal")}
	if err := (tuiGate{load: c.load, act: c.act, in: strings.NewReader("y\n"), out: &out}).check(context.Background()); err == nil || !strings.Contains(err.Error(), "still open after 10s") {
		t.Fatalf("a TUI that stays open = %v, want the console to quit with why", err)
	}
}

// The console's reads tell it, once per TUI, that the flag went up while
// it was open, so it ends its own TUI and asks.
func TestTheConsoleSeesATUIOpenWhileItRuns(t *testing.T) {
	t.Parallel()
	c := &gateCore{}
	seen := make(chan struct{}, 1)
	load := watchTUI(c.load, seen)
	if _, err := load(context.Background()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-seen:
		t.Fatal("seen a TUI with none open")
	default:
	}
	open := shopTUI
	c.open = &open
	for range 3 {
		if _, err := load(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	select {
	case <-seen:
	default:
		t.Fatal("the open TUI was not seen")
	}
}
