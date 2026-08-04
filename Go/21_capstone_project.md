# Dự án Capstone: REST API Quản Lý Task Có Phân Quyền

## Mục tiêu
Ghép toàn bộ kiến thức từ Bài 2 → Bài 20 thành 1 dự án hoàn chỉnh, đủ để đưa vào CV: API quản lý "Task" với 2 role (`admin`, `user`), kiến trúc phân tầng chuẩn, có test, có Docker.

**Kiến thức áp dụng:** [Bài 16](./16_clean_architecture.md) (kiến trúc), [Bài 17](./17_database.md) (database), [Bài 18](./18_rest_api.md) (REST API), [Bài 19](./19_auth.md) (auth/RBAC), [Bài 20](./20_logging_config_deploy.md) (logging/deploy).

## Yêu cầu chức năng

- `user` đăng ký/đăng nhập, tạo/xem/sửa/xóa **task của chính mình**.
- `admin` xem/sửa/xóa **mọi task**, xem danh sách toàn bộ user.
- Mật khẩu hash bằng bcrypt, xác thực bằng JWT.
- Có test cho tầng service (mock repository).
- Chạy được bằng `docker-compose up`.

## Cấu trúc thư mục

```
taskapi/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── user.go          # struct User + interface UserRepository
│   │   └── task.go          # struct Task + interface TaskRepository
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── task_handler.go
│   │   └── admin_handler.go
│   ├── service/
│   │   ├── auth_service.go
│   │   └── task_service.go
│   ├── repository/
│   │   ├── user_postgres.go
│   │   └── task_postgres.go
│   ├── middleware/
│   │   ├── auth.go          # xác thực JWT
│   │   ├── rbac.go          # phân quyền theo role
│   │   └── logging.go
│   └── router/
│       └── router.go
├── pkg/
│   ├── jwtutil/
│   │   └── jwt.go
│   └── hash/
│       └── bcrypt.go
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_tasks_table.up.sql
│   └── 000002_create_tasks_table.down.sql
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Dockerfile
└── docker-compose.yml
```

## Bước 1: Bootstrap project

```powershell
mkdir taskapi; cd taskapi
go mod init example.com/taskapi
go get github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5/stdlib
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/joho/godotenv
go get github.com/stretchr/testify
```

`.env.example`:
```
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/taskdb?sslmode=disable
JWT_SECRET=change-me-in-production
```

## Bước 2: Domain models (`internal/domain/`)

```go
// internal/domain/user.go
package domain

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("không tìm thấy")
var ErrEmailExists = errors.New("email đã tồn tại")

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	Role         Role
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]*User, error)
}
```

```go
// internal/domain/task.go
package domain

import (
	"context"
	"errors"
)

var ErrForbidden = errors.New("không có quyền thực hiện")

type Task struct {
	ID      int
	OwnerID int
	Title   string
	Done    bool
}

type TaskRepository interface {
	Create(ctx context.Context, t *Task) error
	FindByID(ctx context.Context, id int) (*Task, error)
	ListByOwner(ctx context.Context, ownerID int) ([]*Task, error)
	ListAll(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id int) error
}
```

