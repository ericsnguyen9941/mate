package notifications

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

type recordingSender struct {
	messages []Message
	err      error
}

func (s *recordingSender) Send(_ context.Context, m Message) error {
	if s.err != nil {
		return s.err
	}
	s.messages = append(s.messages, m)
	return nil
}

func TestNewsReceiptsSurviveRestartAndResumeReportsAfterClear(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notifications.json")
	at := time.Date(2026, 10, 20, 10, 0, 0, 0, time.UTC)
	sender := &recordingSender{}
	snap := Snapshot{Complete: true, Projects: []Project{{Name: "shop", Summary: "2 crew làm"}}}
	tick := func(d time.Duration) {
		t.Helper()
		if err := Tick(context.Background(), path, snap, sender, at.Add(d)); err != nil {
			t.Fatal(err)
		}
	}
	tick(0)
	tick(9 * time.Minute)
	if len(sender.messages) != 0 {
		t.Fatal("reported before 10 minutes")
	}
	tick(10 * time.Minute)
	if len(sender.messages) != 1 || sender.messages[0].Title != "Mate · báo cáo 10 phút" {
		t.Fatalf("messages: %+v", sender.messages)
	}
	snap.Projects[0].News = []News{{Key: "shop/1", Revision: "opened", Title: "#1 · cần key"}}
	tick(11 * time.Minute)
	if len(sender.messages) != 2 {
		t.Fatal("news was not immediate")
	}
	// Tick loads receipts anew on every call, as a second console or restart does.
	tick(12 * time.Minute)
	tick(45 * time.Minute)
	if len(sender.messages) != 2 {
		t.Fatal("duplicate news or heartbeat while news waiting")
	}
	snap.Projects[0].News[0].Revision = "ask again"
	tick(46 * time.Minute)
	if len(sender.messages) != 3 {
		t.Fatal("new ask was lost")
	}
	if sender.messages[1].ID != sender.messages[2].ID {
		t.Fatal("same news should replace its delivered notification")
	}
	snap.Projects[0].News = nil
	tick(55 * time.Minute)
	if len(sender.messages) != 3 {
		t.Fatal("quiet interval should start at last outstanding news")
	}
	tick(56 * time.Minute)
	if len(sender.messages) != 4 {
		t.Fatal("report did not resume after news cleared")
	}
}

func TestFailedDeliveryIsRetriedWithoutLosingOtherReceipts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notifications.json")
	at := time.Now()
	snap := Snapshot{Complete: true, Projects: []Project{{Name: "shop", News: []News{{Key: "1", Revision: "a", Title: "one"}}}}}
	sender := &recordingSender{err: errors.New("permission denied")}
	if err := Tick(context.Background(), path, snap, sender, at); err == nil {
		t.Fatal("expected delivery error")
	}
	sender.err = nil
	if err := Tick(context.Background(), path, snap, sender, at.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	snap.Projects[0].News = append(snap.Projects[0].News, News{Key: "2", Revision: "a", Title: "two"})
	sender.err = errors.New("native app unavailable")
	if err := Tick(context.Background(), path, snap, sender, at.Add(2*time.Minute)); err == nil {
		t.Fatal("expected delivery error")
	}
	sender.err = nil
	if err := Tick(context.Background(), path, snap, sender, at.Add(3*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if len(sender.messages) != 2 {
		t.Fatalf("want exactly two successfully accepted news, got %+v", sender.messages)
	}
}

func TestUnreadableNewsSuppressReportsButKnownNewsStillNotify(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notifications.json")
	at := time.Now()
	sender := &recordingSender{}
	snap := Snapshot{Complete: false, Projects: []Project{{Name: "shop"}}}
	for _, d := range []time.Duration{0, time.Hour} {
		if err := Tick(context.Background(), path, snap, sender, at.Add(d)); err != nil {
			t.Fatal(err)
		}
	}
	if len(sender.messages) != 0 {
		t.Fatal("reported no news when news could not be read")
	}
	snap.Projects[0].News = []News{{Key: "1", Revision: "a", Title: "readable news"}}
	if err := Tick(context.Background(), path, snap, sender, at.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if len(sender.messages) != 1 {
		t.Fatal("a different unreadable project must not hide known news")
	}
}

func TestConcurrentConsolesDeliverOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notifications.json")
	at := time.Now()
	sender := &recordingSender{}
	snap := Snapshot{Complete: true, Projects: []Project{{Name: "shop", News: []News{{Key: "1", Revision: "a", Title: "one"}}}}}
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			if err := Tick(context.Background(), path, snap, sender, at); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if len(sender.messages) != 1 {
		t.Fatalf("delivered %d times", len(sender.messages))
	}
}

func TestCorruptReceiptDoesNotFloodOldNews(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notifications.json")
	if err := os.WriteFile(path, []byte("broken"), 0o600); err != nil {
		t.Fatal(err)
	}
	sender := &recordingSender{}
	if err := Tick(context.Background(), path, Snapshot{Complete: true}, sender, time.Now()); err == nil {
		t.Fatal("corrupt receipt must be visible")
	}
	if len(sender.messages) != 0 {
		t.Fatal("sent with corrupt receipts")
	}
}

// Work moves alert once per move from the first tick on, never replay the
// board's history, and do not count as news waiting for the captain.
func TestWorkMovesAlertOncePerMoveAfterTheFirstTick(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notifications.json")
	at := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	sender := &recordingSender{}
	old := Work{Key: "shop/work/1", At: at.Add(-time.Hour), Title: "Xong #1 · cũ"}
	snap := Snapshot{Complete: true, Projects: []Project{{Name: "shop", Summary: "yên", Work: []Work{old}}}}
	tick := func(d time.Duration) {
		t.Helper()
		if err := Tick(context.Background(), path, snap, sender, at.Add(d)); err != nil {
			t.Fatal(err)
		}
	}
	tick(0)
	if len(sender.messages) != 0 {
		t.Fatalf("history replayed: %+v", sender.messages)
	}
	added := Work{Key: "shop/work/2", At: at.Add(time.Minute), Title: "Nhận việc #2 · thêm", Body: "Mở Mate", URL: "http://x/office/#mate/shop/work"}
	snap.Projects[0].Work = []Work{old, added}
	tick(2 * time.Minute)
	tick(3 * time.Minute)
	if len(sender.messages) != 1 {
		t.Fatalf("messages: %+v", sender.messages)
	}
	m := sender.messages[0]
	if m.Title != "Mate · shop · Work" || m.Body != "Nhận việc #2 · thêm\nMở Mate" || m.URL != added.URL {
		t.Fatalf("message: %+v", m)
	}
	snap.Projects[0].Work[1].At = at.Add(4 * time.Minute)
	tick(5 * time.Minute)
	if len(sender.messages) != 2 || sender.messages[1].ID != m.ID {
		t.Fatalf("a new move should replace the row's notification: %+v", sender.messages)
	}
	// A work move delays the report like news, but does not hold it off.
	tick(14 * time.Minute)
	if len(sender.messages) != 2 {
		t.Fatal("reported within 10 minutes of a work move")
	}
	tick(15 * time.Minute)
	if len(sender.messages) != 3 || sender.messages[2].Title != "Mate · báo cáo 10 phút" {
		t.Fatalf("report did not follow: %+v", sender.messages)
	}
}
