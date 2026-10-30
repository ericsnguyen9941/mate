package timeline

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/db"
	"github.com/nguyenngocanh94/mate/internal/gitx"
)

// Commit is one commit of a crew's branch that the default branch does not
// have. The record separator and field separator are control characters no
// commit subject can contain, so a subject with a newline in it is still one
// record.
type Commit struct {
	SHA       string
	Committed time.Time
	Subject   string
	Files     []string
}

const (
	commitRecordSep = "\x1e"
	commitFieldSep  = "\x1f"
	commitFormat    = commitRecordSep + "%H" + commitFieldSep + "%cI" + commitFieldSep + "%s" + commitFieldSep
)

// ingestGit records what a crew actually produced.
//
// `git log <default>..<branch>` - two dots, the list of commits that belong
// to this branch alone, which is the same spelling `mate diff` uses for its
// commit list and for the same reason (docs/mvp.md section 7: three dots
// there is the diff, two dots here is the list).
//
// `merge.done` is the other half: a crew is merged when it has been closed
// `finished` and its branch tip is now an ancestor of the default branch.
// There is no file that records a merge - `mate merge` types into no pane,
// so `sent.log` is silent about it (docs/mvp.md section 7) - so this is
// derived from git itself, which is also what makes it survive a rebuild.
func (p *pass) ingestGit(ctx context.Context) error {
	cfg, err := p.ing.ws.LoadProject(p.project)
	if err != nil {
		// A project whose yaml is unreadable is not a project with no
		// commits; nothing is concluded and the next pass asks again.
		return nil //nolint:nilerr
	}
	git := p.boundedGit()
	facts := p.ing.gitFacts()

	for _, crew := range p.crews {
		branch := strings.TrimSpace(crew.Meta[MetaBranch])
		if branch == "" {
			continue
		}
		// Each crew's commits are in the repo it worked in (docs/mvp.md
		// M9). A crew whose repo cannot be resolved has nothing to read.
		repoCfg, err := cfg.CrewRepo(crew.Meta)
		if err != nil {
			continue
		}
		repo, base := p.ing.ws.RepoDir(repoCfg.Path), repoCfg.DefaultBranch
		if p.gitSettled(ctx, crew) {
			continue
		}
		// A closed crew stays closed (crewstate.Declare), so what git
		// answered about its branch and its commits after it closed is not
		// asked again; see gitFacts.
		closed := crew.state().Closed()
		branchKey := repoRef{repo: repo, ref: branch}
		exists := false
		if !closed || !facts.gone[branchKey] {
			exists, err = git.BranchExists(ctx, repo, branch)
			if err != nil {
				continue
			}
			if closed && !exists {
				facts.gone[branchKey] = true
			}
		}
		if exists {
			commits, err := commitsBetween(ctx, git, repo, base, branch)
			if err == nil {
				for _, commit := range commits {
					p.commitEvent(crew, commit, branch, repo)
				}
			}
		}
		// The commits the crew's own transcript shows it making, whether or
		// not any poll ever caught the branch ahead. They are keyed by the
		// full sha, so a commit both rules find is one event.
		if err := p.ingestSightings(ctx, git, facts, crew, repo, branch, closed); err != nil {
			return err
		}
		if err := p.mergeEvent(ctx, git, crew, repo, base, branch, exists); err != nil {
			return err
		}
	}
	return nil
}

// gitSettled reports whether git has nothing left to tell about a crew: it
// closed more than settledFor ago and its merge is decided - recorded, or
// never coming because it did not close `finished`. A finished crew still
// waiting for its pull request is asked on every pass, however old.
//
// gitFacts already spares a warm pass most of this, but it lives in memory:
// measured 2026-10-30 (CON-140), once the transcripts were left unread, a
// console start still spent ~10 s waiting on one git call or more per closed
// crew of the workspace.
func (p *pass) gitSettled(ctx context.Context, crew crewRecord) bool {
	if !crew.state().Closed() {
		return false
	}
	stopped := metaTime(crew.Meta, MetaStoppedAt)
	if stopped.IsZero() || p.now.Sub(stopped) < settledFor {
		return false
	}
	if crew.state() != crewstate.StateFinished {
		return true
	}
	_, merged := p.w.eventID(ctx, dedup(KindMergeDone, crew.ActorID))
	return merged
}

