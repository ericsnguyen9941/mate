package matehost

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// HerdrSource is the --source the host reports its lifecycle under.
const HerdrSource = "mate"

// HerdrReporter makes the host a named Herdr agent from inside its own pane
// and keeps its state current. Measured 2026-10-15 (docs/evidence/
// headless-mate-probe-2026-10-15.md): `pane report-agent` from the process
// the pane runs lists the pane as an agent for as long as that process
// lives, `agent rename` gives it the name `agent get` and `agent attach`
// find, and the agent is gone once the process exits. So liveness, attach
// and the console's rows treat a headless Mate like any other.
type HerdrReporter struct {
	Binary string
	Pane   string
	Kind   string
	Agent  string

	mu      sync.Mutex
	pending chan string
	once    sync.Once
}

// NewHerdrReporter reads the pane from the environment Herdr gives every
// pane. ok is false outside Herdr: a host run by hand in a plain terminal
// reports nothing.
func NewHerdrReporter(kind, agent string) (*HerdrReporter, bool) {
	pane := os.Getenv("HERDR_PANE_ID")
	if pane == "" {
		return nil, false
	}
	bin := os.Getenv("HERDR_BIN_PATH")
	if bin == "" {
		bin = "herdr"
	}
	return &HerdrReporter{Binary: bin, Pane: pane, Kind: kind, Agent: agent}, true
}

func (r *HerdrReporter) run(args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, r.Binary, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("herdr %s: %v: %s", strings.Join(args[:2], " "), err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Register reports the host idle and names it.
func (r *HerdrReporter) Register() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.run("pane", "report-agent", "--source", HerdrSource, "--agent", r.Kind, "--state", "idle", r.Pane); err != nil {
		return err
	}
	if r.Agent == "" {
		return nil
	}
	return r.run("agent", "rename", r.Pane, r.Agent)
}

// Report implements Reporter without waiting for Herdr: a turn starts
// whatever Herdr takes to answer under load. Reports go out in order, and
// a state superseded before it went out is skipped. A failure is not the
// turn's: it is dropped.
func (r *HerdrReporter) Report(state string) {
	r.once.Do(func() {
		r.pending = make(chan string, 1)
		go func() {
			for s := range r.pending {
				r.mu.Lock()
				_ = r.run("pane", "report-agent", "--source", HerdrSource, "--agent", r.Kind, "--state", s, r.Pane)
				r.mu.Unlock()
			}
		}()
	})
	for {
		select {
		case r.pending <- state:
			return
		default:
		}
		select {
		case <-r.pending:
		default:
		}
	}
}

// Release gives the pane's lifecycle back to Herdr as the host exits.
func (r *HerdrReporter) Release() {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = r.run("pane", "release-agent", "--source", HerdrSource, "--agent", r.Kind, r.Pane)
}
