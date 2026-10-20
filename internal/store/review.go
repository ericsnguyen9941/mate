package store

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// A review (`mate review`, docs/mvp.md) is a scout Crew whose brief asks for
// a `## Verdict` section in its report. The command records what it was
// asked to review in review.meta beside that report, so the file is also
// what tells a reader the Crew is a reviewer.

const reviewMetaFileName = "review.meta"

// CrewReviewMeta is `projects/<project>/crews/<crew>/review.meta`, present
// only for a reviewer Crew.
func (w *Workspace) CrewReviewMeta(project, crew string) string {
	return filepath.Join(w.CrewDir(project, crew), reviewMetaFileName)
}

// ReviewVerdict is the first line of a report's `## Verdict` section - by
// the review brief one of pass, changes-requested or incomplete - or ""
// when the report has no such section.
func ReviewVerdict(report string) string {
	_, tail, ok := strings.Cut("\n"+report, "\n## "+VerdictSection+"\n")
	if !ok {
		return ""
	}
	return strings.TrimSpace(strings.SplitN(strings.TrimSpace(tail), "\n", 2)[0])
}

// A review report lists its material findings under `## Findings`, one
// top-level list item per finding (or the word none); what does not block
// goes under `## Notes`, which nothing here reads. A scoped review (`mate
// review --since`) answers the earlier review's findings under `## Prior
// findings`: one list item per finding, starting with `F<n> closed` or
// `F<n> open`, where n is the finding's position in the list the brief
// gave, followed by the evidence. Both sections are what `mate review
// --since` reads to build the next brief and what `--check` reads to
// accept a scoped review, so the brief teaches them verbatim.

const (
	VerdictSection       = "Verdict"
	FindingsSection      = "Findings"
	PriorFindingsSection = "Prior findings"
	FindingClosed        = "closed"
	FindingOpen          = "open"
)

var priorFindingPattern = regexp.MustCompile(`(?i)^F(\d+)\s*[:\-–—]?\s*(closed|open)\b\s*[:\-–—]?\s*`)

// ReviewFindings are the material findings a report leaves standing, in
// the order a later brief numbers them: every item under `## Findings`,
// then every prior finding the report marks open (its text after the
// status word). A report with neither section has no findings.
func ReviewFindings(report string) []string {
	var out []string
	for _, it := range reportItems(report, FindingsSection) {
		if !strings.EqualFold(it, "none") {
			out = append(out, it)
		}
	}
	for _, it := range reportItems(report, PriorFindingsSection) {
		if m := priorFindingPattern.FindStringSubmatch(it); m != nil && strings.EqualFold(m[2], FindingOpen) {
			if text := strings.TrimSpace(it[len(m[0]):]); text != "" {
				out = append(out, text)
			} else {
				out = append(out, it)
			}
		}
	}
	return out
}

// ReviewPriorFindings is each prior finding's status as the report marks it
// under `## Prior findings`, keyed by the finding's number: "closed" or
// "open". Items that do not follow the convention are not in the map.
func ReviewPriorFindings(report string) map[int]string {
	out := map[int]string{}
	for _, it := range reportItems(report, PriorFindingsSection) {
		m := priorFindingPattern.FindStringSubmatch(it)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		out[n] = strings.ToLower(m[2])
	}
	return out
}

// reportItems are the top-level list items of the report's `## <name>`
// section, each joined with its indented continuation lines, fences
// skipped. A list item indented deeper than the section's first one is a
// sub-bullet of the item above it, not an item of its own (CON-106): a
// finding with three supporting bullets is one finding. A section whose
// body has no list marker is one item.
func reportItems(report, name string) []string {
	body, ok := reportSection(report, name)
	if !ok {
		return nil
	}
	var items []string
	fenced := false
	top := -1
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") || strings.HasPrefix(strings.TrimSpace(line), "~~~") {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		text, isItem := stripListMarker(trimmed)
		if isItem && top < 0 {
			top = indentWidth(line)
		}
		if isItem && indentWidth(line) <= top {
			items = append(items, text)
		} else if len(items) > 0 {
			items[len(items)-1] += " " + trimmed
		} else {
			items = append(items, trimmed)
		}
	}
	return items
}

