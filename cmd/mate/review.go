package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/mateassets"
	"github.com/nguyenngocanh94/mate/internal/spawn"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// Review uses a normal scout Crew with its own transcript and token ledger.
func cmdReview(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("review", flag.ContinueOnError)
	fs.SetOutput(stderr)
	workspace := fs.String("workspace", "", "workspace directory")
	id := fs.String("id", "", "new reviewer crew id")
	check := fs.String("check", "", "validate an existing review and print its summary")
	model := fs.String("model", "", "reviewer's model from crew dispatch")
	effort := fs.String("effort", "", "reviewer's reasoning effort")
	kind := fs.String("harness", "", "reviewer's harness from crew dispatch")
	budget := fs.Int64("budget", 0, "reviewer's token budget (default: the dispatch table's review_budget, if any)")
	since := fs.String("since", "", "review only the change since this earlier reviewer's commit, closing its findings")
	if err := fs.Parse(reorderArgs(fs, args)); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return newUsageError("usage: mate review <project> <crew> --id <reviewer> [--since <prior-reviewer>] [--harness <kind>] [--model <model>] [--effort <effort>] [--budget <tokens>] | --check <reviewer>")
	}
	if *budget < 0 {
		return newUsageError("mate review: --budget wants a positive number of tokens")
	}
	w, err := resolveWorkspace(*workspace)
	if err != nil {
		return err
	}
	project, crew := fs.Arg(0), fs.Arg(1)
	ctx := context.Background()
	if *check != "" {
		fingerprint, summary, err := checkReview(ctx, w, gitx.New(), project, crew, *check)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "review valid for %s\n%s\n", fingerprint["head"], summary)
		return nil
	}
	if *id == "" {
		return newUsageError("mate review: choose a new --id for the reviewer")
	}
	if err := store.ValidateCrewID(*id); err != nil {
		return err
	}
	if meta, err := w.ReadCrewMeta(project, *id); err != nil {
		return err
	} else if len(meta) > 0 {
		return fmt.Errorf("reviewer %s already exists; choose a new --id", *id)
	}
	fp, err := reviewFingerprint(ctx, w, gitx.New(), project, crew)
	if err != nil {
		return err
	}
	var scope *reviewScope
	if *since != "" {
		scope, err = reviewSince(ctx, w, gitx.New(), project, crew, *since, fp)
		if err != nil {
			return err
		}
		scope.record(fp)
	}
	e, err := harness.ParseEffort(*effort)
	if err != nil {
		return err
	}
	m, err := harness.ParseModel(*model)
	if err != nil {
		return err
	}
	var h harness.Kind
	if *kind != "" {
		h, err = harnesses.Parse(*kind)
		if err != nil {
			return err
		}
	}
	// The change is measured before the table is read: a small one takes
	// the table's review_small profile (CON-74). The table is read before
	// anything is written, so a malformed one leaves no review record
	// behind.
	stat, err := reviewDiffStat(ctx, w, gitx.New(), project, crew, fp, scope)
	if err != nil {
		return err
	}
	profile, note, err := applyReviewDispatch(w, spawn.SpawnCrewRequest{Harness: h, Model: m, Effort: e, Budget: *budget}, stat)
	if err != nil {
		return err
	}
	h, m, e = profile.Harness, profile.Model, profile.Effort
	if note != "" {
		fmt.Fprintf(stderr, "note: %s\n", note)
	}
	path, err := reviewRecordPath(w, project, *id)
	if err != nil {
		return err
	}
	fp[store.MetaReviewContract] = store.ReviewContract
	if err := store.WriteMeta(path, fp); err != nil {
		return err
	}
	task := "Review " + crew + " at " + fp["head"]
	if scope != nil {
		task += " since " + scope.Since
	}
	req := spawn.SpawnCrewRequest{Project: project, Crew: *id, Repo: fp["repo"], Harness: h, Model: m, Effort: e, Budget: profile.Budget, Scout: true, Task: task, BriefText: reviewBrief(w, project, crew, fp, scope, stat)}
	result, err := spawn.SpawnCrew(ctx, w, liveDeps(w, stderr), req)
	if err != nil {
		return err
	}
	// The task knows its reviewers (CON-75): `mate usage` and the dashboard
	// roll the reviewer's tokens into the task from this key. The reviewer
	// is already running, so a failure here is a note, not a refusal.
	if err := w.AddCrewReview(project, crew, *id); err != nil {
		fmt.Fprintf(stderr, "note: %s's meta does not record review %s: %v\n", crew, *id, err)
	}
	if scope != nil {
		fmt.Fprintf(stdout, "review started: %s for %s at %s, scoped to the change since %s (%s, %d finding(s))\n", result.Crew, crew, fp["head"], scope.Since, scope.Reviewed, len(scope.Findings))
	} else {
		fmt.Fprintf(stdout, "review started: %s for %s at %s\n", result.Crew, crew, fp["head"])
	}
	fmt.Fprintln(stdout, "turn: end your turn; the review arrives through the crew digest")
	return nil
}

