# Bài 12: Docker Compose nâng cao

## Mục tiêu
- Hiểu và sử dụng được các cấu hình nâng cao trong Docker Compose để quản lý dependency (sự phụ thuộc) giữa các service.
- Biết cách cấu hình `healthcheck` để đảm bảo trạng thái hoạt động thực sự của ứng dụng thay vì chỉ trạng thái start của container.
- Sử dụng hiệu quả `profiles` để chạy các nhóm service riêng biệt theo từng ngữ cảnh.
- Tái sử dụng cấu hình YAML nhằm giảm trùng lặp mã bằng thuộc tính `extends` và cơ chế anchor/alias (`&`, `*`).
- Quản lý linh hoạt cấu hình đa môi trường (dev, staging, prod) với Override files và lệnh `docker compose -f`.
- Nắm vững cách quản lý biến môi trường, scale service và thiết lập các `restart` policy.

## 1. Quản lý thứ tự khởi động với `depends_on` và `healthcheck`

### 1.1 Hạn chế của `depends_on` cơ bản
Khi ứng dụng của bạn có nhiều thành phần, ví dụ backend API phụ thuộc vào database PostgreSQL, bạn không thể để API gọi vào database khi nó chưa sẵn sàng. Mặc định, Docker Compose khởi động mọi thứ đồng thời. Bạn dùng `depends_on` để quy định thứ tự.

Tuy nhiên, `depends_on` mặc định **chỉ đảm bảo thứ tự start container**. Tức là Compose gửi lệnh start tới database, rồi gửi ngay lệnh start tới API. Nhưng PostgreSQL cần vài giây để khởi tạo bộ nhớ và sẵn sàng nhận kết nối, dẫn đến API vẫn có thể bị crash vì kết nối thất bại ngay từ giây đầu tiên.

### 1.2 Kết hợp `healthcheck` và `condition: service_healthy`
Để giải quyết bài toán trên, Docker cung cấp `healthcheck`. Bạn định nghĩa cách kiểm tra ứng dụng (ví dụ: dùng lệnh `pg_isready`), và `depends_on` sẽ đợi cho đến khi bài kiểm tra sức khỏe báo cáo thành công (`service_healthy`).

> 💡 **Tip:** Tương tự như việc bạn chờ bác sĩ khám và cấp giấy chứng nhận sức khỏe trước khi được phép tham gia một giải đấu thể thao.

```yaml
# docker-compose.yml
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secretpassword
      POSTGRES_DB: mydb
    # Khai báo cách kiểm tra sức khỏe của container
    healthcheck:
      # test là lệnh chạy bên trong container để kiểm tra trạng thái
      test: ["CMD-SHELL", "pg_isready -U admin -d mydb"]
      interval: 10s    # Tần suất kiểm tra (cứ 10 giây một lần)
      timeout: 5s      # Thời gian tối đa chờ phản hồi của mỗi lần test
      retries: 5       # Thất bại liên tiếp 5 lần mới đánh dấu là 'unhealthy'
      start_period: 5s # Thời gian ân hạn ban đầu để app kịp khởi động

  api:
    image: my-backend-api:v1.0
    ports:
      - "8080:8080"
    depends_on:
      db:
        # Chờ đến khi service 'db' báo cáo trạng thái 'healthy' mới bắt đầu start API
        condition: service_healthy
```

**Kết quả:**
```
[+] Running 2/2
 ✔ Container db   Healthy                                               11.5s
 ✔ Container api  Started                                                0.5s
```
Container API sẽ chỉ được khởi động an toàn khi PostgreSQL đã hoàn toàn sẵn sàng nhận kết nối ở cổng 5432.

## 2. Profiles — Phân nhóm service theo ngữ cảnh

Trong quá trình phát triển, bạn có thể có rất nhiều service trong cùng một file `docker-compose.yml`: Frontend, Backend, Database, Redis, Elasticsearch, Kafka, và các công cụ Test. Nhưng đôi khi bạn chỉ cần chạy Backend và Database, không muốn tốn tài nguyên chạy Kafka hay UI.

Docker Compose giải quyết việc này thông qua **`profiles`**.

