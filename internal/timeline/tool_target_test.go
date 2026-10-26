package timeline

import (
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// A Skill call is about the skill it loads. Without the name as its target
// the action row says only "Skill", and a reader cannot tell which one.
func TestSkillCallTargetIsTheSkillItLoads(t *testing.T) {
	call := harness.TranscriptToolCall{ToolName: "Skill", CommandClass: harness.CommandOther, InputJSON: `{"skill":"token-review","args":"shop"}`}
	target, summary, class := toolTarget(harness.Kind("claude"), call)
	if target != "token-review" {
		t.Fatalf("target %q, want the skill's name", target)
	}
	if summary != `{"skill":"token-review","args":"shop"}` || class != harness.CommandOther {
		t.Fatalf("summary %q class %q: the input and the class are the call's own", summary, class)
	}
}

// A delegation is about the kind of agent it went to. Claude's Agent tool
// names it as subagent_type; without the kind as its target the action row
// says only "Agent" and a description, and a page cannot badge the crew with
// the backend-reviewer it used (CON-124).
func TestAgentCallTargetIsTheKindOfAgent(t *testing.T) {
	input := `{"description":"Review the backend diff","prompt":"Review backend/ against AGENTS.md","subagent_type":"backend-reviewer"}`
	call := harness.TranscriptToolCall{ToolName: "Agent", CommandClass: harness.CommandAgent, InputJSON: input}
	target, summary, class := toolTarget(harness.Kind("claude"), call)
	if target != "backend-reviewer" {
		t.Fatalf("target %q, want the kind of agent", target)
	}
	if summary != input || class != harness.CommandAgent {
		t.Fatalf("summary %q class %q: the input and the class are the call's own", summary, class)
	}
	if got := agentType(call); got != "backend-reviewer" {
		t.Fatalf("agentType %q, want backend-reviewer", got)
	}
}

// A delegation whose input names no kind keeps its description as the
// target, as before, and has no agent: a description must not pass for a
// kind of agent on a badge.
func TestAgentCallWithoutAKindKeepsItsDescription(t *testing.T) {
	call := harness.TranscriptToolCall{ToolName: "Task", CommandClass: harness.CommandAgent,
		InputJSON: `{"description":"Find the checkout page","prompt":"Where is the checkout page?"}`}
	target, _, _ := toolTarget(harness.Kind("claude"), call)
	if target != "Find the checkout page" {
		t.Fatalf("target %q, want the description", target)
	}
	if got := agentType(call); got != "" {
		t.Fatalf("agentType %q, want none", got)
	}
	shell := harness.TranscriptToolCall{ToolName: "Bash", CommandClass: harness.CommandShell, InputJSON: `{"command":"ls","subagent_type":"x"}`}
	if got := agentType(shell); got != "" {
		t.Fatalf("agentType of a shell call %q, want none: only a delegation has one", got)
	}
}
