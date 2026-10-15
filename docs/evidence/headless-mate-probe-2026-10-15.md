# Headless Mate probes (2026-10-15, CON-78)

Claude Code 2.1.292, Herdr 0.9.3, macOS. Load average around 1000 during the run (three hellovietnam crews building).

## Probe 1: Herdr hosts a non-harness process as an agent

Throwaway session `hlprobe`, workspace `w1`, tab `w1:t2`, pane `w1:p2`.

- `herdr agent start --kind` accepts only known harnesses, so the host cannot start through `runtime.StartAgent`.
- `herdr pane run w1:p2 /tmp/hlprobe-pane.sh` types the command into the pane's shell. The script sees `HERDR_PANE_ID=w1:p2`, `HERDR_SESSION`, and `HERDR_SOCKET_PATH`.
- The script runs `herdr pane report-agent --source mate --agent mate-probe --state idle "$HERDR_PANE_ID"`. `agent list` then shows `{"agent":"mate-probe","agent_status":"idle","pane_id":"w1:p2",…}`. It stays listed for as long as the process lives (checked 10 s and 30 s later).
- A report sent from outside the pane while the pane ran only its shell was listed once and gone 3 s later. The report must come from the process the pane runs.
- `agent get mate-probe` returns `agent_not_found`: the `--agent` label is the kind, not the name. After `herdr agent rename w1:p2 mate-hl`, `agent get mate-hl`, `agent wait w1:p2 --timeout 2000` and `agent list` all return the agent with `"name":"mate-hl"`. `ListAgents` maps `agent` to `Kind` and `name` to `Name`.
- `report-agent --state working` from outside the pane, once the pane owns the agent, moves `agent_status` to `working`.
- Killing the script removes the agent from `agent list` within 4 s.

Conclusion: the host reports itself from inside its pane and is renamed to the Mate's agent name. Liveness (`query.ReadLiveness`), `InspectAgent` and `agent attach` then work unchanged, and a dead host is not listed.

## Probe 2: Claude hooks under `claude -p`

Settings file holding the three Mate hook events (each appends its stdin payload to a log), plus the captain's `env` (proxy). Flags: `-p --setting-sources project --settings <file> --model haiku --output-format stream-json --verbose --dangerously-skip-permissions`, nested-session variables unset, prompt on stdin.

The first attempt kept `CLAUDECODE` and the other nested-session variables of the calling session (`harness.NestedSessionEnv`) and retried `401 authentication_failed` ten times. The same run with them unset authenticated.

Turn 1, `--session-id <uuid>`, prompt `⟦mate⟧ digest: test. Reply with exactly: ok1`:
- SessionStart `source=startup`
- UserPromptSubmit `prompt` = `⟦mate⟧ digest: test. Reply with exactly: ok1` (marker intact)
- Stop `last_assistant_message` = `ok1`
- stream `result`: `{"result":"ok1","session_id":"899693ea-…","is_error":false,"num_turns":1,"duration_ms":1698}`

Turn 2, `--resume <same uuid>`, prompt `Reply with exactly: ok2`:
- SessionStart `source=resume` (fires on every resumed turn)
- UserPromptSubmit, Stop `ok2`
- `result` `ok2`, same session id, `duration_ms` 1781

Stream event types seen: `system` (`hook_started`, `hook_response`, `init`, `thinking_tokens`, `api_retry`), `assistant`, `user`, `result`.

Conclusion: the existing hooks keep `sent.log` and `mate.meta` working under `-p`. SessionStart(resume) runs `mate hook mate-session` on every turn, so its output (live state) would be added to every turn unless the host suppresses it.
