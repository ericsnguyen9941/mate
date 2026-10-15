package store

import (
	"reflect"
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
