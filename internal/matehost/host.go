package matehost

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness/claude"
)

// Sources of a turn.
const (
	SourceCaptain = "captain"
	SourceApp     = "app"
)

// Item is one message and what came of it: what the REPL shows and the
// history keeps.
type Item struct {
	At     time.Time `json:"at"`
	Source string    `json:"source"`
	// Label is an app line's tag (Label); empty for the captain.
	Label string `json:"label,omitempty"`
	Text  string `json:"text"`
	Reply string `json:"reply,omitempty"`
	// Silent is a reply that had nothing for the captain (Silent).
	Silent      bool   `json:"silent,omitempty"`
	Interrupted bool   `json:"interrupted,omitempty"`
	Error       string `json:"error,omitempty"`
	Done        bool   `json:"done,omitempty"`
}

// Visible reports whether the REPL shows the item's reply.
func (it Item) Visible() bool { return !it.Silent }

// State is what the host tells the rest of mate.
type State struct {
	Busy bool `json:"busy"`
	// Queued counts captain messages waiting for the current turn.
	Queued int `json:"queued"`
	// Draft is what the captain has typed and not sent.
	Draft         string    `json:"draft,omitempty"`
	TurnSource    string    `json:"turn_source,omitempty"`
	TurnStartedAt time.Time `json:"turn_started_at,omitzero"`
	Stopping      bool      `json:"stopping,omitempty"`
}

// Event is what the host tells its REPL.
type Event struct {
	// Item is set when a turn starts (Done false) or ends (Done true).
	Item *Item
	// Tool is set while a turn runs: the tool the Mate just called.
	Tool string
	// Queued is set when a captain message was queued behind a turn.
	Queued string
	// State is the host's state after the event.
	State State
}

// Reporter tells Herdr what the host is doing (herdrReporter).
type Reporter interface {
	Report(state string)
}

// Host serialises the Mate's messages: one turn at a time, captain
// messages queued behind it, app lines refused while it runs.
type Host struct {
	runner   *Runner
	history  string
	reporter Reporter
	now      func() time.Time
	notify   func(Event)

	mu       sync.Mutex
	run      *Run
	cur      *Item
	started  time.Time
	queue    []string
	draft    string
	stopping bool
	idle     chan struct{}
}

// Options configure a Host.
type Options struct {
	// History is the file items are appended to; empty keeps none.
	History  string
	Reporter Reporter
	Now      func() time.Time
}

// New makes a host for runner.
func New(runner *Runner, opts Options) *Host {
	h := &Host{runner: runner, history: opts.History, reporter: opts.Reporter, now: opts.Now, idle: make(chan struct{})}
	if h.now == nil {
		h.now = time.Now
	}
	close(h.idle)
	return h
}

// HistoryPath is the history file of the Mate in mateDir.
func HistoryPath(mateDir string) string { return filepath.Join(mateDir, HistoryName) }

// Subscribe sets the one listener for events; the REPL.
func (h *Host) Subscribe(fn func(Event)) {
	h.mu.Lock()
	h.notify = fn
	h.mu.Unlock()
}

// Errors a delivery or a message can be refused with.
var (
	ErrStopping = errors.New("the Mate is stopping")
	ErrQueued   = errors.New("the captain has messages waiting")
)

// Submit is a captain message: it runs now, or after the turn in flight.
func (h *Host) Submit(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	h.mu.Lock()
	if h.stopping {
		h.mu.Unlock()
		return
	}
	if h.run != nil || len(h.queue) > 0 {
		h.queue = append(h.queue, text)
		ev := Event{Queued: text, State: h.stateLocked()}
		h.mu.Unlock()
		h.emit(ev)
		return
	}
	h.mu.Unlock()
	if err := h.start(SourceCaptain, text); errors.Is(err, ErrBusy) {
		// Another process holds the turn lock: wait in the queue and retry.
		h.mu.Lock()
		h.queue = append([]string{text}, h.queue...)
		ev := Event{Queued: text, State: h.stateLocked()}
		h.mu.Unlock()
		h.emit(ev)
		go h.retryQueue()
	} else if err != nil {
		h.finishWithError(SourceCaptain, text, err)
	}
}

// Deliver is an app line. It is refused rather than queued while anything
// runs or waits, so the outbox keeps the line and its own retry and wedged
// rules (docs/mvp.md M21). It returns once the turn has started.
func (h *Host) Deliver(text string) error {
	h.mu.Lock()
	switch {
	case h.stopping:
		h.mu.Unlock()
		return ErrStopping
	case h.run != nil:
		h.mu.Unlock()
		return ErrBusy
	case len(h.queue) > 0:
		h.mu.Unlock()
		return ErrQueued
	}
	h.mu.Unlock()
	return h.start(SourceApp, text)
}

// start runs one turn. Only one start can win: run is checked again under
// the lock it is set under.
func (h *Host) start(source, text string) error {
	h.mu.Lock()
	if h.run != nil {
		h.mu.Unlock()
		return ErrBusy
	}
	run, err := h.runner.Start(text)
	if err != nil {
		h.mu.Unlock()
		return err
	}
	item := &Item{At: h.now(), Source: source, Text: text}
	if source == SourceApp {
		item.Label = Label(text)
	}
	h.run, h.cur, h.started = run, item, item.At
	h.idle = make(chan struct{})
	snapshot := *item
	ev := Event{Item: &snapshot, State: h.stateLocked()}
	h.mu.Unlock()
	h.report("working")
	h.emit(ev)
	go h.wait(run, item)
	return nil
}

