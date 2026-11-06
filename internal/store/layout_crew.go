package store

import (
	"fmt"
	"path/filepath"
)

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

// CrewStopLock is `projects/<project>/crews/<crew>.stop.lock`, held for the
// whole of one `crew stop`: from before the agent is stopped until the meta
// says the task is over.
//
// The observer reads it. Between the two moments Herdr no longer has the
// agent while the meta still names it, and without the lock the observer
// reads that window as the agent having vanished: every stop of
// hellovietnam's crews up to 2026-10-13 left a `runtime_lost` incident open.
func (p Paths) CrewStopLock(project, crew string) string {
	return filepath.Join(p.CrewsDir(project), crew+".stop.lock")
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

// CrewReviewMeta is `projects/<project>/crews/<crew>/review.meta`, present
// only for a reviewer Crew.
func (p Paths) CrewReviewMeta(project, crew string) string {
	return filepath.Join(p.CrewDir(project, crew), reviewMetaFileName)
}

// CrewHarnessProfile is `projects/<project>/crews/<crew>/harness-profile.json`.
func (p Paths) CrewHarnessProfile(project, crew string) string {
	return filepath.Join(p.CrewDir(project, crew), "harness-profile.json")
}

// CrewSessionInput is `projects/<project>/crews/<crew>/session-<n>.md`: the
// opening input of the crew's n-th session (CON-73), composed by a fresh
// relaunch out of the brief, the last hand-back, the Mate's message and
// the diff since the last review. Like the hand-back it outlives the Crew.
func (p Paths) CrewSessionInput(project, crew string, n int) string {
	return filepath.Join(p.CrewDir(project, crew), fmt.Sprintf("session-%d.md", n))
}

// WorktreeDir is `<root>/.worktrees/<project>-<crew>`, the git worktree a crew
// runs in.
func (p Paths) WorktreeDir(project, crew string) string {
	return filepath.Join(p.WorktreesDir(), project+"-"+crew)
}
