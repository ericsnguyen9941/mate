package chain_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/screen/chain"
)

type stateStub struct {
	reading screen.StateReading
	err     error
	calls   int
}

func (s *stateStub) Observe(context.Context, harness.ScreenProfile, string) (screen.Observation, error) {
	return screen.Observation{Composer: screen.ComposerUnknown, Source: s.reading.Source}, nil
}

func (s *stateStub) ObserveState(context.Context, harness.ScreenProfile, string) (screen.StateReading, error) {
	s.calls++
	return s.reading, s.err
}

func TestChainStateStubFirstJevWhenTheStubCannot(t *testing.T) {
	ctx := screen.WithCaller(context.Background(), screen.CallerSettle)
	jevSure := screen.StateReading{State: screen.PaneHooksReview, Confidence: 0.93, Source: "jev"}
	cases := []struct {
		name     string
		fixture  screen.StateReading
		jev      screen.StateReading
		jevErr   error
		want     screen.PaneState
		jevCalls int
		logLines int
	}{
		{"stub reads it", screen.StateReading{State: screen.PaneReady, Confidence: 1, Source: "fixture"}, jevSure, nil, screen.PaneReady, 0, 0},
		{"stub cannot, jev sure", screen.StateReading{State: screen.PaneUnknown, Source: "fixture"}, jevSure, nil, screen.PaneHooksReview, 1, 1},
		{"stub cannot, jev unsure", screen.StateReading{State: screen.PaneUnknown, Source: "fixture"}, screen.StateReading{State: screen.PaneTrustDir, Confidence: 0.6, Source: "jev"}, nil, screen.PaneUnknown, 1, 1},
		{"stub cannot, jev fails", screen.StateReading{State: screen.PaneUnknown, Source: "fixture"}, screen.StateReading{}, errors.New("boom"), screen.PaneUnknown, 1, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jev := &stateStub{reading: c.jev, err: c.jevErr}
			var lines []string
			ch := chain.New(jev, &stateStub{reading: c.fixture}, 0.85,
				chain.WithLog(func(l string) error { lines = append(lines, l); return nil }))
			r, err := ch.(screen.StateObserver).ObserveState(ctx, codexScreen, "pane text")
			if err != nil || r.State != c.want || jev.calls != c.jevCalls || len(lines) != c.logLines {
				t.Fatalf("reading %+v err %v, %d jev calls, lines %q", r, err, jev.calls, lines)
			}
			if c.logLines == 1 && !strings.Contains(lines[0], " state="+string(c.jev.State)) && c.jevErr == nil {
				t.Fatalf("log line %q does not name Jev's state", lines[0])
			}
			if c.want == screen.PaneUnknown && c.jevCalls == 1 && r.Reason == "" {
				t.Fatalf("unknown reading has no reason: %+v", r)
			}
		})
	}
}

func TestChainStateAsksJevOncePerScreenAMinute(t *testing.T) {
	now := time.Date(2026, 10, 13, 6, 0, 0, 0, time.UTC)
	jev := &stateStub{reading: screen.StateReading{State: screen.PaneHooksReview, Confidence: 0.93, Source: "jev"}}
	fix := &stateStub{reading: screen.StateReading{State: screen.PaneUnknown, Source: "fixture"}}
	ch := chain.New(jev, fix, 0.85, chain.WithClock(func() time.Time { return now })).(screen.StateObserver)
	for i := 0; i < 2; i++ {
		if _, err := ch.ObserveState(context.Background(), codexScreen, "same"); err != nil {
			t.Fatal(err)
		}
		now = now.Add(10 * time.Second)
	}
	if jev.calls != 1 {
		t.Fatalf("jev asked %d times", jev.calls)
	}
	now = now.Add(chain.TTL)
	_, _ = ch.ObserveState(context.Background(), codexScreen, "same")
	if jev.calls != 2 {
		t.Fatalf("jev asked %d times after TTL", jev.calls)
	}
}
