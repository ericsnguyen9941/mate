package fswatch

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

// within is how long a test waits for a change it expects: the window and
// a slow machine's scheduling.
const within = 5 * time.Second

// testDelay is the window the tests fold bursts into.
const testDelay = 100 * time.Millisecond

type recorder struct {
	mu    sync.Mutex
	paths []string
	c     chan string
}

func newRecorder() *recorder { return &recorder{c: make(chan string, 64)} }

func (r *recorder) on(p string) {
	r.mu.Lock()
	r.paths = append(r.paths, p)
	r.mu.Unlock()
	r.c <- p
}

func (r *recorder) want(t *testing.T, path string) {
	t.Helper()
	select {
	case got := <-r.c:
		if got != path {
			t.Fatalf("changed %s, want %s", got, path)
		}
	case <-time.After(within):
		t.Fatalf("%s: no change within %s", path, within)
	}
}

func (r *recorder) none(t *testing.T, wait time.Duration) {
	t.Helper()
	select {
	case got := <-r.c:
		t.Fatalf("changed %s, want nothing", got)
	case <-time.After(wait):
	}
}

func newWatcher(t *testing.T, r *recorder) *Watcher {
	t.Helper()
	w, err := New(Options{Delay: testDelay, OnChange: r.on})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = w.Close() })
	return w
}

