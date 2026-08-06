# Bài 20: Giới thiệu Orchestration

## Mục tiêu
- Hiểu được khái niệm Orchestration và tại sao nó lại cần thiết khi triển khai hệ thống có nhiều container.
- Nắm bắt kiến trúc và các lệnh cơ bản của Docker Swarm (Manager, Worker, Service, Stack).
- Có cái nhìn tổng quan về Kubernetes (K8s) và các khái niệm cơ bản (Pod, Deployment, Service).
- So sánh sự khác biệt giữa Docker Swarm và Kubernetes để chọn công cụ phù hợp.
- Biết cách bật và sử dụng Kubernetes tích hợp sẵn trong Docker Desktop.
- Nắm bắt các công cụ orchestration phổ biến khác trên thị trường (Nomad, ECS, Cloud Run).

## 1. Tại sao cần Container Orchestration?

Khi hệ thống của bạn phát triển từ một vài container trên một máy chủ (host) sang hàng chục, hàng trăm container chạy trên nhiều máy chủ khác nhau, việc quản lý thủ công trở thành "cơn ác mộng". 

> 💡 **Tip:** Hãy tưởng tượng bạn là nhạc trưởng (Orchestrator) của một dàn nhạc giao hưởng (hệ thống). Nếu không có bạn điều phối, mỗi nhạc công (container) sẽ chơi theo một nhịp điệu riêng, dẫn đến một mớ âm thanh hỗn độn.

**Container Orchestration giải quyết các vấn đề:**
- **Triển khai (Deployment):** Tự động hóa quá trình đưa container lên các máy chủ khác nhau.
- **Tự động mở rộng (Scaling):** Tăng/giảm số lượng container dựa trên tải thực tế (load).
- **Cân bằng tải (Load Balancing):** Phân phối traffic đồng đều đến các container.
- **Tự phục hồi (Self-healing):** Khởi động lại hoặc thay thế các container bị lỗi mà không cần can thiệp thủ công.
- **Quản lý tài nguyên:** Sắp xếp container vào các node (máy chủ) có đủ tài nguyên (CPU, RAM).

## 2. Docker Swarm - Giải pháp Orchestration Native của Docker

Docker Swarm là công cụ orchestration được tích hợp sẵn ngay trong Docker Engine, giúp bạn biến một nhóm các máy chủ Docker thành một cụm (cluster) thống nhất.

### Kiến trúc cơ bản
- **Manager Nodes:** Chịu trách nhiệm quản lý cluster, duy trì trạng thái hệ thống, lên lịch (scheduling) và phục vụ API của Swarm.
- **Worker Nodes:** Nhiệm vụ duy nhất là nhận lệnh từ Manager và chạy các container.

### Các lệnh Swarm cơ bản

Khởi tạo một Swarm cluster trên máy Manager:
```bash
# Khởi tạo Swarm
docker swarm init --advertise-addr <IP_CỦA_MANAGER>
```

**Kết quả:**
```
Swarm initialized: current node (abc123xyz) is now a manager.

To add a worker to this swarm, run the following command:
    docker swarm join --token SWMTKN-1-... <IP_CỦA_MANAGER>:2377
```

Tham gia vào Swarm với tư cách là Worker (sử dụng lệnh được sinh ra từ `init`):
```bash
docker swarm join --token <TOKEN> <IP_CỦA_MANAGER>:2377
```

### Services và Scaling

Trong Swarm, bạn không quản lý các container riêng lẻ, mà quản lý các **Services**. Một Service định nghĩa trạng thái mong muốn (ví dụ: cần 3 bản sao của image `nginx`).

```bash
# Tạo một service với 3 bản sao (replicas)
docker service create --name web_server --replicas 3 -p 8080:80 nginx:alpine
```

Kiểm tra trạng thái service:
```bash
# Xem danh sách các services
docker service ls

# Xem chi tiết các container (tasks) đang chạy cho service web_server
docker service ps web_server
```

Tự động mở rộng (Scaling):
```bash
# Tăng số lượng bản sao lên 5
docker service scale web_server=5
```

### Rolling Updates

Khi có phiên bản ứng dụng mới, Swarm hỗ trợ cập nhật dần dần (rolling update) để không làm gián đoạn dịch vụ:
```bash
# Cập nhật image cho service mà không downtime
docker service update --image nginx:1.21-alpine web_server
```

### Stack Deploy

Giống như Docker Compose dùng cho 1 máy, **Docker Stack** dùng để triển khai một hệ thống gồm nhiều service lên toàn bộ Swarm cluster thông qua file `docker-compose.yml`.

```bash
# Triển khai hệ thống bằng file compose
docker stack deploy -c docker-compose.yml my_app_stack
```

## 3. Tổng quan về Kubernetes (K8s)

Dù Docker Swarm rất tốt, **Kubernetes (K8s)** (do Google phát triển) mới là "vị vua" thống trị mảng Container Orchestration hiện nay.

