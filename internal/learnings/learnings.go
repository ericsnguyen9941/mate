// Package learnings reads the predictions a Mate's learning records make
// and closes them with evidence (CON-131).
//
// A learning record (store.CrewLearning) is the Mate's prose, one per
// crew, written by the task-learning skill. Its `## Predictions` section
// is the part a later task can test: one bullet per prediction, ending in
// `outcome: open` until a comparable task supplies evidence:
//
//   - p1: a demo launch argument saves the screenshot round; next check HU-53 — outcome: open
//   - p2: router-level acceptance avoids the UI-test dispute; next check the next integration brief — outcome: improved (HU-114 r1 pass, no dispute, 2026-10-28)
//
// On mihuv2 (2026-10-27/28) thirteen records held predictions and none was
// ever closed: a lesson marked `captured` was never checked. `mate
// learnings open` is what new-task reads to carry a prediction into the
// next brief, and `mate learnings close` is how task-learning records what
// that task showed, rewriting only the outcome of the one bullet.
package learnings

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// OutcomeOpen is a prediction no task has tested yet.
const OutcomeOpen = "open"

// Outcomes are the words a prediction is closed with: the later task's
// evidence showed the predicted saving, showed none, could not tell, or was
// not a comparable task (scope, size, model, environment or coverage).
var Outcomes = []string{"improved", "unchanged", "inconclusive", "not comparable"}

// Prediction is one bullet of a record's `## Predictions` section.
type Prediction struct {
	// Crew is the record it lives in, Path the record's file.
	Crew string
	Path string
	// ID is the item's number, `1` for `1.`; a plain bullet is numbered by
	// its position, so every prediction has an id to close it by. It is
	// printed and accepted as `#1` too.
	ID string
	// Text is the bullet between the id and `outcome:`, the claim with its
	// baseline and next check as the Mate wrote them.
	Text string
	// Outcome is the word after `outcome:`; Evidence what the closing task
	// put in parentheses after it, "" while open.
	Outcome  string
	Evidence string
}

// String is the line `mate learnings open` prints for a prediction.
func (p Prediction) String() string {
	s := fmt.Sprintf("%s #%s: %s", p.Crew, p.ID, p.Text)
	if p.Outcome != OutcomeOpen {
		s += " — outcome: " + p.Outcome
		if p.Evidence != "" {
			s += " (" + p.Evidence + ")"
		}
	}
	return s
}

// bullet is a prediction line: a numbered or plain list marker, the text,
// an optional separator, then `outcome:` and the rest of the line.
var bullet = regexp.MustCompile(`^(?:[-*]|(\d+)[.)])\s+(.*?)\s*(?:[-—–|;·]+\s*)?outcome:\s*(.*?)\s*$`)

// Open is every open prediction of a project's records, records in name
// order, bullets in file order. A record without the section has none.
func Open(w *store.Workspace, project string) ([]Prediction, error) {
	all, err := List(w, project)
	if err != nil {
		return nil, err
	}
	var open []Prediction
	for _, p := range all {
		if p.Outcome == OutcomeOpen {
			open = append(open, p)
		}
	}
	return open, nil
}

// List is every prediction of a project's records, open or closed.
func List(w *store.Workspace, project string) ([]Prediction, error) {
	entries, err := os.ReadDir(w.LearningsDir(project))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			names = append(names, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	sort.Strings(names)
	var out []Prediction
	for _, crew := range names {
		ps, err := Read(w, project, crew)
		if err != nil {
			return nil, err
		}
		out = append(out, ps...)
	}
	return out, nil
}

// Read is one record's predictions.
func Read(w *store.Workspace, project, crew string) ([]Prediction, error) {
	path := w.CrewLearning(project, crew)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []Prediction
	for _, l := range sectionLines(string(data)) {
		p, ok := parseBullet(l.text)
		if !ok {
			continue
		}
		p.Crew, p.Path = crew, path
		if p.ID == "" {
			p.ID = strconv.Itoa(len(out) + 1)
		}
		out = append(out, p)
	}
	return out, nil
}

// ErrUnknownPrediction is Close on an id the record has no bullet for.
var ErrUnknownPrediction = errors.New("learnings: no such prediction")

// ErrBadOutcome is Close with a word that is not one of Outcomes.
var ErrBadOutcome = errors.New("learnings: outcome must be one of " + strings.Join(Outcomes, ", "))

// Close rewrites one prediction's outcome in its record, keeping the bullet
// up to `outcome:` verbatim: `outcome: <outcome> (<evidence>, <date>)`. A
// closed prediction may be closed again when later evidence changes the
// verdict; the previous one is overwritten, the record's history is git's.
func Close(w *store.Workspace, project, crew, id, outcome, evidence string, today time.Time) (Prediction, error) {
	outcome = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(outcome, "-", " ")))
	known := false
	for _, o := range Outcomes {
		known = known || o == outcome
	}
	if !known {
		return Prediction{}, fmt.Errorf("%w; got %q", ErrBadOutcome, outcome)
	}
	id = strings.TrimLeft(strings.TrimSpace(id), "#")
	evidence = strings.TrimSpace(strings.ReplaceAll(evidence, "\n", " "))
	if evidence == "" {
		return Prediction{}, errors.New("learnings: --evidence names the task, revision and what it showed")
	}
	path := w.CrewLearning(project, crew)
	data, err := os.ReadFile(path)
	if err != nil {
		return Prediction{}, err
	}
	lines := strings.Split(string(data), "\n")
	n := 0
	for _, l := range sectionLines(string(data)) {
		p, ok := parseBullet(l.text)
		if !ok {
			continue
		}
		n++
		if p.ID == "" {
			p.ID = strconv.Itoa(n)
		}
		if p.ID != id {
			continue
		}
		at := strings.Index(lines[l.index], "outcome:")
		p.Outcome, p.Evidence = outcome, evidence+", "+today.Format("2006-01-02")
		lines[l.index] = lines[l.index][:at] + "outcome: " + p.Outcome + " (" + p.Evidence + ")"
		if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
			return Prediction{}, err
		}
		p.Crew, p.Path = crew, path
		return p, nil
	}
	return Prediction{}, fmt.Errorf("%w: %s #%s in %s", ErrUnknownPrediction, crew, id, path)
}

type sectionLine struct {
	index int
	text  string
}

// sectionLines are the lines of the record's `## Predictions` section,
// with their index in the file, up to the next heading.
func sectionLines(text string) []sectionLine {
	heading := "## " + store.LearningPredictionsSection
	var out []sectionLine
	in := false
	for i, l := range strings.Split(text, "\n") {
		if strings.HasPrefix(l, "## ") || strings.HasPrefix(l, "# ") {
			if in {
				break
			}
			in = strings.TrimSpace(l) == heading
			continue
		}
		if in {
			out = append(out, sectionLine{i, l})
		}
	}
	return out
}

func parseBullet(line string) (Prediction, bool) {
	m := bullet.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return Prediction{}, false
	}
	p := Prediction{ID: m[1], Text: strings.TrimSpace(m[2])}
	rest := strings.TrimSpace(m[3])
	if i := strings.Index(rest, " ("); i > 0 && strings.HasSuffix(rest, ")") {
		p.Outcome, p.Evidence = rest[:i], rest[i+2:len(rest)-1]
	} else {
		p.Outcome = rest
	}
	p.Outcome = strings.ToLower(strings.TrimSpace(p.Outcome))
	if p.Outcome == "" {
		p.Outcome = OutcomeOpen
	}
	return p, true
}
