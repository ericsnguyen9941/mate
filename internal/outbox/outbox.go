// Package outbox is the one sender into a Mate's composer (docs/mvp.md task
// 30): every line mate types into a Mate on its own initiative - an
// `[assign]` the captain pressed, a digest the auto daemon built - is queued
// in `mate/.outbox` (internal/store/outbox.go) and delivered from here.
//
// # Why a queue
//
// The Mate's composer is shared with the captain and is busy most of the
// time. Task 24 measured `[assign]` on a supervising Mate refused four or
// five times in a row with `agent is mid-turn`, because section 9's poll
// loop keeps the Mate inside a tool call for most of every cycle; the
// captain was left pressing the button again until it went through. The
// auto daemon had the same problem and its own private retry and `wedged`
// stopwatch for it. Both now enqueue here and return, and this package
// retries every DefaultInterval until the composer is empty.
//
// # What it never does
//
// It never hands a line to the harness's own queue (`herdr agent prompt`,
// send.Options.QueueWhileBusy): a line queued behind a running turn is one
// nobody can prove was read, which is the failure mvp.md section 7 records.
// It never types into a composer that is not empty: send.Send classifies
// first and refuses Busy, Pending and Unknown. And it never types one item
// twice: an item is attempted only under the outbox's writers' lock, marked
// sent in the same critical section that delivered it, and - for the one
// window a crash can leave open, between sent.log and the outbox rewrite -
// looked for in sent.log before it is typed again.
//
// # Wedged
//
// A queued item older than WedgedAfter at a failed attempt opens a `wedged`
// incident whose crew is `mate` (mvp.md section 4b); the next verified send
// resolves it. The clock is the oldest queued item's own `at`, on disk, so a
// console restart neither resets it nor opens a second incident: the file is
// the state, and the open line is written only when none is open. This is
// the one implementation of that rule; internal/autopilot no longer has its
// own.
package outbox

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/loopstat"
	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/observability"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/screen/fixture"
	"github.com/nguyenngocanh94/mate/internal/send"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// MateCrew is the crew field of the `wedged` incident this package opens.
// The Mate is not a crew, but `incidents.log` has one identity per line - a
// (crew, kind) pair - and a wedged composer is a finding about the Mate's
// pane, so this is the name it is filed under.
const MateCrew = "mate"

// WedgedKind is the incident kind (box.IncidentWedged's spelling; this
// package does not import internal/box).
const WedgedKind = "wedged"

// Defaults of mvp.md task 30.
const (
	// DefaultInterval is how often the loop retries the oldest queued item
	// of each project.
	DefaultInterval = 2 * time.Second
	// DefaultWedgedAfter is how long an item may wait before `wedged`.
	DefaultWedgedAfter = 5 * time.Minute
	// DefaultIdleInterval is the loop's pause while no project has an item
	// queued (CON-270): it then runs on a wake (Wake, an Enqueue or Offer
	// in this process), and this is only the safety net for a wake lost -
	// a write by a mate process with no core to tell, a notice dropped on
	// the way.
	DefaultIdleInterval = time.Minute
)

// DefaultDeliverTimeout bounds one delivery. A delivery holds the outbox's
// writers' lock from the handle lookup to the verified enter, so a Herdr
// that never answers would otherwise hold every sender and enqueuer of the
// project with it. A healthy delivery is well inside it: Jev's own deadline
// (notice.Timeout, 8s) plus the reads, the settle and the enter retries of
// send.Send.
const DefaultDeliverTimeout = 20 * time.Second

// Runtime is send.Runtime: read a pane, type, press keys, wait. It
// structurally cannot start, stop or prompt anything.
type Runtime interface {
	send.Runtime
}

var _ Runtime = runtime.Adapter(nil)

// HandleFunc resolves the Herdr handle and harness of one project's Mate.
// spawn.MateOutbox wires spawn.MateHandle into it, so this package imports
// nothing that can start or stop an agent. An error means the line has
// nowhere to go yet: the Mate is not running, or Herdr is not answering.
type HandleFunc func(ctx context.Context, project string) (runtime.AgentHandle, harness.Kind, error)

// Clock is where the sender reads the time; tests pass a fake one.
type Clock interface {
	Now() time.Time
}

// Sleeper is the pause between passes and inside a send.
type Sleeper interface {
	Sleep(ctx context.Context, d time.Duration) error
}

