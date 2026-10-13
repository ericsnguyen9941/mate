package fixture_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness/claude"
	"github.com/nguyenngocanh94/mate/internal/harness/codex"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/screen/fixture"
)

func stateCapture(t *testing.T, kind, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "harness", kind, "testdata", "startup", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// The stub names a state only for what a harness profile measured; every
// other screen is unknown, which is where Jev is asked.
func TestFixtureStateFromTheMeasuredStartup(t *testing.T) {
	o := fixture.New().(screen.StateObserver)
	cases := []struct {
		kind, file string
		want       screen.PaneState
	}{
		{"codex", "codex-0.162.1-hooks-review.txt", screen.PaneHooksReview},
		{"codex", "codex-0.162.1-loading-composer.txt", screen.PaneReady},
		{"codex", "codex-0.156.1-trust-dialog.txt", screen.PaneTrustDir},
		{"codex", "codex_update_dialog.txt", screen.PaneUpdate},
		{"claude", "claude-2.1.270-trust-dialog.txt", screen.PaneTrustDir},
		{"claude", "claude-2.1.270-ready.txt", screen.PaneReady},
		{"claude", "claude-2.1.285-bypass-dialog.txt", screen.PaneBypass},
	}
	for _, c := range cases {
		profile := codex.Codex{}.Screen()
		if c.kind == "claude" {
			profile = claude.Claude{}.Screen()
		}
		r, err := o.ObserveState(context.Background(), profile, stateCapture(t, c.kind, c.file))
		if err != nil || r.State != c.want || r.Source != fixture.Source || r.Confidence != 1 {
			t.Errorf("%s: %+v %v, want %s from the fixture at 1", c.file, r, err, c.want)
		}
	}
	r, _ := o.ObserveState(context.Background(), codex.Codex{}.Screen(), "welcome to something nobody measured\n  a) yes\n  b) no\n")
	if r.State != screen.PaneUnknown || r.Confidence != 0 {
		t.Errorf("unmeasured screen = %+v, want unknown at 0", r)
	}
}
