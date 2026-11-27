package agy

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/telemetry"
	_ "modernc.org/sqlite" // agy's conversations are SQLite; pure Go, as internal/db
)

// TranscriptAgy is the format of an agy conversation database, as the
// timeline's cursor and fact tables record it.
const TranscriptAgy harness.TranscriptFormat = "agy_steps"

// An agy conversation is conversations/<id>.db, SQLite in WAL mode. Table
// `steps` holds one row per step, keyed by `idx`; `metadata` is a protobuf
// with no published schema. Measured on agy 1.5.18 (2026-11-27) against
// `agy -p --output-format stream-json` on the same conversation:
//
//   - step_type 15 is a model response. metadata field 9 is its usage:
//     9.2 input tokens without the cache (stream-json input_tokens), 9.5
//     cache reads (absent when 0), 9.3 output tokens with thinking, 9.9
//     thinking (absent when 0), 9.10 visible output (9.3 = 9.9 + 9.10).
//     9.1 is a model number agy does not name; the Crew's meta names it.
//   - step_type 132 is a tool call: metadata field 4 is {1 call id, 2 tool
//     name, 3 arguments JSON}.
//   - metadata field 1 is when the step started and field 7 when it was
//     done, each {1 seconds, 2 nanoseconds}; on a response, 7 - 1 was the
//     stream-json duration_seconds.
//   - status 3 is done; 6 and 7 are a step that stopped early, and a
//     response in either still carries the usage it was billed.
//
// Other step types (17, 101) carry no usage and are counted as skipped. A
// step whose metadata does not decode costs only itself: a response becomes
// a usage failure and any other step is skipped, so a new shape never stops
// the rest of the conversation from being read. A conversation is read
// whole on every change, so a step's position is its idx and ConsumedBytes
// is one past the last idx.
const (
	stepUserInput = 14
	stepResponse  = 15
	stepToolCall  = 132

	statusDone = 3
)

// agyTranscripts is agy's TranscriptSource.
type agyTranscripts struct {
	// home is the profile's Home; empty is agy's own.
	home string
}

// agyRuleCwdMap is the rule a conversation is located by: the conversation
// agy's map names for the Crew's worktree, written at or after the launch.
const agyRuleCwdMap = "agy.last_conversation"

// Locate implements TranscriptSource. agy names its conversation only once
// it runs, so a Crew usually has no session id: the map agy keeps of the
// last conversation each directory ran names it, as at stop.
func (s agyTranscripts) Locate(req harness.TranscriptLocateRequest) (harness.TranscriptLocation, string) {
	home := req.Root
	if home == "" {
		var err error
		if home, err = dataHome(s.home); err != nil {
			return harness.TranscriptLocation{}, harness.LocateNotFound
		}
	}
	id := strings.TrimSpace(req.SessionID)
	if id == "" {
		id = agySessions{home: home}.AtStop(req.Cwd, req.LaunchedAt, "")
	}
	if id == "" {
		return harness.TranscriptLocation{}, harness.LocateNoSession
	}
	path := ConversationPath(home, id)
	if fi, err := os.Stat(path); err != nil || !fi.Mode().IsRegular() {
		return harness.TranscriptLocation{}, harness.LocateNotFound
	}
	return harness.TranscriptLocation{Path: path, SessionID: id, Rule: agyRuleCwdMap}, ""
}

// Read implements TranscriptSource. The database is opened read-only, so
// steps still in agy's WAL are read too.
func (agyTranscripts) Read(req harness.TranscriptReadRequest) (harness.TranscriptBatch, error) {
	steps, err := readSteps(req.Path)
	if err != nil {
		return harness.TranscriptBatch{}, err
	}
	return stepsBatch(steps), nil
}

// step is one row of `steps`.
type step struct {
	idx      int64
	kind     int64
	status   int64
	metadata []byte
}

