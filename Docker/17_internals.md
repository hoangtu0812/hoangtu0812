# Bài 17: Docker Internals — Docker hoạt động thế nào bên trong

## Mục tiêu
- Hiểu được các công nghệ cốt lõi của Linux Kernel đứng sau Docker (Namespaces, cgroups, UnionFS).
- Nắm bắt cơ chế quản lý tài nguyên và cô lập môi trường của Container.
- Hiểu kiến trúc Container Runtime (containerd, runc) và chuẩn OCI.
- Biết cách sử dụng `docker inspect` và khám phá cấu trúc lưu trữ nội bộ của Docker.
- Phân biệt rõ ràng và sâu sắc giữa Virtual Machine (VM) và Container.

---

## 1. So sánh: Virtual Machine (VM) vs Container

Trước khi đi sâu vào hệ thống bên trong, hãy làm rõ sự khác biệt lớn nhất giữa máy ảo và container. Đây là nền tảng để hiểu tại sao Docker lại nhẹ và nhanh.

### Virtual Machine (Hypervisor)
Máy ảo (VM) thực hiện ảo hóa ở mức phần cứng. Mỗi VM chạy một hệ điều hành khách (Guest OS) hoàn chỉnh (có kernel riêng) trên đỉnh một phần mềm gọi là Hypervisor (như VMware, VirtualBox, Hyper-V).
- **Ưu điểm:** Mức độ cô lập hoàn toàn, tính bảo mật rất cao do có Kernel riêng biệt.
- **Nhược điểm:** Dung lượng khổng lồ (tính bằng GB), khởi động chậm (vài chục giây đến vài phút), tiêu tốn nhiều tài nguyên thừa cho Guest OS.

### Container (Kernel Sharing)
Container thực hiện ảo hóa ở mức hệ điều hành. Các container chia sẻ chung kernel của hệ điều hành máy chủ (Host OS) và không cần Guest OS.
- **Ưu điểm:** Kích thước siêu nhẹ (tính bằng MB), khởi động trong tích tắc (tính bằng milliseconds), tận dụng tối đa tài nguyên máy chủ.
- **Nhược điểm:** Mức độ cô lập không tuyệt đối như VM vì dùng chung kernel. Một lỗi nghiêm trọng ở mức kernel có thể làm sập toàn bộ các container trên host.

> 💡 **Tip:** Nếu Virtual Machine giống như xây nhiều ngôi nhà riêng biệt, mỗi nhà có móng, hệ thống điện nước (Kernel) riêng; thì Container giống như các căn hộ trong một tòa nhà chung cư, nơi mọi căn hộ đều chia sẻ chung một hệ thống điện nước tòa nhà nhưng không gian sinh hoạt thì độc lập.

---

## 2. Linux Namespaces — Cô lập môi trường

Làm sao để các container chia sẻ chung Kernel nhưng không "nhìn thấy" hay can thiệp vào tiến trình của nhau? Docker sử dụng một tính năng tuyệt vời của Linux kernel gọi là **Namespaces**. Namespaces tạo ra ảo giác rằng mỗi container có một hệ thống máy tính hoàn toàn độc lập.

Dưới đây là 6 loại Namespaces cốt lõi:
- **PID (Process ID) namespace:** Cô lập không gian ID của các tiến trình. Trong container, tiến trình chạy ứng dụng chính của bạn có thể mang PID = 1. Nhưng nhìn từ ngoài Host OS, nó chỉ là một tiến trình bình thường với PID = 3456.
- **Network namespace:** Cung cấp cho mỗi container một stack mạng riêng (địa chỉ IP độc lập, bảng định tuyến riêng, các interface riêng như `eth0`).
- **Mount namespace:** Cô lập hệ thống tập tin. Container có root filesystem `/` riêng biệt mà không thấy được ổ cứng của Host OS trừ khi được mount.
- **UTS (UNIX Timesharing System) namespace:** Cho phép mỗi container có một hostname và domain name riêng biệt.
- **IPC (Inter-Process Communication) namespace:** Cô lập các cơ chế giao tiếp giữa các tiến trình (như shared memory). Các container không thể đọc bộ nhớ chung của nhau.
- **User namespace:** Ánh xạ user/group IDs. Một user có quyền `root` bên trong container có thể được cấu hình chỉ là một user bình thường (unprivileged) trên Host OS để tăng cường bảo mật.

