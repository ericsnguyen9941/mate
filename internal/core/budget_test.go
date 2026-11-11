package core

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// The core's budget (CON-263, epic CON-253): what the two things it does on
// every tick cost on a workspace with a long history and little going on.
// On 2026-11-11 `mate core run` was found at 1.4 GB RSS and ~70% of a CPU
// for days: one steady timeline ingest pass over the captain's 1308 crews
// cost 4.9 s and 1.38 GB allocated every 5 s poll, and one query.Load
// ~100 ms and 62 MB every 2 s a UI read. Nobody noticed, because nothing
// measured it. This test does, over a generated workspace shaped like that
// one - a thousand crews closed long ago, a few at work - and fails once a
// steady pass or a load allocates past its budget.
//
// The hard assertion is on allocation, which is the same on any machine:
// bytes and objects of one call, read from runtime.MemStats around it (no
// test of this package runs in parallel, so the counters are the call's
// within what the database driver's goroutines allocate meanwhile). The
// time budgets only catch a call gone from a fraction of a second to many
// seconds: on the same laptop at a load average of 275 (other worktrees'
// `make check`), query.Load took 2 s instead of 75 ms.
//
// Measured 2026-11-11 on main (eeed0ee, Apple M-series) over this fixture:
// a steady pass 67.6 MB in 703k objects and 576 ms, a query.Load 21.6 MB
// in 128k objects and 75 ms (Wave 0 had brought the real workspace's pass
// to 0.68 s and 78 MB). Each allocation budget is about twice that. With
// internal/timeline, internal/db and internal/harness/codex as they were
// before Wave 0 (fbef89a), the same steady pass allocated 2972 MB: this
// test would have failed on it.
//
// With CON-259's seals (the core-idle wave, 2026-11-11) a steady pass is
// 2.6 MB in 24k objects and 15 ms over the same fixture: a sealed crew
// costs two stats, so the pass no longer grows with the closed crews and
// its budget is a fixed one, about twice that, the same with a tenth of the
// crews under the race detector (0.7 MB in 7k objects there). query.Load
// (23.9 MB in 144k objects) still grows with them. Any change that makes
// either cheaper should tighten its budget to the new measurement with the
// same headroom, so the next regression is caught from there.
//
// Under the race detector, which runs this code some ten times slower and
// for which the catch-up alone took a minute, the fixture has a tenth of
// the closed crews and query.Load's budgets shrink with it (budgetShare).
const (
	budgetClosedCrews = 1000 / raceShrink
	budgetLiveCrews   = 4
	// One closed crew in so many has no rollout, and one in so many was
	// discarded (budgetWorkspace).
	budgetUnlocatedEvery = 50
	budgetDiscardedEvery = 10

	// One steady ingest pass, however many closed crews: they are sealed.
	budgetIngestBytes   = 6 << 20
	budgetIngestObjects = 60_000
	budgetIngestTime    = 5 * time.Second
	// One query.Load over the same.
	budgetLoadBytes   = 45 << 20
	budgetLoadObjects = 260_000
	budgetLoadTime    = 3 * time.Second
)

// budgetShare is the part of a thousand crews' allocation budget that is
// the fixture's: query.Load's cost grows with the closed crews.
func budgetShare(n uint64) uint64 { return n * budgetClosedCrews / 1000 }

func TestSteadyPassStaysWithinBudget(t *testing.T) {
	f := budgetWorkspace(t)
	ctx := context.Background()
	// Catch-up, then one more pass for what the first one settled (the
	// merges it recorded end git's questions about a crew): the third is
	// the pass every poll runs from then on.
	for i := 0; i < 2; i++ {
		if err := f.ingest.Ingest(ctx); err != nil {
			t.Fatalf("catch-up pass %d: %v", i+1, err)
		}
	}
	var turns int
	if err := f.db.SQL().QueryRow(`SELECT count(*) FROM turn`).Scan(&turns); err != nil {
		t.Fatal(err)
	}
	if turns < f.rollouts {
		// The ingest must have read every crew's transcript once, or a
		// steady pass that read them all again would cost nothing here.
		t.Fatalf("the catch-up recorded %d turns from %d rollouts: the fixture's transcripts were not read", turns, f.rollouts)
	}

	ingest := measure(t, func() error { return f.ingest.Ingest(ctx) })
	load := measure(t, func() error {
		snap, err := query.Load(ctx, f.ws, testHarnesses, nil)
		if err != nil {
			return err
		}
		open, closed := 0, 0
		for _, p := range snap.Projects {
			open += len(p.Crews)
			closed += p.ClosedCrews
		}
		if open != budgetLiveCrews || closed != budgetClosedCrews {
			return fmt.Errorf("query.Load shows %d crews at work and %d closed, want the fixture's %d and %d", open, closed, budgetLiveCrews, budgetClosedCrews)
		}
		return nil
	})
	t.Logf("%d closed crews: a steady ingest pass %s; query.Load %s", budgetClosedCrews, ingest, load)
	ingest.check(t, "a steady ingest pass", budgetIngestBytes, budgetIngestObjects, budgetIngestTime)
	load.check(t, "query.Load", budgetShare(budgetLoadBytes), budgetShare(budgetLoadObjects), budgetLoadTime)
}

