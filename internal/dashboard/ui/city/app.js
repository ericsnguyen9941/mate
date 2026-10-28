import * as THREE from './vendor/three.js';
import { OrbitControls } from './vendor/three.js';
import { createCity, lights, PALETTE, harnessName, formatTokens } from './city.js';
import { createTentRoom } from './room.js';
import { readHistory } from './history.js';

// Agent City's page: the workspace drawn as a city (city.js), kept live from
// the dashboard's read-only API the way Mate Office is - one long poll on
// /api/events, and a reload of what it reads when an event lands. Hover a
// mascot, a tent or City Hall for what it is; click a mascot or City Hall to
// open it in Mate Office, a tent (or open #tent=<project>/<repo>) to step
// into it and play back what its agents did, with a bar to scrub through it
// (room.js, history.js).

const BOARD_DAYS = 7;
const $ = (s) => document.querySelector(s);
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
// Mascots face +z; keep the view in front of them.
const CITY_LIMITS = { minDistance: 8, minAzimuthAngle: -1.1, maxAzimuthAngle: 1.1 };
const ROOM_LIMITS = { minDistance: 6, maxDistance: 34, minAzimuthAngle: -0.9, maxAzimuthAngle: 0.9 };

const enc = encodeURIComponent;
function api(path) {
  // Read-only by construction: the only request this page makes is a GET
  // to its own origin.
  return fetch(path, { method: 'GET', headers: { accept: 'application/json' } })
    .catch(() => {
      throw new Error('the dashboard is not answering');
    })
    .then((r) =>
      r.json().then(
        (body) => {
          if (!r.ok) throw new Error(body.reason || body.error || `the dashboard answered ${r.status}`);
          return body;
        },
        () => {
          throw new Error(`the dashboard answered ${r.status} with something that is not JSON`);
        },
      ),
    );
}

// ---- State -------------------------------------------------------------------

let scene = null;
let city = null;
let shape = ''; // the projects and repos the city was built for
let home = null;
let ws = null; // GET /api/workspace
const tasks = new Map(); // project → task rows
const boardStatsOf = new Map(); // project → the signboard's numbers, or null
let lastEventID = 0;
const start = performance.now();
const now = () => (performance.now() - start) / 1000; // the city's clock, in seconds

function setLive(state, text) {
  const pill = $('#live');
  pill.dataset.state = state;
  pill.querySelector('.text').textContent = text;
}

// The project's repos as the city draws them: none means one tent named after it.
const tentRepos = (p) => (p.repos?.length ? p.repos : [p.name]);
// Whether a task row worked in the tent: its repo's, or the project's first
// tent when its repo is unknown or no longer in project.yaml.
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
  home = { target: new THREE.Vector3(0, 1, 1), position: new THREE.Vector3(-20, 62, 70).multiplyScalar(k) };
  controls.maxDistance = 170 * k;
  if (!room) {
    controls.target.copy(home.target);
    camera.position.copy(home.position);
  }
  $('#hidden').textContent = city.hidden ? `${city.hidden} more project${city.hidden === 1 ? '' : 's'} than the city has blocks` : '';
}

// What has been read. The city is drawn as soon as /api/workspace answers;
// each project fills in as its own page arrives, and the board and the news,
// which are slower and change less, are read beside them and never hold the
// city up.
let board = null; // GET /api/board
let boardAt = -Infinity; // when it was last asked for (performance.now())
let boardBusy = false;
const BOARD_EVERY = 5 * 60e3; // the board counts days; minutes-old numbers are fine
const loaded = new Set(); // projects whose crews are placed; a new crew there says its name

function applyProject(p, t) {
  const rows = tasks.get(p.name);
  if (!rows) return;
  city.setCrews(p.name, rows.filter((r) => r.state !== 'gone'), t, { fresh: loaded.has(p.name) });
  for (const repo of tentRepos(p)) city.setMerged(p.name, repo, rows.filter((r) => r.role === 'ship' && r.merged_at && onTent(p, repo)(r)).length, t);
  loaded.add(p.name);
  if (room?.project === p.name) refreshRoom();
}

