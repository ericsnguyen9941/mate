package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/query"
)

// actionSnapshot is the console's snapshot for the fixture project: a
// running Mate and the fixture crew, handed back, with the availability the
// query layer would publish for them.
func actionSnapshot() query.Snapshot {
	avail := func(pairs ...any) []query.ActionAvailability {
		var out []query.ActionAvailability
		for i := 0; i < len(pairs); i += 2 {
			out = append(out, query.ActionAvailability{Action: pairs[i].(string), Available: pairs[i+1].(bool), Reason: "recorded " + pairs[i].(string)})
		}
		return out
	}
	return query.Snapshot{Projects: []query.ProjectNode{{
		ProjectID: fixtureProject, Name: fixtureProject, Mode: query.ModeManual,
		Mate: query.MateNode{Actions: avail("start", false, "stop", true, "resume", false, "onboard", false,
			"restart_mate", true, "clear_composer", true)},
		Crews: []query.CrewNode{{
			CrewID: fixtureCrew, Status: query.CrewWaitMate,
			Worktree: query.KnownField(query.WorktreeValue{Branch: fixtureBranch}),
			Repo:     query.KnownField(query.RepoValue{DefaultBranch: "main"}),
			Actions:  avail("stop", true, "repair", false, "diff", true, "restart_crew", true, "merge", true),
		}},
	}}}
}

type actRecorder struct {
	mu  sync.Mutex
	got []ActionRequest
	err error
	// entered is told when an action starts; block holds it until closed.
	entered chan struct{}
	block   chan struct{}
}

func (a *actRecorder) act(_ context.Context, req ActionRequest) (string, error) {
	if a.entered != nil {
		a.entered <- struct{}{}
	}
	if a.block != nil {
		<-a.block
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.got = append(a.got, req)
	if a.err != nil {
		return "", a.err
	}
	return "did " + req.Action, nil
}

func actingServer(t *testing.T, f *fixture, rec *actRecorder) *httptest.Server {
	t.Helper()
	snap := actionSnapshot()
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

func post(t *testing.T, ts *httptest.Server, body string, header map[string]string) (int, ErrorResponse, ActionResult) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/actions", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatal(err)
	}
	var e ErrorResponse
	var ok ActionResult
	_ = json.Unmarshal(raw, &e)
	_ = json.Unmarshal(raw, &ok)
	return resp.StatusCode, e, ok
}

func TestActionsAreTheConsoleSheetsWithItsAvailability(t *testing.T) {
	f := newFixture(t)
	ts := actingServer(t, f, &actRecorder{})
	resp, err := ts.Client().Get(ts.URL + "/api/projects/" + fixtureProject + "/actions")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var got ActionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.ReadOnly != "" {
		t.Fatalf("read_only = %q on the console's dashboard", got.ReadOnly)
	}
	var mate []string
	for _, b := range got.Mate {
		mate = append(mate, b.Action)
	}
	if strings.Join(mate, " ") != "start stop restart_mate mode clear_composer" {
		t.Fatalf("mate buttons = %v", mate)
	}
	if start := got.Mate[0]; start.Available || start.Reason != "recorded start" {
		t.Errorf("start = %+v, want refused with the query layer's reason", start)
	}
	if mode := got.Mate[3]; mode.Label != "Mode → auto" || !mode.Available {
		t.Errorf("mode = %+v, want the flip to auto offered on a manual project", mode)
	}
	if got.Mate[1].Confirm == "" || got.Mate[4].Confirm != "" {
		t.Errorf("stop must ask first and clear composer must not: %+v %+v", got.Mate[1], got.Mate[4])
	}
	crew := got.Crews[fixtureCrew]
	if len(crew) != 4 {
		t.Fatalf("crew buttons = %+v, want stop, restart, merge, diff", crew)
	}
	merge := crew[2]
	if merge.Action != "merge" || !merge.Available || !strings.Contains(merge.Confirm, "Merge "+fixtureBranch+" into main") {
		t.Errorf("merge = %+v, want it offered and naming both branches", merge)
	}
}

