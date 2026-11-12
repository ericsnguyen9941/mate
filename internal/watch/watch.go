package watch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/loopstat"
	"github.com/nguyenngocanh94/mate/internal/observability"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/send"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Defaults of docs/mvp.md task 18. The poll is short because the console
// draws the health column from it; the stale threshold is long because the
// cost of the two mistakes is not symmetric - a crew wrongly called stale
// puts a false item in the human's inbox, while a crew called stale a minute
// late costs a minute.
const (
	// DefaultPollInterval is the pause between rounds.
	DefaultPollInterval = 5 * time.Second
	// WakeGap is the least time between the starts of two rounds when a
	// Wake cut the pause short (CON-260): a burst of writes costs one early
	// round, not one round each. It is half the poll interval, not less,
	// while every round still ingests the whole workspace (0.68 s on the
	// captain's 1308 crews, 2026-11-11): steady writes then cost at most
	// twice the rounds they did. The UI does not wait on it - a notice drops
	// the snapshot at once - and CON-268, which makes a round read only what
	// changed, is what may bring it down.
	WakeGap = 2500 * time.Millisecond
	// DefaultStaleAfter is how long the pane and the status file must both
	// be unchanged before `stale` opens.
	DefaultStaleAfter = 3 * time.Minute
	// DefaultHerdrTimeout bounds each Herdr question the observer asks: a
	// crew's handle, its agent, its pane, and the session. It is the
	// console's own bound on a Herdr question (cmd/mate livenessBound).
	DefaultHerdrTimeout = 3 * time.Second
	// DefaultNoTurnAfter is how long a working crew may go after launch
	// with no turn in the token ledger before `no_turn` opens (CON-74).
	// Longer than a sign-in or a trust dialog takes to clear, shorter than
	// a Mate should wait on a reviewer that never started.
	DefaultNoTurnAfter = 5 * time.Minute
	// DefaultSweepInterval is how often a round examines every project,
	// whatever the change feed said (CON-267). Between sweeps a round
	// examines only the projects that can change unheard - an open crew,
	// a live Mate - and those a notice named (Notices). The sweep is the
	// safety net under what neither covers: a git ref moved by hand, a
	// pricing.yaml or project.yaml edited in an editor, a notice lost to
	// a core restarting. A minute is the longest such a change waits to be
	// read; a sweep costs what every round cost before (~130 ms on the
	// captain's workspace, 2026-11-12), so once a minute is ~0.2% of a CPU.
	DefaultSweepInterval = time.Minute
)

// Runtime is the slice of runtime.Adapter the observer needs. It is narrow
// on purpose: this package can ask Herdr what it has and read a pane, and it
// structurally cannot start, stop, type into or prompt anything.
type Runtime interface {
	InspectAgent(ctx context.Context, handle runtime.AgentHandle) (runtime.ObservedAgent, error)
	// ReadAgentStyled, not ReadAgent: the health column reports the
	// composer state, and telling a harness's own faint suggestion from a
	// person's unsubmitted line needs the attributes
	// (internal/screen/fixture/classify.go's faintPlaceholder).
	ReadAgentStyled(ctx context.Context, handle runtime.AgentHandle, source harness.ReadSource, lines int) (string, error)
}

var _ Runtime = runtime.Adapter(nil)

// Rounds is a Runtime that keeps Herdr's answers for one round
// (runtime.RoundCache, CON-264): one `agent list` per session answers every
// crew's InspectAgent, and the timeline ingest's after it. Poll begins a
// round on it before it asks anything, so no answer outlives the round that
// asked. A Runtime that is not one asks Herdr every time, as before.
type Rounds interface {
	BeginRound()
}

var _ Rounds = (*runtime.RoundCache)(nil)

// HandleFunc resolves the Herdr handle recorded for one crew, and the screen
// profile of the harness it runs. It is a seam rather than a direct call to
// internal/spawn so the observer imports no code that can start or stop an
// agent; cmd/mate wires spawn.CrewHandle into it.
//
// An error means "mate could not resolve a handle to look at" - the crew has
// no agent recorded, or the Herdr session is not running. That is not a
// finding: the observer skips the crew for this round and concludes nothing.
type HandleFunc func(ctx context.Context, project, crew string) (runtime.AgentHandle, harness.ScreenProfile, error)

// Clock is where the observer reads the time. Tests pass a fake one so every
// threshold in this package can be crossed without waiting.
type Clock interface {
	Now() time.Time
}

// Sleeper is the pause between polls, interruptible by the context. Tests
// pass one that returns immediately, or one that ends the run.
type Sleeper interface {
	Sleep(ctx context.Context, d time.Duration) error
}

// Round is what one round examined (CON-267): every project on a sweep,
// else the projects it read. The timeline ingest, the budget and the
// no-turn checks after it read the same projects and no others: what the
// round did not examine has not changed since it last was, or will be
// read by the next sweep.
type Round struct {
	// Sweep is a round that examined every project and trusts nothing it
	// remembered of the files (DefaultSweepInterval).
	Sweep bool
	// Projects are the projects the round examined, in workspace.yaml's
	// order; every project on a sweep.
	Projects []string
}

// Ingest is the derived timeline's writer (internal/timeline.Ingester,
// IngestScope). The observer calls it at the end of every poll, because
// the observer is the process that already holds the workspace and is
// already awake every five seconds - and because M5 makes it the single
// writer of the database.
//
// It is an interface rather than the concrete type so this package keeps
// importing nothing that can start an agent, and so a Watcher with no
// timeline (every test that predates M5) still runs.
type Ingest interface {
	Ingest(ctx context.Context, round Round) error
}

// BudgetChecker is mvp.md M5 task 27's budget check
// (internal/timeline.Ingester.CheckBudgetsScope): at the end of a poll, for
// every project of the round with a `budget:` block in project.yaml, it
// opens a `budget` incident on whichever crew (or the project) has crossed
// a limit.
//
// It is a seam for the same reason Ingest is: this package must not import
// internal/db or parse project.yaml's price arithmetic itself, and a
// Watcher with none (every test and every build that predates M5) simply
// checks no budgets, which is the honest behaviour for a workspace with no
// timeline database.
type BudgetChecker interface {
	CheckBudgets(ctx context.Context, round Round) error
}

// TurnChecker is the observer's no-turn check
// (internal/timeline.Ingester.CheckNoTurnScope, CON-74): at the end of a
// poll, it opens a `no_turn` incident on every working crew of the round's
// projects launched more than after ago whose actor has no turn in the
// ledger, and resolves it once a turn is recorded. A seam for the same
// reason BudgetChecker is.
type TurnChecker interface {
	CheckNoTurn(ctx context.Context, round Round, after time.Duration) error
}

// Notices is what the core's change feed heard since the observer last
// asked (CON-260): the projects a writer's notice named, and whether one
// named the workspace's own files or the core acted itself.
type Notices struct {
	// Projects are the projects a notice named.
	Projects map[string]bool
	// All is a notice about the workspace's own files (workspace.yaml,
	// pricing), or an action the core ran, whose writes tell no feed: the
	// next round is a sweep.
	All bool
}

// NoticeFunc hands over what the feed heard since its last call, and
// starts a new set. internal/core wires its change feed in (Hear); this
// package imports nothing of it.
type NoticeFunc func() Notices

// SessionFunc asks whether the workspace's terminal runtime is up, without
// starting one. It is a seam for the same reason HandleFunc is: the
// observer imports nothing that can start an agent. nil means the runtime
// answered; an error describes why it could not be reached, and is what the
// runtime notice carries.
//
// It is called once per poll, independent of whether any crew exists: a
// workspace with a Mate and no open crews must still hear that Herdr is
// gone. Asking per crew cannot, because without a crew there is no handle
// to resolve and so no question to ask.
type SessionFunc func(ctx context.Context) error

