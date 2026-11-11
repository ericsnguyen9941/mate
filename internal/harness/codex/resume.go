package codex

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// codexRolloutHeadBytes bounds the read of a rollout's first record. The
// session_meta record carries Codex's base instructions and measured about
// 18 KiB on codex-cli 0.154.0 (2026-09-24); the bound leaves room for that to
// grow without ever reading a multi-megabyte rollout whole.
const codexRolloutHeadBytes = 256 * 1024

// CodexRolloutPath finds the rollout Codex wrote for sessionID under
// sessionsDir (CodexSessionsDir). Codex names every rollout
// `rollout-<local time>-<session id>.jsonl`, so the file name alone answers
// the question and no rollout is opened.
//
// It is what `mate start` asks before it launches `codex resume <id>`: an id
// with no rollout makes codex-cli 0.154.0 print "No saved session found with
// ID ..." and exit to the shell, which Herdr reports only as an agent-start
// timeout a minute later (measured 2026-09-24).
func CodexRolloutPath(sessionsDir, sessionID string) (string, bool) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionsDir == "" || sessionID == "" {
		return "", false
	}
	suffix := "-" + sessionID + ".jsonl"
	var found string
	_ = filepath.WalkDir(sessionsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || found != "" {
			return nil //nolint:nilerr // an unreadable branch is not an answer
		}
		if !d.IsDir() && strings.HasPrefix(d.Name(), "rollout-") && strings.HasSuffix(d.Name(), suffix) {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	return found, found != ""
}

// CodexRolloutCandidatesSince lists the rollouts under sessionsDir whose day
// directory (`YYYY/MM/DD`, Codex's local date) is not before since's local
// date minus one day, each with its validated session_meta. The day bound is
// what keeps this cheap on a machine with months of Codex history. Codex
// files a rollout under the local date it created it (codex-cli 0.154.0:
// `2026/10/08/rollout-2026-10-08T00-35-39-...` for a session_meta stamped
// 2026-10-07T17:40Z at +07:00), and it creates it after the launch, so in
// one clock zone the launch's own date is already enough; the extra day
// absorbs a Codex whose zone is behind mate's. A zero since bounds nothing.
// The candidates go to AdoptCodexRollout, which applies the real rule (cwd
// and a session_meta timestamp at or after launch).
func CodexRolloutCandidatesSince(sessionsDir string, since time.Time) ([]CodexRolloutCandidate, error) {
	floor := since.Local().AddDate(0, 0, -1)
	floorDay := time.Date(floor.Year(), floor.Month(), floor.Day(), 0, 0, 0, 0, time.Local)
	var out []CodexRolloutCandidate
	err := filepath.WalkDir(sessionsDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if path == sessionsDir {
				return walkErr
			}
			return nil //nolint:nilerr // an unreadable branch is skipped, not fatal
		}
		if d.IsDir() {
			if path != sessionsDir && dayDirBefore(sessionsDir, path, floorDay) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasPrefix(d.Name(), "rollout-") || !strings.HasSuffix(d.Name(), ".jsonl") {
			return nil
		}
		meta, ok := rolloutMetas.read(path, d)
		if !ok {
			return nil
		}
		out = append(out, CodexRolloutCandidate{Path: path, Meta: meta})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Path < out[b].Path })
	return out, nil
}

// dayDirBefore reports whether dir is a `YYYY`, `YYYY/MM` or `YYYY/MM/DD`
// directory under root that ends before floor. Anything that does not parse
// as such is walked, never skipped: an unknown layout must not hide a rollout.
func dayDirBefore(root, dir string, floor time.Time) bool {
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		return false
	}
	parts := strings.Split(rel, string(filepath.Separator))
	nums := make([]int, 0, 3)
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return false
		}
		nums = append(nums, n)
	}
	switch len(nums) {
	case 1:
		return nums[0] < floor.Year()
	case 2:
		end := time.Date(nums[0], time.Month(nums[1])+1, 1, 0, 0, 0, 0, time.Local)
		return !end.After(floor)
	case 3:
		end := time.Date(nums[0], time.Month(nums[1]), nums[2]+1, 0, 0, 0, 0, time.Local)
		return !end.After(floor)
	}
	return false
}

// rolloutMetas is every rollout's session_meta this process has parsed.
// The ingest adopts a rollout for each Codex crew it has not located yet,
// on every poll, and each adoption walks every rollout since the crew's
// launch: reading each one's head again cost 1.2 GB of allocation per
// five-second pass on a workspace with about twenty such crews (CON-255).
var rolloutMetas = &rolloutMetaCache{byPath: map[string]rolloutMeta{}}

// rolloutMetaCache keeps a rollout's parsed session_meta by path. A rollout
// is only ever appended to, so its first record does not change once it is
// whole; a file smaller than when it was parsed is not the same file and is
// read again. A head that does not parse yet is not kept: Codex may still
// be writing it.
type rolloutMetaCache struct {
	mu     sync.Mutex
	byPath map[string]rolloutMeta
}

type rolloutMeta struct {
	size int64
	meta CodexSessionMeta
}

func (c *rolloutMetaCache) read(path string, d os.DirEntry) (CodexSessionMeta, bool) {
	info, err := d.Info()
	if err != nil {
		return CodexSessionMeta{}, false
	}
	c.mu.Lock()
	held, ok := c.byPath[path]
	c.mu.Unlock()
	if ok && info.Size() >= held.size {
		return held.meta, true
	}
	head, err := readFileHead(path, codexRolloutHeadBytes)
	if err != nil {
		return CodexSessionMeta{}, false
	}
	meta, err := ParseCodexSessionMeta(head)
	if err != nil {
		return CodexSessionMeta{}, false
	}
	c.mu.Lock()
	c.byPath[path] = rolloutMeta{size: info.Size(), meta: meta}
	c.mu.Unlock()
	return meta, true
}

// readFileHead reads at most n bytes from the start of path, and no more
// than the file holds: most rollouts are far smaller than the bound.
func readFileHead(path string, n int) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if info, err := f.Stat(); err == nil && info.Size() < int64(n) {
		n = int(info.Size())
	}
	buf := make([]byte, n)
	read, err := io.ReadFull(f, buf)
	if read == 0 && err != nil {
		return nil, err
	}
	return buf[:read], nil
}
