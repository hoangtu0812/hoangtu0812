# Bài 18: CI/CD với Docker

## Mục tiêu
- Hiểu rõ khái niệm CI/CD (Continuous Integration / Continuous Delivery) và lý do Docker lại là giải pháp hoàn hảo cho quá trình này.
- Biết cách thiết lập GitHub Actions để tự động build và push Docker image.
- Sử dụng hiệu quả Multi-stage build và Docker secrets trong môi trường CI.
- Nắm vững các chiến lược đánh thẻ (tagging strategy) cho Docker image.
- Phân biệt được Docker-in-Docker (DinD) và Docker socket mounting.
- Nắm được luồng triển khai (deployment pipeline) từ build, test, push đến deploy.

---

## 1. CI/CD là gì và Tại sao nên dùng Docker?

### Khái niệm CI/CD
- **CI (Continuous Integration - Tích hợp liên tục):** Là quá trình tự động hóa việc gộp code của nhiều lập trình viên vào một nhánh chính. Mỗi lần gộp code, hệ thống sẽ tự động build và chạy test để phát hiện lỗi sớm nhất có thể.
- **CD (Continuous Delivery/Deployment - Phân phối/Triển khai liên tục):** Là quá trình tự động hóa việc đưa ứng dụng đã vượt qua bài test lên môi trường staging hoặc production.

> 💡 **Tip:** Hãy tưởng tượng CI/CD như một dây chuyền lắp ráp ô tô tự động. Mỗi khi có linh kiện mới (code), dây chuyền tự động lắp ráp (build), kiểm tra chất lượng (test), và xuất xưởng (deploy) mà không cần con người can thiệp.

### Tại sao Docker lại lý tưởng cho CI/CD?
Môi trường truyền thống thường gặp lỗi "chạy được trên máy tôi nhưng lỗi trên server" (It works on my machine). Docker giải quyết triệt để vấn đề này bằng cách đóng gói mọi thứ thành một khối thống nhất (image). Docker mang lại:
1. **Sự nhất quán:** Image được test trong CI chính xác là image sẽ chạy trên production.
2. **Khả năng tái lập:** Có thể dễ dàng chạy lại quy trình build mọi lúc mọi nơi với cùng một kết quả.
3. **Triển khai nhanh chóng:** Khởi động container mất chưa tới vài giây so với vài phút của máy ảo.

---

## 2. CI/CD Pipeline Cơ Bản: Build → Test → Push → Deploy

Một đường ống (pipeline) triển khai tiêu chuẩn với Docker thường bao gồm 4 bước:
1. **Build:** Lấy code mới nhất, sử dụng Dockerfile để tạo Docker image.
2. **Test:** Chạy container từ image vừa tạo và thực thi các bài unit test / integration test.
3. **Push:** Đẩy image đã test thành công lên Docker Registry (như Docker Hub, GitHub Container Registry).
4. **Deploy:** Cập nhật server/cluster (ví dụ Kubernetes, Docker Swarm) để sử dụng image mới nhất này.

> 📌 **Ghi nhớ:** Luôn build image trước, test trên chính image đó, sau đó mới push. Không bao giờ build lại (rebuild) ở bước deploy để tránh sự sai lệch.

---

## 3. GitHub Actions với Docker

GitHub Actions là một công cụ CI/CD tích hợp sẵn trong GitHub. Bạn định nghĩa các "workflow" thông qua các file YAML.

### Tự động Build và Push Image

Dưới đây là ví dụ đầy đủ một file `.github/workflows/docker.yml`.

```yaml
# .github/workflows/docker.yml
name: Docker CI/CD Pipeline

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

jobs:
  build-test-push:
    runs-on: ubuntu-latest

    steps:
      # Bước 1: Clone mã nguồn từ repository
      - name: Checkout code
        uses: actions/checkout@v3

      # Bước 2: Thiết lập Docker Buildx (giúp build nhanh hơn và cache tốt hơn)
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      # Bước 3: Đăng nhập vào Docker Hub sử dụng Secrets
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_HUB_USERNAME }}
          password: ${{ secrets.DOCKER_HUB_ACCESS_TOKEN }}

      # Bước 4: Build và chạy Test (dùng target 'test' nếu có multi-stage)
      - name: Build and Test
        run: |
          docker build --target test -t my-app:test .
          docker run my-app:test

      # Bước 5: Build và Push image chính thức lên Docker Hub
      - name: Build and Push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }} # Không push nếu là PR
          tags: mydockerhubuser/my-app:latest,mydockerhubuser/my-app:${{ github.sha }}
          cache-from: type=registry,ref=mydockerhubuser/my-app:buildcache
          cache-to: type=registry,ref=mydockerhubuser/my-app:buildcache,mode=max
```

