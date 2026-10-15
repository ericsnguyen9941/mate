package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// usageWorkspace builds a workspace with a Mate, one crew declared
// `working`, and a hand-seeded ledger: the crew's turns total
// 90000+6000+300+1200 = 97500 tokens on a model that is not yet priced. The
// Mate gets a smaller total on the same unpriced model.
func usageWorkspace(t *testing.T) *store.Workspace {
	t.Helper()
	w := liveCrewWorkspace(t, "shop")
	if err := w.WriteMateMeta("shop", map[string]string{"harness": "claude"}); err != nil {
		t.Fatalf("WriteMateMeta: %v", err)
	}
	if err := w.WriteCrewMeta("shop", "k3", map[string]string{
		"harness": "codex", "task": "add a Buy button", "state": "spawned",
	}); err != nil {
		t.Fatalf("WriteCrewMeta: %v", err)
	}
	if err := w.AppendStatus("shop", "k3", "working: writing the button"); err != nil {
		t.Fatalf("AppendStatus: %v", err)
	}

	handle, err := db.Open(w)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer handle.Close()

	crewActor := timeline.CrewActorID("shop", "k3")
	mateActor := timeline.MateActorID("shop")
	// seedActorTurnAndPricing writes the crew's `task` row itself (it must,
	// since `v_task_ledger` is a join against `task`); the Mate gets none,
	// which matches production - a Mate's turns belong to the project, not
	// to any one task (docs/timeline.md).
	seedActorTurnAndPricing(t, handle, crewActor, "shop", "crew", "unpriced-model", 90_000, 6_000, 0, 1_200, 200_000)
	seedActorTurnAndPricing(t, handle, mateActor, "shop", "mate", "unpriced-model", 400, 0, 0, 100, 200_000)
	return w
}

func TestUsagePrintsTheLedgerWithMateFirstAndATotalsFooter(t *testing.T) {
	w := usageWorkspace(t)
	out := runCLI(t, "usage", "shop", "--workspace", w.Root())

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("usage printed %d line(s):\n%s", len(lines), out)
	}
	header := lines[0]
	for _, col := range []string{"ID", "STATE", "MODEL", "EFFORT", "TURNS", "IN", "CACHE-READ", "CACHE-WRITE", "OUT", "TOTAL", "COST", "CTX", "CTX-AVG", "CTX%", "ASKED", "WAITED"} {
		if !strings.Contains(header, col) {
			t.Fatalf("header %q is missing column %q", header, col)
		}
	}
	if !strings.Contains(lines[1], "mate") {
		t.Fatalf("mate row is not first:\n%s", out)
	}
	if !strings.Contains(lines[2], "k3") || !strings.Contains(lines[2], "working") {
		t.Fatalf("crew row missing or wrong state:\n%s", out)
	}
	// The model is unpriced: every cost cell, including the footer's, reads
	// "?" rather than a fabricated $0.00 (mvp.md M5: a missing price is not
	// a price of zero).
	if !strings.Contains(lines[1], "?") || !strings.Contains(lines[2], "?") {
		t.Fatalf("an unpriced model must show cost as ?:\n%s", out)
	}
	footer := lines[3]
	if !strings.Contains(footer, "total:") || !strings.Contains(footer, "? cost") {
		t.Fatalf("footer = %q, want a totals line with ? cost", footer)
	}
	// Numbers humanised: 90000+6000+0+1200 = 97200 -> 97.2k.
	if !strings.Contains(lines[2], "97.2k") {
		t.Fatalf("crew TOTAL is not humanised:\n%s", lines[2])
	}
}

