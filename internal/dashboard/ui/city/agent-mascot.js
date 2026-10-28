import * as THREE from './vendor/three.js';
import { RoundedBoxGeometry, SVGLoader, mergeGeometries } from './vendor/three.js';
import { HARNESSES } from './harnesses.js';
import { Bubble, BUBBLE_STYLE } from './bubble.js';

// AgentMascot: a rounded matte block with a harness logo on its front face.
// It has no face or limbs; behaviour is motion only.
//
// Every pose is a pure function of (timeline, t). A video scene sets the
// timeline once and calls update(t) for any t, in any order. A web page calls
// play(state) as things happen and update(clockSeconds) each frame; play()
// stamps the segment with the last time update() saw, so both paths go
// through the same code.
//
// Units: size 1 is a body 1 wide. Place the mascot with root.position and
// turn it with `heading` (radians, 0 = logo facing +Z); keep the root itself
// unrotated so packets and aiming stay in its frame.

export { HARNESSES, BUBBLE_STYLE };

const BODY = { width: 1, height: 1.14, depth: 0.46, radius: 0.2 };
const LOGO_SIZE = 0.56; // largest logo dimension, as a share of body width
const HOVER = 0.1; // resting gap between body and floor
const PIVOT = 0.36 * BODY.height; // tilt and lean pivot, above the body's base
const BLEND = 0.35; // seconds to blend from one segment into the next

export const SEND = { emitAt: 0.62, flight: 0.9, duration: 2.0, turn: 0.75 };
const RECEIVE_DURATION = 1.4;
const COMPLETE = { duration: 1.6, landAt: 0.72, ring: 0.95 };
const SPEECH = { cps: 26, hold: 1.6, pop: 0.28, out: 0.2, typeDelay: 0.12 };
const THINK_BUBBLE = { delay: 0.25 }; // the thought bubble opens this long into thinking
const SLEEP_Z = { delay: 0.6, cycle: 2.7, count: 3, out: 0.3 }; // Zs rise this long into sleep, each cycle seconds

// One "Z" texture shared by every mascot.
let zTexture = null;
function sleepTexture() {
  if (zTexture) return zTexture;
  const c = document.createElement('canvas');
  c.width = c.height = 128;
  const ctx = c.getContext('2d');
  ctx.font = '900 104px ui-sans-serif, system-ui, sans-serif';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.lineJoin = 'round';
  ctx.lineWidth = 14;
  ctx.strokeStyle = '#FFFDF8';
  ctx.strokeText('Z', 64, 68);
  ctx.fillStyle = '#6B655C';
  ctx.fillText('Z', 64, 68);
  zTexture = new THREE.CanvasTexture(c);
  zTexture.colorSpace = THREE.SRGBColorSpace;
  return zTexture;
}

const EASE = {
  inOut: (x) => x * x * x * (x * (x * 6 - 15) + 10),
  out: (x) => 1 - (1 - x) ** 3,
  in: (x) => x * x * x,
  sine: (x) => 0.5 - 0.5 * Math.cos(Math.PI * x),
};
const clamp01 = (x) => Math.min(1, Math.max(0, x));
const TAU = Math.PI * 2;
const wrap = (a) => Math.atan2(Math.sin(a), Math.cos(a));

// Keyframe track: keys are [time, value, ease?], ease shaping the stretch
// that ends at that key.
function track(u, keys) {
  if (u <= keys[0][0]) return keys[0][1];
  for (let i = 1; i < keys.length; i++) {
    const [t1, v1, ease = 'inOut'] = keys[i];
    if (u <= t1) {
      const [t0, v0] = keys[i - 1];
      return v0 + (v1 - v0) * EASE[ease]((u - t0) / (t1 - t0));
    }
  }
  return keys.at(-1)[1];
}

const REST = { y: HOVER, yaw: 0, pitch: 0, roll: 0, stretch: 1, glow: 0 };
const KEYS = Object.keys(REST);

