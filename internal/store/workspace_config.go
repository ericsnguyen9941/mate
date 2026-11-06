package store

import (
	"errors"
	"fmt"
	"maps"
	"os"

	"gopkg.in/yaml.v3"
)

// ErrEmptyWorkspaceFile is LoadConfig's refusal of a workspace.yaml that is
// empty: no document in it at all (zero bytes, or only whitespace or
// comments), `{}`, or a null document. Such a file reads as a workspace with
// no projects on layout 1, and the next save through the handle would write
// that over the file (CON-218). mate writes the file whole, so an empty one was emptied
// from outside: by an editor that truncates it before writing, a shell
// redirect, or such a write that hit a full disk.
var ErrEmptyWorkspaceFile = errors.New("store: workspace.yaml is empty")

// Defaults are the harnesses a project uses when it does not say otherwise.
// They are what workspace.yaml records and nothing more: an empty field
// means the workspace names none, and the binary's harness registry decides
// (harness.Registry.Default) where the harness is resolved. The store knows
// no harness by name.
type Defaults struct {
	MateHarness string `yaml:"mate_harness"`
	CrewHarness string `yaml:"crew_harness"`
}

// ProjectRef is one row of the project list in workspace.yaml: the project
// name. Its repos live in the project's own project.yaml (docs/mvp.md M9);
// the `repo:` a pre-M9 workspace.yaml carries here is ignored on read and
// dropped on the next save.
type ProjectRef struct {
	Name string `yaml:"name"`
}

// WorkspaceConfig is `.mate/workspace.yaml`.
type WorkspaceConfig struct {
	Version int `yaml:"version"`
	// Layout is the directory layout the workspace is on: 2 for a project
	// that is a directory under the root with its repos under it, absent
	// (read as 1) for repos beside `.mate/`, which only `mate migrate`
	// opens.
	Layout int `yaml:"layout,omitempty"`
	// Session is the Herdr session name of this workspace. It is derived from
	// the absolute path once, at Init, and then stored, so moving or
	// re-resolving the workspace never renames a live session.
	Session string `yaml:"session"`
	// Root is the absolute path the workspace had when it was last written.
	// It is what tells a console that the workspace was copied or moved: the
	// path in the file is no longer the path it was opened at, and the
	// absolute paths mate wrote into briefs are the old ones.
	Root     string   `yaml:"root,omitempty"`
	Defaults Defaults `yaml:"defaults"`
	// Crew is what every Crew of the workspace is launched with, whatever
	// its project. Absent is the defaults below.
	Crew CrewLaunch `yaml:"crew,omitempty"`
	// Lanes sizes the lanes repos declare (RepoConfig.Lane, CON-118):
	// how many crews may hold each at once. A lane not listed here is
	// DefaultLaneCap wide.
	Lanes    map[string]int `yaml:"lanes,omitempty"`
	Projects []ProjectRef   `yaml:"projects"`
}

// DefaultLaneCap is how many crews may hold a lane workspace.yaml does not
// size: two, because on the captain's machine a third concurrent heavy
// build made every crew wait (CON-118).
const DefaultLaneCap = 2

// LaneCap is how many crews may hold the lane at once.
func (w *Workspace) LaneCap(lane string) int {
	if n := w.snapshot().Lanes[lane]; n > 0 {
		return n
	}
	return DefaultLaneCap
}

// CrewLaunch is the `crew:` block of workspace.yaml.
type CrewLaunch struct {
	// Autocompact is the context size, in tokens, at which a Claude Crew
	// compacts its conversation (`--autocompact`). 0 or absent is
	// DefaultCrewAutocompact; a negative number launches without the flag,
	// so the harness's own threshold applies.
	Autocompact int `yaml:"autocompact,omitempty"`
	// ContextHigh is the context size, in tokens, past which the observer
	// files a `context_high` incident on a Crew that is still working
	// (CON-73): the Mate is told to open a fresh session at its next
	// hand-back rather than run a fix round in a bloated one. 0 or absent
	// is DefaultCrewContextHigh; a negative number turns the check off.
	ContextHigh int `yaml:"context_high,omitempty"`
}

