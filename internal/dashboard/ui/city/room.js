import * as THREE from './vendor/three.js';
import { RoundedBoxGeometry } from './vendor/three.js';
import { AgentMascot, handoff } from './agent-mascot.js';
import { HARNESSES } from './harnesses.js';
import { PALETTE, LiveTexture, matte, textTexture, valanceGeometry, harnessLook, harnessName, formatTokens } from './city.js';

// One repo's tent seen up close: a big open tent (front and roof cut away)
// with a task board, a round table where the crews sit with their laptops,
// and the project's Mate between the two. It replays one real episode that
// episode.js read from the dashboard: the Mate hands a task to a ship crew,
// the crew reports status, asks the Mate and gets its answers, hands back, a
// reviewer reads the diff and gives its verdict, the Mate merges. Bubbles
// show the lines shortened; `events` carries them in full with their real
// times for the page's feed. update(t) is a pure function of t, so the
// replay loops by calling it with t modulo `duration`.

export const ROOM_LOOP = 18; // a room with nothing to replay
const STEP = 3.0; // seconds per replayed event
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

function summary(ep) {
  const s = ep.ship;
  return [s.id, `${formatTokens(s.tokens)} tokens`, `${s.turns} turns`, s.questions ? `${s.questions} question${s.questions === 1 ? '' : 's'}` : null, s.pr ? `PR #${s.pr}` : null].filter(Boolean).join(' · ');
}

function fit(ctx, text, w) {
  if (ctx.measureText(text).width <= w) return text;
  let s = text;
  while (s && ctx.measureText(`${s}…`).width > w) s = s.slice(0, -1);
  return `${s.trimEnd()}…`;
}