function applyHall() {
  city.setHall((ws.projects ?? []).map((p) => {
    const bp = board?.projects?.find((x) => x.name === p.name);
    const working = (tasks.get(p.name) ?? []).filter((r) => r.state !== 'gone').length;
    return { name: p.name, running: Boolean(p.mate?.running), working, tokens: bp?.days?.at(-1)?.tokens?.total ?? p.mate?.tokens_today ?? 0 };
  }));
}

function applyBoard(t) {
  for (const p of ws.projects ?? []) {
    const bp = board?.projects?.find((x) => x.name === p.name);
    boardStatsOf.set(p.name, bp ? boardStats(bp, board.days) : null);
    city.setBoard(p.name, boardStatsOf.get(p.name), t);
  }
  applyHall();
}

/**
 * Read the workspace, redraw the city if its projects or repos changed, then
 * read the pages of `only` (a set of project names; every project when
 * null), applying each as it arrives.
 */
async function loadAll(only = null) {
  const wsBody = await api('/api/workspace');
  ws = wsBody;
  // Polling from the id the workspace was built at can only read an event twice, never miss one.
  lastEventID = wsBody.last_event_id ?? 0;
  const projects = wsBody.projects ?? [];
  const next = JSON.stringify(projects.map((p) => [p.name, p.repos ?? []]));
  const t = now();
  if (next !== shape) {
    shape = next;
    buildCity(projects);
    loaded.clear();
    for (const p of projects) applyProject(p, t);
    if (board) applyBoard(t);
    only = null;
  }
  for (const p of projects) city.setMate(p.name, p.mate ?? {}, t);
  applyHall();
  setLive('live', 'live');
  $('#boot').hidden = true;
  await Promise.all(projects.filter((p) => !only || only.has(p.name)).map((p) =>
    api(`/api/projects/${enc(p.name)}`).then((body) => {
      if (ws !== wsBody) return;
      tasks.set(p.name, body.tasks ?? []);
      applyProject(p, now());
    }, () => {}),
  ));
  if (ws === wsBody) applyHall();
}

function refreshBoard() {
  if (boardBusy || performance.now() - boardAt < BOARD_EVERY) return;
  boardBusy = true;
  boardAt = performance.now();
  api(`/api/board?days=${BOARD_DAYS}`)
    .then((body) => {
      board = body;
      if (city) applyBoard(now());
    }, () => (boardAt = -Infinity))
    .finally(() => (boardBusy = false));
}

let newsBusy = false;
function refreshNews() {
  if (newsBusy) return;
  newsBusy = true;
  api('/api/news')
    .then((body) => city?.setMailbox(waitingNews(body), now()), () => {})
    .finally(() => (newsBusy = false));
}

function boardStats(bp, days) {
  const s = { days, tokens: 0, finished: 0, passed: 0 };
  for (const d of bp.days ?? []) {
    s.tokens += d.tokens?.total ?? 0;
    s.finished += d.closed?.finished ?? 0;
    s.passed += d.reviews?.pass ?? 0;
  }
  return s;
}

// News the captain owes an answer: unread or read and waiting.
function waitingNews(news) {
  let n = 0;
  for (const p of news?.projects ?? []) for (const x of p.news ?? []) if (x.status === 'new' || x.status === 'waiting') n++;
  return n;
}

// One long poll on /api/events, as the admin UI and Office keep themselves
// live. An event reloads the workspace and the pages of the projects it
// names; every answer, quiet or not, reads the news again (a file the REPL
// writes, not an event) and the board once it is BOARD_EVERY old.
async function poll() {
  for (;;) {
    try {
      const body = await api(`/api/events?since=${lastEventID}&wait=20`);
      const id = body.last_event_id ?? 0;
      const moved = [...(body.events ?? []), ...(body.now ?? [])];
      // A smaller id is a `mate reindex` under us: the old cursor will never come round again.
      if (id < lastEventID) await loadAll();
      else if (moved.length) {
        const names = new Set(moved.map((x) => x.project));
        await loadAll(names.has('') || names.has(undefined) ? null : names);
      } else lastEventID = id;
      setLive('live', 'live');
      refreshNews();
      refreshBoard();
    } catch (e) {
      setLive('stale', e.message);
      await new Promise((r) => setTimeout(r, 3000));
    }
  }
}

