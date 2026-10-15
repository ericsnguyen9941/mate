package harness

// HeadlessTurns is what a harness offers to run a Mate headless (CON-78):
// one non-interactive run per message on one session. A profile that
// implements it can be started with `mate mate start --mode headless`; one
// that does not is refused. Claude Code is the only one today.
type HeadlessTurns interface {
	// TurnArgs is one run's argv after the program: the launch args
	// (Launcher.Build) with the session flag replaced, resuming the
	// session when opened and creating it otherwise.
	TurnArgs(launch []string, sessionID string, opened bool) []string
	// LaunchSessionID is the session the launch args name.
	LaunchSessionID(launch []string) string
	// SessionOpened reports whether a run has already created the session.
	SessionOpened(cwd, sessionID string) bool
	// ParseTurnEvent reads one line of a run's output.
	ParseTurnEvent(line []byte) (TurnEvent, bool)
}

// TurnEvent is one line of a headless run's output, reduced to what a host
// shows while the run goes and keeps when it ends.
type TurnEvent struct {
	// Delta is reply text streamed since the previous event.
	Delta string
	// Said is a whole message's text, for a run that streams none.
	Said string
	// Calls are the tools the agent called in this event, in order.
	Calls []ToolCall
	// Results are the calls this event reports finished.
	Results []ToolResult
	// Context is the context, in tokens, of the request this event
	// answers; zero when the event does not say.
	Context int
	// Final marks the run's result; Text, SessionID, IsError, DurationMS
	// and Tokens are its.
	Final     bool
	Text      string
	SessionID string
	IsError   bool
	// DurationMS is the run's own measure of how long it took.
	DurationMS int64
	// Tokens are the run's tokens that were not cache reads: new input,
	// cache writes and output.
	Tokens int
}

// ToolCall is one tool the agent called.
type ToolCall struct {
	ID   string
	Name string
	// Detail is the one line that says what the call does: a command, a
	// file, a pattern. Empty when the input names none.
	Detail string
}

// ToolResult is the end of one call.
type ToolResult struct {
	ID      string
	IsError bool
}
