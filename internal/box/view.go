package box

import (
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Cursor is a set of per-file read offsets: one for sent.log, one for
// incidents.log, and one per crew status file. LoadSince uses it to read
// only what is new since a previous Load or LoadSince; a caller (the
// observer) keeps the Cursor a View returns and passes it back on the next
// poll.
type Cursor struct {
	Sent int64
	// Incidents is the resume offset of incidents.log. Only the entries
	// after it are returned, but the whole file is always read: whether an
	// incident is still open is decided by the last line of its (crew,
	// kind) pair, which can sit on either side of the cursor.
	Incidents int64
	Status    map[string]int64
}

// View is a project's communication, merged and sorted by time.
type View struct {
	// Entries is every message, status line and incident, sorted by time
	// then by (file, offset) - see the package doc for why status lines can
	// only be ordered that coarsely.
	Entries []Entry
	// ByCrew groups Entries by Crew, in the same relative order, dropping
	// entries with no crew (a message to Mate's own pane, for instance).
	ByCrew map[string][]Entry
	// Attention is the subset that needs a human or Mate's eyes: a status
	// entry whose state is Attention, or any incident. Newest first.
	Attention []Entry
	// Closed is the set of crews whose meta declares a terminal state
	// (`state=finished|failed`, or a pre-4b `stopped_at` with no state):
	// the Mate or the captain closed them with `mate crew stop`, or their
	// spawn failed. Their lines stay in Entries as history; Inbox skips
	// them, because nobody can answer a crew that is gone.
	Closed map[string]bool
	// Manual is true when the project's `.auto` flag was absent at load: the
	// captain chose manual mode with the console's `m` key, so no digest
	// will tell the Mate anything and Inbox keeps a crew's latest
	// `wait-mate` for the captain (docs/mvp.md M18).
	Manual bool
	// Cursor resumes a later LoadSince from exactly where this View ended.
	Cursor Cursor

	gen uint64
}

// Summary is a header digest of a View.
type Summary struct {
	// Crews is the number of distinct crews present in the view.
	Crews int
	// Awaiting is the number of crews whose latest status entry is an
	// Attention state.
	Awaiting int
	// LastAt is the time of the most recent entry, zero if there is none.
	LastAt time.Time
}

// Load reads every crew status file, sent.log and incidents.log for project
// from the start, and returns the full View.
func Load(ws Workspace, project string) (View, error) {
	return load(ws, project, Cursor{})
}

// LoadSince reads only what is new since a previous Load or LoadSince, using
// since's per-file offsets. The returned View's Entries hold only the new
// lines - a caller wanting the full history keeps its own accumulation
// across calls.
func LoadSince(ws Workspace, project string, since Cursor) (View, error) {
	return load(ws, project, since)
}

func load(ws Workspace, project string, cursor Cursor) (View, error) {
	// A Load reads through the memo (memo.go); a LoadSince reads from its
	// own cursor.
	var memoized *loadReader
	if cursor.Sent == 0 && cursor.Incidents == 0 && len(cursor.Status) == 0 {
		memoized = newLoadReader(ws, project)
	}

	var (
		crews []string
		err   error
	)
	if memoized != nil {
		crews, err = memoized.listCrews()
	} else {
		crews, err = listCrews(ws, project)
	}
	if err != nil {
		return View{}, err
	}

	newCursor := Cursor{Status: make(map[string]int64, len(crews))}
	closed := make(map[string]bool)
	type statusFile struct {
		crew, path string
		mtime      time.Time
		lines      []store.StatusEntry
	}
	var statuses []statusFile
	n := 0

	for i, crew := range crews {
		var (
			lines      []store.StatusEntry
			next       int64
			info       os.FileInfo
			statusPath string
		)
		if memoized != nil {
			if memoized.crewClosed(i) {
				closed[crew] = true
			}
			lines, next, info, err = memoized.status(i)
			statusPath = memoized.next.crews.status[i]
		} else {
			if meta, err := ws.ReadCrewMeta(project, crew); err == nil && crewstate.Declare(crewstate.Declaration{Meta: meta}).Closed() {
				closed[crew] = true
			}
			from := int64(0)
			if cursor.Status != nil {
				from = cursor.Status[crew]
			}
			lines, next, err = ws.ReadStatus(project, crew, from)
		}
		if err != nil {
			return View{}, err
		}
		newCursor.Status[crew] = next
		if len(lines) == 0 {
			continue
		}
		if statusPath == "" {
			statusPath = ws.CrewStatus(project, crew)
		}
		mtime := time.Now()
		if info == nil {
			info, _ = os.Stat(statusPath)
		}
		if info != nil {
			mtime = info.ModTime()
		}
		statuses = append(statuses, statusFile{crew: crew, path: statusPath, mtime: mtime, lines: lines})
		n += len(lines)
	}

	var (
		sentLines []store.SentEntry
		sentNext  int64
	)
	if memoized != nil {
		sentLines, sentNext, err = memoized.sent()
	} else {
		sentLines, sentNext, err = ws.ReadSent(project, cursor.Sent)
	}
	if err != nil {
		return View{}, err
	}
	newCursor.Sent = sentNext

	var (
		incidentLines     []store.IncidentEntry
		incidentEntryList []Entry
	)
	if memoized != nil {
		incidentLines, newCursor.Incidents, err = memoized.incidents()
	} else {
		incidentEntryList, newCursor.Incidents, err = loadIncidents(ws, project, cursor.Incidents)
	}
	if err != nil {
		return View{}, err
	}

	manual := !ws.Auto(project)
	if memoized != nil {
		if v, ok := memoized.heldView(manual); ok {
			memoized.done(v)
			return v, nil
		}
		incidentEntryList = incidentEntries(ws.IncidentsLog(project), incidentLines, 0)
	}

	entries := make([]Entry, 0, n+len(sentLines)+len(incidentEntryList))
	for _, f := range statuses {
		target := store.CrewTarget(f.crew)
		for _, l := range f.lines {
			entries = append(entries, Entry{
				At:     f.mtime,
				Source: SourceCrew,
				Target: target,
				Kind:   KindStatus,
				Crew:   f.crew,
				Text:   l.Line,
				Ref:    Ref{File: f.path, Offset: l.Offset},
			})
		}
	}
	sentPath := ws.SentLog(project)
	for _, s := range sentLines {
		entries = append(entries, Entry{
			At:     s.Time,
			Source: Source(s.Source),
			Target: s.Target,
			Kind:   KindMessage,
			Crew:   crewFromTarget(s.Target),
			Text:   s.Text,
			Ref:    Ref{File: sentPath, Offset: s.Offset},
		})
	}
	entries = append(entries, incidentEntryList...)

	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if !a.At.Equal(b.At) {
			return a.At.Before(b.At)
		}
		if a.Ref.File != b.Ref.File {
			return a.Ref.File < b.Ref.File
		}
		return a.Ref.Offset < b.Ref.Offset
	})
	for i := range entries {
		entries[i].Seq = i
	}

	byCrew := make(map[string][]Entry)
	for _, e := range entries {
		if e.Crew == "" {
			continue
		}
		byCrew[e.Crew] = append(byCrew[e.Crew], e)
	}

	var attention []Entry
	for _, e := range entries {
		if e.Kind == KindIncident || (e.Kind == KindStatus && Attention(ParseStatus(e.Text).State)) {
			attention = append(attention, e)
		}
	}
	for i, j := 0, len(attention)-1; i < j; i, j = i+1, j-1 {
		attention[i], attention[j] = attention[j], attention[i]
	}

	v := View{Entries: entries, ByCrew: byCrew, Attention: attention, Closed: closed, Manual: manual, Cursor: newCursor}
	if memoized != nil {
		// The View is shared from here on (memo.go): no slice may have
		// room an append would write into.
		v.Entries = slices.Clip(v.Entries)
		v.Attention = slices.Clip(v.Attention)
		for crew, list := range byCrew {
			byCrew[crew] = slices.Clip(list)
		}
		v.gen = generation.Add(1)
		memoized.done(v)
	}
	return v, nil
}