// ingestSightings records the commits a crew's transcript shows it making.
// A sighting the timeline already has is not looked up again: the commit it
// names is immutable and its event is on the table. Only a sighting nobody has
// recorded costs a `git log -1`.
func (p *pass) ingestSightings(ctx context.Context, git gitx.Git, facts *gitFacts,
	crew crewRecord, repo, branch string, closed bool) error {

	sightings := p.commitSightings[crew.ActorID]
	if len(sightings) == 0 {
		return nil
	}
	recorded, err := p.recordedCommits(ctx, crew.ActorID)
	if err != nil {
		return err
	}
	for _, sighting := range sightings {
		if recorded.has(sighting.sha) {
			continue
		}
		key := repoRef{repo: repo, ref: sighting.sha}
		if closed && facts.missing[key] {
			continue
		}
		commit, ok := commitByRef(ctx, git, repo, sighting.sha)
		if !ok {
			// Remembered only when git itself says the commit is not
			// there; a git that could not run, or a pass that ran out of
			// time, has said nothing.
			if closed {
				if found, err := git.RevisionExists(ctx, repo, sighting.sha); err == nil && !found {
					facts.missing[key] = true
				}
			}
			continue
		}
		p.commitEvent(crew, commit, branch, repo)
	}
	return nil
}

// recordedShas is the full shas of the commits the timeline has for one crew.
type recordedShas []string

// has reports whether ref - a full sha, or the short one `git commit` echoed -
// names one of them.
func (r recordedShas) has(ref string) bool {
	if len(ref) < 7 {
		return false
	}
	for _, sha := range r {
		if strings.HasPrefix(sha, ref) {
			return true
		}
	}
	return false
}

// recordedCommits is every commit recorded for a crew: by an earlier pass,
// or by this one a moment ago.
func (p *pass) recordedCommits(ctx context.Context, crewActorID string) (recordedShas, error) {
	var shas recordedShas
	for _, e := range p.b.events {
		if e.Kind != KindGitCommitted || e.ActorID != crewActorID {
			continue
		}
		if sha, ok := e.Payload["sha"].(string); ok && sha != "" {
			shas = append(shas, sha)
		}
	}
	rows, err := p.tx.QueryContext(ctx,
		`SELECT json_extract(payload, '$.sha') FROM event WHERE actor_id = ? AND kind = ?`,
		crewActorID, KindGitCommitted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sha sql.NullString
		if err := rows.Scan(&sha); err != nil {
			return nil, err
		}
		if sha.Valid && sha.String != "" {
			shas = append(shas, sha.String)
		}
	}
	return shas, rows.Err()
}

// gitFacts are answers git gave about closed crews that the ingest keeps for
// the life of the Ingester instead of asking for again every five seconds.
// Measured 2026-10-14 on a live workspace of 105 closed crews: asking cost
// 2.2 s of git per pass, inside the write transaction, and recorded nothing.
//
// Only answers about git are kept, never what the timeline holds: what is
// recorded is read from the table inside the pass's transaction, so a pass
// that rolls back or a `mate reindex` that empties the table is never
// mistaken for one that wrote.
//
// Both answers are negative, and both are about a crew that can no longer
// work - a closed crew stays closed. The one thing given up: a branch
// recreated under a closed crew's name while this process runs is not read.
// Anything committed there is not the crew's work.
type gitFacts struct {
	// gone is the closed crews' branches git said are not there.
	gone map[repoRef]bool
	// missing is the closed crews' transcript commits git cannot resolve.
	missing map[repoRef]bool
}

