package spawn_test

import (
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// A nil *store.Workspace handed to ListCrews arrives inside the CrewReader
// interface, which is not nil itself. ListCrews must still refuse it rather
// than panic in the store (plan section 3, rule 4).
func TestListCrewsRefusesANilWorkspace(t *testing.T) {
	for name, w := range map[string]spawn.CrewReader{
		"nil interface":        nil,
		"nil *store.Workspace": (*store.Workspace)(nil),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := spawn.ListCrews(w, "demo")
			if err == nil || !strings.Contains(err.Error(), "a workspace is required") {
				t.Fatalf("ListCrews = %v, want the workspace-required error", err)
			}
		})
	}
}
