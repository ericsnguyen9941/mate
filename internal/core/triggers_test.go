package core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/notifications"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// What wakes the engine's loops instead of a clock (CON-270): the outbox
// loop, the notifications and the reading of workspace.yaml.

// wakeSender is an outbox loop that counts its wakes.
type wakeSender struct{ wakes, nudges atomic.Int32 }

func (*wakeSender) Start(context.Context) error { return nil }
func (*wakeSender) Stop()                       {}
func (s *wakeSender) Wake()                     { s.wakes.Add(1) }
func (s *wakeSender) Nudge()                    { s.nudges.Add(1) }

// An `outbox` notice wakes the outbox loop; a `sent` or `mate` one nudges
// it (it runs only while a line waits); the rest leave it be. Every notice
// is a change the notifications look at.
func TestChangedWakesTheLoopsItConcerns(t *testing.T) {
	ws := shopWorkspace(t)
	s := &wakeSender{}
	e := New(Options{Workspace: ws, Sender: s, Notifier: senderFunc(func(context.Context, notifications.Message) error { return nil })})
	for _, tc := range []struct {
		kind          store.ChangeKind
		wakes, nudges int32
	}{
		{store.ChangeOutbox, 1, 0},
		{store.ChangeSent, 1, 1},
		{store.ChangeMate, 1, 2},
		{store.ChangeCrewStatus, 1, 2},
		{store.ChangeNews, 1, 2},
	} {
		<-drained(e.notify.wake)
		e.Changed(coreapi.Change{Project: "shop", Kind: string(tc.kind)})
		if s.wakes.Load() != tc.wakes || s.nudges.Load() != tc.nudges {
			t.Fatalf("after %s: %d wakes, %d nudges; want %d, %d", tc.kind, s.wakes.Load(), s.nudges.Load(), tc.wakes, tc.nudges)
		}
		if len(e.notify.wake) != 1 {
			t.Fatalf("a %s notice did not wake the notifications", tc.kind)
		}
	}
}

// drained empties a wake channel and returns a closed one.
func drained(c chan struct{}) <-chan struct{} {
	select {
	case <-c:
	default:
	}
	done := make(chan struct{})
	close(done)
	return done
}

// workspace.yaml is read on a `workspace` notice, and when the file's size
// or modification time changed - a write no notice told of - and not
// otherwise.
func TestConfigReloadsOnANoticeAndWhenTheFileChanges(t *testing.T) {
	ws := shopWorkspace(t)
	e := New(Options{Workspace: ws})
	r := e.reload
	r.reload()
	other, err := store.Open(ws.Root())
	if err != nil {
		t.Fatal(err)
	}
	if r.look() {
		t.Fatal("the file was read again though it had not changed")
	}

	if err := other.AddProject("lab", store.ProjectConfig{}); err != nil {
		t.Fatal(err)
	}
	if hasProject(ws, "lab") {
		t.Fatal("the shared workspace saw the project before any read")
	}
	if !r.look() || !hasProject(ws, "lab") {
		t.Fatal("a changed file was not read back")
	}
	if r.look() {
		t.Fatal("the file was read twice for one change")
	}

	if err := other.AddProject("dock", store.ProjectConfig{}); err != nil {
		t.Fatal(err)
	}
	e.Changed(coreapi.Change{Kind: string(store.ChangeWorkspace)})
	if !hasProject(ws, "dock") {
		t.Fatal("a workspace notice did not read the file back")
	}
	if r.look() {
		t.Fatal("the file a notice read was read again at the next look")
	}

	// A notice of another kind reads nothing.
	if err := other.AddProject("yard", store.ProjectConfig{}); err != nil {
		t.Fatal(err)
	}
	e.Changed(coreapi.Change{Project: "shop", Kind: string(store.ChangeOutbox)})
	if hasProject(ws, "yard") {
		t.Fatal("an outbox notice read workspace.yaml")
	}

	// The loop looks on its own clock.
	r.every = 5 * time.Millisecond
	if err := r.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	if err := other.AddProject("pier", store.ProjectConfig{}); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "the loop to read the changed file", func() bool { return hasProject(ws, "pier") })
}

