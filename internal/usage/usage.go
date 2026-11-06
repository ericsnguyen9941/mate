// Package usage builds the pages `mate usage` prints (mvp.md M5 task 27):
// the token ledger of a project, and one task's model calls, per call or
// per session. cmd/mate's cmdUsage reads the flags, opens the database and
// prints what Write writes (CON-168). `--why`, what a task's calls were
// spent on, and the usage section of a crew's learning record are in
// internal/usage/why (CON-234): they read through dashboard.Server, and
// this package imports nothing of dashboard's, so dashboard can import it.
//
// Every read of the timeline goes through internal/timeline/read; this
// package holds no SQL. No page adds a measurement of its own: every
// number is the ledger's or the projection's.
package usage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/harness"
)

// Deps are what a usage page needs from its caller.
type Deps struct {
	// Harnesses are the harnesses this binary launches. The ledger's
	// autocompact warning names only Crews whose harness takes the
	// threshold (harness.Info.Autocompact); a zero Registry names none.
	Harnesses harness.Registry
}

// Request is one `mate usage` page other than `--why` (package why).
// Write does not check the flags' combinations; cmdUsage refuses a wrong
// one before it opens the database.
type Request struct {
	Project string
	// Crew names one task; empty is the project's ledger.
	Crew string
	// Top keeps the Mate and the Top largest tasks by TASK-TOTAL; zero
	// keeps every task. Only the ledger reads it.
	Top int
	// Sessions prints Crew's model calls per harness session (CON-73).
	Sessions bool
}

// Write writes the page req names to out, reading the timeline through
// handle and the workspace through w.
func Write(ctx context.Context, out io.Writer, w Workspace, handle *db.DB, deps Deps, req Request) error {
	switch {
	case req.Crew == "":
		return writeLedger(ctx, out, w, handle, deps.Harnesses, req.Project, req.Top)
	case req.Sessions:
		return writeCrewSessions(ctx, out, handle, req.Project, req.Crew)
	default:
		return writeCrewTurns(ctx, out, handle, req.Project, req.Crew)
	}
}

// The helpers below are private copies of cmd/mate's, kept here so the
// usage pages do not depend on cmd/mate; they are part of those pages, so
// a change to one is a change to `mate usage`'s output. Package why keeps
// its own copies of formatAge and waitedWord, for `--why`'s head.

// dashIfEmpty keeps a table cell from reading as a missing column.
func dashIfEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// formatAge is an age in its largest whole unit, rounded down.
func formatAge(d time.Duration) string {
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
