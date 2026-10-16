package matehost

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// News in the REPL (the design's topics): a badge in the header and a tray
// over the composer while the captain owes an answer, ^T for the list, ^N
// for the next one owed, and a screen per news where the captain reads it,
// picks an option or writes an answer, and Esc goes back to the chat where
// it was. An answer waits undoDelay before it is sent, so `u` can take it
// back.

const (
	undoDelay   = 5 * time.Second
	confirmFor  = 5 * time.Second
	listMaxRows = 12
	// unclosedAfter is how long a news may be worked after the captain's
	// answer before the list says the Mate has not closed it: the captain
	// owes nothing, the Mate likely forgot `mate news close`.
	unclosedAfter = 30 * time.Minute
)

var (
	colAttention   = lipgloss.AdaptiveColor{Light: "130", Dark: "214"}
	styleAttention = lipgloss.NewStyle().Foreground(colAttention)
	styleSelected  = lipgloss.NewStyle().Reverse(true)
)

type viewMode int

const (
	viewChat viewMode = iota
	viewList
	viewNews
)

// pendingAnswer is an answer waiting out its undo delay.
type pendingAnswer struct {
	id    int
	text  string
	label string
	at    time.Time
}

// newsRef is a captain message in the chat that answers a news: `#3 …`.
var newsRef = regexp.MustCompile(`^#(\d+)\s+(?s)(.+)$`)

func (m Model) findNews(id int) (News, bool) { return findNews(m.news, id) }

// owed are the news waiting on the captain, oldest first.
func (m Model) owed() []News {
	var out []News
	for _, n := range m.news {
		if n.Owed() {
			out = append(out, n)
		}
	}
	return out
}

// listed are the news of the list's three groups: owed, being handled,
// closed today.
func (m Model) listed() (owed, working, closed []News) {
	today := m.now().Format("2006-01-02")
	for _, n := range m.news {
		switch {
		case n.Owed():
			owed = append(owed, n)
		case n.Status == NewsWorking:
			working = append(working, n)
		case n.Closed() && !n.Archived && n.At.Format("2006-01-02") == today:
			closed = append(closed, n)
		}
	}
	newest := func(l []News) {
		sort.SliceStable(l, func(i, j int) bool { return l[i].At.After(l[j].At) })
	}
	newest(owed)
	newest(working)
	newest(closed)
	return
}

// listFlat are the rows the list selects: the news closed today only
// while Tab has unfolded them.
func (m Model) listFlat() []News {
	a, b, c := m.listed()
	out := append(append([]News(nil), a...), b...)
	if m.listClosed {
		out = append(out, c...)
	}
	return out
}

// --- moving between views ---

// openNews shows news id, keeping the chat's draft and scroll to come back
// to.
func (m *Model) openNews(id int) {
	n, ok := m.findNews(id)
	if !ok {
		return
	}
	m.stash()
	m.mode, m.newsID, m.newsScroll = viewNews, id, 0
	m.restore()
	m.optSel = 0
	for i, o := range n.Options {
		if o.Recommended {
			m.optSel = i
		}
	}
	m.focusOptions = len(n.Options) > 0 && n.Owed() && len(m.input) == 0
	if n.Status == NewsNew {
		go m.host.SeenNews(id)
	}
}

func (m *Model) closeView() {
	m.stash()
	m.mode, m.newsID = viewChat, 0
	m.restore()
}

// stash and restore keep a draft per view: the chat's, and each news's.
func (m *Model) stash() {
	if m.drafts == nil {
		m.drafts = map[int]string{}
	}
	key := 0
	if m.mode == viewNews {
		key = m.newsID
	}
	m.drafts[key] = string(m.input)
}

func (m *Model) restore() {
	key := 0
	if m.mode == viewNews {
		key = m.newsID
	}
	r := []rune(m.drafts[key])
	m.input, m.cursor = r, len(r)
}

