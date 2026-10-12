package claude

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

func testServers() []harness.MCPServer {
	return []harness.MCPServer{
		{Name: "tracker-work", Type: "http", URL: "https://mcp.tracker.example/mcp"},
		{Name: "local-tools", Type: "stdio", Command: []string{"node", "server.js"}, Cwd: "/plugins/local/bin"},
		{Name: "bare", Type: "stdio", Command: []string{"/usr/local/bin/bare-mcp"}},
	}
}

// The file is exactly Claude Code's --mcp-config shape: http gets type and
// url, stdio gets command and args (a cwd reached through /bin/sh, since
// Claude 2.1.287 has none), and nothing else.
func TestMCPConfigIsExact(t *testing.T) {
	got, err := MCPConfig(testServers())
	if err != nil {
		t.Fatal(err)
	}
	const want = `{
  "mcpServers": {
    "bare": {
      "command": "/usr/local/bin/bare-mcp",
      "args": []
    },
    "local-tools": {
      "command": "/bin/sh",
      "args": [
        "-c",
        "cd \"$1\" && shift && exec \"$@\"",
        "mate-mcp-cwd",
        "/plugins/local/bin",
        "node",
        "server.js"
      ]
    },
    "tracker-work": {
      "type": "http",
      "url": "https://mcp.tracker.example/mcp"
    }
  }
}
`
	if string(got) != want {
		t.Fatalf("MCPConfig =\n%s\nwant\n%s", got, want)
	}

	// No key beyond the documented ones, at any level.
	var doc map[string]map[string]map[string]any
	if err := json.Unmarshal(got, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc) != 1 || doc["mcpServers"] == nil {
		t.Fatalf("top-level keys = %v, want only mcpServers", keys(doc))
	}
	allowed := map[string][]string{
		"tracker-work": {"type", "url"},
		"local-tools":  {"args", "command"},
		"bare":         {"args", "command"},
	}
	for name, entry := range doc["mcpServers"] {
		if k := keys(entry); !slices.Equal(k, allowed[name]) {
			t.Errorf("server %s has keys %v, want %v", name, k, allowed[name])
		}
	}
	for _, banned := range []string{"env", "headers", "token", "cwd"} {
		if strings.Contains(string(got), `"`+banned+`"`) {
			t.Errorf("config carries %q:\n%s", banned, got)
		}
	}

	empty, err := MCPConfig(nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(empty) != "{\n  \"mcpServers\": {}\n}\n" {
		t.Fatalf("empty config = %q", empty)
	}
}

func keys[V any](m map[string]V) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// A server Claude cannot be given as described is refused, not dropped.
func TestMCPConfigRefusesWhatClaudeCannotTake(t *testing.T) {
	for name, s := range map[string][]harness.MCPServer{
		"no name":        {{Type: "http", URL: "https://x"}},
		"duplicate":      {{Name: "a", Type: "http", URL: "https://x"}, {Name: "a", Type: "http", URL: "https://y"}},
		"http no url":    {{Name: "a", Type: "http"}},
		"http command":   {{Name: "a", Type: "http", URL: "https://x", Command: []string{"x"}}},
		"stdio no cmd":   {{Name: "a", Type: "stdio"}},
		"stdio empty":    {{Name: "a", Type: "stdio", Command: []string{""}}},
		"stdio url":      {{Name: "a", Type: "stdio", Command: []string{"x"}, URL: "https://x"}},
		"relative cwd":   {{Name: "a", Type: "stdio", Command: []string{"x"}, Cwd: "bin"}},
		"unknown type":   {{Name: "a", Type: "sse", URL: "https://x"}},
		"type left open": {{Name: "a", URL: "https://x"}},
	} {
		if _, err := MCPConfig(s); err == nil {
			t.Errorf("%s: MCPConfig accepted %+v", name, s)
		}
	}
}

// The stdio cwd wrapper really starts the command in its directory, with
// its arguments intact, whatever they contain.
func TestMCPCwdShellRunsTheCommandInItsDirectory(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "out")
	harnesstest.WriteAbs(t, dir, "probe.sh", "pwd -P > \"$1\"; printf '%s|' \"$@\" >> \"$1\"\n")
	argv := append(append(append([]string(nil), mcpCwdShell...), dir), "/bin/sh", "./probe.sh", out, "a b", `$HOME;"x"`)
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = t.TempDir()
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v: %s", err, b)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if want := dir + "\n" + out + `|a b|$HOME;"x"|`; string(got) != want {
		t.Fatalf("server saw %q, want %q", got, want)
	}
}

