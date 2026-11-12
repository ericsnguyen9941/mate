package timeline

import (
	"path/filepath"
	"sort"
)

// The files a poll found written past mate (CON-266): each open Crew's
// `.status`, which the Crew appends to itself, and the transcript the pass
// located for it and for the Mate, which its harness writes. The core
// watches them (internal/core/filewatch.go) so it hears of a write at once
// rather than at its next poll; the ingest is the one that knows where the
// transcripts are, so it says. Each committed poll of a project replaces
// that project's list.

// LiveKind is what a live file is to its agent.
type LiveKind string

const (
	// LiveStatus is a Crew's `crews/<c>.status`.
	LiveStatus LiveKind = "status"
	// LiveTranscript is a Crew's or the Mate's harness transcript.
	LiveTranscript LiveKind = "transcript"
)

// LiveFile is one file an open Crew or a running Mate writes.
type LiveFile struct {
	Project string
	// Crew is the Crew whose file it is; empty for the Mate's.
	Crew string
	Kind LiveKind
	// Path is clean and absolute as the pass read it.
	Path string
}

// LiveFiles is every project's live files as its latest committed poll
// found them, sorted by path. Safe to call from any goroutine.
func (i *Ingester) LiveFiles() []LiveFile {
	i.mu.Lock()
	defer i.mu.Unlock()
	var out []LiveFile
	for _, files := range i.live {
		out = append(out, files...)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Path < out[b].Path })
	return out
}

// Polled receives once a poll (Ingest) has ended since it was last
// received from, whatever it found: LiveFiles may have moved.
func (i *Ingester) Polled() <-chan struct{} {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.polled == nil {
		i.polled = make(chan struct{}, 1)
	}
	return i.polled
}

// notePolled ends a poll: the projects no longer in the workspace lose
// their live files, and Polled is signalled.
func (i *Ingester) notePolled(projects map[string]bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for p := range i.live {
		if !projects[p] {
			delete(i.live, p)
		}
	}
	if i.polled == nil {
		i.polled = make(chan struct{}, 1)
	}
	select {
	case i.polled <- struct{}{}:
	default:
	}
}

// recordLive keeps the live files of a committed poll's pass: each Crew not
// closed, its status file and its located transcript, and the Mate's
// transcript while its meta says it runs.
func (i *Ingester) recordLive(p *pass) {
	var files []LiveFile
	add := func(crew string, kind LiveKind, path string) {
		if path != "" {
			files = append(files, LiveFile{Project: p.project, Crew: crew, Kind: kind, Path: filepath.Clean(path)})
		}
	}
	if p.mate.Present && p.mateTranscript != "" {
		started, stopped := metaTime(p.mate.Meta, MetaStartedAt), metaTime(p.mate.Meta, MetaStoppedAt)
		if stopped.IsZero() || started.After(stopped) {
			add("", LiveTranscript, p.mateTranscript)
		}
	}
	for _, crew := range p.crews {
		if crew.state().Closed() {
			continue
		}
		add(crew.ID, LiveStatus, i.ws.CrewStatus(p.project, crew.ID))
		add(crew.ID, LiveTranscript, p.located[crew.ActorID])
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.live == nil {
		i.live = map[string][]LiveFile{}
	}
	i.live[p.project] = files
}
