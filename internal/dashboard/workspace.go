package dashboard

import (
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Workspace is the part of *store.Workspace that the server reads: each
// project's box (query.BoxWorkspace, for query.LoadBox), plus the project
// list, project.yaml, the Mate's files and the paths of a crew's brief,
// hand-back, report and review record. Every method only reads, which is
// the store side of the dashboard's claim that it cannot change the
// workspace. It has no LoadConfig: re-reading workspace.yaml is the
// engine's job (CON-238, plan docs/plans/store-split-2026-11-06.md
// section 3).
type Workspace interface {
	query.BoxWorkspace
	CrewBrief(project, crew string) string
	CrewHandback(project, crew string) string
	CrewIsScout(project, crew string, meta map[string]string) bool
	CrewReport(project, crew string) string
	CrewReviewMeta(project, crew string) string
	Defaults() store.Defaults
	LoadProject(name string) (store.ProjectConfig, error)
	MateDir(project string) string
	MateHeadless(project string) bool
	Project(name string) (store.ProjectRef, bool)
	Projects() []store.ProjectRef
	ReadMateMeta(project string) (map[string]string, error)
	ReviewLinks(project string) (reviewersOf map[string][]string, reviewOf map[string]string, err error)
	Root() string
}

var _ Workspace = (*store.Workspace)(nil)

// noWorkspace reports whether ws holds no workspace at all: a nil
// interface, or a nil *store.Workspace inside one, which `ws == nil` alone
// does not see (plan section 3, rule 4).
func noWorkspace(ws Workspace) bool {
	if ws == nil {
		return true
	}
	w, ok := ws.(*store.Workspace)
	return ok && w == nil
}
