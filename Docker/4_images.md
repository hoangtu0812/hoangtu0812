# Bài 4: Làm việc với Image

## Mục tiêu
- Hiểu rõ Docker Image là gì, cấu trúc phân lớp (layered filesystem) của nó.
- Tải, liệt kê và xem thông tin chi tiết của các image.
- Quản lý thẻ (tag) cho image và đẩy (push) image lên Docker Hub.
- Tìm hiểu lịch sử hình thành (history) của một image.
- Xóa và dọn dẹp các image không còn sử dụng để giải phóng ổ cứng.

---

## 1. Docker Image và Kiến trúc phân lớp (Layers)

### Docker Image là gì?
Hãy tưởng tượng Docker Image như một **bản thiết kế (blueprint)** hoặc một chiếc **khuôn đúc bánh**. Từ một cái khuôn (Image), bạn có thể đúc ra hàng nghìn chiếc bánh (Container) giống hệt nhau. Docker Image là một mẫu *chỉ-đọc (read-only)* chứa mọi thứ cần thiết để chạy ứng dụng: mã nguồn, thư viện, biến môi trường, và các tệp cấu hình.

### Kiến trúc phân lớp (Layered Filesystem)
Khác với các file nén thông thường, Docker Image được tạo thành từ nhiều **lớp (layers)** xếp chồng lên nhau. Mỗi khi bạn thêm một câu lệnh cài đặt hay copy một file vào image, Docker tạo ra một layer mới.

> 💡 **Tip:** Các image khác nhau có thể **chia sẻ chung các layer**. Ví dụ, nếu bạn có 5 image cùng xây dựng từ nền tảng Ubuntu, lớp Ubuntu gốc sẽ chỉ được lưu một lần trên ổ cứng của bạn. Điều này giúp tối ưu hóa không gian lưu trữ và tốc độ tải mạng cực kỳ hiệu quả!

---

## 2. Tìm kiếm và Tải Image từ Registry

### Lệnh `docker pull`
Để lấy một image từ kho chứa (registry, thường là Docker Hub), chúng ta sử dụng lệnh `docker pull`.

```bash
# Tải image nginx phiên bản alpine siêu nhẹ
docker pull nginx:alpine
```

**Kết quả:**
```
alpine: Pulling from library/nginx
a0d0a0d46f8b: Pull complete 
c7a9772653fc: Pull complete 
d88f6ab63a0a: Pull complete 
Digest: sha256:d55ff771a396e9da8912dc1ea354c4f34606ea2b8812c3f761dcce1fbf41e053
Status: Downloaded newer image for nginx:alpine
docker.io/library/nginx:alpine
```
*(Lưu ý: Mỗi dòng "Pull complete" đại diện cho một layer đang được tải về.)*

---

## 3. Liệt kê và Xem chi tiết Image

### Lệnh `docker images` / `docker image ls`
Sau khi tải, bạn cần kiểm tra xem image đã nằm trên máy chưa. Bạn có thể dùng một trong hai lệnh dưới đây, chúng hoàn toàn tương đương nhau.

```bash
# Liệt kê tất cả các image đang có trong máy
docker image ls
```

**Kết quả:**
```
REPOSITORY   TAG       IMAGE ID       CREATED        SIZE
nginx        alpine    49a2a9e22754   3 days ago     42.6MB
ubuntu       latest    ba6acccedd29   2 weeks ago    72.8MB
```

### Lệnh `docker image inspect`
Nếu bạn muốn xem toàn bộ "nội tạng" của image (như các biến môi trường cấu hình sẵn, hệ điều hành, hay cấu trúc thư mục), hãy dùng lệnh `inspect`.

```bash
# Xem chi tiết cấu hình của image
docker image inspect nginx:alpine
```

**Kết quả:**
```json
[
    {
        "Id": "sha256:49a2a9e22754...",
        "RepoTags": [
            "nginx:alpine"
        ],
        "Architecture": "amd64",
        "Os": "linux",
        ...
    }
]
```

---

## 4. Quản lý Tag và Lịch sử Image

### Tag là gì?
Tag (thẻ) đóng vai trò như số phiên bản. Định dạng tiêu chuẩn là `tên-image:tag`. Nếu bạn không chỉ định tag, Docker sẽ tự động dùng tag mặc định là `latest`. 

> ⚠️ **Lưu ý:** Tag `latest` không có nghĩa là phiên bản mới nhất về mặt kỹ thuật, nó chỉ là một cái tên mặc định. Hãy luôn chỉ định rõ phiên bản (ví dụ: `1.21.0`) trên các hệ thống thực tế (Production) để tránh những bất ngờ không đáng có.

### Lệnh `docker tag`
Dùng để tạo một tên hoặc tag mới trỏ đến cùng một image ID gốc.

```bash
# Gắn thêm một tag 'v1' cho image nginx:alpine hiện tại
docker tag nginx:alpine my-nginx:v1
```

