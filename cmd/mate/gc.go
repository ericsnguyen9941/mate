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
	// Shutdown shuts down the booted simulators of idle crews' worktrees
	// (CON-119). Nil skips that pass.
	Shutdown func(ctx context.Context, targets []teardown.Target, dry bool) ([]string, error)
	// Now is the clock the idle pass measures a crew's silence by; nil is
	// time.Now.
	Now func() time.Time
}

func liveGCDeps() gcDeps {
	run := process.ExecRunner{}
	return gcDeps{
		Collect: teardown.Live().Collect,
		Orphans: func(ctx context.Context, root string) []teardown.Target {
			return teardown.Orphans(ctx, root, run, teardown.XcodeDerivedData())
		},
		Shutdown: teardown.Simulators{Runner: run}.Shutdown,
	}
}

// gcIdleAfter is how long a crew that waits on the Mate, the captain or an
// incident keeps its simulator booted. A crew pla-90 that had finished at
// 12:47 still held 1.8 GB and 180 processes at 15:40 (CON-119); a crew that
// comes back boots it again in seconds.
const gcIdleAfter = 30 * time.Minute

// gcIdle is one crew whose simulator the idle pass shuts down.
type gcIdle struct {
	Project, Crew, State string
	Age                  time.Duration
	Target               teardown.Target
}

func (c gcIdle) String() string {
	return fmt.Sprintf("crew %s/%s %s %s", c.Project, c.Crew, c.State, compactDuration(c.Age))
}

// compactDuration is `3h`, `2h15m` or `45m`, never seconds.
func compactDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h, m := int(d.Hours()), int(d.Minutes())%60
	switch {
	case h > 0 && m > 0:
		return fmt.Sprintf("%dh%dm", h, m)
	case h > 0:
		return fmt.Sprintf("%dh", h)
	default:
		return fmt.Sprintf("%dm", m)
	}
}

// gcIdleReport is what the idle pass did: one line per simulator, shut
// down or only listed, and what it could not shut down.
type gcIdleReport struct {
	Crews   []gcIdle
	Done    []string
	Pending []string
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
	idle := runGCIdle(context.Background(), w, liveGCDeps(), *applyFlag)
	fmt.Fprint(stdout, gcReport(targets, r, idle, *applyFlag))
	return nil
}

// runGC collects for every worktree that is gone: those of closed crews,
// and those docker or Xcode still name under `.worktrees/`.
func runGC(ctx context.Context, w *store.Workspace, d gcDeps, apply bool) ([]teardown.Target, teardown.Report) {
	targets := gcTargets(ctx, w, d)
	return targets, d.Collect(ctx, targets, !apply)
}

// runGCIdle shuts down the booted simulators of crews that have waited
// longer than gcIdleAfter (CON-119). Without a Shutdown it does nothing.
func runGCIdle(ctx context.Context, w *store.Workspace, d gcDeps, apply bool) gcIdleReport {
	if d.Shutdown == nil {
		return gcIdleReport{}
	}
	crews := gcIdleCrews(w, d)
	if len(crews) == 0 {
		return gcIdleReport{}
	}
	targets := make([]teardown.Target, 0, len(crews))
	for _, c := range crews {
		targets = append(targets, c.Target)
	}
	done, err := d.Shutdown(ctx, targets, !apply)
	r := gcIdleReport{Crews: crews}
	for _, line := range done {
		r.Done = append(r.Done, line+gcIdleOwner(crews, line))
	}
	if err != nil {
		r.Pending = append(r.Pending, err.Error())
	}
	return r
}

// gcIdleOwner names the crew whose worktree the simulator line carries,
// so the report says who left it booted.
func gcIdleOwner(crews []gcIdle, line string) string {
	for _, c := range crews {
		if c.Target.NamedAfter(strings.TrimPrefix(line, "simulator ")) {
			return " (" + c.String() + ")"
		}
	}
	return ""
}

