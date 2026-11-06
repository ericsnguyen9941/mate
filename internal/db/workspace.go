package db

import "github.com/nguyenngocanh94/mate/internal/store"

// Workspace is the part of *store.Workspace that Path, Open and OpenRead
// read: the `.mate/` directory and the boundary check that keeps a symlink
// at `mate.db` from moving the database out of the workspace. Neither
// method writes; the file they name is this package's own (CON-240, plan
// docs/plans/store-split-2026-11-06.md section 3).
type Workspace interface {
	Resolve(path string) (string, error)
	StateDir() string
}

var _ Workspace = (*store.Workspace)(nil)
