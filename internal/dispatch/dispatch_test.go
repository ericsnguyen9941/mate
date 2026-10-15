package dispatch

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/harness/claude"
	"github.com/nguyenngocanh94/mate/internal/harness/codex"
)

func writeTable(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "crew-dispatch.json")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// firstmate's schema: rules with a natural-language `when`, a `use` that is
// one profile or an array of alternatives, an optional `why`, and an
// optional `default` of the same shape.
func TestLoadReadsFirstmatesSchema(t *testing.T) {
	p := writeTable(t, `{
  "rules": [
    {"when": "a trivial mechanical edit", "use": {"harness": "claude", "model": "haiku", "effort": "low"}, "why": "cheap and fast"},
    {"when": "a big or ambiguous feature", "use": [
      {"harness": "claude", "model": "opus", "effort": "high"},
      {"harness": "codex", "model": "gpt-5.5", "effort": "xhigh"}
    ]}
  ],
  "default": {"harness": "codex", "effort": "medium"}
}`)
	tbl, ok, err := Load(p, catalog.Default())
	if err != nil || !ok {
		t.Fatalf("Load: ok=%v err=%v", ok, err)
	}
	if len(tbl.Rules) != 2 || tbl.Rules[0].Why != "cheap and fast" || len(tbl.Rules[1].Use) != 2 {
		t.Fatalf("rules = %+v", tbl.Rules)
	}
	if got := tbl.Rules[0].Use[0]; got != (Profile{Harness: claude.KindClaude, Model: "haiku", Effort: harness.EffortLow}) {
		t.Fatalf("first profile = %+v", got)
	}
	if len(tbl.Default) != 1 || tbl.Default[0].Harness != codex.KindCodex || tbl.Default[0].Model != "" {
		t.Fatalf("default = %+v", tbl.Default)
	}
}

// No file is no table: spawns go on as before. It is not an error.
func TestAMissingTableIsNotAnError(t *testing.T) {
	_, ok, err := Load(filepath.Join(t.TempDir(), "crew-dispatch.json"), catalog.Default())
	if ok || err != nil {
		t.Fatalf("Load of a missing file: ok=%v err=%v", ok, err)
	}
}

