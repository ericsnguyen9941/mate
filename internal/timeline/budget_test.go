package timeline_test

import (
	"context"
	"os"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// The fixture crew's turns total 23272+209152+0+1544 = 233968 tokens
// (timeline_test.go's codexTotal* constants: input net of cache, per
// codexTurns's doc comment), all on model gpt-5.6-terra.

// A crew_tokens ceiling below the fixture crew's total opens a budget
// incident on that crew, and CheckBudgets is idempotent: running it twice
// leaves exactly one open incident, never a second one (mvp.md section 4b -
// a budget incident is a fact, recorded once, never resolved).
func TestCheckBudgetsOpensAnIncidentOnceWhenCrewTokensIsCrossed(t *testing.T) {
	f := newFixture(t)
	f.ingest(t)

	if err := f.ws.SaveProject(fixtureProject, store.ProjectConfig{
		Repos:  []store.RepoConfig{{Path: fixtureProject + "/" + fixtureProject, DefaultBranch: "main"}},
		Budget: &store.BudgetConfig{CrewTokens: 200000},
	}); err != nil {
		t.Fatalf("SaveProject: %v", err)
	}

	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets (second run): %v", err)
	}

	view, err := box.Load(f.ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	open := box.OpenIncidents(view, fixtureCrew)
	if len(open) != 1 {
		t.Fatalf("open incidents for %s = %+v, want exactly 1", fixtureCrew, open)
	}
	if open[0].Kind != box.IncidentBudget {
		t.Fatalf("incident kind = %q, want %q", open[0].Kind, box.IncidentBudget)
	}
	if open[0].Text != "234k tokens of 200k tokens" {
		t.Fatalf("incident text = %q, want the humanised total-of-limit shape", open[0].Text)
	}
}

// A crew whose spend has not crossed crew_tokens gets no incident at all.
func TestCheckBudgetsOpensNothingUnderTheLimit(t *testing.T) {
	f := newFixture(t)
	f.ingest(t)

	if err := f.ws.SaveProject(fixtureProject, store.ProjectConfig{
		Repos:  []store.RepoConfig{{Path: fixtureProject + "/" + fixtureProject, DefaultBranch: "main"}},
		Budget: &store.BudgetConfig{CrewTokens: 10_000_000},
	}); err != nil {
		t.Fatalf("SaveProject: %v", err)
	}
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	view, err := box.Load(f.ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	if open := box.OpenIncidents(view, fixtureCrew); len(open) != 0 {
		t.Fatalf("open incidents = %+v, want none: the crew is under budget", open)
	}
}

// A crew_usd ceiling only fires once the model actually has a real price -
// the seeded placeholder (price 0) must never look like a crossing.
func TestCheckBudgetsCrewUSDIgnoresAnUnpricedModel(t *testing.T) {
	f := newFixture(t)
	f.ingest(t)

	if err := f.ws.SaveProject(fixtureProject, store.ProjectConfig{
		Repos:  []store.RepoConfig{{Path: fixtureProject + "/" + fixtureProject, DefaultBranch: "main"}},
		Budget: &store.BudgetConfig{CrewUSD: 0.01},
	}); err != nil {
		t.Fatalf("SaveProject: %v", err)
	}
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	view, err := box.Load(f.ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	if open := box.OpenIncidents(view, fixtureCrew); len(open) != 0 {
		t.Fatalf("open incidents = %+v, want none: the model is still unpriced", open)
	}

	// Now price it for real, above the ceiling, and re-ingest so `pricing`
	// picks it up.
	if err := f.ws.SavePricing(store.PricingConfig{Models: []store.PricingModel{
		{Model: "gpt-5.6-terra", InputPerM: 1000000},
	}}); err != nil {
		t.Fatalf("SavePricing: %v", err)
	}
	f.ingest(t)
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets after pricing: %v", err)
	}
	view, err = box.Load(f.ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	open := box.OpenIncidents(view, fixtureCrew)
	if len(open) != 1 || open[0].Kind != box.IncidentBudget {
		t.Fatalf("open incidents = %+v, want one budget incident now that the model is priced", open)
	}
}

// project_usd is checked against every crew of the project summed, and is
// filed under "mate" rather than under any one crew.
func TestCheckBudgetsProjectUSDIsFiledUnderMate(t *testing.T) {
	f := newFixture(t)
	if err := f.ws.SavePricing(store.PricingConfig{Models: []store.PricingModel{
		{Model: "gpt-5.6-terra", InputPerM: 1000000},
	}}); err != nil {
		t.Fatalf("SavePricing: %v", err)
	}
	f.ingest(t)

	if err := f.ws.SaveProject(fixtureProject, store.ProjectConfig{
		Repos:  []store.RepoConfig{{Path: fixtureProject + "/" + fixtureProject, DefaultBranch: "main"}},
		Budget: &store.BudgetConfig{ProjectUSD: 0.01},
	}); err != nil {
		t.Fatalf("SaveProject: %v", err)
	}
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}

	view, err := box.Load(f.ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	if open := box.OpenIncidents(view, fixtureCrew); len(open) != 0 {
		t.Fatalf("open incidents for the crew = %+v, want none: this is a project-wide ceiling", open)
	}
	open := box.OpenIncidents(view, "mate")
	if len(open) != 1 || open[0].Kind != box.IncidentBudget {
		t.Fatalf("open incidents for mate = %+v, want one budget incident", open)
	}
}

// A project with no `budget:` block at all is checked and finds nothing -
// CheckBudgets must not error just because most projects never configure one.
func TestCheckBudgetsNoBudgetConfiguredIsANoOp(t *testing.T) {
	f := newFixture(t)
	f.ingest(t)
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	view, err := box.Load(f.ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	if open := box.OpenIncidents(view, fixtureCrew); len(open) != 0 {
		t.Fatalf("open incidents = %+v, want none: no budget is configured", open)
	}
}

// crewOnlyWorking rewrites the fixture crew's status file so it has never
// handed back: the captured run ends in `wait-mate:`, which is exactly
// what the half-budget rule must not fire after.
func crewOnlyWorking(t *testing.T, f *fixture, lines ...string) {
	t.Helper()
	text := statusWorking + "\n"
	for _, l := range lines {
		text += l + "\n"
	}
	if err := os.WriteFile(f.ws.CrewStatus(fixtureProject, fixtureCrew), []byte(text), 0o644); err != nil {
		t.Fatalf("rewrite status: %v", err)
	}
}

func setCrewBudget(t *testing.T, f *fixture, budget string) {
	t.Helper()
	meta, err := f.ws.ReadCrewMeta(fixtureProject, fixtureCrew)
	if err != nil {
		t.Fatal(err)
	}
	meta[spawn.MetaBudget] = budget
	if err := f.ws.WriteCrewMeta(fixtureProject, fixtureCrew, meta); err != nil {
		t.Fatal(err)
	}
}

func openBudgetTexts(t *testing.T, f *fixture) []string {
	t.Helper()
	view, err := box.Load(f.ws, fixtureProject)
	if err != nil {
		t.Fatalf("box.Load: %v", err)
	}
	var out []string
	for _, inc := range box.OpenIncidents(view, fixtureCrew) {
		if inc.Kind != box.IncidentBudget {
			t.Fatalf("incident kind = %q, want %q", inc.Kind, box.IncidentBudget)
		}
		out = append(out, inc.Text)
	}
	return out
}

// A crew past half its own `budget=` (CON-76) with no hand-back yet gets
// one half-budget incident, marked by box.BudgetHalfPrefix, and never a
// second one for the same threshold however many polls follow. The
// fixture crew's 234k tokens are past half of 400k.
func TestCheckBudgetsFilesTheHalfIncidentOnceBeforeAHandBack(t *testing.T) {
	f := newFixture(t)
	crewOnlyWorking(t, f)
	setCrewBudget(t, f, "400000")
	f.ingest(t)

	for range 2 {
		if err := f.ing.CheckBudgets(context.Background()); err != nil {
			t.Fatalf("CheckBudgets: %v", err)
		}
	}
	got := openBudgetTexts(t, f)
	if len(got) != 1 || got[0] != "50%: 234k tokens of 400k tokens, no hand-back yet" {
		t.Fatalf("open budget incidents = %q, want exactly the half one", got)
	}
	if percent, rest := box.BudgetThreshold(got[0]); percent != 50 || rest != "234k tokens of 400k tokens, no hand-back yet" {
		t.Fatalf("BudgetThreshold = %d %q", percent, rest)
	}
}

// Once the crew has handed back - `pr-open:` or `wait-mate:` - half the
// budget is not worth an incident: the Mate is already looking at it. The
// captured run ends in `wait-mate:`; a `pr-open:` reads the same.
func TestCheckBudgetsFilesNoHalfIncidentAfterAHandBack(t *testing.T) {
	for name, lines := range map[string][]string{
		"wait-mate": nil, // the fixture's own status file, ending in wait-mate
		"pr-open":   {"pr-open: https://github.com/shop/shop/pull/7"},
	} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			if lines != nil {
				crewOnlyWorking(t, f, lines...)
			}
			setCrewBudget(t, f, "400000")
			f.ingest(t)
			if err := f.ing.CheckBudgets(context.Background()); err != nil {
				t.Fatalf("CheckBudgets: %v", err)
			}
			if got := openBudgetTexts(t, f); len(got) != 0 {
				t.Fatalf("open budget incidents = %q, want none after a hand-back", got)
			}
		})
	}
}

// The two thresholds are filed independently: a crew that already carries
// the half incident still gets the full one when it crosses its budget,
// and neither repeats. The budget is lowered between polls to stand in for
// the spend rising.
func TestCheckBudgetsFilesHalfThenFullOnceEach(t *testing.T) {
	f := newFixture(t)
	crewOnlyWorking(t, f)
	setCrewBudget(t, f, "400000")
	f.ingest(t)
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	setCrewBudget(t, f, "200000")
	for range 2 {
		if err := f.ing.CheckBudgets(context.Background()); err != nil {
			t.Fatalf("CheckBudgets: %v", err)
		}
	}
	got := openBudgetTexts(t, f)
	want := []string{"50%: 234k tokens of 400k tokens, no hand-back yet", "234k tokens of 200k tokens"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("open budget incidents = %q, want %q", got, want)
	}
}

// A crew's own `budget=` is what it is held to, before the project's
// crew_tokens: a project ceiling the crew is past opens nothing when the
// crew's rule sized it larger.
func TestCheckBudgetsPrefersTheCrewsOwnBudget(t *testing.T) {
	f := newFixture(t)
	crewOnlyWorking(t, f)
	setCrewBudget(t, f, "1000000")
	if err := f.ws.SaveProject(fixtureProject, store.ProjectConfig{
		Repos:  []store.RepoConfig{{Path: fixtureProject + "/" + fixtureProject, DefaultBranch: "main"}},
		Budget: &store.BudgetConfig{CrewTokens: 200000},
	}); err != nil {
		t.Fatalf("SaveProject: %v", err)
	}
	f.ingest(t)
	if err := f.ing.CheckBudgets(context.Background()); err != nil {
		t.Fatalf("CheckBudgets: %v", err)
	}
	if got := openBudgetTexts(t, f); len(got) != 0 {
		t.Fatalf("open budget incidents = %q, want none: 234k is under half of the crew's own 1M", got)
	}
}
