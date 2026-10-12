package plugin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var skillNameRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// Manifest is a validated mate-plugin.yaml. Every block is optional.
type Manifest struct {
	// File is the absolute manifest path.
	File   string
	Name   string
	MCP    *MCP
	Skills []SkillFile
	Open   *Open
}

// MCP describes the package's MCP server. Cwd, if set, is absolute.
type MCP struct {
	Type    string
	URL     string
	Command []string
	Cwd     string
}

// SkillFile is one skill of the package. Path is absolute and inside the
// package; Content already has {instance_id} substituted.
type SkillFile struct {
	Path    string
	Name    string
	Content []byte
}

// Open is the package's link. URL is a template: only {config.<key>}
// literals are substituted per project binding.
type Open struct {
	Title string
	URL   string
}

type rawManifest struct {
	APIVersion string   `yaml:"apiVersion"`
	Name       string   `yaml:"name"`
	MCP        *rawMCP  `yaml:"mcp"`
	Skills     []string `yaml:"skills"`
	Open       *rawOpen `yaml:"open"`
}

type rawMCP struct {
	Type    string   `yaml:"type"`
	URL     string   `yaml:"url"`
	Command []string `yaml:"command"`
	Cwd     string   `yaml:"cwd"`
}

type rawOpen struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

// loadManifest reads and validates <pkg>/mate-plugin.yaml for instance id,
// including every skill file it lists.
func loadManifest(pkg, id string) (*Manifest, error) {
	file := filepath.Join(pkg, ManifestFile)
	fail := func(format string, args ...any) error {
		return &ConfigError{File: file, Instance: id, Msg: fmt.Sprintf(format, args...)}
	}
	data, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fail("package has no %s", ManifestFile)
	}
	if err != nil {
		return nil, fail("%v", err)
	}
	var raw rawManifest
	if err := decodeStrict(data, &raw); err != nil {
		return nil, fail("%v", err)
	}
	if raw.APIVersion != manifestAPIVersion {
		return nil, fail("apiVersion must be %s, got %q", manifestAPIVersion, raw.APIVersion)
	}
	if strings.TrimSpace(raw.Name) == "" {
		return nil, fail("name is required")
	}
	m := &Manifest{File: file, Name: raw.Name}

	if r := raw.MCP; r != nil {
		mcp := &MCP{Type: r.Type, URL: r.URL, Command: r.Command}
		switch r.Type {
		case "http":
			if len(r.Command) > 0 || r.Cwd != "" {
				return nil, fail("mcp: type http takes url, not command or cwd")
			}
			if err := checkURL(r.URL); err != nil {
				return nil, fail("mcp.url: %v", err)
			}
		case "stdio":
			if r.URL != "" {
				return nil, fail("mcp: type stdio takes command, not url")
			}
			if len(r.Command) == 0 || strings.TrimSpace(r.Command[0]) == "" {
				return nil, fail("mcp: type stdio requires command argv")
			}
			if r.Cwd != "" {
				mcp.Cwd = r.Cwd
				if !filepath.IsAbs(mcp.Cwd) {
					mcp.Cwd = filepath.Join(pkg, mcp.Cwd)
				}
				mcp.Cwd = filepath.Clean(mcp.Cwd)
				if st, err := os.Stat(mcp.Cwd); err != nil || !st.IsDir() {
					return nil, fail("mcp.cwd %s is not a directory", mcp.Cwd)
				}
			}
		default:
			return nil, fail("mcp.type must be http or stdio, got %q", r.Type)
		}
		m.MCP = mcp
	}

	seen := map[string]string{}
	for _, rel := range raw.Skills {
		s, err := loadSkill(pkg, rel, id)
		if err != nil {
			return nil, fail("skill %q: %v", rel, err)
		}
		if prev, dup := seen[s.Name]; dup {
			return nil, fail("skills %q and %q share the name %q", prev, rel, s.Name)
		}
		seen[s.Name] = rel
		m.Skills = append(m.Skills, s)
	}

	if o := raw.Open; o != nil {
		if strings.TrimSpace(o.Title) == "" {
			return nil, fail("open.title is required")
		}
		if o.URL == "" {
			return nil, fail("open.url is required")
		}
		// Check variable syntax now; the URL itself is checked here when it
		// has no variables, and per binding after substitution otherwise.
		_, vars, err := substitute(o.URL, func(string) (string, bool) { return "", true })
		if err == nil && vars == 0 {
			err = checkURL(o.URL)
		}
		if err != nil {
			return nil, fail("open.url: %v", err)
		}
		m.Open = &Open{Title: o.Title, URL: o.URL}
	}
	return m, nil
}

// loadSkill resolves rel inside pkg, refusing paths or symlinks that leave
// the package, and reads it with {instance_id} substituted.
func loadSkill(pkg, rel, id string) (SkillFile, error) {
	if rel == "" || filepath.IsAbs(rel) {
		return SkillFile{}, errors.New("must be a path relative to the package root")
	}
	clean := filepath.Clean(filepath.FromSlash(rel))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return SkillFile{}, errors.New("leaves the package")
	}
	name := filepath.Base(filepath.Dir(clean))
	if !skillNameRe.MatchString(name) {
		return SkillFile{}, errors.New("must sit in a named directory, such as skills/<name>/SKILL.md")
	}
	path := filepath.Join(pkg, clean)
	realPkg, err := filepath.EvalSymlinks(pkg)
	if err != nil {
		return SkillFile{}, err
	}
	real, err := filepath.EvalSymlinks(path)
	if errors.Is(err, os.ErrNotExist) {
		return SkillFile{}, fmt.Errorf("%s does not exist", path)
	}
	if err != nil {
		return SkillFile{}, err
	}
	if r, err := filepath.Rel(realPkg, real); err != nil || r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
		return SkillFile{}, fmt.Errorf("%s resolves outside the package", path)
	}
	if st, err := os.Stat(real); err != nil || !st.Mode().IsRegular() {
		return SkillFile{}, fmt.Errorf("%s is not a regular file", path)
	}
	data, err := os.ReadFile(real)
	if err != nil {
		return SkillFile{}, err
	}
	content, err := expandSkill(data, id)
	if err != nil {
		return SkillFile{}, err
	}
	return SkillFile{Path: path, Name: name, Content: content}, nil
}
