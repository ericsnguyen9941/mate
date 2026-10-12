package claude

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

// mateClaudeTools is the --tools allowlist the Mate profile carries in
// claude.go. It is copied rather than shared on purpose: this probe measures
// the flags as they stood when it was written, and a change to the launcher
// should not silently change what an old piece of evidence measured.
const mateClaudeTools = "Bash,Read,Write,Edit,Glob,Grep,Skill"

const probeToolName = "mcp__probe-inst__read_probe"

// probePrompts are tried in order until the fixture records a tools/call;
// the brief allows at most two prompts beyond the first.
var probePrompts = []string{
	"Call the tool " + probeToolName + ` with id "abc123" and reply with exactly the text it returns, verbatim, nothing else.`,
	"Use your " + probeToolName + ` tool now. Pass {"id": "abc123"}. Then print the tool's result verbatim.`,
	`You have an MCP tool named read_probe from the server probe-inst. Invoke it with the argument id set to "abc123" and output only its result.`,
}

// TestLiveClaudeMCPProbe measures whether a Claude Code launch carrying the
// Mate profile's flags (claude.go, RoleMate) discovers and calls a tool from
// an MCP server given with --mcp-config. Opt in with MATE_LIVE=1: it builds
// scripts/mcpfixture, runs `claude -p` headless in an empty temp directory
// several times and logs, per variant, the redacted argv, the exit code, the
// tools and MCP servers the session reported at init, the result event, and
// the requests the fixture received. Evidence:
// docs/evidence/claude-mcp-probe-2026-10-12.md.
//
// It measures and does not assert the answer: a variant whose tool is never
// called passes. It fails only when claude cannot be run or the fixture never
// receives initialize, because then nothing was measured.
//
// Variants: (A) the Mate flags as they are, --tools included; (B) without
// --tools; (C) --tools with the MCP tool's name appended. A fourth subtest
// starts a session with --session-id under C's flags and calls the tool again
// through --resume with the same --mcp-config.
func TestLiveClaudeMCPProbe(t *testing.T) {
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
	fixture := filepath.Join(tmp, "mcpfixture")
	build := exec.Command("go", "build", "-o", fixture, "github.com/nguyenngocanh94/mate/scripts/mcpfixture")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build fixture: %v\n%s", err, out)
	}

	p := &mcpProbe{claude: claudePath, fixture: fixture, tmp: tmp, home: os.Getenv("HOME")}
	p.globalServers = p.listGlobalServers(t)

	variants := []struct {
		name  string
		tools []string // nil drops --tools entirely
	}{
		{"A_mate_flags_as_is", []string{"--tools", mateClaudeTools}},
		{"B_no_tools_flag", nil},
		{"C_tools_plus_mcp_name", []string{"--tools", mateClaudeTools + "," + probeToolName}},
	}
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			for i, prompt := range probePrompts {
				label := fmt.Sprintf("prompt%d", i+1)
				if p.run(t, v.name, label, "", v.tools, prompt, nil, "abc123") {
					break
				}
			}
		})
	}

	t.Run("D_resume_with_C_flags", func(t *testing.T) {
		tools := []string{"--tools", mateClaudeTools + "," + probeToolName}
		session := newClaudeSessionID(t)
		// Resume looks the session up under the cwd's project directory, so
		// both turns share one cwd.
		cwd := filepath.Join(p.tmp, "D-cwd")
		first := p.run(t, "D", "fresh", cwd, tools, probePrompts[0], []string{"--session-id", session}, "abc123")
		resumePrompt := "Call the tool " + probeToolName + ` again, this time with id "resume456", and reply with exactly the text it returns.`
		second := p.run(t, "D", "resume", cwd, tools, resumePrompt, []string{"--resume", session}, "resume456")
		t.Logf("D verdict: fresh called=%v, resumed called=%v", first, second)
	})
}

type mcpProbe struct {
	claude        string
	fixture       string
	tmp           string
	home          string
	globalServers []string
}

// listGlobalServers records the MCP servers the captain's user/plugin config
// would give a plain launch, so the probe can show whether any of them reached
// a --strict-mcp-config session. Only names are logged.
func (p *mcpProbe) listGlobalServers(t *testing.T) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, p.claude, "mcp", "list")
	cmd.Dir = p.tmp
	cmd.Env = probeEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("claude mcp list failed (%v); global servers not measured:\n%s", err, p.redact(string(out)))
		return nil
	}
	var names []string
	for _, line := range strings.Split(string(out), "\n") {
		// "<name>: <command or url> - <status>"; names may contain ':'.
		name, rest, ok := strings.Cut(line, ": ")
		if !ok || !strings.Contains(rest, " - ") || strings.HasPrefix(line, "Checking") {
			continue
		}
		names = append(names, name)
	}
	t.Logf("global MCP servers from `claude mcp list`: %q", names)
	return names
}

