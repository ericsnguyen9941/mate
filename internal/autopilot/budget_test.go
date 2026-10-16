package autopilot_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/autopilot"
	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// budgetSpy records what the daemon's budget enforcement (CON-94) asked
// for: lines typed into a crew's pane and crews stopped.
type budgetSpy struct {
	lines   []string // "<crew>: <text>"
	stops   []string
	lineErr error
}

func (s *budgetSpy) actions() *autopilot.BudgetActions {
	return &autopilot.BudgetActions{
		Line: func(_ context.Context, _, crew, text string) error {
			if s.lineErr != nil {
				return s.lineErr
			}
			s.lines = append(s.lines, crew+": "+text)
			return nil
		},
		Stop: func(_ context.Context, _, crew string) (string, error) {
			s.stops = append(s.stops, crew)
			return "agent stopped, unlanded work kept", nil
		},
	}
}

// budgetFixture is the daemon fixture with a spy for the budget actions and
// one working crew k1 with a live agent and a 2M budget.
func budgetFixture(t *testing.T) (*fixture, *budgetSpy) {
	t.Helper()
	f := newFixture(t)
	spy := &budgetSpy{}
	deps := f.deps(nil)
	deps.Budget = spy.actions()
	f.pilot = autopilot.New(f.ws, deps)
	f.crewMeta("k1", map[string]string{"state": "spawned", "agent": "crew-k1", "pane": "w1:p2", store.MetaBudget: "2000000"})
	f.status("k1", "working: building the endpoint")
	return f, spy
}

func (f *fixture) meta(crew string) map[string]string {
	f.t.Helper()
	meta, err := f.ws.ReadCrewMeta(project, crew)
	if err != nil {
		f.t.Fatalf("ReadCrewMeta: %v", err)
	}
	return meta
}

// Under the budget - nothing filed, or only the 50% incident - the daemon
// types nothing into the crew and stops nothing.
func TestBudgetUnderTheCeilingSendsNothing(t *testing.T) {
	f, spy := budgetFixture(t)
	f.mustTick()
	f.incident("k1", box.IncidentBudget, store.IncidentOpen, box.BudgetHalfPrefix+"1.1M tokens of 2M tokens, no hand-back yet")
	f.mustTick()
	if len(spy.lines) != 0 || len(spy.stops) != 0 {
		t.Fatalf("lines %q, stops %q; want nothing under the budget", spy.lines, spy.stops)
	}
}

// At 100% the crew gets exactly one "budget reached" line, however many
// ticks follow, in manual mode as in auto, and the incident stays open for
// the Mate.
func TestBudgetReachedTypesOneLineIntoTheCrew(t *testing.T) {
	f, spy := budgetFixture(t)
	f.incident("k1", box.IncidentBudget, store.IncidentOpen, "2.1M tokens of 2M tokens")
	for range 3 {
		f.mustTick()
	}
	want := "k1: " + autopilot.BudgetReachedLine
	if len(spy.lines) != 1 || spy.lines[0] != want {
		t.Fatalf("lines = %q, want exactly %q", spy.lines, want)
	}
	if !strings.Contains(autopilot.BudgetReachedLine, "stop with wait-mate") {
		t.Fatalf("the line does not ask for a clean stop: %q", autopilot.BudgetReachedLine)
	}
	if f.meta("k1")[store.MetaBudgetLine] == "" {
		t.Fatalf("budget_line= not recorded: %v", f.meta("k1"))
	}
	if len(spy.stops) != 0 {
		t.Fatalf("stops = %q at 100%%, want none", spy.stops)
	}
	open := box.OpenIncidents(f.view(), "k1")
	if len(open) != 1 || open[0].Kind != box.IncidentBudget {
		t.Fatalf("open incidents = %+v, want the budget incident still open", open)
	}
}

