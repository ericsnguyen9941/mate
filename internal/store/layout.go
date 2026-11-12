package store

import (
	"path/filepath"
)

// Names of the directories and files of docs/mvp.md section 3. They are here
// so no other package spells them.
const (
	StateDirName     = ".mate"
	WorktreesDirName = ".worktrees"

	workspaceFileName = "workspace.yaml"
	workspaceDocName  = "WORKSPACE.md"
	pricingFileName   = "pricing.yaml"
	envFileName       = ".env"
	jevLogName        = "jev.log"
	projectsDirName   = "projects"
	projectFileName   = "project.yaml"
	projectDocName    = "PROJECT.md"
	sentLogName       = "sent.log"
	incidentsLogName  = "incidents.log"
	mateDirName       = "mate"
	learningsDirName  = "learnings"
	mateMetaName      = "mate.meta"
	autoFlagName      = ".auto"
	manualHoldName    = ".manual"
	autoCursorName    = ".auto-cursor"
	appNotesName      = "app-notes.log"
	routeLogName      = "route.jsonl"
	wakeCursorName    = ".wake-cursor"
	outboxName        = ".outbox"
	outboxLockName    = ".outbox.lock"
	memoryFileName    = "memory.md"
	backlogFileName   = "backlog.md"
	crewsDirName      = "crews"
	briefFileName     = "brief.md"
	reportFileName    = "report.md"
	handbackFileName  = "handback.md"
	crewDocName       = "CREW.md"
)

// Paths names every file and directory of a workspace from its resolved
// root. Workspace embeds it, so each helper is called on the workspace
// handle as before (`ws.CrewMeta(p, c)`). Its methods are pure string joins:
// they do not touch the filesystem and they do not validate names, because
// the methods that write validate first and resolve the result through the
// workspace boundary. Paths holds only the root, so a helper cannot reach
// the boundary check, the file writers or workspace.yaml (CON-237).
type Paths struct {
	root string
}

// Root is the resolved absolute path of the workspace.
func (p Paths) Root() string { return p.root }

// StateDir is `<root>/.mate`.
func (p Paths) StateDir() string { return filepath.Join(p.root, StateDirName) }

// GCLog is `.mate/gc.log`, where the console's periodic `mate gc` writes
// what it removed and what it could not (CON-114).
func (p Paths) GCLog() string { return filepath.Join(p.StateDir(), "gc.log") }

// WorkspaceFile is `<root>/.mate/workspace.yaml`.
func (p Paths) WorkspaceFile() string { return filepath.Join(p.StateDir(), workspaceFileName) }

// WorkspaceDoc is `<root>/.mate/WORKSPACE.md`, the user's rules for every Mate.
func (p Paths) WorkspaceDoc() string { return filepath.Join(p.StateDir(), workspaceDocName) }

// WorkspaceCrewDoc is `<root>/.mate/CREW.md`, the captain's standing rules
// for every Crew in the workspace, appended to the end of every brief
// (docs/mvp.md M7). Optional.
func (p Paths) WorkspaceCrewDoc() string { return filepath.Join(p.StateDir(), crewDocName) }

// PricingFile is `<root>/.mate/pricing.yaml`, unused before the token monitor.
func (p Paths) PricingFile() string { return filepath.Join(p.StateDir(), pricingFileName) }

// EnvFile is `<root>/.mate/.env`, the workspace's settings for this machine
// (docs/mvp.md section 3): `KEY=VALUE` lines that LoadEnv reads. Optional;
// it holds switches and paths, never a secret itself.
func (p Paths) EnvFile() string { return filepath.Join(p.StateDir(), envFileName) }

// jevLog is `<root>/.mate/jev.log`: one line per request the Jev observer
// chain sent (internal/screen/chain), only while `.mate/.env` turns the
// chain on (a readable key, MATE_JEV not off or fixture). It never holds
// screen text or a key.
func (p Paths) jevLog() string { return filepath.Join(p.StateDir(), jevLogName) }

// ScreensDir is `<root>/.mate/screens`: the text of every screen the Jev
// observer chain was asked to name a pane state for, one file per screen
// (docs/plans/pane-state-jev-2026-10-13.md section 6).
func (p Paths) ScreensDir() string { return filepath.Join(p.StateDir(), "screens") }

// WorktreesDir is `<root>/.worktrees`.
func (p Paths) WorktreesDir() string { return filepath.Join(p.root, WorktreesDirName) }

// RepoDir is the absolute path of a registered repository, from the relative
// path stored in workspace.yaml.
func (p Paths) RepoDir(repo string) string { return filepath.Join(p.root, repo) }

// ConsoleLockFile is `.mate/console.lock`.
func (p Paths) ConsoleLockFile() string {
	return filepath.Join(p.StateDir(), consoleLockName)
}

// CoreLockFile is `.mate/core.lock`.
func (p Paths) CoreLockFile() string { return filepath.Join(p.StateDir(), coreLockName) }

// CoreSocket is `.mate/core.sock`, where `mate core` serves its HTTP routes.
func (p Paths) CoreSocket() string { return filepath.Join(p.StateDir(), coreSocketName) }

// CoreLogFile is `.mate/core.log`, the output of a `mate core` started in
// the background (CON-213).
func (p Paths) CoreLogFile() string { return filepath.Join(p.StateDir(), coreLogName) }

// tuiLockFile is `.mate/tui.lock`.
func (p Paths) tuiLockFile() string { return filepath.Join(p.StateDir(), tuiLockName) }

// recoverLockFile is `.mate/recover.lock`.
func (p Paths) recoverLockFile() string {
	return filepath.Join(p.StateDir(), recoverLockName)
}

// migrateLockFile is `.mate/migrate.lock`, held for the whole of one
// `mate migrate`.
func (p Paths) migrateLockFile() string {
	return filepath.Join(p.StateDir(), migrateLockName)
}

// MigrateLog is `.mate/migrate.log`, the append-only record of every
// directory `mate migrate` renamed and every repo it finished moving.
func (p Paths) MigrateLog() string {
	return filepath.Join(p.StateDir(), migrateLogName)
}
