// Package loopstat records what each of the core's loops costs in wall time
// (CON-263): how many runs, the last, the mean and the longest, and when the
// last one started. `mate core stats` reads it from the running core.
//
// On 2026-11-11 `mate core run` was found at 1.4 GB RSS and ~70% of a CPU
// for days: one timeline ingest pass took 4.9 s of every 5 s poll, and
// nothing said so. A loop records each run as it ends; reading is the
// reader's cost. A run costs two clock reads and one uncontended lock, and
// no allocation: nothing here reads the memory statistics, whose read
// stops the world. Allocation is not split between loops: in a process
// where they run side by side, no counter says which one allocated what.
//
// It is a leaf: the loops of watch, outbox and autopilot record into it as
// well as the core's own, and cmd/mate hands the same Recorder to all of
// them and to the engine.
package loopstat

import (
	"sync"
	"time"
)

// Recorder is the loops of one process, in the order they were first
// named. Safe for concurrent use. A nil *Recorder hands out nil loops,
// which record nothing.
type Recorder struct {
	mu     sync.Mutex
	loops  []*Loop
	byName map[string]*Loop
}

// New is a Recorder with no loop yet.
func New() *Recorder { return &Recorder{byName: map[string]*Loop{}} }

// Loop is the loop named name, made the first time it is asked for: a loop
// named before it ever runs is read as zero runs, which says it is wired.
func (r *Recorder) Loop(name string) *Loop {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if l, ok := r.byName[name]; ok {
		return l
	}
	l := &Loop{name: name}
	r.byName[name] = l
	r.loops = append(r.loops, l)
	return l
}

// Read is every loop's figures now, in the order they were first named.
func (r *Recorder) Read() []Stat {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	loops := append([]*Loop(nil), r.loops...)
	r.mu.Unlock()
	out := make([]Stat, 0, len(loops))
	for _, l := range loops {
		out = append(out, l.read())
	}
	return out
}

// Loop is one loop's figures. Its methods are safe on a nil *Loop, which
// records nothing, so a part built without a Recorder needs no check.
type Loop struct {
	name string

	mu      sync.Mutex
	runs    uint64
	total   time.Duration
	last    time.Duration
	max     time.Duration
	lastAt  time.Time
	running int
	since   time.Time
}

// Begin says a run starts now and returns when, for End:
//
//	defer l.End(l.Begin())
func (l *Loop) Begin() time.Time {
	if l == nil {
		return time.Time{}
	}
	now := time.Now()
	l.mu.Lock()
	if l.running == 0 {
		l.since = now
	}
	l.running++
	l.mu.Unlock()
	return now
}

// End records the run Begin started at start.
func (l *Loop) End(start time.Time) {
	if l == nil {
		return
	}
	took := time.Since(start)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.runs++
	l.total += took
	l.last = took
	l.lastAt = start
	if took > l.max {
		l.max = took
	}
	if l.running > 0 {
		l.running--
	}
}

func (l *Loop) read() Stat {
	l.mu.Lock()
	defer l.mu.Unlock()
	s := Stat{Name: l.name, Runs: l.runs, Total: l.total, Last: l.last, Max: l.max, LastAt: l.lastAt}
	if l.runs > 0 {
		s.Mean = l.total / time.Duration(l.runs)
	}
	if l.running > 0 {
		since := l.since
		s.RunningSince = &since
	}
	return s
}

// Stat is one loop's figures, in wall time: a run that waited on git,
// Herdr or a lock counts the wait. internal/coreapi's LoopStats carries
// them to `mate core stats`.
type Stat struct {
	Name string
	// Runs is how many runs have ended.
	Runs uint64
	// Last is the last run's, Mean the mean of all, Max the longest and
	// Total their sum.
	Last, Mean, Max, Total time.Duration
	// LastAt is when the last run started; zero before the first ends.
	LastAt time.Time
	// RunningSince is when the run in flight started, nil when none is:
	// a loop stuck in one run shows here and not in Runs.
	RunningSince *time.Time
}
