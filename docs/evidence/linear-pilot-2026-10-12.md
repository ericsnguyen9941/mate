# Linear plugin pilot — CHUA CHAY

Template. Captain fills it in after running the pilot from
`examples/plugins/README.md`. Nothing below has been observed yet.

- Package: `examples/plugins/linear`
- Tool names in the skill come from third-party catalogs; Linear's docs
  (https://linear.app/docs/mcp, 2026-10-12) list none. All are "verified: no".

| Field | Value |
| --- | --- |
| Date | |
| Claude version (`claude --version`) | |
| Instance id and project | |
| `mate plugin doctor` output | |
| `/mcp` status for the instance | |
| Setup time (add to first read) | |

## Gates

Each is yes / no with what was seen. Only a stdio server through `claude -p`
has been measured (`docs/evidence/claude-mcp-probe-2026-10-12.md`); these are
the first observations of an http server, its OAuth login and the
interactive pane.

| Gate | yes / no | Seen |
| --- | --- | --- |
| `/mcp` in the Mate's pane shows the http server `<instance-id>` connected | | |
| Login survives `mate mate stop <project>` then `mate mate start <project>` | | |
| Login survives a resumed session (`--resume`) | | |
| The skill `<instance-id>-workflow` is listed in the Mate's skills | | |

## Real tool list from `/mcp`

Compare with the skill table (`list_issues`, `get_issue`, `save_issue`,
`save_comment`, `list_projects`).

| Tool name seen | In skill table? | Notes |
| --- | --- | --- |
| | | |

## One read

- Prompt given to the Mate:
- Tool called:
- Result (no secrets):

## One write (captain chooses)

- Issue ID:
- Action (comment or status change):
- Tool called:
- Visible in Linear? (yes/no):

## Errors encountered

## Pilot measures (plan section 8)

From `docs/plans/plugins-and-office-2026-10-12.md` section 8, over the
pilot's 10-20 real issues.

| Measure | Value | Notes |
| --- | --- | --- |
| Auth or connection failures (count, each with date and cause) | | |
| Stale or duplicate tasks (count, issue IDs) | | |
| Request to merge, per issue (time) | | |
| Times the captain had to open Linear because the Office link was not enough (each: what was missing) | | |
| Did the captain follow the backlog more easily? (yes / no, why) | | |

## Verdict

- Skill table needs changes: yes / no
- Notes:
