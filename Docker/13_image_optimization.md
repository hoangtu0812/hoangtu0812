# Bài 13: Tối ưu hóa Image

## Mục tiêu
- Hiểu tầm quan trọng của việc tối ưu kích thước Docker Image.
- Biết cách lựa chọn base image phù hợp (alpine, slim, scratch).
- Nắm vững các kỹ thuật tận dụng cache và giảm số lượng layer.
- Sử dụng `.dockerignore` và Multi-stage build hiệu quả.
- Biết cách sử dụng các công cụ phân tích image (`docker history`, `dive`).

## 1. Tại sao kích thước Image lại quan trọng?

Tối ưu hóa Docker image không chỉ là vấn đề thẩm mỹ, nó ảnh hưởng trực tiếp đến hiệu suất và bảo mật của hệ thống:
- **Tốc độ Pull/Push**: Image nhỏ hơn giúp giảm thời gian tải lên (push) registry và kéo về (pull) máy chủ, đặc biệt quan trọng trong quy trình CI/CD.
- **Chi phí lưu trữ**: Image nhỏ giúp tiết kiệm dung lượng đĩa trên server và registry.
- **Bảo mật (Attack Surface)**: Image chứa ít thư viện và công cụ hơn (như shell, curl, wget) sẽ giảm thiểu nguy cơ bị tấn công. Ít code hơn = Ít lỗ hổng hơn.
- **Thời gian khởi động**: Các image nhỏ gọn có thể khởi động nhanh hơn một chút, hữu ích trong các kiến trúc Serverless hoặc Auto-scaling.

## 2. Chọn Base Image phù hợp

Lựa chọn base image là bước đầu tiên và quan trọng nhất. Cùng một ứng dụng, nhưng chọn sai base image có thể khiến kích thước tăng gấp 10 lần.

- **Standard (bookworm, ubuntu)**: Đầy đủ công cụ, kích thước lớn (hàng trăm MB). Chỉ nên dùng khi đang phát triển hoặc cần nhiều dependency phức tạp.
- **Slim (`python:3.9-slim`, `node:18-slim`)**: Đã loại bỏ nhiều công cụ không cần thiết nhưng vẫn giữ lại libc chuẩn (Debian). Kích thước trung bình (khoảng 100-150MB). Phù hợp cho hầu hết các ứng dụng.
- **Alpine (`python:3.9-alpine`, `node:18-alpine`)**: Rất nhỏ gọn (vài MB) vì sử dụng `musl libc` thay vì `glibc`. Kích thước cực nhỏ. Thích hợp nếu bạn biết chắc ứng dụng không gặp vấn đề tương thích với `musl`.
- **Scratch (`scratch`)**: Một image hoàn toàn trống rỗng (0 bytes). Rất tuyệt vời cho các ứng dụng đã được biên dịch tĩnh (static binary) như Go hoặc Rust.

> 💡 **Tip:** Nếu bạn không chắc chắn, hãy bắt đầu với `slim`. `alpine` đôi khi gây lỗi biên dịch với các thư viện Python (như numpy, pandas) do thiếu C compiler và glibc.

## 3. Tối ưu Cache Layer và Thứ tự Instruction

Docker build các image theo từng lớp (layer) tương ứng với mỗi lệnh trong Dockerfile. Nếu một layer không thay đổi, Docker sẽ dùng lại cache của nó.

### Quy tắc "Ít thay đổi trước, nhiều thay đổi sau"
Hãy đặt những lệnh ít bị thay đổi lên đầu Dockerfile.

**Sai lầm phổ biến:**
```dockerfile
COPY . /app
RUN npm install
```
Nếu bạn đổi chỉ 1 dòng code trong file `.js`, lệnh `COPY .` sẽ làm mất cache. Do đó, `npm install` sẽ phải chạy lại từ đầu, rất tốn thời gian.

