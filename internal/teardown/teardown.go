// Package teardown removes what a crew's worktree started outside the
// worktree itself before mate gives the worktree back (CON-111): the repo's
// own teardown command, then processes still running in it, the docker
// compose stack started from it, the Xcode DerivedData built from it and the
// simulators named after it.
//
// Ownership is proven, never guessed. A resource belongs to a worktree when
// docker or Xcode recorded the worktree's path for it, or when its name
// carries the worktree's directory name whole; anything else - the
// captain's primary checkouts, a base simulator, a hand-made one - is never
// touched.
//
// Nothing here blocks a stop. A sweep that fails is reported as pending, so
// the task still closes and what was left behind is named.
package teardown

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/process"
)

// DefaultHookTimeout bounds the repo's teardown command. A compose stack or
// a simulator delete takes seconds; five minutes is room for a slow docker
// without letting a hung command hold the stop lock forever.
const DefaultHookTimeout = 5 * time.Minute

// Target is one worktree being given back.
type Target struct {
	// Worktree is the worktree's absolute path. It need not exist any more:
	// the sweeps match on the path docker and Xcode recorded.
	Worktree string
	// Hook is the repo's `teardown:` command, run with `sh -c` in the
	// worktree. Empty, or a worktree that is already gone, runs nothing.
	Hook string
	// Env is added to the hook's environment.
	Env []string
	// Log is where the hook's output is written. Empty discards it.
	Log string
}

// Name is the worktree's directory name, `<project>-<crew>`: unique per
// crew, and what repos name per-worktree resources after.
func (t Target) Name() string {
	return filepath.Base(filepath.Clean(t.Worktree))
}

// roots are the worktree's path as given and, when it still exists, with
// symlinks resolved: docker and Xcode record whichever form they were
// started with.
func (t Target) roots() []string {
	root := filepath.Clean(t.Worktree)
	out := []string{root}
	if real, err := filepath.EvalSymlinks(root); err == nil && real != root {
		out = append(out, real)
	}
	return out
}

