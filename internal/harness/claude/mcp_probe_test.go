package claude

import "testing"

// TestMCPFixtureCalledWith pins the verdict TestLiveClaudeMCPProbe records:
// only a read_probe tools/call whose arguments.id matches counts, wherever
// else the id appears in the line.
func TestMCPFixtureCalledWith(t *testing.T) {
	// A line as Claude Code 2.1.287 sent it on 2026-10-12.
	logged := `{"method":"tools/call","params":{"name":"read_probe","arguments":{"id":"abc123"},"_meta":{"claudecode/toolUseId":"toolu_013xhiZCk1C9iSnJi7Asd5tW","progressToken":2}},"jsonrpc":"2.0","id":2}`
	cases := []struct {
		name string
		log  string
		want bool
	}{
		{"real call", `{"method":"initialize","jsonrpc":"2.0","id":0}` + "\n" + logged + "\n", true},
		{"id only in _meta", `{"method":"tools/call","params":{"name":"read_probe","arguments":{"id":"other"},"_meta":{"note":"abc123"}},"jsonrpc":"2.0","id":2}`, false},
		{"id as the JSON-RPC id", `{"method":"tools/call","params":{"name":"read_probe","arguments":{}},"jsonrpc":"2.0","id":"abc123"}`, false},
		{"other tool", `{"method":"tools/call","params":{"name":"x","arguments":{"id":"abc123"}},"jsonrpc":"2.0","id":2}`, false},
		{"not a call", `{"method":"tools/list","params":{"arguments":{"id":"abc123"}},"jsonrpc":"2.0","id":1}`, false},
		{"garbage", "not json abc123\n", false},
		{"empty", "", false},
	}
	for _, tc := range cases {
		if got := fixtureCalledWith(tc.log, "abc123"); got != tc.want {
			t.Errorf("%s: fixtureCalledWith = %v, want %v", tc.name, got, tc.want)
		}
	}
}
