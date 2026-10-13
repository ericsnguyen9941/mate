package codex

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

func TestMCPConfigRoundTripAndResume(t *testing.T) {
	t.Setenv(CodexHomeEnv, t.TempDir())
	servers := []harness.MCPServer{{Name: "a.b-quote\"", Type: "stdio", Command: []string{"./server", "quote\"\n", "$(not-executed)"}, Cwd: "/plugin dir"}, {Name: "http", Type: "http", URL: "https://example.com/mcp"}}
	doc, err := MCPConfig(servers)
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Servers map[string]struct {
			URL, Command, Cwd string
			Args              []string
			Enabled           bool
		} `toml:"mcp_servers"`
	}
	if _, err := toml.Decode("mcp_servers="+string(doc), &decoded); err != nil {
		t.Fatal(err)
	}
	got := decoded.Servers[servers[0].Name]
	if got.Command != servers[0].Command[0] || !slices.Equal(got.Args, servers[0].Command[1:]) || got.Cwd != servers[0].Cwd || !got.Enabled {
		t.Fatalf("decoded: %+v", got)
	}
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	req.MCPServers = servers
	prep, err := (Codex{}).Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	harnesstest.WritePrepared(t, prep)
	var fresh []string
	for _, id := range []string{"", "01a0d260-cd47-77d2-bee7-46d98aa0461a"} {
		spec, err := (Codex{}).Build(context.Background(), harness.AgentSpec{Role: harness.RoleMate, Cwd: req.Cwd, ContextPath: prep.ContextPath, ResumeSessionID: id, MCPServers: servers})
		if err != nil {
			t.Fatal(err)
		}
		args := spec.Args()
		if id != "" {
			if args[0] != "resume" || args[len(args)-1] != id {
				t.Fatal(args)
			}
			args = args[1 : len(args)-1]
		}
		if fresh == nil {
			fresh = args
		} else if !reflect.DeepEqual(fresh, args) {
			t.Fatalf("resume changed overrides: %q -> %q", fresh, args)
		}
	}
}
func TestMCPConfigInvalidAndEmpty(t *testing.T) {
	for _, servers := range [][]harness.MCPServer{
		{{Name: "x", Type: "http", URL: "https://example.com", Command: []string{"cmd"}}},
		{{Name: "x", Type: "stdio"}},
		{{Name: "x", Type: "stdio", Command: []string{"cmd"}, Cwd: "relative"}},
		{{Name: "x", Type: "http", URL: "https://example.com"}, {Name: "x", Type: "http", URL: "https://other.com"}},
	} {
		if _, err := MCPConfig(servers); err == nil {
			t.Fatalf("accepted %+v", servers)
		}
	}
	args, err := MCPArgs(nil)
	if err != nil || !slices.Equal(args, []string{"-c", "mcp_servers={}"}) {
		t.Fatalf("empty %q %v", args, err)
	}
	// The operator's config is read by Codex, never rewritten by the adapter.
	home := t.TempDir()
	t.Setenv(CodexHomeEnv, home)
	path := filepath.Join(home, "config.toml")
	original := []byte("[mcp_servers.global]\ncommand='/bin/true'\n")
	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatal(err)
	}
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	prep, err := (Codex{}).Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	harnesstest.WritePrepared(t, prep)
	spec, err := (Codex{}).Build(context.Background(), harness.AgentSpec{Role: harness.RoleMate, Cwd: req.Cwd, ContextPath: prep.ContextPath})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(spec.Args(), "mcp_servers={}") {
		t.Fatal("empty Mate did not reset table")
	}
	got, _ := os.ReadFile(path)
	if string(got) != string(original) {
		t.Fatal("modified operator config")
	}
}
