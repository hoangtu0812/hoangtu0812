# Bài 3: Toán tử & Luồng điều khiển

## Mục tiêu
- Sử dụng thành thạo toán tử số học/so sánh/logic.
- Viết `if/else`, `switch`, các dạng `for` trong Go.
- Hiểu `break`, `continue`, nhãn (label), và vì sao Go không có `while`/`do-while` riêng.

## 1. Toán tử

```go
a, b := 10, 3

// Số học
fmt.Println(a+b, a-b, a*b, a/b, a%b) // 13 7 30 3 1

// So sánh
fmt.Println(a > b, a == b, a != b) // true false true

// Logic
isAdult := true
hasID := false
fmt.Println(isAdult && hasID, isAdult || hasID, !isAdult) // false true false
```

Lưu ý: Go **không có** toán tử `++`/`--` dùng như biểu thức (chỉ dùng như statement độc lập: `a++`), và không hỗ trợ toán tử ba ngôi `? :`.

## 2. `if` / `else`

```go
score := 75

if score >= 90 {
	fmt.Println("Xuất sắc")
} else if score >= 70 {
	fmt.Println("Khá")
} else {
	fmt.Println("Cần cố gắng")
}

// if với statement khởi tạo (rất phổ biến trong Go, đặc biệt với error)
if v, err := someFunc(); err != nil {
	fmt.Println("Lỗi:", err)
} else {
	fmt.Println("Giá trị:", v)
}
```

`v` và `err` chỉ tồn tại trong phạm vi khối `if/else` này — đây là idiom rất hay gặp khi gọi hàm trả về `(value, error)`.

## 3. `switch`

```go
day := 3

switch day {
case 1, 7:
	fmt.Println("Cuối tuần")
case 2, 3, 4, 5, 6:
	fmt.Println("Ngày thường")
default:
	fmt.Println("Không hợp lệ")
}

// switch không điều kiện = if/else if dài thay thế gọn hơn
switch {
case score >= 90:
	fmt.Println("Xuất sắc")
case score >= 70:
	fmt.Println("Khá")
default:
	fmt.Println("Cần cố gắng")
}
```

Khác với C/Java, Go **tự động break** sau mỗi `case` (không rơi xuống case tiếp theo). Muốn "rơi xuống" thì dùng từ khóa `fallthrough`.

## 4. Vòng lặp `for` (Go chỉ có `for`, không có `while`/`do-while`)

```go
// Dạng đầy đủ (giống C)
for i := 0; i < 5; i++ {
	fmt.Println(i)
}

// Dạng "while" — chỉ có điều kiện
n := 0
for n < 5 {
	fmt.Println(n)
	n++
}

// Vòng lặp vô hạn
for {
	fmt.Println("Chạy mãi...")
	break // phải có điều kiện dừng
}

// Duyệt qua slice/map bằng range
nums := []int{10, 20, 30}
for index, value := range nums {
	fmt.Println(index, value)
}

// Chỉ cần index
for i := range nums {
	fmt.Println(i)
}

// Chỉ cần value, bỏ qua index bằng _
for _, v := range nums {
	fmt.Println(v)
}
```

## 5. `break`, `continue`, và nhãn (label)

```go
// break/continue với vòng lặp lồng nhau cần label để thoát đúng vòng ngoài
outer:
for i := 0; i < 3; i++ {
	for j := 0; j < 3; j++ {
		if j == 1 {
			continue outer // bỏ qua vòng i hiện tại, sang i tiếp theo
		}
		if i == 2 {
			break outer // thoát hẳn cả 2 vòng lặp
		}
		fmt.Println(i, j)
	}
}
```

## Ví dụ đầy đủ: FizzBuzz

```go
package main

import "fmt"

func main() {
	for i := 1; i <= 100; i++ {
		switch {
		case i%15 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}
}
```

## Bài tập

1. **FizzBuzz 1–100**: dùng code mẫu trên làm nền, tự gõ lại (không copy-paste).
2. **Kiểm tra số nguyên tố**: viết hàm `isPrime(n int) bool`, in ra tất cả số nguyên tố từ 2 đến 100.
3. **Bảng cửu chương**: dùng 2 vòng `for` lồng nhau, in bảng cửu chương từ 1 đến 9 (mỗi bảng cách nhau 1 dòng trống).
4. **Đoán số (nâng cao)**: cho trước số bí mật `secret := 42`, dùng vòng lặp giả lập việc "đoán dần" từ danh sách `guesses := []int{10, 50, 42, 100}`, dùng `break` với label để dừng ngay khi đoán đúng và in ra số lần đoán.

### Gợi ý bài 2
```go
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
```

## Tiếp theo
→ [Bài 4: Hàm (Functions)](./4_functions.md)
