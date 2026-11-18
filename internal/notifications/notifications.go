// Package notifications delivers captain news, work board moves, the
// Mate's replies to app lines and periodic recorded-state reports.
// It does not call a model or infer that an agent is stuck from silence.
package notifications

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/platform"
)

const ReportInterval = 10 * time.Minute

type Message struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url,omitempty"`
}

type Sender interface {
	Send(context.Context, Message) error
}

type News struct{ Key, Revision, Title string }

// Work is one move of a work board row (CON-158) the captain is told about:
// At is when the row moved, and a later move of the same Key alerts again.
type Work struct {
	Key   string
	At    time.Time
	Title string
	Body  string
	URL   string
}

// Reply is what the Mate told the captain in its reply to an app line,
// such as a crew's hand-back (CON-280): At is when the reply ended, and
// Text its first words.
type Reply struct {
	At   time.Time
	Text string
	URL  string
}

type Project struct {
	Name    string
	News    []News  // Only news waiting for the captain.
	Work    []Work  // Each open row's latest move.
	Replies []Reply // The Mate's latest replies, oldest first.
	Summary string
}
type Snapshot struct {
	Projects []Project
	// Complete is false if any news source could not be read. In that case
	// known news can still notify, but a "no news" report would be a lie.
	Complete bool
	URL      string
}

type state struct {
	Delivered map[string]string `json:"delivered"`
	ReportAt  time.Time         `json:"report_at"`
	// WorkFrom is when work board moves started alerting: the first tick
	// that saw a receipts file without it. Older moves never alert, so the
	// first run does not replay a board's history.
	WorkFrom time.Time `json:"work_from,omitzero"`
	// RepliesFrom is WorkFrom for the Mate's replies.
	RepliesFrom time.Time `json:"replies_from,omitzero"`
}

// NextReport is when the periodic report falls due by the receipts at
// path: ReportInterval after the last report, news or work move a Tick
// recorded. ok is false while there is no receipt yet - the first Tick
// starts the interval - or the receipts do not read. A Tick at that time
// sends the report unless news is owed or unreadable then, which start the
// interval again (CON-270: the core ticks on a change and on this
// deadline, not every five seconds).
func NextReport(path string) (at time.Time, ok bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return time.Time{}, false
	}
	var s state
	if err := json.Unmarshal(b, &s); err != nil || s.ReportAt.IsZero() {
		return time.Time{}, false
	}
	return s.ReportAt.Add(ReportInterval), true
}

// Tick serializes consoles sharing a workspace, and persists receipts only
// after the native notification API accepts a message. The caller retries
// failures; neither an unavailable notifier nor a restart loses a news.
func Tick(ctx context.Context, path string, snap Snapshot, sender Sender, now time.Time) error {
	lock, err := platform.OpenNoFollow(path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer lock.Close()
	locked, err := platform.TryLock(lock)
	if err != nil {
		return err
	}
	if !locked {
		return nil
	}
	defer platform.Unlock(lock)
	s := state{Delivered: map[string]string{}, ReportAt: now}
	b, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("notification receipts: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if s.Delivered == nil {
		s.Delivered = map[string]string{}
	}
	if s.WorkFrom.IsZero() {
		s.WorkFrom = now
	}
	if s.RepliesFrom.IsZero() {
		s.RepliesFrom = now
	}
	save := func() error {
		b, err := json.Marshal(s)
		if err != nil {
			return err
		}
		f, err := os.CreateTemp(filepath.Dir(path), ".notifications-*")
		if err != nil {
			return err
		}
		defer os.Remove(f.Name())
		if _, err := f.Write(b); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		return os.Rename(f.Name(), path)
	}
	owed := false
	for _, p := range snap.Projects {
		for _, n := range p.News {
			owed = true
			if s.Delivered[n.Key] == n.Revision {
				continue
			}
			digest := sha256.Sum256([]byte(path + "\x00" + n.Key))
			msg := Message{ID: fmt.Sprintf("news-%x", digest[:12]), Title: "Mate · " + p.Name + " · cần anh xem", Body: n.Title + "\nMở Mate, nhấn Ctrl+T để xem và trả lời news.", URL: snap.URL}
			if err := sender.Send(ctx, msg); err != nil {
				return err
			}
			s.Delivered[n.Key] = n.Revision
			s.ReportAt = now
			if err := save(); err != nil {
				return err
			}
		}
	}
	// A work move is a report, not news: it delays the next periodic report
	// but does not hold it off.
	for _, p := range snap.Projects {
		for _, w := range p.Work {
			rev := w.At.UTC().Format(time.RFC3339Nano)
			if !w.At.After(s.WorkFrom) || s.Delivered[w.Key] == rev {
				continue
			}
			digest := sha256.Sum256([]byte(path + "\x00" + w.Key))
			msg := Message{ID: fmt.Sprintf("work-%x", digest[:12]), Title: "Mate · " + p.Name + " · Work", Body: w.Title + "\n" + w.Body, URL: w.URL}
			if w.Body == "" {
				msg.Body = w.Title
			}
			if err := sender.Send(ctx, msg); err != nil {
				return err
			}
			s.Delivered[w.Key] = rev
			s.ReportAt = now
			if err := save(); err != nil {
				return err
			}
		}
	}
	// A reply is a report too. Only a project's newest one alerts, saying
	// how many came with it, so a core that was down does not replay them;
	// it replaces the project's last one in the notification centre.
	for _, p := range snap.Projects {
		key := p.Name + "/reply"
		from := s.RepliesFrom
		if t, err := time.Parse(time.RFC3339Nano, s.Delivered[key]); err == nil && t.After(from) {
			from = t
		}
		var fresh []Reply
		for _, r := range p.Replies {
			if r.At.After(from) {
				fresh = append(fresh, r)
			}
		}
		if len(fresh) == 0 {
			continue
		}
		last := fresh[len(fresh)-1]
		body := last.Text
		if len(fresh) > 1 {
			body += fmt.Sprintf("\n(+%d tin trước đó)", len(fresh)-1)
		}
		digest := sha256.Sum256([]byte(path + "\x00" + key))
		msg := Message{ID: fmt.Sprintf("reply-%x", digest[:12]), Title: "Mate · " + p.Name + " · Mate báo", Body: body + "\nMở Mate để đọc hết.", URL: last.URL}
		if err := sender.Send(ctx, msg); err != nil {
			return err
		}
		s.Delivered[key] = last.At.UTC().Format(time.RFC3339Nano)
		s.ReportAt = now
		if err := save(); err != nil {
			return err
		}
	}
	// Begin a fresh quiet interval after the last poll with outstanding news.
	// Unreadable news also suppress a report and reset that interval.
	if owed || !snap.Complete {
		s.ReportAt = now
		return save()
	}
	if now.Sub(s.ReportAt) >= ReportInterval && len(snap.Projects) > 0 {
		lines := make([]string, 0, len(snap.Projects))
		for _, p := range snap.Projects {
			lines = append(lines, p.Name+": "+p.Summary)
		}
		digest := sha256.Sum256([]byte(path))
		msg := Message{ID: fmt.Sprintf("report-%x", digest[:12]), Title: "Mate · báo cáo 10 phút", Body: "Không có news chờ anh.\n" + strings.Join(lines, "\n"), URL: snap.URL}
		if err := sender.Send(ctx, msg); err != nil {
			return err
		}
		s.ReportAt = now
	}
	return save()
}
