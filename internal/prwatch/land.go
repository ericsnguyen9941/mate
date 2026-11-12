package prwatch

import (
	"context"
	"fmt"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/outbox"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// landLeftover are the words `mate land` prints on a line that left something
// for the Mate to do: a record step that did not go through, or a judgment
// the command does not make (cmd/mate/land.go).
var landLeftover = []string{"yourself", "not written", "not moved", "not recorded", "not noted", "predictions:"}

// landForMate lands a crew whose pull request was merged on GitHub, the way
// the Mate would with `mate land` (CON-194). On mihuv2 (2026-11-04..10) 49
// pr-merged lines opened a turn and none of them led to a decision: the Mate
// ran `mate land`, or nothing at all because it had landed the crew itself
// a moment before.
//
// It returns quiet when the Mate need not hear of it now: the landing went
// through and some other crew is still at work, so a turn is coming anyway.
// What the landing printed for the Mate then goes into `mate/app-notes.log`,
// which the next digest carries. Otherwise it returns the line to wake the
// Mate with: the landing failed, or nothing else is running and the work has
// to be moved on.
//
// Only an auto project is landed. A manual one is the captain's, and gets the
// line it always got.
func (d Deps) landForMate(ctx context.Context, project, crew, url string, pr string) (line string, quiet bool) {
	if d.Land == nil || !d.WS.Auto(project) {
		return "", false
	}
	if meta, err := d.WS.ReadCrewMeta(project, crew); err == nil && crewstate.Declare(crewstate.Declaration{Meta: meta}).Closed() {
		// The Mate landed it itself while the watcher was looking.
		d.logf("%s is closed already; the Mate landed it", crew)
		return "", true
	}
	out, err := d.Land(ctx, project, crew)
	if err != nil {
		d.logf("mate land %s %s: %v", project, crew, err)
		return "", false
	}
	left, next := landReport(out)
	d.logf("landed %s: %s", crew, strings.Join(append(left, next), "; "))

	view, err := box.Load(d.WS, project)
	working := false
	if err == nil {
		for other := range view.ByCrew {
			if other != "" && other != crew && other != outbox.MateCrew && box.AtWork(view, other) {
				working = true
				break
			}
		}
	}
	what := fmt.Sprintf("landed %s (%s merged on GitHub as %s)", crew, url, pr)
	if len(left) > 0 {
		what += "; left for you: " + strings.Join(left, "; ")
	}
	if working {
		if err := d.WS.AppendAppNote(project, store.AppNote{Time: d.now(), Text: what}); err == nil {
			return "", true
		}
	}
	line = "pr-merged: the app " + what
	if next != "" {
		line += "; " + next
	}
	return line + "; no other crew is working, so keep the work moving", false
}

// landReport splits what `mate land` printed into the lines that left
// something for the Mate and its closing `next:` line.
func landReport(out string) (left []string, next string) {
	for _, l := range strings.Split(out, "\n") {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "next:") {
			next = l
			continue
		}
		for _, w := range landLeftover {
			if strings.Contains(l, w) {
				left = append(left, l)
				break
			}
		}
	}
	return left, next
}
