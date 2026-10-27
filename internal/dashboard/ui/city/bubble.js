import * as THREE from './vendor/three.js';

// Speech and thought bubbles drawn on a canvas and shown as camera-facing
// sprites. A bubble's anchor is the tip of its tail; the mascot puts it just
// above its head. Size is set per camera just before drawing, so text reads the
// same near or far: at scale 1 the font is 2.4% of the viewport height.
//
// Bubbles that would overlap on screen are arranged once per rendered frame:
// first slid sideways (the tail stays on the head, the body moves), then, if
// that is not enough, the one further up the screen is lifted on a longer
// tail. The arrangement depends only on what is on screen in that frame, so a
// video can seek to any time and get the same picture.

const FONT = '600 44px ui-sans-serif, system-ui, -apple-system, "Segoe UI", sans-serif';
const PAD_X = 34;
const PAD_Y = 24;
const LINE = 58;
const MAX_TEXT = 500; // wrap width in layout px
const TAIL = 24;
const LIFT_MAX = 420; // tail room reserved below the body for lifting
const RADIUS = 32;
const RES = 2; // canvas px per layout px
const SCREEN = 0.024 / 44; // viewport heights per layout px at scale 1
const GAP = 0.012; // viewport heights kept between arranged bubbles
const THOUGHT = { w: 172, h: 84, tail: 40 };

export const BUBBLE_STYLE = { background: '#FFFDF8', color: '#24211D', border: 'rgba(36,33,29,0.16)', dots: '#6B655C' };

const live = new Set();
const arrangedFrame = new WeakMap(); // camera → renderer frame last arranged for
let ids = 0;

function canvas(w, h) {
  const c = typeof document !== 'undefined' ? document.createElement('canvas') : new OffscreenCanvas(w, h);
  c.width = w;
  c.height = h;
  return c;
}

// Greedy word wrap; a word wider than the line is broken by character.
function wrap(ctx, text) {
  const lines = [];
  for (const para of text.split('\n')) {
    let line = '';
    for (const word of para.split(/\s+/).filter(Boolean)) {
      const next = line ? `${line} ${word}` : word;
      if (ctx.measureText(next).width <= MAX_TEXT) { line = next; continue; }
      if (line) lines.push(line);
      line = '';
      for (const ch of word) {
        if (ctx.measureText(line + ch).width > MAX_TEXT && line) { lines.push(line); line = ''; }
        line += ch;
      }
    }
    lines.push(line);
  }
  return lines;
}

function sceneOf(object) {
  while (object.parent) object = object.parent;
  return object;
}

const ndc = new THREE.Vector3();

function arrange(scene, camera) {
  const aspect = camera.isPerspectiveCamera ? camera.aspect : (camera.right - camera.left) / (camera.top - camera.bottom);
  const items = [];
  for (const b of live) {
    if (!b.sprite.visible || sceneOf(b.sprite) !== scene) continue;
    ndc.setFromMatrixPosition(b.sprite.matrixWorld).project(camera);
    const u = 2 * SCREEN * b.screenScale(camera); // viewport heights (×2, as NDC) per layout px
    const [lo, hi] = b.tailRange();
    items.push({ b, u, ax: ndc.x * aspect, ay: ndc.y, w: b.w * u, hb: b.hb * u, tail: b.tailBase * u, minDx: (b.w * b.base - hi) * u, maxDx: (b.w * b.base - lo) * u, dx: 0, dy: 0 });
  }
  items.sort((p, q) => p.ax - q.ax || p.b.id - q.b.id);
  const rect = (i) => {
    const l = i.ax - i.b.base * i.w + i.dx;
    const bot = i.ay + i.tail + i.dy;
    return { l, r: l + i.w, bot, top: bot + i.hb };
  };
  const gap = 2 * GAP;
  for (let pass = 0; pass < 12; pass++) {
    let moved = false;
    for (let i = 0; i < items.length; i++) {
      for (let j = i + 1; j < items.length; j++) {
        const A = items[i];
        const B = items[j];
        const a = rect(A);
        const c = rect(B);
        if (a.r + gap <= c.l || c.r + gap <= a.l || a.top + gap <= c.bot || c.top + gap <= a.bot) continue;
        moved = true;
        const need = a.r + gap - c.l; // A sits left of B by anchor
        const roomA = A.dx - A.minDx;
        const roomB = B.maxDx - B.dx;
        if (roomA + roomB >= need) {
          let da = Math.min(need / 2, roomA);
          let db = need - da;
          if (db > roomB) { db = roomB; da = need - db; }
          A.dx -= da;
          B.dx += db;
        } else {
          const [under, over] = B.ay >= A.ay ? [A, B] : [B, A];
          over.dy += rect(under).top + gap - rect(over).bot;
        }
      }
    }
    if (!moved) break;
  }
  for (const i of items) i.b.place(camera, i.dx / i.u, i.dy / i.u);
}

