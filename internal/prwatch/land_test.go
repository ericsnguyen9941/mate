package prwatch_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/github"
	"github.com/nguyenngocanh94/mate/internal/prwatch"
)

// landOutput is what `mate land` prints for a crew whose pull request was
// merged on GitHub, on a project with no tracker configured.
const landOutput = `crew stop k3: stopped
learning: written from the hand-back (Learning: none)
backlog: k3 → Done
tracker: HU-7 not updated (no ` + "`tracker:`" + ` block in plugins.yaml); mark it done yourself
next: 1 crew(s) open: k5 · queued: 2 (first: HU-8 board)`

// landed is a watch over a merged pull request in auto mode, with a Land
// that records its calls and answers out, err.
func (e *env) landed(out string, err error) (calls *int) {
	e.t.Helper()
	e.gh.answers = []github.Result{prJSON("MERGED", e.mergeSHA)}
	if err := e.ws.SetMode(project, true); err != nil {
		e.t.Fatal(err)
	}
	n := 0
	d := e.deps()
	d.Land = func(_ context.Context, p, c string) (string, error) {
		if p != project || c != crew {
			e.t.Fatalf("Land(%s, %s), want (%s, %s)", p, c, project, crew)
		}
		n++
		return out, err
	}
	if err := prwatch.Watch(context.Background(), d, project, crew, prURL); err != nil {
		e.t.Fatalf("Watch: %v", err)
	}
	return &n
}

func (e *env) notes() []string {
	e.t.Helper()
	notes, _, err := e.ws.ReadAppNotes(project, 0)
	if err != nil {
		e.t.Fatal(err)
	}
	var out []string
	for _, n := range notes {
		out = append(out, n.Text)
	}
	return out
}

// With another crew still at work, a pull request merged on GitHub is
// landed by the app and the Mate hears of it in its next digest, not in a
// turn of its own (CON-194).
func TestAMergedPullRequestIsLandedQuietlyWhileAnotherCrewWorks(t *testing.T) {
	e := newEnv(t)
	if err := e.ws.AppendStatus(project, "k5", "working: building"); err != nil {
		t.Fatal(err)
	}
	calls := e.landed(landOutput, nil)

	if *calls != 1 {
		t.Fatalf("Land ran %d time(s), want once", *calls)
	}
	if items := e.prItems(); len(items) != 0 {
		t.Fatalf("outbox = %+v, want no line for the Mate", items)
	}
	notes := e.notes()
	if len(notes) != 1 || !strings.Contains(notes[0], "landed k3 ("+prURL+" merged on GitHub as") ||
		!strings.Contains(notes[0], "left for you: tracker: HU-7 not updated") {
		t.Fatalf("app notes = %q", notes)
	}
}

// With nobody else working, the landing is done but the Mate is woken to move
// the work on, with what the landing printed.
func TestAMergedPullRequestWakesTheMateWhenNothingElseRuns(t *testing.T) {
	e := newEnv(t)
	e.landed(landOutput, nil)

	items := e.prItems()
	if len(items) != 1 {
		t.Fatalf("outbox = %+v, want one line for the Mate", items)
	}
	for _, want := range []string{"pr-merged: the app landed k3", "mark it done yourself", "next: 1 crew(s) open", "keep the work moving"} {
		if !strings.Contains(items[0].Text, want) {
			t.Fatalf("line %q does not carry %q", items[0].Text, want)
		}
	}
	if notes := e.notes(); len(notes) != 0 {
		t.Fatalf("app notes = %q, want none when the Mate is woken", notes)
	}
}

// A landing that fails leaves the Mate the line it always got.
func TestAFailedLandingLeavesItToTheMate(t *testing.T) {
	e := newEnv(t)
	if err := e.ws.AppendStatus(project, "k5", "working: building"); err != nil {
		t.Fatal(err)
	}
	e.landed("", errors.New("refused"))

	items := e.prItems()
	if len(items) != 1 || !strings.Contains(items[0].Text, "then run mate land "+project+" "+crew) {
		t.Fatalf("outbox = %+v, want the old pr-merged line", items)
	}
}

// A manual project is the captain's: nothing is landed for it.
func TestAManualProjectIsNotLanded(t *testing.T) {
	e := newEnv(t, prJSON("MERGED", "x"))
	e.gh.answers[0] = prJSON("MERGED", e.mergeSHA)
	if err := e.ws.SetMode(project, false); err != nil {
		t.Fatal(err)
	}
	d := e.deps()
	d.Land = func(context.Context, string, string) (string, error) {
		t.Fatal("Land ran for a manual project")
		return "", nil
	}
	if err := prwatch.Watch(context.Background(), d, project, crew, prURL); err != nil {
		t.Fatalf("Watch: %v", err)
	}
	if items := e.prItems(); len(items) != 1 || !strings.Contains(items[0].Text, "then run mate land") {
		t.Fatalf("outbox = %+v, want the old pr-merged line", items)
	}
}
