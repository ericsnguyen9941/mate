// Package dispatch reads the captain's crew dispatch table,
// `.mate/crew-dispatch.json`: which harness, model and effort suit which
// kind of task. The schema is firstmate's (config/crew-dispatch.json):
//
//	{
//	  "rules": [
//	    {"when": "<a kind of task, in words>",
//	     "use": {"harness": "claude", "model": "haiku", "effort": "low"},
//	     "why": "<optional rationale>",
//	     "budget": 2000000}
//	  ],
//	  "default": {"harness": "codex", "model": "gpt-5.5", "effort": "medium"},
//	  "review": {"harness": "claude", "model": "sonnet", "effort": "medium"},
//	  "review_budget": 3000000,
//	  "rotate": ["claude", "claude", "claude", "codex", "pi"]
//	}
//
// `review` is mate's addition to firstmate's schema: the profile `mate
// review` launches a reviewer with when the Mate names no model or effort,
// so a review is not sized like the change it reviews (CON-68).
//
// `budget` and `review_budget` are mate's too (CON-76): the token budget a
// task of that rule's kind is sized at, which the Mate passes to `mate crew
// spawn --budget` and the observer checks the crew's spend against (half
// of it without a hand-back, and all of it). Both are optional; a rule
// without one leaves the crew on the project's `budget:` block, if any.
//
// `rotate` is mate's as well (CON-77): a round-robin of harnesses the
// captain wants spawns to take turns on, at table level or per rule (the
// rule's wins for a task of its kind). Each spawn that takes the harness
// the rotation names next records `rotation_slot=` in its meta, and `mate
// crew dispatch` reads those to say which harness is next (Rotation.Next).
// It is a suggestion the Mate reads, never a gate: the captain's explicit
// `--harness` wins and leaves the position where it was.
//
// `use`, `default` and `review` take one profile or an array of alternatives. The
// rules are natural language: the Mate matches a task to a rule with its
// own judgment and passes concrete --harness, --model and --effort to
// `mate crew spawn`. Nothing here matches a rule. This package only checks
// that every profile is one the app can launch, so a table with a typo is
// reported rather than silently launched around. A workspace without the
// file is governed by the built-in table (builtin.go).
package dispatch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// FileName is the table's name under `.mate/`.
const FileName = "crew-dispatch.json"

// ErrInvalid wraps every reason a table is refused.
var ErrInvalid = errors.New("invalid crew dispatch table")

// Profile is one concrete launch: a harness, and optionally a model and an
// effort. Empty model or effort is the harness's own default.
type Profile struct {
	Harness harness.Kind
	Model   string
	Effort  harness.Effort
}

// Flags is the profile as `mate crew spawn` flags.
func (p Profile) Flags() string {
	parts := []string{"--harness", string(p.Harness)}
	if p.Model != "" {
		parts = append(parts, "--model", p.Model)
	}
	if p.Effort != "" {
		parts = append(parts, "--effort", string(p.Effort))
	}
	return strings.Join(parts, " ")
}

// Rule is one kind of task and the profiles that suit it.
type Rule struct {
	When string
	Use  []Profile
	Why  string
	// Budget is the token budget a task of this kind is sized at, 0 when
	// the rule names none (the project's `budget:` then applies).
	Budget int64
	// Rotate is the rule's own harness rotation, nil when it has none and
	// the table's applies.
	Rotate Rotation
}

// Table is a loaded, checked dispatch table. Path is empty and BuiltIn true
// for the table compiled into mate (builtin.go).
type Table struct {
	Path    string
	BuiltIn bool
	Rules   []Rule
	Default []Profile
	// Review is the reviewer profile of `mate review`; empty when the
	// table names none.
	Review []Profile
	// ReviewBudget is the token budget a review is sized at, 0 when the
	// table names none.
	ReviewBudget int64
	// Rotate is the table's harness rotation, nil when it names none.
	Rotate Rotation
}

// Rotation is a round-robin of harnesses: position 1 is the first entry,
// and after the last the first comes round again. Entries may repeat
// (three turns of claude, then one of codex).
type Rotation []harness.Kind

// String lists the rotation the way `crew dispatch` prints it.
func (r Rotation) String() string {
	names := make([]string, len(r))
	for i, k := range r {
		names[i] = string(k)
	}
	return strings.Join(names, ", ")
}

// Equal reports whether two rotations name the same harnesses in the same
// order.
func (r Rotation) Equal(o Rotation) bool {
	if len(r) != len(o) {
		return false
	}
	for i := range r {
		if r[i] != o[i] {
			return false
		}
	}
	return true
}

