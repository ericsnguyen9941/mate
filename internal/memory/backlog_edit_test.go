package memory

import (
	"strings"
	"testing"
)

func TestBacklogMovesPreserveVerbatimObligationsAndMultilineEvidence(t *testing.T) {
	original := BacklogHeader() + "\n## Notes\nDo not rewrite this.\n"
	got, _, err := EditBacklog(original, "add", "checkout", BacklogQueued, `promised "Giữ cách cũ?" Waits on captain.`)
	if err != nil {
		t.Fatal(err)
	}
	got = strings.Replace(got, `Waits on captain.`, `Waits on captain.`+"\n  evidence: crews/checkout/report.md", 1)
	moved, _, err := EditBacklog(got, "move", "checkout", BacklogInFlight, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(moved, `promised "Giữ cách cũ?" Waits on captain.`+"\n  evidence: crews/checkout/report.md") || !strings.Contains(moved, "## Notes\nDo not rewrite this.") {
		t.Fatal(moved)
	}
	if _, _, err := EditBacklog(moved, "add", "checkout", BacklogQueued, "replacement"); err == nil {
		t.Fatal("duplicate silently replaced obligation")
	}
}

func TestBacklogDoneArchivesOldEntriesWithoutLosingThem(t *testing.T) {
	text := BacklogHeader()
	for _, id := range []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven"} {
		text += "- [" + id + "] resolved\n"
	}
	text, _, _ = EditBacklog(text, "add", "twelve", BacklogInFlight, "ship it")
	next, archive, err := EditBacklog(text, "done", "twelve", BacklogDone, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(archive, "eleven") || !strings.Contains(next, "twelve") {
		t.Fatalf("next=%s archive=%s", next, archive)
	}
}

// With rewrite, a move or done leaves the one pointer line `- <id>: <text>`
// where the entry's body was (CON-77): the tracker holds the content.
func TestBacklogRewriteLeavesAPointerLine(t *testing.T) {
	text := "# Backlog\n\n## In flight\n- k2: fix the flaky test\n  evidence: logs\n\n## Held for the captain\n\n## Queued\n- k3: pick a colour\n\n## Done\n"
	next, _, err := EditBacklogEntry(text, "move", "k3", BacklogInFlight, "HVN-12", true)
	if err != nil {
		t.Fatal(err)
	}
	want := "# Backlog\n\n## In flight\n- k3: HVN-12\n- k2: fix the flaky test\n  evidence: logs\n\n## Held for the captain\n\n## Queued\n\n## Done\n"
	if next != want {
		t.Fatalf("move with rewrite =\n%s\nwant\n%s", next, want)
	}
	next, _, err = EditBacklogEntry(next, "done", "k2", BacklogDone, "HVN-7", true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(next, "## Done\n- k2: HVN-7\n") || strings.Contains(next, "evidence: logs") {
		t.Fatalf("done with rewrite =\n%s", next)
	}
	if _, _, err := EditBacklogEntry(text, "move", "k3", BacklogInFlight, "", true); err == nil {
		t.Fatal("a rewrite with no text was accepted")
	}
	// Without rewrite the body travels whole, as before.
	next, _, err = EditBacklogEntry(text, "done", "k2", BacklogDone, "ignored", false)
	if err != nil || !strings.Contains(next, "## Done\n- k2: fix the flaky test\n  evidence: logs\n") {
		t.Fatalf("done without rewrite = %v\n%s", err, next)
	}
}

// TestBacklogMovesAnEntryOutOfTheRetiredHeldSection: a backlog written
// before CON-136 still has `## Held for the captain`; its promises move
// to Queued with backlog move, and the section is no longer required.
func TestBacklogMovesAnEntryOutOfTheRetiredHeldSection(t *testing.T) {
	text := "# Backlog\n\n## In flight\n\n## " + BacklogHeld + "\n- reinstall, 2026-10-28: promised \"I will reinstall it.\" Waits on: the captain asking.\n\n## Queued\n\n## Done\n"
	got, _, err := EditBacklog(text, "move", "reinstall", BacklogQueued, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "## Queued\n- reinstall, 2026-10-28: promised") {
		t.Fatal(got)
	}
	if strings.Contains(BacklogHeader(), BacklogHeld) {
		t.Errorf("a fresh backlog still has ## %s", BacklogHeld)
	}
}
