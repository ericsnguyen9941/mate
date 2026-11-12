package timeline

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// What a poll that is not a sweep may skip within a project it reads
// (CON-268), each under a key of file stats like seal.go's: size and mtime,
// stat'ed before the read, so a write racing the read leaves an old key on
// new content and the next poll reads again. Like the seals, the memos
// live in the Ingester's memory, are kept only once the pass that read
// them commits, and a Reindex forgets them. A sweep (Scope.Sweep, once a
// minute from the observer) trusts none of them and reads everything the
// seals leave, so whatever a key misses is read again within the minute.
//
// Measured 2026-11-12 on a copy of the captain's workspace, a steady pass
// over every project cost ~150 ms, and the CPU profile put a quarter of it
// in a COUNT(*) over `transition` per project (projectScene), a fifth in
// re-ranking sessions (rankSessions: the three live Mates have ~90
// sessions and thousands of turns each, and the Mate's current session is
// batched every pass), an eighth in walking the Mates' unchanged
// transcripts for tool commands (rememberToolCommand) - work a pass
// repeated though nothing had moved.
type memos struct {
	// mates is, by project, what the last committed poll read of the
	// project's Mate.
	mates map[string]*mateMemo
	// git is, by crew actor id, the key git was last asked under with
	// every answer in (gitKey).
	git map[string]string
	// projected is the projects whose scene projection has run on this
	// database (projectScene).
	projected map[string]bool
	// crews is, by project, the crews the last poll read: every crew it
	// did not find sealed. The budget, context and no-turn checks of a
	// poll that is not a sweep read only these (scopeCrews).
	crews map[string][]crewRecord
	// gitDirs is, by repository directory, the directory its refs are in;
	// "" when it could not be found (gitDir).
	gitDirs map[string]string
}

// statPair is a file's size and modification time; a file that is not
// there is a size of -1.
type statPair struct{ size, mtime int64 }

func statOf(path string) statPair {
	size, mtime := statKey(path)
	return statPair{size, mtime}
}

func (s statPair) String() string { return fmt.Sprintf("%d@%d", s.size, s.mtime) }

// mateMemo is what one pass read of a project's Mate: `mate.meta` and the
// archived sessions' directory as they were stat'ed before the read, and
// every transcript of the Mate the pass read or found recorded at rest.
// A Mate whose files all still match is quiet: the next poll skips its
// session rows - and with them rankSessions over its many sessions - its
// started and stopped events, the locate and the read of its transcripts
// and the walk of their tool calls. Nothing else of a pass depends on
// them: a Mate's own shell commands date no status line and its commits
// are no crew's (statusClock, commitSightings).
//
// A Mate whose transcript could not be located is never quiet: the file
// the harness's own rules look for may appear without mate.meta moving
// (a Mate before its first Stop hook).
type mateMemo struct {
	meta, sessions statPair
	transcripts    map[string]statPair
	unlocated      bool
}

// beginMate stats the Mate's files for this pass, before mate.meta is
// read, and says whether the pass may skip the Mate: a poll that is not a
// sweep, whose project's last committed poll read the Mate under the same
// stats.
func (p *pass) beginMate() {
	if !p.poll {
		return
	}
	ws := p.ing.ws
	p.mateRead = &mateMemo{
		meta:        statOf(ws.MateMeta(p.project)),
		sessions:    statOf(filepath.Join(ws.MateDir(p.project), "sessions")),
		transcripts: map[string]statPair{},
	}
	if !p.memo {
		return
	}
	last := p.ing.memos.mates[p.project]
	if last == nil || last.unlocated || last.meta != p.mateRead.meta || last.sessions != p.mateRead.sessions {
		return
	}
	for path, st := range last.transcripts {
		if statOf(path) != st {
			return
		}
	}
	p.mateQuiet = true
	p.mateRead.transcripts = last.transcripts
}

// mateTranscript notes a transcript of the Mate this pass reads, stat'ed
// before it is read.
func (p *pass) mateTranscript(path string) {
	if p.mateRead != nil && path != "" {
		p.mateRead.transcripts[path] = statOf(path)
	}
}

// mateUnlocated notes that a transcript of the Mate could not be located.
func (p *pass) mateUnlocated() {
	if p.mateRead != nil {
		p.mateRead.unlocated = true
	}
}

