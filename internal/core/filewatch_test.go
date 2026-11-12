package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/fswatch"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// The watch on the files written past the change feed (CON-266).

// noticeWithin bounds how long a write takes to come back as a notice: the
// watch's window (300 ms) and a slow machine.
const noticeWithin = 5 * time.Second

// fakeLive is the ingest's live files as a test sets them.
type fakeLive struct {
	mu     sync.Mutex
	files  []timeline.LiveFile
	polled chan struct{}
}

func newFakeLive(files ...timeline.LiveFile) *fakeLive {
	return &fakeLive{files: files, polled: make(chan struct{}, 1)}
}

func (l *fakeLive) LiveFiles() []timeline.LiveFile {
	l.mu.Lock()
	defer l.mu.Unlock()
	return slices.Clone(l.files)
}

func (l *fakeLive) Polled() <-chan struct{} { return l.polled }

// poll is a poll that found files live.
func (l *fakeLive) poll(files ...timeline.LiveFile) {
	l.mu.Lock()
	l.files = files
	l.mu.Unlock()
	l.polled <- struct{}{}
}

// resolved is w's root with its symlinks resolved (macOS's /var is
// /private/var), so the paths the test writes are the ones watched.
func resolvedWorkspace(t *testing.T) *store.Workspace {
	t.Helper()
	w := shopWorkspace(t)
	root, err := filepath.EvalSymlinks(w.Root())
	if err != nil {
		t.Fatal(err)
	}
	w, err = store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

// crewFiles are an open crew's status file and transcript, created.
func crewFiles(t *testing.T, w *store.Workspace, crew string) []timeline.LiveFile {
	t.Helper()
	if err := w.WriteCrewMeta("shop", crew, map[string]string{"harness": "codex", "state": "spawned"}); err != nil {
		t.Fatal(err)
	}
	status := w.CrewStatus("shop", crew)
	transcript := filepath.Join(w.Root(), "transcripts", crew+".jsonl")
	for _, p := range []string{status, transcript} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := writeFile(p, ""); err != nil {
			t.Fatal(err)
		}
	}
	return []timeline.LiveFile{
		{Project: "shop", Crew: crew, Kind: timeline.LiveStatus, Path: status},
		{Project: "shop", Crew: crew, Kind: timeline.LiveTranscript, Path: transcript},
	}
}

func appendTo(t *testing.T, p, line string) {
	t.Helper()
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatal(err)
	}
}

// noticed waits for a notice of want in the feed.
func noticed(t *testing.T, sub *ChangeSub, want coreapi.Change) {
	t.Helper()
	deadline := time.After(noticeWithin)
	for {
		select {
		case <-sub.C():
			if sub.Take()[DirtyKey{Project: want.Project, Crew: want.Crew}][want.Kind] {
				return
			}
		case <-deadline:
			t.Fatalf("no notice of %+v within %s", want, noticeWithin)
		}
	}
}

// watched is the paths the part watches now.
func watched(f *fileWatch) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.w.(*fswatch.Watcher).Paths()
}

