# mate

mate điều phối một đội agent lập trình trên máy của bạn. Bạn là **captain**: bạn
nói chuyện với một **Mate** cho mỗi project. Mate chia việc, giao cho các **Crew**
chạy song song, mỗi Crew một git worktree riêng. Mate cho review độc lập, rồi
merge, hoặc để Crew mở PR. mate là một binary Go: core của workspace, console
trong terminal, dashboard web và CLI cho agent gọi.

```
captain ──chat──▶ Mate (mỗi project một Mate, chạy Claude Code)
                    ├─ mate crew spawn ──▶ Crew A (worktree, nhánh riêng)
                    ├─ mate crew spawn ──▶ Crew B
                    ├─ mate review ──────▶ Crew reviewer (scout, chỉ đọc)
                    └─ mate land (merge local), hoặc Crew mở PR
```

## Nền tảng hỗ trợ

| Nền tảng | Tình trạng |
| --- | --- |
| macOS Apple silicon | Chạy hằng ngày, đủ tính năng. |
| Linux x86_64 / arm64 | Build và test được (`make check-linux`). Chưa ai chạy thật hằng ngày. |
| Windows qua WSL2 | Cài mọi thứ trong WSL2 như trên Linux; workspace dưới `~`, không dưới `/mnt/c`; mở Office ở `http://localhost:7777` từ trình duyệt Windows. Chưa kiểm chứng ([docs/install.md](docs/install.md)). |
| Windows native | Không có. Code dùng syscall chỉ có trên Unix. |

## Cần cài trước