// ---- Inside a tent ------------------------------------------------------------

let room = null;
let roomT = 0;
let cityView = null;
let roomToken = 0;
const roomScene = new THREE.Scene();
roomScene.background = new THREE.Color(PALETTE.paper);
lights(roomScene, renderer);
const panel = $('#room');
const feed = $('#feed');
const KIND = { spawn: 'spawn', status: 'status', ask: 'ask', answer: 'answer', handback: 'hand back', review: 'review', verdict: 'verdict', merge: 'merge', close: 'close', sleep: 'idle' };

// Histories read so far, by tent: kept while the project's task rows are the
// same ones, so a second visit, or a click after a hover, opens at once.
const histories = new Map(); // tentKey → { rows, promise }
function historyFor(key) {
  const tent = city?.tents.get(key);
  const p = ws?.projects.find((x) => x.name === tent?.project);
  if (!tent || !p) return Promise.resolve({ since: '', tasks: [], events: [] });
  const rows = tasks.get(tent.project) ?? [];
  const have = histories.get(key);
  if (have && have.rows === rows) return have.promise;
  const promise = readHistory(api, tent.project, rows, onTent(p, tent.repo));
  histories.set(key, { rows, promise });
  promise.catch(() => histories.delete(key));
  return promise;
}

// Playback: 'live' runs at 1x and stays at now once it gets there, 'play'
// runs at `speed` until now, then goes live; 'pause' holds.
let mode = 'live';
let speed = 1;
const SPEEDS = [1, 4, 16, 64];

async function enterRoom(key) {
  const tent = city?.tents.get(key);
  if (!tent) return;
  const token = ++roomToken;
  if (!room) cityView = { target: controls.target.clone(), position: camera.position.clone() };
  document.body.classList.add('in-room');
  tip.style.display = 'none';
  $('#room-title').textContent = tent.repo;
  $('#room-dept').textContent = `project ${tent.project} · loading…`;
  feed.replaceChildren();
  resize();
  if (location.hash !== `#tent=${key}`) history.replaceState(null, '', `#tent=${key}`);
  const next = await buildRoom(key, token);
  if (!next) return;
  // Open on the last few events, playing up to now.
  setMode('live');
  roomT = next.events.at(-3)?.at ?? 0;
  flight = null;
  controls.target.copy(next.view.target);
  camera.position.copy(next.view.position);
  Object.assign(controls, ROOM_LIMITS);
}

// Read the tent's history and put up its room in place of the one showing;
// null when another tent was asked for meanwhile.
async function buildRoom(key, token) {
  const tent = city.tents.get(key);
  const p = ws.projects.find((x) => x.name === tent.project);
  let h = { since: '', tasks: [], events: [] };
  try {
    // A click before the project's page arrived reads it first.
    if (!tasks.has(tent.project)) tasks.set(tent.project, (await api(`/api/projects/${enc(tent.project)}`)).tasks ?? []);
    h = await historyFor(key);
  } catch (e) {
    $('#room-dept').textContent = `project ${tent.project} · ${e.message}`;
  }
  if (token !== roomToken) return null;
  const was = room;
  if (was) {
    roomScene.remove(was.group);
    was.dispose();
  }
  room = createTentRoom({ project: tent.project, repo: tent.repo, floor: city.blocks.get(tent.project).floor, mateHarness: p.mate?.harness, mateState: p.mate?.running ? p.mate.state : 'asleep', history: h });
  room.key = key;
  roomScene.add(room.group);
  if (!was) resize();
  shown = -1;
  $('#room-dept').textContent = `project ${tent.project} · ${room.summary}`;
  $('#room-h3').textContent = room.since ? `WHAT ITS AGENTS DID · SINCE ${room.since}` : 'WHAT ITS AGENTS DID';
  feed.replaceChildren(...room.events.map((ev, i) => {
    const li = document.createElement('li');
    const kind = document.createElement('i');
    kind.textContent = KIND[ev.kind];
    kind.dataset.kind = ev.kind;
    const by = document.createElement('b');
    by.textContent = ev.clock ? `${ev.who} · ${ev.clock}` : ev.who;
    const text = document.createElement('span');
    text.textContent = ev.text;
    li.append(kind, by, text);
    li.dataset.i = i;
    return li;
  }));
  drawTicks();
  return room;
}

