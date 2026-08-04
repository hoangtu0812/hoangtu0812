# Bài 6: Struct & Method

## Mục tiêu
- Định nghĩa và khởi tạo struct.
- Gắn method vào struct (value receiver vs pointer receiver).
- Dùng embedding để tái sử dụng field/method (Go không có class/kế thừa).

## 1. Định nghĩa struct

```go
type Person struct {
	Name string
	Age  int
}

func main() {
	// Khởi tạo có tên field (khuyến khích, rõ ràng)
	p1 := Person{Name: "Ben", Age: 25}

	// Khởi tạo theo thứ tự field (dễ lỗi khi struct nhiều field, hạn chế dùng)
	p2 := Person{"Alice", 30}

	// Khởi tạo rỗng rồi gán sau
	var p3 Person
	p3.Name = "Bob"
	p3.Age = 28

	fmt.Println(p1, p2, p3)
}
```

## 2. Method — hàm gắn với một kiểu (receiver)

```go
type Rectangle struct {
	Width, Height float64
}

// Value receiver: nhận BẢN SAO của Rectangle, không sửa được struct gốc
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Pointer receiver: nhận ĐỊA CHỈ, có thể sửa struct gốc
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func main() {
	rect := Rectangle{Width: 10, Height: 5}
	fmt.Println(rect.Area())      // 50
	fmt.Println(rect.Perimeter()) // 30

	rect.Scale(2) // Go tự động lấy &rect khi gọi method có pointer receiver
	fmt.Println(rect.Area()) // 200 — rect đã bị thay đổi thật sự
}
```

## 3. Khi nào dùng value receiver, khi nào dùng pointer receiver?

Quy tắc thực tế:
- Dùng **pointer receiver** khi method cần **sửa đổi** dữ liệu của struct, hoặc struct **lớn** (tránh copy tốn kém), hoặc để đồng nhất với các method khác cùng kiểu đã dùng pointer.
- Dùng **value receiver** khi struct nhỏ, bất biến (immutable behavior mong muốn), hoặc để an toàn tuyệt đối không có side-effect.
- **Không trộn lẫn** tùy tiện trong cùng 1 kiểu — nếu 1 method cần pointer receiver, các method khác của kiểu đó cũng nên nhất quán dùng pointer receiver.

## 4. Embedding — cách Go "tái sử dụng" code thay vì kế thừa

Go không có `class`/`extends`. Thay vào đó, struct có thể **nhúng (embed)** struct khác, tự động "thừa hưởng" field và method:

```go
type Animal struct {
	Name string
}

func (a Animal) Describe() string {
	return "Tôi là " + a.Name
}

type Dog struct {
	Animal // embedding — không đặt tên field, chỉ ghi kiểu
	Breed  string
}

func main() {
	d := Dog{
		Animal: Animal{Name: "Rex"},
		Breed:  "Golden Retriever",
	}

	fmt.Println(d.Name)        // truy cập trực tiếp field của Animal
	fmt.Println(d.Describe())  // gọi trực tiếp method của Animal
	fmt.Println(d.Breed)
}
```

Đây gọi là "composition" (hợp thành) thay vì "inheritance" (kế thừa) — triết lý Go: **"Accept interfaces, return structs. Favor composition over inheritance."**

Nếu `Dog` muốn override `Describe()`, chỉ cần định nghĩa method `Describe()` riêng cho `Dog` — nó sẽ che (shadow) method của `Animal`.

## 5. So sánh struct

```go
type Point struct{ X, Y int }

p1 := Point{1, 2}
p2 := Point{1, 2}
fmt.Println(p1 == p2) // true — struct so sánh được bằng == nếu MỌI field đều comparable
```

## Ví dụ đầy đủ

```go
package main

import "fmt"

type Animal struct {
	Name string
}

func (a Animal) Describe() string {
	return fmt.Sprintf("%s là một loài động vật", a.Name)
}

type Dog struct {
	Animal
	Breed string
}

func (d Dog) Describe() string {
	return fmt.Sprintf("%s là chó giống %s", d.Name, d.Breed)
}

type Cat struct {
	Animal
	Indoor bool
}

func main() {
	dog := Dog{Animal: Animal{Name: "Rex"}, Breed: "Poodle"}
	cat := Cat{Animal: Animal{Name: "Mimi"}, Indoor: true}

	fmt.Println(dog.Describe())  // override: "Rex là chó giống Poodle"
	fmt.Println(cat.Describe())  // dùng của Animal: "Mimi là một loài động vật"
}
```

## Bài tập

1. **`Rectangle`**: viết struct `Rectangle{Width, Height float64}` với method `Area()`, `Perimeter()` (value receiver).
2. **Embedding `Animal`**: viết `Animal` với method `Describe()`, nhúng vào `Dog` và `Cat`. Cho `Dog` override `Describe()`, để `Cat` dùng bản gốc của `Animal`.
3. **Giải thích pointer receiver**: viết 2 phiên bản method `Scale` — 1 dùng value receiver, 1 dùng pointer receiver. Gọi cả 2 và in kết quả để chứng minh bằng code: value receiver KHÔNG làm thay đổi struct gốc, pointer receiver thì có.
4. **Nâng cao**: viết struct `BankAccount{Balance float64}` với method `Deposit(amount float64)` và `Withdraw(amount float64) error` (trả lỗi nếu rút quá số dư) — cả 2 đều phải dùng pointer receiver vì có thay đổi state.

## Tiếp theo
→ [Bài 7: Con trỏ](./7_pointers.md)
