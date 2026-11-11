package runtime_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/runtime"
)

// roundFake is a fake Herdr with n idle agents in session "lab", each in a
// pane of its own, and the handles a caller recorded for them.
func roundFake(n int) (*runtime.Fake, []runtime.AgentHandle) {
	fake := runtime.NewFake()
	session := runtime.SessionHandle{Name: "lab"}
	handles := make([]runtime.AgentHandle, 0, n)
	for i := 0; i < n; i++ {
		h := runtime.AgentHandle{
			Session: session,
			Name:    fmt.Sprintf("crew-%d", i),
			Tab:     runtime.TabHandle{Session: session, PaneID: fmt.Sprintf("w1:p%d", i)},
		}
		fake.PutAgent(h, runtime.AgentIdle)
		fake.SetReadOutput(h, "› idle\n")
		handles = append(handles, h)
	}
	return fake, handles
}

func count(calls []string, op string) int {
	n := 0
	for _, c := range calls {
		if c == op {
			n++
		}
	}
	return n
}

// TestRoundCacheAsksOnceARound: every InspectAgent of a round, however many
// agents and however many times each is asked, is one `agent list`; the
// session is one `session list`; and the next round asks again.
func TestRoundCacheAsksOnceARound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fake, handles := roundFake(10)
	cache := runtime.NewRoundCache(fake)
	id, err := runtime.ParseWorkspaceID("ws_test1")
	if err != nil {
		t.Fatal(err)
	}
	spec := runtime.SessionSpec{WorkspaceID: id, Name: "lab", ConfigHome: t.TempDir()}

	for pass := 0; pass < 2; pass++ { // the watcher, then the ingest
		for _, h := range handles {
			if _, running, err := cache.LookupSession(ctx, spec); err != nil || !running {
				t.Fatalf("LookupSession = %v, %v", running, err)
			}
			got, err := cache.InspectAgent(ctx, h)
			if err != nil {
				t.Fatalf("InspectAgent(%s): %v", h.Name, err)
			}
			if got.Handle.Name != h.Name || !got.HasRevision || !got.LiveHandleOK {
				t.Fatalf("InspectAgent(%s) = %+v", h.Name, got)
			}
		}
	}
	if len(fake.Calls) != 2 || count(fake.Calls, "LookupSession") != 1 || count(fake.Calls, "ListAgents") != 1 {
		t.Fatalf("calls in one round = %v, want one session lookup and one agent list", fake.Calls)
	}

	cache.BeginRound()
	if _, err := cache.InspectAgent(ctx, handles[0]); err != nil {
		t.Fatal(err)
	}
	if count(fake.Calls, "ListAgents") != 2 {
		t.Fatalf("calls = %v, want the next round to list again", fake.Calls)
	}
}

// TestRoundCacheAsksForAnAgentTheListDoesNotName: absence from the list is
// not taken as absence. The agent is asked for directly, once a round, and
// that answer - agent_not_found or the agent - is the caller's, as it was
// before the list.
func TestRoundCacheAsksForAnAgentTheListDoesNotName(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fake, handles := roundFake(2)
	cache := runtime.NewRoundCache(fake)

	fake.DropAgent(handles[0])
	for i := 0; i < 2; i++ {
		if _, err := cache.InspectAgent(ctx, handles[0]); !runtime.IsAgentNotFound(err) {
			t.Fatalf("InspectAgent of a gone agent = %v, want agent_not_found", err)
		}
	}
	if count(fake.Calls, "InspectAgent") != 1 {
		t.Fatalf("calls = %v, want one direct ask for the missing agent", fake.Calls)
	}

	// A list that leaves out an agent Herdr still has (Fake.ListOmit) does
	// not make it gone.
	cache.BeginRound()
	fake.ListOmit = map[string]bool{"lab/crew-1": true}
	got, err := cache.InspectAgent(ctx, handles[1])
	if err != nil || got.Handle.Name != "crew-1" {
		t.Fatalf("InspectAgent of an agent the list omits = %+v, %v; want it found", got, err)
	}
}

// TestRoundCacheAnswersTheRoundWithAFailedList: a list Herdr could not
// answer is the answer for every agent of the round, rather than one more
// bounded wait per agent.
func TestRoundCacheAnswersTheRoundWithAFailedList(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fake, handles := roundFake(3)
	fake.ListErr = errors.New("herdr: connection refused")
	cache := runtime.NewRoundCache(fake)
	for _, h := range handles {
		if _, err := cache.InspectAgent(ctx, h); !errors.Is(err, fake.ListErr) {
			t.Fatalf("InspectAgent(%s) = %v, want the list's error", h.Name, err)
		}
	}
	if len(fake.Calls) != 1 {
		t.Fatalf("calls = %v, want the one failed list", fake.Calls)
	}
}

// TestRoundCacheSeesTheAgentAgainstTheRecordedHandle: an agent listed in
// another pane than the one recorded is not where the caller looks, as
// Herdr's `agent get` reports it; and the pane revision comes through.
func TestRoundCacheSeesTheAgentAgainstTheRecordedHandle(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fake, handles := roundFake(1)
	cache := runtime.NewRoundCache(fake)
	recorded := handles[0]
	recorded.RawID = "k3"
	recorded.Tab.PaneID = "w1:p9"
	got, err := cache.InspectAgent(ctx, recorded)
	if err != nil {
		t.Fatal(err)
	}
	if got.LiveHandleOK || got.Handle.Tab.PaneID != "w1:p0" || got.Handle.RawID != "k3" {
		t.Fatalf("InspectAgent of a moved agent = %+v, want the live pane over the recorded handle, not live", got)
	}
	direct, err := fake.InspectAgent(ctx, handles[0])
	if err != nil {
		t.Fatal(err)
	}
	if !got.HasRevision || got.Revision != direct.Revision {
		t.Fatalf("revision = %d (%v), want %d as `agent get` reports it", got.Revision, got.HasRevision, direct.Revision)
	}
}
