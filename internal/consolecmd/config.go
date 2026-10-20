// Package consolecmd describes user-configured console commands. It knows
// paths and argv, not any particular editor or application.
package consolecmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Binding struct {
	Key     string   `yaml:"key"`
	Scope   string   `yaml:"scope"`
	Label   string   `yaml:"label"`
	Command []string `yaml:"command"`
	Cwd     string   `yaml:"cwd,omitempty"`
	// Mode is terminal (a managed host tab) or external (a GUI launcher).
	Mode string `yaml:"mode,omitempty"`
}

type Config struct {
	Bindings []Binding `yaml:"bindings"`
	// Dashboard says where the console serves its read-only dashboard.
	Dashboard Dashboard `yaml:"dashboard,omitempty"`
}

// Dashboard is the console's dashboard bind. Addr empty is the dashboard's
// own default, loopback on port 7777. A routable Addr (0.0.0.0:7777 for the
// LAN) also needs AllowRemote, because the timeline quotes everything the
// captain typed; without it the console says so on its status line and
// serves no dashboard rather than quietly serving loopback.
type Dashboard struct {
	Addr        string `yaml:"addr,omitempty"`
	AllowRemote bool   `yaml:"allow_remote,omitempty"`
}

// Load reads optional console.yaml. Unknown fields and extra YAML documents
// are errors, so a misspelled setting cannot silently select another launch.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)
	if err := d.Decode(&cfg); err != nil && err != io.EOF {
		return Config{}, fmt.Errorf("%s: %w", path, err)
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return Config{}, fmt.Errorf("%s: expected one YAML document", path)
	}
	if addr := cfg.Dashboard.Addr; strings.TrimSpace(addr) != "" {
		if _, _, err := net.SplitHostPort(addr); err != nil {
			return Config{}, fmt.Errorf("%s: dashboard.addr %q is not a host:port address", path, addr)
		}
	}
	for i := range cfg.Bindings {
		b := &cfg.Bindings[i]
		if b.Mode == "" {
			b.Mode = "terminal"
		}
		if b.Key == "" || strings.TrimSpace(b.Key) != b.Key ||
			(b.Scope != "crew" && b.Scope != "project") || b.Label == "" ||
			len(b.Command) == 0 || strings.TrimSpace(b.Command[0]) == "" ||
			(b.Mode != "terminal" && b.Mode != "external") {
			return Config{}, fmt.Errorf("%s: binding %d requires key, scope (crew/project), label, command argv and mode (terminal/external)", path, i+1)
		}
		// Check variable spelling before a key is ever pressed. Crew-only
		// variables cannot be used by project bindings.
		vars := map[string]string{"workspace_dir": "/workspace", "project_dir": "/workspace/project"}
		if b.Scope == "crew" {
			vars["crew_dir"], vars["report_path"], vars["worktree_dir"] = "/crew", "/crew/report.md", "/worktree"
		}
		if _, _, err := b.Resolve(vars); err != nil {
			return Config{}, fmt.Errorf("%s: binding %d: %w", path, i+1, err)
		}
	}
	return cfg, nil
}

// Resolve substitutes path variables without shell parsing or recursive
// expansion. A path containing spaces, $, or braces remains one argument.
func (b Binding) Resolve(vars map[string]string) ([]string, string, error) {
	args := make([]string, len(b.Command))
	for i, arg := range b.Command {
		s, err := expand(arg, vars)
		if err != nil {
			return nil, "", err
		}
		args[i] = s
	}
	cwd := b.Cwd
	if cwd == "" {
		cwd = "{project_dir}"
		if b.Scope == "crew" {
			cwd = "{crew_dir}"
		}
	}
	cwd, err := expand(cwd, vars)
	if err != nil {
		return nil, "", err
	}
	if !filepath.IsAbs(cwd) {
		cwd = filepath.Join(vars["workspace_dir"], cwd)
	}
	return args, cwd, nil
}

// unavailable says why a variable that exists for the scope has no value.
var unavailable = map[string]string{
	"report_path":  "the crew may not have written a report yet",
	"worktree_dir": "the crew has no worktree on disk; it may have finished and been cleaned up",
}

func expand(s string, vars map[string]string) (string, error) {
	var out strings.Builder
	for {
		i := strings.IndexByte(s, '{')
		if i < 0 {
			out.WriteString(s)
			return out.String(), nil
		}
		out.WriteString(s[:i])
		s = s[i+1:]
		end := strings.IndexByte(s, '}')
		if end < 0 {
			return "", fmt.Errorf("unclosed path variable")
		}
		name := s[:end]
		value, ok := vars[name]
		if !ok {
			return "", fmt.Errorf("unknown or unavailable path variable {%s}", name)
		}
		if value == "" {
			return "", fmt.Errorf("{%s} is unavailable; %s", name, unavailable[name])
		}
		out.WriteString(value)
		s = s[end+1:]
	}
}
