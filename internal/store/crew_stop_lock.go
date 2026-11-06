package store

import (
	"context"
	"os"
	"time"
)

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
