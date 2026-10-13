# Đào chất liệu

Chất liệu là một sự kiện có thật trong repo mà dev ngoài kia sẽ thấy đáng gửi cho đồng nghiệp hoặc đáng cãi lại. mate có lợi thế là mọi quyết định, mọi lần đo, mọi lần hỏng đều được ghi lại có ngày. Campaign không cần bịa gì, chỉ cần đào.

## Mạch chất liệu

| Nơi | Chứa gì | Thường thành |
| --- | --- | --- |
| `docs/research/` | Phân tích có số liệu: token của Mate, chi phí điều phối, nhận diện trạng thái TUI, cách firstmate prompt và nhớ | Con số khó tin, bài học về chi phí |
| `docs/evidence/` | Nghiệm thu chạy thật, có log. Postmortem: Enter bị nuốt, Herdr chết. Đo pi theo hợp đồng harness | Thú tội, so sánh có đo, nhật ký build |
| `docs/plans/` | Phương án và lý do chọn: registry harness, Jev làm bộ quan sát, layout workspace | Quyết định ngược số đông |
| `docs/mvp.md` §2 | Các quyết định đã chốt của mate, mỗi điều là một quan điểm | Quyết định ngược số đông, bài học kiến trúc |
| `docs/console-design.md` | Triết lý console "20% terminal butler" | Quan điểm về UX công cụ dev |
| `git log` | Hành trình: viết lại từ v1 (repo archive `mate-v1`), đổi tên từ matev2 ngày 2026-09-24, hàng trăm commit trong vài tuần | Nhật ký build, câu chuyện dài cho arc |
| `sites/mate/dist/` | Câu thương hiệu và cách mate tự giới thiệu | Câu kết, nhận diện |

Lấy ngày và số commit bằng `git log`. Đừng dựa vào bảng này cho mấy con số đó.

## Một mục chất liệu

Mỗi mục trong `material.md` của campaign gồm:

- **id** ngắn, ví dụ `M07-token-co-dinh`. Bài Threads và video tham chiếu mục bằng id này.
- **nguồn:** `file:line`, kèm ngày của tài liệu.
- **sự kiện:** con số hoặc chuyện đã xảy ra, chép đúng như nguồn.
- **góc:** vì sao dev sẽ gửi đi hoặc cãi lại. Một câu.
- **bài học:** điều người xem làm được ngay. Một câu. Mục không rút ra được bài học vẫn dùng được cho Threads, nhưng không dùng cho video.
- **dạng:** dạng mồi lửa phù hợp trong `threads.md`, và có hợp làm video không.
- **cần hỏi captain:** có đưa tên project thật, nội dung transcript hay ảnh workspace thật lên không.
- **đã dùng:** campaign và ngày đăng, điền sau khi đăng.

## Kiểm độ tươi

Docs ghi lại trạng thái ở một ngày cụ thể. Trước khi dùng, kiểm lại sự kiện với code và docs hiện tại:

- Con số đo trước một bản sửa: viết nó như chuyện đã qua ("hôm 27/9 Mate đốt…") và nói thêm điều đã đổi sau đó, nếu có.
- Tính năng trong plan chưa ship: ghi là "đang làm", không mô tả như đã có.
- Đối chiếu với `docs/mvp.md` để chắc tính năng còn đúng như mô tả.