function mix(a, b, k) {
  const out = {};
  for (const key of KEYS) out[key] = a[key] + (b[key] - a[key]) * k;
  return out;
}

// State poses. `u` is seconds since the state began (looping states add the
// mascot's phase); `aim` is the yaw (relative to heading) toward the state's
// target or source, 0 when there is none; `since` is the true time since the
// state began, Infinity for the opening segment.
export const STATES = {
  sleeping: {
    loop: true,
    pose(u, aim, since) {
      const k = EASE.inOut(clamp01(since / 1.2)); // settle down over 1.2 s
      const b = Math.sin((TAU * u) / 4.6); // slow, deep breaths
      return {
        ...REST,
        y: HOVER + (-0.075 + 0.012 * b) * k,
        stretch: 1 + (-0.03 + 0.022 * b) * k,
        pitch: (0.16 + 0.025 * b) * k,
        roll: 0.05 * k,
      };
    },
  },
  idle: {
    loop: true,
    pose(u) {
      const s = Math.sin((TAU * u) / 3.2);
      return { ...REST, y: HOVER + 0.035 * s, stretch: 1 + 0.012 * Math.cos((TAU * u) / 3.2), roll: 0.015 * Math.sin((TAU * u) / 6.4) };
    },
  },
  thinking: {
    loop: true,
    pose(u) {
      return {
        ...REST,
        y: HOVER + 0.025 * Math.sin((TAU * u) / 4),
        roll: 0.11 * Math.sin((TAU * u) / 4.8),
        yaw: 0.35 * Math.sin((TAU * u) / 7.5),
        pitch: -0.05 - 0.03 * Math.sin((TAU * u) / 4.8 + 1),
        glow: 0.04 + 0.04 * Math.sin((TAU * u) / 2.4),
      };
    },
  },
  working: {
    loop: true,
    pose(u) {
      const beat = 0.48;
      const hop = Math.sin(Math.PI * ((u / beat) % 1)); // 0 at each tap, 1 mid-beat
      return {
        ...REST,
        y: HOVER * 0.6 + 0.045 * hop,
        stretch: 0.975 + 0.05 * hop,
        roll: 0.03 * Math.sin((Math.PI * u) / beat),
        pitch: 0.045,
      };
    },
  },
  sending: {
    duration: SEND.duration,
    pose(u, aim) {
      return {
        ...REST,
        // Turn most of the way, not all: a full turn hides the logo from a front camera.
        yaw: aim * track(u, [[0, 0], [0.3, 1], [1.55, 1], [2.0, 0]]),
        pitch: track(u, [[0, 0], [0.3, -0.1], [0.55, 0.28, 'out'], [0.95, 0.1], [1.4, 0]]),
        stretch: track(u, [[0, 1], [0.3, 0.96], [0.55, 1.04, 'out'], [0.9, 1]]),
        y: HOVER + track(u, [[0, 0], [0.3, 0.02], [0.55, -0.03], [1.0, 0]]),
        glow: track(u, [[0.5, 0], [0.62, 0.25, 'out'], [1.0, 0]]),
      };
    },
  },
  receiving: {
    duration: RECEIVE_DURATION,
    pose(u, aim) {
      return {
        ...REST,
        yaw: aim * track(u, [[0, 0.5], [0.3, 0.7], [1.0, 0.5], [1.4, 0]]),
        pitch: track(u, [[0, 0], [0.12, -0.16, 'out'], [0.4, 0.02], [0.62, 0.2], [0.85, -0.03], [1.05, 0.05], [1.4, 0]]),
        stretch: track(u, [[0, 1], [0.12, 0.95, 'out'], [0.35, 1.01], [0.6, 1]]),
        y: HOVER + track(u, [[0, 0], [0.12, 0.025, 'out'], [0.45, 0]]),
        roll: track(u, [[0, 0], [0.12, 0.03, 'out'], [0.5, 0]]),
        glow: track(u, [[0, 0], [0.1, 0.3, 'out'], [0.7, 0]]),
      };
    },
  },
  completed: {
    duration: COMPLETE.duration,
    pose(u) {
      return {
        ...REST,
        y: HOVER + track(u, [[0, 0], [0.18, -0.04], [0.45, 0.42, 'out'], [0.72, -0.03, 'in'], [0.9, 0.06, 'out'], [1.1, 0]]),
        stretch: track(u, [[0, 1], [0.18, 0.86], [0.3, 1.1, 'out'], [0.6, 1], [0.72, 0.9, 'in'], [0.86, 1.03], [1.1, 1]]),
        roll: track(u, [[0.3, 0], [0.45, 0.05], [0.65, -0.03], [0.9, 0]]),
        glow: track(u, [[0.3, 0], [0.5, 1, 'out'], [1.5, 0]]),
      };
    },
  },
};