type repoRef struct{ repo, ref string }

func (i *Ingester) gitFacts() *gitFacts {
	if i.settledGit == nil {
		i.settledGit = &gitFacts{gone: map[repoRef]bool{}, missing: map[repoRef]bool{}}
	}
	return i.settledGit
}

func (p *pass) commitEvent(crew crewRecord, commit Commit, branch, repo string) {
	p.b.event(pendingEvent{
		Dedup:   dedup(KindGitCommitted, crew.ActorID, commit.SHA),
		Project: p.project, At: commit.Committed, ActorID: crew.ActorID, Kind: KindGitCommitted,
		TaskActor: crew.ActorID,
		Payload: map[string]any{
			"sha": commit.SHA, "short": shortSHA(commit.SHA), "subject": commit.Subject,
			"branch": branch, "files": commit.Files,
		},
		RefPath: repo,
	})
}

// mergeEvent decides whether a crew's branch has landed.
//
// The evidence is git's, not a log's: the crew is closed `finished` and its
// last commit is reachable from the default branch. A branch `mate merge`
// deleted on its way out is handled by the same rule from the other side -
// the branch is gone and the commits it carried are in `default`, so the last
// commit this ingest recorded for the crew is the one to test.
//
// Which is why a crew's commits are read from its transcript as well as from
// `git log`: measured 2026-09-20, a crew committed six seconds before the
// captain merged, no five-second poll fell inside that window, and a merge
// that could only be proved from a branch nobody had read was not proved at
// all. The transcript keeps the sha for ever.
func (p *pass) mergeEvent(ctx context.Context, git gitx.Git, crew crewRecord,
	repo, base, branch string, branchExists bool) error {

	if crew.state() != crewstate.StateFinished {
		return nil
	}
	// A merge already on the table is decided: the event is keyed by the
	// crew alone and the task row keeps its pointer, so asking git again
	// could only produce the same event.
	if _, recorded := p.w.eventID(ctx, dedup(KindMergeDone, crew.ActorID)); recorded {
		return nil
	}
	tip := ""
	if branchExists {
		sha, err := git.HeadCommit(ctx, repo, branch)
		if err != nil {
			return nil //nolint:nilerr
		}
		tip = sha
	} else {
		sha, err := p.lastRecordedCommit(ctx, crew.ActorID)
		if err != nil {
			return err
		}
		tip = sha
	}
	if tip == "" {
		return nil
	}
	landed, err := git.IsAncestor(ctx, repo, tip, base)
	if err != nil || !landed {
		return nil //nolint:nilerr
	}
	at := metaTime(crew.Meta, MetaStoppedAt)
	if at.IsZero() {
		at = p.now
	}
	// The merge closes the crew, so the person who ran it is the person who
	// closed it: `mate merge` runs `crew stop` itself on success
	// (docs/mvp.md M4). linkCauses names which turn or gesture it was.
	mergeDedup := dedup(KindMergeDone, crew.ActorID)
	p.b.event(pendingEvent{
		Dedup:   mergeDedup,
		Project: p.project, At: at, ActorID: p.mate.ActorID, Subject: crew.ActorID,
		Kind: KindMergeDone, TaskActor: crew.ActorID,
		Payload: map[string]any{
			"crew": crew.ID, "branch": branch, "into": base,
			"sha": tip, "short": shortSHA(tip),
		},
		RefPath: repo,
	})
	for i := range p.b.tasks {
		if p.b.tasks[i].CrewActorID == crew.ActorID {
			p.b.tasks[i].MergedDedup = mergeDedup
		}
	}
	return nil
}

