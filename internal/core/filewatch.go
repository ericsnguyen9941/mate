package core

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/fswatch"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// The watch on the files written past the change feed (CON-266, epic
// CON-253). The feed hears what mate writes; a Crew's own `echo >>
// $MATE_STATUS`, the harnesses' transcripts and the captain's edits of
// workspace.yaml or a PROJECT.md it does not (docs/architecture.md section
// 2). The engine watches those files (internal/fswatch) and puts each
// change in the same path as a writer's notice, Engine.Changed: the feed's
// subscribers, the shared snapshot and the observer react to it as to one.
//
// What is watched follows the workspace: workspace.yaml, each Project's
// PROJECT.md, and what the timeline ingest's latest poll found live (each
// open Crew's `.status` and transcript, the running Mate's transcript;
// timeline.Ingester.LiveFiles). The set is applied again after every poll
// - which also retries a file that was missing - and on every notice about
// the workspace's own files. Without an ingest (another process records
// the timeline) only workspace.yaml and the PROJECT.md files are watched.
//
// It is a hint, never the record, and it lengthens no poll: when the
// system gives no watch the engine runs on without it (`mate core stats`
// says so), and the polls find every change as before. Polling less is
// CON-262's.

// fileWatcher is what the engine needs of an fswatch.Watcher; a test fails
// or fakes it.
type fileWatcher interface {
	Set(paths []string)
	Stats() fswatch.Stats
	Close() error
}

// liveFiles is the timeline ingest as the watch reads it.
type liveFiles interface {
	LiveFiles() []timeline.LiveFile
	Polled() <-chan struct{}
}

// openFileWatcher is fswatch.New; a test replaces it on the part.
func openFileWatcher(o fswatch.Options) (fileWatcher, error) {
	w, err := fswatch.New(o)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// fileWatch is the engine's part that runs the watch.
type fileWatch struct {
	ws *store.Workspace
	// live is the ingest; nil without one.
	live liveFiles
	// feed is the engine's change feed, whose notices about the
	// workspace's own files make the set be applied again.
	feed *Changes
	// changed is Engine.Changed.
	changed func(coreapi.Change)
	open    func(fswatch.Options) (fileWatcher, error)

	mu sync.Mutex
	w  fileWatcher
	// set is each watched path and the notice its change is.
	set map[string]coreapi.Change
	// off is why there is no watch; empty while it runs or before Start.
	off string

	cancel context.CancelFunc
	done   chan struct{}
}

// Start opens the watch and applies the set. A system that gives no watch
// does not fail it: the part then runs nothing and says why (stats).
func (f *fileWatch) Start(ctx context.Context) error {
	w, err := f.open(fswatch.Options{OnChange: f.onChange})
	if err != nil {
		f.mu.Lock()
		f.off = err.Error()
		f.mu.Unlock()
		return nil
	}
	f.mu.Lock()
	f.w = w
	f.mu.Unlock()
	sub := f.feed.Subscribe()
	var polled <-chan struct{}
	if f.live != nil {
		polled = f.live.Polled()
	}
	f.apply()
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	f.cancel, f.done = cancel, done
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case <-polled:
				f.apply()
			case <-sub.C():
				if _, ok := sub.Take()[DirtyKey{}]; ok {
					f.apply()
				}
			}
		}
	}()
	return nil
}

// Stop ends the loop, then the watch: no notice comes from it once Stop
// returns.
func (f *fileWatch) Stop() {
	if f.cancel != nil {
		f.cancel()
		<-f.done
	}
	f.mu.Lock()
	w := f.w
	f.mu.Unlock()
	if w != nil {
		_ = w.Close()
	}
}

// files is the set as the workspace and the latest poll have it now.
func (f *fileWatch) files() map[string]coreapi.Change {
	set := map[string]coreapi.Change{
		filepath.Clean(f.ws.WorkspaceFile()): {Kind: string(store.ChangeWorkspace)},
	}
	for _, ref := range f.ws.Projects() {
		set[filepath.Clean(f.ws.ProjectDoc(ref.Name))] = coreapi.Change{Project: ref.Name, Kind: string(store.ChangeProject)}
	}
	if f.live != nil {
		for _, lf := range f.live.LiveFiles() {
			kind := store.ChangeTranscript
			if lf.Kind == timeline.LiveStatus {
				kind = store.ChangeCrewStatus
			}
			set[lf.Path] = coreapi.Change{Project: lf.Project, Crew: lf.Crew, Kind: string(kind)}
		}
	}
	return set
}

// apply makes the watch's set the one files says.
func (f *fileWatch) apply() {
	set := f.files()
	paths := make([]string, 0, len(set))
	for p := range set {
		paths = append(paths, p)
	}
	f.mu.Lock()
	f.set = set
	w := f.w
	f.mu.Unlock()
	if w != nil {
		// Set tells nothing itself, but a file under it may be told of
		// meanwhile, and onChange takes mu.
		w.Set(paths)
	}
}

// onChange is a watched path that changed: its notice goes to the engine.
func (f *fileWatch) onChange(path string) {
	f.mu.Lock()
	c, ok := f.set[path]
	f.mu.Unlock()
	if ok {
		f.changed(c)
	}
}

// stats is what `mate core stats` prints of the watch.
func (f *fileWatch) stats() *coreapi.FileWatchStats {
	f.mu.Lock()
	w, off := f.w, f.off
	f.mu.Unlock()
	if w == nil {
		if off == "" {
			off = "not started"
		}
		return &coreapi.FileWatchStats{Off: off}
	}
	st := w.Stats()
	return &coreapi.FileWatchStats{
		Watched: st.Watched, Missing: st.Missing, Events: st.Events, Changes: st.Changes,
		Errors: st.Errors, LastError: st.LastError,
	}
}
