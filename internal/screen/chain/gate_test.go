package chain_test

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/notice"
	"github.com/nguyenngocanh94/mate/internal/screen/chain"
	"github.com/nguyenngocanh94/mate/internal/screen/fixture"
	"github.com/nguyenngocanh94/mate/internal/screen/jev"
)

const moduleRoot = "../../.."

// harnessOf is the harness a corpus file was captured from, the way
// scripts/jeveval labels it: its directory under internal/harness, else
// the prefix of its name. A notice fixture belongs to none.
func harnessOf(reg harness.Registry, file string) (harness.Kind, bool) {
	if rest, ok := strings.CutPrefix(file, "internal/harness/"); ok {
		k, _, _ := strings.Cut(rest, "/")
		return harness.Kind(k), true
	}
	if strings.Contains(file, "#") {
		return "", false
	}
	for _, k := range reg.Kinds() {
		if strings.HasPrefix(filepath.Base(file), string(k)+"_") {
			return k, true
		}
	}
	return "", false
}

// A regression pin on the 2026-10-08 corpus (threshold 0.85, Jev asked
// only about a screen the fixture cannot read, never for its composer):
// every harness screen, answered by Jev from its cassette when asked, goes
// through the chain. A screen the fixture reads is the fixture's reading
// whole, with no request and no log line; on the rest the composer stays
// the fixture's, so "dangerous" - a composer called empty that the fixture
// does not read empty - is zero by construction. The counts pin how often
// Jev is still asked and how often its answer is used. The notice fixtures
// have no harness and so no screen profile; the chain is never asked about
// a screen without one.
func TestChainGateOnTheCorpus(t *testing.T) {
	replay := jev.Replay{Dir: filepath.Join(moduleRoot, "scripts", "jeveval", "testdata", "cassette")}
	screens, err := replay.Screens(moduleRoot)
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	observer := chain.New(jev.New(notice.New("test-key").WithTransport(replay)), fixture.New(), 0.85,
		chain.WithLog(func(line string) error { lines = append(lines, line); return nil }))
	reg := catalog.Default()
	ctx := context.Background()
	counts := map[string]int{}
	var dangerous []string
	chained, noHarness := 0, 0
	for _, s := range screens {
		kind, ok := harnessOf(reg, s.File)
		if !ok {
			noHarness++
			continue
		}
		p, err := reg.Lookup(kind)
		if err != nil {
			t.Fatalf("%s: %v", s.File, err)
		}
		profile := p.Screen()
		det, err := fixture.New().Observe(ctx, profile, s.Text)
		if err != nil {
			t.Fatal(err)
		}
		before := len(lines)
		obs, err := observer.Observe(ctx, profile, s.Text)
		if err != nil {
			t.Fatalf("%s: %v", s.File, err)
		}
		chained++
		if det.Confidence == 1 {
			counts["fixture read it"]++
			if len(lines) != before || obs != det {
				t.Errorf("%s: the fixture read it, yet the chain returned %+v and logged %d lines", s.File, obs, len(lines)-before)
			}
			continue
		}
		if len(lines) != before+1 {
			t.Fatalf("%s: %d log lines for a screen the fixture cannot read", s.File, len(lines)-before)
		}
		logged, err := chain.ParseLogLine(lines[len(lines)-1])
		if err != nil {
			t.Fatal(err)
		}
		switch logged.Fallback {
		case "":
			counts["jev"]++
		default:
			counts["fixture ("+logged.Fallback+")"]++
		}
		if obs.Composer != det.Composer {
			dangerous = append(dangerous, fmt.Sprintf("%s: composer %s, fixture %s (%s)", s.File, obs.Composer, det.Composer, obs.Reason))
		}
		if obs.Draft != det.Draft || obs.Startup != det.Startup || obs.Evidence != det.Evidence || obs.Deterministic != det.Composer {
			t.Errorf("%s: %+v does not carry the fixture's %+v", s.File, obs, det)
		}
	}
	t.Logf("chain(jev cassette, fixture, 0.85) on %d harness screens (%d notice fixtures have no harness):", chained, noHarness)
	t.Logf("  %-28s %3d", "fixture read it, jev not asked", counts["fixture read it"])
	t.Logf("  %-28s %3d", "jev", counts["jev"])
	t.Logf("  %-28s %3d", "fixture ("+chain.FallbackThreshold+")", counts["fixture ("+chain.FallbackThreshold+")"])
	t.Logf("  %-28s %3d", "fixture ("+chain.FallbackError+")", counts["fixture ("+chain.FallbackError+")"])
	t.Logf("  %-28s %3d", "dangerous", len(dangerous))
	if chained != 78 || noHarness != 12 {
		t.Fatalf("%d harness screens and %d without a harness, want the corpus's 78 and 12", chained, noHarness)
	}
	// A request the cassette cannot answer falls back to the fixture, which
	// is never dangerous by this test's measure: the gate would pass
	// without measuring Jev at all.
	if n := counts["fixture ("+chain.FallbackError+")"]; n != 0 {
		t.Fatalf("%d requests failed: the gate measured the fixture, not the chain", n)
	}
	if counts["fixture read it"] != 76 || counts["jev"]+counts["fixture ("+chain.FallbackThreshold+")"] != 2 {
		t.Errorf("counts %v, want Jev asked about the corpus's 2 screens the fixture cannot read and no other", counts)
	}
	for _, d := range dangerous {
		t.Errorf("dangerous: %s", d)
	}
}