// Soft round contact shadow, built from data so it needs no canvas or image.
function shadowTexture() {
  const n = 64;
  const data = new Uint8Array(n * n * 4);
  for (let y = 0; y < n; y++) {
    for (let x = 0; x < n; x++) {
      const r = Math.hypot((x + 0.5) / n - 0.5, (y + 0.5) / n - 0.5) * 2;
      data[(y * n + x) * 4 + 3] = Math.round(255 * (1 - EASE.inOut(clamp01(r))) ** 1.6);
    }
  }
  const texture = new THREE.DataTexture(data, n, n);
  texture.needsUpdate = true;
  return texture;
}

// SVG → flat logo geometry, LOGO_SIZE wide at most, centred on the origin,
// facing +Z. A colour array becomes a top-to-bottom vertex gradient.
function logoGeometry(svg, color) {
  const parts = [];
  for (const path of new SVGLoader().parse(svg).paths) {
    if (path.userData?.style?.fill === 'none') continue;
    for (const shape of path.toShapes()) parts.push(new THREE.ShapeGeometry(shape, 16));
  }
  const geometry = mergeGeometries(parts);
  for (const part of parts) part.dispose();
  geometry.scale(1, -1, 1); // SVG y points down
  const index = geometry.index.array;
  for (let i = 0; i < index.length; i += 3) [index[i], index[i + 2]] = [index[i + 2], index[i]];
  geometry.computeBoundingBox();
  const box = geometry.boundingBox;
  const span = Math.max(box.max.x - box.min.x, box.max.y - box.min.y);
  geometry.translate(-(box.min.x + box.max.x) / 2, -(box.min.y + box.max.y) / 2, 0);
  geometry.scale(LOGO_SIZE / span, LOGO_SIZE / span, 1);
  if (Array.isArray(color)) {
    const stops = color.map((c) => new THREE.Color(c));
    const position = geometry.attributes.position;
    const colors = new Float32Array(position.count * 3);
    const c = new THREE.Color();
    for (let i = 0; i < position.count; i++) {
      const k = clamp01(0.5 - position.getY(i) / LOGO_SIZE) * (stops.length - 1);
      const j = Math.min(Math.floor(k), stops.length - 2);
      c.copy(stops[j]).lerp(stops[j + 1], k - j).toArray(colors, i * 3);
    }
    geometry.setAttribute('color', new THREE.BufferAttribute(colors, 3));
  }
  return geometry;
}

const PACKETS = 3;

