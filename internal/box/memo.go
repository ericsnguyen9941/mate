package box

import (
	"os"
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// memo keeps what Load last read of each file, by path, so a Load of a
// project whose files have not moved reads none of them again (CON-259).
//
// Measured 2026-11-11 on the captain's workspace (1308 crew directories, a
// handful alive): one query.Load cost ~100 ms and 62 MB of allocation, 88%
// of it in box.Load - parsing every closed crew's `.meta` to learn it is
// closed, every `.status`, and a 2 MB `sent.log` - every two seconds, for
// files that had not changed in weeks. Now a file whose size and mtime are
// what they were costs one stat; an append-only file that grew (`.status`,
// `sent.log`) is read from where the last read stopped; a `.meta` that moved
// is parsed again.
//
// Box writes nothing and the memo changes nothing a reader can see: it holds
// exactly what the store's readers would return for a file at that size and
// mtime, which is the identity the timeline's transcript cache and its
// telemetry cursor already trust. Every writer of these files appends under
// the store's lock or replaces the file whole (store.WriteCrewMeta), so a
// change always moves the size or the mtime. LoadSince keeps its own cursor
// and does not use the memo.
//
// It is per process and bounded by the workspace: each Load keeps only the
// files it read, so a crew whose files are gone is forgotten by the next
// Load of its project.
//
// A Load none of whose files moved - the same crews, every file at the size
// and mtime the last one saw, the same mode - returns the View the last one
// built instead of merging, sorting and grouping the same lines again
// (CON-271). The View is shared with the memo and every caller that got it:
// its slices and maps are read only, and each slice is clipped to its length
// so an append copies rather than writes into the shared array.
var memo = &fileMemo{projects: map[string]*projectMemo{}}

type fileMemo struct {
	mu       sync.Mutex
	projects map[string]*projectMemo // by the project's crews directory
}

// projectMemo is one project's files as the last Load saw them. A Load
// builds a new one and swaps it in whole, so two Loads that run at once each
// see a consistent one and the later swap wins.
type projectMemo struct {
	crews     dirMemo
	status    map[string]linesMemo[store.StatusEntry]
	closed    map[string]closedMemo
	sent      linesMemo[store.SentEntry]
	incidents linesMemo[store.IncidentEntry]
	// view is what the Load that kept this memo built from exactly these
	// files, nil when one of them could not be pinned to a stamp.
	view *View
}

// dirMemo is the crews directory's `.status` names at the directory's
// stamp. Creating, removing or renaming a file in a directory moves its
// mtime, and every writer of a crew file does one of those. It keeps each
// crew's two paths beside its id, so a Load that finds nothing moved does
// not spell them again.
type dirMemo struct {
	at     stamp
	ok     bool
	crews  []string
	meta   []string
	status []string
}

// generation numbers every View a memoized Load builds, so a reader can
// tell a View it has seen from a new one (View.Generation).
var generation atomic.Uint64

// stamp is a file's identity for the memo: its size and mtime. A file that
// is not there has none (ok false).
type stamp struct {
	size, mtime int64
}

func statStamp(path string) (stamp, os.FileInfo, bool) {
	fi, err := os.Stat(path)
	if err != nil {
		return stamp{}, nil, false
	}
	return stamp{size: fi.Size(), mtime: fi.ModTime().UnixNano()}, fi, true
}

// linesMemo is an append-only file's parsed lines and the offset its last
// read stopped at.
type linesMemo[T any] struct {
	at    stamp
	lines []T
	next  int64
}

// read returns the file's lines from the start, reading only what the memo
// does not hold. readFrom is the store's reader. The returned slice is
// shared with the memo and must not be written to.
func (m linesMemo[T]) read(now stamp, readFrom func(int64) ([]T, int64, error)) (linesMemo[T], error) {
	if m.at == now && m.next > 0 {
		return m, nil
	}
	from, prior := int64(0), []T(nil)
	if m.next > 0 && now.size > m.at.size && now.size >= m.next {
		// Grown: the lines held are still the file's first ones.
		from, prior = m.next, m.lines
	}
	lines, next, err := readFrom(from)
	if err != nil {
		return linesMemo[T]{}, err
	}
	// A full slice expression, so the append copies instead of writing
	// into an array a concurrent Load may still be reading.
	all := append(prior[:len(prior):len(prior)], lines...)
	return linesMemo[T]{at: now, lines: all, next: next}, nil
}

// closedMemo is whether a crew's `.meta` declares it closed.
type closedMemo struct {
	at     stamp
	closed bool
}

func (f *fileMemo) project(key string) *projectMemo {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.projects[key]
}

func (f *fileMemo) keep(key string, p *projectMemo) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.projects[key] = p
}

// loadReader reads one project's files for a Load through the memo.
type loadReader struct {
	ws      Workspace
	project string
	held    *projectMemo
	next    *projectMemo
	// moved is set by every read that did not find its file as held:
	// the View held cannot be the answer.
	moved bool
	// volatile is set by a read whose answer is not pinned to a stamp (a
	// listed status file that is gone by the time it is read): the View
	// built is not kept.
	volatile bool
}

