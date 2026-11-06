package matehost

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

// While the TUI flag is up the host runs no turn (CON-227): the captain's
// sends are refused in the captain's words, an app line and a reload as a
// busy Mate's are, and what was already on its way waits in the queue. A
// session change goes through; once the flag is down the queue runs on the
// session the TUI left.
func TestHostRunsNoTurnWhileTheTUIIsOpen(t *testing.T) {
	f := newFixture(t)
	var open atomic.Bool
	f.host.paused = open.Load
	open.Store(true)

	if _, err := f.host.Say("hi"); !errors.Is(err, ErrTUI) {
		t.Fatalf("Say while the TUI is open = %v, want ErrTUI", err)
	}
	if _, err := f.host.Answer(1, "yes"); !errors.Is(err, ErrTUI) {
		t.Fatalf("Answer while the TUI is open = %v, want ErrTUI", err)
	}
	if got := failed(ErrTUI); got.Code != "tui" || got.Error != "Mate đang mở TUI" {
		t.Fatalf("the refusal on the socket = %+v", got)
	}
	if err := f.host.Deliver(Marker + "digest: x"); !errors.Is(err, ErrTUI) {
		t.Fatalf("Deliver while the TUI is open = %v, want ErrTUI", err)
	}
	if resp := f.host.answer(request{Op: OpDeliver, Text: Marker + "digest: x"}); resp.OK || !resp.Busy {
		t.Fatalf("deliver op while the TUI is open = %+v, want a busy refusal the outbox keeps the line for", resp)
	}
	next := f.host.runner.Launch
	next.SessionID = "s3"
	if err := f.host.Reload(next); !errors.Is(err, ErrTUI) {
		t.Fatalf("Reload while the TUI is open = %v, want ErrTUI", err)
	}

	// A message held before the TUI opened comes due meanwhile: it waits.
	f.host.Submit("queued before")
	time.Sleep(1500 * time.Millisecond)
	if s := f.host.State(); s.Busy || s.Queued != 1 {
		t.Fatalf("state while the TUI is open = %+v, want one queued and no turn", s)
	}
	if _, err := os.Stat(filepath.Join(f.dir, "argv.log")); err == nil {
		t.Fatal("a turn ran while the TUI was open")
	}
	if err := f.host.SetSession("s2"); err != nil {
		t.Fatalf("SetSession while the TUI is open: %v", err)
	}

	open.Store(false)
	items := f.done(t, 1)
	if items[0].Text != "queued before" || items[0].Error != "" {
		t.Fatalf("item = %+v", items[0])
	}
	if argv := f.read(t, "argv.log"); !strings.Contains(argv, "s2") || strings.Contains(argv, "s1") {
		t.Fatalf("argv = %q, want the turn on s2, the TUI's session", argv)
	}
	if err := f.host.Deliver(Marker + "digest: y"); err != nil {
		t.Fatalf("Deliver once the TUI closed: %v", err)
	}
	f.done(t, 2)
}

// fakeTUI stands in for the harness's TUI: it writes its argv and the
// variables the test asks about, and does what $TUI_DO says - a session
// change as the Mate's hooks record it ("moved"), a long run ("sleep"),
// or a run that ends on SIGTERM ("term").
const fakeTUI = `#!/bin/sh
dir=$(dirname "$0")
echo "$*" > "$dir/tui-argv.log"
echo "CLAUDECODE=${CLAUDECODE:-} MATE_HOST_TURN=${MATE_HOST_TURN:-} MATE_PROJECT_ID=${MATE_PROJECT_ID:-}" > "$dir/tui-env.log"
echo "said in the TUI" >> "$dir/transcript-$(echo "$*" | sed 's/.*--resume //')"
case "$TUI_DO" in
moved) echo s2 > "$dir/meta-session";;
sleep) sleep 30;;
term) trap 'echo term > "$dir/tui-term.log"; exit 0' TERM; while :; do sleep 0.1; done;;
esac
`

type resumeArgs struct{}

func (resumeArgs) TUIArgs(launch []string, id string, _ bool) []string {
	return append(append([]string(nil), launch...), "--resume", id)
}

// tuiFlag is a TUIFlag over a lock file, as store's tui.lock is one.
type tuiFlag struct {
	release func()
	started atomic.Int64
	gone    atomic.Bool
}

func (f *tuiFlag) Started(pid int) error { f.started.Store(int64(pid)); return nil }
func (f *tuiFlag) Release()              { f.gone.Store(true); f.release() }

type tuiFixture struct {
	dir   string
	flag  *tuiFlag
	sigs  chan os.Signal
	notes bytes.Buffer
	run   TUIRun

	mu    sync.Mutex
	moved []string
	saved []Launch
}