func TestUsageShowsARealCostOnceTheModelIsPriced(t *testing.T) {
	w := usageWorkspace(t)
	if err := w.SavePricing(store.PricingConfig{Models: []store.PricingModel{
		{Model: "unpriced-model", InputPerM: 10, OutputPerM: 20, ContextWindow: 200_000},
	}}); err != nil {
		t.Fatalf("SavePricing: %v", err)
	}
	// Reload the price into the already-seeded pricing row: this test seeds
	// the ledger by hand rather than through the ingest, so it upserts the
	// price the same way `timeline.Ingester.ingestPricing` would.
	handle, err := db.Open(w)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	if _, err := handle.SQL().Exec(
		`UPDATE pricing SET input_per_m = 10, output_per_m = 20 WHERE model = 'unpriced-model'`); err != nil {
		handle.Close()
		t.Fatalf("price the model: %v", err)
	}
	handle.Close()

	out := runCLI(t, "usage", "shop", "--workspace", w.Root())
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	// Crew: 90000 input @ $10/M + 1200 output @ $20/M = 0.90 + 0.024 = $0.924.
	if !strings.Contains(lines[2], "$0.92") {
		t.Fatalf("crew row does not show the priced cost:\n%s", lines[2])
	}
	footer := lines[len(lines)-1]
	if strings.Contains(footer, "? cost") {
		t.Fatalf("footer still says unpriced once the model is priced:\n%s", footer)
	}
}

func TestUsageWithACrewArgumentPrintsPerTurnRows(t *testing.T) {
	w := usageWorkspace(t)
	out := runCLI(t, "usage", "shop", "k3", "--workspace", w.Root())
	if !strings.Contains(out, "TURN") || !strings.Contains(out, "MODEL") || !strings.Contains(out, "unpriced-model") {
		t.Fatalf("per-turn output missing expected columns/model:\n%s", out)
	}
	if strings.Contains(out, "mate") {
		t.Fatalf("per-turn output for k3 must not include the mate's turns:\n%s", out)
	}
}

func TestUsageUnknownProjectIsAUsageError(t *testing.T) {
	w := usageWorkspace(t)
	var stdout, stderr strings.Builder
	err := run([]string{"usage", "nope", "--workspace", w.Root()}, &stdout, &stderr)
	if err == nil {
		t.Fatal("usage on an unregistered project must fail")
	}
}

func TestUsageAndConsoleKeepAbsoluteContextWithoutAModelWindow(t *testing.T) {
	w := usageWorkspace(t)
	handle, err := db.Open(w)
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Close()
	if _, err := handle.SQL().Exec(`DELETE FROM pricing`); err != nil {
		t.Fatal(err)
	}
	row, err := mateLedgerRow(context.Background(), handle, w, "shop")
	if err != nil {
		t.Fatal(err)
	}
	if !row.Context.Valid || row.Context.Int64 != 400 || row.CtxPct.Valid {
		t.Fatalf("absolute context requires no pricing: %+v", row)
	}
	tokens, ok := mateTokens(handle, "shop")
	if !ok || tokens.ContextTokens == nil || *tokens.ContextTokens != 400 || tokens.ContextPct != nil {
		t.Fatalf("console lost absolute context: %+v", tokens)
	}
	var out strings.Builder
	printLedgerTable(&out, []ledgerRow{row})
	if !strings.Contains(out.String(), "CACHE-READ") || !strings.Contains(out.String(), "CACHE-WRITE") || !strings.Contains(out.String(), "CTX") {
		t.Fatalf("usage must distinguish cache buckets and context:\n%s", out.String())
	}
}

// Below the totals, one line per model gives calls per task and context
// per call over every task, and a Claude Crew whose context passed the
// workspace's crew autocompact threshold is named (CON-68). A Codex Crew
// above it is not: the threshold is Claude's flag.
func TestUsagePrintsCallsPerTaskAndWarnsAboveTheCompactThreshold(t *testing.T) {
	w := usageWorkspace(t)
	for crew, h := range map[string]string{"big": "claude", "wide": "codex"} {
		if err := w.WriteCrewMeta("shop", crew, map[string]string{"harness": h, "task": crew, "state": "spawned"}); err != nil {
			t.Fatal(err)
		}
	}
	handle, err := db.Open(w)
	if err != nil {
		t.Fatal(err)
	}
	seedActorTurnAndPricing(t, handle, timeline.CrewActorID("shop", "big"), "shop", "crew", "claude-opus-5-5", 1_000, 250_000, 0, 500, 1_000_000)
	seedActorTurnAndPricing(t, handle, timeline.CrewActorID("shop", "wide"), "shop", "crew", "gpt-6-luna", 1_000, 200_000, 0, 500, 272_000)
	handle.Close()

	out := runCLI(t, "usage", "shop", "--workspace", w.Root())
	for _, want := range []string{
		"by model and effort (each task's last call):",
		"  claude-opus-5-5: 1 task(s), 1.0 calls/task, 251k context/call",
		"  unpriced-model: 1 task(s), 1.0 calls/task, 96k context/call",
		"warning: 1 crew(s) passed the 160k crew autocompact threshold (crew.autocompact in .mate/workspace.yaml): big 251k\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("usage output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "wide 201k") {
		t.Errorf("a Codex Crew was warned about Claude's threshold:\n%s", out)
	}
}

