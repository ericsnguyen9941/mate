package outbox_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/outbox"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// The loop wakes on events rather than every two seconds (CON-270): with
// nothing queued it waits for a wake - an `outbox` notice (Wake), an
// Enqueue or Offer in this process - or its one-minute safety net; with a
// line queued it retries on the old two-second cadence.

// stepSleeper is the loop's pause, one at a time: each Sleep says what it
// was asked on asked, then waits for a release, or for its ctx to end - a
// wake, or Stop. The pauses inside a delivery pass at once.
type stepSleeper struct {
	asked   chan time.Duration
	release chan struct{}
}

func newStepSleeper() *stepSleeper {
	return &stepSleeper{asked: make(chan time.Duration, 64), release: make(chan struct{})}
}

func (s *stepSleeper) Sleep(ctx context.Context, d time.Duration) error {
	if d != outbox.DefaultInterval && d != outbox.DefaultIdleInterval {
		// A delivery's own settle and retries (send.Send), at memory speed.
		return ctx.Err()
	}
	s.asked <- d
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.release:
		return nil
	}
}

// pause is the next pause the loop takes, which follows the pass before it.
func (s *stepSleeper) pause(t *testing.T) time.Duration {
	t.Helper()
	select {
	case d := <-s.asked:
		return d
	case <-time.After(5 * time.Second):
		t.Fatal("the loop ran no pass")
		return 0
	}
}

// still checks the loop is not woken: no pass and no pause for a while.
func (s *stepSleeper) still(t *testing.T) {
	t.Helper()
	select {
	case d := <-s.asked:
		t.Fatalf("the loop ran a pass (then paused %s) with nothing to wake it", d)
	case <-time.After(100 * time.Millisecond):
	}
}

// loop starts a sender loop over its own handle whose pauses sleeper takes,
// until the test ends.
func (f *fixture) loop(sleeper outbox.Sleeper) *outbox.Sender {
	f.t.Helper()
	ws, err := store.Open(f.ws.Root())
	if err != nil {
		f.t.Fatal(err)
	}
	d := f.deps()
	d.Sleeper = sleeper
	s := outbox.New(ws, d)
	s.Start(context.Background())
	f.t.Cleanup(s.Stop)
	return s
}

// Nothing queued: the first pass compacts a quiet outbox (CON-258), then
// the loop waits a minute, and a Nudge - a composer that may have emptied -
// does not wake it, for there is nothing to deliver. An Enqueue by another
// Sender of this process over the same workspace wakes it at once, and it
// delivers.
func TestAQuietLoopSleepsUntilAnEnqueueInThisProcess(t *testing.T) {
	f := newFixture(t)
	path := filepath.Join(f.ws.MateDir(project), ".outbox")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	at := start.Format(time.RFC3339Nano)
	old := `{"id":1,"at":"` + at + `","source":"digest","key":"d1","text":"x","state":"sent","sent_at":"` + at + `","cursor":{"/w/a.status":7}}` + "\n"
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	sleeper := newStepSleeper()
	s := f.loop(sleeper)
	if d := sleeper.pause(t); d != outbox.DefaultIdleInterval {
		t.Fatalf("pause with nothing queued = %s, want %s", d, outbox.DefaultIdleInterval)
	}
	if data, _ := os.ReadFile(path); strings.Contains(string(data), "cursor") {
		t.Fatalf("outbox after the first pass = %s, want it compacted", data)
	}
	s.Nudge()
	sleeper.still(t)

	f.assign(f.sender(), "crews/k3.status@0", resolveLine)
	if d := sleeper.pause(t); d != outbox.DefaultInterval {
		t.Fatalf("pause after the pass that delivered = %s, want %s: the queue is looked at once more", d, outbox.DefaultInterval)
	}
	if typed := f.typed(); len(typed) != 1 {
		t.Fatalf("typed %#v, want the line delivered on the wake", typed)
	}
	sleeper.release <- struct{}{}
	if d := sleeper.pause(t); d != outbox.DefaultIdleInterval {
		t.Fatalf("pause once the queue is empty = %s, want %s", d, outbox.DefaultIdleInterval)
	}
	sleeper.still(t)
}

// An `outbox` notice - a line queued by another mate process - is Wake.
func TestWakeRunsAPassAtOnce(t *testing.T) {
	f := newFixture(t)
	sleeper := newStepSleeper()
	s := f.loop(sleeper)
	sleeper.pause(t)
	// Queued without waking the loop, as another process's write is.
	if err := f.ws.UpdateOutbox(project, start, func(items []store.OutboxItem) ([]store.OutboxItem, bool, error) {
		return append(items, store.OutboxItem{ID: 1, At: start, Source: store.OutboxSourceAssign, Key: "k", Text: resolveLine, State: store.OutboxQueued}), true, nil
	}); err != nil {
		t.Fatal(err)
	}
	sleeper.still(t)
	s.Wake()
	sleeper.pause(t)
	if typed := f.typed(); len(typed) != 1 {
		t.Fatalf("typed %#v, want the line delivered on the wake", typed)
	}
}

