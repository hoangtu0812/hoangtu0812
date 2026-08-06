# Bài 6: .dockerignore & Build Context

## Mục tiêu
- Hiểu rõ khái niệm Build Context và những dữ liệu Docker gửi cho daemon khi thực hiện quá trình build.
- Giải thích được tại sao một build context quá lớn sẽ làm chậm tốc độ build và phình to Docker image.
- Nắm vững cú pháp và cách sử dụng file `.dockerignore`.
- Phân biệt rõ ràng chức năng và sự khác nhau giữa `.dockerignore` và `.gitignore`.
- Nắm bắt các best practices để tối ưu hoá build context, tránh rò rỉ dữ liệu nhạy cảm.

---

## 1. Build Context là gì?

### Giao tiếp giữa Docker Client và Docker Daemon
Khi bạn chạy lệnh `docker build -t my-app .`, dấu `.` ở cuối cùng chính là **Build Context** (ngữ cảnh build). Dấu chấm này chỉ định thư mục hiện tại chứa các file, tài nguyên cần thiết để tạo nên Docker image.

> 💡 **Tip:** Hãy tưởng tượng Docker Client giống như một người trợ lý đi chợ mua nguyên liệu, còn Docker Daemon là một vị đầu bếp chuyên nghiệp (thậm chí có thể đang ở một máy chủ đám mây từ xa). Khi bạn yêu cầu "Nấu món ăn này" (build), người trợ lý phải gói **toàn bộ** nguyên liệu có trong nhà bếp của bạn (Build Context) và gửi chuyển phát nhanh cho người đầu bếp.

Về mặt kỹ thuật, Docker CLI sẽ tự động gom và nén toàn bộ các file, thư mục bên trong ngữ cảnh này thành một file lưu trữ (`.tar`) và gửi thẳng cho Docker daemon trước khi bắt đầu thực thi bất kỳ lệnh nào trong `Dockerfile`.

```bash
# Đầu ra thường thấy khi build image
# Dòng đầu tiên luôn thể hiện việc gửi context này
Sending build context to Docker daemon  1.52MB
Step 1/5 : FROM node:18-alpine
...
```

**Kết quả:**
```
Quá trình gửi context diễn ra nhanh chóng nếu thư mục gọn nhẹ, nhưng sẽ rất chậm chạp nếu thư mục chứa hàng GB dữ liệu rác.
```

---

## 2. Tại sao Build Context lớn lại làm chậm quá trình build?

### Vấn đề của việc "đóng gói mù quáng"
Nếu thư mục dự án của bạn đang chứa các file dung lượng khổng lồ không liên quan đến việc vận hành ứng dụng (ví dụ: thư mục `node_modules` khi phát triển local, thư mục `.git`, video test, log files hệ thống), Docker theo mặc định vẫn sẽ đóng gói toàn bộ và gửi chúng đi.

**Hậu quả để lại:**
1. **Chậm quá trình build:** Việc nén dữ liệu khổng lồ và gửi qua API (hoặc mạng nếu daemon ở xa) tốn rất nhiều thời gian. Đôi khi máy bị treo chỉ vì bước "Sending build context...".
2. **Lãng phí tài nguyên:** Daemon phải tiêu tốn RAM, CPU và ổ cứng tạm thời để tiếp nhận và giải nén khối lượng dữ liệu dư thừa này.
3. **Image phình to không kiểm soát:** Nếu bên trong `Dockerfile` của bạn có sử dụng lệnh `COPY . .`, toàn bộ rác và thư mục không cần thiết trong context sẽ được bê nguyên xi vào image cuối cùng.

---

## 3. File `.dockerignore` là gì?

### Cách loại bỏ hành lý thừa trước chuyến đi
`.dockerignore` là một file văn bản trơn được đặt ở **thư mục gốc** của ngữ cảnh build (cùng cấp với file `Dockerfile`). Nhiệm vụ của nó là nói cho Docker biết những file hoặc thư mục nào **không được phép** đưa vào Build Context.

Cú pháp của nó tương tự và kế thừa phong cách của file `.gitignore`.

### Cú pháp thông dụng để cấu hình

```dockerfile
# 1. Bỏ qua thư mục quản lý mã nguồn Git (rất quan trọng)
.git/

# 2. Bỏ qua tất cả các file có đuôi .log
*.log

# 3. Bỏ qua thư mục dependencies nặng nề của Node.js (sẽ được cài lại trong Docker)
node_modules/

# 4. Bỏ qua thư mục cache biên dịch của Python
__pycache__/
*.pyc

# 5. Bỏ qua các file chứa biến môi trường nhạy cảm, mật khẩu
.env
secret.key
```