// Deps are the observer's collaborators. Runtime and Handle are required;
// everything else has a default.
type Deps struct {
	Runtime Runtime
	Handle  HandleFunc
	// Session checks the workspace session itself, once per poll. Nil means
	// the session is not asked, and the observer stands a runtime notice
	// only for a crew's Herdr call that ran out HerdrTimeout (a build or
	// test that predates it).
	Session SessionFunc
	Clock   Clock
	Sleeper Sleeper
	// Timeline records what the files say into `.mate/mate.db`. Nil
	// means no timeline is kept, which is a workspace opened by a build that
	// has none and not an error.
	Timeline Ingest
	// Budget checks every project's `budget:` block against the ledger
	// Timeline just wrote, at the end of every poll. Nil means no budgets
	// are checked, the honest behaviour when there is no timeline database
	// to check them against.
	Budget BudgetChecker
	// Turns checks every open crew for a turn in the ledger Timeline just
	// wrote, at the end of every poll. Nil means no such check, as for
	// Budget.
	Turns TurnChecker
	// PollInterval is the pause between rounds; zero means
	// DefaultPollInterval.
	PollInterval time.Duration
	// StaleAfter is the quiet period that opens `stale`; zero means
	// DefaultStaleAfter.
	StaleAfter time.Duration
	// NoTurnAfter is how long after launch a working crew with no turn in
	// the ledger opens `no_turn`; zero means DefaultNoTurnAfter.
	NoTurnAfter time.Duration
	// HerdrTimeout bounds each Herdr call of a round; zero means
	// DefaultHerdrTimeout.
	HerdrTimeout time.Duration
	// Observer reads each pane snapshot that differs from the last one it
	// read; the health column is its Observation's composer, kept while the
	// pane does not change. It is asked in the background, so a slow answer
	// never holds the poll loop (observe.go). Nil means the fixture
	// observer (internal/screen/fixture), asked in the round.
	Observer screen.Observer
	// Pass records each round the loop runs, the timeline's ingest
	// included, for `mate core stats` (CON-263). Nil records nothing.
	Pass *loopstat.Loop
	// Notices is the change feed (CON-267). Nil means none: every round
	// examines every project and is a sweep, as before the feed; a
	// Watcher that hears one examines between sweeps only the projects
	// that can change unheard and those a notice named.
	Notices NoticeFunc
	// SweepInterval is how often a round that hears the feed examines
	// every project anyway; zero means DefaultSweepInterval.
	SweepInterval time.Duration
	// SharedConfig says the caller keeps the workspace handle's
	// workspace.yaml current itself - internal/core's engine (CON-270) -
	// so a round does not read the file again. False re-reads it every
	// round.
	SharedConfig bool
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

func (d Deps) pollInterval() time.Duration {
	if d.PollInterval > 0 {
		return d.PollInterval
	}
	return DefaultPollInterval
}

func (d Deps) staleAfter() time.Duration {
	if d.StaleAfter > 0 {
		return d.StaleAfter
	}
	return DefaultStaleAfter
}

func (d Deps) noTurnAfter() time.Duration {
	if d.NoTurnAfter > 0 {
		return d.NoTurnAfter
	}
	return DefaultNoTurnAfter
}

func (d Deps) sweepInterval() time.Duration {
	if d.SweepInterval > 0 {
		return d.SweepInterval
	}
	return DefaultSweepInterval
}

func (d Deps) herdrTimeout() time.Duration {
	if d.HerdrTimeout > 0 {
		return d.HerdrTimeout
	}
	return DefaultHerdrTimeout
}

// Watcher observes the open crews of a workspace.
//
// One Watcher polls from one goroutine: Start owns it, and a caller driving
// Poll itself (a test) is that goroutine. Snapshot and Health may be called
// from any other goroutine at any time.
//
// The workspace handle is this package's own: store.Workspace caches
// workspace.yaml and re-reads it on LoadConfig, so the observer opens a
// second handle rather than sharing the console's (cmd/mate does the
// opening).
type Watcher struct {
	ws   Workspace
	deps Deps

	// crews is the per-crew memory of the previous rounds. Only the polling
	// goroutine touches it.
	crews map[CrewRef]*observation
	// unanswered are the crews of this round whose Herdr call ran out its
	// bound. Only the polling goroutine touches it.
	unanswered []CrewRef
	// looks is, by project, what the last round that examined it found
	// (CON-267), and swept when the last sweep began. Only the polling
	// goroutine touches them.
	looks map[string]look
	swept time.Time
	// logs is, by project, what the rounds have read of incidents.log, and
	// closed the crews whose meta a round found closed, by the stat it had
	// (openIncidents, openCrews). Only the polling goroutine touches them.
	logs   map[string]*incidentLog
	closed map[string]map[string]metaStat

	mu     sync.Mutex
	health map[CrewRef]Health
	// runtimeNotice is the standing word about a Herdr the observer could
	// not reach, and runtimeNoticeAt when it was recorded. Empty means the
	// last round that asked got an answer. It is a copy of the round's own
	// read failure, not a conclusion about any crew (mvp.md decision 8).
	runtimeNotice   string
	runtimeNoticeAt time.Time
	cancel          context.CancelFunc
	done            chan struct{}

	// observing counts the observer calls in flight (observe.go).
	observing sync.WaitGroup

	// wake cuts the pause between two rounds short (Wake); it holds one
	// wake at most.
	wake chan struct{}
}

// observation is what the previous polls established about one crew.
type observation struct {
	// screen is the hash of the last pane snapshot, empty before the first
	// readable one.
	screen string
	// statusAt is the resume offset of the crew's status file.
	statusAt int64
	// verb is the last status verb the crew wrote, empty if it has written
	// nothing.
	verb string
	// changedAt is the last time the pane or the status file moved. It is
	// set on the first observation too: the observer counts quiet from when
	// it started looking, never from a past it did not see.
	changedAt time.Time
	// composer is the last composer classification, and composerSince
	// when it last changed. A busy pane redraws its spinner every poll, so
	// its quiet clock is always near zero; "busy for how long" is this
	// clock instead.
	composer      send.ComposerState
	composerSince time.Time
	// observed is the hash of the snapshot composer was read from, empty
	// before the observer has read one.
	observed string
	// reading is the observer's reading of the snapshot observed names.
	reading reading
	// asking is the observer call in flight for this crew, nil when none
	// (observe.go).
	asking *asking
	// stateOf is the hash of the snapshot state was named for, empty before
	// the stale rule has asked (observe.go's stillWorking).
	stateOf string
	state   screen.PaneState
}

// New builds an observer over a workspace. It does not poll until Start or
// Poll is called.
func New(ws Workspace, deps Deps) *Watcher {
	return &Watcher{
		ws:     ws,
		deps:   deps,
		crews:  make(map[CrewRef]*observation),
		looks:  make(map[string]look),
		logs:   make(map[string]*incidentLog),
		closed: make(map[string]map[string]metaStat),
		health: make(map[CrewRef]Health),
		wake:   make(chan struct{}, 1),
	}
}

// Wake starts the next round now rather than at the end of the pause,
// though no sooner than WakeGap after the last one started: a mate process
// wrote a record (the core's change feed, CON-260). It never blocks, and a
// wake that comes during a round starts the next one once that round ends.
// The rounds the pause times still come; a Watcher with Deps.Sleeper (a
// test's clock) is not woken.
func (w *Watcher) Wake() {
	select {
	case w.wake <- struct{}{}:
	default:
	}
}

// Hear makes f the Watcher's change feed (Deps.Notices). It is how the
// engine, which owns the feed, hands it to an observer cmd/mate built; call
// it before Start.
func (w *Watcher) Hear(f NoticeFunc) { w.deps.Notices = f }

// Start begins polling in its own goroutine until Stop or a cancelled
// context. Calling it twice is a no-op.
func (w *Watcher) Start(ctx context.Context) {
	w.mu.Lock()
	if w.cancel != nil {
		w.mu.Unlock()
		return
	}
	runCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	done := make(chan struct{})
	w.done = done
	w.mu.Unlock()
	go w.run(runCtx, done)
}

// Stop ends the polling goroutine and waits for the round in flight and
// for the observer calls it started, which the cancelled context ends. It
// is safe to call on a Watcher that was never started.
func (w *Watcher) Stop() {
	w.mu.Lock()
	cancel, done := w.cancel, w.done
	w.cancel, w.done = nil, nil
	w.mu.Unlock()
	if cancel == nil {
		return
	}
	cancel()
	<-done
	w.AwaitObserver()
}

// run closes done, the channel Start made, not w.done: Stop clears that
// field, possibly before this goroutine has first run.
func (w *Watcher) run(ctx context.Context, done chan struct{}) {
	defer close(done)
	for {
		// A failed round is not fatal: Herdr may be restarting, a project
		// directory may be half-written. The next round asks again, and
		// what the observer could not read it does not conclude from.
		started := time.Now()
		start := w.deps.Pass.Begin()
		_ = w.Poll(ctx)
		w.deps.Pass.End(start)
		if err := w.pause(ctx, started); err != nil {
			return
		}
		if ctx.Err() != nil {
			return
		}
	}
}

// pause waits out the poll interval, or, once Wake is called, until
// WakeGap after the round that started at started. The interval's own
// rounds come as before: a wake only brings one forward.
func (w *Watcher) pause(ctx context.Context, started time.Time) error {
	if w.deps.Sleeper != nil {
		return w.deps.sleep(ctx, w.deps.pollInterval())
	}
	timer := time.NewTimer(w.deps.pollInterval())
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	case <-w.wake:
	}
	gap := min(time.Until(started.Add(WakeGap)), w.deps.pollInterval())
	if gap > 0 {
		if err := w.deps.sleep(ctx, gap); err != nil {
			return err
		}
	}
	// What woke it during the gap was written before the round about to
	// start, which reads it: one more wake would be one more round for
	// nothing. A wake during that round is kept.
	select {
	case <-w.wake:
	default:
	}
	return nil
}

