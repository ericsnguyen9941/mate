package store

import "os"

// Auto reports whether the project is in auto mode: always, unless the
// captain holds it in manual with the console's `m` key (Held). Auto is the
// default (CON-115, docs/mvp.md M19): a new project's Mate hears its crews
// from the first hand-back, with no key to press first. `mate/.auto` is
// still written by SetMode, for readers that date the switch by its mtime,
// but it no longer decides; a project with no flag at all used to sit in
// manual with no way out but the `m` key (signal, 2026-10-20).
func (w *Workspace) Auto(project string) bool {
	if err := ValidateProjectName(project); err != nil {
		return false
	}
	return !w.Held(project)
}

// Held reports whether the captain holds the project in manual mode, which
// is the presence of `mate/.manual`.
func (w *Workspace) Held(project string) bool {
	if err := ValidateProjectName(project); err != nil {
		return false
	}
	_, err := os.Stat(w.ManualHold(project))
	return err == nil
}

// SetMode is the captain's own choice of mode, the console's `m` key: auto
// turns `.auto` on and releases the hold; manual turns it off and holds it
// off, so the captain's choice is recorded. The auto daemon releases the hold
// through here too, once the captain has been away for QuietAfter. The hold
// is released before auto goes on and set before it goes off, so no reader
// between the two writes sees a state the captain did not choose.
func (w *Workspace) SetMode(project string, auto bool) error {
	if err := ValidateProjectName(project); err != nil {
		return err
	}
	if auto {
		if err := w.removeFlag(w.ManualHold(project)); err != nil {
			return err
		}
		return w.setAuto(project, true)
	}
	if err := w.writeFile(w.ManualHold(project), nil, 0o644); err != nil {
		return err
	}
	return w.setAuto(project, false)
}

func (w *Workspace) removeFlag(path string) error {
	if _, err := w.resolve(path); err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// setAuto creates or removes `mate/.auto`. It records the mode, it does not
// decide it: only the hold does (Auto). Choosing a mode is SetMode.
func (w *Workspace) setAuto(project string, on bool) error {
	if err := ValidateProjectName(project); err != nil {
		return err
	}
	path := w.AutoFlag(project)
	if !on {
		return w.removeFlag(path)
	}
	return w.writeFile(path, nil, 0o644)
}
