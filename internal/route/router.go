package route

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/autopilot"
	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/prwatch"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// The states a hand-back's record in `mate/route.jsonl` moves through.
const (
	// StateShadow records what the router would have done, in shadow
	// mode. It takes nothing from the Mate.
	StateShadow = "shadow"
	// StateMate records that the router left the hand-back to the Mate.
	StateMate = "mate"
	// StateClaimed takes the hand-back from the Mate; StateDone and
	// StateFailed record what came of the action.
	StateClaimed = "claimed"
	StateDone    = "done"
	StateFailed  = "failed"
)

// ClaimExpiry is how long a claim stands without a done or failed record
// after it. A core that stopped between claim and action leaves one; once
// it expires the hand-back goes to the Mate.
const ClaimExpiry = 15 * time.Minute

// Record is one line of `mate/route.jsonl`.
type Record struct {
	At time.Time `json:"at"`
	// Ref is the hand-back's status line (autopilot.Item.Ref).
	Ref   string `json:"ref"`
	State string `json:"state"`
	Decision
	// Out is what the action printed, or why it failed.
	Out string `json:"out,omitempty"`
}

// Exec runs the actions, as the commands do (cmd/mate wires them).
type Exec struct {
	// Land lands crew, after the passing review of reviewer when it is
	// not empty, and returns what `mate land` printed.
	Land func(ctx context.Context, project, crew, reviewer string) (string, error)
	// Send types text into crew's pane, as `mate send` does.
	Send func(ctx context.Context, project, crew, text string) error
}

// Router is the autopilot's Router over one workspace.
type Router struct {
	ws    *store.Workspace
	judge Judge
	exec  Exec
	now   func() time.Time

	mu sync.Mutex
	// plans is what Plan decided, per project and Ref, so Jev is asked
	// once per hand-back.
	plans map[string]map[string]Decision
}

var _ autopilot.Router = (*Router)(nil)

// New builds a Router. judge may be nil: nothing Jev would be asked about
// is then taken from the Mate. now nil is time.Now.
func New(ws *store.Workspace, judge Judge, exec Exec, now func() time.Time) *Router {
	if now == nil {
		now = time.Now
	}
	return &Router{ws: ws, judge: judge, exec: exec, now: now, plans: map[string]map[string]Decision{}}
}

// mode is the project's autoroute setting; a project.yaml that does not
// read is off.
func (r *Router) mode(project string) string {
	cfg, err := r.ws.LoadProject(project)
	if err != nil {
		return store.AutorouteOff
	}
	return cfg.AutorouteMode()
}

// Plan decides about every hand-back among items it has not decided yet.
func (r *Router) Plan(ctx context.Context, project string, items []autopilot.Item) {
	if r.mode(project) == store.AutorouteOff {
		return
	}
	cfg, err := r.ws.LoadProject(project)
	if err != nil {
		return
	}
	for _, item := range items {
		if item.Kind != autopilot.ItemWaitMate {
			continue
		}
		r.mu.Lock()
		_, done := r.plans[project][item.Ref()]
		r.mu.Unlock()
		if done {
			continue
		}
		d := Decide(ctx, r.inspect(project, item.Crew, cfg), r.judge)
		if ctx.Err() != nil {
			return
		}
		r.mu.Lock()
		if r.plans[project] == nil {
			r.plans[project] = map[string]Decision{}
		}
		r.plans[project][item.Ref()] = d
		r.mu.Unlock()
	}
}

// inspect reads what Decide needs about one crew.
func (r *Router) inspect(project, crew string, cfg store.ProjectConfig) Input {
	in := Input{Crew: crew, Kind: KindShip, Review: cfg.Review}
	meta, _ := r.ws.ReadCrewMeta(project, crew)
	in.PRMerged = meta[crewstate.MetaPRState] == crewstate.PRStateMerged
	in.Reviewed = strings.TrimSpace(meta[store.MetaReviews]) != ""
	if rec, err := store.ReadMeta(r.ws.CrewReviewMeta(project, crew)); err == nil && strings.TrimSpace(rec["crew"]) != "" {
		in.Kind, in.ReviewOf = KindReview, strings.TrimSpace(rec["crew"])
		in.ReportPath = r.ws.CrewReport(project, crew)
		in.Text = read(in.ReportPath)
		return in
	}
	if r.ws.CrewIsScout(project, crew, meta) {
		in.Kind = KindScout
		in.Text = read(r.ws.CrewReport(project, crew))
		return in
	}
	in.Text = read(r.ws.CrewHandback(project, crew))
	return in
}

