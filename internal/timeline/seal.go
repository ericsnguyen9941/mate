package timeline

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/crewstate"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// A sealed crew is one the ingest has finished with: it is closed, closed for
// longer than settledFor, and every conclusion the files and git can give
// about it is already on the table - its transcript read to the end (or
// provably not there), every status line dated, its merge recorded or
// provably never coming. A poll then skips it whole: no `.meta` read, no
// harness profile, no transcript locate, no git, no actor, session or task
// upsert, no rankSessions (CON-259).
//
// Measured 2026-11-11 on the captain's workspace (1308 crew directories, a
// handful alive): after wave 0 of CON-253 a steady pass still cost 0.64 s and
// 78 MB of allocation, and nearly all of it was the closed crews - flush
// re-upserting their rows ~18%, ReadCrewHarnessProfiles ~12%, rankSessions
// ~11%, locating their transcripts ~9% (Codex adoption for crews that never
// adopted a rollout ~5%), git ~5%. Work in proportion to history, not to
// change.
//
// The seal lives in the Ingester's memory, keyed by the stat of the crew's
// `.meta` and `.status`, and not in the database:
//
//   - What can reopen a crew is a file write by another process - `mate crew
//     relaunch` rewriting the meta of a legacy crew closed by stopped_at
//     alone, `mate pr watch` recording a pull request's end, a person editing
//     it - and that process may not write the database (internal/db has one
//     writer). The file's own stat is therefore the key whatever holds it,
//     and a seal in a table would buy only the first pass after a restart,
//     which settledCrew and gitSettled (CON-140) already keep cheap.
//   - A table is a migration and a SchemaVersion bump, which makes every
//     reader on an older binary refuse the database until it is restarted.
//   - Like gitFacts, it is forgotten with the process, and it is never a
//     claim about what the timeline holds that a rollback could falsify: a
//     pass records its seals only after its transaction commits, and
//     Reindex - which empties every derived table - forgets them and reads
//     every crew again.
//
// The key is stat'ed before the meta is read, so a write racing the read
// leaves an old key on new content: the next pass sees a stat that does not
// match and reads the crew again. Size and mtime are the same identity the
// transcript cache and telemetry_cursor already trust. Every write of a meta
// replaces it whole (store.WriteCrewMeta) and a relaunch rewrites half its
// keys, so a reopened crew always reads as changed; an append to `.status`
// always grows it.
type sealKey struct {
	metaSize, metaMTime     int64
	statusSize, statusMTime int64
}

// crewSealKey stats a crew's two files. A file that is not there is a size
// of -1, so a status file appearing later is a change.
func crewSealKey(metaPath, statusPath string) sealKey {
	var k sealKey
	k.metaSize, k.metaMTime = statKey(metaPath)
	k.statusSize, k.statusMTime = statKey(statusPath)
	return k
}

func statKey(path string) (int64, int64) {
	fi, err := os.Stat(path)
	if err != nil {
		return -1, 0
	}
	return fi.Size(), fi.ModTime().UnixNano()
}

// unlocatedCrew is a closed crew's transcript locate that failed past
// settledFor, remembered against the crew's key. A Codex crew that never
// adopted a rollout was adopted again on every poll - a walk of every
// rollout since its launch - though its agent was gone and no rollout of its
// cwd could appear any more (CON-255 kept each rollout's head, CON-259 stops
// the walk). It is a negative answer about files outside the database, like
// gitFacts, so it is kept as soon as it is seen, not at commit.
type unlocatedCrew struct {
	key    sealKey
	reason string
}

// listCrews reads the project's crews for a pass: every `crews/<id>.meta`,
// except, on a poll, the crews sealed under a key their files still match.
// Those are only named (p.sealed), so a log line addressed to one still
// finds its actor and a digest that names one still links to it.
//
// A rebuild (Reindex) reads every crew: its pass runs on emptied tables.
func (p *pass) listCrews() error {
	ws := p.ing.ws
	dir := ws.CrewsDir(p.project)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	held := p.ing.sealedIn(p.project)
	p.sealed = map[string]sealKey{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		id, ok := strings.CutSuffix(e.Name(), ".meta")
		if !ok || store.ValidateCrewID(id) != nil {
			continue
		}
		path := filepath.Join(dir, e.Name())
		key := crewSealKey(path, ws.CrewStatus(p.project, id))
		actorID := CrewActorID(p.project, id)
		if p.poll {
			if sealed, ok := held[actorID]; ok && sealed == key {
				p.sealed[id] = key
				continue
			}
		}
		meta, err := ws.ReadCrewMeta(p.project, id)
		if err != nil {
			// A half-written meta is not evidence of anything; the next pass
			// reads it again.
			continue
		}
		p.crews = append(p.crews, crewRecord{ID: id, ActorID: actorID, Meta: meta, Path: path, key: key})
	}
	sort.Slice(p.crews, func(i, j int) bool { return p.crews[i].ID < p.crews[j].ID })
	return nil
}

// knownCrew reports whether the pass knows crew by its meta: read this pass,
// or sealed.
func (p *pass) knownCrew(crew string) bool {
	if _, ok := p.sealed[crew]; ok {
		return true
	}
	for _, known := range p.crews {
		if known.ID == crew {
			return true
		}
	}
	return false
}

