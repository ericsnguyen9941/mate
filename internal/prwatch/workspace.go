package prwatch

import "github.com/nguyenngocanh94/mate/internal/store"

// Workspace is the part of *store.Workspace that the pull-request watcher
// and its sweep use: the projects and each project.yaml, each crew's
// metadata, status file and `.prwatch` paths, the repo directory and the
// Mate's outbox. Its writes through the store are AppendStatus and
// UpdateCrewMeta. This package writes and removes the `.prwatch` pid file
// itself, at the path CrewPRWatch gives, and its Starter opens the
// watcher's log at CrewPRWatchLog (CON-241, plan
// docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	AppendStatus(project, crew, line string) error
	CrewIDs(project string) ([]string, error)
	CrewPRWatch(project, crew string) string
	CrewPRWatchLog(project, crew string) string
	CrewStatus(project, crew string) string
	CrewsDir(project string) string
	LoadProject(name string) (store.ProjectConfig, error)
	Projects() []store.ProjectRef
	ReadCrewMeta(project, crew string) (map[string]string, error)
	ReadOutbox(project string) ([]store.OutboxItem, error)
	ReadStatus(project, crew string, from int64) ([]store.StatusEntry, int64, error)
	RepoDir(repo string) string
	Root() string
	UpdateCrewMeta(project, crew string, set map[string]string) error
}

var _ Workspace = (*store.Workspace)(nil)

// noWorkspace reports whether ws holds no workspace at all: a nil
// interface, or a nil *store.Workspace inside one. A caller holding a nil
// *store.Workspace passes the second, and `ws == nil` alone does not see it
// (plan section 3, rule 4).
func noWorkspace(ws Workspace) bool {
	if ws == nil {
		return true
	}
	w, ok := ws.(*store.Workspace)
	return ok && w == nil
}
