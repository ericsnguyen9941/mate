package spawn

import (
	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// CrewReader is the part of *store.Workspace that ListCrews reads: the
// project's box (box.Workspace, for box.Load), plus whether the project is
// registered and which crews have a `.meta`. Every method only reads. Only
// ListCrews and the crewIDs it calls take it; the rest of spawn writes and
// keeps *store.Workspace (CON-240, plan docs/plans/store-split-2026-11-06.md
// section 3).
type CrewReader interface {
	box.Workspace
	CrewIDs(project string) ([]string, error)
	Project(name string) (store.ProjectRef, bool)
}

var _ CrewReader = (*store.Workspace)(nil)

// noWorkspace reports whether w holds no workspace at all: a nil
// interface, or a nil *store.Workspace inside one, which `w == nil` alone
// does not see (plan section 3, rule 4).
func noWorkspace(w CrewReader) bool {
	if w == nil {
		return true
	}
	ws, ok := w.(*store.Workspace)
	return ok && ws == nil
}
