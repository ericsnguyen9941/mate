# mate core

The core is a workspace's engine as a process of its own (CON-203): the
observer and the timeline, the outbox and the auto daemon, recovery, Mate
Office and Agent City, and the notifications. One engine runs per
workspace. `docs/architecture.md` section 2 says what the engine is made of;
this file is how to run it.

## Commands

Every command finds the workspace the way the other `mate` commands do: the
argument or `--workspace <dir>`, else `MATE_WORKSPACE`, else the nearest
directory above the current one that holds `.mate/`.

| Command | What it does |
| -- | -- |
| `mate`, `mate .`, `mate <dir>` | Runs the core in the foreground, like a dev server. Prints the block below, then the core's log. Ctrl+C stops it. |
| `mate core run [<dir>]` | The same. |
| `mate core start [<dir>]` | Runs the core in the background: in a session of its own, with no terminal, its output appended to `.mate/core.log`. Waits up to 10 s for it to answer on `.mate/core.sock`, then prints `mate core: started (pid N) · <url>`. |
| `mate core stop [<dir>]` | Sends the core SIGTERM and waits up to 60 s; a core still running then gets SIGKILL, and its socket is removed. |
| `mate core restart [<dir>]` | `stop`, then `start`. |
| `mate core status [<dir>]` | The pid, owner, version and build, start time, socket and the Office/City URLs. Exit 0 when an engine runs, 1 when none does. |
| `mate core stats [<dir>]` | What each loop of the running core has cost since it started, and its memory, GC and CPU: see "What the core costs" below. Exit 1 when none runs. |
| `mate console [<dir>]` | Opens the console, a client of the core (CON-214): see "The console" below. |

`mate` opens no browser; it prints the URLs:

```
mate core · pid 4121 · ~/work/ws
Office  http://127.0.0.1:7777/office/
City    http://127.0.0.1:7777/city/
Herdr   mate-7b09de89
Log     /Users/me/work/ws/.mate/core.log
Ctrl+C to stop
```

`Herdr` is the workspace's Herdr session (`session:` in `workspace.yaml`),
the one `herdr session attach <name>` takes. `Log` is the absolute path of
`.mate/core.log`: in the foreground the core's output goes to the terminal
and is appended to that file too, so a core's log is in the same place
however it was started. `mate core status`, and `mate` beside a running
core, print both lines as well.

A dashboard that also reaches the LAN (`dashboard.allow_remote` in
`console.yaml`) adds a `LAN` line. When the port is taken the block says
`Office/City: 127.0.0.1:7777 is in use` instead of the two URLs, and the core
runs on its socket alone. Office and City always run while the core runs;
there is no flag to turn them off.

SIGINT, SIGTERM and SIGHUP stop the core within the engine's bound (45 s,
CON-210); a second signal ends it at once.

### When an engine already runs

| What runs | `mate`, `mate <dir>` | `mate core start` | `mate core stop` |
| -- | -- | -- | -- |
| nothing | runs the core | starts the core | `mate core: not running`, exit 0 |
| the core | prints `mate core: already running` and the status, exit 0 | `mate core: already running (pid N, since HH:MM)`, exit 0 | stops it |
| a console's engine (a console built before CON-214) | `mate: the engine runs inside mate console (pid N); quit that console first`, exit 1 | the same refusal, exit 1 | refuses: quit that console |

`mate core run` beside a running core says `already running` and exits 1.

## The console

`mate console` runs no engine (CON-214). It opens in this order:

1. It takes `.mate/console.lock` (one console per workspace, CON-162). A
   second console whose captain answers `N` exits and touches no core.
2. It connects to `.mate/core.sock`. With no core answering it starts one
   as `mate core start` does and waits up to 10 s; past that it says
   `mate console: mate core did not answer within 10s; see .mate/core.log`
   and exits 1.
3. It reads `GET /api/core`. A core of another build is asked about before
   the TUI draws: `mate core (pid N) runs a different build of mate.
   Restart it? [Y/n]`. `n` exits 1 with ``mate console: the console and mate
   core must be the same build; run `mate core restart` when ready``; any
   other answer restarts the core and opens. The snapshot crosses the
   socket as JSON by Go's field names, so the two must be one build.

The console reads `/api/snapshot` and `/api/core` every 2 s and sends its
actions to `/api/act`, each with `X-Mate-Client: console`. The core counts
a console as open for 10 s after its last request: with none, a Manual
project's queued outbox lines wait; an Auto project's are delivered either
way. The status line says the core's notices (`/api/core`); the URL of
Office/City in a terminal the console cannot split comes from there too.

Quitting the console stops no core. An action still running when the
console quits is neither cancelled nor waited for: the console says
`<action> keeps running in mate core`. When the core goes away while the
console is open, the header says `stale` and the console tries again every
2 s; it never starts the core again itself, so it does not race a `mate
core stop` that meant it. `.mate/console.yaml`'s `dashboard` and
`.mate/.env` (Jev) are read by the core when it starts: change them, then
`mate core restart`.

