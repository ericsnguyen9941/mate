package box

import (
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/store"
)

func memoFixture(t *testing.T) *store.Workspace {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "shop", "shop"), 0o755); err != nil {
		t.Fatal(err)
	}
	w, err := store.Init(dir, store.Defaults{})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.AddProject("shop", store.ProjectConfig{Repos: []store.RepoConfig{{Path: "shop/shop"}}}); err != nil {
		t.Fatal(err)
	}
	return w
}

// coldLoad is Load with nothing held: what Load returns when it reads
// every file. No test of this package runs in parallel, so swapping the
// process's memo out for the call is safe.
func coldLoad(t *testing.T, w *store.Workspace) View {
	t.Helper()
	held := memo
	memo = &fileMemo{projects: map[string]*projectMemo{}}
	defer func() { memo = held }()
	v, err := Load(w, "shop")
	if err != nil {
		t.Fatal(err)
	}
	v.gen = 0
	return v
}

func memoLoad(t *testing.T, w *store.Workspace) View {
	t.Helper()
	v, err := Load(w, "shop")
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// CON-271: a Load none of whose files moved returns the View the last one
// built (the same Generation); after each kind of change the View is the
// one a Load reading every file builds.
func TestLoadKeepsTheViewUntilAFileMoves(t *testing.T) {
	w := memoFixture(t)
	at := time.Date(2026, 11, 12, 9, 0, 0, 0, time.UTC)
	steps := []struct {
		name   string
		change func() error
	}{
		{"nothing yet", func() error { return nil }},
		{"a crew at work", func() error {
			if err := w.WriteCrewMeta("shop", "k1", map[string]string{"state": "spawned"}); err != nil {
				return err
			}
			return w.AppendStatus("shop", "k1", "working: reading")
		}},
		{"a status grown", func() error { return w.AppendStatus("shop", "k1", "wait-mate: ready") }},
		{"a second crew", func() error { return w.AppendStatus("shop", "k2", "needs-decision: A or B?") }},
		{"a message", func() error {
			return w.AppendSent("shop", store.SentEntry{Time: at, Source: store.SourceMate, Target: store.CrewTarget("k2"), Text: "A"})
		}},
		{"an incident", func() error {
			return w.AppendIncident("shop", store.IncidentEntry{Time: at, Crew: "k1", Kind: "stale", State: store.IncidentOpen, Text: "quiet"})
		}},
		{"the incident resolved", func() error {
			return w.AppendIncident("shop", store.IncidentEntry{Time: at.Add(time.Minute), Crew: "k1", Kind: "stale", State: store.IncidentResolved})
		}},
		{"manual mode", func() error { return os.WriteFile(w.ManualHold("shop"), nil, 0o644) }},
		{"a meta rewritten closed", func() error {
			return w.WriteCrewMeta("shop", "k1", map[string]string{"state": "finished"})
		}},
		{"a meta gone", func() error { return os.Remove(w.CrewMeta("shop", "k1")) }},
		{"auto mode", func() error { return os.Remove(w.ManualHold("shop")) }},
		{"a status file gone", func() error { return os.Remove(w.CrewStatus("shop", "k2")) }},
		{"sent.log rewritten shorter", func() error { return os.WriteFile(w.SentLog("shop"), nil, 0o644) }},
	}
	var last uint64
	for _, step := range steps {
		if err := step.change(); err != nil {
			t.Fatalf("%s: %v", step.name, err)
		}
		first := memoLoad(t, w)
		if first.Generation() == 0 || first.Generation() == last {
			t.Fatalf("%s: generation %d after a change (last %d), want a new one", step.name, first.Generation(), last)
		}
		last = first.Generation()
		second := memoLoad(t, w)
		if second.Generation() != first.Generation() {
			t.Fatalf("%s: nothing moved, yet the View was built again (generation %d, then %d)", step.name, first.Generation(), second.Generation())
		}
		second.gen = 0
		if want := coldLoad(t, w); !reflect.DeepEqual(second, want) {
			t.Fatalf("%s: the memoized View\n%+v\nis not the cold one\n%+v", step.name, second, want)
		}
	}
	since, err := LoadSince(w, "shop", Cursor{Sent: 1})
	if err != nil {
		t.Fatal(err)
	}
	if since.Generation() != 0 {
		t.Fatalf("LoadSince's View has generation %d, want none", since.Generation())
	}
}

// A View the memo holds is shared: an append to any of its slices copies
// rather than writes into what another caller holds.
func TestAKeptViewHasNoRoomToAppendInto(t *testing.T) {
	w := memoFixture(t)
	for _, line := range []string{"working: a", "needs-decision: b?", "blocked: c"} {
		if err := w.AppendStatus("shop", "k1", line); err != nil {
			t.Fatal(err)
		}
	}
	v := memoLoad(t, w)
	if cap(v.Entries) != len(v.Entries) || cap(v.Attention) != len(v.Attention) || cap(v.ByCrew["k1"]) != len(v.ByCrew["k1"]) {
		t.Fatalf("caps %d/%d %d/%d %d/%d: a kept slice has room an append would write into",
			cap(v.Entries), len(v.Entries), cap(v.Attention), len(v.Attention), cap(v.ByCrew["k1"]), len(v.ByCrew["k1"]))
	}
}

// Loads that run at once while a crew writes share the memo safely (run
// under -race), and each View is whole: its Seqs run 0..n-1.
func TestConcurrentLoadsShareTheMemo(t *testing.T) {
	w := memoFixture(t)
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(stop)
		for i := 0; i < 50; i++ {
			_ = w.AppendStatus("shop", "k1", "working: step")
			_ = w.AppendSent("shop", store.SentEntry{Time: time.Now(), Source: store.SourceMate, Target: store.CrewTarget("k1"), Text: "go"})
		}
	}()
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				v, err := Load(w, "shop")
				if err != nil {
					t.Error(err)
					return
				}
				for i, e := range v.Entries {
					if e.Seq != i {
						t.Errorf("entry %d has Seq %d", i, e.Seq)
						return
					}
				}
			}
		}()
	}
	wg.Wait()
	got := memoLoad(t, w)
	got.gen = 0
	if want := coldLoad(t, w); !reflect.DeepEqual(got, want) {
		t.Fatal("after concurrent loads the memoized View is not the cold one")
	}
}
