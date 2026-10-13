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
	"flag"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultProtocol = "2025-06-18"
	toolName        = "read_probe"
)

func main() {
	httpAddr := flag.String("http", "", "serve streamable HTTP on this address instead of stdio")
	pidFile := flag.String("pid-file", "", "write fixture PID for shutdown assertions")
	oauth := flag.Bool("oauth", false, "require mock OAuth on HTTP; lab only")
	addressFile := flag.String("address-file", "", "write the actual HTTP URL here")
	flag.Parse()
	if *pidFile != "" {
		if err := os.WriteFile(*pidFile, []byte(strconv.Itoa(os.Getpid())), 0600); err != nil {
			panic(err)
		}
	}

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
	if *httpAddr != "" {
		listener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			panic(err)
		}
		if *addressFile != "" {
			if err := os.WriteFile(*addressFile, []byte("http://"+listener.Addr().String()+"/mcp"), 0600); err != nil {
				panic(err)
			}
		}
		base := "http://" + listener.Addr().String()
		var mu sync.Mutex
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if *oauth {
				write := func(v any) { w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(v) }
				switch r.URL.Path {
				case "/.well-known/oauth-protected-resource/mcp", "/.well-known/oauth-protected-resource":
					write(map[string]any{"resource": base + "/mcp", "authorization_servers": []string{base}, "scopes_supported": []string{"read"}, "bearer_methods_supported": []string{"header"}})
					return
				case "/.well-known/oauth-authorization-server":
					write(map[string]any{"issuer": base, "authorization_endpoint": base + "/authorize", "token_endpoint": base + "/token", "registration_endpoint": base + "/register", "response_types_supported": []string{"code"}, "grant_types_supported": []string{"authorization_code", "refresh_token"}, "code_challenge_methods_supported": []string{"S256"}, "token_endpoint_auth_methods_supported": []string{"none"}})
					return
				case "/register":
					var d map[string]any
					_ = json.NewDecoder(r.Body).Decode(&d)
					d["client_id"] = "mate-fixture-client"
					d["client_id_issued_at"] = time.Now().Unix()
					write(d)
					return
				case "/authorize":
					redirect, err := url.Parse(r.URL.Query().Get("redirect_uri"))
					if err != nil || redirect.Scheme != "http" || (redirect.Hostname() != "localhost" && redirect.Hostname() != "127.0.0.1") {
						http.Error(w, "invalid lab callback", 400)
						return
					}
					q := redirect.Query()
					q.Set("code", "fixture-code")
					q.Set("state", r.URL.Query().Get("state"))
					redirect.RawQuery = q.Encode()
					http.Redirect(w, r, redirect.String(), 302)
					return
				case "/token":
					write(map[string]any{"access_token": "mate-fixture-token", "refresh_token": "mate-fixture-refresh", "token_type": "Bearer", "expires_in": 3600, "scope": "read"})
					return
				}
				if r.Header.Get("Authorization") != "Bearer mate-fixture-token" {
					w.Header().Set("WWW-Authenticate", `Bearer resource_metadata="`+base+`/.well-known/oauth-protected-resource/mcp", scope="read"`)
					w.WriteHeader(401)
					return
				}
			}
			if r.Method == "DELETE" {
				w.WriteHeader(204)
				return
			}
			if r.Method != "POST" {
				w.WriteHeader(405)
				return
			}
			b, err := io.ReadAll(io.LimitReader(r.Body, 4<<20))
			if err != nil {
				http.Error(w, "read", 400)
				return
			}
			if logFile != nil {
				mu.Lock()
				_, _ = logFile.Write(append(b, '\n'))
				mu.Unlock()
			}
			resp, ok := handle(b)
			if !ok {
				w.WriteHeader(202)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		if err := http.Serve(listener, handler); err != nil {
			panic(err)
		}
		return
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
		if *p.Arguments.ID == "error" {
			return okReply(req.ID, map[string]any{"isError": true, "content": []any{map[string]any{"type": "text", "text": "probe:error:failed"}}}), true
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
