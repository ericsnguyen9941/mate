# Paired review pilot

CON-134 measures reviewer behaviour independently of the new report schema.
See [the preregistered protocol](../../docs/research/pstack-review-evaluation-method-2026-10-29.md).

A is the exact full/scoped `reviewBrief` output from mate main
`9e6a6836ac6c48ad1f551f45dd63810126a947c8`, exported before edits. Only workspace
paths, SHAs and the prior finding are substituted in `testdata/baseline-*.md`.
B calls `mateassets.RenderReviewBrief`, the same renderer as `mate review`.
Both get the same small standalone harness wrapper; this measures the review
brief through `codex exec`, not Herdr lifecycle, budget stopping or handoff.

```sh
# Check fixtures and their private, discriminating probes without a model.
python3 scripts/revieweval/run.py --out /tmp/review-prepare --prepare-only

# Six pairs, twelve fresh sessions. Output directory must not already exist.
python3 scripts/revieweval/run.py --out /tmp/review-pilot \
  --model gpt-6-luna --effort medium --timeout 300
```

Requires Go, an authenticated Codex CLI and macOS `sandbox-exec`. The pilot does
not install a binary, alter user repositories or merge. It creates temporary
repositories, pins target SHAs, verifies their visible tests pass, and checks
private oracles fail at each defective target and pass before the regression.
The fixture corpus, gold, command order and rubric are frozen before model runs.

The candidate's OS sandbox denies reading the protocol, source, output directory,
or another candidate's checkout. Metadata reads stay permitted because Node
requires `fstat` on inherited stdout/stderr handles. A denied-read assertion checks
the manifest is unreadable before each run. macOS rejects nested sandboxes, so the external profile is the sole sandbox:
it denies writes outside each candidate lab (plus /dev/null and /dev/tty), and
Codex uses its documented externally-sandboxed mode. Network is enabled equally
for both variants so local HTTP unit tests can bind. TMPDIR and zsh TMPPREFIX
point inside the candidate scratch directory so heredoc probes can run. Each run has a fresh
CODEX_HOME with a protected auth-file copy, removed after the run; session
rollouts are retained for evidence. The staging credential variables are
removed. Configuration, rules and project instruction discovery are disabled
equally for both variants.

Order: S1 A/B, S2 B/A, S3 A/B, S4 B/A, S5 A/B, S6 B/A. Each session has the same
300-second wall limit and no token cap; exact tokens are measured from usage
receipts. Reports, commands, usage, timings, target SHA, final diff and probe
sources are retained. The caller must archive results before cleaning up the
temporary checkouts. Never count an empty response or timeout as pass.

Score neutral report IDs against the preregistered gold before joining the A/B
mapping. Report semantic verdict, true defects, false blockers, unsupported pass,
observed/static/inherited/unverified evidence and scope separately from schema
conformance. A startup failure with no model turn is a lab failure: retain it and
rerun both members of the pair, not just the failing treatment. Preserve all runs
when changing the treatment. Six pairs are a pilot, not a statistical estimate of
reviewer quality on real projects.

Check a candidate passing report with the production validator, separately from
semantic scoring:

```sh
go run ./scripts/revieweval --check-report /path/to/report.md
# Scoped review with one prior finding:
go run ./scripts/revieweval --check-report /path/to/report.md --prior-count 1
```