// nextOwed opens the owed news after the one open now, oldest first, and
// wraps around.
func (m *Model) nextOwed() {
	owed := m.owed()
	if len(owed) == 0 {
		m.note = "không có news nào chờ anh"
		return
	}
	next := owed[0]
	if m.mode == viewNews {
		for i, n := range owed {
			if n.ID == m.newsID {
				next = owed[(i+1)%len(owed)]
				break
			}
		}
	}
	m.openNews(next.ID)
}

// --- answering ---

// answer queues an answer to news id behind the undo delay.
func (m *Model) answer(id int, text, label string) {
	if p := m.pending; p != nil {
		m.flush()
	}
	m.pending = &pendingAnswer{id: id, text: text, label: label, at: m.now()}
}

// flush sends the pending answer.
func (m *Model) flush() {
	p := m.pending
	if p == nil {
		return
	}
	m.pending = nil
	m.confirm, m.confirmAt = fmt.Sprintf("✓ #%d đã trả lời · Mate đang xử lý", p.id), m.now()
	host := m.host
	go func() {
		if err := host.AnswerNews(p.id, p.text); err != nil {
			host.emitNote(err)
		}
	}()
}

// tickNews sends an answer whose undo delay is over.
func (m *Model) tickNews() {
	if p := m.pending; p != nil && m.now().Sub(p.at) >= undoDelay {
		m.flush()
	}
}

func (m *Model) choose(n News, i int) {
	o := n.Options[i]
	m.answer(n.ID, fmt.Sprintf("[%s] %s", o.Key, o.Text), "Đã chọn "+strings.ToUpper(o.Key))
}

// --- keys ---

// newsKey handles a key in the list or a news; handled is false for keys
// the composer takes.
func (m Model) newsKey(k tea.KeyMsg) (Model, bool) {
	switch k.Type {
	case tea.KeyCtrlT:
		if m.mode == viewList {
			m.mode = m.listFrom
		} else {
			m.listFrom, m.mode, m.listSel = m.mode, viewList, 0
		}
		return m, true
	case tea.KeyCtrlN:
		m.nextOwed()
		return m, true
	}
	if p := m.pending; p != nil && k.Type == tea.KeyRunes && string(k.Runes) == "u" && (len(m.input) == 0 || m.focusOptions) {
		m.pending = nil
		m.note = fmt.Sprintf("đã hoàn tác câu trả lời cho #%d", p.id)
		return m, true
	}
	switch m.mode {
	case viewList:
		return m.listKey(k), true
	case viewNews:
		return m.newsViewKey(k)
	}
	return m, false
}

func (m Model) listKey(k tea.KeyMsg) Model {
	flat := m.listFlat()
	switch k.Type {
	case tea.KeyEsc:
		m.mode = m.listFrom
	case tea.KeyTab:
		m.listClosed = !m.listClosed
		m.listSel = min(m.listSel, max(0, len(m.listFlat())-1))
	case tea.KeyUp:
		m.listSel = max(0, m.listSel-1)
	case tea.KeyDown:
		m.listSel = min(max(0, len(flat)-1), m.listSel+1)
	case tea.KeyEnter:
		if m.listSel < len(flat) {
			m.mode = m.listFrom
			m.openNews(flat[m.listSel].ID)
		}
	case tea.KeyRunes:
		switch r := string(k.Runes); {
		case r == "j":
			m.listSel = min(max(0, len(flat)-1), m.listSel+1)
		case r == "k":
			m.listSel = max(0, m.listSel-1)
		case r == "x":
			if m.listSel < len(flat) && flat[m.listSel].Closed() {
				go m.host.ArchiveNews(flat[m.listSel].ID)
			}
		case len(r) == 1 && r >= "1" && r <= "9":
			if i, _ := strconv.Atoi(r); i <= len(flat) {
				m.mode = m.listFrom
				m.openNews(flat[i-1].ID)
			}
		}
	}
	return m
}

