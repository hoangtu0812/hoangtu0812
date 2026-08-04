# Bài 16: Thiết kế package & Clean Architecture cơ bản

## Mục tiêu
- Hiểu mô hình phân tầng handler → service → repository.
- Dùng interface để dependency injection thủ công (không cần framework DI).
- Tránh import cycle khi project lớn dần.

## 1. Vì sao cần phân tầng?

Khi project chỉ có vài chục dòng, viết tất cả trong `main.go` không sao. Nhưng khi có nhiều nghiệp vụ, nhiều nguồn dữ liệu (DB, cache, API ngoài), cần tách rõ trách nhiệm để:
- Dễ test (mock từng tầng độc lập — xem lại [Bài 14](./14_testing.md)).
- Dễ thay đổi implementation (đổi Postgres sang MySQL mà không sửa business logic).
- Nhiều người làm việc song song không giẫm chân nhau.

## 2. Mô hình 3 tầng chuẩn

```
Handler (HTTP)  →  Service (business logic)  →  Repository (data access)
     ↓                      ↓                           ↓
  nhận request        xử lý nghiệp vụ              CRUD với DB
  gọi service          validate, tính toán         không biết gì về HTTP
  trả response         KHÔNG biết gì về HTTP
```

**Nguyên tắc dòng phụ thuộc:** Handler phụ thuộc Service, Service phụ thuộc **interface** Repository — không tầng nào phụ thuộc ngược lại.

## 3. Domain layer — nơi định nghĩa interface

```go
// internal/domain/user.go
package domain

import "context"

type User struct {
	ID    int
	Name  string
	Email string
	Role  string // "admin" hoặc "user"
}

// Interface định nghĩa Ở ĐÂY (domain), KHÔNG ở tầng repository
// Service sẽ phụ thuộc vào interface này, không phụ thuộc trực tiếp Postgres/MySQL
type UserRepository interface {
	FindByID(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
}
```

## 4. Repository layer — implement interface

```go
// internal/repository/user_postgres.go
package repository

import (
	"context"
	"database/sql"

	"example.com/myapi/internal/domain"
)

type userPostgresRepo struct {
	db *sql.DB
}

// Constructor trả về INTERFACE, không trả về struct cụ thể
// -> caller (service) chỉ biết tới domain.UserRepository, không biết implementation là Postgres
func NewUserPostgresRepo(db *sql.DB) domain.UserRepository {
	return &userPostgresRepo{db: db}
}

func (r *userPostgresRepo) FindByID(ctx context.Context, id int) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email, role FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
// ... FindByEmail, Create tương tự
```

## 5. Service layer — business logic

```go
// internal/service/user_service.go
package service

import (
	"context"
	"fmt"

	"example.com/myapi/internal/domain"
)

type UserService struct {
	repo domain.UserRepository // phụ thuộc INTERFACE, không phụ thuộc Postgres
}

// Dependency injection thủ công: "tiêm" repo vào qua constructor
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user %d: %w", id, err)
	}
	return user, nil
}
```

## 6. Handler layer — HTTP

```go
// internal/handler/user_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/myapi/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
```

## 7. Nối tất cả lại — `main.go` là nơi duy nhất biết implementation cụ thể

```go
// cmd/api/main.go
func main() {
	db := connectPostgres() // implementation cụ thể CHỈ xuất hiện ở đây

	userRepo := repository.NewUserPostgresRepo(db)   // trả về interface domain.UserRepository
	userService := service.NewUserService(userRepo)   // service không biết đó là Postgres
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	http.ListenAndServe(":8080", mux)
}
```

Nếu sau này đổi sang MySQL hoặc dùng mock để test, chỉ cần thay dòng `repository.NewUserPostgresRepo(db)` — `UserService` và `UserHandler` **không cần sửa 1 dòng nào**.

## 8. Tránh import cycle

Quy tắc: tầng thấp hơn (domain) **không bao giờ** import tầng cao hơn (repository, service, handler). Luồng import chỉ đi 1 chiều:

```
handler → service → domain
repository → domain
main.go → tất cả
```

Nếu bạn thấy `domain` cần import `repository` để dùng 1 struct nào đó — đó là dấu hiệu struct đó nên nằm ở `domain`, không phải `repository`.

## Bài tập

1. **Tách 3 tầng cho `Rectangle`**: lấy ví dụ `shapes.Rectangle` (Bài 9), thiết kế lại thành `domain` (struct + interface `ShapeRepository` giả lập lưu trữ hình dạng), `service` (tính tổng diện tích của danh sách hình), `handler` (giả lập, chỉ cần in ra console thay vì HTTP thật).
2. **Mock repository**: viết 1 implementation `InMemoryUserRepo` (dùng `map[int]*User` thay vì database thật) implement cùng interface `UserRepository`, dùng nó để test `UserService` mà không cần kết nối DB (liên hệ [Bài 14](./14_testing.md)).
3. **Vẽ sơ đồ**: vẽ (bằng ASCII hoặc mô tả bằng lời) luồng phụ thuộc của 1 dự án bạn tưởng tượng có 2 tầng nghiệp vụ (`User` và `Order`), chỉ rõ tầng nào phụ thuộc tầng nào.

## Tiếp theo
→ [Bài 17: Làm việc với Database](./17_database.md)
