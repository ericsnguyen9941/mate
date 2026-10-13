# Codex and pi MCP implementation evidence

Measured 2026-10-13 in the isolated `fm-lab-mcp-20261013` Herdr session.
This is fixture evidence for a draft PR. **Real Linear login and read pilot is
pending on both harnesses.** No external service compatibility is implied by
the mock OAuth issuer.

## Versions and launch path

| Component | Version |
| --- | --- |
| Codex | codex-cli 0.162.1 |
| pi | 1.0.4 |
| Herdr | 0.9.3 |
| Go | 1.27.1 |

The live test calls real `spawn.StartMate`, resolves a workspace plugin, starts
interactive panes through Herdr with production launch flags, calls the tool,
then uses `StopMate` and `StartMate` to resume. It does not substitute print mode
or treat a server listing as a successful connection. Provider credentials are
copied to disposable homes; mock MCP login caches stay there. Browser openers
are stubbed only in the login subprocess. Operator config and caches are not
modified. No credential is included in this document.

Codex receives `-c 'mcp_servers={"probe-inst"={...}}'` on fresh and resume.
The root table replaces global MCP entries and stale fields without rewriting
`config.toml`; an empty workspace table also overrides globals. The installed
parser retains command, argv and package cwd. A Crew without servers keeps its
existing launch behavior.

After PR review, pi receives `-e <mate-extension.mjs> --skill <Mate skills>`
alongside `--no-extensions --no-skills --no-approve --offline`. The extension
uses native `registerMcpServer` with `exposure: "direct"`, including native cwd,
and calls the installed runtime's exported `createMcpExtension` with a loader
returning no file-defined servers. Global/project MCP config is ignored, with
or without workspace plugins. Native pi still owns OAuth and child lifecycle;
no third-party bridge is downloaded and the auth home is not rewritten.
The measured pi API is 1.0.4. Earlier pi versions are unsupported; the extension
checks for the native registration API, rather than claiming all versions work.

## Reproduce

Run from the implementation worktree with an explicitly owned Herdr lab:

```sh
env -u CLAUDECODE herdr --session fm-lab-mcp-20261013 server
env MATE_LIVE=1 MATE_HERDR_LIVE_SESSION=fm-lab-mcp-20261013 TMPDIR=/private/tmp \
  go test ./internal/spawn -run '^TestLivePluginMCPStartResume$' -v -timeout 30m
env MATE_LIVE=1 go test ./internal/harness/codex -run '^TestLiveCodexMCPConfigParser$' -v
make check
make build
bin/mate --version
```

pi model for this proof is `openai-codex/gpt-6.1-sol` with low effort.
The fixture implements stateless streamable HTTP with JSON responses, stdio,
and a deliberately local mock OAuth issuer (metadata, DCR, authorization
callback, token). Its token constants are mock data. This does not test SSE,
Linear consent/scopes, enterprise restrictions, or external refresh policies.

## Observed assertions

The six subtests cover Codex/pi × stdio/HTTP/mock OAuth. Every fresh and resume
turn must produce a fixture request for its unique id and an exact tool result
`probe:<harness>-<transport>-<turn>:ok` in the harness transcript. The session id
must remain identical after stop/resume. The fixture receives tool name `read_probe`; the prompt requests
`mcp__probe_inst__read_probe` and pi uses the normalized `probe_inst` namespace.

Fresh turns also request `id=error`: fixture returns MCP `isError: true` with
`probe:error:failed`; the exact text must reach the real transcript. After each
stdio stop the recorded fixture PID must cease to exist within ten seconds.
HTTP fixtures are owned and terminated by test cleanup.

For pi, each sentinel prompt yields exactly one app prompt and one Mate answer
in `sent.log`. The `MCP-RECALL-CANARY` seeded in memory appears in the transcript
on fresh/resume. This proves the hook extension is loaded alongside MCP under
normal launch flags. Skill paths/content and explicit loading flags are covered
by spawn/adapter tests; a model's use of a plugin workflow skill is not asserted
by the live fixture test.

Native mock OAuth login uses the same server name/URL and isolated home as the
pane. Codex uses `mcp login --no-browser --oauth-client-registration dcr` with the
same root override; file credential storage is selected in the lab only. pi
uses native `pi mcp login` with a temporary lab global config, removed before
registered-server launch. Both must reuse their native caches on fresh/resume.

## Original implementation run results

Before review changes, the live suite: **PASS**, 283.37 seconds, all six subtests; the final source
includes tool-error and PID cleanup assertions. Session ids below were observed
unchanged across fresh/resume:

| Harness | Transport | Fresh/resume session | Result |
| --- | --- | --- | --- |
| Codex | stdio | `01a1332c-0aa2-7023-9e80-2d64d6387a36` | PASS |
| Codex | HTTP | `01a1332c-d976-7173-ad5b-fea82e86fde4` | PASS |
| Codex | mock OAuth | `01a1332d-9984-7393-9749-9c74061733b1` | PASS |
| pi | stdio | `15677e0a-53fc-45ca-8867-5edab9182942` | PASS |
| pi | HTTP | `f20e354c-677e-43b3-bc8f-a01c145e9c37` | PASS |
| pi | mock OAuth | `5142bf04-aa8f-4451-b94e-88a08aec3784` | PASS |

Installed Codex parser test: PASS. `make check`: PASS, with 76 explicitly
expected live skips. `make build`: PASS; worktree binary reports
`mate d55cced-dirty (commit d55cced)` before committing. The live lab builds its
own binary from the same source for the pane hooks. A subsequent extension
formatting/output-byte-bound change passed the embedded Node replay in the
complete offline check; live behavior was measured before that formatting.

## Offline regressions and before/after

Before implementation, Codex refused any MCP server and pi refused the Mate
role. Positive Codex Prepare and `TestPiPreparesAMateWithHooksAndSkills` tests
were run red before removing these refusals. They now pass, alongside:

- Codex TOML roundtrip, quoted/dotted names, preserved argv/cwd, empty table,
  identical fresh/resume config, and unchanged operator config.
- pi native schema, invalid names and normalized namespace collisions, generated
  file ownership/tampering, launch intent, and hook event replay.
- Spawn materialization for Claude/Codex/pi, Crew isolation, stale skill cleanup,
  and pi preflight independent of global config.
- Codex 0.162.1 hook review capture (`Unnamed hook · new`) with own/foreign
  command matching; accepting only the exact Mate hook remains mandatory.

`make check` includes formatting, vet and the complete offline suite. Its
expected live skips are printed separately and provide no live evidence.
The explicit live fixture run is the evidence for the six combinations above.

## PR review follow-up

- Restore draft state: review found the PR had been marked ready after creation.
- Match MCP isolation across harnesses: pi now uses native `createMcpExtension`
  with an empty file-config loader and only registers workspace servers. The
  renderer is pure; Prepare/Build/CheckContext/doctor no longer read global MCP.
  No special collision preflight remains, since global entries cannot load.
- Both Codex and pi refuse MCP servers supplied to a Crew with a usage error.
- Surface hook stderr even on exit 0, preserving any useful recall stdout and
  giving a recall fallback if diagnostics have no stdout. Hook wait is capped
  at five seconds; Node replay checks stderr visibility and a hung-child timeout.
- Repoint launch config through a same-directory temporary file plus rename;
  check all JSON decode errors explicitly. Regression checks ownership, retained
  fields and absence of temporary files.
- Keep recall at `before_agent_start`. Live probing before the first prompt
  showed that pi intentionally does not create a session file until a user or
  assistant message (`SessionManager._hasConversation`). Custom entries do not
  force persistence. A positive empty-session resume assertion failed, so the
  regression now checks the correct fresh restart with the explicit
  missing-session note. No empty-session content is lost. The Go launch records
  the UUID immediately; transcript metadata arrives with the first recall hook.
- pi 1.0.4 source dispatches interactive `/compact` before `session.prompt`, and
  extension commands such as `/mcp` before emitting `input`. Slash skills and
  prompt templates do emit `input`, since they are actual model prompts. Do not
  blanket-filter every slash input or drop valid skill prompt receipts.

The follow-up live test adds exact-name/global-server marker traps to all three
pi transports and tests stop/start with a resume request before a first prompt, plus management
slash commands creating no prompt/answer receipts.

Follow-up results: the full six-transport fixture suite passed in 274.20 seconds.
After removing the ineffective custom-entry persistence attempt, the final pi
suite passed all three transports in 124.81 seconds. The final
`TestLivePiMateBeforeFirstPrompt` passed in 24.20 seconds: `/session`, `/mcp`,
`/compact` made no user/app/Mate receipt, global subprocess markers stayed absent
without workspace plugins, and the missing-session fresh-start note was explicit.

| Harness | Transport | Follow-up fresh/resume session | Result |
| --- | --- | --- | --- |
| Codex | stdio | `01a13356-79e0-7091-adbb-4ada7f647feb` | PASS |
| Codex | HTTP | `01a13357-3a4a-7020-84d7-f822a6e1cddb` | PASS |
| Codex | mock OAuth | `01a13358-0277-7b90-9b0e-b83be57d94be` | PASS |
| pi | stdio | `dc5d0d45-6b62-44f7-96a5-3287c8587be4` | PASS |
| pi | HTTP | `8c9c4842-e6d8-419e-a73c-68ffc01dd439` | PASS |
| pi | mock OAuth | `7b8ee010-9b90-4d6b-926e-82b77c52d6fc` | PASS |

Reproduce the additional proofs with:

```sh
env MATE_LIVE=1 MATE_HERDR_LIVE_SESSION=fm-lab-mcp-20261013 TMPDIR=/private/tmp \
  go test ./internal/spawn -run '^TestLivePluginMCPStartResume$/^pi$' -v -timeout 15m
env MATE_LIVE=1 MATE_HERDR_LIVE_SESSION=fm-lab-mcp-20261013 TMPDIR=/private/tmp \
  go test ./internal/spawn -run '^TestLivePiMateBeforeFirstPrompt$' -v -timeout 5m
```

Final follow-up `make check`: PASS (77 expected live skips, reported separately).
Final `make build`: PASS; worktree binary reports
`mate 226bfff-dirty (226bfff, 2026-10-13T03:46:43Z)` before the follow-up commit.
Only the owned Herdr lab is stopped/deleted after testing.

## Pending acceptance

- Real Linear OAuth login, a real read, stop/resume and another read on both
  Codex and pi. This is the user's explicit draft-PR gate.
- pi live compact/inbox/stow. Compact/reload/fork/new event mapping has offline
  replay coverage; only fresh/resume recall and prompt/answer hooks are live here.
- Live unavailable-server, timeout and cancellation behavior. The native clients
  own these paths; successful reads and stdio cleanup do not prove all failures.
- Older pi refusal before composer readiness: the extension has an API guard,
  but there is no separate binary-version probe before opening a pane.
- Installing the new binary into the operator PATH. Build is local to the
  worktree; the installed user binary is not changed by this PR.