// DefaultCrewAutocompact is a Crew's compaction window when workspace.yaml
// names none. Measured on hellovietnam 2026-10-06..14 (CON-68): Claude Crews
// averaged 162k tokens of context per call and never compacted, and the four
// that passed 300k (the largest 434k) cost 138M tokens between them. Claude
// Code 2.1.292 compacts at min(window*CLAUDE_AUTOCOMPACT_PCT_OVERRIDE/100,
// window-13000), so the override can only lower it: the window is the one
// knob. 160000 compacted at 147k, about 74k above a Crew's baseline; with
// CrewDeniedTools trimming that baseline, 200000 compacts at 187k, which
// leaves room for most tasks to finish without compacting and keeps every
// call under 200k, where long-context pricing would start.
const DefaultCrewAutocompact = 200000

// DefaultCrewContextHigh is the context size past which a working Crew is
// reported (`context_high`) when workspace.yaml names none. hellovietnam's
// top ten tasks (2026-10-15, CON-73) ran their fix rounds at 250-434k of
// context without compacting; the threshold sits above the autocompact one and
// below where those rounds started. It moved to 200000 with the autocompact
// window (CON-96), so a Crew is not reported before it could compact.
const DefaultCrewContextHigh = 200000

// snapshot is the configuration as last loaded or saved. Every reader shares
// it, so it is only read, never changed. A Workspace that neither Init nor
// OpenForMigrate made has none and reads as empty, as the zero Workspace did.
func (w *Workspace) snapshot() *WorkspaceConfig {
	if cfg := w.cfg.Load(); cfg != nil {
		return cfg
	}
	return &WorkspaceConfig{}
}

// clone is a copy of cfg that shares no slice or map with it. A field added
// to WorkspaceConfig that holds a reference must be copied here too
// (TestStoreConfigSharesNothingWithTheWorkspace).
func (cfg *WorkspaceConfig) clone() WorkspaceConfig {
	out := *cfg
	out.Projects = append([]ProjectRef(nil), cfg.Projects...)
	out.Lanes = maps.Clone(cfg.Lanes)
	return out
}

// Config is the loaded workspace.yaml. The copy is the caller's; changes reach
// disk only through the methods on Workspace.
func (w *Workspace) Config() WorkspaceConfig { return w.snapshot().clone() }

// Session is the Herdr session name of this workspace.
func (w *Workspace) Session() string { return w.snapshot().Session }

// RecordedRoot is the root workspace.yaml last recorded, "" for a workspace
// written before the file carried one.
func (w *Workspace) RecordedRoot() string { return w.snapshot().Root }

// SetRoot records the workspace's current root in workspace.yaml.
func (w *Workspace) SetRoot(root string) error {
	return w.updateConfig(func(next *WorkspaceConfig) error {
		next.Root = root
		return nil
	})
}

// SetSession gives the workspace another Herdr session name. It is for a
// workspace copied beside its original, which must never share the original's
// session; nothing else renames a session.
func (w *Workspace) SetSession(name string) error {
	return w.updateConfig(func(next *WorkspaceConfig) error {
		next.Session = name
		return nil
	})
}

// Layout is the workspace's directory layout: 2 when workspace.yaml says so,
// 1 for a file with no `layout:`.
func (w *Workspace) Layout() int {
	layout := w.snapshot().Layout
	if layout == 0 {
		return 1
	}
	return layout
}

// LayoutOld reports a workspace not yet on layout 2, which only
// OpenForMigrate opens.
func (w *Workspace) LayoutOld() bool { return w.Layout() < layoutProjectDirs }

// Defaults are the workspace-wide harness defaults.
func (w *Workspace) Defaults() Defaults { return w.snapshot().Defaults }

