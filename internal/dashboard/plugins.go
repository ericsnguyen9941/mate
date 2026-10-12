package dashboard

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/nguyenngocanh94/mate/internal/plugin"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// The three configuration statuses an app can have (docs/dashboard.md
// section 3b). None of them is a claim about login or tool availability:
// that is unknown to the dashboard, and the Mate's /mcp is where to look.
const (
	pluginConfigured  = "configured"
	pluginConfigError = "configuration error"
	pluginUnsupported = "harness unsupported"
)

// PluginsResponse is GET /api/plugins?project=<p>.
type PluginsResponse struct {
	envelope
	Project string      `json:"project"`
	Apps    []PluginApp `json:"apps"`
	// Error is why plugins.yaml as a whole is refused; empty when it is
	// valid.
	Error string `json:"error,omitempty"`
}

// PluginApp is one plugin instance bound to the project. The only URL in it
// is Open.URL: the MCP endpoint and any binding config stay out.
type PluginApp struct {
	Instance string      `json:"instance"`
	Name     string      `json:"name"`
	Status   string      `json:"status"`
	Harness  string      `json:"harness"`
	Open     *PluginOpen `json:"open"`
	MCP      *PluginMCP  `json:"mcp"`
	Skills   []string    `json:"skills"`
	Error    *string     `json:"error"`
}

type PluginOpen struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type PluginMCP struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// handlePlugins reads plugins.yaml on every request (the route is not
// cached: a file edit is not an event) and calls nothing: no network, no
// MCP, no exec, no write.
func (s *Server) handlePlugins(ctx context.Context, r *http.Request, env envelope) (any, error) {
	project := r.URL.Query().Get("project")
	if _, ok := s.ws.Project(project); !ok {
		return nil, errNotFound{reason: unknownProject(s.ws, project)}
	}
	harness, err := s.mateHarness(ctx, project)
	if err != nil {
		return nil, err
	}
	out := PluginsResponse{envelope: env, Project: project, Apps: []PluginApp{}}

	cfg, err := plugin.Load(s.ws.Root())
	if err != nil {
		out.Error = err.Error()
		var ce *plugin.ConfigError
		if !errors.As(err, &ce) || ce.Instance == "" {
			return out, nil
		}
		// One instance is wrong, and a Mate start refuses the whole file
		// for it. Each app of this project says so on its own row; the
		// one at fault carries the problem itself.
		for _, id := range boundInstances(s.ws.Root(), project) {
			msg := "the Mate will not start with this plugins.yaml: " + out.Error
			if id == ce.Instance {
				msg = out.Error
			}
			out.Apps = append(out.Apps, PluginApp{Instance: id, Name: id, Status: pluginConfigError,
				Harness: harness, Skills: []string{}, Error: &msg})
		}
		if ce.Project == project && !hasApp(out.Apps, ce.Instance) {
			msg := out.Error
			out.Apps = append(out.Apps, PluginApp{Instance: ce.Instance, Name: ce.Instance, Status: pluginConfigError,
				Harness: harness, Skills: []string{}, Error: &msg})
		}
		return out, nil
	}

	// Which harness can take MCP servers is the harness profile's to say,
	// not this package's (the harness ratchet): Deps.PluginHarness answers,
	// and an unknown harness or no answer calls nothing unsupported. Only
	// an app with an MCP server can be unsupported; a skills-only or
	// open-only app stays configured.
	var refusal string
	if harness != "unknown" && s.deps.PluginHarness != nil {
		if err := s.deps.PluginHarness(harness); err != nil {
			refusal = err.Error()
		}
	}
	var unsupported []string
	for _, in := range cfg.Instances {
		b, ok := in.Projects[project]
		if !ok || !in.Enabled || in.Manifest == nil {
			continue
		}
		m := in.Manifest
		app := PluginApp{Instance: in.ID, Name: m.Name, Status: pluginConfigured, Harness: harness, Skills: []string{}}
		if m.Open != nil {
			app.Open = &PluginOpen{Title: m.Open.Title, URL: b.OpenURL}
		}
		if m.MCP != nil {
			app.MCP = &PluginMCP{Name: in.ID, Type: m.MCP.Type}
			if refusal != "" {
				msg := refusal
				app.Status, app.Error = pluginUnsupported, &msg
				unsupported = append(unsupported, in.ID)
			}
		}
		for _, sk := range m.Skills {
			app.Skills = append(app.Skills, sk.Name)
		}
		out.Apps = append(out.Apps, app)
	}
	// A Mate start refuses MCP servers its harness cannot take, so while
	// one app needs them the Mate does not start and the configured apps
	// beside it do not reach it either. They stay configured and say so.
	if len(unsupported) > 0 {
		msg := "the Mate will not start while " + strings.Join(unsupported, ", ") + " needs MCP servers: " + refusal
		for i := range out.Apps {
			if out.Apps[i].Status == pluginConfigured {
				out.Apps[i].Error = &msg
			}
		}
	}
	return out, nil
}

// mateHarness is the harness the project's Mate runs, or would run: the
// running Mate's own (from `actor`, filled from mate.meta), else
// workspace.yaml's defaults.mate_harness, which a new start uses, else
// "unknown".
func (s *Server) mateHarness(ctx context.Context, project string) (string, error) {
	facts, _, err := s.actorFacts(ctx, timeline.MateActorID(project))
	if err != nil {
		return "", err
	}
	if facts.running() && facts.harness != "" {
		return strings.ToLower(facts.harness), nil
	}
	if def := strings.ToLower(strings.TrimSpace(s.ws.Defaults().MateHarness)); def != "" {
		return def, nil
	}
	return "unknown", nil
}

// boundInstances reads plugins.yaml without validating it, for the ids of
// the enabled instances that bind project, in file order. It is only asked
// once plugin.Load has refused the file, to say which apps that refusal
// stops; anything it cannot read gives no ids.
func boundInstances(root, project string) []string {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(plugin.ConfigFile)))
	if err != nil {
		return nil
	}
	var doc yaml.Node
	if yaml.Unmarshal(data, &doc) != nil || len(doc.Content) == 0 {
		return nil
	}
	instances := mapValue(doc.Content[0], "instances")
	if instances == nil {
		return nil
	}
	var ids []string
	for i := 0; i+1 < len(instances.Content); i += 2 {
		in := instances.Content[i+1]
		if e := mapValue(in, "enabled"); e != nil && e.Value == "false" {
			continue
		}
		if mapValue(mapValue(in, "projects"), project) != nil {
			ids = append(ids, instances.Content[i].Value)
		}
	}
	return ids
}

// mapValue is the value under key in a YAML mapping, or nil.
func mapValue(n *yaml.Node, key string) *yaml.Node {
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

func hasApp(apps []PluginApp, instance string) bool {
	for _, a := range apps {
		if a.Instance == instance {
			return true
		}
	}
	return false
}
