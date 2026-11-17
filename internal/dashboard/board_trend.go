package dashboard

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/timeline"
	"github.com/nguyenngocanh94/mate/internal/timeline/read"
)

// trendSlots are the slot widths /api/board/trend cuts time into, each with
// how many slots it shows when the page names no range.
var trendSlots = map[string]int{"1h": 48, "6h": 28, "1d": 14, "1w": 13}

// trendMaxBuckets caps one answer, so an hourly slot over a year is refused
// rather than drawn.
const trendMaxBuckets = 500

// slotStart is the start of the slot t falls in, in loc: the hour, the six
// hours from midnight, the day, or the week from Monday.
func slotStart(t time.Time, slot string, loc *time.Location) time.Time {
	t = t.In(loc)
	y, m, d := t.Date()
	switch slot {
	case "1h":
		return time.Date(y, m, d, t.Hour(), 0, 0, 0, loc)
	case "6h":
		return time.Date(y, m, d, t.Hour()-t.Hour()%6, 0, 0, 0, loc)
	case "1w":
		back := (int(t.Weekday()) + 6) % 7
		return time.Date(y, m, d-back, 0, 0, 0, 0, loc)
	}
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

// slotNext is the slot after the one starting at t. Days and weeks step by
// the calendar, so a daylight-saving day is still one slot.
func slotNext(t time.Time, slot string) time.Time {
	y, m, d := t.Date()
	switch slot {
	case "1h":
		return time.Date(y, m, d, t.Hour()+1, 0, 0, 0, t.Location())
	case "6h":
		return time.Date(y, m, d, t.Hour()+6, 0, 0, 0, t.Location())
	case "1w":
		return time.Date(y, m, d+7, 0, 0, 0, 0, t.Location())
	}
	return time.Date(y, m, d+1, 0, 0, 0, 0, t.Location())
}

// handleBoardTrend is GET /api/board/trend?slot=&from=&to=&count=&project=: the
// board's tokens and finished work cut into slots of an hour, six hours,
// a day or a week, so the Office can chart the work at any grain and zoom
// into one slot. Each slot also carries the tasks started in it (Started),
// so work is set against what it cost from when it began. from and to are RFC 3339 instants; each
// is moved to the start of its slot, to's slot is the last one, and
// without from the window is count slots (or the slot's default) ending
// at to, which defaults to now. Slots
// are the server clock's zone, as the board's days are. project narrows
// it to one project.
func (s *Server) handleBoardTrend(ctx context.Context, r *http.Request, env envelope) (any, error) {
	q := r.URL.Query()
	slot := strings.TrimSpace(q.Get("slot"))
	if slot == "" {
		slot = "1d"
	}
	count, ok := trendSlots[slot]
	if !ok {
		return nil, errBadRequest{reason: "slot must be 1h, 6h, 1d or 1w"}
	}
	if raw := strings.TrimSpace(q.Get("count")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > trendMaxBuckets {
			return nil, errBadRequest{reason: "count must be 1 to 500 slots"}
		}
		count = n
	}
	now := s.deps.now()
	loc := now.Location()
	instant := func(name string) (time.Time, error) {
		raw := strings.TrimSpace(q.Get(name))
		if raw == "" {
			return time.Time{}, nil
		}
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return time.Time{}, errBadRequest{reason: name + " must be an RFC 3339 time"}
		}
		return t, nil
	}
	from, err := instant("from")
	if err != nil {
		return nil, err
	}
	to, err := instant("to")
	if err != nil {
		return nil, err
	}
	if to.IsZero() {
		to = now
	}
	end := slotNext(slotStart(to, slot, loc), slot)
	var starts []time.Time
	if from.IsZero() {
		start := slotStart(to, slot, loc)
		for i := 0; i < count; i++ {
			starts = append([]time.Time{start}, starts...)
			start = slotStart(start.Add(-time.Second), slot, loc)
		}
	} else {
		if !from.Before(end) {
			return nil, errBadRequest{reason: "from must be before to"}
		}
		for t := slotStart(from, slot, loc); t.Before(end); t = slotNext(t, slot) {
			if len(starts) == trendMaxBuckets {
				return nil, errBadRequest{reason: "too many slots: widen the slot or narrow the range"}
			}
			starts = append(starts, t)
		}
	}
	project := strings.TrimSpace(q.Get("project"))

	out := BoardTrendResponse{envelope: env, Slot: slot, Zone: now.Format("-07:00"), Project: project,
		From: starts[0].Format(time.RFC3339), To: end.Format(time.RFC3339), Buckets: make([]BoardTrendBucket, len(starts))}
	for i, t := range starts {
		next := end
		if i+1 < len(starts) {
			next = starts[i+1]
		}
		out.Buckets[i] = BoardTrendBucket{From: t.Format(time.RFC3339), To: next.Format(time.RFC3339)}
	}
	// bucketOf is the slot a stored timestamp falls in, or -1 outside.
	bucketOf := func(ts string) int {
		if ts == "" {
			return -1
		}
		t := db.ParseTime(ts)
		if t.IsZero() || t.Before(starts[0]) || !t.Before(end) {
			return -1
		}
		// The last start not after t; the slots are few enough to walk.
		i := len(starts) - 1
		for i > 0 && t.Before(starts[i]) {
			i--
		}
		return i
	}

	turns, err := read.TurnTokensBetween(ctx, s.db.SQL(), db.FormatTime(starts[0]), db.FormatTime(end))
	if err != nil {
		return nil, err
	}
	for _, t := range turns {
		if project != "" && t.Project != project {
			continue
		}
		i := bucketOf(t.At)
		if i < 0 {
			continue
		}
		total := t.Input + t.CacheRead + t.CacheWrite + t.Output
		b := &out.Buckets[i]
		b.Tokens += total
		switch t.Kind {
		case timeline.ActorMate:
			b.MateTokens += total
		case timeline.ActorCrew:
			b.CrewTokens += total
		}
	}

	for _, ref := range s.ws.Projects() {
		if project != "" && ref.Name != project {
			continue
		}
		tasks, err := s.tasks(ctx, ref.Name, now)
		if err != nil {
			if out.Errors == nil {
				out.Errors = map[string]string{}
			}
			out.Errors[ref.Name] = err.Error()
		}
		for _, t := range tasks {
			if t.Review != "" {
				continue
			}
			finished := t.Closed && closedWord(t) == "finished"
			if i := bucketOf(t.ClosedAt); i >= 0 && finished {
				out.Buckets[i].Done++
				out.Buckets[i].DoneTokens += t.TokensWithReviews
			}
			if i := bucketOf(t.SpawnedAt); i >= 0 {
				c := &out.Buckets[i].Started
				c.Tasks++
				c.Tokens += t.TokensWithReviews
				c.WorkMs += t.AgeMs
				switch {
				case finished:
					c.Done++
				case !t.Closed:
					c.Open++
				}
			}
		}
	}
	return out, nil
}