// Generation names the View a Load built: two Views with the same non-zero
// Generation are the same one, built from the same files at the same size
// and mtime, so a reader that derived something from the first may keep it
// for the second (CON-271). A View from LoadSince has none (0).
func (v View) Generation() uint64 { return v.gen }

// listCrews returns the crew ids of a project, from the `<id>.status`
// filenames under `crews/`. A project with no crews dir yet has none.
func listCrews(ws Workspace, project string) ([]string, error) {
	dir := ws.CrewsDir(project)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var crews []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if id, ok := strings.CutSuffix(name, ".status"); ok {
			crews = append(crews, id)
		}
	}
	sort.Strings(crews)
	return crews, nil
}

// crewFromTarget extracts a crew id from a sent.log target ("crew:<id>"),
// or "" for the mate target - a message to Mate is not tied to a crew.
func crewFromTarget(target string) string {
	if id, ok := strings.CutPrefix(target, "crew:"); ok {
		return id
	}
	return ""
}

// incidentText packs an incident's kind into Entry.Text the same way a
// status line packs its verb, so Entry needs no incident-specific field:
// "<kind>: <text>". ParseIncidentText reverses it.
func incidentText(kind IncidentKind, text string) string {
	return string(kind) + ": " + text
}

// ParseIncidentText reverses incidentText, splitting a KindIncident Entry's
// Text back into its IncidentKind and message.
func ParseIncidentText(raw string) (IncidentKind, string) {
	k, t, ok := strings.Cut(raw, ": ")
	if !ok {
		return "", raw
	}
	return IncidentKind(k), t
}

