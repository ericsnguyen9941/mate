package timeline_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A crew that delegates to a kind of agent leaves it in its tool.called
// event, named as `agent`, and the story says which kind; a page badges the
// crew with it (CON-124). The transcript is a Claude one in which the crew
// hands a review to its project's backend-reviewer and then speaks again, so
// the delegating message group is closed and parsed.
func TestADelegationNamesTheKindOfAgentItWentTo(t *testing.T) {
	f := newFixture(t)
	rec := func(v map[string]any) string {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	usage := map[string]any{"input_tokens": 5, "cache_creation_input_tokens": 100, "cache_read_input_tokens": 200, "output_tokens": 40}
	lines := []string{
		rec(map[string]any{"type": "user", "uuid": "u1", "parentUuid": nil, "isSidechain": false, "promptId": "p1",
			"timestamp": "2026-09-19T10:44:20.000Z", "sessionId": "sess-crew",
			"message": map[string]any{"role": "user", "content": "Add the Buy button, then have the backend reviewed."}}),
		rec(map[string]any{"type": "assistant", "uuid": "a1", "parentUuid": "u1", "isSidechain": false,
			"timestamp": "2026-09-19T10:44:25.000Z", "sessionId": "sess-crew",
			"message": map[string]any{"id": "msg_1", "role": "assistant", "model": "claude-fable-5-1", "type": "message",
				"stop_reason": "tool_use", "usage": usage,
				"content": []map[string]any{{"type": "tool_use", "id": "toolu_1", "name": "Agent",
					"input": map[string]any{"description": "Review the backend diff", "prompt": "Review backend/ against AGENTS.md", "subagent_type": "backend-reviewer"}}}}}),
		rec(map[string]any{"type": "user", "uuid": "u2", "parentUuid": "a1", "isSidechain": false,
			"timestamp": "2026-09-19T10:44:55.000Z", "sessionId": "sess-crew",
			"message": map[string]any{"role": "user", "content": []map[string]any{{"type": "tool_result", "tool_use_id": "toolu_1", "content": "No findings.", "is_error": false}}}}),
		rec(map[string]any{"type": "assistant", "uuid": "a2", "parentUuid": "u2", "isSidechain": false,
			"timestamp": "2026-09-19T10:45:00.000Z", "sessionId": "sess-crew",
			"message": map[string]any{"id": "msg_2", "role": "assistant", "model": "claude-fable-5-1", "type": "message",
				"stop_reason": "end_turn", "usage": usage,
				"content": []map[string]any{{"type": "text", "text": "The backend reviewer found nothing."}}}}),
	}
	path := filepath.Join(f.root, "crew-claude.jsonl")
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	meta, err := f.ws.ReadCrewMeta(fixtureProject, fixtureCrew)
	if err != nil {
		t.Fatal(err)
	}
	meta["harness"], meta["transcript"], meta["session_id"] = "claude", path, "sess-crew"
	if err := f.ws.WriteCrewMeta(fixtureProject, fixtureCrew, meta); err != nil {
		t.Fatal(err)
	}
	f.ingest(t)

	var agent, target string
	err = f.db.SQL().QueryRow(`SELECT json_extract(payload, '$.agent'), json_extract(payload, '$.target') FROM event
		WHERE actor_id = ? AND kind = 'tool.called' AND json_extract(payload, '$.class') = 'agent'`, f.crewActor()).Scan(&agent, &target)
	if err != nil {
		t.Fatalf("the delegation's tool.called: %v", err)
	}
	if agent != "backend-reviewer" || target != "backend-reviewer" {
		t.Fatalf("agent %q target %q, want the kind of agent in both", agent, target)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM action WHERE actor_id = ? AND tool = 'Agent' AND target = 'backend-reviewer'`, f.crewActor()); n != 1 {
		t.Fatalf("%d action row(s) targeting backend-reviewer, want 1", n)
	}
	if story := f.narrate(t); !strings.Contains(story, "buybtn delegates to a backend-reviewer agent") {
		t.Fatalf("the story does not say which kind of agent:\n%s", story)
	}
}
