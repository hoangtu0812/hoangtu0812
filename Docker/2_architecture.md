# Bài 2: Kiến trúc Docker

## Mục tiêu
- Hiểu rõ các thành phần cấu tạo nên Docker Engine.
- Nắm bắt mô hình Client-Server của Docker.
- Phân biệt sự khác nhau giữa Docker Image và Docker Container.
- Nắm được luồng hoạt động cơ bản từ lúc kéo Image đến khi chạy Container.
- Làm quen với các khái niệm mở rộng như Containerd, runc và tiêu chuẩn OCI.

## 1. Kiến trúc Client-Server của Docker
Docker sử dụng mô hình kiến trúc Client-Server, trong đó các thành phần giao tiếp với nhau qua các quy tắc rõ ràng.

```text
+-------------------+       +-----------------------+       +-------------------+
|   Docker Client   | ----> |     Docker Daemon     | ----> |  Docker Registry  |
|  (docker build,   |  API  | (quản lý container,   |       |   (Docker Hub)    |
|   docker run)     |       |  image, network...)   |       |                   |
+-------------------+       +-----------------------+       +-------------------+
```

- **Docker Client:** Là công cụ dòng lệnh (CLI) chính mà bạn tương tác. Khi bạn gõ các lệnh như `docker run` hoặc `docker build`, Client sẽ đóng gói và gửi các yêu cầu này đến Docker Daemon qua REST API.
- **Docker Daemon (dockerd):** Là quá trình chạy nền trên máy host (máy chủ của bạn). Nó lắng nghe các yêu cầu từ Docker Client thông qua REST API và chịu trách nhiệm xây dựng, chạy, và quản lý toàn bộ các đối tượng của Docker như Images, Containers, Networks và Volumes.
- **REST API:** Là cầu nối giao tiếp chuẩn hóa giữa Client và Daemon. Chúng có thể nằm trên cùng một máy tính vật lý hoặc trên các máy khác nhau kết nối qua mạng Internet/LAN.

## 2. Docker Engine
Docker Engine là phần lõi của toàn bộ hệ sinh thái Docker, bao gồm 3 thành phần chính:
1. Một process chạy nền liên tục (background process) gọi là **Docker Daemon**.
2. Một giao diện lập trình ứng dụng **REST API** cho phép các chương trình khác tương tác với Daemon.
3. Một giao diện dòng lệnh **Docker Client** (CLI) dùng để người dùng ra lệnh.

> 💡 **Tip:** Đa số người dùng cá nhân sẽ cài đặt phần mềm Docker Desktop. Phần mềm này đã bao gồm sẵn toàn bộ Docker Engine (Client, Daemon, API) và một giao diện đồ họa (GUI) để dễ dàng quản lý.

## 3. Docker Registry (Docker Hub)
Docker Registry là "kho chứa" các Docker Images.
- **Docker Hub:** Là một public registry mặc định lớn nhất thế giới, nơi bất kỳ ai cũng có thể lưu trữ và tải xuống các image. Tại đây chứa hàng triệu image có sẵn từ các nhà cung cấp chính thức (như Ubuntu, Nginx, MySQL, Python) cũng như từ cộng đồng.
- Khi bạn dùng lệnh `docker pull` hoặc `docker run`, nếu image bạn cần chưa có sẵn trên máy của bạn (máy host), Docker sẽ tự động kết nối mạng và tải nó về từ Docker Hub.
- Ngoài ra, các công ty thường tự xây dựng các Private Registry (kho chứa nội bộ riêng) để bảo mật mã nguồn.

## 4. Docker Image vs Docker Container
Sự khác biệt giữa Image và Container là một trong những khái niệm quan trọng bậc nhất mà bạn cần nắm vững trước khi đi sâu vào Docker.

- **Docker Image (Bản thiết kế):** Là một mẫu (template) ở chế độ chỉ đọc (read-only). Nó chứa tất cả các hướng dẫn và thành phần để tạo ra một Docker container, bao gồm hệ điều hành cắt giảm, mã nguồn, thư viện, dependencies và các công cụ cần thiết cho ứng dụng.
- **Docker Container (Thực thể hoạt động):** Là một phiên bản (instance) đang chạy thực tế của một Image. Từ một Image duy nhất, bạn có thể tạo ra hàng chục, hàng trăm Container hoàn toàn độc lập với nhau.

