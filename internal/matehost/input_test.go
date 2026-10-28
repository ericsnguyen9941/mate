package matehost

import (
	"io"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestOpenSequence(t *testing.T) {
	for in, want := range map[string]bool{
		"":                      false,
		"abc":                   false,
		"\x1b":                  true,
		"ab\x1b":                true,
		"\x1bO":                 true,
		"\x1bOA":                false,
		"\x1b[":                 true,
		"\x1b[<64;35":           true,
		"\x1b[<64;35;39M":       false,
		"\x1b[<64;35;39M\x1b[<": true,
		"\x1b[A":                false,
		"\x1b[200~paste":        false,
		"\x1bx":                 false,
	} {
		if got := openSequence([]byte(in)); got != want {
			t.Errorf("openSequence(%q) = %v, want %v", in, got, want)
		}
	}
}

// inputLog is a model that keeps the keys and wheel notches it was sent.
type inputLog struct {
	keys  []string
	wheel int
}

func (l inputLog) Init() tea.Cmd { return nil }

func (l inputLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		l.keys = append(l.keys, msg.String())
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonWheelUp {
			l.wheel++
		}
	}
	return l, nil
}

func (l inputLog) View() string { return "" }

// runSplit feeds a Bubble Tea program two wheel notches whose bytes arrive
// in two writes, split inside the first report, and returns what it saw.
func runSplit(t *testing.T, wrap func(*os.File) io.Reader) inputLog {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	p := tea.NewProgram(inputLog{}, tea.WithInput(wrap(r)), tea.WithOutput(io.Discard), tea.WithoutSignalHandler())
	done := make(chan tea.Model, 1)
	go func() {
		m, _ := p.Run()
		done <- m
	}()
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte("\x1b[<64;35"))
	time.Sleep(5 * time.Millisecond)
	w.Write([]byte(";39M\x1b[<64;35;39M"))
	time.Sleep(300 * time.Millisecond)
	p.Quit()
	return (<-done).(inputLog)
}

func TestSeqInputKeepsSplitMouseReport(t *testing.T) {
	// The 5ms gap between the two writes is a sleep, which a loaded machine
	// stretches past the 50ms production wait (CON-132); the wait is not
	// what is tested, the joining is.
	wait := seqWait
	seqWait = 2 * time.Second
	t.Cleanup(func() { seqWait = wait })
	// The bare pipe shows the bug: the split report types into the model.
	bare := runSplit(t, func(f *os.File) io.Reader { return f })
	if len(bare.keys) == 0 {
		t.Fatalf("bare input: no keys leaked (wheel %d); the split no longer reproduces the bug", bare.wheel)
	}
	got := runSplit(t, func(f *os.File) io.Reader { return &seqInput{f} })
	if len(got.keys) != 0 || got.wheel != 2 {
		t.Fatalf("seqInput: keys %q, wheel %d; want no keys, wheel 2", got.keys, got.wheel)
	}
}
