# Bài 7: Volume & Bind Mount

## Mục tiêu
- Hiểu được vấn đề lưu trữ dữ liệu bền vững (persistent data) trong môi trường container.
- Nắm rõ 3 loại cơ chế lưu trữ (storage) trong Docker: Volume, Bind mount, và tmpfs mount.
- Thành thạo các câu lệnh tạo, quản lý và dọn dẹp Docker Volume.
- Phân biệt được sự khác biệt giữa Volume và Bind mount, biết khi nào nên dùng loại nào.
- Ứng dụng thực tế: Thiết lập lưu trữ an toàn cho cơ sở dữ liệu và thiết lập môi trường lập trình (live-reload).

---

## 1. Vấn đề: Dữ liệu mất khi container bị xóa

Theo mặc định, tất cả các tệp tin được tạo ra bên trong một container đang chạy sẽ được lưu trữ trên một **lớp ghi (writable layer)** tạm thời của container đó. Cơ chế này dẫn đến một số vấn đề nghiêm trọng khi triển khai thực tế:

1. **Tính tạm thời:** Khi container dừng lại hoặc bị xóa đi do lỗi hoặc cập nhật, lớp ghi này sẽ hoàn toàn bị xóa sạch cùng với container. Bạn sẽ mất toàn bộ dữ liệu quan trọng như hồ sơ khách hàng, đơn hàng, hay nhật ký hệ thống (logs).
2. **Khó truy cập:** Rất khó để lấy dữ liệu ra khỏi container đang chạy để sao lưu (backup) hoặc phân tích.
3. **Hiệu suất ghi kém:** Việc ghi dữ liệu trực tiếp vào lớp ghi của container đòi hỏi một trình điều khiển lưu trữ (storage driver) để quản lý, làm cho tốc độ đọc/ghi dữ liệu chậm hơn so với việc ghi trực tiếp xuống hệ thống tệp của máy chủ.

> 💡 **Tip:** Hãy tưởng tượng container như một chiếc máy tính công cộng thuê theo giờ tại quán net. Khi bạn trả máy (xóa container), toàn bộ dữ liệu trên ổ cứng sẽ bị hệ thống reset đóng băng. Để lưu lại tài liệu bài tập quan trọng (như Database hay Logs), bạn bắt buộc phải mang theo một "ổ cứng di động" (Volume/Bind mount) để cắm vào máy tính đó.

---

## 2. Ba loại Storage trong Docker

Để giải quyết vấn đề lưu trữ dữ liệu bền vững, Docker cung cấp 3 phương pháp chính để mount (gắn) dữ liệu từ máy chủ (host) vào bên trong hệ thống tệp của container:

1. **Volume:** Tùy chọn tốt nhất và được khuyên dùng nhất. Dữ liệu được lưu trữ trong một khu vực an toàn trên máy host do Docker quản lý độc quyền.
2. **Bind Mount:** Cho phép bạn chỉ định một thư mục bất kỳ trên máy host để gắn vào container. Rất linh hoạt nhưng mang lại một số rủi ro bảo mật và tính tương thích.
3. **tmpfs Mount:** Chỉ lưu trữ dữ liệu tạm thời trên bộ nhớ RAM của máy host, hoàn toàn không ghi xuống ổ cứng vật lý.

---

## 3. Docker Volume: Lưu trữ tiêu chuẩn

Volume là cơ chế chuẩn mực để lưu trữ dữ liệu bền vững. Khi tạo một volume, nó được Docker lưu trong thư mục `/var/lib/docker/volumes/` trên hệ điều hành Linux (hoặc thư mục tương đương trên máy ảo Docker của Windows/Mac). Các tiến trình bên ngoài Docker không được phép can thiệp vào thư mục này.

### Quản lý Docker Volume

Bạn có thể dễ dàng tạo và quản lý volume thông qua CLI:

```bash
# Tạo một volume mới có tên là 'my-db-data'
docker volume create my-db-data

# Liệt kê danh sách tất cả các volumes đang tồn tại
docker volume ls

# Xem thông tin chi tiết cấu hình và đường dẫn thực của volume
docker volume inspect my-db-data

# Xóa một volume không còn dùng đến
# (Lưu ý: Bạn không thể xóa volume đang được gắn vào một container)
docker volume rm my-db-data
```

**Kết quả khi chạy lệnh docker volume ls:**
```
DRIVER    VOLUME NAME
local     my-db-data
```

### Mount Volume vào Container

Có hai cú pháp chính để gắn volume vào container: `-v` (cú pháp ngắn gọn, truyền thống) và `--mount` (cú pháp mới, tường minh, được Docker khuyến khích sử dụng).

```bash
# Sử dụng cú pháp -v (volume_name:container_path)
docker run -d --name mydb -v my-db-data:/var/lib/mysql mysql:8.0

# Sử dụng cú pháp --mount (khuyên dùng vì cấu trúc rõ ràng, khó nhầm lẫn)
docker run -d \
  --name mydb2 \
  --mount source=my-db-data,target=/var/lib/mysql \
  mysql:8.0
```

> 📌 **Ghi nhớ:** Khi dùng Volume, bạn không cần quan tâm nó nằm chính xác ở thư mục nào trên ổ cứng máy host, Docker sẽ tự động lo việc cấp phát, gắn kết và dọn dẹp không gian đó.

---

## 4. Bind Mount: Kết nối trực tiếp hệ thống tệp