// cost is what one call allocated and how long it took.
type cost struct {
	bytes, objects uint64
	took           time.Duration
}

func (c cost) String() string {
	return fmt.Sprintf("%.1f MB in %d objects, %s", float64(c.bytes)/(1<<20), c.objects, c.took.Round(time.Millisecond))
}

func (c cost) check(t *testing.T, what string, bytes, objects uint64, took time.Duration) {
	t.Helper()
	if c.bytes > bytes {
		t.Errorf("%s allocated %.1f MB, over its budget of %.1f MB (CON-263)", what, float64(c.bytes)/(1<<20), float64(bytes)/(1<<20))
	}
	if c.objects > objects {
		t.Errorf("%s allocated %d objects, over its budget of %d (CON-263)", what, c.objects, objects)
	}
	if limit := took * raceTime; c.took > limit {
		t.Errorf("%s took %s, over its budget of %s (CON-263)", what, c.took.Round(time.Millisecond), limit)
	}
}

// measure runs call once and reads what it allocated. ReadMemStats stops
// the world, which a test can afford and the core's own recording cannot.
func measure(t *testing.T, call func() error) cost {
	t.Helper()
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	start := time.Now()
	err := call()
	took := time.Since(start)
	runtime.ReadMemStats(&after)
	if err != nil {
		t.Fatal(err)
	}
	return cost{bytes: after.TotalAlloc - before.TotalAlloc, objects: after.Mallocs - before.Mallocs, took: took}
}

type budgetFixture struct {
	ws     *store.Workspace
	db     *db.DB
	ingest *timeline.Ingester
	// rollouts is how many crews have one, each of a single turn.
	rollouts int
}

