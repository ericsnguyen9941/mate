package outbox

import (
	"time"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// Workspace is the part of *store.Workspace that a Sender uses: the
// projects and their mode (LoadConfig re-reads workspace.yaml into the
// handle's memory), the Mate's outbox, meta and logs, the maintenance lock
// a delivery takes, and the TUI lock that says the Mate's session is open
// in its harness's TUI (TUIRunning only reads it). Its writes are
// UpdateOutbox, AppendSent, AppendIncident and WriteAutoCursor (CON-241,
// plan docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	AppendIncident(project string, entry store.IncidentEntry) error
	AppendSent(project string, entry store.SentEntry) error
	Auto(project string) bool
	LoadConfig() error
	MateDir(project string) string
	MateHeadless(project string) bool
	Projects() []store.ProjectRef
	ReadIncidents(project string, from int64) ([]store.IncidentEntry, int64, error)
	ReadMateMeta(project string) (map[string]string, error)
	ReadOutbox(project string) ([]store.OutboxItem, error)
	ReadSent(project string, from int64) ([]store.SentEntry, int64, error)
	SentLogSize(project string) (int64, error)
	TryMateMaintenance(project string, exclusive bool) (release func(), acquired bool, err error)
	TUIRunning() (store.TUIHolder, bool, error)
	UpdateOutbox(project string, now time.Time, change func(items []store.OutboxItem) ([]store.OutboxItem, bool, error)) error
	WriteAutoCursor(project string, cursor map[string]int64) error
}

var _ Workspace = (*store.Workspace)(nil)