func (h *Host) wait(run *Run, item *Item) {
	reply := run.Wait(func(ev claude.StreamEvent) {
		for _, t := range ev.Tools {
			h.mu.Lock()
			state := h.stateLocked()
			h.mu.Unlock()
			h.emit(Event{Tool: t, State: state})
		}
	})
	item.Reply = strings.TrimSpace(reply.Text)
	item.Interrupted = reply.Interrupted
	item.Silent = item.Source == SourceApp && Silent(item.Reply) && reply.Err == nil
	if reply.Err != nil {
		item.Error = reply.Err.Error()
	} else if reply.IsError && item.Reply != "" {
		item.Error = item.Reply
	}
	item.Done = true
	h.appendHistory(*item)

	h.mu.Lock()
	h.run, h.cur = nil, nil
	var next string
	if len(h.queue) > 0 && !h.stopping {
		next, h.queue = h.queue[0], h.queue[1:]
	}
	close(h.idle)
	done := *item
	ev := Event{Item: &done, State: h.stateLocked()}
	h.mu.Unlock()
	if next == "" {
		h.report("idle")
	}
	h.emit(ev)
	if next != "" {
		if err := h.start(SourceCaptain, next); errors.Is(err, ErrBusy) {
			h.mu.Lock()
			h.queue = append([]string{next}, h.queue...)
			h.mu.Unlock()
			go h.retryQueue()
		} else if err != nil {
			h.finishWithError(SourceCaptain, next, err)
		}
	}
}

// retryQueue starts the queue's head once the turn lock is free again; it
// is held by another process (a second host, or a test) only briefly.
func (h *Host) retryQueue() {
	for {
		time.Sleep(time.Second)
		h.mu.Lock()
		if h.stopping || h.run != nil || len(h.queue) == 0 {
			h.mu.Unlock()
			return
		}
		next := h.queue[0]
		h.queue = h.queue[1:]
		h.mu.Unlock()
		switch err := h.start(SourceCaptain, next); {
		case err == nil:
			return
		case errors.Is(err, ErrBusy):
			h.mu.Lock()
			h.queue = append([]string{next}, h.queue...)
			h.mu.Unlock()
		default:
			h.finishWithError(SourceCaptain, next, err)
			return
		}
	}
}

func (h *Host) finishWithError(source, text string, err error) {
	item := Item{At: h.now(), Source: source, Text: text, Error: err.Error(), Done: true}
	h.appendHistory(item)
	h.mu.Lock()
	ev := Event{Item: &item, State: h.stateLocked()}
	h.mu.Unlock()
	h.emit(ev)
}

// Interrupt stops the turn in flight, if any.
func (h *Host) Interrupt() bool {
	h.mu.Lock()
	run := h.run
	h.mu.Unlock()
	if run == nil {
		return false
	}
	run.Interrupt()
	return true
}

// SetDraft records what the captain has typed and not sent.
func (h *Host) SetDraft(text string) {
	h.mu.Lock()
	h.draft = text
	h.mu.Unlock()
}

// ClearDraft drops the captain's draft; the REPL picks it up as an event.
func (h *Host) ClearDraft() {
	h.mu.Lock()
	h.draft = ""
	ev := Event{State: h.stateLocked()}
	h.mu.Unlock()
	h.emit(ev)
}

// State is the host's state now.
func (h *Host) State() State {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.stateLocked()
}

func (h *Host) stateLocked() State {
	s := State{Busy: h.run != nil, Queued: len(h.queue), Draft: h.draft, Stopping: h.stopping}
	if h.cur != nil {
		s.TurnSource, s.TurnStartedAt = h.cur.Source, h.started
	}
	return s
}

// Stop refuses every new message and waits up to wait for the turn in
// flight; force, or a turn still running when wait is up, kills it. Queued
// captain messages are dropped. It reports whether a turn had to be killed.
func (h *Host) Stop(wait time.Duration, force bool) (killed bool) {
	h.mu.Lock()
	h.stopping = true
	h.queue = nil
	run, idle := h.run, h.idle
	ev := Event{State: h.stateLocked()}
	h.mu.Unlock()
	h.emit(ev)
	if run == nil {
		return false
	}
	if !force {
		select {
		case <-idle:
			return false
		case <-time.After(wait):
		}
	}
	run.Kill()
	<-idle
	return true
}

func (h *Host) emit(ev Event) {
	h.mu.Lock()
	fn := h.notify
	h.mu.Unlock()
	if fn != nil {
		fn(ev)
	}
}

func (h *Host) report(state string) {
	if h.reporter != nil {
		h.reporter.Report(state)
	}
}

func (h *Host) appendHistory(it Item) {
	if h.history == "" {
		return
	}
	data, err := json.Marshal(it)
	if err != nil {
		return
	}
	f, err := os.OpenFile(h.history, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(data, '\n'))
}

// History returns the last n items of a history file, oldest first. A
// missing file has none; a line that does not parse is skipped.
func History(path string, n int) []Item {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	var items []Item
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 8*1024*1024)
	for sc.Scan() {
		var it Item
		if json.Unmarshal(sc.Bytes(), &it) == nil {
			items = append(items, it)
		}
	}
	if n > 0 && len(items) > n {
		items = items[len(items)-n:]
	}
	return items
}