export class AgentMascot {
  /**
   * @param {object} options
   * @param {'claude'|'codex'|'pi'} [options.harness]
   * @param {string} [options.color] body colour, overriding the harness's
   * @param {string} [options.logo] SVG markup, overriding the harness's
   * @param {string|string[]} [options.logoColor] fill, or gradient stops top to bottom
   * @param {string} [options.accent] glow, packet and ring colour
   * @param {number} [options.glow] peak emissive strength of the light pulse
   * @param {number} [options.size] body width in scene units
   * @param {number[]|THREE.Vector3} [options.position]
   * @param {number} [options.heading] yaw in radians, 0 = logo toward +Z
   * @param {number} [options.phase] seconds added to looping states so a crowd does not move in step
   * @param {boolean} [options.shadow] draw the soft contact shadow
   * @param {number} [options.bubbleScale] speech bubble size; 1 sets the font at 2.4% of the viewport height
   * @param {number} [options.bubbleAlign] 0 centres the bubble over the head, 1 grows it rightward, -1 leftward
   * @param {number} [options.bubbleNear] camera distance past which bubbles shrink with distance (default: never)
   * @param {boolean} [options.keep] keep the whole timeline, for a page that scrubs back through it; by default
   *   segments and lines more than ten seconds older than the newest are dropped
   */
  constructor({ harness = 'claude', color, logo, logoColor, accent, glow, size = 1, position = [0, 0, 0], heading = 0, phase = 0, shadow = true, bubbleScale = 1, bubbleAlign = 0, bubbleNear = Infinity, keep = false } = {}) {
    this.size = size;
    this.keep = keep;
    this.heading = heading;
    this.phase = phase;
    this.segments = [{ state: 'idle', start: -Infinity }];
    this.lastT = 0;
    this.speeches = [];
    this.bubbleScale = bubbleScale;

    this.root = new THREE.Group();
    this.root.name = 'agent-mascot';
    this.root.position.fromArray(position.isVector3 ? position.toArray() : position);
    this.root.userData.mascot = this;

    this.ground = new THREE.Group();
    this.ground.scale.setScalar(size);
    this.root.add(this.ground);
    this.rig = new THREE.Group();
    this.rig.scale.setScalar(size);
    this.root.add(this.rig);
    this.hover = new THREE.Group();
    this.rig.add(this.hover);
    this.tilt = new THREE.Group();
    this.tilt.position.y = PIVOT;
    this.hover.add(this.tilt);
    this.squash = new THREE.Group();
    this.squash.position.y = -PIVOT;
    this.tilt.add(this.squash);

    this.bodyMaterial = new THREE.MeshStandardMaterial({ roughness: 0.6, metalness: 0 });
    this.body = new THREE.Mesh(new RoundedBoxGeometry(BODY.width, BODY.height, BODY.depth, 6, BODY.radius), this.bodyMaterial);
    this.body.name = 'body';
    this.body.position.y = BODY.height / 2;
    this.body.castShadow = this.body.receiveShadow = true;
    this.squash.add(this.body);

    this.logoMaterial = new THREE.MeshStandardMaterial({ roughness: 0.55, metalness: 0, polygonOffset: true, polygonOffsetFactor: -2, polygonOffsetUnits: -2 });
    this.logo = new THREE.Mesh(new THREE.BufferGeometry(), this.logoMaterial);
    this.logo.name = 'logo';
    this.logo.position.set(0, BODY.height / 2, BODY.depth / 2 + 0.002);
    this.squash.add(this.logo);

    this.shadowMaterial = new THREE.MeshBasicMaterial({ color: '#000000', map: shadowTexture(), transparent: true, depthWrite: false });
    this.shadow = new THREE.Mesh(new THREE.PlaneGeometry(1.5, 0.95), this.shadowMaterial);
    this.shadow.name = 'contact-shadow';
    this.shadow.rotation.x = -Math.PI / 2;
    this.shadow.position.y = 0.003;
    this.shadow.renderOrder = -1;
    this.shadow.visible = shadow;
    this.ground.add(this.shadow);

    this.ringMaterial = new THREE.MeshBasicMaterial({ transparent: true, depthWrite: false, side: THREE.DoubleSide });
    this.ring = new THREE.Mesh(new THREE.RingGeometry(0.47, 0.5, 96), this.ringMaterial);
    this.ring.name = 'pulse-ring';
    this.ring.rotation.x = -Math.PI / 2;
    this.ring.position.y = 0.006;
    this.ground.add(this.ring);

    this.packetMaterial = new THREE.MeshStandardMaterial({ roughness: 0.5, metalness: 0, emissiveIntensity: 0.25 });
    const packetGeometry = new RoundedBoxGeometry(0.2, 0.2, 0.2, 4, 0.06);
    this.packets = Array.from({ length: PACKETS }, (_, i) => {
      const packet = new THREE.Mesh(packetGeometry, this.packetMaterial);
      packet.name = `task-packet-${i}`;
      packet.castShadow = true;
      this.root.add(packet);
      return packet;
    });

    this.bubble = new Bubble({ align: bubbleAlign, near: bubbleNear });
    this.root.add(this.bubble.sprite);

    this.zs = Array.from({ length: SLEEP_Z.count }, (_, i) => {
      const z = new THREE.Sprite(new THREE.SpriteMaterial({ map: sleepTexture(), transparent: true, depthWrite: false }));
      z.name = `sleep-z-${i}`;
      z.visible = false;
      this.root.add(z);
      return z;
    });

    this.setHarness(harness, { color, logo, logoColor, accent, glow });
    this.update(0);
  }

