// Command mcpfixture is a minimal stdio MCP server with exactly one tool,
// read_probe. It exists so a live test can tell whether a Claude Code launch
// discovers and calls an MCP tool under Mate's real launch flags
// (internal/harness/claude/mcp_probe_live_test.go).
//
// The transport is newline-delimited JSON-RPC 2.0 on stdin/stdout. It answers
// initialize (echoing the client's protocolVersion, 2025-06-18 if none), ping,
// tools/list and tools/call; notifications get no reply and any other request
// gets -32601. With MCP_FIXTURE_LOG set, every received line is appended to
// that file, so the test can read which requests arrived.
package main

import (
	"bufio"
	"encoding/json"
	"os"
)

const (
	defaultProtocol = "2025-06-18"
	toolName        = "read_probe"
)

func main() {
	var logFile *os.File
	if path := os.Getenv("MCP_FIXTURE_LOG"); path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			os.Stderr.WriteString("mcpfixture: open log: " + err.Error() + "\n")
			os.Exit(1)
		}
		defer f.Close()
		logFile = f
	}
	out := bufio.NewWriter(os.Stdout)
	if err := serve(bufio.NewReader(os.Stdin), out, logFile); err != nil {
		os.Stderr.WriteString("mcpfixture: " + err.Error() + "\n")
		os.Exit(1)
	}
}

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

// serve reads requests until EOF. Each reply is flushed as one line so the
// client never waits on a buffered answer.
func serve(in *bufio.Reader, out *bufio.Writer, logFile *os.File) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := trimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		if logFile != nil {
			if _, err := logFile.Write(append(append([]byte(nil), line...), '\n')); err != nil {
				return err
			}
		}
		resp, ok := handle(line)
		if !ok {
			continue
		}
		b, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		if _, err := out.Write(append(b, '\n')); err != nil {
			return err
		}
		if err := out.Flush(); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// handle returns the reply to one line, and false for a notification.
func handle(line []byte) (response, bool) {
	var req request
	if err := json.Unmarshal(line, &req); err != nil {
		return errorReply(json.RawMessage("null"), -32700, "parse error"), true
	}
	if len(req.ID) == 0 || string(req.ID) == "null" {
		// notifications/initialized and any other notification: no reply.
		return response{}, false
	}
	switch req.Method {
	case "initialize":
		var p struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		_ = json.Unmarshal(req.Params, &p)
		version := p.ProtocolVersion
		if version == "" {
			version = defaultProtocol
		}
		return okReply(req.ID, map[string]any{
			"protocolVersion": version,
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "mcpfixture", "version": "0.1.0"},
		}), true
	case "ping":
		return okReply(req.ID, map[string]any{}), true
	case "tools/list":
		return okReply(req.ID, map[string]any{"tools": []any{map[string]any{
			"name":        toolName,
			"description": "Read the probe record with the given id and return its status line.",
			"inputSchema": map[string]any{
				"type":       "object",
				"properties": map[string]any{"id": map[string]any{"type": "string"}},
				"required":   []string{"id"},
			},
		}}}), true
	case "tools/call":
		var p struct {
			Name      string `json:"name"`
			Arguments struct {
				ID *string `json:"id"`
			} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return errorReply(req.ID, -32602, "invalid params"), true
		}
		if p.Name != toolName {
			return errorReply(req.ID, -32602, "unknown tool: "+p.Name), true
		}
		if p.Arguments.ID == nil {
			return errorReply(req.ID, -32602, "missing required argument: id"), true
		}
		return okReply(req.ID, map[string]any{
			"content": []any{map[string]any{"type": "text", "text": "probe:" + *p.Arguments.ID + ":ok"}},
		}), true
	default:
		return errorReply(req.ID, -32601, "method not found: "+req.Method), true
	}
}

func okReply(id json.RawMessage, result any) response {
	return response{JSONRPC: "2.0", ID: id, Result: result}
}

func errorReply(id json.RawMessage, code int, msg string) response {
	return response{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: msg}}
}

func trimSpace(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r' || b[len(b)-1] == ' ' || b[len(b)-1] == '\t') {
		b = b[:len(b)-1]
	}
	for len(b) > 0 && (b[0] == ' ' || b[0] == '\t') {
		b = b[1:]
	}
	return b
}
