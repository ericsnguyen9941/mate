// What a tent's agents did, read from the dashboard's API, for the tent room
// to play back (CON-127): every task on the repo spawned in the last `days`
// days or still open, each crew's status lines, its questions and the
// answers, the reviewers' verdicts, the merges and the closes, in the order
// they happened. One request for the project's lines
// (/api/projects/{p}/history), beside the task rows the city already holds.

export const HISTORY_DAYS = 7;
const FALLBACK = 3; // a quiet tent shows its last few merged tasks instead

// Bubble-sized text: URLs become "PR #n", first sentence, at most n chars.
export function short(text, n = 72) {
  let s = String(text ?? '')
    .replace(/https?:\/\/\S+\/pull\/(\d+)\S*/g, 'PR #$1')
    .replace(/\/\S+\/(report|handback)\.md/g, '$1.md')
    .replace(/\s+/g, ' ')
    .trim();
  const stop = s.search(/[.;:](\s|$)/);
  if (stop > 20) s = s.slice(0, stop);
  return s.length > n ? `${s.slice(0, n - 1).trimEnd()}…` : s;
}
const long = (text) => short(text, 280);
const prOf = (lines) => [...lines].reverse().map((s) => s.text.match(/\/pull\/(\d+)/)?.[1]).find(Boolean);

// Midnight UTC `days` days ago: a since that stays the same all day, so the
// server's per-event cache answers every reader with the same key.
export function sinceOf(days, now = Date.now()) {
  const d = new Date(now - days * 864e5);
  d.setUTCHours(0, 0, 0, 0);
  return d.toISOString().replace('.000Z', 'Z');
}

/**
 * The task rows a tent plays back, and the time they start from.
 * @param {object[]} rows the project's task rows
 * @param {(task: object) => boolean} onTent whether a task worked in the tent
 */
export function tentTasks(rows, onTent, days = HISTORY_DAYS) {
  const mine = rows.filter((t) => onTent(t) && t.spawned_at);
  let since = sinceOf(days);
  let kept = mine.filter((t) => t.spawned_at >= since || !t.closed);
  if (!kept.some((t) => t.role === 'ship')) {
    // Nothing new here: the last few merged tasks and their reviews.
    const last = mine.filter((t) => t.role === 'ship' && t.merged_at).sort((a, b) => b.merged_at.localeCompare(a.merged_at)).slice(0, FALLBACK);
    const names = new Set(last.flatMap((t) => [t.crew, ...(t.reviews ?? [])]));
    kept = [...kept, ...mine.filter((t) => names.has(t.crew) && !kept.includes(t))];
  }
  const first = kept.map((t) => t.spawned_at).sort()[0];
  if (first && first < since) since = sinceOf(0, Date.parse(first));
  return { since, tasks: kept };
}

/**
 * Every event of the tent's tasks, oldest first. An event is
 * { at, crew, who, kind, text, short } where crew is the task it belongs to,
 * who is 'mate' or the crew that acted, and kind one of spawn, status, ask,
 * answer, handback, review, verdict, merge, close.
 */
export function buildEvents(tasks, crews) {
  const events = [];
  const add = (at, crew, who, kind, text, extra = {}) => at && events.push({ at, crew, who, kind, text: long(text), short: short(text), ...extra });
  const shipOf = new Map();
  for (const t of tasks) for (const r of t.reviews ?? []) shipOf.set(r, t.crew);
  for (const t of tasks) {
    const h = crews[t.crew] ?? { status_lines: [], questions: [] };
    const lines = h.status_lines ?? [];
    const review = t.role === 'review';
    if (review) {
      const of = shipOf.get(t.crew);
      add(t.spawned_at, t.crew, 'mate', 'review', `Ask ${t.harness || 'a reviewer'} to review ${of ?? t.text}`, { of });
    } else add(t.spawned_at, t.crew, 'mate', 'spawn', `Spawn ${t.crew}: ${t.text}`, { short: `Spawn ${t.crew}` });
    // A decision is the question asked by the needs-decision line it came from, else the next one asked.
    const questions = [...(h.questions ?? [])];
    const verdictAt = review && t.review ? lines.findLastIndex((s) => s.verb === 'wait-mate') : -1;
    lines.forEach((s, i) => {
      if (i === verdictAt) return;
      if (s.verb === 'needs-decision') {
        const k = questions.findIndex((q) => q.asked_event_id === s.event_id);
        const [q] = questions.splice(k >= 0 ? k : 0, 1);
        add(s.at, t.crew, t.crew, 'ask', q?.text ?? s.text);
        if (q?.answered_at && q.answer) add(q.answered_at, t.crew, 'mate', 'answer', q.answered_by === 'mate' ? q.answer : `Captain decided: ${q.answer}`);
      } else add(s.at, t.crew, t.crew, s.verb === 'wait-mate' ? 'handback' : 'status', s.text);
    });
    if (review && t.review) {
      const s = lines[verdictAt];
      add(s?.at ?? t.closed_at ?? t.spawned_at, t.crew, t.crew, 'verdict', s ? `${t.review}: ${s.text}` : t.review, { pass: t.review === 'pass' });
    }
    if (t.merged_at) {
      const pr = prOf(lines);
      add(t.merged_at, t.crew, 'mate', 'merge', pr ? `Merge ${t.crew}, PR #${pr}` : `Merge ${t.crew}`, { pr: pr ? Number(pr) : null });
    } else if (t.closed && !(review && t.review)) {
      add(t.closed_at, t.crew, 'mate', 'close', `Close ${t.crew}: ${(t.close_state || 'closed').replace(/_/g, ' ')}`);
    }
  }
  // Ties keep the order above: a spawn before the crew's first line, a verdict before the merge.
  return events.map((e, i) => [e, i]).sort((a, b) => a[0].at.localeCompare(b[0].at) || a[1] - b[1]).map(([e]) => e);
}

/**
 * A tent's history.
 * @param {(path: string) => Promise<object>} api GET on the dashboard
 * @param {string} project
 * @param {object[]} rows the project's task rows (GET /api/projects/{p})
 * @param {(task: object) => boolean} onTent
 * @returns {Promise<{ since: string, tasks: object[], events: object[] }>}
 */
export async function readHistory(api, project, rows, onTent, days = HISTORY_DAYS) {
  const { since, tasks } = tentTasks(rows, onTent, days);
  if (!tasks.length) return { since, tasks, events: [] };
  const body = await api(`/api/projects/${encodeURIComponent(project)}/history?since=${encodeURIComponent(since)}`);
  return { since, tasks, events: buildEvents(tasks, body.crews ?? {}) };
}
