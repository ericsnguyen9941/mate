package timeline_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/gitx"
)

// TestAFailingPassMarksItsProjectStale: a pass that cannot commit leaves its
// project's timeline behind the files, and the next reader of the timeline
// has to be told so (CON-139). The mark stays, dated by the first failure,
// until a pass commits.
func TestAFailingPassMarksItsProjectStale(t *testing.T) {
	f := newFixture(t)
	if _, err := f.db.SQL().Exec(`CREATE TRIGGER refuse BEFORE INSERT ON event
		BEGIN SELECT RAISE(ABORT, 'refused for the test'); END`); err != nil {
		t.Fatal(err)
	}
	ing := f.boundedIngester(gitx.ExecRunner{}, 0)
	before := time.Now()
	if err := ing.Ingest(context.Background()); err == nil {
		t.Fatal("Ingest = nil with every event insert refused")
	}
	stale := ing.Stale()
	if len(stale) != 1 || stale[0].Project != fixtureProject {
		t.Fatalf("Stale = %+v, want %s alone", stale, fixtureProject)
	}
	since := stale[0].Since
	if since.Before(before) || !strings.Contains(stale[0].Error, "refused for the test") {
		t.Fatalf("Stale = %+v, want the failure dated now and named", stale[0])
	}

	_ = ing.Ingest(context.Background())
	if again := ing.Stale(); len(again) != 1 || !again[0].Since.Equal(since) {
		t.Fatalf("a second failing pass moved the mark: %+v, want since %s", again, since)
	}

	if _, err := f.db.SQL().Exec(`DROP TRIGGER refuse`); err != nil {
		t.Fatal(err)
	}
	if err := ing.Ingest(context.Background()); err != nil {
		t.Fatalf("Ingest after the inserts were allowed again: %v", err)
	}
	if after := ing.Stale(); len(after) != 0 {
		t.Fatalf("Stale = %+v after a pass committed, want none", after)
	}
}
