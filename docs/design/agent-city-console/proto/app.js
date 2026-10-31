import * as THREE from '../city/vendor/three.js';
import { OrbitControls } from '../city/vendor/three.js';
import { createCity, lights, PALETTE, harnessName, formatTokens } from '../city/city.js';
import { createTentRoom } from '../city/room.js';
import { readHistory } from '../city/history.js';
import { HARNESSES } from '../city/harnesses.js';
import * as fake from './fake.js';

// Agent City as the console (design: ../index.html). The scene and the tent
// playback are the real city.js and room.js; this file adds what the design
// asks for on top:
//
// - places: City Hall and its mailbox (topics, projects, add/delete), a
//   block's plaza (its Mate: chat, inbox, actions), a crew (actions, peek,
//   diff), a tent (playback, the crews in it);
// - one panel at a time, docked right on a wide screen and a bottom sheet
//   on a phone, the camera flown so the thing stays in view beside it;
// - pins: DOM buttons over the scene for only what needs the captain
//   (topic waiting, ready to merge, stuck) and what is running, so they
//   read at any distance and Tab reaches them;
// - every action as the console's: fixed buttons, a reason under the ones
//   that cannot run, a question with Cancel focused, Running …, one line back.

const $ = (s) => document.querySelector(s);
const enc = encodeURIComponent;
const phone = () => innerWidth <= 760;
const clock = (iso) => (iso ? new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '?');
const ago = (iso) => {
  if (!iso) return '?';
  const m = Math.round((Date.now() - Date.parse(iso)) / 60e3);
  return m < 1 ? 'now' : m < 60 ? `${m}m` : m < 1440 ? `${Math.round(m / 60)}h` : `${Math.round(m / 1440)}d`;
};
const unk = (v, f = (x) => x) => (v == null ? '?' : f(v));

// h('div', { class: 'x', onclick }, 'text', child…): text is always a text node.
function h(tag, attrs = {}, ...kids) {
  const el = document.createElement(tag);
  for (const [k, v] of Object.entries(attrs ?? {})) {
    if (v == null || v === false) continue;
    if (k.startsWith('on')) el.addEventListener(k.slice(2), v);
    else if (k === 'class') el.className = v;
    else if (k === 'dataset') Object.assign(el.dataset, v);
    else if (k in el && typeof v !== 'string') el[k] = v;
    else el.setAttribute(k, v === true ? '' : v);
  }
  for (const k of kids.flat()) if (k != null && k !== false) el.append(k instanceof Node ? k : document.createTextNode(String(k)));
  return el;
}

// ---- Renderer, camera ---------------------------------------------------------

const canvas = $('#c');
const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
renderer.setPixelRatio(Math.min(devicePixelRatio, 2));
renderer.outputColorSpace = THREE.SRGBColorSpace;
renderer.toneMapping = THREE.NeutralToneMapping;
renderer.shadowMap.enabled = true;
renderer.shadowMap.type = THREE.PCFShadowMap;
const camera = new THREE.PerspectiveCamera(30, 1, 0.5, 600);
const controls = new OrbitControls(camera, canvas);
controls.enableDamping = true;
controls.maxPolarAngle = 1.35;
const CITY_LIMITS = { minDistance: 8, minAzimuthAngle: -1.1, maxAzimuthAngle: 1.1 };
const ROOM_LIMITS = { minDistance: 6, maxDistance: 34, minAzimuthAngle: -0.9, maxAzimuthAngle: 0.9 };

// ---- State -------------------------------------------------------------------

const S = {
  ws: null,
  tasks: new Map(), // project → task rows
  inbox: new Map(), // project → inbox rows
  actions: new Map(), // project → GET /actions
  stale: '',
  staleSince: null,
  sel: null, // { kind: 'mate'|'crew'|'hall'|'tent'|'system'|'goto'|'add'|'lot', project, crew, tab, topic, back }
  busy: new Map(), // row key → label running
  results: new Map(), // row key → { kind, text, at }
  drafts: new Map(), // textarea id → text
  pending: new Map(), // project or topic id → { text, left, timer } sending with undo
  turns: new Map(), // project → { phase, tool, timer }
  folded: new Map(), // project → show earlier turns
  inboxAll: false,
  forks: new Map(), // topic id → { messages, proposal }
  timelineStuck: false,
  herdrDown: false,
  flashes: [], // { anchor, kind, glyph, text, until }
  newBlock: null,
};
const rowKey = (project, crew) => (crew ? `${project}/${crew}` : project);
const projectOf = (name) => S.ws?.projects.find((p) => p.name === name);
const taskOf = (project, crew) => S.tasks.get(project)?.find((t) => t.crew === crew);
const openCrews = (project) => (S.tasks.get(project) ?? []).filter((t) => !t.closed);

let scene = null;
let city = null;
let shape = '';
let home = null;
const start = performance.now();
const now = () => (performance.now() - start) / 1000;
const tentRepos = (p) => (p.repos?.length ? p.repos : [p.name]);
function onTent(p, repo) {
  const repos = tentRepos(p);
  return (t) => t.repo === repo || (repo === repos[0] && !repos.includes(t.repo));
}

function buildCity(projects) {
  if (city) {
    scene.remove(city.group);
    city.dispose();
    scene.environment?.dispose();
  }
  city = createCity(projects.map((p) => ({ name: p.name, repos: p.repos ?? [], mate: { harness: p.mate?.harness } })));
  scene = new THREE.Scene();
  scene.background = new THREE.Color(PALETTE.paper);
  const k = city.size / 3;
  scene.fog = new THREE.Fog(PALETTE.paper, 170 * k, 320 * k);
  lights(scene, renderer, 27 * k);
  scene.add(city.group);
  scene.add(ring);
  scene.add(crate);
  // A phone is tall and narrow: stand further back so the whole island fits across.
  home = { target: new THREE.Vector3(0, 1, 1), position: new THREE.Vector3(-20, 62, 70).multiplyScalar(k * (phone() ? 1.9 : 1)) };
  controls.maxDistance = 170 * k;
  if (!room && !S.sel) {
    controls.target.copy(home.target);
    camera.position.copy(home.position);
  }
}

// The selection ring on the ground, and the crate a merge carries.
const ring = new THREE.Mesh(new THREE.TorusGeometry(0.75, 0.06, 8, 48), new THREE.MeshBasicMaterial({ color: PALETTE.coral, toneMapped: false }));
ring.rotation.x = -Math.PI / 2;
ring.visible = false;
const crate = new THREE.Mesh(new THREE.BoxGeometry(0.45, 0.35, 0.45), new THREE.MeshStandardMaterial({ color: '#C99A62', roughness: 0.8 }));
crate.castShadow = true;
crate.visible = false;
let crateFlight = null;

function applyProject(p, t) {
  const rows = S.tasks.get(p.name);
  if (!rows || !city) return;
  city.setCrews(p.name, rows.filter((r) => r.state !== 'gone'), t, { fresh: true });
  for (const repo of tentRepos(p)) city.setMerged(p.name, repo, rows.filter((r) => r.role === 'ship' && r.merged_at && onTent(p, repo)(r)).length, t);
}

let board = null;
async function loadAll() {
  const body = await fake.api('/api/workspace');
  S.ws = body;
  const projects = body.projects;
  const next = JSON.stringify(projects.map((p) => [p.name, p.repos ?? [], p.mate?.harness]));
  const t = now();
  const rebuilt = next !== shape;
  if (rebuilt) {
    shape = next;
    buildCity(projects);
  }
  await Promise.all(projects.map(async (p) => {
    const [page, acts] = await Promise.all([fake.api(`/api/projects/${enc(p.name)}`), fake.api(`/api/projects/${enc(p.name)}/actions`)]);
    S.tasks.set(p.name, page.tasks);
    S.inbox.set(p.name, page.inbox);
    S.actions.set(p.name, acts);
  }));
  for (const name of [...S.tasks.keys()]) if (!projectOf(name)) S.tasks.delete(name);
  for (const p of projects) {
    city.setMate(p.name, p.mate ?? {}, t);
    applyProject(p, t);
  }
  board ??= await fake.api('/api/board?days=7');
  for (const p of projects) {
    const bp = board.projects.find((x) => x.name === p.name);
    city.setBoard(p.name, bp ? { days: 7, tokens: bp.days.reduce((n, d) => n + d.tokens.total, 0), finished: bp.days.reduce((n, d) => n + d.closed.finished, 0), passed: bp.days.reduce((n, d) => n + d.reviews.pass, 0) } : null, t);
  }
  city.setHall(projects.map((p) => ({ name: p.name, running: Boolean(p.mate?.running), working: openCrews(p.name).length, tokens: p.mate?.tokens_today ?? 0 })));
  city.setMailbox(fake.topics.filter((x) => x.status === 'waiting').length, t);
  if (!S.stale) setLive('live', 'live');
  if (room) refreshRoom();
  render();
  return rebuilt;
}

function setLive(state, text) {
  const pill = $('#live');
  pill.dataset.state = state;
  pill.querySelector('.g').textContent = state === 'stale' ? '▲' : '●';
  pill.querySelector('.t').textContent = text;
}

// One long poll on /api/events, like the real page.
async function poll() {
  for (;;) {
    try {
      await fake.api(`/api/events?since=0&wait=20`);
      if (S.stale) {
        S.stale = '';
        S.staleSince = null;
      }
      await loadAll();
    } catch (e) {
      if (!S.stale) S.staleSince = new Date().toISOString();
      S.stale = e.message;
      setLive('stale', `stale · ${e.message}`);
      render();
      await new Promise((r) => setTimeout(r, 1500));
    }
  }
}

// ---- What needs the captain ---------------------------------------------------------
//
// Three kinds earn a pin and a row in the rail, and only these: a topic
// waiting on the captain, a crew waiting on merge, an incident. A crew's
// question to its Mate is the Mate's to answer; it shows only when the
// project is manual and nobody took it (then the captain is the only one
// who can), and as a quiet pin. One pin per thing; the most urgent wins.

const KINDS = {
  incident: { glyph: '!', rank: 0 },
  topic: { glyph: '?', rank: 1 },
  merge: { glyph: '⇢', rank: 2 },
  decision: { glyph: '◆', rank: 3 },
};

function attention() {
  const out = [];
  const waiting = fake.topics.filter((x) => x.status === 'waiting' && projectOf(x.project));
  for (const tp of waiting) out.push({ kind: 'topic', anchor: 'hall', title: tp.question, sub: `${tp.project} · the Mate asks you · ${ago(tp.at)}`, open: () => select({ kind: 'hall', tab: 'topics', topic: tp.id }) });
  for (const p of S.ws?.projects ?? []) {
    for (const t of openCrews(p.name)) {
      const anchor = `crew:${p.name}/${t.crew}`;
      if (t.incident) out.push({ kind: 'incident', anchor, title: `${t.crew} stuck ${ago(t.since)}`, sub: `${p.name} · ${t.status_text}`, short: 'stuck', open: () => select({ kind: 'crew', project: p.name, crew: t.crew, tab: 'peek' }) });
      if (t.status === 'wait-mate') out.push({ kind: 'merge', anchor, title: `${t.crew} ready to merge`, sub: `${p.name} · ${t.branch} → main`, short: 'merge', open: () => select({ kind: 'crew', project: p.name, crew: t.crew }) });
      if (t.status === 'needs-decision' && p.mode === 'manual' && !t.assigned) out.push({ kind: 'decision', anchor, title: `${t.crew} asks its Mate`, sub: `${p.name} is manual · nobody took it`, short: 'asks', open: () => select({ kind: 'crew', project: p.name, crew: t.crew, tab: 'status' }) });
    }
  }
  return out.sort((a, b) => KINDS[a.kind].rank - KINDS[b.kind].rank);
}

