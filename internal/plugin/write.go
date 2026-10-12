package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// AddRequest is one explicit binding of an instance to a project, as
// `mate plugin add` writes it.
type AddRequest struct {
	Instance string
	// Path is the package root: absolute, or relative to the workspace root.
	Path    string
	Project string
	// Config is the binding's opaque config; it replaces any config the
	// binding had.
	Config map[string]string
	// Key is the optional console shortcut; empty removes one.
	Key string
}

// AddResult says what AddBinding wrote.
type AddResult struct {
	// File is the absolute path of plugins.yaml.
	File string
	// Created is true when plugins.yaml did not exist before.
	Created bool
	// NewInstance is true when the instance was added to the file.
	NewInstance bool
	// Replaced is true when the project already had a binding of this
	// instance, now overwritten.
	Replaced bool
	// Disabled is true when the instance is enabled: false; the binding is
	// written but gives nothing until it is enabled.
	Disabled bool
	// Manifest is the package's manifest, skills included.
	Manifest *Manifest
}

// MissingConfigError is a binding whose config lacks keys the package's
// open.url substitutes.
type MissingConfigError struct {
	Instance string
	Project  string
	// URL is the open.url template, Manifest the file declaring it.
	URL      string
	Manifest string
	Keys     []string
}

func (e *MissingConfigError) Error() string {
	return fmt.Sprintf("instance %s: project %s: %s open.url %q needs config %s, which the binding does not set",
		e.Instance, e.Project, e.Manifest, e.URL, strings.Join(e.Keys, ", "))
}

// AddBinding writes req into <workspaceRoot>/.mate/plugins.yaml: it adds
// the instance when the file has none of that id, and adds or overwrites
// the binding of req.Project, keeping every other instance and project.
// An existing instance keeps its path; a different req.Path is refused.
// The result is validated exactly as a Mate start reads it (Load) before
// anything is written, so a refusal leaves the file untouched.
//
// The file is decoded to a YAML node tree, changed and encoded again with
// two-space indentation: instance and project order and most comments are
// kept, blank lines and hand-chosen indentation are not.
func AddBinding(workspaceRoot string, req AddRequest) (AddResult, error) {
	root, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return AddResult{}, err
	}
	stateDir := filepath.Join(root, filepath.Dir(filepath.FromSlash(ConfigFile)))
	if st, err := os.Stat(stateDir); err != nil || !st.IsDir() {
		return AddResult{}, fmt.Errorf("%s is not a workspace: it has no %s directory", root, filepath.Base(stateDir))
	}
	file := filepath.Join(root, filepath.FromSlash(ConfigFile))
	fail := func(format string, args ...any) error {
		return &ConfigError{File: file, Instance: req.Instance, Project: req.Project, Msg: fmt.Sprintf(format, args...)}
	}
	if !instanceIDRe.MatchString(req.Instance) {
		return AddResult{}, fail("instance id must match [a-z0-9][a-z0-9-]*")
	}
	if strings.TrimSpace(req.Project) == "" || strings.TrimSpace(req.Project) != req.Project {
		return AddResult{}, fail("project name must be non-empty without surrounding spaces")
	}
	if strings.TrimSpace(req.Path) == "" {
		return AddResult{}, fail("package path is required")
	}

	res := AddResult{File: file}
	data, err := os.ReadFile(file)
	mode := os.FileMode(0o644)
	switch {
	case errors.Is(err, os.ErrNotExist):
		res.Created = true
		data = []byte("apiVersion: " + configAPIVersion + "\ninstances: {}\n")
	case err != nil:
		return AddResult{}, &ConfigError{File: file, Msg: err.Error()}
	default:
		if st, err := os.Stat(file); err == nil {
			mode = st.Mode().Perm()
		}
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return AddResult{}, &ConfigError{File: file, Msg: err.Error()}
	}
	if len(doc.Content) == 0 {
		// An empty file: start it as a new one.
		if err := yaml.Unmarshal([]byte("apiVersion: "+configAPIVersion+"\ninstances: {}\n"), &doc); err != nil {
			return AddResult{}, err
		}
	}
	top := doc.Content[0]
	if top.Kind != yaml.MappingNode {
		return AddResult{}, &ConfigError{File: file, Msg: "the file is not a YAML mapping; fix it by hand"}
	}
	instances := mappingField(top, "instances")
	if instances.Kind != yaml.MappingNode {
		return AddResult{}, &ConfigError{File: file, Msg: "instances is not a mapping; fix it by hand"}
	}

	pkg := resolvePath(root, req.Path)
	inst := mappingValue(instances, req.Instance)
	if inst == nil {
		res.NewInstance = true
		inst = mapping()
		setField(inst, "path", str(pkg))
		setField(instances, req.Instance, inst)
	} else {
		if inst.Kind != yaml.MappingNode {
			return AddResult{}, fail("the instance entry is not a mapping; fix it by hand")
		}
		p := mappingValue(inst, "path")
		if p == nil || strings.TrimSpace(p.Value) == "" {
			return AddResult{}, fail("the instance has no path; fix it by hand")
		}
		old := resolvePath(root, p.Value)
		if !samePath(old, pkg) {
			return AddResult{}, fail("the instance already uses package %s, not %s; to move it, change instances.%s.path in the file by hand, or add the package under another instance id",
				old, pkg, req.Instance)
		}
		pkg = old
		if e := mappingValue(inst, "enabled"); e != nil && e.Value == "false" {
			res.Disabled = true
		}
	}

	m, err := loadManifest(pkg, req.Instance)
	if err != nil {
		if st, serr := os.Stat(pkg); serr != nil || !st.IsDir() {
			return AddResult{}, fail("package path %s is not a directory", pkg)
		}
		return AddResult{}, err
	}
	res.Manifest = m
	if m.Open != nil {
		if missing := missingKeys(m.Open.URL, req.Config); len(missing) > 0 {
			return AddResult{}, &MissingConfigError{Instance: req.Instance, Project: req.Project,
				URL: m.Open.URL, Manifest: m.File, Keys: missing}
		}
	}

	projects := mappingField(inst, "projects")
	if projects.Kind != yaml.MappingNode {
		return AddResult{}, fail("projects is not a mapping; fix it by hand")
	}
	binding := mapping()
	if len(req.Config) > 0 {
		cfg := mapping()
		keys := make([]string, 0, len(req.Config))
		for k := range req.Config {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			setField(cfg, k, str(req.Config[k]))
		}
		setField(binding, "config", cfg)
	}
	if req.Key != "" {
		setField(binding, "key", str(req.Key))
	}
	if len(binding.Content) == 0 {
		binding.Style = yaml.FlowStyle
	}
	res.Replaced = mappingValue(projects, req.Project) != nil
	setField(projects, req.Project, binding)

	var out bytes.Buffer
	enc := yaml.NewEncoder(&out)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return AddResult{}, err
	}
	if err := enc.Close(); err != nil {
		return AddResult{}, err
	}
	if _, err := parseConfig(root, file, out.Bytes(), true); err != nil {
		return AddResult{}, err
	}
	if err := writeAtomic(file, out.Bytes(), mode); err != nil {
		return AddResult{}, err
	}
	return res, nil
}

