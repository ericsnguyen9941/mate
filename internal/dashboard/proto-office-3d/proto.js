// PROTOTYPE (throwaway). Three looks for a 3D Mate Office on one page,
// switchable with ?variant=A|B|C. See index.html.
import * as THREE from 'three';
import { OutlineEffect } from './vendor/effects/OutlineEffect.js';
import { GLTFLoader } from './vendor/loaders/GLTFLoader.js';
import { RoundedBoxGeometry } from './vendor/geometries/RoundedBoxGeometry.js';

// ------------------------------------------------------------ data (static)
// What the real page reads from /api: rooms, Mates, crews and their states.
const OFFICE = {
  rooms: [
    { id: 'plantrip', tokens: '447.8M', mate: { harness: 'claude', state: 'working' }, cabinet: 84,
      crews: [
        { name: 'pla-102-past-trips', harness: 'claude', state: 'working' },
        { name: 'pla-104-map-pins', harness: 'codex', state: 'stuck' },
      ] },
    { id: 'mihuv2', tokens: '295.3M', mate: { harness: 'claude', state: 'idle' }, cabinet: 115,
      crews: [
        { name: 'hu-35-washi-tilt', harness: 'claude', state: 'working' },
        { name: 'hu-35-r1', harness: 'pi', state: 'waiting_review' },
        { name: 'hu-88-staging', harness: 'codex', state: 'walking' },
      ] },
    { id: 'signal', tokens: '13M', mate: { harness: 'claude', state: 'asleep' }, cabinet: 7, crews: [] },
  ],
};

const STATE = {
  working:        { word: 'working',        color: '#2FB37A', bubble: null },
  waiting_review: { word: 'waiting review', color: '#E8A93A', bubble: '📄' },
  walking:        { word: 'to the Mate',    color: '#3B82F6', bubble: '?' },
  stuck:          { word: 'stuck',          color: '#E5484D', bubble: '!' },
  idle:           { word: 'idle',           color: '#8B95A1', bubble: null },
  asleep:         { word: 'asleep',         color: '#8E7CC3', bubble: 'z' },
};

const VARIANTS = {
  A: { name: 'Diorama toon · isometric, viền đen', build: buildA },
  B: { name: 'Voxel playground · Crossy Road', build: buildB },
  C: { name: 'Clay studio · render kiểu Blender', build: buildC },
};

// ------------------------------------------------------------ shared bits
const params = new URLSearchParams(location.search);
let key = (params.get('variant') || 'A').toUpperCase();
if (!VARIANTS[key]) key = 'A';

const renderer = new THREE.WebGLRenderer({ antialias: true });
renderer.setPixelRatio(Math.min(devicePixelRatio, 2));
renderer.setSize(innerWidth, innerHeight);
document.body.appendChild(renderer.domElement);

const loader = new GLTFLoader();
const MASCOT = {};
async function loadMascots() {
  for (const id of ['claude', 'codex', 'pi']) {
    try {
      const g = await loader.loadAsync(`./models/${id}.glb`);
      const obj = g.scene;
      const box = new THREE.Box3().setFromObject(obj);
      const size = box.getSize(new THREE.Vector3());
      const s = 1 / Math.max(size.x, size.y, size.z);
      obj.scale.setScalar(s);
      box.setFromObject(obj);
      const c = box.getCenter(new THREE.Vector3());
      obj.position.set(-c.x, -box.min.y, -c.z);
      const wrap = new THREE.Group();
      wrap.add(obj);
      MASCOT[id] = wrap;
    } catch (e) {
      console.warn('mascot', id, e);
    }
  }
}

// A mascot head of `size` units tall; falls back to a coloured squircle.
function head(harness, size, matFn) {
  const tpl = MASCOT[harness];
  let h;
  if (tpl) {
    h = tpl.clone(true);
    h.traverse((o) => { if (o.isMesh) { o.castShadow = true; const m = o.material.clone(); if ('metalness' in m) { m.metalness = 0; m.roughness = 0.45; } o.material = matFn ? matFn(m) : m; } });
  } else {
    const col = { claude: '#D97757', codex: '#1F1F1F', pi: '#2E7A6A' }[harness] || '#999';
    h = new THREE.Mesh(new RoundedBoxGeometry(1, 1, 1, 4, 0.25), matFn ? matFn(new THREE.MeshStandardMaterial({ color: col })) : new THREE.MeshStandardMaterial({ color: col }));
    h.position.y = 0.5;
    h.castShadow = true;
    const g = new THREE.Group(); g.add(h); h = g;
  }
  h.scale.setScalar(size);
  return h;
}

