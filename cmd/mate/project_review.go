package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// cmdProjectReview implements `mate project review <project> [on|off]`
// (CON-101): whether each task's work is independently reviewed before the
// Mate lands it. Without a value it prints the setting; with one it records
// the captain's answer in project.yaml, which `mate merge` reads.
func cmdProjectReview(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("project review", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate project review <project> [on|off] [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() < 1 || fs.NArg() > 2 {
		fs.Usage()
		return newUsageError("mate project review: want <project> and optionally on or off")
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
		if !store.ValidReview(value) {
			fs.Usage()
			return newUsageErrorf("mate project review: want on or off, not %q", value)
		}
		cfg.Review = value
		if err := w.SaveProject(name, cfg); err != nil {
			return err
		}
	}
	fmt.Fprintln(stdout, reviewSettingLine(name, cfg.Review))
	return nil
}

// reviewSettingLine says what the setting means for landing a task, in the
// words `wake` also prints.
func reviewSettingLine(project, review string) string {
	switch review {
	case store.ReviewOn:
		return fmt.Sprintf("review: on - %s's tasks land after an independent review passes (`mate merge --review <reviewer>`)", project)
	case store.ReviewOff:
		return fmt.Sprintf("review: off - %s's tasks land after you check the hand-back against the brief (`mate merge` without --review)", project)
	default:
		return fmt.Sprintf("review: undecided - ask the captain whether %s's tasks get an independent review, then record it with `mate project review %s on|off`", project, project)
	}
}
