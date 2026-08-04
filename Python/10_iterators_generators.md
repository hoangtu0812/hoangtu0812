# Bài 10: Iterator & Generator

## Mục tiêu
- Hiểu protocol `__iter__`/`__next__` đứng sau mọi vòng `for`.
- Viết generator bằng `yield`.
- Lazy evaluation — vì sao nó tiết kiệm bộ nhớ.

## 1. Iterator protocol — điều gì xảy ra khi bạn viết `for x in items`

```python
numbers = [1, 2, 3]
it = iter(numbers)     # gọi __iter__() -> trả về iterator
print(next(it))          # 1 — gọi __next__()
print(next(it))          # 2
print(next(it))          # 3
print(next(it))          # StopIteration exception -> vòng for tự bắt lỗi này để dừng lại
```

`for x in items:` thực chất là cú pháp đường (syntactic sugar) cho vòng lặp gọi `iter()` rồi `next()` liên tục tới khi `StopIteration`.

## 2. Tự viết iterator (dùng class)

```python
class CountUp:
	def __init__(self, start: int, end: int):
		self.current = start
		self.end = end

	def __iter__(self):
		return self

	def __next__(self):
		if self.current > self.end:
			raise StopIteration
		value = self.current
		self.current += 1
		return value

for n in CountUp(1, 5):
	print(n)  # 1 2 3 4 5
```

## 3. Generator — cách viết iterator NGẮN GỌN HƠN NHIỀU bằng `yield`

```python
def count_up(start: int, end: int):
	current = start
	while current <= end:
		yield current   # "tạm dừng" hàm tại đây, trả giá trị ra ngoài, nhớ trạng thái để lần next() sau tiếp tục
		current += 1

for n in count_up(1, 5):
	print(n)  # 1 2 3 4 5
```

Hàm có `yield` tự động trở thành **generator function** — gọi nó không chạy code ngay, mà trả về 1 generator object (đã tự cài sẵn `__iter__`/`__next__`).

## 4. Generator cho dữ liệu vô hạn / lớn — điểm mạnh so với list

```python
def fibonacci():
	a, b = 0, 1
	while True:  # vô hạn — KHÔNG thể làm điều này với 1 list thông thường
		yield a
		a, b = b, a + b

fib = fibonacci()
first_10 = [next(fib) for _ in range(10)]
print(first_10)  # [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]
```

```python
def read_large_file(path: str):
	"""Đọc file lớn từng dòng, không load toàn bộ vào RAM."""
	with open(path) as f:
		for line in f:
			yield line.strip()

# Xử lý file hàng GB mà không tốn RAM tương ứng — mỗi lần chỉ giữ 1 dòng trong bộ nhớ
for line in read_large_file("huge_log.txt"):
	if "ERROR" in line:
		print(line)
```

## 5. `yield from` — ủy quyền cho generator khác

```python
def inner():
	yield 1
	yield 2

def outer():
	yield from inner()  # tương đương "for x in inner(): yield x"
	yield 3

print(list(outer()))  # [1, 2, 3]
```

## 6. Generator expression (đã học sơ ở [Bài 6](./6_strings_comprehensions.md))

```python
squares = (x**2 for x in range(10))  # generator, không phải list
print(sum(squares))  # dùng được ngay vì sum() tự duyệt qua generator
```

## Ví dụ đầy đủ

```python
def batch(iterable, size: int):
	"""Chia 1 iterable thành từng nhóm (batch) kích thước cố định — hữu ích khi gọi API theo lô."""
	batch_items = []
	for item in iterable:
		batch_items.append(item)
		if len(batch_items) == size:
			yield batch_items
			batch_items = []
	if batch_items:  # trả nốt phần còn lại nếu không chia hết
		yield batch_items

if __name__ == "__main__":
	numbers = range(1, 11)
	for group in batch(numbers, 3):
		print(group)
	# [1, 2, 3]
	# [4, 5, 6]
	# [7, 8, 9]
	# [10]
```

## Bài tập

1. **Generator Fibonacci**: viết `fibonacci()` như trên, lấy 10 số đầu bằng `next()` trong vòng lặp/list comprehension.
2. **Iterator tự định nghĩa**: viết class `CountUp` như ví dụ, hoặc viết class `Countdown` đếm ngược.
3. **`batch()`**: dùng code mẫu trên làm nền, tự viết lại — chia 1 list bất kỳ thành các nhóm kích thước N.
4. **Nâng cao**: viết generator `read_large_file(path)` đọc file theo dòng, kết hợp với 1 generator khác lọc dòng chứa từ khóa cụ thể — chứng minh 2 generator có thể "chuỗi" (chain) vào nhau mà không tốn thêm bộ nhớ trung gian.

## Tiếp theo
→ [Bài 11: Decorator & Context Manager](./11_decorators_context_managers.md)