// Text as a sprite drawn on a canvas: a pill with the name, optional dot.
function label(text, opts = {}) {
  const { bg = 'rgba(255,255,255,.92)', fg = '#24211d', dot = null, scale = 1, bold = false } = opts;
  const c = document.createElement('canvas');
  const ctx = c.getContext('2d');
  const font = `${bold ? 700 : 600} 44px ui-sans-serif, system-ui, sans-serif`;
  ctx.font = font;
  const w = Math.ceil(ctx.measureText(text).width) + 70 + (dot ? 40 : 0), h = 88;
  c.width = w; c.height = h;
  ctx.font = font;
  ctx.fillStyle = bg;
  ctx.beginPath(); ctx.roundRect(0, 0, w, h, 44); ctx.fill();
  let x = 34;
  if (dot) { ctx.fillStyle = dot; ctx.beginPath(); ctx.arc(x + 10, h / 2, 13, 0, Math.PI * 2); ctx.fill(); x += 40; }
  ctx.fillStyle = fg; ctx.textBaseline = 'middle'; ctx.fillText(text, x, h / 2 + 2);
  const tex = new THREE.CanvasTexture(c); tex.colorSpace = THREE.SRGBColorSpace;
  const sp = new THREE.Sprite(new THREE.SpriteMaterial({ map: tex, depthTest: false, transparent: true }));
  sp.scale.set((w / h) * 0.6 * scale, 0.6 * scale, 1);
  sp.renderOrder = 10;
  return sp;
}

// A speech bubble above a head: "?", "!", "z" or an emoji.
function bubble(glyph, color) {
  const c = document.createElement('canvas'); c.width = c.height = 128;
  const ctx = c.getContext('2d');
  ctx.fillStyle = '#fff'; ctx.beginPath(); ctx.arc(64, 56, 46, 0, Math.PI * 2); ctx.fill();
  ctx.beginPath(); ctx.moveTo(50, 95); ctx.lineTo(64, 122); ctx.lineTo(78, 95); ctx.fill();
  ctx.lineWidth = 6; ctx.strokeStyle = color; ctx.beginPath(); ctx.arc(64, 56, 46, 0, Math.PI * 2); ctx.stroke();
  ctx.fillStyle = color; ctx.font = '700 58px ui-sans-serif, system-ui, sans-serif'; ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
  ctx.fillText(glyph, 64, 60);
  const tex = new THREE.CanvasTexture(c); tex.colorSpace = THREE.SRGBColorSpace;
  const sp = new THREE.Sprite(new THREE.SpriteMaterial({ map: tex, depthTest: false, transparent: true }));
  sp.scale.setScalar(1.1); sp.renderOrder = 11;
  return sp;
}

// Animation registry: each entry is fn(t).
const anims = [];
function animate(fn) { anims.push(fn); }

// A person: body + mascot head, posed by state. mk(geomW,geomH,geomD,color) builds a block.
function person({ harness, state, shirt, mk, headSize = 0.8, bodyH = 0.7 }) {
  const g = new THREE.Group();
  const body = mk(0.5, bodyH, 0.32, shirt);
  body.position.y = bodyH / 2 + 0.25;
  const legs = mk(0.42, 0.25, 0.28, '#3b3631'); legs.position.y = 0.125;
  const h = head(harness, headSize); h.position.y = bodyH + 0.28;
  g.add(legs, body, h);
  const st = STATE[state] || STATE.idle;
  if (st.bubble) { const b = bubble(st.bubble, st.color); b.position.y = bodyH + headSize + 0.9; g.add(b); animate((t) => { b.position.y = bodyH + headSize + 0.9 + Math.sin(t * 2.2) * 0.06; }); }
  const parts = { body, head: h, legs };
  switch (state) {
    case 'working': animate((t) => { parts.body.rotation.x = 0.18 + Math.sin(t * 6) * 0.03; parts.head.rotation.x = 0.25; }); break;
    case 'stuck': animate((t) => { g.rotation.z = Math.sin(t * 14) * 0.03; parts.head.position.y = bodyH + 0.28 + Math.abs(Math.sin(t * 7)) * 0.05; }); break;
    case 'asleep': parts.body.rotation.x = 0.45; parts.head.rotation.x = 0.7; parts.head.position.y = bodyH + 0.1; break;
    case 'waiting_review': animate((t) => { parts.head.rotation.y = Math.sin(t * 1.3) * 0.5; }); break;
    case 'idle': animate((t) => { g.position.y = Math.sin(t * 1.5) * 0.02; }); break;
    case 'walking': animate((t) => { parts.legs.rotation.x = Math.sin(t * 9) * 0.5; g.position.y = Math.abs(Math.sin(t * 9)) * 0.08; }); break;
  }
  g.userData.parts = parts;
  return g;
}

const SHIRTS = ['#8EA6BF', '#C4A07A', '#B48DA0', '#9BB096', '#C7A46A', '#89A9A3', '#B9917F', '#A9B57E', '#9FA6C7'];
function shirt(name) { let h = 0; for (const ch of name) h = (h * 31 + ch.charCodeAt(0)) >>> 0; return SHIRTS[h % SHIRTS.length]; }

