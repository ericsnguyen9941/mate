package core

import (
	"context"
	"os"
	"runtime/metrics"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/loopstat"
	"github.com/nguyenngocanh94/mate/internal/platform"
	"github.com/nguyenngocanh94/mate/internal/query"
)

// What the core costs (CON-263): every loop records its runs in one
// loopstat.Recorder (Options.Stats), and GET coreapi.PathStats reads it
// with the process's memory and CPU, for `mate core stats`. On 2026-11-11
// a core had spent days at 1.4 GB RSS and ~70% of a CPU, its ingest pass
// taking 4.9 s of every 5 s poll, and nothing said so.

// The loops' names, as `mate core stats` prints them. The caller's loops
// record under the first five; the engine names the rest itself.
const (
	// LoopObserver is a whole observer round (watch.Deps.Pass): every
	// crew's pane, then the ingest, the budget and the no-turn checks.
	LoopObserver = "observer round"
	// LoopIngest is the timeline ingest pass inside each round.
	LoopIngest = "timeline ingest"
	// LoopOutbox is one pass of the outbox sender (outbox.Deps.Pass).
	LoopOutbox = "outbox pass"
	// LoopAutopilot is one tick of the auto daemon (autopilot.Deps.Pass),
	// its sweep included.
	LoopAutopilot = "autopilot tick"
	// LoopSweep is the sweep at the start of each tick: the pull requests
	// (prwatch) and what closed tasks left behind (gc).
	LoopSweep = "PR sweep + gc"
	// LoopReload is each read of workspace.yaml (configReloader, CON-270).
	LoopReload = "workspace.yaml reload"
	// LoopSnapshot is one computation of the snapshot every UI reads,
	// shared or fresh (Engine.Snapshot, Engine.FreshSnapshot).
	LoopSnapshot = "snapshot"
	// LoopRecovery is recovery's one run at the start.
	LoopRecovery = "recovery"
	// LoopHostWatch is one look at every headless Mate's host.
	LoopHostWatch = "host watch"
	// LoopNotify is one notification tick, past the snapshot it reads.
	LoopNotify = "notifications"
)

// timedSnapshot is compute recorded in stat.
func timedSnapshot(stat *loopstat.Loop, compute func(context.Context) (query.Snapshot, error)) func(context.Context) (query.Snapshot, error) {
	return func(ctx context.Context) (query.Snapshot, error) {
		defer stat.End(stat.Begin())
		return compute(ctx)
	}
}

// Stats is GET coreapi.PathStats: each loop's runs since the engine was
// built, and the process as runtime/metrics and the system count it now.
// It stops nothing to read them.
func (e *Engine) Stats() coreapi.CoreStats {
	s := coreapi.CoreStats{
		PID: e.holder.PID, Started: e.holder.Started, At: e.now(),
		Loops: []coreapi.LoopStats{}, Process: readProcess(),
	}
	if s.PID == 0 {
		// An engine with no core.lock (a test's) is still this process.
		s.PID = os.Getpid()
	}
	if e.files != nil {
		s.Files = e.files.stats()
	}
	for _, l := range e.stats.Read() {
		s.Loops = append(s.Loops, coreapi.LoopStats{
			Name: l.Name, Runs: l.Runs, Last: l.Last, Mean: l.Mean, Max: l.Max,
			Total: l.Total, LastAt: l.LastAt, RunningSince: l.RunningSince,
		})
	}
	return s
}

// processMetrics are the runtime/metrics readProcess samples, in the order
// it reads them.
var processMetrics = []string{
	"/memory/classes/heap/objects:bytes",
	"/memory/classes/total:bytes",
	"/gc/heap/allocs:bytes",
	"/gc/heap/allocs:objects",
	"/gc/cycles/total:gc-cycles",
	"/cpu/classes/gc/total:cpu-seconds",
	"/cpu/classes/total:cpu-seconds",
	"/sched/goroutines:goroutines",
}

// readProcess is the process now. runtime/metrics reads without stopping
// the world, which runtime.ReadMemStats does not; a metric this Go does
// not have reads as 0. The CPU figures of runtime/metrics are estimates
// the runtime refreshes at each collection.
func readProcess() coreapi.ProcessStats {
	samples := make([]metrics.Sample, len(processMetrics))
	for i, name := range processMetrics {
		samples[i].Name = name
	}
	metrics.Read(samples)
	u64 := func(i int) uint64 {
		if samples[i].Value.Kind() == metrics.KindUint64 {
			return samples[i].Value.Uint64()
		}
		return 0
	}
	f64 := func(i int) float64 {
		if samples[i].Value.Kind() == metrics.KindFloat64 {
			return samples[i].Value.Float64()
		}
		return 0
	}
	p := coreapi.ProcessStats{
		HeapBytes: u64(0), GoBytes: u64(1), AllocBytes: u64(2), AllocObjects: u64(3),
		GCCycles: u64(4), Goroutines: u64(7),
	}
	if total := f64(6); total > 0 {
		p.GCCPUFraction = f64(5) / total
	}
	if u, err := platform.ProcessUsage(); err == nil {
		p.CPU, p.PeakRSS, p.RSS = u.CPU, u.PeakRSS, u.RSS
	}
	return p
}
