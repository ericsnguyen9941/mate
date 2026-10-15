package timeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// projectBudgetCrew is the crew name a project-wide budget crossing is filed
// under: there is no one crew to blame for the project's total, so it goes
// under the same name the daemon's own `wedged` incident already uses for a
// finding that is about the Mate rather than any one crew (mvp.md section
// 4b, internal/autopilot.MateCrew - not imported here to avoid a needless
// dependency between two leaf-ish packages that would otherwise never need
// each other).
const projectBudgetCrew = "mate"

// CheckBudgets is mvp.md M5 task 27's budget check, run by the observer at
// the end of every poll (watch.Deps.Budget): for every project, it opens a
// `budget` incident on whichever crew - or, for the project-wide ceiling,
// on projectBudgetCrew - has crossed a configured limit.
//
// A crew's token limit is its own `budget=` (the dispatch rule's budget
// `mate crew spawn --budget` recorded, CON-76) when it has one, else the
// project's `budget.crew_tokens`; `crew_usd` and `project_usd` are the
// project's alone. Against the token limit there are two thresholds: all
// of it, and half of it while the crew has not handed back (no `pr-open`
// or `wait-mate` line yet) - the early warning for a task that is spending
// its size away without delivering. The half incident carries
// box.BudgetHalfPrefix so every reader tells the two apart.
//
// It reads the ledger this same poll's Ingest just wrote, so a crossing
// becomes visible in incidents.log one poll after the turn that caused it -
// the same lag every other incident kind already has behind the transcript
// it is read from.
//
// The same call runs CON-73's context check (checkProjectContext): it is
// the other thing the ledger tells the observer about a crew, read from
// the same poll's turns.
//
// A budget incident is never resolved (mvp.md section 4b, decision
// 2026-09-20: "nó là một sự thật" - it is a fact about spend that already
// happened, not a condition that can clear). CheckBudgets therefore only
// ever appends one `open` line per (crew, threshold): a crew already
// carrying one for a threshold is left alone even if it has spent more
// since, because the record already says what needs saying; a crew
// carrying the full one is not checked again at all.
func (i *Ingester) CheckBudgets(ctx context.Context) error {
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
		if err := i.checkProjectBudget(ctx, ref.Name); err != nil {
			errs = append(errs, fmt.Errorf("timeline: budget check %s: %w", ref.Name, err))
		}
		if err := i.checkProjectContext(ctx, ref.Name); err != nil {
			errs = append(errs, fmt.Errorf("timeline: context check %s: %w", ref.Name, err))
		}
	}
	return errors.Join(errs...)
}

