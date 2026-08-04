# Lộ Trình Học Go (Golang) Chi Tiết

>**Mỗi bài đã có file chi tiết riêng** (lý thuyết + code mẫu + bài tập) — bấm vào tiêu đề từng bài để mở. Dự án capstone nằm ở [21_capstone_project.md](./21_capstone_project.md).

---

## Giai đoạn 0 — Giới thiệu & Cài đặt

### Bài 0: Giới thiệu về Go
- Go là gì, ra đời để giải quyết vấn đề gì (Google, tốc độ biên dịch, concurrency, static binary).
- So sánh nhanh với ngôn ngữ bạn đã biết (Python/JS/ABAP): typing tĩnh, không OOP kiểu class, error là giá trị chứ không phải exception.
- Hệ sinh thái: `go` CLI, `GOPATH` vs Go Modules, thư viện chuẩn mạnh (`net/http`, `encoding/json`...).

### Bài 1: Cài đặt & Hello World *(đã có — [1_get_started.md](./1_get_started.md))*
- Cài Go, VS Code + extension Go, `go mod init`, `go run`, `go build`.
- **Bài tập:** Sửa Hello World in ra tên bạn; thử `go build` rồi chạy file `.exe` tạo ra; thử `go fmt` và `go vet`.

---

## Giai đoạn 1 — Go Cơ Bản (Fundamentals)

### [Bài 2: Biến, kiểu dữ liệu, hằng số](./2_variables_types.md)
- `var`, `:=`, kiểu cơ bản (int, float64, string, bool, rune, byte), zero value, `const`, `iota`.
- **Bài tập:**
  1. Khai báo các biến bằng cả `var` và `:=`, in ra kiểu bằng `%T`.
  2. Viết chương trình đổi độ C ↔ F.
  3. Dùng `iota` để tạo enum ngày trong tuần.

### [Bài 3: Toán tử & luồng điều khiển](./3_control_flow.md)
- Toán tử số học/so sánh/logic, `if/else`, `switch` (kể cả switch không điều kiện), `for` (3 dạng, kể cả vô hạn), `break`/`continue`/`goto`.
- **Bài tập:**
  1. FizzBuzz 1–100.
  2. Kiểm tra số nguyên tố.
  3. In bảng cửu chương bằng vòng lặp lồng nhau.

### [Bài 4: Hàm (Functions)](./4_functions.md)
- Khai báo hàm, nhiều giá trị trả về, named return, biến đối số (`...int`), hàm là first-class citizen (closures), hàm ẩn danh, `defer`.
- **Bài tập:**
  1. Hàm `divide(a, b int) (int, error)` trả lỗi khi chia cho 0.
  2. Hàm `sum(nums ...int) int`.
  3. Viết closure đếm số lần gọi (counter generator).

### [Bài 5: Mảng, Slice, Map](./5_slices_maps.md)
- Array cố định vs slice động, `append`, `make`, `len`/`cap`, slice của slice (chia sẻ underlying array), map: khai báo, CRUD, kiểm tra key tồn tại (`v, ok := m[k]`).
- **Bài tập:**
  1. Viết hàm loại bỏ phần tử trùng trong slice.
  2. Đếm tần suất từ trong một câu bằng `map[string]int`.
  3. Viết hàm đảo ngược slice tại chỗ (in-place).

### [Bài 6: Struct & Method](./6_structs_methods.md)
- Định nghĩa struct, embedding (nhúng struct để "kế thừa"), method với receiver (value vs pointer receiver), so sánh struct.
- **Bài tập:**
  1. Struct `Rectangle` với method `Area()`, `Perimeter()`.
  2. Struct `Animal` nhúng vào `Dog`, `Cat` để tái sử dụng field/method.
  3. Giải thích bằng code khi nào nên dùng pointer receiver.

### [Bài 7: Con trỏ (Pointers)](./7_pointers.md)
- `&`, `*`, truyền giá trị vs truyền địa chỉ, khi nào Go tự động dereference (method call).
- **Bài tập:**
  1. Viết hàm `swap(a, b *int)` hoán đổi 2 số.
  2. Viết hàm sửa field của struct thông qua pointer, so sánh kết quả nếu không dùng pointer.

### [Bài 8: Xử lý lỗi (Error Handling)](./8_error_handling.md)
- Interface `error`, `errors.New`, `fmt.Errorf` với `%w`, `errors.Is`/`errors.As`, custom error type, `panic`/`recover` (và khi nào KHÔNG nên dùng).
- **Bài tập:**
  1. Viết custom error `ValidationError` chứa tên field lỗi.
  2. Dùng `errors.Is` để phân biệt lỗi "not found" với lỗi khác.
  3. Viết hàm dùng `recover` để chương trình không crash khi chia cho 0.

