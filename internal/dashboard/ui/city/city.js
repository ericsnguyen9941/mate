import * as THREE from './vendor/three.js';
import { RoundedBoxGeometry, RoomEnvironment } from './vendor/three.js';
import { AgentMascot } from './agent-mascot.js';
import { HARNESSES } from './harnesses.js';
import { CHAT_STYLE, GAP as CHAT_GAP, lineFor } from './chatter.js';

// Agent City: the workspace as a floating diorama island laid out as a grid
// of blocks (3×3, or 5×5 past eight projects). The centre block is City Hall,
// where the captain works; its mailbox holds the news the Mates opened for
// the captain, its tower screen the workspace at a glance. Every other block
// is a project (or a park), its Mate on the coral pad in the plaza. Each repo
// in the project's project.yaml is an open market tent where the crews on it
// work in view of the street; a small sign in front gives its merged task
// count, which a merge ticks up. A signboard at each project's front corner
// shows the board's last days: tokens, finished tasks, passed reviews.
//
// The city is drawn once from the workspace's projects and repos, then kept
// live by the setters (setCrews, setMate, setMerged, setBoard, setHall,
// setMailbox), which app.js calls whenever the dashboard's API moves. Every
// change is scheduled on the page clock t that update(t) is called with, so a
// mascot's play() calls always go forward in time.

export const PALETTE = {
  paper: '#F2EFE8',
  island: '#ECE6D9',
  cliff: '#D9CEBB',
  road: '#D2CABB',
  line: '#FAF7F0',
  coral: '#E77859',
  ink: '#24211D',
  muted: '#6B655C',
  glass: '#7D8999',
  door: '#5A5148',
  pad: '#FFFDF8',
  park: '#CFE4C0',
  hall: '#F4DCCF',
  leaf: ['#A6C99D', '#BBD5A6', '#92BE93'],
  trunk: '#B89B7D',
  walls: ['#FBF7EF', '#F3E4D4', '#E2EBF3', '#E7EFE3'],
  floors: ['#BFE3D2', '#C9DBF5', '#F6E3B4', '#E3D3F0', '#F8D9C4', '#D5E8C4', '#F3D1DC', '#CDE3E8'],
};

const FONT = 'ui-sans-serif, system-ui, sans-serif';
const BLOCK = 11; // block width
const ROAD = 3; // street width
const PITCH = BLOCK + ROAD; // block centre to block centre
const KERB = 0.14; // block height above the street
const EAVE = 2.0; // tent eaves; high enough to see crews under the valance
const RISE = 1.05; // eaves to ridge
const POP = 0.7; // seconds a merge's sign pop lasts
const SIGN_W = 1.25; // repo sign in front of each tent
const SIGN_H = SIGN_W * (220 / 520);
const MATE_SIZE = 1.2;
const CREW_SIZE = 0.9;
const SPOTS = 16; // places a tent has for crews; more crews than this stand in the street
const TAU = Math.PI * 2;
const ease = (x) => x * x * x * (x * (x * 6 - 15) + 10);
const clamp01 = (x) => Math.min(1, Math.max(0, x));
const backOut = (x) => 1 + 2.2 * (x - 1) ** 3 + 1.2 * (x - 1) ** 2;

// A harness the mascot has no preset for (grok, a new one) is a plain grey
// mascot with a ring for a logo.
const GENERIC = {
  color: '#ECE8E1',
  logoColor: '#8E857A',
  accent: '#8E857A',
  logo: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 3a9 9 0 1 0 0 18a9 9 0 1 0 0-18zm0 4a5 5 0 1 1 0 10a5 5 0 1 1 0-10z"/></svg>',
};
export const harnessLook = (harness) => (HARNESSES[harness] ? { harness } : { harness, ...GENERIC });
export const harnessName = (harness) => HARNESSES[harness]?.name ?? (harness ? harness[0].toUpperCase() + harness.slice(1) : 'Agent');

// Scene state (docs/mvp.md M5) → the mascot's pose.
const POSE = {
  arriving: 'idle',
  at_desk_working: 'working',
  walking_to_ceo: 'thinking',
  waiting_at_ceo: 'thinking',
  waiting_review: 'idle',
  idle: 'idle',
  blocked: 'thinking',
  asleep: 'sleeping',
  leaving: 'idle',
};
export const poseOf = (state) => POSE[state] ?? 'idle';
// What a crew says as it moves into a state, so a change catches the eye.
const SAYS = { at_desk_working: 'on it', waiting_at_ceo: 'asking the Mate', walking_to_ceo: 'asking the Mate', waiting_review: 'in review', blocked: 'blocked', leaving: 'done' };

export function lights(scene, renderer, extent = 27) {
  const pmrem = new THREE.PMREMGenerator(renderer);
  scene.environment = pmrem.fromScene(new RoomEnvironment(), 0.04).texture;
  scene.environmentIntensity = 0.5;
  scene.add(new THREE.HemisphereLight('#ffffff', '#d9d0bf', 1.15));
  const sun = new THREE.DirectionalLight('#fff3e2', 2.3);
  sun.position.set(20, 40, 26).multiplyScalar(extent / 27);
  sun.castShadow = true;
  sun.shadow.mapSize.set(4096, 4096);
  Object.assign(sun.shadow.camera, { left: -extent, right: extent, top: extent, bottom: -extent, near: 1, far: 120 * (extent / 27) });
  sun.shadow.bias = -0.0003;
  sun.shadow.normalBias = 0.02;
  scene.add(sun);
  return pmrem;
}

export function rng(seed) {
  let s = seed >>> 0;
  return () => ((s = (s * 1664525 + 1013904223) >>> 0) / 2 ** 32);
}

export function hash(text) {
  let h = 2166136261;
  for (const c of text) h = Math.imul(h ^ c.charCodeAt(0), 16777619);
  return h >>> 0;
}

function canvas(w, h) {
  const c = document.createElement('canvas');
  c.width = w;
  c.height = h;
  return c;
}

// A text label as a texture; returns it with its width/height ratio.
export function textTexture(text, { font = `700 96px ${FONT}`, color = PALETTE.ink, background = null, pad = 36, radius = 28, letterSpacing = 0 } = {}) {
  const probe = canvas(1, 1).getContext('2d');
  probe.font = font;
  probe.letterSpacing = `${letterSpacing}px`;
  const m = probe.measureText(text);
  const h = Math.ceil(m.actualBoundingBoxAscent + m.actualBoundingBoxDescent + 2 * pad);
  const w = Math.ceil(m.width + 2 * pad);
  const c = canvas(w, h);
  const ctx = c.getContext('2d');
  if (background) {
    ctx.fillStyle = background;
    ctx.beginPath();
    ctx.roundRect(0, 0, w, h, radius);
    ctx.fill();
  }
  ctx.font = font;
  ctx.letterSpacing = `${letterSpacing}px`;
  ctx.fillStyle = color;
  ctx.fillText(text, pad, pad + m.actualBoundingBoxAscent);
  const texture = new THREE.CanvasTexture(c);
  texture.colorSpace = THREE.SRGBColorSpace;
  texture.anisotropy = 8;
  return { texture, aspect: w / h };
}

