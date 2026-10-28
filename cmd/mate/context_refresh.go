package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/memory"
	"github.com/nguyenngocanh94/mate/internal/outbox"
	"github.com/nguyenngocanh94/mate/internal/send"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

const defaultRefreshContext int64 = 150000
const refreshRetryAfter = 30 * time.Minute

// contextRefresh changes sessions only at a verified checkpoint. Queued work
// remains in the durable outbox. Ordinary restart retains its resume semantics.
func contextRefresh(ctx context.Context, w *store.Workspace, deps spawn.Deps, project string, automatic bool) (bool, error) {
	now := time.Now()
	if deps.Now != nil {
		now = deps.Now()
	}
	cfg, err := w.LoadProject(project)
	if err != nil {
		return false, err
	}
	meta, err := w.ReadMateMeta(project)
	if err != nil {
		return false, err
	}
	session := meta[spawn.MetaSessionID]
	// Automatic refresh judges the Mate quiet by the answers its Stop hook
	// logs to sent.log (refreshQuiet). A harness whose turn-end evidence is
	// not that log - Codex's is its rollout - or that declares none gives
	// it nothing to judge by, so it is never refreshed automatically; the
	// captain's `mate mate refresh` still can.
	if automatic && (cfg.Mate.RefreshContext < 0 || !logsAnswers(deps.Harnesses, meta[spawn.MetaHarness]) || !w.Auto(project) || w.Held(project)) {
		return false, nil
	}
	if session == "" {
		return false, nil
	}
	if automatic {
		quiet, err := refreshQuiet(w, project, now)
		if err != nil || !quiet {
			return false, err
		}
		handle, err := db.OpenRead(w)
		if err != nil {
			return false, nil
		}
		usage, known, err := handle.LatestContext(ctx, timeline.MateActorID(project), session)
		handle.Close()
		if err != nil {
			return false, err
		}
		limit := cfg.Mate.RefreshContext
		if limit == 0 {
			limit = defaultRefreshContext
		}
		if !known || usage.Tokens < limit {
			return false, nil
		}
	}
	release, locked, err := w.TryMateMaintenance(project, true)
	if err != nil || !locked {
		return false, err
	}
	defer release()
	if err := refreshQueueEmpty(w, project); err != nil {
		if automatic {
			return false, nil
		}
		return false, err
	}
	requestPath, err := refreshMeta(w, project, "refresh-request")
	if err != nil {
		return false, err
	}
	prior, err := store.ReadMeta(requestPath)
	if err != nil {
		return false, err
	}
	if automatic {
		if at, err := time.Parse(time.RFC3339Nano, prior["attempted_at"]); err == nil && now.Sub(at) < refreshRetryAfter {
			return false, nil
		}
	}
	// Busy/pending is a deferral, not a failed stow: do not spend the retry
	// cooldown before an LLM request was even possible. Stow rechecks this.
	composer, err := mateComposer(ctx, w, deps, project)
	if err != nil {
		if automatic {
			return false, nil
		}
		return false, err
	}
	if composer.State != send.StateEmpty {
		if automatic {
			return false, nil
		}
		return false, fmt.Errorf("refresh postponed: Mate composer is %s", send.ComposerLabel(composer.State))
	}
	nonce := rand.Text()
	req := map[string]string{"nonce": nonce, "session": session, "attempted_at": now.UTC().Format(time.RFC3339Nano)}
	if err := store.WriteMeta(requestPath, req); err != nil {
		return false, err
	}
	// The fresh start must be able to come up before the Mate is asked to
	// stow, let alone stopped: plugins that do not resolve for this project,
	// or MCP servers its harness cannot be given, would leave it with no
	// Mate at all. Asked after the request is written, so an automatic
	// refresh waits out refreshRetryAfter rather than retrying every tick.
	preflight := func() error {
		if err := spawn.PreflightMatePlugins(w, deps, project, harness.Kind(meta[spawn.MetaHarness])); err != nil {
			return fmt.Errorf("refresh postponed: a fresh Mate could not start: %w", err)
		}
		return nil
	}
	if err := preflight(); err != nil {
		return false, err
	}
	// Where sent.log ends now, so the check after the stow sees only what
	// was said during it. Writers append whole lines under the lock, so the
	// size is the offset a full read would end at.
	offset, err := w.SentLogSize(project)
	if err != nil {
		return false, err
	}
	binary := deps.Binary
	if binary == "" {
		binary, err = os.Executable()
		if err != nil {
			return false, err
		}
	}
	line := memory.StowLine + "; before ending, preserve unresolved questions verbatim, promises, decisions, active crews and next actions, then run " + shellWord(binary) + " checkpoint " + shellWord(project) + " " + nonce
	stowed, err := consoleOutbox(w, deps).Stow(ctx, project, outbox.StowOptions{RequireEmpty: true, RequireCompletion: true, Text: line})
	if err != nil {
		return false, err
	}
	if !stowed.Stowed {
		return false, fmt.Errorf("refresh postponed: %s", stowed.Outcome())
	}
	if err := verifyCheckpoint(w, project, nonce, session); err != nil {
		return false, err
	}
	// A captain prompt during stow takes precedence, even if it has already
	// been answered. Do not discard a conversation that changed since stow.
	entries, _, err := w.ReadSent(project, offset)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.Source == store.SourceUser && e.Target == store.TargetMate {
			return false, fmt.Errorf("refresh postponed: captain spoke during stow")
		}
	}
	if err := refreshQueueEmpty(w, project); err != nil {
		return false, err
	}
	current, err := w.ReadMateMeta(project)
	if err != nil {
		return false, err
	}
	if current[spawn.MetaSessionID] != session {
		return false, fmt.Errorf("refresh postponed: Mate session changed")
	}
	if automatic && (!w.Auto(project) || w.Held(project)) {
		return false, fmt.Errorf("refresh postponed: captain took manual control")
	}
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	// Again, since the stow took a turn in which plugins.yaml could change.
	if err := preflight(); err != nil {
		return false, err
	}
	if err := w.ArchiveMateSession(project, current); err != nil {
		return false, err
	}
	// A headless Mate changes session in its own pane: stopping it closed
	// the pane, and every attach to it - the console's stage among them -
	// ended with "terminal ... not found" (CON-133).
	switch _, err := spawn.ReloadMate(ctx, w, deps, project); {
	case err == nil:
		// The host switched between turns, so no turn writes the old
		// transcript any more.
		if err := w.FreezeMateSession(project, current); err != nil {
			return false, fmt.Errorf("Mate refreshed with checkpoint saved; archive failed: %w", err)
		}
		err = w.AppendSent(project, store.SentEntry{Source: store.SourceApp, Target: store.TargetMate, Text: "context refreshed from checkpoint; previous session " + session})
		return true, err
	case errors.As(err, new(*matehost.BusyError)):
		return false, fmt.Errorf("refresh postponed: %w", err)
	case !errors.Is(err, spawn.ErrNoReload):
		return false, err
	}
	stopped, err := spawn.StopMate(ctx, w, deps, project)
	if err != nil {
		return false, err
	}
	// Only a confirmed stop freezes the transcript. A Mate whose Herdr
	// session was not reported running was cleared on inference alone, so
	// its archive stays as written above, on the live transcript and not
	// finalized: the timeline keeps its last message group pending rather
	// than flushing a turn a writer may still extend.
	if stopped.Confirmed {
		if err := w.FreezeMateSession(project, current); err != nil {
			return false, fmt.Errorf("Mate stopped with checkpoint saved; archive failed: %w", err)
		}
	}
	_, err = spawn.StartMate(ctx, w, deps, spawn.StartRequest{Project: project, Harness: harness.Kind(meta[spawn.MetaHarness]), Fresh: true})
	if err != nil {
		return false, fmt.Errorf("checkpoint saved; fresh start failed: %w", err)
	}
	err = w.AppendSent(project, store.SentEntry{Source: store.SourceApp, Target: store.TargetMate, Text: "context refreshed from checkpoint; previous session " + session})
	return true, err
}

