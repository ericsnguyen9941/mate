package matehost

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func TestNewsLifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), NewsName)
	now := time.Date(2026, 10, 15, 14, 9, 0, 0, time.Local)
	id, err := OpenNews(path, News{Title: "Dark mode", Crew: "e1", Options: []NewsOption{{Key: "a", Text: "env"}, {Key: "b", Text: "config"}}}, now)
	if err != nil || id != 1 {
		t.Fatalf("open = %d, %v", id, err)
	}
	if id, _ := OpenNews(path, News{Title: "PR #12"}, now); id != 2 {
		t.Fatalf("second id = %d", id)
	}
	if _, err := OpenNews(path, News{}, now); err == nil {
		t.Fatal("a news without a title opened")
	}
	status := func(id int) string {
		list, err := LoadNews(path)
		if err != nil {
			t.Fatal(err)
		}
		n, _ := findNews(list, id)
		return n.Status
	}
	steps := []struct {
		op   newsEvent
		want string
	}{
		{newsEvent{Op: opSeen, ID: 1}, NewsWaiting},
		{newsEvent{Op: opUpdate, ID: 1, Text: "CI xanh"}, NewsWaiting},
		{newsEvent{Op: opAnswer, ID: 1, Text: "[a] env"}, NewsWorking},
		{newsEvent{Op: opFailed, ID: 1}, NewsWaiting},
		{newsEvent{Op: opAnswer, ID: 1, Text: "[a] env"}, NewsWorking},
		{newsEvent{Op: opAsk, ID: 1, Text: "hỏi lại"}, NewsNew},
		{newsEvent{Op: opClose, ID: 1, Text: "PR đã merge"}, NewsDone},
	}
	for i, s := range steps {
		s.op.At = now.Add(time.Duration(i+1) * time.Minute)
		if err := newsOp(path, s.op, true); err != nil {
			t.Fatal(err)
		}
		if got := status(1); got != s.want {
			t.Fatalf("after %s: status %q, want %q", s.op.Op, got, s.want)
		}
	}
	if err := UpdateNews(path, 1, "late", false, now); err == nil {
		t.Error("a closed news took an update")
	}
	if err := CloseNews(path, 2, "PR closed", true, now); err != nil || status(2) != NewsDropped {
		t.Fatalf("drop = %v, %q", err, status(2))
	}
	if err := CloseNews(path, 9, "", false, now); err == nil {
		t.Error("closed a news that does not exist")
	}
	list, _ := LoadNews(path)
	if n := list[0]; n.Answer != "[a] env" || len(n.Notes) != 3 || n.Notes[2].Text != "PR đã merge" {
		t.Fatalf("news 1 = %+v", n)
	}
}

// The captain's answer is a turn that names the news, belongs to it, and
// moves it to working; an app turn that opens a news belongs to it.
func TestHostAnswersAndOpensNews(t *testing.T) {
	f := newFixture(t)
	path := filepath.Join(f.dir, NewsName)
	if _, err := OpenNews(path, News{Title: "Dark mode"}, time.Now().Add(-time.Minute)); err != nil {
		t.Fatal(err)
	}
	f.host.RefreshNews()
	if err := f.host.AnswerNews(1, "[a] env"); err != nil {
		t.Fatal(err)
	}
	items := f.done(t, 1)
	if it := items[0]; it.News != 1 || it.Text != "[a] env" || it.Source != SourceCaptain {
		t.Fatalf("answer item = %+v", it)
	}
	if prompt := f.read(t, "prompts.log"); !strings.Contains(prompt, "[news #1 · Dark mode] [a] env") {
		t.Fatalf("prompt = %q", prompt)
	}
	list, _ := LoadNews(path)
	if list[0].Status != NewsWorking {
		t.Fatalf("status = %q", list[0].Status)
	}
	if err := f.host.AnswerNews(7, "x"); err == nil {
		t.Error("answered a news that does not exist")
	}

	if err := f.host.Deliver(Marker + "digest: slow"); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenNews(path, News{Title: "PR #12"}, time.Now()); err != nil {
		t.Fatal(err)
	}
	items = f.done(t, 2)
	if it := items[1]; it.Source != SourceApp || it.News != 2 {
		t.Fatalf("app item = %+v", it)
	}
}

func newsModel(w, h int) Model {
	m := sampleModel(w, h)
	at := m.now().Add(-4 * time.Minute)
	m.news = []News{
		{ID: 1, Title: "Chọn nhà cung cấp đăng nhập: Google", Crew: "e1-login-resume", Status: NewsDone, OpenedAt: at.Add(-time.Hour), At: at},
		{ID: 3, Title: "Dark mode không áp dụng trên bản clone", Crew: "e1-settings-resume", Kind: "needs-decision", Status: NewsWaiting, OpenedAt: at, At: at,
			Context:  []string{`Crew báo: "dark mode never applies on my clone…"`, "Theme đọc từ biến THEME trong .env.local. File này không có trong git, nên bản clone mới luôn rơi về light."},
			Question: "Sửa theo cách nào cũng nhỏ, nhưng sẽ đổi cách cả team cấu hình máy dev, nên em hỏi anh:",
			Options:  []NewsOption{{Key: "a", Text: "Thêm THEME=system vào .env.example và fallback theo giao diện hệ thống", Recommended: true}, {Key: "b", Text: "Bỏ biến môi trường, chuyển cấu hình theme vào config/theme.ts"}, {Key: "c", Text: "Chỉ ghi chú vào README, không sửa code"}}},
		{ID: 4, Title: "PR #12 sẵn sàng review: thanh toán MoMo", Crew: "e2-payments", Status: NewsNew, OpenedAt: m.now().Add(-time.Minute), At: m.now().Add(-time.Minute)},
	}
	return m
}

