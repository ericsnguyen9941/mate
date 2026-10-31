// Fake workspace for the Agent City console prototype. It answers the same
// paths the real dashboard serves (docs/dashboard.md), in the same shapes,
// plus section 7b's actions (CON-141): GET /api/projects/{p}/actions and
// POST /api/actions. Nothing leaves the page.
//
// 3 projects, 6 open crews in every state the design has to show, 2 topics
// waiting on the captain, 1 incident.

const NOW = Date.now();
const iso = (minAgo) => new Date(NOW - minAgo * 60e3).toISOString().replace(/\.\d+Z$/, 'Z');
const tok = (total) => ({ input: 0, cache_read: 0, cache_write: 0, output: 0, thinking: 0, total });

// Knobs the prototype's demo drawer turns.
export const demo = {
  stale: false, // /api/events stops answering
  readOnly: false, // a dashboard without a console behind it
  next: 'auto', // the next action: auto | refuse | fail
};

let eventID = 412;
const listeners = new Set();
export const onChange = (fn) => listeners.add(fn);
function bump(project) {
  eventID++;
  for (const fn of listeners) fn(project);
}

// ---- The workspace -----------------------------------------------------------

export const ws = {
  root: '/Users/captain/workspace',
  projects: [
    {
      name: 'plantrip', repos: ['plantrip-ios', 'plantrip-api'], mode: 'auto', inbox_waiting: 1, path: '~/workspace/plantrip',
      mate: { harness: 'claude', running: true, state: 'at_desk_working', since: iso(2), tokens_today: 3_412_000, context_pct: 41, context_tokens: 123_000 },
    },
    {
      name: 'mihuv2', repos: ['mihuv2'], mode: 'manual', inbox_waiting: 2, path: '~/workspace/mihuv2',
      mate: { harness: 'codex', running: true, state: 'idle', since: iso(9), tokens_today: 1_208_000, context_pct: 63, context_tokens: 168_000 },
    },
    {
      name: 'landing', repos: ['landing-site'], mode: 'manual', inbox_waiting: 0, path: '~/workspace/landing',
      mate: { harness: '', running: false, state: '', tokens_today: null, context_pct: null, context_tokens: null, created: false },
    },
  ],
};

const crew = (o) => ({
  role: 'ship', closed: false, harness: 'claude', model: 'opus-4.7', effort: 'high', tokens: tok(0), turns: 0,
  ...o,
});

