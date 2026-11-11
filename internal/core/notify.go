package core

import (
	"context"
	"fmt"
	neturl "net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/loopstat"
	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/notifications"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// NativeNotifier is this machine's notifier, asked once: nil on a machine
// that cannot notify (Linux without notify-send, a system mate has no
// notifier for), which then runs no notification loop at all.
func NativeNotifier() notifications.Sender {
	if notifications.Supported() {
		return notifications.Native{}
	}
	return nil
}

// notifyLoop sends the captain a native notification for news waiting on
// them, and reports recorded state every 10 minutes without news
// (internal/notifications), for as long as the engine runs, including while
// the console's terminal is backgrounded.
//
// It reads the engine's shared snapshot rather than loading the workspace
// itself: a second query.Load every five seconds cost as much again as the
// snapshot, which the console, Office and City already keep fresh
// (CON-256). The snapshot's Herdr questions are bounded (livenessBound), so
// a Herdr that hangs delays no notification past them. It never waits on
// the Mate to deliver a notification.
type notifyLoop struct {
	root     string
	snapshot func(context.Context) (query.Snapshot, error)
	sender   notifications.Sender
	// url is the dashboard's, read as the loop starts: a click opens the
	// captain's desk there.
	url func() string
	// stat records each tick past the shared snapshot it waits for, which
	// is the snapshot's own cost (CON-263); nil records nothing.
	stat *loopstat.Loop

	// cancel and done are set by Start and read by Stop: the engine calls
	// each once, Stop on a goroutine it starts after Start has returned.
	cancel context.CancelFunc
	done   chan struct{}

	mu     sync.Mutex
	err    error
	status string
}

// Start begins the loop. It never fails the engine: a workspace whose
// `.mate/.env` says MATE_NOTIFICATIONS=off runs no loop and says nothing,
// and one that cannot run it keeps the reason for Engine.Notifications.
func (n *notifyLoop) Start(ctx context.Context) error {
	err := n.start(ctx)
	n.mu.Lock()
	defer n.mu.Unlock()
	n.err = err
	return nil
}

// notificationsOn is `.mate/.env`'s MATE_NOTIFICATIONS: on unless it says
// off.
func notificationsOn(ws *store.Workspace) (bool, error) {
	env, err := ws.LoadEnv()
	if err != nil {
		return false, err
	}
	switch value := strings.ToLower(strings.TrimSpace(env["MATE_NOTIFICATIONS"])); value {
	case "off", "false", "0", "no":
		return false, nil
	case "", "on", "true", "1", "yes":
		return true, nil
	default:
		return false, fmt.Errorf("MATE_NOTIFICATIONS must be on or off")
	}
}

func (n *notifyLoop) start(ctx context.Context) error {
	ws, err := store.Open(n.root)
	if err != nil {
		return err
	}
	if on, err := notificationsOn(ws); err != nil || !on {
		return err
	}
	path, err := ws.Resolve(filepath.Join(ws.StateDir(), "notifications.json"))
	if err != nil {
		return err
	}
	if _, err := ws.Resolve(path + ".lock"); err != nil {
		return err
	}
	dashboardURL := n.url()
	ctx, cancel := context.WithCancel(ctx)
	// The goroutine closes the channel made here, never n.done (the
	// CON-207 race).
	done := make(chan struct{})
	n.cancel, n.done = cancel, done
	go func() {
		defer close(done)
		timer := time.NewTimer(0)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
			}
			snap, err := n.snapshot(ctx)
			if err == nil {
				start := n.stat.Begin()
				err = notifications.Tick(ctx, path, notificationSnapshot(ctx, ws, snap, dashboardURL, time.Now()), n.sender, time.Now())
				n.stat.End(start)
			}
			n.setStatus(err)
			delay := 5 * time.Second
			if err != nil {
				delay = time.Minute
			}
			timer.Reset(delay)
		}
	}()
	return nil
}

// Stop ends the loop and waits for a tick in flight.
func (n *notifyLoop) Stop() {
	if n.cancel == nil {
		return
	}
	n.cancel()
	<-n.done
	n.cancel, n.done = nil, nil
}

func (n *notifyLoop) setStatus(err error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.status = ""
	if err != nil {
		n.status = "notifications: " + err.Error()
	}
}

func (n *notifyLoop) read() (string, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.status, n.err
}

