# Video ngắn: bài học

Video ngắn đăng trên FB Reels và TikTok. Lõi của video là giá trị: mỗi video dạy **một bài học**, và người xem dùng được ngay, kể cả khi không bao giờ cài mate. mate xuất hiện như bằng chứng và bối cảnh ("tôi học được điều này khi xây mate"), không phải như quảng cáo. Người xem follow vì muốn học thêm từ người này.

Số liệu về nền tảng nằm ở `docs/research/social-algorithms-2026-10-09.md` (§2–4, Facebook và Instagram) và `docs/research/tiktok-algorithm-2026-10-09.md`.

## Hình dạng một bài học

1. **Hook, 3 giây đầu.** Đưa kết quả hoặc con số ra trước, ví dụ "66,5M token trong 3 ngày." Câu hỏi và giải thích để sau. Cả Reels lẫn TikTok đều đo việc người xem có lướt qua sớm hay không.
2. **Vấn đề:** một câu, gọi đúng nỗi đau của người đã dùng coding agent.
3. **Bằng chứng từ mate:** con số, log, sơ đồ hoặc ảnh console, lấy từ `material.md`.
4. **Bài học:** điều người xem làm được ngay, nói như một quy tắc. Ví dụ: "Cô lập cấu hình harness trước khi đo token."
5. **Kết:** nhắc lại bài học và tên series. Không xin follow, không xin comment từ khoá.

Một video chỉ một bài học. Có hai bài học thì làm hai video.

## Thông số

- Khung 9:16, 1080×1920. Chữ và chi tiết quan trọng nằm trong safe zone, tránh vùng nút bấm của app.
- Độ dài 30–90 giây. Meta vẫn ưu tiên short form. TikTok đẩy video trên 1 phút nhưng không công bố độ dài tối ưu, nên mỗi campaign thử cả dưới và trên 1 phút rồi so thời gian xem.
- **Luôn có voiceover và nhạc nền.** Reel tắt tiếng bị hạ hạng, và TikTok Creator Academy nêu đích danh dạng "hình full khung + voiceover". Giọng đọc là giọng TTS do tool riêng của captain tạo ra.
  - Tool đó chưa được tích hợp. Trong lúc chờ, mỗi video có một `SCRIPT.md` cho từng ngôn ngữ: lời đọc chia theo cảnh, có thời lượng ước tính.
  - Không tạo giọng bằng TTS của HyperFrames hay `/media-use`. Chừa sẵn track voiceover trong composition, timing khớp `SCRIPT.md`, để thay file audio vào sau.
  - Viết lời đọc để nghe, không để đọc: câu ngắn, mỗi câu một ý, con số nói theo cách người ta nói ("sáu mươi sáu triệu token"). Chữ trên màn hình là phần nhấn, không chép nguyên lời đọc.
  - Nếu giọng TTS là clone giọng một người thật, kể cả giọng captain, thì gắn nhãn AI trên TikTok. Giọng TTS chung chung thì không cần.
- Phụ đề hoặc chữ trên màn hình phải to và đọc được khi tắt tiếng. Reel chỉ toàn chữ bị đánh giá thấp, và TikTok xếp "slide video" vào nhóm chất lượng thấp. Vì vậy mỗi cảnh cần một chuyển động có chủ đích: sơ đồ đang dựng, terminal đang chạy, con số đang đếm.
- **Mỗi nền tảng một file render riêng** từ cùng project HyperFrames. Không tải video từ TikTok lên Reels: watermark của app khác bị hạ hạng, và tính nguyên bản là điều kiện để được gợi ý.

## Dựng bằng HyperFrames

- Luôn vào qua skill `/hyperframes` trước. Skill đó định tuyến:
  - bài học có giải thích, 60–90 giây → `/faceless-explainer`;
  - một con số hoặc một khoảnh khắc ngắn, cũng dùng làm media cho bài Threads → `/motion-graphics`;
  - câu chuyện của một PR cụ thể trong repo → `/pr-to-video`.
