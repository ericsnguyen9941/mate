# Thuật toán For You của TikTok (tính đến 2026-10-09)

Ngày: 2026-10-09. Mục đích: chuẩn bị kênh video ngắn cho mate (CLI Go mã nguồn mở điều phối AI coding agent). Mục tiêu là **personal brand của người làm mate**, không phải GitHub star. Người làm không lên hình. Threads là hub (viral-first, xem [social-algorithms-2026-10-09.md](social-algorithms-2026-10-09.md)). TikTok và Facebook Reels chạy video "bài học" (value-first) làm bằng HyperFrames (motion graphics HTML→MP4, faceless), có bản tiếng Việt và tiếng Anh.

File này chỉ nói về TikTok. Phần Threads/Instagram/Facebook nằm ở file anh em ở trên.

Quy ước nguồn:

- Nguồn chính: newsroom.tiktok.com, Help Center (support.tiktok.com), Transparency/Safety Center (tiktok.com/transparency, tiktok.com/safety), Community Guidelines (tiktok.com/community-guidelines), Creator Academy (tiktok.com/creator-academy), Help Center quảng cáo (ads.tiktok.com/help).
- Help Center và Creator Academy render bằng JavaScript. Nội dung được đọc bằng trình duyệt headless ngày 2026-10-09. Ngày ghi trong ngoặc là "Last updated" hiển thị trên trang; trang nào không ghi ngày thì ghi "không ghi ngày".
- Điều gì chỉ thấy ở nguồn thứ cấp thì ghi **[chưa xác minh từ nguồn chính]**. Suy luận của người viết ghi **[suy luận]**.
- TikTok **không công bố trọng số bằng số** cho tín hiệu nào. Các con số kiểu "phải đạt 70% completion" đều không có nguồn chính (xem §10).

---

## Tóm tắt nhanh

