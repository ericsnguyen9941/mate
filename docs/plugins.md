# Plugins

A plugin is a local package directory with a `mate-plugin.yaml` manifest. It
may declare an MCP server, skills and an "open" link; every block is
optional. A workspace binds packages to projects in `.mate/plugins.yaml`. Only
a project's Mate receives them, when it starts; Crews never do. Mate has no
package manager: it downloads, locks, caches and updates nothing. Design and
rulings: [plans/plugins-and-office-2026-10-12.md](plans/plugins-and-office-2026-10-12.md).

```yaml
# /path/to/tracker-plugin/mate-plugin.yaml
apiVersion: mate.plugin/v1
name: Tracker
mcp:
  type: http                 # or stdio, with command: [argv] and optional cwd
  url: https://mcp.tracker.example/mcp
skills:
  - skills/workflow/SKILL.md # {instance_id} is replaced by the instance id
open:
  title: Open Tracker
  url: "{config.project_url}"
```

```yaml
# <workspace>/.mate/plugins.yaml
apiVersion: mate.plugins/v1
instances:
  tracker-work:              # instance id: also the MCP server's name
    path: /path/to/tracker-plugin
    projects:
      shop:
        config:
          project_url: https://tracker.example/project/shop
        key: t               # optional console shortcut on the project row
```

## Install a local package in four steps

```sh
mate plugin add /path/to/tracker-plugin --as tracker-work --project shop \
  --set project_url=https://tracker.example/project/shop [--key t]
mate plugin doctor tracker-work --project shop
mate mate stop shop && mate mate start shop
```

Then, in the Mate's pane, run `/mcp` and log in to `tracker-work` (only when
the package has an MCP server), and ask the Mate for one read through its
tools. That read, not `doctor`, is the proof the tools work.

- `add` is the captain's (a Mate or Crew is refused). It writes one explicit
  binding. `--set` is the binding's opaque config; every key the package's
  `open.url` substitutes must be set, or nothing is written. Re-running it
  for the same instance and project replaces that binding and keeps all
  others. An instance's `path` changes only by hand. The file is validated as
  a Mate start reads it before it is written. It is rewritten through a YAML
  node tree: order and most comments survive, blank lines and hand-chosen
  indentation do not.
- `doctor [<instance-id>] [--project <p>]` is offline and prints one
  `ok`/`fail` line per check: the file loads, package, manifest, each skill
  file, the stdio command (a name with `/` is resolved from the server's
  `cwd` or the package root; a bare name on `PATH`; an http server is not
  contacted), `open.url`, console key, project, and whether the Mate's
  harness (the running Mate's, else the workspace default) can take MCP
  servers. It never logs in or calls a tool. Exit 1 on any failure.
- `list` prints instances, packages and bindings; no Mate needs to run.

Editing `plugins.yaml` by hand is fine. A change reaches a Mate on its next
start, and the console or Office on reopen. `enabled: false` keeps an
instance in the file but gives nothing.
