# Bài 17: Làm việc với Database

## Mục tiêu
- Kết nối PostgreSQL/MySQL từ Go bằng `database/sql`.
- Hiểu connection pool.
- Biết lựa chọn giữa SQL thuần (`database/sql`) và ORM (GORM).
- Quản lý migration bằng `golang-migrate`.

## 1. `database/sql` — package chuẩn

Go không có driver database built-in — `database/sql` là interface chung, cần thêm driver cụ thể:

```powershell
go get github.com/jackc/pgx/v5/stdlib   # driver PostgreSQL (khuyến khích, hiệu năng tốt)
# hoặc
go get github.com/go-sql-driver/mysql   # driver MySQL
```

```go
import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // import với "_" vì chỉ cần side-effect (đăng ký driver)
)

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", "postgres://user:password@localhost:5432/mydb?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// sql.Open KHÔNG thực sự kết nối ngay — phải Ping để kiểm tra
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Cấu hình connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
```

`*sql.DB` **không phải** 1 kết nối đơn — nó là 1 **pool** kết nối, tự quản lý mở/đóng/tái sử dụng connection. Nên tạo 1 lần khi khởi động app, dùng chung xuyên suốt (không mở/đóng liên tục).

## 2. CRUD cơ bản

```go
// CREATE
func createUser(ctx context.Context, db *sql.DB, name, email string) (int, error) {
	var id int
	err := db.QueryRowContext(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		name, email,
	).Scan(&id)
	return id, err
}

// READ 1 dòng
func getUserByID(ctx context.Context, db *sql.DB, id int) (*User, error) {
	var u User
	err := db.QueryRowContext(ctx,
		"SELECT id, name, email FROM users WHERE id = $1", id,
	).Scan(&u.ID, &u.Name, &u.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %d: %w", id, ErrNotFound)
	}
	return &u, err
}

// READ nhiều dòng
func listUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // LUÔN đóng rows sau khi dùng

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err() // kiểm tra lỗi phát sinh trong lúc duyệt, không chỉ lỗi ban đầu
}

// UPDATE
func updateUserEmail(ctx context.Context, db *sql.DB, id int, email string) error {
	result, err := db.ExecContext(ctx, "UPDATE users SET email = $1 WHERE id = $2", email, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// DELETE
func deleteUser(ctx context.Context, db *sql.DB, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
```

**Luôn dùng tham số hóa** (`$1`, `$2`...) thay vì nối chuỗi SQL trực tiếp — nối chuỗi trực tiếp gây lỗ hổng **SQL Injection** (OWASP Top 10).

```go
// SAI — dễ bị SQL injection
query := "SELECT * FROM users WHERE email = '" + email + "'"

// ĐÚNG — driver tự escape giá trị an toàn
db.QueryContext(ctx, "SELECT * FROM users WHERE email = $1", email)
```

## 3. Transaction

```go
func transferBalance(ctx context.Context, db *sql.DB, fromID, toID int, amount float64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op nếu đã Commit thành công, an toàn để luôn defer

	if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID); err != nil {
		return err
	}

	return tx.Commit()
}
```

## 4. SQL thuần vs ORM (GORM)

| | `database/sql` thuần | GORM |
|---|---|---|
| Kiểm soát query | Tuyệt đối, SQL viết tay | Sinh tự động, có thể tối ưu chưa tốt |
| Tốc độ học | Cần biết SQL vững | Nhanh lên code, ít cần SQL |
| Phù hợp | Team quen SQL, cần hiệu năng cao | Prototype nhanh, CRUD đơn giản |

Ví dụ GORM:

```go
import "gorm.io/gorm"

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"uniqueIndex"`
}

func main() {
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&User{}) // tự tạo/cập nhật bảng theo struct

	db.Create(&User{Name: "Ben", Email: "ben@example.com"})

	var user User
	db.First(&user, "email = ?", "ben@example.com")

	db.Model(&user).Update("Name", "Ben Pham")
	db.Delete(&user)
}
```

**Gợi ý cho dự án capstone:** dù dùng GORM hay `database/sql` thuần, luôn giấu chi tiết này sau interface `UserRepository` (xem [Bài 16](./16_clean_architecture.md)) để phần còn lại của app không phụ thuộc trực tiếp.

## 5. Migration với `golang-migrate`

```powershell
migrate create -ext sql -dir migrations -seq create_users_table
```

Tạo ra 2 file:
```sql
-- migrations/000001_create_users_table.up.sql
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	role VARCHAR(50) NOT NULL DEFAULT 'user',
	password_hash VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now()
);
```

```sql
-- migrations/000001_create_users_table.down.sql
DROP TABLE users;
```

Chạy migration: `migrate -path migrations -database "postgres://..." up`

## Bài tập

1. **Kết nối DB**: cài PostgreSQL (hoặc dùng Docker: `docker run -e POSTGRES_PASSWORD=pass -p 5432:5432 postgres`), viết hàm `connectDB()` kết nối và `Ping()` thành công.
2. **CRUD `users`**: viết migration tạo bảng `users` (`id`, `name`, `email`, `role`), viết đủ 4 hàm CRUD như ví dụ trên bằng `database/sql`.
3. **Transaction**: viết hàm giả lập chuyển tiền giữa 2 tài khoản trong bảng `accounts`, đảm bảo dùng transaction để cả 2 update cùng thành công hoặc cùng rollback (thử cố tình gây lỗi ở bước 2 để kiểm chứng rollback hoạt động).
4. **Nâng cao**: viết `userPostgresRepo` implement interface `domain.UserRepository` từ [Bài 16](./16_clean_architecture.md), viết test cho nó cần 1 DB test riêng (hoặc dùng thư viện `sqlmock` để mock `*sql.DB`).

## Tiếp theo
→ [Bài 18: HTTP Server & REST API](./18_rest_api.md)
