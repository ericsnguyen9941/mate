package agy

import (
	"strings"
	"unicode/utf8"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

// agyScreen is agy's ScreenProfile, measured on Antigravity CLI 1.5.10 in a
// 120x40 pane read as visible text, in the default (inline) rendering mode
// (testdata/screens).
//
// The composer is the row block between the last two full-width rules at
// the bottom of the pane, its first row opening with `> `. Scrollback
// repeats `> ` on the user's own lines, under a shorter rule and with no
// rule below them, so only the bottom block is the composer. A slash-command
// menu draws below the lower rule, and its rows open with `> ` too.
//
// A turn in flight leaves the composer empty and draws, above it, a braille
// spinner, two spaces and a status ending in "..." ("Generating...",
// "Editing files...", "Running command..."). The footer then reads
// "esc to cancel", but an open slash-command menu draws the same footer, so
// the spinner is the only evidence of a turn.
type agyScreen struct{}

// Kind implements harness.ScreenProfile.
func (agyScreen) Kind() harness.Kind { return KindAgy }

// ReadSource implements ScreenProfile.
func (agyScreen) ReadSource() harness.ReadSource { return harness.ReadVisible }

// ClassifyStartup implements ScreenProfile. Ready is the composer, empty or
// holding its placeholder, and no turn in flight. The launch pre-trusts the
// worktree (Build's --settings), so the folder-trust dialog is not answered;
// a screen that shows it is unrecognised.
func (agyScreen) ClassifyStartup(screen string) harness.StartupScreen {
	return agyStartup().Classify(screen)
}

// StartupAnswer implements ScreenProfile.
func (agyScreen) StartupAnswer(dialog harness.StartupScreen) (harness.StartupDialogAnswer, error) {
	return agyStartup().Answer(dialog)
}

// StartupTargetSelected implements ScreenProfile.
func (agyScreen) StartupTargetSelected(dialog harness.StartupScreen, screen string) bool {
	return agyStartup().TargetSelected(dialog, screen)
}

// ReadyScreen implements ScreenProfile: the smallest screen that classifies
// as an idle empty composer. A fake runtime shows it for a clean start.
func (agyScreen) ReadyScreen() string {
	return "───\n" + agyGlyph + "\n───\n"
}

// ComposerGlyph implements ScreenProfile.
func (agyScreen) ComposerGlyph() string { return agyGlyph }

// Composer implements ScreenProfile.
func (agyScreen) Composer(lines []string) (string, bool) {
	rows, ok := locateComposer(lines)
	if !ok {
		return "", false
	}
	return strings.Join(nonEmpty(rows), " "), true
}

// ComposerPlaceholders implements ScreenProfile. A fresh or resumed session
// shows the permission mode's hint in the empty composer; after the first
// turn it holds the glyph alone.
func (agyScreen) ComposerPlaceholders() []string { return []string{agyPlaceholder} }

// ComposerRows implements ScreenProfile.
func (agyScreen) ComposerRows(lines []string) ([]string, bool) {
	return locateComposer(lines)
}

// Busy implements ScreenProfile.
func (agyScreen) Busy(lines []string) (string, bool) {
	return busyEvidence(lines)
}

const (
	// agyGlyph is the prompt agy draws at the start of the composer.
	agyGlyph = ">"
	// agyPlaceholder is what the empty composer holds before the first
	// turn, under --dangerously-skip-permissions as under the auto mode.
	agyPlaceholder = "Auto mode: edits & safe commands auto-approved (shift+tab to cycle)"
	// agyStatusEnd ends the status line of a turn in flight.
	agyStatusEnd = "..."
	// agyRule is the rune agy draws its rules with (U+2500).
	agyRule = '─'
	// agyBoxScan bounds how far above the lower rule the upper one is
	// looked for.
	agyBoxScan = harness.BusyTailLines
)

// locateComposer returns the composer's rows, the first without its glyph,
// each trimmed. The composer is between the last rule on the screen and the
// rule above it, and its first row opens with the glyph. It is fail-closed:
// a shell after agy exited, the trust dialog, or a capture with no rule
// below the prompt has no composer.
func locateComposer(lines []string) ([]string, bool) {
	bottom := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if isRule(lines[i]) {
			bottom = i
			break
		}
	}
	if bottom < 0 {
		return nil, false
	}
	top := -1
	for i := bottom - 1; i >= 0 && bottom-i <= agyBoxScan; i-- {
		if isRule(lines[i]) {
			top = i
			break
		}
	}
	if top < 0 || top >= bottom-1 {
		return nil, false
	}
	first := strings.TrimSpace(lines[top+1])
	if first != agyGlyph && !strings.HasPrefix(first, agyGlyph+" ") {
		return nil, false
	}
	rows := []string{strings.TrimSpace(strings.TrimPrefix(first, agyGlyph))}
	for _, line := range lines[top+2 : bottom] {
		rows = append(rows, strings.TrimSpace(line))
	}
	return rows, true
}

// isRule reports whether a line is a horizontal rule and nothing else.
func isRule(line string) bool {
	t := strings.TrimSpace(line)
	if t == "" {
		return false
	}
	for _, r := range t {
		if r != agyRule {
			return false
		}
	}
	return true
}

func nonEmpty(rows []string) []string {
	var out []string
	for _, r := range rows {
		if r != "" {
			out = append(out, r)
		}
	}
	return out
}

// busyEvidence is the last trimmed line that is agy's in-flight status: one
// braille cell (U+2800–U+28FF), a space, and a status ending in "...".
func busyEvidence(lines []string) (string, bool) {
	evidence := ""
	for _, line := range lines {
		t := strings.TrimSpace(line)
		r, size := utf8.DecodeRuneInString(t)
		if size == 0 || r < 0x2800 || r > 0x28FF {
			continue
		}
		rest := t[size:]
		if strings.HasPrefix(rest, " ") && strings.HasSuffix(rest, agyStatusEnd) {
			evidence = t
		}
	}
	return evidence, evidence != ""
}

// agyStartup is agy's startup profile. No dialog is answered, so Answer
// refuses every one. The composer predicate is the empty composer, or its
// placeholder, with no turn in flight.
func agyStartup() harness.StartupProfile {
	return harness.StartupProfile{
		Kind: KindAgy,
		Composer: func(lines []string) bool {
			if _, busy := busyEvidence(lines); busy {
				return false
			}
			rows, ok := locateComposer(lines)
			if !ok {
				return false
			}
			content := strings.Join(nonEmpty(rows), " ")
			return content == "" || content == agyPlaceholder
		},
	}
}
