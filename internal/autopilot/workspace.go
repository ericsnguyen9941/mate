package autopilot

import (
	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Workspace is the part of *store.Workspace that a Pilot uses: what
// box.Load reads (the tick and the budget rule pass the workspace to it),
// plus workspace.yaml (LoadConfig re-reads it into the handle's memory),
// the manual hold, the digest cursor and the Mate's outbox. Its writes are
// SetMode and AppendSent (rearming Auto after the manual hold lapses) and
// UpdateCrewMeta (the budget stamps) (CON-241, plan
// docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	box.Workspace
	AppendSent(project string, entry store.SentEntry) error
	LoadConfig() error
	ManualHold(project string) string
	Projects() []store.ProjectRef
	ReadAutoCursor(project string) (map[string]int64, error)
	ReadOutbox(project string) ([]store.OutboxItem, error)
	SetMode(project string, auto bool) error
	UpdateCrewMeta(project, crew string, set map[string]string) error
}

var _ Workspace = (*store.Workspace)(nil)
