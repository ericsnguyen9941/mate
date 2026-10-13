---
name: mate-marketing
description: Marketing cho mate theo giọng của mate. Dùng khi lên campaign, viết bài Threads, hoặc làm video ngắn FB Reels/TikTok (HyperFrames) để xây tên tuổi cá nhân của captain.
---

# mate marketing

Campaign của mate xây **tên tuổi cá nhân của captain**: người làm mate và chạy một crew AI mỗi ngày. Sao GitHub và lượt cài không phải là mục tiêu. Người xem cần nhớ đến *người này* và quay lại xem người này làm gì tiếp. mate là chất liệu kể chuyện, còn nhân vật chính là captain.

Hai động cơ, mỗi động cơ có một lõi riêng:

- **Threads là trung tâm, lõi là viral.** Một bài là một *mồi lửa*: một khẳng định có bằng chứng, đủ sắc để dev muốn reply, cãi hoặc bấm vào profile. Playbook: [references/threads.md](references/threads.md).
- **Video ngắn (FB Reels + TikTok), lõi là giá trị.** Một video là một *bài học*: người xem mang về được một điều dùng ngay, kể cả khi không bao giờ cài mate. Dựng bằng HyperFrames. Playbook: [references/short-video.md](references/short-video.md).

Mọi thứ đều có hai bản: tiếng Việt cho dev Việt và tiếng Anh cho dev quốc tế. Bản thứ hai được **viết lại** cho người đọc của nó, không dịch từ bản kia.

Captain không lộ mặt. Gương mặt thay thế là giọng văn, quan điểm, và nhận diện hình ảnh (✱ coral, theme terminal của mate). Video có giọng đọc TTS (`vieneu-say`), chọn theo ngữ cảnh từng video.

## Giọng của mate

Đây là phần mọi bài, mọi video đều dùng.

- **Người nói là captain**, ngôi thứ nhất ("tôi" / "I"). Không dùng giọng doanh nghiệp kiểu "chúng tôi", "đội ngũ mate".
- **Đo trước, nói sau.** Mỗi khẳng định đi kèm một con số, một log hoặc một quyết định có trong repo. Đây là tính cách gốc của project: docs ghi "Đây là đề xuất, chưa đổi dòng code nào", evidence ghi rõ chỗ chưa chứng minh được. Bài đăng cũng vậy: nói đúng mức bằng chứng cho phép.
- **Kể chỗ hỏng nhiều hơn khoe tính năng.** Ví dụ: Enter bị nuốt khi gửi chữ vào Codex, Mate đốt 66,5M token trong 3 ngày, 87K token cố định ở mọi lần gọi mà 50K không phải của mate. Giọng khô, thẳng, tự giễu nhẹ, không hype và không dùng chữ "cách mạng".
- **Ẩn dụ con tàu là gia vị**, mỗi bài nhiều nhất một lần: captain cầm lái, Mate là trưởng nhóm, crew lên đường, lên tàu. Câu thương hiệu có sẵn: "Bạn có ý tưởng. Mate có crew." và "Lợi thế riêng của một team nhỏ."
- **Từ vựng thương hiệu:** "mate" viết thường là sản phẩm, "Mate" viết hoa là vai trò, "crew" là agent thực thi. Thuật ngữ kỹ thuật giữ tiếng Anh như trong docs (worktree, harness, token, brief).
- **Bản tiếng Việt** viết câu ngắn, như dev nói với dev. **Bản tiếng Anh** viết theo giọng dev builder quốc tế, gọn, và dùng ví dụ họ quen.

Ranh giới sự thật (guardrail cứng):

- Chỉ nói những gì code mate làm được hôm nay. Nguồn là `docs/mvp.md` và evidence.
- Không có số liệu khách hàng, không testimonial. Không ngụ ý mate liên kết hay được Anthropic bảo trợ; mate là "dự án độc lập, xây để làm việc cùng Claude Code".
- Ảnh từ `sites/mate/dist/docs/img/` là dữ liệu demo và phải ghi "dữ liệu demo". Timeline của dashboard chứa nguyên văn prompt của captain, nên không đưa lên.
- Nếu muốn dùng tên project thật, nội dung transcript hay workspace thật, hỏi captain một lần và ghi câu trả lời vào `CAMPAIGN.md`.

## Quy trình lên campaign

Campaign nằm ở `sites/campaigns/<YYYY-MM-DD>-<slug>/`. Thư mục `sites/` bị gitignore: đây là không gian marketing local, giống `sites/videos/`.

1. **Nạp bối cảnh.**
   - Đọc số liệu thuật toán ở `docs/research/social-algorithms-2026-10-09.md` (Threads, Instagram, Facebook) và `docs/research/tiktok-algorithm-2026-10-09.md`. Đây là nguồn duy nhất cho mọi nhận định về nền tảng; điều gì không có ở đó thì ghi là suy luận.
   - Đọc các campaign trước trong `sites/campaigns/`, phần log kết quả, để không lặp chất liệu đã dùng.
   - Hỏi captain những gì `CAMPAIGN.md` chưa ghi: khoảng thời gian, nhịp đăng, một hay hai tài khoản Threads, và một hay hai tài khoản TikTok.
   - Xong khi `CAMPAIGN.md` có mục tiêu, khoảng thời gian, nhịp đăng của từng ngôn ngữ và các quyết định trên.
2. **Đào chất liệu** theo [references/material.md](references/material.md) và ghi vào `material.md`.
   - Xong khi số mục đủ cho cả lịch cộng thêm khoảng một phần ba dự phòng.
   - Mỗi mục có nguồn `file:line`, con số hoặc sự kiện cụ thể, và lý do một dev sẽ gửi bài này cho đồng nghiệp hoặc muốn cãi lại.
3. **Dựng arc và lịch** trong `CAMPAIGN.md`.
   - Arc là một câu chuyện nhiều tuần, ví dụ "một tuần sống với crew AI" hay "những lần crew làm hỏng việc".
   - Mỗi video có một bài Threads neo: bài Threads là mồi lửa của cùng bài học đó. Threads khơi hội thoại, video trả giá trị.
   - Xong khi mỗi slot trong lịch có ngày, nền tảng, ngôn ngữ và mục chất liệu đã chọn.
4. **Viết bài Threads** theo [references/threads.md](references/threads.md), lưu vào `threads/vi/` và `threads/en/`.
5. **Dựng video** theo [references/short-video.md](references/short-video.md).
6. **Soát sự thật.**
   - Mở lại nguồn của từng con số, tên và khẳng định trong mọi bài và mọi khung video, rồi đối chiếu với code và docs *hiện tại*. Docs có ngày; một con số đúng hôm 27/9 có thể đã cũ.
   - Xong khi mọi khẳng định có dòng nguồn khớp và không còn vi phạm guardrail nào. Bảng soát nằm cuối `CAMPAIGN.md`.
7. **Sau khi đăng**, ghi số liệu theo bài vào phần log của `CAMPAIGN.md`: reach tới non-follower, reply, send/share, follow mới, lượt vào profile, và quốc gia (insights theo quốc gia). Cứ hai tuần đọc lại log và chỉnh arc. Dữ liệu 4–6 tuần đầu dùng để chốt chuyện một hay hai tài khoản Threads.