**Cách tối ưu:** Tách riêng việc copy file danh sách dependency.
```dockerfile
# Copy file requirements trước
COPY package.json package-lock.json /app/
RUN npm install

# Copy source code sau
COPY . /app/
```

## 4. Giảm số lượng Layer và Dọn dẹp rác

Mỗi lệnh `RUN`, `COPY`, `ADD` đều tạo ra một layer mới. Lệnh `RUN` thứ hai không thể xóa dữ liệu đã được lưu ở layer trước đó.

### Gom các lệnh `RUN` bằng `&&` và xóa cache ngay lập tức
```dockerfile
# TỒI TỆ: Tạo ra 3 layer, layer cuối xóa file nhưng kích thước tổng không giảm!
RUN apt-get update
RUN apt-get install -y curl
RUN rm -rf /var/lib/apt/lists/*

# TỐT: Gom thành 1 layer, xóa rác trước khi kết thúc layer
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl && \
    rm -rf /var/lib/apt/lists/*
```

## 5. Nhắc lại: Multi-stage Build và `.dockerignore`

- **.dockerignore**: Tương tự `.gitignore`, nó ngăn các file không cần thiết (như `.git`, `node_modules`, `tests`) được gửi vào Docker daemon và xuất hiện trong image.
- **Multi-stage Build**: Dùng nhiều block `FROM` trong 1 Dockerfile. Chỉ copy những file thực thi đã biên dịch từ stage "build" sang stage "production", bỏ lại toàn bộ source code và các công cụ biên dịch. (Xem lại Bài 10).

## 6. Phân tích Image với `docker history` và `dive`

### Kiểm tra kích thước
```bash
docker image inspect my-app:latest --format='{{.Size}}'
```

### Xem lịch sử các layer
```bash
docker history my-app:latest
```
Lệnh này sẽ liệt kê từng layer, dung lượng của chúng và câu lệnh đã tạo ra layer đó.

### Sử dụng công cụ `dive`
`dive` là một công cụ mã nguồn mở tuyệt vời giúp bạn xem chi tiết nội dung của từng layer trong image.
```bash
# Cài đặt và chạy dive (ví dụ qua docker)
docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock wagoodman/dive:latest my-app:latest
```

## 7. Ví dụ: Từ 800MB xuống 50MB

**Trước khi tối ưu (800MB):**
```dockerfile
FROM golang:1.20
WORKDIR /app
COPY . .
RUN go build -o myapp main.go
CMD ["./myapp"]
```

**Sau khi tối ưu bằng Multi-stage build và Scratch (10MB):**
```dockerfile
# Stage 1: Build
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
# Biên dịch tĩnh (static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp main.go

# Stage 2: Production
FROM scratch
COPY --from=builder /app/myapp /myapp
ENTRYPOINT ["/myapp"]
```

## Bài tập

1. **Phân tích Image**: Chạy lệnh `docker history nginx:alpine` và so sánh kết quả với `docker history nginx:latest`. Ghi nhận sự khác biệt về số lượng layer và kích thước tổng thể.
2. **Tối ưu Dockerfile Python**: Viết một Dockerfile cho ứng dụng Python. Đảm bảo copy `requirements.txt` và chạy `pip install` trước khi copy toàn bộ source code. Sử dụng base image `python:3.9-slim`. Xóa cache pip bằng `--no-cache-dir`.
3. **Gom lệnh RUN**: Sửa lại đoạn Dockerfile sau để tối ưu layer:
   ```dockerfile
   RUN apt-get update
   RUN apt-get install -y wget
   RUN apt-get install -y unzip
   RUN rm -rf /var/lib/apt/lists/*
   ```
4. **Cài đặt và thử nghiệm `dive`**: Sử dụng công cụ `dive` để kiểm tra một image bất kỳ (ví dụ `alpine`) trên máy bạn và xem cấu trúc file hệ thống trong các layer.

## Tiếp theo
→ [Bài 14: Bảo mật Container](./14_security.md)
