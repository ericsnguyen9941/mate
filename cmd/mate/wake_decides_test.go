package main

import (
	"os"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/brief"
)

// writeAskBrief gives ask-k4 a brief whose `## Open decisions` is body, in
// the shape `mate crew spawn` renders (only `# Task` is read).
func writeAskBrief(t *testing.T, dir, path, body string) {
	t.Helper()
	text := brief.RoleHeading + "\nYou are a ship Crew.\n\n" + brief.TaskHeading + "\n" +
		"## Captain's words\nclean up the API\n\n## Open decisions\n" + body + "\n\n# Setup\nnothing\n"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Replays mihu 2026-10-14 (CON-93): the Crew stopped on an open decision
// its brief gives to the captain. The wake says so where the Mate is
// looking, with what to do instead of answering.
func TestWakeSaysACaptainsDecisionIsTheCaptains(t *testing.T) {
	w := wakeFixture(t)
	writeAskBrief(t, w.CrewDir("shop", "ask-k4"), w.CrewBrief("shop", "ask-k4"),
		"- Rename the users table? decides: mate\n- Keep or remove the old endpoint once nothing calls it. decides: captain")
	text := runWake(t, w)
	for _, want := range []string{
		"needs-decision ask-k4: keep the old endpoint",
		"\n  decides: captain (Keep or remove the old endpoint",
		"\n  next: ask the captain as news; send the Crew nothing until they answer\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("wake lacks %q:\n%s", want, text)
		}
	}
}

func TestWakeSaysAMatesDecisionIsTheMates(t *testing.T) {
	w := wakeFixture(t)
	writeAskBrief(t, w.CrewDir("shop", "ask-k4"), w.CrewBrief("shop", "ask-k4"),
		"- Keep or remove the old endpoint. decides: mate\n- Pick the colour of the banner. decides: captain")
	text := runWake(t, w)
	if !strings.Contains(text, "\n  decides: mate (Keep or remove the old endpoint. decides: mate)\n") {
		t.Errorf("wake does not give the decision to the Mate:\n%s", text)
	}
	if strings.Contains(text, "Held for the captain") {
		t.Errorf("a Mate's decision was sent to the captain:\n%s", text)
	}
}

// No entry shares a word with the question, or there is no brief at all:
// the wake says it does not know and how many entries there are to read.
func TestWakeSaysWhenItCannotTellWhoDecides(t *testing.T) {
	w := wakeFixture(t)
	if text := runWake(t, w); !strings.Contains(text, "\n  decides: unknown (open decisions: 0)\n") {
		t.Errorf("wake without a brief:\n%s", text)
	}
	writeAskBrief(t, w.CrewDir("shop", "ask-k4"), w.CrewBrief("shop", "ask-k4"),
		"- Pick the colour of the banner. decides: captain\n- Rename the users table. decides: mate")
	if text := runWake(t, w); !strings.Contains(text, "\n  decides: unknown (open decisions: 2)\n") {
		t.Errorf("wake with no matching entry:\n%s", text)
	}
}