function renderRail(items) {
  const top = items.filter((x) => x.kind !== 'decision');
  $('#needs-n').textContent = top.length ? `Needs you · ${top.length}` : 'Nothing needs you';
  $('#needs').dataset.tone = top.length ? 'need' : 'ok';
  const rail = $('#rail');
  rail.replaceChildren(h('h2', {}, top.length ? `Needs you · ${top.length}` : 'Needs you'),
    ...(items.length ? items.map((x) => h('button', { type: 'button', class: `k-${x.kind}`, onclick: x.open },
      h('span', { class: 'g', 'aria-hidden': 'true' }, KINDS[x.kind].glyph), h('b', {}, x.title), h('small', {}, x.sub))) : [h('p', { class: 'empty' }, 'Nothing waits on you.')]));
}

// ---- Pins ---------------------------------------------------------------------------

const pinsEl = $('#pins');
const pinEls = new Map();
function anchorPos(anchor, out) {
  if (!city) return null;
  if (anchor === 'hall') {
    city.hall[0].getWorldPosition(out);
    return out.setY(out.y + 2.1);
  }
  const [kind, rest] = anchor.split(':');
  if (kind === 'crew') {
    const m = city.mascots.get(rest);
    return m ? out.copy(m.root.position).setY(m.root.position.y + 1.45) : null;
  }
  if (kind === 'mate') {
    const m = city.mascots.get(`${rest}:mate`);
    return m ? out.copy(m.root.position).setY(m.root.position.y + 1.85) : null;
  }
  if (kind === 'tent') {
    const b = city.tents.get(rest);
    return b ? out.set(b.x - b.w / 2 + 0.7, 1.6, b.z + b.d / 2 + 0.45) : null;
  }
  if (kind === 'block') {
    const b = city.blocks.get(rest);
    return b ? out.set(b.centre.x, 3.2, b.centre.z) : null;
  }
  return null;
}

function wantedPins() {
  const pins = new Map(); // anchor → pin
  const put = (anchor, pin) => {
    const have = pins.get(anchor);
    if (!have || pin.rank < have.rank) pins.set(anchor, pin);
  };
  const items = attention();
  const topics = items.filter((x) => x.kind === 'topic').length;
  if (topics) put('hall', { kind: 'topic', rank: 1, glyph: '?', long: `${topics} waiting on you`, short: String(topics), label: `City Hall: ${topics} topics waiting on you`, open: () => select({ kind: 'hall', tab: 'topics' }) });
  for (const x of items) if (x.anchor !== 'hall') put(x.anchor, { kind: x.kind, rank: KINDS[x.kind].rank + 1, glyph: KINDS[x.kind].glyph, long: x.title, short: x.short, label: x.title, open: x.open });
  for (const p of S.ws?.projects ?? []) {
    if (p.mate?.created === false) put(`mate:${p.name}`, { kind: 'empty', rank: 5, glyph: '+', long: 'No Mate yet', short: 'no Mate', label: `${p.name} has no Mate yet`, open: () => select({ kind: 'mate', project: p.name }) });
  }
  // Running and just-finished actions sit over their thing and win.
  for (const [key, label] of S.busy) {
    const [project, crew] = key.split('/');
    put(crew ? `crew:${key}` : projectOf(project) ? `mate:${project}` : 'hall', { kind: 'run', rank: -2, glyph: '⟳', long: `Running ${label.replace(/…$/, '').toLowerCase()}…`, short: '…', label: `Running ${label}` , open: () => select(crew ? { kind: 'crew', project, crew } : { kind: 'mate', project }) });
  }
  const t = performance.now();
  S.flashes = S.flashes.filter((f) => f.until > t);
  for (const f of S.flashes) put(f.anchor, { kind: f.kind, rank: -1, glyph: f.glyph, long: f.text, short: f.glyph, label: f.text, shake: f.kind === 'bad' || f.kind === 'refused', open: f.open });
  return pins;
}

let pinsKey = '';
function renderPins() {
  const want = room ? new Map() : wantedPins();
  const key = JSON.stringify([...want].map(([a, p]) => [a, p.kind, p.long]));
  if (key === pinsKey) return;
  pinsKey = key;
  for (const [a, el] of pinEls) if (!want.has(a)) (el.remove(), pinEls.delete(a));
  for (const [a, p] of want) {
    const el = h('button', { type: 'button', class: `pin k-${p.kind}${p.shake ? ' shake' : ''}`, 'aria-label': p.label, onclick: p.open },
      h('span', { class: 'g', 'aria-hidden': 'true' }, p.glyph), h('span', { class: 'w' }, p.long));
    el.dataset.short = p.short ?? '';
    if (pinEls.has(a)) pinEls.get(a).replaceWith(el);
    else pinsEl.append(el);
    pinEls.set(a, el);
  }
}

const v = new THREE.Vector3();
function placePins() {
  const w = canvas.clientWidth;
  const hh = canvas.clientHeight;
  for (const [a, el] of pinEls) {
    if (!anchorPos(a, v)) {
      el.style.display = 'none';
      continue;
    }
    const dist = v.distanceTo(camera.position);
    v.project(camera);
    const off = v.z > 1;
    el.style.display = off ? 'none' : '';
    el.style.transform = `translate(${((v.x + 1) / 2) * w}px, ${((1 - v.y) / 2) * hh}px) translate(-50%, calc(-100% - 6px))`;
    // Far away a pin keeps its glyph and drops its words; the rail and the label still say them.
    el.classList.toggle('far', dist > 70 && !el.classList.contains('k-run'));
    el.classList.toggle('dim', Boolean(S.sel) && !el.matches(':focus') && !isSelected(a));
  }
}
const isSelected = (a) => S.sel && ((S.sel.kind === 'crew' && a === `crew:${S.sel.project}/${S.sel.crew}`) || (S.sel.kind === 'mate' && a === `mate:${S.sel.project}`) || (S.sel.kind === 'hall' && a === 'hall'));

function flash(anchor, kind, glyph, text, ms = 3500) {
  S.flashes.push({ anchor, kind, glyph, text, until: performance.now() + ms });
  setTimeout(renderPins, ms + 50);
}

// ---- Camera ------------------------------------------------------------------------

let flight = null;
function fly(target, position) {
  flight = { from: [controls.target.clone(), camera.position.clone()], to: [target, position], t0: performance.now() };
}
function focusOn(sel) {
  if (!city || room) return;
  const p = new THREE.Vector3();
  let dist = [0, 24, 32];
  if (sel.kind === 'crew' && anchorPos(`crew:${sel.project}/${sel.crew}`, p)) p.y -= 1.2;
  else if (sel.kind === 'mate' && anchorPos(`mate:${sel.project}`, p)) p.y -= 1.6;
  else if (sel.kind === 'hall' || sel.kind === 'system') {
    anchorPos('hall', p);
    p.y = 1.5;
    dist = [0, 16, 24];
  } else if (sel.kind === 'lot' && sel.at) {
    p.copy(sel.at);
    dist = [0, 16, 22];
  } else if (sel.kind === 'block' && anchorPos(`block:${sel.project}`, p)) {
    p.y = 0.5;
    dist = [0, 16, 22];
  } else return;
  fly(p.clone(), p.clone().add(new THREE.Vector3(...dist)));
  ring.visible = sel.kind === 'crew' || sel.kind === 'mate';
  if (ring.visible) {
    ring.position.set(p.x, 0.2, p.z);
    ring.scale.setScalar(sel.kind === 'mate' ? 1.15 : 0.85);
  }
}

// ---- Selection, the dock --------------------------------------------------------------

const dock = $('#dock');
const dbody = $('#dbody');
let lastFocus = null;

function select(sel, { fly: doFly = true } = {}) {
  if (room && sel.kind !== 'crew' && sel.kind !== 'tent' && sel.kind !== 'mate') leaveRoom();
  if (!S.sel) lastFocus = document.activeElement;
  // Chatter never talks over the thing you are working with.
  for (const m of city?.mascots.values() ?? []) if (m.chatAt === Infinity) m.chatAt = undefined;
  const id = sel.kind === 'crew' ? `${sel.project}/${sel.crew}` : sel.kind === 'mate' ? `${sel.project}:mate` : null;
  if (id && city?.mascots.get(id)) city.mascots.get(id).chatAt = Infinity;
  S.sel = sel;
  document.body.classList.add('docked');
  if (phone() && !document.body.classList.contains('sheet-full')) document.body.classList.remove('sheet-peek');
  resize();
  if (doFly) focusOn(sel);
  render();
  $('#dock-title').focus({ preventScroll: true });
}

function closeDock() {
  if (room && S.sel?.kind === 'crew' && S.sel.back) return select(S.sel.back);
  if (room) return leaveRoom();
  if (S.sel?.back) return select(S.sel.back);
  for (const m of city?.mascots.values() ?? []) if (m.chatAt === Infinity) m.chatAt = undefined;
  S.sel = null;
  ring.visible = false;
  document.body.classList.remove('docked', 'sheet-peek', 'sheet-full');
  resize();
  render();
  lastFocus?.focus?.();
}
$('#close').addEventListener('click', closeDock);
$('#grab').addEventListener('click', () => {
  const b = document.body.classList;
  if (b.contains('sheet-peek')) b.remove('sheet-peek');
  else if (b.contains('sheet-full')) (b.remove('sheet-full'), b.add('sheet-peek'));
  else b.add('sheet-full');
  resize();
});

function render() {
  renderRail(attention());
  renderPins();
  renderSystemChip();
  if (!S.sel) return;
  // Keep what the captain is typing, and where they are, across a redraw.
  const focusId = document.activeElement?.id;
  const caret = document.activeElement?.selectionStart;
  const scroll = dbody.scrollTop;
  const view = VIEWS[S.sel.kind]?.() ?? { title: '?', body: [] };
  $('#crumbs').replaceChildren(...(view.crumbs ?? []).flatMap((c, i) => [...(i ? [h('span', { 'aria-hidden': 'true' }, '›')] : []), c.go ? h('button', { type: 'button', onclick: c.go }, c.text) : h('span', {}, c.text)]));
  $('#dock-title').textContent = view.title;
  $('#dock-sub').textContent = view.sub ?? '';
  dbody.replaceChildren(...view.body.filter(Boolean));
  dbody.scrollTop = scroll;
  if (focusId && focusId !== 'dock-title') {
    const el = document.getElementById(focusId);
    if (el) {
      el.focus({ preventScroll: true });
      if (caret != null && el.setSelectionRange) el.setSelectionRange(caret, caret);
    }
  }
}

// ---- Shared pieces of a panel ----------------------------------------------------------