// reviewScope is what `--since` adds to a review: the earlier review it
// continues, the commit that review saw and the findings it left standing.
// The record keeps all three so `--check` can hold the new review to them;
// base and head are still the new commit's, so `mate merge --review` pins
// the commit this reviewer saw, not the earlier one.
type reviewScope struct {
	Since    string
	Reviewed string
	Findings []string
	// PriorReport fingerprints the earlier report the findings came from,
	// so a report edited after the scoped review started makes it stale.
	PriorReport string
}

const (
	recordSince       = "since"
	recordReviewed    = "reviewed"
	recordPriorReport = "prior_report"
)

func (s *reviewScope) record(fp map[string]string) {
	fp[recordSince] = s.Since
	fp[recordReviewed] = s.Reviewed
	fp[recordPriorReport] = s.PriorReport
}

// reviewSince reads the earlier review's record and report and refuses
// every case in which "the change since" is not a well-defined, smaller
// target: no record, a record for another crew, a report without a
// verdict, a head the earlier review already saw, or a head that does not
// descend from it (the branch was rewritten, so the diff would not be
// "since" anything).
func reviewSince(ctx context.Context, w *store.Workspace, g gitx.Git, project, crew, prior string, fp map[string]string) (*reviewScope, error) {
	path, err := reviewRecordPath(w, project, prior)
	if err != nil {
		return nil, err
	}
	rec, err := store.ReadMeta(path)
	if err != nil {
		return nil, err
	}
	if len(rec) == 0 {
		return nil, fmt.Errorf("review --since: %s has no review record in %s; name a reviewer started by mate review", prior, project)
	}
	if rec["crew"] != crew {
		return nil, fmt.Errorf("review --since: %s reviewed crew %s, not %s", prior, rec["crew"], crew)
	}
	if rec["head"] == "" {
		return nil, fmt.Errorf("review --since: record of %s names no head commit", prior)
	}
	data, err := os.ReadFile(w.CrewReport(project, prior))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("review --since: %s has no report yet; wait for its hand-back or review in full", prior)
		}
		return nil, err
	}
	if store.ReviewVerdict(string(data)) == "" {
		return nil, fmt.Errorf("review --since: report of %s has no ## Verdict; wait for its hand-back or review in full", prior)
	}
	if rec["head"] == fp["head"] {
		return nil, fmt.Errorf("review --since: nothing to review; %s is still at %s, the commit %s saw. Use mate review %s %s --check %s", crew, fp["head"], prior, project, crew, prior)
	}
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil {
		return nil, err
	}
	cfg, err := w.LoadProject(project)
	if err != nil {
		return nil, err
	}
	repo, err := cfg.CrewRepo(meta)
	if err != nil {
		return nil, err
	}
	ok, err := g.IsAncestor(ctx, w.RepoDir(repo.Path), rec["head"], fp["head"])
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("review --since: %s at %s does not descend from %s, the commit %s saw; the branch was rewritten, review it in full", crew, fp["head"], rec["head"], prior)
	}
	return &reviewScope{
		Since:       prior,
		Reviewed:    rec["head"],
		Findings:    store.ReviewFindings(string(data)),
		PriorReport: fmt.Sprintf("%x", sha256.Sum256(data)),
	}, nil
}

// reviewDiffStat measures the change the reviewer will read: the crew's
// branch against its base, or against the commit the prior review saw for
// a --since review, since that is the diff the scoped brief asks for.
func reviewDiffStat(ctx context.Context, w *store.Workspace, g gitx.Git, project, crew string, fp map[string]string, scope *reviewScope) (*reviewStat, error) {
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil {
		return nil, err
	}
	cfg, err := w.LoadProject(project)
	if err != nil {
		return nil, err
	}
	repo, err := cfg.CrewRepo(meta)
	if err != nil {
		return nil, err
	}
	base := fp["base"]
	if scope != nil {
		base = scope.Reviewed
	}
	entries, err := g.DiffNumstat(ctx, w.RepoDir(repo.Path), base, fp["head"])
	if err != nil {
		return nil, err
	}
	return newReviewStat(entries), nil
}

