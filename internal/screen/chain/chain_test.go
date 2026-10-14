package chain_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/claude"
	"github.com/nguyenngocanh94/mate/internal/harness/codex"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/screen/chain"
	"github.com/nguyenngocanh94/mate/internal/screen/fixture"
)

// stub answers every screen with one Observation or one error, and counts.
type stub struct {
	obs   screen.Observation
	err   error
	calls int
}

func (s *stub) Observe(context.Context, harness.ScreenProfile, string) (screen.Observation, error) {
	s.calls++
	return s.obs, s.err
}

var (
	claudeScreen = claude.Claude{}.Screen()
	codexScreen  = codex.Codex{}.Screen()
)

// fixtureDraft is a fixture reading that is sure of a draft.
var fixtureDraft = screen.Observation{Composer: screen.ComposerDraft, Deterministic: screen.ComposerDraft, Draft: "half typed", Evidence: "❯ half typed",
	Dialog: screen.DialogNone, Startup: harness.StartupScreenUnrecognized, Highlight: -1, Confidence: 1, Source: "fixture"}

// fixtureUnread is a fixture reading of a screen no stub names: the only
// kind the chain asks Jev about.
var fixtureUnread = screen.Observation{Composer: screen.ComposerUnknown, Deterministic: screen.ComposerUnknown, Dialog: screen.DialogUnknown,
	Startup: harness.StartupScreenUnrecognized, Highlight: -1, Confidence: 0, Source: "fixture", Reason: "nothing recognised"}

func jevSays(composer screen.ComposerState, dialog screen.DialogKind, conf float64) screen.Observation {
	return screen.Observation{Composer: composer, Dialog: dialog, Highlight: -1, Confidence: conf, Source: "jev",
		Notice: "quota_warning", Reason: "composer and dialog answers"}
}

func TestChainRules(t *testing.T) {
	fixtureTrust := screen.Observation{Composer: screen.ComposerUnknown, Deterministic: screen.ComposerUnknown, Evidence: "", Dialog: screen.DialogTrust,
		Startup: harness.StartupScreenTrustDialog, Highlight: 1, Confidence: 1, Source: "fixture"}
	tests := []struct {
		name       string
		fixture    screen.Observation
		jev        screen.Observation
		jevErr     error
		asked      bool
		want       screen.Observation
		wantReason string
	}{
		{
			name: "the fixture reads a draft: its reading whole, jev not asked", fixture: fixtureDraft,
			jev: jevSays(screen.ComposerEmpty, screen.DialogNone, 0.99), want: fixtureDraft,
		},
		{
			name: "the fixture reads an empty composer: jev not asked, so it cannot veto", fixture: screen.Observation{
				Composer: screen.ComposerEmpty, Deterministic: screen.ComposerEmpty, Evidence: "›", Dialog: screen.DialogNone,
				Startup: harness.StartupScreenReady, Highlight: -1, Confidence: 1, Source: "fixture"},
			jev: jevSays(screen.ComposerBusy, screen.DialogNone, 0.99),
			want: screen.Observation{Composer: screen.ComposerEmpty, Deterministic: screen.ComposerEmpty, Evidence: "›", Dialog: screen.DialogNone,
				Startup: harness.StartupScreenReady, Highlight: -1, Confidence: 1, Source: "fixture"},
		},
		{
			name: "the fixture recognises a dialog: its reading whole, jev not asked", fixture: fixtureTrust,
			jev: jevSays(screen.ComposerEmpty, screen.DialogNone, 0.99), want: fixtureTrust,
		},
		{
			name: "jev fails: the fixture's reading whole", fixture: fixtureUnread, jevErr: errors.New("Jev unavailable (HTTP 503)"),
			asked: true, want: fixtureUnread, wantReason: "jev: Jev unavailable (HTTP 503); nothing recognised",
		},
		{
			name: "jev below the threshold: the fixture's reading whole", fixture: fixtureUnread,
			jev: jevSays(screen.ComposerEmpty, screen.DialogNone, 0.84), asked: true, want: fixtureUnread,
			wantReason: "jev: confidence 0.84 below 0.85 (composer and dialog answers); nothing recognised",
		},
		{
			name: "jev calls a screen the fixture cannot read empty: the composer stays unknown", fixture: fixtureUnread,
			jev: jevSays(screen.ComposerEmpty, screen.DialogNone, 0.99), asked: true,
			want: screen.Observation{Composer: screen.ComposerUnknown, Deterministic: screen.ComposerUnknown, Dialog: screen.DialogNone,
				Notice: "quota_warning", Startup: harness.StartupScreenUnrecognized, Highlight: -1, Confidence: 0.99, Source: "jev"},
			wantReason: "composer and dialog answers",
		},
		{
			name: "jev names a dialog the fixture does not recognise: no highlight", fixture: fixtureUnread,
			jev: jevSays(screen.ComposerUnknown, screen.DialogHooksReview, 0.95), asked: true,
			want: screen.Observation{Composer: screen.ComposerUnknown, Deterministic: screen.ComposerUnknown, Dialog: screen.DialogHooksReview,
				Notice: "quota_warning", Startup: harness.StartupScreenUnrecognized, Highlight: -1, Confidence: 0.95, Source: "jev"},
			wantReason: "composer and dialog answers",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jev := &stub{obs: tc.jev, err: tc.jevErr}
			c := chain.New(jev, &stub{obs: tc.fixture}, 0.85)
			got, err := c.Observe(context.Background(), claudeScreen, "screen")
			if err != nil {
				t.Fatal(err)
			}
			if asked := jev.calls == 1; asked != tc.asked {
				t.Errorf("jev asked %d times, want asked=%v", jev.calls, tc.asked)
			}
			if tc.wantReason == "" {
				tc.wantReason = tc.want.Reason
			}
			if got.Reason != tc.wantReason {
				t.Errorf("reason %q, want %q", got.Reason, tc.wantReason)
			}
			got.Reason, tc.want.Reason = "", ""
			if got != tc.want {
				t.Errorf("got  %+v\nwant %+v", got, tc.want)
			}
		})
	}
}