// logsAnswers reports whether the recorded harness declares verified
// turn-end evidence that mate's Stop hook logs each answer to sent.log.
func logsAnswers(reg harness.Registry, recorded string) bool {
	kind, err := reg.Parse(recorded)
	if err != nil {
		return false
	}
	profile, err := reg.Lookup(kind)
	if err != nil {
		return false
	}
	c := profile.Capabilities().TurnEnd
	return c.Verified() && c.Impl.LogsAnswers()
}

func shellWord(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'" }

func refreshQueueEmpty(w *store.Workspace, project string) error {
	items, err := w.ReadOutbox(project)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.Queued() {
			return fmt.Errorf("refresh postponed: %s item %s is queued", item.Source, strconv.FormatInt(item.ID, 10))
		}
	}
	return nil
}

func refreshQuiet(w *store.Workspace, project string, now time.Time) (bool, error) {
	entries, _, err := w.ReadSent(project, 0)
	if err != nil {
		return false, err
	}
	var asked, answered time.Time
	for _, e := range entries {
		if e.Source == store.SourceUser && e.Target == store.TargetMate {
			asked = e.Time
		}
		if e.Source == store.SourceMate && e.Target == store.SourceUser {
			answered = e.Time
		}
	}
	return !answered.IsZero() && !answered.Before(asked) && now.Sub(answered) >= store.QuietAfter, nil
}

// mateComposer is the Mate's composer: read from its pane, or, for a
// headless Mate (CON-78), as its host states it.
func mateComposer(ctx context.Context, w *store.Workspace, deps spawn.Deps, project string) (send.Classification, error) {
	if host, ok := spawn.HostClient(w, project); ok {
		st, err := host.State(ctx)
		if err != nil {
			return send.Classification{}, err
		}
		switch {
		case st.Draft != "":
			return send.Classification{State: send.StatePending, Pending: st.Draft}, nil
		case st.Busy || st.Queued > 0:
			return send.Classification{State: send.StateBusy}, nil
		}
		return send.Classification{State: send.StateEmpty}, nil
	}
	h, kind, err := spawn.MateHandle(ctx, w, deps, project)
	if err != nil {
		return send.Classification{}, err
	}
	screens, err := screenOf(deps, kind)
	if err != nil {
		return send.Classification{}, err
	}
	screen, err := deps.Runtime.ReadAgentStyled(ctx, h, screens.ReadSource(), send.DefaultLines)
	if err != nil {
		return send.Classification{}, err
	}
	return send.ClassifyComposer(screens, screen), nil
}
