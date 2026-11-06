package prwatch_test

import (
	"context"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/outbox"
	"github.com/nguyenngocanh94/mate/internal/prwatch"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// A nil *store.Workspace set on Deps arrives inside the Workspace
// interface, which is not nil itself. Watch must still refuse it rather
// than panic in the store on its first read (plan section 3, rule 4).
func TestWatchRefusesANilWorkspace(t *testing.T) {
	for name, ws := range map[string]prwatch.Workspace{
		"nil interface":        nil,
		"nil *store.Workspace": (*store.Workspace)(nil),
	} {
		t.Run(name, func(t *testing.T) {
			d := prwatch.Deps{WS: ws, Outbox: outbox.New(nil, outbox.Deps{})}
			err := prwatch.Watch(context.Background(), d, project, crew, prURL)
			if err == nil || !strings.Contains(err.Error(), "a workspace and an outbox are required") {
				t.Fatalf("Watch = %v, want the workspace-required error", err)
			}
		})
	}
}