- **Nhận diện hình ảnh dùng lại, không vẽ mới:** `sites/videos/mate-promo/frame.md` (Code editorial). Gồm nền kem `#F5F3EB`, mực `#242723`, đúng một điểm coral `#E77859` mỗi khung, terminal xanh đậm `#242B24`, Be Vietnam Pro, Playfair Display nghiêng cho một cụm nhấn, IBM Plex Mono cho code, dấu ✱. Đây là gương mặt của captain trên video, nên mọi video phải nhận ra được là cùng một người làm.
- Nguồn hình:
  - `sites/mate/dist/docs/img/` (console, dashboard, Office) là dữ liệu demo, gắn nhãn "Minh hoạ · dữ liệu demo".
  - Ảnh chụp workspace thật thì captain phải đồng ý theo từng video.
  - Timeline của dashboard không được lên hình.
- Project đặt ở `sites/videos/<campaign-slug>-<NN>-<vi|en>/`. Trong `BRIEF.md` ghi:
  - `aspect: 1080x1920`;
  - `language`;
  - `destination: reels+tiktok`;
  - `message` là đúng một câu bài học;
  - mục chất liệu đã dùng.

## Luật riêng của TikTok

Nguồn: `docs/research/tiktok-algorithm-2026-10-09.md` §11.

- **Nhãn AI.** Motion graphics của HyperFrames và giọng TTS chung chung không cần nhãn AIGC. Clone giọng của chính captain thì gắn nhãn. Trước khi đăng, kiểm tra file xuất ra có nhúng C2PA không: nếu có, TikTok tự gắn nhãn AI và không gỡ được.
- **Content disclosure.** Video quảng bá mate như một sản phẩm phải bật content disclosure. Video bài học thuần, nơi mate chỉ là bối cảnh, thì không cần. Đây là suy luận của research; đọc Branded Content Policy trước loạt đầu.
- **Search.** TikTok đọc keyword ở ba chỗ: lời nói, chữ trên màn hình và caption. Đặt cụm từ người xem sẽ gõ vào cả ba chỗ, nhưng không nhồi.
- **Trước và sau khi đăng.** Chạy Content Check Lite trong TikTok Studio trước khi đăng. Sau khi đăng, xem analytics có video nào bị đánh dấu không đủ điều kiện For You không.

## Hai nhánh

- Dựng bản VN trước. Bản EN là một project riêng, chữ trên màn hình viết lại cho người đọc quốc tế. Bố cục và chuyển động có thể giữ nguyên.
- Facebook có AI dịch và lồng tiếng Reels có tiếng Việt, cần từ 1.000 follower. Khi đủ điều kiện, captain có thể thử cách này thay cho một bản EN riêng trên Facebook. Bản riêng vẫn là mặc định cho TikTok.
- Caption khi đăng mang từ khoá mà người xem sẽ gõ để tìm, ví dụ "AI coding agent", "Claude Code", "Codex", "token". Viết riêng cho từng ngôn ngữ.
- **TikTok VN:** Việt Nam không có STEM feed, nên chỉ còn For You và Search; Search vì thế càng quan trọng. Burn-in phụ đề tiếng Việt ngay trong HyperFrames, vì chưa chắc TikTok tự tạo phụ đề tiếng Việt.
- **TikTok EN:** research đề xuất tách một tài khoản tiếng Anh riêng, đặt app sang tiếng Anh, vì ngôn ngữ của tác giả là một tín hiệu. Không dùng VPN hay giả vị trí; TikTok cấm các mẹo vượt hệ thống. Cơ hội lớn nhất là STEM feed, chỉ nhận video tiếng Anh: video bài học thuần, có dẫn chứng, **không quảng bá mate**, gắn #STEMTok. Video giới thiệu mate thì không gắn tag này. Ghi nguồn trên màn hình cho mọi con số.

## Nối về Threads

Mỗi video có một bài neo trên Threads (xem `threads.md`). Clip 5–10 giây cắt từ video có thể làm media cho bài neo. Ở dạng này video là bằng chứng, còn hội thoại diễn ra trên Threads.
