package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var instanceIDRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// Config is a validated plugins.yaml with the manifest of every enabled
// instance loaded.
type Config struct {
	// File is the absolute path of plugins.yaml; empty when it is absent.
	File      string
	Instances []Instance
}

// Instance is one entry under instances, in file order.
type Instance struct {
	ID string
	// Path is the absolute package root.
	Path    string
	Enabled bool
	// Manifest is nil for a disabled instance: its package is not read.
	Manifest *Manifest
	Projects map[string]Binding
}

// Binding binds an instance to one project. Config values are opaque.
type Binding struct {
	Config map[string]string
	Key    string
	// OpenURL is the package's open.url with this binding's config
	// substituted; empty when the package declares no open link or the
	// instance is disabled.
	OpenURL string
}

type rawConfig struct {
	APIVersion string                 `yaml:"apiVersion"`
	Instances  map[string]rawInstance `yaml:"instances"`
}

type rawInstance struct {
	Path     string                 `yaml:"path"`
	Enabled  *bool                  `yaml:"enabled"`
	Projects map[string]*rawBinding `yaml:"projects"`
}

type rawBinding struct {
	Config map[string]string `yaml:"config"`
	Key    string            `yaml:"key"`
}

// Load reads and validates <workspaceRoot>/.mate/plugins.yaml and every
// enabled instance's manifest and skills. A missing plugins.yaml is an empty
// Config. Unknown fields, extra YAML documents and invalid values are
// *ConfigError naming the file and, where known, the instance and project.
func Load(workspaceRoot string) (Config, error) {
	return readConfig(workspaceRoot, true)
}

// ReadConfig reads <workspaceRoot>/.mate/plugins.yaml and checks its own
// shape without reading any package: every Manifest is nil and every
// OpenURL empty. It lists what the file says even when a package is
// missing or broken; Load is what a Mate start runs.
func ReadConfig(workspaceRoot string) (Config, error) {
	return readConfig(workspaceRoot, false)
}

func readConfig(workspaceRoot string, packages bool) (Config, error) {
	root, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return Config{}, err
	}
	file := filepath.Join(root, filepath.FromSlash(ConfigFile))
	data, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, &ConfigError{File: file, Msg: err.Error()}
	}
	return parseConfig(root, file, data, packages)
}

// parseConfig validates data as the plugins.yaml at file in the workspace
// root. With packages, every enabled instance's manifest, skills, open link
// and console key are checked as well.
func parseConfig(root, file string, data []byte, packages bool) (Config, error) {
	fail := func(instance, project, format string, args ...any) error {
		return &ConfigError{File: file, Instance: instance, Project: project, Msg: fmt.Sprintf(format, args...)}
	}

	lay, err := readLayout(file, data)
	if err != nil {
		return Config{}, err
	}
	var raw rawConfig
	if err := decodeStrict(data, &raw); err != nil {
		instance, project := lay.locate(err)
		return Config{}, fail(instance, project, "%v", err)
	}
	if raw.APIVersion != configAPIVersion {
		return Config{}, fail("", "", "apiVersion must be %s, got %q", configAPIVersion, raw.APIVersion)
	}

	cfg := Config{File: file}
	keyOwner := map[string]map[string]string{} // project -> key -> instance
	for _, id := range lay.ids {
		ri := raw.Instances[id]
		if !instanceIDRe.MatchString(id) {
			return Config{}, fail(id, "", "instance id must match [a-z0-9][a-z0-9-]*")
		}
		if strings.TrimSpace(ri.Path) == "" {
			return Config{}, fail(id, "", "path is required")
		}
		in := Instance{ID: id, Path: ri.Path, Enabled: ri.Enabled == nil || *ri.Enabled, Projects: map[string]Binding{}}
		if !filepath.IsAbs(in.Path) {
			in.Path = filepath.Join(root, in.Path)
		}
		in.Path = filepath.Clean(in.Path)
		if in.Enabled && packages {
			if st, err := os.Stat(in.Path); err != nil || !st.IsDir() {
				return Config{}, fail(id, "", "package path %s is not a directory", in.Path)
			}
			if in.Manifest, err = loadManifest(in.Path, id); err != nil {
				return Config{}, err
			}
		}
		for _, project := range sortedKeys(ri.Projects) {
			rb := ri.Projects[project]
			if rb == nil {
				rb = &rawBinding{}
			}
			if strings.TrimSpace(project) == "" || strings.TrimSpace(project) != project {
				return Config{}, fail(id, project, "project name must be non-empty without surrounding spaces")
			}
			b := Binding{Config: rb.Config, Key: rb.Key}
			if b.Config == nil {
				b.Config = map[string]string{}
			}
			if rb.Key != "" && strings.TrimSpace(rb.Key) != rb.Key {
				return Config{}, fail(id, project, "key %q has surrounding spaces", rb.Key)
			}
			if in.Enabled && packages {
				if b.Key != "" && in.Manifest.Open == nil {
					return Config{}, fail(id, project, "key %q needs an open link, but %s declares none", b.Key, filepath.Join(in.Path, ManifestFile))
				}
				if in.Manifest.Open != nil {
					if b.OpenURL, err = expandURL(in.Manifest.Open.URL, b.Config); err != nil {
						return Config{}, fail(id, project, "open.url: %v", err)
					}
				}
				if b.Key != "" {
					owners := keyOwner[project]
					if owners == nil {
						owners = map[string]string{}
						keyOwner[project] = owners
					}
					if other, dup := owners[b.Key]; dup {
						return Config{}, fail(id, project, "duplicate key %q, already used by instance %s", b.Key, other)
					}
					owners[b.Key] = id
				}
			}
			in.Projects[project] = b
		}
		cfg.Instances = append(cfg.Instances, in)
	}
	return cfg, nil
}

