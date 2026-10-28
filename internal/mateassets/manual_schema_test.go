package mateassets

import (
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/brief"
	"github.com/nguyenngocanh94/mate/internal/facts"
	"github.com/nguyenngocanh94/mate/internal/memory"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// The text the Mate acts on is held to the code that enforces it: the brief
// schema to internal/brief, command output to the commands that print it,
// memory lines to internal/memory. Each contract is checked in the one
// skill that owns it, by name, so a skill can be rewritten freely as long
// as what it teaches stays true.

// The skills that own each contract the tests below hold to the binary.
const (
	// briefSkill teaches the brief schema, shows a checked example, and
	// names a crew; it also shapes, splits and spawns a task.
	briefSkill = "new-task"
	// commandsSkill shows the output of project facts, remember and recall.
	commandsSkill = "mate-commands"
	// memorySkill shows the memory entry shape remember writes.
	memorySkill = "stow"
	// stowSkill quotes the stow line the app sends before a restart.
	stowSkill = "stow"
	// wakeSkill says what to do when an app event wakes the Mate.
	wakeSkill = "wake"
)

func renderedManual(t *testing.T) string {
	t.Helper()
	return manualWith(t, fixedParams())
}

func manualWith(t *testing.T, p Params) string {
	t.Helper()
	got, err := Render(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(got)
}

// skillText is the named built-in skill rendered with fixedParams.
func skillText(t *testing.T, name string) string {
	t.Helper()
	return skillWith(t, name, fixedParams())
}

func skillWith(t *testing.T, name string, p Params) string {
	t.Helper()
	got, err := RenderSkill(name, p)
	if err != nil {
		t.Fatal(err)
	}
	return string(got)
}

// everything is the manual and every built-in skill rendered with p, for
// checks that must hold wherever the Mate could read the words.
func everything(t *testing.T, p Params) map[string]string {
	t.Helper()
	out := map[string]string{"AGENTS.md": manualWith(t, p)}
	for _, name := range p.Skills() {
		out[name] = skillWith(t, name, p)
	}
	return out
}

// part is the text under the heading line heading, up to the next heading
// of the same or a higher level, so a check about one part cannot pass on
// words that only appear in another. Brief sections quoted inside code
// fences ("## Build") do not end a part.
func part(t *testing.T, text, heading string) string {
	t.Helper()
	start := strings.Index(text, "\n"+heading+"\n")
	if start < 0 {
		if !strings.HasPrefix(text, heading+"\n") {
			t.Fatalf("heading %q not found", heading)
		}
		start = -1
	}
	level := len(heading) - len(strings.TrimLeft(heading, "#"))
	rest := text[start+1+len(heading)+1:]
	fenced := false
	offset := 0
	for _, line := range strings.SplitAfter(rest, "\n") {
		if strings.HasPrefix(line, "```") {
			fenced = !fenced
		}
		if !fenced && strings.HasPrefix(line, "#") {
			l := len(line) - len(strings.TrimLeft(line, "#"))
			if l <= level && strings.HasPrefix(line[l:], " ") {
				return rest[:offset]
			}
		}
		offset += len(line)
	}
	return rest
}

// schemaItem is a list item that opens with a `## Name`.
var schemaItem = regexp.MustCompile("(?m)^- `## ([^`]+)` - ")

// TestBriefSkillSchemaMatchesTheCheck: the brief skill teaches exactly the
// sections `brief check` enforces, in the check's order, with every marker
// the check looks for.
func TestBriefSkillSchemaMatchesTheCheck(t *testing.T) {
	text := skillText(t, briefSkill)
	var got []string
	for _, m := range schemaItem.FindAllStringSubmatch(text, -1) {
		got = append(got, m[1])
	}
	want := brief.AllSections()
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("%s teaches sections %q, brief check enforces %q", briefSkill, got, want)
	}
	for _, marker := range []string{brief.OutOfScopeMarker, brief.VerifyMarker, brief.UnknownMarker, "decides: captain", "decides: mate", "`" + brief.NoneWord + "`"} {
		if !strings.Contains(text, marker) {
			t.Errorf("%s never teaches the marker %q", briefSkill, marker)
		}
	}
}

// TestBriefSkillExamplePassesTheCheck: the one example brief the Mate is
// shown is a brief `crew spawn` accepts.
func TestBriefSkillExamplePassesTheCheck(t *testing.T) {
	example := markdownBlock(t, briefSkill, skillText(t, briefSkill))
	if ps := brief.Check(example, brief.Ship); len(ps) != 0 {
		t.Fatalf("%s's example fails brief check: %v\n%s", briefSkill, ps, example)
	}
	if strings.Contains(example, "ESP32") || strings.Contains(example, "Mua ngay") {
		t.Error("the example is the buyesp32 brief; it should be a neutral one")
	}
}

// TestBriefSkillTeachesTheBriefCommands: the brief skill names the brief
// commands, the scout flag and the `# Task` the file is; the command
// reference names project facts.
func TestBriefSkillTeachesTheBriefCommands(t *testing.T) {
	text := skillText(t, briefSkill)
	for _, want := range []string{"brief check ", "brief append ", "--scout", brief.TaskHeading} {
		if !strings.Contains(text, want) {
			t.Errorf("%s does not document %q", briefSkill, want)
		}
	}
	if !strings.Contains(skillText(t, commandsSkill), "project facts <project>") {
		t.Errorf("%s does not document project facts", commandsSkill)
	}
}

// noRepoParams and twoRepoParams are fixedParams for a project with no repo
// and one with two (docs/mvp.md M9).
func noRepoParams() Params {
	p := fixedParams()
	p.Repos = nil
	return p
}

func twoRepoParams() Params {
	p := fixedParams()
	p.Repos = append(p.Repos, RepoParams{Name: "api", Path: "/ws/api", DefaultBranch: "develop"})
	return p
}

// TestFactsExampleIsTheRealOutput: the project facts example is what
// `project facts` prints, byte for byte: the first repo's block while its
// branch is empty, or the no-repo lines for a project with none.
func TestFactsExampleIsTheRealOutput(t *testing.T) {
	for name, p := range map[string]Params{"none": noRepoParams(), "one": fixedParams(), "two": twoRepoParams()} {
		t.Run(name, func(t *testing.T) {
			want := strings.Join(facts.NoRepoLines(p.ProjectName), "\n")
			if len(p.Repos) > 0 {
				r := p.Repos[0]
				want = strings.Join(facts.Facts{Project: p.ProjectName, RepoName: r.Name, Repo: r.Path, DefaultBranch: r.DefaultBranch}.Lines(), "\n")
			}
			if !strings.Contains(skillWith(t, commandsSkill, p), "```text\n"+want+"\n```") {
				t.Fatalf("%s's project facts example is not the real output:\n%s", commandsSkill, want)
			}
		})
	}
}

// identityHeading is the manual's part that holds the hard rules and the
// project's repos.
const identityHeading = "## Identity and hard rules"

// TestManualListsTheRepos: the manual lists every repo with its absolute
// path and default branch, or says there is none and what that means.
func TestManualListsTheRepos(t *testing.T) {
	sec := part(t, manualWith(t, twoRepoParams()), identityHeading)
	for _, want := range []string{
		"| Repo | Path | Default branch |",
		"| `shop` | `/ws/shop` | `main` |",
		"| `api` | `/ws/api` | `develop` |",
		"exactly one repo",
		"`/usr/local/bin/mate project repo list shop`",
	} {
		if !strings.Contains(sec, want) {
			t.Errorf("the manual with two repos does not say %q", want)
		}
	}
	none := part(t, manualWith(t, noRepoParams()), identityHeading)
	for _, want := range []string{"no repo yet", "no Crew can be spawned", "add it yourself"} {
		if !strings.Contains(none, want) {
			t.Errorf("the manual with no repo does not say %q", want)
		}
	}
	if strings.Contains(none, "| Repo |") {
		t.Error("the manual with no repo still prints a repo table")
	}
}

// TestSkillsTeachOneCrewOneRepo: the skills teach --repo, the rule it
// follows, splitting work that spans repos into several Crews, and the
// <repo>:<branch>@<sha> anchor; nothing names a repo the project lacks.
func TestSkillsTeachOneCrewOneRepo(t *testing.T) {
	p := twoRepoParams()
	for name, wants := range map[string][]string{
		commandsSkill: {"<repo>:<branch>@<sha>", "shop:main@none"},
		briefSkill:    {"--repo <repo>", "required", "one Crew per repo", "blocked-by:"},
		memorySkill:   {"<repo>:<branch>@<sha>", "shop:main@3f2a91c"},
	} {
		text := skillWith(t, name, p)
		for _, want := range wants {
			if !strings.Contains(text, want) {
				t.Errorf("%s does not say %q", name, want)
			}
		}
	}
	for file, text := range everything(t, noRepoParams()) {
		for _, gone := range []string{"shop:main@", "/ws/shop"} {
			if strings.Contains(text, gone) {
				t.Errorf("%s of a project with no repo names %q", file, gone)
			}
		}
	}
}

// headingMention is a backquoted `## Name` anywhere in prose.
var headingMention = regexp.MustCompile("`## ([^`]+)`")

// templateHeadings are the `##` headings the crew template itself asks a
// Crew to write, in a hand-back or a report, the sections the review brief
// (cmd/mate/review.go) asks a reviewer to write, and the sections of the
// Mate's own files (internal/memory); they are not brief sections. Notes is
// a review section no code reads (CON-106), so it has no store constant.
var templateHeadings = func() map[string]bool {
	m := map[string]bool{"Deviations from Build": true, "Still open": true, "Durable facts": true, "For the next crew": true,
		store.VerdictSection: true, store.FindingsSection: true, store.PriorFindingsSection: true, "Notes": true,
		store.LearningUsageSection: true, store.LearningPredictionsSection: true}
	for _, set := range [][]string{memory.Sections, memory.ProjectSections, memory.BacklogSections} {
		for _, s := range set {
			m[s] = true
		}
	}
	return m
}()

// TestEveryBriefHeadingMentionIsKnown: wherever the manual, a skill or the
// crew template names a `## ` heading, it is one the check enforces or one
// the template defines, so a typo in prose is a test failure rather than a
// Mate writing a section the check will not recognise.
func TestEveryBriefHeadingMentionIsKnown(t *testing.T) {
	known := map[string]bool{"Summary": true}
	for _, s := range brief.AllSections() {
		known[s] = true
	}
	texts := everything(t, fixedParams())
	for _, p := range []BriefParams{fixedBriefParams(), fixedScoutBriefParams()} {
		got, err := RenderBrief(p)
		if err != nil {
			t.Fatal(err)
		}
		texts["brief"] += string(got)
	}
	for file, text := range texts {
		for _, m := range headingMention.FindAllStringSubmatch(text, -1) {
			if !known[m[1]] && !templateHeadings[m[1]] {
				t.Errorf("%s mentions `## %s`, which is neither a brief section nor a template heading", file, m[1])
			}
		}
	}
}

// TestNothingParaphrasesTheCaptain: the pre-M7 rule told the Mate to write
// the task "in your own words and not by pasting the captain's message";
// M7 reverses it.
func TestNothingParaphrasesTheCaptain(t *testing.T) {
	for file, text := range everything(t, fixedParams()) {
		for _, gone := range []string{"not by pasting", "in your own words", "The template cannot supply that path"} {
			if strings.Contains(text, gone) {
				t.Errorf("%s still says %q", file, gone)
			}
		}
	}
	if !strings.Contains(skillText(t, briefSkill), "verbatim, always") {
		t.Errorf("%s does not say the captain's words are copied verbatim, always", briefSkill)
	}
}

// TestCrewsAreNamedAfterTheirTask: the captain reads a crew's id in the
// console, in crew list and in its branch, so the Mate names it after the
// task and is never taught a counter by example.
func TestCrewsAreNamedAfterTheirTask(t *testing.T) {
	const rule = "no counters such as `k3`, `p1` or `m1`"
	text := skillText(t, briefSkill)
	for _, want := range []string{"it names the task", "`fix-cart-total`", rule, "2 to 24 characters"} {
		if !strings.Contains(text, want) {
			t.Errorf("%s does not say %q", briefSkill, want)
		}
	}
	counter := regexp.MustCompile("\\b[kpm][0-9]+\\b")
	for file, text := range everything(t, fixedParams()) {
		for _, m := range counter.FindAllString(strings.Replace(text, rule, "", 1), -1) {
			t.Errorf("%s still uses the counter id %q in an example", file, m)
		}
	}
}

// TestDeviationsAreSettledBeforeReview pins CON-107: three of the eight
// changes-requested reviews on plantrip and mihuv2 (17-19/10) caught a
// choice the Crew had made and disclosed under Deviations instead of
// asking. The Mate settles a deviation before it starts a reviewer, and the
// Crew asks before it builds one or records its reading as the captain's.
func TestDeviationsAreSettledBeforeReview(t *testing.T) {
	text := skillText(t, "review-delivery")
	for _, want := range []string{"A failed step 2 is settled before anyone reviews", "A reviewer never decides a deviation", "`on`: once step 2 holds"} {
		if !strings.Contains(text, want) {
			t.Errorf("review-delivery does not say %q", want)
		}
	}
	crew, err := RenderBrief(fixedBriefParams())
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"never report it as a deviation afterwards", "quote their words as they are"} {
		if !strings.Contains(string(crew), want) {
			t.Errorf("crew brief does not say %q", want)
		}
	}
}