## What the core costs

On 2026-11-11 a core was found at 1.4 GB RSS and ~70% of a CPU after days
of it: its timeline ingest took 4.9 s of every 5 s poll, and nothing said
so (epic CON-253). `mate core stats` asks the running core on its socket
(`GET /api/core/stats`, CON-263):

```
mate core · pid 8624 · ~/work/ws · up 35.72s
loop                   runs  last   mean   max    busy  last run
observer round         8     20ms   21ms   32ms   0%    552ms ago
timeline ingest        8     0.8ms  0.9ms  1.4ms  0%    532ms ago
outbox pass            18    0.1ms  0.2ms  0.4ms  0%    1.68s ago
autopilot tick         1     0.2ms  0.2ms  0.2ms  0%    35.71s ago
PR sweep + gc          1     0.0ms  0.0ms  0.0ms  0%    35.71s ago
snapshot               15    10ms   20ms   75ms   1%    2.05s ago
recovery               1     18ms   18ms   18ms   0%    35.71s ago
workspace.yaml reload  17    0.1ms  0.2ms  0.6ms  0%    1.71s ago
host watch             8     0.0ms  0.0ms  0.0ms  0%    702ms ago
notifications          0     -      -      -      -     never

heap        3.3 MB live · 18.5 MB mapped by Go · peak RSS 29.9 MB
allocated   6.7 MB in 68.3k objects since start · 192 KB/s
gc          3 cycles · 0.0% of the CPU
cpu         127ms · 0% of one core since start
goroutines  35
```

