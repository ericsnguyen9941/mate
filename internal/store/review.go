package store

import (
	"path/filepath"
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
	_, tail, ok := strings.Cut("\n"+report, "\n## Verdict\n")
	if !ok {
		return ""
	}
	return strings.TrimSpace(strings.SplitN(strings.TrimSpace(tail), "\n", 2)[0])
}