func read(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

// Claim takes out of items every hand-back the app acts on now, or acted on
// already, and records each decision once. In shadow mode it records and
// takes nothing.
func (r *Router) Claim(project string, items []autopilot.Item) (kept, claimed []autopilot.Item) {
	mode := r.mode(project)
	if mode == store.AutorouteOff {
		return items, nil
	}
	states := r.states(project)
	r.mu.Lock()
	plans := r.plans[project]
	r.mu.Unlock()
	for _, item := range items {
		ref := item.Ref()
		if item.Kind != autopilot.ItemWaitMate {
			kept = append(kept, item)
			continue
		}
		switch states[ref] {
		case StateClaimed, StateDone:
			// The app has it, or had it: the Mate hears from the note.
			continue
		case StateFailed, StateMate, StateShadow:
			kept = append(kept, item)
			continue
		}
		d, ok := plans[ref]
		if !ok {
			kept = append(kept, item)
			continue
		}
		switch {
		case mode == store.AutorouteShadow:
			r.record(project, Record{Ref: ref, State: StateShadow, Decision: d})
			kept = append(kept, item)
		case d.Action == ActionMate:
			r.record(project, Record{Ref: ref, State: StateMate, Decision: d})
			kept = append(kept, item)
		default:
			r.record(project, Record{Ref: ref, State: StateClaimed, Decision: d})
			claimed = append(claimed, item)
		}
	}
	return kept, claimed
}

// Execute acts on each claimed hand-back and records what came of it. What
// the app did goes to the Mate as a note in its next digest; an action that
// failed leaves the hand-back to the Mate, at the next tick.
func (r *Router) Execute(ctx context.Context, project string, claimed []autopilot.Item) {
	r.mu.Lock()
	plans := r.plans[project]
	r.mu.Unlock()
	for _, item := range claimed {
		d, ok := plans[item.Ref()]
		if !ok {
			continue
		}
		out, note, err := r.act(ctx, project, d)
		rec := Record{Ref: item.Ref(), State: StateDone, Decision: d, Out: out}
		if err != nil {
			rec.State, rec.Out = StateFailed, strings.TrimSpace(out+"\n"+err.Error())
		}
		r.record(project, rec)
		if err == nil {
			_ = r.ws.AppendAppNote(project, note)
		}
	}
}

// act runs one decision and returns what it printed and the note that tells
// the Mate.
func (r *Router) act(ctx context.Context, project string, d Decision) (string, store.AppNote, error) {
	note := store.AppNote{Time: r.now()}
	switch d.Action {
	case ActionLand:
		if r.exec.Land == nil {
			return "", note, fmt.Errorf("no land action wired")
		}
		out, err := r.exec.Land(ctx, project, d.Target, d.Reviewer)
		if err != nil {
			return out, note, err
		}
		left, next := prwatch.LandReport(out)
		text := fmt.Sprintf("landed %s from its hand-back (%s)", d.Target, d.Reason)
		if len(left) > 0 {
			text += "; left for you: " + strings.Join(left, "; ")
		}
		if next != "" {
			text += "; " + next
		}
		note.Wake, note.Text = true, text
		return out, note, nil
	case ActionSend:
		if r.exec.Send == nil {
			return "", note, fmt.Errorf("no send action wired")
		}
		if err := r.exec.Send(ctx, project, d.Target, d.Text); err != nil {
			return "", note, err
		}
		note.Text = fmt.Sprintf("sent %s the findings of review %s to fix (%s)", d.Target, d.Crew, d.Reason)
		return "", note, nil
	}
	return "", note, fmt.Errorf("nothing to do for action %q", d.Action)
}

// states is the latest state of every hand-back the log records. A claim
// older than ClaimExpiry with nothing after it counts as failed.
func (r *Router) states(project string) map[string]string {
	out := map[string]string{}
	at := map[string]time.Time{}
	for _, l := range strings.Split(read(r.ws.RouteLog(project)), "\n") {
		var rec Record
		if strings.TrimSpace(l) == "" || json.Unmarshal([]byte(l), &rec) != nil {
			continue
		}
		out[rec.Ref], at[rec.Ref] = rec.State, rec.At
	}
	now := r.now()
	for ref, state := range out {
		if state == StateClaimed && now.Sub(at[ref]) > ClaimExpiry {
			out[ref] = StateFailed
		}
	}
	return out
}

// record appends one line to `mate/route.jsonl`. The core is its one writer,
// one tick at a time, so a plain append is enough.
func (r *Router) record(project string, rec Record) {
	rec.At = r.now()
	b, err := json.Marshal(rec)
	if err != nil {
		return
	}
	f, err := os.OpenFile(r.ws.RouteLog(project), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(b, '\n'))
}
