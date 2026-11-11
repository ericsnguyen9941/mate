package loopstat

import (
	"sync"
	"testing"
	"time"
)

func TestLoopCountsRuns(t *testing.T) {
	r := New()
	l := r.Loop("ingest")
	if r.Loop("ingest") != l {
		t.Fatal("a second Loop of the same name is another loop")
	}
	r.Loop("snapshot")
	start := l.Begin()
	if got := r.Read()[0]; got.RunningSince == nil || !got.RunningSince.Equal(start) || got.Runs != 0 {
		t.Fatalf("in flight = %+v, want running since %s and no run", got, start)
	}
	time.Sleep(2 * time.Millisecond)
	l.End(start)
	l.End(l.Begin())

	stats := r.Read()
	if len(stats) != 2 || stats[0].Name != "ingest" || stats[1].Name != "snapshot" {
		t.Fatalf("Read = %+v, want ingest then snapshot", stats)
	}
	got := stats[0]
	if got.Runs != 2 || got.RunningSince != nil {
		t.Fatalf("ingest = %+v, want 2 runs and none in flight", got)
	}
	if got.Max < 2*time.Millisecond || got.Last > got.Max || got.Mean != got.Total/2 || got.LastAt.Before(start) {
		t.Fatalf("ingest = %+v: max, last, mean or last_at off", got)
	}
	if stats[1].Runs != 0 || stats[1].Mean != 0 {
		t.Fatalf("snapshot = %+v, want named and never run", stats[1])
	}
}

func TestNilRecordsNothing(t *testing.T) {
	var r *Recorder
	l := r.Loop("x")
	l.End(l.Begin())
	if r.Read() != nil {
		t.Fatal("a nil Recorder read something")
	}
}

// A run records nothing on the heap: the loops call it every poll.
func TestARunAllocatesNothing(t *testing.T) {
	l := New().Loop("x")
	if n := testing.AllocsPerRun(100, func() { l.End(l.Begin()) }); n != 0 {
		t.Fatalf("a run allocated %v times, want 0", n)
	}
}

func TestConcurrentRuns(t *testing.T) {
	r := New()
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				l := r.Loop("snapshot")
				l.End(l.Begin())
				_ = r.Read()
			}
		}()
	}
	wg.Wait()
	if got := r.Read()[0]; got.Runs != 800 || got.RunningSince != nil {
		t.Fatalf("snapshot = %+v, want 800 runs and none in flight", got)
	}
}
