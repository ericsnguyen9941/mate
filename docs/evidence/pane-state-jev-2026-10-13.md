# Jev đọc trạng thái pane — eval, 2026-10-13

Một lượt chạy thật với API TypeSafe, model `jev-1.13.0`, từ máy phát triển, đo câu hỏi trạng thái pane (`pane_state`, 14 nhãn) trên 38 màn hình khởi động và injection của repo. Một request cho mỗi màn hình, tuần tự, deadline 8 giây, không retry, không request nào lỗi. Nội dung gửi đi là `jev.Prepare` (đã che credential); key không in ra, không commit.

Prompt hash của câu hỏi (`jev.StatePromptHash()`): `eca0358eee41`. Câu hỏi hoặc nhãn đổi thì hash đổi và cassette phải ghi lại.

Tái hiện không cần mạng:

```sh
go run ./scripts/jeveval -question state -replay -cassette scripts/jeveval/testdata/cassette-state
```

Ghi lại với mạng: bỏ `-replay`. In corpus và nhãn: `-question state -list`.

## Corpus và nhãn

Corpus: `internal/harness/*/testdata/startup/` (35 màn hình) và `scripts/jeveval/testdata/screens/` (3 màn hình injection). Nhãn là trạng thái `fixture.ObserveState` đọc ra. Màn hình fixture đọc là `unknown` (13 màn hình) vẫn được hỏi nhưng không chấm: đó là bộ phân loại khởi động từ chối kết luận, không phải nhãn.

Hai điểm cần đọc kèm bảng. Màn hình `claude_injection_draft_says_empty` và `codex_injection_busy_output_says_empty` được dựng cho câu hỏi composer; ở đây nhãn fixture của chúng chỉ nói về khởi động (`unknown` và `ready`), không nói composer đang nháp hay đang bận. Vì vậy dòng 38 (fixture `ready`, Jev `busy` 0.590) tính là "sai" theo nhãn nhưng nhiều khả năng Jev đúng; nhãn mới là chỗ yếu.

## Kết quả từng màn hình

