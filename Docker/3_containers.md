# Bài 3: Làm việc với Container

## Mục tiêu
- Hiểu được vòng đời chi tiết của một Docker container từ lúc khởi tạo đến khi bị hủy.
- Nắm vững các lệnh cơ bản và nâng cao để quản lý container (tạo, chạy, dừng, xóa).
- Biết cách tương tác trực tiếp với môi trường bên trong container.
- Học cách gỡ lỗi (debug) thông qua việc xem log và kiểm tra cấu hình chi tiết của container.
- Thực hiện sao chép dữ liệu tĩnh giữa máy host và container.

---

## 1. Vòng đời của một Container

### Các trạng thái chính
Hãy tưởng tượng container như một chương trình máy tính. Nó không tồn tại vĩnh viễn mà có một vòng đời cụ thể. Một Docker container sẽ trải qua các trạng thái khác nhau:
- **Created (Đã tạo)**: Container đã được tạo ra từ image nhưng chưa bắt đầu chạy các tiến trình bên trong. Nó giống như việc bạn đã cài đặt ứng dụng nhưng chưa mở lên.
- **Running (Đang chạy)**: Container đang thực thi các tiến trình. Ứng dụng của bạn đang hoạt động bình thường.
- **Paused (Tạm dừng)**: Tất cả các tiến trình trong container đang bị "đóng băng". Trạng thái bộ nhớ được giữ nguyên, nhưng nó không tiêu thụ CPU.
- **Stopped/Exited (Đã dừng)**: Container đã dừng hoạt động, tiến trình chính đã kết thúc. Tuy nhiên, hệ thống file của nó vẫn còn tồn tại.
- **Removed (Đã xóa)**: Container đã bị xóa hoàn toàn khỏi hệ thống Docker. Tất cả dữ liệu không được mount ra ngoài sẽ bị mất.

---

## 2. Tạo và chạy Container với `docker run`

### Chi tiết lệnh `docker run`
Lệnh `docker run` là trái tim của việc sử dụng Docker. Khi bạn gọi lệnh này, Docker thực tế thực hiện hai việc: `docker create` (tạo container) và `docker start` (chạy container).

### Các cờ (flags) quan trọng
- `-d` (--detach): Chạy container ngầm (background mode). Bạn sẽ nhận lại container ID và tiếp tục sử dụng terminal.
- `-it`: Thường đi liền với nhau (`-i` cho interactive và `-t` để cấp phát pseudo-TTY). Cần thiết khi bạn muốn mở một shell tương tác bên trong container.
- `--name`: Đặt tên tùy chỉnh cho container. Nếu không có cờ này, Docker sẽ tạo ra một cái tên ngẫu nhiên.
- `-p` (--publish): Ánh xạ cổng. Cú pháp `host_port:container_port`. Nó mở "cánh cửa" để từ bên ngoài (máy của bạn) có thể truy cập vào ứng dụng bên trong container.
- `--rm`: Tự động dọn dẹp và xóa container ngay lập tức sau khi tiến trình bên trong kết thúc. Rất hữu ích cho các tác vụ chạy một lần.

```bash
# Ví dụ 1: Chạy một web server Nginx chạy ngầm, tên là "my-website"
docker run -d --name my-website -p 8080:80 nginx

# Ví dụ 2: Tạo một môi trường Ubuntu tạm thời để vọc vạch
docker run -it --rm ubuntu bash
```

**Kết quả:**
```
12a3f8b9e0... (Chuỗi ID dài được trả về nếu bạn sử dụng cờ -d)
root@container_id:/# (Dấu nhắc lệnh của Ubuntu nếu dùng -it)
```

> 💡 **Tip:** Luôn sử dụng cờ `--name` để dễ dàng quản lý và gọi tên container trong các lệnh sau này.

---

## 3. Xem danh sách các Container

### Sử dụng lệnh `docker ps`
Để biết được hệ thống của bạn đang có những container nào, hãy sử dụng lệnh `docker ps`.

```bash
# Liệt kê danh sách các container ĐANG chạy ở thời điểm hiện tại
docker ps

# Liệt kê TẤT CẢ container ở mọi trạng thái (đang chạy, đã dừng, đã tạo...)
docker ps -a
```

**Kết quả:**
```
CONTAINER ID   IMAGE     COMMAND                  CREATED          STATUS          PORTS                  NAMES
12a3f8b9e0xx   nginx     "/docker-entrypoint.…"   10 minutes ago   Up 10 minutes   0.0.0.0:8080->80/tcp   my-website
```

---

## 4. Quản lý trạng thái Container

