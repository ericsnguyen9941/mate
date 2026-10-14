package outbox_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/outbox"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/send"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// hangingRuntime is the fake Herdr with a pane read that, while hang is set,
// never answers until its context ends.
type hangingRuntime struct {
	*runtime.Fake
	hang atomic.Bool
}

func (h *hangingRuntime) ReadAgentStyled(ctx context.Context, handle runtime.AgentHandle, source harness.ReadSource, lines int) (string, error) {
	if h.hang.Load() {
		<-ctx.Done()
		return "", ctx.Err()
	}
	return h.Fake.ReadAgentStyled(ctx, handle, source, lines)
}

// TestADeliveryThatOutrunsItsBoundIsARefusal: delivery holds the outbox's
// writers' lock, so a Herdr that never answers would hold every other sender
// and enqueuer of the project with it. The delivery ends at DeliverTimeout,
// and what it leaves is what any refusal leaves: the item queued, tried
// once, the cause recorded, nothing typed - and the lock free for the next
// attempt, which delivers.
func TestADeliveryThatOutrunsItsBoundIsARefusal(t *testing.T) {
	f := newFixture(t)
	rt := &hangingRuntime{Fake: f.rt}
	rt.hang.Store(true)
	deps := f.deps()
	deps.Runtime = rt
	deps.DeliverTimeout = 50 * time.Millisecond
	ws, err := store.Open(f.ws.Root())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	s := outbox.New(ws, deps)
	f.assign(s, "crews/k3.status@0", resolveLine)

	type result struct {
		attempt attemptResult
		err     error
	}
	done := make(chan result, 1)
	begin := time.Now()
	go func() {
		a, err := s.Attempt(context.Background(), project)
		done <- result{attemptResult{a.Refusal, a.Delivered}, err}
	}()
	var got result
	select {
	case got = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Attempt did not return: a hung Herdr held the delivery and the outbox lock")
	}
	if took := time.Since(begin); took > 2*time.Second {
		t.Fatalf("Attempt took %s with a 50ms delivery bound", took)
	}
	if got.err != nil {
		t.Fatalf("Attempt: %v (a refusal is not an error)", got.err)
	}
	if got.attempt.delivered || !errors.Is(got.attempt.refusal, context.DeadlineExceeded) {
		t.Fatalf("attempt = %+v, want a refusal for the deadline", got.attempt)
	}
	items := f.items()
	if len(items) != 1 || !items[0].Queued() || items[0].Attempts != 1 || items[0].LastRefusal == "" {
		t.Fatalf("outbox = %+v, want the line queued, tried once, with its refusal", items)
	}
	if typed := f.typed(); len(typed) != 0 {
		t.Fatalf("typed %#v in a delivery that never read the pane", typed)
	}

	rt.hang.Store(false)
	f.drain(s)
	if typed := f.typed(); len(typed) != 1 || typed[0] != send.Marker+resolveLine {
		t.Fatalf("typed %#v once Herdr answered, want the line once", typed)
	}
}

type attemptResult struct {
	refusal   error
	delivered bool
}