func (m Model) newsViewKey(k tea.KeyMsg) (Model, bool) {
	n, ok := m.findNews(m.newsID)
	if !ok {
		m.closeView()
		return m, true
	}
	switch k.Type {
	case tea.KeyEsc:
		m.closeView()
		return m, true
	case tea.KeyPgUp:
		m.scrollView(max(1, m.viewRows()-2))
		return m, true
	case tea.KeyPgDown:
		m.scrollView(-max(1, m.viewRows()-2))
		return m, true
	}
	if n.Closed() {
		// A closed news is read, not answered: no composer takes the key.
		switch k.Type {
		case tea.KeyUp:
			m.scrollView(1)
		case tea.KeyDown:
			m.scrollView(-1)
		}
		return m, true
	}
	switch k.Type {
	case tea.KeyTab:
		if m.focusOptions {
			m.focusOptions = false
		} else if len(n.Options) > 0 && len(m.input) == 0 {
			m.focusOptions = true
		}
		return m, true
	case tea.KeyEnter:
		if k.Alt {
			return m, false
		}
		if m.focusOptions {
			m.choose(n, m.optSel)
			return m, true
		}
		text := strings.TrimSpace(string(m.input))
		if text == "" {
			return m, true
		}
		m.setInput(nil)
		m.answer(n.ID, text, "Đã gửi")
		return m, true
	}
	if !m.focusOptions {
		return m, false
	}
	switch k.Type {
	case tea.KeyUp:
		m.optSel = max(0, m.optSel-1)
	case tea.KeyDown:
		m.optSel = min(len(n.Options)-1, m.optSel+1)
	case tea.KeyRunes:
		key := strings.ToLower(string(k.Runes))
		for i, o := range n.Options {
			if strings.ToLower(o.Key) == key {
				m.optSel = i
				m.choose(n, i)
				return m, true
			}
		}
		// Any other letter starts a written answer.
		m.focusOptions = false
		return m, false
	}
	return m, true
}

// chatNewsRef sends a chat message that names a news (`#3 chọn a`) to that
// news instead, behind the same undo delay.
func (m *Model) chatNewsRef(text string) bool {
	g := newsRef.FindStringSubmatch(text)
	if g == nil {
		return false
	}
	id, _ := strconv.Atoi(g[1])
	n, ok := m.findNews(id)
	if !ok {
		return false
	}
	if n.Closed() {
		m.note = fmt.Sprintf("#%d đã xong nên tin này gửi cho Mate như tin thường", id)
		return false
	}
	answer := strings.TrimSpace(g[2])
	for _, o := range n.Options {
		if strings.EqualFold(answer, o.Key) || strings.EqualFold(answer, "chọn "+o.Key) {
			answer = fmt.Sprintf("[%s] %s", o.Key, o.Text)
		}
	}
	m.answer(id, answer, fmt.Sprintf("Đã ghi vào #%d", id))
	return true
}

// --- drawing ---

// newsStatus is how a news reads at now. A news worked unclosedAfter past
// the captain's answer reads as the Mate not having closed it.
func newsStatus(n News, now time.Time) (glyph, label string, style lipgloss.Style) {
	switch n.Status {
	case NewsNew:
		return "●", "mới", styleAttention.Bold(true)
	case NewsWaiting:
		return "●", "chờ anh", styleAttention
	case NewsWorking:
		if now.Sub(n.At) >= unclosedAfter {
			return "⠋", "Mate chưa đóng", styleAttention
		}
		return "⠋", "Mate xử lý", styleWorking
	case NewsDropped:
		return "✓", "tự đóng", styleMuted
	default:
		return "✓", "xong", styleSuccess
	}
}

func age(d time.Duration) string {
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%dp", max(1, int(d.Minutes())))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dg", int(d.Hours()))
	default:
		return fmt.Sprintf("%dn", int(d.Hours()/24))
	}
}

// badge is the header's count of news owed.
func (m Model) badge() string {
	owed := m.owed()
	if len(owed) == 0 {
		return ""
	}
	style := styleAttention
	for _, n := range owed {
		if n.Status == NewsNew {
			style = style.Bold(true)
		}
	}
	return style.Render(fmt.Sprintf("● %d chờ anh", len(owed)))
}

