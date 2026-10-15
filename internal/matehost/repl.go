package matehost

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// The REPL is what the captain sees in the Mate's pane: their messages and
// the Mate's replies. An app line is only a dim tag; its reply shows under
// the tag unless it was silent.

var (
	styleDim     = lipgloss.NewStyle().Faint(true)
	styleCaptain = lipgloss.NewStyle().Bold(true)
	styleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	styleTag     = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Faint(true)
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type hostMsg Event

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// Model is the REPL's bubbletea model.
type Model struct {
	host    *Host
	project string

	items   []Item
	running *Item
	tool    string
	queued  []string
	state   State

	input  []rune
	width  int
	height int
	scroll int
	frame  int
	note   string
}

// NewModel shows history (oldest first) and then what the host does.
func NewModel(host *Host, project string, history []Item) Model {
	return Model{host: host, project: project, items: history}
}

// Program runs the REPL; the host's events reach it through Send.
func Program(m Model) *tea.Program {
	p := tea.NewProgram(m, tea.WithAltScreen())
	m.host.Subscribe(func(ev Event) { p.Send(hostMsg(ev)) })
	return p
}

func (m Model) Init() tea.Cmd { return tick() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tickMsg:
		m.frame++
		return m, tick()
	case hostMsg:
		m.apply(Event(msg))
	case error:
		m.note = msg.Error()
	case tea.KeyMsg:
		return m.key(msg)
	}
	return m, nil
}

func (m *Model) apply(ev Event) {
	m.state = ev.State
	if ev.Tool != "" {
		m.tool = ev.Tool
	}
	if ev.Queued != "" {
		m.queued = append(m.queued, ev.Queued)
	}
	if it := ev.Item; it != nil {
		if !it.Done {
			cp := *it
			m.running, m.tool = &cp, ""
			if it.Source == SourceCaptain && len(m.queued) > 0 && m.queued[0] == it.Text {
				m.queued = m.queued[1:]
			}
		} else {
			m.running, m.tool = nil, ""
			m.items = append(m.items, *it)
			m.scroll = 0
		}
	}
	if m.state.Draft == "" && len(m.input) > 0 && ev.Item == nil && ev.Tool == "" && ev.Queued == "" {
		// The console cleared the draft (ClearDraft).
		m.input = nil
	}
	if m.state.Stopping {
		m.queued = nil
	}
}

func (m Model) key(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.note = ""
	switch k.Type {
	case tea.KeyEnter:
		text := strings.TrimSpace(string(m.input))
		if text == "" {
			return m, nil
		}
		m.input = nil
		m.host.SetDraft("")
		if !m.state.Busy && len(m.queued) == 0 {
			// Shown at once; the host's start event replaces it.
			m.running = &Item{At: time.Now(), Source: SourceCaptain, Text: text}
		}
		go m.host.Submit(text)
		return m, nil
	case tea.KeyEsc, tea.KeyCtrlC:
		if m.host.Interrupt() {
			m.note = "interrupting…"
		} else if k.Type == tea.KeyCtrlC {
			m.note = "this pane stays open; stop the Mate from the console"
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	case tea.KeyCtrlU:
		m.input = nil
	case tea.KeyPgUp:
		m.scroll += max(1, m.height/2)
		return m, nil
	case tea.KeyPgDown:
		m.scroll = max(0, m.scroll-max(1, m.height/2))
		return m, nil
	case tea.KeySpace:
		m.input = append(m.input, ' ')
	case tea.KeyRunes:
		m.input = append(m.input, k.Runes...)
	default:
		return m, nil
	}
	m.host.SetDraft(string(m.input))
	return m, nil
}

func (m Model) View() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	wrap := lipgloss.NewStyle().Width(w)
	var lines []string
	add := func(s string) { lines = append(lines, strings.Split(wrap.Render(s), "\n")...) }
	for _, it := range m.items {
		m.renderItem(it, add)
		lines = append(lines, "")
	}
	if it := m.running; it != nil {
		if it.Source == SourceCaptain {
			add(styleCaptain.Render("› " + it.Text))
		} else {
			add(styleTag.Render("· " + it.Label))
		}
		doing := "thinking"
		if m.tool != "" {
			doing = m.tool + "…"
		}
		add(styleDim.Render(spinnerFrames[m.frame%len(spinnerFrames)] + " " + doing))
		lines = append(lines, "")
	}
	for _, q := range m.queued {
		add(styleDim.Render("› " + q + "  (queued)"))
	}

	footer := []string{styleDim.Render(strings.Repeat("─", w))}
	footer = append(footer, strings.Split(wrap.Render("› "+string(m.input)+"█"), "\n")...)
	hint := "enter send · esc interrupt · pgup/pgdn scroll"
	if m.note != "" {
		hint = m.note
	} else if m.state.Stopping {
		hint = "stopping…"
	}
	footer = append(footer, styleDim.Render(fmt.Sprintf("mate %s · %s", m.project, hint)))

	h := m.height
	if h <= 0 {
		h = 24
	}
	room := max(1, h-len(footer))
	end := max(0, len(lines)-m.scroll)
	start := max(0, end-room)
	body := lines[start:end]
	for len(body) < room {
		body = append([]string{""}, body...)
	}
	return strings.Join(append(body, footer...), "\n")
}

func (m Model) renderItem(it Item, add func(string)) {
	if it.Source == SourceCaptain {
		add(styleCaptain.Render("› " + it.Text))
	} else {
		add(styleTag.Render("· " + it.Label))
		if it.Silent {
			return
		}
	}
	switch {
	case it.Error != "":
		add(styleError.Render(it.Error))
	case it.Interrupted && it.Reply == "":
		add(styleDim.Render("(interrupted)"))
	default:
		add(it.Reply)
	}
}
