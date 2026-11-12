package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/watch"
)

// A round examines what can change (CON-267): every project with an open
// crew or a live Mate, every project a notice named, and every project on a
// sweep - and nothing else.

// countingWorkspace counts the reads of each project's incidents.log, the
// first thing a round reads of a project it examines.
type countingWorkspace struct {
	*store.Workspace
	mu    sync.Mutex
	reads map[string]int
}

func (c *countingWorkspace) ReadIncidents(project string, from int64) ([]store.IncidentEntry, int64, error) {
	c.mu.Lock()
	c.reads[project]++
	c.mu.Unlock()
	return c.Workspace.ReadIncidents(project, from)
}

func (c *countingWorkspace) readsOf(project string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reads[project]
}

// rounds records what each round handed the timeline, the budget and the
// no-turn checks.
type rounds struct {
	mu                 sync.Mutex
	ingest, budget, nt []watch.Round
}

func (r *rounds) Ingest(_ context.Context, round watch.Round) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ingest = append(r.ingest, round)
	return nil
}

func (r *rounds) CheckBudgets(_ context.Context, round watch.Round) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.budget = append(r.budget, round)
	return nil
}

func (r *rounds) CheckNoTurn(_ context.Context, round watch.Round, _ time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nt = append(r.nt, round)
	return nil
}

func (r *rounds) last(t *testing.T) watch.Round {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.ingest) == 0 {
		t.Fatal("no round reached the timeline")
	}
	last := r.ingest[len(r.ingest)-1]
	if !reflect.DeepEqual(r.budget[len(r.budget)-1], last) || !reflect.DeepEqual(r.nt[len(r.nt)-1], last) {
		t.Fatalf("the checks got rounds %+v and %+v, the timeline %+v: they read what the round examined",
			r.budget[len(r.budget)-1], r.nt[len(r.nt)-1], last)
	}
	return last
}

// feed is a change feed the test fills.
type feed struct {
	mu   sync.Mutex
	next watch.Notices
}

func (f *feed) note(project string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.next.Projects == nil {
		f.next.Projects = map[string]bool{}
	}
	f.next.Projects[project] = true
}

func (f *feed) noteAll() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.next.All = true
}

func (f *feed) take() watch.Notices {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := f.next
	f.next = watch.Notices{}
	return n
}

// scoped is the fixture's workspace with a second project, "quiet", whose
// one crew closed long ago and which has no Mate, observed by a Watcher
// that hears a feed.
type scoped struct {
	*fixture
	cw     *countingWorkspace
	feed   *feed
	rounds *rounds
}

