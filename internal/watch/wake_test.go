package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/watch"
)

// roundIngest tells each round's end: the timeline is recorded last in a
// round.
type roundIngest struct{ rounds chan time.Time }

func (r roundIngest) Ingest(context.Context) error {
	r.rounds <- time.Now()
	return nil
}

func awaitRound(t *testing.T, rounds <-chan time.Time, within time.Duration) time.Time {
	t.Helper()
	select {
	case at := <-rounds:
		return at
	case <-time.After(within):
		t.Fatalf("no round within %s", within)
		return time.Time{}
	}
}

// wakeGapFloor is WakeGap less the time a round takes: the gap is timed
// from a round's start, the test sees rounds end.
const wakeGapFloor = watch.WakeGap - 200*time.Millisecond

// A Wake (the core's change feed, CON-260) starts the next round long
// before the poll interval ends, though no sooner than WakeGap after the
// last round started; a burst of wakes is one early round.
func TestWakeStartsTheNextRoundEarly(t *testing.T) {
	f := newFixture(t)
	rounds := make(chan time.Time, 16)
	deps := f.deps()
	deps.Clock = nil
	deps.PollInterval = time.Hour
	deps.Timeline = roundIngest{rounds}
	w := watch.New(f.ws, deps)
	w.Start(context.Background())
	t.Cleanup(w.Stop)

	first := awaitRound(t, rounds, 2*time.Second)
	for range 10 {
		w.Wake()
	}
	second := awaitRound(t, rounds, 3*time.Second)
	if gap := second.Sub(first); gap < wakeGapFloor {
		t.Fatalf("the woken round came %s after the last, want WakeGap at least", gap)
	}
	select {
	case <-rounds:
		t.Fatal("a burst of wakes ran more than one early round")
	case <-time.After(1500 * time.Millisecond):
	}

	// One more wake, one more round.
	w.Wake()
	awaitRound(t, rounds, 3*time.Second)
}
