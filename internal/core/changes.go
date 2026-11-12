package core

import (
	"sync"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/watch"
)

// The change feed (CON-260, epic CON-253): every mate process that writes a
// record under `.mate/` tells the core which one (internal/store's
// NoteWrite, sent by coreclient.Notify to POST coreapi.PathChanged on
// `.mate/core.sock`), so the core can react to a change instead of only
// finding it on its next poll. The shared snapshot is dropped, so the next
// read of any UI sees the change, and the observer starts its next round
// at once instead of waiting out its 5 s - but for a transcript - and
// examines the Projects the notices named, a transcript's included
// (hearObserver, CON-267). Three loops run on it instead of a short clock
// (CON-270, Engine.Changed): workspace.yaml is read back on a `workspace`
// notice, the outbox loop wakes on an `outbox` one, and the notifications
// on any but a transcript. Each keeps a slow tick under it - a notice is a
// hint. What a writer that bypasses the store did (a Crew's own `echo >>
// $MATE_STATUS`, a harness's transcript; docs/architecture.md section 2)
// is in a Project the observer examines every round anyway, one with an
// open Crew or a live Mate, is seen since CON-266 by the engine's own
// watch on those files (filewatch.go), whose changes come in through
// Changed too, and is read in any case by the observer's once-a-minute
// sweep.

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

// Hearer is a loop that reads the feed's dirty set itself, once a round:
// the observer (watch.Watcher.Hear), which examines between its sweeps only
// the projects that can change unheard and those a notice named (CON-267).
type Hearer interface {
	Hear(watch.NoticeFunc)
}

// hearObserver subscribes the observer to the feed. What it hears each
// round is the Projects the notices named; a notice about the workspace's
// own files, or an action the engine ran since the last round, makes the
// round a sweep. The core's own writes send no notice (cmd/mate sets no
// notifier in the core's process): its loops' writes are about Projects
// the observer examines anyway - a Mate that is live, a Crew that is open
// - but an action from a UI can start a Mate or register a Project the
// observer would otherwise meet only at its next sweep. The talk with a
// Mate and the actions that only read (ownLock) are not counted: they are
// a live Mate's.
func (e *Engine) hearObserver() {
	h, ok := e.o.Watcher.(Hearer)
	if !ok {
		return
	}
	sub := e.changes.Subscribe()
	h.Hear(func() watch.Notices {
		dirty := sub.Take()
		n := watch.Notices{Projects: make(map[string]bool, len(dirty)), All: e.acted.Swap(false)}
		for key := range dirty {
			if key.Project == "" {
				n.All = true
				continue
			}
			n.Projects[key.Project] = true
		}
		return n
	})
}

// OutboxWaker is the outbox loop as the change feed wakes it (CON-270):
// outbox.Sender's Wake, for something queued, and Nudge, for a Mate whose
// composer may have emptied, which wakes it only while a line waits.
type OutboxWaker interface {
	Wake()
	Nudge()
}

// Changes is the engine's change feed, for a loop to Subscribe to.
func (e *Engine) Changes() *Changes { return e.changes }

// Changed is a writer's notice (POST coreapi.PathChanged) or a change the
// engine's file watch saw (filewatch.go). A notice about the workspace's
// own files reads workspace.yaml back first, before anything is woken
// (configReloader). Then it is put in the feed, the shared snapshot is
// dropped so the next read of any UI sees the change, and the loops it
// concerns are woken (CON-270): the observer and the notifications by any
// but a transcript (wakesObserver); the outbox loop by an `outbox` notice
// (something queued), and by a `sent` or `mate` one while a line waits
// (the Mate's composer may have emptied).
func (e *Engine) Changed(c coreapi.Change) {
	if c.Kind == string(store.ChangeWorkspace) && e.reload != nil {
		e.reload.reload()
	}
	e.changes.Note(c)
	e.snaps.invalidate()
	if s, ok := e.o.Sender.(OutboxWaker); ok {
		switch store.ChangeKind(c.Kind) {
		case store.ChangeOutbox:
			s.Wake()
		case store.ChangeSent, store.ChangeMate:
			s.Nudge()
		}
	}
	if !wakesObserver(c) {
		return
	}
	if w, ok := e.o.Watcher.(Waker); ok {
		w.Wake()
	}
	// A transcript changes nothing a notification says, and would tick
	// them as often as the observer.
	if e.notify != nil {
		e.notify.Wake()
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
