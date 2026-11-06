package store

import "path/filepath"

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

// OutboxFile is `projects/<project>/mate/.outbox`.
func (p Paths) OutboxFile(project string) string {
	return filepath.Join(p.MateDir(project), outboxName)
}

func (p Paths) outboxLockFile(project string) string {
	return filepath.Join(p.MateDir(project), outboxLockName)
}

// MemoryFile is `projects/<project>/mate/memory.md`.
func (p Paths) MemoryFile(project string) string {
	return filepath.Join(p.MateDir(project), memoryFileName)
}

// BacklogFile is `projects/<project>/mate/backlog.md`.
func (p Paths) BacklogFile(project string) string {
	return filepath.Join(p.MateDir(project), backlogFileName)
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
