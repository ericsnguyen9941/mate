---
name: {instance_id}-workflow
description: Use when the captain assigns, mentions or asks about a Linear issue, or when a Crew starts, blocks or finishes work that has a Linear issue ID.
---

# Linear workflow

This workspace makes Linear available to every Mate. You dispatch work only
for your Mate project and the scope the captain assigns.
Crew never pick up issues on their own.

## Tools

The Linear MCP server is registered as `{instance_id}`, so its tools appear as
`mcp__{instance_id}__<tool>`. Linear's official documentation
(https://linear.app/docs/mcp, read 2026-10-12) says the server has tools for
finding, creating and updating issues, projects and comments, but does not list
tool names. Run `/mcp` and read the real names before the first call; do not
guess them.

| Purpose | Candidate name | Verified |
| --- | --- | --- |
| Find or list issues, filter by assignee, state, label | `list_issues` | no (third-party catalogs 2026-10-12, not Linear docs) |
| Read one issue with comments and relations | `get_issue` | no (same) |
| Create or update an issue (state, assignee) | `save_issue` | no (same) |
| Create or update a comment on an issue | `save_comment` | no (same) |
| Find a project by name | `list_projects` | no (same) |

If a name differs from this table, use the name `/mcp` shows and tell the
captain so the table can be fixed.

## Finding the project

Use the issue ID, the captain's instructions or PROJECT.md to identify the
Linear team/project for the task. Resolve it with the available search tools.
If the scope is ambiguous, ask the captain before selecting or changing issues.
The shared plugin's `project_url` is an app-opening link; it may point to the
workspace home and does not select a task scope. Work only in the assigned
scope, including when looking for `mate-ready` issues.

## Rules

- Handle only issues the captain assigned to you or that carry the label
  `mate-ready`. Do not take an issue that already has an owner or a running
  Crew.
- Before spawning a Crew, read the issue, its comments and its blockers.
- Put the issue ID and the acceptance criteria in the Crew brief.
- Comment on the issue when work starts, when it is blocked and when it is
  handed over. Each comment carries evidence or the PR link.
- Crew state is independent of the Linear status. Do not mirror one onto the
  other.
- When a Crew finishes, review it and run the project's delivery process. Move
  the issue to Done only after acceptance and merge under the project's policy.
- Linear does not lock issues for you. One dispatcher is the rule that prevents
  double work.
- If a Linear call fails, note the issue ID and tell the captain the sync is
  unconfirmed. A failed write is not a success. A comment payload might look
  like {"issueId": "<ID>", "body": "Started, Crew spawned"}, but that shape
  is a guess, not read from Linear: take the real parameters from the tool's
  schema in `/mcp`.
