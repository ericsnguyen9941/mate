// Package chain is the screen.Observer that reads every snapshot with the
// fixture observer and asks Jev only about a snapshot the fixture cannot
// read (docs/plans/jev-observer-2026-10-08.md, ruling of 2026-10-14).
//
// Jev is the last resort for a screen a harness update drew and no stub
// knows yet, not a second reader of screens the stubs already name. A
// snapshot the fixture reads is its Observation whole, with no request.
// Of one it cannot read, Jev may name the dialog and the notice; the
// composer stays the fixture's unknown, so nothing is ever typed or
// confirmed on Jev's word. A send never asks Jev at all: it types only into
// a composer the fixture reads empty, which Jev's answer cannot change.
//
// Every screen the fixture cannot read is kept for the captain to write a
// stub from (WithScreens), whether Jev was asked or not.
package chain

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/screen"
)

// TTL is how long one answer from Jev stands for the same screen of the
// same harness.
const TTL = 60 * time.Second

// Why the chain returned the fallback's Observation, as the log names it.
const (
	FallbackError     = "error"
	FallbackThreshold = "below-threshold"
)

// Option configures a Chain.
type Option func(*Chain)

// WithClock sets the clock the dedup cache ages by and the log stamps with.
func WithClock(now func() time.Time) Option { return func(c *Chain) { c.now = now } }

// WithLog makes the chain hand one formatted line (log.go) to appendLine
// for every request it sends to the primary observer. The caller owns the
// file and its lock; an append that fails is dropped, never an
// observation's failure.
func WithLog(appendLine func(line string) error) Option {
	return func(c *Chain) { c.appendLine = appendLine }
}

// WithScreens makes the chain keep the text of every screen the fixture
// could not read under dir (screens.go), at most once per harness and
// screen within TTL. A failure to write is dropped.
func WithScreens(dir string) Option { return func(c *Chain) { c.screensDir = dir } }

// Chain is the chained observer.
type Chain struct {
	primary, fallback screen.Observer
	threshold         float64
	now               func() time.Time
	appendLine        func(string) error
	screensDir        string

	mu      sync.Mutex
	cache   map[cacheKey]cached
	states  map[cacheKey]cachedState
	kept    map[cacheKey]time.Time
	breaker breaker
}

type cacheKey struct {
	kind string
	hash [sha256.Size]byte
}

type cached struct {
	obs screen.Observation
	err error
	at  time.Time
}