func reviewRecordPath(w *store.Workspace, project, reviewer string) (string, error) {
	if err := requireProject(w, project); err != nil {
		return "", err
	}
	if err := store.ValidateCrewID(reviewer); err != nil {
		return "", err
	}
	return w.Resolve(w.CrewReviewMeta(project, reviewer))
}

func reviewFingerprint(ctx context.Context, w *store.Workspace, g gitx.Git, project, crew string) (map[string]string, error) {
	if err := requireProject(w, project); err != nil {
		return nil, err
	}
	if err := store.ValidateCrewID(crew); err != nil {
		return nil, err
	}
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil {
		return nil, err
	}
	cfg, err := w.LoadProject(project)
	if err != nil {
		return nil, err
	}
	repo, err := cfg.CrewRepo(meta)
	if err != nil {
		return nil, err
	}
	if meta[spawn.MetaBranch] == "" {
		return nil, fmt.Errorf("review: crew %s has no branch", crew)
	}
	if rel := meta[spawn.MetaWorktree]; rel != "" {
		wt := filepath.Join(w.Root(), filepath.FromSlash(rel))
		if _, err := os.Stat(wt); err == nil {
			dirty, err := g.IsDirty(ctx, wt)
			if err != nil {
				return nil, err
			}
			if dirty > 0 {
				return nil, fmt.Errorf("review: %s has uncommitted changes", crew)
			}
		}
	}
	head, err := g.HeadCommit(ctx, w.RepoDir(repo.Path), meta[spawn.MetaBranch])
	if err != nil {
		return nil, err
	}
	// The base is where the branch grew from, not wherever the local
	// default branch happens to point (CON-91): a crew rebased onto an
	// origin/main the local main never pulled would otherwise be reviewed
	// with every pull request merged in between.
	base, err := g.BranchBase(ctx, w.RepoDir(repo.Path), repo.DefaultBranch, head)
	if err != nil {
		return nil, err
	}
	out := map[string]string{"crew": crew, "repo": repo.Name, "head": head, "base": base}
	for key, path := range map[string]string{"brief": w.CrewBrief(project, crew), "handback": filepath.Join(filepath.Dir(w.CrewBrief(project, crew)), "handback.md")} {
		path, err = w.Resolve(path)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		out[key] = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	return out, nil
}

func reviewBrief(w *store.Workspace, project, crew string, fp map[string]string, scope *reviewScope, stat *reviewStat) string {
	p := mateassets.ReviewParams{Crew: crew, Head: fp["head"], Base: fp["base"], BriefPath: w.CrewBrief(project, crew)}
	p.HandbackPath = filepath.Join(filepath.Dir(p.BriefPath), "handback.md")
	from := p.Base
	if scope != nil {
		from = scope.Reviewed
		p.Prior = &mateassets.PriorReview{ID: scope.Since, Head: scope.Reviewed, ReportPath: w.CrewReport(project, scope.Since), Findings: scope.Findings}
	}
	p.BaseDescription = reviewBaseLines(p.Base, p.Head, from, stat)
	return mateassets.RenderReviewBrief(p)
}

func checkReview(ctx context.Context, w *store.Workspace, g gitx.Git, project, crew, reviewer string) (map[string]string, string, error) {
	path, err := reviewRecordPath(w, project, reviewer)
	if err != nil {
		return nil, "", err
	}
	recorded, err := store.ReadMeta(path)
	if err != nil {
		return nil, "", err
	}
	current, err := reviewFingerprint(ctx, w, g, project, crew)
	if err != nil {
		return nil, "", err
	}
	// A base the reviewed head does not descend from is not where the
	// branch grew from, so the diff the reviewer read was not the crew's
	// own (CON-91). Checked before staleness so the refusal says so.
	if err := checkReviewBase(ctx, w, g, project, crew, reviewer, recorded); err != nil {
		return nil, "", err
	}
	for key, value := range current {
		if recorded[key] != value {
			return nil, "", fmt.Errorf("review %s is stale: %s changed; request a new review", reviewer, key)
		}
	}
	crews, err := spawn.ListCrews(w, project)
	if err != nil {
		return nil, "", err
	}
	ready := false
	for _, c := range crews {
		if c.Crew == reviewer && (c.State == "wait-mate" || c.State == "finished" || c.State == spawn.CrewStateReviewed) {
			ready = true
		}
	}
	if !ready {
		return nil, "", fmt.Errorf("review %s has not handed back", reviewer)
	}
	data, err := os.ReadFile(w.CrewReport(project, reviewer))
	if err != nil {
		return nil, "", err
	}
	if store.ReviewVerdict(string(data)) != "pass" {
		return nil, "", fmt.Errorf("review %s did not pass; read %s", reviewer, w.CrewReport(project, reviewer))
	}
	if contract := recorded[store.MetaReviewContract]; contract != "" {
		if contract != store.ReviewContract {
			return nil, "", fmt.Errorf("review %s uses unsupported contract %q", reviewer, contract)
		}
		if err := store.ValidateReviewPass(string(data)); err != nil {
			return nil, "", fmt.Errorf("review %s is incomplete: %w", reviewer, err)
		}
	}
	if recorded[recordSince] != "" {
		if err := checkScopedReview(w, project, reviewer, recorded, string(data)); err != nil {
			return nil, "", err
		}
	}
	summary, err := reportSummary(string(data))
	if err != nil {
		return nil, "", err
	}
	return current, summary, nil
}

// checkScopedReview holds a `--since` review to the review it continues:
// that review still has its verdict and the report the checklist came from,
// and the new report marks every one of its findings closed. Anything less
// is incomplete, named by the finding that is still open.
func checkScopedReview(w *store.Workspace, project, reviewer string, recorded map[string]string, report string) error {
	prior := recorded[recordSince]
	data, err := os.ReadFile(w.CrewReport(project, prior))
	if err != nil {
		return fmt.Errorf("review %s is incomplete: the report of %s it continues cannot be read: %w", reviewer, prior, err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(data)) != recorded[recordPriorReport] {
		return fmt.Errorf("review %s is stale: the report of %s changed after the review started; request a new review", reviewer, prior)
	}
	if store.ReviewVerdict(string(data)) == "" {
		return fmt.Errorf("review %s is incomplete: review %s it continues has no verdict", reviewer, prior)
	}
	findings := store.ReviewFindings(string(data))
	if recorded[store.MetaReviewContract] == store.ReviewContract {
		if err := store.ValidatePriorFindingClosures(report, len(findings)); err != nil {
			return fmt.Errorf("review %s is incomplete: %w", reviewer, err)
		}
	}
	marks := store.ReviewPriorFindings(report)
	for i, f := range findings {
		n := i + 1
		switch marks[n] {
		case store.FindingClosed:
		case store.FindingOpen:
			return fmt.Errorf("review %s is incomplete: prior finding F%d (%s) is still open", reviewer, n, findingExcerpt(f))
		default:
			return fmt.Errorf("review %s is incomplete: prior finding F%d (%s) is not marked closed under ## %s", reviewer, n, findingExcerpt(f), store.PriorFindingsSection)
		}
	}
	return nil
}

func findingExcerpt(f string) string {
	f = strings.Join(strings.Fields(f), " ")
	if r := []rune(f); len(r) > 80 {
		return string(r[:77]) + "..."
	}
	return f
}

// reviewBaseLines tells the reviewer what the base is and how big the diff
// it is asked to read really is (CON-91), so a reviewer never has to work
// out for itself which part of a diff is somebody else's merged work. from
// is the commit the diff starts at: the base for a full review, the
// reviewed commit for a --since one.
func reviewBaseLines(base, head, from string, stat *reviewStat) string {
	line := fmt.Sprintf("- Base %s is the merge base of the crew's branch with the default branch (and origin's copy of it, when the repository has one): the commit the branch grew from, not the default branch's current head. Commits other work landed there since are not in this review.", base)
	if stat != nil {
		line += fmt.Sprintf("\n- git diff %s...%s is %d file(s), %d changed line(s). That is the whole target; if what you see is bigger, stop and report it instead of reviewing it.", from, head, stat.Files, stat.Lines)
	}
	return line
}

// checkReviewBase refuses a review record whose base is not an ancestor of
// its head. Such a base was read off a default branch the crew's branch
// did not grow from, which is how a record written before CON-91 could
// send a reviewer a diff full of other people's merged work.
func checkReviewBase(ctx context.Context, w *store.Workspace, g gitx.Git, project, crew, reviewer string, recorded map[string]string) error {
	base, head := recorded["base"], recorded["head"]
	if base == "" || head == "" {
		return fmt.Errorf("review %s is invalid: its record names no base or head; request a new review", reviewer)
	}
	meta, err := w.ReadCrewMeta(project, crew)
	if err != nil {
		return err
	}
	cfg, err := w.LoadProject(project)
	if err != nil {
		return err
	}
	repo, err := cfg.CrewRepo(meta)
	if err != nil {
		return err
	}
	ok, err := g.IsAncestor(ctx, w.RepoDir(repo.Path), base, head)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("review %s is invalid: base %s is not an ancestor of head %s, so the reviewer did not read the crew's own diff; request a new review", reviewer, base, head)
	}
	return nil
}