**Kết quả:**
```
Mỗi khi bạn push code lên nhánh `main`, GitHub Actions sẽ tự động chạy các bước trên.
Bạn có thể xem log chi tiết từng bước trong tab "Actions" của GitHub repository.
Nếu có bất kỳ bài test nào thất bại, quy trình sẽ dừng lại và đánh dấu X đỏ.
```

> ⚠️ **Lưu ý:** Tuyệt đối không hardcode mật khẩu Docker Hub trong file YAML. Hãy sử dụng `Settings > Secrets and variables > Actions` của repository để tạo `DOCKER_HUB_USERNAME` và `DOCKER_HUB_ACCESS_TOKEN`.

---

## 4. Multi-stage Build và Tagging Strategy trong CI/CD

### Tích hợp Multi-stage Build
Multi-stage build cho phép bạn gộp cả quá trình build code và chạy test vào trong cùng một `Dockerfile`. Ở ví dụ YAML phía trên, chúng ta đã dùng cờ `--target test` để chỉ dừng lại ở stage có tên là `test` để chạy kiểm thử.

### Chiến lược đánh thẻ (Tagging Strategy)
Trong CI/CD, việc đặt tag cho image rất quan trọng để quản lý phiên bản:
- **Git SHA:** Thêm tag bằng mã commit của Git (VD: `my-app:a1b2c3d`). Giúp bạn dễ dàng truy vết xem image này được build từ dòng code nào.
- **Branch name:** Tag theo tên nhánh (VD: `my-app:staging`, `my-app:develop`).
- **Semantic version:** Các tag như `v1.0.0`, `v1.0.1` dựa trên Git tags.
- **latest:** Chỉ nên áp dụng tag `latest` cho code trên nhánh `main` hoặc `master`.

---

## 5. Môi trường Docker CI: DinD vs Socket Mounting

Khi chạy CI/CD (như GitLab CI hay Jenkins), các tác vụ (runner) thường chạy bên trong các Docker container. Vậy làm sao để chạy lệnh `docker build` bên trong một Docker container khác? Bạn có 2 cách:

### Docker-in-Docker (DinD)
- Khởi chạy một Docker daemon hoàn toàn độc lập bên trong CI container.
- Ưu điểm: Cách ly hoàn toàn, an toàn hơn, không ảnh hưởng đến host.
- Nhược điểm: Phức tạp, chậm hơn do phải cấu hình lại từ đầu mỗi lần chạy, yêu cầu quyền `privileged`.

### Docker Socket Mounting
- Gắn (mount) trực tiếp file socket của Docker host `/var/run/docker.sock` vào trong CI container.
- Cú pháp: `-v /var/run/docker.sock:/var/run/docker.sock`
- Ưu điểm: Nhanh gọn, chia sẻ image cache với host.
- Nhược điểm: Kém bảo mật, CI container có toàn quyền kiểm soát Docker host.

> Giới thiệu ngắn về GitLab CI: Tương tự như GitHub Actions, GitLab CI sử dụng file `.gitlab-ci.yml`. Để build Docker trong GitLab CI, bạn thường dùng service `docker:dind`.

---

## Bài tập

1. **Khởi tạo Workflow:** Tạo một repository GitHub mới, viết một file `Dockerfile` đơn giản (chạy Nginx hiển thị trang HTML) và tạo file `.github/workflows/docker.yml` để build nó.
2. **Bảo mật Credentials:** Cấu hình GitHub Secrets cho Docker Hub account của bạn và chỉnh sửa workflow ở bài 1 để đăng nhập vào Docker Hub.

3. **Cấu hình Tagging:** Cập nhật workflow để tag Docker image của bạn bằng cả `latest` và biến `${{ github.sha }}`. Kiểm tra trên Docker Hub xem 2 tag này có xuất hiện sau khi push code không.

4. **Xây dựng CI/CD Pipeline hoàn chỉnh:** Viết một ứng dụng Node.js hoặc Python có chứa Unit Test. Viết một `Dockerfile` sử dụng Multi-stage (có target `test`). Cấu hình GitHub Actions để khi tạo Pull Request, nó chỉ chạy test (không push). Khi Merge vào `main`, nó chạy test và push image lên registry.

---

## Tiếp theo
Bài 19: Docker trong Production → [Bài 19: Docker trong Production](./19_production.md)
