package matehost

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// forkTurns forks the fake's session s1 on a first run and resumes the fork
// after, as claude.ForkTurnArgs does.
type forkTurns struct{}

func (forkTurns) ForkTurnArgs(_ []string, fork string, allow []string, _ string) ([]string, error) {
	args := append([]string{"--allowedTools"}, allow...)
	if fork == "" {
		return append(args, "--resume", "s1", "--fork-session", "-p"), nil
	}
	return append(args, "--resume", fork, "-p"), nil
}

func newSessionFixture(t *testing.T) *fixture {
	t.Helper()
	f := newFixture(t)
	f.host.forks = forkTurns{}
	f.host.allow = []string{"Bash(git diff*)"}
	f.host.mateBin = "/bin/mate"
	f.host.topicsDir = filepath.Join(f.dir, TopicsDir)
	f.opened["s1"] = true
	id, err := OpenNews(f.host.newsPath, News{Title: "Xác nhận xoá", Question: "Nhận cả XOA không?"}, time.Now())
	if err != nil || id != 1 {
		t.Fatalf("open news: %d %v", id, err)
	}
	f.host.RefreshNews()
	return f
}

// sessionItems waits for n finished items in news id's session.
func (f *fixture) sessionItems(t *testing.T, id, n int) []Item {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		f.host.mu.Lock()
		var items []Item
		if s, ok := f.host.sessions[id]; ok {
			items = append(items, s.items...)
		}
		f.host.mu.Unlock()
		if len(items) >= n {
			return items
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("fewer than %d session turns finished", n)
	return nil
}

// A session forks the Mate's session on its first run, told who it is and
// what the topic asks, and resumes the fork after; the fork is kept in the
// news log. It runs beside a Mate turn instead of waiting for the lock.
func TestSessionForksThenResumesBesideTheMate(t *testing.T) {
	f := newSessionFixture(t)
	f.host.Submit("slow việc của Mate")
	if err := f.host.OpenSession(1); err != nil {
		t.Fatal(err)
	}
	if err := f.host.AskSession(1, "câu 3 nói gì?"); err != nil {
		t.Fatal(err)
	}
	first := f.sessionItems(t, 1, 1)[0]
	if f.host.State().Busy == false {
		t.Fatal("the Mate's slow turn ended before the session's; the session waited for it")
	}
	prompts := f.read(t, "prompts.log")
	if !strings.Contains(prompts, "⟦topic #1⟧") || !strings.Contains(prompts, "Question: Nhận cả XOA không?") ||
		!strings.Contains(prompts, "/bin/mate news propose p 1") || first.Reply != "re: câu 3 nói gì?" {
		t.Fatalf("first session run was not told what it is (reply %q):\n%s", first.Reply, prompts)
	}
	if first.Text != "câu 3 nói gì?" || first.News != 1 {
		t.Fatalf("session item = %+v, want the captain's words only", first)
	}
	if err := f.host.AskSession(1, "vậy chọn a"); err != nil {
		t.Fatal(err)
	}
	items := f.sessionItems(t, 1, 2)
	if items[1].Reply != "re: vậy chọn a" {
		t.Fatalf("second run reply = %q, want the captain's words alone", items[1].Reply)
	}
	argv := f.read(t, "argv.log")
	if !strings.Contains(argv, "--resume s1 --fork-session") || !strings.Contains(argv, "--resume f1 -p") || !strings.Contains(argv, "--allowedTools Bash(git diff*)") {
		t.Fatalf("session runs did not fork then resume the fork:\n%s", argv)
	}
	list, _ := LoadNews(f.host.newsPath)
	if list[0].Session != "f1" {
		t.Fatalf("news #1 session = %q, want the fork f1", list[0].Session)
	}
	f.done(t, 1)
}

func TestSessionLimitBusyAndClose(t *testing.T) {
	f := newSessionFixture(t)
	for i := 0; i < MaxSessions; i++ {
		if _, err := OpenNews(f.host.newsPath, News{Title: "khác"}, time.Now()); err != nil {
			t.Fatal(err)
		}
	}
	f.host.RefreshNews()
	for id := 1; id <= MaxSessions; id++ {
		if err := f.host.OpenSession(id); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.host.OpenSession(MaxSessions + 1); !errors.Is(err, ErrTooManySessions) {
		t.Fatalf("a %dth session: %v", MaxSessions+1, err)
	}
	if err := f.host.AskSession(MaxSessions+1, "x"); !errors.Is(err, ErrNoSession) {
		t.Fatalf("asking a session never opened: %v", err)
	}
	if err := f.host.AskSession(1, "slow một chút"); err != nil {
		t.Fatal(err)
	}
	if err := f.host.AskSession(1, "nữa"); !errors.Is(err, ErrSessionBusy) {
		t.Fatalf("a second message while the session answers: %v", err)
	}
	f.sessionItems(t, 1, 1)
	f.host.CloseSession(1)
	if err := f.host.OpenSession(MaxSessions + 1); err != nil {
		t.Fatalf("closing one leaves no room: %v", err)
	}
	f.host.CloseSession(2)
	if err := f.host.OpenSession(1); err != nil {
		t.Fatal(err)
	}
	// Opened again, it has what was said and resumes the same fork.
	f.host.mu.Lock()
	s := f.host.sessions[1]
	fork, said := s.fork, len(s.items)
	f.host.mu.Unlock()
	if fork != "f1" || said != 1 {
		t.Fatalf("reopened session: fork %q, %d items; want f1 and 1", fork, said)
	}
}

// A run that answers in the Mate's own session did not fork: keeping its
// id would send the next run into the Mate's session.
func TestSessionNeverTakesTheMateSession(t *testing.T) {
	f := newSessionFixture(t)
	f.host.forks = sameSession{}
	if err := f.host.OpenSession(1); err != nil {
		t.Fatal(err)
	}
	if err := f.host.AskSession(1, "hỏi"); err != nil {
		t.Fatal(err)
	}
	f.sessionItems(t, 1, 1)
	list, _ := LoadNews(f.host.newsPath)
	f.host.mu.Lock()
	fork := f.host.sessions[1].fork
	f.host.mu.Unlock()
	if fork != "" || list[0].Session != "" {
		t.Fatalf("kept the Mate's own session as the fork: %q / %q", fork, list[0].Session)
	}
}

type sameSession struct{}

func (sameSession) ForkTurnArgs([]string, string, []string, string) ([]string, error) {
	return []string{"-p"}, nil
}

func TestSessionNeedsAMateSessionToFork(t *testing.T) {
	f := newSessionFixture(t)
	f.opened["s1"] = false
	if err := f.host.OpenSession(1); err != nil {
		t.Fatal(err)
	}
	if err := f.host.AskSession(1, "hỏi"); err == nil {
		t.Fatal("a session forked a Mate that has no session yet")
	}
}

func TestSessionClosesAfterIdling(t *testing.T) {
	f := newSessionFixture(t)
	now := time.Now()
	f.host.now = func() time.Time { return now }
	if err := f.host.OpenSession(1); err != nil {
		t.Fatal(err)
	}
	now = now.Add(SessionIdle - time.Minute)
	f.host.reapSessions()
	if f.openSessions() != 1 {
		t.Fatal("a session closed before it idled long enough")
	}
	now = now.Add(2 * time.Minute)
	f.host.reapSessions()
	if f.openSessions() != 0 {
		t.Fatal("an idle session stayed open")
	}
}

func TestForkAllowReadsAndProposesOnly(t *testing.T) {
	allow := strings.Join(ForkAllow("/opt/homebrew/bin/mate", "plantrip"), " ")
	for _, want := range []string{"Bash(git diff*)", "Bash(mate report plantrip *)", "Bash(/opt/homebrew/bin/mate news propose plantrip *)", "Bash(mate backlog plantrip)"} {
		if !strings.Contains(allow, want) {
			t.Errorf("fork may not run %q: %s", want, allow)
		}
	}
	for _, never := range []string{"crew spawn", "merge", "send", "news close", "news open", "backlog plantrip *", "rm ", "Write", "Edit"} {
		if strings.Contains(allow, never) {
			t.Errorf("fork may run %q: %s", never, allow)
		}
	}
}

func (f *fixture) openSessions() int {
	f.host.mu.Lock()
	defer f.host.mu.Unlock()
	return len(f.host.sessions)
}