// TestBuilderRunsTheReviewersChecks pins CON-108: four of the eight
// changes-requested reviews on plantrip and mihuv2 (17-19/10) caught what
// the Crew would have caught had it run the reviewer's check - a re-run
// claim never re-run, a gate run only on a warm simulator, UI built to a
// reference image never compared to it, a fix that broke what the earlier
// check had caught. The brief's verify: is the reviewer's check in a clean
// environment, every hand-back claim carries its command, a fix reruns the
// check that caught it, and the finding becomes a standing check.
func TestBuilderRunsTheReviewersChecks(t *testing.T) {
	for skill, wants := range map[string][]string{
		briefSkill: {"the check an independent reviewer would run, in a clean environment", "runs from the clean state the repo provides", "compares the render to the reference, with a tolerance", "re-runs it, then `git diff --exit-code`"},
		stowSkill:  {"could have caught with a cheaper check goes under `## How to work here` as a standing check"},
	} {
		text := skillText(t, skill)
		for _, want := range wants {
			if !strings.Contains(text, want) {
				t.Errorf("%s does not say %q", skill, want)
			}
		}
	}
	crew, err := RenderBrief(fixedBriefParams())
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Every other claim in the hand-back, not only the acceptance rows, carries the command you ran and its output, or says `unverified`", "rerun the check that caught it"} {
		if !strings.Contains(string(crew), want) {
			t.Errorf("crew brief does not say %q", want)
		}
	}
}