*(Lệnh này không tạo ra image mới tốn thêm dung lượng, nó chỉ tạo một "nhãn dán" mới lên image cũ).*

### Lệnh `docker history`
Lệnh này giúp bạn xem image đã được xây dựng như thế nào qua từng layer.

```bash
# Xem các bước đã tạo nên image
docker history nginx:alpine
```

**Kết quả:**
```
IMAGE          CREATED        CREATED BY                                      SIZE      COMMENT
49a2a9e22754   3 days ago     CMD ["nginx" "-g" "daemon off;"]                0B        buildkit.dockerfile.v0
<missing>      3 days ago     EXPOSE map[80/tcp:{}]                           0B        buildkit.dockerfile.v0
<missing>      3 days ago     COPY /usr/share/nginx/html /usr/share/ngin…   1.2MB     buildkit.dockerfile.v0
...
```

---

## 5. Đẩy Image lên Registry (Push)

Sau khi tạo tag, bạn có thể đẩy image của mình lên một Registry (như Docker Hub) để chia sẻ với người khác hoặc để triển khai (deploy). Để đẩy lên tài khoản Docker Hub của bạn, tag của image phải bắt đầu bằng `tên-tài-khoản/`.

```bash
# Giả sử tài khoản Docker Hub của bạn là hoangtu0812
docker tag nginx:alpine hoangtu0812/my-nginx:v1

# Đăng nhập vào Docker Hub (chỉ cần làm một lần)
docker login

# Đẩy image lên mạng
docker push hoangtu0812/my-nginx:v1
```

**Kết quả:**
```
The push refers to repository [docker.io/hoangtu0812/my-nginx]
e5e4a8b7c3d1: Pushed 
d1a9b2c3d4e5: Pushed 
v1: digest: sha256:abcd... size: 1234
```

---

## 6. Xóa và Dọn dẹp Image

### Lệnh `docker rmi` / `docker image rm`
Bạn dùng lệnh này để xóa đi các image không còn sử dụng. 

```bash
# Xóa image thông qua tên và tag
docker rmi my-nginx:v1
```

**Kết quả:**
```
Untagged: my-nginx:v1
```

> 📌 **Ghi nhớ:** Khi bạn dùng `docker rmi` xóa một image có nhiều tag, Docker ban đầu chỉ xóa cái tag (Untagged). Nó chỉ thực sự xóa dữ liệu ổ cứng (Deleted) khi chiếc tag cuối cùng của Image ID đó bị gỡ bỏ, và không có container nào đang chạy dùng image đó.

### Lệnh `docker image prune`
Khi bạn build image nhiều lần, sẽ xuất hiện các image không có tên và không có tag (được hiển thị là `<none>:<none>`), gọi là các *dangling images*. Để dọn dẹp tất cả chúng một cách an toàn và giải phóng dung lượng, hãy dùng prune.

```bash
# Xóa tất cả các image 'mồ côi' (dangling)
docker image prune
```

**Kết quả:**
```
WARNING! This will remove all dangling images.
Are you sure you want to continue? [y/N] y
Deleted Images:
deleted: sha256:f8b...
Total reclaimed space: 152MB
```

---

## Bài tập

1. **Bài 1:** Bạn hãy tải về phiên bản chính thức mới nhất của hệ điều hành `ubuntu`. Sau khi tải xong, hãy dùng lệnh liệt kê để xác minh image đó đã có trên máy tính của bạn và ghi nhận dung lượng (size) của nó.

2. **Bài 2:** Chọn image `ubuntu` bạn vừa tải. Sử dụng lệnh `docker tag` để tạo cho nó một cái tên mới: `he-dieu-hanh:test`. Sau đó, hãy dùng `docker history` để xem danh sách các layer bên trong image này.

3. **Bài 3:** Thử nghiệm tạo lỗi. Bạn hãy chạy một container ngầm từ image `ubuntu` bằng lệnh: `docker run -d ubuntu sleep 1000`. Ngay khi container đang chạy, bạn hãy thử xóa thẳng image ubuntu bằng lệnh `docker rmi ubuntu`. Bạn nhận được thông báo lỗi gì? (Lỗi này rất quan trọng để hiểu cách Docker bảo vệ các dữ liệu đang hoạt động).

4. **Bài 4:** Bạn hãy đóng giả quá trình tải ứng dụng lên môi trường Production. Hãy đổi tên một image bất kỳ trên máy bạn để gắn với tên tài khoản Docker ảo của bạn (ví dụ: `my-user/app:v2`). Liệt kê chi tiết cấu hình JSON của image đó ra màn hình bằng lệnh thích hợp. Cuối cùng, thực hiện dọn dẹp: Hãy xóa sạch tất cả các tag bạn vừa tạo cho bài thực hành này.


---

## Tiếp theo
→ [Bài 5: Xây dựng Image với Dockerfile](./5_dockerfile_basics.md)
