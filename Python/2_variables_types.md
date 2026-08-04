# Bài 2: Biến, Kiểu dữ liệu, Toán tử

## Mục tiêu
- Khai báo biến trong Python (typing động — khác hẳn Go).
- Nắm các kiểu dữ liệu cơ bản.
- Toán tử, f-string, phân biệt `is` vs `==`.

## 1. Khai báo biến — không cần từ khóa

```python
name = "Ben"       # str
age = 25            # int
height = 1.75       # float
is_active = True    # bool
nothing = None       # None — tương đương nil của Go

print(name, age, height, is_active, nothing)
```

Không có `var`/`:=` — Python suy luận kiểu tự động và **kiểu có thể đổi khi gán lại** (khác Go: biến Go giữ nguyên kiểu suốt vòng đời):

```python
x = 5        # x là int
x = "five"   # giờ x là str — HOÀN TOÀN hợp lệ trong Python (dynamic typing)
```

Đây vừa là điểm mạnh (linh hoạt) vừa là nguồn lỗi runtime phổ biến nhất khi chuyển từ Go sang — không có compiler bắt lỗi kiểu giúp bạn, nên cẩn trọng hơn (xem [Bài 13: Type Hints](./13_type_hints.md) để lấy lại phần nào sự an toàn kiểu).

## 2. Kiểu dữ liệu cơ bản

```python
i = 42          # int — không giới hạn kích thước như int64 (Python tự động dùng số lớn tùy ý)
f = 3.14        # float — luôn là double-precision (giống float64 của Go)
s = "hello"     # str — immutable, hỗ trợ Unicode đầy đủ
b = True        # bool
n = None        # None — giá trị "rỗng", tương đương nil

print(type(i), type(f), type(s), type(b), type(n))
# <class 'int'> <class 'float'> <class 'str'> <class 'bool'> <class 'NoneType'>
```

## 3. Toán tử

```python
a, b = 10, 3

# Số học
print(a + b, a - b, a * b, a / b, a // b, a % b, a ** b)
# 13 7 30 3.333... 3 1 1000
# Lưu ý: / luôn trả về float; // là chia lấy nguyên (floor division); ** là lũy thừa

# So sánh
print(a > b, a == b, a != b)  # True False True

# Logic — dùng từ khóa, KHÔNG dùng &&/||/!
print(True and False, True or False, not True)  # False True False
```

## 4. `is` vs `==`

```python
a = [1, 2, 3]
b = [1, 2, 3]
c = a

print(a == b)  # True — so sánh GIÁ TRỊ (nội dung giống nhau)
print(a is b)  # False — so sánh DANH TÍNH (2 object khác nhau trong bộ nhớ)
print(a is c)  # True — c và a trỏ tới CÙNG 1 object
```

**Quy tắc:** dùng `==` để so sánh giá trị (hầu hết trường hợp), chỉ dùng `is` khi so sánh danh tính object — phổ biến nhất là `x is None` (không dùng `x == None`).

## 5. f-string — format chuỗi (quan trọng, dùng liên tục)

```python
name = "Ben"
age = 25

# f-string (Python 3.6+) — cách chuẩn, tương đương fmt.Sprintf của Go
message = f"Tôi tên {name}, {age} tuổi"
print(message)

# Có thể nhúng biểu thức, format số
pi = 3.14159
print(f"Pi làm tròn 2 chữ số: {pi:.2f}")  # Pi làm tròn 2 chữ số: 3.14
print(f"Năm sau tôi {age + 1} tuổi")
```

## 6. Ép kiểu (type conversion)

```python
x = "42"
y = int(x)     # str -> int
z = float(y)   # int -> float
s = str(z)     # float -> str

print(int("10") + int("20"))  # 30
# int("abc") sẽ raise ValueError — luôn cẩn trọng khi convert dữ liệu từ input/API
```

## Ví dụ đầy đủ: Đổi độ C ↔ F

```python
def celsius_to_fahrenheit(c: float) -> float:
	return c * 9 / 5 + 32

def fahrenheit_to_celsius(f: float) -> float:
	return (f - 32) * 5 / 9

if __name__ == "__main__":
	celsius = 30
	print(f"{celsius}°C = {celsius_to_fahrenheit(celsius):.1f}°F")

	fahrenheit = 98.6
	print(f"{fahrenheit}°F = {fahrenheit_to_celsius(fahrenheit):.1f}°C")
```

## Bài tập

1. **Khai báo & `type()`**: khai báo ít nhất 5 biến đủ kiểu (int, float, str, bool, None), in ra kiểu bằng `type()`.
2. **Đổi độ C ↔ F**: viết lại ví dụ trên, tự gõ không copy.
3. **`is` vs `==`**: viết code minh họa rõ trường hợp `a == b` là `True` nhưng `a is b` là `False`; giải thích bằng comment.
4. **Type conversion an toàn**: viết hàm `safe_int(s: str) -> int | None` dùng `try/except` (xem trước ở [Bài 8](./8_exceptions.md), hoặc dùng `str.isdigit()`) để convert an toàn, trả `None` nếu không convert được thay vì crash.

## Tiếp theo
→ [Bài 3: Luồng điều khiển](./3_control_flow.md)
