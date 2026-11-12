package core

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/coreclient"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/watch"
)

// The change feed (CON-260): a writer's notice lands in every subscriber's
// dirty set, drops the snapshot and wakes the observer.

// Every subscriber gets every notice, woken once however many came, and
// Take hands over what came since the last Take.
func TestChangesMarkAndWakeEverySubscriber(t *testing.T) {
	f := &Changes{}
	a, b := f.Subscribe(), f.Subscribe()
	f.Note(coreapi.Change{Project: "shop", Crew: "fix-cart", Kind: "meta"})
	f.Note(coreapi.Change{Project: "shop", Crew: "fix-cart", Kind: "status"})
	f.Note(coreapi.Change{Project: "shop", Kind: "sent"})
	want := Dirty{
		{Project: "shop", Crew: "fix-cart"}: {"meta": true, "status": true},
		{Project: "shop"}:                   {"sent": true},
	}
	for name, s := range map[string]*ChangeSub{"a": a, "b": b} {
		select {
		case <-s.C():
		default:
			t.Fatalf("%s was not woken", name)
		}
		select {
		case <-s.C():
			t.Fatalf("%s was woken twice for one wait", name)
		default:
		}
		if got := s.Take(); !reflect.DeepEqual(got, want) {
			t.Fatalf("%s took %v, want %v", name, got, want)
		}
		if got := s.Take(); len(got) != 0 {
			t.Fatalf("%s's second Take = %v, want nothing", name, got)
		}
	}
}

// Notices from many writers at once, a subscriber taking meanwhile: no
// race (make test-race), nothing lost - each key is in some Take.
func TestChangesConcurrentWritersLoseNothing(t *testing.T) {
	f := &Changes{}
	s := f.Subscribe()
	const writers, each = 8, 200
	var wg sync.WaitGroup
	for i := range writers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range each {
				f.Note(coreapi.Change{Project: "shop", Crew: fmt.Sprintf("c%d-%d", i, j), Kind: "status"})
			}
		}()
	}
	// A late subscriber, while the notices come.
	late := make(chan *ChangeSub, 1)
	go func() { late <- f.Subscribe() }()
	seen := map[DirtyKey]bool{}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	for finished := false; !finished; {
		select {
		case <-s.C():
		case <-done:
			finished = true
		}
		for k := range s.Take() {
			seen[k] = true
		}
	}
	for k := range s.Take() {
		seen[k] = true
	}
	<-late
	if len(seen) != writers*each {
		t.Fatalf("saw %d crews, want %d", len(seen), writers*each)
	}
}

// wakingWatcher is a watcher that counts its wakes.
type wakingWatcher struct {
	fakeWatcher
	wakes atomic.Int32
}

func (w *wakingWatcher) Wake() { w.wakes.Add(1) }

// A notice drops the snapshot, so the next read computes one that sees the
// write, and wakes the observer.
func TestChangedDropsSnapshotAndWakesWatcher(t *testing.T) {
	w := &wakingWatcher{}
	e := New(Options{Watcher: w})
	defer e.Stop()
	var computed atomic.Int32
	e.snaps.compute = func(context.Context) (query.Snapshot, error) {
		computed.Add(1)
		return query.Snapshot{}, nil
	}
	sub := e.Changes().Subscribe()
	ctx := context.Background()
	for range 2 {
		if _, err := e.Snapshot(ctx); err != nil {
			t.Fatal(err)
		}
	}
	if got := computed.Load(); got != 1 {
		t.Fatalf("two reads within 2s computed %d snapshots, want 1", got)
	}
	e.Changed(coreapi.Change{Project: "shop", Crew: "fix-cart", Kind: "status"})
	if _, err := e.Snapshot(ctx); err != nil {
		t.Fatal(err)
	}
	if got := computed.Load(); got != 2 {
		t.Fatalf("a read after a notice computed %d snapshots in all, want 2", got)
	}
	if got := w.wakes.Load(); got != 1 {
		t.Fatalf("the watcher was woken %d times, want 1", got)
	}
	want := Dirty{{Project: "shop", Crew: "fix-cart"}: {"status": true}}
	if got := sub.Take(); !reflect.DeepEqual(got, want) {
		t.Fatalf("dirty = %v, want %v", got, want)
	}
}

// awaitDirty waits for key to be in what s takes, within two seconds.
func awaitDirty(t *testing.T, s *ChangeSub, key DirtyKey, kind string) {
	t.Helper()
	deadline := time.After(2 * time.Second)
	for {
		if s.Take()[key][kind] {
			return
		}
		select {
		case <-s.C():
		case <-deadline:
			t.Fatalf("no %s notice for %+v within 2s", kind, key)
		}
	}
}