### [Bài 9: Package & Go Modules](./9_packages_modules.md)
- Cấu trúc package, export (chữ hoa) vs unexported, `go mod init/tidy`, semantic import, tổ chức nhiều file trong 1 package.
- **Bài tập:** Tách chương trình Bài 6 (Rectangle) thành package riêng `shapes`, import và dùng từ `main`.

---

## Giai đoạn 2 — Go Trung Cấp (Intermediate)

### [Bài 10: Interface](./10_interfaces.md)
- Interface implicit (không cần khai báo "implements"), interface rỗng `interface{}`/`any`, type assertion, type switch.
- **Bài tập:**
  1. Interface `Shape` với `Area()`, cho `Circle`, `Rectangle` cùng implement, viết hàm `PrintArea(s Shape)`.
  2. Viết hàm nhận `interface{}` và dùng type switch để xử lý nhiều kiểu khác nhau.

### [Bài 11: Goroutines & Channels](./11_goroutines_channels.md)
- `go func()`, channel (`chan`), buffered vs unbuffered, `select`, `sync.WaitGroup`, `sync.Mutex`, tránh race condition (`go run -race`).
- **Bài tập:**
  1. Chạy 5 goroutine in số 1–5 song song, dùng `WaitGroup` để chờ.
  2. Producer-consumer đơn giản bằng channel.
  3. Viết bộ đếm an toàn với `sync.Mutex`, kiểm tra bằng `-race`.

### [Bài 12: Context](./12_context.md)
- `context.Context`: `WithCancel`, `WithTimeout`, `WithValue`, dùng để hủy goroutine/luồng xử lý request.
- **Bài tập:** Viết hàm giả lập gọi API mất 3s nhưng có `context.WithTimeout(1s)` để hủy sớm.

### [Bài 13: Xử lý JSON & File I/O](./13_json_file_io.md)
- `encoding/json` (Marshal/Unmarshal, struct tag), đọc/ghi file với `os`, `bufio`, `io`.
- **Bài tập:**
  1. Định nghĩa struct `User` với json tag, marshal/unmarshal qua lại.
  2. Đọc file `.txt` theo từng dòng và đếm số dòng.
  3. Ghi danh sách `User` ra file `users.json`.

### [Bài 14: Testing trong Go](./14_testing.md)
- Package `testing`, `go test`, table-driven test, `testify` (assert/require), benchmark cơ bản, test coverage.
- **Bài tập:**
  1. Viết unit test table-driven cho hàm `divide` (Bài 4).
  2. Chạy `go test -cover` và đạt >80% coverage cho package `shapes`.

### [Bài 15: Generics (Go 1.18+)](./15_generics.md)
- Type parameter, constraint (`any`, `comparable`, constraint tự định nghĩa), khi nào nên/không nên dùng generics.
- **Bài tập:** Viết hàm generic `Map[T, U any](s []T, f func(T) U) []U` và `Filter`.

---

## Giai đoạn 3 — Go Nâng Cao (Advanced)

### [Bài 16: Thiết kế package & clean architecture cơ bản](./16_clean_architecture.md)
- Nguyên tắc tổ chức code: tách theo layer (handler → service → repository), dependency injection thủ công qua interface, tránh import cycle.

### [Bài 17: Làm việc với Database](./17_database.md)
- `database/sql`, driver (vd. `pgx`/`mysql`), connection pool, hoặc dùng ORM (GORM), migration (golang-migrate).
- **Bài tập:** Kết nối tới PostgreSQL/MySQL, viết CRUD cho bảng `users`.

### [Bài 18: HTTP Server & REST API](./18_rest_api.md)
- `net/http` thuần, sau đó framework (Gin/Fiber/Echo — chọn 1), routing, middleware, request binding & validation (`validator`).
- **Bài tập:** Viết API `GET/POST /todos` bằng `net/http` thuần trước, sau đó viết lại bằng Gin để so sánh.

### [Bài 19: Authentication & Authorization](./19_auth.md)
- Hash password (`bcrypt`), JWT (tạo/verify token), middleware xác thực, phân quyền theo role (RBAC) hoặc permission-based.

### [Bài 20: Logging, Config, Deployment cơ bản](./20_logging_config_deploy.md)
- Structured logging (`log/slog` hoặc `zap`), quản lý config bằng `.env`/`viper`, Dockerize ứng dụng Go, graceful shutdown.

---

## Giai đoạn 4 — Dự Án: REST API có Phân Quyền (Capstone)