func TestADashboardWithoutTheConsoleIsReadOnly(t *testing.T) {
	f := newFixture(t)
	var got ActionsResponse
	f.get(t, "/api/projects/"+fixtureProject+"/actions", http.StatusOK, &got)
	if got.ReadOnly == "" || len(got.Mate) != 0 {
		t.Fatalf("actions = %+v, want read-only and no buttons", got)
	}
	resp, err := f.http.Client().Post(f.http.URL+"/api/actions", "application/json", strings.NewReader(`{"action":"stop","project":"`+fixtureProject+`"}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("POST on a read-only dashboard = %d, want 403", resp.StatusCode)
	}
}

func TestAPressRunsTheConsolesActionAndAnswersItsLine(t *testing.T) {
	f := newFixture(t)
	rec := &actRecorder{}
	ts := actingServer(t, f, rec)
	status, e, ok := post(t, ts, `{"action":"merge","project":"`+fixtureProject+`","crew":"`+fixtureCrew+`"}`, nil)
	if status != http.StatusOK || ok.Line != "did merge" {
		t.Fatalf("merge = %d %q %+v", status, ok.Line, e)
	}
	want := ActionRequest{Action: "merge", Project: fixtureProject, Crew: fixtureCrew}
	if len(rec.got) != 1 || rec.got[0] != want {
		t.Fatalf("act got %+v, want %+v", rec.got, want)
	}

	rec.err = errors.New("send: the composer is not empty")
	status, e, _ = post(t, ts, `{"action":"clear_composer","project":"`+fixtureProject+`"}`, nil)
	if status != http.StatusUnprocessableEntity || e.Error != "send: the composer is not empty" {
		t.Fatalf("a failed action = %d %+v, want 422 with the service's own words", status, e)
	}
}

func TestAPressIsRefusedBeforeAnythingRuns(t *testing.T) {
	f := newFixture(t)
	rec := &actRecorder{}
	ts := actingServer(t, f, rec)
	for _, tc := range []struct {
		name   string
		body   string
		header map[string]string
		status int
	}{
		{"unavailable", `{"action":"start","project":"` + fixtureProject + `"}`, nil, http.StatusConflict},
		{"not offered", `{"action":"retry","project":"` + fixtureProject + `"}`, nil, http.StatusBadRequest},
		{"remove not allowed by the snapshot", `{"action":"remove_project","project":"` + fixtureProject + `"}`, nil, http.StatusConflict},
		{"repair is not wired", `{"action":"repair","project":"` + fixtureProject + `","crew":"` + fixtureCrew + `"}`, nil, http.StatusBadRequest},
		{"unknown project", `{"action":"stop","project":"nosuch"}`, nil, http.StatusNotFound},
		{"closed crew", `{"action":"stop","project":"` + fixtureProject + `","crew":"gone"}`, nil, http.StatusNotFound},
		{"another origin", `{"action":"stop","project":"` + fixtureProject + `"}`, map[string]string{"Origin": "https://evil.example"}, http.StatusForbidden},
		{"a form", `{"action":"stop","project":"` + fixtureProject + `"}`, map[string]string{"Content-Type": "text/plain"}, http.StatusUnsupportedMediaType},
	} {
		if status, e, _ := post(t, ts, tc.body, tc.header); status != tc.status {
			t.Errorf("%s: status = %d (%+v), want %d", tc.name, status, e, tc.status)
		}
	}
	if len(rec.got) != 0 {
		t.Fatalf("act ran %+v; every press above must be refused first", rec.got)
	}
	if status, _, _ := post(t, ts, `{"action":"stop","project":"`+fixtureProject+`"}`, map[string]string{"Origin": ts.URL}); status != http.StatusOK {
		t.Fatalf("a press from the page's own origin = %d, want 200", status)
	}
}

func TestASecondPressOnTheSameTargetWaitsForTheFirst(t *testing.T) {
	f := newFixture(t)
	rec := &actRecorder{entered: make(chan struct{}, 1), block: make(chan struct{})}
	ts := actingServer(t, f, rec)
	body := `{"action":"restart_mate","project":"` + fixtureProject + `"}`
	first := make(chan int)
	go func() {
		status, _, _ := post(t, ts, body, nil)
		first <- status
	}()
	<-rec.entered
	status, e, _ := post(t, ts, `{"action":"clear_composer","project":"`+fixtureProject+`"}`, nil)
	if status != http.StatusConflict || !strings.Contains(e.Error, "restart_mate is still running") {
		t.Fatalf("the second press = %d %+v, want refused while the first runs", status, e)
	}
	close(rec.block)
	if status := <-first; status != http.StatusOK {
		t.Fatalf("the first press = %d", status)
	}
}
