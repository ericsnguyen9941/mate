package spawn_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/harness/claude"
	"github.com/nguyenngocanh94/mate/internal/matehost"
	"github.com/nguyenngocanh94/mate/internal/process"
	"github.com/nguyenngocanh94/mate/internal/runtime"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// TestLiveHeadlessMate is CON-78's proof: a headless Claude Mate in a real
// Herdr pane. Herdr lists the host under the Mate's name, two app lines
// are two `claude -p` turns on one session (the second resumed), the
// Mate's own hooks record both answers in sent.log, and the stop is the
// host's own clean exit.
func TestLiveHeadlessMate(t *testing.T) {
	requireLive(t)
	session, configHome := liveLabSession(t)

	// The host's socket lives in the Mate directory, so the workspace sits
	// under /private/tmp: short enough for a unix socket path, and with no
	// symlink for Herdr's resolved pane cwd to disagree with.
	root, err := os.MkdirTemp("/private/tmp", "mh")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(root) })
	w, err := store.Init(root, store.Defaults{})
	if err != nil {
		t.Fatalf("store.Init: %v", err)
	}
	useLabSession(t, w, session)
	if w, err = store.Open(root); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(w.ProjectHome("shop"), "shop")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	liveGit(t, repo, "init", "-b", "main")
	liveGit(t, repo, "config", "user.email", "mate-test@example.com")
	liveGit(t, repo, "config", "user.name", "mate test")
	liveGit(t, repo, "commit", "--allow-empty", "-m", "init")
	if err := w.AddProject("shop", store.ProjectConfig{
		Mate:  store.MateConfig{Model: "haiku", Effort: "low"},
		Repos: []store.RepoConfig{{Path: repo, DefaultBranch: "main"}},
	}); err != nil {
		t.Fatalf("AddProject: %v", err)
	}

	names := runtime.NewMemoryNameRegistry()
	rt := runtime.NewHerdr(process.ExecRunner{})
	rt.Names = names
	rt.StartServer = func(context.Context, string) error {
		t.Fatal("the lab session is provisioned by the runner; this test must not start a Herdr server")
		return nil
	}
	deps := spawn.Deps{
		Harnesses:        catalog.Default(),
		Runtime:          rt,
		Names:            names,
		ConfigHome:       configHome,
		Binary:           binaryPath(t),
		ReadinessTimeout: 90 * time.Second,
		StartTimeout:     60 * time.Second,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	started, err := spawn.StartMate(ctx, w, deps, spawn.StartRequest{Project: "shop", Harness: claude.KindClaude, Mode: spawn.ModeHeadless})
	if err != nil {
		t.Fatalf("StartMate: %v", err)
	}
	t.Cleanup(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer stopCancel()
		_, _ = spawn.StopMate(stopCtx, w, deps, "shop")
	})
	if !started.Headless || started.Agent != "mate-shop" {
		t.Fatalf("started = %+v", started)
	}
	handle := runtime.AgentHandle{Session: runtime.SessionHandle{Name: session, ConfigHome: configHome}, Name: started.Agent, RawID: "shop"}
	if _, err := rt.InspectAgent(ctx, handle); err != nil {
		t.Fatalf("Herdr does not list the host as %s: %v", started.Agent, err)
	}

	host := matehost.Dial(w.MateDir("shop"))
	for i, word := range []string{"pong", "pang"} {
		line := matehost.Marker + "ping: run no commands and reply with exactly the word " + word
		if err := host.Deliver(ctx, line); err != nil {
			t.Fatalf("deliver %d: %v", i, err)
		}
		for {
			st, err := host.State(ctx)
			if err != nil {
				t.Fatalf("state: %v", err)
			}
			if !st.Busy {
				break
			}
			select {
			case <-ctx.Done():
				t.Fatal("the turn never ended")
			case <-time.After(time.Second):
			}
		}
		items, err := host.History(ctx, 1)
		if err != nil || len(items) != 1 {
			t.Fatalf("history = %+v, %v", items, err)
		}
		if it := items[0]; it.Error != "" || !strings.Contains(strings.ToLower(it.Reply), word) {
			t.Fatalf("turn %d = %+v, want a reply with %q", i, it, word)
		}
	}
	entries, _, err := w.ReadSent("shop", 0)
	if err != nil {
		t.Fatal(err)
	}
	answers := 0
	for _, e := range entries {
		if e.Source == store.SourceMate && e.Target == store.SourceUser {
			answers++
		}
	}
	if answers < 2 {
		t.Fatalf("sent.log has %d answers from the Mate, want both turns' Stop hooks: %+v", answers, entries)
	}
	meta := readMeta(t, w, "shop")
	if meta[spawn.MetaSessionID] != started.SessionID || meta[spawn.MetaTranscript] == "" {
		t.Fatalf("meta = %v, want the one session and its transcript", meta)
	}

	stopped, err := spawn.StopMate(ctx, w, deps, "shop")
	if err != nil {
		t.Fatalf("StopMate: %v", err)
	}
	if !stopped.Confirmed || stopped.Forced != "" {
		t.Fatalf("stop = %+v, want the host's clean exit", stopped)
	}
}
