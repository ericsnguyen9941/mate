package claude

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
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
	// Tools are the names of an assistant event's tool_use blocks, in order.
	Tools []string
	// SessionID is the event's `session_id`.
	SessionID string
	// IsError is a result event's `is_error`.
	IsError bool
	// RetryStatus is an api_retry event's `error_status`.
	RetryStatus int
}

// Result reports whether the event ends the run.
func (e StreamEvent) Result() bool { return e.Type == "result" }

type streamLine struct {
	Type        string `json:"type"`
	Subtype     string `json:"subtype"`
	SessionID   string `json:"session_id"`
	Result      string `json:"result"`
	IsError     bool   `json:"is_error"`
	ErrorStatus int    `json:"error_status"`
	Message     struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
			Name string `json:"name"`
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
		ev.Text = l.Result
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
			}
		}
		ev.Text = strings.Join(parts, "\n")
	}
	return ev, true
}

// StreamTurnFlags are what turns a Mate launch into one headless turn.
var StreamTurnFlags = []string{"-p", "--output-format", "stream-json", "--verbose"}

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
