package spawn_test

import (
	"context"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/brief/brieftest"
	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/agy"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/spawn"
)

// An agy Crew's brief is an always-on rule nested in the worktree
// (.agents/rules/mate-brief.md, CON-278): spawn writes it there and keeps it
// out of git by its path, so the Crew's branch never carries it, and a stop
// clears the composer with ctrl+u before it exits.
func TestSpawnCrewOnAgy(t *testing.T) {
	w := crewWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	res, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project: "shop", Crew: "k3", BriefText: brieftest.Ship("work"),
		Harness: agy.KindAgy, Model: "gemini-4-hydrogen-high", Effort: harness.EffortHigh,
	})
	if err != nil {
		t.Fatalf("SpawnCrew: %v", err)
	}
	if res.Harness != agy.KindAgy || !res.BriefDelivered {
		t.Fatalf("spawn = %+v", res)
	}
	brief, err := os.ReadFile(res.BriefPath)
	if err != nil {
		t.Fatal(err)
	}
	rule, err := os.ReadFile(agy.RulePath(res.Worktree))
	if err != nil {
		t.Fatalf("the brief's rule is not in the worktree: %v", err)
	}
	if !strings.HasSuffix(string(rule), string(brief)) || !strings.HasPrefix(string(rule), "---\ntrigger: always_on\n") {
		t.Errorf("rule =\n%s\nwant the brief under always_on frontmatter", rule)
	}
	if status := git(t, res.Worktree, "status", "--porcelain"); strings.TrimSpace(status) != "" {
		t.Errorf("the rule shows in the Crew's git status:\n%s", status)
	}
	if ignored := git(t, res.Worktree, "check-ignore", agy.RuleRelPath); strings.TrimSpace(ignored) != agy.RuleRelPath {
		t.Errorf("check-ignore %s = %q", agy.RuleRelPath, ignored)
	}
	args := rt.StartLaunches[0].Args()
	if !slices.Contains(args, "--dangerously-skip-permissions") || !slices.Contains(args, "gemini-4-hydrogen-high") {
		t.Errorf("args = %q", args)
	}

	if _, err := spawn.StopCrew(context.Background(), w, deps, "shop", "k3", false); err != nil {
		t.Fatalf("StopCrew: %v", err)
	}
	if len(rt.ExitClearKeys) != 1 || !slices.Equal(rt.ExitClearKeys[0], []string{"ctrl+u"}) || !slices.Equal(rt.ExitPrompts, []string{"/exit"}) {
		t.Fatalf("stop pressed %q then typed %q, want ctrl+u then /exit", rt.ExitClearKeys, rt.ExitPrompts)
	}
}