func TestREPLNewsScreens(t *testing.T) {
	for _, w := range []int{80, 100} {
		m := newsModel(w, 30)
		m.host = New(&Runner{}, Options{})
		chat := ansi.Strip(m.View())
		for _, want := range []string{"● 2 chờ anh", "● mới  #4 e2-payments: PR #12 sẵn sàng review", "+1 · ^N news"} {
			if !strings.Contains(chat, want) {
				t.Errorf("chat@%d lacks %q:\n%s", w, want, chat)
			}
		}
		press := func(k tea.KeyMsg) {
			next, _ := m.key(k)
			m = next.(Model)
		}
		press(tea.KeyMsg{Type: tea.KeyCtrlT})
		list := ansi.Strip(m.View())
		for _, want := range []string{"News", "2 chờ anh · 0 Mate xử lý · 1 xong", "CHỜ ANH", "▸ ● #4", "XONG HÔM NAY", "✓ #1"} {
			if !strings.Contains(list, want) {
				t.Errorf("list@%d lacks %q:\n%s", w, want, list)
			}
		}
		wheel := func(b tea.MouseButton) {
			next, _ := m.Update(tea.MouseMsg{Button: b, Action: tea.MouseActionPress})
			m = next.(Model)
		}
		wheel(tea.MouseButtonWheelDown)
		if m.listSel != 1 || m.scroll != 0 {
			t.Fatalf("wheel down in the list: sel %d, chat scroll %d", m.listSel, m.scroll)
		}
		wheel(tea.MouseButtonWheelUp)
		wheel(tea.MouseButtonWheelUp)
		if m.listSel != 0 {
			t.Fatalf("wheel up in the list: sel %d", m.listSel)
		}
		press(tea.KeyMsg{Type: tea.KeyDown})
		press(tea.KeyMsg{Type: tea.KeyEnter})
		// A short pane, so the news does not fit and scrolls.
		m.height = 16
		wheel(tea.MouseButtonWheelUp)
		if m.newsScroll != wheelLines || m.scroll != 0 || len(m.input) != 0 {
			t.Fatalf("wheel up in a news: news scroll %d, chat scroll %d, input %q", m.newsScroll, m.scroll, string(m.input))
		}
		for range 50 {
			wheel(tea.MouseButtonWheelUp)
		}
		if m.newsScroll != m.maxScroll() {
			t.Fatalf("news scrolled to %d past its top %d", m.newsScroll, m.maxScroll())
		}
		m.height, m.newsScroll = 30, 0
		news := ansi.Strip(m.View())
		for _, want := range []string{"News #3", "CHỜ ANH · needs-decision · crew e1-settings-resume", "Ngữ cảnh", "│ Crew báo", "▸ [a] Thêm THEME=system", "khuyên dùng", "[c] Chỉ ghi chú", "Trả lời #3", "a/b/c chọn ngay"} {
			if !strings.Contains(news, want) {
				t.Errorf("news@%d lacks %q:\n%s", w, want, news)
			}
		}
		for i, l := range strings.Split(m.View(), "\n") {
			if ansi.StringWidth(l) > w {
				t.Errorf("news@%d line %d too wide: %q", w, i, ansi.Strip(l))
			}
		}
		press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")})
		if p := m.pending; p == nil || p.id != 3 || p.text != "[b] Bỏ biến môi trường, chuyển cấu hình theme vào config/theme.ts" {
			t.Fatalf("pending = %+v", m.pending)
		}
		if v := ansi.Strip(m.View()); !strings.Contains(v, "✓ Đã chọn B · u hoàn tác (5s)") {
			t.Errorf("no undo line:\n%s", v)
		}
		press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("u")})
		if m.pending != nil {
			t.Fatal("u did not undo")
		}
		press(tea.KeyMsg{Type: tea.KeyTab})
		for _, r := range "làm a" {
			press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		}
		press(tea.KeyMsg{Type: tea.KeyEsc})
		if m.mode != viewChat || len(m.input) != 0 {
			t.Fatalf("Esc left mode %d input %q", m.mode, string(m.input))
		}
		press(tea.KeyMsg{Type: tea.KeyCtrlN})
		if m.mode != viewNews || m.newsID != 3 || string(m.input) != "làm a" {
			t.Fatalf("^N opened #%d with draft %q", m.newsID, string(m.input))
		}
		if w == 100 && testing.Verbose() {
			t.Logf("\n%s\n\n%s\n\n%s", chat, list, news)
		}
	}
}

func TestREPLChatNewsReference(t *testing.T) {
	m := newsModel(100, 30)
	m.host = New(&Runner{}, Options{})
	for _, r := range "#3 chọn a" {
		next, _ := m.key(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = next.(Model)
	}
	next, _ := m.key(tea.KeyMsg{Type: tea.KeyEnter})
	m = next.(Model)
	if p := m.pending; p == nil || p.id != 3 || !strings.HasPrefix(p.text, "[a] Thêm THEME") {
		t.Fatalf("pending = %+v", m.pending)
	}
	if v := ansi.Strip(m.View()); !strings.Contains(v, "✓ Đã ghi vào #3 · u hoàn tác") {
		t.Errorf("no confirmation in the tray:\n%s", v)
	}
}
