# Crew Claude với `--setting-sources project` (CON-68, 2026-10-14)

Câu hỏi: crew Claude nhận `--setting-sources project` có mất skill project
của worktree không, và có thật bỏ skill toàn cục của captain không.

Máy captain, Claude Code 2.1.287, thư mục tạm `/tmp/ssprobe` (git init, một
skill project `.claude/skills/probe-local/SKILL.md`), `ANTHROPIC_BASE_URL` và
`CLAUDECODE` bỏ khỏi env như `probeEnv` của `mcp_probe_live_test.go`.

```
claude -p --setting-sources project --autocompact 160000 --model haiku \
  "List the names of every skill available to you via the Skill tool, comma separated, nothing else."
```

```
probe-local, dataviz, update-config, keybindings-help, code-review, simplify, fewer-permission-prompts, loop, schedule, claude-api, workflow-authoring, run, plugin-authoring, init, security-review
```

Cùng lệnh, không có `--setting-sources`:

```
chrome-devtools-axi, code-review, codebase-design, diagnosing-bugs, domain-modeling, gh-axi, git-guardrails-claude-code, grilling, lavish, migrate-to-shoehorn, pr, prototype, quota-axi, research, scaffold-exercises, setup-pre-commit, tdd, terminal-browser, tiktok-growth, wizard, writing-for-agents, probe-local, code-review:code-review, claude-md-management:revise-claude-md, frontend-design:fronten…
```

Kết luận: skill project vẫn nạp; skill trong `~/.claude/skills` và skill của
plugin người dùng không còn; skill đi kèm Claude Code vẫn còn. `--autocompact
160000` được CLI nhận. Chưa đo: một crew thật compact ở ngưỡng 160000 (cần chạy
một task dài; đo bằng `mate usage <project>` sau 10 task).

File settings riêng của crew (`--settings crews/<id>/settings.json`) vẫn nạp khi có
`--setting-sources project`: một file chỉ chứa `{"env":{"MATE_PROBE_FLAG":"from-flag-settings"}}`
truyền bằng `--settings`, cùng `--setting-sources project --dangerously-skip-permissions`,
và `echo $MATE_PROBE_FLAG` trong phiên in `from-flag-settings`.

## 2026-10-15: crew vẫn đi qua cliproxy sau khi chép `apiKeyHelper` và `env`

Review PR #59 chỉ ra `--setting-sources project` bỏ `apiKeyHelper` và `env` của
`~/.claude/settings.json`, nên crew rời cliproxy (127.0.0.1:8317) và đi thẳng
api.anthropic.com bằng login keychain. Bản sửa: `CrewClaudeSettings` chép hai
khoá đó từ settings.json của config dir vào `crews/<id>/settings.json`.

Đo trên máy captain, Claude Code 2.1.292. File settings sinh bằng chính
`claude.Claude{}.Prepare` (RoleCrew) của bản sửa đọc `~/.claude/settings.json` thật:
khoá `apiKeyHelper`, `autoMemoryEnabled`, `env` (`ANTHROPIC_BASE_URL`,
`ANTHROPIC_DEFAULT_{HAIKU,OPUS,SONNET}_MODEL`). Lệnh chạy giống argv crew, env đã
bỏ `ANTHROPIC_BASE_URL`, `ANTHROPIC_API_KEY`, `CLAUDECODE`:

```
claude -p "Reply with the single word ok." --dangerously-skip-permissions \
  --settings <file> --setting-sources project --autocompact 160000 --model haiku
```

lsof theo pid mỗi 0,25 s trong lúc chạy:

| `--settings` | kết nối tới 127.0.0.1:8317 | kết quả |
|---|---|---|
| file của bản sửa | có (`127.0.0.1:55107->127.0.0.1:8317`), kèm cảnh báo "claude.ai connectors are disabled because ANTHROPIC_API_KEY or another auth source is set" | `ok` |
| `{"autoMemoryEnabled": false}` (như trước bản sửa) | 0 lần | `ok` |

Các kết nối :443 có ở cả hai lượt. Chưa đo: crew thật spawn dưới Herdr bằng
binary mới. Lượt này phải làm trước khi cài binary.
