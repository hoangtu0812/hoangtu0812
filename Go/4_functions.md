# Bài 4: Hàm (Functions)

## Mục tiêu
- Khai báo hàm, hiểu multiple return values và named return.
- Dùng biến đối số (variadic), closure, hàm ẩn danh.
- Hiểu `defer` và thứ tự thực thi của nó.

## 1. Khai báo hàm cơ bản

```go
func add(a int, b int) int {
	return a + b
}

// Rút gọn khi 2 tham số cùng kiểu
func subtract(a, b int) int {
	return a - b
}
```

## 2. Nhiều giá trị trả về (multiple return values)

Đây là idiom cực kỳ phổ biến trong Go, đặc biệt với cặp `(value, error)`:

```go
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("không thể chia cho 0")
	}
	return a / b, nil
}

func main() {
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Lỗi:", err)
		return
	}
	fmt.Println("Kết quả:", result)
}
```

## 3. Named return values

```go
// result và err được khai báo sẵn, "return" trần sẽ trả về giá trị hiện tại của chúng
func divide(a, b int) (result int, err error) {
	if b == 0 {
		err = fmt.Errorf("không thể chia cho 0")
		return // tương đương return result, err
	}
	result = a / b
	return
}
```

Named return giúp code rõ ràng hơn khi hàm có nhiều giá trị trả về, nhưng dùng quá nhiều có thể làm giảm độ rõ ràng — chỉ nên dùng khi tên có ý nghĩa bổ sung ngữ cảnh.

## 4. Biến đối số (Variadic functions)

```go
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	fmt.Println(sum(1, 2, 3))       // 6
	fmt.Println(sum())              // 0
	nums := []int{1, 2, 3, 4}
	fmt.Println(sum(nums...))       // "mở" slice ra thành nhiều đối số bằng ...
}
```

## 5. Hàm là first-class citizen: closure & hàm ẩn danh

Hàm trong Go có thể được gán cho biến, truyền như tham số, trả về từ hàm khác:

```go
// Hàm ẩn danh gán cho biến
square := func(x int) int {
	return x * x
}
fmt.Println(square(5)) // 25

// Closure: hàm "nhớ" biến trong scope bao quanh nó
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	counter := makeCounter()
	fmt.Println(counter()) // 1
	fmt.Println(counter()) // 2
	fmt.Println(counter()) // 3

	counter2 := makeCounter() // instance độc lập, có count riêng
	fmt.Println(counter2())   // 1
}
```

## 6. `defer`

`defer` hoãn thực thi một lời gọi hàm cho đến khi hàm chứa nó return. Dùng nhiều để dọn dẹp tài nguyên (đóng file, đóng kết nối DB...).

```go
func readFile() {
	fmt.Println("Mở file")
	defer fmt.Println("Đóng file") // sẽ chạy CUỐI CÙNG, dù có return sớm hay panic
	fmt.Println("Đọc dữ liệu")
}
// Output:
// Mở file
// Đọc dữ liệu
// Đóng file
```

Nhiều `defer` chạy theo thứ tự **LIFO** (vào sau ra trước):

```go
func main() {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")
}
// Output: 3 2 1
```

## Ví dụ đầy đủ

```go
package main

import "fmt"

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("chia cho 0 không hợp lệ")
	}
	return a / b, nil
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	if result, err := divide(10, 0); err != nil {
		fmt.Println("Lỗi:", err)
	} else {
		fmt.Println("Kết quả:", result)
	}

	fmt.Println("Tổng:", sum(1, 2, 3, 4, 5))

	counter := makeCounter()
	for i := 0; i < 3; i++ {
		fmt.Println("Đếm:", counter())
	}
}
```

## Bài tập

1. **`divide` an toàn**: viết hàm `divide(a, b int) (int, error)` trả lỗi khi `b == 0`, gọi thử với cả 2 trường hợp hợp lệ/không hợp lệ.
2. **`sum` biến đối số**: viết hàm `sum(nums ...int) int`, thử gọi cả với danh sách số rời rạc và với 1 slice có sẵn (dùng `...`).
3. **Closure counter**: viết `makeCounter()` như trên, tạo 2 counter độc lập, chứng minh chúng không ảnh hưởng lẫn nhau.
4. **`defer` thứ tự**: viết hàm có 3 lệnh `defer` in ra số 1, 2, 3 — dự đoán trước thứ tự output rồi chạy để kiểm tra.
5. **Nâng cao**: viết hàm `applyTwice(f func(int) int, x int) int` áp dụng hàm `f` lên `x` hai lần liên tiếp (vd `applyTwice(square, 3)` → `81`).

## Tiếp theo
→ [Bài 5: Mảng, Slice, Map](./5_slices_maps.md)
