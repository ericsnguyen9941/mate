package spawn_test

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/brief/brieftest"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/spawn"
)

// A fresh relaunch (CON-73) composes the new session's opening input out
// of exactly four parts in a fixed order - the brief, the last hand-back,
// the Mate's message, the diff since the last review - writes it to
// crews/<id>/session-<n>.md and points the new agent at that file, with
// the note where a relaunch note goes.
func TestRelaunchCrewFreshComposesTheFourPartsInOrder(t *testing.T) {
	w := crewWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	res, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project: "shop", Crew: "k3", BriefText: brieftest.Ship("add the healthcheck"),
	})
	if err != nil {
		t.Fatalf("SpawnCrew: %v", err)
	}
	// The crew committed, handed back and opened a PR in its first session.
	if err := os.WriteFile(filepath.Join(res.Worktree, "health.go"), []byte("package shop\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, res.Worktree, "add", "health.go")
	git(t, res.Worktree, "-c", "user.email=k3@example.com", "-c", "user.name=k3", "commit", "-q", "-m", "add healthcheck")
	if err := os.WriteFile(w.CrewHandback("shop", "k3"), []byte("## For the next crew\n- the endpoint is /healthz\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status := "working: adding\npr-open: https://github.com/x/shop/pull/7\n"
	if err := os.WriteFile(w.CrewStatus("shop", "k3"), []byte(status), 0o644); err != nil {
		t.Fatal(err)
	}

	rec := &promptRecorder{Fake: rt}
	deps.Runtime = rec
	message := "Review F1: /healthz returns 200 on a closed database; check the pool first. Already tried: nothing yet."
	again, err := spawn.RelaunchCrew(context.Background(), w, deps, "shop", "k3", "Already tried: nothing yet",
		spawn.RelaunchOptions{Fresh: &spawn.FreshSession{Message: message}})
	if err != nil {
		t.Fatalf("RelaunchCrew: %v", err)
	}
	if again.Ordinal != 2 {
		t.Fatalf("Ordinal = %d, want 2 (the second session)", again.Ordinal)
	}
	if again.InputPath != w.CrewSessionInput("shop", "k3", 2) {
		t.Fatalf("InputPath = %q, want %q", again.InputPath, w.CrewSessionInput("shop", "k3", 2))
	}
	if again.Pane == res.Pane || !again.Stopped {
		t.Fatalf("the resting pane was not replaced: pane %q (was %q), stopped %v", again.Pane, res.Pane, again.Stopped)
	}

	data, err := os.ReadFile(again.InputPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	briefText, _ := os.ReadFile(res.BriefPath)
	parts := []string{
		"# Part 1: your brief",
		strings.TrimSpace(string(briefText)),
		"# Part 2: your last hand-back",
		"- the endpoint is /healthz",
		"# Part 3: the Mate's message",
		message,
		"# Part 4: change since the last review (git diff --stat main...HEAD)",
		"health.go",
	}
	at := -1
	for _, p := range parts {
		i := strings.Index(text, p)
		if i < 0 {
			t.Fatalf("session input lacks %q:\n%s", p, text)
		}
		if i < at {
			t.Fatalf("session input has %q out of order:\n%s", p, text)
		}
		at = i
	}
	if strings.Count(text, "# Part ") != 4 {
		t.Fatalf("session input has %d parts, want exactly 4:\n%s", strings.Count(text, "# Part "), text)
	}

	if len(rec.prompts) != 1 {
		t.Fatalf("prompts = %q, want exactly one", rec.prompts)
	}
	prompt := rec.prompts[0]
	if !strings.Contains(prompt, again.InputPath) || !strings.Contains(prompt, "session 2") {
		t.Fatalf("prompt %q does not point at the session input", prompt)
	}
	if !strings.Contains(prompt, "Progress so far: Already tried: nothing yet.") {
		t.Fatalf("prompt %q does not carry the note where a relaunch note goes", prompt)
	}
	if strings.Contains(prompt, "\n") {
		t.Fatalf("prompt is not one line: %q", prompt)
	}

	meta, err := w.ReadCrewMeta("shop", "k3")
	if err != nil {
		t.Fatal(err)
	}
	if meta[spawn.MetaSessions] != "2" {
		t.Fatalf("meta sessions = %q, want 2", meta[spawn.MetaSessions])
	}
	if meta[spawn.MetaSessionStatusOffset] != strconv.Itoa(len(status)) {
		t.Fatalf("meta session_status_offset = %q, want %d (the status file's size at relaunch)", meta[spawn.MetaSessionStatusOffset], len(status))
	}
	if got, _ := os.ReadFile(w.CrewStatus("shop", "k3")); string(got) != status {
		t.Fatalf("status file = %q; a fresh relaunch must not rewrite it", got)
	}
}

// The fourth part is the diff from the commit the latest review saw when
// the caller names one, and says so in its heading.
func TestRelaunchCrewFreshDiffsFromTheReviewedCommit(t *testing.T) {
	w := crewWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	res, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project: "shop", Crew: "k3", BriefText: brieftest.Ship("work"),
	})
	if err != nil {
		t.Fatalf("SpawnCrew: %v", err)
	}
	reviewed := strings.TrimSpace(git(t, res.Worktree, "rev-parse", "HEAD"))
	if err := os.WriteFile(filepath.Join(res.Worktree, "after.txt"), []byte("later\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, res.Worktree, "add", "after.txt")
	git(t, res.Worktree, "-c", "user.email=k3@example.com", "-c", "user.name=k3", "commit", "-q", "-m", "after review")

	again, err := spawn.RelaunchCrew(context.Background(), w, deps, "shop", "k3", "",
		spawn.RelaunchOptions{Fresh: &spawn.FreshSession{Message: "F2 open", Reviewed: reviewed}})
	if err != nil {
		t.Fatalf("RelaunchCrew: %v", err)
	}
	data, err := os.ReadFile(again.InputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "git diff --stat "+reviewed+"...HEAD") {
		t.Fatalf("session input does not name the reviewed range:\n%s", data)
	}
	if !strings.Contains(string(data), "after.txt") {
		t.Fatalf("session input's diff lacks the change since the review:\n%s", data)
	}
	if !strings.Contains(string(data), "# Part 2: your last hand-back\n\n(no hand-back on disk)") {
		t.Fatalf("a missing hand-back must be said, not left blank:\n%s", data)
	}
}

// Every relaunch that opens a new conversation counts a session; the plain
// relaunch (brief pointer alone) counts too and writes no session input.
func TestRelaunchCrewCountsSessions(t *testing.T) {
	w := crewWorkspace(t, "shop")
	rt := runtime.NewFake()
	deps := fakeDeps(t, rt)
	if _, err := spawn.SpawnCrew(context.Background(), w, deps, spawn.SpawnCrewRequest{
		Project: "shop", Crew: "k3", BriefText: brieftest.Ship("work"),
	}); err != nil {
		t.Fatalf("SpawnCrew: %v", err)
	}
	meta, _ := w.ReadCrewMeta("shop", "k3")
	if spawn.SessionOrdinal(meta) != 1 {
		t.Fatalf("spawn: sessions = %q, want 1", meta[spawn.MetaSessions])
	}
	for want := 2; want <= 3; want++ {
		res, err := spawn.RelaunchCrew(context.Background(), w, deps, "shop", "k3", "")
		if err != nil {
			t.Fatalf("RelaunchCrew #%d: %v", want, err)
		}
		if res.Ordinal != want || res.InputPath != "" {
			t.Fatalf("relaunch #%d: Ordinal %d InputPath %q, want %d and none", want, res.Ordinal, res.InputPath, want)
		}
		meta, _ = w.ReadCrewMeta("shop", "k3")
		if spawn.SessionOrdinal(meta) != want {
			t.Fatalf("relaunch #%d: sessions = %q, want %d", want, meta[spawn.MetaSessions], want)
		}
	}
	// A record from before the key existed is the first session.
	if spawn.SessionOrdinal(map[string]string{}) != 1 || spawn.SessionOrdinal(map[string]string{spawn.MetaSessions: "x"}) != 1 {
		t.Fatal("SessionOrdinal must read a missing or unreadable count as 1")
	}
}
