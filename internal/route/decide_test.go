package route

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/notice"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// fakeJudge answers every question with one label, and remembers the text.
type fakeJudge struct {
	label string
	conf  float64
	err   error
	asked []string
}

func (j *fakeJudge) answer(text string) (notice.Result, error) {
	j.asked = append(j.asked, text)
	return notice.Result{Label: j.label, Confidence: j.conf}, j.err
}

func (j *fakeJudge) ClassifyHandback(_ context.Context, text string) (notice.Result, error) {
	return j.answer(text)
}

func (j *fakeJudge) ClassifyFindings(_ context.Context, text string) (notice.Result, error) {
	return j.answer(text)
}

const cleanHandback = `Closes HU-489

## For the next crew
Pick flow archives.
Learning: none (HU-489).

| criterion | pass/fail | command | output |
| --- | --- | --- | --- |
| Pick archives | pass | make pre-merge | passed |
| Gate green | pass | make check-zone | ok |

## Deviations from Build
Updated two docs files for stale wording.

## Still open
The trash glyph lives outside this zone, so I left it.
`

func ship(text string) Input {
	return Input{Crew: "hu-489", Kind: KindShip, Text: text, Review: store.ReviewOff}
}

func TestDecide(t *testing.T) {
	routine := &fakeJudge{label: notice.LabelRoutine, conf: .93}
	for _, tc := range []struct {
		name   string
		in     Input
		judge  Judge
		action Action
		reason string
	}{
		{"a clean hand-back Jev finds routine is landed", ship(cleanHandback), routine, ActionLand, "2 criteria pass"},
		{"Jev's veto leaves it to the Mate", ship(cleanHandback), &fakeJudge{label: notice.LabelNeedsJudgment, conf: .9}, ActionMate, "Jev vetoes: needs_judgment"},
		{"a routine label below the threshold is a veto", ship(cleanHandback), &fakeJudge{label: notice.LabelRoutine, conf: .7}, ActionMate, "Jev vetoes: routine at 0.70"},
		{"Jev not answering is a veto", ship(cleanHandback), &fakeJudge{err: errors.New("timeout")}, ActionMate, "Jev did not answer"},
		{"no Jev clears nothing", ship(cleanHandback), nil, ActionMate, "Jev is not configured"},
		{"a failing criterion goes to the Mate", ship(strings.Replace(cleanHandback, "| Gate green | pass", "| Gate green | fail", 1)), routine, ActionMate, "not marked pass"},
		{"no criteria table goes to the Mate", ship("## Still open\nnone"), routine, ActionMate, "no criteria table"},
		{"review on is the Mate's", Input{Crew: "k", Kind: KindShip, Text: cleanHandback, Review: store.ReviewOn}, routine, ActionMate, "review is on"},
		{"an undecided review setting is the Mate's", Input{Crew: "k", Kind: KindShip, Text: cleanHandback}, routine, ActionMate, "has not chosen"},
		{"a crew under review is the Mate's", Input{Crew: "k", Kind: KindShip, Text: cleanHandback, Review: store.ReviewOff, Reviewed: true}, routine, ActionMate, "review was started"},
		{"a scout's report is the Mate's", Input{Crew: "s", Kind: KindScout, Text: "## Deliverable\nx"}, routine, ActionMate, "scout"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := Decide(context.Background(), tc.in, tc.judge)
			if d.Action != tc.action || !strings.Contains(d.Reason, tc.reason) {
				t.Fatalf("Decide = %+v, want %s with reason containing %q", d, tc.action, tc.reason)
			}
			if d.Action == ActionLand && d.Target != tc.in.Crew {
				t.Fatalf("land target = %q, want the crew itself", d.Target)
			}
		})
	}
}

// Jev reads only what the app has not: the deviations and what is open.
func TestDecideAsksJevAboutDeviationsAndOpenItemsOnly(t *testing.T) {
	j := &fakeJudge{label: notice.LabelRoutine, conf: .9}
	Decide(context.Background(), ship(cleanHandback), j)
	if len(j.asked) != 1 {
		t.Fatalf("asked Jev %d time(s)", len(j.asked))
	}
	want := "## Deviations\nUpdated two docs files for stale wording.\n\n## Still open\nThe trash glyph lives outside this zone, so I left it."
	if j.asked[0] != want {
		t.Fatalf("asked Jev\n%s\nwant\n%s", j.asked[0], want)
	}
}

const passingReview = `## Verdict
pass

## Findings
none

## Verification gaps
none
`

const changesReview = `## Verdict
changes-requested

## Findings
- The toast text is hardcoded in BoardLayer; read it from BoardText.
- The gate ran at ff6b60d, not the head of the branch.

## Verification gaps
none
`

func TestDecideReview(t *testing.T) {
	inScope := &fakeJudge{label: notice.LabelInScopeFix, conf: .9}
	review := func(text string) Input {
		return Input{Crew: "hu-4-r1", Kind: KindReview, Text: text, ReviewOf: "hu-4", ReportPath: "/ws/crews/hu-4-r1/report.md"}
	}
	d := Decide(context.Background(), review(passingReview), inScope)
	if d.Action != ActionLand || d.Target != "hu-4" || d.Reviewer != "hu-4-r1" {
		t.Fatalf("passing review: %+v, want land hu-4 after hu-4-r1", d)
	}
	if len(inScope.asked) != 0 {
		t.Fatalf("asked Jev about a passing review with no findings")
	}

	d = Decide(context.Background(), review(changesReview), inScope)
	if d.Action != ActionSend || d.Target != "hu-4" || !strings.Contains(d.Text, "/ws/crews/hu-4-r1/report.md") {
		t.Fatalf("changes requested: %+v, want the findings sent to hu-4", d)
	}

	d = Decide(context.Background(), review(changesReview), &fakeJudge{label: notice.LabelNeedsJudgment, conf: .9})
	if d.Action != ActionMate {
		t.Fatalf("findings Jev vetoes: %+v, want the Mate", d)
	}

	d = Decide(context.Background(), review("## Verdict\nincomplete\n"), inScope)
	if d.Action != ActionMate {
		t.Fatalf("incomplete review: %+v, want the Mate", d)
	}
}
