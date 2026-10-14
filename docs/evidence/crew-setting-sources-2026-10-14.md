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