// indentWidth is how far line is indented, a tab counting as four spaces.
func indentWidth(line string) int {
	n := 0
	for _, r := range line {
		switch r {
		case ' ':
			n++
		case '\t':
			n += 4
		default:
			return n
		}
	}
	return n
}

var listMarkerPattern = regexp.MustCompile(`^(?:[-*+]|\d+[.)])\s+`)

func stripListMarker(line string) (string, bool) {
	loc := listMarkerPattern.FindStringIndex(line)
	if loc == nil {
		return line, false
	}
	return strings.TrimSpace(line[loc[1]:]), true
}

// reportSection is the text under the report's `## <name>` heading up to
// the next `#` heading outside a fence.
func reportSection(report, name string) (string, bool) {
	var lines []string
	found, fenced := false, false
	for _, line := range strings.Split(report, "\n") {
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			fenced = !fenced
		}
		if !fenced && strings.EqualFold(strings.TrimSpace(line), "## "+name) {
			if found {
				break
			}
			found = true
			continue
		}
		if found && !fenced && (strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "# ")) {
			break
		}
		if found {
			lines = append(lines, line)
		}
	}
	if !found {
		return "", false
	}
	return strings.Join(lines, "\n"), true
}

// MetaReviews is the implementing crew's `.meta` key naming the reviewers
// `mate review` started for it, comma-separated and in the order they were
// started (CON-75). It is what lets a ledger roll a task's reviewers into
// the task: before it, "Review <crew> at <sha>" in the reviewer's task text
// was the only link, and summing them was hand work.
const MetaReviews = "reviews"

// AddCrewReview records reviewer on crew's `.meta` under MetaReviews,
// once. The crew's meta must exist; a crew that was already stopped keeps
// its reviewers (clearCrewRunMeta keeps the key).
func (w *Workspace) AddCrewReview(project, crew, reviewer string) error {
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil {
		return err
	}
	ids := splitReviews(meta[MetaReviews])
	for _, id := range ids {
		if id == reviewer {
			return nil
		}
	}
	ids = append(ids, reviewer)
	return w.UpdateCrewMeta(project, crew, map[string]string{MetaReviews: strings.Join(ids, ",")})
}

func splitReviews(value string) []string {
	var out []string
	for _, id := range strings.Split(value, ",") {
		if id = strings.TrimSpace(id); id != "" {
			out = append(out, id)
		}
	}
	return out
}

// ReviewLinks is every task→reviewer link of a project, from both records
// that carry one: the implementing crew's `reviews=` and the reviewer's own
// review.meta (`crew=`), which reviewers started before `reviews=` existed
// still have. reviewersOf maps a crew to its reviewers, sorted; reviewOf
// maps a reviewer to the crew it reviewed.
func (w *Workspace) ReviewLinks(project string) (reviewersOf map[string][]string, reviewOf map[string]string, err error) {
	ids, err := w.CrewIDs(project)
	if err != nil {
		return nil, nil, err
	}
	reviewersOf, reviewOf = map[string][]string{}, map[string]string{}
	link := func(crew, reviewer string) {
		if crew == "" || reviewer == "" || crew == reviewer {
			return
		}
		for _, have := range reviewersOf[crew] {
			if have == reviewer {
				return
			}
		}
		reviewersOf[crew] = append(reviewersOf[crew], reviewer)
		if _, ok := reviewOf[reviewer]; !ok {
			reviewOf[reviewer] = crew
		}
	}
	for _, id := range ids {
		meta, err := w.ReadCrewMeta(project, id)
		if err != nil {
			return nil, nil, err
		}
		for _, reviewer := range splitReviews(meta[MetaReviews]) {
			link(id, reviewer)
		}
		if _, statErr := os.Stat(w.CrewReviewMeta(project, id)); statErr == nil {
			if rec, err := ReadMeta(w.CrewReviewMeta(project, id)); err == nil {
				link(strings.TrimSpace(rec["crew"]), id)
			}
		}
	}
	for crew := range reviewersOf {
		sort.Strings(reviewersOf[crew])
	}
	return reviewersOf, reviewOf, nil
}