// ------------------------------------------------------------ variant A
// Diorama toon: isometric orthographic camera, MeshToonMaterial with a
// three-step gradient, black outlines, saturated candy palette, rooms as
// floating slabs with two cut-away walls.
function buildA() {
  const scene = new THREE.Scene();
  scene.background = new THREE.Color('#FFF1D6');
  const grad = new THREE.DataTexture(new Uint8Array([90, 170, 255]), 3, 1, THREE.RedFormat);
  grad.minFilter = grad.magFilter = THREE.NearestFilter; grad.needsUpdate = true;
  const toon = (c) => new THREE.MeshToonMaterial({ color: c, gradientMap: grad });
  const mk = (w, h, d, c) => { const m = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), toon(c)); return m; };
  const toonify = (m) => toon(m.color ? m.color.clone() : '#999');

  scene.add(new THREE.HemisphereLight('#ffffff', '#ffb08a', 1.2));
  const sun = new THREE.DirectionalLight('#fff', 2.2); sun.position.set(6, 12, 4); scene.add(sun);

  const PAL = [{ floor: '#7FD8BE', wall: '#C9F2E4', trim: '#2FB37A' }, { floor: '#8EC5FF', wall: '#D6E9FF', trim: '#3B82F6' }, { floor: '#FFC29B', wall: '#FFE4D1', trim: '#F2742B' }];
  const world = new THREE.Group(); scene.add(world);
  const RW = 11, RD = 9, GAP = 4;
  const total = OFFICE.rooms.length * RW + (OFFICE.rooms.length - 1) * GAP;

  OFFICE.rooms.forEach((room, i) => {
    const p = PAL[i % PAL.length];
    const g = new THREE.Group();
    g.position.x = -total / 2 + RW / 2 + i * (RW + GAP);
    const floor = mk(RW, 0.6, RD, p.floor); floor.position.y = -0.3; g.add(floor);
    const base = mk(RW + 0.4, 0.3, RD + 0.4, '#5B4636'); base.position.y = -0.75; g.add(base);
    const wallB = mk(RW, 4.2, 0.35, p.wall); wallB.position.set(0, 2.1, -RD / 2 + 0.17); g.add(wallB);
    const wallL = mk(0.35, 4.2, RD, p.wall); wallL.position.set(-RW / 2 + 0.17, 2.1, 0); g.add(wallL);
    const band = mk(RW, 0.5, 0.4, p.trim); band.position.set(0, 3.9, -RD / 2 + 0.17); g.add(band);
    // Window on the back wall.
    const win = mk(3.2, 2, 0.1, '#BFEFFF'); win.position.set(2.6, 2.4, -RD / 2 + 0.4); g.add(win);
    const frame = mk(3.5, 2.3, 0.08, '#fff'); frame.position.set(2.6, 2.4, -RD / 2 + 0.36); g.add(frame);
    // Plant.
    const pot = mk(0.6, 0.6, 0.6, '#E07A5F'); pot.position.set(RW / 2 - 1, 0.3, -RD / 2 + 1); g.add(pot);
    const leaf = new THREE.Mesh(new THREE.SphereGeometry(0.75, 10, 8), toon('#3FA34D')); leaf.position.set(RW / 2 - 1, 1.3, -RD / 2 + 1); g.add(leaf);
    // Mate's big desk by the back wall.
    const mdesk = mk(4.2, 0.25, 1.6, '#F2A65A'); mdesk.position.set(-1.5, 1.05, -RD / 2 + 2.2); g.add(mdesk);
    for (const dx of [-1.9, 1.9]) { const leg = mk(0.2, 1, 1.4, '#C97B3B'); leg.position.set(-1.5 + dx, 0.5, -RD / 2 + 2.2); g.add(leg); }
    for (const dx of [-1.3, 1.3]) { const mon = mk(1.1, 0.75, 0.1, '#2B2B2B'); mon.position.set(-1.5 + dx, 1.6, -RD / 2 + 1.9); g.add(mon); const scr = mk(0.95, 0.6, 0.04, room.mate.state === 'asleep' ? '#444' : '#9BE7FF'); scr.position.set(-1.5 + dx, 1.6, -RD / 2 + 1.84); g.add(scr); }
    const mate = person({ harness: room.mate.harness, state: room.mate.state, shirt: p.trim, mk, headSize: 1.15, bodyH: 0.9 });
    mate.position.set(-1.5, 0, -RD / 2 + 3.1); mate.rotation.y = Math.PI; g.add(mate);
    const ml = label(`Mate · ${room.id}`, { dot: STATE[room.mate.state].color, bold: true }); ml.position.set(-1.5, 3.1, -RD / 2 + 3.1); g.add(ml);
    // Crew desks.
    room.crews.forEach((c, j) => {
      const x = -RW / 2 + 2.2 + j * 3.1, z = 1.6;
      const desk = mk(2.2, 0.2, 1.1, '#F7C98B'); desk.position.set(x, 0.95, z); g.add(desk);
      for (const dx of [-0.95, 0.95]) { const leg = mk(0.15, 0.9, 0.9, '#D3984E'); leg.position.set(x + dx, 0.45, z); g.add(leg); }
      const lap = mk(0.8, 0.5, 0.06, '#2B2B2B'); lap.position.set(x, 1.3, z - 0.25); lap.rotation.x = -0.2; g.add(lap);
      const glow = mk(0.68, 0.4, 0.03, c.state === 'working' ? '#9BE7FF' : '#555'); glow.position.set(x, 1.3, z - 0.21); glow.rotation.x = -0.2; g.add(glow);
      const who = person({ harness: c.harness, state: c.state, shirt: shirt(c.name), mk });
      who.position.set(x, 0, z + 0.9); who.rotation.y = Math.PI; g.add(who);
      if (c.state === 'walking') {
        animate((t) => { const k = (Math.sin(t * 0.6) + 1) / 2; who.position.set(x + (-1.5 - x) * k, who.position.y, z + 0.9 + (-RD / 2 + 4.2 - (z + 0.9)) * k); who.rotation.y = Math.PI + Math.atan2(-1.5 - x, -RD / 2 + 4.2 - z - 0.9) * 0.4; });
      }
      const nl = label(c.name, { dot: STATE[c.state].color, scale: 0.85 }); nl.position.set(x, 2.5, z + 0.9); g.add(nl);
    });
    // Filing cabinet.
    const cab = mk(1.2, 2.4, 1, '#4DA8A0'); cab.position.set(RW / 2 - 1.1, 1.2, 1.6); g.add(cab);
    for (let k = 0; k < 3; k++) { const dr = mk(1.05, 0.55, 0.08, '#6FCBC3'); dr.position.set(RW / 2 - 1.1, 0.5 + k * 0.7, 2.12); g.add(dr); const hd = mk(0.3, 0.08, 0.06, '#1F4B48'); hd.position.set(RW / 2 - 1.1, 0.5 + k * 0.7, 2.18); g.add(hd); }
    const cl = label(`${room.cabinet} done`, { bg: '#1F4B48', fg: '#fff', scale: 0.8 }); cl.position.set(RW / 2 - 1.1, 3, 1.6); g.add(cl);
    const rl = label(`${room.id}  ·  ${room.tokens} tok`, { bg: p.trim, fg: '#fff', bold: true, scale: 1.2 }); rl.position.set(-RW / 2 + 2.6, 4.9, -RD / 2); g.add(rl);
    world.add(g);
  });
  // Hallway in front: captain desk and whiteboard.
  const hall = mk(total + 4, 0.5, 4.5, '#F6D7A6'); hall.position.set(0, -0.45, RD / 2 + 3.4); world.add(hall);
  const cdesk = mk(3.2, 0.25, 1.4, '#E77859'); cdesk.position.set(0, 0.9, RD / 2 + 3.2); world.add(cdesk);
  const cap = person({ harness: 'captain', state: 'idle', shirt: '#24211d', mk, headSize: 0.7, bodyH: 0.9 }); cap.position.set(0, -0.2, RD / 2 + 4.3); cap.rotation.y = Math.PI; world.add(cap);
  const capl = label('Captain · you', { bold: true, bg: '#24211d', fg: '#fff' }); capl.position.set(0, 2.9, RD / 2 + 4.3); world.add(capl);
  const wb = mk(3.6, 2.2, 0.15, '#fff'); wb.position.set(-7, 1.8, RD / 2 + 1.5); world.add(wb);
  [0.6, 1.1, 0.8, 1.5, 1.2].forEach((hh, k) => { const bar = mk(0.4, hh, 0.1, ['#2FB37A', '#3B82F6', '#E8A93A', '#E5484D', '#8E7CC3'][k]); bar.position.set(-8.2 + k * 0.6, 0.9 + hh / 2, RD / 2 + 1.4); world.add(bar); });

  const cam = new THREE.OrthographicCamera(-1, 1, 1, -1, 0.1, 200);
  cam.position.set(30, 26, 30); cam.lookAt(0, 0.5, 1.5);
  const fit = () => { const a = innerWidth / innerHeight, h = 13; cam.left = -h * a; cam.right = h * a; cam.top = h; cam.bottom = -h; cam.updateProjectionMatrix(); };
  fit();
  const fx = new OutlineEffect(renderer, { defaultThickness: 0.0032, defaultColor: [0.16, 0.1, 0.12], defaultAlpha: 1 });
  animate((t) => { world.rotation.y = Math.sin(t * 0.25) * 0.06; });
  return { scene, cam, fit, render: () => fx.render(scene, cam), legend: ['Mỗi project một phòng mở nắp', 'Harness = đầu (icon app)', 'Viền đen, bóng 3 nấc', 'Màu kẹo, nền kem'] };
}