// A refused send is not recorded, so the next tick tries again.
func TestBudgetReachedRetriesARefusedLine(t *testing.T) {
	f, spy := budgetFixture(t)
	f.incident("k1", box.IncidentBudget, store.IncidentOpen, "2.1M tokens of 2M tokens")
	spy.lineErr = errors.New("composer has a draft")
	if err := f.tick(); err == nil {
		t.Fatal("Tick reported no error for a refused line")
	}
	if f.meta("k1")[store.MetaBudgetLine] != "" {
		t.Fatal("a refused line was recorded as sent")
	}
	spy.lineErr = nil
	f.mustTick()
	if len(spy.lines) != 1 {
		t.Fatalf("lines = %q, want the retry to go through once", spy.lines)
	}
}

// A crew that already handed back has stopped at a clean point: no line.
func TestBudgetReachedSkipsACrewThatHandedBack(t *testing.T) {
	f, spy := budgetFixture(t)
	f.status("k1", "wait-mate: done, see handback.md")
	f.incident("k1", box.IncidentBudget, store.IncidentOpen, "2.1M tokens of 2M tokens")
	f.mustTick()
	if len(spy.lines) != 0 {
		t.Fatalf("lines = %q, want none for a resting crew", spy.lines)
	}
}

// At the hard ceiling the crew is stopped once, without asking anyone, and
// the incident - "budget hard ceiling" - reaches the Mate in the digest.
func TestBudgetHardCeilingStopsTheCrewAndTellsTheMate(t *testing.T) {
	f, spy := budgetFixture(t)
	f.auto(true)
	f.incident("k1", box.IncidentBudget, store.IncidentOpen, "2.1M tokens of 2M tokens")
	f.incident("k1", box.IncidentBudget, store.IncidentOpen, box.BudgetHardPrefix+"4.1M tokens of 2M tokens, over twice the budget; the app stops the crew")
	f.mustTick()
	f.mustTick()
	if len(spy.stops) != 1 || spy.stops[0] != "k1" {
		t.Fatalf("stops = %q, want k1 stopped exactly once", spy.stops)
	}
	if len(spy.lines) != 0 {
		t.Fatalf("lines = %q, want no budget line to a crew being stopped", spy.lines)
	}
	if got := f.meta("k1")[store.MetaBudgetStop]; !strings.Contains(got, "agent stopped") {
		t.Fatalf("budget_stop = %q, want what the stop did", got)
	}
	digest := f.requireOneDigest()
	if !strings.Contains(digest, "k1 budget hard ceiling: 4.1M tokens of 2M tokens") {
		t.Fatalf("digest does not carry the hard ceiling: %q", digest)
	}
}

// An extension resolves the open budget incidents and clears budget_line
// (cmd/mate's `crew budget --extend`); the observer files the next one
// against the new ceiling, and the daemon acts on that one again.
func TestBudgetExtensionReArmsTheLineAtTheNewCeiling(t *testing.T) {
	f, spy := budgetFixture(t)
	f.incident("k1", box.IncidentBudget, store.IncidentOpen, "2.1M tokens of 2M tokens")
	f.mustTick()

	// What `mate crew budget --extend 3000000` writes.
	meta := f.meta("k1")
	meta[store.MetaBudget] = "3000000"
	meta[store.MetaBudgetExtensions] = "1"
	delete(meta, store.MetaBudgetLine)
	f.crewMeta("k1", meta)
	f.incident("k1", box.IncidentBudget, store.IncidentResolved, "budget extended to 3M tokens: tests left")
	f.mustTick()
	if len(spy.lines) != 1 {
		t.Fatalf("lines = %q, want nothing more until the new ceiling is crossed", spy.lines)
	}

	f.incident("k1", box.IncidentBudget, store.IncidentOpen, "3.2M tokens of 3M tokens")
	f.mustTick()
	if len(spy.lines) != 2 {
		t.Fatalf("lines = %q, want a second line at the new ceiling", spy.lines)
	}
}
