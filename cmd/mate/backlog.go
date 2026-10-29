package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/memory"
	"github.com/nguyenngocanh94/mate/internal/plugin"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// cmdBacklog reads the live crew table or edits one durable backlog entry.
func cmdBacklog(args []string, stdout, stderr io.Writer) error {
	if len(args) > 0 && (args[0] == "add" || args[0] == "move" || args[0] == "done") {
		return cmdBacklogEdit(args, stdout, stderr)
	}
	fs := flag.NewFlagSet("backlog", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate backlog <project> [--all] [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	allFlag := fs.Bool("all", false, "include closed crews (those `crew stop` has run on)")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return newUsageError("mate backlog: want exactly 1 argument: <project>")
	}
	w, err := resolveWorkspace(*workspaceFlag)
	if err != nil {
		return err
	}
	return printBacklog(w, fs.Arg(0), *allFlag, time.Now(), stdout)
}

// backlogRow is one line of the table: the crew's declared state (the same
// derivation query.Load uses, by way of spawn.ListCrews, both of which
// resolve through crewstate.Declare) plus the recorded facts a Mate needs
// to reconcile its own notes against.
type backlogRow struct {
	ID     string
	State  string
	Age    time.Duration
	HasAge bool
	Branch string
	Task   string
	Last   string
}

// printBacklog is cmdBacklog's core, kept apart from flag parsing so a test
// can drive it with a fixed clock instead of time.Now().
//
// It asks spawn.ListCrews for the crews and their declared state - the one
// place that ordering is implemented (mvp.md section 4b) - and reads each
// crew's own `.meta` again only for started_at, which ListCrews's
// CrewSummary does not carry. Nothing here re-derives a state.
func printBacklog(w *store.Workspace, project string, all bool, now time.Time, stdout io.Writer) error {
	summaries, err := spawn.ListCrews(w, project)
	if err != nil {
		return err
	}
	if len(summaries) == 0 {
		fmt.Fprintf(stdout, "no crews for %s\n", project)
		return nil
	}

	rows := make([]backlogRow, 0, len(summaries))
	closed := 0
	for _, s := range summaries {
		if s.Closed {
			closed++
			if !all {
				continue
			}
		}
		age, hasAge := crewAge(w, project, s.Crew, now)
		rows = append(rows, backlogRow{
			ID:     s.Crew,
			State:  s.State,
			Age:    age,
			HasAge: hasAge,
			Branch: dashIfEmpty(s.Branch),
			Task:   dashIfEmpty(s.Task),
			Last:   dashIfEmpty(s.Note),
		})
	}
	open := len(summaries) - closed

	// Oldest first: the crew that has been open longest leads the table,
	// because that is the one a Mate reconciling its backlog after a
	// restart most needs to see is still accounted for. A crew whose age
	// could not be determined sorts after every crew whose age is known,
	// and among those, by id, so the order stays fixed across runs.
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].HasAge != rows[j].HasAge {
			return rows[i].HasAge
		}
		if rows[i].HasAge {
			return rows[i].Age > rows[j].Age
		}
		return rows[i].ID < rows[j].ID
	})

	tw := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tSTATE\tAGE\tBRANCH\tTASK\tLAST")
	for _, r := range rows {
		age := "-"
		if r.HasAge {
			age = formatAge(r.Age)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", r.ID, r.State, age, r.Branch, r.Task, r.Last)
	}
	if err := tw.Flush(); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "%d open, %d closed\n", open, closed)

	var attention []string
	for _, r := range rows {
		switch crewstate.State(r.State) {
		case crewstate.StateWaitMate, crewstate.StateNeedsDecision, crewstate.StateBlocked:
			attention = append(attention, r.ID)
		}
	}
	if len(attention) > 0 {
		fmt.Fprintf(stdout, "attention: %s\n", strings.Join(attention, ", "))
	}
	return nil
}

// crewAge reads a crew's own `.meta` for started_at= and reports how long
// ago that was. A meta that cannot be read, or carries no started_at, or
// carries one that does not parse as RFC3339, reports hasAge=false rather
// than guessing: an unknown age is a dash in the table, never a zero.
func crewAge(w *store.Workspace, project, crew string, now time.Time) (time.Duration, bool) {
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil {
		return 0, false
	}
	raw := strings.TrimSpace(meta[spawn.MetaStartedAt])
	if raw == "" {
		return 0, false
	}
	startedAt, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return 0, false
	}
	age := now.Sub(startedAt)
	if age < 0 {
		age = 0
	}
	return age, true
}

