package dashboard

import (
	"net/http"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// A day slot is the board's day: same tokens, split the same way.
func TestBoardTrendDaysAreTheBoardsDays(t *testing.T) {
	f := newFixture(t)
	var board BoardResponse
	f.get(t, "/api/board?days=14", http.StatusOK, &board)
	var got BoardTrendResponse
	f.get(t, "/api/board/trend?slot=1d", http.StatusOK, &got)

	if got.Slot != "1d" || len(got.Buckets) != 14 || got.Zone != "+00:00" {
		t.Fatalf("slot = %q, %d buckets, zone %q", got.Slot, len(got.Buckets), got.Zone)
	}
	if got.From != "2026-09-06T00:00:00Z" || got.To != "2026-09-20T00:00:00Z" {
		t.Fatalf("window = %s .. %s", got.From, got.To)
	}
	for i, b := range got.Buckets {
		d := board.Total[i]
		if b.From[:10] != d.Date || b.Tokens != d.Tokens.Total || b.MateTokens != d.MateTokens || b.CrewTokens != d.CrewTokens {
			t.Errorf("bucket %d = %+v, want the board's %s: %d (%d/%d)", i, b, d.Date, d.Tokens.Total, d.MateTokens, d.CrewTokens)
		}
	}
	if got.Buckets[13].Tokens == 0 {
		t.Error("the fixture's day spent nothing")
	}
}

// An hour slot puts the fixture's turns, all between 10:43 and 10:47, in
// the 10:00 hour, and a finished task in the hour it closed.
func TestBoardTrendHoursZoomIntoADay(t *testing.T) {
	f := newFixture(t)
	writer, err := db.Open(f.ws)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.SQL().Exec(`UPDATE task SET close_state = 'finished', closed_at = ? WHERE crew_actor_id = ?`,
		db.FormatTime(fixtureNow.Add(-3*time.Hour)), timeline.CrewActorID(fixtureProject, fixtureCrew)); err != nil {
		t.Fatal(err)
	}
	writer.Close()

	var day BoardTrendResponse
	f.get(t, "/api/board/trend?slot=1d&from=2026-09-19T00:00:00Z&to=2026-09-19T00:00:00Z", http.StatusOK, &day)
	var got BoardTrendResponse
	f.get(t, "/api/board/trend?slot=1h&from=2026-09-19T00:00:00Z&to=2026-09-19T23:00:00Z", http.StatusOK, &got)
	if len(day.Buckets) != 1 || len(got.Buckets) != 24 || got.To != "2026-09-20T00:00:00Z" {
		t.Fatalf("%d day buckets, %d hours to %s", len(day.Buckets), len(got.Buckets), got.To)
	}
	var tokens int64
	for i, b := range got.Buckets {
		tokens += b.Tokens
		if b.Tokens != 0 && i != 10 {
			t.Errorf("hour %d spent %d, want every token at 10:00", i, b.Tokens)
		}
		if b.Done != 0 && i != 7 {
			t.Errorf("hour %d has %d done, want the task in 07:00", i, b.Done)
		}
	}
	if tokens != day.Buckets[0].Tokens || tokens == 0 {
		t.Errorf("hours sum to %d, want the day's %d", tokens, day.Buckets[0].Tokens)
	}
	if b := got.Buckets[7]; b.Done != 1 || b.DoneTokens == 0 || b.DoneTokens != day.Buckets[0].DoneTokens {
		t.Errorf("07:00 = %+v, want one done with the day's %d done tokens", b, day.Buckets[0].DoneTokens)
	}
}

func TestBoardTrendSlotsAndRanges(t *testing.T) {
	f := newFixture(t)
	for _, c := range []struct {
		query   string
		buckets int
		from    string
	}{
		{"slot=1h", 48, "2026-09-17T11:00:00Z"},
		{"slot=6h", 28, "2026-09-12T12:00:00Z"},
		{"slot=1w", 13, "2026-06-22T00:00:00Z"}, // a Monday
		{"slot=1w&from=2026-09-16T08:00:00Z", 1, "2026-09-14T00:00:00Z"},
		{"slot=1d&project=nope", 14, "2026-09-06T00:00:00Z"},
		{"slot=1d&count=90", 90, "2026-06-22T00:00:00Z"},
		{"slot=1w&count=2", 2, "2026-09-07T00:00:00Z"},
	} {
		var got BoardTrendResponse
		f.get(t, "/api/board/trend?"+c.query, http.StatusOK, &got)
		if len(got.Buckets) != c.buckets || got.From != c.from {
			t.Errorf("%s: %d buckets from %s, want %d from %s", c.query, len(got.Buckets), got.From, c.buckets, c.from)
		}
		if c.query == "slot=1d&project=nope" {
			for _, b := range got.Buckets {
				if b.Tokens != 0 {
					t.Errorf("an unknown project spent %d", b.Tokens)
				}
			}
		}
	}
	for _, q := range []string{"slot=2h", "slot=1d&from=yesterday", "slot=1d&from=2026-09-19T12:00:00Z&to=2026-09-18T00:00:00Z",
		"slot=1h&from=2025-01-01T00:00:00Z", "slot=1d&count=0", "slot=1d&count=501"} {
		f.get(t, "/api/board/trend?"+q, http.StatusBadRequest, nil)
	}
}

// The tasks a slot started are counted there whenever they end: the
// fixture's ship spawned at 10:44 and closed finished at 13:44 is in the
// 10:00 hour with three hours of work, though it closed in the 13:00 one.
func TestBoardTrendCountsWorkFromItsStart(t *testing.T) {
	f := newFixture(t)
	var project ProjectResponse
	f.get(t, "/api/projects/"+fixtureProject, http.StatusOK, &project)
	var cost int64
	for _, task := range project.Tasks {
		if task.Crew == fixtureCrew {
			cost = task.TokensWithReviews
		}
	}
	closedAt := fixtureCrewStart.Add(3 * time.Hour)
	writer, err := db.Open(f.ws)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.SQL().Exec(`UPDATE task SET close_state = 'finished', closed_at = ? WHERE crew_actor_id = ?`,
		db.FormatTime(closedAt), timeline.CrewActorID(fixtureProject, fixtureCrew)); err != nil {
		t.Fatal(err)
	}
	writer.Close()

	var got BoardTrendResponse
	f.get(t, "/api/board/trend?slot=1h&from=2026-09-19T00:00:00Z&to=2026-09-19T23:00:00Z", http.StatusOK, &got)
	want := BoardCohort{Tasks: 1, Done: 1, Tokens: cost, WorkMs: (3 * time.Hour).Milliseconds()}
	if c := got.Buckets[10].Started; c != want || cost == 0 {
		t.Errorf("10:00 started = %+v, want %+v", c, want)
	}
	if got.Buckets[13].Done != 1 || got.Buckets[13].Started.Tasks != 0 {
		t.Errorf("13:00 = %+v, want the close there and no start", got.Buckets[13])
	}

	// Still open, it is counted open with its time so far.
	writer, err = db.Open(f.ws)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.SQL().Exec(`UPDATE task SET close_state = NULL, closed_at = NULL WHERE crew_actor_id = ?`,
		timeline.CrewActorID(fixtureProject, fixtureCrew)); err != nil {
		t.Fatal(err)
	}
	writer.Close()
	f.get(t, "/api/board/trend?slot=1h&from=2026-09-19T00:00:00Z&to=2026-09-19T23:00:00Z", http.StatusOK, &got)
	if c := got.Buckets[10].Started; c.Tasks != 1 || c.Open != 1 || c.Done != 0 || c.WorkMs != fixtureNow.Sub(fixtureCrewStart).Milliseconds() {
		t.Errorf("an open task started = %+v, want open with %d ms so far", c, fixtureNow.Sub(fixtureCrewStart).Milliseconds())
	}
}
