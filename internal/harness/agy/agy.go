// Package agy is the Antigravity CLI (agy 1.5.10) as a mate harness: its
// Profile, launch, screens and sessions. Everything that knows agy lives
// here; the contract it implements is internal/harness.
//
// agy runs Crews only. It loads hooks from plugins, but mate has not
// measured a hook command against `mate hook`, and a Mate's memory and inbox
// rest on hooks; a Mate runs on Claude Code anyway (CON-236). Every number
// here was measured on agy 1.5.10 in a 120x40 Herdr 0.9.5 pane, read as
// visible text, on 2026-11-18 (CON-278).
package agy

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/observability"
)

// KindAgy is agy's kind, and what Herdr is told at `agent start --kind`.
const KindAgy harness.Kind = "agy"

// Agy is the Profile of agy 1.5.10.
type Agy struct {
	// Home is agy's CLI data directory, which holds its conversations and
	// the cwd-to-conversation map. Empty means agy's own,
	// ~/.gemini/antigravity-cli. The launch does not set it: the sign-in
	// lives beside it, so a Crew uses the operator's.
	Home string
}

// Kind implements Profile.
func (Agy) Kind() harness.Kind { return KindAgy }

// Info implements Profile.
func (Agy) Info() harness.Info {
	return harness.Info{
		Name:        "Antigravity",
		RuntimeKind: string(KindAgy),
		// agy reads a workspace's rules, skills and workflows from `.agents/`
		// at its root, and `GEMINI.md` beside `AGENTS.md`.
		ConfigDir:       ".agents",
		InstructionFile: "GEMINI.md",
		// agy has no brand mark that was measured as one cell, so it draws
		// the same letter in every alphabet.
		Icon: harness.Icon{Nerd: "A", Unicode: "A", ASCII: "A"},
		Documents: []harness.Document{
			{Path: RuleRelPath, Role: "agy_rule"},
			{Path: "GEMINI.md", Role: "repo_gemini"},
		},
		// agy 1.5.10 --help: "--effort ... (low|medium|high|xhigh|max)".
		Efforts:      []harness.Effort{harness.EffortLow, harness.EffortMedium, harness.EffortHigh, harness.EffortXHigh, harness.EffortMax},
		AdapterNotes: adapterNotes,
	}
}

// adapterNotes is agy's section of the crew-recovery skill.
//
//go:embed adapter.md
var adapterNotes string

// Launcher implements Profile.
func (a Agy) Launcher() harness.Launcher { return a }

// Screen implements Profile.
func (Agy) Screen() harness.ScreenProfile { return agyScreen{} }

func measured(proof string) harness.Evidence {
	return harness.Evidence{Version: "agy 1.5.10", Measured: "2026-11-18 (Herdr 0.9.5, visible 120x40)", Proof: proof}
}

// hooksReason is why agy's Hooks are not verified, and so why it runs no Mate.
const hooksReason = "agy loads PreToolUse, Stop and SessionStart hooks from a plugin's hooks.json, but mate has not measured a hook command whose stdin is the payload `mate hook` writes, so no HookInstaller is implemented"

// measuredTranscript is when agy's conversation database was measured
// against its own stream-json usage (CON-279).
func measuredTranscript(proof string) harness.Evidence {
	return harness.Evidence{Version: "agy 1.5.18", Measured: "2026-11-27 (agy -p --output-format stream-json)", Proof: proof}
}

