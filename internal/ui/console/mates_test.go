package console

import (
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/query"
)

// matesTree is sampleTree with a running Mate on its second Project too, so
// the Mates list has two Mates to tell apart.
func matesTree() query.Snapshot {
	tree := sampleTree()
	ledger := tree.Projects[0].Mate
	ledger.Designated.Value.MateID = "mate_ledger"
	ledger.AgentName = query.KnownField("mate-ledger-worker")
	ledger.Binding.Value.AgentName = "mate-ledger-worker"
	tree.Projects[1].Mate = ledger
	return tree
}

func TestVListsOneMateRowPerProjectAndNoCrews(t *testing.T) {
	m := loaded(t, matesTree(), nil)
	m, _ = send(t, m, key("v"))
	if m.cur().kind != frameMates || len(m.stack) != 1 {
		t.Fatalf("stack after v = %+v, want the Mates list alone", m.stack)
	}
	rows := m.currentRows()
	if len(rows) != 2 || rows[0].kind != rowMate || rows[1].kind != rowMate {
		t.Fatalf("rows = %+v, want two Mate rows", rows)
	}
	frame := renderFrame(t, m)
	for _, want := range []string{"acme " + unicodeGlyphs.Dot + " mates", unicodeGlyphs.Mate + " payments-api", unicodeGlyphs.Mate + " ledger-worker", "v projects"} {
		if !strings.Contains(frame, want) {
			t.Errorf("the Mates list does not draw %q:\n%s", want, frame)
		}
	}
	_, h := bxSlot(t, m, slotList).body()
	for _, l := range m.listBody(m.plan(), h).lines {
		if text := l.render(40, m.p); strings.Contains(text, unicodeGlyphs.Crew) {
			t.Errorf("the Mates list draws a crew: %q", text)
		}
	}
}

// One Enter (or click) shows another Project's Mate: no Esc, no opening
// the Project first.
func TestEnterOnTheMatesListShowsThatProjectsMate(t *testing.T) {
	spy := &stageSpy{}
	m := loaded(t, matesTree(), nil).WithStage(spy.fn)
	m, _ = send(t, m, key("v"))
	m, _ = send(t, m, key("down"))
	m, cmd := send(t, m, key("enter"))
	if cmd == nil {
		t.Fatal("Enter on a Mate row returned no Cmd")
	}
	m, _ = send(t, m, cmd())
	if len(spy.calls) != 1 {
		t.Fatalf("stage calls = %d, want 1", len(spy.calls))
	}
	got := spy.calls[0]
	if got.Kind != StageMate || got.ID != "mate_ledger" || got.AgentName != "mate-ledger-worker" || got.ProjectID != "proj_01J9M4H7K9L1M3N5P7Q9R1S3TU" {
		t.Fatalf("stage target = %+v, want the ledger-worker Mate", got)
	}
	if m.cur().kind != frameMates || len(m.stack) != 1 {
		t.Fatalf("stack after the stage = %+v, want the Mates list", m.stack)
	}
	if m.staged.name != "ledger-worker" {
		t.Fatalf("staged = %+v, want the ledger-worker Mate", m.staged)
	}
}

func TestClickOnTheMatesListShowsThatMate(t *testing.T) {
	spy := &stageSpy{}
	m := loaded(t, matesTree(), nil).WithStage(spy.fn)
	m, _ = send(t, m, key("v"))
	m, cmd := send(t, m, press(6, bxRowLine(t, m, 1)))
	if cmd == nil {
		t.Fatal("the click showed nothing")
	}
	cmd()
	if len(spy.calls) != 1 || spy.calls[0].ID != "mate_ledger" {
		t.Fatalf("stage calls = %+v, want the ledger-worker Mate", spy.calls)
	}
}

// A Project without a Mate keeps its row and says why it cannot be shown.
func TestTheMatesListRefusesAProjectWithoutAMate(t *testing.T) {
	spy := &stageSpy{}
	m := loaded(t, sampleTree(), nil).WithStage(spy.fn)
	m, _ = send(t, m, key("v"))
	m, _ = send(t, m, key("down"))
	m, cmd := send(t, m, key("enter"))
	if cmd != nil || len(spy.calls) != 0 {
		t.Fatalf("a Project with no Mate was staged: %+v", spy.calls)
	}
	if !strings.Contains(m.msg.text, "this Project has no Mate") {
		t.Fatalf("status = %q, want the no-Mate refusal", m.msg.text)
	}
}

// v keeps the Project the reader was on, both ways and from inside one.
func TestVKeepsTheSelectedProject(t *testing.T) {
	m := loaded(t, matesTree(), nil)
	m, _ = send(t, m, key("down"))
	m, _ = send(t, m, key("v"))
	if r, _ := m.selectedRow(); r.id != "proj_01J9M4H7K9L1M3N5P7Q9R1S3TU" {
		t.Fatalf("selection on the Mates list = %+v, want ledger-worker", r)
	}
	if p := m.currentProject(); p.Name != "ledger-worker" {
		t.Fatalf("currentProject = %q, want the selected Mate's Project", p.Name)
	}
	m, _ = send(t, m, key("v"))
	if r, _ := m.selectedRow(); m.cur().kind != frameWorkspace || r.kind != rowProject || r.id != "proj_01J9M4H7K9L1M3N5P7Q9R1S3TU" {
		t.Fatalf("back on the Projects: frame %v row %+v, want ledger-worker", m.cur().kind, r)
	}

	m = projectFrame(t, matesTree())
	m, _ = send(t, m, key("v"))
	if r, _ := m.selectedRow(); m.cur().kind != frameMates || len(m.stack) != 1 || r.id != "proj_01J9M1F8K2Q7C4H6N0R3V5T8YZ" {
		t.Fatalf("v inside payments-api: stack %+v row %+v, want its Mate on the Mates list", m.stack, r)
	}
}

// The Mates list survives a re-read and stays the root a jump returns to.
func TestTheMatesListIsTheRootUntilVAgain(t *testing.T) {
	m := loaded(t, matesTree(), nil)
	m, _ = send(t, m, key("v"))
	m, _ = send(t, m, key("down"))
	m, cmd := m.startLoad()
	m, _ = send(t, m, cmd())
	if r, _ := m.selectedRow(); m.cur().kind != frameMates || r.id != "proj_01J9M4H7K9L1M3N5P7Q9R1S3TU" {
		t.Fatalf("after a re-read: frame %v row %+v, want ledger-worker on the Mates list", m.cur().kind, r)
	}
	crew := matesTree().Projects[0].Crews[1].CrewID
	m, ok := m.jumpToCrew(crew)
	if !ok || len(m.stack) != 2 || m.stack[0].kind != frameMates {
		t.Fatalf("jump to a crew: stack %+v, want Mates -> Project", m.stack)
	}
	m, _ = send(t, m, key("esc"))
	if r, _ := m.selectedRow(); m.cur().kind != frameMates || r.id != "proj_01J9M1F8K2Q7C4H6N0R3V5T8YZ" {
		t.Fatalf("Esc from the crew's Project: frame %v row %+v, want payments-api on the Mates list", m.cur().kind, r)
	}
}
