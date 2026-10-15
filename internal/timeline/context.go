package timeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// checkProjectContext is CON-73's context check, run by CheckBudgets beside
// the budget one: for every open crew of a project whose latest status is
// not a hand-back, it files a `context_high` incident once when the
// crew's latest `context_tokens_after` is above the workspace's
// `crew.context_high`, and resolves the incident once that context is back
// under it (a fresh session) or the crew is closed.
//
// The incident never cuts a crew: a fix round in a 300k context is a fact
// about spend to come, and the Mate decides what to do with it at the
// crew's next hand-back (`mate wake` says so). A crew that has already
// handed back is left alone - the Mate's `mate send` opens its fresh
// session itself - and a closed crew carries no open incident.
func (i *Ingester) checkProjectContext(ctx context.Context, project string) error {
	threshold := int64(i.ws.CrewContextHigh())
	if threshold <= 0 {
		return nil
	}
	crews, err := listCrewMetas(i.ws, project)
	if err != nil {
		return err
	}
	view, err := box.Load(i.ws, project)
	if err != nil {
		return err
	}
	var errs []error
	for _, crew := range crews {
		open := false
		for _, inc := range box.OpenIncidents(view, crew.ID) {
			if inc.Kind == box.IncidentContextHigh {
				open = true
			}
		}
		latest, err := i.latestContext(ctx, crew.ActorID)
		if err != nil {
			return err
		}
		closed := view.Closed[crew.ID]
		resting := false
		if st, ok := box.LatestStatus(view, crew.ID); ok && st.State == box.StateWaitMate {
			resting = true
		}
		switch {
		case open && (closed || (latest.Valid && latest.Int64 <= threshold)):
			why := "the crew is closed"
			if !closed {
				why = fmt.Sprintf("its latest context is %s tokens, under %s",
					query.HumanizeTokens(latest.Int64), query.HumanizeTokens(threshold))
			}
			errs = append(errs, i.ws.AppendIncident(project, store.IncidentEntry{
				Time: i.deps.now(), Crew: crew.ID, Kind: string(box.IncidentContextHigh),
				State: store.IncidentResolved, Text: why,
			}))
		case !open && !closed && !resting && latest.Valid && latest.Int64 > threshold:
			errs = append(errs, i.ws.AppendIncident(project, store.IncidentEntry{
				Time: i.deps.now(), Crew: crew.ID, Kind: string(box.IncidentContextHigh),
				State: store.IncidentOpen, Text: ContextHighText(latest.Int64, threshold),
			}))
		}
	}
	return errors.Join(errs...)
}

// ContextHighText is the text of a `context_high` incident: the crew's
// latest context and the threshold it passed, humanised the way the
// ledger prints them.
func ContextHighText(latest, threshold int64) string {
	return fmt.Sprintf("%s tokens of context, over %s", query.HumanizeTokens(latest), query.HumanizeTokens(threshold))
}

// latestContext is the context size after the actor's most recent turn
// that recorded one; invalid when none did.
func (i *Ingester) latestContext(ctx context.Context, actorID string) (sql.NullInt64, error) {
	var latest sql.NullInt64
	err := i.d.SQL().QueryRowContext(ctx, `
		SELECT context_tokens_after FROM turn
		 WHERE actor_id = ? AND context_tokens_after > 0
		 ORDER BY started_at DESC, ordinal DESC LIMIT 1`, actorID).Scan(&latest)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.NullInt64{}, nil
	}
	return latest, err
}
