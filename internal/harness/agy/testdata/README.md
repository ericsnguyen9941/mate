# agy captures

Measured on agy 1.5.10 (Antigravity CLI) in a 120x40 Herdr 0.9.5 pane, read as visible text, on 2026-11-18 (CON-278).
The rendering mode was agy's default (inline), and the model `gemini-3.8-flash-low`.
The account, host and home directory in the screens are replaced with `operator@example.com`, `operator@host` and `/Users/operator`.

## screens

- `trust.txt` — the folder-trust dialog of a launch with no `--settings`. Build trusts the worktree, so a Crew should never see it.
- `empty.txt` — the composer after startup, holding the auto-mode hint, under the "Introducing Gemini 4 Hydrogen" banner.
- `draft-pong.txt` — `Reply with exactly: pong` typed and not sent.
- `busy-generating.txt` — that turn in flight: `⣽  Generating...` above an empty composer, footer `esc to cancel`.
- `idle-after-turn.txt` — the turn finished: `pong` in the scrollback and the composer holding `>` alone.
- `busy-editing.txt`, `busy-running-command.txt` — a turn that wrote a file and ran a shell command, with no permission prompt.
- `idle-after-tools.txt` — that turn finished.
- `draft-long.txt` — a draft that wraps onto a second row, which starts with two spaces.
- `exit-armed-draft.txt` — one ctrl+c on that draft: the draft stays and the footer reads `press ctrl+c again to exit`. The second press exited with the draft lost.
- `resume.txt` — `--conversation <id>` of that session: the earlier scrollback and the composer holding the hint again.
- `draft-slash-menu.txt` — `/quit` typed: the slash menu opens below the composer, and the footer reads `esc to cancel`, as during a turn.

## transcript

- `transcript/conversation.db` — the `steps` table of two conversations run with `agy -p --output-format stream-json` on agy 1.5.18 (2026-11-27), with every other table and column dropped, so no prompt, payload or path is kept. Steps 0–3 read a file with a tool: a response of input 15098, a `read_file` call, a response of input 2180 with 13044 cache reads. Steps 4–5 are the second conversation, renumbered: a response with effort high, output 168 of which 161 thinking. Each figure is the one stream-json printed for that call.
