package autopilot

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// DefaultQuietAfter is how long the captain has to leave a manual project's
// Mate alone before the daemon turns auto mode back on (store.QuietAfter).
const DefaultQuietAfter = store.QuietAfter

// RearmText is the sent.log line the daemon appends when it turns auto mode
// back on.
func RearmText(quiet time.Duration) string {
	return fmt.Sprintf("auto mode on: the captain has not typed to the Mate or pressed m for %s", quiet)
}

// talk is the captain's last prompt to the Mate in one project's sent.log,
// read incrementally from offset.
type talk struct {
	offset int64
	asked  time.Time
}

func (t *talk) read(ws *store.Workspace, project string) error {
	entries, next, err := ws.ReadSent(project, t.offset)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Source == store.SourceUser && e.Target == store.TargetMate && e.Time.After(t.asked) {
			t.asked = e.Time
		}
	}
	t.offset = next
	return nil
}

// rearm turns auto mode back on for a project the captain held in manual
// with the `m` key, once they have been away for DefaultQuietAfter (docs/mvp.md
// M21). Manual is the captain at the Mate, so the clock is their last sign
// of being there: a prompt to the Mate, or the `m` press itself (the hold's
// mtime). Whether the Mate has answered does not matter: the outbox types a
// digest only into an empty composer, so a turn the captain is waiting on is
// never interrupted. A project with no `.auto` and no hold - never switched
// on, or a flag removed by hand - is left alone. The cursor, which only a
// delivered digest moves, hands the Mate whatever is still open once auto
// is back.
func (p *Pilot) rearm(project string, now time.Time) (bool, error) {
	info, err := os.Stat(p.ws.ManualHold(project))
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	t := p.talk[project]
	if t == nil {
		t = &talk{}
		p.talk[project] = t
	}
	if err := t.read(p.ws, project); err != nil {
		return false, err
	}
	last := info.ModTime()
	if t.asked.After(last) {
		last = t.asked
	}
	if now.Sub(last) < DefaultQuietAfter {
		return false, nil
	}
	if err := p.ws.SetMode(project, true); err != nil {
		return false, err
	}
	if err := p.ws.AppendSent(project, store.SentEntry{Time: now, Source: store.SourceApp, Target: store.TargetMate, Text: RearmText(DefaultQuietAfter)}); err != nil {
		return true, err
	}
	return true, nil
}
