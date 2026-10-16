package brief

import "testing"

func TestParseOpenDecisions(t *testing.T) {
	if got := ParseOpenDecisions("## Open decisions\nnone\n"); got != nil {
		t.Fatalf("none = %+v, want nil", got)
	}
	got := ParseOpenDecisions("## Open decisions\n- Which checkout page? decides: captain\n- Test names. decides: Mate\n- No decider here\n")
	want := []OpenDecision{
		{Text: "Which checkout page? decides: captain", Decides: "captain"},
		{Text: "Test names. decides: Mate", Decides: "mate"},
		{Text: "No decider here"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestMatchDecisionPicksTheEntrySharingMostWords(t *testing.T) {
	ds := []OpenDecision{
		{Text: "Which checkout page the Buy button opens. decides: captain", Decides: "captain"},
		{Text: "Test file names. decides: mate", Decides: "mate"},
	}
	if d, ok := MatchDecision(ds, "what is the checkout page URL for the Buy button?"); !ok || d.Decides != "captain" {
		t.Fatalf("match = %+v, %v; want the checkout entry", d, ok)
	}
	if _, ok := MatchDecision(ds, "the build is red"); ok {
		t.Fatal("matched a question that shares no word")
	}
	tie := []OpenDecision{
		{Text: "Banner colour. decides: captain", Decides: "captain"},
		{Text: "Banner font. decides: mate", Decides: "mate"},
	}
	if _, ok := MatchDecision(tie, "what about the banner?"); ok {
		t.Fatal("a tie between deciders matched")
	}
}