func (i *Ingester) checkProjectBudget(ctx context.Context, project string) error {
	cfg, err := i.ws.LoadProject(project)
	if err != nil {
		return err
	}
	crews, err := listCrewMetas(i.ws, project)
	if err != nil {
		return err
	}
	budget := cfg.Budget
	if budget == nil {
		budget = &store.BudgetConfig{}
	}
	if *budget == (store.BudgetConfig{}) && !anyCrewBudget(crews) {
		// Nothing to hold anyone to; most projects never configure one.
		return nil
	}

	view, err := box.Load(i.ws, project)
	if err != nil {
		return err
	}
	// filed is which thresholds the crew already carries an open budget
	// incident for: the record of a crossing is written once.
	filed := func(crew string) (half, full bool) {
		for _, inc := range box.OpenIncidents(view, crew) {
			if inc.Kind != box.IncidentBudget {
				continue
			}
			if percent, _ := box.BudgetThreshold(inc.Text); percent == 50 {
				half = true
			} else {
				full = true
			}
		}
		return half, full
	}

	for _, crew := range crews {
		half, full := filed(crew.ID)
		if full {
			continue
		}
		limit := crewBudget(crew.Meta, budget.CrewTokens)
		if limit == 0 && budget.CrewUSD == 0 {
			continue
		}
		tokens, cost, err := i.crewSpend(ctx, crew.ActorID)
		if err != nil {
			return err
		}
		text := ""
		switch {
		case limit > 0 && tokens > limit:
			text = fmt.Sprintf("%s tokens of %s tokens",
				query.HumanizeTokens(tokens), query.HumanizeTokens(limit))
		case budget.CrewUSD > 0 && cost.Valid && cost.Float64 > budget.CrewUSD:
			text = fmt.Sprintf("%s of %s", query.HumanizeCost(cost.Float64), query.HumanizeCost(budget.CrewUSD))
		case limit > 0 && !half && tokens*2 > limit && !view.Closed[crew.ID] && !handedBack(view, crew.ID):
			// Half the budget spent and nothing handed back yet. A closed
			// crew is past warning about, and one that has handed back
			// has the Mate's eyes on it already.
			text = box.BudgetHalfPrefix + fmt.Sprintf("%s tokens of %s tokens, no hand-back yet",
				query.HumanizeTokens(tokens), query.HumanizeTokens(limit))
		}
		if text == "" {
			continue
		}
		if err := i.ws.AppendIncident(project, store.IncidentEntry{
			Time: i.deps.now(), Crew: crew.ID, Kind: string(box.IncidentBudget),
			State: store.IncidentOpen, Text: text,
		}); err != nil {
			return err
		}
	}

	if _, full := filed(projectBudgetCrew); budget.ProjectUSD > 0 && !full {
		cost, err := i.projectSpend(ctx, project)
		if err != nil {
			return err
		}
		if cost.Valid && cost.Float64 > budget.ProjectUSD {
			text := fmt.Sprintf("%s of %s", query.HumanizeCost(cost.Float64), query.HumanizeCost(budget.ProjectUSD))
			if err := i.ws.AppendIncident(project, store.IncidentEntry{
				Time: i.deps.now(), Crew: projectBudgetCrew, Kind: string(box.IncidentBudget),
				State: store.IncidentOpen, Text: text,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// crewBudget is the token limit one crew is held to: its own `budget=`
// when the meta carries a positive one, else the project's crew_tokens
// (0 when neither). A meta value that is not a number is ignored rather
// than failing the poll: the record is the Mate's, and a bad one costs the
// crew its own limit, not every crew its check.
func crewBudget(meta map[string]string, projectTokens int64) int64 {
	if raw := strings.TrimSpace(meta[store.MetaBudget]); raw != "" {
		if n, err := strconv.ParseInt(raw, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return projectTokens
}

// anyCrewBudget reports whether any crew carries its own `budget=`.
func anyCrewBudget(crews []crewRecord) bool {
	for _, crew := range crews {
		if crewBudget(crew.Meta, 0) > 0 {
			return true
		}
	}
	return false
}

// handedBack reports whether the crew has ever written a hand-back line:
// `wait-mate`, `pr-open`, or any other verb box reads as StateWaitMate. It
// is "ever", not "latest": a crew sent back to work after handing back is
// one the Mate already looked at, so the half-budget warning has nothing
// to add.
func handedBack(view box.View, crew string) bool {
	for _, e := range view.ByCrew[crew] {
		if e.Kind == box.KindStatus && box.ParseStatus(e.Text).State == box.StateWaitMate {
			return true
		}
	}
	return false
}

// crewSpend is one crew's whole-task token total and cost, read the same way
// `v_task_ledger` computes them: cost only counts a turn whose model has a
// real (non-zero) price, so a model that is only listed in pricing.yaml as
// a placeholder never crosses a dollar budget.
func (i *Ingester) crewSpend(ctx context.Context, actorID string) (int64, sql.NullFloat64, error) {
	var tokens int64
	if err := i.d.SQL().QueryRowContext(ctx, `
		SELECT COALESCE(SUM(input_tokens+cache_read_tokens+cache_write_tokens+output_tokens), 0)
		  FROM turn WHERE actor_id = ?`, actorID).Scan(&tokens); err != nil {
		return 0, sql.NullFloat64{}, err
	}
	var cost sql.NullFloat64
	if err := i.d.SQL().QueryRowContext(ctx, `
		SELECT SUM(
			u.input_tokens       * COALESCE(p.input_per_m, 0) / 1000000.0 +
			u.cache_read_tokens  * COALESCE(p.cache_read_per_m, 0) / 1000000.0 +
			u.cache_write_tokens * COALESCE(p.cache_write_per_m, 0) / 1000000.0 +
			u.output_tokens      * COALESCE(p.output_per_m, 0) / 1000000.0)
		  FROM turn u JOIN pricing p ON p.model = u.model
		 WHERE u.actor_id = ?
		   AND (p.input_per_m > 0 OR p.cache_read_per_m > 0 OR p.cache_write_per_m > 0 OR p.output_per_m > 0)`,
		actorID).Scan(&cost); err != nil {
		return 0, sql.NullFloat64{}, err
	}
	return tokens, cost, nil
}

// projectSpend is every crew of a project summed: the cost a `project_usd`
// ceiling is checked against.
func (i *Ingester) projectSpend(ctx context.Context, project string) (sql.NullFloat64, error) {
	var cost sql.NullFloat64
	row := i.d.SQL().QueryRowContext(ctx, `
		SELECT SUM(
			u.input_tokens       * COALESCE(p.input_per_m, 0) / 1000000.0 +
			u.cache_read_tokens  * COALESCE(p.cache_read_per_m, 0) / 1000000.0 +
			u.cache_write_tokens * COALESCE(p.cache_write_per_m, 0) / 1000000.0 +
			u.output_tokens      * COALESCE(p.output_per_m, 0) / 1000000.0)
		FROM turn u
		JOIN actor a ON a.id = u.actor_id
		JOIN pricing p ON p.model = u.model
		WHERE a.project = ? AND a.kind = ?
		  AND (p.input_per_m > 0 OR p.cache_read_per_m > 0 OR p.cache_write_per_m > 0 OR p.output_per_m > 0)`,
		project, ActorCrew)
	if err := row.Scan(&cost); err != nil {
		return sql.NullFloat64{}, err
	}
	return cost, nil
}
