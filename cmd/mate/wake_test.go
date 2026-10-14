package main

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// wakeFixture is a project with a ship that handed back on a branch one
// commit ahead (k3), its reviewer, a Crew asking a question, and one still
// working.
func wakeFixture(t *testing.T) *store.Workspace {
	t.Helper()
	f := newDiffFixture(t)
	w := f.w
	commitInWorktree(t, f.worktree, "README.md", "shop\nmore\n", "add a line")
	if err := os.MkdirAll(w.CrewDir("shop", "k3"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(w.CrewHandback("shop", "k3"), []byte("## For the next crew\nREADME now has two lines.\nThe build is `make`.\n\n| criterion | pass/fail | command | output |\n| --- | --- | --- | --- |\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := w.AppendStatus("shop", "k3", "working: adding the line"); err != nil {
		t.Fatal(err)
	}
	if err := w.AppendStatus("shop", "k3", "wait-mate: ready in branch mate/k3"); err != nil {
		t.Fatal(err)
	}
	for crew, meta := range map[string]map[string]string{
		"rev-k3":   {spawn.MetaHarness: "claude", spawn.MetaModel: "sonnet", spawn.MetaEffort: "medium", spawn.MetaKind: "scout", spawn.MetaState: spawn.CrewStateSpawned, spawn.MetaTask: "Review k3"},
		"ask-k4":   {spawn.MetaHarness: "codex", spawn.MetaState: spawn.CrewStateSpawned, spawn.MetaTask: "ask"},
		"busy-k5":  {spawn.MetaHarness: "codex", spawn.MetaState: spawn.CrewStateSpawned, spawn.MetaTask: "busy"},
		"done-old": {spawn.MetaHarness: "codex", spawn.MetaState: "finished", spawn.MetaTask: "old"},
	} {
		if err := w.WriteCrewMeta("shop", crew, meta); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(w.CrewDir("shop", "rev-k3"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := store.WriteMeta(w.CrewReviewMeta("shop", "rev-k3"), map[string]string{"crew": "k3"}); err != nil {
		t.Fatal(err)
	}
	for crew, line := range map[string]string{
		"rev-k3":   "working: reading the diff",
		"ask-k4":   "needs-decision: keep the old endpoint, or remove it now that nothing calls it? Removing it breaks the mobile app v1.2.",
		"busy-k5":  "working: running the test suite",
		"done-old": "wait-mate: long done",
	} {
		if err := w.AppendStatus("shop", crew, line); err != nil {
			t.Fatal(err)
		}
	}
	return w
}

func runWake(t *testing.T, w *store.Workspace) string {
	t.Helper()
	text, err := wake(context.Background(), w, gitx.New(), "shop", time.Date(2026, 10, 14, 22, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("wake: %v", err)
	}
	return text
}

// One wake holds what the Mate would otherwise learn with tail, cat, peek
// and crew list: the question in full before the hand-back, the hand-back's
// branch, its section for the next crew and its review, then the new status
// lines and every open Crew. A closed Crew is in none of it.
func TestWakePrintsWhatToActOnInOrder(t *testing.T) {
	text := runWake(t, wakeFixture(t))
	for _, want := range []string{
		"== 1. To act on (2) ==\n",
		"needs-decision ask-k4: keep the old endpoint, or remove it now that nothing calls it? Removing it breaks the mobile app v1.2.\n",
		"wait-mate k3: ready in branch mate/k3\n",
		"  branch mate/k3: 1 commit(s) beyond main; 1 file changed, 1 insertion(+)\n",
		"(6 lines), ## For the next crew:\n    README now has two lines.\n    The build is `make`.\n  reviews:",
		"  reviews: rev-k3 pending\n",
		"== 2. New status since your last wake (4) ==\n",
		"busy-k5 working: running the test suite\n",
		"== 3. Open crews (4) ==\n",
		"rev-k3 working · claude sonnet/medium · reading the diff\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("wake lacks %q:\n%s", want, text)
		}
	}
	if strings.Index(text, "needs-decision ask-k4") > strings.Index(text, "wait-mate k3") {
		t.Errorf("a question must come before a hand-back:\n%s", text)
	}
	if strings.Contains(text, "done-old") {
		t.Errorf("a closed Crew was shown:\n%s", text)
	}
}

// Part 2 shows only what is new since the last wake; part 1 keeps an item
// until it is dealt with, however many wakes have shown it.
func TestWakeShowsOnlyNewStatusButKeepsOpenItems(t *testing.T) {
	w := wakeFixture(t)
	runWake(t, w)
	if err := w.AppendStatus("shop", "busy-k5", "working: fixing a failing test"); err != nil {
		t.Fatal(err)
	}
	text := runWake(t, w)
	if !strings.Contains(text, "== 2. New status since your last wake (1) ==\nbusy-k5 working: fixing a failing test\n") {
		t.Errorf("second wake did not show exactly the new line:\n%s", text)
	}
	if !strings.Contains(text, "== 1. To act on (2) ==") {
		t.Errorf("open items vanished after one wake:\n%s", text)
	}
}

// A review Crew's hand-back is its verdict, and the reviewed Crew's line
// says the same.
func TestWakeShowsAReviewsVerdict(t *testing.T) {
	w := wakeFixture(t)
	if err := os.WriteFile(w.CrewReport("shop", "rev-k3"), []byte("## Summary\nok\n## Verdict\npass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := w.AppendStatus("shop", "rev-k3", "wait-mate: review done"); err != nil {
		t.Fatal(err)
	}
	text := runWake(t, w)
	for _, want := range []string{
		"wait-mate rev-k3: review done\n  review of k3: verdict pass",
		"  reviews: rev-k3 pass\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("wake lacks %q:\n%s", want, text)
		}
	}
}

func TestMarkdownSectionStopsAtTheNextHeadingAndBounds(t *testing.T) {
	text := "# T\n## Summary\nx\n## For the next crew\na\nb\nc\n## Other\nz\n"
	if got := markdownSection("## For the next crew\na\n\n| c | d |\n", "For the next crew", 10); got != "a" {
		t.Errorf("section before a table = %q", got)
	}
	if got := markdownSection(text, "For the next crew", 10); got != "a\nb\nc" {
		t.Errorf("section = %q", got)
	}
	if got := markdownSection(text, "For the next crew", 2); got != "a\nb\n(+1 more lines in the section)" {
		t.Errorf("bounded section = %q", got)
	}
	if got := markdownSection(text, "Missing", 10); got != "" {
		t.Errorf("missing section = %q", got)
	}
}