// missingKeys lists, sorted and once each, the {config.<key>} variables of
// tmpl that values does not set.
func missingKeys(tmpl string, values map[string]string) []string {
	seen := map[string]bool{}
	var keys []string
	_, _, _ = substitute(tmpl, func(key string) (string, bool) {
		if _, ok := values[key]; !ok && !seen[key] {
			seen[key] = true
			keys = append(keys, key)
		}
		return "", true
	})
	sort.Strings(keys)
	return keys
}

func resolvePath(root, p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Join(root, p)
	}
	return filepath.Clean(p)
}

// samePath compares two absolute paths, and failing that the paths they
// resolve to, so a symlinked spelling of the same package is the same.
func samePath(a, b string) bool {
	if a == b {
		return true
	}
	ra, errA := filepath.EvalSymlinks(a)
	rb, errB := filepath.EvalSymlinks(b)
	return errA == nil && errB == nil && ra == rb
}

func writeAtomic(file string, data []byte, mode os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(file), "."+filepath.Base(file)+".*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), file)
}

func mapping() *yaml.Node { return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"} }

func str(v string) *yaml.Node { return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v} }

// mappingValue is the value of key in mapping m, or nil.
func mappingValue(m *yaml.Node, key string) *yaml.Node {
	if m == nil || m.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}

// mappingField is the mapping under key in m, created at the end of m when
// key is absent or null. A flow-style mapping such as {} becomes block style.
func mappingField(m *yaml.Node, key string) *yaml.Node {
	v := mappingValue(m, key)
	if v == nil || (v.Kind == yaml.ScalarNode && v.Tag == "!!null") {
		n := mapping()
		setField(m, key, n)
		return n
	}
	if v.Kind == yaml.MappingNode && v.Style == yaml.FlowStyle {
		v.Style = 0
	}
	return v
}

// setField sets key in mapping m to v, in place when key exists so its
// position and comments stay, and appended otherwise.
func setField(m *yaml.Node, key string, v *yaml.Node) {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			old := m.Content[i+1]
			v.LineComment, v.FootComment = old.LineComment, old.FootComment
			m.Content[i+1] = v
			return
		}
	}
	m.Content = append(m.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}, v)
}