// A send never asks Jev, whatever the fixture reads: it types only into a
// composer the fixture reads empty, and Jev cannot make one empty. A screen
// the fixture cannot read is still kept for the captain, once a minute.
func TestChainNeverAsksJevForASend(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 10, 14, 9, 0, 0, 0, time.UTC)
	jev := &stub{obs: jevSays(screen.ComposerDraft, screen.DialogPermission, 0.99)}
	var lines []string
	c := chain.New(jev, &stub{obs: fixtureUnread}, 0.85, chain.WithClock(func() time.Time { return now }),
		chain.WithScreens(dir), chain.WithLog(func(l string) error { lines = append(lines, l); return nil }))
	ctx := screen.WithCaller(context.Background(), screen.CallerSend)
	for range 3 {
		got, err := c.Observe(ctx, claudeScreen, "❯ ok, từ nay xoá luôn\n  em để cho anh 2 news:")
		if err != nil {
			t.Fatal(err)
		}
		if got.Composer != screen.ComposerUnknown || got.Source != "fixture" || !strings.HasPrefix(got.Reason, "jev: not asked for a send") {
			t.Fatalf("got %+v, want the fixture's unknown, not asked", got)
		}
		now = now.Add(2 * time.Second)
	}
	if jev.calls != 0 || len(lines) != 0 {
		t.Fatalf("%d requests and log %q for a send", jev.calls, lines)
	}
	files, _ := filepath.Glob(filepath.Join(dir, "claude", "*.txt"))
	if len(files) != 1 {
		t.Fatalf("files = %v, want the unread screen kept", files)
	}
	b, _ := os.ReadFile(files[0])
	for _, want := range []string{"# seen: 1\n", "# caller: send\n", "# stub: unrecognized\n", "# jev: not asked\n", "em để cho anh 2 news:"} {
		if !strings.Contains(string(b), want) {
			t.Errorf("record lacks %q:\n%s", want, b)
		}
	}
}

// One screen of one harness is asked of Jev once a minute however often it
// is observed; another screen, another harness or an older answer asks
// again. Every request, and only a request, is logged.
func TestChainAsksJevOncePerScreenAMinute(t *testing.T) {
	now := time.Date(2026, 10, 8, 9, 0, 0, 0, time.UTC)
	jev := &stub{obs: jevSays(screen.ComposerEmpty, screen.DialogNone, 0.9)}
	var lines []string
	c := chain.New(jev, &stub{obs: fixtureUnread}, 0.85,
		chain.WithClock(func() time.Time { return now }),
		chain.WithLog(func(line string) error { lines = append(lines, line); return nil }))
	observe := func(profile harness.ScreenProfile, pane string) screen.Observation {
		t.Helper()
		obs, err := c.Observe(context.Background(), profile, pane)
		if err != nil {
			t.Fatal(err)
		}
		return obs
	}
	steps := []struct {
		advance time.Duration
		profile harness.ScreenProfile
		pane    string
		calls   int
	}{
		{0, claudeScreen, "a", 1},
		{59 * time.Second, claudeScreen, "a", 1},
		{0, claudeScreen, "b", 2},
		{0, codexScreen, "a", 3},
		{time.Second, claudeScreen, "a", 4},
		{0, claudeScreen, "a", 4},
	}
	for i, s := range steps {
		now = now.Add(s.advance)
		if obs := observe(s.profile, s.pane); obs.Source != "jev" {
			t.Fatalf("step %d: %+v, want jev's reading, cached or not", i, obs)
		}
		if jev.calls != s.calls || len(lines) != s.calls {
			t.Fatalf("step %d: %d calls, %d log lines, want %d", i, jev.calls, len(lines), s.calls)
		}
	}
	if !strings.Contains(lines[2], " codex ") || !strings.Contains(lines[0], " claude ") {
		t.Fatalf("log lines do not name the harness: %q", lines)
	}
}