// tempFile is a file with text in a directory of its own. The directory is
// resolved, so the path a watch reports is the one the test gave.
func tempFile(t *testing.T, name, text string) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func appendLine(t *testing.T, p, line string) {
	t.Helper()
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

// An append, as a Crew's `echo >> $MATE_STATUS` makes, is told.
func TestAppendIsTold(t *testing.T) {
	r := newRecorder()
	w := newWatcher(t, r)
	p := tempFile(t, "k3.status", "")
	w.Set([]string{p})
	if st := w.Stats(); st.Watched != 1 || st.Missing != 0 {
		t.Fatalf("Stats = %+v, want one watched", st)
	}
	appendLine(t, p, "working: on it")
	r.want(t, p)
	r.none(t, 3*testDelay)
}

// Many writes in one window are one change.
func TestBurstIsOneChange(t *testing.T) {
	r := newRecorder()
	w := newWatcher(t, r)
	p := tempFile(t, "session.jsonl", "")
	w.Set([]string{p})
	for i := range 20 {
		appendLine(t, p, "line")
		if i%5 == 0 {
			time.Sleep(5 * time.Millisecond)
		}
	}
	r.want(t, p)
	r.none(t, 3*testDelay)
	if st := w.Stats(); st.Changes != 1 || st.Events < 2 {
		t.Fatalf("Stats = %+v, want many events and one change", st)
	}
}

// A file replaced by a rename (workspace.yaml, written whole) is told, and
// watched again: the next write of the new file is told too.
func TestReplacedFileIsWatchedAgain(t *testing.T) {
	r := newRecorder()
	w := newWatcher(t, r)
	p := tempFile(t, "workspace.yaml", "projects: []\n")
	w.Set([]string{p})
	for round := range 3 {
		tmp := p + ".tmp"
		if err := os.WriteFile(tmp, []byte("projects: [shop]\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(tmp, p); err != nil {
			t.Fatal(err)
		}
		r.want(t, p)
		appendLine(t, p, "# round")
		r.want(t, p)
		if st := w.Stats(); st.Watched != 1 {
			t.Fatalf("round %d: Stats = %+v, want the file watched", round, st)
		}
	}
}

// A file removed is told and missing; written again, the next Set watches
// it.
func TestRemovedFileIsMissingUntilSet(t *testing.T) {
	r := newRecorder()
	w := newWatcher(t, r)
	p := tempFile(t, "PROJECT.md", "# shop\n")
	w.Set([]string{p})
	if err := os.Remove(p); err != nil {
		t.Fatal(err)
	}
	r.want(t, p)
	if st := w.Stats(); st.Watched != 0 || st.Missing != 1 {
		t.Fatalf("Stats = %+v, want it missing", st)
	}
	if err := os.WriteFile(p, []byte("# shop again\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	w.Set([]string{p})
	if st := w.Stats(); st.Watched != 1 {
		t.Fatalf("Stats = %+v, want it watched again", st)
	}
	appendLine(t, p, "more")
	r.want(t, p)
}

// A path left out of the set is no longer told; a new one is.
func TestSetDropsAndAdds(t *testing.T) {
	r := newRecorder()
	w := newWatcher(t, r)
	a := tempFile(t, "a.status", "")
	b := tempFile(t, "b.status", "")
	missing := filepath.Join(filepath.Dir(a), "not-yet.jsonl")
	w.Set([]string{a, missing})
	if st := w.Stats(); st.Watched != 1 || st.Missing != 1 {
		t.Fatalf("Stats = %+v, want a watched and one missing", st)
	}
	w.Set([]string{b})
	if got := w.Paths(); len(got) != 1 || got[0] != b {
		t.Fatalf("Paths = %v, want only %s", got, b)
	}
	appendLine(t, a, "dropped")
	appendLine(t, b, "added")
	r.want(t, b)
	r.none(t, 3*testDelay)
}

// A file swapped while nothing was said (an event lost) is found at the
// next Set, told and watched again.
func TestSetFindsAReplacedFile(t *testing.T) {
	r := newRecorder()
	w := newWatcher(t, r)
	p := tempFile(t, "k3.status", "")
	w.Set([]string{p})
	// Swap the file under the watch behind its back: drop the watch first,
	// as a lost event would leave it.
	w.mu.Lock()
	_ = w.b.Remove(p)
	w.mu.Unlock()
	tmp := p + ".new"
	if err := os.WriteFile(tmp, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmp, p); err != nil {
		t.Fatal(err)
	}
	w.Set([]string{p})
	r.want(t, p)
	appendLine(t, p, "after")
	r.want(t, p)
}

type failingBackend struct{ backend }

func (failingBackend) Add(string) error { return errors.New("no more watches") }

// A system that refuses a watch leaves the path missing, said in Stats;
// one that gives no watcher at all is New's error.
func TestRefusedWatchesAreSaid(t *testing.T) {
	saved := newBackend
	t.Cleanup(func() { newBackend = saved })
	newBackend = func() (backend, error) {
		b, err := saved()
		if err != nil {
			return nil, err
		}
		return failingBackend{b}, nil
	}
	w := newWatcher(t, newRecorder())
	p := tempFile(t, "k3.status", "")
	w.Set([]string{p})
	if st := w.Stats(); st.Watched != 0 || st.Missing != 1 || st.Errors != 1 || st.LastError != "no more watches" {
		t.Fatalf("Stats = %+v, want the refusal said", st)
	}

	newBackend = func() (backend, error) { return nil, errors.New("inotify: too many instances") }
	if _, err := New(Options{}); err == nil {
		t.Fatal("New = nil with no watcher to be had")
	}
}

// Close leaves no goroutine and no descriptor behind, and a second Close
// does nothing.
func TestCloseLeavesNothing(t *testing.T) {
	fds := openFDs(t)
	goroutines := runtime.NumGoroutine()
	r := newRecorder()
	w, err := New(Options{Delay: testDelay, OnChange: r.on})
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, name := range []string{"a", "b", "c", "d", "e"} {
		paths = append(paths, tempFile(t, name, ""))
	}
	w.Set(paths)
	appendLine(t, paths[0], "x")
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	_ = w.Close()
	deadline := time.Now().Add(within)
	for runtime.NumGoroutine() > goroutines && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := runtime.NumGoroutine(); got > goroutines {
		t.Fatalf("%d goroutines after Close, %d before New", got, goroutines)
	}
	if fds >= 0 {
		if got := openFDs(t); got > fds {
			t.Fatalf("%d descriptors open after Close, %d before New", got, fds)
		}
	}
}

// openFDs counts the process's open descriptors; -1 where the system does
// not list them.
func openFDs(t *testing.T) int {
	t.Helper()
	for _, dir := range []string{"/proc/self/fd", "/dev/fd"} {
		if entries, err := os.ReadDir(dir); err == nil {
			// ReadDir's own descriptor is in the list on both systems.
			return len(entries)
		}
	}
	return -1
}

// The event Op names what the flush makes of it.
func TestAttributeChangeAloneTellsNothing(t *testing.T) {
	r := newRecorder()
	w := newWatcher(t, r)
	p := tempFile(t, "k3.status", "")
	w.Set([]string{p})
	w.mu.Lock()
	w.note(p, fsnotify.Chmod)
	w.mu.Unlock()
	w.wake()
	r.none(t, 3*testDelay)
}
