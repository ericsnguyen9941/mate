package matehost

import (
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

// The composer is a rune slice and a cursor into it. Lines are split by
// '\n'; the readline keys move and delete within the line the cursor is on.

func (m *Model) insert(r rune) {
	in := make([]rune, 0, len(m.input)+1)
	in = append(in, m.input[:m.cursor]...)
	in = append(in, r)
	m.input = append(in, m.input[m.cursor:]...)
	m.cursor++
}

func (m *Model) cut(from, to int) {
	if from >= to {
		return
	}
	m.input = append(m.input[:from:from], m.input[to:]...)
	m.cursor = from
}

// edit applies a key that changes the composer and reports whether it did.
func (m *Model) edit(k tea.KeyMsg) bool {
	switch k.Type {
	case tea.KeySpace:
		m.insert(' ')
	case tea.KeyBackspace:
		if m.cursor == 0 {
			return false
		}
		if k.Alt {
			// ⌥⌫: the word before the cursor.
			m.cut(m.wordStart(), m.cursor)
			break
		}
		m.cut(m.cursor-1, m.cursor)
	case tea.KeyDelete, tea.KeyCtrlD:
		if m.cursor >= len(m.input) {
			return false
		}
		m.cut(m.cursor, m.cursor+1)
	case tea.KeyLeft, tea.KeyCtrlB:
		m.cursor = max(0, m.cursor-1)
		return false
	case tea.KeyRight, tea.KeyCtrlF:
		m.cursor = min(len(m.input), m.cursor+1)
		return false
	case tea.KeyCtrlA:
		m.cursor = m.lineStart()
		return false
	case tea.KeyCtrlE:
		m.cursor = m.lineEnd()
		return false
	case tea.KeyCtrlW:
		m.cut(m.wordStart(), m.cursor)
	case tea.KeyCtrlU:
		// The whole composer, every line: ⌘⌫ sends it in most macOS
		// terminals.
		if len(m.input) == 0 {
			return false
		}
		m.cut(0, len(m.input))
	case tea.KeyCtrlK:
		m.cut(m.cursor, m.lineEnd())
	default:
		return false
	}
	return true
}

// wordStart is where the word before the cursor begins, spaces before it
// included.
func (m Model) wordStart() int {
	i := m.cursor
	for i > 0 && unicode.IsSpace(m.input[i-1]) {
		i--
	}
	for i > 0 && !unicode.IsSpace(m.input[i-1]) {
		i--
	}
	return i
}

func (m Model) lineStart() int {
	i := m.cursor
	for i > 0 && m.input[i-1] != '\n' {
		i--
	}
	return i
}

func (m Model) lineEnd() int {
	i := m.cursor
	for i < len(m.input) && m.input[i] != '\n' {
		i++
	}
	return i
}

// lineUp and lineDown move the cursor to the same column of the line above
// or below, or to its end when that line is shorter.
func (m *Model) lineUp() {
	start := m.lineStart()
	if start == 0 {
		m.cursor = 0
		return
	}
	col := m.cursor - start
	prevEnd := start - 1
	prevStart := prevEnd
	for prevStart > 0 && m.input[prevStart-1] != '\n' {
		prevStart--
	}
	m.cursor = min(prevStart+col, prevEnd)
}

func (m *Model) lineDown() {
	end := m.lineEnd()
	if end >= len(m.input) {
		m.cursor = len(m.input)
		return
	}
	col := m.cursor - m.lineStart()
	nextStart := end + 1
	nextEnd := nextStart
	for nextEnd < len(m.input) && m.input[nextEnd] != '\n' {
		nextEnd++
	}
	m.cursor = min(nextStart+col, nextEnd)
}
