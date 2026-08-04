# Bài 8: Xử lý lỗi (Error Handling)

## Mục tiêu
- Hiểu triết lý "error là giá trị" của Go (khác exception của Java/Python).
- Tạo, wrap, và kiểm tra lỗi (`errors.New`, `fmt.Errorf`, `errors.Is`, `errors.As`).
- Tự định nghĩa custom error type.
- Hiểu `panic`/`recover` và khi nào (hiếm khi) nên dùng.

## 1. Interface `error`

Go không có `try/catch`. Lỗi trong Go chỉ là một **giá trị bình thường** thỏa interface:

```go
type error interface {
	Error() string
}
```

Bất kỳ kiểu nào có method `Error() string` đều là một `error`.

```go
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("không thể chia cho 0")
	}
	return a / b, nil
}

func main() {
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("Lỗi:", err)
		return // idiom: xử lý lỗi ngay, return sớm — tránh nested if sâu
	}
	fmt.Println("Kết quả:", result)
}
```

**Idiom chuẩn của Go**: luôn kiểm tra `if err != nil` ngay sau khi gọi hàm trả về error, KHÔNG bỏ qua lỗi (không dùng `_` để lờ lỗi trừ khi thực sự chắc chắn không cần).

## 2. `fmt.Errorf` và wrap lỗi với `%w`

```go
func readConfig() error {
	err := loadFile("config.yaml")
	if err != nil {
		// %w "bọc" lỗi gốc, giữ nguyên lỗi gốc để truy vết sau này
		return fmt.Errorf("không thể đọc config: %w", err)
	}
	return nil
}
```

## 3. `errors.Is` và `errors.As`

```go
var ErrNotFound = errors.New("không tìm thấy")

func findUser(id int) error {
	if id == 0 {
		return fmt.Errorf("tìm user %d: %w", id, ErrNotFound)
	}
	return nil
}

func main() {
	err := findUser(0)

	// errors.Is: kiểm tra lỗi có "chứa" (wrap) một lỗi cụ thể hay không
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Xử lý riêng cho trường hợp not found")
	}
}
```

`errors.As` dùng khi cần lấy ra một **custom error type** cụ thể từ chuỗi lỗi đã wrap:

```go
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field %q: %s", e.Field, e.Msg)
}

func validate(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Msg: "không được âm"}
	}
	return nil
}

func main() {
	err := validate(-5)

	var valErr *ValidationError
	if errors.As(err, &valErr) {
		fmt.Println("Field lỗi:", valErr.Field) // "age"
	}
}
```

## 4. Custom error type

Tự định nghĩa error type khi bạn cần đính kèm thêm dữ liệu (không chỉ là 1 chuỗi thông báo) — như ví dụ `ValidationError` ở trên. Rất hữu ích khi viết API để trả về mã lỗi HTTP tương ứng dựa trên loại lỗi.

## 5. `panic` và `recover`

`panic` dừng ngay luồng thực thi bình thường (tương tự exception không bắt được). Chỉ nên dùng cho lỗi **thực sự không thể phục hồi** (bug lập trình, dữ liệu bất biến bị vi phạm) — **KHÔNG dùng panic để xử lý lỗi nghiệp vụ thông thường** (đó là việc của `error`).

```go
func safeDivide(a, b int) (result int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Đã recover từ panic:", r)
			result = 0
		}
	}()

	result = a / b // chia cho 0 với số nguyên sẽ panic runtime error
	return
}

func main() {
	fmt.Println(safeDivide(10, 2)) // 5
	fmt.Println(safeDivide(10, 0)) // in log recover, trả về 0 thay vì crash
}
```

**Quy tắc thực tế:** trong code nghiệp vụ, luôn ưu tiên trả `error`. `recover` chủ yếu dùng ở tầng cao nhất (vd middleware HTTP) để chương trình không crash hoàn toàn khi có bug bất ngờ ở tầng dưới.

## Ví dụ đầy đủ

```go
package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("không tìm thấy")

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field %q: %s", e.Field, e.Msg)
}

func findUser(id int) (string, error) {
	if id < 0 {
		return "", &ValidationError{Field: "id", Msg: "không được âm"}
	}
	if id == 0 {
		return "", fmt.Errorf("findUser: %w", ErrNotFound)
	}
	return "Ben", nil
}

func main() {
	for _, id := range []int{-1, 0, 1} {
		name, err := findUser(id)
		if err != nil {
			var valErr *ValidationError
			switch {
			case errors.As(err, &valErr):
				fmt.Println("Lỗi validate:", valErr.Field, "-", valErr.Msg)
			case errors.Is(err, ErrNotFound):
				fmt.Println("Không tìm thấy user với id:", id)
			default:
				fmt.Println("Lỗi khác:", err)
			}
			continue
		}
		fmt.Println("Tìm thấy:", name)
	}
}
```

## Bài tập

1. **Custom error**: viết `ValidationError{Field, Msg string}` như trên, viết hàm `validateAge(age int) error` trả về lỗi này khi `age < 0` hoặc `age > 150`.
2. **`errors.Is`**: định nghĩa `ErrNotFound = errors.New(...)`, viết hàm giả lập tra cứu (map) trả về lỗi wrap `ErrNotFound` khi không tìm thấy key, dùng `errors.Is` để phân biệt với lỗi khác.
3. **`recover`**: viết hàm chia 2 số dùng `defer`+`recover` để không crash khi chia cho 0, trả về `(0, error)` thay vì để panic thoát ra ngoài.
4. **Nâng cao**: viết hàm `processOrders(orders []int) []error` xử lý danh sách đơn hàng (số lượng), thu thập TẤT CẢ lỗi gặp phải (không dừng ở lỗi đầu tiên) vào 1 slice `[]error`, dùng `errors.Join` (Go 1.20+) để gộp lại nếu muốn trả về 1 lỗi duy nhất.

## Tiếp theo
→ [Bài 9: Package & Go Modules](./9_packages_modules.md)