// ------------------------------------------------------------ variant B
// Voxel playground: perspective camera high above, one open floor of
// pastel tiles, everything built from cubes, carpets instead of walls,
// crews hop like Crossy Road, strong sun shadows.
function buildB() {
  const scene = new THREE.Scene();
  scene.background = new THREE.Color('#F5F3EB');
  scene.fog = new THREE.Fog('#F5F3EB', 40, 70);
  renderer.shadowMap.enabled = true; renderer.shadowMap.type = THREE.PCFShadowMap;
  const mk = (w, h, d, c) => { const m = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), new THREE.MeshLambertMaterial({ color: c })); m.castShadow = m.receiveShadow = true; return m; };
  scene.add(new THREE.HemisphereLight('#fff', '#d8c9b8', 1.1));
  const sun = new THREE.DirectionalLight('#fff', 2.4); sun.position.set(-10, 20, 8); sun.castShadow = true;
  sun.shadow.mapSize.set(2048, 2048); Object.assign(sun.shadow.camera, { left: -26, right: 26, top: 26, bottom: -26, near: 1, far: 60 }); scene.add(sun);

  const world = new THREE.Group(); scene.add(world);
  // Tiles.
  const tileA = new THREE.MeshLambertMaterial({ color: '#EFE9DC' }), tileB = new THREE.MeshLambertMaterial({ color: '#E6DFCF' });
  const tg = new THREE.BoxGeometry(2, 0.4, 2);
  for (let x = -12; x <= 12; x++) for (let z = -7; z <= 7; z++) { const t = new THREE.Mesh(tg, (x + z) % 2 ? tileA : tileB); t.position.set(x * 2, -0.2, z * 2); t.receiveShadow = true; world.add(t); }
  const CARPET = ['#B7E4C7', '#BBD7FF', '#FFD6A5'];
  const ACCENT = ['#2FB37A', '#3B82F6', '#F2742B'];
  const RW = 14;
  OFFICE.rooms.forEach((room, i) => {
    const cx = -RW * 1.1 + i * RW * 1.1;
    const carpet = mk(RW, 0.25, 10, CARPET[i]); carpet.position.set(cx, 0.12, -1); carpet.receiveShadow = true; world.add(carpet);
    // Room sign: a tall voxel totem.
    const sign = mk(1, 5, 1, ACCENT[i]); sign.position.set(cx - RW / 2 + 0.6, 2.5, -6.5); world.add(sign);
    const sl = label(`${room.id} · ${room.tokens}`, { bg: ACCENT[i], fg: '#fff', bold: true, scale: 1.3 }); sl.position.set(cx - RW / 2 + 0.6, 6.2, -6.5); world.add(sl);
    // Mate's desk: a big coral block with two monitors.
    const md = mk(5, 1.2, 2, '#E77859'); md.position.set(cx, 0.85, -4); world.add(md);
    for (const dx of [-1.4, 1.4]) { const mon = mk(1.4, 1, 0.3, '#24211d'); mon.position.set(cx + dx, 1.95, -4.3); world.add(mon); const scr = mk(1.2, 0.8, 0.1, room.mate.state === 'asleep' ? '#555' : '#9BE7FF'); scr.position.set(cx + dx, 1.95, -4.1); world.add(scr); }
    const mate = person({ harness: room.mate.harness, state: room.mate.state, shirt: ACCENT[i], mk, headSize: 1.35, bodyH: 1 }); mate.position.set(cx, 0.25, -2.6); mate.rotation.y = Math.PI; world.add(mate);
    const ml = label(`Mate · ${room.id}`, { dot: STATE[room.mate.state].color, bold: true }); ml.position.set(cx, 3.6, -2.6); world.add(ml);
    room.crews.forEach((c, j) => {
      const x = cx - 4 + j * 4, z = 1.5;
      const desk = mk(2.6, 1, 1.4, '#F7C98B'); desk.position.set(x, 0.75, z); world.add(desk);
      const lap = mk(1, 0.7, 0.2, '#24211d'); lap.position.set(x, 1.6, z - 0.3); world.add(lap);
      const scr = mk(0.85, 0.55, 0.05, c.state === 'working' ? '#9BE7FF' : '#666'); scr.position.set(x, 1.6, z - 0.17); world.add(scr);
      const who = person({ harness: c.harness, state: c.state, shirt: shirt(c.name), mk, headSize: 1.05, bodyH: 0.8 }); who.position.set(x, 0.25, z + 1.4); who.rotation.y = Math.PI; world.add(who);
      // Crossy hop.
      animate((t) => { const ph = (t * 1.6 + j) % 1; who.scale.y = 1 - Math.max(0, Math.sin(ph * Math.PI * 2)) * 0.08; });
      if (c.state === 'walking') animate((t) => { const k = (Math.sin(t * 0.6) + 1) / 2; who.position.x = x + (cx - x) * k; who.position.z = z + 1.4 + (-1.4 - (z + 1.4)) * k; });
      const nl = label(c.name, { dot: STATE[c.state].color, scale: 0.85 }); nl.position.set(x, 3.1, z + 1.4); world.add(nl);
    });
    // Cabinet: a stack of drawers.
    for (let k = 0; k < 4; k++) { const d = mk(1.6, 0.7, 1.4, k % 2 ? '#4DA8A0' : '#5FBDB5'); d.position.set(cx + RW / 2 - 1.3, 0.6 + k * 0.72, 2.5); world.add(d); }
    const cl = label(`${room.cabinet} done`, { bg: '#1F4B48', fg: '#fff', scale: 0.8 }); cl.position.set(cx + RW / 2 - 1.3, 4.3, 2.5); world.add(cl);
    // Plants.
    const pot = mk(0.8, 0.8, 0.8, '#C97B3B'); pot.position.set(cx + RW / 2 - 1.3, 0.65, -5); world.add(pot);
    const leaf = mk(1.3, 1.3, 1.3, '#3FA34D'); leaf.position.set(cx + RW / 2 - 1.3, 1.7, -5); world.add(leaf);
  });
  // Captain down the hall.
  const cd = mk(4, 1.1, 1.8, '#24211d'); cd.position.set(0, 0.8, 8); world.add(cd);
  const cap = person({ harness: 'captain', state: 'idle', shirt: '#E77859', mk, headSize: 1, bodyH: 1 }); cap.position.set(0, 0.25, 9.6); cap.rotation.y = Math.PI; world.add(cap);
  const capl = label('Captain · you', { bold: true, bg: '#24211d', fg: '#fff' }); capl.position.set(0, 3.7, 9.6); world.add(capl);
  const board = mk(5, 3, 0.4, '#fff'); board.position.set(-9, 2, 8.5); world.add(board);
  [0.6, 1.1, 0.8, 1.5, 1.2].forEach((hh, k) => { const bar = mk(0.5, hh, 0.2, ['#2FB37A', '#3B82F6', '#E8A93A', '#E5484D', '#8E7CC3'][k]); bar.position.set(-10.6 + k * 0.8, 0.9 + hh / 2, 8.2); world.add(bar); });

  const cam = new THREE.PerspectiveCamera(38, innerWidth / innerHeight, 0.1, 200);
  cam.position.set(0, 30, 26); cam.lookAt(0, 0, 1);
  animate((t) => { cam.position.x = Math.sin(t * 0.15) * 3; cam.lookAt(0, 0, 1); });
  return { scene, cam, fit: () => { cam.aspect = innerWidth / innerHeight; cam.updateProjectionMatrix(); }, render: () => renderer.render(scene, cam), legend: ['Sàn mở, thảm thay tường', 'Mọi thứ là khối vuông', 'Crew nhún như Crossy Road', 'Nền giấy + cam san hô của landing'] };
}

