package hook

import "github.com/nguyenngocanh94/mate/internal/store"

// Workspace is the part of *store.Workspace that the hook handlers use: the
// project's sent.log, which they append to, and the Mate's meta, which they
// read and rewrite (CON-241, plan docs/plans/store-split-2026-11-06.md
// section 3).
type Workspace interface {
	AppendSent(project string, entry store.SentEntry) error
	ReadMateMeta(project string) (map[string]string, error)
	WriteMateMeta(project string, meta map[string]string) error
}

var _ Workspace = (*store.Workspace)(nil)
