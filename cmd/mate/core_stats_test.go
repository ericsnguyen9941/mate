package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/core"
	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// fakeCoreStats is what the fake core answers GET coreapi.PathStats with:
// a core up two hours whose ingest takes most of every poll, the snapshot
// in flight, and notifications wired but never run.
func fakeCoreStats(pid int, started time.Time) coreapi.CoreStats {
	at := started.Add(2 * time.Hour)
	running := at.Add(-1500 * time.Millisecond)
	return coreapi.CoreStats{
		PID: pid, Started: started, At: at,
		Loops: []coreapi.LoopStats{
			{Name: core.LoopObserver, Runs: 1440, Last: 4900 * time.Millisecond, Mean: 4800 * time.Millisecond, Max: 9100 * time.Millisecond, Total: 6912 * time.Second, LastAt: at.Add(-3 * time.Second)},
			{Name: core.LoopIngest, Runs: 1440, Last: 4700 * time.Millisecond, Mean: 4600 * time.Millisecond, Max: 8800 * time.Millisecond, Total: 6624 * time.Second, LastAt: at.Add(-3 * time.Second)},
			{Name: core.LoopSnapshot, Runs: 3600, Last: 98 * time.Millisecond, Mean: 101 * time.Millisecond, Max: 2 * time.Second, Total: 363600 * time.Millisecond, LastAt: at.Add(-2 * time.Second), RunningSince: &running},
			{Name: core.LoopNotify},
		},
		Process: coreapi.ProcessStats{
			HeapBytes: 1400 << 20, GoBytes: 1600 << 20, AllocBytes: 7 << 40, AllocObjects: 81_000_000_000,
			GCCycles: 9411, GCCPUFraction: 0.031, Goroutines: 84,
			CPU: 84 * time.Minute, PeakRSS: 1500 << 20,
		},
		Files: &coreapi.FileWatchStats{Watched: 7, Missing: 1, Events: 312, Changes: 40},
	}
}

func TestPrintCoreStats(t *testing.T) {
	started := time.Date(2026, 11, 11, 7, 0, 0, 0, time.UTC)
	var b bytes.Buffer
	printCoreStats(&b, "~/workspace", fakeCoreStats(4121, started))
	want := `mate core · pid 4121 · ~/workspace · up 2h00m
loop             runs  last   mean   max    busy  last run
observer round   1440  4.90s  4.80s  9.10s  96%   3.00s ago
timeline ingest  1440  4.70s  4.60s  8.80s  92%   3.00s ago
snapshot         3600  98ms   101ms  2.00s  5%    running for 1.50s
notifications    0     -      -      -      -     never

heap        1.4 GB live · 1.6 GB mapped by Go · peak RSS 1.5 GB
allocated   7.0 TB in 81.0G objects since start · 1019 MB/s
gc          9411 cycles · 3.1% of the CPU
cpu         1h24m · 70% of one core since start
goroutines  84
files       7 watched · 1 missing · 312 events → 40 changes
`
	if got := b.String(); got != want {
		t.Fatalf("printCoreStats =\n%s\nwant\n%s", got, want)
	}
}

// A core whose system gave no watch says so, and that the polls run on.
func TestStatsFilesOff(t *testing.T) {
	got := statsFiles(coreapi.FileWatchStats{Off: "too many open files"})
	if want := "off: too many open files (the polls find every change)"; got != want {
		t.Fatalf("statsFiles = %q, want %q", got, want)
	}
	got = statsFiles(coreapi.FileWatchStats{Watched: 3, Events: 1200, Changes: 9, Errors: 2, LastError: "no space left on device"})
	if want := "3 watched · 1.2k events → 9 changes · 2 errors, last: no space left on device"; got != want {
		t.Fatalf("statsFiles = %q, want %q", got, want)
	}
}

func TestStatsUnits(t *testing.T) {
	for d, want := range map[time.Duration]string{
		400 * time.Microsecond:  "0.4ms",
		712 * time.Millisecond:  "712ms",
		4120 * time.Millisecond: "4.12s",
		185 * time.Second:       "3m05s",
		134 * time.Minute:       "2h14m",
	} {
		if got := statsDuration(d); got != want {
			t.Errorf("statsDuration(%s) = %q, want %q", d, got, want)
		}
	}
	for n, want := range map[uint64]string{512: "512 B", 2048: "2.0 KB", 62 << 20: "62.0 MB", 140 << 20: "140 MB"} {
		if got := statsBytes(n); got != want {
			t.Errorf("statsBytes(%d) = %q, want %q", n, got, want)
		}
	}
	for n, want := range map[uint64]string{950: "950", 12_300: "12.3k", 4_500_000: "4.5M", 81_000_000_000: "81.0G"} {
		if got := statsCount(n); got != want {
			t.Errorf("statsCount(%d) = %q, want %q", n, got, want)
		}
	}
}

// `mate core stats` asks the running core on its socket; with no core it
// is errCoreNotRunning, exit 1 as for status.
func TestCoreStatsAsksTheRunningCore(t *testing.T) {
	t.Parallel()
	w := coreWorkspace(t)
	f := newFakeCores(t, "serve")
	var b bytes.Buffer
	if err := f.stats(w, &b); !errors.Is(err, errCoreNotRunning) {
		t.Fatalf("stats with no core = %v, want %v", err, errCoreNotRunning)
	}
	pid, _ := f.startOK(t, w)
	b.Reset()
	if err := f.stats(w, &b); err != nil {
		t.Fatalf("stats: %v\n%s", err, b.String())
	}
	out := b.String()
	for _, want := range []string{"mate core · pid " + itoa(pid), "timeline ingest  1440  4.70s  4.60s", "goroutines  84"} {
		if !strings.Contains(out, want) {
			t.Fatalf("stats said\n%s\nwant %q in it", out, want)
		}
	}
	if err := f.stop(w, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
}

// An engine inside a console serves no socket: stats says so.
func TestCoreStatsBesideAConsoleEngine(t *testing.T) {
	w := coreWorkspace(t)
	lock, err := w.LockCore(store.CoreHolder{PID: 4242, Owner: core.OwnerConsole, Started: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Unlock()
	err = liveCoreControl().stats(w, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "inside mate console (pid 4242)") {
		t.Fatalf("stats beside a console's engine = %v", err)
	}
}