// Context per call is a mean over the calls that recorded a context: a task
// whose other calls recorded none must not weigh its mean by them.
func TestUsageContextPerCallCountsOnlyCallsThatRecordedOne(t *testing.T) {
	w := usageWorkspace(t)
	for _, crew := range []string{"lean", "full"} {
		if err := w.WriteCrewMeta("shop", crew, map[string]string{"harness": "claude", "task": crew, "state": "spawned"}); err != nil {
			t.Fatal(err)
		}
	}
	handle, err := db.Open(w)
	if err != nil {
		t.Fatal(err)
	}
	seedActorTurnAndPricing(t, handle, timeline.CrewActorID("shop", "lean"), "shop", "crew", "claude-sonnet-5-5", 0, 100_000, 0, 0, 1_000_000)
	seedActorTurnAndPricing(t, handle, timeline.CrewActorID("shop", "full"), "shop", "crew", "claude-sonnet-5-5", 0, 300_000, 0, 0, 1_000_000)
	now := db.FormatTime(time.Now())
	for i := range 3 {
		a := timeline.CrewActorID("shop", "lean")
		if _, err := handle.SQL().Exec(
			`INSERT INTO turn(id, actor_id, session_id, started_at, ended_at, model, input_tokens, cache_read_tokens, cache_write_tokens, output_tokens, context_tokens_after)
			 VALUES (?, ?, ?, ?, ?, 'claude-sonnet-5-5', 0, 0, 0, 10, 0)`,
			fmt.Sprintf("%s#t%d", a, i), a, a+"#s", now, now); err != nil {
			t.Fatal(err)
		}
	}
	handle.Close()

	out := runCLI(t, "usage", "shop", "--workspace", w.Root())
	if want := "  claude-sonnet-5-5: 2 task(s), 2.5 calls/task, 200k context/call"; !strings.Contains(out, want) {
		t.Errorf("usage output lacks %q:\n%s", want, out)
	}
}

// The ledger names each row's model and effort (CON-75): the last call's
// when the turn recorded them, else the launch profile in `.meta` - which
// is where a Claude crew's effort always comes from, since its transcript
// never says. The per-model line keys on both, so Opus high and Sonnet
// medium are two lines.
func TestUsageShowsModelAndEffortFromTheTurnOrTheLaunchProfile(t *testing.T) {
	w := usageWorkspace(t)
	if err := w.UpdateCrewMeta("shop", "k3", map[string]string{"model": "claude-opus-5-5", "effort": "high"}); err != nil {
		t.Fatal(err)
	}
	out := runCLI(t, "usage", "shop", "--workspace", w.Root())
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	fields := strings.Fields(lines[2])
	if fields[0] != "k3" || fields[2] != "unpriced-model" || fields[3] != "high" {
		t.Fatalf("crew row = %q, want the turn's model and the meta's effort", lines[2])
	}
	if !strings.Contains(out, "unpriced-model high: 1 task(s)") {
		t.Fatalf("the per-model line does not key on effort:\n%s", out)
	}
	turns := runCLI(t, "usage", "shop", "k3", "--workspace", w.Root())
	if !strings.Contains(strings.Split(turns, "\n")[0], "EFFORT") {
		t.Fatalf("per-turn header lacks EFFORT:\n%s", turns)
	}
}
