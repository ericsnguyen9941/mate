package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nguyenngocanh94/mate/internal/core"
	"github.com/nguyenngocanh94/mate/internal/coreapi"
	"github.com/nguyenngocanh94/mate/internal/coreclient"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// `mate core stats` (CON-263): what the running core has cost since it
// started, loop by loop, and what the process holds now. On 2026-11-11 a
// core had spent days at 1.4 GB RSS and ~70% of a CPU, its timeline ingest
// taking 4.9 s of every 5 s poll, and nothing showed it; this does, in one
// command, without a profile.

// coreStatsWait bounds the question: the core answers it from memory.
const coreStatsWait = 5 * time.Second

// askCoreStats is GET coreapi.PathStats on the core's socket.
func askCoreStats(ctx context.Context, socket string) (coreapi.CoreStats, error) {
	client, err := coreclient.Dial(socket)
	if err != nil {
		return coreapi.CoreStats{}, err
	}
	defer client.Close()
	return client.Stats(ctx)
}

// stats is `mate core stats`: errCoreNotRunning when no engine runs on ws;
// an engine inside a console serves no socket to ask.
func (c coreControl) stats(ws *store.Workspace, stdout io.Writer) error {
	h, held, err := ws.CoreRunning()
	if err != nil {
		return fmt.Errorf("mate core: reading %s: %w", ws.CoreLockFile(), err)
	}
	if !held {
		return errCoreNotRunning
	}
	if h.Owner == core.OwnerConsole {
		return fmt.Errorf("mate core: the engine runs inside mate console (%s), which serves no socket to ask", pidWord(h.PID))
	}
	ctx, cancel := context.WithTimeout(context.Background(), coreStatsWait)
	defer cancel()
	st, err := c.askStats(ctx, ws.CoreSocket())
	if err != nil {
		return fmt.Errorf("mate core: %s: %w", ws.CoreSocket(), err)
	}
	printCoreStats(stdout, homeShort(ws.Root()), st)
	return nil
}

// printCoreStats is the page; a core like the one of 2026-11-11 would have
// said:
//
//	mate core · pid 4121 · ~/workspace · up 2h00m
//	loop             runs  last   mean   max    busy  last run
//	observer round   1440  4.90s  4.80s  9.10s  96%   3.00s ago
//	timeline ingest  1440  4.70s  4.60s  8.80s  92%   3.00s ago
//	snapshot         3600  98ms   101ms  2.00s  5%    running for 1.50s
//	notifications    0     -      -      -      -     never
//
//	heap        1.4 GB live · 1.6 GB mapped by Go · peak RSS 1.5 GB
//	allocated   7.0 TB in 81.0G objects since start · 1019 MB/s
//	gc          9411 cycles · 3.1% of the CPU
//	cpu         1h24m · 70% of one core since start
//	goroutines  84
//
// busy is the share of the core's life the loop spent running: a loop
// near 100% never rests between two runs. A loop with no run is wired and
// has not run yet. Allocation is not split by loop: they share one heap.
func printCoreStats(w io.Writer, where string, st coreapi.CoreStats) {
	up := st.At.Sub(st.Started)
	head := fmt.Sprintf("mate core · %s · %s", pidWord(st.PID), where)
	if !st.Started.IsZero() && up > 0 {
		head += " · up " + statsDuration(up)
	}
	fmt.Fprintln(w, head)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "loop\truns\tlast\tmean\tmax\tbusy\tlast run")
	for _, l := range st.Loops {
		busy := "-"
		if up > 0 {
			busy = fmt.Sprintf("%.0f%%", 100*float64(l.Total)/float64(up))
		}
		when := "never"
		switch {
		case l.RunningSince != nil:
			when = "running for " + statsDuration(st.At.Sub(*l.RunningSince))
		case !l.LastAt.IsZero():
			when = statsDuration(st.At.Sub(l.LastAt)) + " ago"
		}
		if l.Runs == 0 {
			fmt.Fprintf(tw, "%s\t0\t-\t-\t-\t-\t%s\n", l.Name, when)
			continue
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n", l.Name, l.Runs,
			statsDuration(l.Last), statsDuration(l.Mean), statsDuration(l.Max), busy, when)
	}
	_ = tw.Flush()

	p := st.Process
	fmt.Fprintln(w)
	tw = tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	memory := []string{statsBytes(p.HeapBytes) + " live", statsBytes(p.GoBytes) + " mapped by Go"}
	if p.RSS > 0 {
		memory = append(memory, "RSS "+statsBytes(p.RSS))
	}
	if p.PeakRSS > 0 {
		memory = append(memory, "peak RSS "+statsBytes(p.PeakRSS))
	}
	fmt.Fprintf(tw, "heap\t%s\n", strings.Join(memory, " · "))
	allocated := fmt.Sprintf("%s in %s objects since start", statsBytes(p.AllocBytes), statsCount(p.AllocObjects))
	if up > 0 {
		allocated += fmt.Sprintf(" · %s/s", statsBytes(uint64(float64(p.AllocBytes)/up.Seconds())))
	}
	fmt.Fprintf(tw, "allocated\t%s\n", allocated)
	fmt.Fprintf(tw, "gc\t%d cycles · %.1f%% of the CPU\n", p.GCCycles, 100*p.GCCPUFraction)
	if p.CPU > 0 {
		cpu := statsDuration(p.CPU)
		if up > 0 {
			cpu += fmt.Sprintf(" · %.0f%% of one core since start", 100*p.CPU.Seconds()/up.Seconds())
		}
		fmt.Fprintf(tw, "cpu\t%s\n", cpu)
	}
	fmt.Fprintf(tw, "goroutines\t%d\n", p.Goroutines)
	if f := st.Files; f != nil {
		fmt.Fprintf(tw, "files\t%s\n", statsFiles(*f))
	}
	_ = tw.Flush()
}

