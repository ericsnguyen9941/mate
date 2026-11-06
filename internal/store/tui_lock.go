package store

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// tuiLockName is the lock `mate mate tui` holds while a headless Mate's
// session is open in the harness's own TUI (CON-227). It is the workspace's
// one TUI flag: while it is held the host of that Mate runs no turn, and
// Office, Agent City and the console only look. One TUI per workspace, so
// a second one is refused whatever its project.
const tuiLockName = "tui.lock"

// ErrTUIOpen is what LockTUI returns while another TUI holds the lock.
var ErrTUIOpen = errors.New("store: a Mate's TUI is open on this workspace")

// TUIHolder is what the TUI holding tui.lock wrote into it.
type TUIHolder struct {
	// Project is the Mate whose session the TUI runs on.
	Project string `json:"project"`
	// PID is `mate mate tui`'s own process: the one a force close signals.
	PID int `json:"pid"`
	// Started is when the lock was taken.
	Started time.Time `json:"started"`
	// HarnessPID is the harness's process in the terminal, once started.
	HarnessPID int `json:"harness_pid,omitempty"`
}

// TUILockFile is `.mate/tui.lock`.
func (w *Workspace) TUILockFile() string { return filepath.Join(w.StateDir(), tuiLockName) }

// TUILock is tui.lock held by this process.
type TUILock struct {
	f *os.File
}

// LockTUI takes the TUI lock without waiting and writes holder into it.
// While another process holds it the answer is ErrTUIOpen, and nothing is
// written. Like core.lock it is advisory and dies with the process, so a
// TUI that was killed clears the flag; the file is never deleted.
func (w *Workspace) LockTUI(holder TUIHolder) (*TUILock, error) {
	f, err := os.OpenFile(w.TUILockFile(), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	ok, err := tryHolderLock(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	if !ok {
		f.Close()
		return nil, ErrTUIOpen
	}
	l := &TUILock{f: f}
	if err := l.Write(holder); err != nil {
		l.Unlock()
		return nil, err
	}
	return l, nil
}

// Write replaces what the lock says about its holder.
func (l *TUILock) Write(holder TUIHolder) error { return writeHolder(l.f, holder) }

// Unlock releases the lock. The file stays, saying who held it last.
func (l *TUILock) Unlock() {
	_ = unlockFile(l.f)
	_ = l.f.Close()
}

// TUIRunning reports whether a TUI holds tui.lock now, and what it wrote
// there; the holder is zero when its record cannot be read.
func (w *Workspace) TUIRunning() (TUIHolder, bool, error) {
	held, err := lockHeld(w.TUILockFile())
	if err != nil || !held {
		return TUIHolder{}, false, err
	}
	var h TUIHolder
	readHeldRecord(w.TUILockFile(), &h)
	return h, true, nil
}