> **File chi tiết đầy đủ (kèm code từng bước): [21_capstone_project.md](./21_capstone_project.md)**

**Đề bài:** Xây dựng API quản lý "Task" (hoặc "Blog") với 2 role: `admin` và `user`. `user` chỉ CRUD task của chính mình; `admin` xem/sửa/xóa mọi task và quản lý user.

### Cấu trúc thư mục đề xuất (theo [golang-standards/project-layout](https://github.com/golang-standards/project-layout))
```
myapi/
├── cmd/
│   └── api/
│       └── main.go              # entry point, khởi tạo server
├── internal/
│   ├── config/                  # load .env / config
│   ├── domain/                  # struct models + interface (User, Task, repository interfaces)
│   ├── handler/                 # HTTP handlers (nhận request, gọi service, trả response)
│   ├── service/                 # business logic
│   ├── repository/               # tương tác database (Postgres/MySQL)
│   ├── middleware/               # JWT auth, RBAC, logging, recover
│   └── router/                  # khai báo route + gắn middleware
├── pkg/                          # code dùng chung, có thể tái sử dụng ở project khác
│   ├── jwtutil/
│   └── hash/
├── migrations/                   # file SQL migration
├── .env
├── go.mod
├── go.sum
└── Dockerfile
```

### Các bước triển khai (chia thành sub-lesson, mỗi bước = 1 buổi học)

1. **Bootstrap project**: `go mod init`, cấu trúc thư mục trên, load config bằng `viper`/`godotenv`.
2. **Domain models**: struct `User`, `Task`, interface `UserRepository`, `TaskRepository`.
3. **Database layer**: kết nối Postgres bằng `pgx`/GORM, viết migration tạo bảng `users`, `tasks` (có cột `role`).
4. **Repository layer**: implement interface ở bước 2 (CRUD thật với DB).
5. **Service layer**: business logic — đăng ký (hash password), đăng nhập (so sánh hash, phát JWT), tạo/sửa/xóa task với kiểm tra quyền sở hữu.
6. **Middleware xác thực**: parse JWT từ header `Authorization: Bearer ...`, gắn `userID`, `role` vào context.
7. **Middleware phân quyền (RBAC)**: middleware `RequireRole("admin")`, hoặc kiểm tra "chỉ chủ sở hữu hoặc admin" ngay trong service.
8. **Handler + Router**: expose route:
   - `POST /auth/register`, `POST /auth/login`
   - `GET/POST /tasks`, `GET/PUT/DELETE /tasks/:id` (yêu cầu JWT)
   - `GET /admin/users` (yêu cầu role admin)
9. **Validation & error response chuẩn hóa**: dùng `go-playground/validator`, format lỗi JSON thống nhất (`{"error": "..."}`).
10. **Testing**: unit test cho service (mock repository qua interface), integration test cho handler bằng `httptest`.
11. **Logging & graceful shutdown**: log request/response, xử lý `SIGINT`/`SIGTERM` để đóng DB pool sạch sẽ.
12. **Docker hóa**: viết `Dockerfile` multi-stage build, `docker-compose.yml` chạy kèm Postgres.
13. **(Mở rộng)** Refresh token, rate limiting, Swagger/OpenAPI docs (`swaggo`), CI đơn giản bằng GitHub Actions chạy `go test`.

### Tiêu chí hoàn thành dự án
- [ ] Đăng ký/đăng nhập trả JWT hợp lệ, password không lưu plaintext.
- [ ] User thường không thể sửa/xóa task của người khác (test bằng Postman/curl).
- [ ] Admin có thể thao tác trên mọi task và xem danh sách user.
- [ ] Có ít nhất 5 test case (unit + integration) chạy pass qua `go test ./...`.
- [ ] Chạy được bằng `docker-compose up`.

---

## Gợi ý cách học
- Mỗi bài học: đọc lý thuyết ngắn → gõ lại code mẫu (không copy-paste) → làm hết bài tập → tự giải thích lại bằng lời cho chính mình (Feynman technique).
- Sau Giai đoạn 1, bắt đầu làm song song 1-2 bài LeetCode/Exercism bằng Go mỗi tuần để luyện cú pháp.
- Sau Giai đoạn 2, đọc source code thư viện chuẩn (`net/http`, `encoding/json`) để học cách người ta thiết kế API.
- Giai đoạn 4 nên làm thật, push lên GitHub, viết README mô tả kiến trúc — đây sẽ là project để show trong CV.

## Tài liệu tham khảo
- Tour of Go: https://go.dev/tour/
- Effective Go: https://go.dev/doc/effective_go
- Go by Example: https://gobyexample.com/
- Standard project layout: https://github.com/golang-standards/project-layout