// Poll runs one round over every open crew of every project: it appends
// whatever incident transitions it found to `incidents.log` and replaces the
// health snapshot.
//
// The error it returns is the joined per-crew and per-project read failures
// of the round; the round itself always finishes. A caller that wants to
// know whether anything was observed reads Snapshot.
func (w *Watcher) Poll(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if rounds, ok := w.deps.Runtime.(Rounds); ok {
		rounds.BeginRound()
	}
	// Re-read workspace.yaml every round: a project registered by another
	// process (or by the console beside this one) must come under watch
	// without a restart. Under the engine the handle is the engine's, which
	// reloads it on a `workspace` notice and when the file changes.
	if !w.deps.SharedConfig {
		if err := w.ws.LoadConfig(); err != nil {
			return err
		}
	}

	now := w.deps.now()
	results := make(map[CrewRef]Health)
	seen := make(map[CrewRef]bool)
	listed := make(map[string]bool)
	w.unanswered = nil
	var errs []error

	round := w.beginRound(now)
	examined := make(map[string]bool)
	for _, ref := range w.ws.Projects() {
		if err := ctx.Err(); err != nil {
			return err
		}
		examined[ref.Name] = true
		if !round.Sweep && !w.due(ref.Name, round) {
			continue
		}
		round.Projects = append(round.Projects, ref.Name)
		busy, err := w.pollProject(ctx, ref.Name, round.Sweep, now, results, seen, listed)
		if err != nil {
			errs = append(errs, err)
		}
		w.looks[ref.Name] = look{busy: busy}
	}
	for project := range w.looks {
		if !examined[project] {
			// No longer in workspace.yaml: registered again, it is
			// examined as a project never seen.
			delete(w.looks, project)
			delete(w.logs, project)
			delete(w.closed, project)
		}
	}
	w.noteRuntime(ctx, now)
	w.commit(results, seen, listed)
	// The timeline is recorded last, after the health snapshot has been
	// swapped in: the ingest reads this round's readings (Watcher.Readings),
	// and any incident this round opened is already in `incidents.log`. A
	// failed ingest joins the round's other errors and never stops the
	// observing - the database is derived, and a round that could not write
	// it still watched the crews.
	if w.deps.Timeline != nil {
		if err := w.deps.Timeline.Ingest(ctx, round.Round); err != nil {
			errs = append(errs, fmt.Errorf("watch: record the timeline: %w", err))
		}
	}
	// The budget check reads the ledger the ingest above just wrote, so it
	// has to run after it: a crossing this poll's turns caused must already
	// be in the database before anything can check it against a limit.
	if w.deps.Budget != nil {
		if err := w.deps.Budget.CheckBudgets(ctx, round.Round); err != nil {
			errs = append(errs, fmt.Errorf("watch: check budgets: %w", err))
		}
	}
	// The no-turn check reads the same ledger, for the same reason.
	if w.deps.Turns != nil {
		if err := w.deps.Turns.CheckNoTurn(ctx, round.Round, w.deps.noTurnAfter()); err != nil {
			errs = append(errs, fmt.Errorf("watch: check turns: %w", err))
		}
	}
	return errors.Join(errs...)
}

// look is what the last round that examined a project found (CON-267).
type look struct {
	// busy is a project that can change without a notice: it has an open
	// crew - whose pane, `echo >> $MATE_STATUS` and transcript tell nobody
	// - or a live Mate, whose transcript grows unheard; or the round could
	// not tell, its listing having failed. Every round examines it, as
	// every round examined every project before CON-267.
	busy bool
}

// pendingRound is a round being run, with the notices it took.
type pendingRound struct {
	Round
	notices Notices
}

