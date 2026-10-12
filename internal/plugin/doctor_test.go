package plugin

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// pkgDir writes a package with the given manifest and files under a fresh
// directory and returns its root.
func pkgDir(t *testing.T, manifest string, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ManifestFile), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		mode := os.FileMode(0o644)
		if strings.HasPrefix(name, "bin/") {
			mode = 0o755
		}
		if err := os.WriteFile(p, []byte(body), mode); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func lines(checks []Check) string {
	var b strings.Builder
	for _, c := range checks {
		b.WriteString(c.String() + "\n")
	}
	return b.String()
}

func failed(checks []Check) []string {
	var out []string
	for _, c := range checks {
		if !c.OK() {
			out = append(out, c.String())
		}
	}
	return out
}

func noPath(string) (string, error) { return "", errors.New("not found") }

func TestDoctorAllOK(t *testing.T) {
	pkgs := packagesDir(t)
	root := workspace(t, `apiVersion: mate.plugins/v1
instances:
  tracker-work:
    path: $PKG/full
    projects:
      shop:
        config: {project_url: "https://tracker.example/p/shop"}
        key: t
  local:
    path: $PKG/stdio
    projects:
      shop:
  notes:
    path: $PKG/no-such
    enabled: false
    projects:
      shop:
`)
	var asked []string
	checks := Doctor(root, DoctorOptions{
		Projects: []string{"shop", "blog"},
		Harness:  func(p string) (string, string) { return "alpha", "workspace default" },
		MCPConfig: func(kind string) error {
			asked = append(asked, kind)
			return nil
		},
		KeyOwner: func(string) string { return "" },
		LookPath: func(name string) (string, error) { return "/usr/bin/" + name, nil },
	})
	want := `plugins.yaml: valid: ok (3 instances; a Mate start loads it)
tracker-work: package $PKG/full: ok
tracker-work: manifest: ok
tracker-work: skill skills/workflow/SKILL.md: ok
tracker-work: mcp http https://mcp.tracker.example/mcp: ok (not contacted)
tracker-work/shop: open.url: ok (https://tracker.example/p/shop)
tracker-work/shop: key t: ok
tracker-work/shop: harness alpha: ok (workspace default)
local: package $PKG/stdio: ok
local: manifest: ok
local: mcp command node: ok (/usr/bin/node)
local/shop: harness alpha: ok (workspace default)
notes: enabled: ok (enabled: false; not read, gives nothing)
`
	if got := lines(checks); got != strings.ReplaceAll(want, "$PKG", pkgs) {
		t.Fatalf("doctor =\n%s\nwant\n%s", got, strings.ReplaceAll(want, "$PKG", pkgs))
	}
	if len(asked) != 2 {
		t.Fatalf("MCPConfig asked %v, want once per MCP binding", asked)
	}
}

func TestDoctorFailures(t *testing.T) {
	broken := pkgDir(t, "apiVersion: mate.plugin/v1\nname: Broken\nskills: [skills/a/SKILL.md, skills/b/SKILL.md]\n",
		map[string]string{"skills/a/SKILL.md": "# a\n"})
	slash := pkgDir(t, "apiVersion: mate.plugin/v1\nname: Slash\nmcp:\n  type: stdio\n  command: [./bin/serve]\n",
		map[string]string{"bin/serve": "#!/bin/sh\n"})
	slashCwd := pkgDir(t, "apiVersion: mate.plugin/v1\nname: SlashCwd\nmcp:\n  type: stdio\n  command: [./serve]\n  cwd: bin\n",
		map[string]string{"bin/x": ""})
	notExec := pkgDir(t, "apiVersion: mate.plugin/v1\nname: NotExec\nmcp:\n  type: stdio\n  command: [lib/run]\n",
		map[string]string{"lib/run": "#!/bin/sh\n"})
	badManifest := pkgDir(t, "apiVersion: mate.plugin/v1\nname: Bad\nbogus: 1\n", nil)
	pkgs := packagesDir(t)

	for _, tc := range []struct {
		name string
		yaml string
		opts DoctorOptions
		want []string
		ok   string // a line the checks must hold
	}{{
		name: "schema",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  x:\n    pth: /x\n",
		want: []string{"plugins.yaml: valid: fail: "},
	}, {
		name: "package path missing",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  gone:\n    path: /no/such/package\n    projects: {shop: {}}\n",
		want: []string{"plugins.yaml: valid: fail: ", "gone: package /no/such/package: fail: not a directory"},
	}, {
		name: "manifest unknown field",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  bad:\n    path: " + badManifest + "\n    projects: {shop: {}}\n",
		want: []string{"plugins.yaml: valid: fail: ", "bad: manifest: fail: "},
	}, {
		name: "skill file missing",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  broken:\n    path: " + broken + "\n    projects: {shop: {}}\n",
		want: []string{"plugins.yaml: valid: fail: ", "broken: skill skills/b/SKILL.md: fail: " + filepath.Join(broken, "skills/b/SKILL.md") + " does not exist"},
	}, {
		name: "open.url missing config",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  board:\n    path: $PKG/openonly\n    projects:\n      shop: {config: {team: core}}\n",
		want: []string{"plugins.yaml: valid: fail: ", "board/shop: open.url: fail: https://board.example/{config.team}/projects/{config.project_id} needs config project_id, which the binding does not set"},
	}, {
		name: "bare command not on PATH",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  local:\n    path: $PKG/stdio\n    projects: {shop: {}}\n",
		opts: DoctorOptions{LookPath: noPath},
		want: []string{"local: mcp command node: fail: node is not on PATH"},
	}, {
		name: "slash command without cwd resolves from the package root",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  slash:\n    path: " + slash + "\n    projects: {shop: {}}\n",
		opts: DoctorOptions{LookPath: noPath},
		ok:   "slash: mcp command ./bin/serve: ok (" + filepath.Join(slash, "bin", "serve") + ")\n",
	}, {
		name: "slash command resolves from cwd",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  slash:\n    path: " + slashCwd + "\n    projects: {shop: {}}\n",
		opts: DoctorOptions{LookPath: noPath},
		want: []string{"slash: mcp command ./serve: fail: " + filepath.Join(slashCwd, "bin", "serve") + " does not exist"},
	}, {
		name: "slash command not executable",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  run:\n    path: " + notExec + "\n    projects: {shop: {}}\n",
		want: []string{"run: mcp command lib/run: fail: " + filepath.Join(notExec, "lib", "run") + " is not executable"},
	}, {
		name: "harness without MCP config",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  tracker-work:\n    path: $PKG/full\n    projects:\n      shop: {config: {project_url: https://t.example/x}}\n",
		opts: DoctorOptions{
			Harness: func(string) (string, string) { return "beta", "running Mate" },
			MCPConfig: func(k string) error {
				return errors.New("the " + k + " harness cannot give a Mate MCP servers: see docs/evidence/x.md")
			},
		},
		want: []string{"tracker-work/shop: harness beta: fail: the beta harness cannot give a Mate MCP servers: see docs/evidence/x.md"},
	}, {
		name: "harness unknown is not a failure",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  tracker-work:\n    path: $PKG/full\n    projects:\n      shop: {config: {project_url: https://t.example/x}}\n",
		opts: DoctorOptions{
			Harness:   func(string) (string, string) { return "", "" },
			MCPConfig: func(string) error { t.Fatal("asked about an unknown harness"); return nil },
		},
	}, {
		name: "key owned by the console",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  board:\n    path: $PKG/openonly\n    projects:\n      shop: {config: {team: a, project_id: b}, key: q}\n",
		opts: DoctorOptions{KeyOwner: func(k string) string { return "key " + k + " is reserved by the console" }},
		want: []string{"board/shop: key q: fail: key q is reserved by the console"},
	}, {
		name: "key bound twice on a project",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  board:\n    path: $PKG/openonly\n    projects:\n      shop: {config: {team: a, project_id: b}, key: q}\n  tracker-work:\n    path: $PKG/full\n    projects:\n      shop: {config: {project_url: https://t.example/x}, key: q}\n",
		want: []string{"plugins.yaml: valid: fail: ", "board/shop: key q: fail: instance tracker-work binds the same key on this project", "tracker-work/shop: key q: fail: instance board binds the same key on this project"},
	}, {
		name: "unknown project",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  notes:\n    path: $PKG/skillsonly\n    projects: {shpo: {}}\n",
		opts: DoctorOptions{Projects: []string{"shop"}},
		want: []string{"notes/shpo: project: fail: not a project of this workspace"},
	}, {
		name: "unknown instance",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  notes:\n    path: $PKG/skillsonly\n    projects: {shop: {}}\n",
		opts: DoctorOptions{Instance: "nope"},
		want: []string{"nope: instance: fail: no instance nope in "},
	}, {
		name: "instance without that project",
		yaml: "apiVersion: mate.plugins/v1\ninstances:\n  notes:\n    path: $PKG/skillsonly\n    projects: {shop: {}}\n",
		opts: DoctorOptions{Instance: "notes", Project: "blog"},
		want: []string{"notes: project blog: fail: the instance has no binding for this project"},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			root := workspace(t, tc.yaml)
			if tc.opts.LookPath == nil {
				tc.opts.LookPath = func(n string) (string, error) { return "/bin/" + n, nil }
			}
			checks := Doctor(root, tc.opts)
			got := failed(checks)
			if len(got) != len(tc.want) {
				t.Fatalf("fails =\n%s\nwant %d:\n%s\nall:\n%s", strings.Join(got, "\n"), len(tc.want), strings.Join(tc.want, "\n"), lines(checks))
			}
			for i, w := range tc.want {
				if !strings.HasPrefix(got[i], strings.ReplaceAll(w, "$PKG", pkgs)) {
					t.Fatalf("fail %d = %q, want prefix %q\nall:\n%s", i, got[i], w, lines(checks))
				}
			}
			if tc.ok != "" && !strings.Contains(lines(checks), tc.ok) {
				t.Fatalf("checks lack %q:\n%s", tc.ok, lines(checks))
			}
		})
	}
}

func TestDoctorFilters(t *testing.T) {
	root := workspace(t, `apiVersion: mate.plugins/v1
instances:
  notes:
    path: $PKG/skillsonly
    projects: {shop: {}, blog: {}}
  board:
    path: $PKG/openonly
    projects:
      blog: {config: {team: a, project_id: b}}
`)
	got := lines(Doctor(root, DoctorOptions{Project: "shop"}))
	if strings.Contains(got, "board") || strings.Contains(got, "/blog") || !strings.Contains(got, "notes: skill skills/beta/SKILL.md: ok") {
		t.Fatalf("--project shop =\n%s", got)
	}
	got = lines(Doctor(root, DoctorOptions{Instance: "board"}))
	if strings.Contains(got, "notes") || !strings.Contains(got, "board/blog: open.url: ok (https://board.example/a/projects/b)") {
		t.Fatalf("board =\n%s", got)
	}
	got = lines(Doctor(root, DoctorOptions{Project: "docs"}))
	if !strings.HasSuffix(got, "plugins.yaml: project docs: ok (no instance binds it)\n") {
		t.Fatalf("--project docs =\n%s", got)
	}
}

func TestDoctorWithoutPluginsFile(t *testing.T) {
	got := lines(Doctor(emptyWorkspace(t), DoctorOptions{}))
	if got != "plugins.yaml: valid: ok (absent: no plugins configured)\n" {
		t.Fatalf("doctor = %q", got)
	}
}