const STATE = {
  at_desk_working: ['▶', 'working', 'ok'],
  waiting_at_ceo: ['◆', 'asking its Mate', 'need'],
  walking_to_ceo: ['◆', 'asking its Mate', 'need'],
  waiting_review: ['⇢', 'waiting on merge', 'need'],
  blocked: ['■', 'blocked', 'warn'],
  asleep: ['◌', 'asleep', 'off'],
  idle: ['○', 'idle', 'off'],
  arriving: ['›', 'arriving', 'ok'],
  leaving: ['‹', 'leaving', 'off'],
};
function stateEl(state, override) {
  const [g, word, tone] = override ?? STATE[state] ?? ['?', state || '?', 'off'];
  return h('span', { class: 'state', dataset: { tone } }, h('span', { class: 'g', 'aria-hidden': 'true' }, g), word);
}
function facts(rows) {
  return h('dl', { class: 'facts' }, rows.map(([k, val]) => h('div', {}, h('dt', {}, k), h('dd', { class: val === '?' ? 'unk' : '' }, val))));
}
function bands(extra = []) {
  const out = [];
  if (S.stale) out.push(h('div', { class: 'band stale', role: 'status' }, h('span', { class: 'g' }, '▲'), h('span', {}, `Stale since ${clock(S.staleSince)}: ${S.stale}. What you see may be old; actions wait until the city is live again.`)));
  return [...out, ...extra];
}
function resultEl(key) {
  const r = S.results.get(key);
  if (!r) return null;
  const glyph = { ok: '✓', refused: '⊘', failed: '✕', info: 'i' }[r.kind];
  const lead = { ok: '', refused: 'Refused, nothing changed: ', failed: 'Ran and failed: ', info: '' }[r.kind];
  return h('div', { class: 'result', dataset: { kind: r.kind }, role: 'status', 'aria-live': 'polite' }, h('span', { class: 'g', 'aria-hidden': 'true' }, glyph),
    h('span', {}, h('b', {}, lead), r.text, h('span', { style: 'opacity:.7' }, ` · ${clock(r.at)}`)));
}

// Why a button cannot run right now, beyond what the API says: the page's own state.
function blockedBy(key) {
  if (S.stale) return `stale since ${clock(S.staleSince)}; wait for live`;
  if (room && mode !== 'live') return 'replaying; go Live to act';
  const running = S.busy.get(key);
  if (running) return `waiting: ${running.replace(/…$/, '')} is running`;
  return '';
}

/**
 * The Actions block: the API's buttons in its order, each with its reason
 * line kept even when empty, so nothing moves when a button changes state.
 */
function actionsBlock({ project, crew, buttons, extra = [], links = [] }) {
  const acts = S.actions.get(project);
  const key = rowKey(project, crew);
  if (acts?.read_only) {
    return [h('h3', { class: 'sec' }, 'Actions'), h('div', { class: 'band ro' }, h('span', { class: 'g' }, '⊘'), h('span', {}, `Read-only: ${acts.read_only.replace(/^this dashboard is read-only: /, '')}.`)), links.length ? h('div', { class: 'links' }, links) : null];
  }
  const running = S.busy.get(key);
  const page = blockedBy(key);
  const els = [...(buttons ?? []), ...extra].map((b) => {
    const isRunning = running && running === b.label;
    const why = isRunning ? 'running…' : b.toggle ? (S.stale ? '' : '') : page || (b.available ? '' : b.reason);
    const disabled = isRunning || (!b.toggle && (Boolean(page) || !b.available));
    const danger = /^(stop|restart|restart_mate|restart_crew|remove_project)$/.test(b.action);
    const primary = b.action === 'merge';
    return h('div', { class: 'act' },
      h('button', { type: 'button', id: `act-${key}-${b.action}`, class: `${isRunning ? 'running' : ''} ${danger ? 'danger' : ''} ${primary ? 'primary' : ''}`, disabled, 'aria-pressed': b.toggle ? String(Boolean(b.pressed)) : null, 'aria-describedby': `why-${key}-${b.action}`, onclick: () => (b.run ? b.run() : press(project, crew, b)) },
        isRunning ? h('span', { class: 'g', 'aria-hidden': 'true' }, '⟳') : null, isRunning ? `Running ${b.label.replace(/…$/, '').toLowerCase()}…` : b.label),
      h('small', { id: `why-${key}-${b.action}` }, why));
  });
  return [h('h3', { class: 'sec' }, 'Actions', crew ? null : h('span', { class: 'r' }, S.ws && projectOf(project)?.mode === 'auto' ? '● auto' : '○ manual')), h('div', { class: 'actions' }, els), resultEl(key), links.length ? h('div', { class: 'links' }, links) : null];
}

// A destructive press asks first, Cancel focused; then runs and reports one line.
const TITLES = {
  'crew:stop': (p, t) => [`Stop crew ${t.crew}?`, `${p}/${t.crew} · worktree and branch ${t.branch}`, 'Stop crew'],
  'crew:restart_crew': (p, t) => [`Restart crew ${t.crew}?`, `${p}/${t.crew} · same worktree, same branch`, 'Restart crew'],
  'crew:merge': (p, t) => [`Merge ${t.branch} → main?`, `${p}/${t.crew} · fast-forward only · the crew finishes`, 'Merge'],
  'mate:stop': (p) => [`Stop the Mate of ${p}?`, `${p} · the Mate only · ${openCrews(p).length} crews keep running`, 'Stop mate'],
  'mate:restart_mate': (p) => [`Restart the Mate of ${p}?`, `${p} · the Mate only · ${openCrews(p).length} crews keep running`, 'Restart mate'],
};
function press(project, crew, b) {
  if (b.action === 'onboard') return select({ kind: 'create', project, back: { kind: 'mate', project } });
  const t = crew ? taskOf(project, crew) : null;
  const title = TITLES[`${crew ? 'crew' : 'mate'}:${b.action}`]?.(project, t);
  const go = () => run(rowKey(project, crew), b.label, { action: b.action, project, crew }, crew ? `crew:${project}/${crew}` : `mate:${project}`);
  if (!b.confirm) return go();
  confirmBox({ title: title?.[0] ?? `${b.label.replace(/…$/, '')}?`, scope: title?.[1], body: b.confirm, yes: title?.[2] ?? b.label.replace(/…$/, ''), onYes: go });
}

async function run(key, label, req, anchor, after) {
  S.busy.set(key, label);
  S.results.delete(key);
  render();
  if (req.action === 'merge') startCrate(req.project, req.crew);
  try {
    const { line } = req.scope === 'workspace' ? await fake.actWorkspace(req.action, req) : await fake.act(req);
    S.results.set(key, { kind: 'ok', text: line, at: new Date().toISOString() });
    const t = req.crew && taskOf(req.project, req.crew);
    if (req.action === 'merge' && t) flash(`tent:${req.project}/${t.repo}`, 'ok', '✓', `merged ${req.crew}`, 5000);
    else flash(anchor, 'ok', '✓', 'done');
    after?.(true);
  } catch (e) {
    const kind = e.status === 422 ? 'failed' : 'refused';
    S.results.set(key, { kind, text: e.message, at: new Date().toISOString() });
    flash(anchor, kind === 'failed' ? 'bad' : 'refused', '✕', kind === 'failed' ? 'failed' : 'refused');
    if (req.action === 'merge') crateBack();
    after?.(false);
  } finally {
    S.busy.delete(key);
    await loadAll();
  }
}

const scrim = $('#scrim');
function confirmBox({ title, scope, body, list, yes, onYes }) {
  const back = document.activeElement;
  const close = () => {
    scrim.classList.remove('on');
    scrim.replaceChildren();
    back?.focus?.();
  };
  const cancel = h('button', { type: 'button', class: 'btn', onclick: close }, 'Cancel');
  const box = h('div', { class: 'confirm', role: 'alertdialog', 'aria-modal': 'true', 'aria-labelledby': 'cf-t', 'aria-describedby': 'cf-b' },
    h('h2', { id: 'cf-t' }, title),
    scope ? h('div', { class: 'scope' }, scope) : null,
    h('p', { id: 'cf-b' }, body),
    list ? h('ul', {}, list.map((x) => h('li', {}, x))) : null,
    h('div', { class: 'row' }, cancel, h('button', { type: 'button', class: 'btn coral', onclick: () => (close(), onYes()) }, yes)));
  box.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') (e.stopPropagation(), close());
    if (e.key === 'Tab') {
      const f = [...box.querySelectorAll('button')];
      const i = f.indexOf(document.activeElement);
      e.preventDefault();
      f[(i + (e.shiftKey ? -1 : 1) + f.length) % f.length].focus();
    }
  });
  scrim.replaceChildren(box);
  scrim.classList.add('on');
  cancel.focus();
}

// A send with a short undo: the text waits `left` seconds, then goes.
function sendLater(id, text, go) {
  const p = { text, left: 5 };
  p.timer = setInterval(() => {
    p.left--;
    if (p.left <= 0) {
      clearInterval(p.timer);
      S.pending.delete(id);
      go(text);
    }
    render();
  }, 1000);
  S.pending.set(id, p);
  render();
}
function undo(id, draftId) {
  const p = S.pending.get(id);
  if (!p) return;
  clearInterval(p.timer);
  S.pending.delete(id);
  S.drafts.set(draftId, p.text);
  render();
  document.getElementById(draftId)?.focus();
}
function textarea(id, attrs = {}) {
  const el = h('textarea', { id, ...attrs, oninput: (e) => S.drafts.set(id, e.target.value) });
  el.value = S.drafts.get(id) ?? '';
  return el;
}

// ---- Views ----------------------------------------------------------------------------

const VIEWS = {};

function tabs(list, current, set) {
  return h('div', { class: 'tabs', role: 'tablist' }, list.map(([id, label, n]) => h('button', { type: 'button', role: 'tab', id: `tab-${id}`, 'aria-selected': String(current === id), onclick: () => set(id) }, label, n ? h('span', { class: 'n' }, String(n)) : null)));
}
const setTab = (tab) => ((S.sel = { ...S.sel, tab }), render());

// The Mate, at its plaza.
VIEWS.mate = () => {
  const p = projectOf(S.sel.project);
  if (!p) return { title: 'Gone', body: [h('p', {}, `${S.sel.project} is no longer in the workspace.`)] };
  const m = p.mate;
  const created = m.created !== false;
  const acts = S.actions.get(p.name);
  const inbox = (S.inbox.get(p.name) ?? []).filter((x) => x.attention);
  const tab = S.sel.tab ?? 'chat';
  const copy = h('button', { type: 'button', class: 'btn', style: 'min-height:32px', onclick: () => {
    const id = `mate-${p.name}-7f3a`;
    navigator.clipboard?.writeText(id).catch(() => {});
    S.results.set(p.name, { kind: 'info', text: `copied agent id ${id}`, at: new Date().toISOString() });
    render();
  } }, 'Copy agent id');
  const state = !created ? stateEl('', ['+', 'no Mate yet', 'off']) : m.running ? stateEl(m.state) : stateEl('', ['◌', 'stopped', 'off']);
  return {
    crumbs: [{ text: 'City' }, { text: p.name, go: () => select({ kind: 'block', project: p.name }) }],
    title: `Mate · ${p.name}`,
    sub: created ? `${harnessName(m.harness)} · plaza of ${p.name}` : 'This block has a plaza but no Mate',
    body: [
      ...bands(),
      facts([
        ['State', state],
        ['Mode', p.mode === 'auto' ? '● auto' : '○ manual'],
        ['Tokens today', unk(m.tokens_today, formatTokens)],
        ['Context', m.context_pct == null ? '?' : `${m.context_pct}% · ${formatTokens(m.context_tokens)}`],
      ]),
      ...actionsBlock({ project: p.name, buttons: acts?.mate, extra: created ? [] : [] }),
      created ? h('div', { class: 'links' }, copy) : null,
      tabs([['chat', 'Chat'], ['inbox', 'Inbox', inbox.length], ['apps', 'Apps'], ['project', 'Project']], tab, setTab),
      ...{ chat: chatView, inbox: inboxView, apps: appsView, project: projectView }[tab](p),
    ],
  };
};

