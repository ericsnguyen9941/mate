package query

import (
	"os"
	"sync"
)

// The core loads the workspace every two seconds while a UI reads it
// (internal/core/snapshot.go), and most of what a load reads has not moved
// since the last one. Measured 2026-11-12 on the captain's workspace, after
// box's own memo (CON-259): one Load cost ~55 ms and ~45 MB, ~18 MB of it
// re-parsing a 6 MB `mate/.outbox` as JSON to learn a handful of assign
// keys, and most of the rest flattening a box of 8,700 entries the same way
// it was flattened two seconds before (CON-271).
//
// So a Load keeps, per file, what it derived from it, keyed by the file's
// size and mtime - the identity box's memo (internal/box/memo.go) and the
// timeline's caches already trust, because every writer of these files
// appends or replaces the file whole, which moves one of the two. A file
// whose stamp is what it was is stat'ed and not read.
//
// What is kept is exactly what the derivation would return again, and it
// is shared with every Snapshot that carries it: a reader never writes to a
// Snapshot's slices or maps (the core already hands one Snapshot to every
// UI), and every kept slice is clipped so an append copies.
//
// The memo is per process and keyed by absolute path, so it is bounded by
// the workspaces the process reads. Concurrent loads (the core's snapshot,
// FreshSnapshot, the dashboard's LoadBox) each read a consistent entry and
// the later write wins.

// stamp is a file's identity for the memo: its size and mtime. A file that
// cannot be stat'ed has none (ok false), and nothing is kept for it.
type stamp struct {
	size, mtime int64
	ok          bool
}

func statStamp(path string) stamp {
	fi, err := os.Stat(path)
	if err != nil {
		return stamp{}
	}
	return stamp{size: fi.Size(), mtime: fi.ModTime().UnixNano(), ok: true}
}

// memoCell keeps one value per path, with the key it was derived at.
type memoCell[K comparable, V any] struct {
	mu sync.Mutex
	m  map[string]memoEntry[K, V]
}

type memoEntry[K comparable, V any] struct {
	key K
	val V
}

// get is the value kept for path at key.
func (c *memoCell[K, V]) get(path string, key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.m[path]
	if !ok || e.key != key {
		var zero V
		return zero, false
	}
	return e.val, true
}

// held is whatever is kept for path, at whichever key.
func (c *memoCell[K, V]) held(path string) (K, V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.m[path]
	return e.key, e.val, ok
}

func (c *memoCell[K, V]) put(path string, key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.m == nil {
		c.m = map[string]memoEntry[K, V]{}
	}
	c.m[path] = memoEntry[K, V]{key: key, val: val}
}

func (c *memoCell[K, V]) forget(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, path)
}

// byStamp is derive(path) through cell: the value kept for path when the
// file's stamp is what it was, else derive's, kept when ok. A file that
// cannot be stat'ed is derived every time and forgotten.
func byStamp[V any](cell *memoCell[stamp, V], path string, derive func() (V, bool)) V {
	now := statStamp(path)
	if now.ok {
		if v, ok := cell.get(path, now); ok {
			return v
		}
	}
	v, ok := derive()
	if ok && now.ok {
		cell.put(path, now, v)
	} else {
		cell.forget(path)
	}
	return v
}