  /** Swap identity in place; overrides fall back to the harness preset. */
  setHarness(harness, { color, logo, logoColor, accent, glow } = {}) {
    const preset = HARNESSES[harness];
    if (!preset && !(logo && color)) throw new Error(`Unknown harness ${harness}; pass logo and color`);
    this.harness = harness;
    this.accent = new THREE.Color(accent ?? preset?.accent ?? '#ffffff');
    this.glow = glow ?? preset?.glow ?? 0.4;
    this.bodyMaterial.color.set(color ?? preset.body);
    this.bodyMaterial.emissive.copy(this.accent);
    const fill = logoColor ?? preset?.logoColor ?? '#ffffff';
    this.logo.geometry.dispose();
    this.logo.geometry = logoGeometry(logo ?? preset.logo, fill);
    this.logoMaterial.vertexColors = Array.isArray(fill);
    this.logoMaterial.color.set(Array.isArray(fill) ? '#ffffff' : fill);
    this.logoMaterial.needsUpdate = true;
    this.ringMaterial.color.copy(this.accent);
    this.packetMaterial.color.copy(this.accent);
    this.packetMaterial.emissive.copy(this.accent);
    return this;
  }

  /**
   * Start `state` at time `at` (default: the last time passed to update()).
   * Segments at or after `at` are replaced. Options:
   * - target: mascot or world position the sending state leans toward and throws at
   * - from: mascot or world position the receiving state turns toward
   * - then: state to fall into after a one-shot ends (default idle)
   * - turn: share of the way sending turns toward its target (default SEND.turn)
   */
  play(state, { at = this.lastT, target, from, then = 'idle', turn } = {}) {
    if (!STATES[state]) throw new Error(`Unknown state ${state}`);
    if (!STATES[then]?.loop) throw new Error(`then must be a looping state, got ${then}`);
    this.segments = this.segments.filter((s) => s.start < at);
    // Keep history only as long as a packet or ring from it could still show.
    while (!this.keep && this.segments.length > 2 && this.segments[1].start < at - 10) this.segments.shift();
    this.segments.push({ state, start: at, target, from, then, turn });
    return this;
  }

  /**
   * Say `text` in a speech bubble from `at` (default: the last update() time).
   * It types out at `cps` characters a second, stays `hold` seconds, then
   * closes; a later line replaces it. Speech runs alongside any state.
   */
  say(text, { at = this.lastT, cps = SPEECH.cps, hold = SPEECH.hold, duration, style = BUBBLE_STYLE } = {}) {
    const typing = Array.from(text).length / cps;
    this.speeches = this.speeches.filter((s) => s.at < at && (this.keep || s.at + s.duration > at - 10));
    this.speeches.push({ text, at, cps, typing, duration: duration ?? SPEECH.typeDelay + typing + hold, style });
    return this;
  }

