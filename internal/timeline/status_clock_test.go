package timeline

import (
	"testing"
	"time"
)

// A status line repeated word for word is dated by the echo that wrote it
// this time, not by the first echo of the same words: ui-evidence-harness's
// second `pr-open` (2026-10-13, offset 1705) was dated 02:29Z, the first
// echo, although it was written after 04:52Z, so the timeline told the
// Mate's wake and the dashboard the line was old news (CON-75).
func TestARepeatedStatusLineIsDatedByItsOwnEcho(t *testing.T) {
	t0 := time.Date(2026, 10, 13, 2, 29, 16, 0, time.UTC)
	t1 := t0.Add(2 * time.Hour)
	t2 := t0.Add(3 * time.Hour)
	mtime := t0.Add(4 * time.Hour)
	p := &pass{statusClock: map[string][]datedCommand{
		"crew": {
			{command: `echo "pr-open: https://x/pull/27" >> $MATE_STATUS`, at: t0},
			{command: `echo "working: upload" >> $MATE_STATUS`, at: t1},
			{command: `echo "pr-open: https://x/pull/27" >> $MATE_STATUS`, at: t2},
		},
	}}

	first, dated, _ := p.dateStatusLine("crew", "pr-open: https://x/pull/27", mtime, time.Time{})
	if !first.Equal(t0) || dated != "transcript.shell" {
		t.Fatalf("first line dated %s by %s, want %s by the first echo", first, dated, t0)
	}
	mid, _, _ := p.dateStatusLine("crew", "working: upload", mtime, first)
	if !mid.Equal(t1) {
		t.Fatalf("middle line dated %s, want %s", mid, t1)
	}
	again, dated, _ := p.dateStatusLine("crew", "pr-open: https://x/pull/27", mtime, mid)
	if !again.Equal(t2) || dated != "transcript.shell" {
		t.Fatalf("repeated line dated %s by %s, want %s by its own echo", again, dated, t2)
	}
	// With no later echo the fallback still never goes backwards.
	late, dated, _ := p.dateStatusLine("crew", "pr-open: https://x/pull/27", mtime, again)
	if !late.Equal(mtime) || dated != "status.mtime" {
		t.Fatalf("a third repeat with no echo dated %s by %s, want the file's mtime", late, dated)
	}
}
