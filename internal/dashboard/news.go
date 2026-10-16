package dashboard

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/matehost"
)

// NewsResponse is GET /api/news[?project=<p>]: every news a headless Mate
// opened for the captain (`mate/news.jsonl`, CON-78), closed ones too, for
// the Office's captain drawer (CON-97). It reads the log through
// matehost.LoadNews on every request and is not cached: seeing or
// answering a news in the REPL is not an event.
type NewsResponse struct {
	envelope
	// Project is the one asked for, empty for the whole workspace.
	Project  string        `json:"project,omitempty"`
	Projects []ProjectNews `json:"projects"`
}

// ProjectNews is one project's news, newest move first.
type ProjectNews struct {
	Project string `json:"project"`
	// Logged is whether the project's Mate has a news log at all: a Mate
	// that never opened one (or is not headless) has none, which is not
	// the same answer as a log with nothing in it.
	Logged bool        `json:"logged"`
	News   []NewsEntry `json:"news"`
	// Error is why the log could not be read; News is empty then.
	Error string `json:"error,omitempty"`
}

// NewsEntry is one news as its log folds now. Status is matehost's word:
// `new` (the captain has not opened it), `waiting` (opened, not answered),
// `working` (answered, the Mate is on it), `done`, `dropped`.
type NewsEntry struct {
	ID         int          `json:"id"`
	Crew       string       `json:"crew,omitempty"`
	Kind       string       `json:"kind,omitempty"`
	Title      string       `json:"title"`
	Context    []string     `json:"context"`
	Question   string       `json:"question,omitempty"`
	Options    []NewsOption `json:"options"`
	Status     string       `json:"status"`
	OpenedAt   string       `json:"opened_at"`
	At         string       `json:"at"`
	SeenAt     string       `json:"seen_at,omitempty"`
	Answer     string       `json:"answer,omitempty"`
	AnsweredAt string       `json:"answered_at,omitempty"`
	Notes      []NewsNote   `json:"notes"`
	Archived   bool         `json:"archived,omitempty"`
}

// NewsOption is one answer the Mate offered.
type NewsOption struct {
	Key         string `json:"key"`
	Text        string `json:"text"`
	Recommended bool   `json:"recommended,omitempty"`
}

// NewsNote is an update the Mate added, or why it closed the news.
type NewsNote struct {
	At   string `json:"at"`
	Text string `json:"text"`
}

func (s *Server) handleNews(_ context.Context, r *http.Request, env envelope) (any, error) {
	project := strings.TrimSpace(r.URL.Query().Get("project"))
	var names []string
	if project != "" {
		if _, ok := s.ws.Project(project); !ok {
			return nil, errNotFound{reason: unknownProject(s.ws, project)}
		}
		names = []string{project}
	} else {
		for _, ref := range s.ws.Projects() {
			names = append(names, ref.Name)
		}
	}
	out := NewsResponse{envelope: env, Project: project, Projects: []ProjectNews{}}
	for _, name := range names {
		out.Projects = append(out.Projects, s.projectNews(name))
	}
	return out, nil
}

// projectNews folds one project's news log. It reads the file and writes
// nothing: the Office shows news, it never marks one seen.
func (s *Server) projectNews(project string) ProjectNews {
	out := ProjectNews{Project: project, News: []NewsEntry{}}
	path := matehost.NewsPath(s.ws.MateDir(project))
	if _, err := os.Stat(path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			out.Error = err.Error()
		}
		return out
	}
	out.Logged = true
	list, err := matehost.LoadNews(path)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].At.After(list[j].At) })
	for _, n := range list {
		e := NewsEntry{
			ID: n.ID, Crew: n.Crew, Kind: n.Kind, Title: n.Title, Context: append([]string{}, n.Context...),
			Question: n.Question, Options: []NewsOption{}, Status: n.Status,
			OpenedAt: db.FormatTime(n.OpenedAt), At: db.FormatTime(n.At), SeenAt: db.FormatTime(n.SeenAt),
			Answer: n.Answer, AnsweredAt: db.FormatTime(n.AnsweredAt), Notes: []NewsNote{}, Archived: n.Archived,
		}
		for _, o := range n.Options {
			e.Options = append(e.Options, NewsOption{Key: o.Key, Text: o.Text, Recommended: o.Recommended})
		}
		for _, x := range n.Notes {
			e.Notes = append(e.Notes, NewsNote{At: db.FormatTime(x.At), Text: x.Text})
		}
		out.News = append(out.News, e)
	}
	return out
}