// LatestStatus returns the most recent status entry of a crew, parsed, and
// whether the crew has any status at all.
func LatestStatus(v View, crew string) (Status, bool) {
	entries := v.ByCrew[crew]
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Kind == KindStatus {
			return ParseStatus(entries[i].Text), true
		}
	}
	return Status{}, false
}

// AtWork reports whether crew is open and still at work: its latest status
// line is not a hand-back. A crew with no status line yet is at work, since it
// was just spawned. The digest uses it to decide whether a hand-back has
// anybody to wait for (CON-272), and the pull request watcher whether a
// landing leaves the Mate with nothing running (CON-194).
func AtWork(v View, crew string) bool {
	if v.Closed[crew] {
		return false
	}
	st, ok := LatestStatus(v, crew)
	return !ok || st.State != StateWaitMate
}

// LastCrewVerb is LastVerb over the status lines v holds for crew. A crew's
// status lines share one time (their file's mtime) and one file, so the
// merge keeps them in file order and the answer is the one LastVerb gives
// over the file itself. On a View from Load that is the whole file; on one
// from LoadSince it is only what was new.
func LastCrewVerb(v View, crew string) State {
	entries := v.ByCrew[crew]
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		if e.Kind != KindStatus || strings.TrimSpace(e.Text) == "" {
			continue
		}
		if st := ParseStatus(e.Text); st.State != StateUnknown {
			return st.State
		}
	}
	return StateUnknown
}

// Summarize builds a header digest of a View.
func Summarize(v View) Summary {
	var s Summary
	s.Crews = len(v.ByCrew)
	for crew := range v.ByCrew {
		if st, ok := LatestStatus(v, crew); ok && Attention(st.State) && !Answered(v, crew) {
			s.Awaiting++
		}
	}
	if n := len(v.Entries); n > 0 {
		s.LastAt = v.Entries[n-1].At
	}
	return s
}
