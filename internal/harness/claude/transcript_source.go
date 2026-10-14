package claude

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/telemetry"
)

// claudeTranscripts is Claude Code's TranscriptSource: its session JSONL
// under the config directory's projects/, tailed while it only grows, with
// usage already per call.
type claudeTranscripts struct {
	// configDir is the profile's CLAUDE_CONFIG_DIR; empty resolves it as
	// transcript discovery always has (EffectiveClaudeConfigDir).
	configDir string
}

// claudeRuleProjects is the rule a Claude transcript is located by:
// `session_id=` plus Claude's own naming rule,
// `<config>/projects/<slug of cwd>/<session-id>.jsonl`.
const claudeRuleProjects = "claude.projects"

// Locate implements TranscriptSource. Claude names its transcript by the
// session id mate gave it at launch, so the id and the cwd are the path.
func (s claudeTranscripts) Locate(req harness.TranscriptLocateRequest) (harness.TranscriptLocation, string) {
	if req.SessionID == "" {
		return harness.TranscriptLocation{}, harness.LocateNoSession
	}
	root := req.Root
	if root == "" {
		configDir, err := EffectiveClaudeConfigDir(s.configDir)
		if err != nil {
			return harness.TranscriptLocation{}, harness.LocateNotFound
		}
		root = filepath.Join(configDir, "projects")
	}
	path := ClaudeTranscriptPath(root, req.Cwd, req.SessionID)
	if fi, err := os.Stat(path); path == "" || err != nil || fi.IsDir() {
		return harness.TranscriptLocation{}, harness.LocateNotFound
	}
	return harness.TranscriptLocation{Path: path, Rule: claudeRuleProjects}, ""
}

// Read implements TranscriptSource. A file that has only grown since a live
// read of it is read from where that read stopped, with its parser state,
// and the tail appended. That is the parser's own incremental contract:
// ParseTranscript withholds the trailing message group, stops ConsumedBytes
// at its first record and pairs NextState with that cursor, so a resume
// re-reads the group whole and carries the open harness turn - a resume
// that got either wrong would charge one turn's tokens to another,
// permanently. Anything else reads the file whole: a prior that stopped at
// a malformed line, one read at rest (which emitted a group a tail would
// meet again), one this source did not produce, or a read at rest, which is
// parsed final to flush that group.
func (s claudeTranscripts) Read(req harness.TranscriptReadRequest) (harness.TranscriptBatch, error) {
	if st, ok := req.Prior.HarnessState.(claudeReadState); ok && req.Grown && !req.AtRest && !st.final && req.Prior.Malformed == nil {
		info, err := os.Stat(req.Path)
		if err != nil {
			return harness.TranscriptBatch{}, err
		}
		if base := req.Prior.ConsumedBytes; base <= info.Size() {
			data, err := readFileFrom(req.Path, base, info.Size())
			if err != nil {
				return harness.TranscriptBatch{}, err
			}
			tail, consumed := claudeParseRange(req.Prior.NextState, data, false)
			return appendClaudeBatch(req.Prior, st, tail, consumed, base), nil
		}
	}
	data, err := os.ReadFile(req.Path)
	if err != nil {
		return harness.TranscriptBatch{}, err
	}
	b, consumed := claudeParseRange(harness.TranscriptParseState{}, data, req.AtRest)
	tidyClaude(&b)
	normaliseClaude(&b)
	b.HarnessState = claudeReadState{lines: consumed.lines, skipped: consumed.skipped, final: req.AtRest}
	return b, nil
}

// claudeReadState is what a Claude read carries to the next read of the
// same file, beside the batch: the parse's view of the bytes before its
// cursor (claudeConsumed), and whether the read was final - a final read
// emitted the group a live one withholds, so it cannot be carried on from.
type claudeReadState struct {
	lines   int
	skipped map[string]int
	final   bool
}

