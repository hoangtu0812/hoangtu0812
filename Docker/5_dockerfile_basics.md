# Bài 5: Dockerfile cơ bản

## Mục tiêu
- Hiểu được Dockerfile là gì và vai trò của nó trong quá trình tạo Docker Image.
- Nắm vững cú pháp cơ bản và các từ khóa (instruction) thông dụng nhất.
- Biết cách viết một Dockerfile hoàn chỉnh cho một ứng dụng Python đơn giản.
- Hiểu cách sử dụng lệnh `docker build` và khái niệm Build Context.
- Nắm được cách tối ưu hóa Dockerfile dựa trên cơ chế Layer Cache.

---

## 1. Dockerfile là gì?

Nếu Docker Image là một "chiếc bánh" đã được nướng chín, thì **Dockerfile** chính là "công thức nấu ăn". Nó là một file văn bản thuần túy chứa một chuỗi các lệnh (instructions) mà Docker sẽ đọc từ trên xuống dưới để tự động hóa quá trình đóng gói ứng dụng của bạn thành một Docker Image.

> 💡 **Tip:** Sử dụng Dockerfile giúp quá trình xây dựng môi trường trở nên minh bạch, có thể chia sẻ được và dễ dàng quản lý phiên bản (version control) bằng Git.

---

## 2. Cú pháp cơ bản và các Instruction quan trọng

Một file Dockerfile luôn bắt đầu bằng lệnh chỉ định Image nền (Base Image). Dưới đây là các lệnh (instructions) bạn sẽ gặp trong hầu hết mọi Dockerfile.

### 2.1. Lệnh FROM
Xác định Base Image để bắt đầu. Mọi Dockerfile hợp lệ đều phải bắt đầu bằng lệnh `FROM`.
```dockerfile
# Sử dụng Python 3.9 bản gọn nhẹ (alpine)
FROM python:3.9-alpine
```

### 2.2. Lệnh LABEL
Thêm metadata (thông tin mô tả) vào image, ví dụ như người tạo, phiên bản, hoặc mô tả.
```dockerfile
LABEL maintainer="hoangtu@example.com"
LABEL version="1.0"
```

### 2.3. Lệnh WORKDIR
Thiết lập thư mục làm việc mặc định cho các lệnh `RUN`, `CMD`, `COPY`, v.v. phía sau nó. Nếu thư mục chưa tồn tại, Docker sẽ tự động tạo.
```dockerfile
# Đặt thư mục làm việc là /app trong container
WORKDIR /app
```

### 2.4. Lệnh COPY và ADD
Cả hai đều dùng để chép file từ máy host (máy của bạn) vào trong image.
- **`COPY`**: Đơn giản là copy file/thư mục. (Được khuyên dùng).
- **`ADD`**: Tương tự `COPY` nhưng có thêm tính năng giải nén file `.tar` hoặc tải file từ một URL.

```dockerfile
# Copy file requirements.txt vào /app
COPY requirements.txt .

# Copy toàn bộ mã nguồn vào /app
COPY . .
```

### 2.5. Lệnh RUN
Thực thi các lệnh trên terminal (bash/sh) trong quá trình **build** Image. Thường dùng để cài đặt thư viện hoặc gói phần mềm.
```dockerfile
RUN pip install -r requirements.txt
```

### 2.6. Lệnh EXPOSE
Khai báo cho Docker biết container sẽ lắng nghe trên cổng (port) nào khi chạy. Lệnh này mang tính chất tài liệu là chính, nó không thực sự "mở" cổng trên máy host.
```dockerfile
EXPOSE 8080
```

### 2.7. Lệnh CMD
Xác định lệnh mặc định sẽ được thực thi khi **chạy Container** (chứ không phải khi build). Chỉ có thể có một lệnh `CMD` có hiệu lực trong Dockerfile (lệnh cuối cùng).
```dockerfile
CMD ["python", "app.py"]
```

> ⚠️ **Lưu ý:** Phân biệt `RUN` và `CMD`:
> - `RUN`: Chạy lúc build image (ví dụ: cài dependencies). Tạo ra một layer mới.
> - `CMD`: Chạy lúc khởi động container (ví dụ: chạy server ứng dụng).

---