- Each loop's figures are wall time since the core started: a run that
  waited on git, Herdr or a lock counts the wait. `busy` is the share of the
  core's life the loop spent running; a loop near 100% never rests. `last
  run` says `running for …` while a run is in flight, and `never` for a loop
  that is wired and has not run (notifications off, a timer not yet due).
- `observer round` is the whole round, `timeline ingest` included; `autopilot
  tick` includes `PR sweep + gc`; `notifications` leaves out the snapshot
  it reads, which `snapshot` counts.
- Allocation is the process's, not split by loop: the loops share one heap
  and no counter says which allocated what. For that, profile the core
  (`MATE_PPROF`, CON-257).
- Recording costs two clock reads a run. Reading reads `runtime/metrics`
  and `getrusage(2)`, never `runtime.ReadMemStats`, which stops the world.
  `RSS` (now) is shown on Linux only; `peak RSS` everywhere.

`TestSteadyPassStaysWithinBudget` (`internal/core/budget_test.go`) keeps
the two costs that grow with a workspace's history, one steady ingest pass
and one `query.Load`, under an allocation budget over a generated workspace
of a thousand closed crews, in `make check`.

## How a command knows the core runs

Not by pid: the pid of a process that died may belong to another process by
now. The engine holds an advisory lock (flock) on `.mate/core.lock` for as
long as it runs, and the system releases it when the process ends, however
it ends. A command tries the lock without waiting: if it gets it (and lets go
at once), no engine runs; if not, one does, and the JSON in the file says
which process. A signal goes only to the pid the lock's current holder
wrote, and only while the lock is held. A `core.lock` left by a core that
died is therefore "not running": `stop` signals nobody, and `start` starts a
new core, which writes its own record over the old one.

## Files in `.mate/`

| File | What it is |
| -- | -- |
| `core.lock` | Held by the running engine. JSON: `pid`, `owner` (`core` or `console`), `started`, `version`, `binary`, `build` (sha256 of the binary when it started), `socket`, `url`, `lan_url`, `tcp_error`. The file stays after the engine ends. |
| `core.sock` | The core's unix socket (0600): Office/City's routes and the clients' API (`/api/core`, `/api/core/stats`, `/api/snapshot`, `/api/act`, CON-212). Removed when the core stops. |
| `tui.lock` | Held by `mate mate-agent tui` while a Mate's session is open in the TUI (CON-227). JSON: `project`, `pid` (`mate mate-agent tui`), `started`, `harness_pid` (Claude, while it runs). The file stays after it ends. |
| `core.log` | The output of every core: `mate core start` sends the core's output there, and a core in the foreground (`mate`, `mate core run`) appends what it prints to the terminal. Before it opens the log, either moves a log larger than 10 MB to `core.log.1`, keeping that one old log. The log of a core that is running is never rotated. |

## The headless Mates' hosts

A headless Mate's host is a process of its own (CON-225), not the core's:
stopping or restarting the core sends it nothing, and a turn in it runs on.
A core that starts finds each host alive on its socket and leaves it be.

While it runs, the core watches every host that `mate.meta` records
(`pid`, no pane) every 5 seconds (CON-226). A host found dead twice in a row,
under the same `pid`, is started again on its session through recovery:
under `.mate/recover.lock` first, then the action lock of its Project, and
only while `mate.meta` still records that `pid`. `core.log` gets a line:

```
mate core: Mate shop: host pid 4242 died; started again (1 of 3 in 10m), resumed its session
```

| Case | What the core does |
| -- | -- |
| The captain stopped the Mate | Nothing: the stop drops `agent`, `pid` and `socket` from `mate.meta`. |
| The host holds its host lock but does not answer | Nothing: it is alive and wedged; stop it. |
| A UI is stopping, starting or restarting that Mate | Nothing this time: that action decides. |
| A Mate in a Herdr pane (`mate.mode: pane`) | Nothing: Herdr reports it. |
| A Manual project | The same as an Auto one: starting the host types nothing into the Mate. |
| A fourth death within 10 minutes | Starts nothing more. A notification, unless `MATE_NOTIFICATIONS=off`, and the notice `Mate shop: host died 4 times in 10m; not restarting` on the console's status line, on the Mate's row and in `mate core status` (when the engine is the core's; a console's engine serves no socket to ask). |
| A host the core did not start (a start by hand, a UI's restart, recovery) | A new count, even before the core gave up. |

The count lives in the core's memory: a restarted core counts from zero, and
its recovery brings back a host that died meanwhile. `mate mate-agent stop` and
`mate mate-agent start` from a shell do not take the core's action lock: a stop
that lands in the instant the core restarts the same host can leave the new
host running, or have it started once more; stop again.

## A Mate's TUI

A headless Mate's session can be opened in Claude Code's own TUI for a
while (CON-227), for slash commands, `/mcp`, plan mode or a long back and
forth. Starting a headless Mate opens a shell pane beside it in the
workspace's Herdr session, `mate-<project>-shell`, cwd the Mate's
directory; `./tui` there runs `mate mate-agent tui <project>`, which works from
any terminal too and needs no core:

1. It takes `.mate/tui.lock`, the workspace's TUI flag. One TUI per
   workspace: a second one, of any project, says `Mate shop đang mở TUI
   (pid 4242, từ 09:30)` and exits 1.
2. It waits for the Mate's turn in flight, saying so, and holds the
   Mate's turn lock. The host stays alive and runs no turn.
3. It runs `claude` with the Mate's launch, `--resume` on the session
   `mate.meta` records, in this terminal.
4. When Claude ends - `/exit`, ^C, a crash, a kill - the Mate's hooks
   have recorded the session it ended in; when that is another, the host
   is moved to it. Then both locks go, the host runs what waited, and the
   outbox delivers.

While the flag is up:

| Who | What it does |
| -- | -- |
| Office, City | Every button but Peek is off: "Mate shop đang mở TUI". A `POST /api/actions` is a 409 with `code: tui`. Reading, Peek and the history work |
| The console | Opening, or when the flag goes up while it is open: `Mate shop đang mở TUI (pid 4242, từ 09:30). Tắt TUI? [y/N]`. `y` has the core send SIGTERM to `mate mate-agent tui` and wait up to 10 s for the flag to go, then the console runs; `N` or Enter says `close the TUI of shop first` and exits 1 |
| `POST /api/act` on `core.sock` | Refused with code `tui`, except `close_tui` and `peek` |
| That Mate's REPL | A message is refused: "Mate đang mở TUI". An app line (a digest, `mate send`) waits in the outbox as for a busy Mate, and goes once the TUI closes |
| That Mate's host | Runs no turn; alive, so the watch below restarts nothing |
| The engine's loops | Run on: the observer, the PR sweep, the auto daemon of the other Projects, the watch over the hosts |

The `pid` the UIs show is Claude's once it runs: `kill -9` of it ends the
TUI, and `mate mate-agent tui` lets the flag go and the host run on.
`mate mate-agent tui` killed outright lets its lock go, but Claude may run on
and write the session: the flag stays up while the `harness_pid` its
record names still runs, and a force close then signals that process.
While the TUI is open its lines wait in the outbox untried, so a TUI kept
open long opens no `wedged` incident.

The flag goes up before the wait for the turn in flight, so the UIs only
look during that wait too.

## After installing a new mate

A running core keeps running the binary it started from. `mate core status`
compares the sha256 of the file at `binary` with the `build` the core
recorded, and says

```
stale: the mate binary changed after this core started; run `mate core restart`
```

when they differ. `mate core restart` stops the old core and starts the
binary installed now. A core in the foreground is restarted with Ctrl+C and
`mate` again.

## mate reindex

`mate reindex` rebuilds `mate.db`, which the engine writes. Stop the core
first (`mate core stop`), run `mate reindex`, then start it again. While an
engine holds the timeline, `mate reindex` refuses and names that process.

## WSL2

- Keep the workspace on the Linux filesystem (under `~`), not under `/mnt/c`:
  the core serves a unix socket, and a Windows drive mounted in WSL2 may not
  hold one. A core that cannot bind its socket does not start.
- Office and City listen on `127.0.0.1` inside WSL2; Windows browsers reach
  them at the same URL while WSL's localhost forwarding is on (the default).
- WSL2 can shut the distribution down some time after its last terminal
  closes, and a core started with `mate core start` ends with it. Check with
  `mate core status` after reopening a terminal, or keep one terminal open.
