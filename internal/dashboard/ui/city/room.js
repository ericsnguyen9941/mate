import * as THREE from './vendor/three.js';
import { RoundedBoxGeometry } from './vendor/three.js';
import { AgentMascot, handoff } from './agent-mascot.js';
import { HARNESSES } from './harnesses.js';
import { PALETTE, LiveTexture, matte, textTexture, valanceGeometry, harnessLook, harnessName, poseOf } from './city.js';

// One repo's tent seen up close: a big open tent (front and roof cut away)
// with a task board, a round table where the crews sit with their laptops,
// and the project's Mate between the two. It plays out what history.js read
// from the dashboard, in order: the Mate spawns crews, they report status,
// ask the Mate and get its answers, hand back; reviewers come, read the diff
// and give their verdicts; the Mate merges or closes. Bubbles show the lines
// shortened; `events` carries them in full with their real times for the
// page's feed. update(t) is a pure function of t, so the page can scrub to
// any time; past `duration` the room shows things as they are now.

export const ROOM_LOOP = 18; // a room with nothing to play
const STEP = 2.5; // room seconds per event
const LEAD = 1; // room seconds before the first one
const W = 15; // floor width
const D = 10.5; // floor depth
const EAVE = 3.6;
const TABLE = { x: 1.9, z: -0.7, r: 1.75, h: 0.72 };
const BOARD = { x: -3.9, z: -2.9, w: 6.2, ry: 0.3 };
const BOARD_PX = [1500, 900];
const CREW_SIZE = 1.0;
const MATE_SIZE = 1.2;
const FONT = 'ui-sans-serif, system-ui, sans-serif';

// 'pla-104-map-pins' → { key: 'PLA-104', title: 'map pins' }
export function ticket(id) {
  const m = id.match(/^([a-z]+-\d+)-?(.*)$/i);
  return m ? { key: m[1].toUpperCase(), title: m[2].replace(/-/g, ' ') } : { key: id, title: '' };
}

function fit(ctx, text, w) {
  if (ctx.measureText(text).width <= w) return text;
  let s = text;
  while (s && ctx.measureText(`${s}…`).width > w) s = s.slice(0, -1);
  return `${s.trimEnd()}…`;
}

function drawTaskBoard(ctx, { repo, color, columns, done, merged, since }, w, h) {
  const font = (weight, px) => `${weight} ${px}px ${FONT}`;
  ctx.fillStyle = PALETTE.pad;
  ctx.beginPath();
  ctx.roundRect(0, 0, w, h, 34);
  ctx.fill();
  ctx.fillStyle = color;
  ctx.beginPath();
  ctx.roundRect(0, 0, w, 120, [34, 34, 0, 0]);
  ctx.fill();
  ctx.fillStyle = PALETTE.ink;
  ctx.textBaseline = 'middle';
  ctx.textAlign = 'left';
  ctx.font = font(800, 62);
  ctx.fillText(fit(ctx, repo, w - 520), 44, 62);
  ctx.textAlign = 'right';
  ctx.fillStyle = 'rgba(36,33,29,0.6)';
  ctx.font = font(800, 34);
  ctx.fillText('TASK BOARD', w - 44, 64);

  const titles = [['doing', 'DOING'], ['review', 'REVIEW'], ['done', 'DONE']];
  const colW = (w - 44 * 2 - 30 * 2) / 3;
  titles.forEach(([id, title], i) => {
    const x = 44 + i * (colW + 30);
    ctx.fillStyle = 'rgba(36,33,29,0.05)';
    ctx.beginPath();
    ctx.roundRect(x, 150, colW, h - 186, 22);
    ctx.fill();
    const cards = id === 'done' ? done : columns[id];
    ctx.textAlign = 'left';
    ctx.fillStyle = PALETTE.muted;
    ctx.font = font(800, 32);
    ctx.fillText(title, x + 24, 196);
    ctx.textAlign = 'right';
    ctx.fillText(String(id === 'done' ? merged : cards.length), x + colW - 24, 196);
    let y = 236;
    if (id === 'done') {
      ctx.fillStyle = PALETTE.pad;
      ctx.beginPath();
      ctx.roundRect(x + 16, y, colW - 32, 150, 18);
      ctx.fill();
      ctx.textAlign = 'left';
      ctx.fillStyle = PALETTE.ink;
      ctx.font = font(800, 76);
      ctx.fillText(String(merged), x + 40, y + 62);
      ctx.fillStyle = PALETTE.muted;
      ctx.font = font(700, 28);
      ctx.fillText(`merged since ${since}`, x + 40, y + 118);
      y += 170;
    }
    if (!cards.length && id === 'doing') {
      ctx.textAlign = 'left';
      ctx.fillStyle = PALETTE.muted;
      ctx.font = font(600, 30);
      ctx.fillText(columns.review.length || done.length ? 'Empty' : 'Nothing yet', x + 28, y + 30);
    }
    for (const card of cards) {
      if (y + 112 > h - 50) break;
      ctx.fillStyle = PALETTE.pad;
      ctx.shadowColor = 'rgba(36,33,29,0.12)';
      ctx.shadowBlur = 10;
      ctx.shadowOffsetY = 3;
      ctx.beginPath();
      ctx.roundRect(x + 16, y, colW - 32, 104, 16);
      ctx.fill();
      ctx.shadowColor = 'transparent';
      ctx.fillStyle = card.accent;
      ctx.beginPath();
      ctx.roundRect(x + 16, y, 14, 104, [16, 0, 0, 16]);
      ctx.fill();
      ctx.textAlign = 'left';
      ctx.fillStyle = PALETTE.ink;
      ctx.font = font(800, 40);
      ctx.fillText(fit(ctx, card.key, colW - 80), x + 48, y + 38);
      ctx.fillStyle = PALETTE.muted;
      ctx.font = font(600, 31);
      ctx.fillText(fit(ctx, card.note ?? card.title, colW - 80), x + 48, y + 80);
      y += 120;
    }
  });
}