export function roundedRect(w, h, r) {
  const s = new THREE.Shape();
  s.moveTo(-w / 2 + r, -h / 2);
  s.lineTo(w / 2 - r, -h / 2);
  s.quadraticCurveTo(w / 2, -h / 2, w / 2, -h / 2 + r);
  s.lineTo(w / 2, h / 2 - r);
  s.quadraticCurveTo(w / 2, h / 2, w / 2 - r, h / 2);
  s.lineTo(-w / 2 + r, h / 2);
  s.quadraticCurveTo(-w / 2, h / 2, -w / 2, h / 2 - r);
  s.lineTo(-w / 2, -h / 2 + r);
  s.quadraticCurveTo(-w / 2, -h / 2, -w / 2 + r, -h / 2);
  return s;
}

export const matte = (color, extra = {}) => new THREE.MeshStandardMaterial({ color, roughness: 0.85, metalness: 0, ...extra });

// Instanced meshes collected while building, flushed into one mesh each.
class Batch {
  constructor(geometry, material, name) {
    Object.assign(this, { geometry, material, name, items: [] });
  }
  add(position, { rotation = [0, 0, 0], scale = [1, 1, 1], color } = {}) {
    this.items.push({ position, rotation, scale, color });
  }
  build({ cast = true, receive = true } = {}) {
    const mesh = new THREE.InstancedMesh(this.geometry, this.material, Math.max(1, this.items.length));
    mesh.name = this.name;
    mesh.count = this.items.length;
    const m = new THREE.Matrix4();
    const q = new THREE.Quaternion();
    const e = new THREE.Euler();
    this.items.forEach((it, i) => {
      m.compose(new THREE.Vector3(...it.position), q.setFromEuler(e.set(...it.rotation)), new THREE.Vector3(...it.scale));
      mesh.setMatrixAt(i, m);
      if (it.color) mesh.setColorAt(i, new THREE.Color(it.color));
    });
    mesh.castShadow = cast;
    mesh.receiveShadow = receive;
    return mesh;
  }
}

// A canvas texture redrawn only when its key changes.
export class LiveTexture {
  constructor(w, h, draw) {
    this.canvas = canvas(w, h);
    this.draw = draw;
    this.key = null;
    this.texture = new THREE.CanvasTexture(this.canvas);
    this.texture.colorSpace = THREE.SRGBColorSpace;
    this.texture.anisotropy = 8;
  }
  set(key) {
    if (key === this.key) return;
    this.key = key;
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    this.draw(ctx, key, this.canvas.width, this.canvas.height);
    this.texture.needsUpdate = true;
  }
}

// Scalloped valance strip of width w hanging down from y = 0, in the xy plane.
export function valanceGeometry(w, depth = 0.2, scallop = 0.34) {
  const shape = new THREE.Shape();
  const n = Math.max(2, Math.round(w / scallop));
  const r = w / n / 2;
  shape.moveTo(-w / 2, 0);
  shape.lineTo(w / 2, 0);
  shape.lineTo(w / 2, -depth + r);
  for (let i = n - 1; i >= 0; i--) shape.absarc(-w / 2 + (2 * i + 1) * r, -depth + r, r, 0, -Math.PI, true);
  shape.lineTo(-w / 2, 0);
  return new THREE.ShapeGeometry(shape, 8);
}

// Shrinks the font until text fits w; returns the font set.
function fitFont(ctx, text, w, weight, size, min) {
  ctx.font = `${weight} ${size}px ${FONT}`;
  while (ctx.measureText(text).width > w && size > min) ctx.font = `${weight} ${--size}px ${FONT}`;
}

// Sign in front of a repo's tent: the repo's name over its merged task count.
function drawRepoSign(ctx, { repo, n, stripe }, w, h) {
  ctx.fillStyle = PALETTE.pad;
  ctx.beginPath();
  ctx.roundRect(0, 0, w, h, 22);
  ctx.fill();
  ctx.fillStyle = stripe;
  ctx.beginPath();
  ctx.roundRect(0, 0, 22, h, [22, 0, 0, 22]);
  ctx.fill();
  ctx.textBaseline = 'alphabetic';
  ctx.textAlign = 'left';
  ctx.fillStyle = PALETTE.muted;
  fitFont(ctx, repo, w - 70, 700, 54, 30);
  ctx.fillText(repo, 50, 76);
  ctx.fillStyle = PALETTE.ink;
  ctx.font = `800 108px ${FONT}`;
  const num = n == null ? '–' : String(n);
  ctx.fillText(num, 48, 188);
  const nw = ctx.measureText(num).width;
  ctx.fillStyle = PALETTE.muted;
  ctx.font = `700 50px ${FONT}`;
  ctx.fillText('merged', 48 + nw + 18, 186);
}

// 553000000 → "553M", 9200000 → "9.2M".
export function formatTokens(x) {
  const [v, unit] = x >= 1e9 ? [x / 1e9, 'B'] : x >= 1e6 ? [x / 1e6, 'M'] : x >= 1e3 ? [x / 1e3, 'k'] : [x, ''];
  return `${v >= 100 || unit === '' ? Math.round(v) : v.toFixed(1).replace(/\.0$/, '')}${unit}`;
}

// Project signboard: name on a band in the block colour, then the numbers.
function drawBoard(ctx, { name, floor, stats }, w, h) {
  const dark = new THREE.Color(floor).offsetHSL(0, 0.05, -0.36).getStyle();
  ctx.fillStyle = PALETTE.pad;
  ctx.beginPath();
  ctx.roundRect(6, 6, w - 12, h - 12, 34);
  ctx.fill();
  ctx.lineWidth = 10;
  ctx.strokeStyle = dark;
  ctx.stroke();
  ctx.save();
  ctx.clip();
  ctx.fillStyle = floor;
  ctx.fillRect(0, 0, w, 124);
  ctx.restore();
  ctx.fillStyle = dark;
  ctx.textBaseline = 'middle';
  ctx.textAlign = 'left';
  fitFont(ctx, name, w - 230, 800, 64, 36);
  ctx.fillText(name, 40, 68);
  ctx.font = `600 30px ${FONT}`;
  ctx.textAlign = 'right';
  ctx.fillText(stats ? `last ${stats.days} days` : 'no board', w - 40, 70);
  const s = stats ?? { tokens: 0, finished: 0, passed: 0 };
  const rows = [['Tokens', formatTokens(s.tokens)], ['Tasks done', String(s.finished)], ['Reviews passed', String(s.passed)], ['Tokens / task', s.finished ? formatTokens(s.tokens / s.finished) : '—']];
  rows.forEach(([k, v], i) => {
    const y = 172 + i * 84;
    ctx.textAlign = 'left';
    ctx.fillStyle = PALETTE.muted;
    ctx.font = `600 40px ${FONT}`;
    ctx.fillText(k, 40, y);
    ctx.textAlign = 'right';
    ctx.fillStyle = i === 1 ? dark : PALETTE.ink;
    ctx.font = `800 52px ${FONT}`;
    ctx.fillText(v, w - 40, y);
    if (i < rows.length - 1) {
      ctx.fillStyle = 'rgba(36,33,29,0.08)';
      ctx.fillRect(40, y + 42, w - 80, 3);
    }
  });
}