// beginRound decides whether this round is a sweep and takes what the
// change feed heard since the last one.
//
// Measured on the captain's workspace (2026-11-12, after CON-259, CON-260
// and CON-264): an observer round cost ~330-500 ms every ~5 s, ~9% of a
// core, though of its 12 projects at most a third had an open crew or a
// live Mate; the others were read whole - incidents.log, every crew's
// meta, the timeline's pass - to learn again that nothing had moved.
// Between sweeps a round now examines a project only when it is busy
// (look), when a notice named it, or when it has never been examined.
//
// A notice taken is read by this round: what a writer wrote before its
// notice is on disk before the round reads, and a notice that comes during
// the round is in the next set.
func (w *Watcher) beginRound(now time.Time) pendingRound {
	if w.deps.Notices == nil {
		return pendingRound{Round: Round{Sweep: true}}
	}
	notices := w.deps.Notices()
	sweep := notices.All || w.swept.IsZero() || now.Sub(w.swept) >= w.deps.sweepInterval()
	if sweep {
		w.swept = now
	}
	return pendingRound{Round: Round{Sweep: sweep}, notices: notices}
}

// due reports whether a round that is not a sweep examines project.
func (w *Watcher) due(project string, round pendingRound) bool {
	last, seen := w.looks[project]
	return !seen || last.busy || round.notices.Projects[project]
}

func (w *Watcher) pollProject(ctx context.Context, project string, sweep bool, now time.Time,
	results map[CrewRef]Health, seen map[CrewRef]bool, listed map[string]bool) (bool, error) {

	open, err := w.openIncidents(project, sweep)
	if err != nil {
		return true, fmt.Errorf("watch: read incidents of %s: %w", project, err)
	}
	crews, err := w.openCrews(project, sweep)
	if err != nil {
		return true, fmt.Errorf("watch: list crews of %s: %w", project, err)
	}
	listed[project] = true
	busy := len(crews) > 0 || w.mateLive(project)

	var errs []error
	for _, crew := range crews {
		if err := ctx.Err(); err != nil {
			return true, err
		}
		ref := CrewRef{Project: project, Crew: crew}
		seen[ref] = true
		if err := w.pollCrew(ctx, ref, now, open, results); err != nil {
			errs = append(errs, fmt.Errorf("watch: observe crew %s of %s: %w", crew, project, err))
		}
	}
	if err := w.resolveClosed(project, crews, now, open); err != nil {
		// An incident left open on a closed crew is tried again next
		// round, not at the next sweep.
		busy = true
		errs = append(errs, fmt.Errorf("watch: resolve incidents of closed crews of %s: %w", project, err))
	}
	return busy, errors.Join(errs...)
}

// mateLive reports whether the project's mate.meta records a Mate that was
// started and not stopped (store.MateRecorded). Its transcript grows with
// every turn and tells no feed, so the timeline must read it every round;
// a stop rewrites mate.meta through the store, which is heard. A meta that
// cannot be read counts as live: the round could not tell.
func (w *Watcher) mateLive(project string) bool {
	meta, err := w.ws.ReadMateMeta(project)
	if err != nil {
		return true
	}
	return store.MateRecorded(meta)
}

// stopped reports whether an agent Herdr no longer has went away because
// the crew was stopped: a `crew stop` holds the crew's stop lock from before
// it stops the agent until its meta closes the task, and a stop that
// finished after this round listed the crew has already closed the meta.
// Either way the agent left on purpose, and that is not `runtime_lost`.
//
// The lock is asked first: a stop releases it only after writing the meta,
// so a free lock and an open meta read in that order cannot straddle one.
func (w *Watcher) stopped(ref CrewRef) (bool, error) {
	stopping, err := w.ws.CrewStopping(ref.Project, ref.Crew)
	if err != nil || stopping {
		return stopping, err
	}
	meta, err := w.ws.ReadCrewMeta(ref.Project, ref.Crew)
	if err != nil {
		return false, err
	}
	return crewstate.Declare(crewstate.Declaration{Meta: meta}).Closed(), nil
}