// Capabilities implements Profile.
func (a Agy) Capabilities() harness.Capabilities {
	return harness.Capabilities{
		GracefulStop: harness.Cap[harness.GracefulStopper]{
			Status: harness.CapVerified,
			Impl:   agyExit{},
			Evidence: measured("ctrl+c does not clear a draft: the first press arms \"press ctrl+c again to exit\" and the second exits with the draft lost; ctrl+u clears the draft, and `/exit` (alias /quit) then Enter exits and prints `agy --conversation=<id>` " +
				"(TestAgyStopsByClearingThenExiting)"),
		},
		Session: harness.Cap[harness.SessionIdentity]{
			Status: harness.CapVerified,
			Impl:   agySessions{home: a.Home},
			Evidence: measured("a launch cannot name its conversation (--conversation with a new UUID warns \"not found\" and starts another); cache/last_conversations.json maps the cwd to the conversation it ran, and --conversation <id> resumed it with the earlier scrollback " +
				"(TestAgySessionAtStopReadsTheCwdMap, TestAgyResumableLooksForTheDatabase)"),
		},
		Hooks:   harness.Cap[harness.HookInstaller]{Status: harness.CapUnknown, Reason: hooksReason},
		TurnEnd: harness.Cap[harness.TurnEndEvidence]{Status: harness.CapUnknown, Reason: "the end of a turn is in the conversation database, but TurnEndEvidence reads a transcript as bytes, which a SQLite database in WAL mode is not; agy runs no Mate, the one agent whose stow reads it"},
		Transcript: harness.Cap[harness.TranscriptSource]{
			Status: harness.CapVerified,
			Impl:   agyTranscripts{home: a.Home},
			Evidence: measuredTranscript("conversations/<id>.db `steps`: a response's metadata field 9 gave the stream-json usage of the same calls (9.2 input 15098 and 2180, 9.5 cache read 13044, 9.3 output 168 = 9.9 thinking 161 + 9.10 visible 7), and a tool call's field 4 its id, name and arguments; " +
				"the database file changes only when agy checkpoints its WAL, so a running Crew's tokens arrive in steps and all of them once it exits (TestAgyTranscriptReadsTheMeasuredUsage)"),
		},
		Quota: harness.Cap[harness.QuotaProvider]{
			Status: harness.CapUnknown,
			Reason: "quota-axi 0.1.69 has an agy row, but it bounds Gemini and Claude/GPT models in separate scopes and none provider-wide, and mate folds only provider-wide scopes (CON-279)",
		},
		MCPConfig: harness.Cap[harness.MCPConfigurer]{Status: harness.CapUnsupported, Reason: mcpUnsupported},
	}
}

// agyExit is agy's GracefulStopper. ctrl+c would arm, and on a second press
// take, the exit with the draft in it, so ctrl+u clears the composer first;
// `/exit` opens the slash menu on that command, and Enter runs it.
type agyExit struct{}

// ClearKeys implements GracefulStopper.
func (agyExit) ClearKeys() []string { return []string{"ctrl+u"} }

// ExitPrompt implements GracefulStopper.
func (agyExit) ExitPrompt() string { return "/exit" }

// agySessions is agy's SessionIdentity. A launch cannot name its
// conversation, so it is read at stop from the map agy keeps of the last
// conversation each directory ran.
type agySessions struct{ home string }

// AtStop implements SessionIdentity. Herdr 0.9.5 records no agent_session
// for agy, so runtimeRef is used only when it ever does. The map's entry
// counts when its conversation's database was written at or after the
// launch: an entry from an earlier agent in the same directory does not.
func (s agySessions) AtStop(cwd string, launched time.Time, runtimeRef string) string {
	if ref := strings.TrimSpace(runtimeRef); ref != "" {
		return ref
	}
	home, err := dataHome(s.home)
	if err != nil {
		return ""
	}
	id, ok := LastConversation(home, cwd)
	if !ok {
		return ""
	}
	fi, err := os.Stat(ConversationPath(home, id))
	if err != nil || !fi.Mode().IsRegular() {
		return ""
	}
	if !launched.IsZero() && fi.ModTime().Before(launched) {
		return ""
	}
	return id
}

// Resumable implements SessionIdentity. agy 1.5.10 answers --conversation
// with an unknown id by warning and starting a new conversation, so the
// database is looked for before anything is launched.
func (s agySessions) Resumable(id string) error {
	home, err := dataHome(s.home)
	if err != nil {
		return fmt.Errorf("cannot look for the agy conversation %s to resume (%v)", id, err)
	}
	parsed, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return &harness.NoSessionError{Session: "the agy conversation " + id, Missing: "it is not a conversation UUID"}
	}
	path := ConversationPath(home, parsed.String())
	if fi, err := os.Stat(path); err != nil || !fi.Mode().IsRegular() {
		return &harness.NoSessionError{Session: "the agy conversation " + id, Missing: path + " does not exist"}
	}
	return nil
}

// dataHome is agy's CLI data directory: the profile's, else
// ~/.gemini/antigravity-cli.
func dataHome(home string) (string, error) {
	dir := strings.TrimSpace(home)
	if dir == "" {
		user, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("agy home: %w", err)
		}
		return filepath.Join(user, ".gemini", "antigravity-cli"), nil
	}
	if !filepath.IsAbs(dir) {
		return "", fmt.Errorf("agy home %q is not absolute", dir)
	}
	return filepath.Clean(dir), nil
}

