# Bài 7: Con trỏ (Pointers)

## Mục tiêu
- Hiểu `&` (lấy địa chỉ) và `*` (dereference / khai báo kiểu con trỏ).
- Phân biệt truyền theo giá trị (pass by value) vs truyền theo địa chỉ.
- Biết khi nào Go tự động dereference giúp bạn.

## 1. Cơ bản về con trỏ

```go
x := 10
p := &x // p là con trỏ trỏ tới địa chỉ của x, kiểu của p là *int

fmt.Println(x)  // 10
fmt.Println(p)  // địa chỉ bộ nhớ, vd 0xc0000140a0
fmt.Println(*p) // 10 — dereference: lấy giá trị TẠI địa chỉ p trỏ tới

*p = 20 // sửa giá trị thông qua con trỏ
fmt.Println(x) // 20 — x cũng thay đổi vì p trỏ trực tiếp tới x
```

**Go không có phép toán con trỏ (pointer arithmetic)** như C — không thể làm `p + 1` để "nhảy" tới ô nhớ kế tiếp. Điều này giúp Go an toàn hơn C rất nhiều.

## 2. Tại sao cần con trỏ? — Pass by value vs pass by reference

Mặc định, Go truyền **mọi tham số theo giá trị** (bản sao). Muốn hàm sửa được biến gốc của caller, phải truyền con trỏ:

```go
// Truyền theo giá trị — KHÔNG thay đổi được biến gốc
func incrementByValue(n int) {
	n++ // chỉ sửa bản sao cục bộ
}

// Truyền theo con trỏ — THAY ĐỔI được biến gốc
func incrementByPointer(n *int) {
	*n++ // sửa giá trị tại địa chỉ được trỏ tới
}

func main() {
	x := 5
	incrementByValue(x)
	fmt.Println(x) // vẫn là 5

	incrementByPointer(&x)
	fmt.Println(x) // 6
}
```

## 3. Con trỏ với struct

```go
type Point struct{ X, Y int }

func movePoint(p *Point, dx, dy int) {
	p.X += dx // Go tự động dereference: (*p).X tương đương p.X
	p.Y += dy
}

func main() {
	pt := Point{X: 1, Y: 1}
	movePoint(&pt, 5, 5)
	fmt.Println(pt) // {6 6}
}
```

Đây chính là lý do method với **pointer receiver** (Bài 6) có thể sửa struct gốc — về bản chất nó cũng chỉ là truyền con trỏ vào hàm.

## 4. `new()` — cấp phát và trả về con trỏ

```go
p := new(int)   // cấp phát 1 int, zero value = 0, p có kiểu *int
*p = 42
fmt.Println(*p) // 42

// Với struct, thường dùng &Struct{} thay vì new() vì linh hoạt hơn
account := &BankAccount{Balance: 100}
```

## 5. Con trỏ nil

```go
var p *int // zero value của con trỏ là nil
if p == nil {
	fmt.Println("p chưa trỏ tới đâu cả")
}

// *p khi p là nil sẽ gây PANIC (nil pointer dereference) — luôn kiểm tra nil trước khi dereference
```

## Ví dụ đầy đủ: Swap 2 số

```go
package main

import "fmt"

func swap(a, b *int) {
	*a, *b = *b, *a
}

type Point struct{ X, Y int }

func (p *Point) MoveBy(dx, dy int) {
	p.X += dx
	p.Y += dy
}

func main() {
	x, y := 1, 2
	fmt.Println("Trước:", x, y)
	swap(&x, &y)
	fmt.Println("Sau:", x, y)

	pt := Point{X: 0, Y: 0}
	pt.MoveBy(3, 4)
	fmt.Println(pt) // {3 4}
}
```

## Bài tập

1. **`swap`**: viết hàm `swap(a, b *int)` hoán đổi 2 số, gọi thử và in trước/sau để kiểm chứng.
2. **Sửa struct qua con trỏ**: viết hàm sửa field của struct thông qua pointer receiver (vd `MoveBy` ở trên); sau đó thử viết bản KHÔNG dùng pointer (value receiver) để so sánh — chứng minh bằng code rằng struct gốc không thay đổi.
3. **Giải thích bằng lời** (viết comment trong code): tại sao slice và map khi truyền vào hàm thì hàm vẫn sửa được "nội dung" dù không dùng con trỏ tường minh? (Gợi ý: slice/map bản thân đã chứa con trỏ tới dữ liệu bên trong, dù bản thân biến slice/map được truyền theo giá trị).
4. **Nâng cao**: viết hàm `safeDereference(p *int) int` trả về giá trị tại `p`, hoặc trả về `0` nếu `p == nil` thay vì để chương trình panic.

## Tiếp theo
→ [Bài 8: Xử lý lỗi (Error Handling)](./8_error_handling.md)