// New chains primary (Jev) in front of fallback (the fixture observer).
// threshold is the confidence, in [0, 1], below which primary's answer is
// not used.
func New(primary, fallback screen.Observer, threshold float64, opts ...Option) screen.Observer {
	c := &Chain{primary: primary, fallback: fallback, threshold: threshold, now: time.Now, cache: map[cacheKey]cached{}, states: map[cacheKey]cachedState{},
		kept: map[cacheKey]time.Time{}}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Observe reads one snapshot, the fallback (the fixture) first:
//
//  1. A snapshot the fallback reads (Confidence 1: a composer, or a startup
//     screen its profile measured) is the fallback's Observation whole. The
//     primary (Jev) is not asked.
//  2. A send never asks the primary. It types only into a composer the
//     fallback reads empty, and the primary never makes one empty, so its
//     answer cannot change what a send does.
//  3. Otherwise the primary is asked. On an error or an answer below the
//     threshold the fallback's Observation is returned whole. Else the
//     fallback's, with the primary's Dialog, Notice, Confidence, Source and
//     Reason, and Highlight -1 when the primary names a dialog: settle
//     refuses such a screen, naming Jev's dialog. The composer stays the
//     fallback's unknown, and Draft, Startup, Evidence and Deterministic
//     are the fallback's, always.
//
// The primary is asked once per harness and screen hash within TTL; a
// repeat is answered from memory. After BreakerFailures failures in a row
// it is not asked at all for BreakerCooldown (breaker.go): the fallback's
// reading is returned whole, Reason "jev: circuit open", even for a screen
// whose earlier answer is still in memory - a cache hit costs no request,
// but an open circuit answers from the fixture alone. A screen the primary
// skips (a skipper with no text to send) is not asked, not counted by the
// breaker and not logged: the fallback's reading is returned whole. The
// call is synchronous, so an answer always belongs to the snapshot that was
// hashed: a pane that changes while Jev is asked is a new snapshot on the
// caller's next read, never this answer's.
func (c *Chain) Observe(ctx context.Context, profile harness.ScreenProfile, pane string) (screen.Observation, error) {
	fix, err := c.fallback.Observe(ctx, profile, pane)
	if err != nil || fix.Confidence == 1 {
		return fix, err
	}
	key := cacheKey{kind: kindOf(profile), hash: sha256.Sum256([]byte(pane))}
	stub := string(fix.Startup)
	if s, ok := c.primary.(skipper); ok && s.Skips(pane) {
		return withReason(fix, "jev: no terminal text to observe; not asked"), nil
	}
	if screen.CallerOf(ctx) == screen.CallerSend {
		c.keep(ctx, key, pane, stub, keptNotAsked)
		return withReason(fix, "jev: not asked for a send"), nil
	}
	ask, trial := c.admit()
	if !ask {
		c.keep(ctx, key, pane, stub, keptCircuitOpen)
		return withReason(fix, circuitOpen), nil
	}
	if !trial {
		if hit, ok := c.lookup(key); ok {
			obs, _ := decide(fix, hit.obs, hit.err, c.threshold)
			return obs, nil
		}
	}
	start := c.now()
	primary, perr := c.primary.Observe(ctx, profile, pane)
	latency := c.now().Sub(start)
	event := c.record(ctx, perr)
	if ctx.Err() == nil {
		c.store(key, cached{obs: primary, err: perr, at: c.now()})
	}
	obs, fallback := decide(fix, primary, perr, c.threshold)
	if c.appendLine != nil {
		_ = c.appendLine(LogLine{Time: start, Kind: key.kind, Hash: fmt.Sprintf("%x", key.hash[:6]), Latency: latency,
			Source: obs.Source, Composer: primary.Composer, Dialog: primary.Dialog, Confidence: primary.Confidence,
			Fixture: fix.Composer, Caller: screen.CallerOf(ctx), Used: used(fallback), Fallback: fallback}.String())
		if event != "" {
			_ = c.appendLine(BreakerLine(c.now(), event))
		}
	}
	answer := keptFailed
	if perr == nil {
		answer = fmt.Sprintf("composer=%s dialog=%s %.2f", primary.Composer, primary.Dialog, primary.Confidence)
	}
	c.keep(ctx, key, pane, stub, answer)
	return obs, nil
}

// skipper is a primary that can say, without a request, that a screen has
// nothing to ask it about (jev.Observer.Skips).
type skipper interface {
	Skips(pane string) bool
}

// decide applies Observe's rule 3 to one pair of readings. It names the
// fallback rule that applied, or "" when the primary's reading stands.
func decide(fix, primary screen.Observation, perr error, threshold float64) (screen.Observation, string) {
	switch {
	case perr != nil:
		return withReason(fix, "jev: "+perr.Error()), FallbackError
	case primary.Confidence < threshold:
		return withReason(fix, fmt.Sprintf("jev: confidence %.2f below %.2f (%s)", primary.Confidence, threshold, primary.Reason)), FallbackThreshold
	}
	out := fix
	out.Dialog, out.Notice = primary.Dialog, primary.Notice
	out.Confidence, out.Source, out.Reason = primary.Confidence, primary.Source, primary.Reason
	if primary.Dialog != screen.DialogNone {
		out.Highlight = -1
	}
	return out, ""
}

// used is the log's word for whether the primary's answer is in the
// reading returned: it is whenever no fallback rule applied.
func used(fallback string) string {
	if fallback == "" {
		return UsedYes
	}
	return UsedNo
}

// withReason is the fallback's Observation with why the primary's was not
// used put first in its Reason.
func withReason(fix screen.Observation, why string) screen.Observation {
	if fix.Reason != "" {
		why += "; " + fix.Reason
	}
	fix.Reason = why
	return fix
}

func (c *Chain) lookup(key cacheKey) (cached, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	hit, ok := c.cache[key]
	if !ok || c.now().Sub(hit.at) >= TTL {
		return cached{}, false
	}
	return hit, true
}

// store remembers one answer and forgets every expired one, so the cache
// holds at most the screens of the last TTL.
func (c *Chain) store(key cacheKey, entry cached) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.cache {
		if entry.at.Sub(v.at) >= TTL {
			delete(c.cache, k)
		}
	}
	c.cache[key] = entry
}

// kindOf names the harness a screen profile belongs to, "-" for none.
func kindOf(profile harness.ScreenProfile) string {
	if profile == nil {
		return "-"
	}
	return string(profile.Kind())
}
