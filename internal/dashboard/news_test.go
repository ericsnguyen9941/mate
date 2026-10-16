package dashboard

import (
	"crypto/sha256"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/matehost"
)

// newsLog is a Mate's news.jsonl as the REPL and `mate news` write it: one
// news the captain has not opened, one opened and waiting, one answered and
// on the Mate, and one closed with a reason.
const newsLog = `{"at":"2026-09-19T10:40:00Z","op":"open","id":1,"crew":"buybtn","kind":"decision","title":"Checkout URL","context":["buybtn asked twice"],"question":"Which page?","options":[{"key":"a","text":"express","recommended":true},{"key":"b","text":"classic"}]}
{"at":"2026-09-19T10:41:00Z","op":"open","id":2,"title":"PR #12 ready"}
{"at":"2026-09-19T10:41:30Z","op":"seen","id":2}
{"at":"2026-09-19T10:42:00Z","op":"open","id":3,"title":"Key for staging"}
{"at":"2026-09-19T10:42:10Z","op":"seen","id":3}
{"at":"2026-09-19T10:43:00Z","op":"answer","id":3,"text":"[a] use the vault one"}
{"at":"2026-09-19T10:30:00Z","op":"open","id":4,"title":"Old flaky test"}
{"at":"2026-09-19T10:31:00Z","op":"seen","id":4}
{"at":"2026-09-19T10:32:00Z","op":"answer","id":4,"text":"skip it"}
{"at":"2026-09-19T10:35:00Z","op":"close","id":4,"text":"PR merged"}
{"at":"2026-09-19T10:36:00Z","op":"archive","id":4}
not json
`

func writeNewsLog(t *testing.T, f *fixture, project, body string) string {
	t.Helper()
	path := matehost.NewsPath(f.ws.MateDir(project))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestNewsListsOpenAndClosedNewsWithTheirState(t *testing.T) {
	f := newFixture(t)
	path := writeNewsLog(t, f, fixtureProject, newsLog)
	before := sha256.Sum256([]byte(newsLog))

	var got NewsResponse
	f.get(t, "/api/news?project="+fixtureProject, http.StatusOK, &got)
	if got.Project != fixtureProject || got.GeneratedAt == "" || len(got.Projects) != 1 {
		t.Fatalf("response = %+v", got)
	}
	p := got.Projects[0]
	if p.Project != fixtureProject || !p.Logged || p.Error != "" {
		t.Fatalf("project = %+v", p)
	}
	var order, states []string
	for _, n := range p.News {
		order = append(order, n.Title)
		states = append(states, n.Status)
	}
	// Newest move first; the closed news stays, it is the history.
	if strings.Join(order, "|") != "Key for staging|PR #12 ready|Checkout URL|Old flaky test" {
		t.Errorf("order = %v", order)
	}
	if strings.Join(states, "|") != "working|waiting|new|done" {
		t.Errorf("states = %v", states)
	}
	byID := map[int]NewsEntry{}
	for _, n := range p.News {
		byID[n.ID] = n
	}
	if n := byID[1]; n.Crew != "buybtn" || n.Question != "Which page?" || len(n.Options) != 2 || !n.Options[0].Recommended ||
		len(n.Context) != 1 || n.SeenAt != "" || n.Answer != "" {
		t.Errorf("news 1 = %+v", n)
	}
	if n := byID[2]; n.SeenAt != "2026-09-19T10:41:30.000000000Z" || n.At != "2026-09-19T10:41:00.000000000Z" {
		t.Errorf("news 2 seen %q, at %q: being seen is not a move", n.SeenAt, n.At)
	}
	if n := byID[3]; n.Answer != "[a] use the vault one" || n.AnsweredAt != "2026-09-19T10:43:00.000000000Z" {
		t.Errorf("news 3 = %+v", n)
	}
	if n := byID[4]; n.Status != "done" || !n.Archived || len(n.Notes) != 1 || n.Notes[0].Text != "PR merged" || n.Answer != "skip it" {
		t.Errorf("news 4 = %+v", n)
	}

	// Read on every request, and never written: the Office does not mark a
	// news seen.
	data, err := os.ReadFile(path)
	if err != nil || sha256.Sum256(data) != before {
		t.Fatalf("the news log changed under a GET: %v", err)
	}
	writeNewsLog(t, f, fixtureProject, newsLog+`{"at":"2026-09-19T10:44:00Z","op":"seen","id":1}`+"\n")
	f.get(t, "/api/news?project="+fixtureProject, http.StatusOK, &got)
	for _, n := range got.Projects[0].News {
		if n.ID == 1 && n.Status != "waiting" {
			t.Errorf("news 1 after a seen = %q, want waiting without a new event", n.Status)
		}
	}
}

func TestNewsWithoutALogIsNotAnError(t *testing.T) {
	f := newFixture(t)
	var got NewsResponse
	f.get(t, "/api/news", http.StatusOK, &got)
	if got.Project != "" || len(got.Projects) != 1 {
		t.Fatalf("response = %+v", got)
	}
	if p := got.Projects[0]; p.Logged || p.Error != "" || p.News == nil || len(p.News) != 0 {
		t.Errorf("a Mate with no news log = %+v", p)
	}
	f.get(t, "/api/news?project=nope", http.StatusNotFound, nil)
}

func TestNewsSaysWhyALogCannotBeRead(t *testing.T) {
	f := newFixture(t)
	path := matehost.NewsPath(f.ws.MateDir(fixtureProject))
	// A directory where the log should be: it exists, and does not read.
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	var got NewsResponse
	f.get(t, "/api/news", http.StatusOK, &got)
	if p := got.Projects[0]; !p.Logged || p.Error == "" || len(p.News) != 0 {
		t.Errorf("an unreadable log = %+v", p)
	}
}