## 3. Ví dụ thực tế: Ứng dụng Python đơn giản

Hãy xem một Dockerfile hoàn chỉnh để đóng gói một ứng dụng web Python.

```dockerfile
# 1. Chọn base image
FROM python:3.9-slim

# 2. Thêm thông tin người viết
LABEL author="Hoang Tu"

# 3. Chuyển đến thư mục làm việc
WORKDIR /app

# 4. Copy file requirements trước (để tối ưu cache)
COPY requirements.txt .

# 5. Cài đặt các thư viện cần thiết
RUN pip install --no-cache-dir -r requirements.txt

# 6. Copy toàn bộ mã nguồn còn lại vào container
COPY . .

# 7. Khai báo port ứng dụng sử dụng
EXPOSE 5000

# 8. Lệnh chạy ứng dụng khi container khởi động
CMD ["python", "main.py"]
```

---

## 4. Lệnh `docker build`

Để biến Dockerfile thành Image, bạn sử dụng lệnh `docker build` trên terminal.

```bash
docker build -t my-python-app:1.0 .
```

**Giải thích các cờ (flags):**
- `-t my-python-app:1.0`: Gắn tên (tag) cho image là `my-python-app` với phiên bản `1.0`.
- `.`: Ký tự dấu chấm đại diện cho **Build Context** (thư mục hiện tại). Docker sẽ gửi toàn bộ file trong thư mục này cho Docker Daemon để tiến hành build.

Nếu file Dockerfile của bạn có tên khác (ví dụ `Dockerfile.dev`), bạn có thể dùng cờ `-f`:
```bash
docker build -f Dockerfile.dev -t my-python-app:dev .
```

---

## 5. Thứ tự instruction ảnh hưởng tới Cache

Docker sử dụng cơ chế **Layer Cache** để tăng tốc quá trình build. Mỗi instruction (`RUN`, `COPY`, `ADD`) tạo ra một "tầng" (layer) mới.
Nếu một instruction không thay đổi, Docker sẽ dùng lại cache của layer đó thay vì chạy lại. **Nhưng nếu một layer bị thay đổi, tất cả các layer phía sau nó sẽ phải build lại từ đầu.**

> 📌 **Ghi nhớ:** Hãy đặt những instruction ít thay đổi lên trên cùng, và những instruction hay thay đổi (như copy mã nguồn) xuống dưới cùng.

Đó là lý do ở ví dụ trên, chúng ta `COPY requirements.txt` và `RUN pip install` trước khi `COPY . .` (copy toàn bộ mã nguồn). Mã nguồn thay đổi liên tục, nhưng các thư viện (`requirements.txt`) thì ít thay đổi hơn. Việc sắp xếp này giúp bạn không phải tải lại các thư viện mỗi khi sửa một dòng code nhỏ.

---

## Bài tập

1. **Cơ bản (Bài tập 1)**: Khởi tạo một thư mục trống. Tạo một file `Dockerfile` sử dụng base image `nginx:alpine`. Thêm lệnh `COPY` để đưa một file `index.html` đơn giản từ máy bạn vào thư mục `/usr/share/nginx/html/` của image.

2. **Bài 2:** Build image từ Dockerfile ở bài 1 với tên `my-nginx-web` và chạy một container từ image này, map port 8080 của máy host vào port 80 của container. Kiểm tra trên trình duyệt.

3. **Bài 3:** Tạo một Dockerfile cho một ứng dụng Node.js giả lập. Hãy đảm bảo bạn cấu trúc Dockerfile sao cho việc copy file `package.json` và chạy `npm install` được thực hiện trước khi copy toàn bộ mã nguồn để tận dụng tối đa layer cache.

4. **Nâng cao (Bài tập 4)**: Viết một Dockerfile sử dụng `ubuntu:20.04`. Cập nhật apt-get, cài đặt `curl`, và đặt biến môi trường `APP_ENV=production`. Thay vì dùng lệnh `CMD` để chạy một file cụ thể, hãy dùng lệnh `CMD ["bash"]` và dùng `docker run -it` để vào thẳng giao diện dòng lệnh của container này, kiểm tra xem curl đã được cài chưa.


---

## Tiếp theo
→ [Bài 6: Tối ưu Build Context với .dockerignore](./6_dockerignore.md)
