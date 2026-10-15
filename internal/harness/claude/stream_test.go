package claude

import (
	"slices"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
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
	ev, ok = ParseStreamLine([]byte(`{"type":"assistant","message":{"usage":{"input_tokens":3,"cache_creation_input_tokens":40,"cache_read_input_tokens":500,"output_tokens":9},"content":[{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"mate status\nmore","description":"x"}},{"type":"tool_use","id":"t2","name":"Read","input":{"file_path":"/a/b.md"}},{"type":"tool_use","id":"t3","name":"TodoWrite","input":{"todos":[]}}]}}`))
	if !ok || ev.Context != 543 || !slices.Equal(ev.Calls, []harness.ToolCall{{ID: "t1", Name: "Bash", Detail: "mate status"}, {ID: "t2", Name: "Read", Detail: "/a/b.md"}, {ID: "t3", Name: "TodoWrite"}}) {
		t.Fatalf("tool calls = %+v", ev)
	}
	ev, ok = ParseStreamLine([]byte(`{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"t1","is_error":true,"content":"boom"}]}}`))
	if !ok || !slices.Equal(ev.Results, []harness.ToolResult{{ID: "t1", IsError: true}}) {
		t.Fatalf("tool results = %+v", ev)
	}
	ev, ok = ParseStreamLine([]byte(`{"type":"stream_event","event":{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Xin "}}}`))
	if !ok || ev.Delta != "Xin " {
		t.Fatalf("delta = %+v", ev)
	}
	ev, ok = ParseStreamLine([]byte(`{"type":"result","subtype":"success","result":"ok","duration_ms":3200,"usage":{"input_tokens":5,"cache_creation_input_tokens":100,"cache_read_input_tokens":90000,"output_tokens":20}}`))
	if !ok || ev.DurationMS != 3200 || ev.Tokens != 125 {
		t.Fatalf("result usage = %+v", ev)
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
	want := []string{"--dangerously-skip-permissions", "--settings", "s.json", "--model", "opus", "--session-id", "new", "-p", "--output-format", "stream-json", "--verbose", "--include-partial-messages"}
	if !slices.Equal(got, want) {
		t.Fatalf("fresh = %q", got)
	}
	got = TurnArgs(launch, "new", true)
	if !slices.Equal(got[len(got)-7:len(got)-5], []string{"--resume", "new"}) {
		t.Fatalf("opened = %q", got)
	}
	if id := LaunchSessionID(launch); id != "old" {
		t.Fatalf("LaunchSessionID = %q", id)
	}
}