// ConversationPath is one conversation's database under agy's data home.
func ConversationPath(home, id string) string {
	return filepath.Join(home, "conversations", id+".db")
}

// LastConversation is the conversation agy last ran in cwd, from
// cache/last_conversations.json. agy keys it by its working directory as it
// was handed it: /private/tmp/x after a shell's `cd`, /tmp/x when Herdr
// started it there (both measured). So cwd is looked up as given and then
// with its symlinks resolved.
func LastConversation(home, cwd string) (string, bool) {
	raw, err := os.ReadFile(filepath.Join(home, "cache", "last_conversations.json"))
	if err != nil {
		return "", false
	}
	var byCwd map[string]string
	if json.Unmarshal(raw, &byCwd) != nil {
		return "", false
	}
	keys := []string{filepath.Clean(cwd)}
	if resolved, err := filepath.EvalSymlinks(cwd); err == nil && resolved != keys[0] {
		keys = append(keys, resolved)
	}
	for _, k := range keys {
		parsed, err := uuid.Parse(strings.TrimSpace(byCwd[k]))
		if err == nil {
			return parsed.String(), true
		}
	}
	return "", false
}

// PaneEnv implements Launcher. agy finds its data home and sign-in from the
// operator's home, so a pane needs nothing pinned.
func (Agy) PaneEnv() ([]harness.EnvVar, error) { return nil, nil }

// mcpUnsupported is why the MCPConfig capability is unsupported: agy reads
// MCP servers from one global file (`agy mcp add`), not from a launch, and
// it cannot run a Mate, the only agent plugins reach.
const mcpUnsupported = "agy reads MCP servers only from its global mcp_config.json, which a launch cannot scope, and it cannot run a Mate"

// refuseMate is the refusal of a Mate on agy, naming the capability it lacks.
func refuseMate() error {
	return observability.NewError(observability.CodeUsage,
		"the Antigravity harness cannot run a Mate: it declares no verified Hooks capability ("+hooksReason+"); it can still run a Crew")
}

// RuleRelPath is where, relative to a Crew's worktree, its brief is written
// as an agy rule. agy has no flag that adds to the system prompt; it loads
// every rule under .agents/rules whose frontmatter says always_on.
const RuleRelPath = ".agents/rules/mate-brief.md"

// ruleFrontmatter makes the rule load on every turn. A rule without it is
// not loaded on its own (measured: the brief's codeword went unanswered
// until the frontmatter was added).
const ruleFrontmatter = "---\ntrigger: always_on\n---\n\n"

// MaxRuleBytes is agy's cap on one rule file, frontmatter included; past it
// agy truncates the rule on a line boundary without a word.
const MaxRuleBytes = 24000

// RulePath is the brief's rule file in a Crew's worktree.
func RulePath(cwd string) string { return filepath.Join(cwd, filepath.FromSlash(RuleRelPath)) }

// Prepare implements Launcher. The brief becomes an always-on rule in the
// worktree, kept out of git. agy names its conversation only once it runs,
// so a fresh launch has no session id.
func (Agy) Prepare(_ context.Context, req harness.PrepareRequest) (harness.Prepared, error) {
	if err := req.Check(KindAgy); err != nil {
		return harness.Prepared{}, err
	}
	if err := harness.RefuseMCPServers(KindAgy, mcpUnsupported, req.MCPServers); err != nil {
		return harness.Prepared{}, err
	}
	if req.Role == harness.RoleMate {
		return harness.Prepared{}, refuseMate()
	}
	brief, err := os.ReadFile(req.ContextPath)
	if err != nil {
		return harness.Prepared{}, err
	}
	rule := RulePath(req.Cwd)
	data := append([]byte(ruleFrontmatter), brief...)
	return harness.Prepared{
		Files:       []harness.LaunchFile{{Path: rule, Data: data, Exclude: true}},
		SessionID:   strings.TrimSpace(req.ResumeSessionID),
		ContextPath: rule,
	}, nil
}