func waitWatched(t *testing.T, f *fileWatch, path string, want bool) {
	t.Helper()
	deadline := time.Now().Add(noticeWithin)
	for slices.Contains(watched(f), path) != want {
		if time.Now().After(deadline) {
			t.Fatalf("%s watched = %v after %s, want %v", path, !want, noticeWithin, want)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// A Crew's own `echo >> $MATE_STATUS`, its transcript, the captain's edit
// of workspace.yaml and the Mate's of PROJECT.md each come to the engine as
// a notice with the Project, the Crew and the kind; a crew that closes is
// no longer watched, one that opens is.
func TestFileWatchFeedsTheEngine(t *testing.T) {
	w := resolvedWorkspace(t)
	doc := w.ProjectDoc("shop")
	if err := writeFile(doc, "# shop\n"); err != nil {
		t.Fatal(err)
	}
	k3 := crewFiles(t, w, "k3")
	live := newFakeLive(k3...)
	e := New(Options{Workspace: w})
	e.files.live = live
	sub := e.Changes().Subscribe()
	if err := e.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(e.Stop)

	appendTo(t, k3[0].Path, "working: on it")
	noticed(t, sub, coreapi.Change{Project: "shop", Crew: "k3", Kind: "status"})
	appendTo(t, k3[1].Path, `{"type":"turn"}`)
	noticed(t, sub, coreapi.Change{Project: "shop", Crew: "k3", Kind: "transcript"})
	appendTo(t, w.WorkspaceFile(), "# the captain was here")
	noticed(t, sub, coreapi.Change{Kind: "workspace"})
	appendTo(t, doc, "more")
	noticed(t, sub, coreapi.Change{Project: "shop", Kind: "project"})

	// k3 closes and m7 opens: the next poll says so.
	m7 := crewFiles(t, w, "m7")
	live.poll(m7...)
	waitWatched(t, e.files, m7[0].Path, true)
	waitWatched(t, e.files, k3[0].Path, false)
	if slices.Contains(watched(e.files), k3[1].Path) {
		t.Fatalf("k3's transcript still watched after it closed: %v", watched(e.files))
	}
	appendTo(t, m7[0].Path, "working: next")
	noticed(t, sub, coreapi.Change{Project: "shop", Crew: "m7", Kind: "status"})

	st := e.Stats().Files
	if st == nil || st.Off != "" || st.Watched != 4 || st.Changes < 5 {
		t.Fatalf("Stats().Files = %+v, want 4 watched and the changes counted", st)
	}
}

// A system that gives no watch neither stops the engine nor its feed: the
// stats say why, and a writer's notice still lands.
func TestFileWatchThatCannotStartLeavesTheEngine(t *testing.T) {
	w := shopWorkspace(t)
	e := New(Options{Workspace: w})
	e.files.open = func(fswatch.Options) (fileWatcher, error) {
		return nil, errors.New("inotify: too many instances")
	}
	sub := e.Changes().Subscribe()
	if err := e.Start(context.Background()); err != nil {
		t.Fatalf("Start = %v, want the engine running without the watch", err)
	}
	t.Cleanup(e.Stop)
	if st := e.Stats().Files; st == nil || st.Off != "inotify: too many instances" {
		t.Fatalf("Stats().Files = %+v, want it off with the reason", st)
	}
	e.Changed(coreapi.Change{Project: "shop", Crew: "k3", Kind: "status"})
	noticed(t, sub, coreapi.Change{Project: "shop", Crew: "k3", Kind: "status"})
}

// Options.NoFileWatch runs no watch at all.
func TestNoFileWatch(t *testing.T) {
	e := New(Options{Workspace: shopWorkspace(t), NoFileWatch: true})
	if e.files != nil {
		t.Fatal("a file watch with NoFileWatch")
	}
	if st := e.Stats().Files; st != nil {
		t.Fatalf("Stats().Files = %+v, want none", st)
	}
}

// Stop ends the watch: nothing of it is told after (internal/fswatch's own
// test counts the goroutines and descriptors Close leaves).
func TestFileWatchStopLeavesNothing(t *testing.T) {
	w := resolvedWorkspace(t)
	k3 := crewFiles(t, w, "k3")
	var mu sync.Mutex
	var after []coreapi.Change
	stopped := false
	f := &fileWatch{
		ws: w, live: newFakeLive(k3...), feed: &Changes{}, open: openFileWatcher,
		changed: func(c coreapi.Change) {
			mu.Lock()
			defer mu.Unlock()
			if stopped {
				after = append(after, c)
			}
		},
	}
	if err := f.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	appendTo(t, k3[0].Path, "working")
	f.Stop()
	mu.Lock()
	stopped = true
	mu.Unlock()
	appendTo(t, k3[0].Path, "working more")
	time.Sleep(2 * fswatch.DefaultDelay)
	mu.Lock()
	defer mu.Unlock()
	if len(after) > 0 {
		t.Fatalf("told %+v after Stop", after)
	}
}
