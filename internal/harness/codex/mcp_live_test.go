package codex

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

func TestLiveCodexMCPConfigParser(t *testing.T) {
	harnesstest.RequireLive(t)
	t.Setenv(CodexHomeEnv, t.TempDir())
	servers := []harness.MCPServer{{Name: "probe.dot-hyphen", Type: "stdio", Command: []string{"./server", "arg with space"}, Cwd: "/plugin dir"}, {Name: "remote", Type: "http", URL: "https://example.com/mcp"}}
	args, err := MCPArgs(servers)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("codex", append(args, "mcp", "get", "probe.dot-hyphen")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("parser: %v %s", err, out)
	}
	for _, value := range []string{"./server", "arg with space", "/plugin dir"} {
		if !strings.Contains(string(out), value) {
			t.Fatalf("lost %q: %s", value, out)
		}
	}
	t.Logf("native parser retained command, args and cwd")
}
