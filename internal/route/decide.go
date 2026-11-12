// Package route takes a crew's hand-back off the Mate when there is nothing
// in it for the Mate to decide (CON-273, CON-274).
//
// On mihuv2 (2026-11-04..10) hand-backs opened 233 of the Mate's turns and
// cost 52% of its tokens, and in about 120 of them the Mate decided nothing:
// it landed what the crew handed back, or sent a review's findings back to
// the crew. Every such turn re-read 150-250k tokens of context.
//
// Three layers decide, cheapest first. The app's own checks come first:
// what kind of crew handed back, its criteria table, the project's review
// setting, a review's verdict. Then Jev, which can only veto: it reads the
// hand-back's deviations and open items, or a review's findings, and a
// label other than routine (or no answer) leaves the hand-back to the Mate.
// The Mate gets everything else, and hears in its next digest of what the
// app did.
//
// project.yaml's `autoroute` turns it on per project: off (the default),
// shadow (decide and record in `mate/route.jsonl`, act on nothing) or on.
package route

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/notice"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Kind is what kind of crew handed back.
type Kind string

const (
	KindShip   Kind = "ship"
	KindReview Kind = "review"
	KindScout  Kind = "scout"
)

// Action is what the router decided.
type Action string

const (
	// ActionMate leaves the hand-back to the Mate, as before.
	ActionMate Action = "mate"
	// ActionLand lands Target with `mate land` (with --review for a
	// passing review).
	ActionLand Action = "land"
	// ActionSend sends a review's findings back to Target, the crew it
	// reviewed.
	ActionSend Action = "send"
)

// Threshold is the confidence Jev's routine label needs before the app acts.
// Below it the label counts as a veto.
const Threshold = 0.85

// Input is what the router knows about one hand-back.
type Input struct {
	Crew string
	Kind Kind
	// Text is the crew's hand-back (handback.md) for a ship, its report
	// (report.md) for a reviewer or a scout.
	Text string
	// ReviewOf is the crew a reviewer reviewed.
	ReviewOf string
	// ReportPath is where a reviewer's report is, for the line sent back.
	ReportPath string
	// Reviewed says some reviewer was already started for this ship crew:
	// the Mate started a review cycle, and it decides how the cycle ends.
	Reviewed bool
	// PRMerged says the crew's pull request was merged already, which the
	// watcher lands (internal/prwatch).
	PRMerged bool
	// Review is project.yaml's review setting.
	Review string
}