// ------------------------------------------------------------ variant C
// Clay studio: a cut-away dollhouse on a plinth, rounded clay materials,
// soft studio light, ACES tone mapping, slow orbit. Reads like a Blender
// still.
function buildC() {
  const scene = new THREE.Scene();
  scene.background = new THREE.Color('#F5F3EB');
  renderer.shadowMap.enabled = true; renderer.shadowMap.type = THREE.PCFSoftShadowMap;
  renderer.toneMapping = THREE.ACESFilmicToneMapping; renderer.toneMappingExposure = 1.05;
  const clay = (c) => new THREE.MeshStandardMaterial({ color: c, roughness: 0.62, metalness: 0 });
  const mk = (w, h, d, c) => { const m = new THREE.Mesh(new RoundedBoxGeometry(w, h, d, 3, Math.min(w, h, d) * 0.22), clay(c)); m.castShadow = m.receiveShadow = true; return m; };
  scene.add(new THREE.HemisphereLight('#fff', '#e0cfc0', 0.9));
  const key = new THREE.DirectionalLight('#fff6ea', 2.6); key.position.set(12, 18, 10); key.castShadow = true;
  key.shadow.mapSize.set(4096, 4096); key.shadow.radius = 6; Object.assign(key.shadow.camera, { left: -24, right: 24, top: 24, bottom: -24, near: 1, far: 60 }); scene.add(key);
  const fill = new THREE.DirectionalLight('#cfe3ff', 0.8); fill.position.set(-14, 8, -6); scene.add(fill);

  const world = new THREE.Group(); scene.add(world);
  const plinth = new THREE.Mesh(new THREE.CylinderGeometry(24, 25, 1.2, 64), clay('#E9E2D4')); plinth.position.y = -0.9; plinth.receiveShadow = true; world.add(plinth);
  const PAL = [{ floor: '#BFD8B8', wall: '#EAF2E6' }, { floor: '#C9D7F2', wall: '#EDF1FA' }, { floor: '#F6C9B8', wall: '#FBEDE7' }];
  const RW = 10.5, RD = 8.5, GAP = 0.6;
  const total = OFFICE.rooms.length * RW + (OFFICE.rooms.length - 1) * GAP;
  OFFICE.rooms.forEach((room, i) => {
    const p = PAL[i];
    const g = new THREE.Group(); g.position.x = -total / 2 + RW / 2 + i * (RW + GAP);
    const floor = mk(RW, 0.5, RD, p.floor); floor.position.y = -0.25; g.add(floor);
    const wallB = mk(RW, 4.5, 0.4, p.wall); wallB.position.set(0, 2.25, -RD / 2 + 0.2); g.add(wallB);
    const wallS = mk(0.4, 4.5, RD, p.wall); wallS.position.set(-RW / 2 + 0.2, 2.25, 0); g.add(wallS);
    const skirt = mk(RW, 0.35, 0.5, '#fff'); skirt.position.set(0, 0.18, -RD / 2 + 0.25); g.add(skirt);
    // Picture + window.
    const win = mk(2.8, 1.9, 0.2, '#CFE9FF'); win.position.set(2.5, 2.6, -RD / 2 + 0.42); g.add(win);
    const rug = mk(5, 0.08, 3.4, '#fff'); rug.position.set(-1, 0.04, 1.4); g.add(rug);
    // Mate desk.
    const md = mk(4, 0.3, 1.5, '#E77859'); md.position.set(-1.2, 1.05, -RD / 2 + 2.1); g.add(md);
    for (const dx of [-1.7, 1.7]) { const leg = mk(0.3, 0.95, 1.2, '#D96849'); leg.position.set(-1.2 + dx, 0.47, -RD / 2 + 2.1); g.add(leg); }
    for (const dx of [-1.1, 1.1]) { const mon = mk(1.1, 0.75, 0.12, '#2A2825'); mon.position.set(-1.2 + dx, 1.65, -RD / 2 + 1.8); g.add(mon); const scr = mk(0.95, 0.6, 0.06, room.mate.state === 'asleep' ? '#555' : '#DFF4FF'); scr.position.set(-1.2 + dx, 1.65, -RD / 2 + 1.73); g.add(scr); }
    const chair = mk(1.1, 1.5, 0.3, '#3C3833'); chair.position.set(-1.2, 1, -RD / 2 + 3.6); g.add(chair);
    const mate = person({ harness: room.mate.harness, state: room.mate.state, shirt: '#F4B6A3', mk, headSize: 1.05, bodyH: 0.85 }); mate.position.set(-1.2, 0, -RD / 2 + 3); mate.rotation.y = Math.PI; g.add(mate);
    const ml = label(`Mate · ${room.id}`, { dot: STATE[room.mate.state].color, bold: true, scale: 0.9 }); ml.position.set(-1.2, 3, -RD / 2 + 3); g.add(ml);
    room.crews.forEach((c, j) => {
      const x = -RW / 2 + 2 + j * 3, z = 1.5;
      const desk = mk(2.1, 0.25, 1.1, '#F3E3C9'); desk.position.set(x, 0.95, z); g.add(desk);
      for (const dx of [-0.85, 0.85]) { const leg = mk(0.2, 0.85, 0.85, '#D9C4A3'); leg.position.set(x + dx, 0.42, z); g.add(leg); }
      const lap = mk(0.8, 0.5, 0.1, '#2A2825'); lap.position.set(x, 1.3, z - 0.25); lap.rotation.x = -0.25; g.add(lap);
      const scr = mk(0.66, 0.38, 0.05, c.state === 'working' ? '#DFF4FF' : '#777'); scr.position.set(x, 1.3, z - 0.2); scr.rotation.x = -0.25; g.add(scr);
      const who = person({ harness: c.harness, state: c.state, shirt: shirt(c.name), mk, headSize: 0.8, bodyH: 0.7 }); who.position.set(x, 0, z + 0.9); who.rotation.y = Math.PI; g.add(who);
      if (c.state === 'walking') animate((t) => { const k = (Math.sin(t * 0.6) + 1) / 2; who.position.x = x + (-1.2 - x) * k; who.position.z = z + 0.9 + (-RD / 2 + 4.3 - (z + 0.9)) * k; });
      const nl = label(c.name, { dot: STATE[c.state].color, scale: 0.75 }); nl.position.set(x, 2.4, z + 0.9); g.add(nl);
    });
    const cab = mk(1.2, 2.3, 1, '#9FBFBC'); cab.position.set(RW / 2 - 1, 1.15, 1.5); g.add(cab);
    for (let k = 0; k < 3; k++) { const hd = mk(0.5, 0.1, 0.1, '#3C3833'); hd.position.set(RW / 2 - 1, 0.5 + k * 0.7, 2.05); g.add(hd); }
    const cl = label(`${room.cabinet} done`, { bg: '#3C3833', fg: '#fff', scale: 0.7 }); cl.position.set(RW / 2 - 1, 2.9, 1.5); g.add(cl);
    const leaf = new THREE.Mesh(new THREE.SphereGeometry(0.7, 24, 16), clay('#7FB77E')); leaf.position.set(RW / 2 - 1, 1.5, -RD / 2 + 1.1); leaf.castShadow = true; g.add(leaf);
    const pot = mk(0.6, 0.7, 0.6, '#D9A066'); pot.position.set(RW / 2 - 1, 0.35, -RD / 2 + 1.1); g.add(pot);
    const rl = label(`${room.id} · ${room.tokens} tok`, { bg: '#24211d', fg: '#fff', bold: true, scale: 1.05 }); rl.position.set(-RW / 2 + 2.4, 5.2, -RD / 2); g.add(rl);
    world.add(g);
  });
  const cd = mk(3.4, 0.3, 1.4, '#24211d'); cd.position.set(0, 0.95, RD / 2 + 4.6); world.add(cd);
  for (const dx of [-1.4, 1.4]) { const leg = mk(0.3, 0.85, 1.1, '#3C3833'); leg.position.set(dx, 0.42, RD / 2 + 4.6); world.add(leg); }
  const cap = person({ harness: 'captain', state: 'idle', shirt: '#E77859', mk, headSize: 0.7, bodyH: 0.9 }); cap.position.set(0, 0, RD / 2 + 5.7); cap.rotation.y = Math.PI; world.add(cap);
  const capl = label('Captain · you', { bold: true, bg: '#E77859', fg: '#fff', scale: 0.9 }); capl.position.set(0, 3, RD / 2 + 5.7); world.add(capl);
  const wb = mk(3.4, 2.1, 0.2, '#fff'); wb.position.set(-7, 1.9, RD / 2 + 1.6); world.add(wb);
  [0.5, 1, 0.7, 1.3, 1.1].forEach((hh, k) => { const bar = mk(0.36, hh, 0.12, ['#2FB37A', '#3B82F6', '#E8A93A', '#E5484D', '#8E7CC3'][k]); bar.position.set(-8.1 + k * 0.55, 1.05 + hh / 2, RD / 2 + 1.45); world.add(bar); });

  const cam = new THREE.PerspectiveCamera(26, innerWidth / innerHeight, 0.1, 200);
  const R = 38;
  animate((t) => { const a = -0.35 + Math.sin(t * 0.12) * 0.22; cam.position.set(Math.sin(a) * R, 24, Math.cos(a) * R); cam.lookAt(0, 1, 0); });
  cam.position.set(Math.sin(-0.35) * R, 24, Math.cos(-0.35) * R); cam.lookAt(0, 1, 0);
  return { scene, cam, fit: () => { cam.aspect = innerWidth / innerHeight; cam.updateProjectionMatrix(); }, render: () => renderer.render(scene, cam), legend: ['Mô hình cắt lát trên bệ tròn', 'Vật liệu đất sét bo góc', 'Ánh sáng studio mềm, ACES', 'Camera xoay chậm, còn nền giấy'] };
}

