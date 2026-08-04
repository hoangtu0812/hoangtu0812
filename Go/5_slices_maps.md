# Bài 5: Mảng, Slice, Map

## Mục tiêu
- Phân biệt array (cố định) và slice (động).
- Dùng thành thạo `append`, `make`, `len`, `cap`.
- Hiểu vì sao slice chia sẻ underlying array (side-effect khi truyền slice).
- CRUD với map, kiểm tra key tồn tại.

## 1. Array (mảng cố định)

```go
var arr [5]int              // [0 0 0 0 0], độ dài CỐ ĐỊNH là 5
arr[0] = 10
nums := [3]int{1, 2, 3}     // khai báo kèm giá trị
letters := [...]string{"a", "b", "c"} // độ dài suy ra từ số phần tử = 3
```

Array trong Go **hiếm khi dùng trực tiếp** trong code thực tế — vì kích thước cố định là 1 phần của kiểu (`[5]int` và `[3]int` là 2 kiểu khác nhau). Thực tế người ta dùng **slice** gần như luôn luôn.

## 2. Slice (mảng động)

Slice là "view" linh động trên một underlying array, gồm 3 thành phần: con trỏ tới dữ liệu, độ dài (`len`), sức chứa (`cap`).

```go
// Khai báo slice literal
nums := []int{1, 2, 3}

// Tạo bằng make(kiểu, len, cap)
s := make([]int, 3)      // [0 0 0], len=3, cap=3
s2 := make([]int, 0, 10) // len=0, cap=10 (tối ưu khi biết trước sẽ append nhiều)

// append: thêm phần tử, tự động cấp phát lại nếu vượt cap
nums = append(nums, 4)          // [1 2 3 4]
nums = append(nums, 5, 6, 7)    // append nhiều giá trị cùng lúc
other := []int{8, 9}
nums = append(nums, other...)   // nối 2 slice, dùng ...

fmt.Println(len(nums), cap(nums))
```

**Quan trọng:** `append` có thể trả về slice mới nếu vượt quá `cap` hiện tại (Go cấp phát mảng mới, copy dữ liệu qua) — vì vậy LUÔN gán lại kết quả: `nums = append(nums, x)`, không bao giờ gọi `append` mà bỏ qua giá trị trả về.

## 3. Slice của slice — chia sẻ underlying array

```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:3] // [2, 3] — nhưng CHIA SẺ cùng mảng gốc với original!

sub[0] = 99
fmt.Println(original) // [1 99 3 4 5] — original cũng bị thay đổi!
```

Đây là điểm dễ gây bug nhất cho người mới: sửa phần tử của slice con có thể ảnh hưởng slice cha nếu chúng còn share underlying array. Muốn tách biệt hoàn toàn, dùng `copy()`:

```go
independent := make([]int, len(sub))
copy(independent, sub) // copy dữ liệu thật sự, không share array nữa
```

## 4. Map

```go
// Khai báo
scores := map[string]int{
	"Alice": 90,
	"Bob":   85,
}

// Hoặc tạo rỗng bằng make
ages := make(map[string]int)
ages["Ben"] = 25

// Đọc giá trị
fmt.Println(scores["Alice"]) // 90
fmt.Println(scores["Carol"]) // 0 (zero value của int, KHÔNG panic dù key không tồn tại)

// Kiểm tra key có tồn tại hay không — idiom "comma ok"
value, ok := scores["Carol"]
if !ok {
	fmt.Println("Carol chưa có điểm")
} else {
	fmt.Println(value)
}

// Xóa key
delete(scores, "Bob")

// Duyệt map (LƯU Ý: thứ tự duyệt map trong Go là NGẪU NHIÊN, không đảm bảo thứ tự)
for name, score := range scores {
	fmt.Println(name, score)
}
```

## Ví dụ đầy đủ: Đếm tần suất từ

```go
package main

import (
	"fmt"
	"strings"
)

func wordFrequency(sentence string) map[string]int {
	words := strings.Fields(strings.ToLower(sentence))
	freq := make(map[string]int)
	for _, w := range words {
		freq[w]++
	}
	return freq
}

func main() {
	sentence := "the quick brown fox jumps over the lazy dog the fox runs"
	freq := wordFrequency(sentence)
	for word, count := range freq {
		fmt.Printf("%s: %d\n", word, count)
	}
}
```

## Bài tập

1. **Loại bỏ phần tử trùng**: viết hàm `removeDuplicates(nums []int) []int` (gợi ý: dùng `map[int]bool` để đánh dấu giá trị đã gặp).
2. **Đếm tần suất từ**: dùng code mẫu trên làm nền, tự viết lại `wordFrequency`, in kết quả theo thứ tự alphabet (gợi ý: lấy keys ra 1 slice rồi `sort.Strings`).
3. **Đảo ngược slice in-place**: viết hàm `reverse(nums []int)` đảo ngược slice ngay trên chính nó (không tạo slice mới), dùng kỹ thuật 2 con trỏ đầu-cuối.
4. **Thử nghiệm chia sẻ underlying array**: viết code minh họa hiện tượng ở mục 3 (sửa slice con ảnh hưởng slice cha), sau đó sửa lại bằng `copy()` để chứng minh đã tách biệt.

### Gợi ý bài 1
```go
func removeDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, n := range nums {
		if !seen[n] {
			seen[n] = true
			result = append(result, n)
		}
	}
	return result
}
```

### Gợi ý bài 3
```go
func reverse(nums []int) {
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
}
```

## Tiếp theo
→ [Bài 6: Struct & Method](./6_structs_methods.md)