// appendClaudeBatch appends a tail parsed from base, the prior batch's
// ConsumedBytes, to that batch. Every position the tail reports is relative
// to base, so offsets are rebased by base and ordinals by the lines the
// prior read consumed. The prior's facts all precede base and the tail's all
// follow it, so the two are concatenated, never merged; what describes the
// end of the file - the cursor, the withheld group, a malformed line, the
// parser state - is the tail's. Usage turns are derived again over the whole
// batch, because a result in the tail can end a call the prior priced.
func appendClaudeBatch(prior harness.TranscriptBatch, st claudeReadState, tail harness.TranscriptBatch, consumed claudeConsumed, base int64) harness.TranscriptBatch {
	lines := st.lines
	for n := range tail.Records {
		tail.Records[n].Offset += base
		tail.Records[n].Ordinal += lines
	}
	for n := range tail.Turns {
		tail.Turns[n].Offset += base
		tail.Turns[n].Ordinal += lines
	}
	for n := range tail.ToolCalls {
		tail.ToolCalls[n].Offset += base
		tail.ToolCalls[n].Ordinal += lines
	}
	for n := range tail.ToolResults {
		tail.ToolResults[n].Offset += base
		tail.ToolResults[n].Ordinal += lines
	}
	for n := range tail.UsageFailures {
		tail.UsageFailures[n].Offset += base
		tail.UsageFailures[n].Ordinal += lines
	}
	if m := tail.Malformed; m != nil {
		shifted := *m
		shifted.Offset += base
		if shifted.EndOffset != 0 {
			shifted.EndOffset += base
		}
		shifted.Ordinal += lines
		tail.Malformed = &shifted
	}
	if tail.OpenSourceRef != "" {
		tail.OpenOffset += base
	}
	skipped := make(map[string]int, len(st.skipped)+len(tail.Skipped))
	carried := make(map[string]int, len(st.skipped)+len(consumed.skipped))
	for kind, count := range st.skipped {
		skipped[kind] += count
		carried[kind] += count
	}
	for kind, count := range tail.Skipped {
		skipped[kind] += count
	}
	for kind, count := range consumed.skipped {
		carried[kind] += count
	}

	b := harness.TranscriptBatch{
		Source: TranscriptClaude,
		// Concat copies: the prior batch may still be held, and appending
		// into its spare capacity would write under it.
		Records:       slices.Concat(prior.Records, tail.Records),
		Turns:         slices.Concat(prior.Turns, tail.Turns),
		ToolCalls:     slices.Concat(prior.ToolCalls, tail.ToolCalls),
		ToolResults:   slices.Concat(prior.ToolResults, tail.ToolResults),
		UsageFailures: slices.Concat(prior.UsageFailures, tail.UsageFailures),
		ConsumedBytes: base + tail.ConsumedBytes,
		TotalBytes:    base + tail.TotalBytes,
		PendingBytes:  tail.PendingBytes,
		OpenSourceRef: tail.OpenSourceRef,
		OpenOffset:    tail.OpenOffset,
		OpenRecords:   tail.OpenRecords,
		Malformed:     tail.Malformed,
		NextState:     tail.NextState,
		Skipped:       skipped,
		HarnessState:  claudeReadState{lines: lines + consumed.lines, skipped: carried},
	}
	b.UsageTurns = claudeUsageTurns(b)
	// The prior's compactions are its records', all before base.
	b.Compactions = slices.Concat(prior.Compactions, claudeCompactions(tail.Records))
	return b
}

// tidyClaude makes an empty fact list nil, which is what a list that never
// had a fact is: a parse that withheld every fact it found leaves an empty
// one, and a tail read cannot tell the two apart.
func tidyClaude(b *harness.TranscriptBatch) {
	if len(b.Records) == 0 {
		b.Records = nil
	}
	if len(b.Turns) == 0 {
		b.Turns = nil
	}
	if len(b.ToolCalls) == 0 {
		b.ToolCalls = nil
	}
	if len(b.ToolResults) == 0 {
		b.ToolResults = nil
	}
}