// ------------------------------------------------------------ run
const clock = new THREE.Clock();
let view;
function mount(k) {
  key = k;
  anims.length = 0;
  renderer.shadowMap.enabled = false; renderer.toneMapping = THREE.NoToneMapping; renderer.toneMappingExposure = 1;
  view = VARIANTS[k].build();
  document.getElementById('vkey').textContent = k;
  document.getElementById('vlabel').textContent = VARIANTS[k].name;
  document.getElementById('vname').textContent = `prototype 3D · ${k}`;
  const lg = document.getElementById('legend'); lg.innerHTML = '';
  for (const s of Object.values(STATE)) { const p = document.createElement('span'); p.className = 'pill'; const i = document.createElement('i'); i.style.background = s.color; p.append(i, s.word); lg.appendChild(p); }
  for (const s of view.legend) { const p = document.createElement('span'); p.className = 'pill'; p.textContent = s; p.style.background = '#24211d'; p.style.color = '#fff'; lg.appendChild(p); }
  const u = new URL(location); u.searchParams.set('variant', k); history.replaceState(null, '', u);
}
function step(dir) { const ks = Object.keys(VARIANTS); mount(ks[(ks.indexOf(key) + dir + ks.length) % ks.length]); }
document.getElementById('prev').onclick = () => step(-1);
document.getElementById('next').onclick = () => step(1);
addEventListener('keydown', (e) => { if (e.target.closest('input,textarea,[contenteditable]')) return; if (e.key === 'ArrowLeft') step(-1); if (e.key === 'ArrowRight') step(1); });
addEventListener('resize', () => { renderer.setSize(innerWidth, innerHeight); view && view.fit(); });

loadMascots().then(() => {
  document.getElementById('boot').remove();
  mount(key);
  renderer.setAnimationLoop(() => { const t = clock.getElapsedTime(); for (const f of anims) f(t); view.render(); });
});