// tray is the line over the composer in the chat: an answer just sent, and
// the newest news owed.
func (m Model) tray(w int) []string {
	var out []string
	if p := m.pending; p != nil {
		left := max(0, int((undoDelay-m.now().Sub(p.at)+time.Second-1)/time.Second))
		out = append(out, styleSuccess.Render(fmt.Sprintf(" ✓ %s · u hoàn tác (%ds)", p.label, left)))
	} else if m.confirm != "" && m.now().Sub(m.confirmAt) < confirmFor {
		out = append(out, styleSuccess.Render(" "+m.confirm))
	}
	owed := m.owed()
	if len(owed) == 0 {
		return out
	}
	n := owed[0]
	for _, x := range owed {
		if x.At.After(n.At) {
			n = x
		}
	}
	glyph, label, style := newsStatus(n, m.now())
	text := fmt.Sprintf(" %s %s  #%d ", glyph, label, n.ID)
	if n.Crew != "" {
		text += n.Crew + ": "
	}
	text += n.Title
	more := ""
	if len(owed) > 1 {
		more = fmt.Sprintf("+%d · ", len(owed)-1)
	}
	right := styleMuted.Render(more + "^N news ")
	text = ansi.Truncate(text, max(4, w-ansi.StringWidth(right)-2), "…")
	return append(out, spread(style.Render(text), right, w))
}

// listBox is the ^T list, drawn over the bottom of the conversation.
func (m Model) listBox(w int) []string {
	owed, working, closed := m.listed()
	inner := w - 4
	row := func(n News, i int) string {
		glyph, label, style := newsStatus(n, m.now())
		mark := "  "
		if i == m.listSel {
			mark = "▸ "
		}
		ageText := age(m.now().Sub(n.At))
		right := fmt.Sprintf("%-14s%4s", label, ageText)
		title := n.Title
		if n.Crew != "" && w >= 100 {
			title = fmt.Sprintf("%-20s %s", ansi.Truncate(n.Crew, 20, "…"), title)
		}
		left := fmt.Sprintf("%s%s #%d %s", mark, style.Render(glyph), n.ID, title)
		left = ansi.Truncate(left, max(4, inner-ansi.StringWidth(right)-2), "…")
		line := spread(left, style.Render(right), inner)
		if i == m.listSel {
			line = styleSelected.Render(ansi.Strip(line) + strings.Repeat(" ", max(0, inner-ansi.StringWidth(line))))
		}
		return line
	}
	var body []string
	i := 0
	group := func(title string, list []News) {
		if len(list) == 0 {
			return
		}
		if len(body) > 0 {
			body = append(body, "")
		}
		body = append(body, styleMuted.Render(title))
		for _, n := range list {
			body = append(body, row(n, i))
			i++
		}
	}
	group("CHỜ ANH", owed)
	group("MATE ĐANG XỬ LÝ", working)
	if m.listClosed {
		group("XONG HÔM NAY", closed)
	} else if len(closed) > 0 {
		if len(body) > 0 {
			body = append(body, "")
		}
		body = append(body, styleMuted.Render(fmt.Sprintf("  ✓ %d xong hôm nay · Tab xem", len(closed))))
	}
	if len(body) == 0 {
		body = []string{styleMuted.Render("Chưa có news nào. Khi crew báo việc cần anh, Mate mở news ở đây.")}
	}
	if len(body) > listMaxRows {
		body = body[:listMaxRows]
	}
	summary := fmt.Sprintf(" %d chờ anh · %d Mate xử lý · %d xong ", len(owed), len(working), len(closed))
	top := styleBorder.Render("┌─ ") + styleBold.Render("News") + " " + styleBorder.Render(strings.Repeat("─", max(0, w-8-ansi.StringWidth(summary)-2))) + styleMuted.Render(summary) + styleBorder.Render("─┐")
	out := []string{top}
	for _, l := range body {
		out = append(out, styleBorder.Render("│ ")+l+strings.Repeat(" ", max(0, inner-ansi.StringWidth(l)))+styleBorder.Render(" │"))
	}
	hint := " ↑↓ chọn · Enter mở · 1–9 nhảy · Tab news xong · Esc về "
	if m.listClosed {
		hint = " ↑↓ chọn · Enter mở · Tab ẩn news xong · x lưu trữ · Esc về "
	}
	return append(out, styleBorder.Render("└─")+styleMuted.Render(hint)+styleBorder.Render(strings.Repeat("─", max(0, w-3-ansi.StringWidth(hint)))+"┘"))
}

