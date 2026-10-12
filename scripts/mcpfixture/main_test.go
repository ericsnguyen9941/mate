package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestServe(t *testing.T) {
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"t","version":"0"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":"three","method":"tools/call","params":{"name":"read_probe","arguments":{"id":"abc123"}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"resources/list"}`,
		`{"jsonrpc":"2.0","id":5,"method":"ping"}`,
		`{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"read_probe","arguments":{}}}`,
		`{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"other","arguments":{"id":"x"}}}`,
		`not json`,
		``,
	}, "\n")
	logPath := filepath.Join(t.TempDir(), "log")
	logFile, err := os.Create(logPath)
	if err != nil {
		t.Fatal(err)
	}
	defer logFile.Close()
	var out bytes.Buffer
	w := bufio.NewWriter(&out)
	if err := serve(bufio.NewReader(strings.NewReader(input)), w, logFile); err != nil {
		t.Fatalf("serve: %v", err)
	}

	var replies []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("reply %q is not JSON: %v", line, err)
		}
		if m["jsonrpc"] != "2.0" {
			t.Errorf("reply %q lacks jsonrpc 2.0", line)
		}
		replies = append(replies, m)
	}
	// Nine non-empty lines, one of them a notification.
	if len(replies) != 8 {
		t.Fatalf("got %d replies, want 8:\n%s", len(replies), out.String())
	}

	init := replies[0]["result"].(map[string]any)
	if init["protocolVersion"] != "2025-03-26" {
		t.Errorf("initialize protocolVersion = %v, want the client's 2025-03-26", init["protocolVersion"])
	}
	if _, ok := init["capabilities"].(map[string]any)["tools"]; !ok {
		t.Errorf("initialize does not advertise tools: %v", init)
	}

	tools := replies[1]["result"].(map[string]any)["tools"].([]any)
	if len(tools) != 1 {
		t.Fatalf("tools/list returned %d tools, want exactly 1", len(tools))
	}
	tool := tools[0].(map[string]any)
	if tool["name"] != "read_probe" {
		t.Errorf("tool name = %v, want read_probe", tool["name"])
	}
	schema, _ := json.Marshal(tool["inputSchema"])
	if want := `{"properties":{"id":{"type":"string"}},"required":["id"],"type":"object"}`; string(schema) != want {
		t.Errorf("inputSchema = %s, want %s", schema, want)
	}

	if replies[2]["id"] != "three" {
		t.Errorf("tools/call reply id = %v, want the string id echoed", replies[2]["id"])
	}
	content := replies[2]["result"].(map[string]any)["content"].([]any)[0].(map[string]any)
	if content["type"] != "text" || content["text"] != "probe:abc123:ok" {
		t.Errorf("tools/call content = %v, want text probe:abc123:ok", content)
	}

	for i, want := range map[int]float64{3: -32601, 5: -32602, 6: -32602, 7: -32700} {
		e, ok := replies[i]["error"].(map[string]any)
		if !ok || e["code"] != want {
			t.Errorf("reply %d = %v, want error %v", i, replies[i], want)
		}
	}
	if r, ok := replies[4]["result"].(map[string]any); !ok || len(r) != 0 {
		t.Errorf("ping reply = %v, want an empty result object", replies[4])
	}

	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(logged), "\n"); n != 9 {
		t.Errorf("log has %d lines, want 9 (every non-empty request):\n%s", n, logged)
	}
	if !strings.Contains(string(logged), `"method":"tools/call"`) {
		t.Errorf("log lacks the tools/call request:\n%s", logged)
	}
}

func TestInitializeDefaultsProtocol(t *testing.T) {
	resp, ok := handle([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`))
	if !ok {
		t.Fatal("initialize got no reply")
	}
	if v := resp.Result.(map[string]any)["protocolVersion"]; v != defaultProtocol {
		t.Errorf("protocolVersion = %v, want %s", v, defaultProtocol)
	}
}
