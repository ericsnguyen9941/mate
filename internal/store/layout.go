package store

import (
	"fmt"
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

// JevLog is `<root>/.mate/jev.log`: one line per request the Jev observer
// chain sent (internal/screen/chain), only while `.mate/.env` turns the
// chain on (a readable key, MATE_JEV not off or fixture). It never holds
// screen text or a key.
func (p Paths) JevLog() string { return filepath.Join(p.StateDir(), jevLogName) }

// ScreensDir is `<root>/.mate/screens`: the text of every screen the Jev
// observer chain was asked to name a pane state for, one file per screen
// (docs/plans/pane-state-jev-2026-10-13.md section 6).
func (p Paths) ScreensDir() string { return filepath.Join(p.StateDir(), "screens") }

// ProjectsDir is `<root>/.mate/projects`.
func (p Paths) ProjectsDir() string { return filepath.Join(p.StateDir(), projectsDirName) }

// ProjectDir is `<root>/.mate/projects/<project>`.
func (p Paths) ProjectDir(project string) string {
	return filepath.Join(p.ProjectsDir(), project)
}

// ProjectHome is `<root>/<project>`, the project's own directory: the user's,
// not mate's state. Every repo of the project lives under it (docs/mvp.md
// section 3). Mate's state for the project stays in ProjectDir.
func (p Paths) ProjectHome(project string) string {
	return filepath.Join(p.root, project)
}

// ProjectFile is `projects/<project>/project.yaml`.
func (p Paths) ProjectFile(project string) string {
	return filepath.Join(p.ProjectDir(project), projectFileName)
}

// ProjectDoc is `projects/<project>/PROJECT.md`.
func (p Paths) ProjectDoc(project string) string {
	return filepath.Join(p.ProjectDir(project), projectDocName)
}

// ProjectCrewDoc is `projects/<project>/CREW.md`, the captain's standing
// rules for every Crew of one project, appended after the workspace's.
// Optional.
func (p Paths) ProjectCrewDoc(project string) string {
	return filepath.Join(p.ProjectDir(project), crewDocName)
}

// SentLog is `projects/<project>/sent.log`.
func (p Paths) SentLog(project string) string {
	return filepath.Join(p.ProjectDir(project), sentLogName)
}

// IncidentsLog is `projects/<project>/incidents.log`, the observer's
// append-only record of the incidents it opened and resolved (docs/mvp.md
// section 4b). Only internal/watch writes it; internal/box reads it.
func (p Paths) IncidentsLog(project string) string {
	return filepath.Join(p.ProjectDir(project), incidentsLogName)
}

// MateDir is `projects/<project>/mate`, the working directory of Mate.
func (p Paths) MateDir(project string) string {
	return filepath.Join(p.ProjectDir(project), mateDirName)
}

// MateMeta is `projects/<project>/mate/mate.meta`.
func (p Paths) MateMeta(project string) string {
	return filepath.Join(p.MateDir(project), mateMetaName)
}

// AutoFlag is `projects/<project>/mate/.auto`; its presence means auto mode.
func (p Paths) AutoFlag(project string) string {
	return filepath.Join(p.MateDir(project), autoFlagName)
}

// ManualHold is `projects/<project>/mate/.manual`; its presence means the
// captain chose manual mode with the console's `m` key, so context refresh
// does not run on its own while the captain holds the Mate. Its mtime is when
// they chose it: the auto daemon releases a hold the captain has left alone
// for QuietAfter (docs/mvp.md M21).
func (p Paths) ManualHold(project string) string {
	return filepath.Join(p.MateDir(project), manualHoldName)
}

// AutoCursorFile is `projects/<project>/mate/.auto-cursor`; it records how
// far the auto daemon has digested (mvp.md section 5).
func (p Paths) AutoCursorFile(project string) string {
	return filepath.Join(p.MateDir(project), autoCursorName)
}

// WakeCursorFile is `projects/<project>/mate/.wake-cursor`; it records how
// far `mate wake` has shown the project's status lines to the Mate.
func (p Paths) WakeCursorFile(project string) string {
	return filepath.Join(p.MateDir(project), wakeCursorName)
}

// MemoryFile is `projects/<project>/mate/memory.md`.
func (p Paths) MemoryFile(project string) string {
	return filepath.Join(p.MateDir(project), memoryFileName)
}

// BacklogFile is `projects/<project>/mate/backlog.md`.
func (p Paths) BacklogFile(project string) string {
	return filepath.Join(p.MateDir(project), backlogFileName)
}

// CrewsDir is `projects/<project>/crews`.
func (p Paths) CrewsDir(project string) string {
	return filepath.Join(p.ProjectDir(project), crewsDirName)
}

// CrewMeta is `projects/<project>/crews/<crew>.meta`.
func (p Paths) CrewMeta(project, crew string) string {
	return filepath.Join(p.CrewsDir(project), crew+".meta")
}

// CrewStatus is `projects/<project>/crews/<crew>.status`, the append-only log
// a crew writes with echo.
func (p Paths) CrewStatus(project, crew string) string {
	return filepath.Join(p.CrewsDir(project), crew+".status")
}

// CrewPRWatch is `projects/<project>/crews/<crew>.prwatch`: the pid of the
// background `mate pr watch` of the crew's pull request, one number per file
// (docs/mvp.md M18).
func (p Paths) CrewPRWatch(project, crew string) string {
	return filepath.Join(p.CrewsDir(project), crew+".prwatch")
}

// CrewPRWatchLog is `projects/<project>/crews/<crew>/prwatch.log`, where the
// detached watcher writes what it did.
func (p Paths) CrewPRWatchLog(project, crew string) string {
	return filepath.Join(p.CrewDir(project, crew), "prwatch.log")
}

// CrewTeardownLog is `projects/<project>/crews/<crew>/teardown.log`, the
// output of the repo's teardown command at the crew's stop (CON-111).
func (p Paths) CrewTeardownLog(project, crew string) string {
	return filepath.Join(p.CrewDir(project, crew), "teardown.log")
}

// CrewDir is `projects/<project>/crews/<crew>`, kept after teardown.
func (p Paths) CrewDir(project, crew string) string {
	return filepath.Join(p.CrewsDir(project), crew)
}

// CrewBrief is `projects/<project>/crews/<crew>/brief.md`.
func (p Paths) CrewBrief(project, crew string) string {
	return filepath.Join(p.CrewDir(project, crew), briefFileName)
}

// CrewReport is `projects/<project>/crews/<crew>/report.md`.
func (p Paths) CrewReport(project, crew string) string {
	return filepath.Join(p.CrewDir(project, crew), reportFileName)
}

// CrewHandback is `projects/<project>/crews/<crew>/handback.md`, the ship
// Crew's own account of its acceptance checks, written before `wait-mate`
// (docs/mvp.md M7). Like the report it lives outside the worktree, so it
// outlives the Crew.
func (p Paths) CrewHandback(project, crew string) string {
	return filepath.Join(p.CrewDir(project, crew), handbackFileName)
}

// CrewLearning is `projects/<project>/mate/learnings/<crew>.md`, the one
// learning record the Mate keeps per crew (skill task-learning, CON-129).
// It lives under the Mate's own directory, not the crew's: the Mate writes
// it, and `mate wake` reads it to say whether a hand-back has one.
func (p Paths) CrewLearning(project, crew string) string {
	return filepath.Join(p.LearningsDir(project), crew+".md")
}

// LearningsDir is `projects/<project>/mate/learnings`, one record per crew.
func (p Paths) LearningsDir(project string) string {
	return filepath.Join(p.MateDir(project), learningsDirName)
}

// LearningPredictionsSection is the `## ` heading of a learning record under
// which the Mate lists what its lessons predict, one bullet each with an
// `outcome:` (CON-131); internal/learnings reads and closes them.
const LearningPredictionsSection = "Predictions"

// LearningUsageSection is the `## ` heading of the one section of a
// learning record the app writes at a crew's close (CON-130): the crew's
// usage as the ledger has it. The Mate's skill names it and leaves it alone.
const LearningUsageSection = "Usage (app)"

// CrewSessionInput is `projects/<project>/crews/<crew>/session-<n>.md`: the
// opening input of the crew's n-th session (CON-73), composed by a fresh
// relaunch out of the brief, the last hand-back, the Mate's message and
// the diff since the last review. Like the hand-back it outlives the Crew.
func (p Paths) CrewSessionInput(project, crew string, n int) string {
	return filepath.Join(p.CrewDir(project, crew), fmt.Sprintf("session-%d.md", n))
}

// WorktreesDir is `<root>/.worktrees`.
func (p Paths) WorktreesDir() string { return filepath.Join(p.root, WorktreesDirName) }

// WorktreeDir is `<root>/.worktrees/<project>-<crew>`, the git worktree a crew
// runs in.
func (p Paths) WorktreeDir(project, crew string) string {
	return filepath.Join(p.WorktreesDir(), project+"-"+crew)
}

// RepoDir is the absolute path of a registered repository, from the relative
// path stored in workspace.yaml.
func (p Paths) RepoDir(repo string) string { return filepath.Join(p.root, repo) }