func readSteps(path string) ([]step, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	dsn := (&url.URL{Scheme: "file", Path: abs, RawQuery: "mode=ro"}).String()
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT idx, step_type, status, metadata FROM steps ORDER BY idx`)
	if err != nil {
		return nil, fmt.Errorf("agy conversation %s: %w", path, err)
	}
	defer rows.Close()
	var out []step
	for rows.Next() {
		var s step
		if err := rows.Scan(&s.idx, &s.kind, &s.status, &s.metadata); err != nil {
			return nil, fmt.Errorf("agy conversation %s: %w", path, err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// stepsBatch normalises a conversation's steps. A tool call belongs to the
// response before it, the one that issued it.
func stepsBatch(steps []step) harness.TranscriptBatch {
	b := harness.TranscriptBatch{Source: TranscriptAgy, Skipped: map[string]int{}}
	var prompt string
	last := -1
	for ordinal, s := range steps {
		md, err := decodeMessage(s.metadata)
		if err != nil {
			if s.kind == stepResponse {
				b.UsageFailures = append(b.UsageFailures, harness.TranscriptUsageFailure{
					Offset: s.idx, Ordinal: ordinal, Reason: harness.MalformedUsageSnapshot, Detail: "metadata: " + err.Error()})
			} else {
				b.Skipped["undecodable_metadata"]++
			}
			continue
		}
		started, ended := md.time(1), md.time(7)
		switch s.kind {
		case stepUserInput:
			prompt = fmt.Sprint(s.idx)
			b.Records = append(b.Records, stepRecord(harness.TranscriptKindUser, s, ordinal, started))
		case stepResponse:
			usage, ok, err := responseUsage(md)
			if err != nil {
				b.UsageFailures = append(b.UsageFailures, harness.TranscriptUsageFailure{
					Offset: s.idx, Ordinal: ordinal, Reason: harness.MalformedUsageSnapshot, Detail: err.Error()})
				continue
			}
			if !ok {
				// Still streaming: no usage until it is done.
				b.Skipped["response_without_usage"]++
				continue
			}
			ref := fmt.Sprint(s.idx)
			outcome := "end_turn"
			if s.status != statusDone {
				outcome = fmt.Sprintf("status_%d", s.status)
			}
			b.Records = append(b.Records, stepRecord(harness.TranscriptKindUsage, s, ordinal, ended))
			b.Turns = append(b.Turns, harness.TranscriptTurn{
				SourceRef: ref, Offset: s.idx, Ordinal: ordinal, OccurredAt: ended,
				StopReason: outcome, Usage: usage, HarnessTurnRef: prompt,
			})
			b.UsageTurns = append(b.UsageTurns, harness.UsageTurn{
				SourceRef: ref, Offset: s.idx, StartedAt: started, EndedAt: ended,
				Outcome: outcome, HarnessTurnRef: prompt, Usage: usage, ContextTokens: usage.ContextTokens(),
				Sample: harness.UsageSample{At: ended, Usage: usage},
			})
			last = len(b.UsageTurns) - 1
		case stepToolCall:
			call, ok := toolCall(md, s, ordinal, started)
			if !ok {
				b.Skipped["tool_call_without_id"]++
				continue
			}
			b.Records = append(b.Records, stepRecord(harness.TranscriptKindAssistant, s, ordinal, started))
			b.ToolCalls = append(b.ToolCalls, call)
			if s.status >= statusDone {
				b.ToolResults = append(b.ToolResults, harness.TranscriptToolResult{
					SourceRef: call.SourceRef, CallID: call.SourceRef, Offset: s.idx, Ordinal: ordinal,
					IsError: s.status != statusDone, CompletedAt: ended,
				})
			}
			if last >= 0 {
				call.TurnSourceRef = b.UsageTurns[last].SourceRef
				b.UsageTurns[last].Calls = append(b.UsageTurns[last].Calls, call)
				if ended.After(b.UsageTurns[last].EndedAt) {
					b.UsageTurns[last].EndedAt = ended
				}
			} else {
				b.UnpricedCalls = append(b.UnpricedCalls, call)
			}
		default:
			b.Skipped[fmt.Sprintf("step_type_%d", s.kind)]++
		}
	}
	// A response that issued tool calls did not end the turn.
	for i := range b.UsageTurns {
		if len(b.UsageTurns[i].Calls) > 0 && b.UsageTurns[i].Outcome == "end_turn" {
			b.UsageTurns[i].Outcome = "tool_use"
		}
	}
	if n := len(steps); n > 0 {
		b.TotalBytes = steps[n-1].idx + 1
	}
	b.ConsumedBytes = b.TotalBytes
	return b
}

func stepRecord(kind harness.TranscriptKind, s step, ordinal int, at time.Time) harness.TranscriptRecord {
	return harness.TranscriptRecord{SourceRef: fmt.Sprint(s.idx), Kind: kind, Offset: s.idx, Ordinal: ordinal, OccurredAt: at}
}

// responseUsage is a response step's usage in the four buckets every
// harness reports, false when the step has none yet. A shape that was not
// measured is an error, not a guess.
func responseUsage(md message) (harness.TokenUsage, bool, error) {
	raw, ok := md.bytes(9)
	if !ok {
		return harness.TokenUsage{}, false, nil
	}
	u, err := decodeMessage(raw)
	if err != nil {
		return harness.TokenUsage{}, false, fmt.Errorf("usage: %w", err)
	}
	input, hasInput := u.varint(2)
	if !hasInput {
		return harness.TokenUsage{}, false, errors.New("usage has no input tokens (9.2)")
	}
	output, _ := u.varint(3)
	thinking, _ := u.varint(9)
	visible, hasVisible := u.varint(10)
	cache, _ := u.varint(5)
	if thinking > output {
		return harness.TokenUsage{}, false, fmt.Errorf("thinking %d exceeds output %d", thinking, output)
	}
	if hasVisible && thinking+visible != output {
		return harness.TokenUsage{}, false, fmt.Errorf("thinking %d plus visible output %d is not output %d", thinking, visible, output)
	}
	return harness.TokenUsage{Input: int64(input), CacheRead: int64(cache), Output: int64(output), Reasoning: int64(thinking)}, true, nil
}

func toolCall(md message, s step, ordinal int, at time.Time) (harness.TranscriptToolCall, bool) {
	raw, ok := md.bytes(4)
	if !ok {
		return harness.TranscriptToolCall{}, false
	}
	c, err := decodeMessage(raw)
	if err != nil {
		return harness.TranscriptToolCall{}, false
	}
	id, _ := c.bytes(1)
	name, _ := c.bytes(2)
	args, _ := c.bytes(3)
	if len(id) == 0 {
		return harness.TranscriptToolCall{}, false
	}
	return harness.TranscriptToolCall{
		SourceRef: string(id), Offset: s.idx, Ordinal: ordinal, ToolName: string(name),
		CommandClass: classifyTool(string(name)), InputJSON: string(args), StartedAt: at,
	}, true
}

// classifyTool maps agy's tool names, as the measured conversations named
// them, to the classes every harness shares.
func classifyTool(name string) harness.CommandClass {
	switch name {
	case "run_command", "command_status":
		return harness.CommandShell
	case "read_file", "view_file":
		return harness.CommandRead
	case "edit_file", "write_file":
		return harness.CommandEdit
	case "search_web", "grep_search", "find_by_name":
		return harness.CommandSearch
	case "call_mcp_tool":
		return harness.CommandMCP
	case "invoke_subagent", "manage_subagents", "browser_subagent":
		return harness.CommandAgent
	}
	return harness.CommandOther
}

// Telemetry implements TranscriptSource. Each response is one response
// fact. The state keeps no header hash: a SQLite database rewrites its
// first page on every commit, so a changed header is not a replaced file
// here, and the conversation is read whole anyway.
func (agyTranscripts) Telemetry(req harness.TelemetryRequest) (harness.TelemetryUpdate, error) {
	state := req.State
	state.HeaderSHA256, state.HeaderBytes = "", 0
	up := harness.TelemetryUpdate{Offset: req.Offset, State: state, Gaps: []string{
		"model_named_by_crew_meta", "tool_output_not_read",
	}}
	if req.Offset == req.Batch.ConsumedBytes && !req.Changed {
		return up, nil
	}
	for _, t := range req.Batch.Turns {
		off := t.Offset
		up.Facts = append(up.Facts, telemetry.Fact{
			Version: telemetry.Version, ID: t.SourceRef, Kind: "response", SourceRef: t.SourceRef, SourceOffset: off,
			OccurredAt: t.OccurredAt, MeasurementKind: "normalized", HarnessTurnRef: t.HarnessTurnRef, ResponseID: t.SourceRef,
			InputTokens: t.Usage.Input, CacheReadTokens: t.Usage.CacheRead, CacheWriteTokens: t.Usage.CacheWrite,
			OutputTokens: t.Usage.Output, ReasoningTokens: t.Usage.Reasoning,
			ContextTokens: t.Usage.ContextTokens(), LedgerRefOffset: &off,
		})
	}
	up.Offset = req.Batch.ConsumedBytes
	return up, nil
}

// message is one decoded protobuf message: each field number's values, in
// order. Only the wire types the measured metadata uses are read.
type message map[int][]field

type field struct {
	varint uint64
	bytes  []byte
	isLen  bool
}

func (m message) varint(n int) (uint64, bool) {
	for _, f := range m[n] {
		if !f.isLen {
			return f.varint, true
		}
	}
	return 0, false
}

func (m message) bytes(n int) ([]byte, bool) {
	for _, f := range m[n] {
		if f.isLen {
			return f.bytes, true
		}
	}
	return nil, false
}

// time is a {1 seconds, 2 nanoseconds} field, or the zero time.
func (m message) time(n int) time.Time {
	raw, ok := m.bytes(n)
	if !ok {
		return time.Time{}
	}
	ts, err := decodeMessage(raw)
	if err != nil {
		return time.Time{}
	}
	sec, ok := ts.varint(1)
	if !ok {
		return time.Time{}
	}
	nanos, _ := ts.varint(2)
	return time.Unix(int64(sec), int64(nanos)).UTC()
}

func decodeMessage(b []byte) (message, error) {
	m := message{}
	for i := 0; i < len(b); {
		key, n := binary.Uvarint(b[i:])
		if n <= 0 {
			return nil, errors.New("truncated field key")
		}
		i += n
		num, wire := int(key>>3), key&7
		if num == 0 {
			return nil, errors.New("field number 0")
		}
		switch wire {
		case 0:
			v, n := binary.Uvarint(b[i:])
			if n <= 0 {
				return nil, errors.New("truncated varint")
			}
			i += n
			m[num] = append(m[num], field{varint: v})
		case 1:
			if i+8 > len(b) {
				return nil, errors.New("truncated fixed64")
			}
			i += 8
		case 2:
			l, n := binary.Uvarint(b[i:])
			if n <= 0 || l > uint64(len(b)-i-n) {
				return nil, errors.New("truncated length-delimited field")
			}
			i += n
			m[num] = append(m[num], field{bytes: b[i : i+int(l)], isLen: true})
			i += int(l)
		case 5:
			if i+4 > len(b) {
				return nil, errors.New("truncated fixed32")
			}
			i += 4
		default:
			return nil, fmt.Errorf("wire type %d", wire)
		}
	}
	return m, nil
}
