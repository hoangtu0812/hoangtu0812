# Bài 20: Logging, Config, Deployment cơ bản

## Mục tiêu
- Structured logging với `log/slog` (thư viện chuẩn từ Go 1.21).
- Quản lý config bằng biến môi trường / file `.env`.
- Dockerize ứng dụng Go, graceful shutdown.

## 1. Structured logging với `log/slog`

`log/slog` (chuẩn từ Go 1.21) thay thế `log` cũ, hỗ trợ log dạng key-value có cấu trúc — dễ parse bởi hệ thống log tập trung (ELK, Loki...).

```go
import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("server đang khởi động", "port", 8080, "env", "production")
	slog.Warn("connection pool gần đầy", "current", 23, "max", 25)
	slog.Error("không thể kết nối database", "error", err, "retry_count", 3)
}
```

Output dạng JSON:
```json
{"time":"2026-08-04T10:00:00Z","level":"INFO","msg":"server đang khởi động","port":8080,"env":"production"}
```

Dùng trong middleware để log mỗi request:

```go
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
```

## 2. Quản lý config

### Cách đơn giản: đọc trực tiếp biến môi trường

```go
import "os"

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
	}
}
```

### File `.env` (dùng `godotenv` để load lúc dev, production thường set env trực tiếp)

```powershell
go get github.com/joho/godotenv
```

```go
import "github.com/joho/godotenv"

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("không tìm thấy file .env, dùng biến môi trường hệ thống")
	}
	cfg := LoadConfig()
	// ...
}
```

```
# .env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
JWT_SECRET=doi-secret-nay-trong-production
```

**Quan trọng:** file `.env` chứa secret KHÔNG được commit vào git — luôn thêm vào `.gitignore` và chỉ commit `.env.example` (không có giá trị thật) làm mẫu.

### Cách nâng cao hơn: `viper` (hỗ trợ nhiều nguồn config: file YAML, env, flag...)

```go
import "github.com/spf13/viper"

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // biến môi trường override giá trị trong file

	viper.ReadInConfig()

	return Config{
		Port:        viper.GetString("port"),
		DatabaseURL: viper.GetString("database_url"),
	}
}
```

## 3. Graceful shutdown

Khi server nhận tín hiệu dừng (`Ctrl+C`, hoặc `SIGTERM` từ Docker/Kubernetes), nên đóng kết nối đang xử lý dở một cách sạch sẽ thay vì cắt ngang đột ngột:

```go
func main() {
	srv := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server lỗi", "error", err)
			os.Exit(1)
		}
	}()
	slog.Info("server đang chạy tại :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // block cho tới khi nhận tín hiệu dừng

	slog.Info("đang tắt server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil { // chờ request đang xử lý hoàn thành, tối đa 10s
		slog.Error("shutdown lỗi", "error", err)
	}
	slog.Info("server đã tắt hoàn toàn")
}
```

## 4. Dockerize — multi-stage build

```dockerfile
# Dockerfile
# Stage 1: build binary
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# Stage 2: image chạy thật, nhỏ gọn (không chứa Go toolchain)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

Multi-stage build giúp image cuối cùng chỉ chứa binary đã build (vài MB) thay vì toàn bộ Go toolchain (~300MB+).

## 5. `docker-compose.yml` — chạy app kèm Postgres

```yaml
version: "3.9"
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://user:pass@db:5432/mydb?sslmode=disable
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

Chạy: `docker-compose up --build`

## Bài tập

1. **Structured logging**: thay toàn bộ `fmt.Println`/`log.Println` trong 1 project nhỏ bạn đã viết trước đó bằng `slog`, thêm middleware log mỗi request (method, path, thời gian xử lý).
2. **Load config từ `.env`**: viết `Config` struct + `LoadConfig()`, tạo file `.env` với `PORT`, `DATABASE_URL`, đảm bảo `.env` nằm trong `.gitignore`.
3. **Graceful shutdown**: thêm graceful shutdown vào server HTTP đã viết ở [Bài 18](./18_rest_api.md), test bằng cách gửi `Ctrl+C` và quan sát log "đang tắt server..." xuất hiện trước khi process thực sự thoát.
4. **Dockerize**: viết `Dockerfile` multi-stage cho project của bạn, build image, chạy `docker run -p 8080:8080 <image>` và gọi thử API từ Postman.

## Tổng kết Giai đoạn 3
Bạn đã có đủ kiến thức để xây dựng 1 REST API hoàn chỉnh: kiến trúc phân tầng, database, HTTP/middleware, authentication/authorization, logging/config/deploy. Bước tiếp theo là ghép TẤT CẢ lại thành 1 dự án thực tế hoàn chỉnh.

## Tiếp theo
→ [Dự án Capstone: REST API có phân quyền](./21_capstone_project.md)
