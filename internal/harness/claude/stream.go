package claude

import (
	"encoding/json"
	"os"
	"slices"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// A headless Mate turn is one `claude -p --output-format stream-json
// --verbose` run (docs/mvp.md, CON-78). Its stdout is one JSON event per
// line; measured on Claude Code 2.1.292 (docs/evidence/headless-mate-probe-
// 2026-10-15.md) the types are `system` (subtypes hook_started,
// hook_response, init, thinking_tokens, api_retry), `assistant`, `user` and
// a final `result`. Only the shapes read here are relied on.

// StreamEvent is one line of a stream-json run, reduced to what a host
// shows while a turn runs and what it keeps when the turn ends.
type StreamEvent struct {
	// Type is the event's `type`.
	Type string
	// Subtype is a system or result event's `subtype`.
	Subtype string
	// Text is an assistant event's text blocks joined, or a result event's
	// `result`.
	Text string
	// Delta is a stream_event's text_delta (--include-partial-messages).
	Delta string
	// Tools are the names of an assistant event's tool_use blocks, in order.
	Tools []string
	// Calls are an assistant event's tool_use blocks, in order.
	Calls []harness.ToolCall
	// Results are a user event's tool_result blocks.
	Results []harness.ToolResult
	// Context is an assistant event's request context in tokens: new
	// input, cache writes and cache reads.
	Context int
	// Tokens is a result event's tokens other than cache reads.
	Tokens int
	// DurationMS is a result event's `duration_ms`.
	DurationMS int64
	// SessionID is the event's `session_id`.
	SessionID string
	// IsError is a result event's `is_error`.
	IsError bool
	// RetryStatus is an api_retry event's `error_status`.
	RetryStatus int
}

// Result reports whether the event ends the run.
func (e StreamEvent) Result() bool { return e.Type == "result" }

type streamUsage struct {
	Input      int `json:"input_tokens"`
	CacheWrite int `json:"cache_creation_input_tokens"`
	CacheRead  int `json:"cache_read_input_tokens"`
	Output     int `json:"output_tokens"`
}

type streamLine struct {
	Type        string       `json:"type"`
	Subtype     string       `json:"subtype"`
	SessionID   string       `json:"session_id"`
	Result      string       `json:"result"`
	IsError     bool         `json:"is_error"`
	ErrorStatus int          `json:"error_status"`
	DurationMS  int64        `json:"duration_ms"`
	Usage       *streamUsage `json:"usage"`
	Event       struct {
		Type  string `json:"type"`
		Delta struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"delta"`
	} `json:"event"`
	Message struct {
		Usage   *streamUsage `json:"usage"`
		Content []struct {
			Type      string          `json:"type"`
			Text      string          `json:"text"`
			ID        string          `json:"id"`
			Name      string          `json:"name"`
			Input     json.RawMessage `json:"input"`
			ToolUseID string          `json:"tool_use_id"`
			IsError   bool            `json:"is_error"`
		} `json:"content"`
	} `json:"message"`
}

// ParseStreamLine reads one stream-json line. A blank line, or one that is
// not a JSON object with a type, is not an event: ok is false.
func ParseStreamLine(line []byte) (StreamEvent, bool) {
	var l streamLine
	if err := json.Unmarshal(line, &l); err != nil || l.Type == "" {
		return StreamEvent{}, false
	}
	ev := StreamEvent{Type: l.Type, Subtype: l.Subtype, SessionID: l.SessionID, IsError: l.IsError, RetryStatus: l.ErrorStatus}
	switch l.Type {
	case "result":
		ev.Text, ev.DurationMS = l.Result, l.DurationMS
		if u := l.Usage; u != nil {
			ev.Tokens = u.Input + u.CacheWrite + u.Output
		}
	case "stream_event":
		if l.Event.Type == "content_block_delta" && l.Event.Delta.Type == "text_delta" {
			ev.Delta = l.Event.Delta.Text
		}
	case "assistant":
		var parts []string
		for _, c := range l.Message.Content {
			switch c.Type {
			case "text":
				if strings.TrimSpace(c.Text) != "" {
					parts = append(parts, c.Text)
				}
			case "tool_use":
				ev.Tools = append(ev.Tools, c.Name)
				ev.Calls = append(ev.Calls, harness.ToolCall{ID: c.ID, Name: c.Name, Detail: toolDetail(c.Name, c.Input)})
			}
		}
		ev.Text = strings.Join(parts, "\n")
		if u := l.Message.Usage; u != nil {
			ev.Context = u.Input + u.CacheWrite + u.CacheRead
		}
	case "user":
		for _, c := range l.Message.Content {
			if c.Type == "tool_result" {
				ev.Results = append(ev.Results, harness.ToolResult{ID: c.ToolUseID, IsError: c.IsError})
			}
		}
	}
	return ev, true
}

// toolDetail is the one line of a tool call's input that says what it does.
func toolDetail(name string, input json.RawMessage) string {
	var in map[string]any
	if json.Unmarshal(input, &in) != nil {
		return ""
	}
	for _, key := range []string{"command", "file_path", "path", "pattern", "skill", "url", "query", "description"} {
		if v, ok := in[key].(string); ok && strings.TrimSpace(v) != "" {
			line, _, _ := strings.Cut(strings.TrimSpace(v), "\n")
			return line
		}
	}
	return ""
}

// StreamTurnFlags are what turns a Mate launch into one headless turn.
var StreamTurnFlags = []string{"-p", "--output-format", "stream-json", "--verbose", "--include-partial-messages"}

// TurnArgs is a headless Mate turn's argv after the program: the Mate's
// launch args (Claude.Build) with its session flag replaced, then
// StreamTurnFlags. A session Claude has already opened is resumed; one it
// has not is created under the same id, so the first turn of a fresh start,
// and a start that resumed a session no turn ever opened, both work.
func TurnArgs(launch []string, sessionID string, opened bool) []string {
	out := make([]string, 0, len(launch)+len(StreamTurnFlags)+2)
	for i := 0; i < len(launch); i++ {
		if launch[i] == "--resume" || launch[i] == "--session-id" {
			i++
			continue
		}
		out = append(out, launch[i])
	}
	if sessionID != "" {
		flag := "--session-id"
		if opened {
			flag = "--resume"
		}
		out = append(out, flag, sessionID)
	}
	return append(out, StreamTurnFlags...)
}

// LaunchSessionID is the session a launch's args name, by either flag.
func LaunchSessionID(launch []string) string {
	for _, flag := range []string{"--resume", "--session-id"} {
		if i := slices.Index(launch, flag); i >= 0 && i+1 < len(launch) {
			return launch[i+1]
		}
	}
	return ""
}

// SessionOpened reports whether Claude has written the transcript of the
// session it would run in cwd under sessionID: whether a turn has opened it.
func SessionOpened(cwd, sessionID string) bool {
	if sessionID == "" {
		return false
	}
	root, err := ClaudeProjectsDir()
	if err != nil {
		return false
	}
	fi, err := os.Stat(ClaudeTranscriptPath(root, cwd, sessionID))
	return err == nil && !fi.IsDir()
}

// TurnArgs implements harness.HeadlessTurns.
func (Claude) TurnArgs(launch []string, sessionID string, opened bool) []string {
	return TurnArgs(launch, sessionID, opened)
}

// LaunchSessionID implements harness.HeadlessTurns.
func (Claude) LaunchSessionID(launch []string) string { return LaunchSessionID(launch) }

// SessionOpened implements harness.HeadlessTurns.
func (Claude) SessionOpened(cwd, sessionID string) bool { return SessionOpened(cwd, sessionID) }

// ParseTurnEvent implements harness.HeadlessTurns.
func (Claude) ParseTurnEvent(line []byte) (harness.TurnEvent, bool) {
	ev, ok := ParseStreamLine(line)
	if !ok {
		return harness.TurnEvent{}, false
	}
	out := harness.TurnEvent{Delta: ev.Delta, Calls: ev.Calls, Results: ev.Results, Context: ev.Context}
	if ev.Type == "assistant" {
		out.Said = ev.Text
	}
	if ev.Result() {
		out.Final, out.Text, out.SessionID, out.IsError = true, ev.Text, ev.SessionID, ev.IsError
		out.DurationMS, out.Tokens = ev.DurationMS, ev.Tokens
	}
	return out, true
}
