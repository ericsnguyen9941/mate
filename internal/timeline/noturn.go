package timeline

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// CheckNoTurn is the observer's ledger check for a crew that never got
// going (CON-74), run at the end of every poll after Ingest
// (watch.Deps.Turns): for every open crew that is still working - its last
// status verb is `working`, or it has written no status at all - launched
// more than after ago, whose actor has no turn in the ledger, it opens a
// `no_turn` incident once. On hellovietnam (2026-10-14) a Codex reviewer
// sat at a sign-in screen with 0 calls and no incident, and the Mate never
// knew.
//
// It resolves the incident as soon as the ledger shows a turn for the
// crew: the condition is "never started", and the first turn ends it. An
// incident still open on a crew that has since closed is resolved by the
// observer's own sweep of closed crews, like every other kind.
//
// A crew whose launch time is not recorded is skipped: there is no moment
// to count from. A crew that asked a question or handed back is skipped
// too, whatever the ledger says: it has spoken, and a ledger that has not
// caught up with it (a transcript not yet located) is the ingest's
// problem, not a crew that never started.
func (i *Ingester) CheckNoTurn(ctx context.Context, after time.Duration) error {
	if !i.d.Writable() {
		return errors.New("timeline: this database handle is read-only")
	}
	if err := i.ws.LoadConfig(); err != nil {
		return err
	}
	var errs []error
	for _, ref := range i.ws.Projects() {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := i.checkProjectNoTurn(ctx, ref.Name, after); err != nil {
			errs = append(errs, fmt.Errorf("timeline: no-turn check %s: %w", ref.Name, err))
		}
	}
	return errors.Join(errs...)
}

func (i *Ingester) checkProjectNoTurn(ctx context.Context, project string, after time.Duration) error {
	crews, err := listCrewMetas(i.ws, project)
	if err != nil {
		return err
	}
	if len(crews) == 0 {
		return nil
	}
	view, err := box.Load(i.ws, project)
	if err != nil {
		return err
	}
	now := i.deps.now()
	for _, crew := range crews {
		if view.Closed[crew.ID] || crew.state().Closed() {
			continue
		}
		open := false
		for _, inc := range box.OpenIncidents(view, crew.ID) {
			if inc.Kind == box.IncidentNoTurn {
				open = true
			}
		}
		turns, err := i.crewTurns(ctx, crew.ActorID)
		if err != nil {
			return err
		}
		if turns > 0 {
			if open {
				if err := i.ws.AppendIncident(project, store.IncidentEntry{
					Time: now, Crew: crew.ID, Kind: string(box.IncidentNoTurn),
					State: store.IncidentResolved, Text: fmt.Sprintf("the ledger now records %d turn(s)", turns),
				}); err != nil {
					return err
				}
			}
			continue
		}
		if open {
			continue
		}
		if verb := box.LastCrewVerb(view, crew.ID); verb != box.StateWorking && verb != box.StateUnknown {
			continue
		}
		launched := metaTime(crew.Meta, MetaLaunchedAt)
		if launched.IsZero() {
			// An older record dates only the spawn; launchTime's slack
			// is for adopting a rollout, not for counting a wait.
			launched = metaTime(crew.Meta, MetaStartedAt)
		}
		if launched.IsZero() {
			continue
		}
		quiet := now.Sub(launched)
		if quiet <= after {
			continue
		}
		if err := i.ws.AppendIncident(project, store.IncidentEntry{
			Time: now, Crew: crew.ID, Kind: string(box.IncidentNoTurn),
			State: store.IncidentOpen, Text: fmt.Sprintf("no turn in the ledger %s after launch", quiet.Round(time.Minute)),
		}); err != nil {
			return err
		}
	}
	return nil
}

// crewTurns is how many turns the ledger holds for one actor.
func (i *Ingester) crewTurns(ctx context.Context, actorID string) (int64, error) {
	var n int64
	err := i.d.SQL().QueryRowContext(ctx, `SELECT COUNT(*) FROM turn WHERE actor_id = ?`, actorID).Scan(&n)
	return n, err
}
