// Package plugin resolves local plugin packages bound to projects in
// <workspace>/.mate/plugins.yaml. It only reads and validates files and returns
// what one project receives: MCP servers, skill contents and links. AddBinding
// is its one writer, of plugins.yaml alone; Doctor checks the file offline.
// It starts nothing, contacts nothing, and knows no particular application.
package plugin

import "strings"

const (
	// ConfigFile is the workspace binding file, relative to the workspace root.
	ConfigFile = ".mate/plugins.yaml"
	// ManifestFile is the package manifest, relative to the package root.
	ManifestFile = "mate-plugin.yaml"

	configAPIVersion   = "mate.plugins/v1"
	manifestAPIVersion = "mate.plugin/v1"
)

// MCPServer is one MCP server a project receives. Name is the instance id;
// Mate adds no alias or namespace. Type is http (URL set) or stdio (Command
// set; Cwd is absolute, the package root unless the manifest names one).
type MCPServer struct {
	Name    string
	Type    string
	URL     string
	Command []string
	Cwd     string
}

// Skill is one skill file of an instance with {instance_id} substituted.
// Name is the base name of the directory holding the skill file.
type Skill struct {
	Instance string
	Name     string
	Content  []byte
}

// Link is an instance's "open" link for one project. Key is the optional
// console shortcut from the project binding.
type Link struct {
	Instance string
	Title    string
	URL      string
	Key      string
}

// Resolved is what one project receives, ordered by instance order in
// plugins.yaml and then by skill order in each manifest.
type Resolved struct {
	Servers []MCPServer
	Skills  []Skill
	Links   []Link
}

// ConfigError is a configuration problem located by file and, where known,
// instance and project.
type ConfigError struct {
	File     string
	Instance string
	Project  string
	Msg      string
}

func (e *ConfigError) Error() string {
	parts := []string{e.File}
	if e.Instance != "" {
		parts = append(parts, "instance "+e.Instance)
	}
	if e.Project != "" {
		parts = append(parts, "project "+e.Project)
	}
	return strings.Join(append(parts, e.Msg), ": ")
}

// Resolve reads <workspaceRoot>/.mate/plugins.yaml and the manifest of every
// enabled instance bound to project, and returns what project receives. A
// missing plugins.yaml yields an empty Resolved. A problem in the file
// itself (YAML, apiVersion, a field outside any instance, a duplicate
// instance id) or in an instance bound to project (its fields, package,
// manifest, skills, open link, or project's binding and key) is a
// *ConfigError. A problem only another project's binding has, or an
// instance project is not bound to, is not: one project's broken binding
// keeps no other project's Mate from starting. Load checks everything.
func Resolve(workspaceRoot, project string) (Resolved, error) {
	cfg, err := readConfig(workspaceRoot, true, project)
	if err != nil {
		return Resolved{}, err
	}
	return cfg.Project(project), nil
}

// Project returns what project receives from an already loaded Config.
func (c Config) Project(project string) Resolved {
	var r Resolved
	for _, in := range c.Instances {
		b, ok := in.Projects[project]
		if !ok || !in.Enabled || in.Manifest == nil {
			continue
		}
		m := in.Manifest
		if m.MCP != nil {
			r.Servers = append(r.Servers, MCPServer{
				Name: in.ID, Type: m.MCP.Type, URL: m.MCP.URL,
				Command: append([]string(nil), m.MCP.Command...), Cwd: m.MCP.Cwd,
			})
		}
		for _, s := range m.Skills {
			r.Skills = append(r.Skills, Skill{Instance: in.ID, Name: s.Name, Content: append([]byte(nil), s.Content...)})
		}
		if m.Open != nil {
			r.Links = append(r.Links, Link{Instance: in.ID, Title: m.Open.Title, URL: b.OpenURL, Key: b.Key})
		}
	}
	return r
}
