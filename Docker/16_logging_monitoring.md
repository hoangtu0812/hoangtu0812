# Bài 16: Logging & Monitoring

## Mục tiêu
- Hiểu và cấu hình các Docker logging drivers (json-file, syslog, journald, fluentd).
- Sử dụng thành thạo lệnh `docker logs` cùng các tham số nâng cao để tra cứu lỗi.
- Nắm bắt cách quản lý kích thước log (log rotation) để bảo vệ không gian ổ đĩa.
- Hiểu về các giải pháp Centralized Logging hiện đại (ELK Stack, Fluentd, Loki).
- Giám sát tài nguyên container bằng `docker stats`, `docker system df`, cAdvisor và Prometheus.
- Áp dụng HEALTHCHECK để hệ thống tự động theo dõi trạng thái ứng dụng.

---

## 1. Docker Logging Drivers là gì?

> 💡 **Tip:** Hãy tưởng tượng logging driver như một người thư ký ghi chép chuyên nghiệp. Mỗi khi container của bạn "nói" điều gì đó (xuất dữ liệu ra `stdout` hoặc báo lỗi ra `stderr`), người thư ký này sẽ quyết định xem nên lưu lời nói đó vào đâu và định dạng như thế nào.

Mặc định, Docker sử dụng `json-file` logging driver. Nghĩa là mọi log sẽ được lưu thành các tập tin dạng JSON nằm trên máy chủ Host. Tuy nhiên, khi hệ thống lớn lên, Docker hỗ trợ nhiều driver khác nhau để đáp ứng nhu cầu:
- **json-file** (Mặc định): Lưu log dưới định dạng JSON ngay trên máy chủ.
- **syslog**: Gửi log trực tiếp đến syslog daemon của hệ điều hành Linux.
- **journald**: Gửi log đến systemd journal, hữu ích cho các hệ thống dùng systemd.
- **fluentd**: Gửi log đến Fluentd daemon (rất phổ biến để thiết lập log tập trung).
- **awslogs**: Đẩy log trực tiếp lên dịch vụ Amazon CloudWatch Logs.

### Cấu hình Logging Driver

Bạn có thể cấu hình logging driver linh hoạt cho từng container hoặc thiết lập mặc định toàn cục cho Docker daemon.

**Cấu hình Per-container (Từng container):**
```bash
# Chạy một container Nginx và gửi log trực tiếp đến máy chủ syslog từ xa
docker run -d --name my-app \
  --log-driver syslog \
  --log-opt syslog-address=udp://192.168.1.100:514 \
  nginx
```

**Cấu hình Global (Toàn hệ thống):**
Sửa tập tin `/etc/docker/daemon.json` (hoặc trong mục cấu hình của Docker Desktop):
```json
{
  "log-driver": "local",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```
*Lưu ý: Sau khi sửa file này, bạn cần khởi động lại Docker daemon (ví dụ: `systemctl restart docker`) để thay đổi có hiệu lực.*

**Kết quả:**
```
Tất cả các container được tạo mới sau khi cấu hình sẽ tự động sử dụng driver 'local' cùng với các giới hạn về kích thước đã thiết lập.
```

---

## 2. Đọc Log với `docker logs`

Lệnh `docker logs` là người bạn đồng hành hàng ngày của lập trình viên và quản trị viên hệ thống để debug và tìm hiểu nguyên nhân lỗi.

```bash
# Xem toàn bộ log của container từ lúc khởi động
docker logs my-app

# Theo dõi log theo thời gian thực (tương tự như lệnh tail -f trên Linux)
docker logs --follow my-app

# Chỉ trích xuất 50 dòng cuối cùng (rất hữu ích với log dài)
docker logs --tail 50 my-app

# Xem log kèm theo nhãn thời gian chuẩn (timestamps)
docker logs --timestamps my-app

# Lọc log theo khoảng thời gian cụ thể (rất mạnh mẽ khi tìm lỗi theo giờ)
docker logs --since 1h my-app
docker logs --until 10m my-app
docker logs --since 2023-10-15T10:00:00 --until 2023-10-15T11:00:00 my-app
```

> ⚠️ **Lưu ý:** Lệnh `docker logs` thông thường chỉ hoạt động với các driver `json-file`, `local`, và `journald`. Nếu bạn cấu hình đẩy log ra ngoài bằng syslog hoặc awslogs, bạn sẽ không thể dùng lệnh này mà phải vào hệ thống đích để xem log.

---

## 3. Quản lý kích thước Log (Log Rotation)

Một lỗi "kinh điển" của người mới dùng Docker là để các container chạy liên tục nhiều tháng mà không giới hạn file log. Hậu quả là file JSON của log có thể phình to lên hàng chục GB và làm đầy 100% dung lượng ổ cứng.

**Giải pháp:** Luôn sử dụng Log Rotation (Xoay vòng log).

```bash
# Chạy container với giới hạn log: Tối đa 10MB mỗi file, giữ lại tối đa 3 file cũ nhất
docker run -d --name app-with-rotation \
  --log-driver json-file \
  --log-opt max-size=10m \
  --log-opt max-file=3 \
  nginx
```

---

## 4. Centralized Logging (Quản lý log tập trung)

Khi triển khai hệ thống với kiến trúc Microservices gồm hàng chục hoặc hàng trăm container nằm rải rác trên nhiều máy chủ, việc SSH vào từng máy để gõ `docker logs` là bất khả thi. Bạn cần các hệ thống Centralized Logging:

