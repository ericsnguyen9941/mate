package dashboard

import (
	"net/http"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// A day opened up adds up to its row on the board: the actors' tokens are
// the day's tokens, split the same way, and the tasks are the ones the
// row counted as spawned.
func TestBoardDayListsWhatTheRowCounts(t *testing.T) {
	f := newFixture(t)
	var board BoardResponse
	f.get(t, "/api/board?days=7", http.StatusOK, &board)
	row := board.Total[6]
	var got BoardDayResponse
	f.get(t, "/api/board/day?date="+row.Date, http.StatusOK, &got)

	if got.Date != row.Date || got.Zone != board.Zone || len(got.Errors) != 0 {
		t.Fatalf("date = %q, zone = %q, errors = %v", got.Date, got.Zone, got.Errors)
	}
	var mate, crew int64
	for _, a := range got.Actors {
		if a.Turns == 0 || a.Project != fixtureProject {
			t.Errorf("actor %+v, want turns in %s", a, fixtureProject)
		}
		switch a.Kind {
		case timeline.ActorMate:
			mate += a.Tokens.Total
			if a.Name != "" {
				t.Errorf("the Mate's name = %q, want empty", a.Name)
			}
		case timeline.ActorCrew:
			crew += a.Tokens.Total
		}
	}
	if mate != row.MateTokens || crew != row.CrewTokens {
		t.Errorf("actors' mate/crew tokens = %d/%d, want the row's %d/%d", mate, crew, row.MateTokens, row.CrewTokens)
	}
	for i := 1; i < len(got.Actors); i++ {
		if got.Actors[i-1].Tokens.Total < got.Actors[i].Tokens.Total {
			t.Errorf("actors not biggest first: %+v", got.Actors)
		}
	}

	opened := 0
	for _, task := range got.Tasks {
		if task.Opened {
			opened++
		}
		if task.Crew == fixtureCrew && (task.DayTokens.Total != task.Tokens.Total || task.DayTurns != task.Turns) {
			t.Errorf("%s day tokens/turns = %d/%d, want its whole %d/%d: it lived only today",
				task.Crew, task.DayTokens.Total, task.DayTurns, task.Tokens.Total, task.Turns)
		}
	}
	if opened != sumCounts(row.Spawned) || opened == 0 {
		t.Errorf("opened tasks = %d, want the row's %d", opened, sumCounts(row.Spawned))
	}
}

func TestBoardDayListsItsIncidents(t *testing.T) {
	f := newFixture(t)
	writer, err := db.Open(f.ws)
	if err != nil {
		t.Fatal(err)
	}
	crew := timeline.CrewActorID(fixtureProject, fixtureCrew)
	for i, c := range []struct {
		kind     string
		at       time.Time
		resolved time.Time
	}{
		{"stale", fixtureNow.Add(-time.Hour), fixtureNow.Add(-30 * time.Minute)},
		{"budget", fixtureNow.Add(-2 * time.Hour), time.Time{}},
		{"budget", fixtureNow.Add(-24 * time.Hour), time.Time{}}, // yesterday
	} {
		var resolved any
		if !c.resolved.IsZero() {
			resolved = db.FormatTime(c.resolved)
		}
		if _, err := writer.SQL().Exec(`INSERT INTO incident(id, project, actor_id, kind, opened_event_id, opened_at, resolved_at) VALUES (?, ?, ?, ?, 0, ?, ?)`,
			"test-"+string(rune('a'+i)), fixtureProject, crew, c.kind, db.FormatTime(c.at), resolved); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var got BoardDayResponse
	f.get(t, "/api/board/day?date="+fixtureNow.Format(time.DateOnly), http.StatusOK, &got)
	if len(got.Incidents) != 2 {
		t.Fatalf("incidents = %+v, want today's two", got.Incidents)
	}
	first, second := got.Incidents[0], got.Incidents[1]
	if first.Kind != "budget" || first.ResolvedAt != "" || first.ActorKind != timeline.ActorCrew || first.Name != fixtureCrew {
		t.Errorf("first = %+v, want the open budget incident of %s", first, fixtureCrew)
	}
	if second.Kind != "stale" || second.ResolvedAt == "" {
		t.Errorf("second = %+v, want the resolved stale one", second)
	}
}

// An empty day is empty lists, not nulls; a date that is not one is a 400.
func TestBoardDayParameter(t *testing.T) {
	f := newFixture(t)
	var got BoardDayResponse
	f.get(t, "/api/board/day?date=2026-01-01", http.StatusOK, &got)
	if got.Actors == nil || got.Tasks == nil || got.Incidents == nil || len(got.Actors)+len(got.Tasks)+len(got.Incidents) != 0 {
		t.Errorf("an empty day = %+v, want empty lists", got)
	}
	for _, q := range []string{"", "?date=yesterday", "?date=2026-13-01"} {
		f.get(t, "/api/board/day"+q, http.StatusBadRequest, nil)
	}
}

func sumCounts(m map[string]int) int {
	n := 0
	for _, v := range m {
		n += v
	}
	return n
}
