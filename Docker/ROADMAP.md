# Lộ Trình Học Docker Chi Tiết

> Cấu trúc giống [Go/ROADMAP.md](../Go/ROADMAP.md), [Python/ROADMAP.md](../Python/ROADMAP.md), [SAP/ROADMAP.md](../SAP/ROADMAP.md), [MachineLearning/ROADMAP.md](../MachineLearning/ROADMAP.md), [Git/ROADMAP.md](../Git/ROADMAP.md): mỗi bài có file chi tiết riêng (lý thuyết + ý nghĩa cơ chế bên trong + lệnh thực hành + bài tập), kết thúc bằng dự án thực hành mô phỏng quy trình làm việc nhóm hoàn chỉnh.

---

## Giai đoạn 0 — Giới thiệu & Cài đặt

### [Bài 1: Giới thiệu Docker & Cài đặt](./1_get_started.md)
- Container là gì, sự khác biệt giữa Virtual Machine (VM) và Container.
- Cài đặt Docker Desktop, kiểm tra cài đặt.
- **Bài tập:** Cài đặt Docker, chạy container "hello-world" để kiểm tra.

---

## Giai đoạn 1 — Docker Cơ Bản (Fundamentals)

### [Bài 2: Kiến trúc Docker](./2_architecture.md)
- Kiến trúc cơ bản: Docker Engine, Client-Server model.
- Khái niệm về Image, Container, Docker Registry.
- **Bài tập:** Chạy một container và xem thông tin docker system.

### [Bài 3: Làm việc với Container](./3_containers.md)
- Vòng đời container: `run`, `stop`, `start`, `rm`.
- Tương tác với container: `exec`, `logs`, xem các process đang chạy.
- **Bài tập:** Tạo và quản lý nginx container, chạy lệnh bash bên trong container.

### [Bài 4: Làm việc với Image](./4_images.md)
- Quản lý Image: `pull`, `build`, `tag`, `push`.
- Khái niệm về layers trong Image.
- **Bài tập:** Tải một image từ Docker Hub, tạo tag mới và quản lý images nội bộ.

### [Bài 5: Dockerfile cơ bản](./5_dockerfile_basics.md)
- Cấu trúc Dockerfile và các lệnh cơ bản: `FROM`, `RUN`, `COPY`, `CMD`, `EXPOSE`.
- **Bài tập:** Viết Dockerfile cho một ứng dụng Node.js/Python đơn giản.

### [Bài 6: .dockerignore & Build Context](./6_dockerignore.md)
- Tối ưu build context bằng `.dockerignore`.
- Ảnh hưởng của kích thước context đến tốc độ build.
- **Bài tập:** Tạo file `.dockerignore` để loại bỏ `node_modules` hoặc thư mục cache khi build image.

---

## Giai đoạn 2 — Docker Trung Cấp (Intermediate)

### [Bài 7: Volume & Bind Mount](./7_volumes.md)
- Cách xử lý dữ liệu bền vững (persistent data).
- Phân biệt Volume (do Docker quản lý) và Bind Mount (file hệ thống máy host).
- **Bài tập:** Tạo database container sử dụng Volume để giữ dữ liệu sau khi xóa container.

### [Bài 8: Docker Networking](./8_networking.md)
- Các loại network: `bridge`, `host`, `none`.
- Giao tiếp giữa các container (container linking, user-defined network).
- **Bài tập:** Tạo mạng bridge riêng và cho hai container (web và db) ping/kết nối với nhau.

### [Bài 9: Environment Variables & Cấu hình](./9_env_config.md)
- Truyền biến môi trường bằng cờ `-e` và file `.env`.
- **Bài tập:** Khởi chạy container với thông tin kết nối database thông qua biến môi trường.

### [Bài 10: Dockerfile nâng cao](./10_dockerfile_advanced.md)
- Tối ưu với Multi-stage build.
- Sử dụng `ARG`, và sự khác biệt giữa `ENTRYPOINT` vs `CMD`.
- **Bài tập:** Viết Dockerfile multi-stage để build và run ứng dụng (ví dụ Go, React).

### [Bài 11: Docker Compose cơ bản](./11_compose_basics.md)
- Giới thiệu Docker Compose, file `docker-compose.yml`.
- Cấu hình services, networks, volumes trong Compose.
- **Bài tập:** Viết `docker-compose.yml` chạy một app web cùng database Redis/Postgres.

