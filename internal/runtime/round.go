package runtime

import (
	"context"
	"sync"
)

// RoundCache is an Adapter that answers the read questions of one observer
// round from one inventory per Herdr session (CON-264).
//
// The core's watcher asked Herdr, per open crew and per 5 s round, for the
// session (`session list`, through spawn.CrewHandle), the agent (`agent
// get`) and the pane (`agent read`), and the timeline ingest asked for the
// session and the agent once more to locate a transcript: up to 5 `herdr`
// processes per crew per round with nothing changing, 51 a round for 10
// idle crews, 31 of them the watcher's (cmd/mate's
// TestWatcherRoundAsksHerdrOnce). Through a RoundCache a round asks `session
// list` once per session spec and `agent list` once per session, and every
// InspectAgent of the round - the watcher's and the ingest's - is answered
// from that list: 12 a round for the same crews, 10 of them the pane reads
// every round still makes (internal/watch pollCrew says why).
//
// An answer, an error included, stands until BeginRound: a list that ran out
// its bound is the answer for every crew of the round, which is what it was
// (each crew asking again would hold the round for one bound per crew). A
// round therefore never sees Herdr change under it, and a reader outside the
// round sees the last round's answers, at most a poll old; the one caller
// that holds a RoundCache is the watcher's wiring in cmd/mate, whose Poll
// calls BeginRound first.
//
// An agent the list does not name is asked with InspectAgent itself, once a
// round: its answer, agent_not_found included, is the agent's, exactly as
// before the list. A list is not proof of absence (Fake.ListOmit), and the
// observer opens runtime_lost on that answer.
//
// Every other method, the pane read included, goes to the wrapped Adapter
// unchanged.
type RoundCache struct {
	Adapter

	mu       sync.Mutex
	sessions map[sessionKey]sessionAnswer
	lists    map[string]listAnswer
	inspects map[string]inspectAnswer
}

var _ Adapter = (*RoundCache)(nil)

// sessionKey is what LookupSession's answer depends on. Names and Collision
// are for a session being started, not looked up.
type sessionKey struct {
	workspace  WorkspaceID
	name       string
	configHome string
}

type sessionAnswer struct {
	handle  SessionHandle
	running bool
	err     error
}

type listAnswer struct {
	agents []ObservedAgent
	err    error
}

type inspectAnswer struct {
	observed ObservedAgent
	err      error
}

// NewRoundCache wraps a for rounds. It answers nothing from memory until
// it has asked a once.
func NewRoundCache(a Adapter) *RoundCache {
	c := &RoundCache{Adapter: a}
	c.BeginRound()
	return c
}

// BeginRound forgets every answer: the next question of each kind asks
// Herdr again.
func (c *RoundCache) BeginRound() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessions = make(map[sessionKey]sessionAnswer)
	c.lists = make(map[string]listAnswer)
	c.inspects = make(map[string]inspectAnswer)
}

// LookupSession implements Adapter, once a round per spec.
func (c *RoundCache) LookupSession(ctx context.Context, spec SessionSpec) (SessionHandle, bool, error) {
	key := sessionKey{workspace: spec.WorkspaceID, name: spec.Name, configHome: spec.ConfigHome}
	c.mu.Lock()
	defer c.mu.Unlock()
	if a, ok := c.sessions[key]; ok {
		return a.handle, a.running, a.err
	}
	handle, running, err := c.Adapter.LookupSession(ctx, spec)
	c.sessions[key] = sessionAnswer{handle: handle, running: running, err: err}
	return handle, running, err
}

// ListAgents implements Adapter, once a round per session. The caller gets
// its own copy of the list.
func (c *RoundCache) ListAgents(ctx context.Context, session SessionHandle) ([]ObservedAgent, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	a := c.list(ctx, session)
	return append([]ObservedAgent(nil), a.agents...), a.err
}

// list is ListAgents without the copy. The caller holds c.mu.
func (c *RoundCache) list(ctx context.Context, session SessionHandle) listAnswer {
	if a, ok := c.lists[session.Name]; ok {
		return a
	}
	agents, err := c.Adapter.ListAgents(ctx, session)
	a := listAnswer{agents: agents, err: err}
	c.lists[session.Name] = a
	return a
}

// InspectAgent implements Adapter from the round's list of the handle's
// session: the agent of the handle's name, seen against the handle the
// caller recorded as Herdr's own `agent get` would see it. A failed list is
// the answer; an agent the list does not name is asked directly, once.
func (c *RoundCache) InspectAgent(ctx context.Context, handle AgentHandle) (ObservedAgent, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	list := c.list(ctx, handle.Session)
	if list.err != nil {
		return ObservedAgent{}, list.err
	}
	for _, listed := range list.agents {
		if listed.Handle.Name == handle.Name {
			return observedAs(handle, listed), nil
		}
	}
	key := handle.Session.Name + "/" + handle.Name
	if a, ok := c.inspects[key]; ok {
		return a.observed, a.err
	}
	observed, err := c.Adapter.InspectAgent(ctx, handle)
	c.inspects[key] = inspectAnswer{observed: observed, err: err}
	return observed, err
}

// observedAs is one listed agent as InspectAgent of the recorded handle
// reports it: the handle the caller recorded, with what Herdr says now laid
// over it, and LiveHandleOK false when the agent is no longer in the
// recorded pane - the rule internal/runtime/herdr applies to `agent get`.
func observedAs(recorded AgentHandle, listed ObservedAgent) ObservedAgent {
	now := listed.Handle
	out := recorded
	if now.Name != "" {
		out.Name = now.Name
	}
	if now.Tab.PaneID != "" {
		out.Tab.PaneID = now.Tab.PaneID
	}
	if now.Tab.TabID != "" {
		out.Tab.TabID = now.Tab.TabID
	}
	if now.Tab.TerminalID != "" {
		out.Tab.TerminalID = now.Tab.TerminalID
	}
	if now.Tab.WorkspaceID != "" {
		out.Tab.WorkspaceID = now.Tab.WorkspaceID
	}
	if out.Kind == "" {
		out.Kind = now.Kind
	}
	moved := recorded.Tab.PaneID != "" && now.Tab.PaneID != "" && recorded.Tab.PaneID != now.Tab.PaneID
	listed.Handle = out
	listed.LiveHandleOK = listed.LiveHandleOK && !moved
	return listed
}
