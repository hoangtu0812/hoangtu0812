# Bài 18: HTTP Server & REST API

## Mục tiêu
- Viết HTTP server bằng `net/http` thuần (Go 1.22+ đã có routing pattern tốt sẵn).
- So sánh và viết lại bằng framework Gin.
- Middleware, request binding & validation.

## 1. `net/http` thuần (Go 1.22+)

Từ Go 1.22, `http.ServeMux` hỗ trợ method + path pattern trực tiếp (không cần framework cho case đơn giản):

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var todos = map[int]*Todo{
	1: {ID: 1, Text: "Học Go", Done: false},
}
var nextID = 2

func listTodos(w http.ResponseWriter, r *http.Request) {
	result := make([]*Todo, 0, len(todos))
	for _, t := range todos {
		result = append(result, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "body không hợp lệ", http.StatusBadRequest)
		return
	}
	t.ID = nextID
	nextID++
	todos[t.ID] = &t

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id")) // r.PathValue là tính năng mới từ Go 1.22
	t, ok := todos[id]
	if !ok {
		http.Error(w, "không tìm thấy", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", listTodos)
	mux.HandleFunc("POST /todos", createTodo)
	mux.HandleFunc("GET /todos/{id}", getTodo)

	log.Println("Server đang chạy tại :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## 2. Middleware với `net/http` thuần

Middleware là hàm bọc quanh handler để xử lý logic chung (logging, auth, CORS...):

```go
type Middleware func(http.Handler) http.Handler

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // gọi handler tiếp theo trong chuỗi
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", listTodos)

	handler := loggingMiddleware(mux) // bọc toàn bộ mux bằng middleware
	http.ListenAndServe(":8080", handler)
}
```

## 3. Framework Gin — cùng ví dụ, viết gọn hơn

```powershell
go get github.com/gin-gonic/gin
```

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text" binding:"required"` // validation ngay trên struct tag
	Done bool   `json:"done"`
}

func main() {
	r := gin.Default() // đã có sẵn logging + recovery middleware

	r.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"todos": []Todo{}})
	})

	r.POST("/todos", func(c *gin.Context) {
		var t Todo
		if err := c.ShouldBindJSON(&t); err != nil { // tự động validate theo tag "binding"
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, t)
	})

	r.GET("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	r.Run(":8080")
}
```

## 4. Middleware trong Gin

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "thiếu token"})
			c.Abort() // dừng chuỗi handler, KHÔNG gọi handler tiếp theo
			return
		}
		c.Next() // cho phép tiếp tục tới handler tiếp theo
	}
}

func main() {
	r := gin.Default()

	// Middleware áp dụng cho 1 group route
	authorized := r.Group("/api")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/profile", getProfileHandler)
	}
}
```

## 5. Validation nâng cao với `go-playground/validator`

Gin đã tích hợp sẵn validator qua struct tag `binding`:

```go
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Age      int    `json:"age" binding:"gte=0,lte=150"`
}
```

## 6. Response format chuẩn hóa

```go
type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func respondSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, APIResponse{Success: true, Data: data})
}

func respondError(c *gin.Context, status int, err error) {
	c.JSON(status, APIResponse{Success: false, Error: err.Error()})
}
```

## Bài tập

1. **API Todo bằng `net/http` thuần**: viết đủ `GET/POST /todos`, `GET/PUT/DELETE /todos/{id}` theo code mẫu, test bằng `curl` hoặc Postman.
2. **Viết lại bằng Gin**: chuyển API trên sang Gin, thêm validation bằng struct tag `binding`.
3. **Middleware logging**: viết middleware log ra `method, path, status code, thời gian xử lý` cho mỗi request (thử cả 2 cách: `net/http` thuần và Gin).
4. **Nâng cao**: chuẩn hóa response theo format `APIResponse{success, data, error}` như trên, áp dụng cho toàn bộ route Todo, viết middleware `recover` bắt panic để trả JSON lỗi 500 thay vì crash server.

## Tiếp theo
→ [Bài 19: Authentication & Authorization](./19_auth.md)
