package outbox_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness/claude"
	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/outbox"
	"github.com/nguyenngocanh94/mate/internal/send"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// headlessFixture is a fixture whose Mate runs headless (CON-78): its
// workspace is under /tmp so the host's socket path fits, and a host
// answers on it. lock decides whether the host's turn lock is free.
func headlessFixture(t *testing.T, lockFree bool) (*fixture, string) {
	t.Helper()
	f := newFixture(t)
	root, err := os.MkdirTemp("/tmp", "ob")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(root) })
	ws, err := store.Init(root, store.Defaults{})
	if err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(ws.ProjectHome(project), project)
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ws.AddProject(project, store.ProjectConfig{Repos: []store.RepoConfig{{Path: repo, DefaultBranch: "main"}}}); err != nil {
		t.Fatal(err)
	}
	f.ws = ws
	mateDir := ws.MateDir(project)
	if err := os.MkdirAll(mateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMateMeta(project, map[string]string{"harness": "claude", "agent": "mate-" + project, store.MetaMateMode: store.MateModeHeadless}); err != nil {
		t.Fatal(err)
	}
	prompts := filepath.Join(mateDir, "prompts.log")
	prog := filepath.Join(mateDir, "claude")
	script := "#!/bin/sh\ncat >> " + prompts + "\necho >> " + prompts + "\necho '{\"type\":\"result\",\"result\":\"—\",\"session_id\":\"s\"}'\n"
	if err := os.WriteFile(prog, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	runner := &matehost.Runner{
		Launch: matehost.Launch{Project: project, Program: prog, Dir: mateDir, SessionID: "s"},
		Turns:  claude.Claude{},
		Lock:   func() (func(), bool, error) { return func() {}, lockFree, nil },
		Opened: func(string, string) bool { return false },
	}
	host := matehost.New(runner, matehost.Options{})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go func() { _ = host.Serve(ctx, matehost.SocketPath(mateDir), nil) }()
	deadline := time.Now().Add(5 * time.Second)
	for matehost.Dial(mateDir).Ping(ctx) != nil {
		if time.Now().After(deadline) {
			t.Fatal("host never answered")
		}
		time.Sleep(10 * time.Millisecond)
	}
	return f, prompts
}

func TestAHeadlessMateGetsTheLineThroughItsHost(t *testing.T) {
	f, prompts := headlessFixture(t, true)
	s := f.sender()
	f.assign(s, "crews/k3.status@0", resolveLine)
	f.drain(s)

	if typed := f.typed(); len(typed) != 0 {
		t.Fatalf("typed %#v into a pane; a headless Mate has none to type into", typed)
	}
	items := f.items()
	if len(items) != 1 || items[0].State != store.OutboxSent {
		t.Fatalf("outbox = %+v, want the line sent", items)
	}
	sent := f.sent()
	if len(sent) != 1 || sent[0].Source != store.SourceApp || sent[0].Text != resolveLine {
		t.Fatalf("sent.log = %+v", sent)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		data, _ := os.ReadFile(prompts)
		if strings.TrimSpace(string(data)) == send.Marker+resolveLine {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("the host ran %q, want the marked line", data)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestABusyHeadlessMateKeepsTheLineQueued(t *testing.T) {
	f, _ := headlessFixture(t, false)
	s := f.sender()
	f.assign(s, "crews/k3.status@0", resolveLine)
	f.drain(s)
	f.clock.Advance(outbox.DefaultInterval)
	f.drain(s)

	items := f.items()
	if len(items) != 1 || !items[0].Queued() || items[0].Attempts != 2 {
		t.Fatalf("outbox = %+v, want one queued item tried twice", items)
	}
	if !strings.Contains(items[0].LastRefusal, "busy") {
		t.Fatalf("refusal = %q, want the host's busy answer", items[0].LastRefusal)
	}
	if sent := f.sent(); len(sent) != 0 {
		t.Fatalf("sent.log = %+v before anything was delivered", sent)
	}
}