export const tasks = {
  plantrip: [
    crew({ crew: 'PLA-61', text: 'Offline map tiles for the trip screen', branch: 'mate/PLA-61', repo: 'plantrip-ios', state: 'at_desk_working', since: iso(4), spawned_at: iso(96), tokens: tok(2_140_000), status: 'working', status_text: 'tiles render; wiring the cache eviction', dirty: true }),
    crew({ crew: 'PLA-58', text: 'Share a trip as a link', branch: 'mate/PLA-58', repo: 'plantrip-api', state: 'waiting_review', since: iso(11), spawned_at: iso(260), tokens: tok(1_870_000), status: 'wait-mate', status_text: 'PR #212 green, review pass, ready to land', review: '', reviews: ['PLA-58-r1'], model: 'sonnet-4.6', effort: 'medium' }),
    crew({ crew: 'PLA-58-r1', role: 'review', review_of: 'PLA-58', text: 'Review PLA-58', branch: 'mate/PLA-58', repo: 'plantrip-api', state: 'gone', closed: true, close_state: 'finished', review: 'pass', spawned_at: iso(70), closed_at: iso(14), harness: 'codex', model: 'gpt-5.5', effort: 'high' }),
    crew({ crew: 'PLA-63', role: 'review', review_of: 'PLA-60', text: 'Review PLA-60', branch: 'mate/PLA-60', repo: 'plantrip-api', state: 'at_desk_working', since: iso(6), spawned_at: iso(8), harness: 'codex', model: 'gpt-5.5', effort: 'high', tokens: tok(310_000), status: 'working', status_text: 'reading the migration' }),
    crew({ crew: 'PLA-55', text: 'Currency picker', branch: 'mate/PLA-55', repo: 'plantrip-ios', state: 'gone', closed: true, close_state: 'finished', spawned_at: iso(60 * 30), merged_at: iso(60 * 26), closed_at: iso(60 * 26) }),
    crew({ crew: 'PLA-52', text: 'Trip cover photos', branch: 'mate/PLA-52', repo: 'plantrip-ios', state: 'gone', closed: true, close_state: 'finished', spawned_at: iso(60 * 54), merged_at: iso(60 * 50), closed_at: iso(60 * 50) }),
    crew({ crew: 'PLA-57', text: 'Rate limit the share endpoint', branch: 'mate/PLA-57', repo: 'plantrip-api', state: 'gone', closed: true, close_state: 'finished', spawned_at: iso(60 * 40), merged_at: iso(60 * 37), closed_at: iso(60 * 37) }),
  ],
  mihuv2: [
    crew({ crew: 'HU-131', text: 'Pay with saved card', branch: 'mate/HU-131', repo: 'mihuv2', state: 'waiting_at_ceo', since: iso(7), spawned_at: iso(55), tokens: tok(880_000), status: 'needs-decision', status_text: 'Stripe or the bank SDK for saved cards? Stripe is two days, the SDK one week.', harness: 'claude', model: 'opus-4.7', effort: 'high' }),
    crew({ crew: 'HU-129', text: 'Gallery drag to reorder', branch: 'mate/HU-129', repo: 'mihuv2', state: 'blocked', since: iso(42), spawned_at: iso(130), tokens: tok(1_430_000), status: 'blocked', status_text: 'xcodebuild test hangs on the simulator lane', incident: 'stuck', harness: 'pi', model: 'kimi-k2', effort: 'medium' }),
    crew({ crew: 'HU-127', text: 'Onboarding copy pass', branch: 'mate/HU-127', repo: 'mihuv2', state: 'asleep', since: iso(95), spawned_at: iso(200), tokens: tok(240_000), status: 'idle', status_text: 'handed back, waiting for the next brief', model: 'haiku-4.5', effort: 'low' }),
    crew({ crew: 'HU-120', text: 'Learning gate', branch: 'mate/HU-120', repo: 'mihuv2', state: 'gone', closed: true, close_state: 'finished', spawned_at: iso(60 * 70), merged_at: iso(60 * 66), closed_at: iso(60 * 66) }),
  ],
  landing: [],
};

const openCrews = (p) => tasks[p].filter((t) => !t.closed);

// Status lines and questions for the tents' playback (section 4b).
function linesFor(t) {
  const start = Date.parse(t.spawned_at);
  const at = (frac) => new Date(start + frac * ((t.closed_at ? Date.parse(t.closed_at) : NOW) - start)).toISOString().replace(/\.\d+Z$/, 'Z');
  const L = (frac, verb, text, event_id) => ({ at: at(frac), verb, text, event_id });
  const script = {
    'PLA-61': [L(0.05, 'working', 'reading the trip screen and the tile API'), L(0.4, 'working', 'tiles render at zoom 12–16'), L(0.85, 'working', t.status_text)],
    'PLA-58': [L(0.05, 'working', 'drafting the share token'), L(0.5, 'working', 'link opens the trip read-only'), L(0.8, 'wait-mate', 'https://github.com/x/plantrip/pull/212 green, ready')],
    'PLA-58-r1': [L(0.2, 'working', 'reading the diff'), L(0.9, 'wait-mate', 'tokens expire, tests cover revoke')],
    'PLA-63': [L(0.3, 'working', t.status_text)],
    'PLA-55': [L(0.1, 'working', 'picker in settings'), L(0.8, 'wait-mate', 'https://github.com/x/plantrip/pull/201')],
    'PLA-52': [L(0.1, 'working', 'unsplash fallback'), L(0.8, 'wait-mate', 'https://github.com/x/plantrip/pull/196')],
    'PLA-57': [L(0.2, 'working', 'token bucket per link'), L(0.8, 'wait-mate', 'https://github.com/x/plantrip/pull/205')],
    'HU-131': [L(0.1, 'working', 'saved card screen drafted'), L(0.8, 'needs-decision', t.status_text, 9001)],
    'HU-129': [L(0.1, 'working', 'long-press drag works'), L(0.5, 'working', 'running the UI tests'), L(0.7, 'blocked', t.status_text)],
    'HU-127': [L(0.2, 'working', 'new copy for 4 screens'), L(0.7, 'wait-mate', 'handed back')],
    'HU-120': [L(0.2, 'working', 'gate reads learnings'), L(0.8, 'wait-mate', 'https://github.com/x/mihuv2/pull/104')],
  };
  const questions = t.crew === 'HU-131' ? [{ asked_event_id: 9001, text: t.status_text }] : [];
  return { status_lines: script[t.crew] ?? [], questions };
}

