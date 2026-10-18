package main

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/matehost"
)

// cmdNews implements `mate news <open|update|close|list> <project> ...`
// (CON-78): how a headless Mate tells the captain something outside the
// chat - a crew's question only the captain can answer, a key, a review, a
// failure - and closes it once settled. The REPL shows the news and turns
// the captain's answer into a `[news #n]` turn.
func cmdNews(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return newUsageError("usage: mate news <open|update|close|list> <project> ...")
	}
	switch args[0] {
	case "open":
		return cmdNewsOpen(args[1:], stdout, stderr)
	case "update":
		return cmdNewsUpdate(args[1:], stdout, stderr)
	case "close":
		return cmdNewsClose(args[1:], stdout, stderr)
	case "list":
		return cmdNewsList(args[1:], stdout, stderr)
	default:
		return newUsageErrorf("unknown news subcommand %q (want open, update, close or list)", args[0])
	}
}

// lineFlags is a repeatable flag that keeps each value in order.
type lineFlags []string

func (l *lineFlags) String() string { return strings.Join(*l, "\n") }

func (l *lineFlags) Set(v string) error {
	*l = append(*l, v)
	return nil
}

func newsLog(workspaceFlag, project string) (string, error) {
	w, err := resolveWorkspace(workspaceFlag)
	if err != nil {
		return "", err
	}
	if err := requireProject(w, project); err != nil {
		return "", err
	}
	return matehost.NewsPath(w.MateDir(project)), nil
}

// optionKeys are the keys an option may have. The REPL's Topic page takes
// the other letters as its own keys, so an answer is never one of them.
const optionKeys = "abcd"

func cmdNewsOpen(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("news open", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, `usage: mate news open <project> --title "<one line>" [--crew <id>] [--kind <kind>] [--context "<line>"]... [--question "<question>"] [--option "a=<choice>"]... [--recommend a] [--workspace <dir>]`)
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	title := fs.String("title", "", "what the news is about, in one line")
	crew := fs.String("crew", "", "the crew it comes from")
	kind := fs.String("kind", "", "needs-decision, blocked, review, budget, failed, ...")
	question := fs.String("question", "", "what the captain is asked")
	recommend := fs.String("recommend", "", "the option you recommend")
	var context, options lineFlags
	fs.Var(&context, "context", "a line of context (repeatable)")
	fs.Var(&options, "option", "an answer as key=text (repeatable)")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return newUsageError("mate news open: want exactly 1 argument: <project>")
	}
	if strings.TrimSpace(*title) == "" {
		return newUsageError("mate news open: --title is required")
	}
	n := matehost.News{Title: *title, Crew: *crew, Kind: *kind, Question: *question, Context: context}
	if len(options) > len(optionKeys) {
		return newUsageErrorf("mate news open: at most %d options, got %d", len(optionKeys), len(options))
	}
	seen := map[string]bool{}
	for _, o := range options {
		key, text, ok := strings.Cut(o, "=")
		key = strings.TrimSpace(key)
		if !ok || !strings.Contains(optionKeys, key) || len(key) != 1 || strings.TrimSpace(text) == "" {
			return newUsageErrorf("mate news open: --option wants key=text with a key from a to d, got %q", o)
		}
		if seen[key] {
			return newUsageErrorf("mate news open: option %s is given twice", key)
		}
		seen[key] = true
		n.Options = append(n.Options, matehost.NewsOption{Key: key, Text: strings.TrimSpace(text), Recommended: key == *recommend})
	}
	if *recommend != "" && !seen[*recommend] {
		return newUsageErrorf("mate news open: --recommend %s names no option", *recommend)
	}
	path, err := newsLog(*workspaceFlag, fs.Arg(0))
	if err != nil {
		return err
	}
	id, err := matehost.OpenNews(path, n, time.Now())
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "#%d\n", id)
	return nil
}

func newsID(s string) (int, error) {
	id, err := strconv.Atoi(strings.TrimPrefix(s, "#"))
	if err != nil || id <= 0 {
		return 0, newUsageErrorf("want a news id like 3 or #3, got %q", s)
	}
	return id, nil
}

func cmdNewsUpdate(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("news update", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, `usage: mate news update <project> <id> "<text>" [--ask] [--workspace <dir>]`)
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	ask := fs.Bool("ask", false, "the text is the news's question, asked again; it replaces the question and its options")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 3 {
		fs.Usage()
		return newUsageError("mate news update: want exactly 3 arguments: <project> <id> <text>")
	}
	id, err := newsID(fs.Arg(1))
	if err != nil {
		return err
	}
	path, err := newsLog(*workspaceFlag, fs.Arg(0))
	if err != nil {
		return err
	}
	return matehost.UpdateNews(path, id, fs.Arg(2), *ask, time.Now())
}

func cmdNewsClose(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("news close", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, `usage: mate news close <project> <id> [--reason "<one line>"] [--dropped] [--workspace <dir>]`)
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	reason := fs.String("reason", "", "how it was settled, in one line")
	dropped := fs.Bool("dropped", false, "what it was about went away before the captain answered")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 2 {
		fs.Usage()
		return newUsageError("mate news close: want exactly 2 arguments: <project> <id>")
	}
	id, err := newsID(fs.Arg(1))
	if err != nil {
		return err
	}
	path, err := newsLog(*workspaceFlag, fs.Arg(0))
	if err != nil {
		return err
	}
	return matehost.CloseNews(path, id, *reason, *dropped, time.Now())
}

func cmdNewsList(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("news list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate news list <project> [--all] [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	all := fs.Bool("all", false, "closed news too")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return newUsageError("mate news list: want exactly 1 argument: <project>")
	}
	path, err := newsLog(*workspaceFlag, fs.Arg(0))
	if err != nil {
		return err
	}
	list, err := matehost.LoadNews(path)
	if err != nil {
		return err
	}
	shown := 0
	for _, n := range list {
		if n.Closed() && !*all {
			continue
		}
		shown++
		crew := ""
		if n.Crew != "" {
			crew = " " + n.Crew + ":"
		}
		fmt.Fprintf(stdout, "#%d %s%s %s\n", n.ID, n.Status, crew, n.Title)
		if n.Answer != "" {
			fmt.Fprintf(stdout, "    captain: %s\n", n.Answer)
		}
	}
	if shown == 0 {
		fmt.Fprintln(stdout, "no open news")
	}
	return nil
}
