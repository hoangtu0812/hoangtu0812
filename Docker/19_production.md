# Bài 19: Docker trong Production

## Mục tiêu
- Hiểu và áp dụng các tiêu chuẩn (checklist) khi đưa container lên môi trường production.
- Cấu hình tự động khởi động lại (Restart policy) và giới hạn tài nguyên (Resource limits).
- Hiểu cách tắt ứng dụng an toàn (Graceful shutdown) và xử lý tiến trình PID 1.
- Nắm vững cách bảo trì hệ thống Docker và sao lưu dữ liệu.

## 1. Production-ready Checklist

Khi chuyển từ môi trường phát triển (development) sang môi trường thực tế (production), bạn cần đảm bảo container của mình an toàn, nhẹ và ổn định. Dưới đây là các tiêu chí quan trọng nhất:

- **Image nhỏ gọn**: Sử dụng các image nền như `alpine`, `distroless`, hoặc `slim`. Tránh cài đặt các công cụ debug (như `curl`, `vim`, `ping`) trên môi trường production. Sử dụng Multi-stage builds để tối ưu kích thước.
- **Không chạy bằng quyền root**: Luôn sử dụng chỉ thị `USER` trong Dockerfile để chỉ định một người dùng không có đặc quyền. Điều này giúp hạn chế thiệt hại nếu container bị tấn công.
- **Sử dụng HEALTHCHECK**: Báo cho Docker biết cách kiểm tra xem ứng dụng của bạn có thực sự đang hoạt động tốt hay không, thay vì chỉ kiểm tra xem tiến trình có đang chạy không.
- **Giới hạn tài nguyên**: Tránh tình trạng một container "ăn" hết CPU hoặc RAM của máy chủ, gây ảnh hưởng đến các ứng dụng khác.

## 2. Restart Policy (Chính sách khởi động lại)

Trong production, bạn không thể lúc nào cũng túc trực để bật lại container nếu nó bị lỗi hoặc khi server khởi động lại. Docker cung cấp các chính sách tự động khởi động lại.

```bash
# Tự động khởi động lại trừ khi bạn chủ động dừng nó bằng lệnh `docker stop`
docker run -d --restart unless-stopped my-app

# Luôn luôn khởi động lại, bất kể lý do dừng là gì
docker run -d --restart always my-app
```

> 💡 **Tip:** Cấu hình `--restart unless-stopped` là lựa chọn phổ biến nhất cho production vì nó tôn trọng quyết định dừng thủ công của quản trị viên, nhưng vẫn đảm bảo container tự bật lại khi server reboot.

## 3. Graceful Shutdown (Tắt ứng dụng an toàn)

Khi bạn chạy `docker stop`, Docker không ngắt điện đột ngột. Nó gửi tín hiệu **SIGTERM** để yêu cầu ứng dụng dừng lại. Ứng dụng có mặc định **10 giây** để dọn dẹp (đóng kết nối database, ghi nốt log) trước khi Docker gửi tín hiệu **SIGKILL** (bắt buộc kết thúc).

```dockerfile
# Sử dụng STOPSIGNAL để thay đổi tín hiệu mặc định nếu ứng dụng của bạn yêu cầu
STOPSIGNAL SIGQUIT
```

Ứng dụng của bạn (Node.js, Python, Java...) cần được viết để bắt (catch) sự kiện `SIGTERM` và thực hiện dọn dẹp tài nguyên.

## 4. Init Process (Tiến trình Init với Tini)

Trong Linux, tiến trình đầu tiên chạy có PID 1. Nó chịu trách nhiệm dọn dẹp các "zombie process" (tiến trình con đã chết nhưng chưa được giải phóng). Nếu ứng dụng của bạn (ví dụ: Node.js, Java) không xử lý tốt PID 1, hệ thống có thể bị tràn bộ nhớ.

Docker cung cấp cờ `--init` để bọc ứng dụng của bạn bằng một tiến trình init nhỏ gọn tên là `tini`.

```bash
docker run -d --init -p 8080:8080 my-app
```

Hoặc trong Dockerfile:
```dockerfile
RUN apk add --no-cache tini
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["node", "app.js"]
```

