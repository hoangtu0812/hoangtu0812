# Bài 1: Giới thiệu Docker & Cài đặt

## Mục tiêu
- Hiểu được khái niệm về Container và Containerization.
- Phân biệt được sự khác nhau giữa Virtual Machine (Máy ảo) và Container.
- Nắm bắt cơ bản về Docker, Docker Engine, Docker Desktop và Docker Hub.
- Cài đặt thành công Docker trên máy tính cá nhân và chạy thử nghiệm container đầu tiên.

## 1. Container là gì? Tại sao cần Containerization?

### Container là gì?
Container (vùng chứa) là một đơn vị phần mềm tiêu chuẩn hóa, đóng gói mã nguồn và tất cả các phụ thuộc (dependencies) của nó để ứng dụng có thể chạy nhanh chóng và đáng tin cậy từ môi trường máy tính này sang môi trường máy tính khác.

### Tại sao cần Containerization?
Trước đây, khi phát triển ứng dụng, một vấn đề rất phổ biến là "Nó chạy được trên máy của tôi nhưng lại lỗi trên máy chủ" (It works on my machine!). Nguyên nhân là do sự khác biệt về môi trường, hệ điều hành, phiên bản thư viện giữa máy tính của lập trình viên và máy chủ thực tế.
Containerization (Công nghệ vùng chứa) giải quyết triệt để vấn đề này bằng cách đóng gói ứng dụng cùng mọi thứ nó cần vào một container độc lập. Bạn có thể mang container này chạy ở bất cứ đâu: laptop cá nhân, máy chủ test, hay trên cloud (AWS, Google Cloud, Azure) mà không sợ lỗi tương thích môi trường.

## 2. So sánh Virtual Machine (VM) và Container

| Tiêu chí | Virtual Machine (Máy ảo) | Container |
|---|---|---|
| **Cấu trúc** | Bao gồm cả một Hệ điều hành khách (Guest OS) hoàn chỉnh, chạy trên Hypervisor. | Chỉ đóng gói ứng dụng và thư viện, chia sẻ Hệ điều hành máy chủ (Host OS). |
| **Kích thước** | Rất lớn (thường tính bằng Gigabyte). | Rất nhỏ (thường tính bằng Megabyte). |
| **Khởi động** | Chậm (mất vài phút để khởi động hệ điều hành). | Nhanh (gần như ngay lập tức, chỉ trong vài giây). |
| **Hiệu suất** | Tiêu tốn nhiều tài nguyên do phải chạy nhiều OS. | Tiết kiệm tài nguyên, tận dụng tối đa phần cứng. |
| **Cách ly (Isolation)** | Cách ly hoàn toàn ở mức phần cứng ảo hóa. | Cách ly ở mức tiến trình (process) trên hệ điều hành. |

## 3. Docker, Docker Engine và Docker Desktop

### Docker là gì?
Docker là một nền tảng phần mềm mã nguồn mở cho phép các lập trình viên xây dựng, triển khai và chạy các ứng dụng trong các container một cách dễ dàng. Nó đã trở thành tiêu chuẩn công nghiệp (de facto standard) cho công nghệ container.

### Docker Engine
Đây là phần cốt lõi của Docker, là một ứng dụng client-server bao gồm:
- **Server (Docker daemon):** Một tiến trình chạy nền (background process) có tên là `dockerd`, quản lý các đối tượng Docker (image, container, network, volume).
- **REST API:** Cung cấp giao diện để các chương trình khác giao tiếp với Docker daemon.
- **Client (Docker CLI):** Công cụ dòng lệnh `docker` mà chúng ta sử dụng để gõ các lệnh điều khiển Docker daemon.

### Docker Desktop
Là một ứng dụng dành cho Windows và Mac (gần đây có cả Linux), cung cấp giao diện đồ họa (GUI) trực quan và đóng gói sẵn Docker Engine, Docker CLI, Docker Compose và các công cụ khác để bạn có thể bắt đầu sử dụng Docker một cách nhanh nhất mà không cần cấu hình phức tạp.

## 4. Docker Hub là gì?
Docker Hub là một dịch vụ registry do Docker cung cấp để tìm kiếm và chia sẻ các Docker image. Nó giống như GitHub nhưng dành cho Docker images. Tại đây, bạn có thể tìm thấy các image chính thức (official) của Ubuntu, Nginx, MySQL, Python,... hoặc chia sẻ image do chính bạn tạo ra cho cộng đồng.

## 5. Hướng dẫn cài đặt Docker

### Trên Windows / MacOS
Cách tốt nhất là cài đặt **Docker Desktop**:
1. Truy cập trang chủ: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
2. Tải về file cài đặt phù hợp (Windows hoặc Mac).
3. Chạy file cài đặt, để các tùy chọn mặc định (trên Windows, hãy đảm bảo WSL 2 được bật).
4. Khởi động lại máy tính (nếu được yêu cầu) và mở ứng dụng Docker Desktop.

### Trên Linux (Ubuntu/Debian)
Bạn có thể cài đặt Docker Engine qua dòng lệnh:

```bash
# Cập nhật danh sách các gói tin
sudo apt-get update

# Cài đặt các gói phụ thuộc
sudo apt-get install ca-certificates curl gnupg

# Thêm khóa GPG chính thức của Docker
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Thiết lập repository
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Cài đặt Docker Engine
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

## 6. Kiểm tra (Verify) cài đặt

Mở Terminal (trên Mac/Linux) hoặc PowerShell (trên Windows) và chạy các lệnh sau:

Kiểm tra phiên bản Docker:
```powershell
docker --version
```
**Kết quả mong đợi:** `Docker version 24.0.5, build ced0996` (Số phiên bản có thể khác)

Chạy thử nghiệm một container cơ bản:
```powershell
docker run hello-world
```
Lệnh này sẽ tải image `hello-world` từ Docker Hub và chạy nó trong một container.
**Kết quả mong đợi:** 
Bạn sẽ thấy một thông báo bắt đầu bằng:
`Hello from Docker!`
`This message shows that your installation appears to be working correctly.`

## Bài tập

1. **Kiểm tra phiên bản**: Mở terminal/powershell, chạy lệnh kiểm tra phiên bản Docker hiện tại trên máy của bạn và ghi chú lại.
2. **Khám phá Docker Desktop**: Nếu bạn dùng Windows/Mac, hãy mở giao diện Docker Desktop, tìm tab "Containers" và "Images". Bạn có thấy image `hello-world` xuất hiện không?
3. **Chạy lại hello-world**: Xóa output cũ trên terminal và chạy lại lệnh `docker run hello-world`. Bạn có thấy tốc độ thực thi nhanh hơn lần đầu tiên không? (Gợi ý: Lần thứ hai không cần tải lại image từ Docker Hub).
4. **Tìm kiếm Image**: Truy cập trang web Docker Hub (hub.docker.com), tìm kiếm image `nginx` và đọc lướt qua phần giới thiệu (description) của nó.

## Tiếp theo
→ [Bài 2: Kiến trúc Docker & Các thành phần cốt lõi](./2_architecture.md)