// resolveClosed resolves the incidents still open on crews of the project
// that have since closed. A closed crew is not watched, so nothing else would
// ever clear them. A crew whose meta cannot be read is left alone: an
// unreadable record says nothing about whether the task is over.
func (w *Watcher) resolveClosed(project string, watched []string, now time.Time, open map[incidentKey]bool) error {
	isWatched := make(map[string]bool, len(watched))
	for _, crew := range watched {
		isWatched[crew] = true
	}
	var keys []incidentKey
	for key, isOpen := range open {
		if isOpen && !isWatched[key.ref.Crew] {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].ref.Crew != keys[j].ref.Crew {
			return keys[i].ref.Crew < keys[j].ref.Crew
		}
		return keys[i].kind < keys[j].kind
	})
	var errs []error
	for _, key := range keys {
		meta, err := w.ws.ReadCrewMeta(project, key.ref.Crew)
		if err != nil || len(meta) == 0 || !crewstate.Declare(crewstate.Declaration{Meta: meta}).Closed() {
			continue
		}
		if err := w.resolveIncident(key.ref, key.kind, now, open,
			fmt.Sprintf("crew %s is closed", key.ref.Crew)); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// pollCrew is the whole contract of mvp.md section 4b for one crew.
func (w *Watcher) pollCrew(ctx context.Context, ref CrewRef, now time.Time,
	open map[incidentKey]bool, results map[CrewRef]Health) error {

	obs := w.observation(ref)

	first := obs.changedAt.IsZero()
	statusMoved, err := w.readStatus(ref, obs)
	if err != nil {
		// The crew's own file is unreadable. Nothing about the pane or the
		// agent has been established either, so nothing is concluded.
		return err
	}
	if first {
		// The first look at a crew reads whatever its status file already
		// holds. That backlog is not movement: a console opened beside a
		// crew that went quiet an hour ago must not report the crew as
		// having just spoken, and must not resolve an incident the previous
		// console left open. The quiet clock starts here instead.
		statusMoved = false
		obs.changedAt = now
	} else if statusMoved {
		obs.changedAt = now
	}

	// Every Herdr call is bounded on its own (Deps.HerdrTimeout): the
	// crews are polled one after another, so one call that never returns
	// would otherwise hold every crew after it, the runtime notice and the
	// timeline.
	hctx, cancel := context.WithTimeout(ctx, w.deps.herdrTimeout())
	handle, screens, err := w.deps.Handle(hctx, ref.Project, ref.Crew)
	cancel()
	if err != nil {
		if w.timedOut(ctx, ref, err) {
			return err
		}
		// No handle to look at: the crew records no agent yet, or the Herdr
		// session is not running. Neither says the crew is in trouble.
		return nil
	}

	hctx, cancel = context.WithTimeout(ctx, w.deps.herdrTimeout())
	_, err = w.deps.Runtime.InspectAgent(hctx, handle)
	cancel()
	if err != nil {
		if !runtime.IsAgentNotFound(err) {
			// Herdr could not answer. "Mate could not look" is not "the
			// agent is gone" (mvp.md decision 8).
			w.timedOut(ctx, ref, err)
			return err
		}
		if stopped, err := w.stopped(ref); err != nil || stopped {
			return err
		}
		if err := w.openIncident(ref, box.IncidentRuntimeLost, now, open,
			fmt.Sprintf("herdr no longer lists agent %s", handle.Name)); err != nil {
			return err
		}
		results[ref] = Health{
			AgentPresent: false,
			Composer:     send.StateUnknown,
			QuietFor:     now.Sub(obs.changedAt),
			ObservedAt:   now,
		}
		return nil
	}
	if err := w.resolveIncident(ref, box.IncidentRuntimeLost, now, open,
		fmt.Sprintf("herdr lists agent %s again", handle.Name)); err != nil {
		return err
	}

	// Every open crew's pane is read every round, on purpose. Nothing Herdr
	// lists says the screen is unchanged: on Herdr 0.9.5 the `revision` of
	// an agent and of a pane was measured not to move with output
	// (2026-11-11, CON-264 review: 0 before and after an echo and a
	// three-line loop that `pane read` showed; crews running for hours list
	// revision 2), and agent_status / state_change_seq do not move while a
	// working harness redraws or an idle one gets a draft. Gating the read
	// on any of them would hold a stale snapshot and starve the stale rule
	// and the classification. Change detection belongs to Herdr's event
	// subscription (CON-265), not to the list.
	hctx, cancel = context.WithTimeout(ctx, w.deps.herdrTimeout())
	screen, err := w.deps.Runtime.ReadAgentStyled(hctx, handle, screens.ReadSource(), send.DefaultLines)
	cancel()
	if err != nil {
		// The agent is there and the pane is not readable: an unread screen
		// is not a quiet one, so the round ends without a verdict.
		w.timedOut(ctx, ref, err)
		return err
	}
	hash := screenHash(screen)
	paneMoved := obs.screen != "" && obs.screen != hash
	held := obs.screen == hash
	obs.screen = hash
	if paneMoved {
		obs.changedAt = now
	}

	// The observer is asked only about a snapshot it has not read and the
	// pane has held for a poll: a pane that has not changed keeps its last
	// reading, so a still pane costs one observer call, and a pane that
	// changes every poll costs none (docs/plans/jev-observer-2026-10-08.md
	// section 4.3).
	// A screen the observer could not read is not a reading: the round ends
	// without a verdict, as for a pane that could not be read, and the next
	// round asks again.
	read, err := w.composer(ctx, obs, screens, screen, hash, held)
	if err != nil {
		return err
	}
	composer := read.composer
	if obs.composerSince.IsZero() || composer != obs.composer {
		obs.composer = composer
		obs.composerSince = now
	}
	quiet := now.Sub(obs.changedAt)

	// The stale incident is decided on the fixture's own composer reading
	// (Deterministic): Jev's composer veto neither keeps it from opening nor
	// resolves it. Before it opens, the pane's state is asked once per still
	// snapshot (stillWorking), stub first and Jev only when the stub cannot
	// name it: a pane between turns that waits on work it left running is
	// not stale (2026-10-13). The column shows the observer's composer.
	deterministic := read.deterministic
	switch {
	case open[incidentKey{ref, box.IncidentStale}]:
		if reason, ok := staleCleared(paneMoved, statusMoved, deterministic); ok {
			if err := w.resolveIncident(ref, box.IncidentStale, now, open, reason); err != nil {
				return err
			}
		}
	case quiet >= w.deps.staleAfter() && deterministic != send.StateBusy && !waiting(obs.verb):
		working, err := w.stillWorking(ctx, obs, screens, screen, hash)
		if err != nil {
			return err
		}
		if working {
			break
		}
		text := fmt.Sprintf("no status line and no pane change for %s; composer %s",
			quiet.Round(time.Second), send.ComposerLabel(deterministic))
		if obs.verb != "" {
			text += fmt.Sprintf("; last status verb %s", obs.verb)
		}
		if err := w.openIncident(ref, box.IncidentStale, now, open, text); err != nil {
			return err
		}
	}

	results[ref] = Health{
		AgentPresent: true,
		Composer:     composer,
		Source:       read.source,
		QuietFor:     quiet,
		ComposerFor:  now.Sub(obs.composerSince),
		ObservedAt:   now,
	}
	return nil
}

// timedOut reports whether err is a Herdr call that ran out its own bound
// rather than the round being cancelled, and records the crew as unanswered
// this round when it is.
func (w *Watcher) timedOut(round context.Context, ref CrewRef, err error) bool {
	if !errors.Is(err, context.DeadlineExceeded) || round.Err() != nil {
		return false
	}
	w.unanswered = append(w.unanswered, ref)
	return true
}

// RuntimeNotice is the observer's standing word about a Herdr it could not
// reach, and when it last failed to. Empty means the last round that asked
// got an answer. It is safe to call from any goroutine.
func (w *Watcher) RuntimeNotice() (string, time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.runtimeNotice, w.runtimeNoticeAt
}

// noteRuntime asks the workspace session once per poll and records the
// answer, or clears the notice when the runtime answered. The check is the
// workspace's own rather than a crew's: with no open crew there is no handle
// to resolve, and without one the notice would clear while Herdr is gone
// (the bug this replaces), never appear at all, or depend on which crews
// happen to be open. "Herdr could not answer" is never an incident
// (TestWatchConcludesNothingWhenHerdrCannotAnswer); this is the separate
// signal a reader needs, because while it stands every file-read row on
// screen is as fresh as its last read and no row can be re-read.
//
// A session that answered with crews whose own Herdr call ran out its bound
// this round still stands a notice, naming them: those crews keep the
// Health of their last answered round (commit), which on its own reads as
// current as any other row.
func (w *Watcher) noteRuntime(ctx context.Context, now time.Time) {
	if ctx.Err() != nil {
		return
	}
	var err error
	if w.deps.Session != nil {
		sctx, cancel := context.WithTimeout(ctx, w.deps.herdrTimeout())
		err = w.deps.Session(sctx)
		cancel()
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	var notice string
	switch {
	case err != nil && errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil:
		notice = fmt.Sprintf("herdr did not answer within %s; health shown is from the last round it did",
			w.deps.herdrTimeout())
	case err != nil:
		notice = noticeLine(err.Error())
		var coded *observability.Error
		if errors.As(err, &coded) && strings.TrimSpace(coded.Message) != "" {
			notice = noticeLine(coded.Message)
		}
	case len(w.unanswered) > 0:
		notice = unansweredNotice(w.unanswered, w.deps.herdrTimeout())
	default:
		w.runtimeNotice, w.runtimeNoticeAt = "", time.Time{}
		return
	}
	w.runtimeNotice, w.runtimeNoticeAt = notice, now
}

// unansweredNotice names the crews whose Herdr call ran out its bound.
func unansweredNotice(refs []CrewRef, bound time.Duration) string {
	names := make([]string, 0, len(refs))
	for _, ref := range refs {
		names = append(names, fmt.Sprintf("%s of %s", ref.Crew, ref.Project))
	}
	if len(names) == 1 {
		return fmt.Sprintf("herdr did not answer within %s for crew %s; its health is from its last answered round",
			bound, names[0])
	}
	return fmt.Sprintf("herdr did not answer within %s for crews %s; their health is from their last answered round",
		bound, strings.Join(names, ", "))
}

// noticeLine collapses whitespace so a multi-line error reads as one notice.
func noticeLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// staleCleared is the resolve side of the stale rule: the pane moved, the
// crew spoke, or the harness is working again.
func staleCleared(paneMoved, statusMoved bool, composer send.ComposerState) (string, bool) {
	switch {
	case statusMoved:
		return "the crew wrote a status line", true
	case paneMoved:
		return "the pane changed", true
	case composer == send.StateBusy:
		return "the composer is busy again", true
	default:
		return "", false
	}
}

// waiting reports whether a status verb means the crew is waiting on a human
// rather than stuck. A crew that asked a question or handed back to the Mate
// is silent on purpose (mvp.md section 4b); calling that stale would put the
// same crew in the inbox twice.
//
// The vocabulary is box's, so there is one spelling of it in the codebase -
// including the legacy `done:`, which box already reads as `wait-mate`.
func waiting(verb string) bool {
	switch box.State(verb) {
	case box.StateNeedsDecision, box.StateWaitMate:
		return true
	default:
		return false
	}
}

// readStatus advances the crew's status cursor and reports whether the crew
// appended anything since the previous round.
func (w *Watcher) readStatus(ref CrewRef, obs *observation) (bool, error) {
	lines, next, err := w.ws.ReadStatus(ref.Project, ref.Crew, obs.statusAt)
	if err != nil {
		return false, err
	}
	obs.statusAt = next
	if len(lines) == 0 {
		return false, nil
	}
	for i := len(lines) - 1; i >= 0; i-- {
		if verb := statusVerb(lines[i].Line); verb != "" {
			obs.verb = verb
			break
		}
	}
	return true, nil
}

// statusVerb is the crew verb of one `state: one line`, parsed by
// internal/box so the observer reads exactly what the console and the CLI
// read - the three verbs of mvp.md section 4b plus the pre-4b spellings box
// maps onto them. A line box does not recognise yields "", and readStatus
// then keeps the verb the crew last actually wrote: a malformed echo is not
// a state change.
func statusVerb(line string) string {
	if st := box.ParseStatus(line); st.State != box.StateUnknown {
		return string(st.State)
	}
	return ""
}

func (w *Watcher) observation(ref CrewRef) *observation {
	obs, ok := w.crews[ref]
	if !ok {
		obs = &observation{}
		w.crews[ref] = obs
	}
	return obs
}

// commit replaces the health snapshot and forgets the crews that are gone.
//
// A crew that is still listed but produced no observation this round keeps
// the one it had: the round failed to look, which is not the same as the
// crew having no health, and Health.ObservedAt already says how old the
// answer is. A crew of a project whose listing failed is kept for the same
// reason.
func (w *Watcher) commit(results map[CrewRef]Health, seen map[CrewRef]bool, listed map[string]bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for ref, h := range w.health {
		if _, fresh := results[ref]; fresh {
			continue
		}
		if seen[ref] || !listed[ref.Project] {
			results[ref] = h
		}
	}
	w.health = results
	for ref := range w.crews {
		if !seen[ref] && listed[ref.Project] {
			delete(w.crews, ref)
		}
	}
}

// incidentKey is one open incident: a crew and a kind, which is all the
// identity an incident has (mvp.md section 4b).
type incidentKey struct {
	ref  CrewRef
	kind box.IncidentKind
}

// openIncidents reads which incidents of a project are currently open: the
// file is the state, so a console that restarted sees exactly what the
// previous one left behind.
//
// The log is append-only, so a round reads only what was appended since
// the last round that read it (CON-267), onto what that round knew; a
// sweep reads it from the start. Measured 2026-11-12 on a copy of the
// captain's workspace: once a round skipped the idle projects, reading the
// busy ones' whole incidents.log and every crew's meta was most of what
// was left of its 10.5 MB.
func (w *Watcher) openIncidents(project string, sweep bool) (map[incidentKey]bool, error) {
	log := w.logs[project]
	if log == nil || sweep {
		log = &incidentLog{open: make(map[incidentKey]bool)}
	}
	entries, next, err := w.ws.ReadIncidents(project, log.at)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		key := incidentKey{CrewRef{Project: project, Crew: e.Crew}, box.IncidentKind(e.Kind)}
		log.open[key] = e.Open()
	}
	log.at = next
	w.logs[project] = log
	return log.open, nil
}

// incidentLog is what the rounds have read of a project's incidents.log:
// up to where, and which incidents it says are open. Its map is the one a
// round hands its rules, so what they append is in it at once, and is read
// back, to the same effect, by the next round.
type incidentLog struct {
	at   int64
	open map[incidentKey]bool
}

func (w *Watcher) openIncident(ref CrewRef, kind box.IncidentKind, now time.Time,
	open map[incidentKey]bool, text string) error {

	key := incidentKey{ref, kind}
	if open[key] {
		return nil
	}
	if err := w.appendIncident(ref, kind, store.IncidentOpen, now, text); err != nil {
		return err
	}
	open[key] = true
	return nil
}

func (w *Watcher) resolveIncident(ref CrewRef, kind box.IncidentKind, now time.Time,
	open map[incidentKey]bool, text string) error {

	key := incidentKey{ref, kind}
	if !open[key] {
		return nil
	}
	if err := w.appendIncident(ref, kind, store.IncidentResolved, now, text); err != nil {
		return err
	}
	open[key] = false
	return nil
}

func (w *Watcher) appendIncident(ref CrewRef, kind box.IncidentKind, state string, now time.Time, text string) error {
	return w.ws.AppendIncident(ref.Project, store.IncidentEntry{
		Time:  now,
		Crew:  ref.Crew,
		Kind:  string(kind),
		State: state,
		Text:  text,
	})
}

// openCrews lists the crews of a project that are still in flight: a
// `crews/<id>.meta` whose declared state is not terminal. The decision is
// crewstate.Declare's, which is the one resolution order in the codebase
// (mvp.md section 4b) and a leaf package that can start nothing - so the
// observer still depends on nothing that can start or stop an agent. An
// incident is deliberately not an input here: `blocked` is a state the
// observer itself produces, and a crew it has flagged is exactly the crew it
// must keep watching so it can clear the flag again.
//
// A meta that cannot be read is skipped rather than watched: an unreadable
// record is not evidence of a running crew.
//
// A meta a round found closed is not read again while its size and mtime
// hold (CON-267), stat'ed before the read as seal.go's keys are: what
// reopens a crew - `mate crew relaunch`, a person's edit - rewrites the
// file whole. A sweep reads every meta. On the captain's workspace a busy
// project holds hundreds of closed crews and a handful of open ones.
func (w *Watcher) openCrews(project string, sweep bool) ([]string, error) {
	entries, err := os.ReadDir(w.ws.CrewsDir(project))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	closed := w.closed[project]
	if closed == nil || sweep {
		closed = make(map[string]metaStat)
	}
	next := make(map[string]metaStat, len(closed))
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		id, ok := strings.CutSuffix(e.Name(), ".meta")
		if !ok || store.ValidateCrewID(id) != nil {
			continue
		}
		st, statted := statOf(e)
		if statted {
			if was, ok := closed[id]; ok && was == st {
				next[id] = st
				continue
			}
		}
		meta, err := w.ws.ReadCrewMeta(project, id)
		if err != nil {
			continue
		}
		if crewstate.Declare(crewstate.Declaration{Meta: meta}).Closed() {
			if statted {
				next[id] = st
			}
			continue
		}
		out = append(out, id)
	}
	w.closed[project] = next
	sort.Strings(out)
	return out, nil
}

// metaStat is a crew meta's size and modification time.
type metaStat struct{ size, mtime int64 }

func statOf(e os.DirEntry) (metaStat, bool) {
	fi, err := e.Info()
	if err != nil {
		return metaStat{}, false
	}
	return metaStat{fi.Size(), fi.ModTime().UnixNano()}, true
}

// screenHash is the pane snapshot reduced to something comparable between
// rounds. The screen itself is not kept: the observer answers "did this
// change", never "what did it say".
func screenHash(screen string) string {
	sum := sha256.Sum256([]byte(screen))
	return hex.EncodeToString(sum[:])
}
