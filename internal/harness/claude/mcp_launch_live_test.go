package claude

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

// TestLiveClaudeMateLaunchMCP runs a Mate launch exactly as Claude{}.Prepare
// and Build lay it out, with one stdio plugin server (scripts/mcpfixture),
// and proves the model calls its tool, fresh and through --resume. Opt in
// with MATE_LIVE=1. Evidence: docs/evidence/claude-mcp-probe-2026-10-12.md,
// "PR 4 verification".
//
// The server is given as a relative command (./mcpfixture) with the
// fixture's directory as its cwd, and the Mate's cwd holds no fixture, so it
// starts only if the launch really runs it in its own directory: Claude
// 2.1.287 ignores a "cwd" key, which is why MCPConfig wraps it in /bin/sh.
// The argv is Build's, unchanged, with -p, --output-format json and
// --verbose added so the turn is headless and its init event readable.
func TestLiveClaudeMateLaunchMCP(t *testing.T) {
	harnesstest.RequireLive(t)
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		t.Fatalf("MATE_LIVE=1 but no claude CLI on PATH: %v", err)
	}
	t.Logf("claude version: %s", claudeVersion(t, claudePath))
	tmp, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	pkg := filepath.Join(tmp, "pkg")
	build := exec.Command("go", "build", "-o", filepath.Join(pkg, "mcpfixture"), "github.com/nguyenngocanh94/mate/scripts/mcpfixture")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build fixture: %v\n%s", err, out)
	}
	mateDir := filepath.Join(tmp, "mate")
	if err := os.MkdirAll(mateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := &mcpProbe{claude: claudePath, tmp: tmp, home: os.Getenv("HOME")}
	t.Cleanup(func() { p.removeProjectDir(t, mateDir) })
	req := harness.PrepareRequest{
		Role: harness.RoleMate, Cwd: mateDir, StateDir: mateDir,
		ContextPath: harnesstest.WriteAbs(t, mateDir, "AGENTS.md", "# Mate\nYou are a test agent. Do exactly what the user asks.\n"),
		// The Mate's hooks call this binary; true answers them with nothing.
		Binary: "/usr/bin/true",
		MCPServers: []harness.MCPServer{
			{Name: "probe-inst", Type: "stdio", Command: []string{"./mcpfixture"}, Cwd: pkg},
		},
	}
	logPath := filepath.Join(tmp, "fixture.log")

	session := ""
	for _, turn := range []struct {
		label, prompt, id string
	}{
		{"fresh", "Call the tool " + probeToolName + ` with id "abc123" and reply with exactly the text it returns, verbatim, nothing else.`, "abc123"},
		{"resume", "Call the tool " + probeToolName + ` again, this time with id "resume456", and reply with exactly the text it returns.`, "resume456"},
	} {
		req.ResumeSessionID = session
		prep, err := Claude{}.Prepare(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
		harnesstest.WritePrepared(t, prep)
		spec, err := Claude{}.Build(context.Background(), harness.AgentSpec{
			Role: harness.RoleMate, Kind: KindClaude, Cwd: mateDir, Launch: prep.Launch,
			ResumeSessionID: session, MCPServers: req.MCPServers,
		})
		if err != nil {
			t.Fatal(err)
		}
		session = prep.SessionID
		args := spec.Args()
		if i := slices.Index(args, "--mcp-config"); i < 0 || args[i+1] != MCPConfigPath(mateDir) {
			t.Fatalf("[%s] Build's argv %q carries no --mcp-config", turn.label, args)
		}
		cfg, _ := os.ReadFile(MCPConfigPath(mateDir))
		t.Logf("[%s] %s: %s", turn.label, MCPConfigFile, p.redact(string(cfg)))
		if err := os.Truncate(logPath, 0); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}

		argv := append([]string{"-p", turn.prompt}, args...)
		argv = append(argv, "--output-format", "json", "--verbose")
		t.Logf("[%s] argv: claude %s", turn.label, p.redact(shellJoin(argv)))
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
		cmd := exec.CommandContext(ctx, claudePath, argv...)
		cmd.Dir = mateDir
		cmd.Env = append(probeEnv(), "MCP_FIXTURE_LOG="+logPath)
		var stdout, stderr strings.Builder
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		runErr := cmd.Run()
		cancel()
		var exitErr *exec.ExitError
		if runErr != nil && !errors.As(runErr, &exitErr) {
			t.Fatalf("[%s] claude could not be run: %v; stderr: %s", turn.label, runErr, p.redact(stderr.String()))
		}
		t.Logf("[%s] exit: %v", turn.label, runErr)

		var events []map[string]any
		if err := json.Unmarshal([]byte(stdout.String()), &events); err != nil {
			t.Fatalf("[%s] stdout is not the --verbose JSON array (%v):\n%s\nstderr: %s", turn.label, err, p.redact(stdout.String()), p.redact(stderr.String()))
		}
		var mcpTools []string
		var servers []any
		for _, e := range events {
			switch e["type"] {
			case "system":
				if e["subtype"] != "init" {
					continue
				}
				for _, raw := range asSlice(e["tools"]) {
					if name, _ := raw.(string); strings.HasPrefix(name, "mcp__") {
						mcpTools = append(mcpTools, name)
					}
				}
				servers = asSlice(e["mcp_servers"])
				b, _ := json.Marshal(e["mcp_servers"])
				t.Logf("[%s] init MCP tools %q, mcp_servers %s", turn.label, mcpTools, b)
			case "assistant":
				msg, _ := e["message"].(map[string]any)
				for _, c := range asSlice(msg["content"]) {
					if block, _ := c.(map[string]any); block["type"] == "tool_use" {
						in, _ := json.Marshal(block["input"])
						t.Logf("[%s] tool_use: %v %s", turn.label, block["name"], in)
					}
				}
			case "result":
				t.Logf("[%s] result: %v (is_error=%v)", turn.label, e["result"], e["is_error"])
			}
		}
		if !slices.Equal(mcpTools, []string{probeToolName}) || len(servers) != 1 {
			t.Errorf("[%s] init MCP tools %q, %d servers; want only %s (no global server)", turn.label, mcpTools, len(servers), probeToolName)
		}
		log := readFile(logPath)
		t.Logf("[%s] fixture received: %q", turn.label, fixtureMethods(log))
		if !fixtureCalledWith(log, turn.id) {
			t.Fatalf("[%s] the fixture got no read_probe call with id %q", turn.label, turn.id)
		}
		t.Logf("[%s] verdict: tool called with id %q", turn.label, turn.id)
	}
}