// readFileFrom reads a file from offset to size, the size a caller has just
// measured. It returns what it could read when the file shrank under it;
// only a read that found nothing at all is an error.
func readFileFrom(path string, offset, size int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data := make([]byte, size-offset)
	n, err := f.ReadAt(data, offset)
	if err != nil && n == 0 && len(data) > 0 {
		return nil, err
	}
	return data[:n], nil
}

// normaliseClaude makes one usage turn per assistant message group, which
// is one API response: the unit Claude restates usage on, and the unit
// ParseTranscript already keys by message id. Claude reports usage per
// call and its input_tokens excludes both cache buckets, so a turn's usage
// is the response's own and needs no arithmetic.
//
// ContextTokens is the input side of that response - fresh input plus both
// cache buckets - because that is what was carried into the call and
// therefore what the context held when it returned.
func normaliseClaude(b *harness.TranscriptBatch) {
	b.UsageTurns, b.UnpricedCalls = claudeUsageTurns(*b), nil
	b.Compactions = claudeCompactions(b.Records)
}

// claudeUsageTurns is the batch's usage turns, one per response.
func claudeUsageTurns(b harness.TranscriptBatch) []harness.UsageTurn {
	results := firstResults(b)
	callsByTurn := map[string][]harness.TranscriptToolCall{}
	for _, call := range b.ToolCalls {
		callsByTurn[call.TurnSourceRef] = append(callsByTurn[call.TurnSourceRef], call)
	}
	var out []harness.UsageTurn
	for _, t := range b.Turns {
		calls := callsByTurn[t.SourceRef]
		ended := t.OccurredAt
		for _, call := range calls {
			if res, ok := results[call.SourceRef]; ok && res.CompletedAt.After(ended) {
				ended = res.CompletedAt
			} else if call.StartedAt.After(ended) {
				ended = call.StartedAt
			}
		}
		out = append(out, harness.UsageTurn{
			SourceRef: t.SourceRef, Offset: t.Offset,
			StartedAt: t.OccurredAt, EndedAt: ended,
			Outcome: t.StopReason, Model: t.Model, HarnessTurnRef: t.HarnessTurnRef,
			Usage: t.Usage, ContextTokens: t.Usage.ContextTokens(), Calls: calls,
			Sample: harness.UsageSample{At: t.OccurredAt, Usage: t.Usage},
		})
	}
	return out
}

// claudeCompactions is the compactions the records mark. Claude marks one on
// the summary record itself (`isCompactSummary`) or with a
// `compact_boundary` system record.
func claudeCompactions(records []harness.TranscriptRecord) []harness.TranscriptCompaction {
	var out []harness.TranscriptCompaction
	for _, rec := range records {
		if !strings.Contains(rec.RawJSON, `"compact_boundary"`) && !strings.Contains(rec.RawJSON, `"isCompactSummary"`) {
			continue
		}
		var probe struct {
			Type             string `json:"type"`
			Subtype          string `json:"subtype"`
			IsCompactSummary bool   `json:"isCompactSummary"`
		}
		if json.Unmarshal([]byte(rec.RawJSON), &probe) != nil {
			continue
		}
		trigger := ""
		switch {
		case probe.IsCompactSummary:
			trigger = "claude.compact_summary"
		case probe.Type == "system" && probe.Subtype == "compact_boundary":
			trigger = "claude.compact_boundary"
		}
		if trigger != "" {
			out = append(out, harness.TranscriptCompaction{Offset: rec.Offset, OccurredAt: rec.OccurredAt, Trigger: trigger})
		}
	}
	return out
}

// firstResults is each tool call's result, the first one the batch carries.
func firstResults(b harness.TranscriptBatch) map[string]harness.TranscriptToolResult {
	out := make(map[string]harness.TranscriptToolResult, len(b.ToolResults))
	for _, res := range b.ToolResults {
		if _, seen := out[res.CallID]; !seen {
			out[res.CallID] = res
		}
	}
	return out
}

