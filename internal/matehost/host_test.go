package matehost

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// fakeClaude answers like `claude -p --output-format stream-json`: it echoes
// its argv and stdin into the result, sleeps while the prompt says "slow",
// answers "—" to a prompt that says "quiet", and fails on "crash".
const fakeClaude = `#!/bin/sh
prompt=$(cat)
echo "$*" >> "$(dirname "$0")/argv.log"
echo "$prompt" >> "$(dirname "$0")/prompts.log"
case "$prompt" in *crash*) echo "boom" >&2; exit 3;; esac
echo '{"type":"system","subtype":"init","session_id":"s1"}'
echo '{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash"}]}}'
case "$prompt" in *slow*) sleep 2;; esac
case "$prompt" in *quiet*) reply='—';; *) reply="re: $prompt";; esac
printf '{"type":"result","subtype":"success","result":"%s","session_id":"s1","is_error":false}\n' "$reply"
`

func fileLock(path string) LockFunc {
	return func() (func(), bool, error) {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
		if err != nil {
			return nil, false, err
		}
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
			f.Close()
			if errors.Is(err, syscall.EWOULDBLOCK) {
				return nil, false, nil
			}
			return nil, false, err
		}
		return func() { _ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN); f.Close() }, true, nil
	}
}

type fixture struct {
	dir    string
	host   *Host
	opened map[string]bool

	mu     sync.Mutex
	events []Event
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	dir := t.TempDir()
	prog := filepath.Join(dir, "claude")
	if err := os.WriteFile(prog, []byte(fakeClaude), 0o755); err != nil {
		t.Fatal(err)
	}
	f := &fixture{dir: dir, opened: map[string]bool{}}
	runner := &Runner{
		Launch: Launch{Project: "p", Program: prog, Args: []string{"--model", "opus", "--session-id", "s1"}, Dir: dir, SessionID: "s1"},
		Lock:   fileLock(filepath.Join(dir, ".turn.lock")),
		Opened: func(_, id string) bool {
			f.mu.Lock()
			defer f.mu.Unlock()
			return f.opened[id]
		},
		Grace: 200 * time.Millisecond,
	}
	f.host = New(runner, Options{History: filepath.Join(dir, HistoryName)})
	f.host.Subscribe(func(ev Event) {
		f.mu.Lock()
		f.events = append(f.events, ev)
		if ev.Item != nil && ev.Item.Done && ev.Item.Error == "" {
			f.opened["s1"] = true
		}
		f.mu.Unlock()
	})
	return f
}

