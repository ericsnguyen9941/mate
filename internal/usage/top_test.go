package usage

import (
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// seedSmallerCrews adds two more tasks to usageWorkspace's ledger, both
// smaller than k3's 97.2k: k1 at 10.5k and k2 at 50.5k.
func seedSmallerCrews(t *testing.T, w *store.Workspace) {
	t.Helper()
	handle, err := db.Open(w)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer handle.Close()
	seedActorTurnAndPricing(t, handle, timeline.CrewActorID("shop", "k1"), "shop", "crew", "unpriced-model", 10_000, 0, 0, 500, 200_000)
	seedActorTurnAndPricing(t, handle, timeline.CrewActorID("shop", "k2"), "shop", "crew", "unpriced-model", 50_000, 0, 0, 500, 200_000)
}

func TestUsageTopKeepsTheMateAndTheLargestTasksAndTotalsEverything(t *testing.T) {
	w := usageWorkspace(t)
	seedSmallerCrews(t, w)
	out := writePage(t, w, Request{Project: "shop", Top: 2})
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 8 {
		t.Fatalf("want header, mate, two tasks, total, the cut notice and the by-model lines; got %d line(s):\n%s", len(lines), out)
	}
	if !strings.HasPrefix(lines[1], "mate") || !strings.HasPrefix(lines[2], "k3") || !strings.HasPrefix(lines[3], "k2") {
		t.Fatalf("want the mate, then k3 and k2 largest first:\n%s", out)
	}
	if strings.Contains(out, "k1") {
		t.Fatalf("the smallest task must be cut:\n%s", out)
	}
	// 97.2k + 50.5k + 10.5k + the Mate's 500: the cut task still counts.
	if !strings.HasPrefix(lines[4], "total:") || !strings.Contains(lines[4], "158.7k total") {
		t.Fatalf("the total line must count every task, shown or not:\n%s", lines[4])
	}
	if !strings.Contains(lines[5], "2 largest of 3 task(s)") {
		t.Fatalf("a cut table must say what it left out:\n%s", lines[5])
	}
}

func TestUsageTopLargerThanTheLedgerChangesNothing(t *testing.T) {
	w := usageWorkspace(t)
	seedSmallerCrews(t, w)
	plain := writePage(t, w, Request{Project: "shop"})
	if top := writePage(t, w, Request{Project: "shop", Top: 10}); top != plain {
		t.Fatalf("--top above the task count must print the plain ledger:\n%s\n--- plain ---\n%s", top, plain)
	}
}
