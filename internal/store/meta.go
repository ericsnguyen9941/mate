package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// A `.meta` file is `key=value`, one key per line, keys `[a-z_]+` and values a
// single line. It is written whole and atomically, and it is small enough that
// the store never edits one in place.

// ReadMeta parses a `.meta` file. A missing file is not an error: it returns an
// empty map, because a crew that has not been recorded yet has no meta.
//
// ReadMeta and WriteMeta are the file format only; they do not check the
// workspace boundary. Use the Workspace methods for anything under `.mate/`.
func ReadMeta(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	out := map[string]string{}
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimRight(scanner.Text(), "\r")
		if text == "" {
			continue
		}
		key, value, ok := strings.Cut(text, "=")
		if !ok {
			return nil, fmt.Errorf("store: %s line %d: no '=' in %q", path, line, text)
		}
		if !metaKeyPattern.MatchString(key) {
			return nil, fmt.Errorf("store: %s line %d: invalid key %q: want [a-z_]+", path, line, key)
		}
		out[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// WriteMeta writes a `.meta` file atomically, keys in sorted order so the file
// is byte-identical for the same map and a diff shows only real changes. A key
// that is not `[a-z_]+`, or a value containing a newline, is refused and
// nothing is written.
func WriteMeta(path string, meta map[string]string) error {
	data, err := encodeMeta(meta)
	if err != nil {
		return err
	}
	return writeFileAtomic(path, data, 0o644)
}

func encodeMeta(meta map[string]string) ([]byte, error) {
	keys := make([]string, 0, len(meta))
	for k := range meta {
		if !metaKeyPattern.MatchString(k) {
			return nil, fmt.Errorf("store: invalid meta key %q: want [a-z_]+", k)
		}
		if strings.ContainsAny(meta[k], "\n\r") {
			return nil, fmt.Errorf("store: meta value for %q contains a newline", k)
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(meta[k])
		b.WriteByte('\n')
	}
	return []byte(b.String()), nil
}

// writeFileAtomic is the boundary-free form of Workspace.writeFile, used by the
// exported file-format helpers.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}

// MetaBudget is the crew `.meta` key holding the crew's own token budget
// (CON-76): the dispatch rule's budget `mate crew spawn --budget` recorded,
// else a copy of the project's crew_tokens at spawn. The observer's budget
// check reads it before the project's block. It lives here, beside
// MetaRepo, because both the writer (internal/spawn) and the reader
// (internal/timeline) already import this package.
const MetaBudget = "budget"

// The crew `.meta` keys of CON-94's budget enforcement. budget_line is when
// the daemon typed the "budget reached" line into the crew's pane, so it is
// typed once per ceiling; budget_stop is what the daemon did at the hard
// ceiling, so it acts once. budget_extensions counts `mate crew budget
// --extend`, budget_reason is the last extension's reason, and
// budget_original the ceiling the crew was spawned with. An extension
// deletes budget_line and budget_stop, which re-arms both against the new
// ceiling.
const (
	MetaBudgetLine       = "budget_line"
	MetaBudgetStop       = "budget_stop"
	MetaBudgetExtensions = "budget_extensions"
	MetaBudgetReason     = "budget_reason"
	MetaBudgetOriginal   = "budget_original"
)

// MetaFollowups is the crew `.meta` key `mate merge --followups` writes
// (CON-74): where the findings a passing review left standing were
// recorded before the merge - a `backlog:<text>` line under the project's
// `## Queued`, or a tracker issue id such as `HVN-12`. It is the record a
// reader of the closed crew (`mate usage`, the Office) can show for "what
// happened to the reviewer's findings".
const MetaFollowups = "followups"

// ReadCrewMeta reads `crews/<crew>.meta`.
func (w *Workspace) ReadCrewMeta(project, crew string) (map[string]string, error) {
	if err := validatePair(project, crew); err != nil {
		return nil, err
	}
	return ReadMeta(w.CrewMeta(project, crew))
}

// WriteCrewMeta writes `crews/<crew>.meta` atomically, inside the workspace.
func (w *Workspace) WriteCrewMeta(project, crew string, meta map[string]string) error {
	if err := validatePair(project, crew); err != nil {
		return err
	}
	data, err := encodeMeta(meta)
	if err != nil {
		return err
	}
	return w.writeFile(w.CrewMeta(project, crew), data, 0o644)
}

// UpdateCrewMeta sets keys in `crews/<crew>.meta`: the file is read again at
// the moment of the write, the keys in set are applied and every other key is
// kept, so a writer that only means to add a key (the pull request watcher)
// does not write back a stale copy of the rest. A crew with no meta is
// refused: nothing here creates one.
func (w *Workspace) UpdateCrewMeta(project, crew string, set map[string]string) error {
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil {
		return err
	}
	if len(meta) == 0 {
		return fmt.Errorf("store: crew %s/%s has no meta to update", project, crew)
	}
	for k, v := range set {
		meta[k] = v
	}
	return w.WriteCrewMeta(project, crew, meta)
}

// CrewIDs are the ids with a `.meta` under the project's `crews/`, sorted. A
// file whose name is not a valid crew id is ignored rather than reported:
// `crews/` is a directory the user can also put things in.
func (w *Workspace) CrewIDs(project string) ([]string, error) {
	if err := ValidateProjectName(project); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(w.CrewsDir(project))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		id, ok := strings.CutSuffix(e.Name(), ".meta")
		if !ok || ValidateCrewID(id) != nil {
			continue
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, nil
}

func validatePair(project, crew string) error {
	if err := ValidateProjectName(project); err != nil {
		return err
	}
	return ValidateCrewID(crew)
}