// closedLongAgo reports whether a crew closed more than settledFor ago. A
// crew closed without a stopped_at (a spawn that failed) is dated by its
// meta's mtime, which is no earlier than the write that closed it.
func (p *pass) closedLongAgo(crew crewRecord) bool {
	if !crew.state().Closed() {
		return false
	}
	stopped := metaTime(crew.Meta, MetaStoppedAt)
	if stopped.IsZero() && crew.key.metaSize >= 0 {
		stopped = time.Unix(0, crew.key.metaMTime)
	}
	return !stopped.IsZero() && p.now.Sub(stopped) >= settledFor
}

// decideSeals runs after the pass has flushed, inside its transaction, so
// what it reads of the tables includes what this pass just wrote. A crew is
// sealed when, closed long ago:
//
//   - its transcript is on the table as it is on disk (recordedAtRest: the
//     telemetry cursor is written with the transcript's facts), or it could
//     not be located at all - an hour after its agent stopped, nothing will
//     write the file the locate is looking for;
//   - every line of its `.status` is read (the cursor is at the end), so no
//     line waits for the transcript's clock;
//   - git has nothing left to tell (gitDone).
//
// A pass that was past its deadline for git and Herdr by the end of its git
// phase seals nothing: a call that ran out of time said nothing, which is
// not an answer, and a sealed crew is never asked again. The next pass
// decides instead.
func (p *pass) decideSeals(ctx context.Context) error {
	if !p.poll || p.externalLate {
		return nil
	}
	for _, crew := range p.crews {
		if !p.closedLongAgo(crew) {
			continue
		}
		if path, ok := p.located[crew.ActorID]; ok {
			recorded, err := p.recordedAtRest(ctx, path)
			if err != nil {
				return err
			}
			if !recorded {
				continue
			}
		} else if _, failed := p.unlocated[crew.ActorID]; !failed {
			// Neither located nor failed: the pass did not get to it.
			continue
		}
		status := p.ing.ws.CrewStatus(p.project, crew.ID)
		if crew.key.statusSize >= 0 {
			read, err := p.readCursor(ctx, status)
			if err != nil {
				return err
			}
			if read < crew.key.statusSize {
				continue
			}
		}
		done, err := p.gitDone(ctx, crew)
		if err != nil {
			return err
		}
		if !done {
			continue
		}
		if p.seals == nil {
			p.seals = map[string]sealKey{}
		}
		p.seals[crew.ActorID] = crew.key
	}
	return nil
}

// gitDone reports whether git can tell nothing more about a closed crew:
// it has no branch, did not close `finished` (only a finished crew's merge is
// looked for, mergeEvent), its merge is recorded, or its merge can never be
// proved - the branch is gone and the timeline holds no commit of it, so
// mergeEvent has no tip to test, and gitFacts already holds that a closed
// crew's branch does not come back. A finished crew whose branch is still
// there, or whose recorded tip has not landed, may still be merged and stays
// unsealed, as gitSettled keeps asking about it.
func (p *pass) gitDone(ctx context.Context, crew crewRecord) (bool, error) {
	branch := strings.TrimSpace(crew.Meta[MetaBranch])
	if branch == "" || crew.state() != crewstate.StateFinished {
		return true, nil
	}
	if _, merged := p.w.eventID(ctx, dedup(KindMergeDone, crew.ActorID)); merged {
		return true, nil
	}
	cfg, err := p.ing.ws.LoadProject(p.project)
	if err != nil {
		return false, nil //nolint:nilerr // an unreadable project.yaml concludes nothing
	}
	repoCfg, err := cfg.CrewRepo(crew.Meta)
	if err != nil {
		return false, nil //nolint:nilerr
	}
	if !p.ing.gitFacts().gone[repoRef{repo: p.ing.ws.RepoDir(repoCfg.Path), ref: branch}] {
		return false, nil
	}
	tip, err := p.lastRecordedCommit(ctx, crew.ActorID)
	if err != nil {
		return false, err
	}
	return tip == "", nil
}

// sealedIn is the project's seals as the last committed pass left them.
func (i *Ingester) sealedIn(project string) map[string]sealKey {
	return i.sealed[project]
}

// recordSeals keeps what a committed pass decided: the crews it skipped
// because they were sealed, and the ones it sealed now. A crew it read and
// did not seal is dropped, and so is a crew whose meta is gone.
func (i *Ingester) recordSeals(p *pass) {
	next := make(map[string]sealKey, len(p.sealed)+len(p.seals))
	for id, key := range p.sealed {
		next[CrewActorID(p.project, id)] = key
	}
	for actorID, key := range p.seals {
		next[actorID] = key
	}
	if i.sealed == nil {
		i.sealed = map[string]map[string]sealKey{}
	}
	i.sealed[p.project] = next
}

// rememberedUnlocated is the reason a closed crew's locate failed on an
// earlier pass, when the crew's files have not changed since.
func (p *pass) rememberedUnlocated(crew crewRecord) (string, bool) {
	if !p.poll || !p.closedLongAgo(crew) {
		return "", false
	}
	u, ok := p.ing.unlocated[crew.ActorID]
	if !ok || u.key != crew.key {
		return "", false
	}
	return u.reason, true
}

func (p *pass) rememberUnlocated(crew crewRecord, reason string) {
	if !p.poll || !p.closedLongAgo(crew) {
		return
	}
	if p.ing.unlocated == nil {
		p.ing.unlocated = map[string]unlocatedCrew{}
	}
	p.ing.unlocated[crew.ActorID] = unlocatedCrew{key: crew.key, reason: reason}
}
