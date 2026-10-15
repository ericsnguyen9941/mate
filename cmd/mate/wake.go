package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nguyenngocanh94/mate/internal/autopilot"
	"github.com/nguyenngocanh94/mate/internal/box"
	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// cmdWake implements `mate wake <project> [--peek] [--workspace <dir>]` (CON-68): the
// one read a Mate makes when a digest or a `resolve:` line wakes it. On
// hellovietnam (2026-10-06..14) a Mate woke 247 times and spent 6.9 model
// calls per wake, most of them a chain of tail, cat, peek, crew list and
// python run to learn what had happened; each call re-read its ~65k token
// prefix. This prints that in one output:
//
//  1. every open item, in the order the Mate acts on them (blocked,
//     needs-decision, wait-mate - event-handling's order), each with what
//     deciding it needs: a question in full, a hand-back's branch summary,
//     its `## For the next crew` section and its reviews;
//  2. the status lines written since the Mate's last wake;
//  3. each open Crew's state.
//
// It reads the same view the digest is built from (internal/box), so a
// wake never disagrees with the digest that caused it. Part 2's cursor is
// its own (`mate/.wake-cursor`), moved by every run but a --peek, so a
// captain who looks at a wake does not take its lines from the Mate's next
// one; part 1 has none, so an item stays listed until it is dealt with,
// whatever was shown before. It asks Herdr nothing: a pane's health is `mate state`'s, and `mate peek`
// is still how a Mate looks at a blocked Crew.
func cmdWake(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("wake", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mate wake <project> [--peek] [--workspace <dir>]")
	}
	workspaceFlag := fs.String("workspace", "", "workspace directory")
	peek := fs.Bool("peek", false, "print the wake without moving the cursor of part 2")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return &usageError{err}
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return newUsageError("mate wake: want exactly 1 argument: <project>")
	}
	w, err := resolveWorkspace(*workspaceFlag)
	if err != nil {
		return err
	}
	text, err := wake(context.Background(), w, gitx.New(), fs.Arg(0), time.Now(), !*peek)
	if err != nil {
		return err
	}
	_, err = io.WriteString(stdout, text)
	return err
}

// Bounds on one wake, so a busy project still prints a page rather than a
// log. Every count printed is the true one; what is cut is named.
const (
	// wakeNewLines bounds part 2.
	wakeNewLines = 20
	// wakeNoteRunes bounds a quoted status note in parts 2 and 3. Part 1
	// quotes a question in full: it is what the Mate decides on.
	wakeNoteRunes = 120
	// wakeNextCrewLines bounds the `## For the next crew` excerpt; the
	// brief asks for at most 10.
	wakeNextCrewLines = 12
)

func wake(ctx context.Context, w *store.Workspace, git gitx.Git, project string, now time.Time, advance bool) (string, error) {
	if err := requireProject(w, project); err != nil {
		return "", err
	}
	view, err := box.Load(w, project)
	if err != nil {
		return "", err
	}
	crews, err := spawn.ListCrews(w, project)
	if err != nil {
		return "", err
	}
	cursor, err := w.ReadWakeCursor(project)
	if err != nil {
		return "", err
	}

	var out strings.Builder
	fmt.Fprintf(&out, "# mate wake %s · %s\n", project, now.Format(recallTimeLayout+" -0700"))
	out.WriteString("Contract: this is what woke you, read once. Act on part 1 in its order; do not tail status files, cat hand-backs, peek or run crew list to find what is already here. Peek a blocked Crew; read a file below only for what its excerpt does not answer.\n")

	items := wakeItems(autopilot.Gather(view, nil, now))
	fmt.Fprintf(&out, "\n== 1. To act on (%d) ==\n", len(items))
	if len(items) == 0 {
		out.WriteString("nothing open\n")
	}
	reviews := reviewsByTarget(w, project, crews)
	for _, it := range items {
		writeWakeItem(ctx, &out, w, git, project, it, reviews)
	}

	fresh, next := newStatusLines(view, cursor)
	fmt.Fprintf(&out, "\n== 2. New status since your last wake (%d) ==\n", len(fresh))
	if len(fresh) == 0 {
		out.WriteString("none\n")
	}
	shown := fresh
	if len(shown) > wakeNewLines {
		shown = shown[len(shown)-wakeNewLines:]
		fmt.Fprintf(&out, "(the %d oldest are not shown; they are in the status files under %s)\n", len(fresh)-wakeNewLines, w.CrewsDir(project))
	}
	for _, e := range shown {
		fmt.Fprintf(&out, "%s %s\n", e.Crew, clipRunes(e.Text, wakeNoteRunes))
	}

	var open []spawn.CrewSummary
	for _, c := range crews {
		if !c.Closed {
			open = append(open, c)
		}
	}
	fmt.Fprintf(&out, "\n== 3. Open crews (%d) ==\n", len(open))
	if len(open) == 0 {
		out.WriteString("none\n")
	}
	for _, c := range open {
		profile := dashIfEmpty(c.Harness)
		if c.Model != "" {
			profile += " " + c.Model
		}
		if c.Effort != "" {
			profile += "/" + c.Effort
		}
		fmt.Fprintf(&out, "%s %s · %s · %s\n", c.Crew, c.State, profile, dashIfEmpty(clipRunes(c.Note, wakeNoteRunes)))
	}

	if advance {
		if err := w.WriteWakeCursor(project, next); err != nil {
			return "", err
		}
	}
	return out.String(), nil
}