export class Bubble {
  /**
   * align: 0 centres the bubble on its anchor, 1 grows it rightward, -1 leftward.
   * near: past this camera distance the bubble shrinks with distance (down to
   * 40%), so far-off chatter stays small; Infinity keeps one size everywhere.
   */
  constructor({ align = 0, near = Infinity } = {}) {
    this.near = near;
    this.id = ids++;
    this.base = 0.5 - 0.3 * Math.sign(align); // where the tail meets the body, as a share of its width
    this.scale = 0;
    this.layoutKey = null;
    this.drawnKey = null;
    this.material = new THREE.SpriteMaterial({ transparent: true, depthTest: false, depthWrite: false });
    this.sprite = new THREE.Sprite(this.material);
    this.sprite.name = 'bubble';
    this.sprite.renderOrder = 10;
    this.sprite.visible = false;
    this.sprite.frustumCulled = false; // its size is only known at draw time
    this.sprite.onBeforeRender = (renderer, scene, camera) => {
      const frame = renderer.info.render.frame;
      if (arrangedFrame.get(camera) === frame) return;
      arrangedFrame.set(camera, frame);
      arrange(scene, camera);
    };
    live.add(this);
  }

  /**
   * Speech: { text, chars, style } shows the first `chars` characters.
   * Thought: { thought: true, phase, style } shows three pulsing dots.
   */
  show(content, scale) {
    const style = content.style ?? BUBBLE_STYLE;
    const key = content.thought ? 'thought' : `say:${content.text}`;
    if (key !== this.layoutKey || style !== this.style) this.layout(key, content, style);
    this.content = content;
    this.scale = scale;
    this.sprite.visible = scale > 0.001;
  }

  hide() {
    this.sprite.visible = false;
  }

  // Size the bubble for its whole text, so it does not grow while typing.
  layout(key, content, style) {
    this.layoutKey = key;
    this.style = style;
    this.drawnKey = null;
    this.thought = Boolean(content.thought);
    if (this.thought) {
      this.w = THOUGHT.w;
      this.hb = THOUGHT.h;
      this.tailBase = THOUGHT.tail;
    } else {
      const probe = canvas(1, 1).getContext('2d');
      probe.font = FONT;
      this.lines = wrap(probe, content.text);
      this.widths = this.lines.map((l) => probe.measureText(l).width);
      this.w = Math.ceil(Math.max(...this.widths, LINE) + 2 * PAD_X);
      this.hb = this.lines.length * LINE + 2 * PAD_Y;
      this.tailBase = TAIL;
    }
    this.hTotal = this.hb + this.tailBase + LIFT_MAX;
    this.canvas = canvas(this.w * RES, this.hTotal * RES);
    this.material.map?.dispose();
    this.material.map = new THREE.CanvasTexture(this.canvas);
    this.material.map.colorSpace = THREE.SRGBColorSpace;
    this.material.map.anisotropy = 4;
    this.material.needsUpdate = true;
  }

  // Where along the body's width the tail may leave it.
  tailRange() {
    const inset = this.thought ? this.hb / 2 : RADIUS + 22;
    return [Math.min(inset, this.w / 2), Math.max(this.w - inset, this.w / 2)];
  }

  // Size for this camera and apply the arranged shift (dx) and lift (dy), in layout px.
  // Scale after distance falloff, for this camera.
  screenScale(camera) {
    if (!Number.isFinite(this.near)) return (this.eff = this.scale);
    const depth = Math.abs(ndc.setFromMatrixPosition(this.sprite.matrixWorld).applyMatrix4(camera.matrixWorldInverse).z);
    return (this.eff = this.scale * Math.min(1, Math.max(0.4, this.near / depth)));
  }

