package autopilot

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// BudgetReachedLine is the line the app types into a crew's pane, after
// the from-app sentinel, once its task has spent its whole budget (CON-94).
// It asks for a clean stop rather than taking one: the crew commits what it
// has and hands back, and the Mate decides from that hand-back whether to
// extend, narrow or drop the task (skill `wake`, its Budget section).
const BudgetReachedLine = "budget reached: commit what you have, write the hand-back with what is done and what is not, stop with wait-mate"

// metaAgent is the crew meta key naming its live agent (spawn.MetaAgent;
// this package does not import internal/spawn). A crew whose meta names no
// agent is not running, and nothing can be typed into it.
const metaAgent = "agent"

// BudgetActions is how the daemon acts on a crew's budget (CON-94). It is a
// seam, not an import of internal/spawn, for the reason Outbox is: this
// package holds no runtime. cmd/mate wires both.
type BudgetActions struct {
	// Line types text into the crew's pane, with the from-app sentinel,
	// and returns once the send is verified. A crew mid-turn takes it into
	// the harness's own queue: a crew past its budget is almost always
	// mid-turn, and the line is for its next step.
	Line func(ctx context.Context, project, crew, text string) error
	// Stop ends a crew at its hard ceiling, without asking anyone, and
	// says in a few words what it did (a full stop, or the agent stopped
	// with the unlanded work kept).
	Stop func(ctx context.Context, project, crew string) (string, error)
}

// enforceBudgets is one tick's budget enforcement for one project, run in
// manual and auto mode alike: a budget is a stop condition the app counts,
// not something a Mate or a captain has to be present for.
//
// It reads the observer's open `budget` incidents
// (internal/timeline.CheckBudgets files them; this package never counts
// tokens) and, for every open crew:
//
//   - at 100% types BudgetReachedLine into the crew's pane once
//     (`budget_line=` in its meta), unless the crew has already handed
//     back: a resting crew has already stopped at a clean point, and gets
//     the line if a fix round starts it again;
//   - at the hard ceiling (twice the budget) stops it once
//     (`budget_stop=`). The incident already says why; the digest carries
//     it to the Mate like every other budget item.
//
// The 50% incident is the Mate's alone (CON-76): nothing is typed.
// An extension (`mate crew budget --extend`) resolves the incidents and
// deletes both meta keys, so each threshold acts again at the new ceiling.
func (p *Pilot) enforceBudgets(ctx context.Context, project string, now time.Time) error {
	actions := p.deps.Budget
	if actions == nil {
		return nil
	}
	view, err := box.Load(p.ws, project)
	if err != nil {
		return err
	}
	reached := make(map[string]int)
	for _, e := range view.Entries {
		if e.Kind != box.KindIncident || e.Incident == nil || e.Incident.Resolved || e.Incident.Kind != box.IncidentBudget {
			continue
		}
		if percent, _ := box.BudgetThreshold(e.Incident.Text); percent > reached[e.Crew] {
			reached[e.Crew] = percent
		}
	}
	crews := make([]string, 0, len(reached))
	for crew, percent := range reached {
		if percent >= 100 && crew != "" && crew != MateCrew && !view.Closed[crew] {
			crews = append(crews, crew)
		}
	}
	sort.Strings(crews)

	var errs []error
	for _, crew := range crews {
		if err := ctx.Err(); err != nil {
			return err
		}
		meta, err := p.ws.ReadCrewMeta(project, crew)
		if err != nil || len(meta) == 0 || crewstate.Declare(crewstate.Declaration{Meta: meta}).Closed() {
			continue
		}
		stamp := now.UTC().Format(time.RFC3339)
		if reached[crew] >= 200 {
			if meta[store.MetaBudgetStop] != "" || actions.Stop == nil {
				continue
			}
			did, err := actions.Stop(ctx, project, crew)
			if err != nil {
				errs = append(errs, fmt.Errorf("budget hard ceiling: stop %s: %w", crew, err))
				continue
			}
			if err := p.ws.UpdateCrewMeta(project, crew, map[string]string{store.MetaBudgetStop: stamp + " " + did}); err != nil {
				errs = append(errs, err)
			}
			continue
		}
		if meta[store.MetaBudgetLine] != "" || actions.Line == nil || meta[metaAgent] == "" || resting(view, crew) {
			continue
		}
		if err := actions.Line(ctx, project, crew, BudgetReachedLine); err != nil {
			// A refused send is retried next tick; the line is recorded
			// only once it went through.
			errs = append(errs, fmt.Errorf("budget reached: tell %s: %w", crew, err))
			continue
		}
		if err := p.ws.UpdateCrewMeta(project, crew, map[string]string{store.MetaBudgetLine: stamp}); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// resting reports whether the crew's latest status line is a hand-back
// (`wait-mate`, `pr-open` and the other verbs box reads as StateWaitMate).
func resting(view box.View, crew string) bool {
	entries := view.ByCrew[crew]
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Kind == box.KindStatus {
			return box.ParseStatus(entries[i].Text).State == box.StateWaitMate
		}
	}
	return false
}
