# Dự án Capstone: Containerize ứng dụng Multi-service hoàn chỉnh
## Mục tiêu
- Tổng hợp và áp dụng toàn bộ kiến thức Docker đã học từ Bài 1 đến Bài 20.
- Containerize một ứng dụng thực tế với kiến trúc microservices (Frontend, Backend, Database, Cache, Reverse Proxy).
- Áp dụng các Best Practices về bảo mật, tối ưu kích thước image và quản lý cấu hình.
- Xây dựng file `docker-compose.yml` hoàn chỉnh cho môi trường development và production.
- Thiết lập CI/CD cơ bản để tự động hóa việc build và push image.

---

## Kiến trúc dự án

Dự án Capstone của chúng ta là một ứng dụng Web phân tán (Multi-service Web App) bao gồm 5 thành phần chính:
1. **Frontend**: Nginx phục vụ các file tĩnh (React/Vue/Angular build) hoặc HTML/JS thuần.
2. **Backend**: Python FastAPI (hoặc Flask/Go) cung cấp REST API.
3. **Database**: PostgreSQL lưu trữ dữ liệu chính.
4. **Cache**: Redis dùng để caching dữ liệu hoặc xử lý hàng đợi.
5. **Reverse Proxy**: Nginx đóng vai trò là API Gateway để định tuyến traffic từ ngoài vào Frontend và Backend.

### Cấu trúc thư mục

```text
capstone-project/
├── .env                  # Lưu trữ biến môi trường (không commit)
├── .env.example          # Mẫu file cấu hình biến môi trường
├── docker-compose.yml    # File Compose chính để chạy toàn bộ stack
├── reverse-proxy/        # Cấu hình Nginx Reverse Proxy
│   ├── Dockerfile
│   └── nginx.conf
├── frontend/             # Mã nguồn Frontend
│   ├── .dockerignore
│   ├── Dockerfile
│   ├── package.json
│   ├── src/
│   └── nginx.conf        # Cấu hình Nginx nội bộ cho Frontend
├── backend/              # Mã nguồn Backend API
│   ├── .dockerignore
│   ├── Dockerfile
│   ├── requirements.txt
│   └── app/
└── db-init/              # Script khởi tạo Database
    └── init.sql
```

---

## Bước 1: Thiết kế kiến trúc ứng dụng

Trước khi bắt tay vào viết code, hãy cùng nhìn lại cách các thành phần giao tiếp với nhau (xem lại **Bài 15: Docker Compose Networking**).
- User gửi request tới cổng 80 của `reverse-proxy`.
- `reverse-proxy` sẽ định tuyến request tới `frontend` nếu URL là `/`, hoặc tới `backend` nếu URL có tiền tố `/api`.
- `backend` sẽ kết nối với `postgres` (Database) để đọc/ghi dữ liệu và `redis` (Cache) để tối ưu hiệu suất.
- Các container sẽ nằm trong cùng một custom bridge network nội bộ, ngoại trừ port 80 của `reverse-proxy` được publish ra ngoài.

> 💡 **Tip:** Việc tách biệt `reverse-proxy` giúp chúng ta dễ dàng thêm SSL/TLS (HTTPS) hoặc load balancing trong tương lai mà không ảnh hưởng tới các service bên trong.

---

## Bước 2: Viết Dockerfile cho Backend

Backend sử dụng Python (FastAPI). Chúng ta sẽ áp dụng kỹ thuật **Multi-stage Build** (nhắc lại ở **Bài 10**) và chạy với quyền **Non-root User** (nhắc lại ở **Bài 18**) để tối ưu và bảo mật.

**Tạo file `backend/Dockerfile`**:

```dockerfile
# --- Stage 1: Builder ---
FROM python:3.11-slim AS builder

# Thiết lập thư mục làm việc
WORKDIR /app

# Cài đặt các thư viện cần thiết để build (nếu có)
RUN apt-get update && apt-get install -y --no-install-recommends gcc && rm -rf /var/lib/apt/lists/*

# Copy file requirements và cài đặt dependencies vào một thư mục tạm
COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt

# --- Stage 2: Production ---
FROM python:3.11-slim

WORKDIR /app

# Tạo non-root user để tăng tính bảo mật
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# Copy dependencies đã được cài đặt từ Stage 1
COPY --from=builder /root/.local /home/appuser/.local

# Cập nhật PATH để dùng các thư viện Python
ENV PATH=/home/appuser/.local/bin:$PATH
# Đảm bảo Python in ra log ngay lập tức mà không bị buffer
ENV PYTHONUNBUFFERED=1

# Copy mã nguồn Backend vào container
COPY ./app ./app

# Phân quyền cho non-root user
RUN chown -R appuser:appgroup /app

# Chuyển sang sử dụng non-root user
USER appuser

# Cổng mặc định
EXPOSE 8000

# Lệnh khởi chạy ứng dụng (ví dụ dùng uvicorn cho FastAPI)
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

> 📌 **Ghi nhớ:** Luôn cung cấp `.dockerignore` để tránh copy các file không cần thiết (ví dụ `__pycache__`, `.git`, `.env`) vào image.

**File `backend/.dockerignore`**:
```text
__pycache__/
*.pyc
.env
venv/
.pytest_cache/
```

---

## Bước 3: Viết Dockerfile cho Frontend

Với Frontend (ví dụ ứng dụng React), ta sử dụng Node.js để build ở stage 1 và Nginx Alpine để serve các file tĩnh ở stage 2.

**Tạo file `frontend/Dockerfile`**:

```dockerfile
# --- Stage 1: Build React App ---
FROM node:18-alpine AS build

WORKDIR /app

# Lợi dụng Layer Caching: Chỉ copy file package trước
COPY package*.json ./
RUN npm ci

# Copy toàn bộ mã nguồn và build
COPY . .
RUN npm run build

# --- Stage 2: Phục vụ file bằng Nginx ---
FROM nginx:1.25-alpine