// run runs one headless turn in cwd (a fresh empty directory when "") and
// logs what it measured. It reports whether the fixture received a tools/call
// carrying id. The fixture log is per run, so a tools/call in it belongs to
// this turn alone.
func (p *mcpProbe) run(t *testing.T, variant, label, cwd string, tools []string, prompt string, extra []string, id string) bool {
	t.Helper()
	runDir := filepath.Join(p.tmp, variant+"-"+label)
	if cwd == "" {
		cwd = filepath.Join(runDir, "cwd")
	}
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(runDir, "fixture.log")
	cfg := map[string]any{"mcpServers": map[string]any{"probe-inst": map[string]any{
		"command": p.fixture, "args": []string{}, "env": map[string]string{"MCP_FIXTURE_LOG": logPath},
	}}}
	cfgPath := filepath.Join(runDir, "mcp.json")
	b, _ := json.Marshal(cfg)
	if err := os.WriteFile(cfgPath, b, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { p.removeProjectDir(t, cwd) })

	args := []string{"-p", prompt}
	args = append(args, extra...)
	args = append(args, "--dangerously-skip-permissions", "--setting-sources", "project", "--strict-mcp-config")
	args = append(args, tools...)
	args = append(args, "--autocompact", "300000", "--mcp-config", cfgPath, "--output-format", "json", "--verbose")
	t.Logf("[%s/%s] %s starting claude (timeout 120s); argv: claude %s", variant, label, time.Now().UTC().Format(time.RFC3339), p.redact(shellJoin(args)))

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, p.claude, args...)
	cmd.Dir = cwd
	cmd.Env = probeEnv()
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	runErr := cmd.Run()
	exit := 0
	var exitErr *exec.ExitError
	switch {
	case ctx.Err() != nil:
		// A timeout is a measured result for this variant, not a reason to
		// stop measuring the others.
		t.Logf("[%s/%s] claude did not finish in 120s; stderr: %s", variant, label, p.redact(stderr.String()))
		methods := fixtureMethods(readFile(logPath))
		t.Logf("[%s/%s] fixture received: %q", variant, label, methods)
		if !contains(methods, "initialize") {
			t.Fatalf("[%s/%s] the fixture never received initialize; the server was not started, so nothing was measured", variant, label)
		}
		called := fixtureCalledWith(readFile(logPath), id)
		t.Logf("[%s/%s] verdict: TIMED OUT after 120s (fixture tools/call with id %q received: %v)", variant, label, id, called)
		return called
	case errors.As(runErr, &exitErr):
		exit = exitErr.ExitCode()
	case runErr != nil:
		t.Fatalf("[%s/%s] claude could not be run: %v", variant, label, runErr)
	}
	t.Logf("[%s/%s] exit code: %d", variant, label, exit)
	if s := strings.TrimSpace(stderr.String()); s != "" {
		t.Logf("[%s/%s] stderr: %s", variant, label, p.redact(s))
	}

	var events []map[string]any
	if err := json.Unmarshal([]byte(stdout.String()), &events); err != nil {
		t.Fatalf("[%s/%s] stdout is not the --verbose JSON array (%v), so claude did not run a turn:\n%s", variant, label, err, p.redact(stdout.String()))
	}
	for _, e := range events {
		switch e["type"] {
		case "system":
			if e["subtype"] != "init" {
				continue
			}
			var mcpTools, builtin []string
			for _, raw := range asSlice(e["tools"]) {
				name, _ := raw.(string)
				if strings.HasPrefix(name, "mcp__") {
					mcpTools = append(mcpTools, name)
				} else {
					builtin = append(builtin, name)
				}
			}
			servers, _ := json.Marshal(e["mcp_servers"])
			t.Logf("[%s/%s] init built-in tools: %q", variant, label, builtin)
			t.Logf("[%s/%s] init MCP tools: %q", variant, label, mcpTools)
			t.Logf("[%s/%s] init mcp_servers: %s", variant, label, servers)
			p.logLeaks(t, variant, label, mcpTools, asSlice(e["mcp_servers"]))
		case "assistant":
			msg, _ := e["message"].(map[string]any)
			for _, c := range asSlice(msg["content"]) {
				block, _ := c.(map[string]any)
				if block["type"] == "tool_use" {
					in, _ := json.Marshal(block["input"])
					t.Logf("[%s/%s] assistant tool_use: %v %s", variant, label, block["name"], in)
				}
			}
		case "user":
			msg, _ := e["message"].(map[string]any)
			for _, c := range asSlice(msg["content"]) {
				block, _ := c.(map[string]any)
				if block["type"] == "tool_result" {
					res, _ := json.Marshal(block["content"])
					t.Logf("[%s/%s] tool_result (is_error=%v): %s", variant, label, block["is_error"], p.redact(string(res)))
				}
			}
		case "result":
			delete(e, "usage")
			delete(e, "modelUsage")
			res, _ := json.Marshal(e)
			t.Logf("[%s/%s] result event: %s", variant, label, p.redact(string(res)))
		}
	}

	log := readFile(logPath)
	methods := fixtureMethods(log)
	t.Logf("[%s/%s] fixture received: %q", variant, label, methods)
	if !contains(methods, "initialize") {
		t.Fatalf("[%s/%s] the fixture never received initialize; the server was not started, so nothing was measured", variant, label)
	}
	for _, line := range strings.Split(log, "\n") {
		if strings.Contains(line, `"tools/call"`) {
			t.Logf("[%s/%s] fixture tools/call: %s", variant, label, line)
		}
	}
	called := fixtureCalledWith(log, id)
	verdict := "tool NOT called"
	if called {
		verdict = "tool called"
	}
	t.Logf("[%s/%s] verdict: %s", variant, label, verdict)
	return called
}