function drawTaskBoard(ctx, { repo, color, columns, done, finished, merged, days }, w, h) {
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
    ctx.fillText(String(id === 'done' ? finished : cards.length), x + colW - 24, 196);
    let y = 236;
    if (id === 'done') {
      ctx.fillStyle = PALETTE.pad;
      ctx.beginPath();
      ctx.roundRect(x + 16, y, colW - 32, 150, 18);
      ctx.fill();
      ctx.textAlign = 'left';
      ctx.fillStyle = PALETTE.ink;
      ctx.font = font(800, 76);
      ctx.fillText(String(finished), x + 40, y + 62);
      ctx.fillStyle = PALETTE.muted;
      ctx.font = font(700, 28);
      ctx.fillText(`tasks done · last ${days} days`, x + 40, y + 118);
      y += 170;
      ctx.fillStyle = PALETTE.ink;
      ctx.font = font(700, 30);
      ctx.fillText(`${merged} merged on ${fit(ctx, repo, colW - 220)}`, x + 28, y + 14);
      y += 50;
    }
    if (!cards.length && id === 'doing') {
      ctx.textAlign = 'left';
      ctx.fillStyle = PALETTE.muted;
      ctx.font = font(600, 30);
      ctx.fillText(columns.review.length || done.length ? 'Empty' : 'No crew yet', x + 28, y + 30);
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
 * @param {object|null} o.episode from episode.js, or null when the repo has no merged task
 * @param {{ id: string, text: string }[]} o.recent earlier merged tasks on the repo
 * @param {number} o.finished tasks done in the project over the board's days
 * @param {number} o.merged tasks merged on the repo
 * @param {number} o.days
 */
export function createTentRoom({ project, repo, floor, mateHarness, episode: ep, recent = [], finished = 0, merged = 0, days = 7 }) {
  const crews = ep ? [ep.ship, ...ep.reviewers] : [];
  const key = ep ? ticket(ep.ship.id).key : '';
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

  // Mascots: crews on the far arc of the table, turned partly to the room;
  // the Mate between the board and the table.
  const mascots = new Map();
  const seats = new Map(); // id → { project, kind, crew?, role, harness, model, effort }
  const padGeo = keep(new THREE.CylinderGeometry(1, 1, 0.08, 40));
  const padMat = keep(matte('#EFE7DA'));
  const laptopBase = keep(new RoundedBoxGeometry(0.62, 0.04, 0.42, 2, 0.015));
  const laptopLid = keep(new RoundedBoxGeometry(0.62, 0.4, 0.03, 2, 0.012));
  const laptopMat = keep(matte('#5A5148'));
  function seat(id, s, kind, x, z, size, heading) {
    const m = new AgentMascot({ ...harnessLook(s.harness), size, position: [x, 0.14, z], heading, phase: (mascots.size * 1.37) % 4 });
    m.id = id;
    group.add(m.root);
    const pad = new THREE.Mesh(padGeo, kind === 'mate' ? keep(matte('#F2A285')) : padMat);
    pad.scale.set(0.62 * size + 0.06, 1, 0.62 * size + 0.06);
    pad.position.set(x, 0.1, z);
    pad.receiveShadow = true;
    group.add(pad);
    mascots.set(id, m);
    seats.set(id, kind === 'mate' ? { project, kind, harness: s.harness } : { project, kind, crew: s.id, role: s.role, harness: s.harness, model: s.model, effort: s.effort, text: s.text });
    return m;
  }
  const n = crews.length;
  const crewMascots = crews.map((c, k) => {
    const phi = n === 1 ? 0 : -1.2 + (2.4 * k) / (n - 1);
    const R = TABLE.r + 0.7;
    const x = TABLE.x + R * Math.sin(phi);
    const z = TABLE.z - R * Math.cos(phi);
    // Laptop in front of the seat, screen toward it.
    const lx = TABLE.x + (TABLE.r - 0.55) * Math.sin(phi);
    const lz = TABLE.z - (TABLE.r - 0.55) * Math.cos(phi);
    const laptop = new THREE.Group();
    laptop.position.set(lx, TABLE.h + 0.07, lz);
    laptop.rotation.y = -phi;
    const lb = new THREE.Mesh(laptopBase, laptopMat);
    laptop.add(lb);
    const lid = new THREE.Mesh(laptopLid, keep(matte('#5A5148', { emissive: accent(c.harness), emissiveIntensity: 0.15 })));
    lid.position.set(0, 0.2, 0.2);
    lid.rotation.x = 0.25;
    laptop.add(lid);
    laptop.traverse((o) => (o.castShadow = true));
    group.add(laptop);
    return seat(`${project}/${c.id}`, c, 'crew', x, z, CREW_SIZE, -0.55 * phi);
  });
  const mate = seat(`${project}:mate`, { harness: mateHarness }, 'mate', -1.3, 1.0, MATE_SIZE, 0.35);

  // ---- Replaying the episode -------------------------------------------------
  // One event per STEP seconds, in the episode's real order; the feed carries
  // the real clock times. A mascot's plays must go in time order, which the
  // event order gives, and each hand-off lands before the next event.
  const events = []; // feed entries: { at, clock, who, kind, text }
  const moves = []; // { at, to: doing|review|done, note }
  const byId = new Map(crews.map((c, k) => [c.id, crewMascots[k]]));
  const shipM = ep && byId.get(ep.ship.id);
  const label = (id) => {
    if (id === 'mate') return 'Mate';
    const c = id === 'ship' ? ep.ship : crews.find((x) => x.id === id);
    return c.role === 'review' ? `${harnessName(c.harness)} · reviewer` : `${harnessName(c.harness)} · ${key}`;
  };
  const clockFmt = new Intl.DateTimeFormat(undefined, { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false });
  const clock = (iso) => clockFmt.format(new Date(iso));
  const say = (m, text, at) => m.say(text, { at, hold: 0.7 });
  const pass = (from, to, at, then, receiverThen) => handoff(from, to, at, { then, receiverThen, turn: 0.6 });
  let duration = ROOM_LOOP;

  mate.play('idle', { at: 0 });
  if (!ep) {
    mate.play('sleeping', { at: 0 });
    events.push({ at: 0, clock: '', who: 'Mate', kind: 'sleep', text: 'No task has been merged on this repo yet, so there is nothing to replay.' });
  } else {
    shipM.play('sleeping', { at: 0 });
    for (const r of ep.reviewers) byId.get(r.id).play('sleeping', { at: 0 });
    ep.events.forEach((ev, i) => {
      const at = 1 + i * STEP;
      const m = ev.who === 'mate' ? mate : ev.who === 'ship' ? shipM : byId.get(ev.who);
      events.push({ at, clock: clock(ev.at), who: label(ev.who), kind: ev.kind, text: ev.text });
      switch (ev.kind) {
        case 'spawn':
          say(mate, ev.short, at);
          pass(mate, shipM, at + 0.8, 'idle', 'working');
          moves.push({ at: at + 2.3, to: 'doing' });
          break;
        case 'status':
          m.play('working', { at });
          say(m, ev.short, at);
          break;
        case 'ask':
          shipM.play('thinking', { at });
          say(shipM, ev.short, at);
          pass(shipM, mate, at + 1.0, 'thinking', 'thinking');
          break;
        case 'answer':
          say(mate, ev.short, at);
          pass(mate, shipM, at + 1.0, 'idle', 'working');
          break;
        case 'handback':
          m.play('completed', { at });
          say(m, ev.short, at);
          pass(m, mate, at + 1.0, m === shipM ? 'idle' : 'thinking', 'thinking');
          break;
        case 'review': {
          const r = byId.get(ev.reviewer) ?? crewMascots[1];
          say(mate, ev.short, at);
          pass(mate, r, at + 0.8, 'idle', 'working');
          moves.push({ at: at + 2.3, to: 'review', note: `${harnessName(crews.find((c) => byId.get(c.id) === r).harness)} reviewing` });
          break;
        }
        case 'verdict': {
          const ok = ev.pass;
          m.play(ok ? 'completed' : 'thinking', { at });
          say(m, ev.short, at);
          pass(m, mate, at + 1.0, 'sleeping', 'thinking');
          moves.push({ at: at + 2.5, to: ok ? 'review' : 'doing', note: ok ? 'review passed' : 'fixing after review' });
          break;
        }
        case 'merge':
          mate.play('completed', { at });
          say(mate, ev.short, at);
          moves.push({ at: at + 0.3, to: 'done', note: ep.ship.pr ? `PR #${ep.ship.pr}` : 'merged' });
          break;
      }
    });
    duration = 1 + (ep.events.length - 1) * STEP + STEP + 3;
  }

  const recentCards = recent.map((r) => ({ key: ticket(r.id).key, title: r.text, accent: '#B8AFA3' }));
  function boardState(t) {
    const columns = { doing: [], review: [] };
    const done = [];
    const last = moves.filter((m) => m.at <= t).at(-1);
    if (ep && last) {
      const card = { key, title: ep.ship.text, accent: accent(ep.ship.harness), note: last.note };
      (last.to === 'done' ? done : columns[last.to]).push(card);
    }
    done.push(...recentCards);
    return JSON.stringify({ repo, color: floor, columns, done, finished, merged, days });
  }

  function update(t) {
    boardTex.set(boardState(t));
    for (const m of mascots.values()) m.update(t);
    group.updateMatrixWorld(true);
  }

  const target = new THREE.Vector3(-0.8, 1.2, -1.0);
  const view = target.clone().add(new THREE.Vector3(Math.sin(0.1), 0.62, Math.cos(0.1)).setLength(27));
  return {
    group,
    project,
    repo,
    mascots,
    seats,
    events,
    duration,
    summary: ep ? summary(ep) : 'no merged task on this repo yet',
    view: { target, position: view },
    update,
    dispose() {
      for (const m of mascots.values()) m.dispose();
      for (const d of disposables) d.dispose();
    },
  };
}
