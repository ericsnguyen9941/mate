package why

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/telemetry"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

func TestUsageWhyExplainsOneTask(t *testing.T) {
	w := usageWorkspace(t)
	handle, err := db.Open(w)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	actor := timeline.CrewActorID("shop", "k3")
	end := time.Now().UTC()
	start := end.Add(-10 * time.Second)
	exit := 1
	fact := telemetry.Fact{Version: 1, ID: "failed-test", ExecutionID: "failed-test", Kind: "execution", SessionID: actor + "#s",
		Command: "go test ./cart", Tool: "exec_command", StartedAt: &start, CompletedAt: &end, Status: "failed", ExitCode: &exit,
		OccurredAt: end, ObservedAt: end, MeasurementKind: "native", SourcePath: "rollout.jsonl", SourceOffset: 7,
		Output: &telemetry.Output{Bytes: 9, SHA256: "failed-output", Preview: "FAIL cart"}}
	payload, err := json.Marshal(fact)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := handle.SQL().Exec(`INSERT INTO event(dedup,project,at,actor_id,kind,payload,ref_path,ref_offset) VALUES(?,?,?,?,?,?,?,?)`,
		"why-test#failed-test", "shop", db.FormatTime(end), actor, "telemetry.execution", string(payload), fact.SourcePath, fact.SourceOffset); err != nil {
		handle.Close()
		t.Fatalf("seed telemetry: %v", err)
	}
	handle.Close()

	out := writePage(t, w, "shop", "k3")
	for _, want := range []string{
		"usage why shop k3",
		// The ledger's own numbers, bucket by bucket; an unpriced model is "?".
		"tokens: 97.2k total = 90k in + 6k cache-read + 0 cache-write + 1.2k out · cost ?",
		"1 model call(s)",
		"tool elapsed 10s",
		"tool invocation sum 10s",
		"unallocated is not model latency",
		"? is unknown",
		"context peak 96k, last 96k",
		// The failed command is named with its exit code and the review hint
		// that says where a fix would live.
		"execution_failure",
		"go test ./cart · exit 1",
		"review: ",
		// Findings overlap, and the digest has to say so where they are read.
		"never add them up",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("--why output is missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "TURN ") {
		t.Errorf("--why must replace the per-call table, not add to it:\n%s", out)
	}
}

func TestUsageWhyOnAnUnknownCrewNamesTheCrew(t *testing.T) {
	w := usageWorkspace(t)
	_, err := tryPage(w, "shop", "nope")
	if err == nil || !strings.Contains(err.Error(), `"nope"`) {
		t.Fatalf("want an error naming the crew, got %v", err)
	}
}

// A crew that has made no model call yet gets the page's head and a line
// saying there is nothing to explain, never an error or an empty page
// (CON-168, Trường hợp đặc biệt).
func TestUsageWhyOfACrewWithNoTurnsSaysSo(t *testing.T) {
	w := usageWorkspace(t)
	if err := w.WriteCrewMeta("shop", "idle", map[string]string{"harness": "codex", "task": "wait", "state": "spawned"}); err != nil {
		t.Fatal(err)
	}
	// The ingest has recorded the crew and its task, but no turn.
	handle, err := db.Open(w)
	if err != nil {
		t.Fatal(err)
	}
	actor := timeline.CrewActorID("shop", "idle")
	for _, stmt := range []string{
		`INSERT INTO actor(id, project, kind, name) VALUES ('` + actor + `', 'shop', 'crew', 'idle')`,
		`INSERT INTO task(crew_actor_id, project, spawned_at) VALUES ('` + actor + `', 'shop', '` + db.FormatTime(time.Now()) + `')`,
	} {
		if _, err := handle.SQL().Exec(stmt); err != nil {
			handle.Close()
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	handle.Close()
	why := writePage(t, w, "shop", "idle")
	if !strings.HasPrefix(why, "usage why shop idle\n") || !strings.HasSuffix(why, "\nno model calls recorded for idle; there is nothing to explain yet\n") {
		t.Fatalf("--why of a crew with no calls:\n%s", why)
	}
}
