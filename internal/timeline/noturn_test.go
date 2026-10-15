package timeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/box"
)

func noTurnIncidents(t *testing.T, f *fixture) (open, resolved int) {
	t.Helper()
	view, err := box.Load(f.ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	for _, e := range view.Entries {
		if e.Kind != box.KindIncident || e.Crew != fixtureCrew || e.Incident == nil || e.Incident.Kind != box.IncidentNoTurn {
			continue
		}
		if e.Incident.Resolved {
			resolved++
		} else {
			open++
		}
	}
	return open, resolved
}

// A working crew launched longer than the threshold ago with no turn in
// the ledger gets one `no_turn` incident, however many polls follow
// (CON-74); once the ledger records its turns the incident is resolved.
// The fixture crew was launched 2m51s before the fixture clock.
func TestCheckNoTurnFilesOnceAndResolvesOnTheFirstTurn(t *testing.T) {
	f := newFixture(t)
	crewOnlyWorking(t, f)
	for range 2 {
		if err := f.ing.CheckNoTurn(context.Background(), 2*time.Minute); err != nil {
			t.Fatalf("CheckNoTurn: %v", err)
		}
	}
	if open, resolved := noTurnIncidents(t, f); open != 1 || resolved != 0 {
		t.Fatalf("no_turn incidents open=%d resolved=%d, want exactly one open", open, resolved)
	}
	view, _ := box.Load(f.ws, fixtureProject)
	if inc := box.OpenIncidents(view, fixtureCrew); len(inc) != 1 || inc[0].Text != "no turn in the ledger 3m0s after launch" {
		t.Fatalf("incident = %+v, want the no-turn text", inc)
	}
	if len(box.BlockingIncidents(view, fixtureCrew)) != 0 {
		t.Fatal("a no_turn incident made the crew blocked")
	}

	f.ingest(t)
	for range 2 {
		if err := f.ing.CheckNoTurn(context.Background(), 2*time.Minute); err != nil {
			t.Fatalf("CheckNoTurn after ingest: %v", err)
		}
	}
	if open, resolved := noTurnIncidents(t, f); open != 0 || resolved != 1 {
		t.Fatalf("no_turn incidents open=%d resolved=%d after the ledger has turns, want the one resolved", open, resolved)
	}
}

// Under the threshold, or once the crew has handed back, nothing is filed.
func TestCheckNoTurnFilesNothingUnderTheThresholdOrAfterAHandBack(t *testing.T) {
	f := newFixture(t)
	crewOnlyWorking(t, f)
	if err := f.ing.CheckNoTurn(context.Background(), 5*time.Minute); err != nil {
		t.Fatalf("CheckNoTurn: %v", err)
	}
	if open, _ := noTurnIncidents(t, f); open != 0 {
		t.Fatalf("filed under the threshold: %d open", open)
	}
	g := newFixture(t) // its own status file ends in wait-mate
	if err := g.ing.CheckNoTurn(context.Background(), 2*time.Minute); err != nil {
		t.Fatalf("CheckNoTurn: %v", err)
	}
	if open, _ := noTurnIncidents(t, g); open != 0 {
		t.Fatalf("filed after a hand-back: %d open", open)
	}
}
