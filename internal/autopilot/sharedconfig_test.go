package autopilot_test

import (
	"context"
	"os"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/autopilot"
)

// Under the engine the daemon shares the engine's handle, which reloads
// workspace.yaml itself (CON-270): a tick with SharedConfig does not read
// the file; without it the tick reads it, as before.
func TestASharedConfigTickDoesNotReadWorkspaceYAML(t *testing.T) {
	f := newFixture(t)
	if err := os.WriteFile(f.ws.WorkspaceFile(), []byte("projects: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := f.pilot.Tick(context.Background()); err == nil {
		t.Fatal("a tick read a broken workspace.yaml without a word")
	}
	d := f.deps(nil)
	d.SharedConfig = true
	if err := autopilot.New(f.ws, d).Tick(context.Background()); err != nil {
		t.Fatalf("a tick over a shared handle read workspace.yaml: %v", err)
	}
}
