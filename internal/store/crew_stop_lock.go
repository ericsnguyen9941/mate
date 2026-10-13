package store

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

// CrewStopLock is `projects/<project>/crews/<crew>.stop.lock`, held for the
// whole of one `crew stop`: from before the agent is stopped until the meta
// says the task is over.
//
// The observer reads it. Between the two moments Herdr no longer has the
// agent while the meta still names it, and without the lock the observer
// reads that window as the agent having vanished: every stop of
// hellovietnam's crews up to 2026-10-13 left a `runtime_lost` incident open.
func (w *Workspace) CrewStopLock(project, crew string) string {
	return filepath.Join(w.CrewsDir(project), crew+".stop.lock")
}

// LockCrewStop waits for the crew's stop lock and returns what releases it.
// The wait ends with ctx. The lock is advisory and dies with the process, so
// a stop that crashed leaves nothing to clean up.
func (w *Workspace) LockCrewStop(ctx context.Context, project, crew string) (unlock func(), err error) {
	if err := os.MkdirAll(w.CrewsDir(project), 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(w.CrewStopLock(project, crew), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	for {
		ok, err := tryLockFile(f, true)
		if err != nil {
			f.Close()
			return nil, err
		}
		if ok {
			return func() {
				_ = unlockFile(f)
				_ = f.Close()
			}, nil
		}
		select {
		case <-ctx.Done():
			f.Close()
			return nil, ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

// CrewStopping reports whether a `crew stop` of this crew is in flight in
// any process. It creates nothing: a crew that was never stopped has no lock
// file, and nobody holds it.
func (w *Workspace) CrewStopping(project, crew string) (bool, error) {
	f, err := os.OpenFile(w.CrewStopLock(project, crew), os.O_RDWR, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer f.Close()
	ok, err := tryLockFile(f, true)
	if err != nil {
		return false, err
	}
	if ok {
		_ = unlockFile(f)
	}
	return !ok, nil
}
