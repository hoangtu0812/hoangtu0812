# Bài 10: Dockerfile nâng cao

## Mục tiêu
- Hiểu và áp dụng được kỹ thuật Multi-stage build để tối ưu dung lượng image.
- Phân biệt và sử dụng đúng ENTRYPOINT và CMD.
- Truyền tham số lúc build với ARG và kết hợp với ENV.
- Cấu hình kiểm tra sức khỏe container bằng HEALTHCHECK.
- Tăng cường bảo mật với USER và quản lý dữ liệu với VOLUME.
- Nắm vững các best practices (thực hành tốt nhất) khi viết Dockerfile.

---

## 1. Multi-stage build (Build nhiều giai đoạn)

### Tại sao cần Multi-stage build?
Khi build ứng dụng (đặc biệt là Go, Java, React, C++), bạn cần rất nhiều công cụ, thư viện (SDK, compiler, node_modules) làm cho kích thước image phình to (có thể lên tới hàng GB). Nhưng khi chạy thực tế, ứng dụng chỉ cần file thực thi (binary) hoặc file tĩnh (HTML/CSS/JS).
Multi-stage build cho phép bạn dùng nhiều lệnh `FROM` trong một Dockerfile. Mỗi `FROM` bắt đầu một giai đoạn (stage) mới. Bạn có thể copy kết quả (artifact) từ stage trước sang stage sau, và bỏ lại tất cả những công cụ build không cần thiết.

### Ví dụ 1: Build ứng dụng Golang
```dockerfile
# Stage 1: Build (Môi trường có đầy đủ công cụ compile)
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
# Biên dịch mã nguồn thành file chạy 'myapp'
RUN go build -o myapp main.go

# Stage 2: Production (Môi trường siêu nhẹ)
FROM alpine:latest
WORKDIR /app
# Chỉ copy file thực thi từ stage 'builder'
COPY --from=builder /app/myapp .
# Lệnh chạy ứng dụng
CMD ["./myapp"]
```

> 💡 **Tip:** Kích thước image cuối cùng sẽ chỉ khoảng vài MB (alpine + myapp) thay vì ~300MB+ của image `golang`.

### Ví dụ 2: Build ứng dụng React (Node.js + Nginx)
```dockerfile
# Stage 1: Cài đặt node_modules và build code React
FROM node:18-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Phục vụ file tĩnh bằng Nginx
FROM nginx:alpine
# Copy thư mục build từ stage trước vào thư mục public của Nginx
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**Kết quả:**
Image được build ra chỉ chứa server Nginx và các file HTML/JS/CSS tĩnh, không hề có NodeJS hay thư mục node_modules nặng nề.

---

## 2. ENTRYPOINT và CMD

Nhiều người thường nhầm lẫn giữa `ENTRYPOINT` và `CMD`. Dù cả hai đều xác định lệnh chạy khi container khởi động, nhưng chúng có vai trò khác nhau.

### Shell form vs Exec form
Docker hỗ trợ hai cách viết:
- **Shell form:** `CMD echo "Hello"` (Chạy dưới quyền sub-process của shell: `/bin/sh -c`)
- **Exec form (Khuyên dùng):** `CMD ["echo", "Hello"]` (Chạy trực tiếp, không qua shell, nhận tín hiệu tắt tốt hơn)

### Sự khác biệt
- **ENTRYPOINT:** Lệnh **cố định**, rất khó bị ghi đè (override) khi `docker run`. Thường dùng để biến container thành một công cụ dòng lệnh (CLI).
- **CMD:** Tham số **mặc định**, rất dễ bị ghi đè khi bạn truyền thêm tham số vào cuối lệnh `docker run`.

### Kết hợp ENTRYPOINT + CMD
Khi dùng chung, `ENTRYPOINT` đóng vai trò là lệnh chính, còn `CMD` đóng vai trò là tham số mặc định cho `ENTRYPOINT`.

```dockerfile
FROM alpine:latest
ENTRYPOINT ["ping", "-c", "3"]
CMD ["google.com"]
```

**Kết quả:**
```bash
# Sẽ chạy: ping -c 3 google.com
docker run my-ping-image 

# Sẽ chạy: ping -c 3 facebook.com (Ghi đè CMD bằng facebook.com)
docker run my-ping-image facebook.com
```

---

## 3. Lệnh ARG (Build-time variables)

`ARG` định nghĩa biến chỉ tồn tại trong lúc **build image** (không khả dụng khi container đang chạy). Rất hữu ích để truyền phiên bản động.

```dockerfile
FROM alpine:latest
# Khai báo ARG có giá trị mặc định
ARG APP_VERSION=1.0.0
# Có thể dùng để tạo thư mục, tải đúng phiên bản
RUN echo "Đang build phiên bản $APP_VERSION" > version.txt