function chatView(p) {
  const m = p.mate;
  const msgs = fake.chats[p.name] ?? [];
  const turn = S.turns.get(p.name);
  const pend = S.pending.get(`chat:${p.name}`);
  const blocked = m.created === false ? 'no Mate yet · Create mate first' : !m.running ? 'the Mate is stopped · Resume mate first' : S.stale ? `stale since ${clock(S.staleSince)}` : S.actions.get(p.name)?.read_only ? 'read-only dashboard' : '';
  // Older turns fold to one line; the last two exchanges stay open.
  const shownFrom = S.folded.get(p.name) ? 0 : Math.max(0, msgs.length - 4);
  const list = [
    shownFrom ? h('button', { type: 'button', class: 'fold', onclick: () => (S.folded.set(p.name, true), render()) }, `Show ${Math.ceil(shownFrom / 2)} earlier turns`) : null,
    ...msgs.slice(shownFrom).map((x) => h('div', { class: `msg ${x.who}` }, x.text,
      h('div', { class: 'meta' }, x.who === 'mate' ? '✓ done' : 'you', clock(x.at), ...(x.tools ?? []).map((t) => h('span', { class: 'tool' }, `⚙ ${t}`))))),
    pend ? h('div', { class: 'msg pending' }, pend.text, h('div', { class: 'meta' }, `Sending in ${pend.left}s`, h('button', { type: 'button', id: 'undo-chat', onclick: () => undo(`chat:${p.name}`, `chat-${p.name}`) }, 'Undo'))) : null,
    turn ? h('div', { class: 'turn', role: 'status' }, h('span', { class: 'g', 'aria-hidden': 'true' }, turn.phase === 'tool' ? '⚙' : '◐'), turn.phase === 'tool' ? `using ${turn.tool}…` : 'thinking…') : null,
  ];
  const ta = textarea(`chat-${p.name}`, { placeholder: blocked ? '' : `Message the Mate of ${p.name}`, disabled: Boolean(blocked), 'aria-label': 'Message to the Mate', onkeydown: (e) => {
    if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) (e.preventDefault(), sendChat(p.name));
  } });
  return [
    h('div', { class: 'chat', 'aria-live': 'polite' }, list),
    h('div', { class: 'composer' }, ta,
      h('div', { class: 'row' }, h('small', {}, blocked || '⌘↵ send · 5 s to undo · Enter for a new line'),
        h('button', { type: 'button', class: 'btn', disabled: !turn, title: turn ? '' : 'nothing to interrupt', onclick: () => interrupt(p.name) }, 'Interrupt'),
        h('button', { type: 'button', class: 'btn solid', disabled: Boolean(blocked) || Boolean(pend), onclick: () => sendChat(p.name) }, 'Send'))),
  ];
}
function sendChat(project) {
  const id = `chat-${project}`;
  const text = (S.drafts.get(id) ?? '').trim();
  if (!text) return;
  S.drafts.set(id, '');
  sendLater(`chat:${project}`, text, (txt) => {
    fake.chats[project].push({ who: 'captain', at: new Date().toISOString(), text: txt });
    const turn = { phase: 'think' };
    S.turns.set(project, turn);
    render();
    turn.timer = setTimeout(() => {
      Object.assign(turn, { phase: 'tool', tool: 'mate status' });
      render();
      turn.timer = setTimeout(() => {
        S.turns.delete(project);
        fake.chats[project].push({ who: 'mate', at: new Date().toISOString(), text: 'On it. PLA-61 is wiring cache eviction; PLA-58 waits for your merge.', tools: ['mate status'] });
        render();
      }, 1800);
    }, 1500);
  });
}
function interrupt(project) {
  const turn = S.turns.get(project);
  if (!turn) return;
  clearTimeout(turn.timer);
  S.turns.delete(project);
  fake.chats[project].push({ who: 'mate', at: new Date().toISOString(), text: 'Interrupted by the captain.' });
  render();
}

const VERB = { 'needs-decision': ['◆', 'decision'], blocked: ['■', 'incident'], 'wait-mate': ['⇢', 'merge'], stuck: ['!', 'incident'], working: ['▶', ''], 'agent gone': ['✕', 'incident'] };
function inboxView(p) {
  const rows = (S.inbox.get(p.name) ?? []).filter((x) => S.inboxAll || x.attention);
  return [
    h('h3', { class: 'sec' }, S.inboxAll ? 'All log' : 'Inbox', h('span', { class: 'r' }, h('span', { class: 'toggle' },
      h('button', { type: 'button', 'aria-pressed': String(!S.inboxAll), onclick: () => ((S.inboxAll = false), render()) }, 'Inbox'),
      h('button', { type: 'button', 'aria-pressed': String(S.inboxAll), onclick: () => ((S.inboxAll = true), render()) }, 'All log')))),
    rows.length ? h('ul', { class: 'list' }, rows.map((x) => {
      const [g, k] = VERB[x.verb] ?? ['·', ''];
      const t = taskOf(p.name, x.crew);
      const assignWhy = x.assigned ? 'assigned to the Mate' : p.mode === 'auto' ? 'auto: the daemon sends it' : blockedBy(p.name);
      return h('li', { class: `click k-${k} ${x.attention ? '' : 'quiet'}`, onclick: (e) => !e.target.closest('button') && select({ kind: 'crew', project: p.name, crew: x.crew, back: { kind: 'mate', project: p.name, tab: 'inbox' } }) },
        h('span', { class: 'g', 'aria-hidden': 'true' }, g),
        h('b', {}, `${x.crew} · ${x.kind === 'incident' ? `incident: ${x.verb}` : x.verb}`),
        x.verb === 'needs-decision' ? h('div', { class: 'side' }, h('button', { type: 'button', disabled: Boolean(assignWhy), onclick: () => assign(p.name, x.crew) }, x.assigned ? 'Assigned' : 'Assign'), h('small', {}, assignWhy)) : h('div', { class: 'side' }, h('small', {}, ago(x.at))),
        h('small', {}, x.text + (t ? '' : '')));
    })) : h('p', { class: 'note' }, S.inboxAll ? 'No log lines yet.' : 'Nothing in the inbox.'),
  ];
}
async function assign(project, crew) {
  S.busy.set(project, 'Assign');
  render();
  const { line } = await fake.assignToMate(project, crew);
  S.busy.delete(project);
  S.results.set(project, { kind: 'ok', text: line, at: new Date().toISOString() });
  await loadAll();
}

function appsView(p) {
  const apps = p.name === 'plantrip' ? [['Linear', 'configured', 'https://linear.app'], ['Beads', 'not configured', null]] : p.name === 'mihuv2' ? [['Linear', 'configured', 'https://linear.app'], ['Fresh', 'configured', null]] : [];
  return [apps.length ? h('ul', { class: 'list' }, apps.map(([name, st, url]) => h('li', {}, h('span', { class: 'g' }, st === 'configured' ? '●' : '○'), h('b', {}, name), h('div', { class: 'side' }, url ? h('a', { href: url, target: '_blank', rel: 'noopener' }, 'Open ↗') : h('small', {}, url === null && st === 'configured' ? 'no page to open' : 'not configured')), h('small', {}, st)))) : h('p', { class: 'note' }, 'No apps or plugins in this project.')];
}

function projectView(p) {
  const n = openCrews(p.name).length;
  const busy = S.busy.get('workspace');
  return [
    facts([['Path', p.path ?? '?'], ['Repos', (p.repos ?? []).join(', ') || 'none yet'], ['Open crews', String(n)], ['Mode', p.mode]]),
    h('h3', { class: 'sec' }, 'Project'),
    h('div', { class: 'actions' }, h('div', { class: 'act' },
      h('button', { type: 'button', class: 'danger', id: `del-${p.name}`, disabled: Boolean(blockedBy('workspace')), onclick: () => deleteProject(p.name) }, busy === `Delete ${p.name}…` ? `Running delete…` : 'Delete project…'),
      h('small', {}, blockedBy('workspace')))),
    resultEl('workspace'),
  ];
}

function deleteProject(name) {
  const p = projectOf(name);
  const mates = p.mate?.running ? 1 : 0;
  const crews = openCrews(name).length;
  confirmBox({
    title: `Delete project ${name}?`,
    scope: `${name} · ${mates} Mate and ${crews} crews stop`,
    body: `${name} leaves this workspace and its block becomes a park.`,
    list: [`Stops ${mates} Mate${mates === 1 ? '' : 's'} and ${crews} crew${crews === 1 ? '' : 's'}.`, `Data stays on disk in ${p.path}; nothing there is deleted.`, 'Adding it back later starts a new Mate session.'],
    yes: 'Delete project',
    onYes: () => {
      const at = city.blocks.get(name)?.centre.clone();
      run('workspace', `Delete ${name}…`, { scope: 'workspace', action: 'remove_project', project: name }, `mate:${name}`, (ok) => {
        if (ok) select({ kind: 'hall', tab: 'projects' }, { fly: false });
        if (ok && at) fly(at.clone().setY(0.5), at.clone().add(new THREE.Vector3(0, 16, 22)));
      });
    },
  });
}

// Create a Mate: pick the harness, then the console's onboard.
VIEWS.create = () => {
  const p = projectOf(S.sel.project);
  const pick = S.drafts.get('harness') ?? 'claude';
  const key = p.name;
  const running = S.busy.get(key);
  return {
    crumbs: [{ text: p.name, go: () => select({ kind: 'mate', project: p.name }) }, { text: 'Create mate' }],
    title: `Create the Mate of ${p.name}`,
    sub: 'It moves onto the plaza of this block and starts by reading the repo.',
    body: [
      ...bands(),
      h('div', { class: 'form' },
        h('fieldset', { class: 'opts', style: 'border:0;padding:0;margin:0' }, h('legend', { class: 'sec', style: 'margin-bottom:6px;font-size:11px;letter-spacing:.08em;color:var(--muted)' }, 'HARNESS'),
          h('div', { class: 'harness' }, ['claude', 'codex', 'pi'].map((x) => h('label', {},
            h('input', { type: 'radio', name: 'harness', value: x, checked: pick === x, onchange: () => (S.drafts.set('harness', x), render()) }),
            logo(x), h('b', {}, HARNESSES[x]?.name ?? x), h('small', { class: 'hint' }, { claude: 'Claude Code', codex: 'Codex CLI', pi: 'pi' }[x]))))),
        h('div', { class: 'act' }, h('button', { type: 'button', class: running ? 'running' : 'primary', id: 'create-go', disabled: Boolean(running) || Boolean(blockedBy(key)), onclick: () => run(key, 'Create mate…', { action: 'onboard', project: p.name, harness: pick }, `mate:${p.name}`, (ok) => ok && select({ kind: 'mate', project: p.name }, { fly: true })) },
          running ? h('span', { class: 'g' }, '⟳') : null, running ? 'Running create…' : `Create mate on ${HARNESSES[pick]?.name ?? pick}`), h('small', {}, running ? 'reading the repo can take a minute' : blockedBy(key))),
        resultEl(key)),
    ],
  };
};
function logo(harness) {
  const hn = HARNESSES[harness];
  const fill = Array.isArray(hn?.logoColor) ? hn.logoColor[1] : hn?.logoColor ?? '#8E857A';
  const svg = (hn?.logo ?? '').replace('<svg ', `<svg fill="${fill}" `);
  return h('img', { src: `data:image/svg+xml,${enc(svg)}`, alt: '' });
}