// formatAge renders a duration the way the AGE column wants it: the
// largest whole unit that keeps the number small enough to read at a
// glance, seconds under a minute up through days. It rounds down, the way
// every other age reading in this codebase does (internal/ui/console's
// shortDuration is the sibling for pane quiet time) - a crew spawned 3
// days and 23 hours ago has been open for 3 days, not 4.
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

// cmdBacklogEdit implements `mate backlog add|move|done <project> <id>
// [--section <section>] [--text <text>] [--tracker <id>]`.
//
// `--tracker <id>` (CON-77) is for a project whose tasks live in a tracker
// plugin (the `tracker:` block of an instance in .mate/plugins.yaml): the
// entry written, or the one moved, becomes the one pointer line
// `- <crew-id>: <tracker-id>`, under `## In flight` unless --section says
// otherwise, and no content is copied - the tracker holds it. Without
// --tracker nothing changes: `## Queued` entries without a tracker id are written as they always were. A project
// with no tracker plugin refuses --tracker, so a pointer never points at
// nothing.
func cmdBacklogEdit(args []string, stdout, stderr io.Writer) error {
	action := args[0]
	fs := flag.NewFlagSet("backlog "+action, flag.ContinueOnError)
	fs.SetOutput(stderr)
	workspace := fs.String("workspace", "", "workspace directory")
	section := fs.String("section", "", "In flight, Queued, or Done (default: Queued; In flight with --tracker)")
	text := fs.String("text", "", "one complete entry; a promise to the captain in their words")
	tracker := fs.String("tracker", "", "the tracker's own id for this task (such as HVN-12); the entry becomes the pointer `- <id>: <tracker-id>`; needs a tracker plugin")
	if err := fs.Parse(reorderArgs(fs, args[1:])); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return newUsageError("usage: mate backlog " + action + " <project> <id> [--section <section>] [--text <text>] [--tracker <tracker-id>] [--workspace <dir>]")
	}
	trackerID := strings.TrimSpace(*tracker)
	if *tracker != "" && (trackerID == "" || strings.ContainsAny(trackerID, " \t\r\n")) {
		return newUsageError("mate backlog " + action + ": --tracker wants one tracker id without spaces, such as HVN-12")
	}
	if *section == "" {
		*section = memory.BacklogQueued
		if trackerID != "" {
			*section = memory.BacklogInFlight
		}
	}
	if action == "done" {
		*section = memory.BacklogDone
	}
	w, err := resolveWorkspace(*workspace)
	if err != nil {
		return err
	}
	if err := requireProject(w, fs.Arg(0)); err != nil {
		return err
	}
	if trackerID == "" {
		if err := w.EditBacklog(fs.Arg(0), action, fs.Arg(1), *section, *text); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "backlog %s: %s → %s\n", action, fs.Arg(1), *section)
		return nil
	}
	in, err := projectTracker(w, fs.Arg(0))
	if err != nil {
		return err
	}
	if err := w.EditBacklogEntry(fs.Arg(0), action, fs.Arg(1), *section, trackerID, true); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "backlog %s: %s → %s (pointer to %s in %s)\n", action, fs.Arg(1), *section, trackerID, in.ID)
	return nil
}

// projectTracker is the tracker plugin instance of project, read from
// .mate/plugins.yaml without loading any package. No tracker is a usage
// error naming the block to add: --tracker is a claim that the tracker
// holds the task, and nothing here can make that true.
func projectTracker(w *store.Workspace, project string) (plugin.Instance, error) {
	cfg, err := plugin.ReadConfig(w.Root())
	if err != nil {
		return plugin.Instance{}, err
	}
	in, ok := cfg.TrackerFor(project)
	if !ok {
		return plugin.Instance{}, newUsageErrorf("mate backlog: project %s has no tracker plugin; --tracker needs a `tracker:` block on the plugin instance in %s (docs/plugins.md)", project, filepath.Join(w.Root(), plugin.ConfigFile))
	}
	return in, nil
}
