package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/harnesstest"
)

// A Claude Mate loads its manual from CLAUDE.md in its cwd, so the launch
// carries no context flag, and its session and settings go on the argv.
func TestClaudePrepareMate(t *testing.T) {
	noCaptainSettings(t)
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	prep, err := Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if prep.SessionID != harnesstest.SessionID || prep.ContextPath != "" {
		t.Fatalf("session %q, context path %q; want the minted id and no context path", prep.SessionID, prep.ContextPath)
	}
	if got := harnesstest.FileData(t, prep, filepath.Join(req.Cwd, "CLAUDE.md")); string(got) != "@AGENTS.md\n" {
		t.Fatalf("CLAUDE.md = %q, want \"@AGENTS.md\\n\"", got)
	}
	settings := ClaudeSettingsPath(req.Cwd)
	want, err := ClaudeSettings(req.Binary)
	if err != nil {
		t.Fatal(err)
	}
	if got := harnesstest.FileData(t, prep, settings); !bytes.Equal(got, want) {
		t.Fatalf("settings = %s, want ClaudeSettings", got)
	}
	harnesstest.WritePrepared(t, prep)
	spec, err := Claude{}.Build(context.Background(), harness.AgentSpec{
		Role: harness.RoleMate, Cwd: req.Cwd, ContextPath: prep.ContextPath, Launch: prep.Launch,
	})
	if err != nil {
		t.Fatal(err)
	}
	if spec.Delivery() != harness.DeliveryCwdManual {
		t.Fatalf("delivery = %q, want %q", spec.Delivery(), harness.DeliveryCwdManual)
	}
	args := spec.Args()
	if i := slices.Index(args, "--session-id"); i < 0 || args[i+1] != harnesstest.SessionID {
		t.Fatalf("args %q carry no --session-id %s", args, harnesstest.SessionID)
	}
	if i := slices.Index(args, "--settings"); i < 0 || args[i+1] != settings {
		t.Fatalf("args %q carry no --settings %s", args, settings)
	}

	// A resumed start names the recorded session and mints none.
	req.ResumeSessionID = "22222222-2222-4222-8222-222222222222"
	req.NewSessionID = func() string { t.Fatal("a resumed launch minted a session id"); return "" }
	prep, err = Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	spec, err = Claude{}.Build(context.Background(), harness.AgentSpec{
		Role: harness.RoleMate, Cwd: req.Cwd, Launch: prep.Launch, ResumeSessionID: req.ResumeSessionID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if prep.SessionID != req.ResumeSessionID || slices.Contains(spec.Args(), "--session-id") || !slices.Contains(spec.Args(), "--resume") {
		t.Fatalf("resume: session %q, args %q", prep.SessionID, spec.Args())
	}
}

// An existing settings file keeps its keys; one that is not JSON refuses the
// launch rather than being overwritten.
func TestClaudePrepareMateKeepsTheSettingsFileItFinds(t *testing.T) {
	noCaptainSettings(t)
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	path := ClaudeSettingsPath(req.Cwd)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	const done = `{"autoMemoryEnabled":true,"hooks":{"SessionStart":[{"hooks":[{"type":"command","command":"x hook mate-session"}]}]}}`
	if err := os.WriteFile(path, []byte(done), 0o644); err != nil {
		t.Fatal(err)
	}
	prep, err := Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if got := harnesstest.FileData(t, prep, path); string(got) != done {
		t.Fatalf("settings = %s, want the file left byte for byte", got)
	}
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := (Claude{}).Prepare(context.Background(), req); err == nil {
		t.Fatal("a settings file that is not JSON must refuse the launch")
	}
}

// A Claude Crew takes its brief as the context flag and its settings from
// its state directory: nothing lands in the worktree.
func TestClaudePrepareCrewWritesNothingInTheWorktree(t *testing.T) {
	noCaptainSettings(t)
	req := harnesstest.PrepareFixture(t, harness.RoleCrew)
	want, err := CrewClaudeSettings(nil)
	if err != nil {
		t.Fatal(err)
	}
	prep, err := Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if len(prep.Files) != 1 || prep.Files[0].Path != filepath.Join(req.StateDir, ClaudeSettingsFile) ||
		!bytes.Equal(prep.Files[0].Data, want) || prep.Files[0].Exclude {
		t.Fatalf("files = %+v, want only the crew settings in the state directory", prep.Files)
	}
	if prep.ContextPath != req.ContextPath || prep.SessionID != harnesstest.SessionID {
		t.Fatalf("context path %q, session %q", prep.ContextPath, prep.SessionID)
	}
	harnesstest.WritePrepared(t, prep)
	spec, err := Claude{}.Build(context.Background(), harness.AgentSpec{
		Role: harness.RoleCrew, Cwd: req.Cwd, ContextPath: prep.ContextPath, Launch: prep.Launch,
	})
	if err != nil {
		t.Fatal(err)
	}
	if spec.Delivery() != harness.DeliveryAppendSystemPromptFile || spec.ContextPath() != req.ContextPath {
		t.Fatalf("delivery %q, context path %q", spec.Delivery(), spec.ContextPath())
	}
}

// --setting-sources project drops the captain's settings.json, so a Crew's
// own settings carry its apiKeyHelper and env over, from the config dir the
// launch uses, and nothing else of it. A captain file that is not JSON
// refuses the launch rather than silently changing the Crew's auth path.
func TestClaudePrepareCrewKeepsTheCaptainsAuthPath(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Setenv(ClaudeConfigDirEnv, configDir)
	captain := filepath.Join(configDir, ClaudeSettingsFile)
	if err := os.WriteFile(captain, []byte(`{"apiKeyHelper":"echo k","env":{"ANTHROPIC_BASE_URL":"http://127.0.0.1:8317"},"hooks":{"Stop":[]},"model":"opus","autoMemoryEnabled":true}`), 0o644); err != nil {
		t.Fatal(err)
	}
	req := harnesstest.PrepareFixture(t, harness.RoleCrew)
	prep, err := Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	keys := settingsKeys(t, prep.Files[0].Data)
	var env map[string]string
	if err := json.Unmarshal(keys["env"], &env); err != nil || string(keys["apiKeyHelper"]) != `"echo k"` || env["ANTHROPIC_BASE_URL"] != "http://127.0.0.1:8317" || len(env) != 1 {
		t.Fatalf("crew settings lost the captain's auth: %s", prep.Files[0].Data)
	}
	if string(keys[AutoMemoryKey]) != "false" || len(keys) != 3 {
		t.Fatalf("crew settings took more than apiKeyHelper and env: %s", prep.Files[0].Data)
	}
	if err := os.WriteFile(captain, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := (Claude{}).Prepare(context.Background(), req); err == nil {
		t.Fatal("a captain settings file that is not JSON must refuse the launch")
	}
}

// noCaptainSettings points HOME at an empty directory, so a launch finds no
// captain settings.json.
func noCaptainSettings(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv(ClaudeConfigDirEnv, "")
}

// A Mate launches with --setting-sources project as a Crew does, so each
// start sets the captain's apiKeyHelper and env in the Mate's settings,
// keeping its hooks and every other key, and a start with nothing new to
// set leaves the file byte for byte.
func TestClaudePrepareMateKeepsTheCaptainsAuthPath(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Setenv(ClaudeConfigDirEnv, configDir)
	if err := os.WriteFile(filepath.Join(configDir, ClaudeSettingsFile), []byte(`{"apiKeyHelper":"echo k","env":{"ANTHROPIC_BASE_URL":"http://127.0.0.1:8317"},"hooks":{"Stop":[]}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	req := harnesstest.PrepareFixture(t, harness.RoleMate)
	path := ClaudeSettingsPath(req.Cwd)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	const mine = `{"autoMemoryEnabled":false,"apiKeyHelper":"old","model":"opus","hooks":{"SessionStart":[{"hooks":[{"type":"command","command":"x hook mate-session"}]}]}}`
	if err := os.WriteFile(path, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}
	prep, err := Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	got := harnesstest.FileData(t, prep, path)
	keys := settingsKeys(t, got)
	var env map[string]string
	if err := json.Unmarshal(keys["env"], &env); err != nil || string(keys["apiKeyHelper"]) != `"echo k"` || env["ANTHROPIC_BASE_URL"] != "http://127.0.0.1:8317" {
		t.Fatalf("mate settings lost the captain's auth: %s", got)
	}
	if string(keys["model"]) != `"opus"` || !bytes.Contains(keys["hooks"], []byte("x hook mate-session")) || bytes.Contains(keys["hooks"], []byte("Stop")) {
		t.Fatalf("mate settings lost their own keys or took the captain's hooks: %s", got)
	}
	if err := os.WriteFile(path, got, 0o644); err != nil {
		t.Fatal(err)
	}
	again, err := Claude{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(harnesstest.FileData(t, again, path), got) {
		t.Fatal("a second start with the same captain settings rewrote the file")
	}
}