// A malformed table is reported and corrected, never selected around.
func TestLoadRefusesAMalformedTable(t *testing.T) {
	for name, body := range map[string]string{
		"not json":           `{"rules": [`,
		"rule without when":  `{"rules": [{"use": {"harness": "codex"}}]}`,
		"rule without use":   `{"rules": [{"when": "x"}]}`,
		"empty use array":    `{"rules": [{"when": "x", "use": []}]}`,
		"profile no harness": `{"rules": [{"when": "x", "use": {"model": "opus"}}]}`,
		"unknown harness":    `{"rules": [{"when": "x", "use": {"harness": "gemini"}}]}`,
		"bad effort":         `{"rules": [{"when": "x", "use": {"harness": "claude", "effort": "hgih"}}]}`,
		"unsupported effort": `{"rules": [{"when": "x", "use": {"harness": "codex", "effort": "max"}}]}`,
		"flag-like model":    `{"rules": [{"when": "x", "use": {"harness": "codex", "model": "-c"}}]}`,
		"unknown field":      `{"rules": [{"when": "x", "use": {"harness": "codex", "modle": "gpt-5.5"}}]}`,
		"empty table":        `{}`,
		"zero budget":        `{"rules": [{"when": "x", "use": {"harness": "codex"}, "budget": 0}]}`,
		"negative budget":    `{"rules": [{"when": "x", "use": {"harness": "codex"}, "budget": -5}]}`,
		"budget as a string": `{"rules": [{"when": "x", "use": {"harness": "codex"}, "budget": "2M"}]}`,
		"bad review budget":  `{"default": {"harness": "codex"}, "review_budget": 0}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, _, err := Load(writeTable(t, body), catalog.Default())
			if !errors.Is(err, ErrInvalid) {
				t.Fatalf("Load accepted a malformed table (err=%v)", err)
			}
			if !strings.Contains(err.Error(), "crew-dispatch.json") {
				t.Fatalf("the error does not name the file: %v", err)
			}
		})
	}
}

// A rule may carry a token budget (CON-76), and so may the review line as
// `review_budget`; a rule without one reads as 0, which spawn and the
// observer take as "the project's budget, if any".
func TestLoadReadsAnOptionalBudgetPerRule(t *testing.T) {
	p := writeTable(t, `{
  "rules": [
    {"when": "a small change", "use": {"harness": "claude", "model": "sonnet"}, "budget": 2000000},
    {"when": "anything else", "use": {"harness": "codex"}}
  ],
  "review": {"harness": "claude", "model": "sonnet", "effort": "medium"},
  "review_budget": 3000000
}`)
	tbl, ok, err := Load(p, catalog.Default())
	if err != nil || !ok {
		t.Fatalf("Load: ok=%v err=%v", ok, err)
	}
	if got := tbl.Rules[0].Budget; got != 2_000_000 {
		t.Fatalf("rule 1 budget = %d, want 2000000", got)
	}
	if got := tbl.Rules[1].Budget; got != 0 {
		t.Fatalf("rule 2 budget = %d, want 0 (none named)", got)
	}
	if tbl.ReviewBudget != 3_000_000 {
		t.Fatalf("review budget = %d, want 3000000", tbl.ReviewBudget)
	}

	// A table that names no review_budget has none.
	tbl, _, err = Load(writeTable(t, `{"default": {"harness": "codex"}}`), catalog.Default())
	if err != nil || tbl.ReviewBudget != 0 {
		t.Fatalf("table without review_budget: %d %v", tbl.ReviewBudget, err)
	}
}

func TestProfileFlags(t *testing.T) {
	p := Profile{Harness: claude.KindClaude, Model: "haiku", Effort: harness.EffortLow}
	if got := p.Flags(); got != "--harness claude --model haiku --effort low" {
		t.Fatalf("flags = %q", got)
	}
	if got := (Profile{Harness: codex.KindCodex}).Flags(); got != "--harness codex" {
		t.Fatalf("flags = %q", got)
	}
}

// The built-in table is the captain's five rules, in their order, each with
// a Claude and a Codex profile of the same weight.
func TestTheBuiltInTableIsTheCaptains(t *testing.T) {
	tbl, err := Resolve(t.TempDir(), catalog.Default())
	if err != nil || !tbl.BuiltIn || tbl.Path != "" {
		t.Fatalf("Resolve without a file: %+v %v", tbl, err)
	}
	flags := func(ps []Profile) []string {
		out := make([]string, len(ps))
		for i, p := range ps {
			out[i] = p.Flags()
		}
		return out
	}
	want := [][]string{
		{"--harness claude --model sonnet --effort medium", "--harness codex --model gpt-6-luna --effort medium",
			"--harness pi --model deepseek/deepseek-flash --effort high"},
		{"--harness claude --model sonnet --effort high", "--harness codex --model gpt-6-luna --effort high",
			"--harness claude --model opus --effort medium", "--harness codex --model gpt-5.6-terra --effort medium"},
		{"--harness claude --model opus --effort high", "--harness codex --model gpt-5.6-terra --effort high"},
		{"--harness claude --model opus --effort high", "--harness codex --model gpt-6-sol --effort high"},
		{"--harness claude --model sonnet --effort medium", "--harness codex --model gpt-6-luna --effort medium"},
		{"--harness claude --model opus --effort medium", "--harness codex --model gpt-5.6-terra --effort medium"},
	}
	if len(tbl.Rules) != len(want) {
		t.Fatalf("%d rules, want %d", len(tbl.Rules), len(want))
	}
	for i, r := range tbl.Rules {
		if got := flags(r.Use); strings.Join(got, "|") != strings.Join(want[i], "|") {
			t.Errorf("rule %d (%s) = %v, want %v", i+1, r.When, got, want[i])
		}
	}
	for i, prefix := range []string{"A ship that is a small", "A ship that is a medium", "A ship that is a large", "A scout that", "A light scout", "A review of another Crew's branch whose own task was a large change"} {
		if !strings.HasPrefix(tbl.Rules[i].When, prefix) {
			t.Errorf("rule %d = %q, want it to start %q", i+1, tbl.Rules[i].When, prefix)
		}
	}
	if got := strings.Join(flags(tbl.Review), "|"); got != "--harness claude --model sonnet --effort medium|--harness codex --model gpt-6-luna --effort medium" {
		t.Errorf("review = %s, want sonnet and luna at medium", got)
	}
	// The captain's suggested budgets (CON-76): the three ship rules and
	// the review line; the scouts and the large-change review name none.
	for i, budget := range []int64{2_000_000, 8_000_000, 20_000_000, 0, 0, 0} {
		if tbl.Rules[i].Budget != budget {
			t.Errorf("rule %d budget = %d, want %d", i+1, tbl.Rules[i].Budget, budget)
		}
	}
	if tbl.ReviewBudget != 3_000_000 {
		t.Errorf("review budget = %d, want 3000000", tbl.ReviewBudget)
	}
}

// A review profile is read like the default and checked like every
// profile; a table may name none.
func TestReviewIsReadAndChecked(t *testing.T) {
	tbl, _, err := Load(writeTable(t, `{"default": {"harness": "codex"}, "review": {"harness": "claude", "model": "haiku", "effort": "low"}}`), catalog.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(tbl.Review) != 1 || tbl.Review[0].Flags() != "--harness claude --model haiku --effort low" {
		t.Fatalf("review = %+v", tbl.Review)
	}
	if _, _, err := Load(writeTable(t, `{"default": {"harness": "codex"}, "review": {"harness": "nope"}}`), catalog.Default()); !errors.Is(err, ErrInvalid) {
		t.Fatalf("an unknown review harness: err = %v, want ErrInvalid", err)
	}
	none, _, err := Load(writeTable(t, `{"default": {"harness": "codex"}}`), catalog.Default())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := none.ReviewFor("codex", false); ok {
		t.Error("a table with no review profile returned one")
	}
}

// ReviewFor answers a named harness only with a profile on it, and the
// workspace's default harness with the first alternative when none runs
// on it.
func TestReviewForFollowsTheReviewersHarness(t *testing.T) {
	tbl, _ := Resolve(t.TempDir(), catalog.Default())
	for _, c := range []struct {
		kind  string
		named bool
		want  string
	}{
		{"codex", true, "--harness codex --model gpt-6-luna --effort medium"},
		{"claude", false, "--harness claude --model sonnet --effort medium"},
		{"pi", false, "--harness claude --model sonnet --effort medium"},
		{"pi", true, ""},
	} {
		p, ok := tbl.ReviewFor(c.kind, c.named)
		got := ""
		if ok {
			got = p.Flags()
		}
		if got != c.want {
			t.Errorf("ReviewFor(%s, %v) = %q, want %q", c.kind, c.named, got, c.want)
		}
	}
}

// A harness is known when it is registered: the same table is refused by a
// registry that lacks one of its harnesses, the built-in one included.
func TestTablesAreCheckedAgainstTheRegistry(t *testing.T) {
	onlyCodex, err := harness.NewRegistry(nil, codex.Codex{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(t.TempDir(), onlyCodex); !errors.Is(err, ErrInvalid) {
		t.Fatalf("built-in table naming claude under a registry without it: err = %v, want ErrInvalid", err)
	}
	p := writeTable(t, `{"default": {"harness": "claude"}}`)
	if _, _, err := Load(p, onlyCodex); !errors.Is(err, ErrInvalid) {
		t.Fatalf("table naming claude under a registry without it: err = %v, want ErrInvalid", err)
	}
	if _, _, err := Load(p, catalog.Default()); err != nil {
		t.Fatalf("the same table under the catalog: %v", err)
	}
}

// A workspace file replaces the built-in table whole; a malformed one is
// an error, not a fall back.
func TestAWorkspaceTableReplacesTheBuiltIn(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(Path(dir), []byte(`{"default": {"harness": "claude"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	tbl, err := Resolve(dir, catalog.Default())
	if err != nil || tbl.BuiltIn || len(tbl.Rules) != 0 || tbl.Path != Path(dir) {
		t.Fatalf("Resolve = %+v, %v", tbl, err)
	}
	if err := os.WriteFile(Path(dir), []byte(`{"rules": []}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(dir, catalog.Default()); !errors.Is(err, ErrInvalid) {
		t.Fatalf("malformed table: err = %v, want ErrInvalid", err)
	}
}

func TestDefaultForPrefersTheWorkspaceHarness(t *testing.T) {
	tbl, _ := Resolve(t.TempDir(), catalog.Default())
	for kind, want := range map[string]string{
		"codex":  "--harness codex --model gpt-6-luna --effort high",
		"claude": "--harness claude --model sonnet --effort high",
		"other":  "--harness codex --model gpt-6-luna --effort high",
	} {
		p, ok := tbl.DefaultFor(kind)
		if !ok || p.Flags() != want {
			t.Errorf("DefaultFor(%s) = %s, want %s", kind, p.Flags(), want)
		}
	}
	if _, ok := (Table{}).DefaultFor("codex"); ok {
		t.Error("a table with no default returned one")
	}
}

// A table may carry a harness rotation (CON-77) at table level or per
// rule; the rule's is its own, and a rule without one reads as nil so the
// table's applies. Entries must be registered harnesses.
func TestLoadReadsAnOptionalRotation(t *testing.T) {
	p := writeTable(t, `{
  "rules": [
    {"when": "a small change", "use": {"harness": "claude", "model": "sonnet"}, "rotate": ["claude", "codex"]},
    {"when": "anything else", "use": {"harness": "codex"}}
  ],
  "default": {"harness": "codex"},
  "rotate": ["Claude", "claude", "claude", "codex", "pi"]
}`)
	tbl, ok, err := Load(p, catalog.Default())
	if err != nil || !ok {
		t.Fatalf("Load: ok=%v err=%v", ok, err)
	}
	if got := tbl.Rotate.String(); got != "claude, claude, claude, codex, pi" {
		t.Fatalf("table rotation = %q", got)
	}
	if got := tbl.Rules[0].Rotate.String(); got != "claude, codex" {
		t.Fatalf("rule 1 rotation = %q", got)
	}
	if tbl.Rules[1].Rotate != nil {
		t.Fatalf("rule 2 rotation = %v, want none", tbl.Rules[1].Rotate)
	}
	if rs := tbl.Rotations(); len(rs) != 2 || rs[0].String() != "claude, codex" || rs[1].String() != "claude, claude, claude, codex, pi" {
		t.Fatalf("rotations in force = %v", rs)
	}

	for name, body := range map[string]string{
		"empty rotation":   `{"default": {"harness": "codex"}, "rotate": []}`,
		"unknown harness":  `{"default": {"harness": "codex"}, "rotate": ["claude", "gemini"]}`,
		"not an array":     `{"default": {"harness": "codex"}, "rotate": "claude"}`,
		"rule unknown":     `{"rules": [{"when": "x", "use": {"harness": "codex"}, "rotate": ["gemini"]}]}`,
		"rule empty array": `{"rules": [{"when": "x", "use": {"harness": "codex"}, "rotate": []}]}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, _, err := Load(writeTable(t, body), catalog.Default())
			if !errors.Is(err, ErrInvalid) || !strings.Contains(err.Error(), "rotate") {
				t.Fatalf("Load accepted a bad rotation, or did not name `rotate`: %v", err)
			}
		})
	}
	// The built-in table has no rotation: the captain opts in.
	if tbl, err := builtIn(catalog.Default()); err != nil || tbl.Rotate != nil || len(tbl.Rotations()) != 0 {
		t.Fatalf("built-in rotation = %v, %v; want none", tbl.Rotate, err)
	}
}

// Next walks the rotation from the slots already taken, and `--skip`
// passes over a harness without deferring its turn.
func TestRotationNext(t *testing.T) {
	r := Rotation{claude.KindClaude, claude.KindClaude, claude.KindClaude, codex.KindCodex, "pi"}
	for _, tc := range []struct {
		consumed int
		skip     []harness.Kind
		want     Step
		skipped  string
		ok       bool
	}{
		{0, nil, Step{claude.KindClaude, 1, 1}, "", true},
		{1, nil, Step{claude.KindClaude, 2, 2}, "", true},
		{3, nil, Step{codex.KindCodex, 4, 4}, "", true},
		{4, nil, Step{"pi", 5, 5}, "", true},
		{5, nil, Step{claude.KindClaude, 1, 6}, "", true},
		{3, []harness.Kind{codex.KindCodex}, Step{"pi", 5, 5}, "codex", true},
		{0, []harness.Kind{claude.KindClaude}, Step{codex.KindCodex, 4, 4}, "claude", true},
		{4, []harness.Kind{"pi", claude.KindClaude}, Step{codex.KindCodex, 4, 9}, "pi, claude", true},
		{0, []harness.Kind{claude.KindClaude, codex.KindCodex, "pi"}, Step{}, "claude, codex, pi", false},
	} {
		step, skipped, ok := r.Next(tc.consumed, tc.skip)
		if ok != tc.ok || step != tc.want || Rotation(skipped).String() != tc.skipped {
			t.Errorf("Next(%d, %v) = %+v skipped %q ok %v; want %+v %q %v", tc.consumed, tc.skip, step, Rotation(skipped).String(), ok, tc.want, tc.skipped, tc.ok)
		}
	}
	if _, _, ok := Rotation(nil).Next(0, nil); ok {
		t.Fatal("an empty rotation has a next")
	}
}

// A `review_small` profile and its threshold are read like `review` and
// checked like every profile (CON-74); a table may name neither, and then
// the threshold is the default and the small profile is absent.
func TestReviewSmallIsReadAndChecked(t *testing.T) {
	tbl, _, err := Load(writeTable(t, `{"default": {"harness": "codex"}, "review_small": {"harness": "claude", "model": "haiku", "effort": "low"}, "review_small_max_lines": 20}`), catalog.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(tbl.ReviewSmall) != 1 || tbl.ReviewSmall[0].Flags() != "--harness claude --model haiku --effort low" || tbl.ReviewSmallMaxLines != 20 {
		t.Fatalf("review_small = %+v, max lines %d", tbl.ReviewSmall, tbl.ReviewSmallMaxLines)
	}
	for name, body := range map[string]string{
		"unknown harness": `{"default": {"harness": "codex"}, "review_small": {"harness": "nope"}}`,
		"zero threshold":  `{"default": {"harness": "codex"}, "review_small_max_lines": 0}`,
		"string":          `{"default": {"harness": "codex"}, "review_small_max_lines": "50"}`,
	} {
		if _, _, err := Load(writeTable(t, body), catalog.Default()); !errors.Is(err, ErrInvalid) {
			t.Errorf("%s: err = %v, want ErrInvalid", name, err)
		}
	}
	none, _, err := Load(writeTable(t, `{"default": {"harness": "codex"}}`), catalog.Default())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := none.ReviewSmallFor("codex", false); ok || none.ReviewSmallMaxLines != DefaultReviewSmallMaxLines {
		t.Errorf("a table with no review_small: profile %v, max lines %d", ok, none.ReviewSmallMaxLines)
	}
	// The built-in table reviews a small change on the lightest profile of
	// each harness, under the default threshold.
	builtin, _ := Resolve(t.TempDir(), catalog.Default())
	if p, ok := builtin.ReviewSmallFor("codex", true); !ok || p.Flags() != "--harness codex --model gpt-6-luna --effort low" {
		t.Errorf("built-in review_small for codex = %+v %v", p, ok)
	}
	if p, ok := builtin.ReviewSmallFor("pi", false); !ok || p.Flags() != "--harness claude --model haiku --effort low" {
		t.Errorf("built-in review_small for an unlisted default harness = %+v %v", p, ok)
	}
	if builtin.ReviewSmallMaxLines != 50 {
		t.Errorf("built-in review_small_max_lines = %d, want 50", builtin.ReviewSmallMaxLines)
	}
}

// IsCodePath is the conservative list the small-change rule relies on:
// languages, shell and build files are code; prose, data, configuration
// and markup are not.
func TestIsCodePath(t *testing.T) {
	for path, want := range map[string]bool{
		"cmd/mate/main.go": true, "web/app.js": true, "lib/x.rb": true, "scripts/run.sh": true, "Makefile": true,
		"deploy/Dockerfile": true, "go.mod": true, "db/schema.sql": true, "src/App.tsx": true,
		"README.md": false, "docs/mvp.md": false, "config.yaml": false, "package.json": false,
		"site/index.html": false, "site/style.css": false, "LICENSE": false, "data.csv": false,
		"assets/mate/skills/x/SKILL.md.tmpl": false,
	} {
		if got := IsCodePath(path); got != want {
			t.Errorf("IsCodePath(%q) = %v, want %v", path, got, want)
		}
	}
}
