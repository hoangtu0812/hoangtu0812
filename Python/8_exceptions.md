# Bài 8: Xử lý lỗi (Exception Handling)

## Mục tiêu
- `try/except/else/finally`.
- Custom exception, `raise ... from`.
- Context manager cơ bản (`with`) — liên hệ [Bài 11](./11_decorators_context_managers.md) để tự viết context manager riêng.

## 1. Khác biệt triết lý so với Go

Go coi **error là giá trị** (trả về, kiểm tra `if err != nil` — xem [Go Bài 8](../Go/8_error_handling.md)). Python dùng **exception** — lỗi làm gián đoạn luồng thực thi bình thường và "nhảy" lên tới nơi có `except` phù hợp gần nhất. Đây là khác biệt lớn nhất giữa 2 ngôn ngữ về xử lý lỗi.

## 2. `try/except` cơ bản

```python
def divide(a: float, b: float) -> float:
	return a / b

try:
	result = divide(10, 0)
except ZeroDivisionError as e:
	print(f"Lỗi: {e}")
else:
	# else chỉ chạy nếu try KHÔNG có lỗi — ít dùng nhưng hữu ích để tách rõ code "happy path"
	print(f"Kết quả: {result}")
finally:
	# finally LUÔN chạy, dù có lỗi hay không — dùng để dọn dẹp tài nguyên
	print("Đã hoàn thành thử chia")
```

## 3. Bắt nhiều loại exception

```python
def parse_and_divide(a_str: str, b_str: str) -> float:
	try:
		a = int(a_str)
		b = int(b_str)
		return a / b
	except ValueError:
		print("Không convert được sang số nguyên")
		raise
	except ZeroDivisionError:
		print("Không thể chia cho 0")
		raise
	except Exception as e:  # bắt "mọi thứ khác" — CHỈ dùng khi thực sự cần, đặt SAU CÙNG
		print(f"Lỗi không xác định: {e}")
		raise
```

**Quy tắc:** bắt exception **cụ thể nhất có thể**, tránh `except:` trần (bắt mọi thứ kể cả `KeyboardInterrupt`) hoặc `except Exception:` quá rộng ở tầng thấp — điều này giấu bug thay vì xử lý nó.

## 4. Custom exception

```python
class ValidationError(Exception):
	def __init__(self, field: str, message: str):
		self.field = field
		self.message = message
		super().__init__(f"{field}: {message}")

def validate_age(age: int) -> None:
	if age < 0 or age > 150:
		raise ValidationError("age", "phải nằm trong khoảng 0-150")

try:
	validate_age(-5)
except ValidationError as e:
	print(f"Lỗi validate field '{e.field}': {e.message}")
```

Custom exception nên kế thừa từ `Exception` (hoặc 1 exception con cháu phù hợp hơn), tương tự việc tự định nghĩa `error` type trong Go ([Go Bài 8](../Go/8_error_handling.md)).

## 5. `raise ... from` — giữ nguyên "chuỗi nguyên nhân" lỗi (tương đương `%w` wrap của Go)

```python
class ConfigError(Exception):
	pass

def load_config():
	try:
		with open("config.yaml") as f:
			return f.read()
	except FileNotFoundError as e:
		raise ConfigError("không tìm thấy file config") from e  # giữ traceback gốc

try:
	load_config()
except ConfigError as e:
	print(e)
	print("Nguyên nhân gốc:", e.__cause__)  # tương đương errors.Unwrap() của Go
```

## 6. Context manager `with` — tự động dọn dẹp resource

```python
# Không dùng with: dễ quên đóng file, đặc biệt khi có exception xảy ra giữa chừng
f = open("data.txt")
data = f.read()
f.close()  # nếu read() lỗi, dòng này KHÔNG chạy -> rò rỉ file handle

# Dùng with: LUÔN tự động đóng file, kể cả khi có exception (tương đương defer f.Close() của Go)
with open("data.txt") as f:
	data = f.read()
# f đã tự động đóng tại đây
```

Bạn sẽ học cách **tự viết** context manager riêng (`__enter__`/`__exit__`, hoặc `@contextmanager`) ở [Bài 11](./11_decorators_context_managers.md).

## 7. Hệ thống phân cấp exception built-in (tham khảo nhanh)

```
BaseException
 ├── SystemExit, KeyboardInterrupt  (thường KHÔNG bắt trừ lý do đặc biệt)
 └── Exception                       (base class cho hầu hết exception nên bắt)
      ├── ValueError    (giá trị sai kiểu logic, vd int("abc"))
      ├── TypeError      (sai kiểu dữ liệu khi thao tác)
      ├── KeyError        (truy cập dict key không tồn tại bằng [])
      ├── IndexError      (truy cập list/tuple index ngoài phạm vi)
      ├── FileNotFoundError
      ├── ZeroDivisionError
      └── ... (và exception tự định nghĩa nên kế thừa Exception)
```

## Ví dụ đầy đủ

```python
class ValidationError(Exception):
	def __init__(self, field: str, message: str):
		self.field = field
		super().__init__(f"{field}: {message}")

def validate_user(name: str, age: int) -> None:
	if not name.strip():
		raise ValidationError("name", "không được để trống")
	if age < 0 or age > 150:
		raise ValidationError("age", "phải trong khoảng 0-150")

def process_users(users: list[tuple[str, int]]) -> list[str]:
	errors = []
	for name, age in users:
		try:
			validate_user(name, age)
		except ValidationError as e:
			errors.append(str(e))
	return errors

if __name__ == "__main__":
	users = [("Ben", 25), ("", 30), ("Alice", -1)]
	errors = process_users(users)
	for err in errors:
		print("Lỗi:", err)
```

## Bài tập

1. **Custom exception**: viết `ValidationError` như trên, viết `validate_age(age)` raise lỗi khi không hợp lệ, bắt và in ra thông báo.
2. **`with open()` an toàn**: viết hàm đọc file, dùng `with` + bắt `FileNotFoundError`, trả về `None` nếu file không tồn tại thay vì crash.
3. **`raise ... from`**: viết hàm wrap 1 exception gốc thành custom exception có ngữ cảnh rõ hơn, in `__cause__` để xem lỗi gốc.
4. **Nâng cao**: viết `process_orders(orders: list) -> list[str]` xử lý danh sách đơn hàng, thu thập TẤT CẢ lỗi (không dừng ở lỗi đầu tiên) vào 1 list thay vì để exception đầu tiên làm crash toàn bộ — liên hệ ý tưởng `errors.Join` ở [Go Bài 8](../Go/8_error_handling.md).

## Tiếp theo
→ [Bài 9: Module, Package & Virtual Environment](./9_modules_packages.md)
