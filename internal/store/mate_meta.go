package store

import "strings"

// ReadMateMeta reads `mate/mate.meta`.
func (w *Workspace) ReadMateMeta(project string) (map[string]string, error) {
	if err := ValidateProjectName(project); err != nil {
		return nil, err
	}
	return ReadMeta(w.MateMeta(project))
}

// WriteMateMeta writes `mate/mate.meta` atomically, inside the workspace.
func (w *Workspace) WriteMateMeta(project string, meta map[string]string) error {
	if err := ValidateProjectName(project); err != nil {
		return err
	}
	data, err := encodeMeta(meta)
	if err != nil {
		return err
	}
	return w.writeFile(w.MateMeta(project), data, 0o644)
}

// MetaMateMode is the mate.meta key recording how the Mate runs, and
// MateModeHeadless its value for a headless host (CON-78). They live here
// because the outbox, stow and the console ask it as well as spawn.
const (
	MetaMateMode     = "mode"
	MateModeHeadless = "headless"
)

// MateHeadless reports whether the Mate of project runs headless now: its
// meta names a live agent and records the headless mode.
func (w *Workspace) MateHeadless(project string) bool {
	meta, err := w.ReadMateMeta(project)
	return err == nil && meta[MetaMateMode] == MateModeHeadless && strings.TrimSpace(meta["agent"]) != ""
}

// MetaMatePID and MetaMateSocket are the mate.meta keys of a headless
// Mate whose host runs as a background process of its own, outside Herdr
// (CON-225): the host's pid and the socket it answers on. Such a meta
// records no pane, tab or Herdr workspace.
const (
	MetaMatePID    = "pid"
	MetaMateSocket = "socket"
)

// MateHostDetached reports whether a mate.meta records a running headless
// Mate whose host is a process of its own (CON-225): the headless mode, an
// agent name, and no pane. Whether it is alive is the host's socket's to
// say, not Herdr's. A headless meta with a pane is a host an older binary
// started in its Herdr pane; Herdr still answers for it.
func MateHostDetached(meta map[string]string) bool {
	return meta[MetaMateMode] == MateModeHeadless && strings.TrimSpace(meta["agent"]) != "" && strings.TrimSpace(meta["pane"]) == ""
}

// MateRecorded reports whether a mate.meta records a Mate that was started
// and not stopped: a pane, or a host of its own (MateHostDetached). A stop
// drops both.
func MateRecorded(meta map[string]string) bool {
	return strings.TrimSpace(meta["pane"]) != "" || MateHostDetached(meta)
}
