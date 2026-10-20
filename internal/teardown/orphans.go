package teardown

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/process"
)

// Orphans are the worktrees under root that something still points at but
// that no longer exist: a docker compose container whose working directory,
// or an Xcode DerivedData folder whose WorkspacePath, lies under
// `root/<name>` while `root/<name>` is gone (CON-114). A worktree that
// exists is never an orphan, so a live crew's resources are never offered
// for collection.
func Orphans(ctx context.Context, root string, run process.Runner, derivedRoot string) []Target {
	roots := []string{filepath.Clean(root)}
	if real, err := filepath.EvalSymlinks(root); err == nil && real != roots[0] {
		roots = append(roots, real)
	}
	seen := map[string]bool{}
	var out []Target
	add := func(path string) {
		for _, r := range roots {
			rel, err := filepath.Rel(r, filepath.Clean(path))
			if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
				continue
			}
			name := strings.SplitN(rel, string(filepath.Separator), 2)[0]
			wt := filepath.Join(roots[0], name)
			if seen[wt] {
				return
			}
			if _, err := os.Stat(filepath.Join(r, name)); err == nil {
				return
			}
			seen[wt] = true
			out = append(out, Target{Worktree: wt})
			return
		}
	}

	if run == nil {
		run = process.ExecRunner{}
	}
	if res, err := run.Run(ctx, process.Spec{Name: "docker", Args: []string{"ps", "-a", "--filter", "label=" + composeDirLabel,
		"--format", `{{.Label "` + composeDirLabel + `"}}`}}); err == nil && res.ExitCode == 0 {
		for _, dir := range strings.Split(string(res.Stdout), "\n") {
			if dir = strings.TrimSpace(dir); dir != "" {
				add(dir)
			}
		}
	}
	if derivedRoot != "" {
		if entries, err := os.ReadDir(derivedRoot); err == nil {
			for _, e := range entries {
				if e.IsDir() {
					if ws := WorkspacePath(filepath.Join(derivedRoot, e.Name(), "info.plist")); ws != "" {
						add(ws)
					}
				}
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Worktree < out[j].Worktree })
	return out
}