// A failed request is remembered like an answer, so a pane Jev cannot read
// is not asked again every poll; a request the caller cancelled is not.
func TestChainRemembersFailuresButNotCancellations(t *testing.T) {
	jev := &stub{err: errors.New("offline")}
	c := chain.New(jev, &stub{obs: fixtureUnread}, 0.85)
	for range 2 {
		if obs, err := c.Observe(context.Background(), claudeScreen, "a"); err != nil || obs.Source != "fixture" {
			t.Fatalf("%+v, %v", obs, err)
		}
	}
	if jev.calls != 1 {
		t.Fatalf("%d calls for one failing screen", jev.calls)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for range 2 {
		if _, err := c.Observe(ctx, claudeScreen, "b"); err != nil {
			t.Fatal(err)
		}
	}
	if jev.calls != 3 {
		t.Fatalf("%d calls, want a cancelled request asked again", jev.calls)
	}
}

// A fixture that fails is the chain's failure: there is nothing to fall
// back to.
func TestChainFailsWhenTheFixtureDoes(t *testing.T) {
	c := chain.New(&stub{obs: jevSays(screen.ComposerEmpty, screen.DialogNone, 1)}, &stub{err: errors.New("broken")}, 0.85)
	if _, err := c.Observe(context.Background(), claudeScreen, "a"); err == nil {
		t.Fatal("no error")
	}
}

// On a Codex trust dialog the fixture recognises, Jev is not asked at all:
// the dialog, its highlight and the composer are the fixture's, so send
// refuses and settle answers the dialog as measured, at no request.
func TestChainKeepsADialogTheFixtureRecognises(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(moduleRoot, "internal", "harness", "codex", "testdata", "startup", "codex-0.156.1-trust-dialog.txt"))
	if err != nil {
		t.Fatal(err)
	}
	jev := &stub{obs: screen.Observation{Composer: screen.ComposerEmpty, Dialog: screen.DialogNone, Highlight: -1,
		Confidence: 0.99, Source: "jev", Reason: "composer empty 0.99, dialog none 0.99"}}
	var lines []string
	c := chain.New(jev, fixture.New(), 0.85, chain.WithLog(func(l string) error { lines = append(lines, l); return nil }))
	got, err := c.Observe(context.Background(), codexScreen, string(data))
	if err != nil {
		t.Fatal(err)
	}
	if got.Dialog != screen.DialogTrust || got.Highlight != 0 || got.Startup != harness.StartupScreenTrustDialog ||
		got.Composer != screen.ComposerUnknown || got.Source != "fixture" {
		t.Fatalf("got %+v, want the fixture's trust dialog with its highlight on option 0", got)
	}
	if jev.calls != 0 || len(lines) != 0 {
		t.Fatalf("%d requests, log %q for a dialog the fixture recognises", jev.calls, lines)
	}
}

// skipsBlank is a failing primary that skips a blank screen, as the jev
// observer does.
type skipsBlank struct{ stub }

func (s *skipsBlank) Skips(pane string) bool { return strings.TrimSpace(pane) == "" }

// A screen the primary skips is the fixture's reading, with no request, no
// log line and no count toward the breaker: blank screens between failures
// do not open the circuit, and a failure after them is still only the
// first in a row.
func TestChainSkipsAScreenThePrimarySkips(t *testing.T) {
	jev := &skipsBlank{stub{err: errors.New("offline")}}
	var lines []string
	c := chain.New(jev, &stub{obs: fixtureUnread}, 0.85, chain.WithLog(func(l string) error { lines = append(lines, l); return nil }))
	for i := range 5 {
		obs, err := c.Observe(context.Background(), claudeScreen, strings.Repeat(" ", i+1))
		if err != nil || obs.Source != "fixture" || !strings.HasPrefix(obs.Reason, "jev: no terminal text to observe; not asked") {
			t.Fatalf("blank screen %d: %+v, %v", i, obs, err)
		}
	}
	if jev.calls != 0 || len(lines) != 0 {
		t.Fatalf("%d requests, log %q for blank screens", jev.calls, lines)
	}
	for _, pane := range []string{"a", "b"} {
		if _, err := c.Observe(context.Background(), claudeScreen, pane); err != nil {
			t.Fatal(err)
		}
	}
	if jev.calls != 2 || len(lines) != 2 || strings.Contains(strings.Join(lines, "\n"), "breaker open") {
		t.Fatalf("%d requests, log %q: want two failed requests and a closed circuit", jev.calls, lines)
	}
}
