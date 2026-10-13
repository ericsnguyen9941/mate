# Thuật toán phân phối của Threads, Instagram và Facebook (tính đến 2026-10-09)

Ngày: 2026-10-09. Mục đích: chuẩn bị marketing mate (CLI Go mã nguồn mở điều phối AI coding agent) trên Threads, Instagram và Facebook. Có hai phiên bản nội dung: tiếng Việt cho dev Việt và tiếng Anh cho dev quốc tế. Mục tiêu là **xây personal brand cho người làm mate**, không phải kéo GitHub star hay lượt cài.

Quy ước nguồn:

- Mỗi khẳng định đều kèm link và ngày. Nguồn chính gồm Meta Transparency Center (system card), about.fb.com, about.instagram.com, creators.instagram.com, creators.facebook.com, engineering.fb.com, Help Center của Meta, và bài đăng của Adam Mosseri hoặc tài khoản @threads.
- Điều gì chỉ thấy ở nguồn thứ cấp thì ghi **[chưa xác minh từ nguồn chính]** kèm link.
- Phần suy luận của người viết ghi **[suy luận]**.
- Meta không công bố trọng số bằng số cho tín hiệu nào. Mọi con số kiểu "sends nặng gấp 3–5 lần like" đều không có nguồn chính (xem §5).

---

## Tóm tắt nhanh

