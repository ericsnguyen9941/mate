package chain

import (
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

type screenRecord struct {
	Kind, Hash, Caller, Stub string
	Jev                      screen.PaneState
	Confidence               float64
	At                       time.Time
	Text                     string // jev.Prepare(pane, "")
}

// recordScreen keeps the text of a screen Jev was asked about for the
// captain to turn into stubs by hand (docs/plans/pane-state-jev-2026-10-13.md
// section 6): dir/<kind>/<hash12>.txt, a header then the prepared text Jev
// was sent. A screen seen again only moves last_seen and seen. Mate never
// reads these files to classify.
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
		for _, line := range strings.Split(string(old), "\n") {
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
	body := fmt.Sprintf("# first_seen: %s\n# last_seen: %s\n# seen: %d\n# harness: %s\n# caller: %s\n# stub: %s\n# jev: %s %.2f\n%s",
		first, at, seen, r.Kind, r.Caller, r.Stub, r.Jev, r.Confidence, r.Text)
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
