package agy

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

const testConversation = "25478b13-a998-4c83-89bd-b2fd2d30ee6b"

func crewRequest(t *testing.T) harness.PrepareRequest {
	t.Helper()
	dir := t.TempDir()
	cwd := filepath.Join(dir, "work")
	state := filepath.Join(dir, "state")
	for _, d := range []string{cwd, state} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	brief := filepath.Join(state, "brief.md")
	if err := os.WriteFile(brief, []byte("# Brief\n\nDo the task.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return harness.PrepareRequest{Role: harness.RoleCrew, Cwd: cwd, StateDir: state, ContextPath: brief}
}

// prepare runs Prepare and writes the files it names, as spawn does.
func prepare(t *testing.T, req harness.PrepareRequest) harness.Prepared {
	t.Helper()
	prep, err := Agy{}.Prepare(context.Background(), req)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	for _, f := range prep.Files {
		if err := os.MkdirAll(filepath.Dir(f.Path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f.Path, f.Data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return prep
}

func build(t *testing.T, req harness.PrepareRequest, model string, effort harness.Effort) harness.LaunchSpec {
	t.Helper()
	prep := prepare(t, req)
	spec, err := Agy{}.Build(context.Background(), harness.AgentSpec{
		Role: harness.RoleCrew, Kind: KindAgy, Cwd: req.Cwd, ContextPath: prep.ContextPath, Launch: prep.Launch,
		Model: model, Effort: effort, TaskPrompt: "go", ResumeSessionID: req.ResumeSessionID,
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	return spec
}

func TestAgyWritesTheBriefAsAnAlwaysOnRule(t *testing.T) {
	req := crewRequest(t)
	prep := prepare(t, req)
	want := filepath.Join(req.Cwd, ".agents", "rules", "mate-brief.md")
	if len(prep.Files) != 1 || prep.Files[0].Path != want || !prep.Files[0].Exclude {
		t.Fatalf("files = %+v, want the one excluded rule %s", prep.Files, want)
	}
	data := string(prep.Files[0].Data)
	if !strings.HasPrefix(data, "---\ntrigger: always_on\n---\n\n# Brief\n") {
		t.Fatalf("rule =\n%s\nwant the always_on frontmatter over the brief", data)
	}
	if prep.ContextPath != want || prep.SessionID != "" || prep.Launch != nil {
		t.Fatalf("prepared = %+v", prep)
	}
}

func TestAgyCrewLaunchSkipsPermissionsAndTrustsTheWorktree(t *testing.T) {
	req := crewRequest(t)
	spec := build(t, req, "gemini-4-hydrogen-high", harness.EffortHigh)
	args := spec.Args()
	if len(args) != 7 || args[0] != "--dangerously-skip-permissions" || args[1] != "--settings" {
		t.Fatalf("args = %q", args)
	}
	if !slices.Equal(args[3:], []string{"--model", "gemini-4-hydrogen-high", "--effort", "high"}) {
		t.Fatalf("args = %q", args)
	}
	var settings struct {
		TrustedWorkspaces []string `json:"trustedWorkspaces"`
	}
	if err := json.Unmarshal([]byte(args[2]), &settings); err != nil {
		t.Fatalf("--settings %q: %v", args[2], err)
	}
	if !slices.Contains(settings.TrustedWorkspaces, req.Cwd) {
		t.Fatalf("trusted %q, want the worktree %s", settings.TrustedWorkspaces, req.Cwd)
	}
	if resolved, err := filepath.EvalSymlinks(req.Cwd); err == nil && !slices.Contains(settings.TrustedWorkspaces, resolved) {
		t.Fatalf("trusted %q, want the resolved worktree %s too", settings.TrustedWorkspaces, resolved)
	}
	if spec.Kind() != "agy" || spec.Delivery() != harness.DeliveryInstructionFile || !spec.ContextRequired() || spec.EffortOmitted() {
		t.Fatalf("spec kind %q delivery %q required %v omitted %v", spec.Kind(), spec.Delivery(), spec.ContextRequired(), spec.EffortOmitted())
	}
	if spec.ContextPath() != RulePath(req.Cwd) {
		t.Fatalf("context path %s, want the rule", spec.ContextPath())
	}
}

func TestAgyPassesEveryEffort(t *testing.T) {
	for _, e := range []harness.Effort{harness.EffortLow, harness.EffortMedium, harness.EffortHigh, harness.EffortXHigh, harness.EffortMax} {
		spec := build(t, crewRequest(t), "", e)
		if args := spec.Args(); !slices.Equal(args[len(args)-2:], []string{"--effort", string(e)}) || spec.EffortOmitted() {
			t.Errorf("effort %s: args %q, omitted %v", e, args, spec.EffortOmitted())
		}
	}
}

func TestAgyResumeNamesTheConversation(t *testing.T) {
	req := crewRequest(t)
	req.ResumeSessionID = testConversation
	spec := build(t, req, "", "")
	if args := spec.Args(); !slices.Equal(args[3:], []string{"--conversation", testConversation}) {
		t.Fatalf("args = %q", args)
	}
	req.ResumeSessionID = "--model"
	prep := prepare(t, req)
	if _, err := (Agy{}).Build(context.Background(), harness.AgentSpec{
		Role: harness.RoleCrew, Kind: KindAgy, Cwd: req.Cwd, ContextPath: prep.ContextPath, ResumeSessionID: req.ResumeSessionID,
	}); err == nil {
		t.Fatal("a resume id that is not a UUID was passed on")
	}
}

func TestAgyRefusesABriefPastTheRuleCap(t *testing.T) {
	req := crewRequest(t)
	if err := os.WriteFile(req.ContextPath, []byte(strings.Repeat("x\n", MaxRuleBytes/2)), 0o644); err != nil {
		t.Fatal(err)
	}
	prep := prepare(t, req)
	_, err := Agy{}.Build(context.Background(), harness.AgentSpec{
		Role: harness.RoleCrew, Kind: KindAgy, Cwd: req.Cwd, ContextPath: prep.ContextPath,
	})
	if !errors.Is(err, harness.ErrContextTooLarge) {
		t.Fatalf("Build of a rule past %d bytes = %v, want ErrContextTooLarge", MaxRuleBytes, err)
	}
}

func TestAgyRefusesAContextOtherThanTheRule(t *testing.T) {
	req := crewRequest(t)
	_, err := Agy{}.Build(context.Background(), harness.AgentSpec{
		Role: harness.RoleCrew, Kind: KindAgy, Cwd: req.Cwd, ContextPath: req.ContextPath,
	})
	if !errors.Is(err, harness.ErrContextRequired) {
		t.Fatalf("Build with the brief itself = %v, want ErrContextRequired", err)
	}
}

func TestAgyRefusesAMateNamingHooks(t *testing.T) {
	req := crewRequest(t)
	req.Role = harness.RoleMate
	_, err := Agy{}.Prepare(context.Background(), req)
	if err == nil || !strings.Contains(err.Error(), "Hooks") {
		t.Fatalf("Prepare of a Mate = %v, want a refusal naming Hooks", err)
	}
}

// writeHome lays out an agy data home: the cwd map and one conversation
// database written at mod.
func writeHome(t *testing.T, cwd, id string, mod time.Time) string {
	t.Helper()
	home := t.TempDir()
	for _, d := range []string{"cache", "conversations"} {
		if err := os.MkdirAll(filepath.Join(home, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	m, _ := json.Marshal(map[string]string{cwd: id, "/elsewhere": "77d6d809-3256-4717-806a-a5803675c88d"})
	if err := os.WriteFile(filepath.Join(home, "cache", "last_conversations.json"), m, 0o644); err != nil {
		t.Fatal(err)
	}
	db := ConversationPath(home, id)
	if err := os.WriteFile(db, []byte("SQLite format 3\x00"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(db, mod, mod); err != nil {
		t.Fatal(err)
	}
	return home
}

func TestAgySessionAtStopReadsTheCwdMap(t *testing.T) {
	cwd := t.TempDir()
	resolved, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		t.Fatal(err)
	}
	launched := time.Now().Add(-time.Minute)
	// agy keys the map by the directory the shell resolved.
	home := writeHome(t, resolved, testConversation, time.Now())
	s := agySessions{home: home}
	if got := s.AtStop(cwd, launched, ""); got != testConversation {
		t.Fatalf("AtStop = %q, want %s", got, testConversation)
	}
	if got := s.AtStop(cwd, launched, "from-herdr"); got != "from-herdr" {
		t.Fatalf("AtStop with a runtime ref = %q", got)
	}
	if got := s.AtStop(filepath.Join(cwd, "other"), launched, ""); got != "" {
		t.Fatalf("AtStop of a directory the map lacks = %q", got)
	}
	// An entry whose conversation was last written before this launch is an
	// earlier agent's.
	if got := s.AtStop(cwd, time.Now().Add(time.Hour), ""); got != "" {
		t.Fatalf("AtStop of a stale entry = %q", got)
	}
}

func TestAgyResumableLooksForTheDatabase(t *testing.T) {
	home := writeHome(t, "/w", testConversation, time.Now())
	s := agySessions{home: home}
	if err := s.Resumable(testConversation); err != nil {
		t.Fatalf("Resumable of a stored conversation: %v", err)
	}
	var missing *harness.NoSessionError
	if err := s.Resumable("3b03c470-291f-4210-abd3-7dc8e18b0ff7"); !errors.As(err, &missing) {
		t.Fatalf("Resumable of an unknown conversation = %v, want a NoSessionError", err)
	}
	if err := s.Resumable("../secrets"); !errors.As(err, &missing) {
		t.Fatalf("Resumable of a non-UUID = %v, want a NoSessionError", err)
	}
}

func TestAgyStopsByClearingThenExiting(t *testing.T) {
	stop := Agy{}.Capabilities().GracefulStop
	if !stop.Verified() {
		t.Fatal("GracefulStop is not verified")
	}
	if keys := stop.Impl.ClearKeys(); !slices.Equal(keys, []string{"ctrl+u"}) {
		t.Fatalf("ClearKeys = %q: ctrl+c would arm the exit with the draft in it", keys)
	}
	if p := stop.Impl.ExitPrompt(); p != "/exit" {
		t.Fatalf("ExitPrompt = %q", p)
	}
}