// TestPassingReviewNotesOpenNoCorrection pins CON-109: a review whose
// Findings is none passes whatever its Notes holds, so the Mate merges it
// rather than send the crew back over a doc nit; a note worth keeping is a
// follow-up at its discretion, and --followups stays for standing Findings.
func TestPassingReviewNotesOpenNoCorrection(t *testing.T) {
	text := skillText(t, "review-delivery")
	for _, want := range []string{"`## Notes` open no correction round", "`mate backlog add` or in the tracker, as you judge", "it needs no `--followups`"} {
		if !strings.Contains(text, want) {
			t.Errorf("review-delivery does not say %q", want)
		}
	}
}

// TestWorkIsNamedAfterItsIssue pins CON-89: with a tracker, a crew id, its
// reviewer's id, its commits and its pull request all carry the issue ID, so
// the captain can follow one issue through the console, git and the tracker.
// Before it, new-task taught words only and contradicted the Linear plugin.
func TestWorkIsNamedAfterItsIssue(t *testing.T) {
	for skill, wants := range map[string][]string{
		briefSkill:        {"**With a tracker**", "`hvn-57-trip-repo`", "`<issue>-f1`", "`<issue>-r1`", "sub-issues", "`- Issue: HVN-57 (<tracker>)`", "**Without a tracker**"},
		"review-delivery": {"`<issue>-r1`", "`<id>-r1`"},
	} {
		text := skillText(t, skill)
		for _, want := range wants {
			if !strings.Contains(text, want) {
				t.Errorf("%s does not say %q", skill, want)
			}
		}
	}
	ship, err := RenderBrief(fixedBriefParams())
	if err != nil {
		t.Fatal(err)
	}
	pr, err := RenderBrief(withPR(fixedBriefParams()))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"`- Issue: <ID> (...)`", "`[HVN-57] Add trip day warnings`"} {
		if !strings.Contains(string(ship), want) {
			t.Errorf("the ship brief does not say %q", want)
		}
	}
	for _, want := range []string{"`[HVN-57] Wire the live trip repository`", "`Closes HVN-57`", "`Part of HVN-57`"} {
		if !strings.Contains(string(pr), want) {
			t.Errorf("the pull-request brief does not say %q", want)
		}
	}
	task := "## Captain's words\nx\n\n## What we already know\n- Issue: HVN-57 (Linear)\n- Unknown: y\n\n## Build\n- z\n- Out of scope: w\n\n## Acceptance\n- a. verify: `true`\n\n## Open decisions\nnone\n"
	if ps := brief.Check(task, brief.Ship); len(ps) != 0 {
		t.Errorf("a brief opening its facts with the issue line fails the check: %v", ps)
	}
}