// budgetWorkspace is two projects holding budgetClosedCrews Codex crews
// closed days ago and budgetLiveCrews at work, each with what a real one
// leaves: its meta, its status lines, a brief, the lines the Mate sent it,
// and a rollout of one turn with two tool calls (testdata/codex-turn.jsonl,
// cut from internal/timeline's captured run) filed in a Codex sessions
// directory as Codex files it, its session_meta naming the crew's session
// and worktree. As on the captain's workspace:
//
//   - a crew's meta names its session but not its rollout (Codex has no hook
//     that records it), so the ingest finds each rollout by its session id
//     once and then by what it recorded;
//   - one closed crew in budgetUnlocatedEvery never got a rollout (it died
//     before Codex wrote one), so every pass still tries to adopt one for it
//     over every rollout filed since its launch: twenty such crews made the
//     1.2 GB pass of CON-255;
//   - one finished crew in budgetDiscardedEvery was discarded, its branch
//     gone and never merged, so every pass asks the timeline for its last
//     commit (CON-254); git answers every other branch there and merged.
//
// Git is a stand-in (budgetGit): no git process runs.
func budgetWorkspace(t *testing.T) budgetFixture {
	t.Helper()
	rollout, err := os.ReadFile(filepath.Join("testdata", "codex-turn.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	head, body, ok := strings.Cut(string(rollout), "\n")
	if !ok {
		t.Fatal("testdata/codex-turn.jsonl has one line")
	}
	var sessionMeta map[string]any
	if err := json.Unmarshal([]byte(head), &sessionMeta); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	ws, err := store.Init(root, store.Defaults{MateHarness: "claude", CrewHarness: "codex"})
	if err != nil {
		t.Fatalf("store.Init: %v", err)
	}
	projects := []string{"shop", "web"}
	for _, p := range projects {
		repo := filepath.Join(ws.ProjectHome(p), p)
		if err := os.MkdirAll(repo, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := ws.AddProject(p, store.ProjectConfig{Repos: []store.RepoConfig{{Path: repo, DefaultBranch: "main"}}}); err != nil {
			t.Fatalf("AddProject %s: %v", p, err)
		}
	}
	sessions := filepath.Join(root, "codex-sessions")
	sent := map[string][]string{}
	rollouts := 0
	// Weeks ago, so every closed crew is past the ingest's settledFor.
	epoch := time.Now().Add(-30 * 24 * time.Hour).UTC().Truncate(time.Second)
	for i := 0; i < budgetClosedCrews+budgetLiveCrews; i++ {
		project := projects[i%len(projects)]
		crew := fmt.Sprintf("t%04d", i)
		live := i >= budgetClosedCrews
		started := epoch.Add(time.Duration(i) * 30 * time.Minute)
		if live {
			started = time.Now().Add(-10 * time.Minute).UTC().Truncate(time.Second)
		}
		worktree := ".worktrees/" + project + "-" + crew
		meta := map[string]string{
			"harness":     "codex",
			"agent":       "crew-" + crew,
			"task":        "Task " + crew + ": add the Buy button to the checkout page",
			"branch":      "mate/" + crew,
			"worktree":    worktree,
			"transcript":  "",
			"started_at":  started.Format(time.RFC3339),
			"launched_at": started.Format(time.RFC3339),
			"state":       "spawned",
		}
		if !live {
			meta["state"] = "finished"
			meta["stopped_at"] = started.Add(20 * time.Minute).Format(time.RFC3339)
			if i%budgetDiscardedEvery == budgetDiscardedEvery/2 {
				meta["branch"] = "discarded/" + crew
			}
		}
		if live || i%budgetUnlocatedEvery != 0 {
			session := fmt.Sprintf("01a0b944-33bb-7503-9db7-%012d", i)
			meta["session_id"] = session
			rollouts++
			// Codex files a rollout under its local date and time, after
			// the launch.
			at := started.Add(time.Minute).Local()
			payload := sessionMeta["payload"].(map[string]any)
			payload["session_id"], payload["id"] = session, session
			payload["cwd"] = filepath.Join(ws.Root(), filepath.FromSlash(worktree))
			sessionMeta["timestamp"] = at.UTC().Format(time.RFC3339Nano)
			line, err := json.Marshal(sessionMeta)
			if err != nil {
				t.Fatal(err)
			}
			name := "rollout-" + at.Format("2006-01-02T15-04-05") + "-" + session + ".jsonl"
			writeBudgetFile(t, filepath.Join(sessions, at.Format("2006/01/02"), name), string(line)+"\n"+body)
		}
		// Written as the store writes them, without its fsync per file:
		// through WriteCrewMeta and AppendStatus the fixture took 7 s.
		keys := make([]string, 0, len(meta))
		for k := range meta {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var m strings.Builder
		for _, k := range keys {
			m.WriteString(k + "=" + meta[k] + "\n")
		}
		writeBudgetFile(t, ws.CrewMeta(project, crew), m.String())
		writeBudgetFile(t, ws.CrewBrief(project, crew), "# "+meta["task"]+"\n\nRead PROJECT.md, change checkout.html, run the tests.\n")
		lines := []string{"working: reading the brief", "working: tests pass", "wait-mate: ready in branch " + meta["branch"]}
		if live {
			lines = lines[:1]
		}
		writeBudgetFile(t, ws.CrewStatus(project, crew), strings.Join(lines, "\n")+"\n")
		sent[project] = append(sent[project],
			sentLine(started, store.SourceApp, store.TargetMate, "[assign] "+meta["task"]),
			sentLine(started.Add(5*time.Minute), store.SourceMate, store.CrewTarget(crew), "Use pages/checkout.html."))
	}
	for _, p := range projects {
		writeBudgetFile(t, ws.SentLog(p), strings.Join(sent[p], ""))
		if entries, _, err := ws.ReadSent(p, 0); err != nil || len(entries) != len(sent[p]) {
			t.Fatalf("ReadSent %s = %d entries, %v; want the %d written", p, len(entries), err, len(sent[p]))
		}
	}
	if meta, err := ws.ReadCrewMeta("web", "t0001"); err != nil || meta["state"] != "finished" {
		t.Fatalf("ReadCrewMeta = %v, %v; want the meta written", meta, err)
	}
	handle, err := db.Open(ws)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = handle.Close() })
	ingest := timeline.New(ws, handle, timeline.Deps{
		Harnesses:       catalog.Default(),
		TranscriptRoots: map[harness.Kind]string{"codex": sessions, "claude": filepath.Join(root, "claude-projects")},
		Git:             gitx.Git{Runner: budgetGit{}},
	})
	return budgetFixture{ws: ws, db: handle, ingest: ingest, rollouts: rollouts}
}

func writeBudgetFile(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// sentLine is one line of `sent.log` as store.AppendSent writes it.
func sentLine(at time.Time, source, target, text string) string {
	return strings.Join([]string{at.Format(time.RFC3339), source, target, text}, "\t") + "\n"
}

// budgetGit answers git as a repository would in which every crew's branch
// is there and merged into main with no commit of its own, but a discarded
// crew's, which is gone: `show-ref` and `merge-base --is-ancestor` say yes,
// `rev-parse <branch>` names a commit, `log` lists nothing.
type budgetGit struct{}

func (budgetGit) Run(_ context.Context, cmd gitx.Command) (gitx.Result, error) {
	args := strings.Join(cmd.Args, " ")
	switch {
	case strings.HasPrefix(args, "show-ref ") && strings.Contains(args, "refs/heads/discarded/"):
		return gitx.Result{ExitCode: 1}, nil
	case len(cmd.Args) == 2 && cmd.Args[0] == "rev-parse":
		return gitx.Result{Stdout: fmt.Sprintf("%x\n", sha1.Sum([]byte(cmd.Args[1])))}, nil
	}
	return gitx.Result{}, nil
}