// Deps are the sender's collaborators. Runtime, Handle and Harnesses are
// needed to deliver; Enqueue and Offer need none of them.
type Deps struct {
	Runtime Runtime
	Handle  HandleFunc
	// Harnesses are the harnesses the Mate may run; a delivery to one that
	// is not registered is refused (send.Deps.Harnesses).
	Harnesses   harness.Registry
	Clock       Clock
	Sleeper     Sleeper
	Interval    time.Duration
	WedgedAfter time.Duration
	// Observer reads the Mate's pane before a delivery is typed
	// (send.Deps.Observer). Nil means the fixture observer.
	Observer screen.Observer
	// DeliverTimeout bounds one delivery; zero means DefaultDeliverTimeout.
	DeliverTimeout time.Duration
	// ConsolePresent says whether a console is open on the workspace. The
	// loop (Drain) types into a Manual project's Mate only while it says
	// true; an Auto project's Mate gets its lines either way (CON-203,
	// decided 2026-11-05). Attempt, called by itself, is not gated. Nil is
	// a console always open.
	ConsolePresent func() bool
	// Pass records each pass the loop runs (Drain), for `mate core stats`
	// (CON-263). Nil records nothing.
	Pass *loopstat.Loop
	// IdleInterval is the loop's pause while nothing is queued; zero means
	// DefaultIdleInterval. Interval is its pause while something is.
	IdleInterval time.Duration
	// SharedConfig says the caller keeps the workspace handle's
	// workspace.yaml current itself - internal/core's engine, which reloads
	// it on a `workspace` notice and when the file changes (CON-270) - so a
	// pass does not read the file again. False re-reads it every pass.
	SharedConfig bool
}

func (d Deps) consolePresent() bool {
	return d.ConsolePresent == nil || d.ConsolePresent()
}

func (d Deps) now() time.Time {
	if d.Clock != nil {
		return d.Clock.Now()
	}
	return time.Now()
}