/**
 * Build the up-close tent for one repo.
 * @param {object} o
 * @param {string} o.project
 * @param {string} o.repo
 * @param {string} o.floor the block colour
 * @param {string} o.mateHarness
 * @param {string} [o.mateState] the Mate's scene state now
 * @param {{ since: string, tasks: object[], events: object[] }} o.history from history.js
 */
export function createTentRoom({ project, repo, floor, mateHarness, mateState, history: { since, tasks, events } }) {
  const stripe = new THREE.Color(floor).offsetHSL(0, -0.05, -0.12).getStyle();
  const accent = (h) => HARNESSES[h]?.accent ?? '#8E857A';

  const group = new THREE.Group();
  group.name = `tent-room:${project}/${repo}`;
  const disposables = [];
  const keep = (x) => (disposables.push(x), x);

  // Floor, back cloth, poles, striped bands with a valance; front and roof open.
  const base = new THREE.Mesh(keep(new RoundedBoxGeometry(W + 1.2, 0.5, D + 1.2, 4, 0.22)), keep(matte(floor)));
  base.position.y = -0.25;
  base.receiveShadow = true;
  group.add(base);
  const cloth = new THREE.Mesh(keep(new RoundedBoxGeometry(W, 0.06, D, 2, 0.03)), keep(matte(PALETTE.pad)));
  cloth.position.y = 0.03;
  cloth.receiveShadow = true;
  group.add(cloth);
  const backCloth = new THREE.Mesh(keep(new THREE.PlaneGeometry(W, EAVE)), keep(matte('#F7F1E6', { side: THREE.DoubleSide })));
  backCloth.position.set(0, EAVE / 2, -D / 2);
  backCloth.receiveShadow = true;
  group.add(backCloth);
  const poleMat = keep(matte(PALETTE.trunk));
  const poleGeo = keep(new THREE.CylinderGeometry(0.09, 0.09, EAVE + 0.5, 12));
  const pennantGeo = keep(new THREE.ShapeGeometry(new THREE.Shape([new THREE.Vector2(0, 0.2), new THREE.Vector2(0.62, 0), new THREE.Vector2(0, -0.2)])));
  const pennantMat = keep(new THREE.MeshStandardMaterial({ color: PALETTE.coral, roughness: 0.9, side: THREE.DoubleSide }));
  for (const sx of [-1, 1]) {
    for (const sz of [-1, 1]) {
      const pole = new THREE.Mesh(poleGeo, poleMat);
      pole.position.set(sx * (W / 2 - 0.1), (EAVE + 0.5) / 2, sz * (D / 2 - 0.1));
      pole.castShadow = true;
      group.add(pole);
      const pennant = new THREE.Mesh(pennantGeo, pennantMat);
      pennant.position.set(sx * (W / 2 - 0.1), EAVE + 0.3, sz * (D / 2 - 0.1));
      group.add(pennant);
    }
  }
  const stripeMats = [keep(matte(stripe)), keep(matte(PALETTE.pad))];
  const bandGeo = keep(new THREE.BoxGeometry(1, 0.62, 0.08));
  const valMat = keep(new THREE.MeshStandardMaterial({ color: stripe, roughness: 0.95, side: THREE.DoubleSide }));
  // Back band runs along x; side bands along z (the front stays open).
  for (const [len, cx, cz, ry] of [[W, 0, -D / 2, 0], [D, -W / 2, 0, Math.PI / 2], [D, W / 2, 0, -Math.PI / 2]]) {
    const n = Math.round(len / 0.7);
    for (let i = 0; i < n; i++) {
      const b = new THREE.Mesh(bandGeo, stripeMats[i % 2]);
      b.scale.x = len / n;
      const u = -len / 2 + (i + 0.5) * (len / n);
      b.position.set(cx + (ry ? 0 : u), EAVE - 0.31, cz + (ry ? u : 0));
      b.rotation.y = ry;
      b.castShadow = true;
      group.add(b);
    }
    const val = new THREE.Mesh(keep(valanceGeometry(len, 0.3, 0.5)), valMat);
    val.position.set(cx, EAVE - 0.6, cz);
    val.rotation.y = ry;
    group.add(val);
  }

  // Zone labels painted on the cloth.
  function floorLabel(text, x, z, height = 0.42) {
    const { texture, aspect } = textTexture(text, { font: `800 110px ${FONT}`, color: 'rgba(107,101,92,0.85)', pad: 18, letterSpacing: 6 });
    keep(texture);
    const m = new THREE.Mesh(keep(new THREE.PlaneGeometry(height * aspect, height)), keep(new THREE.MeshBasicMaterial({ map: texture, transparent: true })));
    m.rotation.x = -Math.PI / 2;
    m.position.set(x, 0.07, z);
    group.add(m);
  }
  floorLabel('TASK BOARD', BOARD.x - 0.7, BOARD.z + 2.3);
  floorLabel('ROUND TABLE', TABLE.x, TABLE.z + TABLE.r + 1.15);

  // Task board on two legs, angled toward the table.
  const boardH = BOARD.w * (BOARD_PX[1] / BOARD_PX[0]);
  const boardGroup = new THREE.Group();
  boardGroup.position.set(BOARD.x, 0, BOARD.z);
  boardGroup.rotation.y = BOARD.ry;
  group.add(boardGroup);
  const legGeo = keep(new THREE.CylinderGeometry(0.06, 0.07, 1, 10));
  const legMat = keep(matte('#8E857A'));
  for (const s of [-1, 1]) {
    const leg = new THREE.Mesh(legGeo, legMat);
    leg.scale.y = 1.0 + boardH;
    leg.position.set(s * (BOARD.w / 2 - 0.3), (1.0 + boardH) / 2, -0.06);
    leg.castShadow = true;
    boardGroup.add(leg);
  }
  const frame = new THREE.Mesh(keep(new RoundedBoxGeometry(BOARD.w + 0.2, boardH + 0.2, 0.1, 3, 0.05)), keep(matte('#E9E3D6')));
  frame.position.y = 0.95 + boardH / 2;
  frame.castShadow = true;
  boardGroup.add(frame);
  const boardTex = new LiveTexture(BOARD_PX[0], BOARD_PX[1], (ctx, state, w, h) => drawTaskBoard(ctx, JSON.parse(state), w, h));
  keep(boardTex.texture);
  const face = new THREE.Mesh(keep(new THREE.PlaneGeometry(BOARD.w, boardH)), keep(new THREE.MeshBasicMaterial({ map: boardTex.texture, transparent: true, toneMapped: false })));
  face.position.set(0, frame.position.y, 0.056);
  boardGroup.add(face);

  // Round table with a laptop at each seat.
  const tableMat = keep(matte('#E8D9C2'));
  const top = new THREE.Mesh(keep(new THREE.CylinderGeometry(TABLE.r, TABLE.r, 0.1, 64)), tableMat);
  top.position.set(TABLE.x, TABLE.h, TABLE.z);
  top.castShadow = top.receiveShadow = true;
  group.add(top);
  const stem = new THREE.Mesh(keep(new THREE.CylinderGeometry(0.16, 0.5, TABLE.h, 24)), keep(matte('#B8A68E')));
  stem.position.set(TABLE.x, TABLE.h / 2, TABLE.z);
  stem.castShadow = true;
  group.add(stem);

  // A plant in the back-right corner.
  const pot = new THREE.Mesh(keep(new THREE.CylinderGeometry(0.42, 0.32, 0.7, 20)), keep(matte('#D9A88C')));
  pot.position.set(W / 2 - 1.2, 0.35, -D / 2 + 1.1);
  pot.castShadow = true;
  group.add(pot);
  const leaves = new THREE.Mesh(keep(new THREE.SphereGeometry(0.85, 20, 14)), keep(matte(PALETTE.leaf[0])));
  leaves.position.set(W / 2 - 1.2, 1.45, -D / 2 + 1.1);
  leaves.scale.y = 1.2;
  leaves.castShadow = true;
  group.add(leaves);


  // ---- The cast -----------------------------------------------------------------
  // Each event takes STEP seconds of the room's clock, oldest first, so a quiet
  // night and a busy hour read at the same pace; the feed and the scrub bar
  // carry the real times. A crew takes a seat when the Mate spawns it and
  // leaves just before the first event after its task closed. Seats are
  // reused, so the table holds as many crews as were ever there at once.
  const n = events.length;
  const tau = (i) => LEAD + i * STEP;
  const duration = n ? tau(n - 1) + STEP : ROOM_LOOP;
  const ats = events.map((e) => e.at);
  // The first event after `iso`, by binary search; n when none came yet.
  function after(iso) {
    let lo = 0;
    let hi = n;
    while (lo < hi) {
      const mid = (lo + hi) >> 1;
      if (ats[mid] <= iso) lo = mid + 1;
      else hi = mid;
    }
    return lo;
  }
  const byCrew = new Map(tasks.map((t) => [t.crew, t]));
  const stays = []; // { task, from, to, seat }, by arrival
  const stayOf = new Map();
  events.forEach((e, i) => {
    const t = byCrew.get(e.crew);
    if ((e.kind !== 'spawn' && e.kind !== 'review') || !t || stayOf.has(e.crew)) return;
    let to = Infinity;
    if (t.closed) {
      const k = after(t.closed_at || t.merged_at || e.at);
      to = k < n ? tau(k) - 0.6 : duration - 0.4;
    }
    const st = { task: t, from: tau(i), to: Math.max(to, tau(i) + 1.6), seat: 0 };
    stays.push(st);
    stayOf.set(e.crew, st);
  });
  const freeAt = [];
  for (const st of stays) {
    let k = freeAt.findIndex((f) => f <= st.from);
    if (k < 0) k = freeAt.push(0) - 1;
    st.seat = k;
    freeAt[k] = st.to + 0.5;
  }

  // Seats on the far arc of the table, turned partly to the room; past seven,
  // a second arc behind the first.
  const mascots = new Map();
  const seats = new Map(); // mascot id → { project, kind, crew?, role, harness, model, effort, text }
  const padGeo = keep(new THREE.CylinderGeometry(1, 1, 0.08, 40));
  const padMat = keep(matte('#EFE7DA'));
  const laptopBase = keep(new RoundedBoxGeometry(0.62, 0.04, 0.42, 2, 0.015));
  const laptopLid = keep(new RoundedBoxGeometry(0.62, 0.4, 0.03, 2, 0.012));
  const laptopMat = keep(matte('#5A5148'));
  const lidMat = keep(matte('#5A5148', { emissive: '#F2EFE8', emissiveIntensity: 0.15 }));
  function pad(x, z, size, mat) {
    const m = new THREE.Mesh(padGeo, mat);
    m.scale.set(0.62 * size + 0.06, 1, 0.62 * size + 0.06);
    m.position.set(x, 0.1, z);
    m.receiveShadow = true;
    group.add(m);
    return m;
  }
  const PER_ARC = 7;
  const chairs = freeAt.map((_, k) => {
    const arc = Math.floor(k / PER_ARC);
    const count = Math.min(PER_ARC, freeAt.length - arc * PER_ARC);
    const j = k % PER_ARC;
    const phi = count === 1 ? 0 : -1.2 + (2.4 * j) / (count - 1);
    const R = TABLE.r + 0.7 + arc * 1.25;
    const x = TABLE.x + R * Math.sin(phi);
    const z = TABLE.z - R * Math.cos(phi);
    const own = stays.filter((st) => st.seat === k);
    const look = own[0].task.harness;
    const m = new AgentMascot({ ...harnessLook(look), size: CREW_SIZE, position: [x, 0.14, z], heading: -0.55 * phi, phase: (k * 1.37) % 4 });
    m.id = `${project}/seat-${k}`;
    m.play('idle', { at: 0 });
    group.add(m.root);
    mascots.set(m.id, m);
    let laptop = null;
    if (!arc) {
      // Laptop in front of the seat, screen toward it.
      laptop = new THREE.Group();
      laptop.position.set(TABLE.x + (TABLE.r - 0.55) * Math.sin(phi), TABLE.h + 0.07, TABLE.z - (TABLE.r - 0.55) * Math.cos(phi));
      laptop.rotation.y = -phi;
      laptop.add(new THREE.Mesh(laptopBase, laptopMat));
      const lid = new THREE.Mesh(laptopLid, lidMat);
      lid.position.set(0, 0.2, 0.2);
      lid.rotation.x = 0.25;
      laptop.add(lid);
      laptop.traverse((o) => (o.castShadow = true));
      group.add(laptop);
    }
    return { m, own, look, laptop, pad: pad(x, z, CREW_SIZE, padMat) };
  });
  const mate = new AgentMascot({ ...harnessLook(mateHarness), size: MATE_SIZE, position: [-1.3, 0.14, 1.0], heading: 0.35, phase: 0.5 });
  mate.id = `${project}:mate`;
  group.add(mate.root);
  mascots.set(mate.id, mate);
  seats.set(mate.id, { project, kind: 'mate', harness: mateHarness });
  pad(-1.3, 1.0, MATE_SIZE, keep(matte('#F2A285')));
  const mascotOf = (crew) => chairs[stayOf.get(crew)?.seat]?.m;

  // ---- Playing it out -------------------------------------------------------------
  // A mascot's plays must go in time order, which the event order gives.
  const feed = []; // { at, clock, iso, who, kind, text, crew }
  const clockFmt = new Intl.DateTimeFormat(undefined, { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false });
  const clock = (iso) => clockFmt.format(new Date(iso));
  const label = (crew) => {
    if (crew === 'mate') return 'Mate';
    const t = byCrew.get(crew);
    return `${harnessName(t?.harness)} · ${ticket(crew).key}${t?.role === 'review' ? ' review' : ''}`;
  };
  const say = (m, text, at) => m.say(text, { at, hold: 0.7 });
  const pass = (from, to, at, then, receiverThen) => to && handoff(from, to, at, { then, receiverThen, turn: 0.6 });
  const notes = new Map(); // crew → [room time, words] for its card on the board
  const note = (crew, at, words) => notes.set(crew, [...(notes.get(crew) ?? []), [at, words]]);
  const merges = new Map(); // ship crew → { at, pr }
  const reviewsOf = new Map(); // ship crew → reviewer stays

  mate.play('idle', { at: 0 });
  if (!n) {
    mate.play('sleeping', { at: 0 });
    feed.push({ at: 0, clock: '', who: 'Mate', kind: 'sleep', text: 'No task has worked on this repo yet.' });
  }
  events.forEach((ev, i) => {
    const at = tau(i);
    const m = ev.who === 'mate' ? mate : mascotOf(ev.who);
    const c = mascotOf(ev.crew);
    feed.push({ at, clock: clock(ev.at), iso: ev.at, who: label(ev.who), kind: ev.kind, text: ev.text, crew: ev.crew });
    switch (ev.kind) {
      case 'spawn':
      case 'review':
        say(mate, ev.short, at);
        pass(mate, c, at + 0.8, 'idle', 'working');
        note(ev.crew, at, ev.kind === 'review' ? 'reviewing' : 'starting');
        if (ev.kind === 'review' && ev.of) reviewsOf.set(ev.of, [...(reviewsOf.get(ev.of) ?? []), stayOf.get(ev.crew)]);
        break;
      case 'status':
        m?.play('working', { at });
        if (m) say(m, ev.short, at);
        note(ev.crew, at, ev.short);
        break;
      case 'ask':
        m?.play('thinking', { at });
        if (m) say(m, ev.short, at);
        if (m) pass(m, mate, at + 1.0, 'thinking', 'thinking');
        note(ev.crew, at, 'asking the Mate');
        break;
      case 'answer':
        say(mate, ev.short, at);
        pass(mate, c, at + 1.0, 'idle', 'working');
        break;
      case 'handback':
      case 'verdict':
        m?.play(ev.kind === 'verdict' && !ev.pass ? 'thinking' : 'completed', { at });
        if (m) say(m, ev.short, at);
        if (m) pass(m, mate, at + 1.0, 'idle', 'thinking');
        note(ev.crew, at, ev.kind === 'verdict' ? (ev.pass ? 'review passed' : 'changes requested') : 'handed back');
        if (ev.kind === 'verdict') {
          const ship = byCrew.get(ev.crew) && events.find((e) => e.kind === 'review' && e.crew === ev.crew)?.of;
          if (ship) note(ship, at, ev.pass ? 'review passed' : 'fixing after review');
        }
        break;
      case 'merge':
        mate.play('completed', { at });
        say(mate, ev.short, at);
        c?.play('completed', { at: at + 0.3 });
        merges.set(ev.crew, { at, pr: ev.pr });
        break;
      case 'close':
        say(mate, ev.short, at);
        break;
    }
  });
  // Live, the open crews and the Mate end posed as they are now.
  const end = Math.max(0, duration - 0.2);
  for (const st of stays) if (!st.task.closed && st.task.state) mascotOf(st.task.crew)?.play(poseOf(st.task.state), { at: end });
  if (n && mateState) mate.play(mateState === 'asleep' ? 'sleeping' : poseOf(mateState), { at: end });

  // ---- The board --------------------------------------------------------------------
  const ships = stays.filter((st) => st.task.role === 'ship');
  const noteAt = (crew, t) => {
    const list = notes.get(crew) ?? [];
    let words = '';
    for (const [at, w] of list) if (at <= t) words = w;
    return words;
  };
  const sinceLabel = new Intl.DateTimeFormat(undefined, { day: '2-digit', month: '2-digit' }).format(new Date(since));
  function boardState(t) {
    const columns = { doing: [], review: [] };
    const done = [];
    for (const st of ships) {
      if (st.from > t) continue;
      const crew = st.task.crew;
      const card = { key: ticket(crew).key, title: st.task.text, accent: accent(st.task.harness) };
      const mg = merges.get(crew);
      if (mg && mg.at <= t) done.push({ ...card, note: mg.pr ? `PR #${mg.pr}` : 'merged', at: mg.at });
      else if (t < st.to) {
        const inReview = (reviewsOf.get(crew) ?? []).some((r) => r && r.from <= t && t < r.to);
        (inReview ? columns.review : columns.doing).push({ ...card, note: noteAt(crew, t) || card.title });
      }
    }
    done.sort((a, b) => b.at - a.at);
    return JSON.stringify({ repo, color: floor, columns, done: done.map(({ at, ...c }) => c), merged: done.length, since: sinceLabel });
  }

  function update(t) {
    for (const ch of chairs) {
      const st = ch.own.find((o) => o.from <= t && t < o.to);
      ch.m.root.visible = ch.pad.visible = Boolean(st);
      if (ch.laptop) ch.laptop.visible = Boolean(st);
      if (st) {
        if (st.task.harness !== ch.look) {
          const look = harnessLook(st.task.harness);
          ch.m.setHarness(look.harness, look);
          ch.look = st.task.harness;
        }
        // Pop in on arrival, shrink away on leaving.
        const k = Math.min(1, (t - st.from) / 0.4, (st.to - t) / 0.4);
        ch.m.root.scale.setScalar(Math.max(0.001, 1 - (1 - k) ** 3));
        const x = st.task;
        seats.set(ch.m.id, { project, kind: 'crew', crew: x.crew, role: x.role, harness: x.harness, model: x.model, effort: x.effort, text: x.text });
        ch.m.update(t);
      } else seats.delete(ch.m.id);
    }
    mate.update(t);
    boardTex.set(boardState(t));
    group.updateMatrixWorld(true);
  }

  // The feed entry showing at room time t: the last one at or before it.
  function indexAt(t) {
    let lo = 0;
    let hi = feed.length;
    while (lo < hi) {
      const mid = (lo + hi) >> 1;
      if (feed[mid].at <= t) lo = mid + 1;
      else hi = mid;
    }
    return lo - 1;
  }

  const shipCount = ships.length;
  const target = new THREE.Vector3(-0.8, 1.2, -1.0);
  const view = target.clone().add(new THREE.Vector3(Math.sin(0.1), 0.62, Math.cos(0.1)).setLength(27));
  return {
    group,
    project,
    repo,
    mascots,
    seats,
    events: feed,
    duration,
    since: n ? sinceLabel : '',
    summary: n ? `${shipCount} task${shipCount === 1 ? '' : 's'} · ${merges.size} merged · since ${sinceLabel}` : 'nothing has worked on this repo yet',
    view: { target, position: view },
    update,
    indexAt,
    dispose() {
      for (const m of mascots.values()) m.dispose();
      for (const d of disposables) d.dispose();
    },
  };
}
