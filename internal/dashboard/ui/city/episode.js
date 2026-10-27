// One real episode for a tent room to replay, read from the dashboard's API:
// the repo's latest merged ship task (preferring, among the last few, one
// that went through review and asked the Mate something), with the crew's
// status lines, its questions and the Mate's answers, the reviewers' status
// lines and verdicts, and the merge.

const CANDIDATES = 6; // latest merged ship tasks looked at
const RECENT = 4; // earlier merged tasks listed under Done
const STATUS_KEEP = 3; // plain status lines kept per agent: first, middle, last

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

const seat = (l) => ({ id: l.crew, harness: l.harness, model: l.model, effort: l.effort, role: l.role, text: l.text });
const prOf = (lines) => [...lines].reverse().map((s) => s.text.match(/\/pull\/(\d+)/)?.[1]).find(Boolean);

// Keeps every event but plain status lines, of which each agent keeps its
// first, middle and last, so a long task still replays in a minute or two.
function thin(events) {
  const byWho = new Map();
  events.forEach((e, i) => e.kind === 'status' && byWho.set(e.who, [...(byWho.get(e.who) ?? []), i]));
  const drop = new Set();
  for (const list of byWho.values()) {
    if (list.length <= STATUS_KEEP) continue;
    const keep = new Set([list[0], list[Math.floor(list.length / 2)], list.at(-1)]);
    for (const i of list) if (!keep.has(i)) drop.add(i);
  }
  return events.filter((_, i) => !drop.has(i));
}

function build(d, reviews) {
  const l = d.ledger;
  const events = [];
  const add = (at, who, kind, text, extra = {}) => at && events.push({ at, who, kind, text: long(text), short: short(text), ...extra });
  add(l.spawned_at, 'mate', 'spawn', `Spawn ${l.crew}: ${l.text}`, { short: `Spawn ${l.crew}` });
  const questions = [...(d.questions ?? [])];
  for (const s of d.status_lines ?? []) {
    if (s.verb === 'needs-decision') {
      const q = questions.shift();
      add(s.at, 'ship', 'ask', q?.text ?? s.text);
      // The Mate answers, or relays the captain's call.
      if (q?.answered_at && q.answer) add(q.answered_at, 'mate', 'answer', q.answered_by === 'mate' ? q.answer : `Captain decided: ${q.answer}`);
    } else add(s.at, 'ship', s.verb === 'wait-mate' ? 'handback' : 'status', s.text);
  }
  for (const r of reviews) {
    const rl = r.ledger;
    add(rl.spawned_at, 'mate', 'review', `Ask ${rl.harness} to review ${l.crew}`, { reviewer: rl.crew });
    // The verdict takes the place of the reviewer's last hand-back, at its time.
    const lines = r.status_lines ?? [];
    const last = rl.review ? lines.findLastIndex((s) => s.verb === 'wait-mate') : -1;
    lines.forEach((s, i) => i !== last && add(s.at, rl.crew, s.verb === 'wait-mate' ? 'handback' : 'status', s.text));
    if (rl.review) add(lines[last]?.at ?? rl.closed_at, rl.crew, 'verdict', lines[last] ? `${rl.review}: ${lines[last].text}` : rl.review, { pass: rl.review === 'pass' });
  }
  const pr = prOf(d.status_lines ?? []);
  add(l.merged_at, 'mate', 'merge', pr ? `Merge PR #${pr}` : `Merge ${l.crew}`);
  events.sort((a, b) => a.at.localeCompare(b.at));
  return {
    ship: { ...seat(l), tokens: l.tokens_with_reviews || l.tokens?.total || 0, turns: l.turns, questions: (d.questions ?? []).length, pr: pr ? Number(pr) : null },
    reviewers: reviews.map((r) => seat(r.ledger)),
    events: thin(events),
  };
}

/**
 * The episode for a tent, from the project's task rows.
 * @param {(path: string) => Promise<object>} api GET on the dashboard
 * @param {string} project
 * @param {object[]} tasks the project's task rows (GET /api/projects/{p})
 * @param {(task: object) => boolean} onTent whether a task worked in this tent
 * @returns {Promise<{ episode: object|null, recent: object[] }>}
 */
export async function readEpisode(api, project, tasks, onTent) {
  const merged = tasks.filter((t) => t.role === 'ship' && t.merged_at && onTent(t)).sort((a, b) => b.merged_at.localeCompare(a.merged_at));
  if (!merged.length) return { episode: null, recent: [] };
  const look = merged.slice(0, CANDIDATES);
  const score = (t) => 2 * ((t.reviews?.length ?? 0) > 0) + (t.question_count > 0);
  const best = Math.max(...look.map(score));
  const pick = look.find((t) => score(t) === best);
  const path = (crew) => `/api/projects/${encodeURIComponent(project)}/tasks/${encodeURIComponent(crew)}`;
  const [detail, ...reviews] = await Promise.all([pick.crew, ...(pick.reviews ?? [])].map((c) => api(path(c))));
  return {
    episode: build(detail, reviews),
    recent: merged.filter((t) => t !== pick).slice(0, RECENT).map((t) => ({ id: t.crew, text: short(t.text, 60) })),
  };
}