# Xóa trang mặc định của Nginx
RUN rm -rf /usr/share/nginx/html/*

# Copy file build từ Stage 1 sang thư mục của Nginx
COPY --from=build /app/build /usr/share/nginx/html

# (Tùy chọn) Copy cấu hình Nginx riêng cho SPA (Single Page Application)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Nginx Alpine mặc định chạy master process bằng root, nhưng worker chạy bằng nginx user.
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

---

## Bước 4: Cấu hình PostgreSQL với Volume và Init script

Như đã học ở **Bài 13 về Data Persistence (Volumes)**, dữ liệu database phải được lưu ở volume để không bị mất khi container bị xóa. Chúng ta cũng sẽ cấu hình để tự chạy file `init.sql` tạo bảng lúc khởi tạo.

**Tạo thư mục `db-init/` và file `init.sql`**:
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, email) VALUES ('admin', 'admin@example.com') ON CONFLICT DO NOTHING;
```

Trong file `docker-compose.yml` sắp tới, ta sẽ mount thư mục này vào `/docker-entrypoint-initdb.d/` của PostgreSQL container.

---

## Bước 5: Cấu hình Redis

Redis đóng vai trò làm cache. Ta sẽ giới hạn memory để tránh Redis chiếm hết tài nguyên của host (xem **Bài 19: Resource Limits**). Redis cũng cần được lưu trữ cấu hình nhỏ để chạy bền bỉ hơn, mặc dù dữ liệu cache có thể tạm thời. Ta sẽ khai báo nó trực tiếp trong Compose file.

---

## Bước 6: Viết docker-compose.yml đầy đủ

File `docker-compose.yml` là "trái tim" của toàn bộ hệ thống. Ở đây, chúng ta sẽ kết hợp mọi kiến thức về networks, volumes, environment variables và depends_on.

**File `docker-compose.yml`**:

```yaml
version: '3.8'

services:
  # 1. PostgreSQL Database
  db:
    image: postgres:15-alpine
    container_name: capstone_db
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-supersecret}
      POSTGRES_DB: ${DB_NAME:-capstone_db}
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./db-init:/docker-entrypoint-initdb.d
    networks:
      - backend-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-postgres} -d ${DB_NAME:-capstone_db}"]
      interval: 10s
      timeout: 5s
      retries: 5

  # 2. Redis Cache
  redis:
    image: redis:7-alpine
    container_name: capstone_redis
    restart: unless-stopped
    command: redis-server --requirepass ${REDIS_PASSWORD:-redispass}
    networks:
      - backend-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 3
    deploy:
      resources:
        limits:
          memory: 256M

  # 3. Backend API
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: capstone_backend
    restart: unless-stopped
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - REDIS_HOST=redis
      - REDIS_PASSWORD=${REDIS_PASSWORD}
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - backend-network
      - frontend-network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  # 4. Frontend Web
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: capstone_frontend
    restart: unless-stopped
    networks:
      - frontend-network
    depends_on:
      - backend

  # 5. Nginx Reverse Proxy
  reverse-proxy:
    build:
      context: ./reverse-proxy
      dockerfile: Dockerfile
    container_name: capstone_proxy
    ports:
      - "80:80"
    restart: unless-stopped
    depends_on:
      - frontend
      - backend
    networks:
      - frontend-network

# Định nghĩa Volumes
volumes:
  pgdata:

# Định nghĩa Networks (Tách biệt mạng frontend và backend để bảo mật)
networks:
  frontend-network:
    driver: bridge
  backend-network:
    driver: bridge
```

> ⚠️ **Lưu ý:** Chú ý kỹ thuật tách mạng (Network Segmentation). `reverse-proxy` và `frontend` nằm ở `frontend-network`. `db` và `redis` nằm ở `backend-network`. `backend` tham gia vào CẢ HAI mạng để làm cầu nối. Điều này giúp tăng cường bảo mật: từ mạng frontend không thể kết nối trực tiếp vào database.

---

## Bước 7: Cấu hình Nginx reverse proxy

Reverse proxy sẽ định tuyến traffic của người dùng tới đúng service.

**Tạo thư mục `reverse-proxy/` và file `nginx.conf`**:
```nginx
server {
    listen 80;
    server_name localhost;

    # Cấu hình log
    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;

    # Định tuyến tới Frontend
    location / {
        proxy_pass http://frontend:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Định tuyến tới Backend API
    location /api/ {
        proxy_pass http://backend:8000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**File `reverse-proxy/Dockerfile`**:
```dockerfile
FROM nginx:1.25-alpine
# Xóa cấu hình mặc định
RUN rm /etc/nginx/conf.d/default.conf
# Thêm cấu hình proxy mới
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

---

## Bước 8: Environment variables và .env files

Quản lý thông tin nhạy cảm là bài học quan trọng (Xem **Bài 14: Biến môi trường**).

**File `.env` ở thư mục gốc**:
```env
DB_USER=myadmin
DB_PASSWORD=SecurePassword123!
DB_NAME=production_db
REDIS_PASSWORD=CacheSecretKey99!
```

Docker Compose sẽ tự động tải các biến này khi chạy lệnh `docker-compose up`.

---

## Bước 9: Health checks cho mọi service

Trong `docker-compose.yml` ở trên, ta đã định nghĩa thuộc tính `healthcheck` cho PostgreSQL và Redis, đồng thời thêm phần `depends_on: ... condition: service_healthy` ở `backend`. Điều này đảm bảo `backend` sẽ CHỈ bắt đầu chạy khi Database và Cache đã sẵn sàng hoàn toàn, tránh lỗi crash app do không kết nối được database khi vừa khởi động (Xem **Bài 17: Healthchecks**).

---

## Bước 10: Tối ưu hóa image (Multi-stage, .dockerignore)

Nhìn lại quá trình, chúng ta đã thực hiện:
- Dùng base image Alpine hoặc Slim (`node:18-alpine`, `python:3.11-slim`) để giảm kích thước.
- Sử dụng `.dockerignore` cho cả Frontend và Backend để không mang rác vào image.
- Chia Stage (Multi-stage) để không lưu lại công cụ build (compiler, gcc) trên môi trường production.

---

## Bước 11: CI/CD Pipeline với GitHub Actions

Để tự động hóa (Xem **Bài 20: Tích hợp CI/CD**), chúng ta định nghĩa một quy trình tự động Build và Push image lên Docker Hub mỗi khi có code mới ở branch `main`.

**Tạo file `.github/workflows/docker-build.yml`**:
```yaml
name: CI/CD Pipeline Build & Push

on:
  push:
    branches:
      - main

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3

      - name: Log in to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and Push Backend
        uses: docker/build-push-action@v4
        with:
          context: ./backend
          push: true
          tags: ${{ secrets.DOCKER_USERNAME }}/capstone-backend:latest

      - name: Build and Push Frontend
        uses: docker/build-push-action@v4
        with:
          context: ./frontend
          push: true
          tags: ${{ secrets.DOCKER_USERNAME }}/capstone-frontend:latest
```
*(Hãy đảm bảo thêm `DOCKER_USERNAME` và `DOCKER_PASSWORD` vào phần Secrets của GitHub repository).*

---

## Bước 12: Production Best Practices

Cuối cùng, dự án của bạn đã bao gồm các Best Practices dành cho môi trường Production:
- **Restart Policy**: `restart: unless-stopped` (Dịch vụ sẽ tự phục hồi nếu bị crash hoặc sau khi server reboot).
- **Logging Limits**: Khai báo JSON file logging với `max-size: "10m"` giúp host không bị đầy ổ cứng do log rác sinh ra qua thời gian.
- **Resource Constraints**: Giới hạn memory cho Redis (`memory: 256M`).
- **Network Segmentation**: Các lớp mạng bị cô lập an toàn, Frontend không có quyền "thấy" Database.

Để chạy dự án, bạn chỉ cần gõ một lệnh:
```bash
docker-compose up -d --build
```

**Kết quả:**
Bạn sẽ thấy 5 container được tạo ra. Mở trình duyệt truy cập `http://localhost`, bạn sẽ thấy giao diện Frontend. Gọi API ở `http://localhost/api/users`, Backend sẽ kết nối Database và trả về dữ liệu.

---

## Bài tập

1. Thay đổi nội dung ứng dụng Frontend và rebuild lại image.
2. Kiểm tra các image vừa được tạo ra bằng `docker images` xem kích thước là bao nhiêu.
3. Triển khai thêm một service "pgadmin" (cho PostgreSQL) và một service "redis-commander" (cho Redis) vào file `docker-compose.yml` để quản lý database qua giao diện UI.
4. Đảm bảo pgadmin và redis-commander chỉ chạy trên `backend-network`.
5. Viết script tự động backup dữ liệu PostgreSQL định kỳ mỗi ngày bằng cách dùng cron chạy lệnh `docker exec pg_dump`.
6. Áp dụng reverse-proxy bảo mật bằng cách cấu hình file SSL (HTTPS) cho Nginx.

---

## Checklist hoàn thành

- [ ] Thiết lập đúng cấu trúc thư mục dự án
- [ ] Viết thành công `Dockerfile` Multi-stage cho Backend
- [ ] Chạy Backend bằng quyền non-root user
- [ ] Cấu hình `.dockerignore` đầy đủ
- [ ] Viết `Dockerfile` Multi-stage cho Frontend
- [ ] Tạo file cấu hình `nginx.conf` cho Reverse Proxy
- [ ] Sử dụng Docker Volume để bảo toàn dữ liệu Postgres
- [ ] Chạy tự động `init.sql` vào database
- [ ] Tạo file `.env` tách biệt cấu hình mật khẩu
- [ ] Phân chia `backend-network` và `frontend-network`
- [ ] Cấu hình `healthcheck` cho Database
- [ ] Cấu hình `depends_on` với `service_healthy` cho Backend
- [ ] Thêm limit resources cho Redis
- [ ] Giới hạn kích thước file log
- [ ] Tích hợp CI/CD tự động lên Docker Hub

---

## Mở rộng (Tùy chọn)

Nếu bạn muốn nâng cấp dự án để đạt chuẩn "Senior", hãy thử thách thêm với các mục sau:
- **Monitoring Stack**: Thêm Prometheus và Grafana vào `docker-compose.yml` để giám sát CPU, RAM của toàn bộ hệ thống.
- **HTTPS & SSL/TLS**: Tích hợp Certbot (Let's Encrypt) vào Nginx Proxy để chạy ứng dụng trên HTTPS (Port 443).
- **Kubernetes Deployment**: Chuyển đổi file `docker-compose.yml` thành các file YAML của Kubernetes (Deployments, Services, Ingress, ConfigMaps, Secrets) và triển khai lên một cụm K3s/Minikube.
- **Blue-Green Deployment**: Cấu hình quy trình triển khai zero-downtime, cập nhật version backend mới mà không làm rớt request nào.

🎉 **Chúc mừng bạn đã hoàn thành khóa học Docker từ A-Z!** 🎉
Bạn đã trang bị cho mình kỹ năng cốt lõi của một DevOps Engineer/Backend Developer hiện đại. Bây giờ, bạn có thể tự tin triển khai bất kỳ dự án nào lên mọi môi trường với Docker.

---

## Tiếp theo
Khóa học tiếp theo: **Làm chủ Kubernetes (K8s) cho người mới bắt đầu** → Đưa các ứng dụng Docker lên một tầm cao mới với khả năng Auto-scaling (tự động mở rộng quy mô) và Self-healing (tự động phục hồi lỗi). Hãy sẵn sàng!
