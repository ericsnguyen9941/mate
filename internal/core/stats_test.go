package core

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/loopstat"
)

// GET coreapi.PathStats (CON-263) is the caller's loops and the engine's
// own in one list, the caller's first as they were wired, with the
// process's memory and CPU; it is the socket's alone.
func TestClientStatsIsEveryLoopsCost(t *testing.T) {
	w := shortWorkspace(t)
	withTimelineDB(t, w)
	stats := loopstat.New()
	observer := stats.Loop(LoopObserver)
	observer.End(observer.Begin())
	e := socketEngine(t, w, Options{Stats: stats})
	if got := tcpDo(t, http.MethodGet, strings.TrimRight(e.Serving().URL, "/")+coreapi.PathStats); got != http.StatusNotFound {
		t.Fatalf("GET %s over TCP = %d, want 404", coreapi.PathStats, got)
	}

	c := dial(t, w)
	ctx := context.Background()
	if _, err := c.Snapshot(ctx); err != nil {
		t.Fatal(err)
	}
	st, err := c.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats over the socket: %v", err)
	}
	if st.PID != os.Getpid() || st.Started.IsZero() || st.At.Before(st.Started) {
		t.Fatalf("Stats = pid %d, started %s, at %s; want this process since Start", st.PID, st.Started, st.At)
	}
	byName := map[string]coreapi.LoopStats{}
	var names []string
	for _, l := range st.Loops {
		byName[l.Name] = l
		names = append(names, l.Name)
	}
	if len(names) == 0 || names[0] != LoopObserver {
		t.Fatalf("loops %v, want the caller's %q first", names, LoopObserver)
	}
	if got := byName[LoopObserver]; got.Runs != 1 || got.LastAt.IsZero() {
		t.Fatalf("%s = %+v, want the caller's one run", LoopObserver, got)
	}
	if got := byName[LoopSnapshot]; got.Runs < 1 || got.Last <= 0 || got.Mean <= 0 || got.Max < got.Last {
		t.Fatalf("%s = %+v, want the snapshot the client read", LoopSnapshot, got)
	}
	// Wired, whether or not it has ticked yet; no recovery was given.
	if _, ok := byName[LoopReload]; !ok {
		t.Fatalf("loops %v, want %q", names, LoopReload)
	}
	if _, ok := byName[LoopRecovery]; ok {
		t.Fatalf("loops %v name recovery, which this engine does not run", names)
	}
	p := st.Process
	if p.HeapBytes == 0 || p.GoBytes < p.HeapBytes || p.AllocBytes < p.HeapBytes || p.AllocObjects == 0 || p.Goroutines == 0 {
		t.Fatalf("process = %+v, want this process's memory and goroutines", p)
	}
	if p.GCCPUFraction < 0 || p.GCCPUFraction > 1 {
		t.Fatalf("GC CPU fraction %v, want a fraction", p.GCCPUFraction)
	}
	if p.CPU <= 0 || p.PeakRSS == 0 {
		t.Fatalf("process = %+v, want the CPU and peak RSS the system counted", p)
	}
}

// An engine given no Recorder records its own loops.
func TestStatsWithoutARecorder(t *testing.T) {
	w := shopWorkspace(t)
	e := New(Options{Workspace: w})
	if _, err := e.FreshSnapshot(context.Background()); err != nil {
		t.Fatal(err)
	}
	st := e.Stats()
	if st.PID != os.Getpid() {
		t.Fatalf("pid %d, want this process's", st.PID)
	}
	for _, l := range st.Loops {
		if l.Name == LoopSnapshot && l.Runs == 1 {
			return
		}
	}
	t.Fatalf("loops %+v, want the one snapshot computed", st.Loops)
}
