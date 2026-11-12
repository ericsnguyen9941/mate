package query

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// memoWorkspace is a workspace under <dir>/ws, so a copy of it under
// another directory has the same name (the snapshot carries it).
func memoWorkspace(t *testing.T) *store.Workspace {
	t.Helper()
	root := filepath.Join(t.TempDir(), "ws")
	ws, err := store.Init(root, store.Defaults{MateHarness: "claude", CrewHarness: "codex"})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "shop", "shop", ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ws.AddProject("shop", store.ProjectConfig{Repos: []store.RepoConfig{{Path: "shop/shop"}}}); err != nil {
		t.Fatal(err)
	}
	return ws
}

// snapshotJSON is the snapshot as the core serves it, without its clock.
func snapshotJSON(t *testing.T, ws *store.Workspace) []byte {
	t.Helper()
	snap, err := Load(context.Background(), ws, testHarnesses, nil)
	if err != nil {
		t.Fatal(err)
	}
	snap.AsOf = time.Time{}
	b, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// coldJSON is snapshotJSON over a copy of ws that no load in this process
// has read - every file's mtime kept, since a status line's time is its
// file's - with the copy's root spelled as ws's.
func coldJSON(t *testing.T, ws *store.Workspace) []byte {
	t.Helper()
	root := filepath.Join(t.TempDir(), filepath.Base(ws.Root()))
	err := filepath.WalkDir(ws.Root(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(ws.Root(), path)
		to := filepath.Join(root, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(to, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(to, data, info.Mode().Perm()); err != nil {
			return err
		}
		return os.Chtimes(to, info.ModTime(), info.ModTime())
	})
	if err != nil {
		t.Fatal(err)
	}
	cold, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.ReplaceAll(snapshotJSON(t, cold), []byte(cold.Root()), []byte(ws.Root()))
}

// CON-271: a load through the memos is the snapshot a load that reads
// every file again would have built, byte for byte, after each kind of
// change a workspace sees.
func TestMemoizedLoadIsTheColdLoad(t *testing.T) {
	ws := memoWorkspace(t)
	at := time.Date(2026, 11, 12, 9, 0, 0, 0, time.UTC)
	mateDir := ws.MateDir("shop")
	if err := os.MkdirAll(mateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMateMeta("shop", map[string]string{"harness": "claude", "agent": "mate-shop", "pane": "p1"}); err != nil {
		t.Fatal(err)
	}
	steps := []struct {
		name   string
		change func() error
	}{
		{"nothing yet", func() error { return nil }},
		{"crews at work", func() error {
			for _, c := range []string{"k1", "k2", "k3"} {
				if err := ws.WriteCrewMeta("shop", c, map[string]string{"harness": "codex", "agent": "crew-" + c, "pane": "p-" + c, "state": "spawned", "task": "Task " + c}); err != nil {
					return err
				}
				if err := ws.AppendStatus("shop", c, "working: reading the brief"); err != nil {
					return err
				}
			}
			return ws.AppendSent("shop", store.SentEntry{Time: at, Source: store.SourceMate, Target: store.CrewTarget("k1"), Text: "go"})
		}},
		{"a question", func() error { return ws.AppendStatus("shop", "k2", "needs-decision: A or B?") }},
		{"an assign queued", func() error {
			return ws.UpdateOutbox("shop", at, func([]store.OutboxItem) ([]store.OutboxItem, bool, error) {
				return []store.OutboxItem{{ID: 1, At: at, Source: store.OutboxSourceAssign, Key: "crews/k2.status@27", Text: "resolve", State: store.OutboxQueued}}, true, nil
			})
		}},
		{"the assign sent", func() error {
			return ws.UpdateOutbox("shop", at, func(items []store.OutboxItem) ([]store.OutboxItem, bool, error) {
				items[0].State, items[0].SentAt = store.OutboxSent, at.Add(time.Minute)
				return items, true, nil
			})
		}},
		{"a second question", func() error { return ws.AppendStatus("shop", "k3", "needs-decision: keep or drop?") }},
		{"a digest delivered", func() error {
			return ws.WriteAutoCursor("shop", map[string]int64{ws.CrewStatus("shop", "k3"): 27})
		}},
		{"an incident", func() error {
			return ws.AppendIncident("shop", store.IncidentEntry{Time: at, Crew: "k3", Kind: "stale", State: store.IncidentOpen, Text: "quiet 20m"})
		}},
		{"the incident resolved", func() error {
			return ws.AppendIncident("shop", store.IncidentEntry{Time: at.Add(time.Minute), Crew: "k3", Kind: "stale", State: store.IncidentResolved})
		}},
		{"news and work", func() error {
			if _, err := matehost.OpenNews(matehost.NewsPath(mateDir), matehost.News{Title: "Epic 1 done", Crew: "k1"}, at); err != nil {
				return err
			}
			_, err := matehost.AddWork(matehost.WorkPath(mateDir), matehost.WorkAssigned, "Epic 1", "", at)
			return err
		}},
		{"a crew closed", func() error {
			return ws.WriteCrewMeta("shop", "k1", map[string]string{"harness": "codex", "agent": "crew-k1", "state": "finished"})
		}},
		{"manual mode", func() error { return os.WriteFile(ws.ManualHold("shop"), nil, 0o644) }},
		{"a status file gone", func() error { return os.Remove(ws.CrewStatus("shop", "k3")) }},
		{"the outbox gone", func() error { return os.Remove(ws.OutboxFile("shop")) }},
	}
	for _, step := range steps {
		if err := step.change(); err != nil {
			t.Fatalf("%s: %v", step.name, err)
		}
		// Twice: the first load after a change reads what moved, the
		// second finds nothing moved.
		for i := 0; i < 2; i++ {
			if got, want := snapshotJSON(t, ws), coldJSON(t, ws); !bytes.Equal(got, want) {
				t.Fatalf("%s, load %d: the memoized snapshot\n%s\nis not the cold one\n%s", step.name, i+1, got, want)
			}
		}
	}
}

// The outbox is parsed once per version: a load over an outbox whose size
// and mtime have not moved does not read it, and one rewritten or removed
// is read again.
func TestLoadDoesNotReparseAnUnmovedOutbox(t *testing.T) {
	ws := memoWorkspace(t)
	if err := ws.AppendStatus("shop", "k2", "needs-decision: A or B?"); err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 11, 12, 9, 0, 0, 0, time.UTC)
	write := func(state string) {
		t.Helper()
		if err := ws.UpdateOutbox("shop", at, func([]store.OutboxItem) ([]store.OutboxItem, bool, error) {
			return []store.OutboxItem{{ID: 1, At: at, Source: store.OutboxSourceAssign, Key: "crews/k2.status@0", Text: "resolve", State: state, SentAt: at}}, true, nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	assigned := func() string {
		t.Helper()
		return LoadBox(ws, "shop").Value.Inbox[0].Assigned.State
	}
	write(store.OutboxQueued)
	if got := assigned(); got != BoxAssignQueued {
		t.Fatalf("assign = %q, want queued", got)
	}
	unreadable(t, ws.OutboxFile("shop"))
	if got := assigned(); got != BoxAssignQueued {
		t.Fatalf("over an unmoved outbox the assign = %q: the outbox was read again", got)
	}
	if err := os.Chmod(ws.OutboxFile("shop"), 0o644); err != nil {
		t.Fatal(err)
	}
	write(store.OutboxSent)
	if got := assigned(); got != BoxAssignSent {
		t.Fatalf("after a rewrite the assign = %q, want sent", got)
	}
	if err := os.Remove(ws.OutboxFile("shop")); err != nil {
		t.Fatal(err)
	}
	if got := assigned(); got != "" {
		t.Fatalf("with no outbox the assign = %q, want none", got)
	}
}

// The news log and the work board are parsed once per version too.
func TestLoadDoesNotReparseUnmovedNewsOrWork(t *testing.T) {
	ws := memoWorkspace(t)
	mateDir := ws.MateDir("shop")
	if err := os.MkdirAll(mateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 11, 12, 9, 0, 0, 0, time.UTC)
	if _, err := matehost.OpenNews(matehost.NewsPath(mateDir), matehost.News{Title: "one"}, at); err != nil {
		t.Fatal(err)
	}
	if _, err := matehost.AddWork(matehost.WorkPath(mateDir), matehost.WorkAssigned, "Epic 1", "", at); err != nil {
		t.Fatal(err)
	}
	counts := func() (int, int) {
		t.Helper()
		snap, err := Load(context.Background(), ws, testHarnesses, nil)
		if err != nil {
			t.Fatal(err)
		}
		m := snap.Projects[0].Mate
		return len(m.News.Value), len(m.Work.Value)
	}
	if n, w := counts(); n != 1 || w != 1 {
		t.Fatalf("news %d, work %d; want 1 and 1", n, w)
	}
	unreadable(t, matehost.NewsPath(mateDir))
	unreadable(t, matehost.WorkPath(mateDir))
	if n, w := counts(); n != 1 || w != 1 {
		t.Fatalf("over unmoved files news %d, work %d; want 1 and 1: they were read again", n, w)
	}
	for _, p := range []string{matehost.NewsPath(mateDir), matehost.WorkPath(mateDir)} {
		if err := os.Chmod(p, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := matehost.OpenNews(matehost.NewsPath(mateDir), matehost.News{Title: "two"}, at.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err := matehost.AddWork(matehost.WorkPath(mateDir), matehost.WorkAssigned, "Epic 2", "", at.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if n, w := counts(); n != 2 || w != 2 {
		t.Fatalf("after an append news %d, work %d; want 2 and 2", n, w)
	}
}

// Loads that run at once - the core's snapshot, FreshSnapshot, the
// dashboard's LoadBox - while the files move share the memos safely (run
// under -race) and each sees a whole box.
func TestConcurrentLoadsShareTheMemos(t *testing.T) {
	ws := memoWorkspace(t)
	if err := ws.AppendStatus("shop", "k1", "working: start"); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			_ = ws.AppendStatus("shop", "k1", "working: step")
			_ = ws.AppendStatus("shop", "k2", "needs-decision: which?")
		}
		close(stop)
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
				snap, err := Load(context.Background(), ws, testHarnesses, nil)
				if err != nil {
					t.Error(err)
					return
				}
				b := snap.Projects[0].Box.Value
				if len(b.Entries) == 0 || b.Entries[len(b.Entries)-1].Seq != len(b.Entries)-1 {
					t.Errorf("a box of %d entries does not end at its last Seq", len(b.Entries))
					return
				}
				_ = LoadBox(ws, "shop")
			}
		}()
	}
	wg.Wait()
	if got, want := snapshotJSON(t, ws), coldJSON(t, ws); !bytes.Equal(got, want) {
		t.Fatal("after concurrent loads the memoized snapshot is not the cold one")
	}
}