  /** The speech showing at time t, with seconds since it began, or null. */
  speechAt(t) {
    let s = null;
    for (const speech of this.speeches) if (speech.at <= t) s = speech;
    return s && t < s.at + s.duration ? { speech: s, u: t - s.at } : null;
  }

  /** Replace the whole timeline: [{state, start, target?, from?, then?}, ...]. */
  setTimeline(segments) {
    this.segments = [{ state: 'idle', start: -Infinity }];
    for (const s of [...segments].sort((a, b) => a.start - b.start)) this.play(s.state, { at: s.start, target: s.target, from: s.from, then: s.then, turn: s.turn });
    return this;
  }

  /** The state showing at the last update(), one-shots resolved to their `then`. */
  get state() {
    return this.stateAt(this.lastT);
  }

  stateAt(t) {
    const seg = this.segmentAt(t);
    const def = STATES[seg.state];
    return def.duration && t - seg.start > def.duration ? seg.then ?? 'idle' : seg.state;
  }

  segmentAt(t) {
    let i = 0;
    while (i + 1 < this.segments.length && this.segments[i + 1].start <= t) i++;
    return this.segments[i];
  }

  /** Pose at absolute time t, blended across segment changes. */
  poseAt(t) {
    let i = 0;
    while (i + 1 < this.segments.length && this.segments[i + 1].start <= t) i++;
    return this.blended(i, t, 3);
  }

  blended(i, t, depth) {
    const pose = this.raw(this.segments[i], t);
    const k = (t - this.segments[i].start) / BLEND;
    if (i === 0 || depth === 0 || k >= 1) return pose;
    return mix(this.blended(i - 1, t, depth - 1), pose, EASE.inOut(clamp01(k)));
  }

  raw(seg, t) {
    const def = STATES[seg.state];
    const u = t - seg.start;
    if (def.loop) return def.pose((Number.isFinite(u) ? u : t) + this.phase, 0, Number.isFinite(u) ? u : Infinity);
    const aim = seg.state === 'sending' ? (seg.turn ?? SEND.turn) * this.aim(seg.target) : this.aim(seg.from);
    if (u <= def.duration) return def.pose(Math.max(0, u), aim);
    const after = u - def.duration;
    return mix(def.pose(def.duration, aim), STATES[seg.then ?? 'idle'].pose(after + this.phase, 0, after), EASE.inOut(clamp01(after / BLEND)));
  }

  // Yaw from the logo's rest direction toward a mascot or world point.
  aim(who) {
    const p = this.localPoint(who);
    if (!p) return 0;
    return wrap(Math.atan2(p.x, p.z) - this.heading);
  }

  // A mascot's body centre, or a world point, in this root's frame.
  localPoint(who) {
    if (!who) return null;
    const world = who instanceof AgentMascot ? who.root.localToWorld(new THREE.Vector3(0, (HOVER + BODY.height / 2) * who.size, 0)) : new THREE.Vector3().fromArray(who.isVector3 ? who.toArray() : who);
    return this.root.worldToLocal(world);
  }

