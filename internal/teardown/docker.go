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

// Sweep implements Sweep. It lists containers, volumes and networks once
// each, whatever the number of targets, and removes by id.
func (d Docker) Sweep(ctx context.Context, targets []Target, dry bool) ([]string, error) {
	containers, err := d.lines(ctx, "ps", "-a", "--filter", "label="+composeProjectLabel,
		"--format", "{{.ID}}\t"+labelFormat(composeProjectLabel)+"\t"+labelFormat(composeDirLabel))
	if err != nil {
		if missing(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("docker: %v", err)
	}
	volumes, err := d.lines(ctx, "volume", "ls", "--filter", "label="+composeProjectLabel, "--format", "{{.Name}}\t"+labelFormat(composeProjectLabel))
	if err != nil {
		return nil, fmt.Errorf("docker: %v", err)
	}
	networks, err := d.lines(ctx, "network", "ls", "--filter", "label="+composeProjectLabel, "--format", "{{.ID}}\t"+labelFormat(composeProjectLabel))
	if err != nil {
		return nil, fmt.Errorf("docker: %v", err)
	}

	// A project is the targets' when its containers all record a working
	// directory inside one of them. One with a container anywhere else is
	// shared and left alone, even when it carries a target's name.
	inside, outside := map[string]bool{}, map[string]bool{}
	byProject := map[string]*composeProject{}
	get := func(name string) *composeProject {
		if byProject[name] == nil {
			byProject[name] = &composeProject{}
		}
		return byProject[name]
	}
	for _, row := range containers {
		f := strings.SplitN(row, "\t", 3)
		if len(f) < 2 || f[1] == "" {
			continue
		}
		dir := ""
		if len(f) == 3 {
			dir = f[2]
		}
		if owns(targets, dir) {
			inside[f[1]] = true
		} else {
			outside[f[1]] = true
		}
		get(f[1]).containers = append(get(f[1]).containers, f[0])
	}
	for _, row := range volumes {
		if id, project, ok := strings.Cut(row, "\t"); ok && project != "" {
			get(project).volumes = append(get(project).volumes, id)
		}
	}
	for _, row := range networks {
		if id, project, ok := strings.Cut(row, "\t"); ok && project != "" {
			get(project).networks = append(get(project).networks, id)
		}
	}
	for _, t := range targets {
		if name := ComposeProjectName(t.Name()); name != "" && byProject[name] != nil {
			inside[name] = true
		}
	}

	var removed, failed []string
	for _, name := range sortedKeys(inside) {
		if outside[name] {
			continue
		}
		n, err := d.removeProject(ctx, byProject[name], dry)
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

// composeProject is what docker holds under one compose project label.
type composeProject struct {
	containers, volumes, networks []string
}

// removeProject removes a project's containers, then volumes, then
// networks, and says what, like "4 containers, 2 volumes, 1 network".
func (d Docker) removeProject(ctx context.Context, p *composeProject, dry bool) (string, error) {
	var counts []string
	for _, kind := range []struct {
		ids  []string
		rm   []string
		noun string
	}{
		{p.containers, []string{"rm", "-f"}, "container"},
		{p.volumes, []string{"volume", "rm", "-f"}, "volume"},
		{p.networks, []string{"network", "rm"}, "network"},
	} {
		if len(kind.ids) == 0 {
			continue
		}
		if !dry {
			if _, err := d.lines(ctx, append(kind.rm, kind.ids...)...); err != nil {
				return strings.Join(counts, ", "), err
			}
		}
		counts = append(counts, plural(len(kind.ids), kind.noun))
	}
	return strings.Join(counts, ", "), nil
}

func labelFormat(label string) string {
	return `{{.Label "` + label + `"}}`
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
