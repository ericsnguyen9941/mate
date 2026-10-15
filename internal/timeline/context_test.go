package timeline_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// setContextHigh writes crew.context_high into the fixture's workspace.yaml;
// CheckBudgets reloads the config at the start of every run.
func setContextHigh(t *testing.T, ws *store.Workspace, n int) {
	t.Helper()
	data, err := os.ReadFile(ws.WorkspaceFile())
	if err != nil {
		t.Fatal(err)
	}
	// A second call replaces the block the first one appended.
	if i := strings.Index(string(data), "crew:\n  context_high:"); i >= 0 {
		data = data[:i]
	}
	data = append(data, fmt.Sprintf("crew:\n  context_high: %d\n", n)...)
	if err := os.WriteFile(ws.WorkspaceFile(), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ws.LoadConfig(); err != nil {
		t.Fatal(err)
	}
}

func openContextHigh(t *testing.T, ws *store.Workspace) []box.Incident {
	t.Helper()
	view, err := box.Load(ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	var out []box.Incident
	for _, inc := range box.OpenIncidents(view, fixtureCrew) {
		if inc.Kind == box.IncidentContextHigh {
			out = append(out, inc)
		}
	}
	return out
}

// A working crew whose latest context is past crew.context_high gets one
// context_high incident (CON-73), and only one however many polls see it.
func TestCheckBudgetsFilesContextHighOnceWhileTheCrewWorks(t *testing.T) {
	f := newFixture(t)
	f.ingest(t)
	// The fixture crew handed back last; a working line after it makes
	// this a crew mid fix round.
	if err := f.ws.AppendStatus(fixtureProject, fixtureCrew, "working: fixing F1"); err != nil {
		t.Fatal(err)
	}
	var latest int64
	if err := f.db.SQL().QueryRow(`SELECT context_tokens_after FROM turn WHERE actor_id = ? AND context_tokens_after > 0 ORDER BY started_at DESC, ordinal DESC LIMIT 1`, f.crewActor()).Scan(&latest); err != nil {
		t.Fatalf("latest context: %v", err)
	}
	if latest <= 1000 {
		t.Fatalf("fixture's latest context is %d; the test needs one above its threshold", latest)
	}
	setContextHigh(t, f.ws, 1000)

	for i := 0; i < 2; i++ {
		if err := f.ing.CheckBudgets(context.Background()); err != nil {
			t.Fatalf("CheckBudgets #%d: %v", i+1, err)
		}
	}
	open := openContextHigh(t, f.ws)
	if len(open) != 1 {
		t.Fatalf("open context_high incidents = %+v, want exactly one", open)
	}
	if want := timeline.ContextHighText(latest, 1000); open[0].Text != want {
		t.Fatalf("incident text = %q, want %q", open[0].Text, want)
	}

	// A fresh session brings the context under the threshold: the
	// incident is resolved, so a later overrun can be filed again.
	setContextHigh(t, f.ws, int(latest)+1)
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets (under): %v", err)
	}
	if open := openContextHigh(t, f.ws); len(open) != 0 {
		t.Fatalf("incident still open with the context under the threshold: %+v", open)
	}
}

// A crew that has handed back is not reported - the Mate's send opens its
// fresh session - and neither is a closed one; the check is off below zero.
func TestCheckBudgetsFilesNoContextHighForARestingOrClosedCrew(t *testing.T) {
	f := newFixture(t)
	f.ingest(t)
	setContextHigh(t, f.ws, 1000)
	// The fixture's last status line is the hand-back.
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	if open := openContextHigh(t, f.ws); len(open) != 0 {
		t.Fatalf("a resting crew got %+v", open)
	}

	meta, err := f.ws.ReadCrewMeta(fixtureProject, fixtureCrew)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.ws.AppendStatus(fixtureProject, fixtureCrew, "working: again"); err != nil {
		t.Fatal(err)
	}
	meta["state"] = "finished"
	if err := f.ws.WriteCrewMeta(fixtureProject, fixtureCrew, meta); err != nil {
		t.Fatal(err)
	}
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	if open := openContextHigh(t, f.ws); len(open) != 0 {
		t.Fatalf("a closed crew got %+v", open)
	}

	meta["state"] = "spawned"
	if err := f.ws.WriteCrewMeta(fixtureProject, fixtureCrew, meta); err != nil {
		t.Fatal(err)
	}
	setContextHigh(t, f.ws, -1)
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	if open := openContextHigh(t, f.ws); len(open) != 0 {
		t.Fatalf("the check filed %+v while turned off", open)
	}
}

// The ledger numbers a crew's sessions (CON-73, schema 4): the fixture's
// one transcript is session 1.
func TestIngestNumbersTheCrewsSessions(t *testing.T) {
	f := newFixture(t)
	f.ingest(t)
	var n int
	if err := f.db.SQL().QueryRow(`SELECT crew_session FROM session WHERE actor_id = ?`, f.crewActor()).Scan(&n); err != nil {
		t.Fatalf("crew_session: %v", err)
	}
	if n != 1 {
		t.Fatalf("crew_session = %d, want 1", n)
	}
	// A second session, with a later first turn, is number 2 after any
	// batch touches the actor's sessions; the first keeps its number.
	if _, err := f.db.SQL().Exec(`INSERT INTO session(id, actor_id, harness_session_id) VALUES (?, ?, 'later')`, f.crewActor()+"#later", f.crewActor()); err != nil {
		t.Fatal(err)
	}
	if _, err := f.db.SQL().Exec(`INSERT INTO turn(id, actor_id, session_id, started_at) VALUES (?, ?, ?, '2027-01-01T00:00:00.000000000Z')`, f.crewActor()+"#later-t", f.crewActor(), f.crewActor()+"#later"); err != nil {
		t.Fatal(err)
	}
	f.ingest(t)
	rows, err := f.db.SQL().Query(`SELECT harness_session_id, crew_session FROM session WHERE actor_id = ? ORDER BY crew_session`, f.crewActor())
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var id string
		var n int
		if err := rows.Scan(&id, &n); err != nil {
			t.Fatal(err)
		}
		got = append(got, fmt.Sprintf("%s=%d", id, n))
	}
	if want := []string{"01a0b944-33bb-7503-9db7-cd51ff61855a=1", "later=2"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("sessions = %v, want %v", got, want)
	}
}