// A Mate with servers gets the file in its cwd and --mcp-config last on its
// argv, fresh and resumed; without servers its argv and files are what they
// always were, and a file left by an earlier start is emptied.
func TestClaudeMateMCPConfigLaunch(t *testing.T) {
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	plain := prepareAndBuild(t, req, "")
	path := MCPConfigPath(req.Cwd)
	if slices.Contains(plain.args, "--mcp-config") || plain.has(path) {
		t.Fatalf("a Mate without plugins got MCP config: args %q files %v", plain.args, plain.files)
	}

	req.MCPServers = testServers()
	withMCP := prepareAndBuild(t, req, "")
	want, _ := MCPConfig(req.MCPServers)
	if got := withMCP.data(path); string(got) != string(want) {
		t.Fatalf("%s = %s, want %s", MCPConfigFile, got, want)
	}
	if !slices.Equal(withMCP.args, append(slices.Clone(plain.args), "--mcp-config", path)) {
		t.Fatalf("args = %q\nwant the plain Mate argv %q plus --mcp-config %s", withMCP.args, plain.args, path)
	}

	resumed := prepareAndBuild(t, req, "22222222-2222-4222-8222-222222222222")
	if i := slices.Index(resumed.args, "--mcp-config"); i < 0 || resumed.args[i+1] != path || !slices.Contains(resumed.args, "--resume") {
		t.Fatalf("resume args = %q, want --resume and --mcp-config %s", resumed.args, path)
	}

	// The plugins are gone: the file is still on disk from the last start,
	// so it is emptied, and the argv no longer names it.
	req.MCPServers = nil
	after := prepareAndBuild(t, req, "")
	if got := after.data(path); string(got) != "{\n  \"mcpServers\": {}\n}\n" {
		t.Fatalf("stale %s rewritten as %q, want no servers", MCPConfigFile, got)
	}
	if !slices.Equal(after.args, plain.args) {
		t.Fatalf("args = %q, want the plain Mate argv %q", after.args, plain.args)
	}
}

type builtLaunch struct {
	args  []string
	files []harness.LaunchFile
}

func (b builtLaunch) has(path string) bool {
	return slices.ContainsFunc(b.files, func(f harness.LaunchFile) bool { return f.Path == path })
}

func (b builtLaunch) data(path string) []byte {
	for _, f := range b.files {
		if f.Path == path {
			return f.Data
		}
	}
	return nil
}

func prepareAndBuild(t *testing.T, req harness.PrepareRequest, resume string) builtLaunch {
	t.Helper()
	req.ResumeSessionID = resume
	prep, err := Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	harnesstest.WritePrepared(t, prep)
	spec, err := Claude{}.Build(context.Background(), harness.AgentSpec{
		Role: harness.RoleMate, Cwd: req.Cwd, Launch: prep.Launch, ResumeSessionID: resume, MCPServers: req.MCPServers,
	})
	if err != nil {
		t.Fatal(err)
	}
	return builtLaunch{args: spec.Args(), files: prep.Files}
}

// A Crew gets no MCP config, and Build refuses servers its Prepare wrote no
// file for.
func TestClaudeMCPServersOnlyReachAPreparedMate(t *testing.T) {
	crew := harnesstest.PrepareFixture(t, harness.RoleCrew)
	crew.MCPServers = testServers()
	if _, err := (Claude{}).Prepare(context.Background(), crew); err == nil {
		t.Fatal("a Crew was prepared with MCP servers")
	}

	spec := manualInCwdSpec(t)
	spec.Role = harness.RoleMate
	spec.MCPServers = testServers()
	if _, err := (Claude{}).Build(context.Background(), spec); err == nil {
		t.Fatal("Build passed MCP servers that no --mcp-config file was written for")
	}

	spec = manualInCwdSpec(t)
	spec.Role = harness.RoleCrew
	launch := spec.Launch.(claudeLaunch)
	launch.mcpConfigPath = MCPConfigPath(spec.Cwd)
	spec.Launch = launch
	harnesstest.WriteAbs(t, spec.Cwd, MCPConfigFile, "{}")
	if _, err := (Claude{}).Build(context.Background(), spec); err == nil {
		t.Fatal("a Crew was built with --mcp-config")
	}

	spec.Role = harness.RoleMate
	if err := os.Remove(launch.mcpConfigPath); err != nil {
		t.Fatal(err)
	}
	if _, err := (Claude{}).Build(context.Background(), spec); err == nil {
		t.Fatal("Build named a --mcp-config file that is not there")
	}
}