# Để biến này tồn tại lúc runtime, kết hợp gán vào ENV
ENV VERSION=$APP_VERSION
```

**Kết quả (Truyền tham số lúc build):**
```bash
docker build --build-arg APP_VERSION=2.5.0 -t myapp:2.5.0 .
```

---

## 4. Kiểm tra sức khỏe với HEALTHCHECK

`HEALTHCHECK` giúp Docker biết ứng dụng bên trong container thực sự có đang hoạt động tốt (Healthy) hay không, thay vì chỉ biết process có đang chạy hay không.

```dockerfile
FROM nginx:alpine
COPY ./html /usr/share/nginx/html

# Cứ mỗi 30 giây, chạy lệnh curl kiểm tra localhost. Timeout sau 3 giây.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
```

> 📌 **Ghi nhớ:** Khi chạy lệnh `docker ps`, trạng thái container sẽ có thêm chữ `(healthy)` hoặc `(unhealthy)`.

---

## 5. Lệnh USER và VOLUME

### USER: Tránh chạy bằng root
Theo mặc định, container chạy bằng tài khoản `root`. Điều này tiềm ẩn rủi ro bảo mật nếu hacker chiếm quyền được container.

```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY . .
RUN npm install
# Đổi sang user 'node' (đã có sẵn trong image node) trước khi chạy app
USER node
CMD ["npm", "start"]
```

### VOLUME: Khai báo điểm gắn dữ liệu
`VOLUME` khai báo rằng thư mục này sẽ dùng để chứa dữ liệu bền vững (persistence).

```dockerfile
FROM postgres:13
# Khai báo sẵn điểm kết nối
VOLUME ["/var/lib/postgresql/data"]
```

> ⚠️ **Lưu ý:** `VOLUME` trong Dockerfile không mount trực tiếp vào một thư mục cụ thể trên máy host. Nó chỉ tạo Anonymous Volume nếu khi chạy bạn không khai báo tham số `-v`.

---

## 6. Best Practices (Thực hành tốt nhất)

1. **Thứ tự layer (Layer caching):** Docker cache từng dòng lệnh. Nếu dòng nào thay đổi, mọi dòng dưới nó sẽ bị chạy lại. **Hãy đặt những lệnh ít thay đổi (như cài dependencies) lên trước, lệnh copy source code xuống sau**.
2. **Giảm số lượng layer:** Kết hợp nhiều lệnh `RUN` bằng `&&` và `\`.
   ```dockerfile
   RUN apt-get update && \
       apt-get install -y curl && \
       rm -rf /var/lib/apt/lists/*
   ```
3. **Dọn dẹp cache:** Luôn xóa cache của các trình quản lý gói (như `apt`, `npm`) trong cùng một lệnh `RUN` để giảm dung lượng file hệ thống.
4. **Luôn dùng .dockerignore:** Loại bỏ `.git`, `node_modules` nội bộ ra khỏi context trước khi `COPY . .`.

---

## Bài tập

1. Viết một Dockerfile sử dụng `ENTRYPOINT` và `CMD`. Khai báo sao cho khi chạy `docker run my-greeting`, màn hình in ra "Hello World". Khi chạy `docker run my-greeting Docker`, màn hình in ra "Hello Docker".

2. Viết một Dockerfile Multi-stage cho dự án Node.js:
- Stage 1 (builder): Cài đặt tất cả các gói kể cả `devDependencies`.
- Stage 2 (production): Chỉ copy mã nguồn, cấu hình cài đặt `production dependencies`, và thiết lập `USER node`.

3. Hãy tối ưu một Dockerfile bạn từng viết bằng cách:
   1. Thêm tham số phiên bản qua `ARG`.
   2. Tạo lệnh kiểm tra tình trạng ứng dụng với `HEALTHCHECK`.
   3. Gom nhiều lệnh `RUN` thành một lớp (layer) duy nhất và làm sạch bộ nhớ tạm.

4. Tự viết một Dockerfile để biên dịch một file C++ "Hello World" bằng `gcc` ở Stage 1. Tại Stage 2, copy file thực thi (binary) đó vào image tên là `scratch` (một image siêu nhẹ của Docker, 0MB dung lượng). Khởi chạy và kiểm tra kết quả!

---

## Tiếp theo
→ [Bài 11: Docker Compose cơ bản](./11_compose_basics.md)
