# Bài 15: Generics (Go 1.18+)

## Mục tiêu
- Viết hàm/kiểu generic với type parameter.
- Dùng constraint (`any`, `comparable`, constraint tự định nghĩa).
- Biết khi nào NÊN và KHÔNG NÊN dùng generics.

## 1. Vấn đề trước khi có generics

Trước Go 1.18, để viết hàm `Max` dùng được cho cả `int` và `float64`, bạn phải viết trùng lặp hoặc dùng `interface{}` (mất type safety, cần type assertion). Generics giải quyết vấn đề này.

## 2. Hàm generic cơ bản

```go
// T là type parameter, ràng buộc bởi constraint (ở đây dùng constraint có sẵn)
func Max[T int | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println(Max(3, 5))         // T được suy luận là int
	fmt.Println(Max(3.5, 2.1))     // T được suy luận là float64
	fmt.Println(Max[float64](3, 5)) // chỉ định tường minh kiểu nếu cần
}
```

## 3. Constraint có sẵn trong package `constraints` / built-in

Từ Go 1.21, package chuẩn `cmp` cung cấp `cmp.Ordered`:

```go
import "cmp"

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
```

`comparable` là constraint built-in cho các kiểu dùng được với `==`/`!=` (dùng khi cần so sánh bằng, vd key của map):

```go
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println(Contains([]int{1, 2, 3}, 2))          // true
	fmt.Println(Contains([]string{"a", "b"}, "c"))    // false
}
```

## 4. Tự định nghĩa constraint

```go
type Number interface {
	int | int64 | float32 | float64
}

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	fmt.Println(Sum([]int{1, 2, 3}))          // 6
	fmt.Println(Sum([]float64{1.5, 2.5}))     // 4.0
}
```

## 5. Generic với slice — `Map` và `Filter` (rất hay dùng thực tế)

```go
func Map[T, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func Filter[T any](s []T, predicate func(T) bool) []T {
	result := []T{}
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6}

	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Println(doubled) // [2 4 6 8 10 12]

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println(evens) // [2 4 6]

	// Map có thể đổi kiểu: []int -> []string
	strs := Map(nums, func(n int) string { return fmt.Sprintf("số %d", n) })
	fmt.Println(strs)
}
```

Go 1.21+ đã có sẵn các hàm tương tự trong package chuẩn `slices` và `maps` (`slices.Contains`, `slices.Sort`...) — nên ưu tiên dùng thư viện chuẩn trước khi tự viết lại.

## 6. Generic struct

```go
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	last := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return last, true
}

func main() {
	intStack := &Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	v, _ := intStack.Pop()
	fmt.Println(v) // 2

	stringStack := &Stack[string]{}
	stringStack.Push("hello")
}
```

## 7. Khi nào NÊN và KHÔNG NÊN dùng generics

**Nên dùng** khi:
- Viết cấu trúc dữ liệu tái sử dụng (Stack, Queue, LinkedList...) hoạt động với nhiều kiểu.
- Viết hàm tiện ích thao tác slice/map chung chung (`Map`, `Filter`, `Reduce`).

**Không nên dùng** khi:
- Chỉ có 1-2 kiểu cụ thể cần xử lý — viết hàm riêng cho từng kiểu thường rõ ràng, dễ đọc hơn.
- Logic nghiệp vụ khác nhau đáng kể giữa các kiểu (generics chỉ hợp khi logic **giống hệt nhau**, chỉ khác kiểu dữ liệu).

## Bài tập

1. **`Max` generic**: viết hàm `Max[T cmp.Ordered](a, b T) T`, thử với `int`, `float64`, `string`.
2. **`Map` và `Filter`**: viết 2 hàm generic này như ví dụ trên, dùng thử để: nhân đôi 1 slice số, và lọc ra các số chẵn.
3. **`Stack[T]` generic**: viết struct `Stack[T any]` với `Push`/`Pop`, test với cả `Stack[int]` và `Stack[string]`.
4. **Nâng cao**: viết hàm generic `Reduce[T, U any](s []T, initial U, f func(U, T) U) U`, dùng nó để tính tổng 1 slice số và để nối 1 slice string thành 1 câu.

## Tổng kết Giai đoạn 2
Bạn đã nắm interface, concurrency (goroutine/channel), context, JSON/file I/O, testing, và generics — đây là bộ công cụ đầy đủ để đọc hiểu hầu hết code Go trong thực tế. Giai đoạn 3 sẽ tập trung vào cách **tổ chức** các công cụ này thành 1 ứng dụng thực tế.

## Tiếp theo
→ [Bài 16: Thiết kế package & Clean Architecture cơ bản](./16_clean_architecture.md)
