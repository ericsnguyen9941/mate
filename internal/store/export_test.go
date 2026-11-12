package store

// The external tests (package store_test) check these through the names
// below. Outside this package's tests nothing calls them, so the package
// does not export them (CON-243). Each line only forwards.

func (p Paths) CrewStopLock(project, crew string) string { return p.crewStopLock(project, crew) }
func (p Paths) JevLog() string                           { return p.jevLog() }
func (p Paths) MigrateLockFile() string                  { return p.migrateLockFile() }
func (p Paths) TUILockFile() string                      { return p.tuiLockFile() }

func (w *Workspace) Config() WorkspaceConfig             { return w.config() }
func (w *Workspace) CoreHolder() (CoreHolder, error)     { return w.coreHolder() }
func (w *Workspace) RelRepo(repo string) (string, error) { return w.relRepo(repo) }
func (w *Workspace) SaveConfig() error                   { return w.saveConfig() }
func (cfg ProjectConfig) SoleRepo() (RepoConfig, error)  { return cfg.soleRepo() }
func DefaultRepoName(repoPath string) string             { return defaultRepoName(repoPath) }
func ValidateRepoName(name string) error                 { return validateRepoName(name) }