### [Bài 12: Docker Compose nâng cao](./12_compose_advanced.md)
- Quản lý thứ tự khởi động với `depends_on`.
- Cấu hình `healthcheck`, sử dụng các `profiles`.
- **Bài tập:** Bổ sung healthcheck cho db, chỉ khởi động app khi db thực sự sẵn sàng.

---

## Giai đoạn 3 — Docker Nâng Cao (Advanced)

### [Bài 13: Tối ưu hóa Image](./13_image_optimization.md)
- Giảm kích thước image (dùng alpine, distroless).
- Tận dụng cache layer để tăng tốc độ build.
- **Bài tập:** Refactor Dockerfile để giảm kích thước image gấp nhiều lần.

### [Bài 14: Bảo mật Container](./14_security.md)
- Chạy container với user non-root.
- Quản lý secrets an toàn, scan image tìm lỗ hổng (vulnerability scanning).
- **Bài tập:** Cấu hình user non-root cho container, chạy scan image bằng công cụ như Trivy/Docker Scout.

### [Bài 15: Docker Registry](./15_registry.md)
- Làm việc chuyên sâu với Docker Hub.
- Tự triển khai Private Registry nội bộ.
- **Bài tập:** Push image lên Docker Hub và cấu hình chạy local registry.

### [Bài 16: Logging & Monitoring](./16_logging_monitoring.md)
- Các logging drivers (json-file, syslog, splunk).
- Theo dõi tài nguyên (`docker stats`), tích hợp Prometheus/Grafana.
- **Bài tập:** Xem tài nguyên tiêu thụ và cấu hình giới hạn memory/cpu cho container.

### [Bài 17: Docker Internals](./17_internals.md)
- Hiểu sâu về cách Docker hoạt động dưới nền (namespaces, cgroups, union filesystem).
- **Bài tập:** Quan sát process tree để hiểu namespaces & cgroups.

### [Bài 18: CI/CD với Docker](./18_cicd.md)
- Tích hợp Docker vào quá trình CI/CD.
- Build, test và push image tự động bằng GitHub Actions.
- **Bài tập:** Viết workflow GitHub Actions cơ bản.

### [Bài 19: Docker trong Production](./19_production.md)
- Các best practices khi mang lên môi trường thực tế.
- Restart policies, giới hạn tài nguyên an toàn.
- **Bài tập:** Chạy ứng dụng giả lập production với `restart: always`.

### [Bài 20: Giới thiệu Orchestration](./20_orchestration.md)
- Nhu cầu về Container Orchestration khi hệ thống lớn.
- Tổng quan Docker Swarm & Kubernetes.
- **Bài tập:** Chạy một Swarm cluster nhỏ (1 node) hoặc Minikube.

---

## Giai đoạn 4 — Dự Án Thực Hành (Capstone)

### [Bài 21: Dự án Capstone](./21_capstone_project.md)
- Containerize một ứng dụng multi-service (Frontend, Backend API, Database, Cache).
- **Bài tập:** Thiết kế Dockerfile hiệu quả cho các service, kết nối bằng `docker-compose.yml`, áp dụng healthcheck và volumes.

---

## Gợi ý cách học
1. **Lý thuyết đi đôi với thực hành:** Hãy luôn mở terminal và tự tay gõ các lệnh để xem kết quả trực tiếp.
2. **Theo dõi Logs & Inspect:** Sử dụng `docker logs` và `docker inspect` thường xuyên để debug khi container không chạy như ý muốn.
3. **Build - Break - Fix:** Đừng ngại làm sai. Thử xóa container, phá file cấu hình và xem thông báo lỗi để hiểu rõ cơ chế bên dưới.
4. **Tối ưu từng chút:** Sau khi ứng dụng chạy được, hãy thử giảm dung lượng image hoặc tăng cường bảo mật.

## Tài liệu tham khảo
- [Trang chủ Docker Docs](https://docs.docker.com/)
- [Docker Hub](https://hub.docker.com/)
- [Dockerfile Reference](https://docs.docker.com/engine/reference/builder/)
- [Docker Compose Reference](https://docs.docker.com/compose/compose-file/)
