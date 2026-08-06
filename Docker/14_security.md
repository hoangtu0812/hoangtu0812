# Bài 14: Bảo mật Container

## Mục tiêu
- Hiểu rõ tầm quan trọng của việc bảo mật trong kiến trúc và môi trường Docker.
- Nắm vững cách xây dựng image và chạy container với quyền non-root.
- Biết cách sử dụng các công cụ để quét lỗ hổng bảo mật (vulnerabilities) cho image.
- Áp dụng các phương pháp quản lý secrets an toàn.
- Sử dụng các tính năng giới hạn tài nguyên, capabilities và read-only filesystem để bảo vệ runtime.
- Nắm được các "Best Practices" để triển khai container một cách an toàn trên môi trường production.

---

## 1. Tại sao bảo mật container quan trọng?

Trong kiến trúc truyền thống dựa trên máy ảo (Virtual Machine - VM), mỗi máy ảo có một hệ điều hành (Kernel) riêng biệt. Điều này cung cấp mức độ cô lập cực kỳ cao. Tuy nhiên, kiến trúc của Docker lại chia sẻ chung Kernel của hệ điều hành máy chủ (host OS) với tất cả các container đang chạy trên đó.

Sự chia sẻ này mang lại hiệu suất vượt trội và tốc độ khởi động nhanh, nhưng đồng thời cũng mở ra một rủi ro bảo mật lớn: **Nếu một container bị xâm nhập và kẻ tấn công có thể "vượt rào" (container breakout), chúng có khả năng tác động trực tiếp lên Kernel của host OS và chiếm quyền điều khiển toàn bộ hệ thống.**

Do đó, bảo mật container không chỉ là việc bảo vệ bản thân ứng dụng bên trong, mà còn là việc xây dựng các lớp phòng thủ (defense in depth) nhằm cô lập thiệt hại nếu có sự cố xảy ra.

> 💡 **Tip:** Bảo mật container là một quy trình liên tục. Bạn phải bảo mật từ lúc viết mã (Code), đóng gói (Build), phân phối (Registry) cho đến khi vận hành (Run).

---

## 2. Chạy container với user non-root

Theo mặc định, tiến trình bên trong container sẽ được chạy với quyền `root`. Mặc dù `root` bên trong container bị giới hạn nhiều quyền so với `root` trên host, nhưng việc chạy ứng dụng bằng quyền cao nhất vẫn tiềm ẩn nguy cơ bảo mật nghiêm trọng.

### Nguy hiểm của việc dùng root
Nếu ứng dụng của bạn có lỗ hổng thực thi mã từ xa (RCE) và nó đang chạy bằng quyền `root`, kẻ tấn công có thể dễ dàng cài đặt thêm phần mềm độc hại, can thiệp vào các tiến trình khác (nếu cấu hình sai) hoặc đọc các file nhạy cảm.

### Giải pháp: Tạo user và phân quyền trong Dockerfile

Để bảo mật, bạn nên luôn luôn tạo một user riêng với quyền hạn thấp và chỉ thị cho Docker chạy ứng dụng bằng user đó thông qua lệnh `USER`.

```dockerfile
# Sử dụng base image nhỏ gọn và an toàn
FROM node:18-alpine

# Cập nhật các gói phần mềm và tạo một nhóm/user không có quyền root
# Tham số -S dùng để tạo system user/group trong Alpine Linux
RUN apk update && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# Thiết lập thư mục làm việc mặc định
WORKDIR /app

# Copy các file config và chuyển quyền sở hữu cho user vừa tạo
COPY --chown=appuser:appgroup package.json package-lock.json ./

# Cài đặt các thư viện (vẫn chạy dưới quyền root vì chưa gọi USER)
RUN npm install --production

# Copy toàn bộ mã nguồn và phân quyền
COPY --chown=appuser:appgroup . .

# Chuyển đổi sang sử dụng user có quyền hạn thấp
USER appuser

# Expose port ứng dụng
EXPOSE 3000

# Chạy ứng dụng
CMD ["node", "app.js"]
```

**Kết quả:**
```
Khi bạn chạy container từ image này, nếu truy cập vào shell của container (`docker exec -it <container> sh`) và gõ lệnh `whoami`, kết quả trả về sẽ là `appuser` thay vì `root`.
```

> ⚠️ **Lưu ý:** Bạn nên sử dụng các cổng lớn hơn 1024 cho ứng dụng (ví dụ 3000 hoặc 8080) vì trong Linux, một user non-root không được phép bind vào các cổng dưới 1024 (như port 80 hoặc 443) trừ khi được cấp quyền đặc biệt.

---

