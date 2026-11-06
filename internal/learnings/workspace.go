package learnings

import "github.com/nguyenngocanh94/mate/internal/store"

// Workspace is the part of *store.Workspace that Open, List, Read and Close
// read: where a project's learning records are. Both methods only name
// paths; the records themselves are this package's to read and Close's to
// rewrite (CON-240, plan docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	CrewLearning(project, crew string) string
	LearningsDir(project string) string
}

var _ Workspace = (*store.Workspace)(nil)