// CrewAutocompact is the compaction threshold a Crew launches with: the
// `crew.autocompact` of workspace.yaml, DefaultCrewAutocompact when it says
// none, and 0 (no flag) when it is negative.
func (w *Workspace) CrewAutocompact() int {
	switch n := w.snapshot().Crew.Autocompact; {
	case n < 0:
		return 0
	case n == 0:
		return DefaultCrewAutocompact
	default:
		return n
	}
}

// CrewContextHigh is the context size past which the observer reports a
// working Crew (`context_high`, CON-73): the `crew.context_high` of
// workspace.yaml, DefaultCrewContextHigh when it says none, and 0 (off)
// when it is negative.
func (w *Workspace) CrewContextHigh() int {
	switch n := w.snapshot().Crew.ContextHigh; {
	case n < 0:
		return 0
	case n == 0:
		return DefaultCrewContextHigh
	default:
		return n
	}
}

// LoadConfig re-reads workspace.yaml from disk. A file that is missing, does
// not parse or is empty (ErrEmptyWorkspaceFile) is an error, and the
// configuration loaded before stays.
func (w *Workspace) LoadConfig() error {
	w.cfgMu.Lock()
	defer w.cfgMu.Unlock()
	data, err := os.ReadFile(w.WorkspaceFile())
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s is missing", ErrNotWorkspace, w.WorkspaceFile())
		}
		return err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("store: %s: %w", w.WorkspaceFile(), err)
	}
	if emptyDocument(&doc) {
		return fmt.Errorf("%w: %s", ErrEmptyWorkspaceFile, w.WorkspaceFile())
	}
	var cfg WorkspaceConfig
	if err := doc.Decode(&cfg); err != nil {
		return fmt.Errorf("store: %s: %w", w.WorkspaceFile(), err)
	}
	if cfg.Session == "" {
		cfg.Session = SessionName(w.root)
	}
	w.cfg.Store(&cfg)
	return nil
}

// emptyDocument reports a parsed workspace.yaml that is empty: no document in
// it at all (zero bytes, or only whitespace or comments), `{}`, or a null
// document (`---`, `~`). A list or a string at the top is not empty; Decode
// refuses it, as before.
func emptyDocument(doc *yaml.Node) bool {
	top := doc
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		top = doc.Content[0]
	}
	switch top.Kind {
	case 0, yaml.DocumentNode:
		// yaml.v3 leaves the node zero for a file with no document in it.
		return true
	case yaml.MappingNode:
		return len(top.Content) == 0
	case yaml.ScalarNode:
		return top.ShortTag() == "!!null"
	}
	return false
}

// SaveConfig writes workspace.yaml atomically.
func (w *Workspace) SaveConfig() error {
	return w.updateConfig(func(*WorkspaceConfig) error { return nil })
}

// updateConfig is every change of workspace.yaml: change edits a copy of the
// configuration, the copy is written to the file, and only then is it the
// one readers see, so a write that fails leaves the configuration as it was.
// cfgMu is held from the read to the swap, so a LoadConfig beside it cannot
// land in between and undo the change (CON-204); work that must not be
// interleaved with the change goes in change too. An error from change saves
// nothing. change must not call LoadConfig or anything that updates the
// configuration: cfgMu is not reentrant.
func (w *Workspace) updateConfig(change func(next *WorkspaceConfig) error) error {
	w.cfgMu.Lock()
	defer w.cfgMu.Unlock()
	next := w.snapshot().clone()
	if err := change(&next); err != nil {
		return err
	}
	if next.Version == 0 {
		next.Version = workspaceVersion
	}
	data, err := yaml.Marshal(next)
	if err != nil {
		return err
	}
	if err := w.writeFile(w.WorkspaceFile(), data, 0o644); err != nil {
		return err
	}
	w.cfg.Store(&next)
	return nil
}
