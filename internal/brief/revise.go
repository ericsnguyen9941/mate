package brief

import (
	"fmt"
	"strings"
)

// ReviseTask replaces only the task specification in a rendered brief.
// The revision is a complete task, not a rendered worker prompt. Captain's
// words must match the current record; new words use AppendCaptainsWords.
// This keeps a stale draft from deleting later captain instructions.
func ReviseTask(current, revision string, kind Kind) (string, error) {
	if problems := Check(revision, kind); len(problems) != 0 {
		return "", &Error{Kind: kind, Problems: problems}
	}
	before, ok := SectionText(current, CaptainsWords)
	after, _ := SectionText(revision, CaptainsWords)
	if !ok || before != after {
		return "", fmt.Errorf("revision must preserve ## %s exactly; copy the current task and use brief append for new captain words", CaptainsWords)
	}
	// Use the same fence-aware boundaries as ExtractTask. Keep the role,
	// setup, worker rules and standing instructions outside # Task intact.
	lines := parseLines(current)
	start, end := -1, len(lines)
	for i, line := range lines {
		if !line.level1 {
			continue
		}
		if start >= 0 {
			end = i
			break
		}
		if line.trimmed == TaskHeading {
			start = i + 1
		}
	}
	if start < 0 {
		return "", fmt.Errorf("current brief has no %s", TaskHeading)
	}
	// Split the original, rather than parsed text, to preserve bytes outside
	// the task (including CRLF in an older brief).
	raw := strings.Split(current, "\n")
	parts := append([]string{}, raw[:start]...)
	parts = append(parts, "", TaskBody(revision), "")
	parts = append(parts, raw[end:]...)
	return strings.Join(parts, "\n"), nil
}
