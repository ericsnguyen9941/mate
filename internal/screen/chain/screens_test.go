package chain_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/notice"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/screen/chain"
	"github.com/nguyenngocanh94/mate/internal/screen/jev"
)

func TestChainKeepsTheScreenJevWasAskedAbout(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 10, 13, 6, 27, 4, 0, time.UTC)
	jev := &stateStub{reading: screen.StateReading{State: screen.PaneHooksReview, Confidence: 0.94, Source: "jev"}}
	fix := &stateStub{reading: screen.StateReading{State: screen.PaneUnknown, Startup: harness.StartupScreenUnrecognized, Source: "fixture"}}
	ch := chain.New(jev, fix, 0.85, chain.WithClock(func() time.Time { return now }), chain.WithScreens(dir))
	ctx := screen.WithCaller(context.Background(), screen.CallerSettle)
	pane := "  Hooks need review\n› 1. Review hooks\n"
	if _, err := ch.(screen.StateObserver).ObserveState(ctx, codexScreen, pane); err != nil {
		t.Fatal(err)
	}
	files, _ := filepath.Glob(filepath.Join(dir, "codex", "*.txt"))
	if len(files) != 1 {
		t.Fatalf("files = %v", files)
	}
	now = now.Add(2 * time.Minute) // past the cache, so Jev is asked again
	if _, err := ch.(screen.StateObserver).ObserveState(ctx, codexScreen, pane); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(files[0])
	got := string(b)
	for _, want := range []string{"# first_seen: 2026-10-13T06:27:04Z", "# last_seen: 2026-10-13T06:29:04Z", "# seen: 2",
		"# harness: codex", "# caller: settle", "# stub: unrecognized", "# jev: dialog_hooks_review 0.94", "› 1. Review hooks"} {
		if !strings.Contains(got, want) {
			t.Errorf("record lacks %q:\n%s", want, got)
		}
	}
}

func TestChainKeepsNoScreenTheStubRead(t *testing.T) {
	dir := t.TempDir()
	fix := &stateStub{reading: screen.StateReading{State: screen.PaneReady, Confidence: 1, Source: "fixture"}}
	ch := chain.New(&stateStub{}, fix, 0.85, chain.WithScreens(dir))
	_, _ = ch.(screen.StateObserver).ObserveState(context.Background(), codexScreen, "ready")
	if files, _ := filepath.Glob(filepath.Join(dir, "*", "*")); len(files) != 0 {
		t.Fatalf("stub-read screen recorded: %v", files)
	}
}

func TestChainKeepsNoScreenWhenJevFailed(t *testing.T) {
	dir := t.TempDir()
	jev := &stateStub{err: context.DeadlineExceeded}
	fix := &stateStub{reading: screen.StateReading{State: screen.PaneUnknown, Source: "fixture"}}
	ch := chain.New(jev, fix, 0.85, chain.WithScreens(dir))
	_, _ = ch.(screen.StateObserver).ObserveState(context.Background(), codexScreen, "x")
	if files, _ := filepath.Glob(filepath.Join(dir, "*", "*")); len(files) != 0 {
		t.Fatalf("failed request recorded: %v", files)
	}
}

func TestChainScreensSurviveConcurrentSettles(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 10, 13, 6, 27, 4, 0, time.UTC)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			jev := &stateStub{reading: screen.StateReading{State: screen.PaneHooksReview, Confidence: 0.94, Source: "jev"}}
			fix := &stateStub{reading: screen.StateReading{State: screen.PaneUnknown, Source: "fixture"}}
			ch := chain.New(jev, fix, 0.85, chain.WithClock(func() time.Time { return now }), chain.WithScreens(dir))
			_, _ = ch.(screen.StateObserver).ObserveState(context.Background(), codexScreen, "same pane")
		}()
	}
	wg.Wait()
	files, _ := filepath.Glob(filepath.Join(dir, "codex", "*"))
	if len(files) != 1 || !strings.HasSuffix(files[0], ".txt") {
		t.Fatalf("files = %v", files)
	}
	b, err := os.ReadFile(files[0])
	if err != nil || !strings.Contains(string(b), "# seen: 8\n") {
		t.Fatalf("record %v:\n%s", err, b)
	}
}

// The text kept for the captain is the text Jev was sent: the Jev client's
// own key, should it be on screen, is redacted there exactly as in the
// request, and never written to disk.
func TestChainKeepsTheScreenRedactedAsSent(t *testing.T) {
	dir := t.TempDir()
	const key = "jev-test-key-1234567890"
	var sent string
	client := notice.New(key).WithTransport(roundTrip(func(r *http.Request) (*http.Response, error) {
		var body struct {
			State string `json:"state"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		sent = body.State
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(stateAnswer(t, "dialog_other", 0.9)))}, nil
	}))
	fix := &stateStub{reading: screen.StateReading{State: screen.PaneUnknown, Startup: harness.StartupScreenUnrecognized, Source: "fixture"}}
	ch := chain.New(jev.New(client), fix, 0.85, chain.WithScreens(dir))
	pane := "export JEV_API_KEY=" + key + "\n› 1. Pick a model\n  2. Keep current\n"
	if _, err := ch.(screen.StateObserver).ObserveState(context.Background(), codexScreen, pane); err != nil {
		t.Fatal(err)
	}
	files, _ := filepath.Glob(filepath.Join(dir, "codex", "*.txt"))
	if len(files) != 1 {
		t.Fatalf("files = %v", files)
	}
	b, _ := os.ReadFile(files[0])
	_, text, _ := strings.Cut(string(b), "# jev: dialog_other 0.90\n")
	if strings.Contains(string(b), key) || !strings.Contains(text, "[redacted]") || text != sent {
		t.Fatalf("kept\n%s\nwant exactly what was sent\n%s", b, sent)
	}
}

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func stateAnswer(t *testing.T, choice string, conf float64) string {
	t.Helper()
	probs := map[string]float64{}
	for _, s := range screen.PaneStates {
		probs[string(s)] = (1 - conf) / float64(len(screen.PaneStates)-1)
	}
	probs[choice] = conf
	b, err := json.Marshal(map[string]any{"model": notice.Model, "answers": map[string]any{
		string(jev.AxisState): map[string]any{"type": "choice", "choice": choice, "confidence": conf, "probabilities": probs}}})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
