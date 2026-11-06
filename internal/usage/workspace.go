package usage

import (
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Workspace is the part of *store.Workspace that Write reads: what the
// ledger's spawn.ListCrews reads (spawn.CrewReader), plus the workspace's
// autocompact threshold, each project.yaml (for the budgets), the Mate's
// metadata and who reviews whom. Every method only reads (CON-240, plan
// docs/plans/store-split-2026-11-06.md section 3). `--why` and the
// learning record's usage section, which read through dashboard.Server,
// are in internal/usage/why (CON-234).
type Workspace interface {
	spawn.CrewReader
	CrewAutocompact() int
	LoadProject(name string) (store.ProjectConfig, error)
	ReadMateMeta(project string) (map[string]string, error)
	ReviewLinks(project string) (reviewersOf map[string][]string, reviewOf map[string]string, err error)
}

var _ Workspace = (*store.Workspace)(nil)
