package spawn_test

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/brief/brieftest"
	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/claude"
	"github.com/nguyenngocanh94/mate/internal/harness/codex"
	"github.com/nguyenngocanh94/mate/internal/harness/pi"
	"github.com/nguyenngocanh94/mate/internal/mateassets"
	"github.com/nguyenngocanh94/mate/internal/observability"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// pluginSecret sits in a binding's config. Config values are opaque to
// Mate and are never written into a launch.
const pluginSecret = "SECRET-7f3a-do-not-write"

// writePluginPackages lays out, under <root>/plugins, a package with an
// http server and a skill (tracker), and one with a stdio server and a
// skill (local). The tracker skill opens with YAML frontmatter.
func writePluginPackages(t *testing.T, root string) {
	t.Helper()
	files := map[string]string{
		"tracker/mate-plugin.yaml":         "apiVersion: mate.plugin/v1\nname: Tracker\nmcp:\n  type: http\n  url: https://mcp.tracker.example/mcp\nskills:\n  - skills/workflow/SKILL.md\n",
		"tracker/skills/workflow/SKILL.md": "---\nname: workflow\ndescription: Work the {instance_id} tracker.\n---\n# Workflow\nCall mcp__{instance_id}__list_issues.\n",
		"local/mate-plugin.yaml":           "apiVersion: mate.plugin/v1\nname: Local\nmcp:\n  type: stdio\n  command: [node, server.js]\n  cwd: bin\nskills:\n  - skills/notes/SKILL.md\n",
		"local/skills/notes/SKILL.md":      "# Notes\n",
		"local/bin/server.js":              "",
	}
	for rel, text := range files {
		path := filepath.Join(root, "plugins", filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// pluginsYAML binds tracker-work (http) and local-tools (stdio) to shop,
// and old-tracker, disabled, too.
func pluginsYAML(trackerEnabled bool) string {
	enabled := "true"
	if !trackerEnabled {
		enabled = "false"
	}
	return `apiVersion: mate.plugins/v1
instances:
  tracker-work:
    path: plugins/tracker
    enabled: ` + enabled + `
    projects:
      shop:
        config:
          api_token: ` + pluginSecret + `
  local-tools:
    path: plugins/local
    projects:
      shop: {}
  old-tracker:
    path: plugins/tracker
    enabled: false
    projects:
      shop: {}
`
}

func writePluginsYAML(t *testing.T, w *store.Workspace, text string) {
	t.Helper()
	if err := os.MkdirAll(w.StateDir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.StateDir(), "plugins.yaml"), []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
}

// launchLog records what every harness was asked to prepare and build.
type launchLog struct {
	mu       sync.Mutex
	prepares []harness.PrepareRequest
	builds   []harness.AgentSpec
}

func (l *launchLog) lastPrepare(t *testing.T) harness.PrepareRequest {
	t.Helper()
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.prepares) == 0 {
		t.Fatal("no launch was prepared")
	}
	return l.prepares[len(l.prepares)-1]
}

func (l *launchLog) lastBuild(t *testing.T) harness.AgentSpec {
	t.Helper()
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.builds) == 0 {
		t.Fatal("no launch was built")
	}
	return l.builds[len(l.builds)-1]
}

type recordingProfile struct {
	harness.Profile
	log *launchLog
}

func (p recordingProfile) Launcher() harness.Launcher {
	return recordingLauncher{p.Profile.Launcher(), p.log}
}

type recordingLauncher struct {
	harness.Launcher
	log *launchLog
}

func (l recordingLauncher) Prepare(ctx context.Context, req harness.PrepareRequest) (harness.Prepared, error) {
	l.log.mu.Lock()
	l.log.prepares = append(l.log.prepares, req)
	l.log.mu.Unlock()
	return l.Launcher.Prepare(ctx, req)
}

func (l recordingLauncher) Build(ctx context.Context, spec harness.AgentSpec) (harness.LaunchSpec, error) {
	l.log.mu.Lock()
	l.log.builds = append(l.log.builds, spec)
	l.log.mu.Unlock()
	return l.Launcher.Build(ctx, spec)
}

// recording is r with every launch recorded in the returned log.
func recording(t *testing.T, r harness.Registry) (harness.Registry, *launchLog) {
	t.Helper()
	log := &launchLog{}
	var profiles []harness.Profile
	for _, k := range r.Kinds() {
		p, err := r.Lookup(k)
		if err != nil {
			t.Fatal(err)
		}
		profiles = append(profiles, recordingProfile{p, log})
	}
	defaults := map[harness.AgentRole]harness.Kind{}
	for _, role := range []harness.AgentRole{harness.RoleMate, harness.RoleCrew} {
		if k, err := r.Default(role); err == nil {
			defaults[role] = k
		}
	}
	out, err := harness.NewRegistry(defaults, profiles...)
	if err != nil {
		t.Fatal(err)
	}
	return out, log
}

func pluginSkillPath(w *store.Workspace, project, dir string) string {
	return filepath.Join(w.MateDir(project), ".claude", "skills", dir, "SKILL.md")
}

func wantServers(root string) []harness.MCPServer {
	return []harness.MCPServer{
		{Name: "tracker-work", Type: "http", URL: "https://mcp.tracker.example/mcp"},
		{Name: "local-tools", Type: "stdio", Command: []string{"node", "server.js"}, Cwd: filepath.Join(root, "plugins", "local", "bin")},
	}
}

// A Mate start resolves the project's plugins: the enabled instances'
// servers reach the harness's Prepare and Build, their skills are written
// with the marker beside the Mate's own, the disabled instance gives
// nothing, and no binding config lands in any file. Claude gets the servers
// as its --mcp-config file; Codex, which cannot take them, refuses the
// start before anything is launched.
func TestStartMateMaterializesPlugins(t *testing.T) {
	for _, kind := range []harness.Kind{claude.KindClaude, codex.KindCodex, pi.KindPi} {
		t.Run(string(kind), func(t *testing.T) {
			w := newWorkspace(t, "shop")
			writePluginPackages(t, w.Root())
			writePluginsYAML(t, w, pluginsYAML(true))
			rt := runtime.NewFake()
			deps := fakeDeps(t, rt)
			deps.Harnesses = codexSessionsIn(t, t.TempDir())
			reg, log := recording(t, deps.Harnesses)
			deps.Harnesses = reg

			_, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "shop", Harness: kind})
			want := wantServers(w.Root())
			if got := log.lastPrepare(t).MCPServers; !reflect.DeepEqual(got, want) {
				t.Errorf("PrepareRequest.MCPServers = %+v\nwant %+v", got, want)
			}
			if err != nil {
				t.Fatalf("StartMate: %v", err)
			}
			if got := log.lastBuild(t).MCPServers; !reflect.DeepEqual(got, want) {
				t.Errorf("servers = %+v, want %+v", got, want)
			}
			if kind == claude.KindClaude {
				cfg := claude.MCPConfigPath(w.MateDir("shop"))
				data, err := os.ReadFile(cfg)
				if err != nil {
					t.Fatal(err)
				}
				if wantCfg, _ := claude.MCPConfig(want); string(data) != string(wantCfg) {
					t.Errorf("%s = %s, want %s", cfg, data, wantCfg)
				}
				argv := rt.StartArgv[len(rt.StartArgv)-1]
				if i := slices.Index(argv, "--mcp-config"); i < 0 || i+1 >= len(argv) || argv[i+1] != cfg {
					t.Errorf("Herdr argv %q does not carry --mcp-config %s", argv, cfg)
				}
			}

			for dir, text := range map[string]string{
				// The tracker skill opens with frontmatter: it stays first, the
				// marker follows its closing line.
				"plugin-tracker-work-workflow": "---\nname: workflow\ndescription: Work the tracker-work tracker.\n---\n<!-- generated by mate plugin tracker-work; do not edit -->\n# Workflow\nCall mcp__tracker-work__list_issues.\n",
				"plugin-local-tools-notes":     "<!-- generated by mate plugin local-tools; do not edit -->\n# Notes\n",
			} {
				got, err := os.ReadFile(filepath.Join(w.MateDir("shop"), regInfoSkillsDir(t, deps.Harnesses, kind), dir, "SKILL.md"))
				if err != nil {
					t.Fatal(err)
				}
				if string(got) != text {
					t.Errorf("%s = %q, want %q", dir, got, text)
				}
			}
			skills, err := os.ReadDir(filepath.Join(w.MateDir("shop"), regInfoSkillsDir(t, deps.Harnesses, kind)))
			if err != nil {
				t.Fatal(err)
			}
			for _, e := range skills {
				if strings.Contains(e.Name(), "old-tracker") {
					t.Errorf("the disabled instance got skill %s", e.Name())
				}
			}
			assertNoSecret(t, append([]string{w.Root(), deps.ConfigHome}, filepath.Dir(deps.Binary))...)
		})
	}
}