func newLoadReader(ws Workspace, project string) *loadReader {
	held := memo.project(ws.CrewsDir(project))
	if held == nil {
		held = &projectMemo{}
	}
	return &loadReader{
		ws: ws, project: project, held: held,
		next: &projectMemo{status: map[string]linesMemo[store.StatusEntry]{}, closed: map[string]closedMemo{}},
	}
}

// listCrews is listCrews through the memo. The reader's crewClosed and
// status take a crew by its index in the list it returns.
func (r *loadReader) listCrews() ([]string, error) {
	dir := r.ws.CrewsDir(r.project)
	now, _, ok := statStamp(dir)
	if ok && r.held.crews.ok && r.held.crews.at == now {
		r.next.crews = r.held.crews
		return r.held.crews.crews, nil
	}
	crews, err := listCrews(r.ws, r.project)
	if err != nil {
		return nil, err
	}
	if !slices.Equal(crews, r.held.crews.crews) {
		r.moved = true
	}
	list := dirMemo{at: now, ok: ok, crews: crews, meta: make([]string, len(crews)), status: make([]string, len(crews))}
	for i, crew := range crews {
		list.meta[i] = filepath.Join(dir, crew+".meta")
		list.status[i] = r.ws.CrewStatus(r.project, crew)
	}
	if !ok {
		r.volatile = true
	}
	r.next.crews = list
	return crews, nil
}

// crewClosed is crewstate.Declare(meta).Closed() for the i-th crew's
// `.meta`. A meta that cannot be read is not closed and is not remembered.
func (r *loadReader) crewClosed(i int) bool {
	crew, path := r.next.crews.crews[i], r.next.crews.meta[i]
	now, _, ok := statStamp(path)
	if !ok {
		return false
	}
	if held, ok := r.held.closed[path]; ok && held.at == now {
		r.next.closed[path] = held
		return held.closed
	}
	r.moved = true
	meta, err := r.ws.ReadCrewMeta(r.project, crew)
	if err != nil {
		r.volatile = true
		return false
	}
	closed := crewstate.Declare(crewstate.Declaration{Meta: meta}).Closed()
	r.next.closed[path] = closedMemo{at: now, closed: closed}
	return closed
}

// status is the i-th crew's status lines from the start, and the file's
// info (nil when it is not there).
func (r *loadReader) status(i int) ([]store.StatusEntry, int64, os.FileInfo, error) {
	crew, path := r.next.crews.crews[i], r.next.crews.status[i]
	now, fi, ok := statStamp(path)
	if !ok {
		r.moved, r.volatile = true, true
		lines, next, err := r.ws.ReadStatus(r.project, crew, 0)
		return lines, next, nil, err
	}
	if r.held.status[path].at != now {
		r.moved = true
	}
	m, err := r.held.status[path].read(now, func(from int64) ([]store.StatusEntry, int64, error) {
		return r.ws.ReadStatus(r.project, crew, from)
	})
	if err != nil {
		return nil, 0, nil, err
	}
	r.next.status[path] = m
	return m.lines, m.next, fi, nil
}

// sent is the project's sent.log from the start.
func (r *loadReader) sent() ([]store.SentEntry, int64, error) {
	now, _, ok := statStamp(r.ws.SentLog(r.project))
	if !ok {
		if r.held.sent.at != (stamp{}) {
			r.moved = true
		}
		return r.ws.ReadSent(r.project, 0)
	}
	if r.held.sent.at != now {
		r.moved = true
	}
	m, err := r.held.sent.read(now, func(from int64) ([]store.SentEntry, int64, error) {
		return r.ws.ReadSent(r.project, from)
	})
	if err != nil {
		return nil, 0, err
	}
	r.next.sent = m
	return m.lines, m.next, nil
}

// incidents is the project's incidents.log from the start.
func (r *loadReader) incidents() ([]store.IncidentEntry, int64, error) {
	now, _, ok := statStamp(r.ws.IncidentsLog(r.project))
	if !ok {
		if r.held.incidents.at != (stamp{}) {
			r.moved = true
		}
		return r.ws.ReadIncidents(r.project, 0)
	}
	if r.held.incidents.at != now {
		r.moved = true
	}
	m, err := r.held.incidents.read(now, func(from int64) ([]store.IncidentEntry, int64, error) {
		return r.ws.ReadIncidents(r.project, from)
	})
	if err != nil {
		return nil, 0, err
	}
	r.next.incidents = m
	return m.lines, m.next, nil
}

// heldView is the View the last Load built, when nothing this Load read
// has moved since: the same files at the same stamps, no file the last one
// read missing now, and the same mode.
func (r *loadReader) heldView(manual bool) (View, bool) {
	held := r.held.view
	if r.moved || r.volatile || held == nil || held.Manual != manual ||
		len(r.next.status) != len(r.held.status) || len(r.next.closed) != len(r.held.closed) {
		return View{}, false
	}
	return *held, true
}

// done keeps what this Load read, and the View it answered with, for the
// next one.
func (r *loadReader) done(v View) {
	if !r.volatile {
		r.next.view = &v
	}
	memo.keep(r.ws.CrewsDir(r.project), r.next)
}