---

## 3. Control Groups (cgroups) — Giới hạn tài nguyên

Nếu Namespaces làm nhiệm vụ "che mắt" (cô lập những gì container nhìn thấy), thì **cgroups** làm nhiệm vụ "trói buộc" (giới hạn những gì container có thể sử dụng). cgroups giám sát và giới hạn lượng tài nguyên phần cứng mà một tiến trình hoặc container được phép tiêu thụ.

- **CPU:** Giới hạn lượng CPU cycle, thời gian CPU, hoặc số lượng core mà container có thể sử dụng (liên hệ với flag `--cpus`).
- **Memory:** Giới hạn RAM (và Swap) tối đa để tránh hiện tượng một container rò rỉ bộ nhớ ngốn hết RAM của máy chủ (liên hệ với flag `--memory`).
- **I/O (Block I/O):** Giới hạn tốc độ đọc/ghi (Read/Write bandwidth hoặc IOPS) vào ổ cứng.
- **Network:** Giới hạn và ưu tiên băng thông mạng.

```bash
# Chạy một container Nginx bị giới hạn ở mức tối đa nửa core CPU và 256MB RAM
docker run -d --name limited-nginx --cpus="0.5" --memory="256m" nginx
```

**Kết quả:**
```
8a9b7c6d5e4f3a2b1c...
```
> ⚠️ **Lưu ý:** Nếu tiến trình trong container `limited-nginx` cố gắng vượt quá 256MB RAM, hệ điều hành (thông qua OOM-Killer) sẽ ngay lập tức "giết" tiến trình đó để bảo vệ Host OS.

---

## 4. Union Filesystem (OverlayFS) và Copy-on-Write

Docker images được tạo thành từ nhiều layer (lớp) xếp chồng lên nhau thông qua công nghệ **Union Filesystem** (hiện tại Docker mặc định sử dụng OverlayFS: `overlay2`).

### Kiến trúc xếp lớp (Layered Architecture)
- **Image layers (Read-Only):** Mỗi lệnh `RUN`, `COPY`, `ADD` trong Dockerfile tạo ra một layer chỉ đọc. Khi nhiều container chạy từ cùng một image, hoặc nhiều image chia sẻ chung một base image (vd: `ubuntu`), chúng dùng chung các layer này. Điều này tiết kiệm rất nhiều dung lượng đĩa cứng.
- **Container layer (Read-Write):** Khi một container khởi tạo từ image, Docker gắn thêm một lớp mỏng gọi là Container layer trên cùng, và lớp này có quyền Đọc-Ghi.

### Cơ chế Copy-on-Write (CoW)
Điều gì xảy ra khi container muốn sửa đổi một file có sẵn trong image?
- Docker áp dụng cơ chế Copy-on-Write. Nó sẽ copy bản gốc của file đó từ Image layer bên dưới lên Container layer (Read-Write) ở trên cùng.
- Mọi sửa đổi sau đó chỉ diễn ra trên bản copy ở Container layer. File gốc ở lớp Read-Only bên dưới vẫn giữ nguyên vẹn để các container khác có thể tiếp tục sử dụng.

```
       +-----------------------+
       |   Container Layer     |  <-- Read-Write (New files, modifications)
       |     (thin layer)      |
       +-----------------------+
       |      Image Layer      |  <-- Read-Only (e.g., ADD app.py)
       +-----------------------+
       |      Image Layer      |  <-- Read-Only (e.g., RUN pip install)
       +-----------------------+
       |   Base Image Layer    |  <-- Read-Only (e.g., Ubuntu 20.04)
       +-----------------------+
```

---

## 5. Kiến trúc Container Runtime (containerd, runc)

Những ngày đầu, toàn bộ logic cốt lõi của Docker là một khối duy nhất (Monolithic). Để hệ sinh thái phát triển, Docker đã tách nhỏ kiến trúc theo chuẩn **OCI (Open Container Initiative)**.

