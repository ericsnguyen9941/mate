package claude

import (
	"slices"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

var _ harness.TUIResume = Claude{}

// The TUI over a headless Mate's session (CON-227) runs the Mate's own
// launch - settings with its hooks, MCP config, model, tools, system
// prompt - on the session the host runs, resumed, and none of a headless
// turn's flags.
func TestTUIArgsAreTheMateLaunchResumedInTheTerminal(t *testing.T) {
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	req.MCPServers = testServers()
	launch := prepareAndBuild(t, req, "").args
	const sid = "33333333-3333-4333-8333-333333333333"

	got := TUIArgs(launch, sid, true)
	for _, flag := range StreamTurnFlags {
		if slices.Contains(got, flag) {
			t.Errorf("TUI args keep the headless flag %s: %q", flag, got)
		}
	}
	// Every other word of the Mate's launch is kept, in its order.
	var want []string
	for i := 0; i < len(launch); i++ {
		if launch[i] == "--resume" || launch[i] == "--session-id" {
			i++
			continue
		}
		want = append(want, launch[i])
	}
	want = append(want, "--resume", sid)
	if !slices.Equal(got, want) {
		t.Fatalf("TUI args = %q\nwant the Mate's launch resumed: %q", got, want)
	}
	for _, must := range []string{"--settings", "--mcp-config", "--model", "--tools", "--strict-mcp-config", "--dangerously-skip-permissions"} {
		if !slices.Contains(got, must) {
			t.Errorf("TUI args have no %s: %q", must, got)
		}
	}

	// A session no turn has opened yet is created under the same id, as a
	// headless turn would create it.
	fresh := TUIArgs(launch, sid, false)
	if i := slices.Index(fresh, "--session-id"); i < 0 || fresh[i+1] != sid || slices.Contains(fresh, "--resume") {
		t.Fatalf("fresh TUI args = %q, want --session-id %s", fresh, sid)
	}
}
