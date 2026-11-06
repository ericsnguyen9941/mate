package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/nguyenngocanh94/mate/internal/memory"
)

// ErrProjectExists is returned by AddProject for a name already registered.
var ErrProjectExists = errors.New("store: project already registered")

// ErrNoProject is returned for a project that is not registered.
var ErrNoProject = errors.New("store: no such project")

// Projects are the registered projects, in the order workspace.yaml lists them.
func (w *Workspace) Projects() []ProjectRef {
	return append([]ProjectRef(nil), w.snapshot().Projects...)
}

// RequireProject refuses a name that is not a valid project name, or one
// workspace.yaml does not register (ErrNoProject). Every command and action
// that works on one project checks it this way first.
func (w *Workspace) RequireProject(name string) error {
	if err := ValidateProjectName(name); err != nil {
		return err
	}
	if _, ok := w.Project(name); !ok {
		return fmt.Errorf("%w: %s", ErrNoProject, name)
	}
	return nil
}

// Project returns the registration of one project.
func (w *Workspace) Project(name string) (ProjectRef, bool) {
	for _, p := range w.snapshot().Projects {
		if p.Name == name {
			return p, true
		}
	}
	return ProjectRef{}, false
}

// AddProject registers a project: it normalises its repos (possibly none),
// refuses a repo another project owns, creates `projects/<name>/` with the
// `mate/` and `crews/` directories, writes project.yaml and appends the
// project to workspace.yaml. Before any of that is written it makes the
// project's own directory, ProjectHome, if it is not there yet: the user's
// directory, where its repos live, kept as it is when it exists.
func (w *Workspace) AddProject(name string, cfg ProjectConfig) error {
	if err := ValidateProjectName(name); err != nil {
		return err
	}
	// One update from the check that the name is free to the save that
	// registers it, so no LoadConfig or other change of this handle lands
	// in between (CON-204).
	return w.updateConfig(func(next *WorkspaceConfig) error {
		if _, ok := w.Project(name); ok {
			return fmt.Errorf("%w: %s", ErrProjectExists, name)
		}
		normalised, err := w.normaliseProject(name, cfg)
		if err != nil {
			return err
		}
		if err := w.checkReposUnclaimed(name, normalised.Repos); err != nil {
			return err
		}
		for _, dir := range []string{w.ProjectHome(name), w.ProjectDir(name), w.MateDir(name), w.CrewsDir(name)} {
			if err := w.mkdirAll(dir); err != nil {
				return err
			}
		}
		if err := w.SaveProject(name, normalised); err != nil {
			return err
		}
		if err := w.seedOnce(w.ProjectCrewDoc(name), projectCrewDocSeed(name)); err != nil {
			return err
		}
		next.Projects = append(next.Projects, ProjectRef{Name: name})
		return nil
	})
}

// EnsureProjectDoc creates `projects/<name>/PROJECT.md` from
// memory.ProjectTemplate if it does not already exist. It never overwrites a
// file the captain or a crew scout has since edited.
func (w *Workspace) EnsureProjectDoc(name string) error {
	if err := ValidateProjectName(name); err != nil {
		return err
	}
	return w.seedOnce(w.ProjectDoc(name), memory.ProjectTemplate(name))
}

// RemoveProject drops a project from workspace.yaml. The state directory under
// `projects/` is left on disk: like `crews/<id>/`, it is history and only the
// user deletes it.
func (w *Workspace) RemoveProject(name string) error {
	return w.updateConfig(func(next *WorkspaceConfig) error {
		for i, p := range next.Projects {
			if p.Name == name {
				next.Projects = append(next.Projects[:i:i], next.Projects[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("%w: %s", ErrNoProject, name)
	})
}

// LoadProject reads `projects/<name>/project.yaml`.
func (w *Workspace) LoadProject(name string) (ProjectConfig, error) {
	if err := ValidateProjectName(name); err != nil {
		return ProjectConfig{}, err
	}
	path := w.ProjectFile(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ProjectConfig{}, fmt.Errorf("%w: %s", ErrNoProject, name)
		}
		return ProjectConfig{}, err
	}
	var file projectFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return ProjectConfig{}, fmt.Errorf("store: %s: %w", path, err)
	}
	cfg := file.ProjectConfig
	if file.LegacyRepo != "" && len(cfg.Repos) == 0 {
		cfg.Repos = []RepoConfig{{
			Name: DefaultRepoName(file.LegacyRepo), Path: file.LegacyRepo, DefaultBranch: file.LegacyDefaultBranch,
		}}
	}
	for i := range cfg.Repos {
		if cfg.Repos[i].DefaultBranch == "" {
			cfg.Repos[i].DefaultBranch = DefaultBranch
		}
	}
	return cfg, nil
}

// projectFile is project.yaml as read: the current shape plus the fields
// older files carry and SaveProject never writes again - the one-repo fields
// every project.yaml had before M9, which LoadProject turns into a one-entry
// Repos, and the `mode` and `yolo` that M19 removed (the Mate picks each
// crew's delivery and merges reviewed work itself), read and ignored.
type projectFile struct {
	ProjectConfig       `yaml:",inline"`
	LegacyRepo          string `yaml:"repo"`
	LegacyDefaultBranch string `yaml:"default_branch"`
	LegacyMode          string `yaml:"mode"`
	LegacyYolo          bool   `yaml:"yolo"`
}

// SaveProject writes `projects/<name>/project.yaml` atomically.
func (w *Workspace) SaveProject(name string, cfg ProjectConfig) error {
	if err := ValidateProjectName(name); err != nil {
		return err
	}
	normalised, err := w.normaliseProject(name, cfg)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(normalised)
	if err != nil {
		return err
	}
	return w.writeFile(w.ProjectFile(name), data, 0o644)
}

// normaliseProject fills the defaults and checks the repos per
// normaliseRepos.
func (w *Workspace) normaliseProject(project string, cfg ProjectConfig) (ProjectConfig, error) {
	repos, err := w.normaliseRepos(project, cfg.Repos)
	if err != nil {
		return ProjectConfig{}, err
	}
	cfg.Repos = repos
	return cfg, nil
}

// RelRepo turns a repository path - absolute, or relative to the workspace
// root - into the path stored in the configuration: relative to the root,
// slash separated, and proven to be inside the workspace. The repository need
// not be a direct child of the root, but it must not be the root itself and
// must not live under `.mate/`.
func (w *Workspace) RelRepo(repo string) (string, error) {
	if repo == "" {
		return "", fmt.Errorf("store: repo path is empty")
	}
	abs := repo
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(w.root, abs)
	}
	resolved, err := w.resolve(abs)
	if err != nil {
		return "", err
	}
	if resolved == w.root {
		return "", fmt.Errorf("store: repo %s is the workspace root", repo)
	}
	if within(w.StateDir(), resolved) {
		return "", fmt.Errorf("store: repo %s is inside %s", repo, StateDirName)
	}
	rel, err := filepath.Rel(w.root, resolved)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}
