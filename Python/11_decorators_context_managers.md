# Bài 11: Decorator & Context Manager

## Mục tiêu
- Hiểu hàm bậc cao, viết decorator cơ bản và decorator có tham số.
- `functools.wraps`.
- Tự viết context manager (`__enter__`/`__exit__` hoặc `@contextmanager`).

## 1. Hàm là first-class citizen (giống Go — xem [Go Bài 4](../Go/4_functions.md))

```python
def shout(text: str) -> str:
	return text.upper() + "!"

def apply(func, value):
	return func(value)

print(apply(shout, "hello"))  # HELLO!
```

## 2. Decorator cơ bản — hàm "bọc" hàm khác

```python
def logging_decorator(func):
	def wrapper(*args, **kwargs):
		print(f"Gọi hàm {func.__name__} với args={args}, kwargs={kwargs}")
		result = func(*args, **kwargs)
		print(f"{func.__name__} trả về {result}")
		return result
	return wrapper

@logging_decorator
def add(a, b):
	return a + b

add(3, 5)
# Gọi hàm add với args=(3, 5), kwargs={}
# add trả về 8
```

`@logging_decorator` phía trên `def add` tương đương `add = logging_decorator(add)` — cú pháp đường (syntactic sugar), khái niệm giống middleware bọc handler trong [Go Bài 18](../Go/18_rest_api.md).

## 3. `functools.wraps` — giữ metadata của hàm gốc

```python
from functools import wraps

def logging_decorator(func):
	@wraps(func)  # BẮT BUỘC nên có — giữ nguyên __name__, __doc__ của hàm gốc
	def wrapper(*args, **kwargs):
		print(f"Gọi {func.__name__}")
		return func(*args, **kwargs)
	return wrapper

@logging_decorator
def add(a, b):
	"""Cộng 2 số."""
	return a + b

print(add.__name__)  # "add" — nếu KHÔNG dùng @wraps sẽ in ra "wrapper", gây khó debug
```

## 4. Decorator có tham số — cần thêm 1 lớp hàm bọc ngoài

```python
import time
from functools import wraps

def retry(times: int):
	def decorator(func):
		@wraps(func)
		def wrapper(*args, **kwargs):
			last_exception = None
			for attempt in range(1, times + 1):
				try:
					return func(*args, **kwargs)
				except Exception as e:
					last_exception = e
					print(f"Lần thử {attempt} thất bại: {e}")
			raise last_exception
		return wrapper
	return decorator

@retry(times=3)
def unstable_call():
	import random
	if random.random() < 0.7:
		raise ConnectionError("Kết nối thất bại")
	return "Thành công"

unstable_call()
```

## 5. Decorator đo thời gian chạy (rất hay dùng thực tế)

```python
import time
from functools import wraps

def timer(func):
	@wraps(func)
	def wrapper(*args, **kwargs):
		start = time.perf_counter()
		result = func(*args, **kwargs)
		elapsed = time.perf_counter() - start
		print(f"{func.__name__} chạy trong {elapsed:.4f}s")
		return result
	return wrapper

@timer
def slow_function():
	time.sleep(1)

slow_function()
```

## 6. Context Manager — `with` (nối tiếp [Bài 8](./8_exceptions.md))

### Cách 1: viết class với `__enter__`/`__exit__`

```python
class Timer:
	def __enter__(self):
		import time
		self.start = time.perf_counter()
		return self  # giá trị này gán cho biến sau "as"

	def __exit__(self, exc_type, exc_value, traceback):
		import time
		elapsed = time.perf_counter() - self.start
		print(f"Khối lệnh chạy trong {elapsed:.4f}s")
		return False  # False = KHÔNG "nuốt" exception, để nó tiếp tục propagate lên trên

with Timer():
	import time
	time.sleep(1)
```

`__exit__` được gọi **dù khối lệnh có raise exception hay không** — tương đương `defer` của Go dùng để dọn dẹp resource ([Go Bài 4](../Go/4_functions.md)).

### Cách 2: `@contextmanager` — ngắn gọn hơn, dùng generator

```python
from contextlib import contextmanager
import time

@contextmanager
def timer():
	start = time.perf_counter()
	try:
		yield  # code trong khối "with" chạy TẠI ĐÂY
	finally:
		elapsed = time.perf_counter() - start
		print(f"Chạy trong {elapsed:.4f}s")

with timer():
	time.sleep(1)
```

Phần trước `yield` = logic của `__enter__`, phần sau `yield` (trong `finally`) = logic của `__exit__`.

## Bài tập

1. **Decorator đo thời gian**: viết `@timer` như ví dụ, áp dụng cho 1 hàm bất kỳ có `time.sleep`.
2. **Decorator retry**: viết `@retry(times=N)`, thử với 1 hàm giả lập thất bại ngẫu nhiên.
3. **Context manager tự viết**: viết context manager quản lý việc mở/đóng "kết nối" giả lập (chỉ cần `print("Mở kết nối")`/`print("Đóng kết nối")`), thử cả 2 cách (class và `@contextmanager`).
4. **Nâng cao**: viết decorator `@cache` đơn giản (dùng dict lưu kết quả theo tham số đầu vào) để tránh tính lại hàm tốn thời gian (đệ quy Fibonacci không tối ưu là ví dụ kinh điển) — sau đó so sánh với `functools.lru_cache` có sẵn.

## Tiếp theo
→ [Bài 12: File I/O & JSON/CSV](./12_file_io_json.md)
