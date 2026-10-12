package main

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/consolecmd"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/ui/console"
)

// pluginFixture is a workspace with projects shop and blog and one
// package with an open link. bindings maps project to the key it binds.
func pluginFixture(t *testing.T, instances map[string]map[string]string) *store.Workspace {
	t.Helper()
	root := t.TempDir()
	w, err := store.Init(root, workspaceDefaults())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"shop", "blog"} {
		repo := filepath.Join(w.ProjectHome(p), p)
		if err := os.MkdirAll(repo, 0o755); err != nil {
			t.Fatal(err)
		}
		initGitRepo(t, repo)
		if err := w.AddProject(p, store.ProjectConfig{Repos: []store.RepoConfig{{Path: repo, DefaultBranch: "main"}}}); err != nil {
			t.Fatal(err)
		}
	}
	pkg := filepath.Join(root, "pkg")
	if err := os.MkdirAll(pkg, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := "apiVersion: mate.plugin/v1\nname: Board\nopen:\n  title: Open board\n  url: \"https://board.example/{config.id}\"\n"
	if err := os.WriteFile(filepath.Join(pkg, "mate-plugin.yaml"), []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	if instances != nil {
		var b strings.Builder
		b.WriteString("apiVersion: mate.plugins/v1\ninstances:\n")
		for inst, projects := range instances {
			b.WriteString("  " + inst + ":\n    path: " + pkg + "\n    projects:\n")
			for p, key := range projects {
				b.WriteString("      " + p + ":\n        config: {id: " + p + "-1}\n        key: " + key + "\n")
			}
		}
		if err := os.WriteFile(filepath.Join(w.StateDir(), "plugins.yaml"), []byte(b.String()), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return w
}

func TestPluginLinksBindAndOpen(t *testing.T) {
	w := pluginFixture(t, map[string]map[string]string{"board": {"shop": "b"}})
	links, err := loadPluginLinks(w, consolecmd.Config{})
	if err != nil {
		t.Fatal(err)
	}
	b := links.bindings()
	if len(b) != 1 || b[0].Key != "b" || b[0].Label != "Open board" || b[0].Scope != "project" || b[0].Role != "plugin-board" {
		t.Fatalf("bindings %+v", b)
	}
	var got consolecmd.Invocation
	runner := func(_ context.Context, inv consolecmd.Invocation, _ io.Reader, _, _ io.Writer) error {
		got = inv
		return nil
	}
	opener := filepath.Join(t.TempDir(), "opener")
	view := consolePluginView(links, opener, runner, nil)
	if err := view(context.Background(), "b", console.StageTarget{ProjectID: "shop"}); err != nil {
		t.Fatal(err)
	}
	if got.Name != opener || len(got.Args) != 1 || got.Args[0] != "https://board.example/shop-1" {
		t.Fatalf("invocation %+v", got)
	}

	got = consolecmd.Invocation{}
	err = view(context.Background(), "b", console.StageTarget{ProjectID: "blog"})
	if err == nil || err.Error() != "no app bound to b for project blog" || got.Name != "" {
		t.Fatalf("unbound project: %v, %+v", err, got)
	}

	failing := func(context.Context, consolecmd.Invocation, io.Reader, io.Writer, io.Writer) error {
		return errors.New("boom")
	}
	if err := consolePluginView(links, opener, failing, nil)(context.Background(), "b", console.StageTarget{ProjectID: "shop"}); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("runner error lost: %v", err)
	}
	if err := consolePluginView(links, "no-such-opener-xyz", runner, nil)(context.Background(), "b", console.StageTarget{ProjectID: "shop"}); err == nil || !strings.Contains(err.Error(), "no-such-opener-xyz") {
		t.Fatalf("missing opener not named: %v", err)
	}
}

func TestPluginLinksSharedKeyOneLinePerInstance(t *testing.T) {
	w := pluginFixture(t, map[string]map[string]string{"board": {"shop": "b", "blog": "b"}})
	links, err := loadPluginLinks(w, consolecmd.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if b := links.bindings(); len(b) != 1 {
		t.Fatalf("same key, same instance, two projects: %+v", b)
	}
	w = pluginFixture(t, map[string]map[string]string{"board": {"shop": "b"}, "wiki": {"blog": "b"}})
	links, err = loadPluginLinks(w, consolecmd.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if b := links.bindings(); len(b) != 2 {
		t.Fatalf("same key, two instances: %+v", b)
	}
}

func TestPluginLinkCollisions(t *testing.T) {
	for _, tc := range []struct {
		name, key string
		cfg       consolecmd.Config
		want      string
	}{
		{"console.yaml", "t", consolecmd.Config{Bindings: []consolecmd.Binding{{Key: "t", Scope: "project", Label: "term", Command: []string{"x"}, Mode: "external"}}}, "already bound"},
		{"reserved", "q", consolecmd.Config{}, "reserved"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := pluginFixture(t, map[string]map[string]string{"board": {"shop": tc.key}})
			_, err := loadPluginLinks(w, tc.cfg)
			if err == nil {
				t.Fatal("collision accepted")
			}
			for _, s := range []string{"plugins.yaml", `"board"`, `"` + tc.key + `"`, tc.want} {
				if !strings.Contains(err.Error(), s) {
					t.Errorf("error %q lacks %q", err, s)
				}
			}
		})
	}
	// A crew-scope console.yaml binding of the same key does not collide.
	w := pluginFixture(t, map[string]map[string]string{"board": {"shop": "t"}})
	cfg := consolecmd.Config{Bindings: []consolecmd.Binding{{Key: "t", Scope: "crew", Label: "c", Command: []string{"x"}, Mode: "external"}}}
	if _, err := loadPluginLinks(w, cfg); err != nil {
		t.Fatal(err)
	}
}

func TestNoPluginsYAMLChangesNothing(t *testing.T) {
	w := pluginFixture(t, nil)
	links, err := loadPluginLinks(w, consolecmd.Config{})
	if err != nil || len(links) != 0 || len(links.bindings()) != 0 {
		t.Fatalf("links %v, %v", links, err)
	}
	if consolePluginView(links, "open", nil, nil) != nil {
		t.Fatal("view installed without plugins")
	}
}

func TestPluginLinkConfigErrorBlocksStartup(t *testing.T) {
	w := pluginFixture(t, nil)
	if err := os.WriteFile(filepath.Join(w.StateDir(), "plugins.yaml"), []byte("bogus: [\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadPluginLinks(w, consolecmd.Config{}); err == nil {
		t.Fatal("bad plugins.yaml accepted")
	}
}

func TestUrlOpenerByOS(t *testing.T) {
	if urlOpener("darwin") != "open" || urlOpener("linux") != "xdg-open" || urlOpener("plan9") != "" {
		t.Fatal("opener table")
	}
}
