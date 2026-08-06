<Bài 8: Docker Networking>
## Mục tiêu
- Hiểu cách thức hoạt động của Docker Networking.
- Nắm vững các loại network driver: `bridge`, `host`, `none`, `overlay`.
- Biết cách ánh xạ cổng (port mapping) từ host vào container.
- Tự tạo và quản lý custom bridge network (User-defined bridge).
- Nắm được cách các container kết nối và phân giải tên miền (DNS) lẫn nhau.

---

## 1. Tổng quan về Docker Networking

Trong Docker, mỗi container khi chạy sẽ có một **network stack riêng** (có địa chỉ IP, bảng định tuyến, interface riêng). Điều này giúp các container hoàn toàn bị cô lập với nhau và với máy host (máy chủ vật lý/ảo đang chạy Docker).

Để các container có thể "nói chuyện" được với nhau hoặc giao tiếp với mạng Internet bên ngoài, Docker cung cấp các cơ chế Networking.

> 💡 **Tip:** Hãy tưởng tượng mỗi container là một ngôi nhà, và Docker Network là hệ thống đường xá, cầu cống giúp các ngôi nhà kết nối với nhau hoặc đi ra ngoài đường lớn (Internet).

---

## 2. Các loại Network Driver trong Docker

Docker sử dụng các "driver" để quyết định cách cấu hình mạng cho container. Dưới đây là các loại phổ biến nhất:

### `bridge` (Mặc định)
Đây là driver mặc định khi bạn khởi chạy một container. Các container sẽ được kết nối với một mạng nội bộ ảo (virtual bridge) có tên là `docker0` trên máy host.
- **Đặc điểm:** Tốt cho các container chạy độc lập cần giao tiếp nội bộ. Tuy nhiên, các container trong default bridge network chỉ có thể giao tiếp thông qua địa chỉ IP, không qua tên (hostname).

### `host`
Container sẽ sử dụng trực tiếp network stack của máy host, không có sự cô lập nào.
- **Đặc điểm:** Nếu container lắng nghe trên cổng 80, cổng 80 của máy host cũng sẽ được sử dụng. Phù hợp cho các ứng dụng cần tối ưu hiệu năng mạng tuyệt đối, nhưng làm giảm đi tính bảo mật và dễ đụng độ (conflict) cổng.

### `none`
Container hoàn toàn không có bất kỳ kết nối mạng nào (loopback interface duy nhất).
- **Đặc điểm:** Dùng khi bạn muốn container chạy trong một môi trường cô lập mạng 100%, hoặc bạn muốn tự cấu hình mạng thủ công bằng các công cụ tùy chỉnh.

### `overlay`
Sử dụng để kết nối nhiều Docker daemon với nhau (thường dùng trong Docker Swarm).
- **Đặc điểm:** Cho phép các container nằm trên các máy host vật lý khác nhau có thể giao tiếp an toàn và liền mạch.

---

## 3. Ánh xạ cổng (Port Mapping)

Mặc định, các ứng dụng chạy bên trong container không thể được truy cập từ bên ngoài máy host. Để làm được điều này, chúng ta cần **port mapping** (ánh xạ cổng) bằng cờ `-p`.

Cú pháp: `-p <host_port>:<container_port>`

### Ví dụ
Chạy Nginx và map cổng 8080 của máy host vào cổng 80 của container:

```bash
# Chạy nginx container và map port
docker run -d --name my_web -p 8080:80 nginx
```

**Kết quả:**
```
Bạn có thể mở trình duyệt và truy cập http://localhost:8080 để xem trang mặc định của Nginx. Dù Nginx bên trong đang chạy ở port 80, nhưng ta có thể tiếp cận qua port 8080 trên host.
```

---

## 4. User-defined Bridge Network

Mặc dù `bridge` mặc định (`docker0`) khá dễ dùng, Docker luôn khuyến cáo bạn nên tự tạo **User-defined bridge network** cho các ứng dụng thực tế.

### Sự khác biệt cốt lõi
Trong User-defined network, các container có thể **tự động phân giải tên của nhau (DNS resolution)**. Nghĩa là thay vì phải nhớ IP của container kia (vốn luôn thay đổi), bạn chỉ cần gọi tên (name) của container đó. Default bridge không hỗ trợ tính năng này!

### Tạo mới và quản lý network

```bash
# Liệt kê tất cả các network hiện có
docker network ls

# Tạo một network mới dạng bridge
docker network create my_custom_net

# Xem chi tiết cấu hình của network (dải IP, các container đang kết nối)
docker network inspect my_custom_net
```

**Kết quả:**
```
NETWORK ID     NAME            DRIVER    SCOPE
5bc2174c8df4   bridge          bridge    local
8612347cb192   host            host      local
992abcxyz123   my_custom_net   bridge    local
0e8dc3a123bc   none            null      local
```

---

## 5. Kết nối Container qua Network (Ví dụ thực tế)

Giả sử bạn cần chạy một ứng dụng web kết nối tới một cơ sở dữ liệu MySQL.

**Bước 1:** Tạo network chung
```bash
docker network create app_net
```

**Bước 2:** Chạy Database container trong `app_net`
```bash
docker run -d \
  --name my_db \
  --network app_net \
  -e MYSQL_ROOT_PASSWORD=secret \
  mysql:5.7
```

**Bước 3:** Chạy Web container trong `app_net` và link tới Database
(Thay vì ứng dụng thật, ta dùng một container ubuntu để test kết nối)
```bash
docker run -it --name my_app --network app_net ubuntu bash

# Bên trong my_app container, thử ping database bằng tên container
apt-get update && apt-get install -y iputils-ping
ping my_db
```

**Kết quả:**
```
PING my_db (172.18.0.2) 56(84) bytes of data.
64 bytes from my_db.app_net (172.18.0.2): icmp_seq=1 ttl=64 time=0.104 ms
```
> 📌 **Ghi nhớ:** Container `my_app` ping thành công `my_db` chỉ thông qua tên `my_db` thay vì IP!

Bạn cũng có thể linh hoạt gắn/ngắt kết nối một container đang chạy vào network:
```bash
# Kết nối container my_web vào app_net
docker network connect app_net my_web

# Ngắt kết nối my_web khỏi app_net
docker network disconnect app_net my_web
```

---

## Bài tập

1. Tạo một network mới có tên `backend_net`.
2. Khởi chạy một container Redis có tên `redis_cache` và kết nối nó vào `backend_net`.
3. Chạy lệnh liệt kê danh sách network để kiểm tra `backend_net` đã xuất hiện.
4. Khởi chạy một container Nginx thứ hai có tên `web_server`, sử dụng network mặc định (không truyền `--network`). Ánh xạ port 9090 trên máy host vào port 80 của container này.
5. Dùng lệnh `docker network connect` để thêm `web_server` vào `backend_net` lúc nãy.
6. Chạy lệnh `docker network inspect backend_net` để kiểm tra rằng giờ đây cả `redis_cache` và `web_server` đều thuộc chung một mạng.
7. Khởi chạy một Alpine container sử dụng driver mạng là `--network host`. Cài đặt và dùng công cụ `curl` (nếu chưa có) để truy cập `http://localhost:9090` (chính là Nginx bạn vừa map ra ở bài Trung bình). Giải thích tại sao Alpine lại dùng được từ khóa `localhost` ở đây.
8. Viết một lệnh để xóa sạch tất cả các network rác (không có container nào đang kết nối) trên hệ thống (Gợi ý: Dùng lệnh prune).

---

## Tiếp theo
Bài 9: Biến môi trường & Cấu hình → [Bài 9: Environment Variables & Config](./9_env_config.md)
