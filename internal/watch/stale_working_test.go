package watch_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/watch"
)

// stillFor advances the clock by d in polls a minute apart, the pane and
// the status file unchanged.
func (f *fixture) stillFor(d time.Duration) {
	f.t.Helper()
	for step := time.Duration(0); step < d; step += time.Minute {
		f.clock.advance(time.Minute)
		f.poll()
		f.w.AwaitObserver()
	}
}

// A crew that ended its turn waiting on a background task reads as working,
// which keeps stale shut - until the pane has held still past the working
// ceiling: then a task that hung opens stale (CON-281, mih-21-engine-ios
// waited ten hours on an export that never finished). The incident stays
// open while the pane holds, and the pane moving resolves it.
func TestWatchOpensStaleOnAWaitingPaneStillPastTheCeiling(t *testing.T) {
	f := newFixture(t)
	says := &statesSay{
		observerSays: observerSays{obs: screen.Observation{Composer: screen.ComposerEmpty, Deterministic: screen.ComposerEmpty, Source: "fixture"}},
		state:        screen.PaneBackground,
	}
	deps := f.deps()
	deps.Observer = says
	f.w = watch.New(f.ws, deps)
	f.appendStatus("k3", "working: exporting the golden images")

	f.poll()
	f.stillFor(watch.DefaultStaleWorkingAfter - time.Minute)
	f.assertIncidents()

	f.stillFor(2 * time.Minute)
	f.assertIncidents("k3 stale open")
	text := f.incidents()[0].Text
	for _, want := range []string{"waiting on a background task", "may be hung", "last status verb working"} {
		if !strings.Contains(text, want) {
			t.Errorf("incident text = %q, want %q in it", text, want)
		}
	}

	f.stillFor(10 * time.Minute)
	f.assertIncidents("k3 stale open")

	f.setScreen(codexIdleScreen + "\nexport finished\n")
	f.clock.advance(5 * time.Second)
	f.poll()
	f.assertIncidents("k3 stale open", "k3 stale resolved")
}

// A busy composer that never redraws - not even its timer - is hung as
// surely: past the ceiling it opens stale, and its own busy reading does
// not resolve the incident it opened.
func TestWatchOpensStaleOnAFrozenBusyPanePastTheCeiling(t *testing.T) {
	f := newFixture(t)
	f.setScreen(codexBusyScreen)
	f.appendStatus("k3", "working: running the suite")

	f.poll()
	f.stillFor(watch.DefaultStaleWorkingAfter - time.Minute)
	f.assertIncidents()

	f.stillFor(5 * time.Minute)
	f.assertIncidents("k3 stale open")
	if text := f.incidents()[0].Text; !strings.Contains(text, "reads as busy") {
		t.Fatalf("incident text = %q, want the busy reading named", text)
	}
}

// project.yaml's stale_working_after sets the ceiling per project; a value
// that does not parse keeps the default rather than lifting the ceiling.
func TestWatchTakesTheWorkingCeilingFromProjectYAML(t *testing.T) {
	for _, c := range []struct {
		value   string
		opensBy time.Duration
	}{
		{"10m", 11 * time.Minute},
		{"soon", watch.DefaultStaleWorkingAfter + time.Minute},
	} {
		t.Run(c.value, func(t *testing.T) {
			f := newFixture(t)
			cfg, err := f.ws.LoadProject("shop")
			if err != nil {
				t.Fatal(err)
			}
			cfg.StaleWorkingAfter = c.value
			if err := f.ws.SaveProject("shop", cfg); err != nil {
				t.Fatal(err)
			}
			f.setScreen(codexBusyScreen)
			f.poll()
			f.stillFor(c.opensBy - 2*time.Minute)
			f.assertIncidents()
			f.stillFor(2 * time.Minute)
			f.assertIncidents("k3 stale open")
		})
	}
}