1. Cả ba nền tảng xếp hạng bằng **dự đoán hành động của từng người xem**, không dùng một điểm "trending" chung. Mỗi bề mặt (Feed, Feed Recommendations, Explore, Reels, Threads Feed) có system card riêng và danh sách dự đoán riêng ([Meta Transparency Center](https://transparency.meta.com/features/explaining-ranking/)).
2. **Instagram**: Mosseri (01/2025) nói ba tín hiệu quan trọng nhất là watch time, likes và sends. Likes nặng hơn một chút với follower, sends nặng hơn một chút với non-follower ([SMT 22/01/2025](https://www.socialmediatoday.com/news/instagram-shares-algorithm-insights-2025/738034/), dẫn [video gốc](https://www.instagram.com/p/DFFyRp-pINJ/)).
3. **Threads**: reply chiếm gần một nửa lượt xem; tài khoản @threads khuyên đăng 2–5 lần/tuần, mỗi bài gắn một topic tag ([@threads, 21/10/2024](https://www.threads.com/@threads/post/DBZbkQHxjeD); [creators.instagram.com/threads](https://creators.instagram.com/threads)). Link không bị hạ hạng, và từ khoảng 05/2025 đã được "rank properly" ([Mosseri 03/06/2025](https://www.threads.com/@mosseri/post/DKdzO8fumTq)).
4. **Facebook**: Pages và profile chế độ chuyên nghiệp đang bị giới hạn số bài có link nếu không trả phí. Meta One (15/09/2026) xếp "links in organic posts and reels" vào gói Advanced ([about.fb.com](https://about.fb.com/news/2026/09/introducing-meta-one-subscription-service-more-features-ai/)).
5. **Tính nguyên bản** là điều kiện để được gợi ý trên cả Instagram (30/04/2026) lẫn Facebook (07/2025, 03/2026). Tài khoản chủ yếu đăng lại nội dung người khác sẽ không được gợi ý tới non-follower.

---

## 0. Khung chung của Meta

- **System card**: Meta công bố một card cho mỗi hệ AI xếp hạng. Mỗi card có các bước (gom inventory → tín hiệu → dự đoán → xếp hạng), danh sách dự đoán và tín hiệu đầu vào chính ([danh mục card](https://transparency.meta.com/features/explaining-ranking/); [Meta AI blog 2023 giới thiệu 22 card](https://ai.meta.com/blog/how-ai-powers-experiences-facebook-instagram-system-cards/)).
- **Nội dung "connected" và "recommended" do hai hệ khác nhau xếp hạng.**
  - Card Instagram Feed ghi rõ nội dung gợi ý "not powered by the AI we describe in this system card" ([IG Feed, cập nhật 29/06/2026](https://transparency.meta.com/features/explaining-ranking/ig-feed/)).
  - Card Facebook Feed nói tương tự và trỏ sang card Feed Recommendations ([FB Feed, cập nhật 06/10/2026](https://transparency.meta.com/features/explaining-ranking/fb-feed/)).
  - Muốn tới được non-follower thì phải xem card của **bề mặt gợi ý**, không phải card Feed.
- **Recommendation Guidelines** đặt chuẩn cao hơn Community Standards. Nội dung được phép đăng vẫn có thể bị loại khỏi gợi ý: ví dụ bạo lực, sản phẩm bị kiểm soát (thuốc lá, vape, dược phẩm) ([Instagram Ranking Explained, 31/05/2023](https://about.instagram.com/blog/announcements/instagram-ranking-explained)). Card FB Feed Recommendations loại "content that goes against our Recommendation Guidelines" ([FB Feed Recs, 24/06/2026](https://transparency.meta.com/features/explaining-ranking/fb-feed-recommendations/)).
- **Content Distribution Guidelines (CDG) – các loại nội dung bị hạ phân phối trên Facebook.**
  - Trang hiện hành chỉ còn 4 nhóm: Clickbait Links, Engagement Bait, Fact-Checked Misinformation, Content Likely Violating Our Community Standards ([Types of content we demote, cập nhật 02/07/2025](https://transparency.meta.com/features/approach-to-ranking/types-of-content-we-demote/)).
  - Tháng 07/2025 Meta **xóa 14 guideline**, trong đó có Ad Farms, Domains with Limited Original Content, Links to Domains and Pages with High "Click-Gap", Unoriginal News Articles, Posts that Indicate Suspicious Virality, Inauthentic Sharing ([CDG changes & corrections](https://transparency.meta.com/features/approach-to-ranking/cdgs-changes-corrections)).
  - ⚠️ Trang "Reducing the distribution of problematic content" (cập nhật 25/04/2025) vẫn liệt kê "Links to websites that are covered with ads, slow to load or broken" và "Limited originality content…". Trang này **có thể đã lỗi thời** so với đợt dọn 07/2025 ([link](https://transparency.meta.com/enforcement/taking-action/lowering-distribution-of-problematic-content/)). Tính nguyên bản giờ được xử lý bằng chính sách originality riêng (§2.4, §3.4).

---

## 1. Threads

### 1.1 Bề mặt đưa nội dung tới non-follower

- **Feed "For you"** trộn bài của người mình follow với nội dung gợi ý. Nếu chọn tab Following thì "you won't see recommended content" ([Threads Feed system card, cập nhật 07/03/2025](https://transparency.meta.com/features/explaining-ranking/ig-threads-feed/)).
- **Topic tag**: mỗi bài gắn tối đa 1 topic, dài tối đa 50 ký tự. Tag giúp người khác "discover your post … in real time" ([Help Center: Tag a topic](https://help.instagram.com/1356090605000312) – nội dung lấy qua kết quả tìm kiếm vì trang render bằng JS; [@threads 07/12/2023](https://www.threads.com/@threads/post/C0j6y3VrD2o): "Why just one? This makes it easier for others who care about that topic to find and read your post").
- Meta (03/2025): "posts that include a topic generally receive more views than those without" ([about.fb.com 03/2025](https://about.fb.com/news/2025/03/new-threads-features-more-personalized-experience-you-control/)).
- **Communities** (ra mắt 02/10/2025) là bước phát triển tiếp của topic tag và custom feed: "We've evolved topic tags and custom feeds…". Tag có biểu tượng ba chấm là tag có community riêng. Đến 15/12/2025 đã có hơn 200 chủ đề ([about.fb.com 02/10/2025, cập nhật 15/12/2025](https://about.fb.com/news/2025/10/introducing-threads-communities-find-your-people/)). Communities ra khỏi beta và có thêm Communities Hub vào 06/2026 ([about.fb.com 06/2026](https://about.fb.com/news/2026/06/meta-launching-new-features-500-million-monthly-threads-users/)). Bài công bố không nói rõ bài trong community có tới được người ngoài community hay không.
- **Phân phối chéo sang Instagram/Facebook**: "eligible posts of yours can be recommended to people on Instagram and Facebook" ([about.fb.com 03/2025](https://about.fb.com/news/2025/03/new-threads-features-more-personalized-experience-you-control/)).
- **Search**: Account Status có mục "What can't be recommended". Theo nguồn thứ cấp, nội dung bị xếp không đủ điều kiện gợi ý sẽ bị loại khỏi For You, Search, topic tag và suggested replies **[chưa xác minh từ nguồn chính]** ([postory.io](https://postory.io/blog/threads-shadowban)).

### 1.2 Tín hiệu theo system card (cập nhật 07/03/2025 – ⚠️ đã 19 tháng, có thể không còn khớp)

Card liệt kê 9 dự đoán ([nguồn](https://transparency.meta.com/features/explaining-ranking/ig-threads-feed/)):

1. Khả năng bạn like bài.
2. Khả năng bạn lướt qua thay vì tương tác.
3. Khả năng bạn bấm vào profile tác giả rồi bấm sang một bài khác.
4. Khả năng bạn bấm vào bài rồi bấm sang một bài khác.
5. Khả năng bạn **viết reply**.
6. Khả năng bạn bấm vào bài.
7. Khả năng bạn **bấm vào profile tác giả**.
8. Khả năng bạn bấm vào bài rồi like.
9. Thời gian bạn xem trang permalink của bài.

Các tín hiệu đầu vào đáng chú ý:

- "Engagement of descendants within level 2": tương tác trong cây reply, tới hai tầng, trong 1–6 giờ gần nhất.
- "Engagement rate of media (likes to impressions ratio)".
- "Engagement of media by viewer country over past week": **quốc gia của người xem là một tín hiệu**.

Card **không** nói dự đoán nào nặng nhất.

**Diễn giải [suy luận]:** 4/9 dự đoán liên quan tới việc bấm vào bài hoặc profile, và nhiều tín hiệu đo độ sôi nổi của cây reply. Nghĩa là Threads thưởng cho bài **tạo ra hội thoại** và **khiến người xem tò mò về tác giả**. Cả hai đều có lợi cho personal brand.

### 1.3 Reply và hội thoại

- @threads (21/10/2024): "Reply! Replies account for almost half of the views on Threads". Cùng bài khuyên "Post often — at least 2 to 5 times a week", theo trend qua topic tag và "Share your unique POV" ([post](https://www.threads.com/@threads/post/DBZbkQHxjeD)).
- Trang creators.instagram.com/threads nhắc lại các ý trên và nói thêm ([link](https://creators.instagram.com/threads), không ghi ngày):
  - "Higher post frequency is linked to higher views per post".
  - "Video, photo, and carousel posts that also include text have performed significantly better".
  - "Top-performing creators post original content to Threads".
  - Tránh "clickbait, engagement bait, or which promotes a contest or giveaway".
- Mosseri, phỏng vấn với Platformer (07/2024): "the sum of all your replies is about as valuable as the sum of … all your posts" **[trích lại qua nguồn thứ cấp, chưa đọc bản gốc]** ([Platformer](https://www.platformer.news/threads-175-million-users-adam-mosseri-interview/)).
- Công cụ quản lý reply (03/2025): reply approvals, giới hạn reply/quote cho follower ([about.fb.com](https://about.fb.com/news/2025/03/new-threads-features-more-personalized-experience-you-control/)).

### 1.4 Link

- Mosseri, 21/11/2024: "We don't downrank links, but we don't place much value on p(click) and people don't like and comment on links much." ([post](https://www.threads.com/@mosseri/post/DCpQG0hz0m8))
- Mosseri, 03/06/2025: "we've been working on making sure links are ranked properly. Links have been working much better for more than a month now." ([post](https://www.threads.com/@mosseri/post/DKdzO8fumTq))
- Threads không có giới hạn số link kiểu Facebook (§3.5). Một số nguồn thứ cấp nói năm 2026 Threads "boost" bài có URL. Không tìm thấy nguồn chính cho chuyện đó **[chưa xác minh từ nguồn chính]** ([posteverywhere.ai](https://posteverywhere.ai/blog/how-the-threads-algorithm-works)).

### 1.5 Đủ điều kiện gợi ý

- Account Status có trên Threads từ 08/05/2025 ([@threads](https://www.threads.com/@threads/post/DJZ0SGLyoOH)). Ảnh chụp kèm bài cho thấy các mục "Removed content", "What can't be recommended", "Content lowered in feed", "Features you can't use". Đường dẫn: Settings > Account > Account status ([@threads](https://www.threads.com/@threads/post/DJZ0SwGyxZw)).

### 1.6 Thay đổi 2024–2026

| Ngày | Thay đổi | Nguồn |
|---|---|---|
| 21/11/2024 | "rebalancing ranking to prioritize content from people you follow … less recommended content". Creator sẽ thấy unconnected reach giảm | [Mosseri](https://www.threads.com/@mosseri/post/DCo4WiOvCB0) |
| 03/2025 | Custom feed làm mặc định; topic hiện ở nhiều chỗ hơn; bài Threads được gợi ý sang IG/FB; insights theo nguồn khám phá | [about.fb.com](https://about.fb.com/news/2025/03/new-threads-features-more-personalized-experience-you-control/) |
| 05/2025 | Link "ranked properly" | [Mosseri](https://www.threads.com/@mosseri/post/DKdzO8fumTq) |
| 06/2025 | Feed fediverse riêng; tìm được user fediverse | [about.fb.com](https://about.fb.com/news/2025/06/its-now-easier-see-more-fediverse-content-threads/), [@threads](https://www.threads.com/@threads/post/DLAnw0cAorN) |
| 01/07/2025 | DM, "Highlighted perspectives". "More than a third of daily Threads users with connections follow mostly different accounts on Threads than on Instagram" | [about.fb.com](https://about.fb.com/news/2025/07/introducing-messaging-highlighted-perspectives-threads/) |
| 02/10/2025 | Communities | [about.fb.com](https://about.fb.com/news/2025/10/introducing-threads-communities-find-your-people/) |
| 11/02/2026 | "Dear Algo": viết bài công khai để chỉnh feed trong 3 ngày (US, NZ, AU, UK) | [about.fb.com](https://about.fb.com/news/2026/02/threads-dear-algo/) |
| 16/06/2026 | 500 triệu MAU; Communities ra khỏi beta; "Your Algo" (chỉnh riêng tư trong 1/3/7 ngày; US, CA, UK, AU, NZ); tag ngôn ngữ bản địa ở JP/KR/TW | [about.fb.com](https://about.fb.com/news/2026/06/meta-launching-new-features-500-million-monthly-threads-users/) |

**Fediverse:** Threads có thể chia sẻ bài ra fediverse. Không có nguồn chính nào nói việc bật fediverse ảnh hưởng tới xếp hạng trong Threads. Theo nguồn thứ cấp, tới cuối 2025 người dùng Threads vẫn chưa reply lại được bài từ fediverse **[chưa xác minh từ nguồn chính]** ([mjtsai.com 12/2025](https://mjtsai.com/blog/2025/12/01/what-happened-with-threads-and-the-fediverse/)).

---

## 2. Instagram

### 2.1 Bề mặt tới non-follower

- **Explore**, **Reels (Reels chaining)**, **Feed Recommendations** (bài gợi ý trong Feed), Search, Suggested Accounts ([danh mục card](https://transparency.meta.com/features/explaining-ranking/)).
- **Trial Reels** (ra mắt 10/12/2024): reel chỉ hiển thị cho non-follower. Khoảng 24 giờ sau xem được số liệu. Có thể tự động chia sẻ cho follower nếu reel đạt ngưỡng lượt xem trong 72 giờ ([creators.instagram.com](https://creators.instagram.com/blog/instagram-trial-reels)). Lịch đăng Trial Reels có từ 04/2026 **[nguồn thứ cấp]** ([SMT](https://www.socialmediatoday.com/news/instagram-allows-creators-to-schedule-trial-reels/816549/)). Điều kiện 200 follower (tài khoản professional) hoặc 1.000 (tài khoản cá nhân) cũng **[chưa xác minh từ nguồn chính]** ([creatorflow.so](https://creatorflow.so/blog/instagram-trial-reels/)).

### 2.2 Dự đoán theo từng system card

| Bề mặt (ngày cập nhật card) | Dự đoán được liệt kê |
|---|---|
| **Explore** (22/06/2026) | Follow tác giả; xem bài ≥5 giây; xem ≥95% video; bấm "Not Interested"; comment; like; reshare; bấm và tương tác; save. Pipeline: retrieval tới 1.500 ứng viên → early ranking còn 100 → late ranking đa nhiệm ([card](https://transparency.meta.com/features/explaining-ranking/ig-explore/)) |
| **Reels chaining** (11/11/2025) | Dùng audio của reel; xem <3 giây; bấm "Interested"; comment; xem >95%; reshare; chia sẻ ra ngoài Instagram; **follow tác giả**… ([card](https://transparency.meta.com/features/explaining-ranking/ig-reels-chaining/)) |
| **Feed Recommendations** (29/06/2026) | "How informative a post is"; bấm "Interested"; thời gian xem Reels sau khi bấm; **bấm vào profile tác giả**; bấm để comment; bỏ qua bài; mở Reel toàn màn hình; bật tiếng… ([card](https://transparency.meta.com/features/explaining-ranking/ig-feed-recommendations/)) |
| **Feed – connected** (29/06/2026) | Vuốt carousel; "enjoyment" (like, comment, reshare, thời gian); bỏ qua; **thời gian ở profile tác giả sau khi xem**; chia sẻ qua DM; thời gian xem; bật tiếng ([card](https://transparency.meta.com/features/explaining-ranking/ig-feed/)) |

Không card nào nêu trọng số.

**Diễn giải [suy luận]:** follow tác giả, vào profile, chia sẻ/DM và xem trọn video lặp lại ở nhiều bề mặt. Đây là những hành động một personal brand cần.

### 2.3 Tín hiệu Instagram nói là quan trọng nhất

- Mosseri, 01/2025: "The top three signals that matter most for ranking are watch time, likes and sends". "Likes are slightly more important for connected content, and sends are slightly more important for unconnected content." Ông khuyên theo dõi average watch time, likes per reach và sends per reach ([SMT 22/01/2025](https://www.socialmediatoday.com/news/instagram-shares-algorithm-insights-2025/738034/), [video gốc](https://www.instagram.com/p/DFFyRp-pINJ/) – trích qua SMT, chưa xem video).
- Mosseri, 01/05/2024: thay đổi xếp hạng gợi ý để "give smaller creators a better chance of breaking through" ([post](https://www.threads.com/@mosseri/post/C6bd8O3xyCj)).
- Bài chính thức "Instagram Ranking Explained" (31/05/2023 – ⚠️ cũ nhưng là bản giải thích dài nhất Instagram từng công bố):
  - Reels dự đoán khả năng reshare, xem hết, like, vào trang audio.
  - Explore dựa trên việc người khác like, comment, share, save "how many and how quickly".
  - Reels bị "less visible" nếu "low-resolution or watermarked … muted or contain borders … majority text, or … already been posted" ([link](https://about.instagram.com/blog/announcements/instagram-ranking-explained)).

### 2.4 Tính nguyên bản (30/04/2026)

- "Original content is work that you wholly created or reflects your unique perspective" ([creators.instagram.com](https://creators.instagram.com/blog/rewarding-original-creators-on-instagram)).
- Aggregator là tài khoản chủ yếu đăng nội dung không nguyên bản trong ảnh, carousel và reels. Họ mất quyền được gợi ý. Có thể lấy lại khi phần lớn bài trong **30 ngày gần nhất (rolling)** là nguyên bản.
- "This update does not impact how we show people content from accounts they follow."
- Số liệu: "75% of recommendations in the US now coming from original posts".
- Instagram khuyên dùng collab, remix, nhãn paid partnership và kiểm tra Account Status.

### 2.5 Định dạng và độ dài

- Reels tối đa 3 phút (Mosseri, 18/01/2025): "We still are focused on short form video over long form" ([post](https://www.threads.com/@mosseri/post/DE-efFqyStv)).
- ⚠️ Có thông tin cũ rằng Instagram chỉ gợi ý reel ≤90 giây tới người chưa follow. Không tìm thấy nguồn chính, và thông tin này có thể đã lỗi thời sau khi nâng giới hạn lên 3 phút **[chưa xác minh từ nguồn chính]** ([medianama 01/2025](https://www.medianama.com/2025/01/223-instagram-reels-3-minutes-us-tiktok-ban/)).
- Carousel: card Feed có dự đoán "swiping through carousels". Theo nguồn thứ cấp, Mosseri (10/2024) nói carousel chưa được vuốt có thể hiện lại ở slide khác **[chưa xác minh từ nguồn chính]** ([alxo.ae](https://alxo.ae/whats-on/instagram-ranking-2026-signals-on-the-record/)).

### 2.6 Bị hạ hoặc loại

- Watermark của app khác, độ phân giải thấp, reel tắt tiếng, có viền, chủ yếu là chữ, đã đăng trước đó ([Ranking Explained 2023](https://about.instagram.com/blog/announcements/instagram-ranking-explained)).
- Nội dung vi phạm Recommendation Guidelines (cùng nguồn).
- Aggregator/unoriginal (§2.4).
- Professional account kiểm tra được trong Account Status xem nội dung có đủ điều kiện gợi ý trong Explore, Reels và Feed Recommendations hay không ([TechCrunch 12/2022](https://techcrunch.com/2022/12/07/instagrams-new-transparency-tools-will-tell-you-if-your-content-is-ineligible-to-be-recommended/) – trang Help Center [Recommendations on Instagram](https://help.instagram.com/313829416281232) không tải được khi viết bài).

### 2.7 Thay đổi 2025–2026

| Ngày | Thay đổi | Nguồn |
|---|---|---|
| 18/01/2025 | Reels tối đa 3 phút | [Mosseri](https://www.threads.com/@mosseri/post/DE-efFqyStv) |
| 22/01/2025 | Watch time, likes, sends là ba tín hiệu hàng đầu | [SMT](https://www.socialmediatoday.com/news/instagram-shares-algorithm-insights-2025/738034/) |
| 10/12/2025 → 15/04/2026 → 10/06/2026 | "Your Algorithm": người dùng xem và chỉnh chủ đề gợi ý. Bắt đầu ở Reels, mở sang Explore, rồi Feed. Một hệ thống chung, gợi ý có nhãn kiểu "From Running". **Chỉ có tiếng Anh** (toàn cầu từ 06/2026) | [about.instagram.com](https://about.instagram.com/blog/announcements/reels-algorithm-control) |
| 30/04/2026 | Originality mở rộng sang ảnh và carousel | [creators.instagram.com](https://creators.instagram.com/blog/rewarding-original-creators-on-instagram) |
| 04/2026 | Insights thêm share rate và skip rate | **[chưa xác minh từ nguồn chính]** ([inro.social](https://www.inro.social/blog/instagram-insights-update)) |

---

## 3. Facebook

### 3.1 Bề mặt tới non-follower

- **Feed Recommendations**: gom "a portion of the public content available on Facebook, which may include photos, videos and links" ([card, 24/06/2026](https://transparency.meta.com/features/explaining-ranking/fb-feed-recommendations/)).
- **Reels**: từ 17/06/2025 mọi video Facebook đều là Reels. "Reels on Facebook will also not have any length or format restrictions". Meta khẳng định thay đổi này "won't change what videos we recommend to you" ([about.fb.com](https://about.fb.com/news/2025/06/making-it-easier-create-videos-facebook/)).
- **Search** dạng lưới toàn màn hình (12/2025) ([about.fb.com 09/12/2025](https://about.fb.com/news/2025/12/making-it-easier-to-create-discover-and-share-content-on-facebook/)).
- Chiều ngược lại: tab **Friends** (27/03/2025, US và Canada) "will now show content just from your Facebook friends, no recommended content" ([about.fb.com](https://about.fb.com/news/2025/03/bringing-magic-of-friends-back-to-facebook/)).

### 3.2 Dự đoán

- **Feed Recommendations** (24/06/2026): thời gian nhìn bài; **share**; bấm "X"; xem lâu hơn; tham gia group; like; **share riêng sang nền tảng khác**; **chụp màn hình**; RSVP; bấm "Show more" ([card](https://transparency.meta.com/features/explaining-ranking/fb-feed-recommendations/)).
- **Reels** (19/06/2026): thời gian xem; xem tiếp sau khi mở toàn màn hình; bấm mở toàn màn hình; like; xem hay bỏ qua; xem hết; yêu cầu bớt nội dung tương tự; ẩn; bắt đầu phiên xem; chuyển sang reel kế tiếp ([card](https://transparency.meta.com/features/explaining-ranking/fb-reels/)).
- **Feed – connected** (06/10/2026): có dự đoán "How likely the user will click on a link that takes the user to a website outside of Facebook and predicts how long the user will stay on that website". Với người đã kết nối, click link là một **kết quả được mô hình hóa**, không phải tín hiệu tiêu cực. Feed connected cũng có "strong intent to interact with their friends/group" ([card](https://transparency.meta.com/features/explaining-ranking/fb-feed/)).
- **Khảo sát "True Interest" cho Reels** (engineering, 14/01/2026): người xem chấm độ hợp sở thích 1–5. Mô hình UTIS được đưa vào cả retrieval lẫn late ranking, "increased niche content delivery while reducing generic, popularity-based recommendations" ([engineering.fb.com](https://engineering.fb.com/2026/01/14/ml-applications/adapting-the-facebook-reels-recsys-ai-model-based-on-user-feedback/)).
- **Friend bubbles** (18/03/2026): Reels mà bạn bè đã like được đưa vào candidate sourcing, xếp hạng theo độ thân thiết giữa người xem và bạn bè ([engineering.fb.com](https://engineering.fb.com/2026/03/18/ml-applications/friend-bubbles-enhancing-social-discovery-on-facebook-reels/)).

### 3.3 Bị hạ

- CDG hiện hành: Clickbait Links, Engagement Bait, Fact-Checked Misinformation, nội dung có khả năng vi phạm Community Standards ([02/07/2025](https://transparency.meta.com/features/approach-to-ranking/types-of-content-we-demote/)). Xem thêm phần xóa 14 guideline ở §0.

### 3.4 Tính nguyên bản

- 14/07/2025: Facebook giảm phân phối video trùng lặp. Tài khoản tái sử dụng nội dung lặp lại mất monetization tạm thời và bị giảm phân phối. Meta khuyên tránh watermark của bên thứ ba, viết caption chất lượng với **dưới 5 hashtag**, và ghép clip không được tính là cải biến có ý nghĩa ([creators.facebook.com](https://creators.facebook.com/blog/combating-unoriginal-content)).
- 13/03/2026: nội dung nguyên bản là "Content filmed or produced directly by a creator or owner of a Profile or Page". Remix chỉ được tính nếu thêm "fresh information, analysis, or substantial improvements" kèm sự hiện diện trên màn hình. Nội dung không nguyên bản bị "deprioritized"; nếu tiếp tục, tài khoản có thể bị "non-recommendable and demonetized". Lượt xem và thời gian xem Reels nguyên bản "approximately doubled" ở nửa cuối 2025 so với cùng kỳ 2024 ([about.fb.com](https://about.fb.com/news/2026/03/rewarding-original-creators-on-facebook/)).

### 3.5 Link trên Facebook

- Không có nguồn chính nào nói Facebook hạ hạng bài chỉ vì có link. Card Feed mô hình hóa việc click link (§3.2), và CDG chỉ hạ **clickbait** link.
- Giới hạn **số lượng** link:
  - Từ 16/12/2025 Meta thử nghiệm giới hạn Pages và profile chế độ chuyên nghiệp không có Meta Verified ở **2 bài có link mỗi tháng**. Link trong comment trên bài của chính mình và link tới nền tảng của Meta không bị giới hạn. Phát ngôn Meta gửi Matt Navarra: "a limited test…" **[nguồn thứ cấp: Mari Smith]** ([marismith.com](https://www.marismith.com/facebook-tests-limiting-link-posts-to-just-2-per-month/)).
  - 15/09/2026, Meta One: gói Advanced (49,99 USD/tháng) có "links in organic posts and reels" ([about.fb.com](https://about.fb.com/news/2026/09/introducing-meta-one-subscription-service-more-features-ai/)). Theo nguồn thứ cấp, giới hạn đã mở rộng sang nhiều Page hơn từ 09/2026, và quá hạn mức thì link hiển thị dạng chữ, không có preview **[chưa xác minh từ nguồn chính]** ([relevantaudience.com](https://www.relevantaudience.com/meta/facebook-pages-charged-for-link-posts-meta-one/), [builtbyfoundry.io](https://www.builtbyfoundry.io/blog/how-many-links-can-you-post-facebook-2026-2-free-4-paid)).

### 3.6 Profile cá nhân, chế độ chuyên nghiệp hay Page

- Chế độ chuyên nghiệp trên profile cá nhân: có insights; nội dung "eligible for increased discovery"; "Who can follow me" chuyển sang công khai; giữ nguyên bạn bè trên cùng một profile ([Help Center](https://www.facebook.com/help/2978174759158711)).
- Meta One có "bold follow button on your reels" và "follow invitations sent automatically to people who engage with your content" ([about.fb.com 15/09/2026](https://about.fb.com/news/2026/09/introducing-meta-one-subscription-service-more-features-ai/)).

---

## 4. Ngôn ngữ và khu vực

| Điểm | Nội dung | Nguồn |
|---|---|---|
| Quốc gia là tín hiệu trên Threads | Card Threads Feed có tín hiệu "Engagement of media by viewer country over past week" | [card 07/03/2025](https://transparency.meta.com/features/explaining-ranking/ig-threads-feed/) |
| Dịch chữ tự động | Caption và bình luận trên Instagram được dịch tự động theo ngôn ngữ của bài và cài đặt ngôn ngữ của người xem ("See Translation") | [Help Center](https://help.instagram.com/249420315434256) (lấy qua kết quả tìm kiếm). Threads có tính năng dịch tương tự từ 2023 **[nguồn thứ cấp]** ([fonearena](https://www.fonearena.com/blog/400182/threads-following-feed-translation.html)) |
| AI dịch và lồng tiếng Reels | Meta AI dịch giọng, giữ âm sắc, có lip-sync tùy chọn, gắn nhãn "Translated with Meta AI". Facebook yêu cầu ≥1.000 follower. Instagram mở cho mọi tài khoản công khai ở nước có Meta AI. **Tiếng Việt được thêm cho Facebook ngày 04/06/2026.** Instagram thêm Pháp, Đức, Ý, Nhật, Hàn ngày 14/07/2026 | [about.fb.com 09/10/2025, cập nhật 04/06/2026 và 14/07/2026](https://about.fb.com/news/2025/10/discover-reels-around-world-meta-ai-translation/); [about.fb.com 06/2026](https://about.fb.com/news/2026/06/creator-assistant-more-languages-for-ai-translations-on-facebook/): "Over half a billion users on Facebook are now watching AI-translated videos weekly" |
| ⚠️ Chưa rõ | Instagram đã hỗ trợ tiếng Việt cho AI translation chưa, và Meta AI có hoạt động ở Việt Nam không: bài cập nhật chỉ ghi tiếng Việt cho Facebook. Cần kiểm tra trực tiếp trong app | — |
| Chỉnh feed theo chủ đề | Instagram "Your Algorithm" chỉ có bằng tiếng Anh. Threads "Dear Algo" và "Your Algo" chỉ có ở US/CA/UK/AU/NZ. Người xem Việt Nam có thể chưa có các công cụ này | [about.instagram.com](https://about.instagram.com/blog/announcements/reels-algorithm-control), [about.fb.com](https://about.fb.com/news/2026/06/meta-launching-new-features-500-million-monthly-threads-users/) |
| Một hay hai tài khoản | **Không tìm thấy nguồn chính** nào của Meta khuyên tách tài khoản theo ngôn ngữ. Theo nguồn thứ cấp, Mosseri (Q&A 07/2026) khuyên bật dịch Reels rồi xem insights theo quốc gia để quyết định đầu tư vào đâu **[chưa xác minh từ nguồn chính]** | [kontentino.com](https://www.kontentino.com/q-and-a/instagram-global-reach/) |

---

## 5. Chưa xác minh hoặc có thể lỗi thời

- "Sends nặng gấp 3–5 lần like": không có nguồn chính. Mosseri chỉ nói "slightly more important" ([SMT](https://www.socialmediatoday.com/news/instagram-shares-algorithm-insights-2025/738034/)).
- "Two-way conversation became a ranked signal in March 2026", "Mosseri công bố trọng số Threads": không có nguồn chính.
- Card Threads Feed (03/2025) và bài Instagram Ranking Explained (05/2023) đã cũ. Card Instagram Reels chaining được cập nhật 11/2025.
- Thông tin về giới hạn link trên Facebook nằm chủ yếu ở nguồn thứ cấp. Meta chưa có trang Help Center hoặc newsroom mô tả trực tiếp giới hạn 2 link.
- Không tìm thấy công bố chính thức nào của Meta về độ dài video tối ưu (ví dụ "15–30 giây"). Mọi con số kiểu đó đều từ bên thứ ba.

---

## 6. Hàm ý cho mate (mục tiêu: personal brand của người làm mate)

### 6.1 Áp dụng cho cả hai nhánh

1. **Tối ưu cho follow và lượt vào profile, không tối ưu cho click link.** Các system card có dự đoán follow tác giả (Explore, Reels), vào hoặc ở lại profile (IG Feed, IG Feed Recs, Threads). Cụ thể:
   - Bio, ảnh đại diện và bài ghim phải trả lời ngay "người này là ai, làm gì với AI coding agent".
   - Dùng cùng một tên và cùng một gương mặt trên cả ba nền tảng **[suy luận]**.
2. **Chỉ đăng nội dung nguyên bản.** Quay màn hình terminal của chính mình, có giọng nói hoặc mặt của mình, và phân tích của mình. Không đăng lại clip hoặc meme AI của người khác. Không dùng video có watermark của app khác ([IG 2023](https://about.instagram.com/blog/announcements/instagram-ranking-explained), [IG 2026](https://creators.instagram.com/blog/rewarding-original-creators-on-instagram), [FB 2026](https://about.fb.com/news/2026/03/rewarding-original-creators-on-facebook/)). Một tài khoản chỉ "tổng hợp tin AI" sẽ có nguy cơ bị xếp là aggregator.
3. **Nội dung đáng gửi cho bạn bè.** Sends và share là dự đoán ở IG Explore, Reels, Feed và FB Feed Recs, gồm cả "privately share" và "screenshot". Mosseri nói sends nặng hơn với non-follower. Ví dụ định dạng **[suy luận]**: "một lệnh/một mẹo tiết kiệm X phút", so sánh trước/sau, sơ đồ kiến trúc nhiều agent mà người xem muốn lưu hoặc gửi cho đồng nghiệp.
4. **Giữ người xem trong 3 giây đầu và ở lại tới cuối.** IG Reels có dự đoán "<3 giây" và ">95%". FB Reels có dự đoán xem hay bỏ qua và xem hết. Demo terminal nên vào kết quả ngay: hiện màn hình nhiều agent đang chạy trước, giải thích sau **[suy luận từ các dự đoán]**. Dùng font terminal lớn và có phụ đề, vì reel tắt tiếng hoặc toàn chữ bị đánh giá thấp ([2023](https://about.instagram.com/blog/announcements/instagram-ranking-explained)).
5. **Không dùng engagement bait** ("comment 'mate' để nhận link"), clickbait hay giveaway. Cả CDG lẫn hướng dẫn Threads đều nêu tên ([CDG](https://transparency.meta.com/features/approach-to-ranking/types-of-content-we-demote/), [Threads](https://creators.instagram.com/threads)).
6. **Kiểm tra Account Status định kỳ** trên Instagram và Threads, mục "What can't be recommended".
7. **Link GitHub**:
   - Threads: để link trong bài khi cần, vì không bị hạ hạng ([Mosseri 2024](https://www.threads.com/@mosseri/post/DCpQG0hz0m8), [2025](https://www.threads.com/@mosseri/post/DKdzO8fumTq)). Nhưng link không kéo tương tác, nên đừng biến nó thành mục tiêu của bài.
   - Facebook (profile chuyên nghiệp hoặc Page): đặt link vào comment, vì link trong comment trên bài của mình được ngoại lệ trong đợt thử nghiệm giới hạn **[nguồn thứ cấp]**.

### 6.2 Nhánh tiếng Việt (dev Việt)

- **Threads là kênh chính [suy luận]**, vì phân phối dựa nhiều vào hội thoại và có tín hiệu quốc gia của người xem. Cách làm:
  - Đăng 2–5 bài/tuần, mỗi bài một topic tag ([@threads](https://www.threads.com/@threads/post/DBZbkQHxjeD)).
  - Dành nhiều thời gian reply hơn đăng bài, vào các thread về AI và lập trình của người Việt. Reply chiếm gần một nửa lượt xem.
  - Kèm ảnh chụp terminal hoặc video ngắn, vì bài có media kèm chữ "performed significantly better" ([creators.instagram.com/threads](https://creators.instagram.com/threads)).
  - Chọn tag tiếng Việt hay tiếng Anh: chưa có dữ liệu. Nên thử cả hai và so insights **[suy luận]**. Kiểm tra xem có community lập trình hoặc AI trên Threads để tham gia không (Meta không công bố danh sách; tag có dấu ba chấm là có community).
- **Facebook**:
  - Dùng **profile cá nhân bật chế độ chuyên nghiệp** thay vì lập Page cho mate. Lý do [suy luận dựa trên Help Center]: chế độ chuyên nghiệp giữ được mạng bạn bè sẵn có (Friends tab, friend bubbles đều ưu tiên bạn bè), mở follow công khai, và cho nội dung "eligible for increased discovery".
  - Reels tiếng Việt quay thật. Bài trong group dev Việt là kênh có thể hiệu quả **[suy luận]**: card Feed có dự đoán "deep/strong intent to engage with posts from a group", và Feed Recs có dự đoán "join a group".
  - Facebook đã có AI dịch Reels sang và từ tiếng Việt (từ 04/06/2026, cần ≥1.000 follower) ([about.fb.com](https://about.fb.com/news/2025/10/discover-reels-around-world-meta-ai-translation/)). Khi đủ điều kiện, một Reel tiếng Việt có thể được dịch sang tiếng Anh. Hiệu quả thực tế chưa được kiểm chứng.
- **Instagram**: Reels demo 30–90 giây hoặc carousel "cách làm". Dùng Trial Reels để thử hook với non-follower mà không làm phiền follower ([creators.instagram.com](https://creators.instagram.com/blog/instagram-trial-reels)). Instagram có thể là kênh yếu nhất với dev Việt **[suy luận, chưa có dữ liệu]**.

### 6.3 Nhánh tiếng Anh (dev quốc tế)

- **Một hay hai tài khoản?** Meta không có hướng dẫn chính thức. Đề xuất **[suy luận]**:
  - **Instagram và Facebook: một tài khoản, đăng Reels bằng tiếng Việt hoặc tiếng Anh rồi bật AI translation.** Personal brand gắn với một người, và follow, profile visit dồn về một chỗ. Meta tự dịch caption và có thể lồng tiếng Reels; Mosseri (nguồn thứ cấp) khuyên bật dịch rồi đọc insights theo quốc gia.
  - **Threads: cân nhắc tách tài khoản tiếng Anh nếu tỷ lệ follower Việt quá áp đảo.** Threads xếp hạng theo hội thoại, và có tín hiệu quốc gia của người xem. Một tài khoản trộn hai ngôn ngữ có thể khiến người xem mỗi nhóm thấy nửa số bài không hợp, dẫn tới dự đoán "scroll past" cao **[suy luận, không có nguồn]**. Cách kiểm chứng: chạy một tài khoản trong 4–6 tuần, xem insights theo quốc gia, rồi mới quyết định.
- **Threads tiếng Anh**: dùng topic tag tiếng Anh phổ biến về AI coding hoặc developer tools, reply trong thread của các dev tool và AI builder nổi tiếng. Viết quan điểm riêng ("unique POV") về cách điều phối nhiều agent, vì đó là mảng mate có chuyên môn thật.
- **Instagram tiếng Anh**: công cụ quan trọng nhất là Trial Reels. Dùng nó để thử nội dung tiếng Anh với non-follower, theo dõi sends per reach và watch time.
- **Facebook tiếng Anh**: kém ưu tiên nhất **[suy luận]**. Feed Facebook đang nghiêng về bạn bè (Friends tab, friend bubbles), mà mạng bạn bè của tác giả chủ yếu là người Việt. Dùng Reels kèm AI translation thay vì làm nội dung tiếng Anh riêng cho Facebook.
- **Không cross-post nguyên văn từ Instagram sang Threads**: "Top-performing creators post original content to Threads", và hơn 1/3 người dùng Threads hằng ngày follow phần lớn những tài khoản khác với trên Instagram ([creators.instagram.com/threads](https://creators.instagram.com/threads), [about.fb.com 07/2025](https://about.fb.com/news/2025/07/introducing-messaging-highlighted-perspectives-threads/)).