  place(camera, dx, dy) {
    const s = this.sprite;
    ndc.setFromMatrixPosition(s.matrixWorld).applyMatrix4(camera.matrixWorldInverse);
    const span = camera.isPerspectiveCamera ? (2 * Math.abs(ndc.z) * Math.tan(THREE.MathUtils.degToRad(camera.fov) / 2)) / camera.zoom : (camera.top - camera.bottom) / camera.zoom;
    const [lo, hi] = this.tailRange();
    const tx = Math.round(Math.min(Math.max(this.w * this.base - dx, lo), hi));
    const tail = Math.round(this.tailBase + Math.min(Math.max(dy, 0), LIFT_MAX));
    const shown = this.hb + tail;
    const map = this.material.map;
    map.repeat.set(1, shown / this.hTotal);
    map.offset.set(0, 1 - shown / this.hTotal);
    s.center.set(tx / this.w, 0);
    const k = span * SCREEN * this.eff;
    s.scale.set(this.w * k, shown * k, 1);
    s.updateMatrixWorld();
    const c = this.content;
    const key = this.thought ? `${tx}|${tail}|${Math.floor(c.phase * 12)}` : `${tx}|${tail}|${c.chars}`;
    if (key !== this.drawnKey) {
      this.drawnKey = key;
      this.draw(tx, tail);
    }
  }

  draw(tx, tail) {
    const { w, hb, style } = this;
    const ctx = this.canvas.getContext('2d');
    ctx.setTransform(RES, 0, 0, RES, 0, 0);
    ctx.clearRect(0, 0, w, this.hTotal);
    ctx.lineWidth = 3;
    ctx.fillStyle = style.background;
    ctx.strokeStyle = style.border;
    const tip = hb + tail - 2;
    const x0 = 2;
    const x1 = w - 2;
    const y0 = 2;
    const b = hb - 2;
    ctx.beginPath();
    if (this.thought) {
      const r = (b - y0) / 2;
      ctx.moveTo(x0 + r, y0);
      ctx.arcTo(x1, y0, x1, b, r);
      ctx.arcTo(x1, b, x0, b, r);
      ctx.arcTo(x0, b, x0, y0, r);
      ctx.arcTo(x0, y0, x1, y0, r);
      ctx.closePath();
      ctx.fill();
      ctx.stroke();
      // Trail of shrinking circles down to the head.
      for (const [f, rad] of [[0.32, 11], [0.78, 6.5]]) {
        ctx.beginPath();
        ctx.arc(tx, b + (tip - b) * f, rad, 0, Math.PI * 2);
        ctx.fill();
        ctx.stroke();
      }
      // Three dots that brighten in turn.
      const phase = Math.floor(this.content.phase * 12) / 12;
      ctx.fillStyle = style.dots;
      for (let i = 0; i < 3; i++) {
        const p = 0.5 + 0.5 * Math.sin(Math.PI * 2 * (phase * 1.25 - i * 0.2));
        ctx.globalAlpha = 0.3 + 0.7 * p;
        ctx.beginPath();
        ctx.arc(w / 2 + (i - 1) * 34, hb / 2 - 3 * p, 8 + 2 * p, 0, Math.PI * 2);
        ctx.fill();
      }
      ctx.globalAlpha = 1;
    } else {
      ctx.moveTo(x0 + RADIUS, y0);
      ctx.arcTo(x1, y0, x1, b, RADIUS);
      ctx.arcTo(x1, b, x0, b, RADIUS);
      ctx.lineTo(tx + 16, b);
      ctx.quadraticCurveTo(tx + 4, b + Math.min(6, tail / 4), tx, tip);
      ctx.quadraticCurveTo(tx - 6, b + Math.min(8, tail / 3), tx - 18, b);
      ctx.arcTo(x0, b, x0, y0, RADIUS);
      ctx.arcTo(x0, y0, x1, y0, RADIUS);
      ctx.closePath();
      ctx.fill();
      ctx.stroke();
      ctx.font = FONT;
      ctx.fillStyle = style.color;
      ctx.textBaseline = 'middle';
      ctx.textAlign = 'left'; // typed text grows rightward from where the full line starts
      let left = this.content.chars;
      this.lines.forEach((line, i) => {
        if (left <= 0) return;
        const glyphs = Array.from(line);
        const part = glyphs.slice(0, left).join('');
        left -= glyphs.length + 1; // the wrapped space counts as one
        ctx.fillText(part, this.lines.length === 1 ? (w - this.widths[0]) / 2 : PAD_X, PAD_Y + LINE * (i + 0.5));
      });
    }
    this.material.map.needsUpdate = true;
  }

  dispose() {
    live.delete(this);
    this.material.map?.dispose();
    this.material.dispose();
  }
}
