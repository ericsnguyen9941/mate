package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/learnings"
)

// cmdLearnings implements `mate learnings open <project>` and `mate
// learnings close <project> <crew> <id> <outcome> --evidence "<text>"`
// (CON-131): the predictions a Mate's learning records make, and the one
// edit that closes one. `open` is what new-task reads before a brief;
// `close` is what task-learning runs for each prediction the brief
// carried. Neither reads anything but the records.
func cmdLearnings(args []string, stdout, stderr io.Writer) error {
	usage := func() {
		fmt.Fprintln(stderr, `usage: mate learnings open <project> | mate learnings close <project> <crew> <id> <outcome> --evidence "<task/revision: what it showed>" [--workspace <dir>]`)
	}
	if len(args) == 0 {
		usage()
		return newUsageError("mate learnings: want open or close")
	}
	switch args[0] {
	case "open":
		return cmdLearningsOpen(args[1:], stdout, stderr)
	case "close":
		return cmdLearningsClose(args[1:], stdout, stderr)
	}
	usage()
	return newUsageError("mate learnings: unknown subcommand " + args[0])
}

func cmdLearningsOpen(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("learnings open", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprintln(stderr, "usage: mate learnings open <project> [--workspace <dir>]") }
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return newUsageError("mate learnings open: want exactly 1 argument: <project>")
	}
	w, err := resolveWorkspace(*workspaceFlag)
	if err != nil {
		return err
	}
	project := fs.Arg(0)
	if err := requireProject(w, project); err != nil {
		return err
	}
	open, err := learnings.Open(w, project)
	if err != nil {
		return err
	}
	if len(open) == 0 {
		fmt.Fprintf(stdout, "no open predictions in %s\n", w.LearningsDir(project))
		return nil
	}
	fmt.Fprintf(stdout, "open predictions in %s (%d); close one with mate learnings close %s <crew> <id> %s --evidence \"<task/revision: what it showed>\"\n",
		w.LearningsDir(project), len(open), project, strings.Join(learnings.Outcomes, "|"))
	for _, p := range open {
		fmt.Fprintln(stdout, p.String())
	}
	return nil
}

func cmdLearningsClose(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("learnings close", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "usage: mate learnings close <project> <crew> <id> %s --evidence \"<task/revision: what it showed>\" [--workspace <dir>]\n", strings.Join(learnings.Outcomes, "|"))
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	evidence := fs.String("evidence", "", "the later task, its revision and what its evidence showed")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 4 {
		fs.Usage()
		return newUsageError("mate learnings close: want exactly 4 arguments: <project> <crew> <id> <outcome>")
	}
	w, err := resolveWorkspace(*workspaceFlag)
	if err != nil {
		return err
	}
	project, crew, id, outcome := fs.Arg(0), fs.Arg(1), fs.Arg(2), fs.Arg(3)
	if err := requireProject(w, project); err != nil {
		return err
	}
	p, err := learnings.Close(w, project, crew, id, outcome, *evidence, time.Now())
	if err != nil {
		if errors.Is(err, learnings.ErrBadOutcome) || errors.Is(err, learnings.ErrUnknownPrediction) || *evidence == "" {
			return &usageError{err}
		}
		return err
	}
	fmt.Fprintf(stdout, "closed %s #%s: outcome %s (%s) in %s\n", p.Crew, p.ID, p.Outcome, p.Evidence, p.Path)
	return nil
}