### Start, Stop và Restart
Giống như máy ảo, bạn có thể bật, tắt, hoặc khởi động lại container bất cứ lúc nào.

```bash
# Yêu cầu container dừng lại một cách nhẹ nhàng
docker stop my-website

# Khởi động lại một container đã bị dừng trước đó
docker start my-website

# Dừng hẳn và ngay lập tức khởi động lại container
docker restart my-website
```

**Kết quả:**
```
my-website
```

---

## 5. Dọn dẹp và Xóa Container

### Xóa an toàn
Những container đã dừng vẫn chiếm dung lượng đĩa. Để giữ hệ thống sạch sẽ, bạn cần xóa chúng đi.

```bash
# Xóa một container (yêu cầu container đó phải đang ở trạng thái dừng)
docker rm my-website

# Lệnh "dọn nhà" mạnh mẽ: Xóa TẤT CẢ các container đang ở trạng thái dừng
docker container prune
```

> ⚠️ **Lưu ý:** Lệnh `docker container prune` không thể hoàn tác, hãy cẩn thận khi sử dụng! Nếu bạn muốn xóa một container đang chạy, bạn phải dùng thêm cờ `-f`: `docker rm -f my-website`.

---

## 6. Tương tác với Container đang chạy

### Lệnh `docker exec`
Lệnh `docker exec` vô cùng giá trị khi bạn cần khắc phục sự cố hoặc thay đổi cấu hình bên trong một container đang hoạt động.

```bash
# Mở một terminal tương tác (bash) ngay bên trong container "my-website" đang chạy
docker exec -it my-website bash

# Hoặc chỉ chạy một lệnh duy nhất rồi trả về kết quả
docker exec my-website ls -la /usr/share/nginx/html
```

**Kết quả:**
```
root@12a3f8b9e0xx:/# 
```

---

## 7. Theo dõi Log và Kiểm tra thông tin

### Inspect và Logs
Khi ứng dụng gặp lỗi, bạn cần kiểm tra log và cấu hình.

```bash
# Xem toàn bộ lịch sử log của container
docker logs my-website

# Theo dõi log theo thời gian thực (như lệnh tail -f)
docker logs -f my-website

# Lấy thông tin cấu hình chi tiết của container (định dạng JSON)
docker inspect my-website
```

**Kết quả:**
```
192.168.1.5 - - [10/Oct/2023:14:32:01 +0000] "GET / HTTP/1.1" 200 615 "-" "Mozilla/5.0..."
```

---

## 8. Sao chép dữ liệu bằng `docker cp`

### Copy file qua lại
Đôi khi bạn muốn copy nhanh một file cấu hình hoặc file log giữa máy thật và container.

```bash
# Copy một file từ máy tính của bạn (Host) vào bên trong Container
docker cp ./custom-index.html my-website:/usr/share/nginx/html/index.html

# Copy toàn bộ một thư mục từ Container ra máy tính của bạn
docker cp my-website:/etc/nginx/conf.d ./my-local-conf
```

**Kết quả:**
```
Successfully copied 2.05kB to my-website:/usr/share/nginx/html/index.html
```

---

## Bài tập

1. **Khởi chạy và Cấu hình**: Chạy một container từ image `nginx`, đặt tên là `web-server`, chạy chế độ ngầm và ánh xạ cổng 80 của container ra cổng 9090 trên máy host của bạn. Dùng trình duyệt truy cập `http://localhost:9090` để kiểm tra.
2. **Quản lý vòng đời**: Dừng container `web-server` bằng lệnh `docker stop`, sau đó khởi động lại nó bằng `docker start` và kiểm tra xem nó đã chạy lại chưa với lệnh `docker ps`.

3. **Kiểm tra Logs và Tương tác**: Dùng lệnh `docker logs` để xem các lượt truy cập vào `web-server`. Sau đó, sử dụng lệnh `docker exec -it web-server bash` để vào bên trong, tìm đến thư mục `/usr/share/nginx/html/` và sửa nội dung file `index.html` bằng lệnh `echo`. Thoát ra và tải lại trang web để xem sự thay đổi.

4. **Sao chép Dữ liệu và Dọn dẹp**: Tạo một file tên là `test.txt` trên máy tính. Dùng lệnh `docker cp` để copy file đó vào `/tmp` bên trong container `web-server`. Dùng `docker exec` chạy lệnh `ls /tmp` để xác nhận. Sau đó, dừng `web-server`, xóa nó bằng `docker rm`, và dùng `docker container prune` để dọn dẹp hệ thống.


---

## Tiếp theo
→ [Bài 4: Làm việc với Docker Image](./4_images.md)
