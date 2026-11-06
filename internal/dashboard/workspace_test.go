package dashboard

import (
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// A nil *store.Workspace set on Options arrives inside the Workspace
// interface, which is not nil itself. New must still refuse it rather than
// build a server that panics on its first read (plan section 3, rule 4).
func TestNewRefusesANilWorkspace(t *testing.T) {
	f := newFixture(t)
	for name, ws := range map[string]Workspace{
		"nil interface":        nil,
		"nil *store.Workspace": (*store.Workspace)(nil),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := New(Options{Workspace: ws, DB: f.read})
			if err == nil || !strings.Contains(err.Error(), "a workspace is required") {
				t.Fatalf("New = %v, want the workspace-required error", err)
			}
		})
	}
}
