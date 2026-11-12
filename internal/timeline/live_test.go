package timeline_test

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// A poll names the files mate does not write that the core may watch: the
// open crew's status file and transcript, and the running Mate's
// transcript; a crew closed drops out at the next poll (CON-266).
func TestLiveFilesFollowOpenCrews(t *testing.T) {
	f := newFixture(t)
	polled := f.ing.Polled()
	f.ingest(t)
	select {
	case <-polled:
	default:
		t.Fatal("Polled not signalled after a poll")
	}
	want := []timeline.LiveFile{
		{Project: fixtureProject, Kind: timeline.LiveTranscript, Path: abs(t, claudeFixture)},
		{Project: fixtureProject, Crew: fixtureCrew, Kind: timeline.LiveStatus, Path: f.ws.CrewStatus(fixtureProject, fixtureCrew)},
		{Project: fixtureProject, Crew: fixtureCrew, Kind: timeline.LiveTranscript, Path: abs(t, codexFixture)},
	}
	sort.Slice(want, func(a, b int) bool { return want[a].Path < want[b].Path })
	if got := f.ing.LiveFiles(); !reflect.DeepEqual(got, want) {
		t.Fatalf("LiveFiles = %+v\nwant %+v", got, want)
	}

	meta, err := f.ws.ReadCrewMeta(fixtureProject, fixtureCrew)
	if err != nil {
		t.Fatal(err)
	}
	meta["state"] = "finished"
	meta["stopped_at"] = fixtureNow.Format(time.RFC3339)
	if err := f.ws.WriteCrewMeta(fixtureProject, fixtureCrew, meta); err != nil {
		t.Fatal(err)
	}
	f.ingest(t)
	if got := f.ing.LiveFiles(); len(got) != 1 || got[0].Crew != "" {
		t.Fatalf("LiveFiles = %+v, want the Mate's transcript alone once the crew closed", got)
	}
}