// TestTheMateNeverWaitsForACrew pins docs/mvp.md task 57: the captain
// expects a Mate never to hold a long turn except to answer them. A live
// Mate on 2026-09-26 and again on 2026-09-27 polled its Crews inside one
// turn for nine minutes while the captain's messages sat queued, because the
// manual taught a `sleep 20; mate state` loop for manual mode.
func TestTheMateNeverWaitsForACrew(t *testing.T) {
	for file, text := range everything(t, fixedParams()) {
		if strings.Contains(text, "sleep 20") {
			t.Errorf("%s still teaches a sleep loop", file)
		}
	}
	// The rule lives in the manual, once, and in no skill.
	const rule = "**You never wait for a Crew.**"
	for file, text := range everything(t, fixedParams()) {
		if n := strings.Count(text, rule); file == "AGENTS.md" && n != 1 || file != "AGENTS.md" && n != 0 {
			t.Errorf("%s states the rule against waiting %d time(s)", file, n)
		}
	}
	text := renderedManual(t)
	for _, want := range []string{
		"after you spawn a Crew or answer one with `mate send`",
		"digests keep arriving after your turn",
		"The captain typing to you does not change the mode.",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("the manual does not say %q", want)
		}
	}
}

// TestManualGivesTheMateItsOwnProject pins the captain's boundary of
// 2026-09-27: everything inside the project is the Mate's to run, repos
// included, and only what is above it is the captain's. A live Mate had
// refused to add a repo the captain named, and handed them a shell line to
// paste, because the manual called `repo add` theirs.
func TestManualGivesTheMateItsOwnProject(t *testing.T) {
	for _, p := range []Params{twoRepoParams(), noRepoParams()} {
		sec := part(t, manualWith(t, p), identityHeading)
		for _, want := range []string{
			"Everything inside this project is yours to run",
			"creating or removing projects",
			"/usr/local/bin/mate project repo add shop <git-url|repo-path> [--name <name>]",
			"nothing is ever pushed",
		} {
			if !strings.Contains(sec, want) {
				t.Errorf("the manual does not say %q", want)
			}
		}
		if strings.Contains(sec, "theirs to run, never yours") {
			t.Error("the manual still calls repo add the captain's")
		}
	}
}