// Inbox: status lines that want someone, and incidents.
function inboxOf(p) {
  const out = [];
  let seq = 1;
  for (const t of openCrews(p)) {
    if (['needs-decision', 'blocked', 'wait-mate'].includes(t.status)) out.push({ seq: seq++, at: t.since, kind: 'status', crew: t.crew, verb: t.status, text: t.status_text, attention: true, assigned: t.assigned ?? false });
    if (t.incident) out.push({ seq: seq++, at: iso(2), kind: 'incident', crew: t.crew, verb: t.incident, text: 'no output for 42m; last command xcodebuild test', attention: true });
  }
  out.push({ seq: seq++, at: iso(30), kind: 'status', crew: openCrews(p)[0]?.crew ?? '', verb: 'working', text: 'progress note', attention: false });
  return out;
}

// ---- News: the topics the Mates opened for the captain ------------------------

export const topics = [
  {
    id: 't-plantrip-12', project: 'plantrip', status: 'waiting', at: iso(18),
    question: 'Ship offline maps in 2.4, or hold them for the tile CDN?',
    context: 'PLA-61 has tiles rendering from the device cache. The CDN contract starts 11/15; until then tiles come from the free tier, capped at 50k requests a day.',
    options: ['Ship in 2.4 on the free tier', 'Hold for the CDN (2.5)', 'Ship behind a flag, off by default'],
    notes: 'I lean to the flag: we learn without risking the cap.',
  },
  {
    id: 't-mihuv2-7', project: 'mihuv2', status: 'waiting', at: iso(6),
    question: 'HU-129 has hung on the simulator lane for 42 minutes. Restart it, or move it off the lane?',
    context: 'The xcode lane holds one crew at a time (workspace.yaml lanes). HU-131 is waiting for the lane too.',
    options: ['Restart HU-129 in its worktree', 'Stop HU-129, re-brief later', 'Leave it, I will look'],
    notes: '',
  },
  {
    id: 't-plantrip-11', project: 'plantrip', status: 'on_it', at: iso(80),
    question: 'Use the new review model for plantrip-api?', context: '', options: ['Yes', 'No'], notes: '', answer: 'Yes',
  },
  {
    id: 't-mihuv2-5', project: 'mihuv2', status: 'earlier', at: iso(60 * 20),
    question: 'Merge HU-120 with the gate on by default?', context: '', options: ['On', 'Off'], notes: '', answer: 'On',
  },
];

// ---- Chat with each Mate -----------------------------------------------------------

export const chats = {
  plantrip: [
    { who: 'captain', at: iso(150), text: 'Pick up the sharing epic next.' },
    { who: 'mate', at: iso(149), text: 'Spawned PLA-58 (share link) and PLA-57 (rate limit). PLA-57 merged.', tools: ['mate spawn', 'mate spawn'] },
    { who: 'captain', at: iso(40), text: 'How long until PLA-58 lands?' },
    { who: 'mate', at: iso(39), text: 'Review passed 14 minutes ago. It is waiting on merge; I left it for you since the epic is yours to close.', tools: ['mate status PLA-58'] },
  ],
  mihuv2: [
    { who: 'captain', at: iso(90), text: 'Keep HU-131 on Stripe unless the bank answers today.' },
    { who: 'mate', at: iso(88), text: 'Noted in ## Captain. HU-131 asked again 7 minutes ago; I have not answered since mode is manual.', tools: ['mate note'] },
  ],
  landing: [],
};

// ---- Board (7 days) ------------------------------------------------------------------

function boardOf(name) {
  const days = [];
  const seed = { plantrip: [3.1, 4, 3], mihuv2: [1.9, 2, 1], landing: [0, 0, 0] }[name];
  for (let d = 6; d >= 0; d--) days.push({ tokens: { total: Math.round(seed[0] * 1e6 * (0.6 + ((d * 37) % 7) / 10)) }, closed: { finished: d % 2 ? seed[1] : seed[1] - 1 }, reviews: { pass: seed[2] } });
  return { name, days };
}

// ---- Actions (section 7b) -----------------------------------------------------------

const running = new Map(); // "project" or "project/crew" → action

function b(action, label, available, reason, confirm) {
  return { action, label, available, ...(reason ? { reason } : {}), ...(confirm ? { confirm } : {}) };
}

