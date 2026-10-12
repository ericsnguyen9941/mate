package dashboard

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// GET /api/plugins?project=<p> (docs/dashboard.md section 3b) is the
// project's plugins.yaml as validated metadata: names, configuration status
// and the open link. It never says an app is logged in or connected - that
// is the Mate's /mcp to answer - and a broken file is an answer, not a 500.

const linearManifest = `apiVersion: mate.plugin/v1
name: Linear
mcp:
  type: http
  url: https://mcp.linear.app/mcp
skills:
  - skills/workflow/SKILL.md
open:
  title: Open Linear
  url: "https://linear.app/{config.workspace}/team/{config.team}"
`

const linearBinding = `  linear-work:
    path: packages/linear
    projects:
      shop:
        config: {workspace: acme, team: SHOP}
`

// writePlugins puts plugins.yaml into root/.mate and the linear package
// under root/packages/linear.
func writePlugins(t *testing.T, root, pluginsYAML string) {
	t.Helper()
	pkg := filepath.Join(root, "packages", "linear")
	if err := os.MkdirAll(filepath.Join(pkg, "skills", "workflow"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		filepath.Join(pkg, "mate-plugin.yaml"):               linearManifest,
		filepath.Join(pkg, "skills", "workflow", "SKILL.md"): "# Workflow\n",
		filepath.Join(root, ".mate", "plugins.yaml"):         pluginsYAML,
	} {
		if err := os.WriteFile(name, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// getPlugins fetches the route and decodes it as plain JSON, so the test
// pins the wire shape rather than this package's own struct.
func getPlugins(t *testing.T, base, project string, status int) map[string]any {
	t.Helper()
	resp, err := http.Get(base + "/api/plugins?project=" + project)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != status {
		t.Fatalf("GET /api/plugins?project=%s = %d, want %d\n%s", project, resp.StatusCode, status, body)
	}
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("decode: %v\n%s", err, body)
	}
	return out
}

func mustJSON(t *testing.T, s string) any {
	t.Helper()
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		t.Fatalf("bad expected JSON: %v", err)
	}
	return v
}

func assertApps(t *testing.T, got map[string]any, want string) {
	t.Helper()
	if w := mustJSON(t, want); !reflect.DeepEqual(got["apps"], w) {
		g, _ := json.MarshalIndent(got["apps"], "", "  ")
		t.Errorf("apps =\n%s\nwant\n%s", g, want)
	}
}

const configuredLinear = `{"instance":"linear-work","name":"Linear","status":"configured","harness":"claude",
	"open":{"title":"Open Linear","url":"https://linear.app/acme/team/SHOP"},
	"mcp":{"name":"linear-work","type":"http"},"skills":["workflow"],"error":null}`

func TestPluginsListsAConfiguredAppWithItsLinkAndNoAuthClaim(t *testing.T) {
	f := newFixture(t)
	writePlugins(t, f.ws.Root(), "apiVersion: mate.plugins/v1\ninstances:\n"+linearBinding)

	got := getPlugins(t, f.http.URL, fixtureProject, http.StatusOK)
	if got["project"] != fixtureProject {
		t.Errorf("project = %v", got["project"])
	}
	if got["generated_at"] == nil || got["last_event_id"] == nil {
		t.Errorf("no envelope: %v", got)
	}
	if _, ok := got["error"]; ok {
		t.Errorf("a valid plugins.yaml has no top-level error: %v", got["error"])
	}
	// The running Mate is Claude, so the app is configured; the MCP URL is
	// not repeated and nothing claims a login.
	assertApps(t, got, "["+configuredLinear+"]")

	// Read on every request, never from the timeline cache: an edit shows
	// without a new event.
	writePlugins(t, f.ws.Root(), "apiVersion: mate.plugins/v1\ninstances:\n"+
		strings.Replace(linearBinding, "    path:", "    enabled: false\n    path:", 1))
	assertApps(t, getPlugins(t, f.http.URL, fixtureProject, http.StatusOK), "[]")
}

func TestPluginsShowsAConfigurationErrorOnTheApp(t *testing.T) {
	f := newFixture(t)
	root := f.ws.Root()
	broken := `  broken:
    path: packages/nothing-here
    projects:
      shop: {}
`
	writePlugins(t, root, "apiVersion: mate.plugins/v1\ninstances:\n"+linearBinding+broken)
	got := getPlugins(t, f.http.URL, fixtureProject, http.StatusOK)
	problem := filepath.Join(root, ".mate", "plugins.yaml") + ": instance broken: package path " +
		filepath.Join(root, "packages", "nothing-here") + " is not a directory"
	if got["error"] != problem {
		t.Errorf("error = %v, want %q", got["error"], problem)
	}
	blocked := "the Mate will not start with this plugins.yaml: " + problem
	want, _ := json.Marshal([]map[string]any{
		{"instance": "linear-work", "name": "linear-work", "status": "configuration error", "harness": "claude",
			"open": nil, "mcp": nil, "skills": []string{}, "error": blocked},
		{"instance": "broken", "name": "broken", "status": "configuration error", "harness": "claude",
			"open": nil, "mcp": nil, "skills": []string{}, "error": problem},
	})
	assertApps(t, got, string(want))

	// A problem in another project's binding stops this project's Mate
	// too, so its app says so rather than reading as configured.
	writePlugins(t, root, "apiVersion: mate.plugins/v1\ninstances:\n"+linearBinding+
		strings.Replace(broken, "shop:", "blog:", 1))
	got = getPlugins(t, f.http.URL, fixtureProject, http.StatusOK)
	apps, _ := got["apps"].([]any)
	if len(apps) != 1 || apps[0].(map[string]any)["status"] != "configuration error" {
		t.Errorf("apps = %v, want linear-work alone with a configuration error", got["apps"])
	}
}

func TestPluginsAnswersABrokenFileWithNoApps(t *testing.T) {
	f := newFixture(t)
	for name, yaml := range map[string]string{
		"not yaml":    "apiVersion: [\n",
		"wrong shape": "apiVersion: mate.plugins/v9\ninstances:\n" + linearBinding,
	} {
		writePlugins(t, f.ws.Root(), yaml)
		got := getPlugins(t, f.http.URL, fixtureProject, http.StatusOK)
		msg, _ := got["error"].(string)
		if !strings.Contains(msg, "plugins.yaml") {
			t.Errorf("%s: error = %v, want the file's problem", name, got["error"])
		}
		assertApps(t, got, "[]")
	}
}

func TestPluginsNamesAnUnknownProjectAndAnUnboundOne(t *testing.T) {
	f := newFixture(t)
	writePlugins(t, f.ws.Root(), "apiVersion: mate.plugins/v1\ninstances:\n"+
		strings.Replace(linearBinding, "      shop:", "      blog:", 1))
	assertApps(t, getPlugins(t, f.http.URL, fixtureProject, http.StatusOK), "[]")

	for _, p := range []string{"nosuch", ""} {
		got := getPlugins(t, f.http.URL, p, http.StatusNotFound)
		if msg, _ := got["error"].(string); msg == "" || !strings.Contains(msg, `"`+p+`"`) {
			t.Errorf("?project=%s: 404 reason %v does not name it", p, got["error"])
		}
	}

	resp, err := http.Post(f.http.URL+"/api/plugins?project="+fixtureProject, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Error("POST /api/plugins was accepted; this API has no write endpoint")
	}
}

// A Mate that is not running is judged by the harness a new start would
// get, workspace.yaml's defaults.mate_harness; with none recorded it is
// unknown, which is no reason to call the app unsupported. Whether a
// harness can take plugins is Deps.PluginHarness's answer, never a name
// this package knows.
func TestPluginsJudgesAStoppedMateByTheDefaultHarness(t *testing.T) {
	refuse := func(kind string) error {
		if kind == "codex" {
			return errors.New("the codex harness cannot give an agent MCP servers")
		}
		return nil
	}
	for _, tc := range []struct {
		def, harness, status, error string
		seam                        func(string) error
	}{
		{"codex", "codex", "harness unsupported", `"the codex harness cannot give an agent MCP servers"`, refuse},
		{"Claude", "claude", "configured", "null", refuse},
		{"", "unknown", "configured", "null", func(string) error { return errors.New("never asked about an unknown harness") }},
		{"codex", "codex", "configured", "null", nil},
	} {
		root := t.TempDir()
		ws, err := store.Init(root, store.Defaults{MateHarness: tc.def})
		if err != nil {
			t.Fatal(err)
		}
		if err := ws.AddProject("shop", store.ProjectConfig{}); err != nil {
			t.Fatal(err)
		}
		if ws, err = store.Open(root); err != nil {
			t.Fatal(err)
		}
		writer, err := db.Open(ws)
		if err != nil {
			t.Fatal(err)
		}
		_ = writer.Close()
		read, err := db.OpenRead(ws)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = read.Close() })
		s, err := New(Options{Workspace: ws, DB: read, Deps: Deps{
			Now: func() time.Time { return fixtureNow }, PluginHarness: tc.seam,
		}})
		if err != nil {
			t.Fatal(err)
		}
		ts := httptest.NewServer(s.Handler())
		t.Cleanup(ts.Close)
		writePlugins(t, root, "apiVersion: mate.plugins/v1\ninstances:\n"+linearBinding)

		want := strings.NewReplacer(`"status":"configured","harness":"claude"`,
			`"status":"`+tc.status+`","harness":"`+tc.harness+`"`,
			`"error":null`, `"error":`+tc.error).Replace(configuredLinear)
		assertApps(t, getPlugins(t, ts.URL, "shop", http.StatusOK), "["+want+"]")
	}
}
