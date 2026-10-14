package store_test

import (
	"context"
	"os"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/store"
)

func TestCrewStoppingIsTrueOnlyWhileTheStopLockIsHeld(t *testing.T) {
	w := newWorkspace(t)
	if err := w.AddProject("shop", store.ProjectConfig{Repos: []store.RepoConfig{{Path: "shop/shop"}}}); err != nil {
		t.Fatal(err)
	}

	stopping, err := w.CrewStopping("shop", "k3")
	if err != nil || stopping {
		t.Fatalf("CrewStopping before any stop = %v, %v; want false", stopping, err)
	}
	if _, err := os.Stat(w.CrewStopLock("shop", "k3")); !os.IsNotExist(err) {
		t.Fatalf("asking created the lock file: %v", err)
	}

	unlock, err := w.LockCrewStop(context.Background(), "shop", "k3")
	if err != nil {
		t.Fatalf("LockCrewStop: %v", err)
	}
	if stopping, err := w.CrewStopping("shop", "k3"); err != nil || !stopping {
		t.Fatalf("CrewStopping while held = %v, %v; want true", stopping, err)
	}
	if stopping, err := w.CrewStopping("shop", "k4"); err != nil || stopping {
		t.Fatalf("another crew's stop lock = %v, %v; want false", stopping, err)
	}

	unlock()
	if stopping, err := w.CrewStopping("shop", "k3"); err != nil || stopping {
		t.Fatalf("CrewStopping after unlock = %v, %v; want false", stopping, err)
	}
}

func TestLockCrewStopWaitsForTheStopInFlight(t *testing.T) {
	w := newWorkspace(t)
	if err := w.AddProject("shop", store.ProjectConfig{Repos: []store.RepoConfig{{Path: "shop/shop"}}}); err != nil {
		t.Fatal(err)
	}
	unlock, err := w.LockCrewStop(context.Background(), "shop", "k3")
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := w.LockCrewStop(ctx, "shop", "k3"); err == nil {
		t.Fatal("a second stop took the lock the first still holds")
	}
}
