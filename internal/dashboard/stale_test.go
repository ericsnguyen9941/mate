package dashboard

import (
	"net/http"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// TestEventsSayWhichTimelineIsStale: the Office polls /api/events, so that is
// where it learns that a project's timeline stopped moving and every scene it
// draws for that project is old (CON-139). A project filter keeps only that
// project's word.
func TestEventsSayWhichTimelineIsStale(t *testing.T) {
	f := newFixture(t)
	f.server.deps.TimelineStale = func() []timeline.Stale {
		return []timeline.Stale{
			{Project: "other", Since: mustTime("2026-09-19T09:00:00Z"), Error: "disk full"},
			{Project: fixtureProject, Since: mustTime("2026-09-19T09:30:43Z"), Error: "refused"},
		}
	}

	var all EventsResponse
	f.get(t, "/api/events?since=0", http.StatusOK, &all)
	if len(all.TimelineStale) != 2 {
		t.Fatalf("timeline_stale = %+v, want both projects", all.TimelineStale)
	}

	var one EventsResponse
	f.get(t, "/api/events?since=0&project="+fixtureProject, http.StatusOK, &one)
	want := StaleTimeline{Project: fixtureProject, Since: "2026-09-19T09:30:43.000000000Z", Error: "refused"}
	if len(one.TimelineStale) != 1 || one.TimelineStale[0] != want {
		t.Fatalf("timeline_stale = %+v, want [%+v]", one.TimelineStale, want)
	}

	f.server.deps.TimelineStale = nil
	var quiet EventsResponse
	f.get(t, "/api/events?since=0", http.StatusOK, &quiet)
	if len(quiet.TimelineStale) != 0 {
		t.Fatalf("timeline_stale = %+v with no ingest to ask", quiet.TimelineStale)
	}
}
