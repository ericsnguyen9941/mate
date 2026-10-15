package matehost

import (
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func sampleModel(w, h int) Model {
	at := time.Date(2026, 10, 15, 14, 2, 0, 0, time.Local)
	items := []Item{
		{At: at, Source: SourceCaptain, Text: "ok anh đã có 2 test user là gmail của anh.", Done: true,
			Reply: "Em đã ghi nhận việc anh thêm 2 Gmail làm test user. Phía Google đủ để bắt đầu rồi. Em đã đánh thức crew để chạy lại luồng OAuth.",
			Tools: []ToolRecord{{Name: "Bash", Detail: "mate wake hellovietnam", Done: true, MS: 900}}, MS: 4100, Tokens: 1200, Context: 84000},
		{At: at.Add(time.Minute), Source: SourceApp, Label: "digest", Text: Marker + "digest: 1 item(s) — e1 budget 50%: 10.1M of 20M — status files under /x; act per AGENTS.md section 10", Reply: "—", Silent: true, Done: true, MS: 2000},
		{At: at.Add(3 * time.Minute), Source: SourceCaptain, Text: "Trong lúc chờ, mình bàn màn Hồ sơ nhé. Anh muốn có ảnh đại diện và tên hiển thị, chưa cần đổi mật khẩu.", Done: true,
			Reply: "Được ạ. Em đề xuất màn **Hồ sơ** gồm:\n- Ảnh đại diện: chọn từ máy, cắt vuông, tối đa 2 MB\n- Tên hiển thị: 2–40 ký tự, cho phép tiếng Việt có dấu\n- Email: chỉ đọc, lấy từ `auth.users`\n\nAnh muốn em tạo task cho crew mới luôn, hay chờ màn đăng nhập merge xong rồi mới làm?", MS: 2800, Tokens: 900},
		{At: at.Add(19 * time.Minute), Source: SourceCaptain, Text: "Tóm tắt giúp anh các crew đang chạy.", Done: true, Error: "exit status 1: API Error: 529 overloaded"},
	}
	m := NewModel(nil, Config{Project: "hellovietnam", Model: "opus"}, items)
	m.now = func() time.Time { return at.Add(20 * time.Minute) }
	m.crews = 3
	m.width, m.height = w, h
	return m
}

func busyModel(w, h int) Model {
	m := sampleModel(w, h)
	m.items = m.items[:3]
	start := m.now().Add(-7 * time.Second)
	m.state = State{Busy: true, TurnSource: SourceCaptain, TurnStartedAt: start}
	m.running = &Item{At: start, Source: SourceCaptain, Text: "Làm luôn đi, tách thành task riêng.",
		Reply: "Em tạo task e1-profile gồm 3 phần như trên và giao cho một crew mới. Crew làm trên nhánh riêng nên không đụng vào",
		Tools: []ToolRecord{{ID: "1", Name: "Bash", Detail: "mate task new e1-profile --from-chat", Done: true, MS: 800}, {ID: "2", Name: "Bash", Detail: "mate crew spawn e1-profile", Started: start.Add(4 * time.Second)}}}
	m.queue = []string{`à nhớ đổi tên hiển thị của app thành "Hello Vietnam" nhé`, "xong thì tóm tắt giúp anh các crew đang chạy"}
	m.input = []rune("Ảnh đại diện lưu ở đâu em,\nSupabase Storage hay S3?")
	m.cursor = len(m.input)
	return m
}

// TestREPLFitsTheWidth draws the design's screens at 80 and 100 columns:
// no line may be wider than the pane, and the screen is exactly its height.
func TestREPLFitsTheWidth(t *testing.T) {
	for _, w := range []int{80, 100, 60} {
		for name, m := range map[string]Model{"idle": sampleModel(w, 30), "busy": busyModel(w, 30)} {
			view := m.View()
			lines := strings.Split(view, "\n")
			if len(lines) != 30 {
				t.Errorf("%s@%d: %d lines, want 30", name, w, len(lines))
			}
			for i, l := range lines {
				if got := ansi.StringWidth(l); got > w {
					t.Errorf("%s@%d line %d is %d wide: %q", name, w, i, got, ansi.Strip(l))
				}
			}
			if os.Getenv("MATE_REPL_SHOW") != "" {
				t.Logf("%s@%d\n%s", name, w, ansi.Strip(view))
			}
		}
	}
}

func TestREPLScreens(t *testing.T) {
	idle := ansi.Strip(sampleModel(100, 40).View())
	for _, want := range []string{"mate · hellovietnam", "crew 3 · ctx 84.0k · opus", "› ok anh đã có 2 test user", "● Em đã ghi nhận", "▸ 1 tool · Bash: mate wake hellovietnam · 4.1s · 1.2k tok", "• Ảnh đại diện", "✗ Lượt này lỗi · exit status 1: API Error: 529 overloaded", "^R thử lại", "○ Mate rảnh"} {
		if !strings.Contains(idle, want) {
			t.Errorf("idle screen lacks %q:\n%s", want, idle)
		}
	}
	if strings.Contains(idle, "· digest ·") {
		t.Errorf("a silent digest entered the conversation:\n%s", idle)
	}
	busy := ansi.Strip(busyModel(100, 40).View())
	for _, want := range []string{"đụng vào▌", "✓ Bash: mate task new e1-profile --from-chat", "Bash: mate crew spawn e1-profile", "Xếp hàng (2)", `1 › à nhớ đổi tên`, "Mate đang làm · Bash: mate crew spawn e1-profile · 7.0s", "^C ngắt", "Supabase Storage hay S3?"} {
		if !strings.Contains(busy, want) {
			t.Errorf("busy screen lacks %q:\n%s", want, busy)
		}
	}
}

func TestREPLComposerKeys(t *testing.T) {
	m := NewModel(New(&Runner{}, Options{}), Config{Project: "p"}, nil)
	m.host.queue = []message{{Source: SourceCaptain, Text: "first", Prompt: "first"}, {Source: SourceCaptain, Text: "second", Prompt: "second"}}
	press := func(k tea.KeyMsg) {
		next, _ := m.key(k)
		m = next.(Model)
	}
	for _, r := range "xin chào" {
		press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	press(tea.KeyMsg{Type: tea.KeyCtrlJ})
	for _, r := range "dòng hai" {
		press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	press(tea.KeyMsg{Type: tea.KeyCtrlW})
	press(tea.KeyMsg{Type: tea.KeyUp})
	press(tea.KeyMsg{Type: tea.KeyCtrlK})
	if got := string(m.input); got != "xin c\ndòng " {
		t.Fatalf("input = %q", got)
	}
	if m.host.State().Draft != "xin c\ndòng " {
		t.Fatalf("draft = %q", m.host.State().Draft)
	}
	press(tea.KeyMsg{Type: tea.KeyBackspace, Alt: true})
	if got := string(m.input); got != "xin \ndòng " {
		t.Fatalf("⌥⌫ left %q", got)
	}
	press(tea.KeyMsg{Type: tea.KeyCtrlU})
	if len(m.input) != 0 || m.host.State().Draft != "" {
		t.Fatalf("^U left %q", string(m.input))
	}
	for _, r := range "lại" {
		press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	press(tea.KeyMsg{Type: tea.KeyCtrlC})
	if len(m.input) != 0 {
		t.Fatalf("^C left %q", string(m.input))
	}
	press(tea.KeyMsg{Type: tea.KeyUp})
	if got := string(m.input); got != "second" || len(m.host.queue) != 1 {
		t.Fatalf("↑ took %q, queue %q", got, m.host.queue)
	}
	press(tea.KeyMsg{Type: tea.KeyCtrlC})
	press(tea.KeyMsg{Type: tea.KeyCtrlX})
	if len(m.host.queue) != 0 {
		t.Fatalf("^X left %q", m.host.queue)
	}
}
