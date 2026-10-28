package learnings

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/store"
)

func workspace(t *testing.T) *store.Workspace {
	t.Helper()
	w, err := store.Init(t.TempDir(), store.Defaults{})
	if err != nil {
		t.Fatalf("store.Init: %v", err)
	}
	return w
}

func write(t *testing.T, w *store.Workspace, crew, body string) {
	t.Helper()
	path := w.CrewLearning("shop", crew)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

const record = `# hu-110 learning
Source: HU-110 PR #90 @ e40b2c5; reviews r1 changes, r2 changes, r3 pass. 2026-10-27.
Learning: captured
- real dismissal cannot be proved by invoking the router → brief rule; next check: HU-114

## Usage (app)
tokens: 1.2M total

## Predictions
1. router-level acceptance avoids the UI-test dispute; next check the next integration brief — outcome: open
- a demo launch argument saves the screenshot round | outcome: Open
3. reviewer targeted verification halves its gate time; baseline r1 1256s — outcome: unchanged (HU-114 r1 1256s, 2026-10-28)
not a bullet, outcome: ignored
`

// A record's predictions are the items of its `## Predictions` section
// only, named by their number or their position; everything after
// `outcome:` is the outcome and, in parentheses, the evidence.
func TestReadParsesTheSection(t *testing.T) {
	w := workspace(t)
	write(t, w, "hu-110", record)
	ps, err := Read(w, "shop", "hu-110")
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 3 {
		t.Fatalf("got %d predictions, want 3: %+v", len(ps), ps)
	}
	want := []Prediction{
		{ID: "1", Text: "router-level acceptance avoids the UI-test dispute; next check the next integration brief", Outcome: "open"},
		{ID: "2", Text: "a demo launch argument saves the screenshot round", Outcome: "open"},
		{ID: "3", Text: "reviewer targeted verification halves its gate time; baseline r1 1256s", Outcome: "unchanged", Evidence: "HU-114 r1 1256s, 2026-10-28"},
	}
	for i, p := range ps {
		if p.Crew != "hu-110" || p.ID != want[i].ID || p.Text != want[i].Text || p.Outcome != want[i].Outcome || p.Evidence != want[i].Evidence {
			t.Errorf("prediction %d = %+v, want %+v", i, p, want[i])
		}
	}
}

// Open walks every record in name order and keeps the open ones; a record
// without the section, or no records at all, is not an error.
func TestOpenListsAcrossRecords(t *testing.T) {
	w := workspace(t)
	if ps, err := Open(w, "shop"); err != nil || len(ps) != 0 {
		t.Fatalf("no records: %v %v", ps, err)
	}
	write(t, w, "hu-110", record)
	write(t, w, "hu-104", "# hu-104\nLearning: none\n")
	write(t, w, "hu-109", "# hu-109\n## Predictions\n- fewer reads with a settled brief; next check HU-112 — outcome: open\n")
	ps, err := Open(w, "shop")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, p := range ps {
		got = append(got, p.String())
	}
	want := []string{
		"hu-109 #1: fewer reads with a settled brief; next check HU-112",
		"hu-110 #1: router-level acceptance avoids the UI-test dispute; next check the next integration brief",
		"hu-110 #2: a demo launch argument saves the screenshot round",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("open = \n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// Close rewrites only the one bullet's outcome, keeps the text before it
// verbatim, dates the evidence, and refuses an unknown id or outcome word.
func TestCloseRewritesOneOutcome(t *testing.T) {
	w := workspace(t)
	write(t, w, "hu-110", record)
	today := time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)
	p, err := Close(w, "shop", "hu-110", "#2", "Improved", "HU-114 PR #93: screenshots from -pickDemo, no tap tool", today)
	if err != nil {
		t.Fatal(err)
	}
	if p.Outcome != "improved" || p.Evidence != "HU-114 PR #93: screenshots from -pickDemo, no tap tool, 2026-10-28" {
		t.Errorf("closed = %+v", p)
	}
	data, _ := os.ReadFile(w.CrewLearning("shop", "hu-110"))
	text := string(data)
	if !strings.Contains(text, "- a demo launch argument saves the screenshot round | outcome: improved (HU-114 PR #93: screenshots from -pickDemo, no tap tool, 2026-10-28)\n") {
		t.Errorf("bullet not rewritten in place:\n%s", text)
	}
	if !strings.Contains(text, "next check the next integration brief — outcome: open\n") || !strings.Contains(text, "## Usage (app)\ntokens: 1.2M total\n") {
		t.Errorf("other lines changed:\n%s", text)
	}
	open, _ := Open(w, "shop")
	if len(open) != 1 || open[0].ID != "1" {
		t.Errorf("open after close = %+v", open)
	}
	if _, err := Close(w, "shop", "hu-110", "9", "improved", "x", today); !errors.Is(err, ErrUnknownPrediction) {
		t.Errorf("unknown id: %v", err)
	}
	if _, err := Close(w, "shop", "hu-110", "1", "better", "x", today); !errors.Is(err, ErrBadOutcome) {
		t.Errorf("bad outcome: %v", err)
	}
	if _, err := Close(w, "shop", "hu-110", "1", "not-comparable", "", today); err == nil {
		t.Error("empty evidence accepted")
	}
}