func newScoped(t *testing.T) *scoped {
	t.Helper()
	f := newFixture(t)
	if err := os.MkdirAll(filepath.Join(f.ws.Root(), "quiet", "quiet"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := f.ws.AddProject("quiet", store.ProjectConfig{Repos: []store.RepoConfig{{Path: "quiet/quiet"}}}); err != nil {
		t.Fatalf("AddProject: %v", err)
	}
	if err := f.ws.WriteCrewMeta("quiet", "old1", map[string]string{"harness": "codex", "state": "finished", "stopped_at": "2026-09-01T10:00:00Z"}); err != nil {
		t.Fatal(err)
	}
	s := &scoped{fixture: f, cw: &countingWorkspace{Workspace: f.ws, reads: map[string]int{}}, feed: &feed{}, rounds: &rounds{}}
	deps := f.deps()
	deps.Notices = s.feed.take
	deps.Timeline, deps.Budget, deps.Turns = s.rounds, s.rounds, s.rounds
	f.w = watch.New(s.cw, deps)
	return s
}

func (s *scoped) next() watch.Round {
	s.t.Helper()
	s.clock.advance(5 * time.Second)
	s.poll()
	return s.rounds.last(s.t)
}

func wantRound(t *testing.T, got watch.Round, sweep bool, projects ...string) {
	t.Helper()
	if got.Sweep != sweep || !reflect.DeepEqual(got.Projects, projects) {
		t.Fatalf("round = %+v, want sweep %v over %v", got, sweep, projects)
	}
}

func TestRoundSkipsAnIdleProjectBetweenSweeps(t *testing.T) {
	s := newScoped(t)
	s.poll()
	// The first round looks at everything: nothing has been seen yet.
	wantRound(t, s.rounds.last(t), true, "shop", "quiet")

	for i := 0; i < 10; i++ {
		// The open crew's project every round, the quiet one never.
		wantRound(t, s.next(), false, "shop")
	}
	if got := s.cw.readsOf("quiet"); got != 1 {
		t.Fatalf("incidents.log of the quiet project read %d times in 11 rounds, want once", got)
	}
	if got := s.cw.readsOf("shop"); got != 11 {
		t.Fatalf("incidents.log of the busy project read %d times in 11 rounds, want 11", got)
	}

	// The sweep comes on its schedule: a minute after the last one.
	s.clock.advance(watch.DefaultSweepInterval - 50*time.Second)
	s.poll()
	wantRound(t, s.rounds.last(t), true, "shop", "quiet")
	wantRound(t, s.next(), false, "shop")
}

func TestNoticeExaminesAProjectTheNextRound(t *testing.T) {
	s := newScoped(t)
	s.poll()
	wantRound(t, s.next(), false, "shop")

	s.feed.note("quiet")
	wantRound(t, s.next(), false, "shop", "quiet")
	// Read once, it is idle again.
	wantRound(t, s.next(), false, "shop")

	// A notice about the workspace's own files, or an action the core ran
	// itself, is a sweep.
	s.feed.noteAll()
	wantRound(t, s.next(), true, "shop", "quiet")
}

// A crew spawned in an idle project is a notice (its meta, written through
// the store); from then on the project has an open crew and is examined
// every round, its crew observed like any other.
func TestOpenCrewKeepsAProjectExaminedEveryRound(t *testing.T) {
	s := newScoped(t)
	s.poll()
	wantRound(t, s.next(), false, "shop")

	if err := s.ws.WriteCrewMeta("quiet", "new1", map[string]string{"harness": "codex", "agent": "crew-new1", "state": "spawned"}); err != nil {
		t.Fatal(err)
	}
	s.feed.note("quiet")
	for i := 0; i < 3; i++ {
		wantRound(t, s.next(), false, "shop", "quiet")
	}
	if _, ok := s.w.Health("quiet", "new1"); !ok {
		t.Fatal("the new crew has no health: its project's rounds did not observe it")
	}

	// Its close is heard; the round that reads it is the project's last.
	if err := s.ws.WriteCrewMeta("quiet", "new1", map[string]string{"harness": "codex", "state": "finished", "stopped_at": "2026-09-18T11:00:30Z"}); err != nil {
		t.Fatal(err)
	}
	s.feed.note("quiet")
	wantRound(t, s.next(), false, "shop", "quiet")
	wantRound(t, s.next(), false, "shop")
}

// A live Mate's transcript grows unheard, so its project is examined every
// round, crews or none; its stop is heard.
func TestLiveMateKeepsAProjectExaminedEveryRound(t *testing.T) {
	s := newScoped(t)
	if err := s.ws.WriteMateMeta("quiet", map[string]string{"harness": "claude", "mode": "headless", "agent": "mate-quiet", "pid": "4242"}); err != nil {
		t.Fatal(err)
	}
	s.poll()
	for i := 0; i < 3; i++ {
		wantRound(t, s.next(), false, "shop", "quiet")
	}
	if err := s.ws.WriteMateMeta("quiet", map[string]string{"harness": "claude", "mode": "headless", "stopped_at": "2026-09-18T11:00:30Z"}); err != nil {
		t.Fatal(err)
	}
	s.feed.note("quiet")
	wantRound(t, s.next(), false, "shop", "quiet")
	wantRound(t, s.next(), false, "shop")
}

// Stale is decided per open crew every round, so hearing a feed cannot
// delay it: with no notice and no sweep in sight, it opens in the same
// round as it does without a feed.
func TestStaleOpensOnTimeBetweenSweeps(t *testing.T) {
	s := newScoped(t)
	deps := s.deps()
	deps.Notices = s.feed.take
	deps.SweepInterval = time.Hour
	deps.Timeline, deps.Budget, deps.Turns = s.rounds, s.rounds, s.rounds
	s.w = watch.New(s.cw, deps)
	s.appendStatus("k3", "working: reading the tests")

	s.poll()
	for i := 0; i < 35; i++ {
		// 2m55s of quiet: not yet.
		wantRound(t, s.next(), false, "shop")
	}
	s.assertIncidents()
	wantRound(t, s.next(), false, "shop")
	s.assertIncidents("k3 stale open")
}

// A Watcher with no feed examines every project every round, each round a
// sweep, as before CON-267.
// A round reads only what was appended to incidents.log since the last one
// (openIncidents), and a crew's meta a round found closed is not read again
// until it changes (openCrews): what another process writes there between
// rounds is still seen by the next one.
func TestRoundsSeeWhatOthersWriteBetweenThem(t *testing.T) {
	s := newScoped(t)
	s.writeMeta("c2", map[string]string{"harness": "codex", "state": "finished", "stopped_at": "2026-09-18T10:00:00Z"})
	s.poll()
	wantRound(t, s.next(), false, "shop")
	if _, ok := s.w.Health("shop", "c2"); ok {
		t.Fatal("a closed crew has health")
	}

	// Reopened behind the rounds' back (`mate crew relaunch`).
	s.recordCrew("c2")
	wantRound(t, s.next(), false, "shop")
	if _, ok := s.w.Health("shop", "c2"); !ok {
		t.Fatal("a crew reopened after a round found it closed is not observed")
	}

	// Another process opened an incident on k3; the pane moving clears it,
	// which the round does only if it read the line.
	if err := s.ws.AppendIncident("shop", store.IncidentEntry{Time: s.clock.Now(), Crew: "k3", Kind: "stale", State: store.IncidentOpen, Text: "by hand"}); err != nil {
		t.Fatal(err)
	}
	s.setScreen(codexBusyScreen)
	wantRound(t, s.next(), false, "shop")
	var k3 []string
	for _, e := range s.incidents() {
		if e.Crew == "k3" {
			k3 = append(k3, e.Kind+" "+e.State)
		}
	}
	if want := []string{"stale open", "stale resolved"}; !reflect.DeepEqual(k3, want) {
		t.Fatalf("k3's incidents = %v, want %v", k3, want)
	}
}

func TestNoFeedSweepsEveryRound(t *testing.T) {
	s := newScoped(t)
	deps := s.deps()
	deps.Timeline, deps.Budget, deps.Turns = s.rounds, s.rounds, s.rounds
	s.w = watch.New(s.cw, deps)
	s.poll()
	for i := 0; i < 3; i++ {
		wantRound(t, s.next(), true, "shop", "quiet")
	}
}