function mateButtons(p) {
  const card = p.mate;
  const created = card.created !== false;
  const start = !created
    ? b('onboard', 'Create mate…', true)
    : card.running ? b('start', 'Start mate', false, 'Mate is recorded running') : b('resume', 'Resume mate', true);
  return [
    start,
    b('stop', 'Stop mate…', created && card.running, !created ? 'no Mate yet' : card.running ? '' : 'Mate is not running',
      `Stop the Mate of ${p.name}? This Mate only; its ${openCrews(p.name).length} crews keep running. The agent stops and the binding is released once the service confirms it.`),
    b('restart_mate', 'Restart mate…', created && card.running, !created ? 'no Mate yet' : card.running ? '' : 'Mate is not running',
      `Restart the Mate of ${p.name}? This Mate only; its crews keep running. The Mate is stopped and started again.`),
    b('mode', `Mode → ${p.mode === 'auto' ? 'manual' : 'auto'}`, created, created ? '' : 'no Mate yet'),
    b('clear_composer', 'Clear composer', created && card.running, !created ? 'no Mate yet' : card.running ? '' : 'Mate is not running'),
    b('explain', 'Explain notice', created && card.running, !created ? 'no Mate yet' : card.running ? '' : 'Mate is not running'),
  ];
}

function crewButtons(p, t) {
  const name = `${p}/${t.crew}`;
  const waiting = t.status === 'wait-mate';
  return [
    b('stop', 'Stop crew…', true, '',
      `Close ${name}? Its agent, tab, worktree and branch go; its brief, report and status file stay. Refused, changing nothing, if the branch is not landed or the worktree is dirty.`),
    b('restart_crew', 'Restart crew…', true, '',
      `Restart ${name}? A live agent is stopped first; then the harness starts again in the same worktree, and its branch and status file stay.`),
    b('merge', 'Merge…', waiting, waiting ? '' : 'not waiting on merge',
      `Merge ${t.branch} → main and finish crew ${t.crew}? Fast-forward only: refused, changing nothing, if main moved.`),
    b('diff', 'Diff', true),
  ];
}

function actionsOf(p) {
  if (demo.readOnly) return { project: p, read_only: 'this dashboard is read-only: actions run in the dashboard `mate console` serves', mate: [], crews: {} };
  const proj = ws.projects.find((x) => x.name === p);
  const crews = {};
  for (const t of openCrews(p)) crews[t.crew] = crewButtons(p, t);
  return { project: p, mate: mateButtons(proj), crews };
}

// The console's one line for each action, and what it does to the fake.
function runAction({ scope, action, project, crew, harness, name, path }) {
  const p = ws.projects.find((x) => x.name === project);
  const t = crew && tasks[project]?.find((x) => x.crew === crew);
  switch (`${scope ?? (crew ? 'crew' : 'mate')}:${action}`) {
    case 'crew:merge':
      Object.assign(t, { state: 'gone', closed: true, close_state: 'finished', merged_at: iso(0), closed_at: iso(0), status: 'done' });
      return `merged ${t.branch} into main (fast-forward 3 commits); crew ${crew} finished`;
    case 'crew:stop':
      if (t.dirty) throw refused(`worktree of ${crew} has 4 uncommitted changes; nothing was closed`);
      if (!t.merged_at && t.status !== 'idle') throw refused(`${t.branch} is not landed on main; nothing was closed`);
      Object.assign(t, { state: 'gone', closed: true, close_state: 'stopped', closed_at: iso(0) });
      return `closed ${project}/${crew}; brief, report and status file kept`;
    case 'crew:restart_crew':
      Object.assign(t, { state: 'at_desk_working', since: iso(0), status: 'working', status_text: 'restarted in its worktree', incident: undefined });
      return `restarted ${crew} (${t.harness}) in its worktree on ${t.branch}`;
    case 'mate:stop':
      Object.assign(p.mate, { running: false, state: '' });
      return `stopped the Mate of ${project}; ${openCrews(project).length} crews keep running`;
    case 'mate:restart_mate':
      Object.assign(p.mate, { running: true, state: 'idle', since: iso(0) });
      return `restarted the Mate of ${project}`;
    case 'mate:resume':
    case 'mate:start':
      Object.assign(p.mate, { running: true, state: 'idle', since: iso(0) });
      return `resumed the Mate of ${project} from its last session`;
    case 'mate:onboard':
      Object.assign(p.mate, { created: true, harness, running: true, state: 'idle', since: iso(0), tokens_today: 0, context_pct: 0, context_tokens: 0 });
      return `created the Mate of ${project} on ${harness}; it is reading the repo`;
    case 'mate:mode':
      p.mode = p.mode === 'auto' ? 'manual' : 'auto';
      return `${project} is ${p.mode}${p.mode === 'auto' ? ': the daemon types inbox answers into the Mate' : ': nothing is typed into the Mate for you'}`;
    case 'mate:clear_composer':
      return `cleared the composer of ${project}'s Mate (it held 1 unsent line)`;
    case 'mate:explain':
      return 'Jev: the Mate is at its prompt, idle; no dialog is open';
    case 'workspace:add_project':
      ws.projects.push({ name, repos: [`${name}`], mode: 'manual', inbox_waiting: 0, path: path || `~/workspace/${name}`, mate: { harness: '', running: false, state: '', tokens_today: null, context_pct: null, context_tokens: null, created: false }, fresh: true });
      tasks[name] = [];
      chats[name] = [];
      return `added project ${name}${path ? ` with repo ${path}` : ' (no repo yet)'}`;
    case 'workspace:remove_project': {
      const i = ws.projects.findIndex((x) => x.name === project);
      const n = openCrews(project).length;
      ws.projects.splice(i, 1);
      delete tasks[project];
      return `removed ${project} from the workspace; stopped ${p.mate.running ? 1 : 0} Mate and ${n} crews; files stay in ${p.path}`;
    }
    default:
      throw refused(`${action} is not offered here`);
  }
}