// City Hall's tower screen: every project's Mate, crews at work and tokens today.
function drawHallScreen(ctx, rows, w, h) {
  ctx.fillStyle = '#2B2723';
  ctx.fillRect(0, 0, w, h);
  ctx.textBaseline = 'middle';
  ctx.textAlign = 'left';
  ctx.fillStyle = PALETTE.coral;
  ctx.font = `800 56px ${FONT}`;
  ctx.fillText('mate', 44, 62);
  ctx.fillStyle = 'rgba(255,253,248,0.55)';
  ctx.font = `600 30px ${FONT}`;
  ctx.textAlign = 'right';
  ctx.fillText('crews · tokens today', w - 44, 64);
  const list = rows.slice(0, 6);
  const step = Math.min(84, (h - 140) / Math.max(1, list.length));
  list.forEach((r, i) => {
    const y = 140 + i * step + step / 2 - 10;
    ctx.fillStyle = r.running ? '#9BD3A8' : 'rgba(255,253,248,0.3)';
    ctx.beginPath();
    ctx.arc(56, y, 11, 0, TAU);
    ctx.fill();
    ctx.textAlign = 'left';
    ctx.fillStyle = '#FFFDF8';
    fitFont(ctx, r.name, w * 0.46, 700, 42, 26);
    ctx.fillText(r.name, 84, y);
    ctx.textAlign = 'right';
    ctx.font = `800 42px ${FONT}`;
    ctx.fillText(formatTokens(r.tokens), w - 44, y);
    ctx.fillStyle = 'rgba(255,253,248,0.65)';
    ctx.font = `700 36px ${FONT}`;
    ctx.fillText(String(r.working), w * 0.68, y);
  });
  if (rows.length > list.length) {
    ctx.textAlign = 'left';
    ctx.fillStyle = 'rgba(255,253,248,0.45)';
    ctx.font = `600 28px ${FONT}`;
    ctx.fillText(`+${rows.length - list.length} more`, 84, h - 30);
  }
}

// n spots [x, z, heading] (world) inside tent b, at least GAP apart; the
// front-left corner behind the sign is kept clear when there is room. Each
// spot depends only on those before it, so the first k of any call agree.
const GAP = 1.8; // wanted spacing; the widest gap found wins when a tent is too small
function scatter(b, n, random) {
  const spots = [];
  const xr = Math.max(0, b.w / 2 - 0.55);
  const zMin = -b.d / 2 + 1.2; // deeper in, the roof hides them from most views
  const zMax = b.d / 2 - 0.6;
  for (let k = 0; k < n; k++) {
    let best = null;
    let bestGap = -1;
    for (let tries = 0; tries < 60; tries++) {
      const x = (2 * random() - 1) * xr;
      const z = zMin + random() * (zMax - zMin);
      if (b.w > 2.6 && x < -b.w / 2 + SIGN_W + 0.1 && z > b.d / 2 - 1.3) continue;
      const gap = Math.min(Infinity, ...spots.map(([sx, sz]) => Math.hypot(sx - x, sz - z)));
      if (gap > bestGap) [best, bestGap] = [[x, z], gap];
      if (gap >= GAP) break;
    }
    best ??= [0, 0];
    spots.push([best[0], best[1], 0.7 * (random() - 0.5)]);
  }
  return spots.map(([x, z, h]) => [b.x + x, b.z + z, h]);
}

// Block slots [col, row] around City Hall for a grid n wide, in the order
// projects take them: ring by ring, the blocks beside the hall's sides before
// those in front of and behind it, nearer rows first.
export function blockSlots(n) {
  const r = (n - 1) / 2;
  const out = [];
  for (let ring = 1; ring <= r; ring++) {
    const ringSlots = [];
    for (let row = ring; row >= -ring; row--) for (let col = -ring; col <= ring; col++) if (Math.max(Math.abs(col), Math.abs(row)) === ring) ringSlots.push([col, row]);
    ringSlots.sort((a, b) => (a[0] === 0) - (b[0] === 0) || b[1] - a[1] || a[0] - b[0]);
    out.push(...ringSlots);
  }
  return out;
}

export const gridSize = (projects) => (projects > 8 ? 5 : 3);
export const tentKey = (project, repo) => `${project}/${repo}`;

/**
 * Build the city for the workspace's projects.
 * @param {{ name: string, repos: string[], mate: { harness: string } }[]} projects
 *   in workspace order; a project with no repos gets one tent named after it
 */