// A crew, where it stands.
VIEWS.crew = () => {
  const { project, crew } = S.sel;
  const t = taskOf(project, crew);
  const p = projectOf(project);
  if (!t || !p) return { title: crew, body: [h('p', {}, 'This crew is no longer in the workspace.'), resultEl(rowKey(project, crew))] };
  if (t.closed) {
    return {
      crumbs: [{ text: project, go: () => select({ kind: 'mate', project }) }, { text: 'crew' }],
      title: `${crew} · ${t.role}`,
      sub: t.text,
      body: [
        facts([['State', stateEl('', t.merged_at ? ['✓', `merged ${clock(t.merged_at)}`, 'ok'] : ['◌', (t.close_state || 'closed').replace(/_/g, ' '), 'off'])], ['Branch', t.branch ?? '?'], ['Tokens', t.tokens?.total ? formatTokens(t.tokens.total) : '?'], ['Closed', clock(t.closed_at)]]),
        resultEl(rowKey(project, crew)),
        h('p', { class: 'note' }, 'This crew is closed, so it has no actions. Its story stays in the tent\'s playback and in the Office cabinet.'),
        h('div', { class: 'links' }, h('a', { href: '#', onclick: (e) => (e.preventDefault(), enterRoom(`${project}/${t.repo}`)) }, `Step into tent ${t.repo}`), h('a', { href: '#', onclick: (e) => e.preventDefault() }, 'Timeline ↗')),
      ],
    };
  }
  const acts = S.actions.get(project)?.crews?.[crew] ?? [];
  const tab = S.sel.tab ?? (t.incident ? 'peek' : t.status === 'needs-decision' ? 'status' : 'status');
  const kit = [harnessName(t.harness), t.model, t.effort].filter(Boolean).join(' · ');
  const buttons = acts.filter((b) => b.action !== 'diff');
  const diffB = acts.find((b) => b.action === 'diff');
  const extra = [
    diffB ? { ...diffB, toggle: true, pressed: tab === 'diff', run: () => setTab(tab === 'diff' ? 'status' : 'diff') } : null,
    { action: 'peek', label: 'Peek', available: true, toggle: true, pressed: tab === 'peek', run: () => setTab(tab === 'peek' ? 'status' : 'peek') },
  ].filter(Boolean);
  const replay = room && mode !== 'live';
  return {
    crumbs: [{ text: project, go: () => select({ kind: 'mate', project }) }, ...(S.sel.back?.kind === 'tent' ? [{ text: `tent ${S.sel.back.repo}`, go: () => select(S.sel.back) }] : []), { text: 'crew' }],
    title: `${crew} · ${t.role}`,
    sub: t.text,
    body: [
      ...bands(replay ? [h('div', { class: 'band replay' }, h('span', { class: 'g' }, '⏸'), h('span', {}, `The tent shows ${$('#when').textContent}; this panel shows now. Actions wait for Live. `), h('button', { type: 'button', class: 'btn', style: 'min-height:32px', onclick: goLive }, 'Go Live'))] : []),
      facts([
        ['State', t.incident ? stateEl('', ['!', `stuck ${ago(t.since)}`, 'warn']) : stateEl(t.state)],
        ['Status line', t.status ? `${t.status}${t.assigned ? ' · assigned' : ''}` : '?'],
        ['Harness', kit || '?'],
        ['Branch', t.branch ?? '?'],
        ['Tokens', t.tokens?.total ? formatTokens(t.tokens.total) : '?'],
        ['Since', t.since ? `${clock(t.since)} · ${ago(t.since)}` : '?'],
      ]),
      ...actionsBlock({ project, crew, buttons, extra, links: [h('a', { href: '#', onclick: (e) => e.preventDefault(), title: 'opens the admin timeline in a new tab' }, 'Timeline ↗'), h('a', { href: '#', onclick: (e) => e.preventDefault(), title: 'opens the admin usage page in a new tab' }, 'Usage ↗')] }),
      tabs([['status', 'Status'], ['peek', 'Peek'], ['diff', 'Diff']], tab, setTab),
      ...{ status: crewStatus, peek: crewPeek, diff: crewDiff }[tab](p, t),
    ],
  };
};
function crewStatus(p, t) {
  const q = t.status === 'needs-decision';
  return [
    h('div', { class: `result`, dataset: { kind: 'info' } }, h('span', { class: 'g' }, VERB[t.status]?.[0] ?? '·'), h('span', {}, h('b', {}, `${t.status ?? '?'}: `), t.status_text ?? '?')),
    q ? h('div', { class: 'form' },
      h('p', { class: 'note' }, p.mode === 'manual' ? `${p.name} is manual: the Mate does not see this until someone gives it. Assign hands the question to the Mate.` : 'Auto: the daemon types this into the Mate.'),
      h('div', { class: 'act' }, h('button', { type: 'button', class: 'btn', disabled: Boolean(t.assigned || p.mode === 'auto' || blockedBy(p.name)), onclick: () => assign(p.name, t.crew) }, t.assigned ? 'Assigned to the Mate' : 'Assign to the Mate'), h('small', {}, t.assigned ? 'the Mate has it' : p.mode === 'auto' ? 'auto: the daemon sends it' : blockedBy(p.name))),
      resultEl(p.name)) : null,
  ];
}
function crewPeek(p, t) {
  return [
    h('div', { class: 'peekmeta' }, h('span', {}, '◉ read-only · captured 12 s ago'), h('button', { type: 'button', class: 'btn', style: 'min-height:32px', onclick: () => render() }, 'Refresh')),
    h('pre', { class: 'peek', 'aria-label': `Screen of ${t.crew}` }, fake.peekOf(p.name, t.crew)),
    h('p', { class: 'note' }, 'No terminal in the browser. To type into this agent, use mate console or ask the Mate.'),
  ];
}
function crewDiff(p, t) {
  const d = fake.diffOf(p.name, t.crew);
  if (!d) return [h('p', {}, '?')];
  const add = d.files.reduce((n, f) => n + f.add, 0);
  const del = d.files.reduce((n, f) => n + f.del, 0);
  return [
    h('div', { class: 'peekmeta' }, h('span', { style: 'font-family:var(--mono)' }, `${d.branch} → ${d.into} · ${d.ahead} ahead, ${d.behind} behind`)),
    h('table', { class: 'files' }, h('tbody', {}, d.files.map((f) => h('tr', {}, h('td', {}, f.path), h('td', { class: 'n add' }, `+${f.add}`), h('td', { class: 'n del' }, `−${f.del}`))),
      h('tr', {}, h('td', {}, h('b', {}, `${d.files.length} files`)), h('td', { class: 'n add' }, `+${add}`), h('td', { class: 'n del' }, `−${del}`)))),
    h('pre', { class: 'hunk' }, d.hunk),
    h('div', { class: 'links' }, h('a', { href: '#', onclick: (e) => e.preventDefault() }, 'Whole diff in Office ↗')),
  ];
}

// City Hall: topics from every Mate, the projects, adding one.
VIEWS.hall = () => {
  const tab = S.sel.tab ?? 'topics';
  const waiting = fake.topics.filter((x) => x.status === 'waiting' && projectOf(x.project));
  if (S.sel.topic) return topicView(fake.topics.find((x) => x.id === S.sel.topic));
  return {
    crumbs: [{ text: 'City' }, { text: 'City Hall' }],
    title: 'City Hall',
    sub: 'Your desk: what the Mates ask you, and the projects of this workspace',
    body: [
      ...bands(),
      tabs([['topics', 'Topics', waiting.length], ['projects', 'Projects']], tab, setTab),
      ...(tab === 'topics' ? topicList() : projectList()),
    ],
  };
};
function topicList() {
  const groups = [['waiting', 'Waiting on you'], ['on_it', 'The Mate is on it'], ['earlier', 'Earlier']];
  return groups.flatMap(([st, label]) => {
    const rows = fake.topics.filter((x) => x.status === st && projectOf(x.project));
    if (!rows.length) return [];
    return [h('div', { class: 'group-h' }, `${label} · ${rows.length}`), h('ul', { class: 'list' }, rows.map((x) => h('li', { class: `click ${st === 'waiting' ? 'k-topic' : 'quiet'}`, tabindex: '0', role: 'button', onclick: () => select({ ...S.sel, topic: x.id }, { fly: false }), onkeydown: (e) => e.key === 'Enter' && select({ ...S.sel, topic: x.id }, { fly: false }) },
      h('span', { class: 'g', 'aria-hidden': 'true' }, st === 'waiting' ? '?' : st === 'on_it' ? '▶' : '✓'), h('b', {}, x.question), h('div', { class: 'side' }, h('small', {}, ago(x.at))), h('small', {}, `${x.project}${x.answer ? ` · you said: ${x.answer}` : ''}`))))];
  });
}
function topicView(x) {
  const pend = S.pending.get(`topic:${x.id}`);
  const pick = S.drafts.get(`opt-${x.id}`);
  const own = (S.drafts.get(`own-${x.id}`) ?? '').trim();
  const answer = own || pick;
  const fork = S.forks.get(x.id);
  const why = S.stale ? `stale since ${clock(S.staleSince)}` : S.actions.get(x.project)?.read_only ? 'read-only dashboard' : !answer ? 'pick an option or type an answer' : '';
  const st = { waiting: ['?', 'Waiting on you', 'need'], on_it: ['▶', 'The Mate is on it', 'ok'], earlier: ['✓', 'Earlier', 'off'] }[x.status];
  return {
    crumbs: [{ text: 'City Hall', go: () => select({ kind: 'hall', tab: 'topics' }, { fly: false }) }, { text: 'topic' }],
    title: x.project,
    sub: `asked ${ago(x.at)} ago by the Mate of ${x.project}`,
    body: [
      ...bands(),
      stateEl('', st),
      h('p', { class: 'topic-q' }, x.question),
      x.context ? h('p', { class: 'topic-ctx' }, x.context) : null,
      x.notes ? h('p', { class: 'note' }, `The Mate: ${x.notes}`) : null,
      x.status !== 'waiting' ? h('div', { class: 'result', dataset: { kind: 'info' } }, h('span', { class: 'g' }, '✓'), h('span', {}, `You answered: ${x.answer}`)) : h('div', { class: 'form' },
        h('fieldset', { class: 'opts' }, h('legend', { class: 'sr' }, 'Options'), x.options.map((o, i) => h('label', {}, h('input', { type: 'radio', name: `opt-${x.id}`, id: `opt-${x.id}-${i}`, checked: pick === o, onchange: () => (S.drafts.set(`opt-${x.id}`, o), render()) }), o))),
        h('label', {}, 'Or your own answer', textarea(`own-${x.id}`, { rows: 2, placeholder: 'Type to answer in your words', oninput: (e) => (S.drafts.set(`own-${x.id}`, e.target.value), clearTimeout(topicView.t), (topicView.t = setTimeout(render, 250))) })),
        pend ? h('div', { class: 'msg pending' }, pend.text, h('div', { class: 'meta' }, `Sending in ${pend.left}s`, h('button', { type: 'button', id: 'undo-topic', onclick: () => undo(`topic:${x.id}`, `own-${x.id}`) }, 'Undo'))) : null,
        h('div', { class: 'actions' },
          h('div', { class: 'act' }, h('button', { type: 'button', class: 'primary', id: `send-${x.id}`, disabled: Boolean(why) || Boolean(pend), onclick: () => sendLater(`topic:${x.id}`, answer, async (a) => {
            S.drafts.delete(`own-${x.id}`);
            const { line } = await fake.answerTopic(x.id, a);
            S.results.set(`topic:${x.id}`, { kind: 'ok', text: line, at: new Date().toISOString() });
            await loadAll();
          }) }, 'Send answer'), h('small', {}, why)),
          h('div', { class: 'act' }, h('button', { type: 'button', id: `fork-${x.id}`, disabled: Boolean(fork), onclick: () => openFork(x) }, 'Discuss privately'), h('small', {}, fork ? 'session open below' : 'a read-only fork of the Mate'))),
        fork ? forkView(x, fork) : null),
      resultEl(`topic:${x.id}`),
    ],
  };
}
function openFork(x) {
  const fork = { messages: [], proposal: null, thinking: true };
  S.forks.set(x.id, fork);
  render();
  setTimeout(() => {
    fork.thinking = false;
    fork.messages.push({ who: 'mate', text: x.project === 'plantrip' ? 'The free tier cap is 50k requests a day; last week peaked at 31k. A flag keeps 2.4 on time and lets you turn it on per country.' : 'The lane is the bottleneck: HU-131 waits behind HU-129. Restarting frees the lane in about 2 minutes; the last hang cleared on restart.' });
    fork.proposal = x.project === 'plantrip' ? 'Ship behind a flag, off by default; turn it on for VN first and watch the tile requests.' : 'Restart HU-129 in its worktree; if it hangs again, stop it and give the lane to HU-131.';
    render();
  }, 1300);
}
function forkView(x, fork) {
  return h('div', { class: 'fork', role: 'region', 'aria-label': 'Private session' },
    h('h4', {}, '◉ Private session · read-only fork of the Mate · nothing here reaches it'),
    ...fork.messages.map((m) => h('div', { class: `msg ${m.who}` }, m.text)),
    fork.thinking ? h('div', { class: 'turn' }, h('span', { class: 'g' }, '◐'), 'thinking…') : null,
    fork.proposal ? h('div', { class: 'proposal' }, h('b', {}, 'Proposed answer'), h('p', { style: 'margin:4px 0 8px' }, fork.proposal),
      h('button', { type: 'button', class: 'btn', onclick: () => {
        S.drafts.set(`own-${x.id}`, fork.proposal);
        render();
        const el = document.getElementById(`own-${x.id}`);
        el?.focus();
      } }, 'Use this answer (you can edit it)')) : null,
    h('div', { class: 'composer', style: 'position:static;margin:0;padding:0;border:0;background:none' }, textarea(`fork-in-${x.id}`, { rows: 2, placeholder: 'Ask the fork', 'aria-label': 'Ask the private session' }),
      h('div', { class: 'row' }, h('small', {}, 'the fork reads; it cannot act'), h('button', { type: 'button', class: 'btn', onclick: () => S.forks.delete(x.id) || render() }, 'Close session'), h('button', { type: 'button', class: 'btn solid', onclick: () => {
        const q = (S.drafts.get(`fork-in-${x.id}`) ?? '').trim();
        if (!q) return;
        S.drafts.set(`fork-in-${x.id}`, '');
        fork.messages.push({ who: 'captain', text: q });
        fork.thinking = true;
        render();
        setTimeout(() => ((fork.thinking = false), fork.messages.push({ who: 'mate', text: 'From the repo: the flag would live in RemoteConfig; no release needed to turn it on.' }), render()), 1200);
      } }, 'Ask'))));
}