| # | File | Trạng thái fixture | Jev | Khớp | Độ tin cậy | Latency ms |
| ---: | --- | --- | --- | --- | ---: | ---: |
| 1 | `claude/startup/claude-2.1.270-ready.txt` | ready | auth_required | no | 0.780 | 384 |
| 2 | `claude/startup/claude-2.1.270-trust-dialog-accept-selected.txt` | dialog_trust_dir | dialog_trust_dir | yes | 1.000 | 212 |
| 3 | `claude/startup/claude-2.1.270-trust-dialog.txt` | dialog_trust_dir | dialog_trust_dir | yes | 1.000 | 212 |
| 4 | `claude/startup/claude-2.1.282-ready.txt` | ready | ready | yes | 0.470 | 210 |
| 5 | `claude/startup/claude-2.1.285-bypass-dialog.txt` | dialog_bypass | dialog_bypass | yes | 1.000 | 210 |
| 6 | `codex/startup/codex-0.154.0-hooks-review.txt` | dialog_hooks_review | dialog_hooks_review | yes | 1.000 | 219 |
| 7 | `codex/startup/codex-0.154.0-ready.txt` | ready | ready | yes | 0.800 | 220 |
| 8 | `codex/startup/codex-0.154.0-resume-ready.txt` | ready | ready | yes | 0.870 | 214 |
| 9 | `codex/startup/codex-0.154.0-resume-update-dialog.txt` | dialog_update | dialog_update | yes | 1.000 | 207 |
| 10 | `codex/startup/codex-0.154.0-trust-dialog-quit-selected.txt` | dialog_trust_dir | dialog_trust_dir | yes | 1.000 | 218 |
| 11 | `codex/startup/codex-0.154.0-trust-dialog.txt` | dialog_trust_dir | dialog_trust_dir | yes | 1.000 | 225 |
| 12 | `codex/startup/codex-0.156.1-hooks-closed-ready.txt` | ready | ready | yes | 0.640 | 262 |
| 13 | `codex/startup/codex-0.156.1-hooks-review-two.txt` | dialog_hooks_review | dialog_hooks_review | yes | 1.000 | 204 |
| 14 | `codex/startup/codex-0.156.1-hooks-review.txt` | dialog_hooks_review | dialog_hooks_review | yes | 1.000 | 204 |
| 15 | `codex/startup/codex-0.156.1-hooks-sessionstart-own-truncated.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 210 |
| 16 | `codex/startup/codex-0.156.1-hooks-sessionstart-own-wrapped.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 234 |
| 17 | `codex/startup/codex-0.156.1-hooks-sessionstart-own.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 205 |
| 18 | `codex/startup/codex-0.156.1-hooks-sessionstart-two-all-trusted.txt` | unknown | dialog_hooks_review | unscored | 0.840 | 208 |
| 19 | `codex/startup/codex-0.156.1-hooks-sessionstart-two-foreign-selected.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 223 |
| 20 | `codex/startup/codex-0.156.1-hooks-sessionstart-two-own-selected.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 210 |
| 21 | `codex/startup/codex-0.156.1-hooks-sessionstart-two-own-trusted.txt` | unknown | dialog_hooks_review | unscored | 0.990 | 224 |
| 22 | `codex/startup/codex-0.156.1-hooks-table-review-two.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 216 |
| 23 | `codex/startup/codex-0.156.1-hooks-table-review.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 210 |
| 24 | `codex/startup/codex-0.156.1-hooks-table-trusted.txt` | unknown | dialog_hooks_review | unscored | 0.920 | 211 |
| 25 | `codex/startup/codex-0.156.1-ready.txt` | ready | ready | yes | 0.770 | 205 |
| 26 | `codex/startup/codex-0.156.1-trust-dialog-quit-selected.txt` | dialog_trust_dir | dialog_trust_dir | yes | 1.000 | 219 |
| 27 | `codex/startup/codex-0.156.1-trust-dialog.txt` | dialog_trust_dir | dialog_trust_dir | yes | 1.000 | 223 |
| 28 | `codex/startup/codex-0.162.1-hooks-review.txt` | dialog_hooks_review | dialog_hooks_review | yes | 1.000 | 208 |
| 29 | `codex/startup/codex-0.162.1-hooks-sessionstart-own-home.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 208 |
| 30 | `codex/startup/codex-0.162.1-hooks-sessionstart-own.txt` | unknown | dialog_hooks_review | unscored | 1.000 | 224 |
| 31 | `codex/startup/codex-0.162.1-loading-composer.txt` | ready | ready | yes | 0.910 | 215 |
| 32 | `codex/startup/codex_update_banner_ready.txt` | ready | ready | yes | 0.980 | 212 |
| 33 | `codex/startup/codex_update_dialog.txt` | dialog_update | dialog_update | yes | 1.000 | 209 |
| 34 | `codex/startup/codex_update_dialog_after_enter.txt` | dialog_trust_dir | dialog_trust_dir | yes | 0.950 | 215 |
| 35 | `codex/startup/codex_update_dialog_skip_selected.txt` | dialog_update | dialog_update | yes | 1.000 | 208 |
| 36 | `claude_injection_draft_says_empty.txt` | unknown | ready | unscored | 0.870 | 208 |
| 37 | `claude_injection_empty_after_trust.txt` | ready | ready | yes | 0.910 | 218 |
| 38 | `codex_injection_busy_output_says_empty.txt` | ready | busy | no | 0.590 | 219 |

## Độ chính xác