- **ELK Stack (Elasticsearch + Logstash + Kibana):** Đây là "ông lớn" trong mảng logging. Logstash chịu trách nhiệm thu thập, lọc và chuẩn hóa log; Elasticsearch index dữ liệu cực nhanh; và Kibana cung cấp giao diện biểu đồ tìm kiếm tuyệt đẹp.
- **Fluentd / Fluent Bit:** Là một dự án mã nguồn mở siêu nhẹ, thường được dùng như một "ống nước" (data collector) để thu thập log từ Docker và đẩy về hệ thống lưu trữ như Elasticsearch hay S3.
- **Loki + Grafana:** Giải pháp hiện đại và tiết kiệm tài nguyên đến từ Grafana. Không giống Elasticsearch (index toàn bộ nội dung text), Loki chỉ index các nhãn (labels) metadata, giúp nó chạy cực kỳ nhẹ và nhanh, rất phù hợp cho kiến trúc Kubernetes/Docker.

---

## 5. Monitoring (Giám sát) Docker Containers

### Giám sát tài nguyên theo thời gian thực
Lệnh cơ bản nhất nhưng vô cùng hữu ích:

```bash
# Giám sát CPU, RAM, Network I/O, Disk I/O theo thời gian thực (cập nhật liên tục)
docker stats

# Xem dung lượng ổ cứng tổng thể mà Docker đang sử dụng (Images, Containers, Volumes, Build Cache)
docker system df

# Xem chi tiết từng thành phần đang ngốn dung lượng (thêm cờ -v)
docker system df -v
```

### Các hệ thống Monitoring chuyên sâu
- **cAdvisor (Container Advisor):** Một daemon mã nguồn mở do Google phát triển. Nó chạy ngầm và liên tục phân tích mức tiêu thụ tài nguyên cũng như hiệu năng của tất cả container đang chạy. Bạn có thể triển khai cAdvisor dưới dạng một Docker container duy nhất.
- **Prometheus + Grafana:** Combo giám sát phổ biến nhất hiện nay. Prometheus đóng vai trò Database dạng Time-Series chuyên kéo (pull) dữ liệu metrics định kỳ từ hệ thống, sau đó Grafana kết nối vào Prometheus để hiển thị các Dashboard Dashboard theo dõi trực quan tuyệt vời.

---

## 6. Giám sát sức khỏe với HEALTHCHECK

Làm sao Docker hoặc Load Balancer biết ứng dụng Web bên trong container thực sự đang xử lý được request (trạng thái HTTP 200) hay đang bị "treo" (dù tiến trình node/python vẫn đang chạy)? Đó là lúc `HEALTHCHECK` xuất hiện.

```dockerfile
# Sử dụng Nginx làm base image
FROM nginx:alpine

# Định nghĩa HEALTHCHECK: Mỗi 30s kiểm tra một lần, cho phép timeout 3s, lỗi 3 lần mới coi là unhealthy
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
```

**Kết quả:**
```
Khi chạy 'docker ps', bạn sẽ thấy trạng thái của container được bổ sung thêm: (health: starting) -> (healthy). 
Nếu curl thất bại 3 lần liên tiếp, trạng thái chuyển thành (unhealthy).
```

> 📌 **Ghi nhớ Best Practices:** 
> 1. Luôn sử dụng log rotation (`max-size` và `max-file`) trong `daemon.json`.
> 2. Ghi log từ ứng dụng ra dạng có cấu trúc (Structured Logging như JSON) để các hệ thống ELK/Loki dễ dàng phân tích.
> 3. Luôn định nghĩa `HEALTHCHECK` trong Dockerfile để Docker Swarm hoặc Kubernetes có thể tự động restart các container "bệnh tật".

---

## Bài tập

1. **Khám phá `docker logs`:** Khởi chạy một container Nginx (`docker run -d -p 8080:80 nginx`). Truy cập `http://localhost:8080` trên trình duyệt vài lần. Sau đó dùng lệnh `docker logs --tail 5 --timestamps <container_id>` để xem các dòng log truy cập mới nhất.
2. **Log Rotation:** Khởi chạy một container dùng lệnh lặp vô hạn in ra text. Cấu hình cờ `--log-opt` để đặt kích thước tối đa (`max-size`) là 1MB và giữ lại (`max-file`) 2 file. Theo dõi thư mục vật lý chứa log trên máy tính xem Docker có tự xóa file cũ không.

3. **Thực hành HEALTHCHECK:** Hãy viết một Dockerfile chạy ứng dụng Python đơn giản. Bổ sung cấu hình `HEALTHCHECK` kiểm tra sự tồn tại của một file `/tmp/healthy`. Khởi chạy container, sau đó dùng `docker exec` vào và xóa file đó đi. Quan sát trạng thái chuyển từ `healthy` sang `unhealthy` thông qua lệnh `docker ps`.

4. **Giám sát trực quan với cAdvisor:** Viết một file `docker-compose.yml` định nghĩa 2 services: một ứng dụng web bất kỳ (như WordPress hoặc Nginx) và một service `google/cadvisor`. Chạy `docker-compose up -d`, truy cập vào giao diện web của cAdvisor (thường ở cổng 8080) và khám phá các biểu đồ thống kê CPU, RAM của ứng dụng web.

---

## Tiếp theo
→ [Bài 17: Docker Internals](./17_internals.md)
