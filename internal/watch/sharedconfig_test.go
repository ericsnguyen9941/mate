package watch_test

import (
	"context"
	"os"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/watch"
)

// Under the engine the observer shares the engine's handle, which reloads
// workspace.yaml itself (CON-270): a round with SharedConfig does not read
// the file, so a file that no longer parses fails no round; without it the
// round reads it, as before.
func TestASharedConfigRoundDoesNotReadWorkspaceYAML(t *testing.T) {
	f := newFixture(t)
	if err := os.WriteFile(f.ws.WorkspaceFile(), []byte("projects: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := f.w.Poll(context.Background()); err == nil {
		t.Fatal("a round read a broken workspace.yaml without a word")
	}
	d := f.deps()
	d.SharedConfig = true
	if err := watch.New(f.ws, d).Poll(context.Background()); err != nil {
		t.Fatalf("a round over a shared handle read workspace.yaml: %v", err)
	}
}
