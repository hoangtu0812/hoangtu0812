# Bài 15: Docker Registry

## Mục tiêu
- Hiểu được Docker Registry là gì và vai trò của nó trong hệ sinh thái Docker.
- Biết cách sử dụng Docker Hub để lưu trữ và phân phối Docker images.
- Nắm vững cách tự host một Private Docker Registry tại nội bộ (on-premise).
- Hiểu rõ chiến lược đánh dấu phiên bản (Tagging Strategy) cho image.
- Nắm được các lựa chọn Registry phổ biến thay thế Docker Hub.

---

## 1. Docker Registry là gì?

### 1.1 Khái niệm cơ bản
Docker Registry là một hệ thống dùng để lưu trữ và phân phối các Docker images. Nếu coi Docker image như một ứng dụng đã được đóng gói sẵn, thì Docker Registry giống như một "chợ ứng dụng" (như App Store hoặc Google Play), nơi bạn có thể tải về (pull) hoặc tải lên (push) các ứng dụng của mình.

> 💡 **Tip:** Đừng nhầm lẫn giữa **Registry** và **Repository**. Registry là toàn bộ hệ thống (ví dụ: Docker Hub), còn Repository là một tập hợp các images có cùng tên nhưng khác thẻ phiên bản (ví dụ: repository `ubuntu` chứa các tag `20.04`, `22.04`).

### 1.2 Official images vs User images
- **Official images:** Là các images chính thức do Docker hoặc các tổ chức uy tín bảo trì (ví dụ: `nginx`, `redis`, `ubuntu`). Chúng thường an toàn, được tối ưu và cập nhật thường xuyên.
- **User images:** Là các images do người dùng cộng đồng hoặc bạn tự tạo ra. Định dạng tên thường là `username/image_name`.

---

## 2. Làm việc với Docker Hub

Docker Hub là public registry mặc định và phổ biến nhất của Docker. Nó cung cấp các kho lưu trữ công khai (public) miễn phí và các kho riêng tư (private) có giới hạn.

### 2.1 Đăng nhập và Đăng xuất
Để push image lên Docker Hub, bạn cần có tài khoản và đăng nhập thông qua CLI:

```bash
# Đăng nhập vào Docker Hub
# Hệ thống sẽ yêu cầu bạn nhập Username và Password
docker login

# Đăng xuất khỏi hệ thống
docker logout
```

**Kết quả:**
```
Login Succeeded
```

### 2.2 Đặt tên (Naming Convention)
Để push một image, tên của nó phải tuân thủ định dạng: `[tên_registry]/[username]/[tên_repository]:[tag]`. Nếu sử dụng Docker Hub, phần `[tên_registry]` có thể bỏ qua.

```bash
# Đổi tên một image hiện có để chuẩn bị push
docker tag my-app:latest hoangtu/my-app:v1.0
```

### 2.3 Push và Pull Image
Sau khi đã đăng nhập và tag image, bạn có thể đẩy nó lên Docker Hub:

```bash
# Đẩy image lên Docker Hub
docker push hoangtu/my-app:v1.0

# Tải image từ Docker Hub về máy
docker pull hoangtu/my-app:v1.0
```

> 📌 **Ghi nhớ:** Docker Hub cũng hỗ trợ **Automated Builds**, cho phép tự động build image mỗi khi có code mới được push lên GitHub hoặc Bitbucket, giúp tích hợp tốt với luồng làm việc tự động.

---

## 3. Tự Host Private Registry

Trong môi trường doanh nghiệp, bạn thường không muốn mã nguồn/ứng dụng (đã đóng gói thành image) nằm trên mạng public. Docker cho phép bạn tự chạy một Registry cục bộ bằng image `registry:2`.

### 3.1 Khởi chạy Private Registry
```bash
# Chạy một private registry ở cổng 5000
docker run -d -p 5000:5000 --name my-registry registry:2
```

**Kết quả:**
```
d1f3b2... (ID của container đang chạy registry)
```

### 3.2 Tương tác với Private Registry
Bạn dùng địa chỉ của registry cục bộ làm tiền tố (prefix) cho tên image.

