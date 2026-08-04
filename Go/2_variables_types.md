# Bài 2: Biến, Kiểu dữ liệu, Hằng số

## Mục tiêu
- Khai báo biến bằng `var` và `:=`, hiểu zero value.
- Nắm các kiểu dữ liệu cơ bản của Go.
- Dùng `const` và `iota` để tạo hằng số / enum.

## 1. Khai báo biến

Go có 2 cách khai báo biến chính:

```go
package main

import "fmt"

func main() {
	// Cách 1: dùng var, có thể khai báo kiểu tường minh
	var name string = "Ben"
	var age int = 25

	// Cách 2: var + suy luận kiểu (type inference)
	var city = "Quang Ngai"

	// Cách 3: short declaration (chỉ dùng được BÊN TRONG hàm)
	country := "Vietnam"

	fmt.Println(name, age, city, country)
}
```

**Lưu ý:**
- `:=` chỉ dùng được trong function body, không dùng được ở package scope.
- Biến khai báo mà không dùng tới sẽ bị lỗi biên dịch (`declared and not used`) — đây là chủ đích thiết kế của Go để tránh code rác.
- Có thể khai báo nhiều biến cùng lúc: `var a, b, c int` hoặc `x, y := 1, "two"`.

## 2. Zero value

Khi khai báo biến bằng `var` mà không gán giá trị, Go tự gán "zero value" theo kiểu:

| Kiểu      | Zero value |
|-----------|------------|
| int/float | `0`        |
| string    | `""`       |
| bool      | `false`    |
| pointer/slice/map/interface/channel/func | `nil` |

```go
var count int      // 0
var price float64  // 0
var label string   // ""
var active bool    // false
fmt.Println(count, price, label, active)
```

## 3. Các kiểu dữ liệu cơ bản

- **Số nguyên**: `int`, `int8/16/32/64`, `uint`, `uint8` (alias `byte`)... `int` có độ rộng phụ thuộc kiến trúc máy (32/64 bit) — dùng `int` mặc định trừ khi có lý do cụ thể để chọn kích thước khác.
- **Số thực**: `float32`, `float64` (mặc định nên dùng `float64`).
- **Chuỗi**: `string` — immutable, là dãy byte UTF-8.
- **Ký tự**: `rune` (alias `int32`, đại diện 1 Unicode code point), `byte` (alias `uint8`).
- **Boolean**: `bool` (`true`/`false`).

```go
var r rune = 'A'      // rune lưu Unicode code point, 'A' = 65
var b byte = 'A'      // byte lưu giá trị 0-255
fmt.Printf("%c %d\n", r, r) // A 65
```

## 4. Kiểm tra kiểu bằng `%T`

```go
x := 42
y := 3.14
z := "hello"
fmt.Printf("%v is %T\n", x, x) // 42 is int
fmt.Printf("%v is %T\n", y, y) // 3.14 is float64
fmt.Printf("%v is %T\n", z, z) // hello is string
```

## 5. Ép kiểu (Type conversion)

Go **không** tự động ép kiểu ngầm định (khác C/JS). Phải convert tường minh:

```go
var i int = 10
var f float64 = float64(i) // bắt buộc convert tường minh
var u uint = uint(f)
```

## 6. Hằng số (`const`) và `iota`

```go
const Pi = 3.14159 // hằng số không thể thay đổi giá trị

// iota: bộ đếm tự tăng, dùng để tạo enum
type Weekday int

const (
	Sunday Weekday = iota // 0
	Monday                // 1
	Tuesday               // 2
	Wednesday             // 3
	Thursday              // 4
	Friday                // 5
	Saturday              // 6
)

func (d Weekday) String() string {
	names := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	return names[d]
}
```

## Ví dụ đầy đủ: Chuyển đổi độ C ↔ F

```go
package main

import "fmt"

func celsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

func fahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func main() {
	var celsius float64 = 30
	fahrenheit := celsiusToFahrenheit(celsius)
	fmt.Printf("%.1f°C = %.1f°F\n", celsius, fahrenheit)

	f := 98.6
	c := fahrenheitToCelsius(f)
	fmt.Printf("%.1f°F = %.1f°C\n", f, c)
}
```

## Bài tập

1. **Khai báo & in kiểu**: Khai báo ít nhất 5 biến với `var` và `:=` (đủ int, float64, string, bool, rune), in ra giá trị và kiểu (`%v`, `%T`) của từng biến.
2. **Đổi độ C ↔ F**: Viết chương trình đọc một giá trị nhiệt độ (hardcode hoặc `fmt.Scanln`) và in ra cả hai chiều quy đổi (dùng code mẫu ở trên làm nền, tự viết lại không copy).
3. **Enum bằng iota**: Dùng `iota` tạo enum `Weekday` (Chủ nhật → Thứ 7), viết hàm `String()` để in tên ngày thay vì số. In thử `fmt.Println(Wednesday)`.
4. **Thử lỗi biên dịch**: Khai báo một biến rồi KHÔNG dùng tới, chạy `go build` để tự quan sát lỗi `declared and not used` — giải thích bằng lời tại sao Go thiết kế như vậy.

### Gợi ý lời giải bài 4
Go buộc mọi biến khai báo phải được dùng để tránh code thừa/rác tích tụ theo thời gian, giúp code dễ đọc và ít bug hơn khi refactor. Nếu cố tình bỏ qua, dùng `_ = bien` để "dùng giả" biến đó (thường chỉ nên dùng khi debug tạm thời).

## Tiếp theo
→ [Bài 3: Toán tử & luồng điều khiển](./3_control_flow.md)
