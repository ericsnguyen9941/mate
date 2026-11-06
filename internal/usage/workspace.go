package usage

import (
	"github.com/nguyenngocanh94/mate/internal/dashboard"
	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Workspace is the part of *store.Workspace that Write and RecordAtClose
// read: what `--why`'s dashboard.Server reads (dashboard.Workspace), what
// the ledger's spawn.ListCrews reads (spawn.CrewReader), where db.OpenRead
// finds `mate.db` (db.Workspace), plus the workspace's autocompact
// threshold and the path of a crew's learning record. Every method only
// reads. RecordAtClose writes the record's usage section itself, not
// through the store, as section 7 of the plan notes (CON-240, plan
// docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	dashboard.Workspace
	spawn.CrewReader
	db.Workspace
	CrewAutocompact() int
	CrewLearning(project, crew string) string
}

var _ Workspace = (*store.Workspace)(nil)