### Kiến trúc K8s
- **Master Node (Control Plane):** Bộ não của K8s (bao gồm API Server, Scheduler, Controller Manager, etcd).
- **Worker Node:** Nơi chạy ứng dụng thực tế.

### Các khái niệm cốt lõi (Core Concepts)
- **Pod:** Đơn vị nhỏ nhất trong K8s. Một Pod có thể chứa 1 hoặc nhiều container chia sẻ chung mạng và lưu trữ. Bạn không chạy container trực tiếp, bạn chạy Pod!
- **Deployment:** Quản lý số lượng Pod, hỗ trợ tự động mở rộng và rolling update (tương đương Service của Swarm).
- **Service:** Cung cấp một IP/Domain tĩnh để kết nối đến một nhóm các Pod đang thay đổi liên tục.

> 📌 **Ghi nhớ:** Xu hướng hiện nay là Kubernetes đang thống trị hoàn toàn các hệ thống lớn, Enterprise. Tuy nhiên, Docker Swarm vẫn là một lựa chọn tuyệt vời cho các dự án vừa và nhỏ vì tính đơn giản, dễ học và thiết lập nhanh.

## 4. So sánh Docker Swarm vs Kubernetes

| Tiêu chí | Docker Swarm | Kubernetes (K8s) |
|---|---|---|
| **Đường cong học tập** | Thấp, dễ học, dùng lại cú pháp Docker/Compose | Rất dốc, nhiều khái niệm phức tạp (Pod, Ingress, RBAC...) |
| **Cài đặt & Thiết lập** | Cực kỳ đơn giản (tích hợp sẵn trong Docker) | Phức tạp (cần công cụ như kubeadm, minikube hoặc dùng Cloud) |
| **Quy mô dự án** | Phù hợp dự án nhỏ đến vừa | Phù hợp dự án lớn, Enterprise, Microservices phức tạp |
| **Mở rộng (Scaling)** | Nhanh nhưng tính năng cơ bản | Chậm hơn một chút nhưng tự động hóa cực mạnh (HPA, VPA) |
| **Hệ sinh thái** | Nhỏ gọn, ít công cụ phụ trợ | Khổng lồ, tiêu chuẩn của ngành công nghiệp (CNCF) |

## 5. Kubernetes trên Docker Desktop

Docker Desktop hỗ trợ chạy một cụm Kubernetes 1-node (Master và Worker trên cùng 1 máy) để bạn tiện học tập và kiểm thử.

**Cách bật:**
1. Mở cài đặt (Settings) của Docker Desktop.
2. Chọn tab **Kubernetes**.
3. Tích chọn **"Enable Kubernetes"** và nhấn "Apply & Restart".
4. Chờ vài phút để Docker tải image và khởi động K8s.

Kiểm tra sau khi bật thành công (trên Terminal):
```bash
# Xem thông tin cluster K8s đang kết nối
kubectl cluster-info

# Xem danh sách node (sẽ thấy 1 node là docker-desktop)
kubectl get nodes
```

## 6. Các giải pháp Orchestration khác

Ngoài Swarm và K8s, thị trường còn có:
- **Nomad (của HashiCorp):** Nhẹ nhàng, đa dụng. Không chỉ orchestrate container mà còn orchestrate các ứng dụng non-container (Java, file thực thi).
- **AWS ECS (Elastic Container Service):** Dịch vụ orchestration độc quyền của Amazon, dễ dùng nếu bạn chỉ chạy trên AWS.
- **Google Cloud Run / AWS Fargate:** Mô hình Serverless Container. Bạn chỉ cần đưa image, cấu hình tài nguyên, hệ thống sẽ tự chạy và tính tiền theo từng millisecond sử dụng mà không cần quản lý node.

## Bài tập
1. **Khởi tạo Swarm**: Biến máy tính hiện tại của bạn (đang chạy Docker) thành một Swarm Manager node bằng lệnh `docker swarm init`. Ghi lại token kết nối để sau này có thể dùng.
2. **Triển khai Service**: Sử dụng `docker service create` để chạy một service có tên là `my_nginx`, sử dụng image `nginx:latest`, publish port `8080` của host vào port `80` của container và thiết lập `--replicas 2`.
3. **Scale Service**: Tiếp tục với service `my_nginx` ở trên, hãy chạy lệnh để tăng số lượng bản sao (replicas) lên thành 4. Dùng lệnh `docker service ls` và `docker service ps my_nginx` để quan sát số lượng container đang chạy.
4. **Dọn dẹp**: Xóa service `my_nginx` (bằng lệnh `docker service rm`), sau đó đưa máy hiện tại rời khỏi Swarm cluster (bằng lệnh `docker swarm leave --force`).

## Tiếp theo
→ [Bài 21: Dự án Capstone](./21_capstone_project.md)
