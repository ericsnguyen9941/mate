package agy

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

func screenLines(t *testing.T, name string) []string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "screens", name))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(string(b), "\n")
}

func TestAgyScreensReadAsMeasured(t *testing.T) {
	s := agyScreen{}
	for _, c := range []struct {
		file    string
		content string
		busy    string
	}{
		{file: "empty.txt", content: agyPlaceholder},
		{file: "idle-after-turn.txt", content: ""},
		{file: "resume.txt", content: agyPlaceholder},
		{file: "draft-pong.txt", content: "Reply with exactly: pong"},
		{file: "draft-slash-menu.txt", content: "/quit"},
		{file: "busy-generating.txt", content: "", busy: "⣽  Generating..."},
		{file: "busy-editing.txt", content: "", busy: "⣽  Editing files..."},
		{file: "busy-running-command.txt", content: "", busy: "⣻  Running command..."},
	} {
		lines := screenLines(t, c.file)
		got, ok := s.Composer(lines)
		if !ok || got != c.content {
			t.Errorf("%s: composer %q, %v; want %q", c.file, got, ok, c.content)
		}
		evidence, busy := s.Busy(lines)
		if evidence != c.busy || busy != (c.busy != "") {
			t.Errorf("%s: busy %q, %v; want %q", c.file, evidence, busy, c.busy)
		}
	}
}

func TestAgyWrappedDraftKeepsEveryRow(t *testing.T) {
	rows, ok := agyScreen{}.ComposerRows(screenLines(t, "draft-long.txt"))
	want := []string{
		"This is a long draft that should wrap across more than one row of the composer because it keeps going and going",
		"well past the width of the pane",
	}
	if !ok || !slices.Equal(rows, want) {
		t.Fatalf("rows = %q, %v; want %q", rows, ok, want)
	}
}

// The slash menu draws "esc to cancel" in the footer, as a turn does, and
// rows opening with the glyph below the composer.
func TestAgySlashMenuIsNotBusy(t *testing.T) {
	lines := screenLines(t, "draft-slash-menu.txt")
	if evidence, busy := (agyScreen{}).Busy(lines); busy {
		t.Fatalf("an open slash menu reads as busy: %q", evidence)
	}
}

func TestAgyScrollbackPromptIsNotTheComposer(t *testing.T) {
	// The scrollback prompt sits under a shorter rule with no rule below it.
	lines := []string{"────────", "> Reply with exactly: pong", "", "  pong", ""}
	if c, ok := (agyScreen{}).Composer(lines); ok {
		t.Fatalf("scrollback read as the composer %q", c)
	}
}

func TestAgyStartupInventsNoDialog(t *testing.T) {
	s := agyScreen{}
	if got := s.ClassifyStartup(strings.Join(screenLines(t, "trust.txt"), "\n")); got != harness.StartupScreenUnrecognized {
		t.Fatalf("trust dialog classifies as %s; the launch pre-trusts, so it is unrecognised", got)
	}
	for _, d := range []harness.StartupScreen{harness.StartupScreenTrustDialog, harness.StartupScreenUpdateDialog, harness.StartupScreenHooksReview} {
		if _, err := s.StartupAnswer(d); err == nil {
			t.Errorf("StartupAnswer(%s) answered a dialog nobody measured", d)
		}
	}
	if got := s.ClassifyStartup(strings.Join(screenLines(t, "busy-generating.txt"), "\n")); got == harness.StartupScreenReady {
		t.Fatal("a turn in flight classifies as ready")
	}
}