// decodeStrict decodes exactly one YAML document, refusing unknown fields
// and duplicate mapping keys.
func decodeStrict(data []byte, out any) error {
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)
	if err := d.Decode(out); err != nil && err != io.EOF {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return errors.New("expected one YAML document")
	}
	return nil
}

// layout records the order of instance ids in plugins.yaml and the lines
// each instance and project binding spans, so decode errors can name them.
type layout struct {
	ids   []string
	spans []span
}

type span struct {
	instance, project string
	start, end        int
}

func readLayout(file string, data []byte) (layout, error) {
	var doc yaml.Node
	var l layout
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return l, &ConfigError{File: file, Msg: err.Error()}
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return l, nil // decodeStrict reports the shape error
	}
	top := doc.Content[0]
	for i := 0; i+1 < len(top.Content); i += 2 {
		if top.Content[i].Value != "instances" || top.Content[i+1].Kind != yaml.MappingNode {
			continue
		}
		m := top.Content[i+1]
		seen := map[string]int{}
		for j := 0; j+1 < len(m.Content); j += 2 {
			k, v := m.Content[j], m.Content[j+1]
			if first, dup := seen[k.Value]; dup {
				return l, &ConfigError{File: file, Instance: k.Value,
					Msg: fmt.Sprintf("duplicate instance id (lines %d and %d)", first, k.Line)}
			}
			seen[k.Value] = k.Line
			l.ids = append(l.ids, k.Value)
			l.spans = append(l.spans, span{instance: k.Value, start: k.Line, end: lastLine(v)})
			if v.Kind != yaml.MappingNode {
				continue
			}
			for f := 0; f+1 < len(v.Content); f += 2 {
				if v.Content[f].Value != "projects" || v.Content[f+1].Kind != yaml.MappingNode {
					continue
				}
				pm := v.Content[f+1]
				for p := 0; p+1 < len(pm.Content); p += 2 {
					l.spans = append(l.spans, span{instance: k.Value, project: pm.Content[p].Value,
						start: pm.Content[p].Line, end: lastLine(pm.Content[p+1])})
				}
			}
		}
	}
	return l, nil
}

func lastLine(n *yaml.Node) int {
	line := n.Line
	for _, c := range n.Content {
		line = max(line, lastLine(c))
	}
	return line
}

var errLineRe = regexp.MustCompile(`line (\d+):`)

// locate names the instance and project whose entry contains the first line
// reported by a yaml error; both are empty outside any instance.
func (l layout) locate(err error) (instance, project string) {
	m := errLineRe.FindStringSubmatch(err.Error())
	if m == nil {
		return "", ""
	}
	line, _ := strconv.Atoi(m[1])
	for _, s := range l.spans {
		if line >= s.start && line <= s.end {
			instance = s.instance
			if s.project != "" {
				project = s.project
			}
		}
	}
	return instance, project
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