func (d Deps) sleep(ctx context.Context, dur time.Duration) error {
	if d.Sleeper != nil {
		return d.Sleeper.Sleep(ctx, dur)
	}
	timer := time.NewTimer(dur)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (d Deps) deliverTimeout() time.Duration {
	if d.DeliverTimeout > 0 {
		return d.DeliverTimeout
	}
	return DefaultDeliverTimeout
}

func (d Deps) observer() screen.Observer {
	if d.Observer != nil {
		return d.Observer
	}
	return fixture.New()
}

func (d Deps) interval() time.Duration {
	if d.Interval > 0 {
		return d.Interval
	}
	return DefaultInterval
}

func (d Deps) idleInterval() time.Duration {
	if d.IdleInterval > 0 {
		return d.IdleInterval
	}
	return DefaultIdleInterval
}

func (d Deps) wedgedAfter() time.Duration {
	if d.WedgedAfter > 0 {
		return d.WedgedAfter
	}
	return DefaultWedgedAfter
}

// Sender delivers a workspace's outboxes. All of its state is on disk, so
// any number of Senders over one workspace agree: the console's loop, the
// immediate attempt an `[assign]` makes, and the auto daemon's own attempt
// right after it enqueues.
//
// Drain re-reads workspace.yaml into the handle, unless Deps.SharedConfig
// says its owner keeps it current. Since CON-204 a handle's configuration
// is swapped whole, so a loop may share the handle.
type Sender struct {
	ws   Workspace
	deps Deps
	// wake cuts the loop's pause short (Wake); it holds one pending wake.
	wake chan struct{}
	// pending is whether the loop's last pass found an item queued: a
	// Nudge wakes it only then.
	pending atomic.Bool

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

// New builds a sender. It sends nothing until Attempt, Drain or Start.
func New(ws Workspace, deps Deps) *Sender {
	return &Sender{ws: ws, deps: deps, wake: make(chan struct{}, 1)}
}

// Wake runs the loop's next pass now rather than at the end of its pause:
// something was queued (an `outbox` notice from another mate process).
// It never blocks; a loop not running yet finds the wake when it starts.
func (s *Sender) Wake() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

// Nudge is Wake while the loop's last pass left something queued, and
// nothing otherwise: the Mate's composer may have emptied (a `sent` or
// `mate` notice - a turn ended, the captain's line was recorded), which
// matters only to a line waiting for it.
func (s *Sender) Nudge() {
	if s.pending.Load() {
		s.Wake()
	}
}

// Request is one line to queue.
type Request struct {
	// Source is store.OutboxSourceAssign, OutboxSourceDigest,
	// OutboxSourceStow or OutboxSourcePR.
	Source string
	// Key is the dedup key; see store.OutboxItem.Key.
	Key string
	// Text is the line without the marker.
	Text string
	// Cursor is the digest's auto cursor to record once it is sent.
	Cursor map[string]int64
}

func (r Request) validate() error {
	if r.Source != store.OutboxSourceAssign && r.Source != store.OutboxSourceDigest && r.Source != store.OutboxSourceStow && r.Source != store.OutboxSourcePR {
		return observability.NewError(observability.CodeUsage, fmt.Sprintf("outbox: unknown source %q", r.Source))
	}
	if strings.TrimSpace(r.Key) == "" {
		return observability.NewError(observability.CodeUsage, "outbox: an item needs a dedup key")
	}
	if strings.TrimSpace(r.Text) == "" {
		return observability.NewError(observability.CodeUsage, "outbox: an item needs a non-empty line")
	}
	if strings.ContainsAny(r.Text, "\r\n") {
		return observability.NewError(observability.CodeUsage,
			"outbox: one line per item; long content goes in a file the Mate is pointed at")
	}
	return nil
}

// Enqueued is what Enqueue did.
type Enqueued struct {
	// Item is the item now in the outbox for the request's key: the new
	// one, or the earlier one that made this a duplicate.
	Item store.OutboxItem
	// Duplicate is true when an item with the same key was already queued
	// or sent, so nothing was added.
	Duplicate bool
}

// Enqueue appends one item unless an item with the same key is already
// queued or sent. A dropped item does not count: it was never delivered.
func (s *Sender) Enqueue(project string, req Request) (Enqueued, error) {
	if err := req.validate(); err != nil {
		return Enqueued{}, err
	}
	var out Enqueued
	now := s.deps.now()
	err := s.ws.UpdateOutbox(project, now, func(items []store.OutboxItem) ([]store.OutboxItem, bool, error) {
		for _, item := range items {
			if item.Key == req.Key && item.State != store.OutboxDropped {
				out = Enqueued{Item: item, Duplicate: true}
				return items, false, nil
			}
		}
		item, err := s.newItem(project, items, req, now)
		if err != nil {
			return items, false, err
		}
		out.Item = item
		return append(items, item), true, nil
	})
	if err == nil && !out.Duplicate {
		s.queued()
	}
	return out, err
}

// Offer is the auto daemon's enqueue: at most one queued item of a source at
// a time, refreshed in place. build runs under the writers' lock, so what it
// reads (the auto cursor) cannot move under it: the cursor moves only when
// an item is marked sent, and that happens under the same lock.
//
// When build reports nothing to say, a queued item of that source is
// dropped - what it reported has been answered meanwhile. When it does have
// something, a queued item of that source with the same key is left alone,
// one with a different key is rewritten with the new line (keeping its id,
// its `at` - which is the wedged clock - and its attempts), and with none
// queued a new item is appended.
func (s *Sender) Offer(project, source string,
	build func() (req Request, ok bool, err error)) (store.OutboxItem, bool, error) {
	var (
		out  store.OutboxItem
		have bool
	)
	now := s.deps.now()
	err := s.ws.UpdateOutbox(project, now, func(items []store.OutboxItem) ([]store.OutboxItem, bool, error) {
		req, ok, err := build()
		if err != nil {
			return items, false, err
		}
		queued := -1
		for i, item := range items {
			if item.Source == source && item.Queued() {
				queued = i
				break
			}
		}
		if !ok {
			if queued < 0 {
				return items, false, nil
			}
			items[queued].State = store.OutboxDropped
			items[queued].LastRefusal = "nothing left to report"
			return items, true, nil
		}
		req.Source = source
		if err := req.validate(); err != nil {
			return items, false, err
		}
		have = true
		if queued >= 0 {
			item := &items[queued]
			if item.Key == req.Key {
				out = *item
				return items, false, nil
			}
			item.Key, item.Text, item.Cursor = req.Key, req.Text, req.Cursor
			out = *item
			return items, true, nil
		}
		item, err := s.newItem(project, items, req, now)
		if err != nil {
			return items, false, err
		}
		out = item
		return append(items, item), true, nil
	})
	if err == nil && have {
		s.queued()
	}
	return out, have, err
}

func (s *Sender) newItem(project string, items []store.OutboxItem, req Request, now time.Time) (store.OutboxItem, error) {
	from, err := s.ws.SentLogSize(project)
	if err != nil {
		return store.OutboxItem{}, err
	}
	return store.OutboxItem{
		ID:          store.NextOutboxID(items),
		At:          now,
		Source:      req.Source,
		Key:         req.Key,
		Text:        req.Text,
		State:       store.OutboxQueued,
		SentLogFrom: from,
		Cursor:      req.Cursor,
	}, nil
}

// Attempt is what one try at a project's oldest queued item did.
type Attempt struct {
	// Item is the item tried, as it stands after the try. It is the zero
	// value when nothing was queued.
	Item store.OutboxItem
	// Delivered is the item verified into the composer on this try (or
	// found already delivered in sent.log).
	Delivered bool
	// Dropped is a digest withdrawn because its project left auto mode.
	Dropped bool
	// Refusal is why the item is still queued.
	Refusal error
	// Agent is the Mate's agent name, when a send was made.
	Agent string
}

// Tried reports whether there was anything queued to try.
func (a Attempt) Tried() bool { return a.Item.ID != 0 }

// Attempt tries the oldest queued item of one project, once.
//
// The error is for what went wrong around a send - the outbox could not be
// read or written, sent.log could not be appended after a verified send -
// never for the send's own refusal, which is Attempt.Refusal: a Mate
// mid-turn is the normal case here, not a failure.
func (s *Sender) Attempt(ctx context.Context, project string) (Attempt, error) {
	var (
		out     Attempt
		postErr error
	)
	now := s.deps.now()
	err := s.ws.UpdateOutbox(project, now, func(items []store.OutboxItem) ([]store.OutboxItem, bool, error) {
		index := -1
		for i, item := range items {
			if item.Queued() {
				index = i
				break
			}
		}
		if index < 0 {
			return items, false, nil
		}
		item := &items[index]
		if item.Source != store.OutboxSourceStow {
			release, acquired, err := s.ws.TryMateMaintenance(project, false)
			if err != nil {
				return items, false, err
			}
			if !acquired {
				return items, false, nil
			}
			defer release()
		}

		// A line sent.log already holds was delivered by a sender that
		// died before it could say so. Typing it again is the one mistake
		// this package exists to rule out.
		if at, ok, err := s.recorded(project, *item); err != nil {
			return items, false, err
		} else if ok {
			postErr = s.settle(project, item, at, "")
			out = Attempt{Item: *item, Delivered: true}
			return items, true, nil
		}

		if s.withdrawn(project, *item) {
			item.State = store.OutboxDropped
			item.LastRefusal = "auto mode was turned off before the digest could be typed"
			out = Attempt{Item: *item, Dropped: true}
			return items, true, nil
		}

		if s.tuiOpen(project) {
			// The Mate's session is open in its TUI (CON-227): the line
			// waits for it to close, as for a busy Mate, but nothing is
			// tried, so no attempt counts and no wedged clock runs out
			// against a TUI the captain keeps open.
			out = Attempt{Item: *item, Refusal: matehost.ErrTUI}
			return items, false, nil
		}
		report, err := s.deliver(ctx, project, *item)
		if errors.Is(err, errWithdrawn) {
			item.State = store.OutboxDropped
			item.LastRefusal = "auto mode was turned off before the digest could be typed"
			out = Attempt{Item: *item, Dropped: true}
			return items, true, nil
		}
		item.Attempts++
		item.TriedAt = now
		if err != nil {
			item.LastRefusal = err.Error()
			out = Attempt{Item: *item, Refusal: err}
			if age := now.Sub(item.At); age >= s.deps.wedgedAfter() {
				postErr = s.openWedged(project, now, fmt.Sprintf(
					"nothing queued for the Mate has been delivered for %s: %v", age.Round(time.Second), err))
			}
			return items, true, nil
		}

		// sent.log first, then the digest cursor, then the item marked
		// sent: a crash after sent.log costs nothing (the next attempt
		// finds the line there), and the reverse order could lose a line
		// the Mate never got. The time is read again: a verified send
		// takes a second or more of settle and enter retries, and the
		// record is of when the composer cleared, not when the try began.
		delivered := s.deps.now()
		recordErr := s.ws.AppendSent(project, store.SentEntry{
			Time:   delivered,
			Source: store.SourceApp,
			Target: store.TargetMate,
			Text:   item.Text,
		})
		postErr = errors.Join(recordErr, s.settle(project, item, delivered, report.Agent))
		out = Attempt{Item: *item, Delivered: true, Agent: report.Agent}
		return items, true, nil
	})
	return out, errors.Join(err, postErr)
}

// settle marks an item delivered: the digest's cursor recorded, the item
// sent, and an open `wedged` resolved.
func (s *Sender) settle(project string, item *store.OutboxItem, at time.Time, agent string) error {
	var errs []error
	if item.Source == store.OutboxSourceDigest && item.Cursor != nil {
		errs = append(errs, s.ws.WriteAutoCursor(project, item.Cursor))
	}
	item.State = store.OutboxSent
	item.SentAt = at
	item.LastRefusal = ""
	who := "the Mate"
	if agent != "" {
		who = agent
	}
	errs = append(errs, s.resolveWedged(project, at, fmt.Sprintf("a queued line was verified into %s's composer", who)))
	return errors.Join(errs...)
}

// tuiOpen reports whether project's Mate has its session open in the
// harness's TUI (CON-227): `.mate/tui.lock` names it.
func (s *Sender) tuiOpen(project string) bool {
	h, held, err := s.ws.TUIRunning()
	return err == nil && held && h.Project == project
}

// errWithdrawn is the one refusal that is not one: the captain took the
// composer over (the Mate's hook deleted `.auto`) while the digest was on its
// way. Nothing was typed, nothing is wedged.
var errWithdrawn = errors.New("auto mode was turned off")

// withdrawn reports a digest whose project is no longer in auto mode. An
// assign is the captain's own keystroke and is never withdrawn this way.
func (s *Sender) withdrawn(project string, item store.OutboxItem) bool {
	return item.Source == store.OutboxSourceDigest && !s.ws.Auto(project)
}

// deliver resolves the Mate's pane and types the line into it, verified,
// within DeliverTimeout. Running out of it is a refusal like any other: the
// caller keeps the item queued, and the lock it holds is released.
func (s *Sender) deliver(ctx context.Context, project string, item store.OutboxItem) (send.Report, error) {
	if s.deps.Runtime == nil || s.deps.Handle == nil {
		return send.Report{}, errors.New("outbox: no runtime is wired to deliver with")
	}
	bound := s.deps.deliverTimeout()
	dctx, cancel := context.WithTimeout(ctx, bound)
	defer cancel()
	report, err := s.deliverWithin(dctx, project, item)
	if err != nil && errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil {
		return report, fmt.Errorf("the delivery did not finish within %s: %w", bound, err)
	}
	return report, err
}

func (s *Sender) deliverWithin(ctx context.Context, project string, item store.OutboxItem) (send.Report, error) {
	if s.ws.MateHeadless(project) {
		return s.deliverHeadless(ctx, project, item)
	}
	handle, kind, err := s.deps.Handle(ctx, project)
	if err != nil {
		return send.Report{}, err
	}
	// Read again as late as it can be read: the handle lookup touches Herdr,
	// and a captain who typed into the Mate meanwhile has already had
	// `.auto` deleted by the Mate's own hook. Answering them a moment later
	// is the takeover auto mode promises not to fight.
	if s.withdrawn(project, item) {
		return send.Report{}, errWithdrawn
	}
	deps := send.Deps{Runtime: s.deps.Runtime, Harnesses: s.deps.Harnesses, Observer: s.deps.Observer}
	if s.deps.Sleeper != nil {
		deps.Sleep = s.deps.sleep
	}
	return send.Send(ctx, deps, handle, kind, item.Text, send.Options{Marker: true})
}

// deliverHeadless hands the line to a headless Mate's host (CON-78). The
// host refuses it while a turn runs or the captain has messages waiting,
// which is the same refusal a busy composer gives: the item stays queued.
// The marker goes with it, so the Mate's prompt hook records it as the
// app's line exactly as it does for a typed one.
func (s *Sender) deliverHeadless(ctx context.Context, project string, item store.OutboxItem) (send.Report, error) {
	if s.withdrawn(project, item) {
		return send.Report{}, errWithdrawn
	}
	meta, err := s.ws.ReadMateMeta(project)
	if err != nil {
		return send.Report{}, err
	}
	text := send.Marker + item.Text
	if err := matehost.Dial(s.ws.MateDir(project)).Deliver(ctx, text); err != nil {
		return send.Report{}, err
	}
	return send.Report{Agent: meta["agent"], Text: text, Typed: true}, nil
}

// recorded looks for the item's own line in sent.log after it was queued.
func (s *Sender) recorded(project string, item store.OutboxItem) (time.Time, bool, error) {
	entries, _, err := s.ws.ReadSent(project, item.SentLogFrom)
	if err != nil {
		return time.Time{}, false, err
	}
	want := flatten(item.Text)
	for _, e := range entries {
		if e.Source == store.SourceApp && e.Target == store.TargetMate && e.Text == want {
			return e.Time, true, nil
		}
	}
	return time.Time{}, false, nil
}

// flatten is sent.log's own normalisation of a line (store.AppendSent).
func flatten(text string) string {
	return strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ", "\t", " ").Replace(text)
}

// wedgedOpen reports whether the last (mate, wedged) line is `open`.
func (s *Sender) wedgedOpen(project string) (bool, error) {
	entries, _, err := s.ws.ReadIncidents(project, 0)
	if err != nil {
		return false, err
	}
	open := false
	for _, e := range entries {
		if e.Crew == MateCrew && e.Kind == WedgedKind {
			open = e.Open()
		}
	}
	return open, nil
}

func (s *Sender) openWedged(project string, now time.Time, text string) error {
	open, err := s.wedgedOpen(project)
	if err != nil || open {
		return err
	}
	return s.ws.AppendIncident(project, store.IncidentEntry{
		Time: now, Crew: MateCrew, Kind: WedgedKind, State: store.IncidentOpen, Text: text,
	})
}

func (s *Sender) resolveWedged(project string, now time.Time, text string) error {
	open, err := s.wedgedOpen(project)
	if err != nil || !open {
		return err
	}
	return s.ws.AppendIncident(project, store.IncidentEntry{
		Time: now, Crew: MateCrew, Kind: WedgedKind, State: store.IncidentResolved, Text: text,
	})
}

// Drain is one pass of the loop: the oldest queued item of every project
// that has one. The error joins the per-project failures; the pass always
// finishes.
func (s *Sender) Drain(ctx context.Context) error {
	_, err := s.drain(ctx)
	return err
}

// drain is Drain, and whether it found an item queued in any project: the
// loop then tries again after Interval, and otherwise waits for a wake.
// A pass that failed before it could look counts as one that found some,
// so the loop asks again on the old cadence.
func (s *Sender) drain(ctx context.Context) (pending bool, _ error) {
	if err := ctx.Err(); err != nil {
		return true, err
	}
	// workspace.yaml is re-read every pass, so a project registered beside
	// this console is delivered to without a restart - unless the owner of
	// the handle keeps it current (SharedConfig).
	if !s.deps.SharedConfig {
		if err := s.ws.LoadConfig(); err != nil {
			return true, err
		}
	}
	var errs []error
	for _, ref := range s.ws.Projects() {
		if err := ctx.Err(); err != nil {
			return true, err
		}
		// A lock-free look first: most passes over most projects find
		// nothing queued, and taking the writers' lock to learn that would
		// contend with a send in flight for no reason.
		items, err := s.ws.ReadOutbox(ref.Name)
		if err != nil {
			errs = append(errs, fmt.Errorf("outbox: %s: %w", ref.Name, err))
			pending = true
			continue
		}
		head, ok := oldestQueued(items)
		if !ok {
			if now := s.deps.now(); store.OutboxCompactable(items, now) {
				if err := s.ws.UpdateOutbox(ref.Name, now, func(items []store.OutboxItem) ([]store.OutboxItem, bool, error) {
					return items, true, nil
				}); err != nil {
					errs = append(errs, fmt.Errorf("outbox: %s: %w", ref.Name, err))
				}
			}
			continue
		}
		// Queued: tried again after Interval whatever this try does. A
		// delivery that empties the queue costs one pass more, which finds
		// nothing and lets the loop rest.
		pending = true
		// With no console open a Manual project's Mate is left alone: its
		// captain took the composer, and nothing is typed into it until a
		// console is back or M21 turns auto mode on again. A digest at the
		// head is still tried, which withdraws it as it would with a
		// console (withdrawn) rather than keeping every line behind it.
		if !s.ws.Auto(ref.Name) && !s.deps.consolePresent() && head.Source != store.OutboxSourceDigest {
			continue
		}
		if _, err := s.Attempt(ctx, ref.Name); err != nil {
			errs = append(errs, fmt.Errorf("outbox: %s: %w", ref.Name, err))
		}
	}
	return pending, errors.Join(errs...)
}

// oldestQueued is the item Attempt would try, when one is queued.
func oldestQueued(items []store.OutboxItem) (store.OutboxItem, bool) {
	for _, item := range items {
		if item.Queued() {
			return item, true
		}
	}
	return store.OutboxItem{}, false
}

// Start runs the loop on its own goroutine until Stop or a cancelled
// context: a pass (Drain) at once - which also compacts a quiet outbox
// (store.OutboxCompactable) -, then one each Interval while some project
// has an item queued, and otherwise one on a wake - Wake, Nudge, an
// Enqueue or Offer by any Sender of this process over the same workspace -
// or after IdleInterval at the latest (CON-270). Calling it twice is a
// no-op.
func (s *Sender) Start(ctx context.Context) {
	s.mu.Lock()
	if s.cancel != nil {
		s.mu.Unlock()
		return
	}
	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	done := make(chan struct{})
	s.done = done
	s.mu.Unlock()
	unregister := s.register()
	go func() {
		defer unregister()
		s.run(runCtx, done)
	}()
}

// Stop ends the loop and waits for the pass in flight.
func (s *Sender) Stop() {
	s.mu.Lock()
	cancel, done := s.cancel, s.done
	s.cancel, s.done = nil, nil
	s.mu.Unlock()
	if cancel == nil {
		return
	}
	cancel()
	<-done
}

// run closes done, the channel Start made, not s.done: Stop clears that
// field, possibly before this goroutine has first run.
func (s *Sender) run(ctx context.Context, done chan struct{}) {
	defer close(done)
	for {
		// A failed pass is not fatal: Herdr may be restarting, a file may
		// be mid-write. Nothing queued was lost, and the next pass asks
		// again.
		start := s.deps.Pass.Begin()
		pending, _ := s.drain(ctx)
		s.pending.Store(pending)
		s.deps.Pass.End(start)
		pause := s.deps.idleInterval()
		if pending {
			pause = s.deps.interval()
		}
		if err := s.pause(ctx, pause); err != nil {
			return
		}
	}
}

// pause sleeps d through Deps' Sleeper, cut short by a wake. A wake is
// taken before pause returns, so the pass after it reads what the waker
// wrote before it woke the loop.
func (s *Sender) pause(ctx context.Context, d time.Duration) error {
	wctx, cancel := context.WithCancel(ctx)
	watched := make(chan struct{})
	go func() {
		defer close(watched)
		select {
		case <-s.wake:
			cancel()
		case <-wctx.Done():
		}
	}()
	_ = s.deps.sleep(wctx, d)
	cancel()
	<-watched
	return ctx.Err()
}

// The loops running in this process, by the workspace they deliver for:
// an Enqueue or Offer by any Sender over the same workspace wakes them
// (CON-270). The core's own enqueuers - the auto daemon's digest, a UI's
// `[assign]`, the PR sweep, a stow - each hold a Sender of their own, and
// the core sends itself no notice (internal/store's notifier is unset in
// it).
var (
	loopsMu sync.Mutex
	loops   = map[string]map[*Sender]bool{}
)

// root is the workspace a Sender's handle is over; "" for a handle that
// does not say, whose loop every enqueue of the process wakes.
func (s *Sender) root() string {
	if r, ok := s.ws.(interface{ Root() string }); ok {
		return r.Root()
	}
	return ""
}

func (s *Sender) register() (unregister func()) {
	root := s.root()
	loopsMu.Lock()
	if loops[root] == nil {
		loops[root] = map[*Sender]bool{}
	}
	loops[root][s] = true
	loopsMu.Unlock()
	return func() {
		loopsMu.Lock()
		delete(loops[root], s)
		if len(loops[root]) == 0 {
			delete(loops, root)
		}
		loopsMu.Unlock()
	}
}

// queued wakes every loop of this process over s's workspace: an item was
// just queued, or a queued digest rewritten.
func (s *Sender) queued() {
	root := s.root()
	loopsMu.Lock()
	var wake []*Sender
	for l := range loops[root] {
		wake = append(wake, l)
	}
	if root != "" {
		for l := range loops[""] {
			wake = append(wake, l)
		}
	}
	loopsMu.Unlock()
	for _, l := range wake {
		l.Wake()
	}
}
