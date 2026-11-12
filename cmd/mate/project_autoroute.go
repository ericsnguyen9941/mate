package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// cmdProjectAutoroute implements `mate project autoroute <project> [off|shadow|on]`
// (CON-273, CON-274): whether the app takes a hand-back off the Mate when its
// checks and Jev find nothing to decide. Without a value it prints the
// setting; with one it records the captain's choice in project.yaml.
func cmdProjectAutoroute(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("project autoroute", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate project autoroute <project> [off|shadow|on] [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() < 1 || fs.NArg() > 2 {
		fs.Usage()
		return newUsageError("mate project autoroute: want <project> and optionally off, shadow or on")
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
		if !store.ValidAutoroute(value) {
			fs.Usage()
			return newUsageErrorf("mate project autoroute: want off, shadow or on, not %q", value)
		}
		cfg.Autoroute = value
		if err := w.SaveProject(name, cfg); err != nil {
			return err
		}
	}
	fmt.Fprintln(stdout, autorouteLine(cfg))
	return nil
}

// autorouteLine says what the setting means for the next hand-back.
func autorouteLine(cfg store.ProjectConfig) string {
	switch cfg.AutorouteMode() {
	case store.AutorouteOn:
		return "autoroute: on - a hand-back whose checks pass and that Jev finds routine is landed by the app, a review's in-scope findings go back to the crew, and the Mate hears of both in its next digest; everything else still goes to the Mate"
	case store.AutorouteShadow:
		return "autoroute: shadow - every hand-back still goes to the Mate; what the app would have done is recorded in mate/route.jsonl"
	}
	return "autoroute: off - every hand-back goes to the Mate"
}
