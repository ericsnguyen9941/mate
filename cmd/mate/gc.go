package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/process"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/teardown"
)

// gcDeps are what a collection needs: the sweeps, and how worktrees that
// are gone but still pointed at are found. Tests replace both.
type gcDeps struct {
	Collect func(ctx context.Context, targets []teardown.Target, dry bool) teardown.Report
	Orphans func(ctx context.Context, root string) []teardown.Target
}

func liveGCDeps() gcDeps {
	run := process.ExecRunner{}
	return gcDeps{
		Collect: teardown.Live().Collect,
		Orphans: func(ctx context.Context, root string) []teardown.Target {
			return teardown.Orphans(ctx, root, run, teardown.XcodeDerivedData())
		},
	}
}

// cmdGC implements `mate gc [--apply] [--workspace <dir>]` (CON-114): what
// closed tasks' worktrees started and a crew stop did not remove - a stop
// from before CON-112, a docker daemon that was down, a simulator that
// would not delete. Without --apply it only lists.
func cmdGC(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("gc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate gc [--apply] [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	applyFlag := fs.Bool("apply", false, "remove what the dry run lists")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return newUsageError("mate gc: takes no arguments")
	}
	w, err := resolveWorkspace(*workspaceFlag)
	if err != nil {
		return err
	}
	targets, r := runGC(context.Background(), w, liveGCDeps(), *applyFlag)
	fmt.Fprint(stdout, gcReport(targets, r, *applyFlag))
	return nil
}

// runGC collects for every worktree that is gone: those of closed crews,
// and those docker or Xcode still name under `.worktrees/`.
func runGC(ctx context.Context, w *store.Workspace, d gcDeps, apply bool) ([]teardown.Target, teardown.Report) {
	targets := gcTargets(ctx, w, d)
	return targets, d.Collect(ctx, targets, !apply)
}

func gcTargets(ctx context.Context, w *store.Workspace, d gcDeps) []teardown.Target {
	seen := map[string]bool{}
	var out []teardown.Target
	add := func(t teardown.Target) {
		key := filepath.Clean(t.Worktree)
		if !seen[key] {
			seen[key] = true
			out = append(out, t)
		}
	}
	for _, ref := range w.Projects() {
		ids, err := w.CrewIDs(ref.Name)
		if err != nil {
			continue
		}
		for _, id := range ids {
			meta, err := w.ReadCrewMeta(ref.Name, id)
			if err != nil || !crewstate.State(meta[spawn.MetaState]).Closed() || meta[spawn.MetaWorktree] == "" {
				continue
			}
			wt := filepath.Join(w.Root(), filepath.FromSlash(meta[spawn.MetaWorktree]))
			if _, err := os.Stat(wt); err == nil {
				continue
			}
			add(teardown.Target{Worktree: wt})
		}
	}
	if d.Orphans != nil {
		for _, t := range d.Orphans(ctx, w.WorktreesDir()) {
			add(t)
		}
	}
	return out
}

// gcReport is what `mate gc` prints: a header, then one line per resource.
func gcReport(targets []teardown.Target, r teardown.Report, apply bool) string {
	var b strings.Builder
	verb := "would remove"
	if apply {
		verb = "removed"
	}
	switch {
	case len(r.Removed) == 0 && len(r.Pending) == 0:
		fmt.Fprintf(&b, "mate gc: %d closed worktree(s), nothing of theirs left\n", len(targets))
	case apply:
		fmt.Fprintf(&b, "mate gc: %d closed worktree(s); removed %d resource(s)\n", len(targets), len(r.Removed))
	default:
		fmt.Fprintf(&b, "mate gc: %d closed worktree(s); dry run, nothing removed - rerun with --apply to remove %d resource(s)\n", len(targets), len(r.Removed))
	}
	for _, line := range r.Removed {
		fmt.Fprintf(&b, "  %s: %s\n", verb, line)
	}
	for _, line := range r.Pending {
		fmt.Fprintf(&b, "  left behind: %s\n", line)
	}
	return b.String()
}

// GCEvery is how often a running console collects: often enough that a
// missed stack does not sit for days, rarely enough that its docker,
// simctl and lsof calls cost nothing anyone notices.
const GCEvery = 6 * time.Hour

// gcSchedule runs a collection from the console's tick: once at the first
// tick, then every GCEvery, never two at a time and never on the tick's
// own time.
type gcSchedule struct {
	mu      sync.Mutex
	running bool
	last    time.Time
	every   time.Duration
	now     func() time.Time
	run     func(ctx context.Context)
}

// Tick starts a collection when one is due. A nil schedule does nothing.
func (g *gcSchedule) Tick(ctx context.Context) {
	if g == nil {
		return
	}
	g.mu.Lock()
	now := g.now()
	if g.running || (!g.last.IsZero() && now.Sub(g.last) < g.every) {
		g.mu.Unlock()
		return
	}
	g.running, g.last = true, now
	g.mu.Unlock()
	go func() {
		defer func() {
			g.mu.Lock()
			g.running = false
			g.mu.Unlock()
		}()
		g.run(ctx)
	}()
}

// consoleGC is the console's periodic `mate gc --apply`. It writes to
// `.mate/gc.log` only when it removed or failed to remove something.
func consoleGC(ws *store.Workspace, d gcDeps) *gcSchedule {
	return &gcSchedule{
		every: GCEvery,
		now:   time.Now,
		run: func(ctx context.Context) {
			ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
			defer cancel()
			targets, r := runGC(ctx, ws, d, true)
			if len(r.Removed) == 0 && len(r.Pending) == 0 {
				return
			}
			f, err := os.OpenFile(ws.GCLog(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				return
			}
			defer f.Close()
			fmt.Fprintf(f, "%s %s", time.Now().Format(time.RFC3339), gcReport(targets, r, true))
		},
	}
}
