package fixture_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/claude"
	"github.com/nguyenngocanh94/mate/internal/harness/codex"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/screen/fixture"
)

func stateCapture(t *testing.T, kind, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "harness", kind, "testdata", "startup", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// The stub names a state only for what a harness profile measured; every
// other screen is unknown, which is where Jev is asked.
func TestFixtureStateFromTheMeasuredStartup(t *testing.T) {
	o := fixture.New().(screen.StateObserver)
	cases := []struct {
		kind, file string
		want       screen.PaneState
	}{
		{"codex", "codex-0.162.1-hooks-review.txt", screen.PaneHooksReview},
		{"codex", "codex-0.162.1-loading-composer.txt", screen.PaneReady},
		{"codex", "codex-0.156.1-trust-dialog.txt", screen.PaneTrustDir},
		{"codex", "codex_update_dialog.txt", screen.PaneUpdate},
		{"claude", "claude-2.1.270-trust-dialog.txt", screen.PaneTrustDir},
		{"claude", "claude-2.1.270-ready.txt", screen.PaneReady},
		{"claude", "claude-2.1.285-bypass-dialog.txt", screen.PaneBypass},
	}
	for _, c := range cases {
		profile := codex.Codex{}.Screen()
		if c.kind == "claude" {
			profile = claude.Claude{}.Screen()
		}
		r, err := o.ObserveState(context.Background(), profile, stateCapture(t, c.kind, c.file))
		if err != nil || r.State != c.want || r.Source != fixture.Source || r.Confidence != 1 {
			t.Errorf("%s: %+v %v, want %s from the fixture at 1", c.file, r, err, c.want)
		}
	}
	r, _ := o.ObserveState(context.Background(), codex.Codex{}.Screen(), "welcome to something nobody measured\n  a) yes\n  b) no\n")
	if r.State != screen.PaneUnknown || r.Confidence != 0 {
		t.Errorf("unmeasured screen = %+v, want unknown at 0", r)
	}
}

// StateOf is the one mapping from a measured startup classification to a
// pane state; an unrecognised screen has no state the stub can name.
func TestStateOfNamesOnlyMeasuredStartups(t *testing.T) {
	if got := fixture.StateOf(harness.StartupScreenHooksReview); got != screen.PaneHooksReview {
		t.Errorf("hooks review = %s", got)
	}
	if got := fixture.StateOf(harness.StartupScreenUnrecognized); got != screen.PaneUnknown {
		t.Errorf("unrecognised = %s, want unknown", got)
	}
}

func injectionCapture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "scripts", "jeveval", "testdata", "screens", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// The stub is the startup classifier plus the fixture's composer read
// (docs/plans/pane-state-jev-2026-10-13.md section 4): a typed line under a
// screen no startup dialog matches is draft, and a turn in flight is busy
// even over a composer the startup classifier calls ready. Ready itself
// comes from the startup classifier alone.
func TestFixtureStateReadsTheComposer(t *testing.T) {
	o := fixture.New().(screen.StateObserver)
	cases := []struct {
		profile harness.ScreenProfile
		file    string
		want    screen.PaneState
	}{
		{claude.Claude{}.Screen(), "claude_injection_draft_says_empty.txt", screen.PaneDraft},
		{codex.Codex{}.Screen(), "codex_injection_busy_output_says_empty.txt", screen.PaneBusy},
		{claude.Claude{}.Screen(), "claude_injection_empty_after_trust.txt", screen.PaneReady},
	}
	for _, c := range cases {
		r, err := o.ObserveState(context.Background(), c.profile, injectionCapture(t, c.file))
		if err != nil || r.State != c.want || r.Source != fixture.Source || r.Confidence != 1 {
			t.Errorf("%s: %+v %v, want %s from the fixture at 1", c.file, r, err, c.want)
		}
	}
}

// Codex marks the highlighted row of its hook lists and tables with its
// composer glyph. Those screens are not a composer holding a draft: the
// stub leaves them unknown, which is where Jev is asked.
func TestFixtureStateLeavesTheCodexHookScreensUnknown(t *testing.T) {
	o := fixture.New().(screen.StateObserver)
	for _, file := range []string{
		"codex-0.156.1-hooks-table-review.txt",
		"codex-0.156.1-hooks-table-review-two.txt",
		"codex-0.156.1-hooks-table-trusted.txt",
		"codex-0.156.1-hooks-sessionstart-own.txt",
		"codex-0.156.1-hooks-sessionstart-own-truncated.txt",
		"codex-0.156.1-hooks-sessionstart-own-wrapped.txt",
		"codex-0.156.1-hooks-sessionstart-two-all-trusted.txt",
		"codex-0.156.1-hooks-sessionstart-two-foreign-selected.txt",
		"codex-0.156.1-hooks-sessionstart-two-own-selected.txt",
		"codex-0.156.1-hooks-sessionstart-two-own-trusted.txt",
		"codex-0.162.1-hooks-sessionstart-own.txt",
		"codex-0.162.1-hooks-sessionstart-own-home.txt",
	} {
		r, err := o.ObserveState(context.Background(), codex.Codex{}.Screen(), stateCapture(t, "codex", file))
		if err != nil || r.State != screen.PaneUnknown || r.Confidence != 0 {
			t.Errorf("%s: %+v %v, want unknown at 0", file, r, err)
		}
	}
}

// Claude idle between turns while a shell and a monitor it started keep
// running (crew ds-components-c, hellovietnam, 2026-10-13): the turn is over
// and the composer holds only a faint suggestion, so a line can be typed
// into it, but the crew is waiting on its own background work. The stub
// names that busy_background; the composer read stays empty, so send still
// types into it.
func TestFixtureStateNamesClaudeBackgroundWork(t *testing.T) {
	pane := captureFile(t, "claude_idle_background.ansi")
	profile := claude.Claude{}.Screen()
	o := fixture.New().(screen.StateObserver)
	r, err := o.ObserveState(context.Background(), profile, pane)
	if err != nil || r.State != screen.PaneBackground || r.Confidence != 1 || r.Source != fixture.Source {
		t.Fatalf("state = %+v %v, want busy_background from the fixture at 1", r, err)
	}
	if cls := fixture.ClassifyComposer(profile, pane); cls.State != screen.ComposerEmpty {
		t.Fatalf("composer = %s (%q), want empty: the turn is over and a line can be typed", cls.State, cls.Evidence)
	}

	r, err = o.ObserveState(context.Background(), profile, capture(t, "claude_empty"))
	if err != nil || r.State == screen.PaneBackground {
		t.Fatalf("idle Claude with nothing running = %+v %v, want not busy_background", r, err)
	}
}
