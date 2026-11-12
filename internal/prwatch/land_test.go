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

// A pull request merged on GitHub is landed by the app, and the Mate hears
// of it, with what the landing left for it, in a note the digest carries
// rather than a line of the watcher's own (CON-194).
func TestAMergedPullRequestIsLandedAndNoted(t *testing.T) {
	e := newEnv(t)
	calls := e.landed(landOutput, nil)

	if *calls != 1 {
		t.Fatalf("Land ran %d time(s), want once", *calls)
	}
	if items := e.prItems(); len(items) != 0 {
		t.Fatalf("outbox = %+v, want no line for the Mate", items)
	}
	notes, _, err := e.ws.ReadAppNotes(project, 0)
	if err != nil || len(notes) != 1 || !notes[0].Wake {
		t.Fatalf("app notes = %+v, %v; want one that wakes the Mate", notes, err)
	}
	for _, want := range []string{"landed k3 (" + prURL + " merged on GitHub as", "left for you: tracker: HU-7 not updated", "next: 1 crew(s) open"} {
		if !strings.Contains(notes[0].Text, want) {
			t.Fatalf("note %q does not carry %q", notes[0].Text, want)
		}
	}
}

// A landing that fails leaves the Mate the line it always got.
func TestAFailedLandingLeavesItToTheMate(t *testing.T) {
	e := newEnv(t)
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