// A line the Mate is too busy for is tried every two seconds, as before,
// and a Nudge then tries it at once.
func TestAQueuedLineRetriesOnTheOldCadence(t *testing.T) {
	f := newFixture(t)
	f.rt.SetReadOutput(f.handle, claudeBusyScreen())
	f.assign(f.sender(), "crews/k3.status@0", resolveLine)
	sleeper := newStepSleeper()
	s := f.loop(sleeper)
	for i := range 3 {
		if d := sleeper.pause(t); d != outbox.DefaultInterval {
			t.Fatalf("pause %d with a line queued = %s, want %s", i, d, outbox.DefaultInterval)
		}
		sleeper.release <- struct{}{}
	}
	sleeper.pause(t)
	f.rt.SetReadOutput(f.handle, claudeScreen(""))
	s.Nudge()
	sleeper.pause(t)
	if typed := f.typed(); len(typed) != 1 {
		t.Fatalf("typed %#v, want the line delivered on the nudge", typed)
	}
	if items := f.items(); items[0].Attempts != 5 {
		t.Fatalf("attempts = %d, want 4 refused and the one that went", items[0].Attempts)
	}
}

// clockSleeper runs the loop's pauses on a virtual clock, at memory speed,
// until it reaches until; it then holds the loop until Stop.
type clockSleeper struct {
	now, until time.Duration
	passes     int
	done       chan struct{}
}

func (s *clockSleeper) Sleep(ctx context.Context, d time.Duration) error {
	if d != outbox.DefaultInterval && d != outbox.DefaultIdleInterval {
		return ctx.Err()
	}
	s.passes++
	if s.now >= s.until {
		close(s.done)
		<-ctx.Done()
		return ctx.Err()
	}
	s.now += d
	return nil
}

// Passes in ten minutes of a fake clock, the figures in the PR (CON-270):
// a quiet workspace runs one pass a minute, where it ran thirty; one with a
// line waiting still runs thirty.
func TestPassesPerMinute(t *testing.T) {
	for _, tc := range []struct {
		name   string
		queued bool
		want   int
	}{
		{"nothing queued", false, 10},
		{"a line waiting", true, 300},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			if tc.queued {
				f.rt.SetReadOutput(f.handle, claudeBusyScreen())
				f.assign(f.sender(), "crews/k3.status@0", resolveLine)
			}
			sleeper := &clockSleeper{until: 10 * time.Minute, done: make(chan struct{})}
			f.loop(sleeper)
			select {
			case <-sleeper.done:
			case <-time.After(20 * time.Second):
				t.Fatal("the loop did not reach ten minutes")
			}
			// The pass at the start, then one per pause.
			if got := sleeper.passes - 1; got != tc.want {
				t.Fatalf("passes in ten minutes = %d, want %d", got, tc.want)
			}
			t.Logf("%s: %d passes in 10 minutes (%.1f a minute)", tc.name, sleeper.passes-1, float64(sleeper.passes-1)/10)
		})
	}
}

// A loop over the engine's handle (SharedConfig) does not read
// workspace.yaml: a project added to the file by another handle is not seen
// until the owner reloads it.
func TestASharedConfigLoopLeavesTheFileToItsOwner(t *testing.T) {
	f := newFixture(t)
	ws, err := store.Open(f.ws.Root())
	if err != nil {
		t.Fatal(err)
	}
	d := f.deps()
	d.SharedConfig = true
	s := outbox.New(ws, d)
	other := "lab"
	if err := f.ws.AddProject(other, store.ProjectConfig{}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Enqueue(other, outbox.Request{Source: store.OutboxSourceAssign, Key: "k", Text: "hello"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Drain(context.Background()); err != nil {
		t.Fatal(err)
	}
	if typed := f.typed(); len(typed) != 0 {
		t.Fatalf("typed %#v: the pass read workspace.yaml itself", typed)
	}
	if err := ws.LoadConfig(); err != nil {
		t.Fatal(err)
	}
	if err := s.Drain(context.Background()); err != nil {
		t.Fatal(err)
	}
	if typed := f.typed(); len(typed) != 1 {
		t.Fatalf("typed %#v after the owner reloaded, want the line", typed)
	}
}