// Build implements Launcher.
func (a Agy) Build(_ context.Context, spec harness.AgentSpec) (harness.LaunchSpec, error) {
	if spec.Kind != "" && spec.Kind != KindAgy {
		return harness.LaunchSpec{}, observability.NewError(observability.CodeUsage, fmt.Sprintf("agy adapter got kind %q", spec.Kind))
	}
	if spec.Launch != nil {
		return harness.LaunchSpec{}, observability.NewError(observability.CodeUsage, fmt.Sprintf("agy adapter got launch data %T that its Prepare did not make", spec.Launch))
	}
	if err := harness.RefuseMCPServers(KindAgy, mcpUnsupported, spec.MCPServers); err != nil {
		return harness.LaunchSpec{}, err
	}
	if spec.Role == harness.RoleMate {
		return harness.LaunchSpec{}, refuseMate()
	}
	cwd := spec.Cwd
	if cwd == "" || !filepath.IsAbs(cwd) {
		return harness.LaunchSpec{}, observability.WrapError(observability.CodeUsage, "agy launch requires the absolute cwd the agent will run in", harness.ErrContextRequired)
	}
	want := RulePath(cwd)
	path := spec.ContextPath
	if path == "" {
		path = want
	}
	if filepath.Clean(path) != filepath.Clean(want) {
		return harness.LaunchSpec{}, observability.WrapError(observability.CodeUsage,
			fmt.Sprintf("agy context path must be the rule %s (got %s)", want, path), harness.ErrContextRequired)
	}
	settings, err := launchSettings(cwd)
	if err != nil {
		return harness.LaunchSpec{}, err
	}
	args := []string{
		// The captain ruled on 2026-09-14 that a Crew asks no permission.
		"--dangerously-skip-permissions",
		// --settings is for this session only and never saved: it trusts the
		// worktree, so no folder dialog stands between the launch and the
		// composer.
		"--settings", settings,
	}
	if resumeID := strings.TrimSpace(spec.ResumeSessionID); resumeID != "" {
		if _, err := uuid.Parse(resumeID); err != nil {
			return harness.LaunchSpec{}, observability.WrapError(observability.CodeUsage,
				fmt.Sprintf("agy resume id %q is not a conversation UUID", spec.ResumeSessionID), harness.ErrContextRequired)
		}
		args = append(args, "--conversation", resumeID)
	}
	if spec.Model != "" {
		args = append(args, "--model", spec.Model)
	}
	if a.Info().SupportsEffort(spec.Effort) {
		args = append(args, "--effort", string(spec.Effort))
	}
	return harness.NewLaunchSpec(harness.LaunchPlan{
		RuntimeKind:     a.Info().RuntimeKind,
		Screen:          a.Screen(),
		Args:            args,
		Cwd:             cwd,
		Env:             spec.Env,
		EnvKeys:         spec.EnvKeys,
		ContextFiles:    []harness.GeneratedFile{{Path: want, Role: "agy_rule"}},
		Delivery:        harness.DeliveryInstructionFile,
		ContextPath:     want,
		ContextRequired: true,
		MaxFileBytes:    MaxRuleBytes,
		TaskPrompt:      spec.TaskPrompt,
		Model:           spec.Model,
		Effort:          spec.Effort,
		EffortOmitted:   spec.Effort != "" && !a.Info().SupportsEffort(spec.Effort),
		Notes: []string{
			"the captain ruled on 2026-09-14 that Crew launches run with no permission prompt; agy is passed --dangerously-skip-permissions",
			"the brief is the always-on rule " + RuleRelPath + ", excluded from git; agy has no flag that adds to the system prompt, and caps a rule at 24000 bytes",
			"--settings trusts the worktree for this session only, so agy draws no folder-trust dialog and saves nothing",
			"agy names its conversation only once it runs; the id is read at stop from cache/last_conversations.json",
		},
	})
}

// launchSettings is the --settings document: the worktree trusted, as given
// and with its symlinks resolved, since agy compares the directory the
// shell resolved.
func launchSettings(cwd string) (string, error) {
	trusted := []string{filepath.Clean(cwd)}
	if resolved, err := filepath.EvalSymlinks(cwd); err == nil && resolved != trusted[0] {
		trusted = append(trusted, resolved)
	}
	b, err := json.Marshal(map[string][]string{"trustedWorkspaces": trusted})
	if err != nil {
		return "", fmt.Errorf("agy settings: %w", err)
	}
	return string(b), nil
}
