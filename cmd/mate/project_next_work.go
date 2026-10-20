package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// cmdProjectNextWork implements `mate project next-work <project> [pick|ask]`
// (CON-120): what the Mate does when a task lands and nothing the captain
// asked for is left. Without a value it prints the setting; with one it
// records the captain's choice in project.yaml.
func cmdProjectNextWork(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("project next-work", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate project next-work <project> [pick|ask] [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() < 1 || fs.NArg() > 2 {
		fs.Usage()
		return newUsageError("mate project next-work: want <project> and optionally pick or ask")
	}
	name := fs.Arg(0)
	w, err := resolveWorkspace(*workspaceFlag)
	if err != nil {
		return err
	}
	if _, ok := w.Project(name); !ok {
		return fmt.Errorf("%w: %s", store.ErrNoProject, name)
	}
	cfg, err := w.LoadProject(name)
	if err != nil {
		return err
	}
	if fs.NArg() == 2 {
		value := fs.Arg(1)
		if !store.ValidNextWork(value) {
			fs.Usage()
			return newUsageErrorf("mate project next-work: want pick or ask, not %q", value)
		}
		cfg.NextWork = value
		if err := w.SaveProject(name, cfg); err != nil {
			return err
		}
	}
	fmt.Fprintln(stdout, nextWorkLine(cfg))
	return nil
}

// nextWorkLine says what the setting means once a task lands, in the words
// recall and `mate merge` also print.
func nextWorkLine(cfg store.ProjectConfig) string {
	if cfg.NextWorkMode() == store.NextWorkAsk {
		return "next work: ask - when nothing the captain asked for is left, ask the captain which work comes next, recommending one"
	}
	return "next work: pick - when nothing the captain asked for is left, dispatch the next ready task yourself and tell the captain in one line; ask only when nothing is ready"
}