// New events on the tent's project: read the history again and keep the
// view where it was - at the same event when scrubbed back, or, live and
// caught up, on the first new event so it plays out.
let refreshTimer = 0;
function refreshRoom() {
  clearTimeout(refreshTimer);
  refreshTimer = setTimeout(async () => {
    if (!room) return;
    const old = room;
    const caughtUp = mode === 'live' && roomT >= old.duration;
    const i = old.indexAt(roomT);
    const at = old.events[i];
    const next = await buildRoom(old.key, roomToken);
    if (!next) return;
    if (caughtUp) roomT = next.events[old.events.length]?.at ?? next.duration;
    else if (at?.iso) {
      const j = next.events.findIndex((e) => e.iso === at.iso && e.crew === at.crew && e.kind === at.kind);
      if (j >= 0) roomT = next.events[j].at + (roomT - at.at);
    }
  }, 800);
}

function leaveRoom() {
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
  resize();
  if (location.hash) history.replaceState(null, '', location.pathname + location.search);
}

// ---- The scrub bar ----------------------------------------------------------------

const seek = $('#seek');
const playBtn = $('#play');
const speedBtn = $('#speed');
const when = $('#when');
let dragging = null; // the mode before a drag
function setMode(next) {
  mode = next;
  playBtn.textContent = mode === 'pause' ? '\u25B6' : '\u275A\u275A';
  playBtn.setAttribute('aria-label', mode === 'pause' ? 'Play' : 'Pause');
  $('#golive').classList.toggle('on', mode === 'live');
}
playBtn.addEventListener('click', () => {
  if (!room) return;
  if (mode !== 'pause') return setMode('pause');
  // Play from the start when already at now.
  if (roomT >= room.duration) roomT = 0;
  setMode('play');
});
speedBtn.addEventListener('click', () => {
  speed = SPEEDS[(SPEEDS.indexOf(speed) + 1) % SPEEDS.length];
  speedBtn.textContent = `${speed}\u00D7`;
  if (mode === 'live' && room && roomT < room.duration) setMode('play');
});
$('#golive').addEventListener('click', () => {
  if (room) roomT = room.duration;
  setMode('live');
});
seek.addEventListener('input', () => {
  if (!room) return;
  dragging ??= mode;
  mode = 'pause';
  roomT = Number(seek.value) * room.duration;
});
seek.addEventListener('change', () => {
  setMode(dragging === 'pause' ? 'pause' : 'play');
  dragging = null;
});
feed.addEventListener('click', (e) => {
  const li = e.target.closest('li[data-i]');
  if (!room || !li) return;
  roomT = room.events[Number(li.dataset.i)].at;
  if (mode === 'live') setMode('play');
});
addEventListener('keydown', (e) => {
  if (e.key === ' ' && room && e.target === document.body) {
    e.preventDefault();
    playBtn.click();
  }
});
// Day marks and merges along the bar.
function drawTicks() {
  const marks = [];
  let day = '';
  let lastDay = -Infinity;
  // A day's label needs about 40px; one closer to the last is left out.
  const gap = 4000 / Math.max(1, seek.clientWidth);
  for (const ev of room.events) {
    if (!ev.iso) continue;
    const x = (100 * ev.at) / room.duration;
    const pos = `${x}%`;
    const d = ev.clock.slice(0, 5);
    if (d !== day && x - lastDay < gap) day = d;
    if (d !== day) {
      lastDay = x;
      day = d;
      const m = document.createElement('span');
      m.className = 'day';
      m.style.left = pos;
      m.textContent = d;
      marks.push(m);
    }
    if (ev.kind === 'merge') {
      const m = document.createElement('i');
      m.className = 'merge';
      m.style.left = pos;
      m.title = ev.text;
      marks.push(m);
    }
  }
  $('#ticks').replaceChildren(...marks);
}
setMode('live');
$('#back').addEventListener('click', leaveRoom);
addEventListener('keydown', (e) => e.key === 'Escape' && leaveRoom());
function fromHash() {
  let key = '';
  try {
    key = decodeURIComponent(location.hash.match(/^#tent=(.+)$/)?.[1] ?? '');
  } catch {}
  if (city?.tents.has(key)) enterRoom(key);
  else if (document.body.classList.contains('in-room')) leaveRoom();
}
addEventListener('hashchange', fromHash);

// ---- Pointer ------------------------------------------------------------------

let flight = null;
function fly(target, position) {
  flight = { from: [controls.target.clone(), camera.position.clone()], to: [target, position], t0: performance.now() };
}
$('#home').addEventListener('click', () => (room ? fly(room.view.target.clone(), room.view.position.clone()) : fly(home.target.clone(), home.position.clone())));

const tip = $('#tip');
const ray = new THREE.Raycaster();
const pointer = new THREE.Vector2();
function aim(e) {
  const r = canvas.getBoundingClientRect();
  pointer.set(((e.clientX - r.left) / r.width) * 2 - 1, -((e.clientY - r.top) / r.height) * 2 + 1);
  ray.setFromCamera(pointer, camera);
}
// What is under the pointer: { mascot: id } | { tent: key } | { hall: true } | null.
function pick(e) {
  aim(e);
  const list = [...(room ? room.mascots : city?.mascots ?? new Map()).values()].filter((m) => m.root.visible);
  const hit = ray.intersectObjects(list.map((m) => m.body))[0];
  if (hit) return { mascot: list.find((m) => m.body === hit.object).id };
  if (room || !city) return null;
  let o = ray.intersectObjects([...[...city.tents.values()].map((b) => b.group), ...city.hall], true)[0]?.object;
  while (o && !o.userData.tent && o.name !== 'mailbox' && !city.hall.includes(o)) o = o.parent;
  if (!o) return null;
  return o.userData.tent ? { tent: o.userData.tent } : { hall: true };
}
const seatOf = (id) => (room ? room.seats : city.seats).get(id);
function officeHref(target) {
  if (target.hall) return '../office/#captain';
  const s = seatOf(target.mascot);
  if (!s) return null;
  return s.kind === 'mate' ? `../office/#mate/${enc(s.project)}` : `../office/#crew/${enc(s.project)}/${enc(s.crew)}`;
}
function showTip(e, heading, detail) {
  const title = document.createElement('b');
  title.textContent = heading;
  const line = document.createElement('span');
  line.textContent = detail;
  tip.replaceChildren(title, line);
  tip.style.display = 'block';
  tip.style.left = `${e.clientX}px`;
  tip.style.top = `${e.clientY}px`;
}
function describe(target) {
  if (target.mascot) {
    const s = seatOf(target.mascot);
    if (!s) return null;
    if (s.kind === 'mate') return [`Mate · ${s.project}`, `${harnessName(s.harness)}${s.state ? ` · ${s.state.replace(/_/g, ' ')}` : ''} · click for its desk in Office`];
    const kit = [harnessName(s.harness), s.model, s.effort].filter(Boolean).join(' · ');
    return [`${s.crew} · ${s.role}`, `${kit}${s.state ? ` · ${s.state.replace(/_/g, ' ')}` : ''} · click to open in Office`];
  }
  if (target.tent) {
    const b = city.tents.get(target.tent);
    return [`${b.repo} · ${b.project}`, `${b.merged ?? 0} merged · ${b.crews.size} at work · click to step in`];
  }
  const p = ws?.projects ?? [];
  const waiting = p.reduce((n, x) => n + (x.inbox_waiting ?? 0), 0);
  const tokens = p.reduce((n, x) => n + (x.mate?.tokens_today ?? 0), 0);
  return ['City Hall · captain', `${waiting} waiting in the inbox · Mates ${formatTokens(tokens)} tokens today · click for the captain's desk`];
}
let hoverTimer = 0;
canvas.addEventListener('pointermove', (e) => {
  const target = pick(e);
  const words = target && describe(target);
  // Read a tent's history while the pointer rests on it, so the click opens at once.
  clearTimeout(hoverTimer);
  if (target?.tent && tasks.has(city.tents.get(target.tent)?.project)) hoverTimer = setTimeout(() => historyFor(target.tent).catch(() => {}), 150);
  canvas.style.cursor = words ? 'pointer' : '';
  if (words) showTip(e, ...words);
  else tip.style.display = 'none';
});
canvas.addEventListener('pointerleave', () => (tip.style.display = 'none'));
// A click is a press that did not drag the view.
let down = null;
canvas.addEventListener('pointerdown', (e) => (down = [e.clientX, e.clientY]));
canvas.addEventListener('click', (e) => {
  if (down && Math.hypot(e.clientX - down[0], e.clientY - down[1]) > 5) return;
  const target = pick(e);
  if (!target) return;
  if (target.tent) {
    location.hash = `#tent=${target.tent}`;
    return;
  }
  const href = officeHref(target);
  if (href) location.href = href;
});

// In a tent the feed panel covers the right of a wide screen, so shift the
// view to centre the room in what is left.
function resize() {
  const w = canvas.clientWidth;
  const h = canvas.clientHeight;
  renderer.setSize(w, h, false);
  camera.aspect = w / h;
  const side = room && w > 900 ? panel.offsetWidth + 20 : 0;
  if (side) camera.setViewOffset(w, h, side / 2, 0, w, h);
  else camera.clearViewOffset();
  camera.updateProjectionMatrix();
}
addEventListener('resize', resize);

// ---- Frame loop -----------------------------------------------------------------

const ease = (x) => x * x * x * (x * (x * 6 - 15) + 10);
let shown = -1;
let last = performance.now();
const drawn = new WeakSet();
renderer.setAnimationLoop((ms) => {
  if (!city) return;
  const dt = (ms - last) / 1000;
  last = ms;
  if (room) {
    if (mode !== 'pause') roomT += Math.min(dt, 0.1) * (mode === 'play' ? speed : 1);
    if (mode === 'play' && roomT >= room.duration) setMode('live');
    room.update(roomT);
    if (dragging === null) seek.value = String(Math.min(1, roomT / room.duration));
    // Light up the feed up to the playhead; its entry is marked.
    const upto = room.indexAt(roomT) + 1;
    if (upto !== shown) {
      const items = feed.children;
      const lo = Math.max(0, Math.min(shown, upto) - 1);
      const hi = shown < 0 ? items.length : Math.min(items.length, Math.max(shown, upto) + 1);
      for (let i = lo; i < hi; i++) {
        items[i].classList.toggle('past', i < upto);
        items[i].classList.toggle('now', i === upto - 1);
      }
      feed.querySelector('.now')?.scrollIntoView({ block: 'nearest', behavior: mode === 'play' && speed > 4 ? 'auto' : 'smooth' });
      shown = upto;
    }
    const ev = room.events[shown - 1];
    const words = mode === 'live' && roomT >= room.duration ? 'now' : ev?.clock ? `${ev.clock} · ${ev.who}` : '';
    if (when.textContent !== words) when.textContent = words;
  } else city.update(now());
  if (flight) {
    const u = ease(Math.min(1, (ms - flight.t0) / 800));
    controls.target.lerpVectors(flight.from[0], flight.to[0], u);
    camera.position.lerpVectors(flight.from[1], flight.to[1], u);
    if (u >= 1) flight = null;
  }
  controls.update();
  renderer.render(room ? roomScene : scene, camera);
  // Marks for the browser's performance timeline: the first frame of the city, and of each tent.
  if (!drawn.has(room ?? city)) {
    drawn.add(room ?? city);
    performance.mark(room ? `city:tent ${room.project}/${room.repo}` : 'city:drawn');
  }
});

Object.assign(controls, CITY_LIMITS);
resize();
loadAll().then(fromHash, (e) => {
  setLive('stale', e.message);
  $('#boot').textContent = `Agent City could not load: ${e.message}`;
}).finally(poll);
refreshNews();
refreshBoard();