function projectList() {
  return [
    h('ul', { class: 'list' }, (S.ws?.projects ?? []).map((p) => {
      const n = openCrews(p.name).length;
      const g = p.mate.created === false ? '+' : p.mate.running ? '●' : '◌';
      return h('li', { class: 'click', tabindex: '0', role: 'button', onclick: (e) => !e.target.closest('button') && select({ kind: 'mate', project: p.name }), onkeydown: (e) => e.key === 'Enter' && e.target === e.currentTarget && select({ kind: 'mate', project: p.name }) },
        h('span', { class: 'g' }, g), h('b', {}, p.name),
        h('div', { class: 'side' }, h('button', { type: 'button', onclick: () => deleteProject(p.name), disabled: Boolean(blockedBy('workspace')) }, 'Delete…')),
        h('small', {}, `${p.mate.created === false ? 'no Mate yet' : p.mate.running ? `Mate ${STATE[p.mate.state]?.[1] ?? p.mate.state}` : 'Mate stopped'} · ${n} crew${n === 1 ? '' : 's'} · ${unk(p.mate.tokens_today, formatTokens)} today · ${p.mode}`));
    })),
    h('div', { class: 'act' }, h('button', { type: 'button', class: 'btn', id: 'add-open', onclick: () => select({ kind: 'add', back: { kind: 'hall', tab: 'projects' } }, { fly: false }), disabled: Boolean(blockedBy('workspace')) }, '+ Add project…'), h('small', {}, blockedBy('workspace') || (city && city.hidden ? '' : 'it takes the next empty lot'))),
    resultEl('workspace'),
  ];
}

// Add a project: a park lot becomes a block.
VIEWS.add = VIEWS.lot = () => {
  const name = (S.drafts.get('add-name') ?? '').trim();
  const path = (S.drafts.get('add-path') ?? '').trim();
  const bad = name && !/^[a-z0-9][a-z0-9-]{0,30}$/.test(name) ? 'lowercase letters, digits and -' : projectOf(name) ? `${name} is already a project` : '';
  const running = S.busy.get('workspace');
  const full = (S.ws?.projects.length ?? 0) >= 24;
  const why = blockedBy('workspace') || (full ? 'the workspace holds 24 projects' : !name ? 'a name is needed' : bad);
  return {
    crumbs: [{ text: 'City Hall', go: () => select({ kind: 'hall', tab: 'projects' }) }, { text: S.sel.kind === 'lot' ? 'empty lot' : 'add project' }],
    title: S.sel.kind === 'lot' ? 'Empty lot' : 'Add project',
    sub: 'A park becomes a block with a tent per repo and a plaza for its Mate.',
    body: [
      ...bands(),
      h('div', { class: 'form' },
        h('label', {}, 'Name', h('input', { type: 'text', id: 'add-name', value: S.drafts.get('add-name') ?? '', autocomplete: 'off', oninput: (e) => (S.drafts.set('add-name', e.target.value), render()) }), bad ? h('span', { class: 'err' }, bad) : h('span', { class: 'hint' }, 'shown on the block; lowercase')),
        h('label', {}, 'Repo path (optional)', h('input', { type: 'text', id: 'add-path', value: S.drafts.get('add-path') ?? '', placeholder: '~/code/my-repo', autocomplete: 'off', oninput: (e) => S.drafts.set('add-path', e.target.value) }), h('span', { class: 'hint' }, 'leave empty to add the repo later')),
        h('div', { class: 'act' }, h('button', { type: 'button', id: 'add-go', class: running ? 'running' : 'primary', disabled: Boolean(running) || Boolean(why), onclick: () => run('workspace', 'Add project…', { scope: 'workspace', action: 'add_project', name, path }, 'hall', async (ok) => {
          if (!ok) return;
          S.drafts.delete('add-name');
          S.drafts.delete('add-path');
          await loadAll();
          flash(`block:${name}`, 'new', '+', `new · ${name}`, 7000);
          select({ kind: 'mate', project: name, tab: 'project' });
          const b = city.blocks.get(name);
          if (b) fly(b.centre.clone().setY(0.5), b.centre.clone().add(new THREE.Vector3(0, 16, 22)));
        }) }, running ? h('span', { class: 'g' }, '⟳') : null, running ? 'Running add…' : 'Add project'), h('small', {}, running ? 'clearing the lot…' : why)),
        resultEl('workspace')),
    ],
  };
};

VIEWS.block = () => {
  const p = projectOf(S.sel.project);
  if (!p) return { title: '?', body: [] };
  return {
    crumbs: [{ text: 'City' }, { text: p.name }],
    title: p.name,
    sub: `${(p.repos ?? []).length} tent${p.repos?.length === 1 ? '' : 's'} · ${openCrews(p.name).length} crews at work`,
    body: [
      h('ul', { class: 'list' },
        h('li', { class: 'click', onclick: () => select({ kind: 'mate', project: p.name }) }, h('span', { class: 'g' }, '◎'), h('b', {}, 'The Mate on the plaza'), h('small', {}, 'chat, inbox, actions')),
        ...tentRepos(p).map((r) => h('li', { class: 'click', onclick: () => enterRoom(`${p.name}/${r}`) }, h('span', { class: 'g' }, '⌂'), h('b', {}, `Tent ${r}`), h('small', {}, 'step in: playback and its crews'))),
        ...openCrews(p.name).map((t) => h('li', { class: 'click', onclick: () => select({ kind: 'crew', project: p.name, crew: t.crew }) }, h('span', { class: 'g' }, STATE[t.state]?.[0] ?? '·'), h('b', {}, t.crew), h('small', {}, `${t.role} · ${STATE[t.state]?.[1] ?? t.state}`)))),
    ],
  };
};

// The system: is what you see current, and what ran without you.
VIEWS.system = () => {
  const rows = [
    ['Dashboard', S.stale ? ['▲', `stale since ${clock(S.staleSince)}: ${S.stale}`, 'warn'] : ['●', 'live', 'ok']],
    ['Timeline', S.timelineStuck ? ['▲', 'stuck · mihuv2 since 14:02 (git timed out)', 'warn'] : ['●', 'current · every project', 'ok']],
    ['Herdr', S.herdrDown ? ['▲', 'not answering since 14:20; Peek and Restart fail until it does', 'warn'] : ['●', 'answering', 'ok']],
  ];
  return {
    crumbs: [{ text: 'City' }, { text: 'System' }],
    title: 'System',
    sub: 'Whether the city is current, and what the auto daemon did',
    body: [
      h('ul', { class: 'list' }, rows.map(([k, s]) => h('li', {}, h('span', { class: 'g' }, s[0]), h('b', {}, k), h('small', {}, s[1])))),
      h('h3', { class: 'sec' }, 'Auto daemon'),
      h('ul', { class: 'list' }, fake.daemonLog.map((x) => h('li', { class: /did not/.test(x.text) ? 'quiet' : '' }, h('span', { class: 'g' }, /did not/.test(x.text) ? '⊘' : '→'), h('b', {}, x.project), h('div', { class: 'side' }, h('small', {}, clock(x.at))), h('small', {}, x.text)))),
    ],
  };
};
function renderSystemChip() {
  const bad = S.stale || S.timelineStuck || S.herdrDown;
  $('#sys').setAttribute('aria-label', S.stale ? `System: stale, ${S.stale}` : bad ? 'System: needs a look' : 'System: live');
  $('#sys').dataset.tone = bad ? 'warn' : 'ok';
  $('#sys .g').textContent = bad ? '▲' : '●';
  $('#sys-t').textContent = S.timelineStuck ? 'Timeline stuck' : S.herdrDown ? 'Herdr down' : 'System';
}
$('#sys').addEventListener('click', () => select({ kind: 'system' }));