// wakeItems is a digest's items in the order a Mate works them: blocked
// first, then questions, then hand-backs, then budgets, each kind oldest
// first as Gather gave them.
func wakeItems(items []autopilot.Item) []autopilot.Item {
	rank := map[autopilot.ItemKind]int{
		autopilot.ItemBlocked: 0, autopilot.ItemNeedsDecision: 1, autopilot.ItemWaitMate: 2, autopilot.ItemBudget: 3,
	}
	out := append([]autopilot.Item(nil), items...)
	sort.SliceStable(out, func(i, j int) bool { return rank[out[i].Kind] < rank[out[j].Kind] })
	return out
}

func writeWakeItem(ctx context.Context, out *strings.Builder, w *store.Workspace, git gitx.Git, project string, it autopilot.Item, reviews map[string][]string) {
	crew := dashIfEmpty(it.Crew)
	switch it.Kind {
	case autopilot.ItemBlocked:
		kind := dashIfEmpty(string(it.IncidentKind))
		fmt.Fprintf(out, "blocked %s: %s, open for %s - %s\n", crew, kind, formatAge(it.QuietFor), it.Text)
		fmt.Fprintf(out, "  next: mate state, then mate peek, then stuck-crew-recovery\n")
	case autopilot.ItemBudget:
		// The observer's text is "<used> tokens of <budget> tokens"; the
		// half-budget one (CON-76) carries box.BudgetHalfPrefix and the
		// reason it was filed. The next lines are the issue's handling:
		// at half, one look, and a line if the Crew is looping or
		// polling; at the budget, the captain decides.
		if percent, rest := box.BudgetThreshold(it.Text); percent == 50 {
			fmt.Fprintf(out, "budget 50%% %s: %s\n", crew, rest)
			fmt.Fprintf(out, "  next: mate peek %s once; if it is looping or polling, send it one line\n", crew)
			return
		}
		fmt.Fprintf(out, "over budget %s: %s\n", crew, it.Text)
		fmt.Fprintf(out, "  next: tell the captain; the Crew keeps running until someone stops it\n")
	case autopilot.ItemNeedsDecision:
		fmt.Fprintf(out, "needs-decision %s: %s\n", crew, it.Text)
	case autopilot.ItemWaitMate:
		fmt.Fprintf(out, "wait-mate %s: %s\n", crew, it.Text)
		writeHandBack(ctx, out, w, git, project, it.Crew, reviews[it.Crew])
	}
}

// writeHandBack is what reviewing a hand-back starts from: the branch it
// left, the file it wrote with that file's `## For the next crew` section,
// and the reviews already recorded for it. A review Crew's own hand-back
// is its verdict.
func writeHandBack(ctx context.Context, out *strings.Builder, w *store.Workspace, git gitx.Git, project, crew string, reviewers []string) {
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil || len(meta) == 0 {
		return
	}
	scout := meta[spawn.MetaKind] == "scout"
	if target, ok := reviewTarget(w, project, crew); ok {
		fmt.Fprintf(out, "  review of %s: verdict %s (report %s)\n", target, reviewVerdictWord(w, project, crew), w.CrewReport(project, crew))
		return
	}
	if !scout {
		fmt.Fprintf(out, "  branch %s\n", branchSummary(ctx, w, git, project, meta))
	}
	file := w.CrewHandback(project, crew)
	if scout {
		file = w.CrewReport(project, crew)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(out, "  %s: not written\n", file)
	} else {
		lines := strings.Count(string(data), "\n")
		section := markdownSection(string(data), "For the next crew", wakeNextCrewLines)
		if section == "" {
			fmt.Fprintf(out, "  %s (%d lines), no ## For the next crew section\n", file, lines)
		} else {
			fmt.Fprintf(out, "  %s (%d lines), ## For the next crew:\n", file, lines)
			for _, l := range strings.Split(section, "\n") {
				fmt.Fprintf(out, "    %s\n", l)
			}
		}
	}
	if len(reviewers) == 0 {
		if !scout {
			out.WriteString("  reviews: none yet\n")
		}
		return
	}
	var parts []string
	for _, r := range reviewers {
		parts = append(parts, r+" "+reviewVerdictWord(w, project, r))
	}
	fmt.Fprintf(out, "  reviews: %s\n", strings.Join(parts, ", "))
}