  /** Pose everything for absolute time t (seconds). Returns the pose. */
  update(t) {
    this.lastT = t;
    this.root.updateWorldMatrix(true, false);
    const pose = this.poseAt(t);
    const talk = this.speechAt(t);
    if (talk) {
      // Small nods while the words are coming out.
      const v = talk.u - SPEECH.typeDelay;
      const env = clamp01(v / 0.15) * clamp01((talk.speech.typing - v) / 0.25);
      pose.pitch += 0.035 * env * Math.sin(TAU * 3.2 * v);
      pose.stretch += 0.01 * env * Math.sin(TAU * 3.2 * v + 1);
    }
    this.rig.rotation.y = this.heading + pose.yaw;
    this.hover.position.y = pose.y;
    this.tilt.rotation.set(pose.pitch, 0, pose.roll);
    const side = 1 / Math.sqrt(pose.stretch);
    this.squash.scale.set(side, pose.stretch, side);
    this.bodyMaterial.emissiveIntensity = this.glow * pose.glow;

    const lift = clamp01((pose.y - HOVER * 0.5) / 0.5);
    this.shadow.scale.setScalar(1 - 0.35 * lift);
    this.shadowMaterial.opacity = 0.3 * (1 - 0.6 * lift);

    this.updateRing(t);
    this.updatePackets(t);
    this.updateBubble(t, talk, pose);
    this.updateSleep(t, pose);
    this.root.updateMatrixWorld(true);
    return pose;
  }

  // When the segment after `seg` starts, or Infinity.
  endOf(seg) {
    const next = this.segments[this.segments.indexOf(seg) + 1];
    return next ? next.start : Infinity;
  }

  updateRing(t) {
    let p = -1;
    for (const seg of this.segments) {
      const land = seg.start + COMPLETE.landAt;
      if (seg.state === 'completed' && land < this.endOf(seg) && t >= land && t < land + COMPLETE.ring) p = (t - land) / COMPLETE.ring;
    }
    this.ring.visible = p >= 0;
    if (p < 0) return;
    this.ring.scale.setScalar(1.1 + 1.1 * EASE.out(p));
    this.ringMaterial.opacity = 0.55 * (1 - p) ** 2;
  }

  updatePackets(t) {
    const flying = [];
    for (const seg of this.segments) {
      if (seg.state !== 'sending' || !seg.target) continue;
      const emit = seg.start + SEND.emitAt;
      if (emit < this.endOf(seg) && t >= emit && t < emit + SEND.flight) flying.push({ seg, p: (t - emit) / SEND.flight });
    }
    this.packets.forEach((packet, i) => {
      const fly = flying[flying.length - 1 - i];
      packet.visible = Boolean(fly);
      if (!fly) return;
      const { seg, p } = fly;
      const yaw = this.heading + (seg.turn ?? SEND.turn) * this.aim(seg.target);
      const start = new THREE.Vector3(0, (HOVER + BODY.height * 0.55) * this.size, (BODY.depth / 2 + 0.16) * this.size).applyAxisAngle(THREE.Object3D.DEFAULT_UP, yaw);
      const end = this.localPoint(seg.target);
      if (seg.target instanceof AgentMascot) end.addScaledVector(end.clone().sub(start).setY(0).normalize(), -(0.5 * seg.target.size + 0.06));
      const k = EASE.sine(p);
      packet.position.lerpVectors(start, end, k);
      packet.position.y += (0.2 + 0.22 * start.distanceTo(end)) * Math.sin(Math.PI * k);
      const pop = p < 0.14 ? EASE.out(p / 0.14) * (1 + 0.25 * Math.sin(Math.PI * p / 0.14)) : p > 0.86 ? 1 - EASE.in((p - 0.86) / 0.14) : 1;
      packet.scale.setScalar(this.size * Math.max(0.001, pop));
      packet.rotation.set(0.5 * Math.sin(Math.PI * k), yaw + TAU * 0.75 * k, 0.3 * k);
    });
  }

