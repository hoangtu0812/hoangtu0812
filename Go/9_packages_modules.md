# Bài 9: Package & Go Modules

## Mục tiêu
- Hiểu cách Go tổ chức code theo package.
- Phân biệt exported (public) vs unexported (private) qua quy tắc viết hoa.
- Dùng `go mod init`/`go mod tidy` để quản lý dependency.
- Tổ chức nhiều file trong cùng 1 package.

## 1. Package là gì?

Mỗi thư mục Go = 1 package (trừ thư mục gốc chứa `main`). Mọi file `.go` trong cùng thư mục phải khai báo cùng tên package ở dòng đầu tiên:

```go
package shapes
```

`package main` là đặc biệt: đánh dấu đây là chương trình thực thi được (phải có hàm `func main()`), khác với package thư viện thông thường.

## 2. Exported vs unexported

Go dùng **chữ cái đầu viết hoa hay thường** để quyết định phạm vi truy cập — không có từ khóa `public`/`private`:

```go
package shapes

// Rectangle: EXPORTED (chữ R hoa) — dùng được từ package khác
type Rectangle struct {
	Width, Height float64
}

// area: unexported (chữ a thường) — chỉ dùng được trong package shapes
func area(r Rectangle) float64 {
	return r.Width * r.Height
}

// Area: EXPORTED — public API của package
func (r Rectangle) Area() float64 {
	return area(r)
}
```

Quy tắc thiết kế: chỉ export những gì thực sự cần dùng từ bên ngoài — giữ phần còn lại unexported để dễ thay đổi implementation sau này mà không phá vỡ code của người dùng package.

## 3. Go Modules — quản lý dependency

```powershell
go mod init example.com/myproject   # tạo go.mod, khai báo module path
go get github.com/gin-gonic/gin     # thêm dependency, cập nhật go.mod + go.sum
go mod tidy                          # dọn dẹp: thêm dependency thiếu, xóa dependency thừa
```

File `go.mod` ví dụ:

```
module example.com/myproject

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
)
```

`go.sum` chứa checksum để đảm bảo tính toàn vẹn của dependency — **không sửa tay**, để `go` tool tự quản lý.

## 4. Import package

```go
import (
	"fmt"                          // thư viện chuẩn
	"example.com/myproject/shapes" // package nội bộ trong cùng module
	"github.com/gin-gonic/gin"     // thư viện bên thứ 3
)
```

## 5. Tổ chức nhiều file trong 1 package

Một package có thể (và nên) được chia thành nhiều file nhỏ theo chức năng — chúng vẫn "thấy" nhau như đang ở cùng 1 file:

```
shapes/
├── rectangle.go   // package shapes
├── circle.go       // package shapes
└── shapes_test.go  // package shapes (hoặc shapes_test cho black-box test)
```

```go
// rectangle.go
package shapes

type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64 { return r.Width * r.Height }
```

```go
// circle.go
package shapes

import "math"

type Circle struct{ Radius float64 }

func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }
```

## Ví dụ đầy đủ: cấu trúc project 2 package

```
myproject/
├── go.mod
├── main.go
└── shapes/
    └── shapes.go
```

```go
// go.mod
module example.com/myproject

go 1.22
```

```go
// shapes/shapes.go
package shapes

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}
```

```go
// main.go
package main

import (
	"fmt"

	"example.com/myproject/shapes"
)

func main() {
	r := shapes.Rectangle{Width: 10, Height: 5}
	fmt.Println("Diện tích:", r.Area())
	fmt.Println("Chu vi:", r.Perimeter())
}
```

Chạy: `go run .`

## Bài tập

1. **Tách package `shapes`**: lấy struct `Rectangle` (Bài 6) ra thành package riêng như ví dụ trên, import và dùng từ `main`.
2. **Thêm `Circle`**: bổ sung `Circle{Radius float64}` vào package `shapes` (file riêng `circle.go`), với method `Area()`.
3. **Thử unexported**: thêm 1 hàm helper không viết hoa chữ cái đầu trong package `shapes`, thử import và gọi hàm đó từ `main` — quan sát lỗi biên tập/biên dịch, giải thích tại sao.
4. **`go mod tidy`**: thêm 1 thư viện bên thứ 3 bất kỳ (vd `github.com/google/uuid`) vào project, chạy `go get` rồi `go mod tidy`, quan sát thay đổi trong `go.mod`/`go.sum`.

## Tổng kết Giai đoạn 1
Bạn đã hoàn thành phần Go cơ bản: biến/kiểu dữ liệu, control flow, hàm, slice/map, struct, con trỏ, error handling, và package. Đây là nền tảng đủ vững để bắt đầu Giai đoạn 2 — các khái niệm khiến Go trở nên đặc biệt: **interface** và **concurrency**.

## Tiếp theo
→ [Bài 10: Interface](./10_interfaces.md)