func newTUIFixture(t *testing.T, do string) *tuiFixture {
	t.Helper()
	dir := t.TempDir()
	prog := filepath.Join(dir, "claude")
	if err := os.WriteFile(prog, []byte(fakeTUI), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "meta-session"), []byte("s1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TUI_DO", do)
	t.Setenv("CLAUDECODE", "1")
	t.Setenv(EnvHostTurn, "7")
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { devNull.Close() })
	x := &tuiFixture{dir: dir, flag: &tuiFlag{}, sigs: make(chan os.Signal, 1)}
	x.run = TUIRun{
		Project: "p",
		Launch:  Launch{Project: "p", Harness: "fake", Program: prog, Args: []string{"--model", "opus", "--session-id", "s1"}, Dir: dir, SessionID: "s1", Env: []string{"MATE_PROJECT_ID=p"}},
		Resume:  resumeArgs{},
		Flag: func() (TUIFlag, error) {
			release, ok, err := fileLock(filepath.Join(dir, "tui.lock"))()
			if err != nil || !ok {
				return nil, errors.New("Mate p đang mở TUI")
			}
			x.flag.release = release
			return x.flag, nil
		},
		TurnLock: fileLock(filepath.Join(dir, ".turn.lock")),
		Session: func() (string, error) {
			data, err := os.ReadFile(filepath.Join(dir, "meta-session"))
			return strings.TrimSpace(string(data)), err
		},
		SaveLaunch: func(l Launch) error { x.mu.Lock(); x.saved = append(x.saved, l); x.mu.Unlock(); return nil },
		MoveHost: func(_ context.Context, id string) error {
			x.mu.Lock()
			x.moved = append(x.moved, id)
			x.mu.Unlock()
			return nil
		},
		// The terminal, as files: a TUI killed while a child of its holds
		// the terminal still ends the run.
		Stdin: devNull, Stdout: devNull, Stderr: devNull,
		Signals: x.sigs,
		Poll:    20 * time.Millisecond,
		Grace:   2 * time.Second,
	}
	x.run.Notes = &x.notes
	return x
}

func (x *tuiFixture) start(t *testing.T) <-chan error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- RunTUI(context.Background(), x.run) }()
	return done
}

func (x *tuiFixture) wait(t *testing.T, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("RunTUI: %v", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("RunTUI did not return")
	}
}

// Both locks are free again: a host could start a turn.
func (x *tuiFixture) released(t *testing.T) {
	t.Helper()
	for _, name := range []string{"tui.lock", ".turn.lock"} {
		release, ok, err := fileLock(filepath.Join(x.dir, name))()
		if err != nil || !ok {
			t.Fatalf("%s is still held after the TUI ended (%v)", name, err)
		}
		release()
	}
}

// The TUI runs the launch on the session mate.meta records, in the
// launch's environment with no nested-session or host-turn variable; when
// the Mate's hooks recorded another session meanwhile, the launch and the
// host move to it before the locks go.
func TestRunTUIResumesTheMateSessionAndMovesTheHostAfter(t *testing.T) {
	x := newTUIFixture(t, "moved")
	x.wait(t, x.start(t))
	argv, _ := os.ReadFile(filepath.Join(x.dir, "tui-argv.log"))
	if got := strings.TrimSpace(string(argv)); got != "--model opus --session-id s1 --resume s1" {
		t.Fatalf("TUI argv = %q", got)
	}
	env, _ := os.ReadFile(filepath.Join(x.dir, "tui-env.log"))
	if got := strings.TrimSpace(string(env)); got != "CLAUDECODE= MATE_HOST_TURN= MATE_PROJECT_ID=p" {
		t.Fatalf("TUI env = %q", got)
	}
	if x.flag.started.Load() == 0 || !x.flag.gone.Load() {
		t.Fatalf("flag: started %d, released %v", x.flag.started.Load(), x.flag.gone.Load())
	}
	x.mu.Lock()
	moved, saved := x.moved, x.saved
	x.mu.Unlock()
	if len(moved) != 1 || moved[0] != "s2" || len(saved) != 1 || saved[0].SessionID != "s2" {
		t.Fatalf("moved %q, saved %+v; want the host and its launch on s2", moved, saved)
	}
	x.released(t)
}

// A TUI that ends on its session leaves the host where it was.
func TestRunTUIOnTheSameSessionMovesNothing(t *testing.T) {
	x := newTUIFixture(t, "")
	x.wait(t, x.start(t))
	if len(x.moved) != 0 || len(x.saved) != 0 {
		t.Fatalf("moved %q, saved %+v; want nothing", x.moved, x.saved)
	}
	x.released(t)
}