class ActionError extends Error {
  constructor(status, message) {
    super(message);
    this.status = status;
  }
}
const refused = (m) => new ActionError(409, m);

/** POST /api/actions: resolves to { line }, rejects with { status, message }. */
export function act(req) {
  const key = req.crew ? `${req.project}/${req.crew}` : req.project || '*';
  if (demo.readOnly) return Promise.reject(new ActionError(403, 'this dashboard is read-only: actions run in the dashboard `mate console` serves'));
  if (running.has(key)) return Promise.reject(new ActionError(409, `${running.get(key)} is still running on ${key}`));
  running.set(key, req.action);
  const how = demo.next;
  demo.next = 'auto';
  const ms = req.action === 'onboard' || req.action === 'add_project' ? 2200 : req.action === 'mode' || req.action === 'clear_composer' ? 500 : 1700;
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      running.delete(key);
      try {
        if (how === 'refuse') throw refused(`${req.action === 'merge' ? 'main moved since the review (2 commits); fast-forward not possible' : 'the console refused it'}; nothing changed`);
        if (how === 'fail') throw new ActionError(422, req.action === 'merge' ? 'git push origin main: ! [remote rejected] main -> main (protected branch hook declined)' : 'herdr: pane 7 did not answer within 10s');
        const line = runAction(req);
        bump(req.project);
        resolve({ line });
      } catch (e) {
        reject(e instanceof ActionError ? e : new ActionError(422, e.message));
      }
    }, ms);
  });
}
// Workspace actions (add, remove a project) go through the same door.
export const actWorkspace = (action, extra) => act({ scope: 'workspace', action, project: '', ...extra });

// A topic answered, the Mate's chat, a private fork.
export function answerTopic(id, answer) {
  const tp = topics.find((x) => x.id === id);
  return new Promise((resolve) => setTimeout(() => {
    Object.assign(tp, { status: 'on_it', answer });
    bump(tp.project);
    resolve({ line: `sent to the Mate of ${tp.project}` });
  }, 900));
}

export function assignToMate(project, crewName) {
  const t = tasks[project].find((x) => x.crew === crewName);
  return new Promise((resolve) => setTimeout(() => {
    t.assigned = true;
    bump(project);
    resolve({ line: `assigned ${crewName}'s question to the Mate of ${project}` });
  }, 700));
}

