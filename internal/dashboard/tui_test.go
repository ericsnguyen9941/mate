package dashboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/query"
)

// tuiSnapshot is actionSnapshot while the Mate of shop has its session
// open in the TUI (CON-227).
func tuiSnapshot() query.Snapshot {
	snap := actionSnapshot()
	snap.TUI = &query.TUIStatus{Project: "shop", PID: 4242, Since: time.Date(2026, 11, 6, 9, 30, 0, 0, time.Local)}
	return snap
}

func tuiServer(t *testing.T, rec *actRecorder) *httptest.Server {
	t.Helper()
	f := newFixture(t)
	snap := tuiSnapshot()
	server, err := New(Options{Workspace: f.ws, DB: f.read, Deps: Deps{
		Now:      func() time.Time { return fixtureNow },
		Snapshot: func(context.Context) (query.Snapshot, error) { return snap, nil },
		Act:      rec.act,
	}})
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(server.Handler())
	t.Cleanup(ts.Close)
	return ts
}

// While a Mate's TUI is open, Office and City only look: every button but
// a peek is off with the Mate's name as its reason, the answer says the
// flag, and a press is a 409 with code tui that runs nothing. A peek runs.
func TestATUIOpenLeavesThePagesOnlyLooking(t *testing.T) {
	rec := &actRecorder{}
	ts := tuiServer(t, rec)
	resp, err := ts.Client().Get(ts.URL + "/api/projects/" + fixtureProject + "/actions")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var got ActionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.TUI != "Mate shop đang mở TUI (pid 4242, từ 09:30)" {
		t.Fatalf("tui = %q", got.TUI)
	}
	all := append(append(append([]Button(nil), got.Mate...), got.MateMore...), got.Crews[fixtureCrew]...)
	all = append(all, got.CrewMore[fixtureCrew]...)
	for _, b := range all {
		if b.Action == actPeek {
			continue
		}
		if b.Available || b.Reason != "Mate shop đang mở TUI" {
			t.Errorf("button %s = %+v, want it off with the TUI's reason", b.Action, b)
		}
	}

	resp, err = ts.Client().Get(ts.URL + "/api/actions")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var every AllActionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&every); err != nil {
		t.Fatal(err)
	}
	if every.TUI == "" {
		t.Fatal("GET /api/actions does not say the TUI is open")
	}

	for _, body := range []string{
		`{"action":"stop","project":"` + fixtureProject + `"}`,
		`{"action":"say","project":"` + fixtureProject + `","text":"hi"}`,
		`{"action":"merge","project":"` + fixtureProject + `","crew":"` + fixtureCrew + `"}`,
	} {
		status, e, _ := post(t, ts, body, nil)
		if status != http.StatusConflict || e.Code != coreapi.CodeTUI || !strings.Contains(e.Error, "Mate shop đang mở TUI") {
			t.Fatalf("%s with the TUI open = %d %+v, want 409 code tui", body, status, e)
		}
	}
	if len(rec.got) != 0 {
		t.Fatalf("acted %+v with the TUI open", rec.got)
	}
}

// An action the core refuses for the TUI - a new project, which no
// snapshot check guards - is a 409 with code tui too.
func TestACoreRefusalForTheTUIIsAConflictWithItsCode(t *testing.T) {
	rec := &actRecorder{}
	rec.fail(&coreapi.TUIError{TUI: query.TUIStatus{Project: "shop", PID: 4242}})
	f := newFixture(t)
	server, err := New(Options{Workspace: f.ws, DB: f.read, Deps: Deps{
		Now:      func() time.Time { return fixtureNow },
		Snapshot: func(context.Context) (query.Snapshot, error) { return actionSnapshot(), nil },
		Act:      rec.act,
	}})
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	status, e, _ := post(t, ts, `{"action":"add_project","project":"landing"}`, nil)
	if status != http.StatusConflict || e.Code != coreapi.CodeTUI {
		t.Fatalf("add_project refused for the TUI = %d %+v, want 409 code tui", status, e)
	}

	// On the core's socket the same refusal carries the flag.
	code, body := clientActError(&coreapi.TUIError{TUI: query.TUIStatus{Project: "shop", PID: 4242}})
	if code != http.StatusConflict || body.Code != coreapi.CodeTUI || body.TUI == nil || body.TUI.PID != 4242 {
		t.Fatalf("socket refusal = %d %+v", code, body)
	}
}
