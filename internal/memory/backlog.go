package memory

// backlog.md's sections (docs/research/firstmate-memory-2026-09-24.md A3,
// B12). The Mate owns its content; backlog add/move/done provide locked edits.
// Bare `mate backlog` is the live table it is reconciled against.
const (
	BacklogInFlight = "In flight"
	BacklogQueued   = "Queued"
	BacklogDone     = "Done"
	// BacklogHeld is a retired section (CON-136): it held the questions
	// sent to the captain, which are news now, and the promises, which are
	// Queued lines. Recall flags a backlog that still has it so the Mate
	// moves its entries.
	BacklogHeld = "Held for the captain"
)

// BacklogSections is backlog.md's sections in file order.
var BacklogSections = []string{BacklogInFlight, BacklogQueued, BacklogDone}

// DoneKeep is how many Done entries backlog.md keeps; older ones move to
// BacklogArchiveName (firstmate `.tasks.toml` `done_keep = 10`).
const DoneKeep = 10

// The Mate's two cold files, never read at a session start.
const (
	ArchiveName        = "memory-archive.md"
	BacklogArchiveName = "backlog-archive.md"
)

// StowLine is the line the app sends the Mate, after the `⟦mate⟧ `
// sentinel, just before it restarts it (docs/mvp.md task 37, B7). The stow
// skill quotes it and a test holds the two together.
const StowLine = "stow: you are about to be restarted; record anything that exists only in this conversation (skill stow), then end your turn"

// BacklogHeader is what a fresh backlog.md starts with.
func BacklogHeader() string {
	s := "# Backlog\n"
	for _, sec := range BacklogSections {
		s += "\n## " + sec + "\n"
	}
	return s
}