func hasProject(ws *store.Workspace, name string) bool {
	for _, p := range ws.Projects() {
		if p.Name == name {
			return true
		}
	}
	return false
}

// countedSnapshot is the engine's snapshot as the notifications read it,
// counting the reads; each read is a new snapshot.
type countedSnapshot struct {
	ws    *store.Workspace
	reads atomic.Uint64
}

func (c *countedSnapshot) get(ctx context.Context) (query.Snapshot, uint64, error) {
	snap, err := query.Load(ctx, c.ws, testHarnesses, nil)
	return snap, c.reads.Add(1), err
}

func (c *countedSnapshot) stays(t *testing.T, want uint64, d time.Duration) {
	t.Helper()
	time.Sleep(d)
	if got := c.reads.Load(); got != want {
		t.Fatalf("%d ticks, want %d", got, want)
	}
}

func testNotifyLoop(ws *store.Workspace, snap *countedSnapshot, sender notifications.Sender) *notifyLoop {
	return &notifyLoop{root: ws.Root(), snapshot: snap.get, sender: sender, url: func() string { return "" },
		gap: 20 * time.Millisecond, safety: time.Hour, retry: time.Hour, wake: make(chan struct{}, 1)}
}

// The notifications tick as they start, then on a change - a notice, or a
// snapshot another reader had computed - no sooner than the gap after the
// last tick, and not on a clock.
func TestNotificationsTickOnAChange(t *testing.T) {
	ws := shopWorkspace(t)
	snap := &countedSnapshot{ws: ws}
	n := testNotifyLoop(ws, snap, senderFunc(func(context.Context, notifications.Message) error { return nil }))
	if err := n.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer n.Stop()
	waitFor(t, "the first tick", func() bool { return snap.reads.Load() == 1 })
	snap.stays(t, 1, 150*time.Millisecond)

	n.Wake()
	waitFor(t, "a tick after a change", func() bool { return snap.reads.Load() == 2 })
	// A burst is one tick.
	for range 5 {
		n.Wake()
	}
	waitFor(t, "a tick after a burst", func() bool { return snap.reads.Load() == 3 })
	snap.stays(t, 3, 150*time.Millisecond)

	// The snapshot the last tick read is no change; another one is.
	n.computed(n.seen.Load())
	snap.stays(t, 3, 100*time.Millisecond)
	n.computed(n.seen.Load() + 100)
	waitFor(t, "a tick after a snapshot computed", func() bool { return snap.reads.Load() == 4 })
}

// The 10-minute report needs no change: the loop ticks when it falls due.
func TestTheReportGoesOutOnItsDeadline(t *testing.T) {
	ws := shopWorkspace(t)
	path := filepath.Join(ws.StateDir(), "notifications.json")
	now := time.Now()
	receipts, err := json.Marshal(map[string]any{"delivered": map[string]string{}, "report_at": now.Add(-notifications.ReportInterval + 300*time.Millisecond), "work_from": now})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, receipts, 0o600); err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var sent []notifications.Message
	snap := &countedSnapshot{ws: ws}
	n := testNotifyLoop(ws, snap, senderFunc(func(_ context.Context, m notifications.Message) error {
		mu.Lock()
		defer mu.Unlock()
		sent = append(sent, m)
		return nil
	}))
	if err := n.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer n.Stop()
	waitFor(t, "the report", func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(sent) == 1
	})
	mu.Lock()
	report := sent[0]
	mu.Unlock()
	if !strings.Contains(report.Title, "báo cáo 10 phút") {
		t.Fatalf("sent %+v, want the report", report)
	}
	// The receipt is saved once the notifier took the report.
	waitFor(t, "the next report ten minutes after this one", func() bool {
		at, ok := notifications.NextReport(path)
		return ok && !at.Before(now.Add(notifications.ReportInterval))
	})
	if got := snap.reads.Load(); got != 2 {
		t.Fatalf("%d ticks, want the first and the one at the deadline", got)
	}
}