// gitKey is what git can tell about one crew depending on, as stats: its
// `.meta` and `.status` (a new status line, a close, a relaunch), the
// transcript this pass located (a commit the crew made shows there), and
// the refs of its branch and of the default branch, loose or packed (a
// tip that moved, a merge that landed). "" when the repository's refs
// cannot be found, which asks git every pass as before.
//
// What it misses: a ref in a reftable, or moved by a git that rewrote a
// ref file to the same size within one mtime tick. The sweep asks again.
func (p *pass) gitKey(crew crewRecord, repo, branch, base string) string {
	dir := p.ing.gitDir(repo)
	if dir == "" {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d@%d %d@%d", crew.key.metaSize, crew.key.metaMTime, crew.key.statusSize, crew.key.statusMTime)
	if path := p.located[crew.ActorID]; path != "" {
		fmt.Fprintf(&b, " t%s", statOf(path))
	}
	fmt.Fprintf(&b, " p%s", statOf(filepath.Join(dir, "packed-refs")))
	for _, ref := range []string{branch, base} {
		fmt.Fprintf(&b, " h%s r%s",
			statOf(filepath.Join(dir, "refs", "heads", filepath.FromSlash(ref))),
			statOf(filepath.Join(dir, "refs", "remotes", filepath.FromSlash(ref))))
	}
	return b.String()
}

// gitDir is the directory a repository's refs live in: its `.git`, or for
// a linked worktree or a submodule the common directory its `.git` file
// leads to. "" when there is none to find.
func (i *Ingester) gitDir(repo string) string {
	if dir, ok := i.memos.gitDirs[repo]; ok {
		return dir
	}
	dir := findGitDir(repo)
	if i.memos.gitDirs == nil {
		i.memos.gitDirs = map[string]string{}
	}
	i.memos.gitDirs[repo] = dir
	return dir
}

func findGitDir(repo string) string {
	dotgit := filepath.Join(repo, ".git")
	fi, err := os.Stat(dotgit)
	if err != nil {
		return ""
	}
	if fi.IsDir() {
		return dotgit
	}
	gitdir := readLine(dotgit, "gitdir:")
	if gitdir == "" {
		return ""
	}
	if !filepath.IsAbs(gitdir) {
		gitdir = filepath.Join(repo, gitdir)
	}
	common := readLine(filepath.Join(gitdir, "commondir"), "")
	if common == "" {
		return gitdir
	}
	if !filepath.IsAbs(common) {
		common = filepath.Join(gitdir, common)
	}
	return filepath.Clean(common)
}

// readLine is the first line of path, after prefix when one is given.
func readLine(path, prefix string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		return ""
	}
	line := strings.TrimSpace(sc.Text())
	if prefix != "" {
		rest, ok := strings.CutPrefix(line, prefix)
		if !ok {
			return ""
		}
		line = strings.TrimSpace(rest)
	}
	return line
}

// scopeCrews is the crews a check of the ledger (budget, context, no turn)
// reads for project: on a sweep every `crews/<id>.meta`, as before; on any
// other poll the crews its pass just read, which are every crew it did not
// find sealed - a sealed crew is closed an hour or more, its spend and its
// turns on the table, so nothing a check concludes about it can change
// until something the sweep reads does. A project no poll has read yet
// reads every meta.
func (i *Ingester) scopeCrews(project string, sweep bool) ([]crewRecord, error) {
	if !sweep {
		if crews, ok := i.memos.crews[project]; ok {
			return crews, nil
		}
	}
	return listCrewMetas(i.ws, project)
}

// anyOpen reports whether any of crews is open by its meta.
func anyOpen(crews []crewRecord) bool {
	for _, crew := range crews {
		if !crew.state().Closed() {
			return true
		}
	}
	return false
}

// noteCrews keeps the crews a poll's pass read, for scopeCrews. They are
// files, not rows, so they are kept whether or not the pass commits.
func (p *pass) noteCrews() {
	if !p.poll {
		return
	}
	if p.ing.memos.crews == nil {
		p.ing.memos.crews = map[string][]crewRecord{}
	}
	p.ing.memos.crews[p.project] = p.crews
}

// recordMemos keeps what a committed poll read under its keys. A pass
// past its deadline for git keeps none of its git keys: a call that ran
// out of time said nothing (decideSeals).
func (i *Ingester) recordMemos(p *pass) {
	if !p.poll {
		return
	}
	m := &i.memos
	if p.mateRead != nil {
		if m.mates == nil {
			m.mates = map[string]*mateMemo{}
		}
		m.mates[p.project] = p.mateRead
	}
	if !p.externalLate {
		for actorID, key := range p.gitAsked {
			if m.git == nil {
				m.git = map[string]string{}
			}
			m.git[actorID] = key
		}
	}
	if p.projected {
		if m.projected == nil {
			m.projected = map[string]bool{}
		}
		m.projected[p.project] = true
	}
}
