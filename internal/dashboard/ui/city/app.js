import * as THREE from './vendor/three.js';
import { OrbitControls } from './vendor/three.js';
import { createCity, lights, PALETTE, harnessName, formatTokens } from './city.js';
import { createTentRoom } from './room.js';
import { readEpisode } from './episode.js';

// Agent City's page: the workspace drawn as a city (city.js), kept live from
// the dashboard's read-only API the way Mate Office is - one long poll on
// /api/events, and a reload of what it reads when an event lands. Hover a
// mascot, a tent or City Hall for what it is; click a mascot or City Hall to
// open it in Mate Office, a tent (or open #tent=<project>/<repo>) to step
// into it and replay its latest merged task (room.js, episode.js).

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
let first = true;
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
  first = true;
}

async function loadAll() {
  const [wsBody, board, news] = await Promise.all([
    api('/api/workspace'),
    api(`/api/board?days=${BOARD_DAYS}`).catch(() => null),
    api('/api/news').catch(() => null),
  ]);
  const projects = wsBody.projects ?? [];
  const bodies = await Promise.all(projects.map((p) => api(`/api/projects/${enc(p.name)}`).catch(() => null)));
  ws = wsBody;
  const next = JSON.stringify(projects.map((p) => [p.name, p.repos ?? []]));
  if (next !== shape) {
    shape = next;
    buildCity(projects);
  }
  const t = now();
  const fresh = !first;
  projects.forEach((p, i) => {
    const body = bodies[i];
    if (body) tasks.set(p.name, body.tasks ?? []);
    const rows = tasks.get(p.name) ?? [];
    city.setMate(p.name, p.mate ?? {}, t);
    city.setCrews(p.name, rows.filter((r) => r.state !== 'gone').map((r) => ({ ...r })), t, { fresh });
    for (const repo of tentRepos(p)) city.setMerged(p.name, repo, rows.filter((r) => r.role === 'ship' && r.merged_at && onTent(p, repo)(r)).length, t);
    const bp = board?.projects?.find((x) => x.name === p.name);
    boardStatsOf.set(p.name, bp ? boardStats(bp, board.days) : null);
    city.setBoard(p.name, boardStatsOf.get(p.name), t);
  });
  city.setHall(projects.map((p) => {
    const bp = board?.projects?.find((x) => x.name === p.name);
    const working = (tasks.get(p.name) ?? []).filter((r) => r.state !== 'gone').length;
    return { name: p.name, running: Boolean(p.mate?.running), working, tokens: bp?.days?.at(-1)?.tokens?.total ?? p.mate?.tokens_today ?? 0 };
  }));
  city.setMailbox(waitingNews(news), t);
  lastEventID = wsBody.last_event_id ?? 0;
  first = false;
  setLive('live', 'live');
  $('#boot').hidden = true;
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
// live. An event reloads everything; a quiet poll still reads the news,
// which is a file the REPL writes, not an event.
async function poll() {
  for (;;) {
    try {
      const body = await api(`/api/events?since=${lastEventID}&wait=20`);
      const id = body.last_event_id ?? 0;
      const moved = body.events?.length || body.now?.length;
      // A smaller id is a `mate reindex` under us: the old cursor will never come round again.
      if (id < lastEventID || moved) await loadAll();
      else {
        lastEventID = id;
        city?.setMailbox(waitingNews(await api('/api/news').catch(() => null)), now());
        setLive('live', 'live');
      }
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
const KIND = { spawn: 'spawn', status: 'status', ask: 'ask', answer: 'answer', handback: 'hand back', review: 'review', verdict: 'verdict', merge: 'merge', sleep: 'idle' };

async function enterRoom(key) {
  const tent = city?.tents.get(key);
  if (!tent) return;
  const token = ++roomToken;
  const p = ws.projects.find((x) => x.name === tent.project);
  if (!room) cityView = { target: controls.target.clone(), position: camera.position.clone() };
  document.body.classList.add('in-room');
  tip.style.display = 'none';
  $('#room-title').textContent = tent.repo;
  $('#room-dept').textContent = `project ${tent.project} · loading…`;
  feed.replaceChildren();
  resize();
  if (location.hash !== `#tent=${key}`) history.replaceState(null, '', `#tent=${key}`);
  let ep = { episode: null, recent: [] };
  try {
    ep = await readEpisode(api, tent.project, tasks.get(tent.project) ?? [], onTent(p, tent.repo));
  } catch (e) {
    $('#room-dept').textContent = `project ${tent.project} · ${e.message}`;
  }
  if (token !== roomToken) return;
  if (room) {
    roomScene.remove(room.group);
    room.dispose();
  }
  const bs = city.blocks.get(tent.project);
  const stats = boardStatsOf.get(tent.project);
  room = createTentRoom({ project: tent.project, repo: tent.repo, floor: bs.floor, mateHarness: p.mate?.harness, episode: ep.episode, recent: ep.recent, finished: stats?.finished ?? 0, merged: tent.merged ?? 0, days: BOARD_DAYS });
  roomScene.add(room.group);
  roomT = 0;
  shown = -1;
  flight = null;
  controls.target.copy(room.view.target);
  camera.position.copy(room.view.position);
  Object.assign(controls, ROOM_LIMITS);
  $('#room-dept').textContent = `project ${tent.project} · ${room.summary}`;
  feed.replaceChildren(...room.events.map((ev) => {
    const li = document.createElement('li');
    const kind = document.createElement('i');
    kind.textContent = KIND[ev.kind];
    kind.dataset.kind = ev.kind;
    const by = document.createElement('b');
    by.textContent = ev.clock ? `${ev.who} · ${ev.clock}` : ev.who;
    const text = document.createElement('span');
    text.textContent = ev.text;
    li.append(kind, by, text);
    return li;
  }));
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
  const list = [...(room ? room.mascots : city?.mascots ?? new Map()).values()];
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
canvas.addEventListener('pointermove', (e) => {
  const target = pick(e);
  const words = target && describe(target);
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

const bar = $('#bar');
const ease = (x) => x * x * x * (x * (x * 6 - 15) + 10);
let shown = -1;
let last = performance.now();
renderer.setAnimationLoop((ms) => {
  if (!city) return;
  const dt = (ms - last) / 1000;
  last = ms;
  if (room) {
    roomT = (roomT + dt) % room.duration;
    room.update(roomT);
    bar.style.transform = `scaleX(${roomT / room.duration})`;
    // Light up the feed up to now; the newest entry is marked.
    const upto = room.events.filter((ev) => ev.at <= roomT).length;
    if (upto !== shown) {
      [...feed.children].forEach((li, i) => {
        li.classList.toggle('past', i < upto);
        li.classList.toggle('now', i === upto - 1);
      });
      feed.querySelector('.now')?.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
      shown = upto;
    }
  } else city.update(now());
  if (flight) {
    const u = ease(Math.min(1, (ms - flight.t0) / 800));
    controls.target.lerpVectors(flight.from[0], flight.to[0], u);
    camera.position.lerpVectors(flight.from[1], flight.to[1], u);
    if (u >= 1) flight = null;
  }
  controls.update();
  renderer.render(room ? roomScene : scene, camera);
});

Object.assign(controls, CITY_LIMITS);
resize();
loadAll().then(fromHash, (e) => {
  setLive('stale', e.message);
  $('#boot').textContent = `Agent City could not load: ${e.message}`;
}).finally(poll);
