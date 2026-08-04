# Bài 13: Xử lý JSON & File I/O

## Mục tiêu
- Marshal/Unmarshal JSON với `encoding/json`, dùng struct tag.
- Đọc/ghi file với `os`, `bufio`, `io`.
- Đây là kỹ năng bắt buộc để xây REST API (request/response body đều là JSON).

## 1. Struct tag cho JSON

```go
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`                    // "-" : KHÔNG BAO GIỜ xuất hiện trong JSON (dùng cho field nhạy cảm)
	Age      int    `json:"age,omitempty"`         // omitempty: bỏ qua field nếu là zero value
}
```

## 2. Marshal (Go struct → JSON)

```go
func main() {
	u := User{ID: 1, Name: "Ben", Email: "ben@example.com", Password: "secret123"}

	data, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
	// {"id":1,"name":"Ben","email":"ben@example.com"}  <- Password bị ẩn nhờ json:"-"

	// MarshalIndent để in đẹp, dễ đọc (thường dùng khi debug, không dùng trong response API thật vì tốn băng thông)
	pretty, _ := json.MarshalIndent(u, "", "  ")
	fmt.Println(string(pretty))
}
```

## 3. Unmarshal (JSON → Go struct)

```go
func main() {
	jsonStr := `{"id":2,"name":"Alice","email":"alice@example.com"}`

	var u User
	err := json.Unmarshal([]byte(jsonStr), &u) // LUÔN truyền con trỏ &u để Unmarshal có thể ghi dữ liệu vào
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", u)
}
```

Field nào trong JSON không khớp field nào trong struct sẽ tự động bị bỏ qua (không lỗi) — hữu ích khi API trả nhiều field hơn bạn cần.

## 4. Làm việc với JSON array và nested object

```go
type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type UserWithAddress struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

func main() {
	jsonStr := `{"name":"Ben","address":{"city":"Quang Ngai","country":"Vietnam"}}`
	var u UserWithAddress
	json.Unmarshal([]byte(jsonStr), &u)
	fmt.Println(u.Address.City) // Quang Ngai

	// Array
	jsonArr := `[{"id":1,"name":"Ben"},{"id":2,"name":"Alice"}]`
	var users []User
	json.Unmarshal([]byte(jsonArr), &users)
	fmt.Println(len(users)) // 2
}
```

## 5. Đọc/ghi file

```go
// Ghi file — cách đơn giản
func writeFile() error {
	content := []byte("Xin chào Go!\nDòng thứ hai.")
	return os.WriteFile("output.txt", content, 0644)
}

// Đọc toàn bộ file — cách đơn giản
func readFile() error {
	data, err := os.ReadFile("output.txt")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// Đọc file theo từng dòng — dùng bufio.Scanner
func readLineByLine(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close() // LUÔN đóng file sau khi dùng xong

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		lineCount++
	}
	fmt.Println("Tổng số dòng:", lineCount)
	return scanner.Err()
}
```

## 6. Ghi JSON ra file

```go
func writeUsersToFile(users []User, path string) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func readUsersFromFile(path string) ([]User, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var users []User
	err = json.Unmarshal(data, &users)
	return users, err
}
```

## Ví dụ đầy đủ

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	users := []User{
		{ID: 1, Name: "Ben", Email: "ben@example.com"},
		{ID: 2, Name: "Alice", Email: "alice@example.com"},
	}

	// Ghi ra file JSON
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("users.json", data, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Đã ghi users.json")

	// Đọc lại từ file
	raw, err := os.ReadFile("users.json")
	if err != nil {
		log.Fatal(err)
	}
	var loaded []User
	if err := json.Unmarshal(raw, &loaded); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Đã đọc lại %d user\n", len(loaded))
	for _, u := range loaded {
		fmt.Printf("- %+v\n", u)
	}
}
```

## Bài tập

1. **Marshal/Unmarshal qua lại**: định nghĩa struct `User` với json tag (bao gồm 1 field `json:"-"` để ẩn), marshal thành JSON string, rồi unmarshal ngược lại thành struct, in ra so sánh.
2. **Đọc file theo dòng**: tạo 1 file `.txt` bất kỳ (nhiều dòng), viết hàm đọc và đếm số dòng bằng `bufio.Scanner`.
3. **Ghi danh sách `User` ra JSON**: dùng code mẫu trên làm nền, tự viết lại — ghi 1 slice `[]User` ra file `users.json`, sau đó đọc lại và in ra.
4. **Nâng cao**: viết struct lồng nhau (`Order` chứa `[]OrderItem`, mỗi `OrderItem` có `ProductName string`, `Quantity int`, `Price float64`), marshal/unmarshal, tính tổng giá trị đơn hàng sau khi unmarshal.

## Tiếp theo
→ [Bài 14: Testing trong Go](./14_testing.md)
