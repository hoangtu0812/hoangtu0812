# Bài 11: Docker Compose cơ bản

## Mục tiêu
- Hiểu được Docker Compose là gì và tại sao chúng ta cần sử dụng nó.
- Nắm vững cấu trúc cơ bản của một file `compose.yaml` (hoặc `docker-compose.yml`).
- Biết cách cấu hình các services, volumes, networks, và environment variables.
- Thành thạo các câu lệnh Docker Compose cơ bản để quản lý vòng đời của một ứng dụng đa container.
- Thực hành triển khai một ứng dụng web kết nối với cơ sở dữ liệu (PostgreSQL/MySQL).

---

## 1. Docker Compose là gì?

Docker Compose là một công cụ giúp bạn định nghĩa và chạy các ứng dụng Docker có nhiều container (multi-container Docker applications). Với Compose, bạn sử dụng một file YAML (thường là `compose.yaml` hoặc `docker-compose.yml`) để cấu hình tất cả các services mà ứng dụng của bạn cần.

> 💡 **Tip:** Từ phiên bản Docker Desktop hiện đại, Docker Compose (v2) đã được tích hợp sẵn. Bạn không cần cài đặt thêm và có thể dùng lệnh `docker compose` thay vì `docker-compose` như trước đây.

### Tại sao cần Docker Compose?
Hãy tưởng tượng bạn có một ứng dụng web cần chạy một container cho Node.js (Web Server) và một container cho PostgreSQL (Database). Nếu dùng lệnh `docker run`, bạn sẽ phải:
1. Tạo một network để chúng giao tiếp.
2. Chạy container Database kèm theo các biến môi trường và volume.
3. Chạy container Web Server, gắn vào cùng network, map port.

Việc gõ đi gõ lại hàng tá lệnh `docker run` với các tham số dài dòng rất dễ gây lỗi và khó quản lý. Docker Compose giải quyết vấn đề này bằng cách lưu trữ toàn bộ cấu hình vào một file duy nhất. Bạn chỉ cần 1 lệnh để khởi động toàn bộ hệ thống!

---

## 2. Cấu trúc file compose.yaml

Một file `compose.yaml` tiêu chuẩn thường bao gồm các thành phần chính sau: `version` (tùy chọn trong bản mới), `services`, `networks`, và `volumes`.

### Các thành phần cấu hình cơ bản cho một Service:
- `image`: Hình ảnh Docker sẽ được sử dụng.
- `build`: Đường dẫn chứa Dockerfile nếu bạn muốn tự build image.
- `ports`: Mapping cổng giữa host và container.
- `volumes`: Ánh xạ thư mục/volume từ máy host vào container để lưu trữ dữ liệu bền vững.
- `environment`: Các biến môi trường cần thiết cho container.
- `networks`: Xác định network mà service tham gia (Docker Compose tự động tạo một network mặc định cho tất cả các service nếu không chỉ định).

---

## 3. Ví dụ thực tế: Web App + Database

Dưới đây là một ví dụ thiết lập một ứng dụng Node.js kết nối với Redis.

```yaml
# compose.yaml
name: my-web-app

services:
  web:
    build: .                 # Build image từ Dockerfile trong thư mục hiện tại
    ports:
      - "8080:8080"          # Map port 8080 của host với port 8080 của container
    environment:
      - NODE_ENV=production
      - REDIS_HOST=redis_db  # Trỏ tới service redis bên dưới
    depends_on:
      - redis_db             # Đảm bảo redis chạy trước khi web khởi động

  redis_db:
    image: redis:alpine      # Dùng image redis có sẵn từ Docker Hub
    volumes:
      - redis_data:/data     # Mount volume để lưu trữ dữ liệu Redis bền vững
    ports:
      - "6379:6379"

volumes:
  redis_data:                # Khai báo named volume
```

> 📌 **Ghi nhớ:** Docker Compose sẽ tự động tạo một mạng ảo (network) và đặt cả `web` lẫn `redis_db` vào đó. Các container có thể gọi nhau thông qua tên service (ở đây `web` gọi `redis_db` dễ dàng qua hostname `redis_db`).

---

## 4. Các lệnh Docker Compose cơ bản

Dưới đây là các lệnh bạn sẽ sử dụng hằng ngày với Docker Compose:

### Khởi động hệ thống
```bash
# Khởi động tất cả các service (chạy ở foreground, in log ra màn hình)
docker compose up

# Khởi động và chạy ngầm (detached mode)
docker compose up -d

# Bắt buộc build lại image trước khi chạy (nếu có thay đổi source code)
docker compose up -d --build
```

**Kết quả:**
```
[+] Running 3/3
 ✔ Network my-web-app_default      Created
 ✔ Container my-web-app-redis_db-1 Started
 ✔ Container my-web-app-web-1      Started
```

### Dừng và xóa hệ thống
```bash
# Dừng các container nhưng không xóa chúng
docker compose stop

# Dừng và xóa container, network (mặc định)
docker compose down

# Dừng, xóa container, network VÀ xóa luôn cả volumes (-v)
docker compose down -v
```

### Theo dõi trạng thái
```bash
# Xem danh sách các container đang chạy bởi compose
docker compose ps

# Xem log của tất cả các service
docker compose logs

# Theo dõi log (follow) của một service cụ thể
docker compose logs -f web
```

### Tương tác với container
```bash
# Chạy một lệnh (ví dụ bash/sh) bên trong một container đang chạy
docker compose exec web sh

# Build lại các image được định nghĩa trong file
docker compose build
```

---

## Bài tập

1. **Khởi tạo Compose file**: Tạo một file `compose.yaml` chạy 1 service `nginx` (image: `nginx:latest`), map port 80 của máy bạn vào port 80 của container.
2. **Thêm Database**: Bổ sung một service `db` sử dụng image `mysql:8.0`. Thiết lập biến môi trường `MYSQL_ROOT_PASSWORD=secret` và map port 3306.
3. **Sử dụng Volume**: Thêm cấu hình volume để dữ liệu của `db` lưu vào thư mục `./mysql-data` trên máy tính của bạn (bind mount) hoặc dùng named volume. Đảm bảo dữ liệu không bị mất đi khi gõ `docker compose down`.
4. **Thực hành các lệnh**: Chạy ứng dụng bằng `docker compose up -d`. Sử dụng lệnh `docker compose ps` để kiểm tra. Sau đó dùng lệnh `docker compose down -v` để dọn dẹp hệ thống.

---

## Tiếp theo
→ [Bài 12: Docker Compose nâng cao](./12_compose_advanced.md)
