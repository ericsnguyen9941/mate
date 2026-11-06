package store

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

// ErrNotWorkspace is returned by Open for a directory with no `.mate/`.
var ErrNotWorkspace = errors.New("store: not a mate workspace")

const (
	// DefaultBranch is what a project falls back to when none is given.
	DefaultBranch = "main"

	// workspaceVersion is the schema version written into workspace.yaml.
	workspaceVersion = 1

	// layoutProjectDirs is the layout where a project is a directory under
	// the workspace and every repo of it lives under that directory
	// (docs/mvp.md section 3, 2026-10-08). Init writes it; a workspace.yaml
	// without `layout:` is layout 1, repos beside `.mate/`.
	layoutProjectDirs = 2
)

// ErrLayoutOld is Open's refusal of a workspace still on layout 1: only
// `mate migrate` and the two stops its refusals name, `mate crew stop` and
// `mate mate stop`, read one (OpenForMigrate).
var ErrLayoutOld = errors.New("this workspace has the old layout (repos beside .mate); run mate migrate first")

// Workspace is an opened workspace: a resolved root plus the parsed
// workspace.yaml. It is the handle every other package holds.
//
// One Workspace may be shared by goroutines: in the console, the refresh
// reloads workspace.yaml while the Office and City handlers, the actions and
// recovery read it (CON-204). Nothing changes the loaded configuration in
// place. LoadConfig and every change of workspace.yaml build a new one and
// swap it in whole, so a reader never waits and never sees parts of two.
// The swaps take turns under one lock, each holding it from its read to its
// swap with the file I/O between, so a reload cannot land inside a change
// and undo it. Two Workspace values on one directory still each hold their
// own copy, and the last one to save wins.
type Workspace struct {
	// Paths is the resolved root and every path helper built from it.
	Paths
	// cfg is the loaded workspace.yaml: set once Init or OpenForMigrate
	// returns, and never changed in place.
	cfg atomic.Pointer[WorkspaceConfig]
	// cfgMu serialises LoadConfig and every change of workspace.yaml.
	cfgMu sync.Mutex
}

// Init creates `.mate/` under workspaceDir and writes a workspace.yaml
// recording defaults, which the caller resolves from its harness registry
// so a new workspace's file says what it runs. It is idempotent for a
// directory that already has one: the existing workspace is opened instead,
// so the stored session name and defaults survive.
func Init(workspaceDir string, defaults Defaults) (*Workspace, error) {
	// A directory that does not exist yet is made, as `git init <dir>`
	// does; Open, which only reads, still refuses one.
	if err := os.MkdirAll(workspaceDir, 0o755); err != nil {
		return nil, fmt.Errorf("store: create workspace %s: %w", workspaceDir, err)
	}
	root, err := resolveRoot(workspaceDir)
	if err != nil {
		return nil, err
	}
	w := &Workspace{Paths: Paths{root: root}}
	if _, err := os.Stat(w.WorkspaceFile()); err == nil {
		return Open(root)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	w.cfg.Store(&WorkspaceConfig{
		Version:  workspaceVersion,
		Layout:   layoutProjectDirs,
		Session:  SessionName(root),
		Root:     root,
		Defaults: defaults,
	})
	if err := w.mkdirAll(w.ProjectsDir()); err != nil {
		return nil, err
	}
	if err := w.SaveConfig(); err != nil {
		return nil, err
	}
	// The captain's rules for every Mate. Every Mate's bootstrap reads this
	// file first (mateassets, section 3), so an init that did not create it
	// left each Mate opening with a "no such file" error on a file its own
	// manual names (measured 2026-09-19, task 24). Seeded once, never
	// rewritten: it is the captain's file from here on.
	if err := os.WriteFile(w.WorkspaceDoc(), []byte(workspaceDocSeed), 0o644); err != nil {
		return nil, err
	}
	// The token price table (mvp.md M5 task 27). Seeded once, at zero, with
	// a comment telling the captain it is theirs to fill in: mate records
	// tokens whether or not anyone has priced them, and an init that left
	// this file missing would make every ledger's cost column silently
	// unreadable rather than honestly "?".
	if err := os.WriteFile(w.PricingFile(), []byte(pricingFileSeed), 0o644); err != nil {
		return nil, err
	}
	// The captain's standing rules for every Crew (docs/mvp.md M7), seeded
	// with only a comment so it is discoverable and still adds nothing to a
	// brief until the captain writes a rule.
	if err := w.seedOnce(w.WorkspaceCrewDoc(), workspaceCrewDocSeed); err != nil {
		return nil, err
	}
	return w, nil
}

// workspaceDocSeed is the WORKSPACE.md a new workspace starts with: a place
// for the captain's rules, empty of rules, so a Mate that reads it learns
// only that nothing workspace-wide has been said yet.
const workspaceDocSeed = `# Workspace rules

Rules here apply to every Mate in this workspace, on top of each project's PROJECT.md.
Nothing is written here yet.
`

// Open resolves workspaceDir, which must contain `.mate/workspace.yaml`, and
// loads the configuration. A workspace still on the old layout is refused
// with ErrLayoutOld.
func Open(workspaceDir string) (*Workspace, error) {
	w, err := OpenForMigrate(workspaceDir)
	if err != nil {
		return nil, err
	}
	if w.LayoutOld() {
		return nil, ErrLayoutOld
	}
	return w, nil
}

// OpenForMigrate is Open for `mate migrate` and the stops it sends the
// captain to (`mate crew stop`, `mate mate stop`): it opens a workspace on
// either layout, the old one included.
func OpenForMigrate(workspaceDir string) (*Workspace, error) {
	root, err := resolveRoot(workspaceDir)
	if err != nil {
		return nil, err
	}
	w := &Workspace{Paths: Paths{root: root}}
	if fi, err := os.Stat(w.StateDir()); err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("%w: %s has no %s directory", ErrNotWorkspace, root, StateDirName)
	}
	if err := w.LoadConfig(); err != nil {
		return nil, err
	}
	return w, nil
}

// resolveRoot makes the workspace root absolute and canonical. A workspace
// reached through a symlink is ordinary (on macOS /tmp is /private/tmp), so
// the link is followed and the target becomes the root; every boundary check
// afterwards is against that resolved path.
func resolveRoot(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return "", fmt.Errorf("store: workspace %s is not a directory", resolved)
	}
	return resolved, nil
}

// SessionName is the Herdr session name for a workspace path: `mate-` plus a
// short stable hash of the absolute path. Init stores the result so it never
// changes for a workspace that already exists.
func SessionName(absPath string) string {
	sum := sha256.Sum256([]byte(absPath))
	return "mate-" + hex.EncodeToString(sum[:4])
}
