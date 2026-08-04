# Bài 10: Interface

## Mục tiêu
- Hiểu interface "implicit" của Go (không cần khai báo `implements`).
- Dùng `interface{}`/`any`, type assertion, type switch.
- Hiểu vì sao interface là công cụ thiết kế quan trọng nhất của Go.

## 1. Interface là gì?

Interface định nghĩa một tập hợp method mà một kiểu phải có để "thỏa mãn" interface đó. Khác Java/C#, Go **không cần** khai báo tường minh "class X implements InterfaceY" — chỉ cần kiểu có đủ method là tự động thỏa mãn interface (gọi là **structural typing** / **duck typing tĩnh**).

```go
type Shape interface {
	Area() float64
}

type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64 { return r.Width * r.Height }

type Circle struct{ Radius float64 }

func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }

// Rectangle và Circle đều tự động thỏa mãn Shape — không cần khai báo gì thêm
func PrintArea(s Shape) {
	fmt.Printf("Diện tích: %.2f\n", s.Area())
}

func main() {
	PrintArea(Rectangle{Width: 10, Height: 5})
	PrintArea(Circle{Radius: 3})
}
```

## 2. Vì sao interface quan trọng?

Interface cho phép viết code **không phụ thuộc vào implementation cụ thể** — hàm `PrintArea` không cần biết gì về `Rectangle` hay `Circle`, chỉ cần biết đối số có method `Area()`. Đây là nền tảng để:
- Viết unit test dễ dàng (mock interface thay vì implementation thật — vd mock database).
- Thiết kế theo layer (handler phụ thuộc vào interface `UserRepository`, không phụ thuộc trực tiếp Postgres) — sẽ áp dụng ở [Bài 16](./16_clean_architecture.md) và [Bài 21](./21_capstone_project.md).

**Nguyên tắc Go**: *"Accept interfaces, return structs"* — hàm nên nhận interface làm tham số (linh hoạt cho caller), nhưng trả về struct cụ thể (rõ ràng cho người gọi).

## 3. Interface rỗng: `interface{}` / `any`

```go
// any là alias của interface{} từ Go 1.18, dùng khi cần "chứa bất kỳ kiểu gì"
func describe(v any) {
	fmt.Printf("Giá trị: %v, Kiểu: %T\n", v, v)
}

func main() {
	describe(42)
	describe("hello")
	describe(3.14)
	describe(Rectangle{Width: 1, Height: 2})
}
```

Lưu ý: lạm dụng `any` làm mất đi lợi ích của static typing — chỉ dùng khi thực sự cần (vd hàm generic trước Go 1.18, hoặc xử lý dữ liệu JSON không rõ schema).

## 4. Type assertion

Khi có giá trị kiểu interface, muốn lấy lại kiểu cụ thể bên trong, dùng type assertion:

```go
var s Shape = Rectangle{Width: 4, Height: 5}

r, ok := s.(Rectangle) // "comma ok" — an toàn, không panic nếu sai kiểu
if ok {
	fmt.Println("Đây là Rectangle với Width =", r.Width)
}

// Dạng không kiểm tra ok — PANIC nếu assertion sai, chỉ dùng khi CHẮC CHẮN về kiểu
r2 := s.(Rectangle)
```

## 5. Type switch

Khi cần xử lý nhiều kiểu cụ thể khác nhau từ 1 giá trị interface:

```go
func describeShape(s Shape) string {
	switch v := s.(type) {
	case Rectangle:
		return fmt.Sprintf("Hình chữ nhật %vx%v", v.Width, v.Height)
	case Circle:
		return fmt.Sprintf("Hình tròn bán kính %v", v.Radius)
	default:
		return "Hình dạng không xác định"
	}
}
```

## 6. Interface chuẩn hay gặp trong thư viện chuẩn Go

```go
// fmt.Stringer — bất kỳ kiểu nào có String() string sẽ tự động dùng khi Print
type Stringer interface {
	String() string
}

// error — đã học ở Bài 8
type error interface {
	Error() string
}

// io.Writer — nền tảng cho ghi file, HTTP response, log...
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

## Ví dụ đầy đủ

```go
package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64 { return r.Width * r.Height }

type Circle struct{ Radius float64 }

func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }

func PrintArea(s Shape) {
	fmt.Printf("Diện tích: %.2f\n", s.Area())
}

func describe(v any) {
	switch val := v.(type) {
	case int:
		fmt.Println("Số nguyên:", val)
	case string:
		fmt.Println("Chuỗi:", val)
	case Shape:
		fmt.Println("Hình dạng có diện tích:", val.Area())
	default:
		fmt.Println("Kiểu không xác định:", val)
	}
}

func main() {
	shapes := []Shape{
		Rectangle{Width: 10, Height: 5},
		Circle{Radius: 3},
	}
	for _, s := range shapes {
		PrintArea(s)
	}

	describe(42)
	describe("hello")
	describe(Rectangle{Width: 1, Height: 1})
}
```

## Bài tập

1. **`Shape` interface**: viết interface `Shape` với `Area()`, implement cho `Rectangle` và `Circle`, viết hàm `PrintArea(s Shape)` dùng chung cho cả 2.
2. **Type switch**: viết hàm `describe(v any)` xử lý ít nhất 4 case khác nhau (int, string, bool, và 1 struct tự định nghĩa) bằng type switch.
3. **`Stringer`**: cho struct `Rectangle` implement interface `fmt.Stringer` (viết method `String() string`), sau đó `fmt.Println(rect)` sẽ tự động dùng `String()` thay vì in mặc định — kiểm chứng bằng code.
4. **Nâng cao**: viết interface `Notifier` với method `Notify(message string) error`, tạo 2 implementation giả `EmailNotifier` và `SMSNotifier` (chỉ cần `fmt.Println` giả lập gửi), viết hàm `SendAlert(n Notifier, msg string)` dùng chung — đây chính là ý tưởng nền tảng cho dependency injection sẽ dùng ở dự án capstone.

## Tiếp theo
→ [Bài 11: Goroutines & Channels](./11_goroutines_channels.md)
