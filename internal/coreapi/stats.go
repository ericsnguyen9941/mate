package coreapi

import "time"

// PathStats is GET on the core's socket only: what the core has cost since
// it started (CON-263), a CoreStats. `mate core stats` reads it.
const PathStats = "/api/core/stats"

// CoreStats is the answer to GET PathStats: each of the engine's loops in
// wall time, and the process's memory, garbage collection and CPU. It is
// read when asked and costs the core nothing between two reads; the
// loops' own recording is two clock reads a run.
type CoreStats struct {
	PID int `json:"pid"`
	// Started is when the engine took core.lock; At is when these figures
	// were read.
	Started time.Time `json:"started"`
	At      time.Time `json:"at"`
	// Loops are the engine's loops in the order they were wired, each named
	// as `mate core stats` prints it. A loop with no run is wired and has
	// not run yet. Allocation is not split between them: they run side by
	// side in one heap, and no counter says which one allocated what.
	Loops   []LoopStats  `json:"loops"`
	Process ProcessStats `json:"process"`
}

// LoopStats is one loop's runs since the core started, in wall time: a run
// that waited on git, Herdr or a lock counts the wait. Durations are
// nanoseconds on the wire.
type LoopStats struct {
	Name string `json:"name"`
	Runs uint64 `json:"runs"`
	// Last is the last run's, Mean the mean of all, Max the longest and
	// Total their sum.
	Last  time.Duration `json:"last_ns"`
	Mean  time.Duration `json:"mean_ns"`
	Max   time.Duration `json:"max_ns"`
	Total time.Duration `json:"total_ns"`
	// LastAt is when the last run started; zero before the first ends.
	LastAt time.Time `json:"last_at"`
	// RunningSince is when the run in flight started, nil when none is: a
	// loop stuck in one run shows here and not in Runs.
	RunningSince *time.Time `json:"running_since,omitempty"`
}

// ProcessStats is the core process as Go's runtime/metrics and the system
// count it, read when asked (no runtime.ReadMemStats, which stops the
// world).
type ProcessStats struct {
	// HeapBytes is the heap the live and not yet swept objects take
	// (/memory/classes/heap/objects:bytes); GoBytes is all the memory the
	// Go runtime has mapped (/memory/classes/total:bytes), what the
	// resident set cannot exceed by much.
	HeapBytes uint64 `json:"heap_bytes"`
	GoBytes   uint64 `json:"go_bytes"`
	// AllocBytes and AllocObjects are everything allocated on the heap
	// since the process started (/gc/heap/allocs:*).
	AllocBytes   uint64 `json:"alloc_bytes"`
	AllocObjects uint64 `json:"alloc_objects"`
	// GCCycles is the collections since the start; GCCPUFraction the share
	// of the CPU time available to the process (GOMAXPROCS times its age)
	// that collecting took, as runtime.MemStats.GCCPUFraction means it.
	GCCycles      uint64  `json:"gc_cycles"`
	GCCPUFraction float64 `json:"gc_cpu_fraction"`
	Goroutines    uint64  `json:"goroutines"`
	// CPU is the user and system time the process has used (getrusage),
	// PeakRSS its largest resident set, RSS its resident set now where the
	// system says it cheaply (Linux); each 0 where it is not known.
	CPU     time.Duration `json:"cpu_ns"`
	PeakRSS uint64        `json:"peak_rss_bytes"`
	RSS     uint64        `json:"rss_bytes"`
}
