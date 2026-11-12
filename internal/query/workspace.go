package query

import (
	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// BoxWorkspace is the part of *store.Workspace that LoadBox reads: what
// box.Load reads, plus the project directory an assign key is relative to,
// the Mate's outbox and the autopilot's digest cursor, and their paths,
// which the memo stats (memo.go). Every method only reads (CON-238, plan
// docs/plans/store-split-2026-11-06.md section 3).
type BoxWorkspace interface {
	box.Workspace
	AutoCursorFile(project string) string
	OutboxFile(project string) string
	ProjectDir(project string) string
	ReadAutoCursor(project string) (map[string]int64, error)
	ReadOutbox(project string) ([]store.OutboxItem, error)
}

// Workspace is the part of *store.Workspace that Load and LoadLive read:
// the box of every project (BoxWorkspace), plus workspace.yaml, each
// project.yaml and each Mate's files. LoadConfig re-reads workspace.yaml
// into the handle's memory and writes nothing to disk; every other method
// only reads.
type Workspace interface {
	BoxWorkspace
	Defaults() store.Defaults
	LoadConfig() error
	LoadProject(name string) (store.ProjectConfig, error)
	MateDir(project string) string
	MateMeta(project string) string
	Projects() []store.ProjectRef
	ReadMateMeta(project string) (map[string]string, error)
	Root() string
}

var (
	_ BoxWorkspace = (*store.Workspace)(nil)
	_ Workspace    = (*store.Workspace)(nil)
)

// noWorkspace reports whether ws holds no workspace at all: a nil
// interface, or a nil *store.Workspace inside one. A caller holding a nil
// *store.Workspace passes the second, and `ws == nil` alone does not see it
// (plan section 3, rule 4). A Workspace is a BoxWorkspace, so both kinds go
// through here.
func noWorkspace(ws BoxWorkspace) bool {
	if ws == nil {
		return true
	}
	w, ok := ws.(*store.Workspace)
	return ok && w == nil
}
