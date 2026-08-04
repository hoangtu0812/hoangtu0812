# Bài 3: Luồng điều khiển

## Mục tiêu
- `if/elif/else`, `for`, `while` trong Python — cú pháp dựa trên **thụt lề (indentation)**, không dùng `{}`.
- `range()`, list comprehension cơ bản, walrus operator `:=`.

## 1. Indentation thay cho `{}` — điểm khác biệt lớn nhất với Go

```python
score = 75

if score >= 90:
	print("Xuất sắc")
elif score >= 70:
	print("Khá")
else:
	print("Cần cố gắng")
```

Python dùng **thụt lề nhất quán** (khuyến khích 4 space hoặc 1 tab, không trộn lẫn) để xác định khối lệnh — không có `{}`, không có `;` bắt buộc cuối dòng. Sai thụt lề = lỗi cú pháp (`IndentationError`).

## 2. Toán tử so sánh liên tiếp — tính năng riêng của Python

```python
age = 25
# Python cho phép viết liền mạch, Go không có cú pháp này
if 18 <= age < 65:
	print("Trong độ tuổi lao động")
```

## 3. Vòng lặp `for` — luôn là `for...in`, không có dạng đếm kiểu C

```python
# range(start, stop, step) — stop KHÔNG bao gồm, giống Go's for i:=0;i<n;i++
for i in range(5):        # 0 1 2 3 4
	print(i)

for i in range(2, 10, 2):  # 2 4 6 8
	print(i)

# Duyệt trực tiếp qua list — không cần index, giống range của Go trên slice
fruits = ["táo", "cam", "chuối"]
for fruit in fruits:
	print(fruit)

# Cần cả index: dùng enumerate() — tương đương "for i, v := range slice" của Go
for i, fruit in enumerate(fruits):
	print(i, fruit)
```

## 4. Vòng lặp `while`

```python
n = 0
while n < 5:
	print(n)
	n += 1   # Python KHÔNG có ++/--, phải dùng += 1 / -= 1
```

## 5. `break`, `continue`, và mệnh đề `else` của vòng lặp (tính năng độc đáo của Python)

```python
for i in range(10):
	if i == 5:
		break
	print(i)
else:
	# else của for/while chỉ chạy nếu vòng lặp KẾT THÚC TỰ NHIÊN (không bị break)
	print("Vòng lặp hoàn thành không bị ngắt")
```

Idiom hay dùng: tìm kiếm trong list, `else` chạy khi "không tìm thấy":

```python
target = 99
numbers = [1, 2, 3, 4, 5]
for n in numbers:
	if n == target:
		print("Tìm thấy!")
		break
else:
	print("Không tìm thấy trong danh sách")
```

## 6. List comprehension cơ bản — thay thế filter/map ngắn gọn

```python
# Cách viết dài (giống vòng lặp Go)
squares = []
for x in range(10):
	squares.append(x ** 2)

# List comprehension — cách viết Python-idiomatic, ngắn gọn hơn nhiều
squares = [x ** 2 for x in range(10)]

# Có điều kiện lọc
evens = [x for x in range(20) if x % 2 == 0]

print(squares)
print(evens)
```

## 7. Walrus operator `:=` (Python 3.8+)

Gán giá trị NGAY TRONG biểu thức điều kiện — tương tự idiom `if v, err := f(); err != nil` của Go:

```python
import random

# Không có walrus: phải gọi hàm 2 lần hoặc tách dòng
value = random.randint(1, 10)
if value > 5:
	print(f"Số lớn: {value}")

# Có walrus: gọn hơn, giá trị được gán VÀ kiểm tra cùng lúc
if (value := random.randint(1, 10)) > 5:
	print(f"Số lớn: {value}")
```

## Ví dụ đầy đủ: FizzBuzz

```python
def fizzbuzz(n: int) -> None:
	for i in range(1, n + 1):
		if i % 15 == 0:
			print("FizzBuzz")
		elif i % 3 == 0:
			print("Fizz")
		elif i % 5 == 0:
			print("Buzz")
		else:
			print(i)

if __name__ == "__main__":
	fizzbuzz(100)
```

## Bài tập

1. **FizzBuzz 1-100**: dùng code mẫu trên làm nền, tự gõ lại.
2. **Số nguyên tố**: viết hàm `is_prime(n: int) -> bool`, in các số nguyên tố từ 2-100 bằng list comprehension: `[n for n in range(2, 101) if is_prime(n)]`.
3. **Bảng cửu chương**: 2 vòng `for` lồng nhau in bảng cửu chương 1-9.
4. **Tìm kiếm dùng `for...else`**: viết hàm tìm số trong list, dùng `for...else` để in "không tìm thấy" khi vòng lặp không `break`.

## Tiếp theo
→ [Bài 4: Hàm (Functions)](./4_functions.md)