| Thứ | Để làm gì | Bắt buộc |
| --- | --- | --- |
| Go 1.25+ | Cài hoặc build mate | Có |
| git | Repo, worktree, merge | Có |
| [Herdr](https://herdr.dev) | Chạy Crew trong pane (agent multiplexer, giống tmux) | Có |
| `claude` (Claude Code), đã đăng nhập | Mate chỉ chạy trên Claude Code | Có |
| `codex`, đã đăng nhập | Harness mặc định cho Crew | Có, trừ khi đổi `crew_harness` |
| `pi`, `grok`, đã đăng nhập | Harness khác cho Crew | Không |
| `gh`, đã `gh auth login` | Crew giao bằng PR | Khi dùng PR |
| [WezTerm](https://wezterm.org), hoặc [Ghostty](https://ghostty.org) 1.3+ trên macOS | Console mở pane của agent bên cạnh | Nên có |

Mặc định Mate chạy `claude`, Crew chạy `codex`. Crew chạy được bằng `claude`,
`codex`, `pi`, `grok`. Mỗi harness tự đăng nhập bằng tài khoản của nó; mate không
chép token giữa các harness. Không có `codex` thì đặt
`defaults.crew_harness: claude` trong `.mate/workspace.yaml`. Trên macOS:
`brew install go git gh herdr`.

## Cài mate

```sh
go install github.com/nguyenngocanh94/mate/cmd/mate@latest
mate --version
```

Module path còn mang tên chủ cũ của repo; GitHub tự chuyển sang
`github.com/ericsnguyen9941/mate`, nên lệnh trên chạy đúng như đã viết. `go install`
đặt `mate` vào `$GOBIN`, hoặc `$(go env GOPATH)/bin` (mặc định `~/go/bin`); thư mục
đó phải nằm trong `PATH`. Hoặc build từ mã nguồn:

```sh
git clone https://github.com/ericsnguyen9941/mate.git
cd mate
make build      # ra bin/mate; chép nó vào một thư mục trong PATH
```

Sau khi cài bản mới, core đang chạy vẫn dùng binary cũ: `mate core status` báo
`stale`, chạy `mate core restart`. Chi tiết: [docs/install.md](docs/install.md).

## Bắt đầu nhanh

### 1. Tạo workspace

Workspace là một thư mục chứa mọi project. State của mate nằm trong `.mate/` ở gốc.

```sh
mate init ~/workspace
```

### 2. Thêm project và repo

Mỗi project là một thư mục trong workspace, repo nằm trong thư mục project
(`~/workspace/shop/api`). Cách dễ nhất là để mate clone:

```sh
cd ~/workspace
mate project add shop
mate project repo add shop https://github.com/you/shop-api.git --name api
```

Repo đã có sẵn ở `~/workspace/<project>/<repo>` thì thêm thẳng:
`mate project add blog ~/workspace/blog/blog`. Nhánh mặc định được tự đoán; ghi đè
bằng `--default-branch`.

### 3. Chạy core

```sh
cd ~/workspace
mate            # hoặc: mate . / mate ~/workspace
```

`mate` chạy **mate core** ở foreground, như một dev server. Core là engine của
workspace: theo dõi agent, giao tin cho Mate, phục hồi sau sự cố, gửi notification,
phục vụ Office và Agent City. Mỗi workspace một core. Nó in:

```
mate core · pid 4121 · ~/workspace
Office  http://127.0.0.1:7777/office/
City    http://127.0.0.1:7777/city/
Herdr   mate-7b09de89
Log     /Users/you/workspace/.mate/core.log
Ctrl+C to stop
```

Core không tự mở trình duyệt. Core đang chạy rồi thì `mate` in status và thoát 0.
Muốn chạy nền: `mate core start|stop|restart|status`. Chi tiết:
[docs/core.md](docs/core.md).

### 4. Start Mate và giao việc

```sh
mate mate-agent start shop           # hoặc trong console: chọn project, nhấn s
mate mate-agent start shop --fresh   # bỏ phiên cũ, bắt đầu phiên mới
```

Mate mặc định chạy **headless**: host của Mate là một process nền ngoài Herdr, mỗi
tin nhắn là một lượt `claude -p`. Core giám sát host: host chết thì core bật lại,
tối đa 3 lần trong 10 phút; quá số đó core báo bạn và để yên.

Rồi nói với Mate như với một tech lead: "làm trang checkout theo issue SHOP-12",
"sửa bug login trên Safari, mở PR vào main". Mate sẽ:

1. viết brief và spawn Crew (`mate crew spawn`), chọn harness, model, effort theo độ khó;
2. theo dõi Crew, tự trả lời những câu nó trả lời được;
3. cho một Crew khác review độc lập (`mate review`) nếu project bật review;
4. nhận bàn giao rồi merge local (`mate land`); Crew giao bằng PR thì Mate theo dõi PR;
5. hết việc bạn giao thì tự lấy task sẵn sàng tiếp theo của tracker
   (`next_work: pick`, mặc định), hoặc hỏi bạn (`next_work: ask`).

Khi cần bạn quyết định, Mate mở một **news** (một câu hỏi) và máy bắn notification.

### 5. Nói chuyện với Mate và xem Crew

| Ở đâu | Cách mở |
| --- | --- |
| Console | `mate console`. Là client của core: mọi hành động đi qua core. Chưa có core thì console tự bật. |
| REPL của Mate | `mate mate-agent repl shop`, trong terminal bất kỳ; console hiện REPL này ở pane bên cạnh. `Ctrl+T`: topic (news) để đọc và trả lời; `Ctrl+L`: bảng Work. |
| Office, Agent City | Trình duyệt, ở `/office/` và `/city/`. Chat, trả lời news, xem bàn giao. |
| Pane trong Herdr | `herdr session attach <session>`; tên session là dòng `Herdr` core in ra. |
| TUI của Claude | `./tui` trong thư mục Mate, hoặc `mate mate-agent tui shop`. |

Console là cột hẹp bên trái. Trong WezTerm, hoặc Ghostty 1.3+ trên macOS, nó mở
pane của agent đang chọn bên phải. Ở terminal khác, Mate chạy headless và bạn theo
dõi trong Office hoặc City.

**TUI của Mate.** Mate headless có thêm pane shell `mate-<p>-shell` trong session
Herdr của workspace, cwd là thư mục của Mate (`.mate/projects/<p>/mate`). `./tui` ở
đó (chỉ cho Mate headless) mở session của Mate trong TUI gốc của Claude: dùng được
slash command, `/mcp`, plan mode. Mỗi workspace mở một TUI một lúc. Trong lúc TUI mở, Office, City và
console chỉ xem. Thoát Claude (`/exit`) thì Mate chạy tiếp như cũ.

### Phím trong console

`?` mở danh sách phím đầy đủ, kể cả phím bạn tự gắn. Hay dùng nhất:

| Phím | Việc |
| --- | --- |
| `↑` `↓` · `enter` | di chuyển · mở agent ở pane bên cạnh, mở project |
| `a` | các action của hàng đang chọn (`x`: stop Mate/Crew; ở hàng project là xoá project) |
| `s` · `n` | start Mate · project mới |
| `m` | đổi project giữa Auto và Manual |
| `w` · `v` | bảng Work của Mate · chỉ hiện project / chỉ hiện Mate |
| `tab` · `esc` · `q` | pane kế · về danh sách · thoát |

Auto (mặc định): cứ 90 giây core gom tin mới trong box thành một digest và giao
cho Mate. Manual: core không gõ digest, để bạn tự nói với Mate; 5 phút không nhắn gì
cho Mate thì tự về Auto.

## Cấu hình

Mọi file nằm trong `<workspace>/.mate/`. Core đọc `.env` và khối `dashboard:` của
`console.yaml` lúc bật, nên sửa xong thì `mate core restart`. Phím tự gắn trong
`console.yaml` thì console đọc lúc mở: sửa xong thì mở lại console.

**`workspace.yaml`**: `defaults.crew_harness` là harness mặc định cho Crew.
`lanes` đặt số Crew được giữ cùng một lane một lúc, ví dụ `lanes: {xcode: 4}`;
lane không khai thì là 2.

**`projects/<p>/project.yaml`**:

```yaml
repos:
  - name: api
    path: shop/api
    default_branch: main
    teardown: make worktree-teardown   # dọn resource khi đóng Crew
    lane: xcode                        # tuỳ chọn
review: "on"          # on | off; để trống thì Mate hỏi bạn trước
next_work: pick       # pick: Mate tự lấy việc tiếp; ask: hỏi bạn
mate:
  mode: headless      # headless (mặc định) | pane: TUI của Claude trong pane Herdr
  refresh_context: 280000   # làm mới context khi vượt ngưỡng (mặc định 150000; -1 tắt)
```

Lệnh tiện: `mate project review <p> on|off`, `mate project next-work <p> pick|ask`,
`mate project repo lane <p> <repo> <lane>`. Review: [docs/review.md](docs/review.md).

**Hướng dẫn cho agent**: `.mate/WORKSPACE.md` là luật chung cho mọi Mate,
`.mate/projects/<p>/PROJECT.md` là luật riêng của project (Mate cũng ghi vào đây),
`.mate/CREW.md` và `.mate/projects/<p>/CREW.md` là luật cho Crew. Không có tracker
mặc định; muốn Mate dùng Linear, Beads, Backlog.md…, xem
[docs/task-management.md](docs/task-management.md).

**`console.yaml`**: gắn phím mở editor hay tool, và mở dashboard ra LAN.

```yaml
bindings:
  - key: e
    scope: crew
    label: code
    command: [zed, "{worktree_dir}"]
    mode: external        # GUI app; terminal = tab trong WezTerm/Ghostty
dashboard:
  addr: 0.0.0.0:7777      # mặc định 127.0.0.1:7777
  allow_remote: true      # bắt buộc khi addr không phải loopback
```

Dashboard không có đăng nhập: ai mở được trang cũng đọc được mọi thứ bạn gõ và bấm
được nút (chat, start, stop, merge). Chỉ mở ra LAN khi bạn chắc.
Xem `mate console --help` và [docs/console-commands.md](docs/console-commands.md).

**Plugin** mang MCP server, skill và link cho mọi Mate trong workspace; Crew không
nhận plugin. `mate plugin add <thư-mục> --as <tên>`, rồi `mate plugin doctor <tên>`,
rồi restart Mate. Xem [docs/plugins.md](docs/plugins.md).

**`.env`**: `MATE_NOTIFICATIONS=off` tắt notification. `MATE_JEV=on` cùng
`MATE_JEV_API_KEY_FILE=<đường dẫn file key>` bật Jev đọc màn hình pane khi bộ đọc
fixture không nhận ra ([docs/jev-notices.md](docs/jev-notices.md)). Core đọc các
biến này từ file, không từ shell.

## Notification

Core gửi notification khi có news chờ bạn, và báo tình hình mỗi 10 phút khi không có
news. Thử bằng `mate notifications test`. Trên macOS, lần đầu mate build một app
Swift nhỏ (cần Xcode Command Line Tools) và xin quyền; bật banner ở **System
Settings → Notifications → Mate**. Bấm vào notification thì mở bàn captain trong
Office. Trên Linux mate dùng `notify-send`. Chi tiết:
[docs/notifications.md](docs/notifications.md).

## Lệnh hay dùng

Phần lớn lệnh `mate` là để Mate và Crew gọi. Captain thường chỉ cần:

| Lệnh | Việc |
| --- | --- |
| `mate init [<dir>]` | Tạo workspace |
| `mate project add\|list\|remove`, `mate project repo add\|list\|remove` | Quản lý project và repo |
| `mate`, `mate <dir>`, `mate core start\|stop\|restart\|status` | Chạy core |
| `mate console [<dir>]` | Mở console |
| `mate mate-agent start\|stop\|status\|refresh <p>` | Điều khiển Mate |
| `mate mate-agent repl\|tui <p>` | Chat với Mate, mở TUI của Mate |
| `mate crew list <p> [--all]` | Xem Crew |
| `mate peek <p> <crew>`, `mate report <p> <crew>` | Màn hình hiện tại, báo cáo của một Crew |
| `mate events <p> --follow --narrate` | Theo dõi timeline của project |
| `mate usage <p> [<crew>]` | Token đã dùng của project hoặc một Crew |
| `mate gc [--apply]` | Liệt kê process còn sót trong worktree đã đóng; `--apply` để dọn |
| `mate migrate [--dry-run]` | Đưa workspace cũ lên layout mới |

Chạy `mate` ngoài workspace để xem đủ danh sách lệnh; `mate --help` không có, và
trong workspace `mate` chạy core. Lệnh thiếu đối số in usage của nó.

## Phát triển

```sh
make build              # bin/mate
make check              # gofmt, go vet, go test qua scripts/gotestreport
MATE_LIVE=1 make test   # thêm test TestLive* chạy Herdr và harness thật
```

Test bị skip mà không khai trong `scripts/gotestreport/expected-skips.txt` làm build
đỏ. Suite xanh khi test live bị skip không chứng minh tính năng chạy với Herdr thật.
Bản đồ package và process: [docs/architecture.md](docs/architecture.md). Thuật ngữ:
[GLOSSARY.md](GLOSSARY.md). Quy ước cho agent làm việc trên repo này:
[AGENTS.md](AGENTS.md).