// newsLines are the screen of one news: what it is, its context, the
// Mate's question and options, its updates, and the exchange about it.
func (m Model) newsLines(w int) []string {
	n, ok := m.findNews(m.newsID)
	if !ok {
		return nil
	}
	glyph, label, style := newsStatus(n, m.now())
	if n.Status == NewsWorking {
		glyph = spinnerFrames[m.frame%len(spinnerFrames)]
	}
	head := []string{strings.ToUpper(label)}
	if n.Kind != "" {
		head = append(head, n.Kind)
	}
	if n.Crew != "" {
		head = append(head, "crew "+n.Crew)
	}
	head = append(head, n.OpenedAt.Format("15:04"))
	out := []string{"", ansi.Truncate(style.Render(" "+glyph+" "+strings.Join(head, " · ")), w-1, "…"), ""}
	for _, l := range wrap(styleBold.Render(n.Title), max(10, w-4)) {
		out = append(out, "   "+l)
	}
	if len(n.Context) > 0 {
		out = append(out, "", " "+styleMuted.Render("Ngữ cảnh"))
		for _, c := range n.Context {
			for _, l := range wrap(c, max(10, w-6)) {
				out = append(out, styleBorder.Render("  │ ")+l)
			}
		}
	}
	if n.Question != "" {
		out = append(out, "")
		out = append(out, replyLines(n.Question, w)...)
	}
	if len(n.Options) > 0 {
		out = append(out, "")
		for i, o := range n.Options {
			out = append(out, m.optionLines(n, i, o, w)...)
		}
	}
	for _, note := range n.Notes {
		out = append(out, "", styleMuted.Render(ansi.Truncate(fmt.Sprintf(" · %s  %s", note.At.Format("15:04"), note.Text), w-1, "…")))
	}
	for i, it := range m.items {
		if it.News != m.newsID {
			continue
		}
		out = append(out, "")
		out = append(out, m.threadLines(it, w, i == len(m.items)-1 && m.running == nil)...)
	}
	if it := m.running; it != nil && it.News == m.newsID {
		out = append(out, "")
		out = append(out, m.threadLines(*it, w, true)...)
	}
	if n.Closed() {
		out = append(out, "", style.Render(fmt.Sprintf(" %s %s · %s", glyph, strings.ToUpper(label[:1])+label[1:], n.At.Format("15:04"))))
	}
	return out
}

// threadLines are one turn of a news's exchange. The turn that opened the
// news shows only the Mate's reply, when it had one.
func (m Model) threadLines(it Item, w int, last bool) []string {
	if it.Source == SourceApp {
		if !it.Done {
			return m.liveLines(it, w)
		}
		if it.Silent || it.Reply == "" {
			return nil
		}
		return append(replyLines(it.Reply, w), metaLine(it, w, false))
	}
	return m.itemLines(it, w, last)
}

func (m Model) optionLines(n News, i int, o NewsOption, w int) []string {
	chosen := strings.HasPrefix(n.Answer, "["+o.Key+"]")
	mark := "  "
	if m.focusOptions && i == m.optSel {
		mark = styleAccent.Render("▸ ")
	} else if chosen {
		mark = styleSuccess.Render("✓ ")
	}
	right := ""
	if o.Recommended {
		right = "khuyên dùng"
	}
	key := styleAccent.Render("[" + o.Key + "]")
	lines := wrap(o.Text, max(10, w-8-len(right)-2))
	out := make([]string, len(lines))
	for j, l := range lines {
		if j == 0 {
			out[j] = spread(" "+mark+key+" "+l, styleMuted.Render(right)+" ", w)
		} else {
			out[j] = "       " + l
		}
	}
	return out
}