// lastRecordedCommit is the newest commit this timeline has for a crew,
// looking at what this pass has just found as well as at the table: a crew
// that commits and is merged inside one poll would otherwise have its commit
// and its merge decided in the same pass, with the merge unable to see it.
func (p *pass) lastRecordedCommit(ctx context.Context, crewActorID string) (string, error) {
	newest := time.Time{}
	sha := ""
	for _, e := range p.b.events {
		if e.Kind != KindGitCommitted || e.ActorID != crewActorID {
			continue
		}
		if value, ok := e.Payload["sha"].(string); ok && !e.At.Before(newest) {
			newest, sha = e.At, value
		}
	}
	var stored sql.NullString
	var storedAt sql.NullString
	err := p.tx.QueryRowContext(ctx,
		`SELECT json_extract(payload, '$.sha'), at FROM event
		  WHERE actor_id = ? AND kind = ? ORDER BY at DESC, id DESC LIMIT 1`,
		crewActorID, KindGitCommitted).Scan(&stored, &storedAt)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if stored.Valid && !db.ParseTime(storedAt.String).Before(newest) {
		return stored.String, nil
	}
	return sha, nil
}

// commitsBetween runs one `git log` and parses it. gitx.LogOneline is the
// review list and carries only a short sha and a subject; the timeline needs
// the commit's own time and the files it touched, so this asks for them
// directly rather than pretending a one-line summary has them.
func commitsBetween(ctx context.Context, git gitx.Git, repo, base, branch string) ([]Commit, error) {
	out, err := git.Log(ctx, repo, base+".."+branch,
		"--no-color", "--reverse", "--name-only", "--format="+commitFormat)
	if err != nil {
		return nil, err
	}
	return parseCommits(out), nil
}

func parseCommits(out string) []Commit {
	var commits []Commit
	for _, record := range strings.Split(out, commitRecordSep) {
		if strings.TrimSpace(record) == "" {
			continue
		}
		fields := strings.SplitN(record, commitFieldSep, 4)
		if len(fields) < 4 {
			continue
		}
		at, err := time.Parse(time.RFC3339, strings.TrimSpace(fields[1]))
		if err != nil {
			continue
		}
		var files []string
		for _, line := range strings.Split(fields[3], "\n") {
			if name := strings.TrimSpace(line); name != "" {
				files = append(files, name)
			}
		}
		commits = append(commits, Commit{
			SHA:       strings.TrimSpace(fields[0]),
			Committed: at,
			Subject:   strings.TrimSpace(fields[2]),
			Files:     files,
		})
	}
	return commits
}

// commitByRef reads one commit by whatever names it - a short sha echoed by
// `git commit` is enough. A ref git cannot resolve is not an error here: a
// crew may have committed in a worktree whose objects are no longer reachable,
// and a commit nobody can look up is one this timeline cannot describe.
func commitByRef(ctx context.Context, git gitx.Git, repo, ref string) (Commit, bool) {
	out, err := git.Log(ctx, repo, ref, "-1", "--no-color", "--name-only", "--format="+commitFormat)
	if err != nil {
		return Commit{}, false
	}
	commits := parseCommits(out)
	if len(commits) != 1 {
		return Commit{}, false
	}
	return commits[0], true
}

func shortSHA(sha string) string {
	if len(sha) <= 7 {
		return sha
	}
	return sha[:7]
}

// boundedGit is the ingest's git with every call held to the pass's external
// deadline (DefaultProjectTimeout): past it a call fails at once, which the
// callers above read as git having said nothing.
func (p *pass) boundedGit() gitx.Git {
	git := p.ing.deps.git()
	if p.external.IsZero() {
		return git
	}
	return gitx.Git{Runner: deadlineRunner{inner: git.Runner, deadline: p.external}}
}

type deadlineRunner struct {
	inner    gitx.Runner
	deadline time.Time
}

func (r deadlineRunner) Run(ctx context.Context, cmd gitx.Command) (gitx.Result, error) {
	ctx, cancel := context.WithDeadline(ctx, r.deadline)
	defer cancel()
	return r.inner.Run(ctx, cmd)
}
