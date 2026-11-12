package watch

import "github.com/nguyenngocanh94/mate/internal/store"

// Workspace is the part of *store.Workspace that a Watcher uses: the
// projects in workspace.yaml (LoadConfig re-reads it into the handle's
// memory), the Mate's mate.meta (whether it is live, CON-267), each crew's
// metadata, status file and stop mark, and the project's incidents.log,
// which AppendIncident is its one write to
// (CON-241, plan docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	AppendIncident(project string, entry store.IncidentEntry) error
	CrewStopping(project, crew string) (bool, error)
	CrewsDir(project string) string
	LoadConfig() error
	Projects() []store.ProjectRef
	ReadCrewMeta(project, crew string) (map[string]string, error)
	ReadIncidents(project string, from int64) ([]store.IncidentEntry, int64, error)
	ReadMateMeta(project string) (map[string]string, error)
	ReadStatus(project, crew string, from int64) ([]store.StatusEntry, int64, error)
}

var _ Workspace = (*store.Workspace)(nil)
