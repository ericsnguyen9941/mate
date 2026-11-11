package store_test

import (
	"path/filepath"
	"slices"
	"sync"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// ChangeOf reads the Project, the Crew and the kind of a write from its
// path (CON-260), and nothing from a path that is no record of the feed's.
func TestChangeOfReadsTheLayout(t *testing.T) {
	root := filepath.FromSlash("/ws")
	sock := filepath.Join(root, ".mate", "core.sock")
	at := func(rel string) string { return filepath.Join(root, ".mate", filepath.FromSlash(rel)) }
	for _, tc := range []struct {
		rel  string
		want store.Change
		ok   bool
	}{
		{"workspace.yaml", store.Change{Kind: store.ChangeWorkspace}, true},
		{"projects/shop/project.yaml", store.Change{Project: "shop", Kind: store.ChangeProject}, true},
		{"projects/shop/sent.log", store.Change{Project: "shop", Kind: store.ChangeSent}, true},
		{"projects/shop/incidents.log", store.Change{Project: "shop", Kind: store.ChangeIncident}, true},
		{"projects/shop/mate/mate.meta", store.Change{Project: "shop", Kind: store.ChangeMate}, true},
		{"projects/shop/mate/.manual", store.Change{Project: "shop", Kind: store.ChangeMate}, true},
		{"projects/shop/mate/.outbox", store.Change{Project: "shop", Kind: store.ChangeOutbox}, true},
		{"projects/shop/mate/news.jsonl", store.Change{Project: "shop", Kind: store.ChangeNews}, true},
		{"projects/shop/mate/work.jsonl", store.Change{Project: "shop", Kind: store.ChangeWork}, true},
		{"projects/shop/mate/learnings/fix-cart.md", store.Change{Project: "shop", Kind: store.ChangeMate}, true},
		{"projects/shop/crews/fix-cart.meta", store.Change{Project: "shop", Crew: "fix-cart", Kind: store.ChangeCrewMeta}, true},
		{"projects/shop/crews/fix-cart.status", store.Change{Project: "shop", Crew: "fix-cart", Kind: store.ChangeCrewStatus}, true},
		{"projects/shop/crews/fix-cart.send.json", store.Change{Project: "shop", Crew: "fix-cart", Kind: store.ChangeCrew}, true},
		{"projects/shop/crews/fix-cart/brief.md", store.Change{Project: "shop", Crew: "fix-cart", Kind: store.ChangeCrew}, true},
		{"projects/shop/crews/fix-cart.stop.lock", store.Change{}, false},
		{"projects/shop/mate/.outbox.lock", store.Change{}, false},
		{"jev.log", store.Change{}, false},
		{"core.lock", store.Change{}, false},
		{"projects/shop", store.Change{}, false},
	} {
		got, ok := store.ChangeOf(at(tc.rel))
		if tc.ok {
			tc.want.Socket = sock
		}
		if ok != tc.ok || got != tc.want {
			t.Errorf("ChangeOf(%s) = %+v, %v; want %+v, %v", tc.rel, got, ok, tc.want, tc.ok)
		}
	}
	if _, ok := store.ChangeOf(filepath.Join(root, "shop", "README.md")); ok {
		t.Error("a file outside .mate is a change")
	}
}

// Every write the store makes tells the notifier once it is done: the
// appends, the atomic replaces and the removal of a flag. Not a parallel
// test: the notifier is the process's.
func TestStoreWritesTellTheNotifier(t *testing.T) {
	w := newWorkspace(t)
	if err := w.AddProject("shop", store.ProjectConfig{}); err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var got []store.Change
	store.SetChangeNotifier(func(c store.Change) {
		if c.Socket != w.CoreSocket() {
			return
		}
		mu.Lock()
		got = append(got, c)
		mu.Unlock()
	})
	t.Cleanup(func() { store.SetChangeNotifier(nil) })
	if err := w.AppendStatus("shop", "fix-cart", "working: cart"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteCrewMeta("shop", "fix-cart", map[string]string{"harness": "codex"}); err != nil {
		t.Fatal(err)
	}
	if err := w.SetMode("shop", false); err != nil {
		t.Fatal(err)
	}
	if err := w.SetMode("shop", true); err != nil {
		t.Fatal(err)
	}
	if err := store.WriteMeta(w.CrewMeta("shop", "fix-cart"), map[string]string{"harness": "pi"}); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, want := range []store.Change{
		{Socket: w.CoreSocket(), Project: "shop", Crew: "fix-cart", Kind: store.ChangeCrewStatus},
		{Socket: w.CoreSocket(), Project: "shop", Crew: "fix-cart", Kind: store.ChangeCrewMeta},
		{Socket: w.CoreSocket(), Project: "shop", Kind: store.ChangeMate},
	} {
		if !slices.Contains(got, want) {
			t.Errorf("no %+v among %+v", want, got)
		}
	}
	// .manual written, .auto removed, then .manual removed and .auto written.
	if n := countKind(got, store.ChangeMate); n < 4 {
		t.Errorf("%d mode changes told, want 4 or more: %+v", n, got)
	}
	if n := countKind(got, store.ChangeCrewMeta); n != 2 {
		t.Errorf("%d meta writes told, want 2", n)
	}
}

func countKind(cs []store.Change, k store.ChangeKind) int {
	n := 0
	for _, c := range cs {
		if c.Kind == k {
			n++
		}
	}
	return n
}
