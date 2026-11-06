package box

import "github.com/nguyenngocanh94/mate/internal/store"

// Workspace is the part of *store.Workspace that Load and LoadSince read:
// the project's mode, its sent.log, its incidents.log and its crews' status
// files and metadata. Every method only reads. A package that passes its
// workspace to Load embeds this interface in its own (CON-238, plan
// docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	Auto(project string) bool
	CrewStatus(project, crew string) string
	CrewsDir(project string) string
	IncidentsLog(project string) string
	ReadCrewMeta(project, crew string) (map[string]string, error)
	ReadIncidents(project string, from int64) ([]store.IncidentEntry, int64, error)
	ReadSent(project string, from int64) ([]store.SentEntry, int64, error)
	ReadStatus(project, crew string, from int64) ([]store.StatusEntry, int64, error)
	SentLog(project string) string
}

var _ Workspace = (*store.Workspace)(nil)