// branchSummary is one line on a ship's branch: commits beyond the default
// branch and git's diff summary. A failure is said in the line rather than
// failing the wake: the rest of what woke the Mate still matters.
func branchSummary(ctx context.Context, w *store.Workspace, git gitx.Git, project string, meta map[string]string) string {
	branch := meta[spawn.MetaBranch]
	cfg, err := w.LoadProject(project)
	if err != nil || branch == "" {
		return dashIfEmpty(branch) + ": unknown (no branch or project config)"
	}
	repoCfg, err := cfg.CrewRepo(meta)
	if err != nil {
		return branch + ": unknown (" + oneLineMsg(err) + ")"
	}
	repo := w.RepoDir(repoCfg.Path)
	base := repoCfg.DefaultBranch
	if ok, err := git.BranchExists(ctx, repo, branch); err != nil || !ok {
		return branch + ": no longer exists"
	}
	ahead, err := git.AheadCount(ctx, repo, branch, base)
	if err != nil {
		return branch + ": unknown (" + oneLineMsg(err) + ")"
	}
	if ahead == 0 {
		return fmt.Sprintf("%s: no commits beyond %s", branch, base)
	}
	stat, err := git.Diff(ctx, repo, base, branch, true)
	if err != nil {
		return fmt.Sprintf("%s: %d commit(s) beyond %s", branch, ahead, base)
	}
	lines := strings.Split(strings.TrimSpace(stat), "\n")
	return fmt.Sprintf("%s: %d commit(s) beyond %s;%s", branch, ahead, base, " "+strings.TrimSpace(lines[len(lines)-1]))
}

// reviewsByTarget maps each reviewed Crew to its reviewers, oldest first by
// the order ListCrews gave.
func reviewsByTarget(w *store.Workspace, project string, crews []spawn.CrewSummary) map[string][]string {
	out := map[string][]string{}
	for _, c := range crews {
		if target, ok := reviewTarget(w, project, c.Crew); ok {
			out[target] = append(out[target], c.Crew)
		}
	}
	return out
}

// reviewTarget is the Crew a review Crew reviews, from the record `mate
// review` wrote.
func reviewTarget(w *store.Workspace, project, crew string) (string, bool) {
	rec, err := store.ReadMeta(w.CrewReviewMeta(project, crew))
	if err != nil || rec["crew"] == "" {
		return "", false
	}
	return rec["crew"], true
}

// reviewVerdictWord is a review's verdict as its report states it, or
// "pending" before the report has one. `mate review --check` is still the
// gate a merge takes; this only says where the review stands.
func reviewVerdictWord(w *store.Workspace, project, reviewer string) string {
	data, err := os.ReadFile(w.CrewReport(project, reviewer))
	if err != nil {
		return "pending"
	}
	if v := store.ReviewVerdict(string(data)); v != "" {
		return v
	}
	return "pending"
}

// newStatusLines is every status line past the cursor, oldest first, and
// the cursor that marks them read. A first wake, with no cursor, shows only
// each open Crew's latest line rather than every line ever written.
func newStatusLines(v box.View, cursor map[string]int64) ([]box.Entry, map[string]int64) {
	next := make(map[string]int64, len(cursor))
	for k, off := range cursor {
		next[k] = off
	}
	var fresh []box.Entry
	latest := map[string]box.Entry{}
	for _, e := range v.Entries {
		if e.Kind != box.KindStatus {
			continue
		}
		if off, ok := next[e.Ref.File]; !ok || e.Ref.Offset > off {
			next[e.Ref.File] = e.Ref.Offset
		}
		if v.Closed[e.Crew] {
			continue
		}
		if len(cursor) == 0 {
			latest[e.Crew] = e
			continue
		}
		if off, ok := cursor[e.Ref.File]; !ok || e.Ref.Offset > off {
			fresh = append(fresh, e)
		}
	}
	if len(cursor) == 0 {
		for _, e := range v.Entries {
			if l, ok := latest[e.Crew]; ok && l.Seq == e.Seq {
				fresh = append(fresh, e)
			}
		}
	}
	return fresh, next
}

// markdownSection is the body of the `## <title>` section of text, at most
// max lines, or "" when there is none. It ends at the next heading or at a
// table row.
func markdownSection(text, title string, max int) string {
	sc := bufio.NewScanner(strings.NewReader(text))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	in := false
	var lines []string
	cut := 0
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "# ") {
			if in {
				break
			}
			in = strings.TrimSpace(strings.TrimPrefix(line, "## ")) == title
			continue
		}
		if !in {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "|") {
			// A hand-back's acceptance table follows the section with no
			// heading of its own (assets/crew/brief.md.tmpl).
			break
		}
		if len(lines) >= max {
			cut++
			continue
		}
		lines = append(lines, line)
	}
	body := strings.TrimSpace(strings.Join(lines, "\n"))
	if body != "" && cut > 0 {
		body += fmt.Sprintf("\n(+%d more lines in the section)", cut)
	}
	return body
}

// clipRunes cuts s to n runes, marking the cut.
func clipRunes(s string, n int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n-1]) + "…"
}
