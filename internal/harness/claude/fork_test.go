package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// A fork run must not act or pass for the Mate (CON-104): the launch's
// bypass, tools, plugins and session go, the fork's settings carry no hook,
// and the first run forks the Mate's session while later ones resume the
// fork.
func TestForkTurnArgsLockTheForkDown(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "settings.json")
	if err := os.WriteFile(src, []byte(`{"apiKeyHelper":"x","env":{"A":"1"},"hooks":{"Stop":[{"hooks":[{"type":"command","command":"mate hook mate-stop"}]}]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	launch := []string{"--dangerously-skip-permissions", "--session-id", "mate-sid", "--settings", src, "--setting-sources", "project", "--strict-mcp-config", "--tools", "Bash,Read,Write,Edit,Glob,Grep,Skill", "--model", "opus", "--mcp-config", "/x/mcp.json"}
	allow := []string{"Bash(git diff*)", "Bash(mate news propose *)"}
	forkDir := filepath.Join(dir, "topics")

	first, err := ForkTurnArgs(launch, "", allow, forkDir)
	if err != nil {
		t.Fatal(err)
	}
	line := strings.Join(first, " ")
	for _, gone := range []string{"--dangerously-skip-permissions", "--session-id", "Write", "Skill", "--mcp-config", src} {
		if strings.Contains(line, gone) {
			t.Errorf("fork args keep %q: %s", gone, line)
		}
	}
	for _, want := range []string{"--strict-mcp-config", "--setting-sources project", "--model opus", "--tools " + ForkTools, "--permission-mode dontAsk", "--allowedTools Bash(git diff*) Bash(mate news propose *)", "--resume mate-sid --fork-session", "-p --output-format stream-json"} {
		if !strings.Contains(line, want) {
			t.Errorf("fork args lack %q: %s", want, line)
		}
	}

	later, err := ForkTurnArgs(launch, "fork-sid", allow, forkDir)
	if err != nil {
		t.Fatal(err)
	}
	if i := slices.Index(later, "--resume"); i < 0 || later[i+1] != "fork-sid" || slices.Contains(later, "--fork-session") {
		t.Fatalf("a later run does not resume the fork: %v", later)
	}

	data, err := os.ReadFile(filepath.Join(forkDir, ForkSettingsFile))
	if err != nil {
		t.Fatal(err)
	}
	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatal(err)
	}
	if _, ok := settings["hooks"]; ok || settings["disableAllHooks"] != true || settings["apiKeyHelper"] != "x" {
		t.Fatalf("fork settings = %v, want the Mate's without hooks and with disableAllHooks", settings)
	}
}

func TestForkTurnArgsNeedASession(t *testing.T) {
	if _, err := ForkTurnArgs([]string{"--model", "opus"}, "", nil, t.TempDir()); err == nil {
		t.Fatal("a launch with no session forked")
	}
}