// Owns reports whether path is the worktree or inside it.
func (t Target) Owns(path string) bool {
	if strings.TrimSpace(path) == "" || strings.TrimSpace(t.Worktree) == "" {
		return false
	}
	path = filepath.Clean(path)
	for _, root := range t.roots() {
		if path == root || strings.HasPrefix(path, root+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

// NamedAfter reports whether name carries the worktree's directory name as a
// whole token: `PlanTrip-plantrip-pla-6` belongs to worktree
// `plantrip-pla-6`, `PlanTrip-plantrip-pla-6-be-r2` does not.
func (t Target) NamedAfter(name string) bool {
	base := t.Name()
	if base == "" || base == "." || base == string(filepath.Separator) {
		return false
	}
	for from := 0; ; {
		i := strings.Index(name[from:], base)
		if i < 0 {
			return false
		}
		i += from
		end := i + len(base)
		// Before it anything but a letter or digit; after it nothing that
		// could continue a crew id, so `pla-6` never claims `pla-6-be-r2`.
		if (i == 0 || !alnum(name[i-1])) && (end == len(name) || !alnum(name[end]) && !strings.ContainsRune("-_.", rune(name[end]))) {
			return true
		}
		from = i + 1
	}
}

func alnum(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}

// Report is what one teardown did.
type Report struct {
	// Hook is "" when no hook ran, "ok" when it exited 0.
	Hook string
	// Removed names each resource that is gone, one line each.
	Removed []string
	// Pending names what could not be removed and why.
	Pending []string
}

// Ran reports whether the teardown did or tried anything.
func (r Report) Ran() bool {
	return r.Hook != "" || len(r.Removed) > 0 || len(r.Pending) > 0
}

// Meta is the crew meta's `cleanup=` value: `done` or `pending: <why>`,
// with what was removed after it. One line.
func (r Report) Meta() string {
	parts := []string{}
	if r.Hook != "" {
		parts = append(parts, "hook "+r.Hook)
	}
	parts = append(parts, r.Removed...)
	state := "done"
	if len(r.Pending) > 0 {
		state = "pending: " + strings.Join(r.Pending, "; ")
	}
	if len(parts) == 0 {
		return oneLine(state)
	}
	return oneLine(state + " (" + strings.Join(parts, ", ") + ")")
}

// Note is what a stop line adds about the teardown: "; also removed:
// ..." and "; left behind: ...", or "" when there was nothing.
func (r Report) Note() string {
	var b strings.Builder
	removed := r.Removed
	if r.Hook == "ok" {
		removed = append([]string{"teardown command ran"}, removed...)
	}
	if len(removed) > 0 {
		b.WriteString("; also removed: " + strings.Join(removed, ", "))
	}
	if len(r.Pending) > 0 {
		b.WriteString("; left behind: " + strings.Join(r.Pending, "; "))
	}
	return oneLine(b.String())
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// A Sweep finds and removes one kind of resource a worktree owns. A tool
// that is not installed is not an error: it cannot have started anything.
type Sweep interface {
	Sweep(ctx context.Context, t Target) (removed []string, err error)
}

// Cleaner runs a target's hook, then its sweeps in order.
type Cleaner struct {
	// Runner runs the hook and the sweeps' commands. Nil is the real one.
	Runner process.Runner
	// HookTimeout bounds the hook; zero is DefaultHookTimeout.
	HookTimeout time.Duration
	// Sweeps run after the hook, each whatever the one before did.
	Sweeps []Sweep
}

// Live is the cleaner the CLI uses: the hook, then every sweep this machine
// has the tools for.
func Live() Cleaner {
	run := process.ExecRunner{}
	home, _ := os.UserHomeDir()
	return Cleaner{
		Runner: run,
		Sweeps: []Sweep{
			Processes{Runner: run},
			Docker{Runner: run},
			DerivedData{Root: filepath.Join(home, "Library", "Developer", "Xcode", "DerivedData")},
			Simulators{Runner: run},
		},
	}
}

// Clean tears the target down and reports what happened. It never fails:
// every failure is a Pending line.
func (c Cleaner) Clean(ctx context.Context, t Target) Report {
	var r Report
	if strings.TrimSpace(t.Worktree) == "" {
		return r
	}
	if hook := strings.TrimSpace(t.Hook); hook != "" {
		if info, err := os.Stat(t.Worktree); err == nil && info.IsDir() {
			if err := c.runHook(ctx, t, hook); err != nil {
				r.Hook = "failed"
				r.Pending = append(r.Pending, "teardown hook: "+err.Error())
			} else {
				r.Hook = "ok"
			}
		}
	}
	for _, s := range c.Sweeps {
		removed, err := s.Sweep(ctx, t)
		r.Removed = append(r.Removed, removed...)
		if err != nil {
			r.Pending = append(r.Pending, err.Error())
		}
	}
	return r
}

func (c Cleaner) runner() process.Runner {
	if c.Runner != nil {
		return c.Runner
	}
	return process.ExecRunner{}
}

func (c Cleaner) runHook(ctx context.Context, t Target, hook string) error {
	timeout := c.HookTimeout
	if timeout <= 0 {
		timeout = DefaultHookTimeout
	}
	hctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	env := append(os.Environ(), t.Env...)
	res, err := c.runner().Run(hctx, process.Spec{Name: "sh", Args: []string{"-c", hook}, Dir: t.Worktree, Env: env})
	if t.Log != "" {
		writeLog(t.Log, hook, res, err)
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%q did not finish in %s", hook, timeout)
		}
		return fmt.Errorf("%q: %v", hook, err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("%q exited %d%s", hook, res.ExitCode, lastLine(res.Stderr))
	}
	return nil
}

func writeLog(path, hook string, res process.Result, runErr error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, "$ %s\n", hook)
	b.Write(res.Stdout)
	b.Write(res.Stderr)
	if runErr != nil {
		fmt.Fprintf(&b, "\nerror: %v\n", runErr)
	} else {
		fmt.Fprintf(&b, "\nexit %d\n", res.ExitCode)
	}
	_ = os.WriteFile(path, []byte(b.String()), 0o644)
}

// lastLine is ": <the last non-empty line>" of a command's error output,
// or "" when it wrote none.
func lastLine(out []byte) string {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if s := strings.TrimSpace(lines[len(lines)-1]); s != "" {
		return ": " + s
	}
	return ""
}

// missing reports whether err is a tool that is not installed.
func missing(err error) bool {
	return errors.Is(err, exec.ErrNotFound)
}

func sortedKeys(m map[string]bool) []string {
	return slices.Sorted(maps.Keys(m))
}