```yaml
# docker-compose.yml
services:
  db:
    image: postgres:15
    # Service không có khai báo profiles sẽ luôn chạy mặc định

  api:
    image: my-api
    depends_on:
      - db
    # Không có profiles -> luôn chạy

  frontend:
    image: my-frontend
    profiles:
      - "ui" # Chỉ khởi động khi profile 'ui' được kích hoạt

  admin-dashboard:
    image: admin-panel
    profiles:
      - "ui"
      - "admin" # Kích hoạt bằng 'ui' hoặc 'admin' đều được

  db-seeder:
    image: seeder-script
    profiles:
      - "testing"
```

**Cách chạy thực tế:**
```bash
# Lệnh cơ bản chỉ chạy db và api (frontend, admin-dashboard, db-seeder sẽ bị bỏ qua)
docker compose up -d

# Chạy backend và toàn bộ UI (db, api, frontend, admin-dashboard)
docker compose --profile ui up -d

# Chạy nhiều profile cùng lúc
docker compose --profile ui --profile testing up -d
```

## 3. Tái sử dụng cấu hình YAML

Khi các service của bạn có cấu hình giống hệt nhau (như giới hạn tài nguyên, biến môi trường dùng chung, hay driver logging), việc copy-paste mã sẽ làm file Compose dài và khó bảo trì. Docker Compose hỗ trợ hai kỹ thuật quan trọng để DRY (Don't Repeat Yourself).

### 3.1 YAML Anchor (`&`) và Alias (`*`)
Đây là tính năng gốc của định dạng file YAML. Anchor đánh dấu một phần nội dung và Alias chèn phần nội dung đó vào nơi khác.

> ⚠️ **Lưu ý:** Thường chúng ta định nghĩa các anchor trong một thuộc tính mở rộng (extension field) bắt đầu bằng `x-`, Docker Compose sẽ bỏ qua các block này khi thực thi.

```yaml
# Định nghĩa các block chung
x-logging-config:
  &default-logging
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "5"

x-common-env:
  &common-env
  TIMEZONE: "Asia/Ho_Chi_Minh"
  DEBUG_MODE: "false"

services:
  web:
    image: nginx
    logging: *default-logging # Gắn toàn bộ nội dung của anchor vào đây

  worker:
    image: my-worker
    environment: *common-env  # Tái sử dụng môi trường

  api:
    image: my-api
    # Bổ sung <<: *alias để 'kế thừa' và mở rộng hoặc ghi đè thuộc tính
    logging:
      <<: *default-logging
      options:
        max-size: "50m" # Ghi đè max-size cho riêng api
```

### 3.2 Thuộc tính `extends`
Cho phép một service kế thừa cấu hình từ một service khác. Service gốc có thể nằm trong cùng một file hoặc một file hoàn toàn khác.

```yaml
# base.yml - Định nghĩa template gốc
services:
  app-base:
    image: node:18-alpine
    restart: always
    environment:
      - NODE_ENV=production
```

```yaml
# docker-compose.yml - Sử dụng lại cấu hình
services:
  auth-service:
    extends:
      file: base.yml
      service: app-base
    ports:
      - "3000:3000"
    command: ["npm", "run", "start:auth"]

  billing-service:
    extends:
      file: base.yml
      service: app-base
    ports:
      - "3001:3000"
    command: ["npm", "run", "start:billing"]
```
Cả `auth-service` và `billing-service` đều sẽ có `image`, `restart` và `NODE_ENV` giống hệt nhau.

## 4. Quản lý Override files và lệnh `docker compose -f`

Một nguyên tắc vàng trong Docker Compose là: **Không sửa file gốc để phục vụ nhu cầu cá nhân của môi trường**. Hãy dùng file ghi đè (Override file).

### 4.1 `docker-compose.override.yml` mặc định
Theo mặc định, nếu trong thư mục có file `docker-compose.override.yml`, Docker Compose sẽ tự động kết hợp (merge) cấu hình từ nó vào file gốc `docker-compose.yml`.

- **docker-compose.yml**: Chứa cấu hình tối giản, chuẩn cho Production/CI.
- **docker-compose.override.yml**: Chứa cấu hình dành cho local dev (ví dụ: mount source code qua volumes để live reload, mở cổng port ra máy host).

> 📌 **Ghi nhớ:** Override file chỉ hoạt động tự động khi bạn gõ lệnh không kèm cờ `-f`. File này thường được `.gitignore` để mỗi lập trình viên có thể cấu hình local riêng.

### 4.2 Ghi đè bằng `docker compose -f`
Để quản lý cụ thể các môi trường như dev, staging, prod, bạn dùng chuỗi tham số `-f`. File sau sẽ ghi đè file trước.

```bash
# Triển khai môi trường staging
docker compose -f docker-compose.yml -f docker-compose.staging.yml up -d

# Triển khai môi trường production (file prod sẽ override file gốc)
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## 5. Scaling, Restart Policies và Environment

### 5.1 Scaling (Mở rộng quy mô)
Bạn có thể tăng số lượng container cho một service dễ dàng bằng flag `--scale`. 

```bash
docker compose up -d --scale web=3
```
**Lưu ý:** Để scale được, service `web` không được định nghĩa `ports: "80:80"` cố định, vì 3 container không thể dùng chung port 80 của máy host. Bạn phải để port máy host ngẫu nhiên: `ports: - "80"` (chỉ map port đích).

### 5.2 Restart Policy
Xác định hành vi của Docker Daemon khi container bị tắt:
- `no` (Mặc định): Không tự động khởi động lại, kể cả khi lỗi.
- `always`: Luôn khởi động lại bằng mọi giá. Thường dùng cho Production.
- `on-failure`: Chỉ khởi động lại nếu container kết thúc với mã lỗi (exit code != 0).
- `unless-stopped`: Giống `always`, ngoại trừ việc nếu bạn chủ động dùng lệnh `docker stop`, container sẽ không tự bật lại khi Docker Daemon khởi động lại máy chủ.

```yaml
services:
  web:
    image: nginx
    restart: unless-stopped
```

### 5.3 Biến môi trường (Environment variables)
Cách Docker Compose nạp biến môi trường theo thứ tự ưu tiên (thấp đến cao):
1. Khai báo `.env` file (Compose tự động load các biến trong file `.env` đặt cùng cấp thư mục).
2. Thuộc tính `env_file` trong YAML.
3. Thuộc tính `environment` trong YAML.
4. Biến được set trực tiếp ở shell terminal lúc chạy lệnh.

```yaml
services:
  app:
    image: my-app
    # Load danh sách biến từ file cụ thể
    env_file:
      - .env.staging
    # Trực tiếp ghi đè biến
    environment:
      - API_KEY=${SECRET_KEY:-default_key} # Sử dụng nội suy biến (interpolation)
```

## Bài tập
1. **Healthcheck Dependency**: Tạo một file `docker-compose.yml` chạy một container Redis và một container Ubuntu (app) dùng lệnh `sleep 3600`. Cấu hình healthcheck cho Redis bằng lệnh `redis-cli ping`. Cấu hình để container app chỉ được start nếu Redis pass healthcheck.
2. **Sử dụng Anchor YAML**: Tạo cấu hình file Compose cho 3 service (backend1, backend2, backend3). Cả 3 service này đều dùng image `nginx:alpine` và có chung cài đặt về `restart: always` và `logging` (kiểu json-file, max-size 5m). Dùng kĩ thuật YAML anchor để tránh lặp code.
3. **Profiles**: Viết file compose có service `web-server` chạy bình thường và một service `admin-panel` (image nginx) chỉ được định nghĩa chạy khi kích hoạt profile tên là `admin`. Chạy lệnh test xem admin-panel có hoạt động đúng theo profile chưa.
4. **Override Files**: Định nghĩa một Nginx service chạy ở cổng 8080 trong file cơ bản `docker-compose.yml`. Tạo thêm file `docker-compose.override.yml` để map một thư mục local chứa file `index.html` của bạn vào đường dẫn `/usr/share/nginx/html` trong container. Chạy `docker compose up -d` (không dùng cờ -f) và kiểm tra xem nội dung ghi đè có tự động được nạp hay không.

## Tiếp theo
→ [Bài 13: Tối ưu hóa Image](./13_image_optimization.md)