// done waits for n finished items and returns them in order.
func (f *fixture) done(t *testing.T, n int) []Item {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		f.mu.Lock()
		var items []Item
		for _, ev := range f.events {
			if ev.Item != nil && ev.Item.Done {
				items = append(items, *ev.Item)
			}
		}
		f.mu.Unlock()
		if len(items) >= n {
			return items
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("fewer than %d turns finished", n)
	return nil
}

func (f *fixture) read(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestHostCaptainTurnCreatesThenResumesSession(t *testing.T) {
	f := newFixture(t)
	f.host.Submit("xin chào")
	f.done(t, 1)
	f.host.Submit("lần hai")
	items := f.done(t, 2)
	if items[0].Reply != "re: xin chào" || items[1].Reply != "re: lần hai" {
		t.Fatalf("replies = %q, %q", items[0].Reply, items[1].Reply)
	}
	argv := strings.Split(strings.TrimSpace(f.read(t, "argv.log")), "\n")
	if len(argv) != 2 {
		t.Fatalf("argv.log = %q", argv)
	}
	if !strings.Contains(argv[0], "--session-id s1") || strings.Contains(argv[0], "--resume") {
		t.Errorf("first turn argv = %q, want --session-id", argv[0])
	}
	if !strings.Contains(argv[1], "--resume s1") || strings.Contains(argv[1], "--session-id") {
		t.Errorf("second turn argv = %q, want --resume", argv[1])
	}
	for _, a := range argv {
		if !strings.HasSuffix(a, "-p --output-format stream-json --verbose") || !strings.HasPrefix(a, "--model opus") {
			t.Errorf("argv = %q", a)
		}
	}
	if got := History(HistoryPath(f.dir), 0); len(got) != 2 || got[1].Text != "lần hai" {
		t.Errorf("history = %+v", got)
	}
}

func TestHostQueuesCaptainAndRefusesAppWhileBusy(t *testing.T) {
	f := newFixture(t)
	f.host.Submit("slow one")
	waitBusy(t, f.host)
	f.host.Submit("second")
	if s := f.host.State(); !s.Busy || s.Queued != 1 {
		t.Fatalf("state = %+v, want busy with one queued", s)
	}
	if err := f.host.Deliver(Marker + "digest: x"); !errors.Is(err, ErrBusy) {
		t.Fatalf("deliver while busy = %v, want ErrBusy", err)
	}
	items := f.done(t, 2)
	if items[0].Text != "slow one" || items[1].Text != "second" {
		t.Fatalf("order = %q, %q", items[0].Text, items[1].Text)
	}
	// One run at a time: the second prompt was read after the first ended.
	if got := f.read(t, "prompts.log"); got != "slow one\nsecond\n" {
		t.Errorf("prompts = %q", got)
	}
}

func TestHostRefusesWhenAnotherProcessHoldsTheTurnLock(t *testing.T) {
	f := newFixture(t)
	release, ok, err := fileLock(filepath.Join(f.dir, ".turn.lock"))()
	if err != nil || !ok {
		t.Fatalf("lock: %v %v", ok, err)
	}
	if err := f.host.Deliver(Marker + "digest: x"); !errors.Is(err, ErrBusy) {
		t.Fatalf("deliver under a held lock = %v, want ErrBusy", err)
	}
	f.host.Submit("waits")
	if s := f.host.State(); s.Queued != 1 || s.Busy {
		t.Fatalf("state = %+v, want the message queued", s)
	}
	release()
	items := f.done(t, 1)
	if items[0].Text != "waits" {
		t.Fatalf("item = %+v", items[0])
	}
}

func TestHostAppTurnLabelAndSilence(t *testing.T) {
	f := newFixture(t)
	if err := f.host.Deliver(Marker + "digest: 1 item(s) quiet"); err != nil {
		t.Fatal(err)
	}
	items := f.done(t, 1)
	if it := items[0]; it.Source != SourceApp || it.Label != "digest" || !it.Silent || it.Visible() {
		t.Fatalf("item = %+v, want a silent digest", it)
	}
	if err := f.host.Deliver(Marker + "resolve: incident stale"); err != nil {
		t.Fatal(err)
	}
	items = f.done(t, 2)
	if it := items[1]; it.Label != "resolve" || it.Silent {
		t.Fatalf("item = %+v, want a visible resolve", it)
	}
}

func TestHostReportsAFailedRun(t *testing.T) {
	f := newFixture(t)
	f.host.Submit("crash now")
	items := f.done(t, 1)
	if !strings.Contains(items[0].Error, "boom") {
		t.Fatalf("error = %q, want the stderr tail", items[0].Error)
	}
	if f.host.State().Busy {
		t.Fatal("still busy after a failed run")
	}
}

func TestHostInterruptAndStop(t *testing.T) {
	f := newFixture(t)
	f.host.Submit("slow")
	waitBusy(t, f.host)
	if !f.host.Interrupt() {
		t.Fatal("interrupt found no turn")
	}
	if it := f.done(t, 1)[0]; !it.Interrupted {
		t.Fatalf("item = %+v, want interrupted", it)
	}

	f.host.Submit("slow again")
	waitBusy(t, f.host)
	f.host.Submit("queued")
	if killed := f.host.Stop(5*time.Second, false); killed {
		t.Fatal("a stop with time to wait killed the turn")
	}
	if err := f.host.Deliver(Marker + "digest: x"); !errors.Is(err, ErrStopping) {
		t.Fatalf("deliver after stop = %v", err)
	}
	time.Sleep(200 * time.Millisecond)
	if got := len(f.done(t, 2)); got != 2 {
		t.Fatalf("the queued message ran after stop: %d turns", got)
	}
}

func TestSocketRoundTrip(t *testing.T) {
	f := newFixture(t)
	sock, err := os.MkdirTemp("/tmp", "mh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(sock)
	path := filepath.Join(sock, SocketName)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stopped := make(chan struct{})
	go func() { _ = f.host.Serve(ctx, path, func() { close(stopped) }) }()
	c := Client{Path: path}
	deadline := time.Now().Add(5 * time.Second)
	for c.Ping(ctx) != nil {
		if time.Now().After(deadline) {
			t.Fatal("socket never answered")
		}
		time.Sleep(20 * time.Millisecond)
	}
	f.host.SetDraft("chưa gửi")
	if s, err := c.State(ctx); err != nil || s.Draft != "chưa gửi" || s.Busy {
		t.Fatalf("state = %+v, %v", s, err)
	}
	if err := c.ClearDraft(ctx); err != nil {
		t.Fatal(err)
	}
	if err := c.Deliver(ctx, Marker+"digest: slow"); err != nil {
		t.Fatal(err)
	}
	var busy *BusyError
	if err := c.Deliver(ctx, Marker+"digest: again"); !errors.As(err, &busy) {
		t.Fatalf("second deliver = %v, want *BusyError", err)
	}
	if s, _ := c.State(ctx); !s.Busy || s.Draft != "" || s.TurnSource != SourceApp {
		t.Fatalf("state = %+v", s)
	}
	f.done(t, 1)
	if items, err := c.History(ctx, 5); err != nil || len(items) != 1 {
		t.Fatalf("history = %+v, %v", items, err)
	}
	if _, err := c.Stop(ctx, time.Second, false); err != nil {
		t.Fatal(err)
	}
	select {
	case <-stopped:
	case <-time.After(5 * time.Second):
		t.Fatal("stop did not call back")
	}
}

func TestTurnEnvDropsNestedSessionAndAddsLaunchEnv(t *testing.T) {
	env := turnEnv([]string{"PATH=/bin", "CLAUDECODE=1", "GONE=x", "MATE_PROJECT_ID=old"},
		Launch{Env: []string{"MATE_PROJECT_ID=p"}, UnsetEnv: []string{"GONE"}})
	got := strings.Join(env, " ")
	if got != "PATH=/bin MATE_PROJECT_ID=p "+EnvHostTurn+"=1" {
		t.Fatalf("env = %q", got)
	}
}

func TestLabel(t *testing.T) {
	for in, want := range map[string]string{
		Marker + "digest: 1 item(s) — x":  "digest",
		Marker + "resolve: incident":      "resolve",
		"pr-merged: #12":                  "pr-merged",
		Marker + "a sentence: with colon": "mate",
		"plain":                           "mate",
	} {
		if got := Label(in); got != want {
			t.Errorf("Label(%q) = %q, want %q", in, got, want)
		}
	}
}

func waitBusy(t *testing.T, h *Host) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !h.State().Busy {
		if time.Now().After(deadline) {
			t.Fatal("host never got busy")
		}
		time.Sleep(10 * time.Millisecond)
	}
	// Let the fake reach its sleep, so the prompt is on disk.
	time.Sleep(300 * time.Millisecond)
}