```bash
# Đánh tag cho image trỏ tới registry cục bộ
docker tag my-app:latest localhost:5000/my-app:latest

# Push image vào registry cục bộ
docker push localhost:5000/my-app:latest

# Pull image từ registry cục bộ (sau khi đã xóa ở máy)
docker image rm localhost:5000/my-app:latest
docker pull localhost:5000/my-app:latest
```

> ⚠️ **Lưu ý:** Theo mặc định, Docker yêu cầu Registry phải sử dụng HTTPS. Nếu Registry của bạn không có chứng chỉ SSL/TLS (như chạy nội bộ qua địa chỉ IP), bạn phải cấu hình **insecure registry** trong file `daemon.json` của Docker để cho phép kết nối HTTP.

---

## 4. Chiến lược đánh dấu phiên bản (Tagging Strategy)

Quản lý tag hợp lý giúp bạn tránh được sự cố triển khai sai phiên bản ứng dụng.

### 4.1 Semantic Versioning và "latest"
- **Semantic versioning (SemVer):** Sử dụng các số định danh `MAJOR.MINOR.PATCH` (ví dụ: `v1.2.3`). Bạn nên gắn nhiều tag cho một phiên bản để dễ quản lý. Ví dụ, khi phát hành bản `1.2.3`, hãy đánh tag cả `1.2` và `1.2.3`.
- **Không phụ thuộc vào `latest`:** Tag `latest` chỉ là tag mặc định, không mang ý nghĩa là phiên bản "mới nhất" thực sự. Nó dễ bị ghi đè và làm mất tính nhất quán trên môi trường Production.

### 4.2 Git SHA tags cho CI/CD
Trong quy trình CI/CD, thay vì dùng số phiên bản, người ta thường tag image bằng mã băm của commit (Git SHA) để đảm bảo tính duy nhất tuyệt đối và dễ dàng truy vết (traceability).
Ví dụ: `hoangtu/my-app:a1b2c3d`

---

## 5. docker manifest và Các Registry thay thế

### 5.1 Lệnh docker manifest (Giới thiệu)
Ngày nay, các thiết bị sử dụng kiến trúc CPU khác nhau (như x86_64, ARM64 trên Apple Silicon). `docker manifest` là công cụ giúp bạn tạo và quản lý "multi-architecture images" – một tag duy nhất (ví dụ `ubuntu:latest`) nhưng hỗ trợ chạy trên nhiều loại chip khác nhau mà không cần cấu hình thêm.

### 5.2 Các Registry thay thế phổ biến
Ngoài Docker Hub, các nhà cung cấp đám mây lớn và nền tảng Git đều cung cấp giải pháp Registry (thường được gọi là Container Registry) mạnh mẽ:
- **GitHub Container Registry (ghcr.io):** Tích hợp cực kỳ tốt với GitHub Actions.
- **AWS ECR (Elastic Container Registry):** Dịch vụ của Amazon Web Services.
- **Google Artifact Registry:** Dành cho hệ sinh thái Google Cloud.
- **Azure ACR (Azure Container Registry):** Của Microsoft Azure.

---

## Bài tập
1. **Đăng nhập và Push cơ bản**: Tạo một tài khoản Docker Hub (nếu chưa có). Đăng nhập qua CLI, tải image `alpine`, tag nó với tên `[username_của_bạn]/my-alpine:v1`, và push lên Docker Hub của bạn.
2. **Private Registry**: Khởi chạy một container từ image `registry:2` ở cổng 5000. Lấy một image bất kỳ trên máy, tag và push vào private registry này. Xác minh bằng cách xóa image cũ và pull lại từ `localhost:5000`.
3. **Multi-tagging**: Đánh dấu một image tự build với các tag: `v2`, `v2.1`, và `v2.1.4`. Đẩy tất cả chúng lên private registry của bạn.
4. **Khám phá Manifest**: Chạy lệnh `docker manifest inspect alpine:latest` và quan sát kết quả để hiểu cách một image hỗ trợ cả kiến trúc `amd64` và `arm64`. (Yêu cầu bật tính năng experimental nếu dùng bản cũ).

---

## Tiếp theo
→ [Bài 16: Logging & Monitoring](./16_logging_monitoring.md)