func notificationSnapshot(ctx context.Context, ws *store.Workspace, snap query.Snapshot, dashboardURL string, now time.Time) notifications.Snapshot {
	out := notifications.Snapshot{Complete: true}
	// The dashboard's URL can be the loopback one followed by the LAN one
	// (dashboardPart); the click opens on this machine, so loopback.
	office := ""
	if fields := strings.Fields(dashboardURL); len(fields) > 0 {
		office = strings.TrimRight(fields[0], "/") + "/office/"
		out.URL = office + "#captain"
	}
	for _, p := range snap.Projects {
		np := notifications.Project{Name: p.Name}
		if p.Mate.News.State == query.Unknown {
			out.Complete = false
		}
		for _, n := range p.Mate.OwedNews() {
			np.News = append(np.News, notifications.News{Key: fmt.Sprintf("%s/news/%d", p.ProjectID, n.ID), Revision: n.At.UTC().Format(time.RFC3339Nano), Title: fmt.Sprintf("#%d · %s", n.ID, n.Title)})
		}
		np.Work = workMoves(p, office)
		np.Summary = notificationSummary(ctx, ws, p, now)
		out.Projects = append(out.Projects, np)
	}
	return out
}

// workMoves is each open work row's latest move the captain hears about
// (CON-158): put on the board, done for the captain to OK, or reopened. A
// progress line does not alert, nor does a wait: the topic it waits on is
// news, which alerts on its own. A click opens the Mate's Work tab.
func workMoves(p query.ProjectNode, office string) []notifications.Work {
	if !p.Mate.Work.IsKnown() {
		return nil
	}
	url := ""
	if office != "" {
		url = office + "#mate/" + neturl.PathEscape(p.ProjectID) + "/work"
	}
	var out []notifications.Work
	for _, w := range p.Mate.Work.Value {
		if !w.Open() {
			continue
		}
		var move query.WorkNote
		for _, n := range w.Notes {
			if n.Op == "add" || n.Op == "done" || n.Op == "reopen" {
				move = n
			}
		}
		if move.Op == "" {
			continue
		}
		var title string
		switch {
		case move.Op == "done":
			title = fmt.Sprintf("Xong #%d · %s, chờ anh OK", w.ID, w.Title)
		case move.Op == "reopen":
			title = fmt.Sprintf("Mở lại #%d · %s", w.ID, w.Title)
		case w.Kind == "proposed":
			title = fmt.Sprintf("Mate đề xuất #%d · %s", w.ID, w.Title)
		default:
			title = fmt.Sprintf("Nhận việc #%d · %s", w.ID, w.Title)
		}
		body := "Mở Mate, nhấn Ctrl+L để xem bảng Work."
		if line := strings.TrimSpace(move.Text); line != "" {
			body = line + "\n" + body
		}
		out = append(out, notifications.Work{Key: fmt.Sprintf("%s/work/%d", p.ProjectID, w.ID), At: move.At, Title: title, Body: body, URL: url})
	}
	return out
}

func notificationSummary(ctx context.Context, ws *store.Workspace, p query.ProjectNode, now time.Time) string {
	counts := map[query.CrewStatus]int{}
	for _, c := range p.Crews {
		counts[c.Status]++
	}
	state := "chưa có Mate"
	if p.Mate.Designated.IsKnown() {
		state = "Mate ghi nhận " + p.Mate.Designated.Value.Status.String()
	}
	if p.Mate.Designated.State == query.Unknown {
		state = "không đọc được trạng thái Mate"
	}
	activity := time.Time{}
	if p.Box.IsKnown() {
		activity = p.Box.Value.LastAt
	}
	if ws.MateHeadless(p.ProjectID) {
		client := matehost.Dial(ws.MateDir(p.ProjectID))
		client.Timeout = time.Second
		if s, err := client.State(ctx); err == nil {
			state = "Mate đang chờ lượt mới"
			if s.Busy {
				state = "Mate đang xử lý từ " + s.TurnStartedAt.Local().Format("15:04")
			}
			if s.Queued > 0 {
				state += fmt.Sprintf(", %d tin đang đợi", s.Queued)
			}
		} else {
			state = "không liên lạc được Mate; cần kiểm tra"
		}
		items := matehost.History(matehost.HistoryPath(ws.MateDir(p.ProjectID)), 1)
		if len(items) > 0 {
			it := items[0]
			at := it.At.Add(time.Duration(it.MS) * time.Millisecond)
			if at.After(activity) {
				activity = at
			}
			if it.Done && it.Error != "" {
				state += "; lượt gần nhất lỗi: " + it.Error
			}
		}
	}
	line := fmt.Sprintf("%s; %d crew làm, %d chờ quyết định, %d bàn giao, %d bị chặn", state,
		counts[query.CrewWorking]+counts[query.CrewSpawned], counts[query.CrewNeedsDecision], counts[query.CrewWaitMate], counts[query.CrewBlocked])
	if !activity.IsZero() {
		age := now.Sub(activity).Truncate(time.Minute)
		if age < 0 {
			age = 0
		}
		line += "; cập nhật ghi nhận " + age.String() + " trước"
	} else {
		line += "; chưa ghi nhận hoạt động"
	}
	return line
}