export const peekOf = (project, crewName) => {
  const t = tasks[project]?.find((x) => x.crew === crewName);
  if (!t) return '';
  const head = t.harness === 'codex' ? '>_ Codex' : t.harness === 'pi' ? 'π pi' : '✻ Claude Code';
  if (t.incident) return `${head}  ·  ${t.branch}\n\n⏺ Bash(xcodebuild test -scheme Mihu -destination 'platform=iOS Simulator,name=iPhone 17')\n  ⎿  Testing started…\n     Test Suite 'GalleryReorderTests' started\n     ▌\n\n✶ Running… (42m 08s · esc to interrupt)`;
  if (t.status === 'needs-decision') return `${head}  ·  ${t.branch}\n\n⏺ mate status needs-decision "${t.status_text}"\n  ⎿  sent to the Mate\n\n⏺ Waiting for the Mate's answer.\n\n> ▌`;
  if (t.status === 'wait-mate') return `${head}  ·  ${t.branch}\n\n⏺ gh pr checks 212\n  ⎿  ✓ build  ✓ test  ✓ lint\n\n⏺ mate status wait-mate "PR #212 green, review pass, ready to land"\n\n> ▌`;
  return `${head}  ·  ${t.branch}\n\n⏺ Edit(TripMap/TileCache.swift)\n  ⎿  Updated with 18 additions and 4 removals\n\n✶ Thinking… (3m 12s · ↓ 4.1k tokens)`;
};

export const diffOf = (project, crewName) => {
  const t = tasks[project]?.find((x) => x.crew === crewName);
  if (!t) return null;
  return {
    branch: t.branch, into: 'main', ahead: 3, behind: 0,
    files: [
      { path: 'api/share/token.go', add: 84, del: 0 },
      { path: 'api/share/handler.go', add: 41, del: 6 },
      { path: 'api/share/token_test.go', add: 120, del: 0 },
      { path: 'api/router.go', add: 3, del: 1 },
    ],
    hunk: '@@ -12,6 +12,9 @@ func Routes(r *mux.Router) {\n   r.HandleFunc("/trips/{id}", trip).Methods("GET")\n+  r.HandleFunc("/trips/{id}/share", share.Create).Methods("POST")\n+  r.HandleFunc("/s/{token}", share.Open).Methods("GET")\n-  r.HandleFunc("/trips/{id}/link", legacyLink)\n',
  };
};

export const daemonLog = [
  { at: iso(3), project: 'plantrip', text: 'typed the answer to PLA-63 ("use the v2 schema") into the Mate' },
  { at: iso(7), project: 'mihuv2', text: 'did not send HU-131\'s question: mihuv2 is manual' },
];

// ---- GET ------------------------------------------------------------------------

const env = () => ({ generated_at: new Date().toISOString(), last_event_id: eventID });

export async function api(path) {
  await new Promise((r) => setTimeout(r, 40));
  const url = new URL(path, 'http://x');
  const m = url.pathname.match(/^\/api\/projects\/([^/]+)(\/[a-z]+)?$/);
  if (url.pathname === '/api/workspace') {
    return { ...env(), root: ws.root, projects: ws.projects.map((p) => ({ ...p, inbox_waiting: inboxOf(p.name).filter((x) => x.attention).length })) };
  }
  if (url.pathname === '/api/events') {
    if (demo.stale) {
      await new Promise((r) => setTimeout(r, 600));
      throw new Error('the dashboard is not answering');
    }
    const since = Number(url.searchParams.get('since'));
    if (since < eventID) return { ...env(), events: [{ project: '' }], now: [] };
    await new Promise((resolve) => {
      const fn = () => {
        listeners.delete(fn);
        resolve();
      };
      listeners.add(fn);
      setTimeout(fn, 4000);
    });
    return { ...env(), events: since < eventID ? [{ project: '' }] : [], now: [] };
  }
  if (url.pathname === '/api/news') {
    return { ...env(), projects: ws.projects.map((p) => ({ name: p.name, news: topics.filter((x) => x.project === p.name).map((x) => ({ ...x, status: x.status === 'waiting' ? 'waiting' : 'answered' })) })) };
  }
  if (url.pathname === '/api/board') return { ...env(), days: 7, projects: ws.projects.map((p) => boardOf(p.name)) };
  if (m) {
    const p = decodeURIComponent(m[1]);
    if (!tasks[p]) throw new Error(`project ${p} is not in the workspace`);
    if (m[2] === '/history') {
      const crews = {};
      for (const t of tasks[p]) crews[t.crew] = linesFor(t);
      return { ...env(), crews };
    }
    if (m[2] === '/actions') return { ...env(), ...actionsOf(p) };
    const proj = ws.projects.find((x) => x.name === p);
    return { ...env(), project: p, mode: proj.mode, mate: proj.mate, tasks: tasks[p], inbox: inboxOf(p) };
  }
  throw new Error(`no fake for ${path}`);
}

export { inboxOf, openCrews };
