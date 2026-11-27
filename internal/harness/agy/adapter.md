| Fact | Value |
| --- | --- |
| Busy signature | a braille spinner, two spaces and a status ending in `...` (`⣽  Generating...`, `⣻  Running command...`, `Editing files...`) above the composer, which stays empty |
| Idle composer | a row starting `>` between two full-width rules, over the footer `? for shortcuts` and `auto · <model> · <effort>`; before the first turn the row holds the hint `Auto mode: edits & safe commands auto-approved (shift+tab to cycle)` |
| Directory trust | the launch trusts the worktree with `--settings` for that session only, so the dialog `Do you trust the contents of this project?` should not appear |
| Instructions | the brief is the always-on rule `.agents/rules/mate-brief.md` in the worktree, excluded from git; the project's `AGENTS.md` and `GEMINI.md` load too |
| Model and effort | `--model` takes an id from `agy models` (`gemini-4-hydrogen-high`); `--effort` low to max is passed as is |

The footer reads `esc to cancel` while a turn runs, and also while a slash-command menu is open, so `mate state` reads the spinner, not the footer.
Scrollback repeats `>` on the Crew's earlier prompts; only the row between the last two rules is the composer.

ctrl+c does not clear a draft: the first press shows `press ctrl+c again to exit`, and the second exits with the draft lost.
Clear a draft with ctrl+u. Steer a Crew by typing into the empty composer and pressing Enter.

agy names its conversation only once it runs. The exit prints `agy --conversation=<id>`, and mate reads the id from agy's map of the last conversation per directory.

The Crew uses the operator's agy home, because the sign-in lives there: the operator's plugins, terminal sandbox and rendering mode apply to it.
mate reads an agy Crew's tokens from its conversation database (`conversations/<id>.db`, found through that map). The database changes only when agy checkpoints it, so a running Crew's tokens arrive in steps and are complete once it exits.

A Mate cannot run on agy.
