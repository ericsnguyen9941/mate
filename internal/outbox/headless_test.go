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

// headlessMate is a fixture whose Mate runs headless (CON-78): its
// workspace is under /tmp so the host's socket path fits, and a host
// answers on it. Its harness is a script that reads the file `mode` in the
// Mate's directory before each turn: `fail` ends the turn in an error,
// `hang` runs until the file `release` appears, anything else finishes.
type headlessMate struct {
	*fixture
	prompts, mateDir string
	runner           *matehost.Runner
	stop             context.CancelFunc
}

func newHeadlessMate(t *testing.T, lockFree bool) *headlessMate {
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
	mode, release := filepath.Join(mateDir, "mode"), filepath.Join(mateDir, "release")
	script := "#!/bin/sh\ncat >> " + prompts + "\necho >> " + prompts + "\n" +
		"case \"$(cat " + mode + " 2>/dev/null)\" in\n" +
		"fail) echo '{\"type\":\"result\",\"is_error\":true,\"result\":\"overloaded\",\"session_id\":\"s\"}' ;;\n" +
		"hang) while [ ! -f " + release + " ]; do sleep 0.05; done; echo '{\"type\":\"result\",\"result\":\"—\",\"session_id\":\"s\"}' ;;\n" +
		"*) echo '{\"type\":\"result\",\"result\":\"—\",\"session_id\":\"s\"}' ;;\n" +
		"esac\n"
	if err := os.WriteFile(prog, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	// A turn left hanging by a killed host ends with the test.
	t.Cleanup(func() { _ = os.WriteFile(release, nil, 0o644) })
	m := &headlessMate{fixture: f, prompts: prompts, mateDir: mateDir, runner: &matehost.Runner{
		Launch: matehost.Launch{Project: project, Harness: "claude", Program: prog, Dir: mateDir, SessionID: "s"},
		Turns:  claude.Claude{},
		Lock:   func() (func(), bool, error) { return func() {}, lockFree, nil },
		Opened: func(string, string) bool { return false },
	}}
	m.startHost()
	return m
}

// startHost serves a new host on the Mate's socket, keeping host.jsonl.
func (m *headlessMate) startHost() {
	m.t.Helper()
	host := matehost.New(m.runner, matehost.Options{History: matehost.HistoryPath(m.mateDir)})
	ctx, cancel := context.WithCancel(context.Background())
	m.stop = cancel
	m.t.Cleanup(cancel)
	go func() { _ = host.Serve(ctx, matehost.SocketPath(m.mateDir), nil) }()
	deadline := time.Now().Add(5 * time.Second)
	for matehost.Dial(m.mateDir).Ping(ctx) != nil {
		if time.Now().After(deadline) {
			m.t.Fatal("host never answered")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// killHost is the host process dying: its socket goes, and the turn in
// flight never reaches host.jsonl.
func (m *headlessMate) killHost() {
	m.t.Helper()
	m.stop()
	deadline := time.Now().Add(5 * time.Second)
	for matehost.Dial(m.mateDir).Ping(context.Background()) == nil {
		if time.Now().After(deadline) {
			m.t.Fatal("host still answers")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (m *headlessMate) mode(mode string) {
	m.t.Helper()
	if err := os.WriteFile(filepath.Join(m.mateDir, "mode"), []byte(mode), 0o644); err != nil {
		m.t.Fatal(err)
	}
}

func headlessFixture(t *testing.T, lockFree bool) (*fixture, string) {
	m := newHeadlessMate(t, lockFree)
	return m.fixture, m.prompts
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

const digestKey = "digest:0123456789abcdef"

// offerDigest is the auto daemon's tick over one status line at offset 42:
// a digest of it while the cursor is short of it, tried at once. It
// reports whether a digest was queued.
func (m *headlessMate) offerDigest(s *outbox.Sender) bool {
	m.t.Helper()
	_, queued, err := s.Offer(project, store.OutboxSourceDigest, func() (outbox.Request, bool, error) {
		cursor, err := m.ws.ReadAutoCursor(project)
		if err != nil || cursor[m.statusFile()] >= 42 {
			return outbox.Request{}, false, err
		}
		return outbox.Request{Key: digestKey, Text: digestLine, Cursor: map[string]int64{m.statusFile(): 42}}, true, nil
	})
	if err != nil {
		m.t.Fatalf("Offer: %v", err)
	}
	if queued {
		if _, err := s.Attempt(context.Background(), project); err != nil {
			m.t.Fatalf("Attempt: %v", err)
		}
	}
	return queued
}

const digestLine = `digest: 1 item(s) — k3 wait-mate: "done" — status files under /crews; act per AGENTS.md section 10`

// finished drains until the newest item is no longer delivered.
func (m *headlessMate) finished(s *outbox.Sender) store.OutboxItem {
	m.t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		m.drain(s)
		items := m.items()
		if last := items[len(items)-1]; last.State != store.OutboxDelivered {
			return last
		}
		if time.Now().After(deadline) {
			m.t.Fatalf("outbox = %+v, the turn never ended", items)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (m *headlessMate) statusFile() string {
	return filepath.Join(m.ws.CrewsDir(project), "k3.status")
}

func (m *headlessMate) cursor() map[string]int64 {
	m.t.Helper()
	cursor, err := m.ws.ReadAutoCursor(project)
	if err != nil {
		m.t.Fatal(err)
	}
	return cursor
}

func (m *headlessMate) turns() int {
	data, _ := os.ReadFile(m.prompts)
	return strings.Count(string(data), send.Marker+digestLine)
}

func TestAHeadlessDigestIsSentWhenItsTurnEnds(t *testing.T) {
	m := newHeadlessMate(t, true)
	m.mode("hang")
	s := m.sender()
	if !m.offerDigest(s) {
		t.Fatal("nothing was queued")
	}
	items := m.items()
	if len(items) != 1 || items[0].State != store.OutboxDelivered || items[0].TurnAt.IsZero() {
		t.Fatalf("outbox = %+v, want the digest delivered with its turn", items)
	}
	if cursor := m.cursor(); len(cursor) != 0 {
		t.Fatalf("cursor = %v while the turn runs", cursor)
	}
	// Another tick while the Mate is on it gathers nothing.
	m.drain(s)
	if m.offerDigest(s) {
		t.Fatalf("a second digest was queued while the first one's turn runs: %+v", m.items())
	}

	if err := os.WriteFile(filepath.Join(m.mateDir, "release"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	item := m.finished(s)
	if item.State != store.OutboxSent {
		t.Fatalf("item = %+v, want sent", item)
	}
	if cursor := m.cursor(); cursor[m.statusFile()] != 42 {
		t.Fatalf("cursor = %v, want the digest's", cursor)
	}
	// CON-194: a turn that ended is not answered with another.
	m.clock.Advance(outbox.DefaultTurnRetryAfter)
	m.offerDigest(s)
	if n := len(m.items()); n != 1 || m.turns() != 1 {
		t.Fatalf("outbox = %+v, %d turns, want the one digest once", m.items(), m.turns())
	}
}

func TestADigestWhoseHostDiedIsOfferedAgain(t *testing.T) {
	m := newHeadlessMate(t, true)
	m.mode("hang")
	s := m.sender()
	m.offerDigest(s)
	m.killHost()
	m.mode("")
	m.startHost()

	item := m.finished(s)
	if item.State != store.OutboxFailed || !strings.Contains(item.LastRefusal, "restarted") {
		t.Fatalf("item = %+v, want failed by the restart", item)
	}
	if cursor := m.cursor(); len(cursor) != 0 {
		t.Fatalf("cursor = %v after a turn that never ended", cursor)
	}
	// Not at once: a minute after the failure.
	if m.offerDigest(s) {
		t.Fatalf("queued again at once: %+v", m.items())
	}
	m.clock.Advance(outbox.DefaultTurnRetryAfter)
	if !m.offerDigest(s) {
		t.Fatalf("not queued again after %s: %+v", outbox.DefaultTurnRetryAfter, m.items())
	}
	if item := m.finished(s); item.State != store.OutboxSent || item.Key != digestKey {
		t.Fatalf("item = %+v, want the same items sent", item)
	}
	if cursor := m.cursor(); cursor[m.statusFile()] != 42 {
		t.Fatalf("cursor = %v", cursor)
	}
	if n := m.turns(); n != 2 {
		t.Fatalf("the Mate was given the digest %d times, want 2", n)
	}
}

func TestADigestTurnThatKeepsFailingStopsAndIsWedged(t *testing.T) {
	m := newHeadlessMate(t, true)
	m.mode("fail")
	s := m.sender()
	for i := 0; i < outbox.MaxTurnFailures; i++ {
		if !m.offerDigest(s) {
			t.Fatalf("try %d not queued: %+v", i+1, m.items())
		}
		if item := m.finished(s); item.State != store.OutboxFailed || !strings.Contains(item.LastRefusal, "overloaded") {
			t.Fatalf("try %d: item = %+v, want failed with the turn's error", i+1, item)
		}
		m.clock.Advance(outbox.DefaultTurnRetryAfter)
	}
	if m.offerDigest(s) {
		t.Fatalf("queued a %dth time: %+v", outbox.MaxTurnFailures+1, m.items())
	}
	incidents := m.incidents()
	if len(incidents) != 1 || incidents[0].Kind != outbox.WedgedKind || !incidents[0].Open() {
		t.Fatalf("incidents = %+v, want one open wedged", incidents)
	}
	if cursor := m.cursor(); len(cursor) != 0 {
		t.Fatalf("cursor = %v", cursor)
	}
}
