package store

import (
	"path/filepath"
	"strings"
	"sync/atomic"
)

// The change feed's writer end (CON-260). Almost every change under
// `.mate/` is written by the mate binary itself - a CLI command, `mate
// hook`, a Mate's or a Crew's mate command, a headless host - yet the core
// was never told and had to poll for all of it. Every write this package
// makes ends in NoteWrite, which tells the process's ChangeNotifier which
// project, crew and kind of record the path is; cmd/mate sets one that
// sends the notice to the workspace's core over `.mate/core.sock`
// (coreclient.Notify). The store imports nothing above it, so the sending
// is injected: with no notifier set - every test, and the core itself,
// which would only be telling itself - a write tells nobody.
//
// A notice is a hint, never the record: the file is. One lost on the way
// (no core, a core restarting, a write that bypasses this package -
// docs/architecture.md section 2 lists them) is caught by the core's polls
// as before.

// ChangeKind is what kind of record a write changed, read from its path.
type ChangeKind string

const (
	// ChangeWorkspace is `workspace.yaml` and the workspace's own docs and
	// prices: the Projects and their repos.
	ChangeWorkspace ChangeKind = "workspace"
	// ChangeProject is a Project's `project.yaml`, `PROJECT.md` or
	// `CREW.md`, or any file of its directory no other kind names.
	ChangeProject ChangeKind = "project"
	// ChangeSent is a line of the Project's `sent.log`.
	ChangeSent ChangeKind = "sent"
	// ChangeIncident is a line of the Project's `incidents.log`.
	ChangeIncident ChangeKind = "incident"
	// ChangeMate is a file of `mate/` no kind below names: `mate.meta`, the
	// mode flags, the cursors, memory and backlog, learnings.
	ChangeMate ChangeKind = "mate"
	// ChangeOutbox is the Mate's `.outbox`.
	ChangeOutbox ChangeKind = "outbox"
	// ChangeNews is the Mate's `news.jsonl` (internal/matehost).
	ChangeNews ChangeKind = "news"
	// ChangeWork is the Mate's `work.jsonl` (internal/matehost).
	ChangeWork ChangeKind = "work"
	// ChangeCrewMeta is a Crew's `.meta`.
	ChangeCrewMeta ChangeKind = "meta"
	// ChangeCrewStatus is a line of a Crew's `.status` that mate wrote
	// (prwatch, merge). A Crew's own `echo >> $MATE_STATUS` bypasses mate
	// and sends nothing.
	ChangeCrewStatus ChangeKind = "status"
	// ChangeCrew is any other file of a Crew: its brief, report, hand-back,
	// review record, session inputs.
	ChangeCrew ChangeKind = "crew"
	// ChangeTranscript is a Crew's or the Mate's harness transcript,
	// outside `.mate/`: no writer notes it, only the core's watch on the
	// files written past the feed (CON-266) does. ChangeOf never reads it.
	ChangeTranscript ChangeKind = "transcript"
)

// Change is one write as the core's change feed hears it.
type Change struct {
	// Socket is the workspace's `.mate/core.sock`, where its core listens.
	Socket string
	// Project is the Project the path is under; empty for ChangeWorkspace.
	Project string
	// Crew is the Crew the path is of; empty but for the crew kinds.
	Crew string
	Kind ChangeKind
}

// ChangeNotifier is told of every write NoteWrite reads a Change from. It
// runs on the writer's goroutine once the write is done, the file lock
// released: it must return at once whatever the core does, and fail
// silently.
type ChangeNotifier func(Change)

var notifier atomic.Pointer[ChangeNotifier]

// SetChangeNotifier makes n the process's notifier; nil sets none. cmd/mate
// sets it once, before any command runs.
func SetChangeNotifier(n ChangeNotifier) {
	if n == nil {
		notifier.Store(nil)
		return
	}
	notifier.Store(&n)
}

// NoteWrite tells the notifier that path, a file under a workspace's
// `.mate/`, was just written or removed. A path that is no record of the
// feed's (a lock, a log only the captain reads, anything outside `.mate/`)
// tells nothing. The store's own writers call it; so does any package that
// writes a record under `.mate/` itself (internal/matehost's news and work).
func NoteWrite(path string) {
	n := notifier.Load()
	if n == nil {
		return
	}
	if c, ok := ChangeOf(path); ok {
		(*n)(c)
	}
}

// ChangeOf reads the Change a write of path is: the workspace from the last
// `.mate` element of path, then the Project, the Crew and the kind from the
// rest, by the layout of layout*.go. ok is false for a path that is no
// record of the feed's.
func ChangeOf(path string) (Change, bool) {
	path = filepath.Clean(path)
	sep := string(filepath.Separator)
	i := strings.LastIndex(path, sep+StateDirName+sep)
	if i < 0 {
		return Change{}, false
	}
	state := path[:i+len(sep+StateDirName)]
	c := Change{Socket: filepath.Join(state, coreSocketName)}
	rest := strings.Split(path[len(state)+1:], sep)
	name := rest[len(rest)-1]
	if strings.HasSuffix(name, ".lock") || strings.Contains(name, ".tmp") {
		return Change{}, false
	}
	if len(rest) == 1 {
		switch name {
		case workspaceFileName, workspaceDocName, pricingFileName, crewDocName:
			c.Kind = ChangeWorkspace
			return c, true
		}
		// jev.log, migrate.log, gc.log, the locks: nothing a loop reads.
		return Change{}, false
	}
	if rest[0] != projectsDirName || len(rest) < 3 {
		return Change{}, false
	}
	c.Project = rest[1]
	switch {
	case len(rest) == 3 && name == sentLogName:
		c.Kind = ChangeSent
	case len(rest) == 3 && name == incidentsLogName:
		c.Kind = ChangeIncident
	case len(rest) == 3:
		c.Kind = ChangeProject
	case rest[2] == mateDirName:
		switch {
		case len(rest) == 4 && name == outboxName:
			c.Kind = ChangeOutbox
		case len(rest) == 4 && name == "news.jsonl":
			c.Kind = ChangeNews
		case len(rest) == 4 && name == "work.jsonl":
			c.Kind = ChangeWork
		default:
			c.Kind = ChangeMate
		}
	case rest[2] == crewsDirName && len(rest) == 4:
		// A crew id has no dot (names.CrewPattern): `<c>.meta`,
		// `<c>.status`, `<c>.send.json`, `<c>.prwatch`.
		crew, ext, _ := strings.Cut(name, ".")
		c.Crew = crew
		switch ext {
		case "meta":
			c.Kind = ChangeCrewMeta
		case "status":
			c.Kind = ChangeCrewStatus
		default:
			c.Kind = ChangeCrew
		}
	case rest[2] == crewsDirName:
		c.Crew, c.Kind = rest[3], ChangeCrew
	default:
		c.Kind = ChangeProject
	}
	return c, true
}
