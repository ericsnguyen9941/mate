package gitx_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/gitx"
)

// commitFile writes one file in dir and commits it, returning the new head.
func commitFile(t *testing.T, dir, name, text, subject string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	run(t, dir, "add", name)
	run(t, dir, "commit", "-m", subject)
	return strings.TrimSpace(run(t, dir, "rev-parse", "HEAD"))
}

// TestBranchBaseIsWhereTheBranchGrewFromEvenWhenLocalMainLags is CON-91:
// two pull requests landed on origin/main, the crew rebased onto it, and
// the local main was never pulled. The base is origin/main's tip, and the
// diff from it is the crew's one file, not the three the local main sees.
func TestBranchBaseIsWhereTheBranchGrewFromEvenWhenLocalMainLags(t *testing.T) {
	repo := newRepo(t)
	g := gitx.New()
	ctx := context.Background()
	localMain := strings.TrimSpace(run(t, repo, "rev-parse", "main"))

	// Without an origin the base is the merge base with main.
	run(t, repo, "checkout", "-q", "-b", "mate/k3")
	commitFile(t, repo, "crew.txt", "crew\n", "crew work")
	base, err := g.BranchBase(ctx, repo, "main", "mate/k3")
	if err != nil {
		t.Fatalf("BranchBase: %v", err)
	}
	if base != localMain {
		t.Fatalf("BranchBase without origin = %s, want main %s", base, localMain)
	}

	// Two merged pull requests on origin/main that local main lacks.
	run(t, repo, "checkout", "-q", "-b", "upstream", "main")
	commitFile(t, repo, "pr1.txt", "one\n", "merged PR 1")
	upstream := commitFile(t, repo, "pr2.txt", "two\n", "merged PR 2")
	run(t, repo, "update-ref", "refs/remotes/origin/main", upstream)
	run(t, repo, "checkout", "-q", "mate/k3")
	run(t, repo, "rebase", "-q", "refs/remotes/origin/main")
	run(t, repo, "branch", "-q", "-D", "upstream")

	base, err = g.BranchBase(ctx, repo, "main", "mate/k3")
	if err != nil {
		t.Fatalf("BranchBase: %v", err)
	}
	if base != upstream {
		t.Fatalf("BranchBase = %s, want origin/main %s (local main is %s)", base, upstream, localMain)
	}
	entries, err := g.DiffNumstat(ctx, repo, base, "mate/k3")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Path != "crew.txt" {
		t.Fatalf("diff from the base = %+v, want only crew.txt", entries)
	}

	// A local main ahead of origin (local delivery, nothing pushed) wins
	// the other way round: the newer of the two is still the base.
	run(t, repo, "checkout", "-q", "main")
	run(t, repo, "merge", "-q", "--ff-only", "mate/k3")
	localAhead := strings.TrimSpace(run(t, repo, "rev-parse", "main"))
	run(t, repo, "checkout", "-q", "-b", "mate/k4")
	commitFile(t, repo, "k4.txt", "k4\n", "k4 work")
	base, err = g.BranchBase(ctx, repo, "main", "mate/k4")
	if err != nil {
		t.Fatal(err)
	}
	if base != localAhead {
		t.Fatalf("BranchBase with local main ahead = %s, want %s", base, localAhead)
	}
}
