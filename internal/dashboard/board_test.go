package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// The board's day is the project page's numbers bucketed by local date: its
// tokens are the Mate's and the tasks' totals, and a task spawned that day
// is counted under its role.
func TestBoardSumsTheProjectIntoItsDay(t *testing.T) {
	f := newFixture(t)
	var project ProjectResponse
	f.get(t, "/api/projects/"+fixtureProject, http.StatusOK, &project)
	var got BoardResponse
	f.get(t, "/api/board?days=7", http.StatusOK, &got)

	if got.Days != 7 || len(got.Dates) != 7 || len(got.Total) != 7 || got.Zone != "+00:00" {
		t.Fatalf("days = %d, dates = %v, total = %d, zone = %q", got.Days, got.Dates, len(got.Total), got.Zone)
	}
	if got.Dates[0] != "2026-09-13" || got.Dates[6] != "2026-09-19" {
		t.Fatalf("dates = %v, want 2026-09-13 .. 2026-09-19", got.Dates)
	}
	if len(got.Projects) != 1 || got.Projects[0].Name != fixtureProject || got.Projects[0].Error != "" {
		t.Fatalf("projects = %+v", got.Projects)
	}

	today := got.Total[6]
	crewTokens := int64(0)
	for _, task := range project.Tasks {
		crewTokens += task.Tokens.Total
	}
	if today.MateTokens != project.Mate.Tokens.Total || today.CrewTokens != crewTokens {
		t.Errorf("mate/crew tokens = %d/%d, want the project page's %d/%d", today.MateTokens, today.CrewTokens, project.Mate.Tokens.Total, crewTokens)
	}
	if today.Tokens.Total != today.MateTokens+today.CrewTokens || today.Tokens.Total == 0 {
		t.Errorf("tokens.total = %d, want mate+crew = %d and not zero", today.Tokens.Total, today.MateTokens+today.CrewTokens)
	}
	if today.Spawned["ship"] != 1 || len(today.Spawned) != 1 {
		t.Errorf("spawned = %v, want one ship (%s)", today.Spawned, fixtureCrew)
	}
	for _, d := range got.Total[:6] {
		if d.Tokens.Total != 0 || len(d.Spawned)+len(d.Closed)+len(d.Reviews)+len(d.Incidents) != 0 {
			t.Errorf("%s = %+v, want an empty day", d.Date, d)
		}
	}
}

// A day is the server clock's zone, not UTC's: at UTC+14 the fixture's
// afternoon is already the next morning.
func TestBoardDaysAreTheClocksZone(t *testing.T) {
	f := newFixture(t)
	zone := time.FixedZone("+14", 14*3600)
	server, err := New(Options{Workspace: f.ws, DB: f.read, Deps: Deps{Now: func() time.Time { return fixtureNow.In(zone) }}})
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	resp, err := ts.Client().Get(ts.URL + "/api/board?days=2")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var got BoardResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Zone != "+14:00" || got.Dates[0] != "2026-09-19" || got.Dates[1] != "2026-09-20" {
		t.Fatalf("zone = %q, dates = %v", got.Zone, got.Dates)
	}
	if got.Total[0].Tokens.Total != 0 || got.Total[1].Tokens.Total == 0 {
		t.Fatalf("tokens by day = %d, %d; want them all on 2026-09-20", got.Total[0].Tokens.Total, got.Total[1].Tokens.Total)
	}
}

func TestBoardCountsIncidentsByKind(t *testing.T) {
	f := newFixture(t)
	writer, err := db.Open(f.ws)
	if err != nil {
		t.Fatal(err)
	}
	crew := timeline.CrewActorID(fixtureProject, fixtureCrew)
	for i, c := range []struct {
		kind string
		at   time.Time
	}{
		{"stale", fixtureNow.Add(-time.Hour)},
		{"budget", fixtureNow.Add(-2 * time.Hour)},
		{"budget", fixtureNow.Add(-24 * time.Hour)},
		{"stale", fixtureNow.Add(-30 * 24 * time.Hour)}, // outside the window
	} {
		if _, err := writer.SQL().Exec(`INSERT INTO incident(id, project, actor_id, kind, opened_event_id, opened_at) VALUES (?, ?, ?, ?, 0, ?)`,
			"test-"+string(rune('a'+i)), fixtureProject, crew, c.kind, db.FormatTime(c.at)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var got BoardResponse
	f.get(t, "/api/board?days=7", http.StatusOK, &got)
	if inc := got.Total[6].Incidents; inc["stale"] != 1 || inc["budget"] != 1 || len(inc) != 2 {
		t.Errorf("today's incidents = %v, want 1 stale and 1 budget", inc)
	}
	if inc := got.Total[5].Incidents; inc["budget"] != 1 || len(inc) != 1 {
		t.Errorf("yesterday's incidents = %v, want 1 budget", inc)
	}
}

func TestBoardDaysParameter(t *testing.T) {
	f := newFixture(t)
	for _, c := range []struct {
		query string
		days  int
	}{{"", 14}, {"?days=abc", 14}, {"?days=0", 14}, {"?days=30", 30}, {"?days=500", 90}} {
		var got BoardResponse
		f.get(t, "/api/board"+c.query, http.StatusOK, &got)
		if got.Days != c.days || len(got.Dates) != c.days {
			t.Errorf("%q: days = %d with %d dates, want %d", c.query, got.Days, len(got.Dates), c.days)
		}
	}
}

// A closed task is worded the way the cabinet words it, without merged,
// and a reviewer by its verdict.
func TestBoardWordsAClosedTask(t *testing.T) {
	for _, c := range []struct {
		task Task
		want string
	}{
		{Task{CloseState: "finished"}, "finished"},
		{Task{CloseState: "finished", MergedAt: "2026-09-19T10:00:00.000000000Z"}, "finished"},
		{Task{CloseState: "failed"}, "discarded"},
		{Task{CloseState: "failed", FailedReason: "herdr refused"}, "didnt_start"},
		{Task{ClosedAt: "2026-09-19T10:00:00.000000000Z"}, "finished"},
	} {
		if got := closedWord(c.task); got != c.want {
			t.Errorf("closedWord(%+v) = %q, want %q", c.task, got, c.want)
		}
	}
	for v, want := range map[string]string{"pass": "pass", "changes-requested": "changes-requested", "pending": "no_verdict", "": "no_verdict"} {
		if got := reviewWord(v); got != want {
			t.Errorf("reviewWord(%q) = %q, want %q", v, got, want)
		}
	}
}
