package prwatch

import (
	"context"
	"fmt"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/crewstate"
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
// It returns quiet when the Mate need not be told by this watcher: it
// landed the crew itself already, or the landing went through and its news,
// with what it left for the Mate, is in `mate/app-notes.log`. That note
// wakes the Mate to move the work on the way a hand-back does: at once when
// no other crew is working, with the next hand-back otherwise
// (internal/autopilot). A landing that failed, and a manual project, which is
// the captain's, get the line they always got.
func (d Deps) landForMate(ctx context.Context, project, crew, url, pr string) (quiet bool) {
	if d.Land == nil || !d.WS.Auto(project) {
		return false
	}
	if meta, err := d.WS.ReadCrewMeta(project, crew); err == nil && crewstate.Declare(crewstate.Declaration{Meta: meta}).Closed() {
		d.logf("%s is closed already; the Mate landed it", crew)
		return true
	}
	out, err := d.Land(ctx, project, crew)
	if err != nil {
		d.logf("mate land %s %s: %v", project, crew, err)
		return false
	}
	left, next := LandReport(out)
	what := fmt.Sprintf("landed %s (%s merged on GitHub as %s)", crew, url, pr)
	if len(left) > 0 {
		what += "; left for you: " + strings.Join(left, "; ")
	}
	if next != "" {
		what += "; " + next
	}
	if err := d.WS.AppendAppNote(project, store.AppNote{Time: d.now(), Wake: true, Text: what}); err != nil {
		d.logf("app note: %v", err)
		return false
	}
	d.logf("%s", what)
	return true
}

// LandReport splits what `mate land` printed into the lines that left
// something for the Mate and its closing `next:` line.
func LandReport(out string) (left []string, next string) {
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
