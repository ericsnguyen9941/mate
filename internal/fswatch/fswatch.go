// Package fswatch tells when one of an exact set of files changes (CON-266,
// epic CON-253). The core watches what writes past its change feed - a
// Crew's `echo >> $MATE_STATUS`, the harnesses' transcripts, the captain's
// own edits of workspace.yaml (docs/architecture.md section 2) - so it hears
// of them at once instead of on its next poll.
//
// It watches files, never a directory: kqueue (macOS) opens one descriptor
// per watched file and reports no append to a file inside a watched
// directory, and fsnotify's kqueue backend would open a descriptor for every
// file of one; inotify (Linux) watches a file as well as a directory. A
// burst of writes to one file is told once per window (Options.Delay). A
// file replaced by a rename, or removed and written again, is watched again
// as the new file: its old watch went with the old file. A file that is
// missing is tried again at each Set.
//
// It is a hint, never the record: a watch that could not be set, an event
// the system dropped (inotify's queue overflow), a filesystem with no
// events (WSL2's /mnt/c) - each only means the core's polls find the change
// instead, as they did before. A Watcher that cannot start says so (New's
// error); the caller runs on without it.
package fswatch

import (
	"errors"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// DefaultDelay is the window a burst of writes to one file is folded into:
// the first write opens it, and the file is told once when it closes. A
// file written without a pause (a transcript while its agent works) is told
// once per window, so the delay bounds both the latency and the rate.
const DefaultDelay = 300 * time.Millisecond

// Options configure a Watcher.
type Options struct {
	// Delay is the coalescing window; DefaultDelay when zero.
	Delay time.Duration
	// OnChange is told each watched path that changed: written, replaced
	// or removed. It runs on the Watcher's goroutine, one path at a time,
	// and never after Close returns; it must not block, nor call Set or
	// Close.
	OnChange func(path string)
}

// Stats is what a Watcher has seen since New.
type Stats struct {
	// Watched is the files watched now; Missing the ones of the set that
	// could not be (absent, or the system refused the watch), tried again
	// at the next Set.
	Watched, Missing int
	// Events is every event the system reported for a watched file;
	// Changes the paths told to OnChange, after the bursts were folded.
	Events, Changes uint64
	// Errors is every error: a watch refused, an event queue that
	// overflowed. LastError is the latest, empty with none.
	Errors    uint64
	LastError string
}

// backend is fsnotify's watcher, as the Watcher uses it; a test fails it.
type backend interface {
	Add(name string) error
	Remove(name string) error
	Close() error
	events() <-chan fsnotify.Event
	errors() <-chan error
}

type fsnotifyBackend struct{ *fsnotify.Watcher }

func (b fsnotifyBackend) events() <-chan fsnotify.Event { return b.Events }
func (b fsnotifyBackend) errors() <-chan error          { return b.Errors }

// newBackend is fsnotify's; a test replaces it.
var newBackend = func() (backend, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return fsnotifyBackend{w}, nil
}

// Watcher watches an exact set of files (Set) and tells each change
// (Options.OnChange). Safe for concurrent use.
type Watcher struct {
	b        backend
	delay    time.Duration
	onChange func(string)
	stop     chan struct{}
	done     chan struct{}
	kick     chan struct{}
	close    sync.Once

	mu sync.Mutex
	// files is the set, each with the file it watches; info is nil for a
	// file not watched now (missing).
	files map[string]*file
	// pending is every path told of and not yet passed on, with when its
	// window closes and what happened to it meanwhile.
	pending map[string]*pendingChange
	stats   Stats
}

type file struct {
	info os.FileInfo
}

type pendingChange struct {
	due time.Time
	op  fsnotify.Op
}

// New starts a Watcher with an empty set. An error means the system gives
// no watch (inotify's instances used up, a system fsnotify does not
// support); nothing of it then runs.
func New(o Options) (*Watcher, error) {
	b, err := newBackend()
	if err != nil {
		return nil, err
	}
	delay := o.Delay
	if delay <= 0 {
		delay = DefaultDelay
	}
	onChange := o.OnChange
	if onChange == nil {
		onChange = func(string) {}
	}
	w := &Watcher{
		b: b, delay: delay, onChange: onChange,
		stop: make(chan struct{}), done: make(chan struct{}), kick: make(chan struct{}, 1),
		files: map[string]*file{}, pending: map[string]*pendingChange{},
	}
	go w.run()
	return w, nil
}

// Set makes paths the set: a path no longer in it is no longer watched, a
// new one is, and one already in it is watched again if it was missing or
// has been replaced since (a change then told). Paths are compared as
// given; the caller cleans them.
func (w *Watcher) Set(paths []string) {
	want := make(map[string]bool, len(paths))
	for _, p := range paths {
		want[p] = true
	}
	w.mu.Lock()
	for p := range w.files {
		if !want[p] {
			if w.files[p].info != nil {
				_ = w.b.Remove(p)
			}
			delete(w.files, p)
			delete(w.pending, p)
		}
	}
	replaced := false
	for p := range want {
		f, ok := w.files[p]
		if !ok {
			f = &file{}
			w.files[p] = f
			w.watch(p, f)
			continue
		}
		if f.info == nil {
			w.watch(p, f)
			continue
		}
		// An event lost, or one the system never sends (a rename over a
		// file inotify watches sends nothing until the old file is
		// freed): the file under the path is another one.
		if now, err := os.Stat(p); err != nil || !os.SameFile(now, f.info) {
			w.note(p, fsnotify.Remove)
			replaced = true
		}
	}
	w.count()
	w.mu.Unlock()
	if replaced {
		w.wake()
	}
}

// watch sets f's watch on p, under mu: missing when p does not stat or the
// system refuses it.
func (w *Watcher) watch(p string, f *file) {
	f.info = nil
	info, err := os.Stat(p)
	if err != nil {
		return
	}
	if err := w.b.Add(p); err != nil {
		w.fail(err)
		return
	}
	f.info = info
}

// count refreshes Watched and Missing, under mu.
func (w *Watcher) count() {
	w.stats.Watched, w.stats.Missing = 0, 0
	for _, f := range w.files {
		if f.info != nil {
			w.stats.Watched++
		} else {
			w.stats.Missing++
		}
	}
}

func (w *Watcher) fail(err error) {
	w.stats.Errors++
	w.stats.LastError = err.Error()
}

// Paths is the set, sorted.
func (w *Watcher) Paths() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]string, 0, len(w.files))
	for p := range w.files {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// Stats is what the Watcher has seen so far.
func (w *Watcher) Stats() Stats {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.stats
}

// Close stops watching and waits for the Watcher's goroutine: OnChange is
// not called once it returns. A second Close does nothing.
func (w *Watcher) Close() error {
	var err error
	w.close.Do(func() {
		close(w.stop)
		<-w.done
		err = w.b.Close()
	})
	return err
}

func (w *Watcher) wake() {
	select {
	case w.kick <- struct{}{}:
	default:
	}
}

// run reads the system's events and passes each path on once its window
// closes.
func (w *Watcher) run() {
	defer close(w.done)
	timer := time.NewTimer(time.Hour)
	timer.Stop()
	defer timer.Stop()
	events, errs := w.b.events(), w.b.errors()
	for {
		select {
		case <-w.stop:
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			w.mu.Lock()
			if _, watched := w.files[ev.Name]; watched {
				w.stats.Events++
				w.note(ev.Name, ev.Op)
			}
			w.mu.Unlock()
		case err, ok := <-errs:
			if !ok {
				return
			}
			if errors.Is(err, fsnotify.ErrEventOverflow) {
				// Events were dropped: the polls find what they said.
				err = errors.New("fswatch: the system's event queue overflowed")
			}
			w.mu.Lock()
			w.fail(err)
			w.mu.Unlock()
		case <-w.kick:
		case <-timer.C:
			for _, p := range w.flush(time.Now()) {
				w.onChange(p)
			}
		}
		if next, ok := w.nextDue(); ok {
			timer.Reset(max(time.Until(next), 0))
		}
	}
}

// note folds op into p's pending change, opening its window if none is
// open, under mu.
func (w *Watcher) note(p string, op fsnotify.Op) {
	pc := w.pending[p]
	if pc == nil {
		pc = &pendingChange{due: time.Now().Add(w.delay)}
		w.pending[p] = pc
	}
	pc.op |= op
}

func (w *Watcher) nextDue() (time.Time, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	var next time.Time
	for _, pc := range w.pending {
		if next.IsZero() || pc.due.Before(next) {
			next = pc.due
		}
	}
	return next, !next.IsZero()
}

// flush takes every pending change whose window has closed: a file
// replaced or removed is watched again as the file now under its path (or
// missing), and the paths to tell are returned. An attribute change alone
// on the same file tells nothing.
func (w *Watcher) flush(now time.Time) []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	var out []string
	for p, pc := range w.pending {
		if pc.due.After(now) {
			continue
		}
		delete(w.pending, p)
		f, ok := w.files[p]
		if !ok {
			continue
		}
		changed := pc.op.Has(fsnotify.Write) || pc.op.Has(fsnotify.Create)
		again := pc.op.Has(fsnotify.Remove) || pc.op.Has(fsnotify.Rename)
		if info, err := os.Stat(p); err != nil || f.info == nil || !os.SameFile(info, f.info) {
			again = true
		}
		if again {
			if f.info != nil {
				// fsnotify may have dropped it already; the error says
				// nothing then.
				_ = w.b.Remove(p)
			}
			w.watch(p, f)
			changed = true
		}
		if changed {
			w.stats.Changes++
			out = append(out, p)
		}
	}
	w.count()
	sort.Strings(out)
	return out
}
