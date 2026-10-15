package spawn

import (
	"fmt"
	"strconv"
	"strings"
)

// A fix round in a bloated session is where a task's tokens go (CON-73:
// hellovietnam's top ten tasks opened their PR in 8-15 minutes and 4-14M
// tokens, then spent two to three times that on review rounds at 250-434k
// of context). So a crew that has handed back (`pr-open`, or `wait-mate`
// for a local ship) is resting, and the Mate's next message to it opens a
// fresh session instead of typing into the old pane: `mate send` and
// `mate brief append` call RelaunchCrew with a FreshSession, which
// composes the new session's opening input out of exactly four parts, in
// this order - the original brief, the crew's last hand-back, the Mate's
// message, and `git diff --stat <reviewed>...HEAD` - writes it to
// `crews/<id>/session-<n>.md`, and points the new agent at that file.

// FreshSession is what a fresh relaunch composes the new session's opening
// input from, beyond what is already on disk.
type FreshSession struct {
	// Message is the Mate's message: the review findings `mate send`
	// carries, or the words `brief append` added to the brief.
	Message string
	// Reviewed is the commit the latest review of the crew saw (the
	// `head=` of its review record); the fourth part is the diff from it
	// to HEAD. Empty means no review has run and the base branch is used.
	Reviewed string
}

// Meta keys of a crew's sessions (CON-73).
const (
	// MetaSessions counts the harness sessions the crew has had: 1 at
	// spawn, one more at every relaunch that starts a new conversation (a
	// resume continues the same one and does not count). The ledger's
	// `session.crew_session` is the same number derived from the turns.
	MetaSessions = "sessions"
	// MetaSessionStatusOffset is the size of the crew's `.status` file when
	// its current session started, written by a fresh relaunch. A hand-back
	// line below it belongs to an earlier session, so a message right
	// after a relaunch reaches the new pane rather than opening a third.
	MetaSessionStatusOffset = "session_status_offset"
)

// SessionOrdinal is the crew's session count as its meta records it: 1 for
// a crew spawned before the key existed or with nothing readable in it.
func SessionOrdinal(meta map[string]string) int {
	n, err := strconv.Atoi(strings.TrimSpace(meta[MetaSessions]))
	if err != nil || n < 1 {
		return 1
	}
	return n
}

// ComposeSessionInput is the opening input of session n of crew: the four
// parts of a fresh session in their fixed order, each under a numbered
// heading so a reader (and a test) can tell where one ends. diffRange is
// the range the stat was taken over, named so the crew knows what it is
// looking at; an empty stat or hand-back says so rather than leaving a
// heading over nothing.
func ComposeSessionInput(crew string, n int, brief, handback, message, diffRange, diffStat string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Session %d of crew %s\n\n", n, crew)
	b.WriteString("Your earlier session was closed on purpose, so this round does not run in a bloated context; its conversation is gone. The worktree and branch are the ones you left, and everything you need is below, in four parts, in this order.\n\n")
	fmt.Fprintf(&b, "# Part 1: your brief\n\n%s\n\n", strings.TrimRight(brief, "\n"))
	if strings.TrimSpace(handback) == "" {
		handback = "(no hand-back on disk)"
	}
	fmt.Fprintf(&b, "# Part 2: your last hand-back\n\n%s\n\n", strings.TrimRight(handback, "\n"))
	fmt.Fprintf(&b, "# Part 3: the Mate's message\n\n%s\n\n", strings.TrimRight(message, "\n"))
	if strings.TrimSpace(diffStat) == "" {
		diffStat = "(no change)"
	}
	fmt.Fprintf(&b, "# Part 4: change since the last review (git diff --stat %s)\n\n```\n%s\n```\n", diffRange, strings.TrimRight(diffStat, "\n"))
	return b.String()
}

// FreshSessionPrompt is the one line a fresh session is sent: where its
// composed input is, what the file holds in what order, and the note the
// caller carried (what was already tried), in the place RelaunchPrompt
// puts a relaunch note.
func FreshSessionPrompt(inputPath string, n int, note string) string {
	note = strings.Join(strings.Fields(note), " ")
	progress := ""
	if note != "" {
		progress = " Progress so far: " + note + "."
	}
	return fmt.Sprintf("This is session %d of your task; your earlier session was closed on purpose and its conversation is gone. Read and follow %s: your brief, your last hand-back, the Mate's message and the change since the last review, in that order.%s Start now.",
		n, inputPath, progress)
}