// End to end: a store write in another process's shoes (the notifier
// cmd/mate sets, coreclient.Notify) reaches `mate core` on its socket, in
// the dirty set, the observer woken. Not a parallel test: it sets the
// store's process-wide notifier.
func TestStoreWriteReachesTheCore(t *testing.T) {
	w := shortWorkspace(t)
	withTimelineDB(t, w)
	watcher := &wakingWatcher{fakeWatcher: fakeWatcher{fakeLoop: fakeLoop{name: "observer", rec: &recorder{}}}}
	e := socketEngine(t, w, Options{Watcher: watcher})
	sub := e.Changes().Subscribe()
	store.SetChangeNotifier(func(c store.Change) {
		coreclient.Notify(c.Socket, coreapi.Change{Project: c.Project, Crew: c.Crew, Kind: string(c.Kind)})
	})
	t.Cleanup(func() { store.SetChangeNotifier(nil) })

	// A writer of its own handle, as a CLI command opens one.
	writer, err := store.Open(w.Root())
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.AppendStatus("shop", "fix-cart", "working: cart"); err != nil {
		t.Fatal(err)
	}
	awaitDirty(t, sub, DirtyKey{Project: "shop", Crew: "fix-cart"}, string(store.ChangeCrewStatus))
	if watcher.wakes.Load() == 0 {
		t.Fatal("the notice did not wake the watcher")
	}
	if err := writer.SetMode("shop", false); err != nil {
		t.Fatal(err)
	}
	awaitDirty(t, sub, DirtyKey{Project: "shop"}, string(store.ChangeMate))
}

// A slow core does not hold a writer: the notice is written and the writer
// goes on, whatever the core then does with it.
func TestNoticeDoesNotWaitForTheCore(t *testing.T) {
	w := shortWorkspace(t)
	withTimelineDB(t, w)
	release := make(chan struct{})
	e := socketEngine(t, w, Options{Watcher: &blockingWakeWatcher{
		fakeWatcher: fakeWatcher{fakeLoop: fakeLoop{name: "observer", rec: &recorder{}}},
		release:     release,
	}})
	// Before the engine's own Stop, which waits for the handlers.
	t.Cleanup(func() { close(release) })
	sub := e.Changes().Subscribe()
	start := time.Now()
	for range 5 {
		coreclient.Notify(w.CoreSocket(), coreapi.Change{Project: "shop", Kind: "sent"})
	}
	if took := time.Since(start); took > 4*coreclient.NotifyWait {
		t.Fatalf("five notices to a core that does not answer took %s", took)
	}
	awaitDirty(t, sub, DirtyKey{Project: "shop"}, "sent")
}

// blockingWakeWatcher's Wake does not return until release closes: the
// handler of every notice hangs.
type blockingWakeWatcher struct {
	fakeWatcher
	release chan struct{}
}

func (w *blockingWakeWatcher) Wake() { <-w.release }

// hearingWatcher is an observer that reads the feed itself (Hearer).
type hearingWatcher struct {
	fakeWatcher
	hear watch.NoticeFunc
}

func (w *hearingWatcher) Hear(f watch.NoticeFunc) { w.hear = f }

// The observer hears the feed once a round (CON-267): the Projects the
// notices named since its last round; a notice about the workspace's own
// files, or an action the engine ran itself - whose writes send no notice -
// makes the round a sweep. The talk with a Mate and a peek do not.
func TestObserverHearsTheFeedAndTheEnginesOwnActions(t *testing.T) {
	w := &hearingWatcher{}
	e := New(Options{Watcher: w, Act: func(context.Context, coreapi.ActionRequest) (string, error) { return "", nil }})
	defer e.Stop()
	if w.hear == nil {
		t.Fatal("the engine did not hand the observer its feed")
	}
	e.Changed(coreapi.Change{Project: "shop", Crew: "fix-cart", Kind: "status"})
	e.Changed(coreapi.Change{Project: "docs", Kind: "sent"})
	if got, want := w.hear(), (watch.Notices{Projects: map[string]bool{"shop": true, "docs": true}}); !reflect.DeepEqual(got, want) {
		t.Fatalf("heard %+v, want %+v", got, want)
	}
	if got := w.hear(); got.All || len(got.Projects) != 0 {
		t.Fatalf("heard %+v again, want nothing", got)
	}
	e.Changed(coreapi.Change{Kind: "workspace"})
	if got := w.hear(); !got.All {
		t.Fatalf("a notice about workspace.yaml heard as %+v, want a sweep", got)
	}
	ctx := context.Background()
	if _, err := e.Act(ctx, coreapi.ActionRequest{Action: coreapi.ActionPeek, Target: "shop", TargetKind: "crew", Crew: "buy"}); err != nil {
		t.Fatal(err)
	}
	if got := w.hear(); got.All {
		t.Fatal("a peek made the next round a sweep")
	}
	if _, err := e.Act(ctx, coreapi.ActionRequest{Action: coreapi.ActionStop, Target: "shop", TargetKind: "crew", Crew: "buy"}); err != nil {
		t.Fatal(err)
	}
	if got := w.hear(); !got.All {
		t.Fatal("an action the engine ran did not make the next round a sweep")
	}
	if got := w.hear(); got.All {
		t.Fatal("one action made two sweeps")
	}
}
