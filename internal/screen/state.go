package screen

import (
	"context"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// PaneState is what a pane is waiting for, from the list mate defines
// (docs/plans/pane-state-jev-2026-10-13.md section 3). An observer names one;
// the caller decides what to do.
type PaneState string

const (
	PaneLoading        PaneState = "loading"
	PaneReady          PaneState = "ready"
	PaneBusy           PaneState = "busy"
	PaneDraft          PaneState = "draft"
	PaneTrustDir       PaneState = "dialog_trust_dir"
	PaneUpdate         PaneState = "dialog_update"
	PaneHooksReview    PaneState = "dialog_hooks_review"
	PaneBypass         PaneState = "dialog_bypass"
	PaneToolPermission PaneState = "dialog_tool_permission"
	PaneOther          PaneState = "dialog_other"
	PaneAuthRequired   PaneState = "auth_required"
	PaneQuotaExhausted PaneState = "quota_exhausted"
	PaneExited         PaneState = "exited"
	PaneUnknown        PaneState = "unknown"
)

// PaneStates lists every state in the order of the spec's table.
var PaneStates = []PaneState{PaneLoading, PaneReady, PaneBusy, PaneDraft, PaneTrustDir, PaneUpdate,
	PaneHooksReview, PaneBypass, PaneToolPermission, PaneOther, PaneAuthRequired, PaneQuotaExhausted,
	PaneExited, PaneUnknown}

// StateReading is one observer's answer about a pane.
type StateReading struct {
	State      PaneState
	Startup    harness.StartupScreen // the fixture's classification; "" from Jev
	Confidence float64
	Source     string // "fixture" or "jev"
	Reason     string
}

// StateObserver names the PaneState of a pane. The startup settle asks it;
// the fixture observer answers from the harness's measured startup
// classifier, and the chain asks Jev only when the fixture answers
// PaneUnknown.
type StateObserver interface {
	ObserveState(ctx context.Context, profile harness.ScreenProfile, pane string) (StateReading, error)
}
