# Bài 4: Hàm (Functions)

## Mục tiêu
- Khai báo hàm với `def`, tham số mặc định, `*args`/`**kwargs`.
- Lambda, closure, docstring.

## 1. Khai báo hàm cơ bản

```python
def add(a, b):
	return a + b

# Với type hint (khuyến khích — xem thêm Bài 13)
def subtract(a: int, b: int) -> int:
	return a - b
```

## 2. Tham số mặc định

```python
def greet(name, greeting="Xin chào"):
	return f"{greeting}, {name}!"

print(greet("Ben"))              # Xin chào, Ben!
print(greet("Ben", "Hi"))        # Hi, Ben!
print(greet(name="Ben", greeting="Hello"))  # gọi bằng keyword argument
```

**Cạm bẫy kinh điển:** KHÔNG dùng mutable object (list, dict) làm giá trị mặc định:

```python
# SAI — list mặc định bị CHIA SẺ giữa các lần gọi vì chỉ khởi tạo 1 LẦN DUY NHẤT
def add_item(item, items=[]):
	items.append(item)
	return items

print(add_item("a"))  # ['a']
print(add_item("b"))  # ['a', 'b'] — KHÔNG như mong đợi!

# ĐÚNG — dùng None làm mặc định, tạo list mới bên trong hàm
def add_item(item, items=None):
	if items is None:
		items = []
	items.append(item)
	return items
```

## 3. Nhiều giá trị trả về — thực chất là trả về 1 tuple

```python
def divide(a, b):
	if b == 0:
		return None, "không thể chia cho 0"
	return a / b, None

result, error = divide(10, 2)
if error:
	print("Lỗi:", error)
else:
	print("Kết quả:", result)
```

Khác Go (multiple return values là tính năng ngôn ngữ), Python chỉ trả về 1 giá trị — nhưng "unpack" tuple khi gán khiến nó *trông* giống multiple return values.

## 4. `*args` và `**kwargs` — biến đối số

```python
# *args: gom các positional argument thừa thành tuple
def total(*args):
	return sum(args)

print(total(1, 2, 3))       # 6
print(total())               # 0

# **kwargs: gom các keyword argument thừa thành dict
def print_info(**kwargs):
	for key, value in kwargs.items():
		print(f"{key}: {value}")

print_info(name="Ben", age=25)

# Kết hợp cả 2 — thứ tự BẮT BUỘC: positional, *args, keyword mặc định, **kwargs
def full_example(a, b, *args, c=10, **kwargs):
	print(a, b, args, c, kwargs)

full_example(1, 2, 3, 4, c=99, d=100)  # 1 2 (3, 4) 99 {'d': 100}
```

## 5. Lambda — hàm ẩn danh 1 dòng

```python
square = lambda x: x ** 2
print(square(5))  # 25

# Dùng phổ biến nhất: làm key cho sort/filter/map
people = [("Ben", 25), ("Alice", 30), ("Carol", 20)]
people.sort(key=lambda p: p[1])  # sắp xếp theo tuổi
print(people)

# Lambda chỉ chứa được 1 BIỂU THỨC, không chứa statement (if/for as statement)
# Nếu logic phức tạp hơn 1 dòng, nên dùng def thay vì cố nhồi vào lambda
```

## 6. Closure

```python
def make_counter():
	count = 0
	def counter():
		nonlocal count   # BẮT BUỘC khai báo nonlocal để sửa biến ở scope ngoài
		count += 1
		return count
	return counter

c1 = make_counter()
print(c1())  # 1
print(c1())  # 2

c2 = make_counter()  # instance độc lập
print(c2())  # 1
```

`nonlocal` tương tự việc Go tự động "bắt" biến trong closure — nhưng Python yêu cầu khai báo tường minh khi bạn **gán lại** (không cần nếu chỉ đọc).

## 7. Docstring — convention document hàm

```python
def divide(a: float, b: float) -> float:
	"""Chia a cho b.

	Raises:
		ZeroDivisionError: nếu b bằng 0.
	"""
	return a / b

print(divide.__doc__)  # in ra docstring
```

## Ví dụ đầy đủ

```python
def divide(a: float, b: float) -> tuple[float | None, str | None]:
	if b == 0:
		return None, "không thể chia cho 0"
	return a / b, None

def total(*nums: int) -> int:
	return sum(nums)

def make_counter():
	count = 0
	def counter():
		nonlocal count
		count += 1
		return count
	return counter

if __name__ == "__main__":
	result, err = divide(10, 0)
	print("Lỗi:", err) if err else print("Kết quả:", result)

	print("Tổng:", total(1, 2, 3, 4, 5))

	counter = make_counter()
	for _ in range(3):
		print("Đếm:", counter())
```

## Bài tập

1. **`divide` an toàn**: viết hàm trả `(result, error)` như ví dụ, thử với b=0 và b khác 0.
2. **`total(*nums)`**: viết hàm biến đối số tính tổng, gọi thử với nhiều số và với list dùng `*` để unpack (`total(*[1,2,3])`).
3. **Closure counter**: viết `make_counter()`, tạo 2 counter độc lập, chứng minh chúng không ảnh hưởng lẫn nhau.
4. **Cạm bẫy mutable default**: viết lại ví dụ `add_item` ở mục 2 với cả 2 phiên bản (sai và đúng), chạy để tận mắt thấy sự khác biệt.

## Tiếp theo
→ [Bài 5: List, Tuple, Set, Dict](./5_collections.md)