// Step is the rotation's next turn: the harness, its position in the
// rotation (1..n) and the slot a spawn taking it records. Slot is the
// running count of turns the project has taken, so position is
// (Slot-1) mod n + 1 and a spawn that records it leaves the next read at
// the position after.
type Step struct {
	Harness  harness.Kind
	Position int
	Slot     int
}

// Next is the turn after consumed slots have been taken, skipping every
// harness in skip (a harness a perishable lesson says fails to start: the
// Mate passes `--skip`). skipped lists the harnesses skipped over, once
// each, in order; ok is false when every entry of the rotation is skipped.
// A skipped turn is passed over, not deferred: the spawn that takes the
// step records its slot, and the skipped positions never come back for
// their turn, because the rotation is a share of turns, not a queue.
func (r Rotation) Next(consumed int, skip []harness.Kind) (step Step, skipped []harness.Kind, ok bool) {
	n := len(r)
	if n == 0 || consumed < 0 {
		return Step{}, nil, false
	}
	slot := consumed + 1
	for i := 0; i < n; i++ {
		pos := (slot-1)%n + 1
		k := r[pos-1]
		if !slices.Contains(skip, k) {
			return Step{Harness: k, Position: pos, Slot: slot}, skipped, true
		}
		if !slices.Contains(skipped, k) {
			skipped = append(skipped, k)
		}
		slot++
	}
	return Step{}, skipped, false
}

// Rotations are the rotations in force, each once: every rule's own, in
// rule order, then the table's. A spawn that names no rule consumes the
// first whose next harness it took.
func (t Table) Rotations() []Rotation {
	var out []Rotation
	add := func(r Rotation) {
		if r == nil {
			return
		}
		for _, have := range out {
			if have.Equal(r) {
				return
			}
		}
		out = append(out, r)
	}
	for _, rule := range t.Rules {
		add(rule.Rotate)
	}
	add(t.Rotate)
	return out
}

// Path is the table's path in a workspace whose `.mate/` directory is dir.
func Path(mateDir string) string { return filepath.Join(mateDir, FileName) }

// Load reads and checks the table at path; a profile's harness must be one
// of harnesses. ok is false, with no error, when there is no file: Resolve
// then falls back to the built-in table.
func Load(path string, harnesses harness.Registry) (Table, bool, error) {
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Table{}, false, nil
	}
	if err != nil {
		return Table{}, false, fmt.Errorf("read %s: %w", path, err)
	}
	t, err := parse(raw, harnesses)
	if err != nil {
		return Table{}, false, fmt.Errorf("%w %s: %v", ErrInvalid, path, err)
	}
	t.Path = path
	return t, true, nil
}

type rawProfile struct {
	Harness string `json:"harness"`
	Model   string `json:"model"`
	Effort  string `json:"effort"`
}

type rawRule struct {
	When   string          `json:"when"`
	Use    json.RawMessage `json:"use"`
	Why    string          `json:"why"`
	Budget *int64          `json:"budget"`
	Rotate json.RawMessage `json:"rotate"`
}

type rawTable struct {
	Rules        []rawRule       `json:"rules"`
	Default      json.RawMessage `json:"default"`
	Review       json.RawMessage `json:"review"`
	ReviewBudget *int64          `json:"review_budget"`
	Rotate       json.RawMessage `json:"rotate"`
}

func parse(raw []byte, harnesses harness.Registry) (Table, error) {
	var rt rawTable
	if err := strictUnmarshal(raw, &rt); err != nil {
		return Table{}, err
	}
	var t Table
	for i, r := range rt.Rules {
		where := fmt.Sprintf("rule %d", i+1)
		if strings.TrimSpace(r.When) == "" {
			return Table{}, fmt.Errorf("%s: `when` is required", where)
		}
		if len(r.Use) == 0 {
			return Table{}, fmt.Errorf("%s: `use` is required", where)
		}
		use, err := profiles(r.Use, where, harnesses)
		if err != nil {
			return Table{}, err
		}
		budget, err := checkBudget(r.Budget, where)
		if err != nil {
			return Table{}, err
		}
		rotate, err := checkRotation(r.Rotate, where, harnesses)
		if err != nil {
			return Table{}, err
		}
		t.Rules = append(t.Rules, Rule{When: strings.TrimSpace(r.When), Use: use, Why: strings.TrimSpace(r.Why), Budget: budget, Rotate: rotate})
	}
	if len(rt.Default) > 0 && string(rt.Default) != "null" {
		def, err := profiles(rt.Default, "default", harnesses)
		if err != nil {
			return Table{}, err
		}
		t.Default = def
	}
	if len(rt.Review) > 0 && string(rt.Review) != "null" {
		review, err := profiles(rt.Review, "review", harnesses)
		if err != nil {
			return Table{}, err
		}
		t.Review = review
	}
	reviewBudget, err := checkBudget(rt.ReviewBudget, "review_budget")
	if err != nil {
		return Table{}, err
	}
	t.ReviewBudget = reviewBudget
	rotate, err := checkRotation(rt.Rotate, "rotate", harnesses)
	if err != nil {
		return Table{}, err
	}
	t.Rotate = rotate
	if len(t.Rules) == 0 && len(t.Default) == 0 {
		return Table{}, errors.New("no rules and no default")
	}
	return t, nil
}

