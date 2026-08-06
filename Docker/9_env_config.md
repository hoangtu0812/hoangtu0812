# Bài 9: Environment Variables & Cấu hình

## Mục tiêu
- Hiểu được tại sao cần sử dụng biến môi trường (Environment Variables) trong container.
- Nắm vững các cách truyền biến môi trường vào container khi chạy.
- Phân biệt rõ ràng giữa `ARG` và `ENV` trong Dockerfile.
- Biết cách quản lý file `.env` và các best practices về bảo mật.
- Thực hành cấu hình một ứng dụng thực tế thông qua biến môi trường.

---

## 1. Tại sao cần Environment Variables?

### Khái niệm cơ bản
Trong phát triển phần mềm, một ứng dụng thường phải chạy trên nhiều môi trường khác nhau như: Development (Phát triển), Staging (Kiểm thử), và Production (Sản xuất). Mỗi môi trường sẽ có các cấu hình khác nhau như thông tin kết nối Database, API Keys, hay cổng (port) hoạt động.

Thay vì phải sửa code hoặc tạo ra nhiều image khác nhau cho mỗi môi trường, bạn chỉ cần một Docker image duy nhất và dùng **Environment Variables (Biến môi trường)** để thay đổi cấu hình linh hoạt.

> 💡 **Tip:** Hãy tưởng tượng biến môi trường giống như các "công tắc" bên ngoài giúp bạn thay đổi hành vi của cỗ máy (container) mà không cần phải mở tung nó ra.

---

## 2. Các cách truyền biến môi trường vào Container

### 2.1 Sử dụng cờ `-e` trong lệnh `docker run`
Đây là cách nhanh nhất để truyền một hoặc vài biến môi trường vào container khi khởi chạy.

```bash
# Khởi chạy một container Nginx và truyền biến môi trường WELCOME_MSG
docker run -d --name web_server -e WELCOME_MSG="Xin chào từ Docker!" nginx:alpine
```

**Kết quả:**
```
Container web_server được tạo với biến môi trường WELCOME_MSG="Xin chào từ Docker!"
```

### 2.2 Sử dụng file cấu hình `--env-file`
Khi bạn có quá nhiều biến, việc dùng `-e` sẽ làm câu lệnh trở nên rất dài và khó quản lý. Lúc này, file `.env` là giải pháp hoàn hảo.

Cú pháp file `.env`:
```env
# file: app.env
DB_HOST=192.168.1.10
DB_PORT=5432
DB_USER=admin
DB_PASS=supersecret
```

Chạy container với file env:
```bash
docker run -d --name app_server --env-file ./app.env alpine env
```

> ⚠️ **Lưu ý:** Tuyệt đối không commit các file `.env` chứa mật khẩu hoặc API key thực tế lên Git. Hãy đưa nó vào `.gitignore` và chỉ cung cấp file `.env.example` làm mẫu.

### 2.3 Lệnh `ENV` trong Dockerfile
Bạn có thể thiết lập biến môi trường mặc định ngay trong quá trình build image bằng chỉ thị `ENV`.

```dockerfile
# Sử dụng base image
FROM alpine:latest

# Thiết lập biến môi trường mặc định
ENV APP_PORT=8080
ENV APP_ENV=production

CMD echo "Ứng dụng chạy ở môi trường $APP_ENV trên cổng $APP_PORT"
```

> 📌 **Ghi nhớ - Độ ưu tiên:** Giá trị truyền vào lúc chạy (`docker run -e` hoặc `--env-file`) sẽ **ghi đè** giá trị mặc định được định nghĩa bằng `ENV` trong Dockerfile.

---

## 3. Phân biệt ARG và ENV

Một trong những nhầm lẫn phổ biến nhất của người mới học Docker là sự khác biệt giữa `ARG` và `ENV`.

- **`ARG` (Build-time variable):** Chỉ tồn tại trong quá trình bạn chạy lệnh `docker build`. Khi image đã build xong, các giá trị này sẽ biến mất. Thường dùng để truyền phiên bản gói phần mềm cần cài đặt.
- **`ENV` (Runtime variable):** Tồn tại cả trong lúc build và khi container đang chạy (`docker run`). Thường dùng cho cấu hình ứng dụng.

### Kết hợp ARG và ENV
Bạn có thể dùng `ARG` để nhận giá trị khi build, sau đó gán nó cho `ENV` để container có thể sử dụng khi chạy.

```dockerfile
FROM alpine:latest

# Khai báo ARG (nhận từ docker build --build-arg VERSION=1.0)
ARG VERSION=latest

# Gán ARG cho ENV để ứng dụng có thể đọc khi chạy
ENV APP_VERSION=$VERSION

CMD echo "Phiên bản ứng dụng: $APP_VERSION"
```

---

## 4. Xem biến môi trường của Container đang chạy

Để kiểm tra xem một container đang có những biến môi trường nào, bạn sử dụng lệnh `docker inspect`.

```bash
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' <container_id_hoac_name>
```

**Kết quả:**
```
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
APP_PORT=8080
APP_ENV=production
```

---

## Bài tập

1. Tạo một container từ image `alpine` chạy lệnh `env`. Truyền vào 2 biến môi trường `MY_NAME` và `MY_AGE` qua cờ `-e` và quan sát kết quả in ra màn hình.
2. Tạo một file tên là `config.env` chứa 3 cấu hình database (HOST, USER, PASS). Chạy một container `ubuntu` và mount file này vào bằng cờ `--env-file`, sau đó chạy lệnh `printenv` bên trong container để kiểm tra.
3. Viết một Dockerfile định nghĩa biến `ENV` mặc định là `TIMEZONE=Asia/Ho_Chi_Minh`. Build image và chạy container để kiểm tra. Sau đó chạy một container mới từ image đó nhưng ghi đè `TIMEZONE=Asia/Tokyo` bằng cờ `-e`.
4. Viết Dockerfile sử dụng `ARG` để nhận `APP_VERSION` khi build. Gán giá trị đó cho một biến `ENV` có cùng tên và in ra màn hình bằng `CMD`. Hãy build image với các tham số `--build-arg` khác nhau và chạy container để kiểm chứng `ENV` đã nhận đúng giá trị của `ARG`.

---

## Tiếp theo
→ [Bài 10: Dockerfile Nâng Cao](./10_dockerfile_advanced.md)