// A second TUI is refused while the first holds the flag.
func TestRunTUIRefusesASecondTUI(t *testing.T) {
	x := newTUIFixture(t, "")
	release, ok, err := fileLock(filepath.Join(x.dir, "tui.lock"))()
	if err != nil || !ok {
		t.Fatal(err)
	}
	defer release()
	if err := RunTUI(context.Background(), x.run); err == nil || !strings.Contains(err.Error(), "đang mở TUI") {
		t.Fatalf("second TUI = %v, want the holder's refusal", err)
	}
	if _, err := os.Stat(filepath.Join(x.dir, "tui-argv.log")); err == nil {
		t.Fatal("a refused TUI ran")
	}
}

// A turn in flight is waited for, and said; the TUI opens once it ends.
func TestRunTUIWaitsForTheTurnInFlight(t *testing.T) {
	x := newTUIFixture(t, "")
	release, ok, err := fileLock(filepath.Join(x.dir, ".turn.lock"))()
	if err != nil || !ok {
		t.Fatal(err)
	}
	done := x.start(t)
	time.Sleep(300 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(x.dir, "tui-argv.log")); err == nil {
		t.Fatal("the TUI opened while a turn held the lock")
	}
	release()
	x.wait(t, done)
	if !strings.Contains(x.notes.String(), "in the middle of a turn") {
		t.Fatalf("notes = %q, want the wait said", x.notes.String())
	}
	x.released(t)
}

// Killed with SIGKILL, the TUI still ends the run: the flag and the turn
// lock go and the host runs on.
func TestRunTUIKilledReleasesTheFlag(t *testing.T) {
	x := newTUIFixture(t, "sleep")
	done := x.start(t)
	pid := waitStarted(t, x.flag)
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	x.wait(t, done)
	if !strings.Contains(x.notes.String(), "TUI ended") {
		t.Fatalf("notes = %q", x.notes.String())
	}
	x.released(t)
}

// SIGTERM to `mate mate tui` (a force close) reaches the TUI, which ends.
func TestRunTUIPassesSIGTERMOn(t *testing.T) {
	x := newTUIFixture(t, "term")
	done := x.start(t)
	waitStarted(t, x.flag)
	time.Sleep(200 * time.Millisecond)
	x.sigs <- syscall.SIGTERM
	x.wait(t, done)
	if data, _ := os.ReadFile(filepath.Join(x.dir, "tui-term.log")); strings.TrimSpace(string(data)) != "term" {
		t.Fatalf("the TUI was not passed SIGTERM: %q", data)
	}
	x.released(t)
}

func waitStarted(t *testing.T, f *tuiFlag) int {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for f.started.Load() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("the TUI never started")
		}
		time.Sleep(10 * time.Millisecond)
	}
	return int(f.started.Load())
}

// The host and a TUI over one Mate (CON-227): while the TUI runs, the host
// queues; killed with SIGKILL, the TUI lets the flag go, and the host's
// next turn resumes the session the TUI wrote, then takes the outbox's
// line.
func TestHostRunsOnTheTUISessionAfterTheTUIIsKilled(t *testing.T) {
	f := newFixture(t)
	x := newTUIFixture(t, "sleep")
	x.run.Launch = f.host.runner.Launch
	x.run.Launch.Program = filepath.Join(x.dir, "claude")
	x.run.TurnLock = fileLock(filepath.Join(f.dir, ".turn.lock"))
	x.run.MoveHost = func(_ context.Context, id string) error { return f.host.SetSession(id) }
	f.opened["s1"] = true
	f.host.paused = func() bool {
		release, ok, err := fileLock(filepath.Join(x.dir, "tui.lock"))()
		if err != nil || ok {
			if ok {
				release()
			}
			return false
		}
		return true
	}
	done := x.start(t)
	pid := waitStarted(t, x.flag)
	f.host.Submit("while the TUI runs")
	time.Sleep(1200 * time.Millisecond)
	if s := f.host.State(); s.Busy || s.Queued != 1 {
		t.Fatalf("state with the TUI open = %+v, want the message queued", s)
	}
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	x.wait(t, done)
	items := f.done(t, 1)
	if items[0].Text != "while the TUI runs" {
		t.Fatalf("item = %+v", items[0])
	}
	transcript, err := os.ReadFile(filepath.Join(x.dir, "transcript-s1"))
	if err != nil || !strings.Contains(string(transcript), "said in the TUI") {
		t.Fatalf("the TUI wrote no turn into s1: %q, %v", transcript, err)
	}
	if argv := f.read(t, "argv.log"); !strings.Contains(argv, "--resume s1") {
		t.Fatalf("the host's next turn = %q, want it to resume s1, the session the TUI wrote", argv)
	}
	if err := f.host.Deliver(Marker + "digest: z"); err != nil {
		t.Fatalf("the outbox's line after the TUI: %v", err)
	}
	f.done(t, 2)
}