// sectionPointer is a pointer into a numbered part of a file: "section 4",
// "§13". Skills and the manual name parts by heading or skill name.
var sectionPointer = regexp.MustCompile(`(?i)\bsections? \d|§\s*\d`)

// numberedHeading is a heading that carries a number, `## 7. Spawn`: the
// mark of a skill cut from a numbered manual (CON-69).
var numberedHeading = regexp.MustCompile(`(?m)^#+ \d+\.`)

// TestSkillsStandAlone: every built-in skill is its own document, not a
// chapter cut from a manual: frontmatter, no template marker left, no
// numbered heading, and no pointer to a numbered section; the manual
// points at no section number either.
func TestSkillsStandAlone(t *testing.T) {
	for file, text := range everything(t, fixedParams()) {
		if file != "AGENTS.md" && !strings.HasPrefix(text, "---\nname: "+file+"\n") {
			t.Errorf("%s does not open with its frontmatter", file)
		}
		if strings.Contains(text, "{{") {
			t.Errorf("%s holds a template marker", file)
		}
		if m := numberedHeading.FindString(text); m != "" {
			t.Errorf("%s has the numbered heading %q", file, m)
		}
		if m := sectionPointer.FindString(text); m != "" {
			t.Errorf("%s points at %q", file, m)
		}
	}
}

// skillMention is a backquoted name in the manual's skill list.
var skillMention = regexp.MustCompile("(?m)^- `([a-z-]+)` - ")

// TestManualListsEverySkillOnce: the manual's skill list has one trigger
// line per built-in skill, and names no skill that is not installed.
func TestManualListsEverySkillOnce(t *testing.T) {
	text := part(t, renderedManual(t), "## Skills")
	var got []string
	for _, m := range skillMention.FindAllStringSubmatch(text, -1) {
		got = append(got, m[1])
	}
	if strings.Join(got, "|") != strings.Join(SkillNames, "|") {
		t.Fatalf("the manual lists skills %q, mate installs %q", got, SkillNames)
	}
}

// skillNamed is how the binary's own output names a skill: "(skill stow)".
var skillNamed = regexp.MustCompile(`skill ([a-z-]+)`)

// TestBinaryOutputNamesInstalledSkills: a skill the app names in a line it
// sends the Mate is one the Mate has.
func TestBinaryOutputNamesInstalledSkills(t *testing.T) {
	for _, line := range []string{memory.StowLine} {
		for _, m := range skillNamed.FindAllStringSubmatch(line, -1) {
			if !slices.Contains(SkillNames, m[1]) {
				t.Errorf("%q names skill %s, which mate does not install", line, m[1])
			}
		}
	}
	for _, name := range RetiredSkillNames {
		if slices.Contains(SkillNames, name) {
			t.Errorf("%s is both installed and retired", name)
		}
	}
}
