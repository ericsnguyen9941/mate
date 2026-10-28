package brief_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/brief"
	"github.com/nguyenngocanh94/mate/internal/brief/brieftest"
)

func TestReviseTaskPreservesContractAndWorkerInstructions(t *testing.T) {
	task, err := brief.AppendCaptainsWords(brieftest.Ship("Show screen states."), "Keep retry manual.", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	task = strings.Replace(task, "## Open decisions\nnone", "## Open decisions\n- Include loading when no user path reaches it? decides: captain", 1)
	prefix := brief.RoleHeading + "\r\nWorker identity.\r\n# Task\r\n"
	suffix := "# Setup\r\nKeep isolation.\r\n# Captain's standing crew rules\r\nKeep these.\r\n"
	current := prefix + task + "\n" + suffix
	revision := strings.Replace(task, "## Build\nShow screen states.", "## Build\nShow empty and error states; keep retry manual.", 1)
	revision = strings.Replace(revision, "- The words above are done exactly as written. verify: read your own pane and the status file.", "- Only reachable states render. verify: exercise the empty and error paths.", 1)
	revision = strings.Replace(revision, "- Include loading when no user path reaches it? decides: captain", "none", 1)
	got, err := brief.ReviseTask(current, revision, brief.Ship)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, prefix) || !strings.HasSuffix(got, suffix) {
		t.Fatalf("worker instructions changed:\n%s", got)
	}
	wantWords, _ := brief.SectionText(current, brief.CaptainsWords)
	gotWords, _ := brief.SectionText(got, brief.CaptainsWords)
	decisions, _ := brief.SectionText(got, brief.OpenDecisions)
	if decisions != "none" || !strings.Contains(got, "Only reachable states render.") {
		t.Fatalf("acceptance and decisions were not replaced together:\n%s", got)
	}
	if gotWords != wantWords || !strings.Contains(got, "Show empty and error states; keep retry manual.") {
		t.Fatalf("wrong revised contract:\n%s", got)
	}
}

func TestReviseTaskRejectsInvalidOrStaleRevision(t *testing.T) {
	task := brieftest.Ship("Show screen states.")
	current := brief.RoleHeading + "\n# Task\n" + task + "\n# Setup\nWorker rules.\n"
	for name, revision := range map[string]string{
		"rewrite captain":     strings.Replace(task, "Show screen states.", "Invent a new flow.", 1),
		"missing acceptance":  strings.Replace(task, "## Acceptance", "## Something else", 1),
		"duplicate build":     task + "\n## Build\nMore scope.\n",
		"worker instructions": current,
		"change kind":         brieftest.Scout("Show screen states.", "Which states?"),
	} {
		t.Run(name, func(t *testing.T) {
			if got, err := brief.ReviseTask(current, revision, brief.Ship); err == nil || got != "" {
				t.Fatalf("accepted invalid revision: %q, %v", got, err)
			}
		})
	}
	later, _ := brief.AppendCaptainsWords(current, "Keep retry manual.", time.Now())
	if _, err := brief.ReviseTask(later, task, brief.Ship); err == nil {
		t.Fatal("a stale revision deleted newer captain words")
	}
}

func TestReviseScoutAndFencedHeadings(t *testing.T) {
	task := brieftest.Scout("Find the state path.\n```markdown\n# Example\n## Build\n```", "Which states are reachable?")
	current := brief.RoleHeading + "\n# Task\n" + task + "\n# Setup\nRules.\n"
	revision := strings.Replace(task, "Which states are reachable?", "Which states are reachable, and from which user path?", 1)
	got, err := brief.ReviseTask(current, revision, brief.Scout)
	if err != nil || !strings.HasSuffix(got, "# Setup\nRules.\n") {
		t.Fatalf("scout revision: %q, %v", got, err)
	}
	if problems := brief.CheckFile(got, brief.Scout); len(problems) != 0 {
		t.Fatal(problems)
	}
}