  // Speech wins over the thought bubble; the thought bubble shows through a
  // thinking segment and closes just after it ends.
  // Zs drift up and to the side from the head while sleeping, fading out just
  // after waking.
  updateSleep(t, pose) {
    let sleep = null;
    for (const seg of this.segments) {
      const end = this.endOf(seg);
      const from = seg.start + SLEEP_Z.delay;
      if (seg.state === 'sleeping' && t >= from && t < end + SLEEP_Z.out) sleep = { s: Number.isFinite(from) ? t - from : t + this.phase, fade: Math.min(Number.isFinite(from) ? clamp01((t - from) / 0.5) : 1, clamp01((end + SLEEP_Z.out - t) / SLEEP_Z.out)) };
    }
    const top = pose.y + BODY.height * pose.stretch;
    this.zs.forEach((z, i) => {
      z.visible = Boolean(sleep);
      if (!sleep) return;
      const p = (((sleep.s + (i * SLEEP_Z.cycle) / SLEEP_Z.count) % SLEEP_Z.cycle) + SLEEP_Z.cycle) % SLEEP_Z.cycle / SLEEP_Z.cycle;
      z.position.set((0.3 + 0.38 * p + 0.06 * Math.sin(TAU * 1.5 * p)) * this.size, (top + 0.12 + 0.8 * p) * this.size, 0);
      z.scale.setScalar((0.15 + 0.17 * p) * this.size);
      z.material.opacity = sleep.fade * Math.sin(Math.PI * p) ** 0.6;
      z.material.rotation = 0.25 - 0.2 * p;
    });
  }

  updateBubble(t, talk, pose) {
    const pop = (u) => (u < SPEECH.pop ? EASE.out(u / SPEECH.pop) * (1 + 0.18 * Math.sin((Math.PI * u) / SPEECH.pop)) : 1);
    this.bubble.sprite.position.set(0, (pose.y + BODY.height * pose.stretch + 0.14) * this.size, 0);
    if (talk) {
      const { speech, u } = talk;
      const chars = Math.min(Array.from(speech.text).length, Math.max(0, Math.floor((u - SPEECH.typeDelay) * speech.cps)));
      const close = EASE.inOut(clamp01((speech.duration - u) / SPEECH.out));
      return this.bubble.show({ text: speech.text, chars, style: speech.style }, this.bubbleScale * pop(u) * close);
    }
    let thought = null;
    for (const seg of this.segments) {
      const end = this.endOf(seg);
      if (seg.state === 'thinking' && t >= seg.start + THINK_BUBBLE.delay && t < end + SPEECH.out) thought = { u: t - seg.start - THINK_BUBBLE.delay, close: EASE.inOut(clamp01((end + SPEECH.out - t) / SPEECH.out)) };
    }
    if (!thought) return this.bubble.hide();
    this.bubble.show({ thought: true, phase: thought.u }, this.bubbleScale * pop(thought.u) * thought.close);
  }

  /** World-space points for framing cameras and checking poses. */
  anchors() {
    const top = this.body.localToWorld(new THREE.Vector3(0, BODY.height / 2, 0)).toArray();
    const logo = this.logo.getWorldPosition(new THREE.Vector3()).toArray();
    const packets = this.packets.filter((p) => p.visible).map((p) => p.getWorldPosition(new THREE.Vector3()).toArray());
    const bubble = this.bubble.sprite.visible ? this.bubble.sprite.getWorldPosition(new THREE.Vector3()).toArray() : null;
    return { top, logo, packets, bubble };
  }

  dispose() {
    this.root.removeFromParent();
    this.body.geometry.dispose();
    this.logo.geometry.dispose();
    this.shadow.geometry.dispose();
    this.shadowMaterial.map.dispose();
    this.ring.geometry.dispose();
    this.packets[0].geometry.dispose();
    this.bubble.dispose();
    for (const z of this.zs) z.material.dispose();
    for (const m of [this.bodyMaterial, this.logoMaterial, this.shadowMaterial, this.ringMaterial, this.packetMaterial]) m.dispose();
  }
}

/**
 * Hand a task from one mascot to another: the sender leans in and throws at
 * `at`, the receiver reacts as the packet lands. Returns the landing time.
 */
export function handoff(sender, receiver, at, { then = 'idle', receiverThen = 'idle', turn } = {}) {
  const land = at + SEND.emitAt + SEND.flight;
  sender.play('sending', { at, target: receiver, then, turn });
  receiver.play('receiving', { at: land - 0.04, from: sender, then: receiverThen });
  return land;
}