// statsFiles is the core's watch on the files written past the change feed
// (CON-266): "7 watched · 1 missing · 312 events → 40 changes", or "off: …"
// when the system gave none.
func statsFiles(f coreapi.FileWatchStats) string {
	if f.Off != "" {
		return "off: " + f.Off + " (the polls find every change)"
	}
	line := fmt.Sprintf("%d watched", f.Watched)
	if f.Missing > 0 {
		line += fmt.Sprintf(" · %d missing", f.Missing)
	}
	line += fmt.Sprintf(" · %s events → %s changes", statsCount(f.Events), statsCount(f.Changes))
	if f.Errors > 0 {
		line += fmt.Sprintf(" · %d errors, last: %s", f.Errors, f.LastError)
	}
	return line
}

// statsDuration is d in the unit that keeps it short: 0.4ms, 712ms, 4.12s,
// 3m05s, 2h14m.
func statsDuration(d time.Duration) string {
	switch {
	case d < 10*time.Millisecond:
		return fmt.Sprintf("%.1fms", float64(d)/float64(time.Millisecond))
	case d < time.Second:
		return fmt.Sprintf("%dms", d.Milliseconds())
	case d < time.Minute:
		return fmt.Sprintf("%.2fs", d.Seconds())
	case d < time.Hour:
		return fmt.Sprintf("%dm%02ds", int(d.Minutes()), int(d.Seconds())%60)
	default:
		return fmt.Sprintf("%dh%02dm", int(d.Hours()), int(d.Minutes())%60)
	}
}

// statsCount is n in thousands, millions or billions: 950, 12.3k, 81.0G.
func statsCount(n uint64) string {
	switch {
	case n < 1e3:
		return fmt.Sprintf("%d", n)
	case n < 1e6:
		return fmt.Sprintf("%.1fk", float64(n)/1e3)
	case n < 1e9:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	default:
		return fmt.Sprintf("%.1fG", float64(n)/1e9)
	}
}

// statsBytes is n in binary units: 512 B, 62.1 MB, 1.21 GB.
func statsBytes(n uint64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	value, prefix := float64(n)/unit, 0
	for value >= unit && prefix < 4 {
		value /= unit
		prefix++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	if value >= 100 {
		return fmt.Sprintf("%.0f %s", value, units[prefix])
	}
	return fmt.Sprintf("%.1f %s", value, units[prefix])
}
