package claude

import (
	"slices"
	"testing"
)

// The lines are the shapes measured on Claude Code 2.1.292
// (docs/evidence/headless-mate-probe-2026-10-15.md).
func TestParseStreamLine(t *testing.T) {
	ev, ok := ParseStreamLine([]byte(`{"type":"result","subtype":"success","result":"ok1","session_id":"899693ea","is_error":false,"num_turns":1}`))
	if !ok || !ev.Result() || ev.Text != "ok1" || ev.SessionID != "899693ea" || ev.IsError {
		t.Fatalf("result = %+v %v", ev, ok)
	}
	ev, ok = ParseStreamLine([]byte(`{"type":"assistant","message":{"content":[{"type":"text","text":"a"},{"type":"tool_use","name":"Bash"},{"type":"text","text":" "},{"type":"tool_use","name":"Read"}]}}`))
	if !ok || ev.Text != "a" || !slices.Equal(ev.Tools, []string{"Bash", "Read"}) {
		t.Fatalf("assistant = %+v", ev)
	}
	ev, ok = ParseStreamLine([]byte(`{"type":"system","subtype":"api_retry","attempt":1,"error_status":401,"error":"authentication_failed"}`))
	if !ok || ev.Subtype != "api_retry" || ev.RetryStatus != 401 {
		t.Fatalf("retry = %+v", ev)
	}
	for _, bad := range []string{"", "not json", `{"no":"type"}`} {
		if _, ok := ParseStreamLine([]byte(bad)); ok {
			t.Errorf("%q parsed as an event", bad)
		}
	}
}

func TestTurnArgs(t *testing.T) {
	launch := []string{"--dangerously-skip-permissions", "--resume", "old", "--settings", "s.json", "--model", "opus"}
	got := TurnArgs(launch, "new", false)
	want := []string{"--dangerously-skip-permissions", "--settings", "s.json", "--model", "opus", "--session-id", "new", "-p", "--output-format", "stream-json", "--verbose"}
	if !slices.Equal(got, want) {
		t.Fatalf("fresh = %q", got)
	}
	got = TurnArgs(launch, "new", true)
	if !slices.Equal(got[len(got)-6:len(got)-4], []string{"--resume", "new"}) {
		t.Fatalf("opened = %q", got)
	}
	if id := LaunchSessionID(launch); id != "old" {
		t.Fatalf("LaunchSessionID = %q", id)
	}
}