// Decision is what the router decided about one hand-back, and why.
type Decision struct {
	Crew   string `json:"crew"`
	Kind   Kind   `json:"kind"`
	Action Action `json:"action"`
	// Target is the crew the action is about: the crew itself when it
	// shipped, the crew it reviewed for a reviewer.
	Target string `json:"target,omitempty"`
	// Reviewer is the reviewer whose passing review a land relies on.
	Reviewer string `json:"reviewer,omitempty"`
	// Text is what ActionSend sends.
	Text   string `json:"text,omitempty"`
	Reason string `json:"reason"`
	// Jev and Confidence are Jev's label and confidence, empty when Jev
	// was not asked.
	Jev        string  `json:"jev,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

// Judge is Jev's two hand-back questions (*notice.Client).
type Judge interface {
	ClassifyHandback(ctx context.Context, text string) (notice.Result, error)
	ClassifyFindings(ctx context.Context, text string) (notice.Result, error)
}

var _ Judge = (*notice.Client)(nil)

// Decide decides about one hand-back. A nil judge vetoes everything Jev
// would be asked about.
func Decide(ctx context.Context, in Input, judge Judge) Decision {
	d := Decision{Crew: in.Crew, Kind: in.Kind, Action: ActionMate}
	switch in.Kind {
	case KindScout:
		d.Reason = "a scout's report: what follows from it is the Mate's to decide"
		return d
	case KindReview:
		return decideReview(ctx, in, judge, d)
	}
	switch {
	case in.PRMerged:
		d.Reason = "its pull request is merged already; the watcher lands it"
		return d
	case in.Review == store.ReviewOn:
		d.Reason = "review is on: the Mate starts the review"
		return d
	case in.Review == "":
		d.Reason = "the captain has not chosen whether this project reviews"
		return d
	case in.Reviewed:
		d.Reason = "a review was started for it: the Mate closes the cycle"
		return d
	}
	rows, pass := criteria(in.Text)
	if rows == 0 {
		d.Reason = "the hand-back has no criteria table"
		return d
	}
	if !pass {
		d.Reason = "a criterion is not marked pass"
		return d
	}
	text := sectionsForJudge(in.Text)
	if !ask(ctx, &d, judge, false, text, notice.LabelRoutine) {
		return d
	}
	d.Action, d.Target = ActionLand, in.Crew
	d.Reason = fmt.Sprintf("%d criteria pass and Jev finds the deviations and open items routine", rows)
	return d
}

func decideReview(ctx context.Context, in Input, judge Judge, d Decision) Decision {
	if in.ReviewOf == "" {
		d.Reason = "a reviewer whose crew is not recorded"
		return d
	}
	d.Target = in.ReviewOf
	verdict := strings.ToLower(store.ReviewVerdict(in.Text))
	findings := store.ReviewFindings(in.Text)
	switch {
	case verdict == "pass" && len(findings) == 0:
		if err := store.ValidateReviewPass(in.Text); err != nil {
			d.Reason = "the review passes but its report does not hold: " + err.Error()
			return d
		}
		d.Action, d.Reviewer = ActionLand, in.Crew
		d.Reason = "the review passes with no findings"
		return d
	case verdict == "changes-requested" && len(findings) > 0:
		text := "## Findings\n- " + strings.Join(findings, "\n- ")
		if !ask(ctx, &d, judge, true, text, notice.LabelInScopeFix) {
			return d
		}
		d.Action = ActionSend
		d.Text = fmt.Sprintf("review %s asks for changes: fix each item under ## Findings in %s inside your task, run the gate again, and hand back", in.Crew, in.ReportPath)
		d.Reason = fmt.Sprintf("the review asks for %d change(s) and Jev finds them in scope", len(findings))
		return d
	}
	d.Reason = fmt.Sprintf("review verdict %q with %d finding(s)", verdict, len(findings))
	return d
}

// ask puts one question to Jev, the findings question or the hand-back one,
// and reports whether it answered want with Threshold confidence. Anything
// else is a veto, recorded on d.
func ask(ctx context.Context, d *Decision, judge Judge, findings bool, text, want string) bool {
	if judge == nil {
		d.Reason = "Jev is not configured, so nothing clears it"
		return false
	}
	q := judge.ClassifyHandback
	if findings {
		q = judge.ClassifyFindings
	}
	res, err := q(ctx, text)
	if err != nil {
		d.Reason = "Jev did not answer: " + err.Error()
		return false
	}
	d.Jev, d.Confidence = res.Label, res.Confidence
	if res.Label != want || res.Confidence < Threshold {
		d.Reason = fmt.Sprintf("Jev vetoes: %s at %.2f", res.Label, res.Confidence)
		return false
	}
	return true
}

// criteriaRow is a row of a hand-back's criteria table: `| criterion |
// pass/fail | ...`.
var criteriaRow = regexp.MustCompile(`(?mi)^\|[^|\n]*\|\s*(pass|fail|failed|partial|n/a|skipped|skip|blocked|unverified|not run)\s*\|`)

// criteria counts the rows of the hand-back's criteria table and reports
// whether every one is pass (or n/a).
func criteria(handback string) (rows int, pass bool) {
	pass = true
	for _, m := range criteriaRow.FindAllStringSubmatch(handback, -1) {
		rows++
		if v := strings.ToLower(m[1]); v != "pass" && v != "n/a" {
			pass = false
		}
	}
	return rows, pass
}

// judgedSections are what Jev reads of a hand-back: what the crew did not
// do as briefed and what it left. The rest, the criteria table above all,
// the app has already read.
var judgedSections = []string{"Deviations", "Still open"}

// sectionsForJudge is the hand-back's judged sections, each under its own
// heading, or "none" for a section the crew left out.
func sectionsForJudge(handback string) string {
	var b strings.Builder
	for _, name := range judgedSections {
		body := section(handback, name)
		if body == "" {
			body = "none"
		}
		fmt.Fprintf(&b, "## %s\n%s\n\n", name, body)
	}
	return strings.TrimSpace(b.String())
}

// section is the body of the first `## <name>...` section, up to the next
// `## `. A heading matches by prefix, so "Deviations" finds "Deviations from
// Build".
func section(text, name string) string {
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(l)), "## "+strings.ToLower(name)) {
			continue
		}
		var body []string
		for _, b := range lines[i+1:] {
			if strings.HasPrefix(strings.TrimSpace(b), "## ") {
				break
			}
			body = append(body, b)
		}
		return strings.TrimSpace(strings.Join(body, "\n"))
	}
	return ""
}
