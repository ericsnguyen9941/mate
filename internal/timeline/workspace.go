package timeline

import (
	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Workspace is the part of *store.Workspace that an Ingester reads: what
// box.Load reads (the budget, context and no-turn rules pass the workspace
// to it), plus workspace.yaml, pricing, each project.yaml, the Mates' and
// crews' files and the review links. AppendIncident is its one write to the
// workspace (the budget, context-high and no-turn incidents); LoadConfig
// re-reads workspace.yaml into the handle's memory. Everything else only
// reads (CON-241, plan docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	box.Workspace
	AppendIncident(project string, entry store.IncidentEntry) error
	AutoFlag(project string) string
	CrewBrief(project, crew string) string
	CrewContextHigh() int
	LoadConfig() error
	LoadPricing() (store.PricingConfig, error)
	LoadProject(name string) (store.ProjectConfig, error)
	ManualHold(project string) string
	MateDir(project string) string
	MateMeta(project string) string
	MateSessionArchives(project string) ([]map[string]string, error)
	Projects() []store.ProjectRef
	ReadCrewHarnessProfiles(project, crew string) ([]store.HarnessProfile, error)
	ReadMateMeta(project string) (map[string]string, error)
	RepoDir(repo string) string
	ReviewLinks(project string) (reviewersOf map[string][]string, reviewOf map[string]string, err error)
	Root() string
}

var _ Workspace = (*store.Workspace)(nil)