export function createCity(projects) {
  const N = gridSize(projects.length);
  const SLOTS = blockSlots(N);
  const ISLAND = N * BLOCK + (N - 1) * ROAD + 2.8;
  const STREETS = Array.from({ length: N - 1 }, (_, k) => (k - (N - 2) / 2) * PITCH);
  const shown = projects.slice(0, SLOTS.length);

  const group = new THREE.Group();
  group.name = 'agent-city';
  const disposables = [];
  const keep = (x) => (disposables.push(x), x);
  const label = (text, height, options) => {
    const { texture, aspect } = textTexture(text, options);
    const mesh = new THREE.Mesh(keep(new THREE.PlaneGeometry(height * aspect, height)), keep(new THREE.MeshBasicMaterial({ map: keep(texture), transparent: true, toneMapped: false })));
    mesh.name = `label:${text}`;
    return mesh;
  };
  const centre = ([col, row]) => new THREE.Vector3(col * PITCH, 0, row * PITCH);

  // Island.
  const cliff = new THREE.Mesh(keep(new RoundedBoxGeometry(ISLAND, 2.6, ISLAND, 5, 1.2)), keep(matte(PALETTE.cliff)));
  cliff.position.y = -1.3 - 0.002;
  cliff.receiveShadow = true;
  group.add(cliff);
  const top = new THREE.Mesh(keep(new THREE.ShapeGeometry(roundedRect(ISLAND - 1.2, ISLAND - 1.2, 1.0), 12)), keep(matte(PALETTE.island)));
  top.rotation.x = -Math.PI / 2;
  top.receiveShadow = true;
  group.add(top);

  // Streets between the blocks, with centre dashes and zebra crossings at the junctions.
  const roadMat = keep(matte(PALETTE.road));
  const roadLen = ISLAND - 1.2;
  for (const s of STREETS) {
    for (const along of ['x', 'z']) {
      const road = new THREE.Mesh(keep(new THREE.PlaneGeometry(along === 'x' ? roadLen : ROAD, along === 'x' ? ROAD : roadLen)), roadMat);
      road.rotation.x = -Math.PI / 2;
      road.position.set(along === 'x' ? 0 : s, 0.004, along === 'x' ? s : 0);
      road.receiveShadow = true;
      group.add(road);
    }
  }
  const marks = new Batch(keep(new THREE.PlaneGeometry(1, 1)), keep(new THREE.MeshStandardMaterial({ color: PALETTE.line, roughness: 0.9 })), 'road-marks');
  const flat = [-Math.PI / 2, 0, 0];
  const nearJunction = (d) => STREETS.some((s) => Math.abs(d - s) < ROAD / 2 + 1.6);
  for (let d = -roadLen / 2 + 1; d < roadLen / 2 - 1; d += 1.25) {
    if (nearJunction(d)) continue;
    for (const s of STREETS) {
      marks.add([d, 0.008, s], { rotation: flat, scale: [0.6, 0.09, 1] });
      marks.add([s, 0.008, d], { rotation: flat, scale: [0.09, 0.6, 1] });
    }
  }
  for (const jx of STREETS) {
    for (const jz of STREETS) {
      for (let k = -2; k <= 2; k++) {
        const o = (k * ROAD) / 6;
        for (const s of [-1, 1]) {
          const d = ROAD / 2 + 0.9;
          marks.add([jx + s * d, 0.008, jz + o], { rotation: flat, scale: [1.2, 0.2, 1] });
          marks.add([jx + o, 0.008, jz + s * d], { rotation: flat, scale: [0.2, 1.2, 1] });
        }
      }
    }
  }
  group.add(marks.build({ cast: false }));

  const windows = new Batch(keep(new THREE.PlaneGeometry(1, 1)), keep(new THREE.MeshStandardMaterial({ color: PALETTE.glass, roughness: 0.35, metalness: 0.1 })), 'windows');
  const padGeo = keep(new THREE.CylinderGeometry(1, 1, 0.08, 48));
  const padMat = keep(matte(PALETTE.pad));
  const coralMat = keep(matte(PALETTE.coral));
  const mateMat = keep(matte('#F2A285'));
  const doorMat = keep(matte(PALETTE.door));
  const random = rng(71);
  const trunks = new Batch(keep(new THREE.CylinderGeometry(0.09, 0.12, 1, 8)), keep(matte(PALETTE.trunk)), 'tree-trunks');
  const crowns = new Batch(keep(new THREE.SphereGeometry(1, 20, 14)), keep(matte('#ffffff', { roughness: 0.95 })), 'tree-crowns');
  function tree(x, z, s = 1) {
    const h = (0.7 + 0.25 * random()) * s;
    trunks.add([x, KERB + h / 2, z], { scale: [s, h, s] });
    const r = (0.55 + 0.18 * random()) * s;
    crowns.add([x, KERB + h + r * 0.75, z], { scale: [r, r * 1.15, r], color: PALETTE.leaf[Math.floor(random() * PALETTE.leaf.length)] });
  }

  function tile(slot, color) {
    const c = centre(slot);
    const mesh = new THREE.Mesh(keep(new RoundedBoxGeometry(BLOCK, KERB * 2, BLOCK, 3, 0.35)), keep(matte(color)));
    mesh.position.set(c.x, 0, c.z);
    mesh.receiveShadow = true;
    group.add(mesh);
    return c;
  }

  function groundText(text, x, z, height, color) {
    const m = label(text, height, { font: `800 120px ${FONT}`, color, pad: 20, letterSpacing: 4 });
    m.material.toneMapped = true;
    m.rotation.x = -Math.PI / 2;
    m.position.set(x, KERB + 0.006, z);
    group.add(m);
    return m;
  }

  function pad(x, y, z, radius, material = padMat) {
    const p = new THREE.Mesh(padGeo, material);
    p.scale.set(radius, 1, radius);
    p.position.set(x, y + 0.04, z);
    p.receiveShadow = true;
    group.add(p);
    return p;
  }

  // Windows on the front of a plain box of width w, depth d, at height y.
  function windowRow(x, y, z, w, d) {
    const cols = Math.max(1, Math.floor((w - 0.5) / 0.85));
    for (let c = 0; c < cols; c++) windows.add([x + (c - (cols - 1) / 2) * 0.85, y, z + d / 2 + 0.006], { scale: [0.46, 0.56, 1] });
  }

  // ---- Repo tents -----------------------------------------------------------

  // Each repo is a market tent: a striped hip roof on four poles, a scalloped
  // valance, a back cloth, open front and sides so its crews show from the
  // street. A small sign in front gives the repo and its merged task count.
  const tents = new Map(); // tentKey → tent record
  const poleMat = keep(matte(PALETTE.trunk));
  const poleGeo = keep(new THREE.CylinderGeometry(0.05, 0.05, 1, 10));
  const roofMat = keep(new THREE.MeshStandardMaterial({ vertexColors: true, roughness: 0.92, side: THREE.DoubleSide }));
  const pennantGeo = keep(new THREE.ShapeGeometry(new THREE.Shape([new THREE.Vector2(0, 0.13), new THREE.Vector2(0.42, 0), new THREE.Vector2(0, -0.13)])));
  const pennantMat = keep(new THREE.MeshStandardMaterial({ color: PALETTE.coral, roughness: 0.9, side: THREE.DoubleSide }));
  const signBoardGeo = keep(new RoundedBoxGeometry(SIGN_W + 0.08, SIGN_H + 0.08, 0.06, 3, 0.025));
  const signBackMat = keep(matte('#E9E3D6'));
  const signFaceGeo = keep(new THREE.PlaneGeometry(SIGN_W, SIGN_H));

  // Hip roof over a W×D eave rectangle at y = 0, its ridge along x at y = rise,
  // in alternating stripes running eave to ridge.
  function tentRoof(W, D, rise, colors) {
    const r = Math.max(0, (W - D) / 2);
    const v = (x, y, z) => new THREE.Vector3(x, y, z);
    const FL = v(-W / 2, 0, D / 2), FR = v(W / 2, 0, D / 2), BR = v(W / 2, 0, -D / 2), BL = v(-W / 2, 0, -D / 2);
    const RL = v(-r, rise, 0), RR = v(r, rise, 0);
    const pos = [];
    const col = [];
    const tint = colors.map((c) => new THREE.Color(c));
    // Eave a→b, ridge p→q; a→b runs left to right seen from outside.
    const face = (a, b, p, q) => {
      const n = Math.max(2, Math.round(a.distanceTo(b) / 0.42));
      for (let i = 0; i < n; i++) {
        const a0 = a.clone().lerp(b, i / n), a1 = a.clone().lerp(b, (i + 1) / n);
        const p0 = p.clone().lerp(q, i / n), p1 = p.clone().lerp(q, (i + 1) / n);
        for (const pt of [a0, a1, p1, a0, p1, p0]) pos.push(pt.x, pt.y, pt.z);
        const c = tint[i % tint.length];
        for (let k = 0; k < 6; k++) col.push(c.r, c.g, c.b);
      }
    };
    face(FL, FR, RL, RR);
    face(BR, BL, RR, RL);
    face(FR, BR, RR, RR);
    face(BL, FL, RL, RL);
    const g = new THREE.BufferGeometry();
    g.setAttribute('position', new THREE.Float32BufferAttribute(pos, 3));
    g.setAttribute('color', new THREE.Float32BufferAttribute(col, 3));
    g.computeVertexNormals();
    return keep(g);
  }

  function repoTent({ project, repo, x, z, w, d, cloth, stripe, floor }) {
    const key = tentKey(project, repo);
    const g = new THREE.Group();
    g.name = `repo:${key}`;
    g.position.set(x, KERB, z);
    group.add(g);
    const ground = new THREE.Mesh(keep(new RoundedBoxGeometry(w, 0.06, d, 2, 0.03)), padMat);
    ground.position.y = 0.03;
    ground.receiveShadow = true;
    g.add(ground);
    const back = new THREE.Mesh(keep(new THREE.PlaneGeometry(w - 0.1, EAVE)), keep(new THREE.MeshStandardMaterial({ color: cloth, roughness: 0.95, side: THREE.DoubleSide })));
    back.position.set(0, EAVE / 2, -d / 2 + 0.06);
    back.receiveShadow = true;
    g.add(back);
    for (const sx of [-1, 1]) {
      for (const sz of [-1, 1]) {
        const pole = new THREE.Mesh(poleGeo, poleMat);
        pole.scale.y = EAVE;
        pole.position.set(sx * (w / 2 - 0.06), EAVE / 2, sz * (d / 2 - 0.06));
        pole.castShadow = true;
        g.add(pole);
      }
    }
    const W = w + 0.3;
    const D = d + 0.3;
    const top = new THREE.Group(); // roof and valance, bounced by a merge
    top.position.y = EAVE;
    g.add(top);
    const roof = new THREE.Mesh(tentRoof(W, D, RISE, [stripe, PALETTE.pad]), roofMat);
    roof.castShadow = roof.receiveShadow = true;
    top.add(roof);
    const valMat = keep(new THREE.MeshStandardMaterial({ color: stripe, roughness: 0.95, side: THREE.DoubleSide }));
    for (const [len, ry, px, pz] of [[W, 0, 0, D / 2], [W, Math.PI, 0, -D / 2], [D, Math.PI / 2, W / 2, 0], [D, -Math.PI / 2, -W / 2, 0]]) {
      const val = new THREE.Mesh(keep(valanceGeometry(len)), valMat);
      val.rotation.y = ry;
      val.position.set(px, 0.005, pz);
      val.castShadow = true;
      top.add(val);
    }
    // Pennant on the ridge.
    const mast = new THREE.Mesh(poleGeo, poleMat);
    mast.scale.y = 0.5;
    mast.position.y = RISE + 0.25;
    top.add(mast);
    const pennant = new THREE.Mesh(pennantGeo, pennantMat);
    pennant.position.set(0.02, RISE + 0.44, 0);
    top.add(pennant);
    // Sign on two stakes at the front-left corner, leaning back.
    const sign = new THREE.Group();
    sign.position.set(-w / 2 + SIGN_W / 2 + 0.05, 0, d / 2 + 0.45);
    g.add(sign);
    for (const s of [-1, 1]) {
      const stake = new THREE.Mesh(poleGeo, poleMat);
      stake.scale.set(0.8, 0.42, 0.8);
      stake.position.set(s * (SIGN_W / 2 - 0.12), 0.21, -0.02);
      sign.add(stake);
    }
    const panel = new THREE.Group();
    panel.position.y = 0.42 + SIGN_H / 2 - 0.04;
    panel.rotation.x = -0.18;
    sign.add(panel);
    const board = new THREE.Mesh(signBoardGeo, signBackMat);
    board.castShadow = true;
    panel.add(board);
    const signTex = new LiveTexture(520, 220, (ctx, n, tw, th) => drawRepoSign(ctx, { repo, n, stripe }, tw, th));
    keep(signTex.texture);
    signTex.set(null);
    const face = new THREE.Mesh(signFaceGeo, keep(new THREE.MeshBasicMaterial({ map: signTex.texture, transparent: true, toneMapped: false })));
    face.position.z = 0.036;
    panel.add(face);
    const record = { key, project, repo, floor, group: g, top, panel, signTex, x, z, w, d, merged: null, popAt: -Infinity, crews: new Map(), spots: null };
    record.spots = scatter(record, SPOTS, rng(hash(key)));
    g.userData.tent = key;
    tents.set(key, record);
    return record;
  }

  // A merge pops the sign and bounces the roof.
  function poseTent(b, t) {
    const k = clamp01((t - b.popAt) / POP);
    const pop = k < 1 ? Math.sin(Math.PI * k) : 0;
    b.panel.scale.setScalar(1 + 0.35 * pop);
    b.top.position.y = EAVE + 0.18 * pop;
    b.top.scale.set(1 + 0.04 * pop, 1 - 0.06 * pop, 1 + 0.04 * pop);
  }

  // ---- Projects ----------------------------------------------------------------

  const mascots = new Map(); // id → AgentMascot; crews by "<project>/<crew>", Mates by "<project>:mate"
  const seats = new Map(); // id → { project, kind, crew?, role, harness, model, effort, state, text }
  const mates = new Map(); // project → { mascot, state, running }
  const boards = new Map(); // project → { tex, panel, nodAt }
  const blocks = new Map(); // project → { floor, tents: [tentKey], centre }
  const crewPad = keep(matte(PALETTE.pad));
  const postMat = keep(matte('#8E857A'));
  const postGeo = keep(new THREE.CylinderGeometry(0.06, 0.07, 1, 10));
  let phase = 0.37;

  function mascot(id, look, x, y, z, size, heading = 0) {
    const m = new AgentMascot({ ...look, size, position: [x, y, z], heading, phase: (phase = (phase * 7.31) % 4), bubbleNear: 30 });
    m.id = id;
    group.add(m.root);
    mascots.set(id, m);
    return m;
  }

  // Signboard on two posts at the block's front-right corner, facing +z.
  function signboard(project, floor, c) {
    const g = new THREE.Group();
    g.name = `signboard:${project}`;
    g.position.set(c.x + 2.4, KERB, c.z + BLOCK / 2 - 0.85);
    group.add(g);
    for (const s of [-1, 1]) {
      const post = new THREE.Mesh(postGeo, postMat);
      post.scale.y = 1.1;
      post.position.set(s * 1.0, 0.55, -0.02);
      post.castShadow = true;
      g.add(post);
    }
    const panel = new THREE.Group();
    panel.position.y = 1.1 + 0.88;
    panel.rotation.x = -0.1;
    g.add(panel);
    const back = new THREE.Mesh(keep(new RoundedBoxGeometry(2.7, 1.86, 0.1, 3, 0.04)), keep(matte('#E9E3D6')));
    back.castShadow = true;
    panel.add(back);
    const tex = new LiveTexture(768, 520, (ctx, key, w, h) => drawBoard(ctx, { name: project, floor, stats: JSON.parse(key) }, w, h));
    keep(tex.texture);
    tex.set('null');
    const face = new THREE.Mesh(keep(new THREE.PlaneGeometry(2.6, 2.6 * (520 / 768))), keep(new THREE.MeshBasicMaterial({ map: tex.texture, transparent: true, toneMapped: false })));
    face.position.z = 0.052;
    panel.add(face);
    boards.set(project, { tex, panel, nodAt: -Infinity });
  }

  shown.forEach((p, i) => {
    const slot = SLOTS[i];
    const floor = PALETTE.floors[i % PALETTE.floors.length];
    const c = tile(slot, floor);
    const names = p.repos.length ? p.repos : [p.name];
    const n = names.length;
    const span = BLOCK - 1.6;
    const w = Math.min(6.2, span / n - 0.45);
    const d = 3.4;
    const stripe = new THREE.Color(floor).offsetHSL(0, -0.05, -0.12).getStyle();
    const keys = names.map((repo, k) => repoTent({ project: p.name, repo, x: c.x + (k - (n - 1) / 2) * (span / n), z: c.z - 1.9, w, d, cloth: PALETTE.walls[k % PALETTE.walls.length], stripe, floor }).key);
    blocks.set(p.name, { floor, tents: keys, centre: c });
    signboard(p.name, floor, c);
    groundText(p.name.toUpperCase(), c.x - 1.6, c.z + BLOCK / 2 - 0.75, 0.55, new THREE.Color(floor).offsetHSL(0, 0.05, -0.3).getStyle());
    pad(c.x, KERB, c.z + 2.9, 0.62 * MATE_SIZE + 0.08, mateMat);
    const m = mascot(`${p.name}:mate`, harnessLook(p.mate?.harness), c.x, KERB + 0.08, c.z + 2.9, MATE_SIZE);
    m.play('sleeping', { at: 0 });
    mates.set(p.name, { mascot: m, state: null, running: false });
    seats.set(m.id, { project: p.name, kind: 'mate', harness: p.mate?.harness ?? '', state: '' });
    // A tree in the front-left corner; the signboard takes the right one.
    tree(c.x - (BLOCK / 2 - 0.8), c.z + BLOCK / 2 - 0.8, 1.0);
  });

  // Parks: the blocks no project takes, lawns with a ring of trees.
  for (const slot of SLOTS.slice(shown.length)) {
    const c = tile(slot, PALETTE.park);
    for (let i = 0; i < 9; i++) {
      const a = (i / 9) * TAU + 0.3;
      tree(c.x + Math.cos(a) * 3.6, c.z + Math.sin(a) * 3.6, 0.9 + 0.35 * random());
    }
    tree(c.x - 1.1, c.z - 0.6, 1.3);
    tree(c.x + 1.2, c.z + 0.9, 1.1);
  }

  // ---- City Hall -----------------------------------------------------------------

  const hall = tile([0, 0], PALETTE.hall);
  const white = keep(matte('#FFFDF8'));
  const hallBody = new THREE.Mesh(keep(new RoundedBoxGeometry(8.2, 3.4, 4.2, 4, 0.2)), white);
  hallBody.position.set(hall.x, KERB + 1.7, hall.z - 1.6);
  hallBody.castShadow = hallBody.receiveShadow = true;
  group.add(hallBody);
  for (let i = 0; i < 2; i++) {
    windowRow(hall.x - 2.9, KERB + 1.1 + i * 1.3, hall.z - 1.6, 2.2, 4.2);
    windowRow(hall.x + 2.9, KERB + 1.1 + i * 1.3, hall.z - 1.6, 2.2, 4.2);
  }
  const portico = new THREE.Mesh(keep(new RoundedBoxGeometry(8.6, 0.34, 1.7, 3, 0.12)), coralMat);
  portico.position.set(hall.x, KERB + 3.3, hall.z + 1.1);
  portico.castShadow = true;
  group.add(portico);
  const columnGeo = keep(new THREE.CylinderGeometry(0.17, 0.2, 3.1, 20));
  for (let i = 0; i < 6; i++) {
    const col = new THREE.Mesh(columnGeo, white);
    col.position.set(hall.x - 3.6 + i * 1.44, KERB + 1.67, hall.z + 1.55);
    col.castShadow = true;
    group.add(col);
  }
  for (let i = 0; i < 3; i++) {
    const step = new THREE.Mesh(keep(new RoundedBoxGeometry(8.8 - i * 0.5, 0.12, 2.8 - i * 0.45, 2, 0.05)), white);
    step.position.set(hall.x, KERB + 0.06 + i * 0.12, hall.z + 0.9 - i * 0.12);
    step.receiveShadow = true;
    group.add(step);
  }
  const hallDoor = new THREE.Mesh(keep(new THREE.ShapeGeometry(roundedRect(1.5, 2.1, 0.3), 6)), doorMat);
  hallDoor.position.set(hall.x, KERB + 1.4, hall.z + 0.51);
  group.add(hallDoor);
  const hallSign = label('CAPTAIN', 0.26, { font: `900 96px ${FONT}`, color: '#FFFDF8', pad: 14, letterSpacing: 14 });
  hallSign.position.set(hall.x, KERB + 3.3, hall.z + 1.96);
  group.add(hallSign);
  // Tower with the workspace screen, a coral dome and a flag.
  const tower = new THREE.Mesh(keep(new RoundedBoxGeometry(3.8, 3.6, 3.0, 4, 0.2)), white);
  tower.position.set(hall.x, KERB + 5.2, hall.z - 1.8);
  tower.castShadow = true;
  group.add(tower);
  const screenTex = new LiveTexture(1024, 640, (ctx, key, w, h) => drawHallScreen(ctx, JSON.parse(key), w, h));
  keep(screenTex.texture);
  screenTex.set('[]');
  const screenW = 3.3;
  const screenH = screenW * (640 / 1024);
  const frame = new THREE.Mesh(keep(new RoundedBoxGeometry(screenW + 0.24, screenH + 0.24, 0.12, 3, 0.08)), keep(matte('#2B2723')));
  frame.position.set(hall.x, KERB + 5.3, hall.z - 1.8 + 1.52);
  group.add(frame);
  const screen = new THREE.Mesh(keep(new THREE.PlaneGeometry(screenW, screenH)), keep(new THREE.MeshBasicMaterial({ map: screenTex.texture, toneMapped: false })));
  screen.position.set(hall.x, KERB + 5.3, hall.z - 1.8 + 1.59);
  group.add(screen);
  const dome = new THREE.Mesh(keep(new THREE.SphereGeometry(1.45, 40, 20, 0, TAU, 0, Math.PI / 2)), coralMat);
  dome.position.set(hall.x, KERB + 7.0, hall.z - 1.8);
  dome.castShadow = true;
  group.add(dome);
  const pole = new THREE.Mesh(keep(new THREE.CylinderGeometry(0.04, 0.04, 1.6, 8)), keep(matte('#8E857A')));
  pole.position.set(hall.x, KERB + 9.2, hall.z - 1.8);
  group.add(pole);
  const flagGeo = keep(new THREE.PlaneGeometry(0.9, 0.55, 12, 1));
  flagGeo.translate(0.45, 0, 0);
  const flag = new THREE.Mesh(flagGeo, keep(matte(PALETTE.coral, { side: THREE.DoubleSide })));
  flag.position.set(hall.x + 0.04, KERB + 9.7, hall.z - 1.8);
  flag.castShadow = true;
  group.add(flag);
  const flagRest = flagGeo.attributes.position.array.slice();
  groundText('CITY HALL', hall.x, hall.z + BLOCK / 2 - 0.9, 0.7, '#C9826A');
  tree(hall.x - BLOCK / 2 + 0.9, hall.z + BLOCK / 2 - 0.9, 1.1);
  tree(hall.x + BLOCK / 2 - 0.9, hall.z + BLOCK / 2 - 0.9, 1.1);

  // Mailbox for the news the Mates opened, right of the steps, with the slot
  // in its lid. Its flag is up while any news waits on the captain.
  const mailbox = new THREE.Group();
  mailbox.name = 'mailbox';
  mailbox.position.set(hall.x + 2.6, KERB, hall.z + 3.4);
  group.add(mailbox);
  const post = new THREE.Mesh(keep(new THREE.CylinderGeometry(0.07, 0.08, 0.9, 10)), postMat);
  post.position.y = 0.45;
  mailbox.add(post);
  const box = new THREE.Mesh(keep(new RoundedBoxGeometry(0.9, 0.62, 0.62, 4, 0.16)), coralMat);
  box.position.y = 1.2;
  box.castShadow = true;
  mailbox.add(box);
  const slot = new THREE.Mesh(keep(new THREE.PlaneGeometry(0.56, 0.08)), keep(matte(PALETTE.ink)));
  slot.position.set(0, 1.512, 0);
  slot.rotation.x = -Math.PI / 2;
  mailbox.add(slot);
  const mailLabel = label('NEWS', 0.16, { font: `900 80px ${FONT}`, color: '#FFFDF8', pad: 8, letterSpacing: 8 });
  mailLabel.position.set(0, 1.2, 0.315);
  mailbox.add(mailLabel);
  const flagArm = new THREE.Group();
  flagArm.position.set(0.47, 1.12, 0.1);
  mailbox.add(flagArm);
  const mailFlag = new THREE.Mesh(keep(new RoundedBoxGeometry(0.05, 0.5, 0.16, 2, 0.02)), keep(matte('#F7D46B')));
  mailFlag.position.y = 0.25;
  flagArm.add(mailFlag);
  const countTex = new LiveTexture(128, 128, (ctx, n, w, h) => {
    ctx.fillStyle = PALETTE.ink;
    ctx.beginPath();
    ctx.arc(w / 2, h / 2, w / 2 - 4, 0, TAU);
    ctx.fill();
    ctx.fillStyle = '#FFFDF8';
    ctx.font = `800 ${n > 99 ? 52 : 72}px ${FONT}`;
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(String(n), w / 2, h / 2 + 4);
  });
  keep(countTex.texture);
  const countChip = new THREE.Sprite(keep(new THREE.SpriteMaterial({ map: countTex.texture, depthWrite: false })));
  countChip.position.set(0.05, 1.95, 0.1);
  countChip.scale.setScalar(0.42);
  countChip.visible = false;
  mailbox.add(countChip);
  const mail = { n: 0, at: -Infinity, from: 0 };

  group.add(windows.build({ cast: false }));

  // Lamps along the streets, none on City Hall's kerb, to keep its front clear.
  const poles = new Batch(keep(new THREE.CylinderGeometry(0.05, 0.06, 1, 8)), keep(matte('#8E857A')), 'lamp-poles');
  const bulbs = new Batch(keep(new THREE.SphereGeometry(0.17, 16, 10)), keep(new THREE.MeshStandardMaterial({ color: '#FFF4DC', emissive: '#FFE3A8', emissiveIntensity: 0.6, roughness: 0.5 })), 'lamp-bulbs');
  const onHall = (x, z) => Math.abs(x) < BLOCK / 2 && Math.abs(z) < BLOCK / 2;
  for (let d = -roadLen / 2 + 2.5; d < roadLen / 2 - 2; d += 5) {
    if (nearJunction(d)) continue;
    for (const s of STREETS) {
      for (const side of [-1, 1]) {
        const off = s + side * (ROAD / 2 + 0.25);
        if (!onHall(d, off)) {
          poles.add([d, 1.0, off], { scale: [1, 2.0, 1] });
          bulbs.add([d, 2.05, off]);
        }
        if (!onHall(off, d)) {
          poles.add([off, 1.0, d], { scale: [1, 2.0, 1] });
          bulbs.add([off, 2.05, d]);
        }
      }
    }
  }
  group.add(trunks.build(), crowns.build(), poles.build(), bulbs.build({ cast: false }));

  // Clouds drift on a slow loop.
  const cloudMat = keep(new THREE.MeshStandardMaterial({ color: '#FFFFFF', roughness: 1, transparent: true, opacity: 0.92 }));
  const puff = keep(new THREE.SphereGeometry(1, 20, 14));
  const spread = ISLAND / 39.8;
  const clouds = [[-18, 14, -16, 1.3], [14, 15, 8, 1.1], [4, 13.5, -20, 1.0], [-17, 12.5, 14, 0.9], [20, 13, -12, 1.0]].map(([x, y, z, s], i) => {
    const c = new THREE.Group();
    for (const [dx, dy, dz, r] of [[0, 0, 0, 1], [1.1, -0.15, 0.1, 0.75], [-1.05, -0.2, -0.05, 0.7], [0.4, 0.45, -0.1, 0.65]]) {
      const p = new THREE.Mesh(puff, cloudMat);
      p.position.set(dx, dy, dz);
      p.scale.set(r, r * 0.8, r * 0.9);
      p.castShadow = true;
      c.add(p);
    }
    c.scale.setScalar(s);
    c.userData = { x: x * spread, y, z: z * spread, i };
    group.add(c);
    return c;
  });

  // ---- Live state -------------------------------------------------------------

  // The tent a crew works in: its repo's, else the project's first.
  function tentFor(project, repo) {
    const b = blocks.get(project);
    if (!b) return null;
    return tents.get(tentKey(project, repo)) ?? tents.get(b.tents[0]);
  }

  function addCrew(tent, c, t, fresh) {
    const used = new Set([...tent.crews.values()].map((x) => x.spot));
    let spot = 0;
    while (used.has(spot)) spot++;
    // Past the tent's spots, crews queue in the street in front of it.
    const [x, z, heading] = spot < SPOTS ? tent.spots[spot] : [tent.x - tent.w / 2 + 0.6 + ((spot - SPOTS) % 6) * 0.9, tent.z + tent.d / 2 + 1.6, 0];
    const id = `${tent.project}/${c.crew}`;
    const m = mascot(id, harnessLook(c.harness), x, KERB + 0.14, z, CREW_SIZE, heading);
    const p = pad(x, KERB + 0.06, z, 0.62 * CREW_SIZE + 0.04, crewPad);
    m.play(poseOf(c.state), { at: t });
    if (fresh) m.say(c.crew, { at: t, hold: 2.2 });
    tent.crews.set(c.crew, { spot, mascot: m, pad: p, state: c.state });
  }

  function removeCrew(tent, crew) {
    const x = tent.crews.get(crew);
    group.remove(x.mascot.root, x.pad);
    x.mascot.dispose();
    mascots.delete(x.mascot.id);
    seats.delete(x.mascot.id);
    tent.crews.delete(crew);
  }

  /**
   * The project's crews now on the scene: [{ crew, repo, harness, state,
   * role, model, effort, text }]. Crews that left are taken off, new ones
   * placed in their repo's tent (saying their name unless this is the first
   * load), and a changed state plays its pose at t.
   */
  function setCrews(project, list, t, { fresh = true } = {}) {
    const b = blocks.get(project);
    if (!b) return;
    const want = new Map(); // tentKey → Map(crew → row)
    for (const c of list) {
      const tent = tentFor(project, c.repo);
      if (!want.has(tent.key)) want.set(tent.key, new Map());
      want.get(tent.key).set(c.crew, c);
    }
    for (const key of b.tents) {
      const tent = tents.get(key);
      const rows = want.get(key) ?? new Map();
      for (const crew of [...tent.crews.keys()]) if (!rows.has(crew)) removeCrew(tent, crew);
      for (const [crew, c] of rows) {
        const have = tent.crews.get(crew);
        if (!have) addCrew(tent, c, t, fresh);
        else if (have.state !== c.state) {
          have.mascot.play(poseOf(c.state), { at: t });
          if (SAYS[c.state]) have.mascot.say(SAYS[c.state], { at: t, hold: 1.6 });
          have.state = c.state;
        }
        seats.set(`${project}/${crew}`, { project, kind: 'crew', tent: key, ...c });
      }
    }
  }

  /** The project's Mate: { harness, running, state }. A Mate not running sleeps. */
  function setMate(project, card, t) {
    const m = mates.get(project);
    if (!m) return;
    const state = card.running ? card.state || 'idle' : 'asleep';
    if (state !== m.state) {
      m.mascot.play(poseOf(state), { at: t });
      m.state = state;
    }
    seats.set(m.mascot.id, { project, kind: 'mate', harness: card.harness, state, running: card.running });
  }

  /** The repo's merged task count; a rise after the first one pops its sign at t. */
  function setMerged(project, repo, n, t) {
    const tent = tents.get(tentKey(project, repo));
    if (!tent) return;
    if (tent.merged != null && n > tent.merged) tent.popAt = t;
    tent.merged = n;
    tent.signTex.set(n);
  }

  /** The project's board numbers: { days, tokens, finished, passed }, or null. */
  function setBoard(project, stats, t) {
    const b = boards.get(project);
    if (!b) return;
    const key = JSON.stringify(stats ?? null);
    const before = b.tex.key;
    b.tex.set(key);
    if (before !== 'null' && before !== key) b.nodAt = t;
  }

  /** City Hall's screen: [{ name, running, working, tokens }]. */
  function setHall(rows) {
    screenTex.set(JSON.stringify(rows));
  }

  /** News waiting on the captain: the flag goes up and the count pops when it rises. */
  function setMailbox(n, t) {
    if (n === mail.n) return;
    mail.from = mail.n;
    if (n > mail.n) mail.at = t;
    mail.n = n;
    if (n) countTex.set(n);
  }

  function updateMailbox(t) {
    const since = t - mail.at;
    const raise = mail.n ? (mail.from ? 1 : ease(clamp01(since / 0.4))) : 0;
    flagArm.rotation.z = -Math.PI / 2 + (Math.PI / 2) * raise;
    countChip.visible = mail.n > 0;
    countChip.scale.setScalar(0.42 * (since < 0.4 ? Math.max(0.001, backOut(since / 0.4)) : 1));
  }

  // Now and then a mascot that is not talking says a line from chatter.js
  // that fits its pose; each keeps its own random rhythm.
  function chat(m, t) {
    m.chatAt ??= t + 3 + Math.random() * CHAT_GAP[1];
    if (t < m.chatAt) return;
    m.chatAt = t + CHAT_GAP[0] + Math.random() * CHAT_GAP[1];
    const line = m.speechAt(t) ? null : lineFor(m.stateAt(t));
    if (line) m.say(line, { at: t, hold: 1.4, style: CHAT_STYLE });
  }

  function update(t) {
    for (const b of tents.values()) poseTent(b, t);
    for (const b of boards.values()) {
      const k = clamp01((t - b.nodAt) / 0.6);
      b.panel.scale.setScalar(1 + 0.08 * Math.sin(Math.PI * k) * (k < 1 ? 1 : 0));
    }
    for (const m of mascots.values()) {
      chat(m, t);
      m.update(t);
    }
    updateMailbox(t);
    for (const c of clouds) {
      const { x, y, z, i } = c.userData;
      c.position.set(x + 2 * Math.sin((TAU * t) / 120 + i * 1.7), y + 0.3 * Math.sin((TAU * t) / 60 + i), z);
    }
    const pos = flagGeo.attributes.position;
    for (let i = 0; i < pos.count; i++) {
      const x = flagRest[i * 3];
      pos.setZ(i, 0.08 * x * Math.sin(TAU * (t / 1.3 - x * 1.1)));
    }
    pos.needsUpdate = true;
    group.updateMatrixWorld(true);
  }

  return {
    group,
    size: N,
    extent: ISLAND / 2 + 2,
    hidden: projects.length - shown.length,
    mascots,
    seats,
    tents,
    blocks,
    hall: [group.getObjectByName('mailbox'), hallBody, tower, portico, screen, dome],
    setCrews,
    setMate,
    setMerged,
    setBoard,
    setHall,
    setMailbox,
    update,
    dispose() {
      for (const m of mascots.values()) m.dispose();
      for (const d of disposables) d.dispose();
    },
  };
}