Bind Mount cho phép bạn gắn trực tiếp một tệp tin hoặc thư mục cụ thể từ máy tính của bạn (máy host) vào bên trong container với đường dẫn tuyệt đối.

### Đặc điểm của Bind Mount
- **Ưu điểm:** Rất hữu ích trong môi trường phát triển (Development). Khi bạn sửa source code trên máy tính, thay đổi sẽ lập tức xuất hiện bên trong container mà không cần phải build lại Docker image (Live reload).
- **Nhược điểm:** Phụ thuộc rất nhiều vào cấu trúc thư mục của hệ điều hành. Ví dụ đường dẫn `C:\Users\` trên Windows sẽ không hoạt động nếu bạn mang container lên server chạy Linux.

### Cú pháp sử dụng Bind Mount

Bạn bắt buộc phải cung cấp đường dẫn tuyệt đối của thư mục trên máy host:

```bash
# Cú pháp -v (absolute_host_path:container_path)
# Sử dụng $(pwd) trên Linux/Mac hoặc ${PWD} trên PowerShell để lấy thư mục hiện tại
docker run -d -p 8080:80 \
  -v $(pwd)/src:/usr/share/nginx/html \
  nginx:latest

# Cú pháp --mount
docker run -d \
  --name dev-web \
  --mount type=bind,source="$(pwd)"/src,target=/app \
  node:18
```

> ⚠️ **Lưu ý:** Với Bind mount, container có quyền đọc, ghi, thậm chí vô tình xóa các file hệ thống quan trọng trên máy host của bạn nếu bạn cấu hình nhầm thư mục. Hãy cực kỳ cẩn thận khi sử dụng cơ chế này!

---

## 5. tmpfs Mount: Bộ nhớ tạm siêu tốc

Khác với hai phương pháp trên, `tmpfs` mount không bao giờ lưu dữ liệu xuống ổ cứng. Nó tạo ra một vùng lưu trữ ngay trên bộ nhớ RAM. 

**Khi nào nên dùng tmpfs?**
- Khi cần xử lý dữ liệu với tốc độ siêu nhanh (RAM nhanh hơn rất nhiều so với SSD).
- Khi bạn có các tệp tin chứa dữ liệu nhạy cảm (chìa khóa bảo mật, mật khẩu tạm thời) và không muốn lưu lại bất cứ dấu vết nào trên hệ thống tệp khi container dừng hoạt động.

---

## 6. So sánh Volume vs Bind Mount & Best Practices

| Khái niệm | Docker Volume | Bind Mount |
|---|---|---|
| **Quyền quản lý** | Hoàn toàn bởi Docker | Trực tiếp trên máy Host (bởi người dùng) |
| **Tính độc lập** | Rất cao, dễ dàng di chuyển, sao lưu giữa các host | Thấp, phụ thuộc chặt chẽ vào hệ điều hành và cấu trúc thư mục máy Host |
| **Bảo mật & Cách ly**| Cao, quy trình thao tác giới hạn qua Docker | Thấp, container có thể sửa/xóa nhầm file trên máy Host |
| **Mục đích chính** | Database, Logs, Dữ liệu Production | Ánh xạ Source code khi lập trình (Live reload) |

**Best Practices (Thực hành tốt nhất):**
- Dùng **Volume** cho môi trường Production, đặc biệt là để bảo vệ dữ liệu của các cơ sở dữ liệu (MySQL, PostgreSQL, MongoDB, Redis...).
- Dùng **Bind Mount** cho môi trường Development, giúp lập trình viên có thể viết code trên IDE yêu thích và ứng dụng tự động cập nhật ngay lập tức.

---

## Bài tập

### 🟢 Cơ bản
1. Tạo một Docker volume có tên là `my_website_data`. Chạy lệnh liệt kê danh sách các volume để xác nhận rằng volume của bạn đã được tạo thành công.
2. Chạy một container Nginx, sử dụng cú pháp `--mount` để gắn volume `my_website_data` vừa tạo vào thư mục `/usr/share/nginx/html` bên trong container.

### 🟡 Trung bình
3. Tạo một thư mục `html_source` trên màn hình Desktop của bạn, bên trong tạo một file `index.html` với nội dung `<h1>Hello từ máy host!</h1>`. Khởi chạy một container Nginx sử dụng **Bind mount** để ánh xạ thư mục `html_source` này vào vị trí `/usr/share/nginx/html` của container Nginx. Mở trình duyệt truy cập http://localhost để xem kết quả. Sau đó, tiếp tục sửa file `index.html` trên Desktop và tải lại trang web để kiểm chứng tính năng Live Reload.

### 🟠 Nâng cao
4. Thử nghiệm sự bền vững của Volume với Database: 
   - Khởi chạy một container MySQL 8.0, cấu hình biến môi trường `MYSQL_ROOT_PASSWORD` và gắn một volume có tên `mysql_db_data` vào `/var/lib/mysql`.
   - Kết nối vào container, tạo một database mới có tên là `test_volume_db`.
   - Buộc xóa (force remove) container MySQL hiện tại.
   - Tạo một container MySQL hoàn toàn mới nhưng tiếp tục sử dụng lại volume `mysql_db_data`.
   - Kiểm tra xem database `test_volume_db` có còn tồn tại trong container mới hay không để rút ra kết luận về Volume.

---

## Tiếp theo
→ [Bài 8: Docker Network - Kết nối các Container](./8_networking.md)