// gcIdleCrews are the open crews whose declared state is a wait - on the
// Mate, on the captain, on an incident - and whose status file has not
// changed for gcIdleAfter. A crew that is working, or only just stopped
// talking, keeps its simulator.
func gcIdleCrews(w *store.Workspace, d gcDeps) []gcIdle {
	now := time.Now
	if d.Now != nil {
		now = d.Now
	}
	var out []gcIdle
	for _, ref := range w.Projects() {
		crews, err := spawn.ListCrews(w, ref.Name)
		if err != nil {
			continue
		}
		for _, c := range crews {
			switch crewstate.State(c.State) {
			case crewstate.StateWaitMate, crewstate.StateNeedsDecision, crewstate.StateBlocked:
			default:
				continue
			}
			meta, err := w.ReadCrewMeta(ref.Name, c.Crew)
			if err != nil || meta[spawn.MetaWorktree] == "" {
				continue
			}
			info, err := os.Stat(w.CrewStatus(ref.Name, c.Crew))
			if err != nil {
				continue
			}
			age := now().Sub(info.ModTime())
			if age < gcIdleAfter {
				continue
			}
			out = append(out, gcIdle{
				Project: ref.Name, Crew: c.Crew, State: c.State, Age: age,
				Target: teardown.Target{Worktree: filepath.Join(w.Root(), filepath.FromSlash(meta[spawn.MetaWorktree]))},
			})
		}
	}
	return out
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

// gcReport is what `mate gc` prints: a header, then one line per resource,
// then one per simulator of an idle crew.
func gcReport(targets []teardown.Target, r teardown.Report, idle gcIdleReport, apply bool) string {
	var b strings.Builder
	verb, stop := "would remove", "would shut down"
	if apply {
		verb, stop = "removed", "shut down"
	}
	head := fmt.Sprintf("%d closed worktree(s)", len(targets))
	if len(idle.Crews) > 0 {
		head += fmt.Sprintf(", %d idle crew(s)", len(idle.Crews))
	}
	switch {
	case len(r.Removed) == 0 && len(r.Pending) == 0 && len(idle.Done) == 0 && len(idle.Pending) == 0:
		fmt.Fprintf(&b, "mate gc: %s, nothing of theirs left\n", head)
	case apply:
		fmt.Fprintf(&b, "mate gc: %s; removed %d resource(s), shut down %d simulator(s)\n", head, len(r.Removed), len(idle.Done))
	default:
		fmt.Fprintf(&b, "mate gc: %s; dry run, nothing removed - rerun with --apply to remove %d resource(s) and shut down %d simulator(s)\n", head, len(r.Removed), len(idle.Done))
	}
	for _, line := range r.Removed {
		fmt.Fprintf(&b, "  %s: %s\n", verb, line)
	}
	for _, line := range r.Pending {
		fmt.Fprintf(&b, "  left behind: %s\n", line)
	}
	for _, line := range idle.Done {
		fmt.Fprintf(&b, "  %s: %s\n", stop, line)
	}
	for _, line := range idle.Pending {
		fmt.Fprintf(&b, "  still booted: %s\n", line)
	}
	return b.String()
}

// GCEvery is how often a running console collects: often enough that a
// missed stack does not sit for days, rarely enough that its docker,
// simctl and lsof calls cost nothing anyone notices.
const GCEvery = 6 * time.Hour

// GCIdleEvery is how often a running console shuts down the simulators of
// crews that have waited gcIdleAfter: a booted simulator costs the machine
// every minute it sits, so this runs far more often than a collection.
const GCIdleEvery = 30 * time.Minute

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

// consoleGC is the console's periodic `mate gc --apply` over gone
// worktrees. It writes to `.mate/gc.log` only when it removed or failed to
// remove something.
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
			appendGCLog(ws, gcReport(targets, r, gcIdleReport{}, true))
		},
	}
}

// consoleGCIdle is the console's idle pass (CON-119): every GCIdleEvery it
// shuts down the simulators of crews that have waited gcIdleAfter, and
// logs only when it shut one down or could not.
func consoleGCIdle(ws *store.Workspace, d gcDeps) *gcSchedule {
	return &gcSchedule{
		every: GCIdleEvery,
		now:   time.Now,
		run: func(ctx context.Context) {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			defer cancel()
			idle := runGCIdle(ctx, ws, d, true)
			if len(idle.Done) == 0 && len(idle.Pending) == 0 {
				return
			}
			appendGCLog(ws, gcReport(nil, teardown.Report{}, idle, true))
		},
	}
}

func appendGCLog(ws *store.Workspace, report string) {
	f, err := os.OpenFile(ws.GCLog(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s %s", time.Now().Format(time.RFC3339), report)
}