// Go to: every place by name, Needs you first. On a phone this is how you get to a crew.
VIEWS.goto = () => {
  const q = (S.drafts.get('goto-q') ?? '').toLowerCase();
  const items = [
    ...attention().map((x) => ({ g: KINDS[x.kind].glyph, k: x.kind, title: x.title, sub: x.sub, go: x.open })),
    { g: '⌂', title: 'City Hall', sub: 'topics · projects', go: () => select({ kind: 'hall' }) },
    ...(S.ws?.projects ?? []).flatMap((p) => [
      { g: '◎', title: `Mate · ${p.name}`, sub: p.mate.created === false ? 'no Mate yet' : harnessName(p.mate.harness), go: () => select({ kind: 'mate', project: p.name }) },
      ...openCrews(p.name).map((t) => ({ g: STATE[t.state]?.[0] ?? '·', title: t.crew, sub: `${p.name} · ${t.role} · ${STATE[t.state]?.[1] ?? t.state}`, go: () => select({ kind: 'crew', project: p.name, crew: t.crew }) })),
      ...tentRepos(p).map((r) => ({ g: '⌂', title: `Tent ${r}`, sub: p.name, go: () => enterRoom(`${p.name}/${r}`) })),
    ]),
  ].filter((x) => !q || `${x.title} ${x.sub}`.toLowerCase().includes(q));
  const input = h('input', { type: 'text', id: 'goto-q', value: S.drafts.get('goto-q') ?? '', placeholder: 'crew, project, tent…', 'aria-label': 'Go to', autocomplete: 'off', oninput: (e) => (S.drafts.set('goto-q', e.target.value), render()), onkeydown: (e) => e.key === 'Enter' && items[0]?.go() });
  return {
    crumbs: [{ text: 'City' }, { text: 'Go to' }],
    title: 'Go to',
    sub: 'The camera flies there; no need to turn it',
    body: [h('div', { class: 'form' }, input), h('ul', { class: 'list' }, items.map((x) => h('li', { class: `click k-${x.k ?? ''}`, tabindex: '0', role: 'button', onclick: x.go, onkeydown: (e) => e.key === 'Enter' && x.go() }, h('span', { class: 'g' }, x.g), h('b', {}, x.title), h('small', {}, x.sub))))],
  };
};
$('#goto').addEventListener('click', () => select({ kind: 'goto' }, { fly: false }));
let nextI = 0;
function nextAttention() {
  const items = attention();
  if (!items.length) return;
  items[nextI++ % items.length].open();
}
$('#needs').addEventListener('click', () => (phone() ? select({ kind: 'goto' }, { fly: false }) : nextAttention()));

// ---- Inside a tent (the real room, with the dock beside it) -----------------------------

let room = null;
let roomT = 0;
let cityView = null;
let roomToken = 0;
const roomScene = new THREE.Scene();
roomScene.background = new THREE.Color(PALETTE.paper);
lights(roomScene, renderer);
const feedEl = h('ol', { id: 'feed' });
const KIND = { spawn: 'spawn', status: 'status', ask: 'ask', answer: 'answer', handback: 'hand back', review: 'review', verdict: 'verdict', merge: 'merge', close: 'close' };

VIEWS.tent = () => {
  const { project, repo } = S.sel;
  const tab = S.sel.tab ?? 'feed';
  const crews = openCrews(project).filter(onTent(projectOf(project) ?? { repos: [] }, repo));
  return {
    crumbs: [{ text: project, go: () => (leaveRoom(), select({ kind: 'block', project })) }, { text: 'tent' }],
    title: repo,
    sub: room ? `${project} · ${room.summary}` : `${project} · loading…`,
    body: [
      ...bands(mode !== 'live' ? [h('div', { class: 'band replay' }, h('span', { class: 'g' }, '⏸'), h('span', {}, `Replaying ${$('#when').textContent}. Crews and actions below are now. `), h('button', { type: 'button', class: 'btn', style: 'min-height:32px', onclick: goLive }, 'Go Live'))] : []),
      tabs([['feed', 'Playback'], ['crews', 'Crews here', crews.length]], tab, setTab),
      tab === 'feed' ? feedEl : h('ul', { class: 'list' }, crews.length ? crews.map((t) => h('li', { class: 'click', tabindex: '0', role: 'button', onclick: () => select({ kind: 'crew', project, crew: t.crew, back: { ...S.sel } }), onkeydown: (e) => e.key === 'Enter' && select({ kind: 'crew', project, crew: t.crew, back: { ...S.sel } }) },
        h('span', { class: 'g' }, t.incident ? '!' : STATE[t.state]?.[0] ?? '·'), h('b', {}, `${t.crew} · ${t.role}`), h('small', {}, t.incident ? `stuck ${ago(t.since)}` : STATE[t.state]?.[1] ?? t.state))) : h('li', {}, h('span', { class: 'g' }, '○'), h('b', {}, 'Nobody is working here now'))),
    ],
  };
};

async function enterRoom(key) {
  const tent = city?.tents.get(key);
  if (!tent) return;
  const token = ++roomToken;
  if (!room) cityView = { target: controls.target.clone(), position: camera.position.clone() };
  document.body.classList.add('in-room');
  S.sel = { kind: 'tent', project: tent.project, repo: tent.repo, tab: 'feed' };
  document.body.classList.add('docked');
  ring.visible = false;
  render();
  const next = await buildRoom(key, token);
  if (!next) return;
  setMode('live');
  roomT = next.events.at(-3)?.at ?? 0;
  flight = null;
  controls.target.copy(next.view.target);
  camera.position.copy(next.view.position);
  Object.assign(controls, ROOM_LIMITS);
  resize();
  render();
}

async function buildRoom(key, token) {
  const tent = city.tents.get(key);
  const p = projectOf(tent.project);
  const h0 = await readHistory(fake.api, tent.project, S.tasks.get(tent.project) ?? [], onTent(p, tent.repo));
  if (token !== roomToken) return null;
  if (room) {
    roomScene.remove(room.group);
    room.dispose();
  }
  room = createTentRoom({ project: tent.project, repo: tent.repo, floor: city.blocks.get(tent.project).floor, mateHarness: p.mate?.harness, mateState: p.mate?.running ? p.mate.state : 'asleep', history: h0 });
  room.key = key;
  roomScene.add(room.group);
  shown = -1;
  feedEl.replaceChildren(...room.events.map((ev, i) => {
    const li = h('li', { dataset: { i: String(i) }, onclick: () => {
      roomT = room.events[i].at;
      if (mode === 'live') setMode('play');
    } }, h('i', {}, KIND[ev.kind] ?? ev.kind), h('b', {}, ev.clock ? `${ev.who} · ${ev.clock}` : ev.who), h('span', {}, ev.text));
    return li;
  }));
  drawTicks();
  return room;
}

let refreshTimer = 0;
function refreshRoom() {
  clearTimeout(refreshTimer);
  refreshTimer = setTimeout(async () => {
    if (!room) return;
    const caughtUp = mode === 'live';
    const next = await buildRoom(room.key, roomToken);
    if (next && caughtUp) roomT = next.duration;
  }, 600);
}

function leaveRoom() {
  if (!room && !document.body.classList.contains('in-room')) return;
  roomToken++;
  document.body.classList.remove('in-room');
  if (room) {
    roomScene.remove(room.group);
    room.dispose();
    room = null;
  }
  Object.assign(controls, CITY_LIMITS, { maxDistance: 170 * ((city?.size ?? 3) / 3) });
  if (cityView) {
    controls.target.copy(cityView.target);
    camera.position.copy(cityView.position);
    cityView = null;
  }
  S.sel = null;
  document.body.classList.remove('docked');
  resize();
  render();
}

// The scrub bar: the real one. Acting is for Live only (design, question 3).
const seek = $('#seek');
const playBtn = $('#play');
const speedBtn = $('#speed');
let mode = 'live';
let speed = 1;
const SPEEDS = [1, 4, 16, 64];
let dragging = null;
function setMode(next) {
  const was = mode;
  mode = next;
  playBtn.textContent = mode === 'pause' ? '▶' : '❚❚';
  playBtn.setAttribute('aria-label', mode === 'pause' ? 'Play' : 'Pause');
  $('#golive').classList.toggle('on', mode === 'live');
  if ((was === 'live') !== (mode === 'live')) render();
}
function goLive() {
  if (room) roomT = room.duration;
  setMode('live');
}
playBtn.addEventListener('click', () => {
  if (!room) return;
  if (mode !== 'pause') return setMode('pause');
  if (roomT >= room.duration) roomT = 0;
  setMode('play');
});
speedBtn.addEventListener('click', () => {
  speed = SPEEDS[(SPEEDS.indexOf(speed) + 1) % SPEEDS.length];
  speedBtn.textContent = `${speed}×`;
  if (mode === 'live' && room && roomT < room.duration) setMode('play');
});
$('#golive').addEventListener('click', goLive);
seek.addEventListener('input', () => {
  if (!room) return;
  dragging ??= mode;
  setMode('pause');
  roomT = Number(seek.value) * room.duration;
});
seek.addEventListener('change', () => {
  setMode(dragging === 'pause' || Number(seek.value) < 0.999 ? 'pause' : 'live');
  dragging = null;
});
function drawTicks() {
  const marks = [];
  let day = '';
  for (const ev of room.events) {
    if (!ev.iso) continue;
    const pos = `${(100 * ev.at) / room.duration}%`;
    const d = ev.clock.slice(0, 5);
    if (d !== day) {
      day = d;
      marks.push(h('span', { class: 'day', style: `left:${pos}` }, d));
    }
    if (ev.kind === 'merge') marks.push(h('i', { class: 'merge', style: `left:${pos}`, title: ev.text }));
  }
  $('#ticks').replaceChildren(...marks);
}

// ---- The merge crate: carried from the crew to the tent's sign -------------------------

function startCrate(project, crew) {
  const m = city?.mascots.get(`${project}/${crew}`);
  const t = taskOf(project, crew);
  const tent = city?.tents.get(`${project}/${t?.repo}`);
  if (!m || !tent) return;
  const from = m.root.position.clone().setY(1.3);
  const to = new THREE.Vector3(tent.x - tent.w / 2 + 0.7, 1.1, tent.z + tent.d / 2 + 0.45);
  crateFlight = { from, to, t0: performance.now(), hold: true };
  crate.visible = true;
}
function crateBack() {
  if (crateFlight) Object.assign(crateFlight, { back: true, t0: performance.now() });
}

// ---- Pointer -----------------------------------------------------------------------------