> 📌 **Ghi nhớ:** Luôn luôn đưa `.env` và các file chứa secret keys vào file `.dockerignore`. Nếu không, lệnh `COPY . .` vô tình sẽ sao chép toàn bộ mật khẩu của hệ thống vào Docker image. Bất kỳ ai tải được image đó đều có thể trích xuất (extract) và xem thông tin bảo mật.

---

## 4. Ví dụ thực tế với các ngôn ngữ phổ biến

### `.dockerignore` cho project Node.js
```dockerfile
node_modules
npm-debug.log
.DS_Store
.git
.gitignore

# Tốt nhất là bỏ luôn cả Dockerfile nếu app không dùng tới
Dockerfile
.dockerignore
```

### `.dockerignore` cho project Python
```dockerfile
__pycache__
*.pyc
*.pyo
*.pyd
.Python
env/
venv/
pip-log.txt
.env
```

> ⚠️ **Lưu ý:** Việc đưa chính `Dockerfile` và `.dockerignore` vào trong danh sách bỏ qua (ignore) là một best practice rất tốt nếu bản thân mã nguồn của bạn không cần đọc hoặc parse nội dung của hai file này khi chạy. Điều này giúp ngăn image cache bị vô hiệu hoá mỗi khi bạn sửa nhẹ `Dockerfile`.

---

## 5. So sánh `.dockerignore` vs `.gitignore`

Người mới làm quen với Docker thường dễ nhầm lẫn hai file cấu hình này vì cú pháp giống hệt nhau.

- **`.gitignore`**: Sinh ra dành cho Git. Nó chặn Git theo dõi sự thay đổi của các file không cần thiết và ngăn không cho đẩy (push) lên remote repository (như GitHub/GitLab).
- **`.dockerignore`**: Sinh ra dành riêng cho Docker. Nó chặn Docker thu thập các file để gửi qua cho Docker Daemon khi bạn ra lệnh build image.

Một file hoàn toàn có thể có mặt trong `.gitignore` nhưng không có trong `.dockerignore` và ngược lại. 
- *Ví dụ:* Bạn không muốn đẩy file cấu hình `.env.development` lên Git (phải ở trong `.gitignore`), nhưng ở một vài tình huống build đặc thù tại local, bạn lại CẦN sao chép file đó vào Docker image (thì lại không đưa vào `.dockerignore`).

---

## 6. Lỗi thường gặp và Best Practices

### Các lỗi phổ biến:
- **Quên tạo file `.dockerignore`:** Đây là lỗi kinh điển nhất. Chạy `docker build` và thấy hệ thống "đơ" tới vài phút chỉ để gửi vài GB thư mục ảo `venv` của Python hay `node_modules` cho daemon.
- **Vô tình COPY secret files:** Không đưa file credentials, SSH keys vào `.dockerignore` dẫn đến hậu quả nghiêm trọng về mặt an toàn thông tin (Security).

### Best Practices (Thực hành tốt nhất):
- **Luôn bắt đầu mọi project Docker với file `.dockerignore`.** Đừng chờ ứng dụng lớn lên mới tạo, hãy làm nó trước cả khi viết `Dockerfile`.
- **Giữ build context nhỏ nhất có thể.** Chỉ gửi những nguyên liệu mà "đầu bếp" thực sự cần thiết.
- **Sử dụng cú pháp ngoại trừ (`!`):** Khi thư mục dự án của bạn có quá nhiều rác, thay vì chặn từng cái, hãy chặn tất cả (`*`) và dùng `!` để chỉ cho phép những thứ bạn cần.

---

## Bài tập

1. Tạo một thư mục mới có tên `docker-context-test`. Bên trong, tạo một file giả có kích thước rất lớn (khoảng 500MB) bằng công cụ có sẵn (như `fsutil` trên Windows hoặc `dd` trên Linux).
2. Tạo một `Dockerfile` cực kỳ đơn giản với nội dung: `FROM alpine:latest`. Sau đó chạy lệnh `docker build -t test-context .` và quan sát ở màn hình terminal xem dòng chữ `Sending build context` tốn thời gian mất bao lâu và hiển thị kích thước gửi là bao nhiêu.
3. Tiếp tục bài 1, bạn hãy tạo thêm một file tên là `.dockerignore` ở cùng thư mục và điền tên của file dung lượng lớn 500MB vào bên trong đó. Chạy lại lệnh `docker build -t test-context .`. Hãy so sánh tốc độ gửi ngữ cảnh giữa lần build này và lần trước.
4. Bạn hãy lên ý tưởng và tự cấu hình một file `.dockerignore` nâng cao có sử dụng ký tự loại trừ (`!`). Yêu cầu: File này phải bỏ qua tất cả tài liệu có đuôi `.md` trong hệ thống dự án, **ngoại trừ** duy nhất file `README.md` chính để đưa vào bên trong image giới thiệu cho người dùng cuối.

---

## Tiếp theo
→ [Bài 7: Volumes & Bind Mounts](./7_volumes.md)
