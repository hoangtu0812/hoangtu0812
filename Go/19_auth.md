# Bài 19: Authentication & Authorization

## Mục tiêu
- Hash password an toàn bằng `bcrypt`.
- Tạo/verify JWT (JSON Web Token).
- Viết middleware xác thực (authentication).
- Phân quyền theo role — RBAC (authorization).

## 1. Authentication vs Authorization — phân biệt rõ

- **Authentication (xác thực)**: bạn là ai? → đăng nhập, kiểm tra password, cấp token.
- **Authorization (phân quyền)**: bạn được làm gì? → kiểm tra role/permission trước khi cho phép hành động.

## 2. Hash password với `bcrypt`

**Không bao giờ** lưu password dạng plaintext trong database.

```powershell
go get golang.org/x/crypto/bcrypt
```

```go
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
	// cost càng cao càng an toàn nhưng càng chậm — 10-12 là mức phổ biến
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
```

## 3. Tạo & verify JWT

```powershell
go get github.com/golang-jwt/jwt/v5
```

```go
import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-doc-tu-env-KHONG-hardcode") // LUÔN đọc từ biến môi trường trong thực tế

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func generateToken(userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func verifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token không hợp lệ: %w", err)
	}
	return claims, nil
}
```

## 4. Flow đăng ký / đăng nhập

```go
func Register(ctx context.Context, repo UserRepository, name, email, password string) (*User, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:         name,
		Email:        email,
		PasswordHash: hashed,
		Role:         "user", // mặc định role thấp nhất, admin phải được cấp thủ công/qua tầng khác
	}

	if err := repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func Login(ctx context.Context, repo UserRepository, email, password string) (string, error) {
	user, err := repo.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials // KHÔNG tiết lộ "email không tồn tại" — tránh dò email hợp lệ
	}

	if !checkPassword(user.PasswordHash, password) {
		return "", ErrInvalidCredentials
	}

	return generateToken(user.ID, user.Role)
}
```

## 5. Middleware xác thực (Gin)

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "thiếu hoặc sai định dạng token"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")

		claims, err := verifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token không hợp lệ hoặc đã hết hạn"})
			c.Abort()
			return
		}

		// Gắn thông tin user vào context để handler phía sau dùng
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
```

## 6. Middleware phân quyền (RBAC)

```go
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "chưa xác thực"})
			c.Abort()
			return
		}

		allowed := false
		for _, r := range roles {
			if userRole == r {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "không đủ quyền truy cập"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	api := r.Group("/api")
	api.Use(AuthMiddleware()) // mọi route trong /api đều cần đăng nhập

	api.GET("/tasks", listMyTasks)         // user thường được phép
	api.POST("/tasks", createTask)

	admin := api.Group("/admin")
	admin.Use(RequireRole("admin")) // CHỈ admin mới vào được nhóm route này
	{
		admin.GET("/users", listAllUsers)
		admin.DELETE("/users/:id", deleteUser)
	}
}
```

## 7. Phân quyền theo chủ sở hữu (ownership) — quan trọng cho dự án capstone

RBAC theo role thôi chưa đủ — cần kiểm tra thêm "user chỉ được sửa task CỦA CHÍNH MÌNH":

```go
func (s *TaskService) UpdateTask(ctx context.Context, requesterID int, requesterRole string, taskID int, updates TaskUpdate) error {
	task, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	// admin luôn được phép; user thường chỉ được sửa task của chính mình
	if requesterRole != "admin" && task.OwnerID != requesterID {
		return ErrForbidden
	}

	return s.repo.Update(ctx, taskID, updates)
}
```

Logic ownership này nên nằm ở **tầng service**, không phải middleware — vì cần biết dữ liệu cụ thể (task thuộc về ai) mới quyết định được, còn middleware chỉ biết role chung chung.

## Bài tập

1. **Hash & verify password**: viết `hashPassword`/`checkPassword` bằng `bcrypt`, test với password đúng/sai.
2. **Tạo & verify JWT**: viết `generateToken`/`verifyToken`, thử với token hết hạn (set `ExpiresAt` trong quá khứ) để xác nhận `verifyToken` trả lỗi đúng như mong đợi.
3. **Middleware auth + RBAC**: viết `AuthMiddleware` và `RequireRole` như trên (dùng `net/http` thuần hoặc Gin tùy bạn chọn ở Bài 18), tạo 3 route: public, cần login, cần role admin — test bằng Postman với các token khác nhau.
4. **Nâng cao — ownership check**: viết `UpdateTask` như ví dụ trên, viết test (dùng mock repository từ [Bài 14](./14_testing.md)) cho 3 case: admin sửa task người khác (được phép), user sửa task của mình (được phép), user sửa task người khác (bị từ chối).

## Tiếp theo
→ [Bài 20: Logging, Config, Deployment](./20_logging_config_deploy.md)
