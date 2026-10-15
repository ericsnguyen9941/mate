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
	// Tools are the tools the agent called in this event, in order.
	Tools []string
	// Final marks the run's result; Text, SessionID and IsError are its.
	Final     bool
	Text      string
	SessionID string
	IsError   bool
}
