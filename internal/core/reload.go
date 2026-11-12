package core

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/loopstat"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// configEvery is how often the engine looks at workspace.yaml's size and
// modification time, the safety net under the `workspace` notice (CON-270).
const configEvery = 30 * time.Second

// configReloader keeps the workspace the engine shares - with Mate Office
// and Agent City, the observer, the outbox loop and the auto daemon - in
// step with workspace.yaml. It is the engine's one reader of the file
// (CON-270): the loops it shares the handle with no longer read it each
// pass (their Deps.SharedConfig).
//
// It reads the file on a `workspace` notice (Engine.Changed: `mate project
// add`, a UI of another process, and, with CON-266, an edit by hand), once
// as it starts, and when a look every configEvery finds the file's size or
// modification time changed since it last read it - a write no notice told
// of. A change the engine makes itself goes through the same handle
// (store's updateConfig swaps it in) and needs no read. The dashboard's
// routes that list the projects (`/api/workspace`, `/api/board`,
// `/api/news`, `/api/work`) read the shared workspace and no snapshot
// (CON-208). It asks Herdr nothing. A file that does not read (missing,
// broken, empty) leaves the configuration loaded before, as LoadConfig
// does, and is read again once it changes.
type configReloader struct {
	ws    *store.Workspace
	every time.Duration
	// stat records each read of the file, not each look (CON-263); nil
	// records nothing.
	stat *loopstat.Loop

	// mu orders the reads, so the stamp is the one of the latest.
	mu sync.Mutex
	// stamp is the file as it was when it was last read; zero before.
	stamp fileStamp

	// cancel and done are set by Start and read by Stop: the engine calls
	// each once, Stop on a goroutine it starts after Start has returned.
	cancel context.CancelFunc
	done   chan struct{}
}

// fileStamp is what a look at the file compares: a write changes its size
// or its modification time, an atomic rename both.
type fileStamp struct {
	ok    bool
	size  int64
	mtime time.Time
}

func stampOf(path string) fileStamp {
	info, err := os.Stat(path)
	if err != nil {
		return fileStamp{}
	}
	return fileStamp{ok: true, size: info.Size(), mtime: info.ModTime()}
}

func configReload(ws *store.Workspace) *configReloader {
	return &configReloader{ws: ws, every: configEvery}
}

// reload reads workspace.yaml now. The stamp is taken before the read, so a
// write that lands during it is read again at the next look.
func (r *configReloader) reload() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stamp = stampOf(r.ws.WorkspaceFile())
	start := r.stat.Begin()
	_ = r.ws.LoadConfig()
	r.stat.End(start)
}

// look reads workspace.yaml when its size or modification time is not what
// it was at the last read, and reports whether it did.
func (r *configReloader) look() bool {
	r.mu.Lock()
	now := stampOf(r.ws.WorkspaceFile())
	same := now.ok == r.stamp.ok && now.size == r.stamp.size && now.mtime.Equal(r.stamp.mtime)
	r.mu.Unlock()
	if same {
		return false
	}
	r.reload()
	return true
}

func (r *configReloader) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	// The goroutine closes the channel made here, never r.done: a Stop
	// right after Start must find it whatever the goroutine has done by
	// then (the CON-207 race).
	done := make(chan struct{})
	r.cancel, r.done = cancel, done
	go func() {
		defer close(done)
		// Once at start: the file may have changed since the handle was
		// opened, and the stamp to compare with is taken here.
		r.reload()
		ticker := time.NewTicker(r.every)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r.look()
			}
		}
	}()
	return nil
}

// Stop ends the loop and waits for a read in flight. A loop that never
// started has nothing to stop.
func (r *configReloader) Stop() {
	if r.cancel == nil {
		return
	}
	r.cancel()
	<-r.done
}
