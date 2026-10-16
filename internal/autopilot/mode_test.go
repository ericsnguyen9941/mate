package autopilot_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/autopilot"
	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/hook"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// A prompt to the Mate does not turn auto off (docs/mvp.md M18), so a crew
// that reports while the captain is talking is still told to the Mate.
// Manual is the captain's `m` key, and it lasts while the captain is with
// the Mate: five minutes without a prompt to it or a press of `m` turn auto
// back on (docs/mvp.md M21).

// captainPrompts runs the Mate's UserPromptSubmit hook on a plain prompt, the
// way the harness does when the captain types to the Mate.
func (f *fixture) captainPrompts(text string) {
	f.t.Helper()
	if err := hook.HandlePrompt(f.ws, project, []byte(`{"prompt":"`+text+`"}`)); err != nil {
		f.t.Fatal(err)
	}
}

// TestACaptainPromptThenWaitMateStillWakesTheMate replays 2026-10-07 09:50:
// the captain messaged the Mate, the crew said wait-mate 24 seconds later,
// and for more than five minutes nobody was told. In auto mode the digest
// reaches the Mate.
func TestACaptainPromptThenWaitMateStillWakesTheMate(t *testing.T) {
	f := newFixture(t)
	f.auto(true)

	f.captainPrompts("how is it going")
	if !f.ws.Auto(project) {
		t.Fatal("a captain prompt turned auto mode off")
	}
	f.clock.Advance(24 * time.Second)
	f.status("fix-cart-total", "wait-mate: branch ready")
	f.clock.Advance(time.Minute)
	f.mustTick()

	typed := f.typed()
	if len(typed) != 1 || !strings.Contains(typed[0], "fix-cart-total wait-mate") {
		t.Fatalf("typed %#v, want one digest carrying the crew's wait-mate", typed)
	}
}

// TestManualWaitMateReachesTheCaptainsInbox is the other half: with manual
// held by the captain nothing is typed, and the inbox carries the crew's
// wait-mate, so one side is always told.
func TestManualWaitMateReachesTheCaptainsInbox(t *testing.T) {
	f := newFixture(t)
	if err := f.ws.SetMode(project, false); err != nil {
		t.Fatal(err)
	}

	f.captainPrompts("how is it going")
	f.clock.Advance(24 * time.Second)
	f.status("fix-cart-total", "wait-mate: branch ready")
	f.clock.Advance(4 * time.Minute)
	f.mustTick()

	if typed := f.typed(); len(typed) != 0 {
		t.Fatalf("typed %#v into a manual project", typed)
	}
	items := box.Inbox(f.view())
	if len(items) != 1 || items[0].Crew() != "fix-cart-total" || items[0].State != box.StateWaitMate {
		t.Fatalf("inbox = %+v, want the crew's wait-mate", items)
	}
}

// holdManual is the captain pressing `m` to manual now, by the fixture's
// clock: the hold's mtime is when they chose it.
func (f *fixture) holdManual() {
	f.t.Helper()
	if err := f.ws.SetMode(project, false); err != nil {
		f.t.Fatal(err)
	}
	now := f.clock.Now()
	if err := os.Chtimes(f.ws.ManualHold(project), now, now); err != nil {
		f.t.Fatal(err)
	}
}

// asked is the captain typing to the Mate at a given moment, the line the
// Mate's UserPromptSubmit hook appends.
func (f *fixture) asked(at time.Time) {
	f.t.Helper()
	if err := f.ws.AppendSent(project, store.SentEntry{Time: at, Source: store.SourceUser, Target: store.TargetMate, Text: "how is it going"}); err != nil {
		f.t.Fatal(err)
	}
}

func hasSent(entries []store.SentEntry, text string) bool {
	for _, e := range entries {
		if e.Text == text {
			return true
		}
	}
	return false
}

// TestManualComesBackToAutoAfterFiveQuietMinutes replays 2026-10-13: the
// captain pressed `m` at 11:51 and left, three crews asked for a decision at
// 20:40, and nobody answered them until the captain came back at 21:25.
func TestManualComesBackToAutoAfterFiveQuietMinutes(t *testing.T) {
	f := newFixture(t)
	f.holdManual()
	f.status("fix-cart-total", "needs-decision: pick A or B")

	f.clock.Advance(autopilot.DefaultQuietAfter - time.Second)
	f.mustTick()
	if f.ws.Auto(project) || !f.ws.Held(project) || len(f.typed()) != 0 {
		t.Fatalf("auto = %v, held = %v, typed %#v under five minutes after the captain chose manual", f.ws.Auto(project), f.ws.Held(project), f.typed())
	}

	f.clock.Advance(time.Second)
	f.mustTick()
	if !f.ws.Auto(project) || f.ws.Held(project) {
		t.Fatalf("auto = %v, held = %v five minutes after the captain left", f.ws.Auto(project), f.ws.Held(project))
	}
	typed := f.typed()
	if len(typed) != 1 || !strings.Contains(typed[0], "fix-cart-total needs-decision") {
		t.Fatalf("typed %#v, want the digest the crew's question is in", typed)
	}
	// CON-93: auto came back without the captain, so the digest that hands
	// the Mate the question says a captain's decision is not its to make.
	if !strings.Contains(typed[0], "a decides: captain choice goes to the captain and waits under Held") {
		t.Fatalf("digest %q does not keep a captain's decision with the captain", typed[0])
	}
	if !hasSent(f.sent(), autopilot.RearmText(autopilot.DefaultQuietAfter)) {
		t.Fatalf("sent.log = %+v, want the switch back to auto recorded", f.sent())
	}
}

// The captain typing to the Mate is the captain being there: it restarts
// the five minutes, whether or not the Mate has answered yet.
func TestACaptainPromptRestartsTheQuietMinutes(t *testing.T) {
	f := newFixture(t)
	f.holdManual()
	f.status("fix-cart-total", "needs-decision: pick A or B")

	f.clock.Advance(4 * time.Minute)
	f.asked(f.clock.Now())
	f.mustTick()
	f.clock.Advance(4 * time.Minute)
	f.mustTick()
	if f.ws.Auto(project) || len(f.typed()) != 0 {
		t.Fatalf("auto = %v, typed %#v four minutes after the captain typed", f.ws.Auto(project), f.typed())
	}

	f.clock.Advance(time.Minute)
	f.mustTick()
	if !f.ws.Auto(project) || len(f.typed()) != 1 {
		t.Fatalf("auto = %v, typed %#v five minutes after the captain typed", f.ws.Auto(project), f.typed())
	}
}

// Choosing auto with the key releases the hold at once.
func TestChoosingAutoReleasesTheHold(t *testing.T) {
	f := newFixture(t)
	f.holdManual()
	if err := f.ws.SetMode(project, true); err != nil {
		t.Fatal(err)
	}
	if f.ws.Held(project) || !f.ws.Auto(project) {
		t.Fatalf("held = %v, auto = %v after choosing auto", f.ws.Held(project), f.ws.Auto(project))
	}
}