## 3. Quét lỗ hổng image (Image Scanning)

Việc sử dụng các base image có sẵn trên Docker Hub là rất tiện lợi, nhưng chúng có thể chứa các thư viện phiên bản cũ bị dính lỗ hổng bảo mật (CVEs). 

### Sử dụng Docker Scout
Bắt đầu từ các phiên bản Docker Desktop mới, Docker cung cấp công cụ `docker scout` giúp phân tích và hiển thị chi tiết các lỗ hổng bên trong image.

```bash
# Quét image và hiển thị danh sách các CVE
docker scout cves nginx:latest

# Phân tích các lỗ hổng và gợi ý phiên bản base image tốt hơn
docker scout recommendations nginx:latest
```

**Kết quả:**
```
Output sẽ hiển thị danh sách các gói phần mềm bị lỗi, mức độ nghiêm trọng (CRITICAL, HIGH, MEDIUM, LOW) và các phiên bản đã được vá lỗi để bạn có thể nâng cấp.
```

### Các công cụ quét bên thứ ba
Bên cạnh Docker Scout, hệ sinh thái container còn có nhiều công cụ nguồn mở và thương mại mạnh mẽ khác:
- **Trivy (Aqua Security)**: Công cụ quét rất nhanh, nhẹ và thường được tích hợp vào các pipeline CI/CD (như GitHub Actions, GitLab CI).
- **Snyk**: Cung cấp nền tảng toàn diện để quét mã nguồn và container image.

---

## 4. Quản lý Secrets an toàn

Secrets là các thông tin nhạy cảm như Mật khẩu cơ sở dữ liệu, API Keys, TLS Certificates. **Một nguyên tắc bất di bất dịch: KHÔNG BAO GIỜ hardcode secrets vào mã nguồn hoặc viết trực tiếp vào Dockerfile.**

Nếu bạn dùng `ENV DB_PASS=secret` trong Dockerfile, bất kỳ ai có quyền truy cập vào image (thông qua lệnh `docker inspect`) đều có thể đọc được mật khẩu này.

### Sử dụng file .env và Bind Mount
Cách phổ biến trong môi trường phát triển là sử dụng biến môi trường được truyền vào khi chạy.

```bash
# Sử dụng file .env khi chạy container
docker run -d --name myapp --env-file ./config/.env myapp:latest
```

### Docker Secrets (Trong Swarm mode)
Nếu bạn triển khai ứng dụng bằng Docker Swarm, bạn có thể tận dụng tính năng Docker Secrets để truyền dữ liệu an toàn. Dữ liệu sẽ được mã hóa và chỉ cấp cho các container thực sự cần nó.

```bash
# Tạo một secret từ chuỗi ký tự
echo "my-super-secret-password" | docker secret create db_pass -

# Tạo service và gán secret vào
docker service create --name mydb --secret db_pass mysql:8
```
Bên trong container, secret này sẽ xuất hiện dưới dạng một file văn bản tại đường dẫn `/run/secrets/db_pass`, Ứng dụng có thể đọc nội dung file này để lấy mật khẩu.

---

## 5. Tăng cường bảo mật lúc chạy (Runtime Security)

Bảo vệ image là chưa đủ, bạn cần siết chặt quyền hạn của container khi nó đang hoạt động.

### Read-only filesystem
Kẻ tấn công thường cố gắng ghi các tập lệnh độc hại (webshell, malware) vào ổ đĩa của container. Chạy container với hệ thống file chỉ đọc (read-only) sẽ ngăn chặn điều này.

```bash
# Mount hệ thống file gốc là read-only.
# Cấp một tmpfs (thư mục tạm lưu trên RAM) cho các thư mục bắt buộc phải ghi dữ liệu.
docker run -d --read-only --tmpfs /tmp myapp:latest
```

### Giới hạn tài nguyên
Nếu không có giới hạn, một container bị lỗi hoặc bị tấn công có thể tiêu thụ toàn bộ CPU và RAM của host, gây ra tình trạng từ chối dịch vụ (DoS) cho toàn bộ hệ thống.

```bash
# Giới hạn RAM ở mức 512MB và chỉ sử dụng tối đa 1.5 CPU cores
docker run -d --memory="512m" --cpus="1.5" myapp:latest
```

### Tránh dùng cờ --privileged
Cờ `--privileged` sẽ cấp cho container quyền truy cập trực tiếp vào tất cả các thiết bị phần cứng (devices) trên host và nới lỏng rất nhiều cơ chế bảo mật của kernel. 
> 📌 **Ghi nhớ:** Tuyệt đối **KHÔNG** sử dụng cờ này trừ khi bạn đang chạy các tác vụ liên quan đến cấp thấp của hệ thống (như chạy Docker in Docker hoặc các công cụ quản lý mạng).

