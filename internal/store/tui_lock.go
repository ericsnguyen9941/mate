package store

import (
	"encoding/json"
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/nguyenngocanh94/mate/internal/platform"
)

// tuiLockName is the lock `mate mate-agent tui` holds while a headless Mate's
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
	// PID is `mate mate-agent tui`'s own process, the lock's holder: the one a
	// force close signals, and which passes the signal on to the TUI.
	PID int `json:"pid"`
	// Started is when the lock was taken.
	Started time.Time `json:"started"`
	// HarnessPID is the harness's process in the terminal while it runs;
	// `mate mate-agent tui` sets it back to zero once the harness has ended.
	HarnessPID int `json:"harness_pid,omitempty"`
	// Orphaned is set by TUIRunning, never written: `mate mate-agent tui` is
	// gone (killed outright) and the lock with it, but the harness it
	// started still runs, and may still write the Mate's session.
	Orphaned bool `json:"-"`
}

// TUILock is tui.lock held by this process.
type TUILock struct {
	f *os.File
}

// LockTUI takes the TUI lock without waiting and writes holder into it.
// While another process holds it the answer is ErrTUIOpen, and nothing is
// written. Like core.lock it is advisory and dies with the process; the
// file is never deleted.
//
// A lock left free by a `mate mate-agent tui` that was killed while its harness
// runs on is still refused (ErrTUIOpen): see TUIRunning.
func (w *Workspace) LockTUI(holder TUIHolder) (*TUILock, error) {
	f, err := os.OpenFile(w.tuiLockFile(), os.O_CREATE|os.O_RDWR, 0o644)
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
	var last TUIHolder
	if readHolder(w.tuiLockFile(), &last) == nil && harnessAlive(last) {
		l.Unlock()
		return nil, ErrTUIOpen
	}
	if err := l.Write(holder); err != nil {
		l.Unlock()
		return nil, err
	}
	return l, nil
}

// Write replaces what the lock says about its holder. The new record is
// written over the old one before the file is cut to its length, so a
// reader never finds it empty.
func (l *TUILock) Write(holder TUIHolder) error {
	data, err := json.Marshal(holder)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if _, err := l.f.WriteAt(data, 0); err != nil {
		return err
	}
	return l.f.Truncate(int64(len(data)))
}

// Unlock releases the lock. The file stays, saying who held it last.
func (l *TUILock) Unlock() {
	_ = unlockFile(l.f)
	_ = l.f.Close()
}

// TUIRunning reports whether a Mate's TUI is open now, and what its holder
// wrote: tui.lock is held, or it is not but the harness its last holder
// started still runs (Orphaned, `mate mate-agent tui` killed outright). The
// holder is zero when the lock is held but its record cannot be read.
func (w *Workspace) TUIRunning() (TUIHolder, bool, error) {
	held, err := lockHeld(w.tuiLockFile())
	if err != nil {
		return TUIHolder{}, false, err
	}
	var h TUIHolder
	if held {
		readHeldRecord(w.tuiLockFile(), &h)
		return h, true, nil
	}
	if readHolder(w.tuiLockFile(), &h) != nil || h.HarnessPID <= 0 {
		return TUIHolder{}, false, nil
	}
	if harnessAlive(h) {
		h.Orphaned = true
		return h, true, nil
	}
	// The orphan has ended: the record stops naming it, so a process that
	// is given its pid later is never taken for it.
	w.forgetHarness(h)
	return TUIHolder{}, false, nil
}

// forgetHarness writes h again without its harness pid, under the lock,
// unless another TUI took it meanwhile.
func (w *Workspace) forgetHarness(h TUIHolder) {
	f, err := os.OpenFile(w.tuiLockFile(), os.O_RDWR, 0o644)
	if err != nil {
		return
	}
	l := &TUILock{f: f}
	if ok, err := tryLockFile(f, true); err != nil || !ok {
		f.Close()
		return
	}
	defer l.Unlock()
	var now TUIHolder
	if readHolder(w.tuiLockFile(), &now) == nil && now.PID == h.PID && now.HarnessPID == h.HarnessPID {
		now.HarnessPID = 0
		_ = l.Write(now)
	}
}

// harnessAlive reports whether the harness h names still runs.
func harnessAlive(h TUIHolder) bool {
	if h.HarnessPID <= 0 {
		return false
	}
	err := platform.Signal(h.HarnessPID, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}
