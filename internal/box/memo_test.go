package box_test

import (
	"os"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/store"
)

func texts(v box.View) []string {
	var out []string
	for _, e := range v.Entries {
		out = append(out, e.Text)
	}
	return out
}

func mustLoad(t *testing.T, w *store.Workspace) box.View {
	t.Helper()
	v, err := box.Load(w, "shop")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return v
}

// CON-259: a Load reads only the files that moved since the last one. The
// unmoved files are made unreadable, size and mtime kept, so a Load that
// reads them fails or loses their lines and one that trusts the memo does
// not.
func TestLoadDoesNotReadAFileThatHasNotMoved(t *testing.T) {
	w := newFixtureWorkspace(t)
	if err := w.AppendStatus("shop", "k3", "working: reading the ticket"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteCrewMeta("shop", "k3", map[string]string{"state": "finished"}); err != nil {
		t.Fatal(err)
	}
	if err := w.AppendSent("shop", store.SentEntry{
		Time: time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC), Source: store.SourceUser, Target: store.TargetMate, Text: "spawn k3",
	}); err != nil {
		t.Fatal(err)
	}
	first := mustLoad(t, w)

	paths := []string{w.CrewStatus("shop", "k3"), w.CrewMeta("shop", "k3"), w.SentLog("shop")}
	for _, path := range paths {
		if err := os.Chmod(path, 0); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		for _, path := range paths {
			_ = os.Chmod(path, 0o644)
		}
	})
	second := mustLoad(t, w)
	if got, want := texts(second), texts(first); len(got) != len(want) || !second.Closed["k3"] {
		t.Fatalf("a Load over unmoved files read %q (closed=%v), want %q (closed)", got, second.Closed["k3"], want)
	}
}

func TestLoadSeesWhatMovedSinceTheLastLoad(t *testing.T) {
	w := newFixtureWorkspace(t)
	if err := w.AppendStatus("shop", "k3", "working: reading the ticket"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteCrewMeta("shop", "k3", map[string]string{"state": "spawned"}); err != nil {
		t.Fatal(err)
	}
	if v := mustLoad(t, w); len(v.Entries) != 1 || v.Closed["k3"] {
		t.Fatalf("first Load: %q closed=%v", texts(v), v.Closed["k3"])
	}

	// Grown files are read from where the last read stopped.
	if err := w.AppendStatus("shop", "k3", "wait-mate: ready"); err != nil {
		t.Fatal(err)
	}
	if err := w.AppendSent("shop", store.SentEntry{
		Time: time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC), Source: store.SourceMate, Target: store.CrewTarget("k3"), Text: "thanks",
	}); err != nil {
		t.Fatal(err)
	}
	// A rewritten meta is parsed again.
	if err := w.WriteCrewMeta("shop", "k3", map[string]string{"state": "finished"}); err != nil {
		t.Fatal(err)
	}
	v := mustLoad(t, w)
	if got := texts(v); len(got) != 3 || got[0] != "thanks" || got[1] != "working: reading the ticket" || got[2] != "wait-mate: ready" {
		t.Fatalf("second Load: %q", got)
	}
	if !v.Closed["k3"] {
		t.Fatal("a meta rewritten to finished still reads open")
	}

	// A crew whose files are gone is gone from the next Load.
	if err := os.Remove(w.CrewStatus("shop", "k3")); err != nil {
		t.Fatal(err)
	}
	if v := mustLoad(t, w); len(v.ByCrew["k3"]) != 1 {
		t.Fatalf("after the status file was removed: %q", texts(v))
	}

	// LoadSince keeps its own cursor and is not the memo's.
	if err := w.AppendStatus("shop", "k3", "working: again"); err != nil {
		t.Fatal(err)
	}
	full := mustLoad(t, w)
	since, err := box.LoadSince(w, "shop", full.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	if len(since.Entries) != 0 {
		t.Fatalf("LoadSince from the latest cursor returned %q", texts(since))
	}
}
