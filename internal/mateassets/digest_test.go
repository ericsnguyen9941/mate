package mateassets

import (
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/autopilot"
	"github.com/nguyenngocanh94/mate/internal/send"
)

// TestWakeSkillTeachesTheOrderAndTheSentinel holds the skill the Mate loads
// when an app event wakes it to what it must decide: the sentinel that
// marks an app event, `mate wake` first, the fixed blocked → needs-decision
// → wait-mate order, and the end of the turn after acting.
//
// The grammar of the digest line itself (the 120-rune excerpt, `+N more`,
// `quiet for`) is a contract between the daemon and `mate wake`, held by
// internal/autopilot's and cmd/mate's wake tests; the Mate decides from what
// `wake` prints, never from the line, so the skill does not teach it.
func TestWakeSkillTeachesTheOrderAndTheSentinel(t *testing.T) {
	text := skillText(t, wakeSkill)
	p := fixedParams()
	for _, want := range []string{send.Marker, "`" + p.MatevBin + " wake " + p.ProjectName + "`", "stop your turn"} {
		if !strings.Contains(text, want) {
			t.Errorf("%s does not say %q", wakeSkill, want)
		}
	}
	// The verbs are the exported ItemKind constants, so a rename in
	// internal/autopilot fails this test at compile time.
	var at []int
	for _, item := range []struct {
		kind  autopilot.ItemKind
		place string
	}{{autopilot.ItemBlocked, "first"}, {autopilot.ItemNeedsDecision, "second"}, {autopilot.ItemWaitMate, "third"}} {
		want := string(item.kind) + "` " + item.place
		i := strings.Index(text, want)
		if i < 0 {
			t.Fatalf("%s does not spell out %q of the blocked/needs-decision/wait-mate order", wakeSkill, want)
		}
		at = append(at, i)
	}
	if !(at[0] < at[1] && at[1] < at[2]) {
		t.Errorf("%s lists blocked/needs-decision/wait-mate out of order (positions %v)", wakeSkill, at)
	}
}
