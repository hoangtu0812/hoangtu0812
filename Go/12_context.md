# Bài 12: Context

## Mục tiêu
- Hiểu vai trò của `context.Context` trong việc hủy/timeout goroutine.
- Dùng `context.WithCancel`, `context.WithTimeout`, `context.WithValue`.
- Biết idiom truyền `ctx` làm tham số đầu tiên của hàm — cực kỳ phổ biến khi viết API/service.

## 1. Tại sao cần Context?

Khi 1 request HTTP bị client hủy (đóng tab), hoặc cần giới hạn thời gian gọi database, bạn cần cách để "báo" cho tất cả goroutine liên quan dừng lại — đó chính là việc của `context.Context`.

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
```

## 2. `context.WithCancel` — hủy thủ công

```go
func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done(): // channel này sẽ đóng khi ctx bị cancel
			fmt.Printf("Worker %d dừng: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("Worker %d đang chạy\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx, 1)

	time.Sleep(2 * time.Second)
	cancel() // báo hiệu tất cả goroutine dùng ctx này nên dừng lại

	time.Sleep(1 * time.Second) // chờ worker in log dừng
}
```

## 3. `context.WithTimeout` — tự động hủy sau khoảng thời gian

```go
func callExternalAPI(ctx context.Context) (string, error) {
	select {
	case <-time.After(3 * time.Second): // giả lập API mất 3s để trả lời
		return "kết quả từ API", nil
	case <-ctx.Done():
		return "", ctx.Err() // context deadline exceeded
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // LUÔN gọi cancel() để giải phóng resource, dù timeout đã tự kích hoạt

	result, err := callExternalAPI(ctx)
	if err != nil {
		fmt.Println("Lỗi:", err) // "context deadline exceeded" sau 1s
		return
	}
	fmt.Println("Kết quả:", result)
}
```

**Quy tắc:** luôn `defer cancel()` ngay sau khi tạo context bằng `WithCancel`/`WithTimeout`/`WithDeadline` để tránh rò rỉ goroutine/resource.

## 4. `context.WithValue` — truyền dữ liệu theo request

```go
type contextKey string

const userIDKey contextKey = "userID"

func handleRequest(ctx context.Context) {
	ctx = context.WithValue(ctx, userIDKey, "user-123")
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		fmt.Println("Đang xử lý request cho user:", userID)
	}
}
```

**Lưu ý:** `context.WithValue` chỉ nên dùng cho dữ liệu **liên quan tới request** (request ID, user ID sau khi xác thực, trace ID cho logging) — **KHÔNG** dùng để truyền tham số nghiệp vụ thông thường (đó nên là tham số hàm bình thường, rõ ràng hơn nhiều).

## 5. Idiom: `ctx` luôn là tham số đầu tiên

Đây là convention chuẩn trong toàn bộ hệ sinh thái Go — sẽ dùng xuyên suốt ở tầng service/repository của dự án capstone:

```go
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*User, error) {
	// ctx được truyền tiếp xuống query DB, để query tự hủy nếu request bị timeout/cancel
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)
	// ...
	return nil, nil
}
```

## Ví dụ đầy đủ

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func callExternalAPI(ctx context.Context) (string, error) {
	select {
	case <-time.After(3 * time.Second):
		return "dữ liệu từ API bên ngoài", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result, err := callExternalAPI(ctx)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Lỗi sau %v: %v\n", elapsed, err)
		return
	}
	fmt.Println("Thành công:", result)
}
```

## Bài tập

1. **Timeout cơ bản**: dùng code mẫu trên làm nền, tự viết lại `callExternalAPI` giả lập mất 3s, gọi với `context.WithTimeout(1s)` để chứng minh nó bị hủy sớm; sau đó đổi timeout lên 5s để chứng minh gọi thành công.
2. **`WithCancel` thủ công**: viết 1 worker chạy loop in log mỗi 500ms, dùng `WithCancel` để dừng nó sau 2 giây từ `main`.
3. **Truyền `ctx` qua nhiều tầng hàm**: viết 3 hàm `handler → service → repository` (giả lập), mỗi hàm nhận `ctx context.Context` làm tham số đầu tiên và truyền tiếp xuống tầng dưới; ở tầng `repository`, kiểm tra `ctx.Err()` trước khi "truy vấn" (giả lập bằng `time.Sleep`).
4. **Nâng cao**: viết hàm chạy 3 goroutine gọi 3 "API" giả lập với thời gian khác nhau (0.5s, 1.5s, 2.5s) dùng chung 1 `context.WithTimeout(2s)`, in ra goroutine nào thành công/goroutine nào bị hủy do timeout.

## Tiếp theo
→ [Bài 13: Xử lý JSON & File I/O](./13_json_file_io.md)
