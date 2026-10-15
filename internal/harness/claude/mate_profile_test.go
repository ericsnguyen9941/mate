package claude

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

func TestClaudeMateHasAnIsolatedProfile(t *testing.T) {
	spec := manualInCwdSpec(t)
	spec.Role = harness.RoleMate
	launch, err := (Claude{}).Build(context.Background(), spec)
	if err != nil {
		t.Fatal(err)
	}
	args := launch.Args()
	for flag, value := range map[string]string{
		"--setting-sources": "project",
		"--tools":           "Bash,Read,Write,Edit,Glob,Grep,Skill",
		"--model":           "opus",
		"--effort":          "medium",
		"--autocompact":     "300000",
	} {
		i := slices.Index(args, flag)
		if i < 0 || i+1 >= len(args) || args[i+1] != value {
			t.Errorf("argv %q: want %s %s", args, flag, value)
		}
	}
	if !slices.Contains(args, "--strict-mcp-config") {
		t.Errorf("Mate inherited global MCP servers: %q", args)
	}
}

func TestClaudeMateKeepsAnExplicitModelAndEffort(t *testing.T) {
	spec := manualInCwdSpec(t)
	spec.Role, spec.Model, spec.Effort = harness.RoleMate, "sonnet", harness.EffortHigh
	launch, err := (Claude{}).Build(context.Background(), spec)
	if err != nil {
		t.Fatal(err)
	}
	args := launch.Args()
	for flag, value := range map[string]string{"--model": "sonnet", "--effort": "high"} {
		i := slices.Index(args, flag)
		if i < 0 || i+1 >= len(args) || args[i+1] != value {
			t.Errorf("explicit profile was lost: %q", args)
		}
	}
}

// Plugins add --mcp-config and change nothing else of the profile: the
// --tools allowlist stays as it is, since it filters built-in tools only and
// an MCP tool is called with it (docs/evidence/claude-mcp-probe-2026-10-12.md),
// and --strict-mcp-config still keeps global servers out.
func TestClaudeMateWithPluginsKeepsItsProfile(t *testing.T) {
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	req.MCPServers = testServers()
	got := prepareAndBuild(t, req, "")
	for flag, value := range map[string]string{
		"--setting-sources": "project",
		"--tools":           "Bash,Read,Write,Edit,Glob,Grep,Skill",
		"--model":           "opus",
		"--effort":          "medium",
		"--autocompact":     "300000",
		"--mcp-config":      MCPConfigPath(req.Cwd),
	} {
		i := slices.Index(got.args, flag)
		if i < 0 || i+1 >= len(got.args) || got.args[i+1] != value {
			t.Errorf("argv %q: want %s %s", got.args, flag, value)
		}
	}
	if !slices.Contains(got.args, "--strict-mcp-config") {
		t.Errorf("Mate with plugins inherited global MCP servers: %q", got.args)
	}
	if n := strings.Count(strings.Join(got.args, "\x00"), "--tools"); n != 1 {
		t.Errorf("argv %q carries --tools %d times", got.args, n)
	}
}

// A Crew drops the captain's user-level settings (global skills, plugins)
// and compacts at the threshold it is given, but keeps every tool and its
// MCP sources: no --tools allowlist and no --strict-mcp-config (CON-68).
func TestClaudeCrewDropsUserSettingsAndCompacts(t *testing.T) {
	req := harnesstest.PrepareFixture(t, harness.RoleCrew)
	prep, err := Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	harnesstest.WritePrepared(t, prep)
	build := func(autocompact int) []string {
		t.Helper()
		launch, err := Claude{}.Build(context.Background(), harness.AgentSpec{
			Role: harness.RoleCrew, Cwd: req.Cwd, ContextPath: prep.ContextPath, Launch: prep.Launch,
			Model: "sonnet", Effort: harness.EffortMedium, Autocompact: autocompact,
		})
		if err != nil {
			t.Fatal(err)
		}
		return launch.Args()
	}

	args := build(160000)
	for flag, value := range map[string]string{
		"--setting-sources": "project",
		"--autocompact":     "160000",
		"--model":           "sonnet",
		"--effort":          "medium",
	} {
		i := slices.Index(args, flag)
		if i < 0 || i+1 >= len(args) || args[i+1] != value {
			t.Errorf("argv %q: want %s %s", args, flag, value)
		}
	}
	for _, flag := range []string{"--tools", "--strict-mcp-config"} {
		if slices.Contains(args, flag) {
			t.Errorf("Crew argv %q carries %s, the Mate's restriction", args, flag)
		}
	}

	if args := build(0); slices.Contains(args, "--autocompact") {
		t.Errorf("Crew with no threshold carries --autocompact: %q", args)
	} else if !slices.Contains(args, "--setting-sources") {
		t.Errorf("Crew with no threshold lost --setting-sources: %q", args)
	}
}

// The Mate's argv is its own profile whatever threshold a caller passes:
// its 300000 is the prefix it has always cached.
func TestClaudeMateIgnoresACrewThreshold(t *testing.T) {
	spec := manualInCwdSpec(t)
	spec.Role, spec.Autocompact = harness.RoleMate, 160000
	launch, err := (Claude{}).Build(context.Background(), spec)
	if err != nil {
		t.Fatal(err)
	}
	args := launch.Args()
	if n := strings.Count(strings.Join(args, "\x00"), "--autocompact"); n != 1 {
		t.Fatalf("argv %q carries --autocompact %d times", args, n)
	}
	if i := slices.Index(args, "--autocompact"); args[i+1] != "300000" {
		t.Errorf("Mate threshold = %s, want 300000", args[i+1])
	}
}
