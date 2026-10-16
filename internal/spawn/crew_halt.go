package spawn

import (
	"context"
	"fmt"

	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/observability"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// HaltCrew stops a crew's agent and closes its tab without closing its
// task: the worktree, the branch and every other meta key stay, and the
// meta stops naming an agent, which is the shape every reader already
// knows as "this crew is not running" (CrewHandle's refusal). It is the
// daemon's answer to a crew at its budget's hard ceiling whose work is not
// landed (CON-94): StopCrew would refuse, and throwing the work away is
// the Mate's or the captain's call, not the app's. `mate crew relaunch`
// starts it again; `mate crew stop` closes it.
//
// It reports whether a live agent was stopped. A closed crew is left
// alone and reports false.
func HaltCrew(ctx context.Context, w *store.Workspace, deps Deps, project, crew string) (bool, error) {
	if w == nil {
		return false, errUsage("spawn: a workspace is required")
	}
	if deps.Runtime == nil {
		return false, errUsage("spawn: a runtime adapter is required")
	}
	if deps.Names == nil {
		deps.Names = runtime.NewMemoryNameRegistry()
	}
	if err := store.ValidateProjectName(project); err != nil {
		return false, err
	}
	if err := store.ValidateCrewID(crew); err != nil {
		return false, err
	}
	// The stop lock is held across the agent stopping and the meta write,
	// for the reason StopCrew holds it: the observer reads it and does not
	// call the missing agent `runtime_lost`.
	unlock, err := w.LockCrewStop(ctx, project, crew)
	if err != nil {
		return false, err
	}
	defer unlock()
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil {
		return false, err
	}
	if len(meta) == 0 {
		return false, observability.NewError(observability.CodeNotFound,
			fmt.Sprintf("no crew %s is recorded for project %s; %s is empty", crew, project, w.CrewMeta(project, crew)))
	}
	if crewstate.State(meta[MetaState]).Closed() {
		return false, nil
	}
	stopped, err := stopCrewAgentForRelaunch(ctx, w, deps, project, crew, meta)
	if err != nil {
		return false, err
	}
	next := make(map[string]string, len(meta))
	for k, v := range meta {
		next[k] = v
	}
	// Dropped together, as StopCrew drops them: one without the other is
	// a half-written record.
	delete(next, MetaAgent)
	delete(next, MetaPane)
	delete(next, MetaTab)
	return stopped, w.WriteCrewMeta(project, crew, next)
}
