# Example plugin packages

A package is a directory holding `mate-plugin.yaml`. It lives outside Mate's
Go build and is bound to a project through `<workspace>/.mate/plugins.yaml`.
See `docs/plugins.md` for the format.

## linear

An HTTP MCP server (`https://mcp.linear.app/mcp`), one skill for the Mate and
an "Open Linear" link. Not yet verified live; see
`docs/evidence/linear-pilot-2026-10-12.md`.

```
mate plugin add ./examples/plugins/linear --as linear-work --project <project> --set project_url=<linear-project-url> --key t
mate plugin doctor linear-work --project <project>
mate mate stop <project>
mate mate start <project>
```

Then, in the Mate's pane, run `/mcp`, log in to `linear-work`, and ask the Mate
for one read, for example "list my open Linear issues in this project".
`doctor` checks offline only; `/mcp` and a live read are the proof.
