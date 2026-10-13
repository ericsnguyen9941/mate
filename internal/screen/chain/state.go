package chain

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/screen/jev"
)

type cachedState struct {
	reading screen.StateReading
	err     error
	at      time.Time
}

// ObserveState names the state of a pane, stub first: when the fixture
// observer names one there is no request to Jev and no log line. Only a
// pane the fixture calls PaneUnknown is put to Jev; its answer stands at or
// above the threshold, else (or on an error) the fixture's PaneUnknown is
// returned with why. The answer is cached per harness and screen for TTL,
// logged, and, with WithScreens, the screen kept for the captain.
func (c *Chain) ObserveState(ctx context.Context, profile harness.ScreenProfile, pane string) (screen.StateReading, error) {
	fb, ok := c.fallback.(screen.StateObserver)
	if !ok {
		return screen.StateReading{}, errors.New("chain: fallback cannot observe pane state")
	}
	fix, err := fb.ObserveState(ctx, profile, pane)
	if err != nil || fix.State != screen.PaneUnknown {
		return fix, err
	}
	pr, ok := c.primary.(screen.StateObserver)
	if !ok {
		return withStateReason(fix, "jev: cannot observe pane state"), nil
	}
	if s, ok := c.primary.(skipper); ok && s.Skips(pane) {
		return withStateReason(fix, "jev: no terminal text to observe; not asked"), nil
	}
	ask, trial := c.admit()
	if !ask {
		return withStateReason(fix, circuitOpen), nil
	}
	key := cacheKey{kind: kindOf(profile), hash: sha256.Sum256([]byte(pane))}
	if !trial {
		if hit, ok := c.lookupState(key); ok {
			out, _ := decideState(fix, hit.reading, hit.err, c.threshold)
			return out, nil
		}
	}
	start := c.now()
	answer, perr := pr.ObserveState(ctx, profile, pane)
	latency := c.now().Sub(start)
	event := c.record(ctx, perr)
	if ctx.Err() == nil {
		c.storeState(key, cachedState{reading: answer, err: perr, at: c.now()})
	}
	out, fallback := decideState(fix, answer, perr, c.threshold)
	hash := fmt.Sprintf("%x", key.hash[:6])
	if c.appendLine != nil {
		_ = c.appendLine(LogLine{Time: start, Kind: key.kind, Hash: hash, Latency: latency, Source: out.Source,
			Confidence: answer.Confidence, Caller: screen.CallerOf(ctx), Used: used(fallback), Fallback: fallback,
			State: answer.State}.String())
		if event != "" {
			_ = c.appendLine(BreakerLine(c.now(), event))
		}
	}
	if c.screensDir != "" && perr == nil {
		_ = recordScreen(c.screensDir, screenRecord{Kind: key.kind, Hash: hash, Caller: screen.CallerOf(ctx),
			Stub: string(fix.Startup), Jev: answer.State, Confidence: answer.Confidence, At: start, Text: jev.Prepare(pane, "")})
	}
	return out, nil
}

// decideState is Jev's reading when it is sure, else the fixture's
// PaneUnknown with the reason, and the fallback rule that applied.
func decideState(fix, answer screen.StateReading, perr error, threshold float64) (screen.StateReading, string) {
	switch {
	case perr != nil:
		return withStateReason(fix, "jev: "+perr.Error()), FallbackError
	case answer.Confidence < threshold:
		return withStateReason(fix, fmt.Sprintf("jev: confidence %.2f below %.2f (%s)", answer.Confidence, threshold, answer.Reason)), FallbackThreshold
	}
	answer.Startup = fix.Startup
	return answer, ""
}

func withStateReason(fix screen.StateReading, why string) screen.StateReading {
	if fix.Reason != "" {
		why += "; " + fix.Reason
	}
	fix.Reason = why
	return fix
}

func (c *Chain) lookupState(key cacheKey) (cachedState, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	hit, ok := c.states[key]
	if !ok || c.now().Sub(hit.at) >= TTL {
		return cachedState{}, false
	}
	return hit, true
}

func (c *Chain) storeState(key cacheKey, entry cachedState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.states {
		if entry.at.Sub(v.at) >= TTL {
			delete(c.states, k)
		}
	}
	c.states[key] = entry
}
