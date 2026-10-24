#!/bin/sh
# Copies three.js 0.186.1 (MIT) and the mascot GLBs from ~/app/threejsvideo.
# vendor/ and models/ are gitignored; run this once before serving.
set -e
cd "$(dirname "$0")"
T=$HOME/app/threejsvideo/node_modules/three
mkdir -p vendor/loaders vendor/utils vendor/effects vendor/geometries models
cp "$T/build/three.module.js" "$T/build/three.core.js" vendor/
cp "$T/examples/jsm/loaders/GLTFLoader.js" vendor/loaders/
cp "$T/examples/jsm/utils/BufferGeometryUtils.js" "$T/examples/jsm/utils/SkeletonUtils.js" vendor/utils/
cp "$T/examples/jsm/effects/OutlineEffect.js" vendor/effects/
cp "$T/examples/jsm/geometries/RoundedBoxGeometry.js" vendor/geometries/
for m in claude codex pi; do cp "$HOME/app/threejsvideo/references/export/$m/$m.glb" models/; done
echo "ok: python3 -m http.server 4180 --directory $(pwd)"