- Trên màn hình được chấm: **23/25 (92,0%)**; 13 màn hình không chấm.
- Hai màn hình sai: dòng 1 `claude-2.1.270-ready.txt` (Jev `auth_required`, 0.780) và dòng 38 (xem trên, Jev `busy`, 0.590). Cả hai dưới 0.85.
- Câu trả lời sai ở độ tin cậy từ 0.85 trở lên: 0 trong 25 màn hình được chấm. Nhưng có 1 màn hình không chấm, `claude_injection_draft_says_empty` (chữ đã gõ trong composer, tức là một bản nháp), mà Jev trả `ready` ở 0.870, trên ngưỡng và rất có thể sai.
- Câu trả lời được chấm ở độ tin cậy từ 0.85 trở lên: 19, cả 19 khớp. Sáu câu được chấm dưới 0.85 (hai sai, bốn đúng), nên ở ngưỡng 0.85 sáu màn hình đó rơi về fixture.
- Nhãn `ready` là chỗ Jev yếu nhất: 10 màn hình, độ tin cậy 0.470, 0.640, 0.770, 0.780 (sai), 0.800, 0.870, 0.910, 0.910, 0.980 và 0.590 (sai). Chỉ bốn trong mười đạt 0.85. Mọi dialog đều 0.95 đến 1.000 (riêng `codex_update_dialog_after_enter` 0.950).
- Latency 38 request: p50 212 ms, p95 262 ms, tối đa 384 ms.

## Câu trả lời của Jev trên màn hình không chấm

Mười hai màn hình hook của Codex (`hooks-sessionstart-*`, `hooks-table-*`) mà bộ phân loại khởi động từ chối: Jev đều trả `dialog_hooks_review`, độ tin cậy 1.000 (chín màn hình), 0.990, 0.920, 0.840 (`two-all-trusted`). Một màn hình injection `claude_injection_draft_says_empty`: Jev `ready` 0.870, trong khi đây là pane có chữ trong composer; nếu nhãn thật là `draft` thì đây là một câu sai ở trên ngưỡng, và mô hình nhãn hiện tại không chứng minh được hay bác bỏ điều đó.

Các màn hình hook ở trạng thái "đã tin cậy" (`two-all-trusted`, `table-trusted`) mà Jev vẫn gọi là `dialog_hooks_review` có thể là Jev đúng hoặc sai tuỳ vào việc dialog còn mở hay không; fixture không có nhãn, nên không kết luận.

## Có đủ tin cậy cho settle ở ngưỡng 0.85 không

Chưa kết luận là có. Điều con số nói:

- Ở ngưỡng 0.85, trên 25 màn hình có nhãn không có câu sai nào lọt qua, và 19 câu được nhận đều đúng. Đó là tín hiệu tốt, nhưng 19 mẫu không đủ để nói tỉ lệ sai dưới ngưỡng thấp hơn một vài phần trăm.
- Điểm mù: 60% số màn hình `ready` (6/10) bị giữ lại dưới 0.85, và một màn hình `ready` thật bị gọi là `auth_required` ở 0.780. Settle dùng Jev sẽ hay rơi về fixture đúng ở trạng thái thường gặp nhất; nó an toàn (không nhận câu sai) nhưng không tiết kiệm.
- Phạm vi: corpus gần như chỉ có Claude và Codex ở khởi động; các trạng thái `busy`, `draft`, `dialog_tool_permission`, `auth_required`, `quota_exhausted`, `exited` không có màn hình nào được chấm. Câu hỏi này chưa được đo ở chính những trạng thái đó.
- Hạn chế nhãn: một lượt, một ngày, một model; nhãn fixture không phủ composer, nên sai lệch ở `draft`/`busy` chưa kiểm tra được (xem dòng injection).

Kết luận: độ tin cậy của Jev ở ngưỡng 0.85 cho settle chưa được chứng minh. Lý do cụ thể nhất là chính một lần trượt đã thấy: một bản nháp (`claude_injection_draft_says_empty`) bị đọc là `ready` ở 0.870, trên ngưỡng, và đó đúng loại sai nguy hiểm nhất (settle sẽ coi pane là sẵn sàng khi có chữ đang nằm trong composer). Thêm vào đó: Jev yếu ở `ready` (6/10 dưới 0.85, một màn bị đọc là `auth_required`), corpus không có màn nào được chấm cho `draft`, `busy`, `dialog_tool_permission`, `auth_required`, `quota_exhausted`, `exited`, và đây chỉ là một lượt chạy. Trên 25 màn hình được chấm không có câu sai nào ở từ 0.85 trở lên, nhưng con số đó không bù được các điểm trên.
