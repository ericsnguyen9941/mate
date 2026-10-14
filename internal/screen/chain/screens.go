package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/screen"
)

var screensMu sync.Mutex

// screenHeaderLines is how many lines recordScreen writes before the text.
const screenHeaderLines = 7

// What the header's jev line says when Jev gave no answer.
const (
	keptNotAsked    = "not asked"
	keptCircuitOpen = "not asked, circuit open"
	keptFailed      = "failed"
)

type screenRecord struct {
	Kind, Hash, Caller, Stub string
	// Jev is Jev's answer as the header shows it, or why there is none.
	Jev  string
	At   time.Time
	Text string // the prepared text Jev was or would have been sent (sentText)
}

// keep records a screen the fixture could not read under WithScreens' dir,
// at most once per harness and screen within TTL: a pane that holds an
// unknown screen is read every poll, and one record a minute says as much.
func (c *Chain) keep(ctx context.Context, key cacheKey, pane, stub, jev string) {
	if c.screensDir == "" {
		return
	}
	now := c.now()
	c.mu.Lock()
	if at, ok := c.kept[key]; ok && now.Sub(at) < TTL {
		c.mu.Unlock()
		return
	}
	for k, at := range c.kept {
		if now.Sub(at) >= TTL {
			delete(c.kept, k)
		}
	}
	c.kept[key] = now
	c.mu.Unlock()
	_ = recordScreen(c.screensDir, screenRecord{Kind: key.kind, Hash: fmt.Sprintf("%x", key.hash[:6]),
		Caller: screen.CallerOf(ctx), Stub: stub, Jev: jev, At: now, Text: sentText(c.primary, pane)})
}

// recordScreen keeps the text of a screen the fixture could not read for
// the captain to turn into stubs by hand
// (docs/plans/pane-state-jev-2026-10-13.md section 6): dir/<kind>/<hash12>.txt,
// a header then the prepared text Jev was or would have been sent. A screen
// seen again only moves last_seen and seen. Mate never reads these files to
// classify.
func recordScreen(dir string, r screenRecord) error {
	screensMu.Lock()
	defer screensMu.Unlock()
	sub := filepath.Join(dir, r.Kind)
	if err := os.MkdirAll(sub, 0o755); err != nil {
		return err
	}
	path := filepath.Join(sub, r.Hash+".txt")
	at := r.At.UTC().Format(time.RFC3339)
	first, seen := at, 1
	if old, err := os.ReadFile(path); err == nil {
		// Only the header is read: the screen text after it may hold
		// lines that look like header lines.
		lines := strings.SplitN(string(old), "\n", screenHeaderLines+1)
		for _, line := range lines[:min(len(lines), screenHeaderLines)] {
			if !strings.HasPrefix(line, "#") {
				break
			}
			if v, ok := strings.CutPrefix(line, "# first_seen: "); ok {
				first = v
			} else if v, ok := strings.CutPrefix(line, "# seen: "); ok {
				if n, err := strconv.Atoi(v); err == nil {
					seen = n + 1
				}
			}
		}
	}
	body := fmt.Sprintf("# first_seen: %s\n# last_seen: %s\n# seen: %d\n# harness: %s\n# caller: %s\n# stub: %s\n# jev: %s\n%s",
		first, at, seen, r.Kind, r.Caller, r.Stub, r.Jev, r.Text)
	tmp, err := os.CreateTemp(sub, ".screen-*")
	if err != nil {
		return err
	}
	_, werr := tmp.WriteString(body)
	if cerr := tmp.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		os.Remove(tmp.Name())
		return werr
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}
