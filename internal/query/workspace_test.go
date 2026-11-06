package query

import (
	"context"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// A nil *store.Workspace handed to Load arrives inside the Workspace
// interface, which is not nil itself. Load must still answer "no workspace"
// rather than panic in the store (plan section 3, rule 4).
func TestLoadSaysNoWorkspaceForANilWorkspace(t *testing.T) {
	for name, ws := range map[string]Workspace{
		"nil interface":        nil,
		"nil *store.Workspace": (*store.Workspace)(nil),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Load(context.Background(), ws, testHarnesses, nil)
			if err == nil || !strings.Contains(err.Error(), "no workspace is open") {
				t.Fatalf("Load = %v, want the no-workspace error", err)
			}
			_, err = LoadLive(context.Background(), ws, testHarnesses, nil, Liveness{})
			if err == nil || !strings.Contains(err.Error(), "no workspace is open") {
				t.Fatalf("LoadLive = %v, want the no-workspace error", err)
			}
		})
	}
}

func TestLoadBoxIsUnknownForANilWorkspace(t *testing.T) {
	for name, ws := range map[string]BoxWorkspace{
		"nil interface":        nil,
		"nil *store.Workspace": (*store.Workspace)(nil),
	} {
		t.Run(name, func(t *testing.T) {
			got := LoadBox(ws, "demo")
			if got.State != Unknown || got.Reason != "no workspace is open" {
				t.Fatalf("LoadBox = %+v, want Unknown with the no-workspace reason", got)
			}
		})
	}
}
