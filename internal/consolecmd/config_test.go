package consolecmd

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestLoadAndResolveCommands(t *testing.T) {
	path := filepath.Join(t.TempDir(), "console.yaml")
	if cfg, err := Load(path); err != nil || len(cfg.Bindings) != 0 {
		t.Fatalf("absent config: %+v, %v", cfg, err)
	}
	data := `bindings:
  - key: e
    scope: crew
    label: report
    command: [code, --reuse-window, "{report_path}"]
    mode: external
  - key: v
    scope: project
    label: files
    command: [nvim, "{project_dir}"]
    cwd: notes
`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"workspace_dir": "/tmp/work space", "project_dir": "/tmp/work space/shop",
		"crew_dir": "/tmp/crew", "report_path": "/tmp/crew/a {$HOME} $(touch nope).md",
	}
	args, dir, err := cfg.Bindings[0].Resolve(vars)
	if err != nil || !slices.Equal(args, []string{"code", "--reuse-window", vars["report_path"]}) || dir != vars["crew_dir"] {
		t.Fatalf("GUI argv or path altered: %q, %s, %v", args, dir, err)
	}
	_, dir, err = cfg.Bindings[1].Resolve(vars)
	if err != nil || dir != "/tmp/work space/notes" || cfg.Bindings[1].Mode != "terminal" {
		t.Fatalf("terminal defaults: %s, %v, %+v", dir, err, cfg.Bindings[1])
	}
	vars["report_path"] = ""
	if _, _, err := cfg.Bindings[0].Resolve(vars); err == nil || !strings.Contains(err.Error(), "report") {
		t.Fatalf("missing report passed to command: %v", err)
	}
	worktree := Binding{Key: "e", Scope: "crew", Label: "code", Command: []string{"zed", "{worktree_dir}"}}
	if _, _, err := worktree.Resolve(map[string]string{"crew_dir": "/tmp/crew", "worktree_dir": ""}); err == nil || !strings.Contains(err.Error(), "no worktree") {
		t.Fatalf("missing worktree passed to command: %v", err)
	}
}

func TestLoadRefusesInvalidCommands(t *testing.T) {
	valid := "bindings:\n  - key: e\n    scope: crew\n    label: report\n    command: [editor, '{report_path}']\n"
	for name, data := range map[string]string{
		"typo":                     valid + "    commmand: [other]\n",
		"scalar argv":              strings.Replace(valid, "[editor, '{report_path}']", "editor report.md", 1),
		"empty command":            strings.Replace(valid, "[editor, '{report_path}']", "[]", 1),
		"bad mode":                 valid + "    mode: gui\n",
		"bad scope":                strings.Replace(valid, "scope: crew", "scope: mate", 1),
		"unknown variable":         strings.Replace(valid, "report_path", "repo_path", 1),
		"crew variable on project": strings.Replace(valid, "scope: crew", "scope: project", 1),
		"unclosed variable":        strings.Replace(valid, "{report_path}", "{report_path", 1),
		"multiple documents":       valid + "---\nbindings: []\n",
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "console.yaml")
			if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := Load(path); err == nil {
				t.Fatal("invalid config accepted")
			}
		})
	}
}
