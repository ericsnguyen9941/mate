package core

import (
	"sync"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// The change feed (CON-260, epic CON-253): every mate process that writes a
// record under `.mate/` tells the core which one (internal/store's
// NoteWrite, sent by coreclient.Notify to POST coreapi.PathChanged on
// `.mate/core.sock`), so the core can react to a change instead of only
// finding it on its next poll. Today two consumers react, and only faster:
// the shared snapshot is dropped, so the next read of any UI sees the
// change, and the observer starts its next round at once instead of
// waiting out its 5 s. No poll is longer for it - a notice is a hint, and
// what a writer that bypasses the store did (a Crew's own `echo >>
// $MATE_STATUS`, a harness's transcript; docs/architecture.md section 2)
// is found by the polls, and since CON-266 by the engine's own watch on
// those files (filewatch.go), whose changes come in through Changed too.
// Polling less is CON-262's, which waits on a ChangeSub.

// Changes is the engine's change feed: each notice is put in the dirty set
// of every subscriber and wakes it. Safe for concurrent use.
type Changes struct {
	mu   sync.Mutex
	subs []*ChangeSub
}

// DirtyKey is what a notice is about: a Project, and a Crew of it or none.
// The workspace's own files are the zero key.
type DirtyKey struct {
	Project string
	Crew    string
}

// Dirty is what changed since a subscriber last took its set: each key,
// with the kinds of record written under it (store's ChangeKind words).
type Dirty map[DirtyKey]map[string]bool

// ChangeSub is one consumer of the feed: a channel that wakes it, and the
// set of what changed since it last took it.
type ChangeSub struct {
	c     chan struct{}
	mu    sync.Mutex
	dirty Dirty
}

// Note puts c in every subscriber's dirty set and wakes each. It never
// blocks: a subscriber not yet woken from an earlier notice is not woken
// twice, and finds both in its set.
func (f *Changes) Note(c coreapi.Change) {
	key := DirtyKey{Project: c.Project, Crew: c.Crew}
	f.mu.Lock()
	subs := f.subs
	f.mu.Unlock()
	for _, s := range subs {
		s.mark(key, c.Kind)
	}
}

// Subscribe is a new consumer of the feed, with nothing dirty yet.
func (f *Changes) Subscribe() *ChangeSub {
	s := &ChangeSub{c: make(chan struct{}, 1), dirty: Dirty{}}
	f.mu.Lock()
	// A new slice, so Note may range over the one it read unlocked.
	f.subs = append(append([]*ChangeSub(nil), f.subs...), s)
	f.mu.Unlock()
	return s
}

func (s *ChangeSub) mark(key DirtyKey, kind string) {
	s.mu.Lock()
	kinds := s.dirty[key]
	if kinds == nil {
		kinds = map[string]bool{}
		s.dirty[key] = kinds
	}
	kinds[kind] = true
	s.mu.Unlock()
	select {
	case s.c <- struct{}{}:
	default:
	}
}

// C receives once something changed since the subscriber was last woken.
func (s *ChangeSub) C() <-chan struct{} { return s.c }

// Take hands over what changed since the last Take and starts a new set.
// A notice that comes after it is in the next set, and wakes C again.
func (s *ChangeSub) Take() Dirty {
	s.mu.Lock()
	defer s.mu.Unlock()
	d := s.dirty
	s.dirty = Dirty{}
	return d
}

// Waker is a loop the feed wakes before its time: the observer
// (watch.Watcher.Wake), whose next round then starts at once.
type Waker interface {
	Wake()
}

// Changes is the engine's change feed, for a loop to Subscribe to.
func (e *Engine) Changes() *Changes { return e.changes }

// Changed is a writer's notice (POST coreapi.PathChanged) or a change the
// engine's file watch saw (filewatch.go): it is put in the feed, the
// shared snapshot is dropped so the next read of any UI sees the change,
// and the observer is woken - but for a transcript (wakesObserver). A
// notice about a Project or one of its Crews wakes it; one about the
// workspace's own files (a Project added) too, as it re-reads
// workspace.yaml each round.
func (e *Engine) Changed(c coreapi.Change) {
	e.changes.Note(c)
	e.snaps.invalidate()
	if !wakesObserver(c) {
		return
	}
	if w, ok := e.o.Watcher.(Waker); ok {
		w.Wake()
	}
}

// wakesObserver is whether a notice brings the observer's next round
// forward. A transcript does not (CON-266 review): an agent writes it
// without a pause while it works, so its notices would run a round every
// watch.WakeGap instead of every poll exactly when the core is busiest,
// and nothing the observer acts on is in it. Its notice still marks the
// Crew or the Mate dirty for the feed's subscribers, whose next round
// reads it, and drops the snapshot. What the observer acts on - a status
// line, workspace.yaml, a PROJECT.md, a meta - wakes it. The rule is by
// kind, so a writer's notice of that kind would not wake it either.
func wakesObserver(c coreapi.Change) bool {
	return c.Kind != string(store.ChangeTranscript)
}
