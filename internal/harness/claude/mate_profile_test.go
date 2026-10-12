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
