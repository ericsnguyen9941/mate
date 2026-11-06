package store

import "path/filepath"

// projectsDir is `<root>/.mate/projects`.
func (p Paths) projectsDir() string { return filepath.Join(p.StateDir(), projectsDirName) }

// ProjectDir is `<root>/.mate/projects/<project>`.
func (p Paths) ProjectDir(project string) string {
	return filepath.Join(p.projectsDir(), project)
}

// ProjectHome is `<root>/<project>`, the project's own directory: the user's,
// not mate's state. Every repo of the project lives under it (docs/mvp.md
// section 3). Mate's state for the project stays in ProjectDir.
func (p Paths) ProjectHome(project string) string {
	return filepath.Join(p.root, project)
}

// ProjectFile is `projects/<project>/project.yaml`.
func (p Paths) ProjectFile(project string) string {
	return filepath.Join(p.ProjectDir(project), projectFileName)
}

// ProjectDoc is `projects/<project>/PROJECT.md`.
func (p Paths) ProjectDoc(project string) string {
	return filepath.Join(p.ProjectDir(project), projectDocName)
}

// ProjectCrewDoc is `projects/<project>/CREW.md`, the captain's standing
// rules for every Crew of one project, appended after the workspace's.
// Optional.
func (p Paths) ProjectCrewDoc(project string) string {
	return filepath.Join(p.ProjectDir(project), crewDocName)
}

// SentLog is `projects/<project>/sent.log`.
func (p Paths) SentLog(project string) string {
	return filepath.Join(p.ProjectDir(project), sentLogName)
}

// IncidentsLog is `projects/<project>/incidents.log`, the observer's
// append-only record of the incidents it opened and resolved (docs/mvp.md
// section 4b). Only internal/watch writes it; internal/box reads it.
func (p Paths) IncidentsLog(project string) string {
	return filepath.Join(p.ProjectDir(project), incidentsLogName)
}
