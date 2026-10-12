package codex

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

// Codex reads AGENTS.override.md at its cwd and nothing else; a Crew's copy
// of its brief is kept out of git, and no session id exists before the
// first prompt.
func TestCodexPrepareCrewCopiesTheBriefIntoTheWorktree(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	req := harnesstest.PrepareFixture(t, harness.RoleCrew)
	prep, err := Codex{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	override := CodexInstructionPath(req.Cwd)
	if len(prep.Files) != 1 || prep.Files[0].Path != override || !prep.Files[0].Exclude ||
		string(prep.Files[0].Data) != "# instructions\n" {
		t.Fatalf("files = %+v, want the brief copied to %s and excluded", prep.Files, override)
	}
	if prep.ContextPath != override || prep.SessionID != "" || prep.Launch != nil {
		t.Fatalf("context path %q, session %q, launch %v", prep.ContextPath, prep.SessionID, prep.Launch)
	}
	harnesstest.WritePrepared(t, prep)
	if _, err := (Codex{}).Build(context.Background(), harness.AgentSpec{Role: harness.RoleCrew, Cwd: req.Cwd, ContextPath: prep.ContextPath}); err != nil {
		t.Fatal(err)
	}
}

func TestCodexPrepareMateNamesItsHooks(t *testing.T) {
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	req.ResumeSessionID = "01a0d260-cd47-77d2-bee7-46d98aa0461a"
	prep, err := Codex{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if prep.SessionID != req.ResumeSessionID {
		t.Fatalf("session = %q, want the resumed one", prep.SessionID)
	}
	if got := harnesstest.FileData(t, prep, CodexHooksPath(req.Cwd)); !bytes.Equal(got, CodexHooks(req.Binary)) {
		t.Fatalf("hooks = %s, want CodexHooks", got)
	}
	if got := harnesstest.FileData(t, prep, CodexInstructionPath(req.Cwd)); string(got) != "# instructions\n" {
		t.Fatalf("override = %q, want the manual", got)
	}
}

// The Codex chain check is one of those: a spec is refused once a parent
// directory's AGENTS.md pushes the chain past the cap after Build, which no
// generic check of the override alone would see.
func TestCodexSpecRefusedOnceItsChainOutgrowsTheCap(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	root := t.TempDir()
	initFakeGit(t, root)
	cwd := filepath.Join(root, "sub")
	if err := os.Mkdir(cwd, 0o755); err != nil {
		t.Fatal(err)
	}
	harnesstest.WriteAbs(t, cwd, CodexOverrideName, "you are a crew")
	spec, err := Codex{MaxChainBytes: 1000}.Build(context.Background(), harness.AgentSpec{Cwd: cwd})
	if err != nil {
		t.Fatal(err)
	}
	harnesstest.WriteAbs(t, root, CodexBaseName, string(bytes.Repeat([]byte("r"), 1000)))
	if err := spec.ValidateRequiredContext(); !errors.Is(err, harness.ErrContextTooLarge) {
		t.Fatalf("err = %v, want ErrContextTooLarge from the chain check", err)
	}
}

// Codex has no proven way to hand an agent MCP servers, so a launch with any
// is refused before it starts, by Prepare and by Build, naming the harness
// and the evidence; no -c mcp_servers is ever built.
func TestCodexRefusesMCPServers(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	servers := []harness.MCPServer{{Name: "tracker-work", Type: "http", URL: "https://mcp.tracker.example/mcp"}}
	check := func(what string, err error) {
		t.Helper()
		if !errors.Is(err, harness.ErrMCPUnsupported) {
			t.Fatalf("%s: err = %v, want ErrMCPUnsupported", what, err)
		}
		for _, want := range []string{"codex", "unsupported", "tracker-work", "docs/evidence/codex-mcp-config-audit-2026-10-12.md"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("%s: %q does not name %q", what, err, want)
			}
		}
	}
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	req.MCPServers = servers
	_, err := Codex{}.Prepare(context.Background(), req)
	check("Prepare", err)

	req.MCPServers = nil
	prep, err := Codex{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	harnesstest.WritePrepared(t, prep)
	spec := harness.AgentSpec{Role: harness.RoleMate, Cwd: req.Cwd, ContextPath: prep.ContextPath}
	if _, err := (Codex{}).Build(context.Background(), spec); err != nil {
		t.Fatalf("without servers: %v", err)
	}
	spec.MCPServers = servers
	_, err = Codex{}.Build(context.Background(), spec)
	check("Build", err)
}
