package claude

import (
	"slices"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// TestStartupDialogLayoutsOfOneScreenShareTheirAnswer holds the rule that
// lets one screen have several measured layouts: the settle presses one
// answer for the screen, so every layout of it must take the same keys to
// the same option position.
func TestStartupDialogLayoutsOfOneScreenShareTheirAnswer(t *testing.T) {
	t.Parallel()
	p := claudeStartup()
	first := map[harness.StartupScreen]harness.StartupDialog{}
	for _, d := range p.Dialogs {
		f, ok := first[d.Screen]
		if !ok {
			first[d.Screen] = d
			continue
		}
		if !slices.Equal(d.SelectKeys, f.SelectKeys) || d.Target != f.Target {
			t.Fatalf("%s layouts disagree: %v/%d vs %v/%d", d.Screen, d.SelectKeys, d.Target, f.SelectKeys, f.Target)
		}
	}
}

// Claude counts the shells and monitors a finished turn left running in the
// footer under its composer: `⏵⏵ bypass permissions on · 1 shell, 1 monitor ·
// ← for agents` (measured 2026-10-13). Only the footer counts: the same
// words in the transcript above the composer are something the agent
// quoted.
func TestClaudeBackgroundReadsTheFooterOnly(t *testing.T) {
	t.Parallel()
	bg, ok := harness.ScreenProfile(claudeScreen{}).(harness.BackgroundScreens)
	if !ok {
		t.Fatal("Claude's screen profile does not read background work")
	}
	rule := "──────────────────────────────"
	running := []string{
		"✻ Crunched for 11s · done 10:55 PM · 1 shell, 1 monitor still running",
		rule, "❯ ", rule,
		"  Opus 5.5 | effort: medium | in: 284.7k | out: 13 | ctx left: 72%",
		"  ⏵⏵ bypass permissions on · 1 shell, 1 monitor · ← for agents",
	}
	if evidence, ok := bg.Background(running); !ok || evidence != "1 shell, 1 monitor" {
		t.Fatalf("Background = %q, %v; want the footer's count", evidence, ok)
	}
	shells := []string{rule, "❯ ", rule, "  ⏵⏵ bypass permissions on · 2 shells · ← for agents"}
	if evidence, ok := bg.Background(shells); !ok || evidence != "2 shells" {
		t.Fatalf("Background = %q, %v; want 2 shells", evidence, ok)
	}
	quoted := []string{
		"  ⎿  ⏵⏵ bypass permissions on · 1 shell, 1 monitor · ← for agents",
		rule, "❯ ", rule,
		"  ⏵⏵ bypass permissions on · ← for agents",
	}
	if evidence, ok := bg.Background(quoted); ok {
		t.Fatalf("Background read %q from the transcript; only the footer counts", evidence)
	}
}