### Quản lý Capabilities
Thay vì chạy với quyền `root` đầy đủ, Linux chia quyền root thành các đặc quyền nhỏ lẻ gọi là "Capabilities". Bạn có thể tuân thủ nguyên tắc quyền tối thiểu (Least Privilege) bằng cách xóa tất cả các quyền và chỉ cấp lại những quyền thực sự cần thiết.

```bash
# Xóa bỏ tất cả quyền đặc biệt (cap-drop=ALL)
# Chỉ thêm lại quyền bind cổng mạng dưới 1024 (NET_BIND_SERVICE)
docker run -d --cap-drop=ALL --cap-add=NET_BIND_SERVICE nginx:latest
```

---

## 6. Docker Content Trust (DCT)

Khi bạn tải (pull) một image từ Internet, làm sao để biết image đó không bị giả mạo? Docker Content Trust (DCT) sử dụng chữ ký điện tử để đảm bảo tính toàn vẹn và xác thực nguồn gốc của image.

Khi bật DCT, Docker sẽ từ chối chạy hoặc tải bất kỳ image nào không có chữ ký hợp lệ từ nhà phát triển.

```bash
# Bật tính năng DCT trên máy của bạn
export DOCKER_CONTENT_TRUST=1

# Lệnh này sẽ thực hiện xác minh chữ ký trước khi tải image
docker pull ubuntu:latest
```

---

## 7. Best Practices Checklist (Danh sách kiểm tra thực hành tốt nhất)

Để đảm bảo container của bạn được bảo mật tối đa trước khi đưa lên môi trường thực tế (Production), hãy kiểm tra danh sách sau:

- [ ] **Sử dụng Minimal Base Images**: Sử dụng các image càng nhỏ càng tốt như `alpine`, `distroless` hoặc `scratch` để giảm diện tích tấn công (Attack Surface).
- [ ] **Sử dụng Multi-stage builds**: Không bao giờ đưa các công cụ biên dịch (compilers, build tools) vào image cuối cùng.
- [ ] **Không chạy ứng dụng bằng root**: Luôn sử dụng lệnh `USER` với tài khoản có quyền hạn thấp.
- [ ] **Giới hạn tài nguyên**: Luôn thiết lập cấu hình `--memory` và `--cpus`.
- [ ] **Không lưu trữ Secrets trong image**: Sử dụng biến môi trường hoặc hệ thống quản lý Secret chuyên dụng (như HashiCorp Vault, AWS Secrets Manager).
- [ ] **Kiểm tra lỗ hổng thường xuyên**: Quét image trong quy trình CI/CD trước khi push lên registry.
- [ ] **Read-only**: Kích hoạt `--read-only` flag nếu ứng dụng của bạn không có nhu cầu ghi dữ liệu vĩnh viễn lên ổ cứng.

---

## Bài tập

1. **Viết Dockerfile an toàn**: Hãy lấy một Dockerfile ứng dụng Python hoặc NodeJS bạn đã viết ở các bài trước. Chỉnh sửa nó để thêm một system user tên là `worker`, đổi quyền sở hữu các file mã nguồn cho `worker`, và thêm lệnh `USER worker` trước lệnh `CMD`.
2. **Giới hạn tài nguyên**: Khởi chạy một container Redis, đặt tên là `secure-redis`, giới hạn RAM tối đa là 128MB và cấu hình cho phép sử dụng tối đa 0.5 CPU.

3. **Quét bảo mật Image**: 
   - Sử dụng lệnh `docker scout cves node:14` để kiểm tra một phiên bản Node.js cũ.
   - Sau đó kiểm tra phiên bản mới `docker scout cves node:20-alpine`. So sánh số lượng lỗ hổng (Vulnerabilities) được tìm thấy giữa hai phiên bản này.

4. **Read-only và Capabilities nâng cao**: 
   - Khởi chạy một container Nginx với tùy chọn `--read-only`. 
   - Nginx yêu cầu ghi dữ liệu vào một số thư mục tạm. Hãy sử dụng tùy chọn `--tmpfs` để cấp quyền ghi tạm thời cho các thư mục `/var/cache/nginx`, `/var/run`, và `/tmp`.
   - Kết hợp thêm cấu hình: Xóa tất cả các quyền root (capabilities) bằng `--cap-drop=ALL` và chỉ cấp lại quyền `NET_BIND_SERVICE`.
   - Kiểm tra xem trang web Nginx có hoạt động bình thường không.

---

## Tiếp theo
→ [Bài 15: Docker Registry](./15_registry.md)