- **Docker Engine / dockerd:** Quản lý toàn bộ API, Volumes, Network, Swarm. Khi bạn gõ lệnh `docker run`, CLI gọi API của `dockerd`.
- **containerd:** Daemon chịu trách nhiệm quản lý vòng đời của container (tải images, tạo, khởi động, dừng). Nó đóng vai trò trung gian.
- **runc:** Là OCI runtime cấp thấp (low-level). `containerd` gọi `runc` để làm một việc duy nhất: tương tác với Linux Kernel (Namespaces, cgroups) để tạo container, rồi nó thoát ra.

```
Docker CLI  --->  dockerd (API)  --->  containerd  --->  runc  --->  (Linux Kernel)
```

---

## 6. Khám phá `/var/lib/docker/` và `docker inspect`

### Cấu trúc thư mục dữ liệu
Trên Linux, mọi dữ liệu nội bộ của Docker nằm ở `/var/lib/docker/`. 
```bash
# Yêu cầu quyền root
sudo ls -l /var/lib/docker/
```

- `containers/`: Lưu trữ config (JSON) và log files của từng container.
- `image/`: Metadata mô tả cấu trúc của các image.
- `overlay2/`: Nơi lưu trữ dữ liệu thực sự của các image layers và container layers.
- `volumes/`: Dữ liệu bền vững của các Docker volumes.

> 📌 **Ghi nhớ:** Không bao giờ dùng lệnh như `rm`, `cp`, `vi` trực tiếp vào trong thư mục này. Luôn sử dụng lệnh CLI của Docker (như `docker system prune`, `docker volume rm`) để dọn dẹp và quản lý.

### docker inspect
Lệnh `docker inspect` cung cấp mọi thông tin metadata chi tiết về một container, bao gồm IP, Mounts, Layer, cgroups...

```bash
# Xem thông tin của container
docker inspect --format '{{.State.Pid}}' limited-nginx
```

**Kết quả:**
```
12345
```
Kết quả trả về chính là PID thực tế của tiến trình Nginx đang chạy trên máy Host.

---

## Bài tập

1. **Khám phá cấu trúc Layer**:
   Tải về image `nginx`: `docker pull nginx`.
   Chạy lệnh `docker inspect nginx` và cuộn xuống tìm đến phần `RootFS.Layers`. Hãy đếm xem image Nginx này được tạo thành từ bao nhiêu layer read-only.
2. **Khám phá PID Namespace**:
   Chạy container: `docker run -d --name my-web nginx`.
   Dùng lệnh `docker top my-web` để xem các PID của Nginx từ góc nhìn của Docker. Sau đó, dùng `docker inspect --format '{{.State.Pid}}' my-web` để tìm PID thực của nó trên Host OS.

3. **Thử nghiệm cgroups (Giới hạn CPU)**:
   Chạy container ubuntu với giới hạn tài nguyên: `docker run -it --rm --cpus="0.5" ubuntu /bin/bash`.
   Bên trong container, cài công cụ test: `apt update && apt install -y stress`.
   Chạy lệnh `stress --cpu 1` để ép CPU làm việc 100%.
   Mở một terminal khác trên máy host và chạy `docker stats`. Bạn sẽ thấy tiến trình trong container chỉ có thể tiêu thụ tối đa khoảng 50% CPU do bị khống chế bởi cgroups. Nhấn Ctrl+C để dừng.

4. **Theo dõi Copy-on-Write (CoW)**:
   Chạy container: `docker run -it --name cow-test ubuntu /bin/bash`. 
   Bên trong container, tạo một file mới: `echo "Hello từ Container" > /root/hello.txt` rồi thoát ra (`exit`).
   Chạy lệnh `docker inspect cow-test` và tìm phần `GraphDriver.Data.UpperDir`. Đường dẫn này dẫn đến thư mục vật lý trên Host OS.
   Hãy dùng quyền root/sudo `sudo cat <đường-dẫn-UpperDir>/root/hello.txt`. Bạn sẽ thấy nội dung file bạn vừa tạo! Đây chính là bằng chứng sống động của lớp Container layer (Read-Write).

---

## Tiếp theo
→ [Bài 18: CI/CD với Docker](./18_cicd.md)
