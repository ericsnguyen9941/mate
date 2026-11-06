package store_test

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// One TUI holds tui.lock at a time, whatever its project (CON-227); its
// record is JSON with `project`, `pid` and `started`, and the flag ends
// when the lock is let go.
func TestOneTUIHoldsTheWorkspaceAtATime(t *testing.T) {
	t.Parallel()
	w := newWorkspace(t)
	if _, held, err := w.TUIRunning(); err != nil || held {
		t.Fatalf("TUIRunning before any TUI = %v, %v; want not held", held, err)
	}
	started := time.Date(2026, 11, 6, 9, 30, 0, 0, time.UTC)
	lock, err := w.LockTUI(store.TUIHolder{Project: "shop", PID: 4242, Started: started})
	if err != nil {
		t.Fatalf("LockTUI: %v", err)
	}
	if _, err := w.LockTUI(store.TUIHolder{Project: "blog", PID: 1}); !errors.Is(err, store.ErrTUIOpen) {
		t.Fatalf("second LockTUI, another project = %v; want ErrTUIOpen", err)
	}
	if err := lock.Write(store.TUIHolder{Project: "shop", PID: 4242, Started: started, HarnessPID: 4243}); err != nil {
		t.Fatalf("Write: %v", err)
	}
	h, held, err := w.TUIRunning()
	if err != nil || !held || h.Project != "shop" || h.PID != 4242 || h.HarnessPID != 4243 || !h.Started.Equal(started) {
		t.Fatalf("TUIRunning = %+v, %v, %v; want shop's record", h, held, err)
	}
	data, err := os.ReadFile(w.TUILockFile())
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatalf("tui.lock is not JSON: %v\n%s", err, data)
	}
	for _, key := range []string{"project", "pid", "started"} {
		if _, ok := fields[key]; !ok {
			t.Errorf("tui.lock has no %q: %s", key, data)
		}
	}

	lock.Unlock()
	if _, held, err := w.TUIRunning(); err != nil || held {
		t.Fatalf("TUIRunning after Unlock = %v, %v; want not held", held, err)
	}
	again, err := w.LockTUI(store.TUIHolder{Project: "blog", PID: 7})
	if err != nil {
		t.Fatalf("LockTUI after Unlock: %v", err)
	}
	again.Unlock()
}

// `mate mate tui` killed outright lets the lock go, but the harness it
// started may run on and write the Mate's session: the flag stays up,
// Orphaned, until that process ends; then the record forgets its pid.
func TestATUIWhoseHolderWasKilledStaysOpenWhileItsHarnessRuns(t *testing.T) {
	t.Parallel()
	w := newWorkspace(t)
	harness := exec.Command("sleep", "30")
	if err := harness.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = harness.Process.Kill(); _ = harness.Wait() })
	lock, err := w.LockTUI(store.TUIHolder{Project: "shop", PID: 4242, Started: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if err := lock.Write(store.TUIHolder{Project: "shop", PID: 4242, HarnessPID: harness.Process.Pid}); err != nil {
		t.Fatal(err)
	}
	lock.Unlock() // the holder is gone; its harness is not

	h, held, err := w.TUIRunning()
	if err != nil || !held || !h.Orphaned || h.Project != "shop" || h.HarnessPID != harness.Process.Pid {
		t.Fatalf("TUIRunning with the harness alive = %+v, %v, %v; want it open, orphaned", h, held, err)
	}
	if _, err := w.LockTUI(store.TUIHolder{Project: "blog", PID: 7}); !errors.Is(err, store.ErrTUIOpen) {
		t.Fatalf("LockTUI beside an orphaned harness = %v; want ErrTUIOpen", err)
	}

	_ = harness.Process.Kill()
	_ = harness.Wait()
	if _, held, err := w.TUIRunning(); err != nil || held {
		t.Fatalf("TUIRunning once the harness ended = %v, %v; want closed", held, err)
	}
	var rec map[string]any
	data, _ := os.ReadFile(w.TUILockFile())
	if err := json.Unmarshal(data, &rec); err != nil || rec["harness_pid"] != nil {
		t.Fatalf("record after the orphan ended = %s; want no harness pid", data)
	}
	again, err := w.LockTUI(store.TUIHolder{Project: "blog", PID: 7})
	if err != nil {
		t.Fatalf("LockTUI once the orphan ended: %v", err)
	}
	again.Unlock()
}