const ray = new THREE.Raycaster();
const pointer = new THREE.Vector2();
const ground = new THREE.Plane(new THREE.Vector3(0, 1, 0), 0);
function pick(e) {
  const r = canvas.getBoundingClientRect();
  pointer.set(((e.clientX - r.left) / r.width) * 2 - 1, -((e.clientY - r.top) / r.height) * 2 + 1);
  ray.setFromCamera(pointer, camera);
  const list = [...(room ? room.mascots : city?.mascots ?? new Map()).values()].filter((m) => m.root.visible);
  const hit = ray.intersectObjects(list.map((m) => m.body))[0];
  if (hit) return { mascot: list.find((m) => m.body === hit.object).id };
  if (room || !city) return null;
  let o = ray.intersectObjects([...[...city.tents.values()].map((b) => b.group), ...city.hall], true)[0]?.object;
  while (o && !o.userData.tent && !city.hall.includes(o)) o = o.parent;
  if (o) return o.userData.tent ? { tent: o.userData.tent } : { hall: true };
  // An empty park lot: a place to add a project.
  const at = ray.ray.intersectPlane(ground, new THREE.Vector3());
  if (!at) return null;
  const pitch = 14;
  const slot = [Math.round(at.x / pitch), Math.round(at.z / pitch)];
  const half = (city.size - 1) / 2;
  if (Math.abs(slot[0]) > half || Math.abs(slot[1]) > half || (slot[0] === 0 && slot[1] === 0)) return null;
  for (const [name, b] of city.blocks) if (Math.round(b.centre.x / pitch) === slot[0] && Math.round(b.centre.z / pitch) === slot[1]) return { block: name };
  return { lot: new THREE.Vector3(slot[0] * pitch, 0.5, slot[1] * pitch) };
}
const tip = h('div', { id: 'tip', style: 'position:fixed;pointer-events:none;background:var(--card);border:1px solid var(--line);border-radius:10px;padding:6px 10px;box-shadow:0 6px 20px rgba(36,33,29,.12);transform:translate(-50%,calc(-100% - 14px));display:none;z-index:4;max-width:min(380px,calc(100vw - 32px))' });
document.body.append(tip);
function describe(tg) {
  const seats = room ? room.seats : city.seats;
  if (tg.mascot) {
    const s = seats.get(tg.mascot);
    if (!s) return null;
    if (s.kind === 'mate') return [`Mate · ${s.project}`, `${harnessName(s.harness)} · click for chat and actions`];
    return [`${s.crew} · ${s.role}`, `${[harnessName(s.harness), s.model, s.effort].filter(Boolean).join(' · ')} · click for actions`];
  }
  if (tg.tent) {
    const b = city.tents.get(tg.tent);
    return [`${b.repo} · ${b.project}`, `${b.merged ?? '?'} merged · ${b.crews.size} at work · click to step in`];
  }
  if (tg.hall) return ['City Hall', 'topics from the Mates · projects · click to open'];
  if (tg.block) return [tg.block, 'click for its Mate, tents and crews'];
  if (tg.lot) return ['Empty lot', 'click to add a project here'];
  return null;
}
canvas.addEventListener('pointermove', (e) => {
  const tg = pick(e);
  const words = tg && describe(tg);
  canvas.style.cursor = words ? 'pointer' : '';
  if (words && !phone()) {
    tip.replaceChildren(h('b', { style: 'display:block' }, words[0]), h('span', { style: 'color:var(--muted)' }, words[1]));
    tip.style.display = 'block';
    tip.style.left = `${e.clientX}px`;
    tip.style.top = `${e.clientY}px`;
  } else tip.style.display = 'none';
});
canvas.addEventListener('pointerleave', () => (tip.style.display = 'none'));
let down = null;
canvas.addEventListener('pointerdown', (e) => (down = [e.clientX, e.clientY]));
canvas.addEventListener('click', (e) => {
  if (down && Math.hypot(e.clientX - down[0], e.clientY - down[1]) > 5) return;
  const tg = pick(e);
  if (!tg) return;
  tip.style.display = 'none';
  if (tg.tent) return enterRoom(tg.tent);
  if (tg.hall) return select({ kind: 'hall' });
  if (tg.block) return select({ kind: 'block', project: tg.block });
  if (tg.lot) return select({ kind: 'lot', at: tg.lot });
  const s = (room ? room.seats : city.seats).get(tg.mascot);
  if (!s) return;
  if (s.kind === 'mate') return room ? null : select({ kind: 'mate', project: s.project });
  select({ kind: 'crew', project: s.project, crew: s.crew, back: room ? { ...S.sel, kind: 'tent' } : null }, { fly: !room });
});

// ---- Keyboard -----------------------------------------------------------------------------

addEventListener('keydown', (e) => {
  if (scrim.classList.contains('on')) return;
  const typing = /^(INPUT|TEXTAREA|SELECT)$/.test(e.target.tagName);
  if (e.key === 'Escape') {
    if (typing) return e.target.blur();
    if (S.sel?.topic) return select({ kind: 'hall', tab: 'topics' }, { fly: false });
    if (S.sel) return closeDock();
    return;
  }
  if (typing) return;
  if (e.key === ' ' && room && (e.target === document.body || e.target === canvas)) (e.preventDefault(), playBtn.click());
  else if ((e.key === 'l' || e.key === 'L') && room) goLive();
  else if (e.key === '/' || e.key === 'g') (e.preventDefault(), select({ kind: 'goto' }, { fly: false }), document.getElementById('goto-q')?.focus());
  else if (e.key === 'n') nextAttention();
  else if (e.key === 'h') select({ kind: 'hall' });
});

// ---- Layout ----------------------------------------------------------------------------------

// The dock covers the right of a wide screen, the sheet the bottom of a phone:
// shift the view so the thing in focus sits in the middle of what is left.
function resize() {
  const w = canvas.clientWidth;
  const hh = canvas.clientHeight;
  renderer.setSize(w, hh, false);
  camera.aspect = w / hh;
  const docked = document.body.classList.contains('docked');
  if (docked && !phone()) camera.setViewOffset(w, hh, (dock.offsetWidth + 20) / 2, 0, w, hh);
  else if (docked && phone()) camera.setViewOffset(w, hh, 0, dock.offsetHeight / 2 - 20, w, hh);
  else camera.clearViewOffset();
  camera.updateProjectionMatrix();
  if (room) $('#scrub').style.bottom = phone() ? `${dock.offsetHeight + 8}px` : '';
}
addEventListener('resize', resize);
new ResizeObserver(resize).observe(dock);

// ---- Frame loop --------------------------------------------------------------------------------

const ease = (x) => x * x * x * (x * (x * 6 - 15) + 10);
let shown = -1;
let last = performance.now();
renderer.setAnimationLoop((ms) => {
  if (!city) return;
  const dt = (ms - last) / 1000;
  last = ms;
  if (room) {
    if (mode !== 'pause') roomT += Math.min(dt, 0.1) * (mode === 'play' ? speed : 1);
    if (mode === 'play' && roomT >= room.duration) setMode('live');
    room.update(roomT);
    if (dragging === null) seek.value = String(Math.min(1, roomT / room.duration));
    const upto = room.indexAt(roomT) + 1;
    if (upto !== shown) {
      [...feedEl.children].forEach((li, i) => {
        li.classList.toggle('past', i < upto);
        li.classList.toggle('now', i === upto - 1);
      });
      if (feedEl.isConnected) feedEl.querySelector('.now')?.scrollIntoView({ block: 'nearest' });
      shown = upto;
    }
    const ev = room.events[shown - 1];
    const words = mode === 'live' && roomT >= room.duration ? 'now' : ev?.clock ? `${ev.clock} · ${ev.who}` : '';
    if ($('#when').textContent !== words) {
      $('#when').textContent = words;
      if (mode !== 'live') render();
    }
  } else city.update(now());
  if (flight) {
    const u = ease(Math.min(1, (ms - flight.t0) / 800));
    controls.target.lerpVectors(flight.from[0], flight.to[0], u);
    camera.position.lerpVectors(flight.from[1], flight.to[1], u);
    if (u >= 1) flight = null;
  }
  if (crateFlight) {
    const u = Math.min(1, (ms - crateFlight.t0) / 1400);
    const k = crateFlight.back ? 1 - ease(u) : ease(u);
    crate.position.lerpVectors(crateFlight.from, crateFlight.to, k);
    crate.position.y += Math.sin(Math.PI * k) * 1.2;
    crate.rotation.y = k * 3;
    if (u >= 1 && (crateFlight.back || !S.busy.size)) {
      crate.visible = false;
      crateFlight = null;
    }
  }
  if (ring.visible) ring.material.opacity = 1;
  controls.update();
  renderer.render(room ? roomScene : scene, camera);
  if (!room) placePins();
});

// ---- Demo drawer (states only; not part of the design) --------------------------------------

$('#d-stale').addEventListener('change', (e) => (fake.demo.stale = e.target.checked));
$('#d-ro').addEventListener('change', async (e) => {
  fake.demo.readOnly = e.target.checked;
  await loadAll();
});
$('#d-tl').addEventListener('change', (e) => ((S.timelineStuck = e.target.checked), render()));
$('#d-herdr').addEventListener('change', (e) => ((S.herdrDown = e.target.checked), render()));
$('#d-next').addEventListener('change', (e) => (fake.demo.next = e.target.value));
const JUMPS = {
  overview: () => closeAll(),
  mate: () => select({ kind: 'mate', project: 'plantrip' }),
  chat: () => select({ kind: 'mate', project: 'plantrip', tab: 'chat' }),
  crew: () => select({ kind: 'crew', project: 'plantrip', crew: 'PLA-58' }),
  merge: () => {
    select({ kind: 'crew', project: 'plantrip', crew: 'PLA-58' });
    setTimeout(() => document.getElementById('act-plantrip/PLA-58-merge')?.click(), 900);
  },
  peek: () => select({ kind: 'crew', project: 'mihuv2', crew: 'HU-129', tab: 'peek' }),
  diff: () => select({ kind: 'crew', project: 'plantrip', crew: 'PLA-58', tab: 'diff' }),
  topics: () => select({ kind: 'hall', tab: 'topics' }),
  topic: () => select({ kind: 'hall', tab: 'topics', topic: 't-plantrip-12' }),
  fork: () => {
    select({ kind: 'hall', tab: 'topics', topic: 't-plantrip-12' });
    openFork(fake.topics[0]);
  },
  tent: () => enterRoom('plantrip/plantrip-api'),
  replay: async () => {
    await enterRoom('plantrip/plantrip-api');
    roomT = room.duration * 0.45;
    setMode('pause');
    setTimeout(() => select({ kind: 'crew', project: 'plantrip', crew: 'PLA-58', back: { kind: 'tent', project: 'plantrip', repo: 'plantrip-api' } }, { fly: false }), 300);
  },
  add: () => select({ kind: 'add', back: { kind: 'hall', tab: 'projects' } }),
  create: () => select({ kind: 'create', project: 'landing', back: { kind: 'mate', project: 'landing' } }),
  delete: () => {
    select({ kind: 'mate', project: 'landing', tab: 'project' });
    setTimeout(() => deleteProject('landing'), 700);
  },
  inbox: () => select({ kind: 'mate', project: 'mihuv2', tab: 'inbox' }),
  system: () => select({ kind: 'system' }),
  goto: () => select({ kind: 'goto' }, { fly: false }),
};
function closeAll() {
  if (room) leaveRoom();
  if (S.sel) {
    S.sel = null;
    closeDock();
  }
  if (home) fly(home.target.clone(), home.position.clone());
}
$('#d-jump').replaceChildren(...Object.keys(JUMPS).map((k) => h('button', { type: 'button', onclick: () => {
  history.replaceState(null, '', `#${k}`);
  JUMPS[k]();
} }, k)));

// ---- Boot -------------------------------------------------------------------------------------

Object.assign(controls, CITY_LIMITS);
resize();
loadAll().then(() => {
  performance.mark('city:drawn');
  const k = location.hash.slice(1);
  if (k === 'stale') {
    $('#d-stale').checked = true;
    fake.demo.stale = true;
  } else if (k === 'readonly') {
    $('#d-ro').checked = true;
    fake.demo.readOnly = true;
    loadAll().then(() => JUMPS.crew());
  } else JUMPS[k]?.();
}, (e) => setLive('stale', e.message)).finally(poll);
