package teardown

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/process"
)

const (
	composeProjectLabel = "com.docker.compose.project"
	composeDirLabel     = "com.docker.compose.project.working_dir"
)

// Docker removes the compose projects a worktree started: every container,
// volume and network carrying the project's label. It removes them by label
// rather than with `docker compose down`, which needs the compose file the
// worktree removal is about to delete.
//
// A project belongs to the worktree when its containers record a working
// directory inside it, or when it carries the name compose derives from the
// worktree's directory and no container of it says otherwise: a reviewer's
// stack whose containers are already gone left only volumes.
type Docker struct {
	Runner process.Runner
}

// Sweep implements Sweep.
func (d Docker) Sweep(ctx context.Context, t Target) ([]string, error) {
	rows, err := d.lines(ctx, "ps", "-a", "--filter", "label="+composeProjectLabel,
		"--format", `{{.Label "`+composeProjectLabel+`"}}`+"\t"+`{{.Label "`+composeDirLabel+`"}}`)
	if err != nil {
		if missing(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("docker: %v", err)
	}
	// A project is the worktree's when its containers all record a working
	// directory inside it. One with a container anywhere else is shared
	// and left alone, even when it carries the worktree's name.
	inside, outside := map[string]bool{}, map[string]bool{}
	for _, row := range rows {
		project, dir, _ := strings.Cut(row, "\t")
		if project == "" {
			continue
		}
		if t.Owns(dir) {
			inside[project] = true
		} else {
			outside[project] = true
		}
	}
	if name := ComposeProjectName(t.Name()); name != "" {
		inside[name] = true
	}
	owned := map[string]bool{}
	for name := range inside {
		if !outside[name] {
			owned[name] = true
		}
	}

	var removed []string
	var failed []string
	for _, name := range sortedKeys(owned) {
		n, err := d.removeProject(ctx, name)
		if err != nil {
			failed = append(failed, fmt.Sprintf("docker %s: %v", name, err))
		}
		if n != "" {
			removed = append(removed, "docker "+name+" ("+n+")")
		}
	}
	if len(failed) > 0 {
		return removed, fmt.Errorf("%s", strings.Join(failed, "; "))
	}
	return removed, nil
}

// removeProject removes everything labelled with the compose project and
// says what, like "4 containers, 2 volumes, 1 network"; "" when there was
// nothing.
func (d Docker) removeProject(ctx context.Context, name string) (string, error) {
	filter := "label=" + composeProjectLabel + "=" + name
	var counts []string
	for _, kind := range []struct {
		list, rm []string
		noun     string
	}{
		{[]string{"ps", "-aq", "--filter", filter}, []string{"rm", "-f"}, "container"},
		{[]string{"volume", "ls", "-q", "--filter", filter}, []string{"volume", "rm", "-f"}, "volume"},
		{[]string{"network", "ls", "-q", "--filter", filter}, []string{"network", "rm"}, "network"},
	} {
		ids, err := d.lines(ctx, kind.list...)
		if err != nil {
			return strings.Join(counts, ", "), err
		}
		if len(ids) == 0 {
			continue
		}
		if _, err := d.lines(ctx, append(kind.rm, ids...)...); err != nil {
			return strings.Join(counts, ", "), err
		}
		counts = append(counts, plural(len(ids), kind.noun))
	}
	return strings.Join(counts, ", "), nil
}

// lines runs `docker args...` and returns its non-empty output lines. A
// non-zero exit - the daemon not running, an id that will not go - is an
// error carrying docker's last word.
func (d Docker) lines(ctx context.Context, args ...string) ([]string, error) {
	run := d.Runner
	if run == nil {
		run = process.ExecRunner{}
	}
	res, err := run.Run(ctx, process.Spec{Name: "docker", Args: args})
	if err != nil {
		return nil, err
	}
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("docker %s exited %d%s", args[0], res.ExitCode, lastLine(res.Stderr))
	}
	var out []string
	for _, l := range strings.Split(string(res.Stdout), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out, nil
}

var composeInvalid = regexp.MustCompile(`[^a-z0-9_-]`)

// ComposeProjectName is the project name `docker compose` derives from a
// directory name: lower-cased, every character other than a-z, 0-9, `_`
// and `-` dropped.
func ComposeProjectName(dir string) string {
	return composeInvalid.ReplaceAllString(strings.ToLower(dir), "")
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
