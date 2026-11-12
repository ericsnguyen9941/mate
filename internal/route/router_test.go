package route

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/autopilot"
	"github.com/nguyenngocanh94/mate/internal/notice"
	"github.com/nguyenngocanh94/mate/internal/store"
)

const project = "shop"

type rig struct {
	t       *testing.T
	ws      *store.Workspace
	r       *Router
	landed  []string
	landErr error
}

func newRig(t *testing.T, mode string) *rig {
	t.Helper()
	ws, err := store.Init(t.TempDir(), store.Defaults{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.AddProject(project, store.ProjectConfig{Repos: []store.RepoConfig{{Path: "shop/shop"}}, Review: store.ReviewOff, Autoroute: mode}); err != nil {
		t.Fatal(err)
	}
	g := &rig{t: t, ws: ws}
	now := time.Date(2026, 11, 12, 10, 0, 0, 0, time.UTC)
	g.r = New(ws, &fakeJudge{label: notice.LabelRoutine, conf: .95}, Exec{
		Land: func(_ context.Context, p, crew, reviewer string) (string, error) {
			g.landed = append(g.landed, crew)
			return "merged\ntracker: HU-1 not updated; mark it done yourself\nnext: queued: 1 (first: HU-2)", g.landErr
		},
	}, func() time.Time { return now })
	return g
}

// handBack writes a ship crew whose hand-back is clean, and returns its
// wait-mate item.
func (g *rig) handBack(crew string) autopilot.Item {
	g.t.Helper()
	if err := g.ws.WriteCrewMeta(project, crew, map[string]string{"state": "spawned", "kind": "ship"}); err != nil {
		g.t.Fatal(err)
	}
	path := g.ws.CrewHandback(project, crew)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		g.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(cleanHandback), 0o644); err != nil {
		g.t.Fatal(err)
	}
	return autopilot.Item{Kind: autopilot.ItemWaitMate, Crew: crew, Text: "PR open", File: g.ws.CrewStatus(project, crew), Offset: 0}
}

func (g *rig) round(items ...autopilot.Item) (kept []autopilot.Item) {
	g.r.Plan(context.Background(), project, items)
	kept, claimed := g.r.Claim(project, items)
	g.r.Execute(context.Background(), project, claimed)
	return kept
}

func (g *rig) log() []Record {
	g.t.Helper()
	var out []Record
	for _, l := range strings.Split(read(g.ws.RouteLog(project)), "\n") {
		if l == "" {
			continue
		}
		var rec Record
		if err := jsonUnmarshal(l, &rec); err != nil {
			g.t.Fatal(err)
		}
		out = append(out, rec)
	}
	return out
}

// In shadow mode the Mate still gets every hand-back, and the log says what
// the app would have done, once.
func TestShadowModeRecordsAndTakesNothing(t *testing.T) {
	g := newRig(t, store.AutorouteShadow)
	item := g.handBack("k3")

	for i := 0; i < 2; i++ {
		if kept := g.round(item); len(kept) != 1 {
			t.Fatalf("round %d kept %d item(s), want the hand-back", i, len(kept))
		}
	}
	if len(g.landed) != 0 {
		t.Fatalf("landed %v in shadow mode", g.landed)
	}
	log := g.log()
	if len(log) != 1 || log[0].State != StateShadow || log[0].Action != ActionLand {
		t.Fatalf("route log = %+v, want one shadow land", log)
	}
}

// In on mode a clean hand-back is landed, once, the Mate never sees it, and
// the note that tells it wakes it to move the work on.
func TestOnModeLandsAndNotes(t *testing.T) {
	g := newRig(t, store.AutorouteOn)
	item := g.handBack("k3")

	if kept := g.round(item); len(kept) != 0 {
		t.Fatalf("kept %+v, want the hand-back taken", kept)
	}
	if kept := g.round(item); len(kept) != 0 {
		t.Fatalf("second round kept %+v", kept)
	}
	if len(g.landed) != 1 || g.landed[0] != "k3" {
		t.Fatalf("landed %v, want k3 once", g.landed)
	}
	notes, _, err := g.ws.ReadAppNotes(project, 0)
	if err != nil || len(notes) != 1 || !notes[0].Wake ||
		!strings.Contains(notes[0].Text, "landed k3 from its hand-back") ||
		!strings.Contains(notes[0].Text, "left for you: tracker: HU-1 not updated") ||
		!strings.Contains(notes[0].Text, "next: queued: 1") {
		t.Fatalf("notes = %+v, %v", notes, err)
	}
	if log := g.log(); len(log) != 2 || log[0].State != StateClaimed || log[1].State != StateDone {
		t.Fatalf("route log = %+v, want claimed then done", log)
	}
}

// A landing that fails gives the hand-back back to the Mate.
func TestAFailedLandingGoesToTheMate(t *testing.T) {
	g := newRig(t, store.AutorouteOn)
	g.landErr = errors.New("merge refused: checks failing")
	item := g.handBack("k3")

	g.round(item)
	if kept := g.round(item); len(kept) != 1 {
		t.Fatalf("kept %+v after a failed landing, want the hand-back for the Mate", kept)
	}
	if len(g.landed) != 1 {
		t.Fatalf("landed %d time(s), want one try", len(g.landed))
	}
	if notes, _, _ := g.ws.ReadAppNotes(project, 0); len(notes) != 0 {
		t.Fatalf("notes = %+v, want none for a failed landing", notes)
	}
}

// Off is the default: nothing is decided or recorded.
func TestOffModeIsInert(t *testing.T) {
	g := newRig(t, "")
	item := g.handBack("k3")
	if kept := g.round(item); len(kept) != 1 || len(g.landed) != 0 || len(g.log()) != 0 {
		t.Fatalf("off mode: kept %d, landed %v, log %+v", len(kept), g.landed, g.log())
	}
}

// Other items pass through untouched.
func TestOtherItemsPassThrough(t *testing.T) {
	g := newRig(t, store.AutorouteOn)
	q := autopilot.Item{Kind: autopilot.ItemNeedsDecision, Crew: "k9", Text: "A or B?", File: "x", Offset: 3}
	if kept := g.round(q); len(kept) != 1 || kept[0].Crew != "k9" {
		t.Fatalf("kept %+v", kept)
	}
}

func jsonUnmarshal(s string, v any) error { return json.Unmarshal([]byte(s), v) }

// A reviewer's hand-back is read as a review of the crew its review.meta
// names: findings Jev finds in scope go back to that crew, and the Mate is
// told without being woken.
func TestOnModeSendsAReviewsFindingsBack(t *testing.T) {
	g := newRig(t, store.AutorouteOn)
	g.r.judge = &fakeJudge{label: notice.LabelInScopeFix, conf: .9}
	var sent []string
	g.r.exec.Send = func(_ context.Context, p, crew, text string) error {
		sent = append(sent, crew+": "+text)
		return nil
	}
	if err := g.ws.WriteCrewMeta(project, "hu-4-r1", map[string]string{"state": "spawned", "kind": "scout"}); err != nil {
		t.Fatal(err)
	}
	report := g.ws.CrewReport(project, "hu-4-r1")
	if err := os.MkdirAll(filepath.Dir(report), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(report, []byte(changesReview), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(g.ws.CrewReviewMeta(project, "hu-4-r1"), []byte("crew=hu-4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	item := autopilot.Item{Kind: autopilot.ItemWaitMate, Crew: "hu-4-r1", File: g.ws.CrewStatus(project, "hu-4-r1")}

	if kept := g.round(item); len(kept) != 0 {
		t.Fatalf("kept %+v, want the review taken", kept)
	}
	if len(sent) != 1 || !strings.HasPrefix(sent[0], "hu-4: review hu-4-r1 asks for changes") || !strings.Contains(sent[0], report) {
		t.Fatalf("sent %q", sent)
	}
	notes, _, _ := g.ws.ReadAppNotes(project, 0)
	if len(notes) != 1 || notes[0].Wake || !strings.Contains(notes[0].Text, "sent hu-4 the findings of review hu-4-r1") {
		t.Fatalf("notes = %+v, want one that rides along", notes)
	}
}
