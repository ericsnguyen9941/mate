package timeline

import (
	"context"

	"github.com/nguyenngocanh94/mate/internal/timeline/scene"
)

// projectScene runs the scene projection of task 26 over the project this
// pass has just written, inside the same transaction: the `transition` rows
// and the events they are derived from commit together, so a reader never
// sees a story the scene has not caught up with.
//
// It runs last because the projection reads `cause_event_id` - a Mate turn is
// `reading(crew)` because the note that started it names that crew - and the
// causality pass fills those links.
//
// A pass that inserted no event moved nobody, so the projection is skipped;
// the exception is a project that has never been projected at all, which is
// what fills `transition` for a database written before this existed.
//
// Whether it has been is asked of the table once per project in the life of
// the Ingester (memos.projected, CON-268): `transition` has no index on
// project, and the COUNT(*) this was scanned the whole table for every
// project on every pass - a quarter of a steady pass's CPU on the
// captain's workspace (2026-11-12).
func (p *pass) projectScene(ctx context.Context) error {
	if p.w.insertedEvents == 0 {
		if p.ing.memos.projected[p.project] {
			p.projected = true
			return nil
		}
		var rows int
		if err := p.tx.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM transition WHERE project = ?)`, p.project).Scan(&rows); err != nil {
			return err
		}
		if rows > 0 {
			p.projected = true
			return nil
		}
	}
	if err := scene.Run(ctx, p.tx, p.project); err != nil {
		return err
	}
	p.projected = true
	return nil
}
