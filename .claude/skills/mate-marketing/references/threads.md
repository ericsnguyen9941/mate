# Threads: mồi lửa

Threads là trung tâm của campaign. Mọi video, mọi bài học đều quay về đây để thành hội thoại.

Bài Threads thắng khi nó tạo ra hai thứ:

- **Cây reply sôi nổi.** Threads đo tương tác trong cây reply tới hai tầng, trong 1–6 giờ đầu.
- **Sự tò mò về tác giả.** Threads có dự đoán việc người xem bấm vào profile, rồi bấm sang bài khác.

Cả hai đều phục vụ tên tuổi cá nhân. Số liệu và nguồn ở `docs/research/social-algorithms-2026-10-09.md` §1.

## Một bài là một mồi lửa

Mồi lửa là một khẳng định đủ sắc để dev muốn trả lời, nhưng có bằng chứng đứng sau nên không gãy khi bị cãi. Các dạng sau hợp với giọng mate. Mỗi bài lấy từ một mục trong `material.md`.

1. **Con số khó tin, có nguồn.** "87K token đầu tiên của mọi lần Mate gọi model là phần cố định. 50K trong đó không phải của mate." Câu sau phải nói người ta đang trả tiền cho cái gì.
2. **Thú tội / postmortem.** Một lần crew hoặc mate làm hỏng việc, nguyên nhân gốc, và điều captain đổi sau đó. Dev reply để kể lần hỏng của chính họ.
3. **Quyết định ngược số đông, kèm lý do.** Ví dụ: "Mate không được sửa code. Tôi chặn bằng cấu trúc (cwd không có code), không bằng lời dặn." Người đồng ý và người phản đối đều reply.
4. **So sánh có đo.** Claude Code, Codex, pi, Grok trên cùng một hợp đồng harness. Viết công bằng và chỉ nói điều evidence cho phép.
5. **Nhật ký build.** Một cảnh trong ngày làm việc với crew, kèm ảnh console hoặc terminal: "Hôm nay Mate tự từ chối merge vì reviewer chưa cho qua."
6. **Câu hỏi thật.** Một vấn đề captain đang thật sự cân nhắc, nêu rõ các phương án đang có. Câu trả lời của cộng đồng có thể vào bài sau.

Bài kết thúc bằng một quan điểm có thể cãi lại, hoặc một câu hỏi mà captain thật sự cần câu trả lời. Không xin like, follow hay comment một từ khoá: đó là engagement bait, và Threads ghi rõ là bị hạ hạng.

## Cấu trúc một bài

- **Dòng đầu chọc đúng pain của người đọc.** Ví dụ: "Agent báo test pass. Bạn merge luôn, đúng không?". Câu này là hook, cùng giọng châm chọc với video.
- **Dòng đầu tự đứng được.** Người lướt chỉ đọc dòng đầu, nên con số hay mệnh đề nằm ngay ở đó, không có câu dẫn.
- **Thân bài ngắn.** Một bài Threads giới hạn 500 ký tự. Chuyện dài hơn thì viết thành chuỗi, mỗi bài trong chuỗi có một ý. Kiểm tra lại giới hạn trên app trước khi đăng.
- **Đúng một topic tag mỗi bài.** Thử cả tag tiếng Việt và tiếng Anh, so insights, rồi ghi tag thắng vào `CAMPAIGN.md`. Tag có dấu ba chấm là tag có community: tham gia các community về lập trình và AI nếu có.
- **Kèm media khi có thể:** ảnh terminal, ảnh console, sơ đồ, hoặc clip từ video của tuần. Theo Threads, bài có media kèm chữ chạy tốt hơn hẳn. Ảnh demo phải ghi "dữ liệu demo".
- **Link GitHub hay landing page là phụ.** Threads không hạ hạng link, nhưng link không kéo reply. Nếu cần thì đặt link ở reply đầu tiên của chính mình, và đừng để link thành lời kêu gọi chính của bài.
- **Nội dung viết riêng cho Threads.** Không đăng lại nguyên văn caption Instagram hay Facebook. Người dùng Threads follow những tài khoản khác với trên Instagram.

## Reply quan trọng ngang đăng bài

Reply chiếm gần một nửa lượt xem trên Threads.

- **Giờ đầu sau khi đăng:** trả lời mọi reply bằng nội dung thật, như thêm một con số hay một chi tiết của lần hỏng. Đó là khoảng cây reply được đo.
- **Mỗi ngày, ngoài bài của mình:** reply có thực chất vào thread của các dev và AI builder lớn, cùng ngôn ngữ với nhánh campaign. Một reply hay là một mồi lửa thu nhỏ: mang theo một con số hoặc kinh nghiệm từ mate, không nhắc link.
- Ghi vào `CAMPAIGN.md` các tài khoản hay reply, để nhánh VN và nhánh EN mỗi nhánh có một danh sách riêng.

## Nhịp và series

- Mỗi ngôn ngữ đăng 2–5 bài/tuần. Đăng thường xuyên đi đôi với nhiều view hơn trên mỗi bài.
- Đặt tên series và lặp lại tên đó, để người xem nhận ra captain qua format khi không có gương mặt. Ví dụ: "Crew hỏng việc #3", "Số liệu thứ Sáu".
- Mỗi video tuần đó có một bài neo trên Threads, đăng cùng ngày hoặc trước một ngày. Bài neo là phiên bản mồi lửa của bài học trong video: video trả lời câu hỏi mà bài neo đặt ra.

## Hai nhánh

- **VN:** dev Việt. Nhắc chuyện quen thuộc với họ, ví dụ chi phí token tính ra tiền Việt hay team nhỏ không có nhiều người review. Thuật ngữ giữ tiếng Anh như trong docs.
- **EN:** dev builder quốc tế. Họ đã thấy hàng trăm demo agent, nên con số và chỗ hỏng mới tạo khác biệt. Viết lại từ chất liệu, không dịch bản VN.
- Một hay hai tài khoản là quyết định của captain, ghi trong `CAMPAIGN.md`. Nếu đang chạy một tài khoản chung, xen kẽ hai ngôn ngữ và theo dõi insights theo quốc gia. Sau 4–6 tuần, đọc dữ liệu đó để quyết có tách hay không.

## Định dạng file bài

Mỗi bài là một file `threads/<vi|en>/<NN>-<slug>.md` gồm:

- ngày dự kiến;
- dạng mồi lửa;
- topic tag;
- media đính kèm;
- mục chất liệu đã dùng (id trong `material.md`);
- nội dung bài, mỗi bài con trong chuỗi tách bằng `---`;
- reply đầu tiên của chính mình, nếu có.