## Bước 3: Migration (`migrations/`)

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	role VARCHAR(20) NOT NULL DEFAULT 'user',
	created_at TIMESTAMP NOT NULL DEFAULT now()
);
```

```sql
-- 000002_create_tasks_table.up.sql
CREATE TABLE tasks (
	id SERIAL PRIMARY KEY,
	owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	title VARCHAR(255) NOT NULL,
	done BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP NOT NULL DEFAULT now()
);
```

## Bước 4-5: Repository & pkg helper

```go
// pkg/hash/bcrypt.go
package hash

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func Check(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

```go
// pkg/jwtutil/jwt.go
package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func Generate(secret string, userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func Verify(secret, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token không hợp lệ")
	}
	return claims, nil
}
```

```go
// internal/repository/task_postgres.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	"example.com/taskapi/internal/domain"
)

type taskPostgresRepo struct{ db *sql.DB }

func NewTaskPostgresRepo(db *sql.DB) domain.TaskRepository {
	return &taskPostgresRepo{db: db}
}

func (r *taskPostgresRepo) Create(ctx context.Context, t *domain.Task) error {
	return r.db.QueryRowContext(ctx,
		"INSERT INTO tasks (owner_id, title, done) VALUES ($1, $2, $3) RETURNING id",
		t.OwnerID, t.Title, t.Done,
	).Scan(&t.ID)
}

func (r *taskPostgresRepo) FindByID(ctx context.Context, id int) (*domain.Task, error) {
	var t domain.Task
	err := r.db.QueryRowContext(ctx,
		"SELECT id, owner_id, title, done FROM tasks WHERE id = $1", id,
	).Scan(&t.ID, &t.OwnerID, &t.Title, &t.Done)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return &t, err
}

func (r *taskPostgresRepo) ListByOwner(ctx context.Context, ownerID int) ([]*domain.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, owner_id, title, done FROM tasks WHERE owner_id = $1", ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.OwnerID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

func (r *taskPostgresRepo) ListAll(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, owner_id, title, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.OwnerID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

func (r *taskPostgresRepo) Update(ctx context.Context, t *domain.Task) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tasks SET title = $1, done = $2 WHERE id = $3", t.Title, t.Done, t.ID)
	return err
}

func (r *taskPostgresRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}
```

`user_postgres.go` viết tương tự (xem lại mẫu ở [Bài 16](./16_clean_architecture.md) và [Bài 17](./17_database.md)).

## Bước 6: Service layer — chứa toàn bộ logic phân quyền ownership

```go
// internal/service/task_service.go
package service

import (
	"context"

	"example.com/taskapi/internal/domain"
)

type TaskService struct {
	repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, ownerID int, title string) (*domain.Task, error) {
	t := &domain.Task{OwnerID: ownerID, Title: title}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// ListTasks: user thường chỉ thấy task của mình, admin thấy tất cả
func (s *TaskService) ListTasks(ctx context.Context, requesterID int, requesterRole domain.Role) ([]*domain.Task, error) {
	if requesterRole == domain.RoleAdmin {
		return s.repo.ListAll(ctx)
	}
	return s.repo.ListByOwner(ctx, requesterID)
}

// UpdateTask: chỉ chủ sở hữu hoặc admin mới được sửa — LOGIC PHÂN QUYỀN CỐT LÕI nằm ở đây
func (s *TaskService) UpdateTask(ctx context.Context, requesterID int, requesterRole domain.Role, taskID int, title string, done bool) (*domain.Task, error) {
	task, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if requesterRole != domain.RoleAdmin && task.OwnerID != requesterID {
		return nil, domain.ErrForbidden
	}

	task.Title = title
	task.Done = done
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, requesterID int, requesterRole domain.Role, taskID int) error {
	task, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	if requesterRole != domain.RoleAdmin && task.OwnerID != requesterID {
		return domain.ErrForbidden
	}
	return s.repo.Delete(ctx, taskID)
}
```

```go
// internal/service/auth_service.go
package service

import (
	"context"

	"example.com/taskapi/internal/domain"
	"example.com/taskapi/pkg/hash"
	"example.com/taskapi/pkg/jwtutil"
)

type AuthService struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewAuthService(repo domain.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	existing, _ := s.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, domain.ErrEmailExists
	}

	hashed, err := hash.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{Name: name, Email: email, PasswordHash: hashed, Role: domain.RoleUser}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", domain.ErrNotFound
	}
	if !hash.Check(user.PasswordHash, password) {
		return "", domain.ErrNotFound // cố tình trả lỗi giống nhau để không lộ email tồn tại hay không
	}
	return jwtutil.Generate(s.jwtSecret, user.ID, string(user.Role))
}
```

## Bước 7-8: Middleware + Handler + Router

```go
// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"example.com/taskapi/pkg/jwtutil"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "thiếu token"})
			c.Abort()
			return
		}
		claims, err := jwtutil.Verify(jwtSecret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token không hợp lệ"})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
```

```go
// internal/middleware/rbac.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "không đủ quyền"})
		c.Abort()
	}
}
```

```go
// internal/handler/task_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/taskapi/internal/domain"
	"example.com/taskapi/internal/service"
)

type TaskHandler struct{ service *service.TaskService }

func NewTaskHandler(s *service.TaskService) *TaskHandler { return &TaskHandler{service: s} }

type createTaskRequest struct {
	Title string `json:"title" binding:"required"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt("userID")

	task, err := h.service.CreateTask(c.Request.Context(), userID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	userID := c.GetInt("userID")
	role := domain.Role(c.GetString("role"))

	tasks, err := h.service.ListTasks(c.Request.Context(), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

type updateTaskRequest struct {
	Title string `json:"title" binding:"required"`
	Done  bool   `json:"done"`
}

func (h *TaskHandler) Update(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("userID")
	role := domain.Role(c.GetString("role"))

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.UpdateTask(c.Request.Context(), userID, role, taskID, req.Title, req.Done)
	switch err {
	case nil:
		c.JSON(http.StatusOK, task)
	case domain.ErrForbidden:
		c.JSON(http.StatusForbidden, gin.H{"error": "bạn không có quyền sửa task này"})
	case domain.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "không tìm thấy task"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
```

```go
// internal/router/router.go
package router

import (
	"github.com/gin-gonic/gin"

	"example.com/taskapi/internal/handler"
	"example.com/taskapi/internal/middleware"
)

func Setup(authH *handler.AuthHandler, taskH *handler.TaskHandler, adminH *handler.AdminHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	api := r.Group("/api")
	api.Use(middleware.Auth(jwtSecret))
	{
		api.GET("/tasks", taskH.List)
		api.POST("/tasks", taskH.Create)
		api.PUT("/tasks/:id", taskH.Update)
		api.DELETE("/tasks/:id", taskH.Delete)

		admin := api.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/users", adminH.ListUsers)
		}
	}

	return r
}
```

## Bước 9: `main.go` — nối tất cả

```go
// cmd/api/main.go
package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"

	"example.com/taskapi/internal/handler"
	"example.com/taskapi/internal/repository"
	"example.com/taskapi/internal/router"
	"example.com/taskapi/internal/service"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		slog.Error("không kết nối được DB", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := repository.NewUserPostgresRepo(db)
	taskRepo := repository.NewTaskPostgresRepo(db)

	authService := service.NewAuthService(userRepo, jwtSecret)
	taskService := service.NewTaskService(taskRepo)

	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)
	adminHandler := handler.NewAdminHandler(userRepo)

	r := router.Setup(authHandler, taskHandler, adminHandler, jwtSecret)

	slog.Info("server đang chạy", "port", port)
	r.Run(":" + port)
}
```

## Bước 10: Test tầng service (không cần DB thật)

```go
// internal/service/task_service_test.go
package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"example.com/taskapi/internal/domain"
	"example.com/taskapi/internal/service"
)

type mockTaskRepo struct {
	tasks map[int]*domain.Task
}

func (m *mockTaskRepo) Create(ctx context.Context, t *domain.Task) error { return nil }
func (m *mockTaskRepo) FindByID(ctx context.Context, id int) (*domain.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return t, nil
}
func (m *mockTaskRepo) ListByOwner(ctx context.Context, ownerID int) ([]*domain.Task, error) { return nil, nil }
func (m *mockTaskRepo) ListAll(ctx context.Context) ([]*domain.Task, error)                  { return nil, nil }
func (m *mockTaskRepo) Update(ctx context.Context, t *domain.Task) error                     { return nil }
func (m *mockTaskRepo) Delete(ctx context.Context, id int) error                             { return nil }

func TestTaskService_UpdateTask_Ownership(t *testing.T) {
	repo := &mockTaskRepo{tasks: map[int]*domain.Task{
		1: {ID: 1, OwnerID: 100, Title: "Task của user 100"},
	}}
	svc := service.NewTaskService(repo)
	ctx := context.Background()

	t.Run("admin sửa task người khác - được phép", func(t *testing.T) {
		_, err := svc.UpdateTask(ctx, 999, domain.RoleAdmin, 1, "Sửa bởi admin", true)
		require.NoError(t, err)
	})

	t.Run("chủ sở hữu sửa task của mình - được phép", func(t *testing.T) {
		_, err := svc.UpdateTask(ctx, 100, domain.RoleUser, 1, "Sửa bởi chủ", true)
		require.NoError(t, err)
	})

	t.Run("user khác sửa task không phải của mình - bị từ chối", func(t *testing.T) {
		_, err := svc.UpdateTask(ctx, 200, domain.RoleUser, 1, "Cố sửa", true)
		assert.ErrorIs(t, err, domain.ErrForbidden)
	})
}
```

## Bước 11-12: Dockerfile & docker-compose (xem đầy đủ ở [Bài 20](./20_logging_config_deploy.md))

Áp dụng y nguyên mẫu `Dockerfile` multi-stage và `docker-compose.yml` ở Bài 20, chỉnh `CMD ["./server"]` trỏ đúng `cmd/api`.

## Bước 13: Test thủ công bằng curl

```powershell
# Đăng ký
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d '{"name":"Ben","email":"ben@test.com","password":"12345678"}'

# Đăng nhập, lấy token
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{"email":"ben@test.com","password":"12345678"}'

# Tạo task (cần token)
curl -X POST http://localhost:8080/api/tasks -H "Authorization: Bearer <token>" -H "Content-Type: application/json" -d '{"title":"Học Go xong dự án capstone"}'
```

## Checklist hoàn thành

- [ ] `POST /auth/register`, `POST /auth/login` hoạt động, password không lưu plaintext (kiểm tra trực tiếp trong DB).
- [ ] `user` không sửa/xóa được task của người khác (test bằng 2 tài khoản khác nhau qua Postman).
- [ ] `admin` thao tác được trên mọi task, xem được `/api/admin/users`.
- [ ] `go test ./...` chạy pass toàn bộ, có ít nhất test case ở Bước 10.
- [ ] `docker-compose up --build` chạy thành công, API truy cập được từ `localhost:8080`.
- [ ] README của project mô tả kiến trúc + hướng dẫn chạy — dùng để show trong CV/GitHub.

## Mở rộng (tùy chọn, sau khi hoàn thành checklist)
- Refresh token + access token ngắn hạn.
- Rate limiting (giới hạn số request/phút theo IP hoặc user).
- Swagger/OpenAPI docs bằng `swaggo/swag`.
- CI đơn giản: GitHub Actions chạy `go vet` + `go test` mỗi lần push.
- Pagination cho `GET /api/tasks` khi dữ liệu lớn.

---
Chúc mừng — hoàn thành bài này nghĩa là bạn đã đi hết [lộ trình](./ROADMAP.md) từ Hello World tới 1 REST API production-ready có phân quyền, đúng như mục tiêu ban đầu.