1. For You xếp hạng bằng **điểm dự đoán cho từng người xem**: khả năng like, share, comment, skip, xem hết, follow tác giả, vào profile, xem trong một khoảng thời gian ([Transparency Center](https://www.tiktok.com/transparency/en/recommendation-system), không ghi ngày). Help Center nói tương tác của người dùng, gồm cả thời gian xem, "generally weighted more heavily than others" ([Help Center](https://support.tiktok.com/en/using-tiktok/exploring-videos/how-tiktok-recommends-content), không ghi ngày).
2. Số follower **không phải tín hiệu trực tiếp** ([newsroom 18/06/2020](https://newsroom.tiktok.com/en-us/how-tiktok-recommends-videos-for-you), ⚠️ đã 6 năm). Ở bước cuối có **luật bắt buộc trộn nội dung của creator cùng khu vực** với người xem ([Transparency](https://www.tiktok.com/transparency/en/recommendation-system); [Creator Academy 26/09/2026](https://www.tiktok.com/creator-academy/en/article/guidelines-recommendation-system-intro)).
3. **Nội dung không nguyên bản hoặc chất lượng thấp sẽ không lên For You.** Community Guidelines bản 08/2026 (hiệu lực 24/09/2026) nêu rõ hai loại: "Reused or unoriginal content posted without creative edits" và "Low-quality or minimally edited content" ([CG](https://www.tiktok.com/community-guidelines/en/integrity-authenticity)).
4. **AIGC**: chỉ bắt buộc gắn nhãn khi nội dung AI trông như thật. Motion graphics, phong cách nghệ thuật và TTS chung chung không cần gắn nhãn. TikTok nói nhãn AI "won't affect the distribution" ([Creator Academy 26/09/2026](https://www.tiktok.com/creator-academy/en/article/ai-generated-content-label)). Tuy vậy, từ 11/2025 người xem có thanh trượt để **giảm AIGC** trong feed của họ ([newsroom 19/11/2025](https://newsroom.tiktok.com/more-ways-to-spot-shape-and-understand-ai-content?lang=en)).
5. **Việt Nam bị loại khỏi STEM feed và Creator Rewards Program.** STEM feed "available globally, except … Vietnam" ([Creator Academy 19/08/2026](https://www.tiktok.com/creator-academy/en/article/what-is-STEM-feed)). CRP chỉ mở ở 8 nước, không có Việt Nam ([Creator Academy 11/09/2026](https://www.tiktok.com/creator-academy/en/article/creator-rewards-program)). Với mục tiêu personal brand thì CRP không quan trọng. STEM feed thì đáng giá cho nhánh tiếng Anh, nhưng bản tiếng Anh vẫn phải đăng từ Việt Nam.

---

## 0. Bản đồ nguồn chính

| Nguồn | Ngày | Nói gì |
|---|---|---|
| [Newsroom "How TikTok recommends videos #ForYou"](https://newsroom.tiktok.com/en-us/how-tiktok-recommends-videos-for-you) | 18/06/2020 ⚠️ | Bài gốc: 3 nhóm tín hiệu, trọng số mạnh/yếu, follower không phải tín hiệu trực tiếp |
| [Help Center "How TikTok recommends content"](https://support.tiktok.com/en/using-tiktok/exploring-videos/how-tiktok-recommends-content) | không ghi ngày (đọc 09/10/2026) | Tín hiệu cho từng bề mặt: For You, Following, Friends, LIVE, Search, Comments, Account recs |
| [Transparency Center "Introduction to the TikTok recommendation system"](https://www.tiktok.com/transparency/en/recommendation-system) | không ghi ngày | Pipeline 5 bước, phụ lục liệt kê tín hiệu và dự đoán |
| [Safety Center "Making your feed For You"](https://www.tiktok.com/safety/en/tools-and-guides/making-your-feed-for-you) | cập nhật 30/03/2026 (metadata `updatedAt`) | Yếu tố chất lượng: ưu tiên creator địa phương, nội dung mới, độ dài, sound |
| [Community Guidelines](https://www.tiktok.com/community-guidelines/en/fyf-standards) | phát hành 25/08/2026, hiệu lực 24/09/2026 | Tiêu chuẩn đủ điều kiện For You (FYF), AIGC, originality |
| [Creator Academy](https://www.tiktok.com/creator-academy/en/article/guidelines-recommendation-system-intro) | 04–09/2026 | Hướng dẫn cho creator: watch time, originality, search, STEM, CRP |

---

## 1. For You xếp hạng thế nào

### 1.1 Pipeline (Transparency Center, không ghi ngày)

Theo [Transparency Center](https://www.tiktok.com/transparency/en/recommendation-system), pipeline có 5 bước:

1. **Selecting**: lấy một tập lớn video đã qua kiểm duyệt và "eligible for recommendation to the For You feed".
2. **Prediction**: ước lượng xác suất like, share, comment, skip… rồi "combining prediction scores … we can arrive at an overall score for each video".
3. **Ranking**: sắp theo điểm từ cao xuống thấp.
4. **Similarity check**: nếu các video đứng đầu quá giống nhau, "for example, they use the same sound", hệ thống thay bằng video khác.
5. **Recommendation rules**: ví dụ "ensure videos in your For You feed include content from creators in your region", và nội dung của tài khoản dưới 16 tuổi không được gợi ý.

Creator Academy nhắc lại pipeline này, và nói thêm rằng bước Selecting "takes into account your location and the region in which a video was posted" ([26/09/2026](https://www.tiktok.com/creator-academy/en/article/guidelines-recommendation-system-intro)).

### 1.2 Tín hiệu và dự đoán được liệt kê

Phụ lục của [Transparency Center](https://www.tiktok.com/transparency/en/recommendation-system):

| Nhóm | Nội dung |
|---|---|
| Về người xem | Khu vực ("aligning recommendations to local standards and cultural norms"); **ngôn ngữ thiết bị**; hệ điều hành; video vừa xem, xem hết, like hoặc skip; quan hệ follow |
| Dự đoán | Like, share, comment, "Not Interested"; **follow tác giả hoặc tương tác với profile**; xem hết, skip, favorite; "spend a certain amount of time viewing a video"; bấm vào soundtrack |
| Về video | Thời điểm đăng; **khu vực đăng**; **"The language setting of the video's author"**; soundtrack; **độ dài**; hashtag |

[Help Center](https://support.tiktok.com/en/using-tiktok/exploring-videos/how-tiktok-recommends-content) (không ghi ngày) chia tín hiệu cho For You như sau:

- User interactions: "Content you like, share, comment on, and watch in full or skip".
- Content information: "Sounds, hashtags, number of video views, and the country in which the content was published".
- User information: "Device settings, language preference, location, time zone and day, and device type".
- Về trọng số: "For most users, user interactions, which may include the time spent watching a video, are generally weighted more heavily than others."

### 1.3 Trọng số TikTok từng nói ra (newsroom 18/06/2020, ⚠️ cũ)

- Tín hiệu mạnh: "whether a user finishes watching a longer video from beginning to end, would receive greater weight".
- Tín hiệu yếu: "whether the video's viewer and creator are both in the same country". Ngôn ngữ, quốc gia và loại thiết bị "receive lower weight … since users don't actively express these as preferences".
- "neither follower count nor whether the account has had previous high-performing videos are direct factors" ([newsroom](https://newsroom.tiktok.com/en-us/how-tiktok-recommends-videos-for-you)).

⚠️ Bài này đã 6 năm. Không có nguồn chính mới hơn nhắc lại ý "follower không phải tín hiệu trực tiếp". Trong khi đó Help Center hiện tại liệt kê "number of video views" là một tín hiệu của For You.

### 1.4 Yếu tố chất lượng (Safety Center, 30/03/2026)

Theo [Making your feed For You](https://www.tiktok.com/safety/en/tools-and-guides/making-your-feed-for-you), hệ thống còn xét:

- "Avoiding repetitive recommendations, such as showing multiple similar posts in a row"
- "**Promoting local creators**"
- "**Prioritizing recent content over older content**"
- "Assessing factors like **the length of the post**, the time it was posted, and which sounds were used"
- Ngoài ra: "We double-check posts that are rising in popularity". Bài đang lên sẽ được kiểm duyệt thêm một lần.

### 1.5 Creator Academy khuyên gì

- [Understanding the Recommendation System (26/09/2026)](https://www.tiktok.com/creator-academy/en/article/guidelines-recommendation-system-intro): ba trọng tâm là **Watch time** ("Trim any lulls and use quick cuts"), **Trends** (cân bằng với bài nguyên bản) và **Engagement**. Trang này cũng ghi "Users spend 96% of their time on TikTok scrolling through the FYF".
- [Driving video traffic: The For You feed (19/08/2026)](https://www.tiktok.com/creator-academy/en/article/video-traffic-the-for-you-feed):
  - "Hook viewers in the first three seconds".
  - "Well-crafted, long-form content thrives on the For You feed".
  - "Choose a specific content niche … Use common video topics, hashtags, and sounds from your community".
  - "Post more often! The more videos you post, the more opportunities you give the recommendation system to understand your content".
- [Showcasing your expertise (19/08/2026)](https://www.tiktok.com/creator-academy/en/article/Showcasing-your-expertise-with-specialized-content):
  - Nội dung chuyên môn "must feature a demonstration or explain a process" và không gồm nội dung "meant only to advertise a product or service to direct users off-app".
  - Mẹo: "tease your finished product first".
  - Các định dạng gợi ý: Solve a problem, Bust a myth, Share a hack, Teach step by step, Respond to a comment.

### 1.6 Công cụ người xem dùng để chỉnh feed

| Công cụ | Ngày | Nguồn |
|---|---|---|
| Refresh For You feed (xóa cá nhân hóa của riêng For You) | trang Help hiện hành không ghi ngày | [Help Center](https://support.tiktok.com/en/account-and-privacy/account-privacy-settings/refresh-your-for-you-feed) |
| **Manage Topics**: chỉnh tần suất của hơn 10 chủ đề; "won't eliminate any topics entirely" | 03/06/2025 | [newsroom](https://newsroom.tiktok.com/more-ways-to-discover-new-content-and-creators-you-love?lang=en) |
| Smart Keyword Filters (dùng AI lọc cả từ đồng nghĩa) | 03/06/2025 | cùng nguồn |
| **Thanh trượt AIGC trong Manage Topics** (bắt đầu thử nghiệm) | 19/11/2025 | [newsroom](https://newsroom.tiktok.com/more-ways-to-spot-shape-and-understand-ai-content?lang=en) |
| Tắt cá nhân hóa: For You thành feed "Popular … from your region and internationally". **Chỉ có ở EU và Kazakhstan** | không ghi ngày | [Help Center](https://support.tiktok.com/en/using-tiktok/exploring-videos/how-tiktok-recommends-content) |

---

## 2. Không đủ điều kiện lên For You (FYF Ineligible)

Nguồn: Community Guidelines bản "2026 August", phát hành 25/08/2026, hiệu lực 24/09/2026 ([fyf-standards](https://www.tiktok.com/community-guidelines/en/fyf-standards)).

- **Nguyên tắc chung**: "We make ineligible for the FYF certain content that may not be suitable for a broad audience." Video không lên For You vẫn có thể được tìm thấy qua search hoặc profile. Creator xem được video nào bị loại trong **TikTok analytics**.
- **Cấp tài khoản**: "accounts that don't break the rules still post a lot of content that's ineligible for the FYF … we may make the account and its content ineligible for the FYF and harder to find" ([accounts-features](https://www.tiktok.com/community-guidelines/en/accounts-features)).
- **Các mục FYF Ineligible liên quan tới mate** ([integrity-authenticity](https://www.tiktok.com/community-guidelines/en/integrity-authenticity)):
  - "Reused or unoriginal content posted without creative edits, such as clips that show someone else's watermark or logo"
  - "Low-quality or minimally edited content, such as short clips made from GIFs only"
  - Engagement bait: "'Like-for-like' promises", "False incentives for gifting or following", "Misleading claims meant to boost views or popularity"
  - Nội dung AI trông như thật nhưng chưa xác nhận là AIGC, trình bày vấn đề công cộng theo cách dễ gây hiểu lầm
  - Thông tin sai gây hại ở mức vừa; "unverified information about crises"
- **Nội dung thương mại không khai báo**: "If you're promoting a product, brand, or business, you must use TikTok's content disclosure setting". Khai báo bắt buộc khi "Promoting your own business, product, or service". Nếu không khai báo, TikTok "may reduce its visibility" hoặc "remove it from the FYF" ([regulated-commercial-activities](https://www.tiktok.com/community-guidelines/en/regulated-commercial-activities), [accounts-features](https://www.tiktok.com/community-guidelines/en/accounts-features)).
- **Link**: CG không hạ hạng video chỉ vì có link. Link vi phạm thì bị gỡ, và có thể bị cấm đăng link tạm thời ([accounts-features, mục External Links](https://www.tiktok.com/community-guidelines/en/accounts-features)). Comment chứa "irrelevant promotions, or links" bị xếp thấp trong phần comment (cùng nguồn).
- **Tự động hóa**: "We strictly prohibit automation tools, scripts, or other tricks designed to bypass our systems". Spam gồm "Using automation to run many accounts or send repetitive content". Cho phép nhiều tài khoản "for fan content or creative expression—but not to deceive" ([integrity-authenticity](https://www.tiktok.com/community-guidelines/en/integrity-authenticity)).
- **Kiểm tra trước khi đăng**: **Content Check Lite** (30/07/2025) trên TikTok Studio web cho biết trước nội dung có khả năng bị loại khỏi For You hay không ([newsroom](https://newsroom.tiktok.com/en-us/new-tools-for-creators)).
- Creator Academy có bài giải thích bằng ngôn ngữ dễ hiểu ([Why is your video not getting recommended, 30/04/2026](https://www.tiktok.com/creator-academy/en/article/video-not-recommended)). Bài nói unoriginal là nội dung "recycled, spammy … borrowed from other places without a fresh spin". Bài cũng nói một video ít view "doesn't necessarily mean it wasn't recommended".

---

## 3. Tính nguyên bản

- [Creator Academy, Originality Policy (26/09/2026)](https://www.tiktok.com/creator-academy/en/article/tiktok-originality-policy) liệt kê bốn loại nội dung không nguyên bản:
  1. "Content copied completely from others"
  2. "largely repurposed from another source without adding any creative edits"
  3. "combined from multiple sources with little or no additional information and/or content value"
  4. "someone else's visible watermark or superimposed logo"
- Hậu quả: "Unoriginal content may be removed from the For You feed". Ngoài ra, CRP dùng "originality" như "a key metric in its rewards formula" (cùng nguồn).
- CG 08/2026: "Content is also ineligible for the FYF if it includes unoriginal or reused material without anything new" ([integrity-authenticity](https://www.tiktok.com/community-guidelines/en/integrity-authenticity)).
- Không nguồn chính nào của TikTok cấm đăng lại **video của chính mình** đã đăng ở nền tảng khác. Các quy định đều nhắm vào "someone else's" watermark/logo và nội dung "copied from others" **[suy luận từ cách diễn đạt]**. Cẩn thận vẫn hơn: xuất file sạch từ HyperFrames cho từng nền tảng, đừng tải lại file đã mang watermark của nền tảng khác.

---

## 4. Nội dung do AI tạo (AIGC)

### 4.1 Khi nào bắt buộc gắn nhãn (CG 08/2026)

- Định nghĩa AIGC: "Any image, video, or audio made or changed by AI. This can include realistic scenes or artistic styles, like anime, cartoons, or paintings" ([integrity-authenticity](https://www.tiktok.com/community-guidelines/en/integrity-authenticity)).
- **Bắt buộc** khi "AI or includes significant edits to show realistic-looking people or scenes". Gắn bằng nhãn AIGC của TikTok, hoặc bằng caption, sticker, watermark rõ ràng. Các trường hợp nêu tên:
  - thay mặt người này bằng người khác;
  - làm như ai đó nói điều họ không nói;
  - thêm hoặc bớt nền, vật thể, người theo cách gây hiểu lầm;
  - "AI-generated audio mimics the voice of a real person".
- **Không cần gắn nhãn**:
  - "small edits like color correction, reframing, or cropping";
  - "Using artistic styles, like anime";
  - "**Using generic text-to-speech (TTS) narration, when the TTS isn't a recognizable voice of a known individual**".
- Nếu không gắn nhãn khi bắt buộc: nội dung "may be removed, restricted, or labeled by our team". Nội dung AI trông như thật mà chưa được xác nhận là AIGC có thể bị loại khỏi For You (§2).
- Help Center "About AI-generated content" (không ghi ngày) **khuyến khích** gắn nhãn cho mọi nội dung "completely generated or significantly edited by AI", nhưng chỉ **bắt buộc** với nội dung trông như thật ([Help Center](https://support.tiktok.com/en/using-tiktok/creating-videos/ai-generated-content)).

### 4.2 Nhãn có làm giảm reach không?

- TikTok nói nhãn **không làm giảm phân phối**:
  - "Turning on the AI-generated content setting won't affect the distribution of your video as long as it doesn't violate our Community Guidelines" ([Help Center](https://support.tiktok.com/en/using-tiktok/creating-videos/ai-generated-content); [Creator Academy 26/09/2026](https://www.tiktok.com/creator-academy/en/article/ai-generated-content-label)).
  - FAQ của Creator Academy: "Applying the label will not have an impact on organic engagement with your content, provided that your content isn't violative and meets our For You feed eligibility standards."
  - Trang in-app: "adding the label to your content won't affect the distribution of your video, and you can't remove the label after posting" ([tns-inapp](https://www.tiktok.com/tns-inapp/pages/ai-generated-content), không ghi ngày).
- **Nhưng người xem có thể tự giảm AIGC.** Ngày 19/11/2025 TikTok thử nghiệm thanh trượt AIGC trong Manage Topics, "to empower people to choose how much AIGC they want to see in their For You feeds", "rather than removing or replacing content in feeds entirely" ([newsroom](https://newsroom.tiktok.com/more-ways-to-spot-shape-and-understand-ai-content?lang=en)). Nguồn chính không nói thanh trượt dựa trên nhãn hay trên bộ phân loại chủ đề. Nếu nó dựa trên nhãn, thì video có nhãn sẽ ít xuất hiện hơn với người đã kéo thanh xuống **[suy luận]**. Theo nguồn thứ cấp, tới 03/2026 tính năng vẫn đang thử nghiệm **[chưa xác minh từ nguồn chính]** ([socialbee](https://socialbee.com/blog/tiktok-updates-news-features/)).
- **Gắn nhãn tự động**: TikTok có thể tự gắn nhãn "AI-generated" khi dùng hiệu ứng AI của TikTok, hoặc khi file tải lên mang **C2PA Content Credentials** ([Help Center](https://support.tiktok.com/en/using-tiktok/creating-videos/ai-generated-content)). Từ 11/2025 TikTok thêm "invisible watermarks" cho nội dung làm bằng công cụ AI của TikTok và nội dung có C2PA. Tới thời điểm đó TikTok đã gắn nhãn "over 1.3 billion videos" ([newsroom 19/11/2025](https://newsroom.tiktok.com/more-ways-to-spot-shape-and-understand-ai-content?lang=en)).
- **Symphony Creative Studio** (công cụ AI của TikTok for Business):
  - "An AI-generated label is added to all exported videos" ([ads Help, 08/2026](https://ads.tiktok.com/help/article/how-to-translate-and-dub-videos-with-symphony-creative-studio?lang=en)).
  - Creator Academy: nội dung làm bằng Creative Studio "will be automatically labeled as AI-generated" ([19/08/2026](https://www.tiktok.com/creator-academy/en/article/symphony)).
  - "Symphony Creative Studio does not have any direct connection with the TikTok algorithm" ([ads Help, 06/2026](https://ads.tiktok.com/help/article/about-symphony-creative-studio?lang=en)).
- STEM feed loại "Misleading or unlabelled AI-generated content" ([Creator Academy 19/08/2026](https://www.tiktok.com/creator-academy/en/article/creating-STEM-content)).

---

## 5. Độ dài video và Creator Rewards Program (CRP)

- Giới hạn upload: "up to 60 minutes long" ([Creator Academy 19/08/2026](https://www.tiktok.com/creator-academy/en/article/tool-one-minute-feature-overview)).
- Độ dài là một tín hiệu: phụ lục Transparency có "The length of the video"; Safety Center có "Assessing factors like the length of the post" (§1.2, §1.4). Không nguồn chính nào nói độ dài nào là tối ưu.
- TikTok đẩy video dài hơn 1 phút:
  - Help Center: "Create content that's longer than one minute" ([How to grow your audience](https://support.tiktok.com/en/using-tiktok/growing-your-audience/how-to-grow-your-audience), không ghi ngày).
  - Creator Academy: "Well-crafted, long-form content thrives on the For You feed" ([19/08/2026](https://www.tiktok.com/creator-academy/en/article/video-traffic-the-for-you-feed)).
- **CRP** ([Creator Academy, 11/09/2026](https://www.tiktok.com/creator-academy/en/article/creator-rewards-program)):
  - **Vẫn yêu cầu** "original, high-quality content **over 1 minute long**", ít nhất 1.000 view từ For You, độ phân giải 1080p trở lên, không phải Duet/Stitch/Photo Mode.
  - Điều kiện tài khoản: personal account, 18 tuổi trở lên, ≥10.000 follower, ≥100.000 view trong 30 ngày.
  - Chỉ có ở **US, UK, DE, JP, KR, FR, MX, BR**. Không có Việt Nam. Qualified view cũng phải đến từ các nước này: xem ≥5 giây, hoặc ≥30 giây nếu đến từ Search.
  - Danh sách "content ineligible for the For You feed will not be eligible for the Creator Rewards Program, including…" có mục "**Low quality**: … split screens, meaningless reactions, low-quality images, or **slide videos**". ⚠️ Câu này mơ hồ: không rõ "slide videos" là tiêu chí loại khỏi For You hay chỉ khỏi CRP. Với video motion graphics thì đây là rủi ro cần nhớ.
  - Additional Reward (từ mùa xuân 2025) thưởng nội dung "Well-crafted", "Engaging" và "**Specialized** content that demonstrates your expertise in a niche".

---

## 6. Search và SEO trên TikTok

- **Search xếp hạng chủ yếu theo mức khớp nội dung**, khác với For You: "For most users, content information, which may include how well the content matches the specific term entered, is generally weighted more heavily than others". Tín hiệu nội dung của Search gồm "How well the content matches your search inquiry, hashtags, and the sound used" ([Help Center](https://support.tiktok.com/en/using-tiktok/exploring-videos/how-tiktok-recommends-content), không ghi ngày).
- [Getting discovered with Search (Creator Academy, 19/08/2026)](https://www.tiktok.com/creator-academy/en/article/Getting-discovered-with-Search):
  - **Nên** dùng "your video's description, voiceovers, captions, and on-screen text to make it clear just what your video's all about". Như vậy caption, chữ trên màn hình và **lời nói** đều được dùng để hiểu video.
  - **Myth 1**: chèn chuỗi keyword dạng chữ nhỏ hoặc trượt khỏi màn hình "can get your video flagged for spam".
  - **Myth 2**: không cần dựng video trong app TikTok mới được lên search ("your contents' search performance won't be impacted"). Video dựng bằng HyperFrames bên ngoài không bị thiệt.
  - **Myth 3**: một danh sách dài hashtag không giúp xếp hạng cao hơn. "Overloading your caption with loosely related keywords or hashtags can actually work against you."
  - Creator Search Insights cho xem chủ đề người dùng đang tìm. Bộ lọc "Searches by followers" cần 1.000 follower.
  - CRP "now rewards videos for their search value" (chỉ áp dụng ở nước có CRP).
- CG 08/2026: Search "prioritize[s] entertaining and informative content in search recommendations" và đưa "content from reliable sources at the top" ([accounts-features](https://www.tiktok.com/community-guidelines/en/accounts-features)).
- Video không lên For You "may still find it through search or by going to a creator's account" ([fyf-standards](https://www.tiktok.com/community-guidelines/en/fyf-standards)).

---

## 7. Ngôn ngữ và khu vực

| Điểm | Nội dung | Nguồn |
|---|---|---|
| Khu vực là tín hiệu và có luật riêng | "ensure videos in your For You feed include content from creators in your region"; "Promoting local creators"; tín hiệu "The region the video was posted from" | [Transparency](https://www.tiktok.com/transparency/en/recommendation-system); [Safety 30/03/2026](https://www.tiktok.com/safety/en/tools-and-guides/making-your-feed-for-you) |
| Ngôn ngữ là tín hiệu ở cả hai phía | Người xem: "The language your device is set to". Video: "The language setting of the video's author". Người dùng mới nhận "Content influenced by your location and language settings" | [Transparency](https://www.tiktok.com/transparency/en/recommendation-system); [Help Center](https://support.tiktok.com/en/using-tiktok/exploring-videos/how-tiktok-recommends-content) |
| Trọng số của quốc gia và ngôn ngữ | Năm 2020 TikTok nói các tín hiệu này có "lower weight". ⚠️ Chưa có nguồn mới hơn | [newsroom 2020](https://newsroom.tiktok.com/en-us/how-tiktok-recommends-videos-for-you) |
| Phụ đề và dịch tự động | 21/07/2022: auto-generated captions; dịch caption, mô tả và text sticker, cho đợt ngôn ngữ đầu là EN, PT, DE, ID, IT, KO, ZH, ES, TR (**không có tiếng Việt**), "available on select videos". ⚠️ Không tìm thấy danh sách ngôn ngữ hiện hành trong nguồn chính. Tiếng Việt có được tạo phụ đề tự động không: **cần kiểm tra trong app** | [newsroom](https://newsroom.tiktok.com/en-us/auto-translations-and-captions); [Help Center Accessibility](https://support.tiktok.com/en/using-tiktok/exploring-videos/accessibility-for-watching-videos) |
| AI dịch và lồng tiếng | Không thấy tính năng lồng tiếng AI cho video organic của creator thường (khác với Meta, xem file anh em §4). Có **Symphony Creative Studio → Translate & dub** (clone giọng hoặc dùng giọng có sẵn, lip-sync tùy chọn), nhưng cần tài khoản TikTok for Business, video xuất ra **bị gắn nhãn AI**, và trang hướng dẫn không công bố danh sách ngôn ngữ | [ads Help 08/2026](https://ads.tiktok.com/help/article/how-to-translate-and-dub-videos-with-symphony-creative-studio?lang=en) |
| STEM feed | "available globally, except Thailand, Turkey, Israel, Central Asia, Japan, Korea, and **Vietnam**". Hashtag **#STEMTok "only recognizes videos in English"**. Nội dung được Common Sense Networks và Poynter duyệt. Không nhận nội dung quảng bá hoặc AIGC không gắn nhãn. Mọi video trong STEM feed cũng lên For You | [What is the STEM feed (19/08/2026)](https://www.tiktok.com/creator-academy/en/article/what-is-STEM-feed); [Creating STEM content (19/08/2026)](https://www.tiktok.com/creator-academy/en/article/creating-STEM-content) |
| Pháp lý Việt Nam | Nghị định 147/2024/NĐ-CP (hiệu lực 25/12/2024): tài khoản mạng xã hội phải xác thực bằng số điện thoại di động Việt Nam thì mới được đăng bài và bình luận. Nguồn là cơ sở dữ liệu luật, chưa đọc văn bản gốc trên cổng Chính phủ | [thuvienphapluat.vn](https://thuvienphapluat.vn/banan/tin-tuc/nghi-dinh-1472024ndcp-chinh-thuc-bat-buoc-xac-thuc-tai-khoan-mang-xa-hoi-bang-so-dien-thoai-12259) |
| Nhiều tài khoản | Được phép "for fan content or creative expression—but not to deceive". Không có hướng dẫn chính thức nào về việc tách tài khoản theo ngôn ngữ | [CG integrity-authenticity](https://www.tiktok.com/community-guidelines/en/integrity-authenticity) |

---

## 8. Thay đổi 2025–2026

| Ngày | Thay đổi | Nguồn |
|---|---|---|
| 11/03/2025 | Family Pairing thêm **Time Away**: phụ huynh chặn teen dùng TikTok theo khung giờ. Không liên quan tới creator | [newsroom](https://newsroom.tiktok.com/en-us/new-ways-we-are-supporting-parents-and-helping-teens-build-balanced-digital-habits) |
| 16/04/2025 | Thử nghiệm **Footnotes** (ghi chú cộng đồng kiểu Community Notes, US), dùng "bridge-based ranking". Bài công bố lấy ví dụ trực tiếp "a complex STEM-related concept" | [newsroom](https://newsroom.tiktok.com/footnotes?lang=en) |
| 03/06/2025 | **Manage Topics**, Smart Keyword Filters, ra mắt guide "Making your feed For You" | [newsroom](https://newsroom.tiktok.com/more-ways-to-discover-new-content-and-creators-you-love?lang=en) |
| 30/07/2025 | Footnotes hiển thị công khai ở US; **Content Check Lite**; Creator Care Mode; Creator Inbox | [newsroom](https://newsroom.tiktok.com/en-us/tiktok-announces-a-suite-product-features); [newsroom](https://newsroom.tiktok.com/en-us/new-tools-for-creators) |
| 19/11/2025 | Thử nghiệm thanh trượt **AIGC** trong Manage Topics; invisible watermark; quỹ AI literacy 2 triệu USD | [newsroom](https://newsroom.tiktok.com/more-ways-to-spot-shape-and-understand-ai-content?lang=en) |
| 22–23/01/2026 | Thành lập **TikTok USDS Joint Venture LLC**: "will retrain, test, and update the content recommendation algorithm on U.S. user data", lưu trên Oracle cloud. Nhờ interoperability, US vẫn có "a global TikTok experience, ensuring U.S. creators can be discovered … on a global scale" | [newsroom 23/01/2026](https://newsroom.tiktok.com/announcement-from-the-new-tiktok-usds-joint-venture-llc?lang=en) |
| 02/2026 | Một nữ phát ngôn viên TikTok nói với NPR rằng thuật toán chưa thay đổi kể từ khi đổi chủ. **[nguồn báo chí trích phát ngôn, không phải thông cáo]** | [NPR 04/02/2026](https://www.npr.org/2026/02/04/nx-s1-5701409/tiktok-censorship-report-epstein) |
| 03/2026 | Theo nguồn thứ cấp, Refresh For You giờ mất vài ngày mới xong. Trang Help hiện hành không nhắc chuyện này **[chưa xác minh từ nguồn chính]** | [socialbee](https://socialbee.com/blog/tiktok-updates-news-features/) |
| 25/08/2026 → 24/09/2026 | **Community Guidelines bản 08/2026**. Bản hiện hành ghi: AIGC miễn nhãn cho TTS chung chung và phong cách nghệ thuật; xử lý vi phạm theo thứ tự "a warning, temporary feature or account restriction, or a full account ban", có "optional policy training" để khôi phục quyền; nội dung thương mại không khai báo "may reduce its visibility". Chưa đối chiếu với bản trước, nên chưa biết điểm nào là mới. Theo nguồn thứ cấp, các điểm này là thay đổi mới **[chưa xác minh từ nguồn chính]** ([socialbee](https://socialbee.com/blog/tiktok-updates-news-features/)) | [CG](https://www.tiktok.com/community-guidelines/en/accounts-features) |

**Thuật toán US có khác thuật toán toàn cầu không?** Chưa có nguồn chính nào công bố lộ trình, thời điểm hoàn tất hay nội dung cụ thể của việc "retrain". Theo nguồn báo chí, ngày 29/05/2026 một thượng nghị sĩ Mỹ đã gửi câu hỏi tới JV và Oracle về các biện pháp bảo vệ của nền tảng ([US News/Reuters 29/05/2026](https://www.usnews.com/news/politics/articles/2026-05-29/democratic-senator-wants-tiktok-us-jv-oracle-to-answer-questions-on-platform-safeguards), trang không tải được khi viết bài, chỉ thấy tiêu đề). Người xem ở Việt Nam dùng TikTok toàn cầu (ROW), nên thay đổi này chỉ chạm tới khán giả US của nhánh tiếng Anh **[suy luận]**.

---

## 9. Video faceless, không có giọng thật, toàn chữ

- **Faceless + voiceover là định dạng TikTok gợi ý chính thức**: "use only full-frame visuals and add an engaging voice-over". Cùng bài khuyên giữ chữ và phụ đề trong **safe zone**, không dùng hình vỡ và phải khớp khung 9:16 ([Optimizing your video layout, 19/08/2026](https://www.tiktok.com/creator-academy/en/article/optimize-video-layout)). Cùng bài cũng nói "Showing your face on screen is a simple way to connect … and hold their attention longer". Có mặt là lợi thế, nhưng không phải điều kiện.
- **Âm thanh**: "When someone opens TikTok, they're diving straight into a full-screen experience, sound-on experience". Nên dùng nhạc có license hoặc âm thanh gốc, ví dụ Commercial Music Library ([The importance of sound, 19/08/2026](https://www.tiktok.com/creator-academy/en/article/the-importance-of-sound)). Soundtrack là tín hiệu (§1.2), và thao tác "Tap the video's soundtrack" là một dự đoán.
- **Video tắt tiếng**: không tìm thấy nguồn chính nào của TikTok nói video tắt tiếng bị hạ hạng. (Instagram thì có nói, xem file anh em §2.3.)
- **TTS**: TikTok có tính năng Text-to-Speech trong app ([Creator Academy, Accessibility 19/08/2026](https://www.tiktok.com/creator-academy/en/article/making-your-videos-more-accessible)). TTS chung chung không cần nhãn AIGC (§4.1). Clone giọng của một người thật thì thuộc diện "AI-generated audio mimics the voice of a real person", phải khai báo.
- **Toàn chữ hoặc dạng slide**: CRP xếp "slide videos" và "low-quality images" vào nhóm Low quality (§5). CG xếp "short clips made from GIFs only" vào FYF Ineligible (§2). Không có nguồn chính nào nói video toàn chữ nhưng có chuyển động bị loại.

---

## 10. Chưa xác minh hoặc có thể lỗi thời

- "Video phải đạt khoảng 70% completion mới viral" và "creator US mất 20–40% reach sau khi JV retrain": chỉ có ở blog thứ cấp, không có nguồn chính **[chưa xác minh từ nguồn chính]** ([shortgen.io](https://shortgen.io/blog/tiktok-algorithm-change-2026), [fastclips.app](https://fastclips.app/learn/tiktok-us-algorithm-joint-venture)).
- "TikTok gắn cờ cả tài khoản AI spam, không chỉ từng video" (2026) **[chưa xác minh từ nguồn chính]** ([noobclaw](https://noobclaw.com/blog/tiktok-ai-spam-detection-manage-topics-2026/)). Điều gần nhất có nguồn chính là quy định cấp tài khoản trong CG (§2).
- Bài newsroom 2020 về trọng số đã cũ (§1.3). Trang Help Center và Transparency không ghi ngày cập nhật.
- Danh sách ngôn ngữ của phụ đề tự động và dịch tự động, và danh sách ngôn ngữ của Symphony dub, đều không có trong nguồn chính đọc được.
- Thanh trượt AIGC: chưa có nguồn chính xác nhận đã ra mắt toàn cầu, hay đã có ở Việt Nam.

---

## 11. Hàm ý cho mate (mục tiêu: personal brand của người làm mate)

### 11.1 Áp dụng cho cả hai nhánh

1. **Tối ưu watch time và tỷ lệ xem hết trước mọi thứ khác.** Đây là tín hiệu TikTok nói rõ nhất (§1.2–1.5). Cụ thể **[suy luận từ Creator Academy]**:
   - 3 giây đầu cho thấy ngay kết quả: màn hình nhiều agent đang chạy, hoặc con số "trước/sau" ("tease your finished product first").
   - Cắt mọi khoảng lặng.
   - Dùng các định dạng "Bust a myth", "Solve a problem", "Teach step by step".
2. **Độ dài: thử cả hai mức, dưới và trên 1 phút.** TikTok đẩy video >1 phút (§5), nhưng không có con số tối ưu chính thức. Một "bài học" về điều phối agent có thể dài 60–90 giây mà vẫn chặt **[suy luận]**. Đo bằng analytics, đừng theo một con số ngoài mạng.
3. **Motion graphics từ HyperFrames không cần nhãn AIGC** nếu không có cảnh hay người trông như thật: đó là đồ họa hoặc phong cách nghệ thuật, không phải AI dựng cảnh thật (§4.1).
   - Giọng đọc:
     - TTS chung chung → không cần nhãn.
     - Clone giọng chính mình → nên gắn nhãn, vì quy định nói "voice of a real person" mà không loại trừ chính creator **[suy luận]**.
     - Ghi âm giọng thật → không phải AIGC.
   - Nhãn AI không làm giảm phân phối, theo lời TikTok. Nhưng người xem có thể kéo thanh AIGC xuống, nên **đừng gắn nhãn khi không cần, và tuyệt đối không bỏ nhãn khi bắt buộc** **[suy luận]**.
   - Kiểm tra xem công cụ TTS hoặc ảnh có nhúng C2PA vào file xuất ra không. Nếu có, TikTok sẽ tự gắn nhãn "AI-generated" và không gỡ được (§4.2).
4. **Đừng để video trông như slide.** "slide videos" và "low-quality images" bị CRP xếp là chất lượng thấp (§5). Mỗi cảnh cần chuyển động có chủ đích, 1080p, 9:16, chữ nằm trong safe zone, có voiceover. Định dạng "full-frame visuals + voice-over" được Creator Academy nêu tên (§9).
5. **Đưa keyword vào cả lời nói, chữ trên màn hình và caption.** TikTok đọc cả ba nguồn để hiểu video (§6). Ví dụ "Claude Code", "AI coding agent", "multi-agent". Dùng ít hashtag và đúng chủ đề. Không nhét keyword dạng chữ nhỏ (bị coi là spam).
6. **Khai báo nội dung thương mại khi video quảng bá mate.** CG yêu cầu khai báo khi "Promoting your own … product". Video "bài học" không nhắc mate như sản phẩm có thể không cần **[suy luận, nên đọc kỹ Branded Content Policy trước khi đăng loạt đầu]**. Nội dung quảng bá cũng bị loại khỏi STEM feed (§7). Chiến lược "value-first" vì thế khớp với cả For You lẫn STEM: bài học là chính, mate chỉ xuất hiện như công cụ dùng để minh họa.
7. **Không engagement bait** ("comment 'mate' để nhận link", "follow để xem phần 2" nếu không có phần 2). "Misleading claims meant to boost views" nằm trong mục FYF Ineligible (§2).
8. **Kiểm tra trước và sau khi đăng**: chạy Content Check Lite (TikTok Studio web) trước, rồi xem analytics có video nào bị đánh dấu không đủ điều kiện For You không (§2).
9. **Đăng đều và giữ một niche.** Creator Academy khuyên "Post more often" và "Choose a specific content niche" (§1.5). Mỗi tuần một chủ đề về điều phối agent, nhiều video, tốt hơn rải rác nhiều chủ đề **[suy luận]**.
10. **Nội dung phải chính xác, có dẫn chứng.** Search ưu tiên "reliable sources". Footnotes (US) có thể gắn ghi chú vào video STEM. STEM feed do Poynter kiểm chứng. Không dùng benchmark không có nguồn **[suy luận]**.

### 11.2 Nhánh tiếng Việt (dev Việt)

- **Lợi thế creator địa phương**: TikTok có luật đưa creator cùng khu vực vào For You và "Promoting local creators" (§1.1, §1.4). Một tài khoản đăng từ Việt Nam, đặt ngôn ngữ tiếng Việt, nói tiếng Việt sẽ được phân phối cho người xem Việt Nam trước **[suy luận từ luật khu vực và tín hiệu ngôn ngữ của tác giả]**.
- **Không có STEM feed ở Việt Nam** (§7). Nội dung dev Việt chỉ có hai đường: For You và Search. Vì vậy **SEO trong Search quan trọng hơn** với nhánh này **[suy luận]**:
  - Đặt tiêu đề trên màn hình và nói theo đúng cụm người Việt hay gõ, ví dụ "cách dùng Claude Code", "AI agent lập trình".
  - Giữ thuật ngữ tiếng Anh nguyên văn, vì dev Việt thường tìm bằng tên công cụ tiếng Anh **[suy luận, nên kiểm tra bằng Creator Search Insights khi đủ 1.000 follower]**.
- **Phụ đề**: chưa xác minh được TikTok có tạo phụ đề tự động tiếng Việt hay không (§7). **Burn-in phụ đề tiếng Việt ngay trong HyperFrames**, vừa chắc chắn hiển thị vừa là "on-screen text" cho Search.
- **CRP không áp dụng** ở Việt Nam. Không cần ép video dài hơn 1 phút vì tiền thưởng; chỉ kéo dài khi analytics cho thấy thời gian xem tăng.
- **Tài khoản**: phải xác thực bằng số điện thoại Việt Nam mới đăng bài được (Nghị định 147, §7). Dùng **personal account**, không dùng Business **[suy luận: Business bị loại khỏi CRP; với nhánh VN thì đây là giữ tùy chọn mở, không quyết định gì]**.
- **Quan hệ với Threads (hub)**: TikTok không có tín hiệu nào được công bố cho lưu lượng đến từ ngoài. Mục tiêu của TikTok VN là **follow và vào profile** (đều là dự đoán trong §1.2). Bio TikTok dẫn về Threads hoặc profile chính **[suy luận]**.

### 11.3 Nhánh tiếng Anh (dev quốc tế)

- **Tách tài khoản tiếng Anh riêng [suy luận]**. Lý do:
  - "The language setting of the video's author" là tín hiệu của video. Một tài khoản trộn hai ngôn ngữ gửi tín hiệu ngôn ngữ lẫn lộn.
  - TikTok cho phép nhiều tài khoản "for … creative expression" (§7).
  - Đặt ngôn ngữ thiết bị hoặc app của tài khoản này sang tiếng Anh.
- **Rào cản khu vực [suy luận]**: tài khoản vẫn đăng từ Việt Nam, nên luật "creators in your region" và tín hiệu "region the video was posted from" có thể khiến lượt phân phối đầu nghiêng về người xem ở Việt Nam hoặc Đông Nam Á biết tiếng Anh. **Không dùng VPN hay giả vị trí**: CG cấm "tricks designed to bypass our systems" (§2). Đo tỷ lệ người xem theo quốc gia trong analytics sau 4–6 tuần rồi mới đánh giá.
- **STEM feed là cơ hội lớn nhất của nhánh này**. Feed có ở US, EU và phần lớn thế giới (trừ Việt Nam), #STEMTok chỉ nhận video tiếng Anh, và "a programmer teaching coding" được nêu đích danh (§7). Điều kiện:
  - mục tiêu học tập rõ ràng;
  - có độ sâu và có dẫn chứng;
  - **không quảng bá** (không kêu gọi mua, không chỉ để dẫn ra ngoài app);
  - AIGC phải có nhãn.

  → Video "bài học" thuần (ví dụ "vì sao chạy song song nhiều agent cần một cơ chế khóa file") có thể đủ điều kiện. Video giới thiệu mate thì không **[suy luận]**. Gắn #STEMTok cho bài học, không gắn cho video quảng bá.
- **Làm bản tiếng Anh gốc bằng HyperFrames, không dùng Symphony dub [suy luận]**:
  - Symphony cần tài khoản Business, gắn nhãn AI cho mọi video xuất ra, và không công bố hỗ trợ tiếng Việt (§7).
  - HyperFrames đã render được bản tiếng Anh riêng.
  - Dùng TTS chung chung hoặc giọng đọc thật để khỏi phải gắn nhãn.
- **Khán giả US**: JV cam kết retrain thuật toán trên dữ liệu US, nhưng chưa có mốc thời gian nào được công bố, và cam kết "global TikTok experience" cho creator nước ngoài (§8). Chưa có cơ sở để đổi chiến lược. Chỉ theo dõi tỷ lệ view từ US trong analytics.
- **Footnotes (US)**: bài tiếng Anh về số liệu hiệu năng agent có thể bị gắn footnote nếu thiếu dẫn chứng **[suy luận]**. Ghi nguồn trên màn hình cho mọi con số.