> 📌 **Ghi nhớ:**
> - Trong lập trình hướng đối tượng (OOP): **Image** giống như **Class** (Lớp), còn **Container** giống như **Object** (Đối tượng).
> - Trong xây dựng thực tế: **Image** giống như **Bản vẽ thiết kế** (Blueprint), còn **Container** là **Ngôi nhà** thực tế được xây lên từ bản vẽ đó.

## 5. Luồng hoạt động cơ bản (Workflow)

Hãy hình dung những gì xảy ra ngầm bên dưới khi bạn gõ một lệnh Docker cơ bản:

```text
[1] Người dùng gõ: docker run nginx
       |
[2] Docker Client gửi yêu cầu "chạy image nginx" qua REST API tới Docker Daemon.
       |
[3] Docker Daemon nhận lệnh và kiểm tra xem image "nginx" đã có sẵn trên máy chưa?
       |
       |--- Không có -> [4] Daemon kết nối internet, Pull image từ Docker Hub.
       |
       |--- Đã có -> (Bỏ qua bước 4).
       |
[5] Docker Daemon tạo và khởi chạy một Container hoàn toàn mới từ Image "nginx".
```

Tóm tắt các lệnh cơ bản trong luồng:
- **`docker pull`**: Chỉ tải Image từ Registry về máy host mà không chạy nó.
- **`docker build`**: Tự động tạo ra một Image mới dựa trên các hướng dẫn bạn viết trong file cấu hình (Dockerfile).
- **`docker run`**: Khởi chạy một Container từ một Image (bao gồm cả việc tải image nếu nó chưa tồn tại cục bộ).

## 6. Containerd, runc và Tiêu chuẩn OCI (Kiến thức nâng cao)
Trong những phiên bản đầu tiên, Docker tự làm mọi việc từ A đến Z (quản lý image, network, khởi chạy process...). Tuy nhiên, hiện nay kiến trúc này đã được chuẩn hóa và tách biệt ra nhiều thành phần để tăng tính module và hiệu năng.

- **OCI (Open Container Initiative):** Là một dự án đặt ra các tiêu chuẩn mở toàn cầu cho cấu trúc của image (Image Spec) và runtime của container (Runtime Spec). Điều này giúp các công cụ khác nhau (như Docker, Podman, Kubernetes) có thể giao tiếp và chạy container của nhau mà không gặp lỗi.
- **runc:** Là một công cụ thực thi tiêu chuẩn OCI ở tầng rất thấp (low-level runtime). Nhiệm vụ duy nhất và quan trọng nhất của nó là tương tác với lõi hệ điều hành (Linux kernel - cụ thể là cgroups và namespaces) để tạo ra các container processes.
- **Containerd:** Là một runtime tầng cao (high-level) đứng làm trung gian giữa Docker Daemon và runc. Nó quản lý vòng đời hoàn chỉnh của các container, điều phối quá trình truyền tải image, quản lý storage và chuyển yêu cầu tạo process xuống cho runc.

```text
Docker Client -> Docker Daemon -> Containerd -> runc -> Container (chạy ứng dụng)
```

## Bài tập
1. **Kiến trúc Client-Server**: Giải thích vai trò của Docker Client và Docker Daemon trong quá trình khởi chạy một container. Chúng giao tiếp với nhau bằng phương thức nào?
2. **So sánh cốt lõi**: Bằng ngôn từ của riêng bạn, hãy giải thích sự khác nhau giữa Docker Image và Docker Container. Hãy lấy một ví dụ thực tế khác với ví dụ trong bài để minh họa.
3. **Phân tích luồng**: Khi bạn mở terminal và gõ lệnh `docker pull ubuntu`, hãy mô tả chi tiết theo từng bước cách các thành phần trong kiến trúc Docker tương tác với nhau để hoàn thành lệnh này.
4. **Nghiên cứu mở rộng**: Dựa vào phần kiến thức nâng cao, hãy giải thích ngắn gọn tại sao việc có một tiêu chuẩn chung như OCI lại quan trọng đối với sự phát triển của hệ sinh thái Container hiện đại?

## Tiếp theo
→ [Bài 3: Quản lý Docker Container](./3_containers.md)