// Telemetry implements TranscriptSource. Claude's transcript carries no
// native execution, process or timing record, so its telemetry is the
// normalised facts its ledger already verified, emitted again whenever the
// ledger moved.
func (s claudeTranscripts) Telemetry(req harness.TelemetryRequest) (harness.TelemetryUpdate, error) {
	up := harness.TelemetryUpdate{Offset: req.Offset, State: req.State, Gaps: []string{
		"child_usage_unavailable", "native_execution_unavailable", "native_process_unavailable", "native_timing_unavailable",
	}}
	if req.Offset != req.Batch.ConsumedBytes || req.Changed {
		up.Facts = NormalizedTelemetry(req.Batch)
		up.Offset = req.Batch.ConsumedBytes
	}
	return up, nil
}

// NormalizedTelemetry supplies Claude and older harnesses with the observations
// already verified by their transcript adapter. No native process/timing claims
// are made for this fallback.
func NormalizedTelemetry(tb harness.TranscriptBatch) []telemetry.Fact {
	var out []telemetry.Fact
	for _, t := range tb.Turns {
		off := t.Offset
		f := telemetry.Fact{Version: telemetry.Version, ID: t.SourceRef, Kind: "response", SourceRef: t.SourceRef, SourceOffset: off, OccurredAt: t.OccurredAt, MeasurementKind: "normalized", HarnessTurnRef: t.HarnessTurnRef, ResponseID: t.SourceRef, Model: t.Model, Text: t.Text, InputTokens: t.Usage.Input, CacheReadTokens: t.Usage.CacheRead, CacheWriteTokens: t.Usage.CacheWrite, OutputTokens: t.Usage.Output, ReasoningTokens: t.Usage.Reasoning, ContextTokens: t.Usage.ContextTokens(), LedgerRefOffset: &off}
		out = append(out, f)
	}
	calls := map[string]harness.TranscriptToolCall{}
	for _, c := range tb.ToolCalls {
		calls[c.SourceRef] = c
		at := c.StartedAt
		out = append(out, telemetry.Fact{Version: telemetry.Version, ID: c.SourceRef, Kind: "tool_call", SourceRef: c.SourceRef, SourceOffset: c.Offset, OccurredAt: at, MeasurementKind: "normalized", ResponseID: c.TurnSourceRef, WrapperRef: c.SourceRef, Tool: c.ToolName, Command: c.InputJSON, StartedAt: &at})
	}
	for _, r := range tb.ToolResults {
		at := r.CompletedAt
		c := calls[r.CallID]
		f := telemetry.Fact{Version: telemetry.Version, ID: r.CallID, Kind: "tool_result", SourceRef: r.SourceRef, SourceOffset: r.Offset, OccurredAt: at, MeasurementKind: "normalized", ResponseID: c.TurnSourceRef, WrapperRef: r.CallID, Tool: c.ToolName, Command: c.InputJSON, CompletedAt: &at, Output: harness.OutputFact(r.OutputText, nil)}
		if r.IsError {
			f.Status = "failed"
		}
		out = append(out, f)
	}
	for _, r := range tb.Records {
		if r.Kind != harness.TranscriptKindUser {
			continue
		}
		var p struct {
			Type     string `json:"type"`
			PromptID string `json:"promptId"`
			Message  struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(r.RawJSON), &p) != nil {
			continue
		}
		if bytes.Contains(p.Message.Content, []byte(`"tool_result"`)) {
			continue
		}
		text := harness.MessageText(p.Message.Content)
		if text != "" {
			out = append(out, telemetry.Fact{Version: telemetry.Version, ID: r.SourceRef, Kind: "prompt", SourceRef: r.SourceRef, SourceOffset: r.Offset, OccurredAt: r.OccurredAt, MeasurementKind: "normalized", HarnessTurnRef: p.PromptID, Text: text})
		}
	}
	return out
}