// The only MCP data a launch carries is a server's name, type, url,
// command and cwd: no env, no header, no credential.
func TestMCPServerCarriesNoCredential(t *testing.T) {
	typ := reflect.TypeOf(harness.MCPServer{})
	var fields []string
	for i := 0; i < typ.NumField(); i++ {
		fields = append(fields, typ.Field(i).Name)
	}
	if want := []string{"Name", "Type", "URL", "Command", "Cwd"}; !reflect.DeepEqual(fields, want) {
		t.Fatalf("harness.MCPServer fields = %v, want %v", fields, want)
	}
}

// assertNoSecret fails if any file under roots other than plugins.yaml
// carries pluginSecret.
func assertNoSecret(t *testing.T, roots ...string) {
	t.Helper()
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && d.Name() == ".git" {
				return filepath.SkipDir
			}
			if !d.Type().IsRegular() || d.Name() == "plugins.yaml" {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.Contains(string(data), pluginSecret) {
				t.Errorf("%s carries a binding's config value", path)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

// A plugin skill destination the user wrote refuses the start, by path,
// before anything is launched, and stays as it was.
func TestStartMateRefusesAPluginSkillConflict(t *testing.T) {
	w := newWorkspace(t, "shop")
	writePluginPackages(t, w.Root())
	writePluginsYAML(t, w, pluginsYAML(true))
	path := pluginSkillPath(w, "shop", "plugin-tracker-work-workflow")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	const mine = "# my own skill\n"
	if err := os.WriteFile(path, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)

	_, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "shop"})
	if !errors.Is(err, mateassets.ErrPluginSkillConflict) {
		t.Fatalf("StartMate err = %v, want a plugin skill conflict", err)
	}
	if !strings.Contains(err.Error(), filepath.Dir(path)) {
		t.Errorf("error does not name %s: %v", filepath.Dir(path), err)
	}
	if observability.ExitCode(err) != observability.ExitStateConflict {
		t.Errorf("exit code = %d, want %d", observability.ExitCode(err), observability.ExitStateConflict)
	}
	if got, _ := os.ReadFile(path); string(got) != mine {
		t.Errorf("the user's skill was overwritten: %q", got)
	}
	if _, err := os.Stat(pluginSkillPath(w, "shop", "plugin-local-tools-notes")); !os.IsNotExist(err) {
		t.Errorf("a refused start still wrote another plugin skill: %v", err)
	}
	if len(rt.StartArgv) != 0 {
		t.Fatalf("a refused start launched %d agents", len(rt.StartArgv))
	}
	if meta := readMeta(t, w, "shop"); len(meta) != 0 {
		t.Fatalf("a refused start left a meta: %v", meta)
	}
}

// Disabling an instance removes its generated skill and its server from
// the next start; the user's skill and the built-in ones are untouched.
func TestStartMateRemovesStalePluginSkills(t *testing.T) {
	w := newWorkspace(t, "shop")
	writePluginPackages(t, w.Root())
	writePluginsYAML(t, w, pluginsYAML(true))
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	reg, log := recording(t, deps.Harnesses)
	deps.Harnesses = reg
	ctx := context.Background()

	if _, err := spawn.StartMate(ctx, w, deps, spawn.StartRequest{Project: "shop"}); err != nil {
		t.Fatalf("StartMate: %v", err)
	}
	stale := pluginSkillPath(w, "shop", "plugin-tracker-work-workflow")
	if _, err := os.Stat(stale); err != nil {
		t.Fatalf("first start did not write %s: %v", stale, err)
	}
	user := pluginSkillPath(w, "shop", "foo")
	if err := os.MkdirAll(filepath.Dir(user), 0o755); err != nil {
		t.Fatal(err)
	}
	const mine = "# foo\n"
	if err := os.WriteFile(user, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}
	builtin := pluginSkillPath(w, "shop", mateassets.SkillNames[0])
	builtinBefore, err := os.ReadFile(builtin)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := spawn.StopMate(ctx, w, deps, "shop"); err != nil {
		t.Fatalf("StopMate: %v", err)
	}

	writePluginsYAML(t, w, pluginsYAML(false))
	if _, err := spawn.StartMate(ctx, w, deps, spawn.StartRequest{Project: "shop", Fresh: true}); err != nil {
		t.Fatalf("second StartMate: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(stale)); !os.IsNotExist(err) {
		t.Errorf("the disabled instance's skill directory is still there: %v", err)
	}
	if _, err := os.Stat(pluginSkillPath(w, "shop", "plugin-local-tools-notes")); err != nil {
		t.Errorf("the still enabled instance's skill is gone: %v", err)
	}
	if got, _ := os.ReadFile(user); string(got) != mine {
		t.Errorf("user skill foo changed: %q", got)
	}
	if got, _ := os.ReadFile(builtin); string(got) != string(builtinBefore) {
		t.Errorf("built-in skill %s changed", mateassets.SkillNames[0])
	}
	want := wantServers(w.Root())[1:]
	if got := log.lastPrepare(t).MCPServers; !reflect.DeepEqual(got, want) {
		t.Errorf("PrepareRequest.MCPServers after disabling = %+v, want %+v", got, want)
	}
	wantCfg, _ := claude.MCPConfig(want)
	if got, _ := os.ReadFile(claude.MCPConfigPath(w.MateDir("shop"))); string(got) != string(wantCfg) {
		t.Errorf("mcp-plugins.json after disabling = %s, want %s", got, wantCfg)
	}
}

// An invalid plugins.yaml refuses the Mate start with its location, before
// anything is rendered or launched.
func TestStartMateRefusesAnInvalidPluginConfig(t *testing.T) {
	w := newWorkspace(t, "shop")
	writePluginsYAML(t, w, "apiVersion: mate.plugins/v1\ninstances:\n  broken:\n    path: plugins/missing\n    projects:\n      shop: {}\n")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)

	_, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "shop"})
	if err == nil {
		t.Fatal("StartMate accepted an invalid plugins.yaml")
	}
	for _, want := range []string{filepath.Join(w.StateDir(), "plugins.yaml"), "instance broken", "not a directory"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not carry %q: %v", want, err)
		}
	}
	if observability.ExitCode(err) != observability.ExitUsage {
		t.Errorf("exit code = %d, want %d", observability.ExitCode(err), observability.ExitUsage)
	}
	if _, err := os.Stat(filepath.Join(w.MateDir("shop"), "AGENTS.md")); !os.IsNotExist(err) {
		t.Errorf("a refused start rendered the manual: %v", err)
	}
	if len(rt.StartArgv) != 0 || len(rt.Tabs) != 0 {
		t.Fatalf("a refused start reached Herdr: %d starts, %d tabs", len(rt.StartArgv), len(rt.Tabs))
	}
}