## 5. Health Checks Chi Tiết

Chỉ thị `HEALTHCHECK` giúp Docker xác định trạng thái thực tế của ứng dụng.

```dockerfile
# Kiểm tra mỗi 30s, chờ tối đa 3s cho mỗi lần kiểm tra, thử lại 3 lần trước khi đánh dấu là unhealthy
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD curl -f http://localhost/health || exit 1
```

Trạng thái health check có thể được sử dụng bởi các công cụ điều phối (Orchestration tools) hoặc proxy (như Traefik, Nginx) để ngừng gửi request đến các container bị lỗi.

## 6. Resource Limits (Giới hạn tài nguyên)

Không bao giờ để container chạy tự do không kiểm soát trong production.

```bash
docker run -d \
  --name web_app \
  --memory="512m" \
  --cpus="1.5" \
  --pids-limit=100 \
  my-app
```
- `--memory="512m"`: Giới hạn RAM tối đa là 512 Megabytes.
- `--cpus="1.5"`: Container chỉ được dùng tối đa 1.5 nhân CPU.
- `--pids-limit=100`: Tránh tình trạng "Fork bomb" bằng cách giới hạn số tiến trình tối đa container có thể tạo ra.

## 7. Container Logging trong Production

Tránh sử dụng log mặc định (json-file) vô tội vạ. Nó có thể làm đầy ổ cứng server của bạn. Hãy thiết lập giới hạn dung lượng log hoặc sử dụng hệ thống log tập trung (như ELK stack, Fluentd).

```bash
# Giới hạn kích thước mỗi file log là 10MB và chỉ lưu tối đa 3 file
docker run -d \
  --log-opt max-size=10m \
  --log-opt max-file=3 \
  my-app
```

## 8. Bảo trì hệ thống và Dọn dẹp (System Maintenance)

Theo thời gian, các image cũ, container đã dừng, và volume không dùng tới sẽ chiếm dụng ổ cứng.

```bash
# Kiểm tra xem Docker đang chiếm bao nhiêu dung lượng
docker system df

# Dọn dẹp an toàn các tài nguyên không sử dụng (dangling images, stopped containers)
docker system prune

# Dọn dẹp cực mạnh: Xóa mọi thứ không có container nào đang dùng (bao gồm volume)
docker system prune -a --volumes
```

> ⚠️ **Lưu ý:** Lệnh `docker system prune -a --volumes` cực kỳ nguy hiểm, hãy cẩn thận khi sử dụng trên production.

## 9. Backup và Recovery (Sao lưu và Phục hồi)

Dữ liệu quan trọng phải nằm trong Volumes. Dưới đây là cách backup một volume:

```bash
# Chạy một container tạm thời, mount volume cần backup, tạo file tar và lưu ra máy host
docker run --rm -v my_database_volume:/data -v $(pwd):/backup ubuntu tar cvf /backup/backup.tar /data
```

Để phục hồi (Restore):
```bash
docker run --rm -v my_database_volume:/data -v $(pwd):/backup ubuntu tar xvf /backup/backup.tar -C /
```

## Bài tập
1. **Khởi chạy an toàn**: Viết lệnh `docker run` chạy một Nginx container nền, giới hạn 256MB RAM, 0.5 CPU, và sẽ tự động khởi động lại trừ khi bị dừng thủ công.
2. **Cấu hình Log**: Thêm tham số vào lệnh ở Bài 1 để cấu hình giới hạn kích thước file log là 5MB và giữ tối đa 5 file.
3. **Health Check Command**: Viết chỉ thị `HEALTHCHECK` trong Dockerfile cho một ứng dụng Node.js chạy ở cổng 3000. Dùng lệnh `wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1`. Cấu hình interval 15s.
4. **Dọn rác hệ thống**: Viết chuỗi lệnh để kiểm tra dung lượng Docker đang sử dụng, sau đó thực hiện lệnh dọn dẹp an toàn những thành phần không sử dụng (không xóa volume).

## Tiếp theo
→ [Bài 20: Giới thiệu Orchestration](./20_orchestration.md)
