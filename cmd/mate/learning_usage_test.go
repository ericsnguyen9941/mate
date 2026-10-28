package main

import (
	"context"
	"os"
	"strings"
	"testing"
)

// At a close the app writes the crew's usage into its learning record in
// `usage --why`'s words (CON-130), starting the record when the Mate has
// not, so no record is without numbers.
func TestRecordUsageAtCloseStartsTheRecordWithTheLedgersNumbers(t *testing.T) {
	w := usageWorkspace(t)
	line := recordUsageAtClose(context.Background(), w, "shop", "k3", "crew stop")
	path := w.CrewLearning("shop", "k3")
	if want := "learning: usage written in " + path; line != want {
		t.Fatalf("line = %q, want %q", line, want)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, want := range []string{
		"# k3 learning\nSource: (task-learning fills this",
		"Learning: pending\n",
		"\n## Usage (app)\nscope: crew k3 and its fresh sessions; reviewers are separate crews · measured at crew stop ",
		"tokens: 97.2k total = 90k in + 6k cache-read + 0 cache-write + 1.2k out · cost ?\n",
		"shape: 1 model call(s)",
		"never add them up",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("record lacks %q:\n%s", want, text)
		}
	}
}

// A record the Mate already wrote keeps every section of its own; the
// usage section is replaced in place, never doubled.
func TestRecordUsageAtCloseReplacesOnlyItsOwnSection(t *testing.T) {
	w := usageWorkspace(t)
	path := w.CrewLearning("shop", "k3")
	if err := upsertLearningSection(path, "k3", learningUsageHeading, "stale numbers\n"); err != nil {
		t.Fatal(err)
	}
	own := "# k3 learning\nSource: PR #9 @ abcdef0\nMetrics: crew only\nLearning: captured\n- a trap → docs/x.md; next check: y\n\n" + learningUsageHeading + "\nstale numbers\n\n## Predictions\n- p1 open\n"
	if err := os.WriteFile(path, []byte(own), 0o644); err != nil {
		t.Fatal(err)
	}
	recordUsageAtClose(context.Background(), w, "shop", "k3", "mate merge")
	recordUsageAtClose(context.Background(), w, "shop", "k3", "crew stop")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if strings.Count(text, learningUsageHeading) != 1 {
		t.Errorf("usage section doubled:\n%s", text)
	}
	if strings.Contains(text, "stale numbers") || strings.Contains(text, "measured at mate merge") {
		t.Errorf("old usage kept:\n%s", text)
	}
	for _, want := range []string{"Source: PR #9 @ abcdef0\n", "Learning: captured\n- a trap → docs/x.md; next check: y\n", "measured at crew stop ", "\n## Predictions\n- p1 open\n"} {
		if !strings.Contains(text, want) {
			t.Errorf("record lost %q:\n%s", want, text)
		}
	}
	if strings.Index(text, "## Usage (app)") > strings.Index(text, "## Predictions") {
		t.Errorf("sections reordered:\n%s", text)
	}
}

// Without a readable timeline the record says unknown and how to fill it,
// rather than nothing.
func TestRecordUsageAtCloseWritesUnknownWithoutATimeline(t *testing.T) {
	w := liveCrewWorkspace(t, "shop")
	line := recordUsageAtClose(context.Background(), w, "shop", "k3", "crew stop")
	if !strings.HasPrefix(line, "learning: usage unknown (db: no timeline at ") {
		t.Fatalf("line = %q", line)
	}
	data, err := os.ReadFile(w.CrewLearning("shop", "k3"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "## Usage (app)\nunknown: db: no timeline at ") || !strings.Contains(string(data), "`mate usage shop k3 --why` once the timeline is readable") {
		t.Errorf("record = %q", data)
	}
}

func TestUpsertLearningSectionAppendsWhenMissing(t *testing.T) {
	path := liveCrewWorkspace(t, "shop").CrewLearning("shop", "k3")
	if err := upsertLearningSection(path, "k3", learningUsageHeading, "first\n"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# k3 learning\nLearning: none\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := upsertLearningSection(path, "k3", learningUsageHeading, "second\n"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if got, want := string(data), "# k3 learning\nLearning: none\n\n## Usage (app)\nsecond\n"; got != want {
		t.Errorf("record = %q, want %q", got, want)
	}
}