// checkBudget reads an optional `budget`: absent is 0 (none), and a value
// must be a positive whole number of tokens - a zero or negative budget is
// a typo, not a way to say none, and is refused like any other bad field.
func checkBudget(raw *int64, where string) (int64, error) {
	if raw == nil {
		return 0, nil
	}
	if *raw <= 0 {
		return 0, fmt.Errorf("%s: `budget` must be a positive number of tokens, not %d", where, *raw)
	}
	return *raw, nil
}

// checkRotation reads an optional `rotate`: absent (or null) is nil, and a
// value must be a non-empty array of registered harness kinds. An empty
// array is a typo, not a way to say none.
func checkRotation(raw json.RawMessage, where string, harnesses harness.Registry) (Rotation, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var names []string
	if err := strictUnmarshal(raw, &names); err != nil {
		return nil, fmt.Errorf("%s: `rotate` must be an array of harness names: %v", where, err)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("%s: `rotate` is an empty array", where)
	}
	out := make(Rotation, 0, len(names))
	for i, name := range names {
		kind, err := harnesses.Parse(name)
		if err != nil {
			return nil, fmt.Errorf("%s: `rotate` entry %d: %v", where, i+1, err)
		}
		out = append(out, kind)
	}
	return out, nil
}

// profiles reads one profile object or a non-empty array of them.
func profiles(raw json.RawMessage, where string, harnesses harness.Registry) ([]Profile, error) {
	var list []rawProfile
	if bytes.HasPrefix(bytes.TrimSpace(raw), []byte("[")) {
		if err := strictUnmarshal(raw, &list); err != nil {
			return nil, fmt.Errorf("%s: %v", where, err)
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("%s: an empty profile array", where)
		}
	} else {
		var one rawProfile
		if err := strictUnmarshal(raw, &one); err != nil {
			return nil, fmt.Errorf("%s: %v", where, err)
		}
		list = []rawProfile{one}
	}
	out := make([]Profile, 0, len(list))
	for j, rp := range list {
		p, err := checkProfile(rp, harnesses)
		if err != nil {
			if len(list) > 1 {
				return nil, fmt.Errorf("%s profile %d: %v", where, j+1, err)
			}
			return nil, fmt.Errorf("%s: %v", where, err)
		}
		out = append(out, p)
	}
	return out, nil
}

// checkProfile is every check a spawn would make, made at load: a known
// harness, one model name, and an effort that harness takes. A harness is
// known when it is registered.
func checkProfile(rp rawProfile, harnesses harness.Registry) (Profile, error) {
	if strings.TrimSpace(rp.Harness) == "" {
		return Profile{}, errors.New("`harness` is required")
	}
	kind, err := harnesses.Parse(strings.TrimSpace(rp.Harness))
	if err != nil {
		return Profile{}, err
	}
	model, err := harness.ParseModel(rp.Model)
	if err != nil {
		return Profile{}, err
	}
	effort, err := harness.ParseEffort(rp.Effort)
	if err != nil {
		return Profile{}, err
	}
	profile, err := harnesses.Lookup(kind)
	if err != nil {
		return Profile{}, err
	}
	if effort != "" && !profile.Info().SupportsEffort(effort) {
		return Profile{}, fmt.Errorf("%s does not take effort %s", kind, effort)
	}
	return Profile{Harness: kind, Model: model, Effort: effort}, nil
}

// strictUnmarshal refuses unknown fields: a misspelt key is a profile that
// silently loses an axis.
func strictUnmarshal(raw []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("trailing data after the JSON value")
	}
	return nil
}