// logLeaks reports any MCP tool or server in the session other than
// probe-inst; under --strict-mcp-config there should be none.
func (p *mcpProbe) logLeaks(t *testing.T, variant, label string, mcpTools []string, servers []any) {
	var leaked []string
	for _, name := range mcpTools {
		if !strings.HasPrefix(name, "mcp__probe-inst__") {
			leaked = append(leaked, name)
		}
	}
	for _, s := range servers {
		m, _ := s.(map[string]any)
		if name, _ := m["name"].(string); name != "probe-inst" {
			leaked = append(leaked, "server:"+name)
		}
	}
	if len(leaked) == 0 {
		t.Logf("[%s/%s] leak check: no MCP tool or server besides probe-inst (global servers on this machine: %q)", variant, label, p.globalServers)
	} else {
		t.Logf("[%s/%s] leak check: LEAKED %q", variant, label, leaked)
	}
}

// removeProjectDir deletes the transcript directory the CLI created for a
// probe cwd. The cwd is under this test's own TempDir, so the directory is
// this test's alone.
func (p *mcpProbe) removeProjectDir(t *testing.T, cwd string) {
	configDir := os.Getenv("CLAUDE_CONFIG_DIR")
	if configDir == "" {
		configDir = filepath.Join(p.home, ".claude")
	}
	dir := filepath.Join(configDir, "projects", ClaudeProjectSlug(cwd))
	if err := os.RemoveAll(dir); err != nil {
		t.Logf("cleanup: left %s in place: %v", dir, err)
	}
}

func (p *mcpProbe) redact(s string) string {
	s = strings.ReplaceAll(s, p.tmp, "<tmp>")
	if p.home != "" {
		s = strings.ReplaceAll(s, p.home, "~")
	}
	return s
}

// probeEnv is the environment of a fresh terminal launch. A test run from
// inside a Claude Code session inherits that session's variables, and with
// --setting-sources project the user settings' apiKeyHelper no longer applies:
// an inherited ANTHROPIC_BASE_URL then sends the keychain login to a proxy
// that answers 401, and the CLI retries until the timeout (measured
// 2026-10-12). The parent session's model overrides and session markers go
// too; CLAUDE_CONFIG_DIR stays.
func probeEnv() []string {
	var out []string
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		switch {
		case name == "ANTHROPIC_BASE_URL",
			strings.HasPrefix(name, "ANTHROPIC_DEFAULT_"),
			name == "CLAUDECODE",
			name == "CLAUDE_PID",
			name == "CLAUDE_EFFORT",
			name == "AI_AGENT",
			strings.HasPrefix(name, "CLAUDE_CODE_"):
			continue
		}
		out = append(out, kv)
	}
	return out
}

func fixtureMethods(log string) []string {
	var methods []string
	sc := bufio.NewScanner(strings.NewReader(log))
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		var m struct {
			Method string `json:"method"`
		}
		if json.Unmarshal(sc.Bytes(), &m) == nil && m.Method != "" {
			methods = append(methods, m.Method)
		}
	}
	return methods
}

// fixtureCalledWith reports whether the fixture log holds a tools/call for
// read_probe whose arguments.id is id. It parses each JSON-RPC line, so an id
// that appears only in _meta or another field does not count.
func fixtureCalledWith(log, id string) bool {
	sc := bufio.NewScanner(strings.NewReader(log))
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		var req struct {
			Method string `json:"method"`
			Params struct {
				Name      string `json:"name"`
				Arguments struct {
					ID string `json:"id"`
				} `json:"arguments"`
			} `json:"params"`
		}
		if json.Unmarshal(sc.Bytes(), &req) != nil {
			continue
		}
		if req.Method == "tools/call" && req.Params.Name == "read_probe" && req.Params.Arguments.ID == id {
			return true
		}
	}
	return false
}

func readFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func asSlice(v any) []any {
	s, _ := v.([]any)
	return s
}

func contains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

func shellJoin(args []string) string {
	quoted := make([]string, len(args))
	for i, a := range args {
		if a == "" || strings.ContainsAny(a, " \t\n\"'{}$*?;&|<>()") {
			quoted[i] = "'" + strings.ReplaceAll(a, "'", `'\''`) + "'"
		} else {
			quoted[i] = a
		}
	}
	return strings.Join(quoted, " ")
}
