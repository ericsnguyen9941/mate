---
name: {instance_id}-workflow
description: Use when the captain assigns, mentions or asks about a Linear issue, or when a Crew starts, blocks or finishes work that has a Linear issue ID.
---

# Linear workflow

This workspace makes Linear available to every Mate. You dispatch work only
for your Mate project and the scope the captain assigns.
Crew never pick up issues on their own.

## How Linear is shaped

Linear is not Jira. Map it this way:

| Linear | Means here | Jira equivalent |
| --- | --- | --- |
| Team | One Mate project; issue IDs carry its key (`HVN-12`) | Project |
| Project | One epic: a deliverable with a goal, target date, milestones | Epic |
| Parent issue with sub-issues | A small epic, or a task split in steps | Story with sub-tasks |
| Initiative | A group of projects | Initiative |
| Cycle | A team's time box, optional | Sprint |

Your team's key is in PROJECT.md or the captain's workspace rules. Work only
inside your own team.

## Tools

The Linear MCP server is registered as `{instance_id}`, so its tools appear as
`mcp__{instance_id}__<tool>`. Names below were read from `/mcp` on 2026-10-12
(Linear's hosted server, 76 tools); that day a Mate read issues, created and
updated an issue and added a comment through them.

| Purpose | Tool |
| --- | --- |
| Find or list issues, filter by team, project, state, label | `list_issues` |
| Read one issue with its description | `get_issue` |
| Create an issue, or update one (state, project, labels, description) | `save_issue` |
| Add or edit a comment on an issue | `save_comment`, read with `list_comments` |
| Find teams, projects, states, labels | `list_teams`, `list_projects`, `list_issue_statuses`, `list_issue_labels` |
| Create or update a project (epic) or milestone | `save_project`, `save_milestone` |

Take parameters from each tool's schema, not from memory. If a call fails
because a name or parameter changed, use what `/mcp` shows and tell the
captain so this table can be fixed.

The server cannot delete issues or projects. To drop work, set the issue's
state to Canceled; if something must really be deleted, tell the captain,
who does it on the web.

## Finding the scope

Use the issue ID, the captain's instructions or PROJECT.md to identify the
team and project for the task. If the scope is ambiguous, ask the captain
before selecting or changing issues. The plugin's `project_url` is an
app-opening link; it may point to the workspace home and does not select a
task scope. Work only in the assigned scope, including when looking for
`mate-ready` issues.

## Rules

- Handle only issues the captain assigned to you or that carry the label
  `mate-ready`. Do not take an issue that already has an owner or a running
  Crew.
- Before spawning a Crew, read the issue, its comments and its blockers.
- Name the task after its issue: the Crew id (and so its branch,
  `mate/<id>`) is the issue ID in lower case, then a short slug of the
  title - `hvn-124-login-oauth` for HVN-124. A Crew id is at most 24
  characters, so keep the slug short; never drop the issue ID to make room.
  Use the same `HVN-124` prefix for the task's backlog line and the pull
  request title. Linear links a branch or PR that carries the ID to the
  issue by itself.
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
  unconfirmed. A failed write is not a success. Read the issue back after a
  write you report.
