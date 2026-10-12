# Linear plugin pilot — setup gates passed 2026-10-12

Run by Claude Code on the captain's delegation, in the captain's real
workspace. Gates below were observed; the pilot measures (section 8 of the
plan) need 10-20 real issues over time and are not filled yet.

- Package: `examples/plugins/linear`, installed at workspace scope (PR #46).

| Field | Value |
| --- | --- |
| Date | 2026-10-12 |
| Claude version (`claude --version`) | 2.1.287 (Claude Code) |
| Mate binary | `mate 5846264` |
| Instance id and projects | `linear-work`, workspace `/Users/erics/workspace`, projects hellovietnam and mihu (both receive it; no `--project`) |
| Linear workspace | `controlplane`, Business plan; teams HVN (hellovietnam), MIHU (mihu), CON (not a Mate project) |
| `mate plugin doctor` output | all `ok`: plugins.yaml valid, package, manifest, skill, `mcp http https://mcp.linear.app/mcp: ok (not contacted)`, open.url, key t, harness claude for workspace default and each running Mate |
| `/mcp` status for the instance | before login `⚠ linear-work needs authentication`; after `✔ linear-work 76 tools` |
| Setup time (add to first read) | about 15 minutes, most of it reaching the captain's real Chrome profile |

## Gates

| Gate | yes / no | Seen |
| --- | --- | --- |
| `/mcp` in the Mate's pane shows the http server `linear-work` connected | yes | `/mcp` → Authenticate opened Linear's OAuth in the default browser; the captain was already signed in to Linear, consent returned "Authentication successful. Connected to linear-work." |
| One login serves every Mate | yes | After logging in from mihu's pane, hellovietnam's `/mcp` showed `Auth: ✔ authenticated` without a second login |
| Login survives `mate mate stop <project>` then `mate mate start <project>` | yes | mihu stopped (stowed) and started; `/mcp` still `✔ linear-work 76 tools` |
| Login survives a resumed session (`--resume`) | yes | that start reported `resumed session f8569914-…`; same `/mcp` result |
| The skill `<instance-id>-workflow` is listed in the Mate's skills | yes | hellovietnam's turn showed `1 skill available`, then `Skill(plugin-linear-work-workflow) Successfully loaded skill` |

Pitfall seen: on an already-authenticated server, `/mcp` → the server's
menu puts `Disable` where `Authenticate` was. A second Enter disabled
`linear-work` in hellovietnam; `Enable` restored it with nothing persisted.

## Real tool list from `/mcp`

76 tools under `mcp__linear-work__`. All five names the skill table guessed
exist: `list_issues`, `get_issue`, `save_issue`, `save_comment`,
`list_projects`. Others used by the updated skill: `list_teams`,
`list_issue_statuses`, `list_issue_labels`, `list_comments`, `save_project`,
`save_milestone`.

No tool deletes an issue, a project or a team, and none creates a team. The
only `delete_*` tools are for attachments, comments, diff comments and
status updates. Other families: initiatives, documents, cycles, milestones,
releases, notifications, diffs/reviews, agent skills
(`list_agent_skills`, `get_agent_skill`).

## One read

- Prompt given to the Mate (mihu): list the teams and projects in workspace
  controlplane with each project's open issue count; read only.
- Tool called: `linear-work`, 4 calls.
- Result: team EricsNguyen (CON), project controlplane, 3 open issues
  (CON-12, CON-15, CON-16), 16 done; matched the web UI (19 issues, 84%).

## Writes

- Comment: hellovietnam's Mate commented on HVN-100 (E0.12, In Progress),
  then read the issue back and found the comment. Visible in Linear: yes.
- Create + update: mihu's Mate created MIHU-1 in Backlog, then moved it to
  Todo and added acceptance criteria, and read it back. Visible in Linear: yes.

## Setup done outside the MCP

The MCP cannot create teams, so a temporary personal API key (revoked after
use) created teams HVN and MIHU over GraphQL and imported hellovietnam's
Beads plan: 17 epics as Linear projects, 113 tasks as issues with state,
priority, labels, acceptance criteria, notes, comments and 127 `blocks`
relations; each issue names its Beads id. Beads stays as an archive.

## Errors encountered

- A long message typed into a Mate's composer through `herdr pane send-text`
  arrived with its head cut off; a shorter resend arrived whole.
- `chrome-devtools-axi` reused a headless bridge from an earlier session
  instead of the captain's Chrome until that bridge was stopped and remote
  debugging was enabled in `chrome://inspect/#remote-debugging`.
- Linear's pricing page says MCP access may draw AI credits; the billing page
  showed `$0.00 remaining` and every call above still succeeded.

## Pilot measures (plan section 8)

Not measured yet: they need 10-20 real issues worked through Mates.

| Measure | Value | Notes |
| --- | --- | --- |
| Auth or connection failures (count, each with date and cause) | 0 so far | |
| Stale or duplicate tasks (count, issue IDs) | | |
| Request to merge, per issue (time) | | |
| Times the captain had to open Linear because the Office link was not enough (each: what was missing) | | |
| Did the captain follow the backlog more easily? (yes / no, why) | | |

## Verdict

- Skill table needs changes: yes, done in the same change (names verified,
  no-delete rule, Linear-vs-Jira mapping).
- Notes: setup gates pass; keep the pilot open for the section 8 measures.