// A Crew gets no plugin: its spawn never reads plugins.yaml, so even an
// invalid one does not stop it, and nothing of a plugin reaches its launch
// or its worktree.
func TestSpawnCrewIgnoresPlugins(t *testing.T) {
	w := crewWorkspace(t, "shop")
	writePluginsYAML(t, w, "not: [valid")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	reg, log := recording(t, deps.Harnesses)
	deps.Harnesses = reg

	res, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project:   "shop",
		Crew:      "k3",
		BriefFile: briefFile(t, w, brieftest.Ship("Add a healthcheck endpoint.\n")),
	})
	if err != nil {
		t.Fatalf("SpawnCrew with an invalid plugins.yaml: %v", err)
	}
	if got := log.lastPrepare(t); got.Role != harness.RoleCrew || got.MCPServers != nil {
		t.Errorf("crew PrepareRequest role %s, MCPServers %+v; want crew and none", got.Role, got.MCPServers)
	}
	if got := log.lastBuild(t).MCPServers; got != nil {
		t.Errorf("crew AgentSpec.MCPServers = %+v, want none", got)
	}
	err = filepath.WalkDir(res.Worktree, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), mateassets.PluginSkillPrefix) {
			t.Errorf("crew worktree has %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartMateReceivesWorkspacePlugin(t *testing.T) {
	w := newWorkspace(t, "future-project")
	writePluginPackages(t, w.Root())
	writePluginsYAML(t, w, "apiVersion: mate.plugins/v1\ninstances:\n  tracker-work:\n    path: plugins/tracker\n")
	deps := fakeDeps(t, runtime.NewFake())
	reg, log := recording(t, deps.Harnesses)
	deps.Harnesses = reg
	_, err := spawn.StartMate(context.Background(), w, deps, spawn.StartRequest{Project: "future-project", Harness: claude.KindClaude})
	if err != nil {
		t.Fatal(err)
	}
	servers := log.lastPrepare(t).MCPServers
	if len(servers) != 1 || servers[0].Name != "tracker-work" {
		t.Fatalf("workspace servers: %+v", servers)
	}
	if _, err := os.Stat(pluginSkillPath(w, "future-project", "plugin-tracker-work-workflow")); err != nil {
		t.Fatalf("workspace skill: %v", err)
	}
}

func regInfoSkillsDir(t *testing.T, r harness.Registry, k harness.Kind) string {
	t.Helper()
	p, err := r.Lookup(k)
	if err != nil {
		t.Fatal(err)
	}
	return p.Info().SkillsDir
}

func TestPreflightPiIgnoresGlobalMCP(t *testing.T) {
	w := newWorkspace(t, "shop")
	writePluginPackages(t, w.Root())
	writePluginsYAML(t, w, pluginsYAML(true))
	home := t.TempDir()
	t.Setenv(pi.AgentDirEnv, home)
	if err := os.WriteFile(filepath.Join(home, "mcp.json"), []byte(`{"mcpServers":{"local_tools":{"url":"https://other.example/mcp"}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	if err := spawn.PreflightMatePlugins(w, deps, "shop", pi.KindPi); err != nil {
		t.Fatalf("preflight = %v", err)
	}
	if len(rt.StartLaunches) != 0 || len(rt.ExitPrompts) != 0 {
		t.Fatal("preflight mutated running agent")
	}
}
