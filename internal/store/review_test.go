package store

import (
	"reflect"
	"strings"
	"testing"
)

func TestReviewFindingsReadsTheFindingsAndTheOpenPriorOnes(t *testing.T) {
	report := "## Summary\nx\n## Prior findings\n- F1 closed: fixed in a.go:3\n- F2 open: b.go:9 still leaks\n  the handle\n- F3: closed, see c.go\n## Findings\n1. d.go:4 shadows err\n2) e.go:1 unused import\n```\n- not a finding, fenced\n```\n## Verdict\nchanges-requested\n"
	want := []string{"d.go:4 shadows err", "e.go:1 unused import", "b.go:9 still leaks the handle"}
	if got := ReviewFindings(report); !reflect.DeepEqual(got, want) {
		t.Fatalf("ReviewFindings = %q, want %q", got, want)
	}
	marks := ReviewPriorFindings(report)
	if !reflect.DeepEqual(marks, map[int]string{1: "closed", 2: "open", 3: "closed"}) {
		t.Fatalf("ReviewPriorFindings = %v", marks)
	}
}

func TestReviewFindingsTreatsNoneAndMissingAsEmpty(t *testing.T) {
	for _, report := range []string{"## Summary\nfine\n## Verdict\npass\n", "## Findings\nnone\n## Verdict\npass\n", "## Findings\n- None\n"} {
		if got := ReviewFindings(report); len(got) != 0 {
			t.Fatalf("%q: findings %q", report, got)
		}
	}
	if got := ReviewPriorFindings("## Prior findings\n- fixed everything\n"); len(got) != 0 {
		t.Fatalf("prose marked a finding: %v", got)
	}
}

// A task knows its reviewers two ways (CON-75): `reviews=` on its own
// meta, written by mate review, and each reviewer's review.meta naming it,
// which reviewers started before `reviews=` existed still have. Both link,
// once, and a crew never reviews itself.
func TestReviewLinksReadBothRecordsOnce(t *testing.T) {
	w, err := Init(t.TempDir(), Defaults{})
	if err != nil {
		t.Fatal(err)
	}
	for id, meta := range map[string]map[string]string{
		"k3":   {"task": "build it"},
		"rev1": {"task": "Review k3 at abc"},
		"rev2": {"task": "Review k3 at def"},
		"solo": {"task": "no review"},
	} {
		if err := w.WriteCrewMeta("shop", id, meta); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.AddCrewReview("shop", "k3", "rev1"); err != nil {
		t.Fatal(err)
	}
	if err := w.AddCrewReview("shop", "k3", "rev1"); err != nil {
		t.Fatal(err)
	}
	meta, _ := w.ReadCrewMeta("shop", "k3")
	if meta[MetaReviews] != "rev1" {
		t.Fatalf("reviews = %q after adding rev1 twice", meta[MetaReviews])
	}
	// rev2 was started before reviews= existed: only its review.meta links it.
	if err := WriteMeta(w.CrewReviewMeta("shop", "rev2"), map[string]string{"crew": "k3", "head": "def"}); err != nil {
		t.Fatal(err)
	}
	if err := WriteMeta(w.CrewReviewMeta("shop", "rev1"), map[string]string{"crew": "k3", "head": "abc"}); err != nil {
		t.Fatal(err)
	}

	reviewersOf, reviewOf, err := w.ReviewLinks("shop")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(reviewersOf["k3"], ","); got != "rev1,rev2" {
		t.Fatalf("reviewers of k3 = %q, want rev1,rev2 from both records", got)
	}
	if reviewOf["rev1"] != "k3" || reviewOf["rev2"] != "k3" || reviewOf["solo"] != "" || len(reviewersOf["solo"]) != 0 {
		t.Fatalf("reviewOf = %v, reviewersOf = %v", reviewOf, reviewersOf)
	}
}
