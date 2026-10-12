# Example plugin packages

A package is a directory holding `mate-plugin.yaml`. It lives outside Mate's
Go build and is installed for a workspace through `<workspace>/.mate/plugins.yaml`.
See `docs/plugins.md` for the format.

## linear

An HTTP MCP server (`https://mcp.linear.app/mcp`), one skill for the Mate and
an "Open Linear" link. Verified live on 2026-10-12 (login, read, create,
update, comment); see `docs/evidence/linear-pilot-2026-10-12.md`. The server
cannot delete issues or projects; Mates cancel instead.

```
mate plugin add ./examples/plugins/linear --as linear-work --set project_url=<linear-workspace-url> --key t
mate plugin doctor linear-work
mate mate stop <project>
mate mate start <project>
```

Then, in the Mate's pane, run `/mcp`, log in to `linear-work`, and ask the Mate
for one read, for example "list my open Linear issues in this project".
`doctor` checks offline only; `/mcp` and a live read are the proof.

Run installation from the workspace, or pass `--workspace <dir>`. Restart each
running Mate that should receive it; future projects inherit it automatically.
The `project_url` config name is retained for compatibility and may point to
the Linear workspace home. The skill selects the team/project from the task.
