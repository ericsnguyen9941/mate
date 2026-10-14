package watch_test

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/watch"
)

// A Herdr that stops answering must cost a round one bounded wait per
// question, never the round itself: the observer polls its crews one after
// another, and every reading after a hung call - the other crews' health,
// the runtime notice, the timeline - waits on it.

// testHerdrTimeout is the bound the tests run with, so a hung call costs
// milliseconds; roundBound is how long a round with one hung call may take
// before the test calls it stalled.
const (
	testHerdrTimeout = 50 * time.Millisecond
	roundBound       = 2 * time.Second
)

// hangingRuntime is the fake Herdr with one call that never answers for
// the agents listed: it returns only when its context ends, as a wedged
// `herdr agent read` killed at the deadline does.
type hangingRuntime struct {
	*runtime.Fake
	mu      sync.Mutex
	inspect map[string]bool
	read    map[string]bool
}

func newHangingRuntime(fake *runtime.Fake) *hangingRuntime {
	return &hangingRuntime{Fake: fake, inspect: map[string]bool{}, read: map[string]bool{}}
}

func (h *hangingRuntime) hang(calls map[string]bool, agent string, on bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	calls[agent] = on
}

func (h *hangingRuntime) hangs(calls map[string]bool, agent string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return calls[agent]
}

func (h *hangingRuntime) InspectAgent(ctx context.Context, handle runtime.AgentHandle) (runtime.ObservedAgent, error) {
	if h.hangs(h.inspect, handle.Name) {
		<-ctx.Done()
		return runtime.ObservedAgent{}, ctx.Err()
	}
	return h.Fake.InspectAgent(ctx, handle)
}

func (h *hangingRuntime) ReadAgentStyled(ctx context.Context, handle runtime.AgentHandle, source harness.ReadSource, lines int) (string, error) {
	if h.hangs(h.read, handle.Name) {
		<-ctx.Done()
		return "", ctx.Err()
	}
	return h.Fake.ReadAgentStyled(ctx, handle, source, lines)
}

// deadlineFixture is the shop fixture with a second open crew, k4, polled
// after k3, a runtime that can hang on k3, and a counted session check.
type deadlineFixture struct {
	*fixture
	rt         *hangingRuntime
	handleHang atomic.Bool
	sessions   atomic.Int32
	sessionErr func(ctx context.Context) error
}

func newDeadlineFixture(t *testing.T) *deadlineFixture {
	t.Helper()
	f := &deadlineFixture{fixture: newFixture(t)}
	f.recordCrew("k4")
	k4 := f.handle
	k4.Name = "crew-k4"
	f.fixture.rt.PutAgent(k4, runtime.AgentIdle)
	f.fixture.rt.SetReadOutput(k4, codexIdleScreen)
	f.rt = newHangingRuntime(f.fixture.rt)

	deps := f.deps()
	deps.Runtime = f.rt
	deps.HerdrTimeout = testHerdrTimeout
	resolve := f.handleFunc()
	deps.Handle = func(ctx context.Context, project, crew string) (runtime.AgentHandle, harness.ScreenProfile, error) {
		if crew == "k3" && f.handleHang.Load() {
			<-ctx.Done()
			return runtime.AgentHandle{}, nil, ctx.Err()
		}
		return resolve(ctx, project, crew)
	}
	deps.Session = func(ctx context.Context) error {
		f.sessions.Add(1)
		if f.sessionErr != nil {
			return f.sessionErr(ctx)
		}
		return nil
	}
	f.w = watch.New(f.ws, deps)
	return f
}

// boundedPoll runs one round and fails the test if it does not come back
// within roundBound, rather than hanging the suite on a stalled round.
func (f *deadlineFixture) boundedPoll() {
	f.t.Helper()
	done := make(chan struct{})
	start := time.Now()
	go func() {
		defer close(done)
		_ = f.w.Poll(context.Background())
	}()
	select {
	case <-done:
	case <-time.After(roundBound):
		f.t.Fatalf("Poll did not return within %s: one hung Herdr call stalled the round", roundBound)
	}
	if took := time.Since(start); took > roundBound {
		f.t.Fatalf("Poll took %s, want under %s", took, roundBound)
	}
}

// TestWatchBoundsEachHerdrCallAndStillObservesTheOtherCrews: whichever of
// the three per-crew Herdr calls hangs, the round ends at the bound, the
// session is still asked, and the crew polled after the hung one gets a
// reading of this round while the hung crew keeps the one it had.
func TestWatchBoundsEachHerdrCallAndStillObservesTheOtherCrews(t *testing.T) {
	for _, call := range []string{"handle", "inspect", "read"} {
		t.Run(call, func(t *testing.T) {
			f := newDeadlineFixture(t)
			f.boundedPoll()
			before, ok := f.health("k3")
			if !ok {
				t.Fatal("no health for k3 before anything hung")
			}

			switch call {
			case "handle":
				f.handleHang.Store(true)
			case "inspect":
				f.rt.hang(f.rt.inspect, "crew-k3", true)
			case "read":
				f.rt.hang(f.rt.read, "crew-k3", true)
			}
			f.clock.advance(5 * time.Second)
			asked := f.sessions.Load()
			f.boundedPoll()

			if f.sessions.Load() == asked {
				t.Fatal("the session was not asked in the round a crew's Herdr call hung")
			}
			k4, ok := f.health("k4")
			if !ok || !k4.ObservedAt.Equal(f.clock.Now()) {
				t.Fatalf("k4 health = %+v (ok=%v), want a reading at %s", k4, ok, f.clock.Now())
			}
			k3, ok := f.health("k3")
			if !ok || !k3.ObservedAt.Equal(before.ObservedAt) {
				t.Fatalf("k3 health = %+v (ok=%v), want the reading of %s kept", k3, ok, before.ObservedAt)
			}
			f.assertIncidents()
		})
	}
}

// TestWatchStandsANoticeWhenAHerdrCallTimesOut: a crew whose reading could
// not be taken in time keeps its last Health, which on its own would look
// as current as any other row. The notice says which crew went unread, and
// clears on the round every call answers.
func TestWatchStandsANoticeWhenAHerdrCallTimesOut(t *testing.T) {
	f := newDeadlineFixture(t)
	f.boundedPoll()
	if notice, _ := f.w.RuntimeNotice(); notice != "" {
		t.Fatalf("RuntimeNotice before anything hung = %q, want none", notice)
	}

	f.rt.hang(f.rt.read, "crew-k3", true)
	f.clock.advance(5 * time.Second)
	f.boundedPoll()
	notice, at := f.w.RuntimeNotice()
	if !strings.Contains(notice, "k3") || !strings.Contains(notice, "shop") || strings.Contains(notice, "k4") {
		t.Fatalf("RuntimeNotice = %q, want it to name crew k3 of shop and only it", notice)
	}
	if !strings.Contains(notice, testHerdrTimeout.String()) {
		t.Fatalf("RuntimeNotice = %q, want it to say how long Herdr was given", notice)
	}
	if !at.Equal(f.clock.Now()) {
		t.Fatalf("RuntimeNotice at %s, want %s", at, f.clock.Now())
	}

	f.rt.hang(f.rt.read, "crew-k3", false)
	f.clock.advance(5 * time.Second)
	f.boundedPoll()
	if notice, at := f.w.RuntimeNotice(); notice != "" || !at.IsZero() {
		t.Fatalf("RuntimeNotice after Herdr answered again = %q at %s, want none", notice, at)
	}
}

// TestWatchBoundsTheSessionCheck: the session question is a Herdr call like
// the others. A hung one ends at the bound and stands a notice that says
// Herdr did not answer, not the bare context error.
func TestWatchBoundsTheSessionCheck(t *testing.T) {
	f := newDeadlineFixture(t)
	f.sessionErr = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}
	f.boundedPoll()
	notice, _ := f.w.RuntimeNotice()
	if !strings.Contains(notice, "herdr did not answer") || strings.Contains(notice, "context deadline exceeded") {
		t.Fatalf("RuntimeNotice = %q, want it to say Herdr did not answer in time", notice)
	}
	if _, ok := f.health("k4"); !ok {
		t.Fatal("no health for k4 in a round whose session check hung")
	}
}
