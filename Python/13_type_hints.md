# Bài 13: Type Hints & Static Checking

## Mục tiêu
- Dùng `typing` để lấy lại 1 phần lợi ích "kiểm tra kiểu lúc biên dịch" mà Go có sẵn.
- `Optional`, `Union`, generic (`TypeVar`).
- Kiểm tra static bằng `mypy`.

## 1. Vì sao cần type hint?

Python là **dynamically typed** — không type hint, IDE và người đọc code không biết hàm nhận/trả về kiểu gì cho tới khi chạy thử. Type hint (từ Python 3.5+) thêm annotation kiểu **chỉ mang tính khai báo** — Python **không tự kiểm tra lúc chạy** (khác Go biên dịch sẽ báo lỗi), cần công cụ riêng như `mypy` để kiểm tra tĩnh.

```python
def add(a: int, b: int) -> int:
	return a + b

# Type hint KHÔNG ngăn được điều này lúc runtime — Python vẫn chạy, chỉ mypy mới báo lỗi
add("2", "3")  # mypy sẽ báo lỗi, nhưng Python interpreter vẫn chạy và trả về "23"
```

## 2. Kiểu cơ bản

```python
name: str = "Ben"
age: int = 25
height: float = 1.75
is_active: bool = True

# List, dict, tuple, set với kiểu phần tử (Python 3.9+, không cần import từ typing nữa)
numbers: list[int] = [1, 2, 3]
scores: dict[str, int] = {"Alice": 90}
point: tuple[int, int] = (1, 2)
unique: set[str] = {"a", "b"}
```

## 3. `Optional` và `Union` — tương đương "có thể nil" của Go

```python
from typing import Optional, Union

def find_user(user_id: int) -> Optional[str]:  # Optional[str] = str | None
	if user_id == 1:
		return "Ben"
	return None

# Từ Python 3.10+, có thể viết ngắn gọn hơn bằng toán tử |
def find_user_v2(user_id: int) -> str | None:
	return "Ben" if user_id == 1 else None

def process(value: Union[int, str]) -> str:  # value có thể là int HOẶC str
	return str(value)

# Python 3.10+ viết gọn: int | str thay vì Union[int, str]
def process_v2(value: int | str) -> str:
	return str(value)
```

## 4. Type hint cho hàm, class, collection lồng nhau

```python
from dataclasses import dataclass

@dataclass
class User:
	id: int
	name: str
	email: str

def get_users() -> list[User]:
	return [User(1, "Ben", "ben@example.com")]

def group_by_domain(users: list[User]) -> dict[str, list[User]]:
	result: dict[str, list[User]] = {}
	for user in users:
		domain = user.email.split("@")[1]
		result.setdefault(domain, []).append(user)
	return result

def apply_discount(items: list[dict[str, float]], rate: float) -> list[dict[str, float]]:
	return [{**item, "price": item["price"] * (1 - rate)} for item in items]
```

## 5. Generic với `TypeVar` (trước khi có generic mới hơn ở Python 3.12)

```python
from typing import TypeVar

T = TypeVar("T")

def first_or_none(items: list[T]) -> T | None:
	return items[0] if items else None

# Python 3.12+ có cú pháp generic mới, gọn hơn TypeVar (tương tự generics của Go — xem Go Bài 15)
def first_or_none_v2[T](items: list[T]) -> T | None:
	return items[0] if items else None
```

## 6. `Protocol` — "interface" của Python (liên hệ trực tiếp interface của Go)

```python
from typing import Protocol

class Shape(Protocol):
	def area(self) -> float: ...

class Rectangle:
	def __init__(self, w: float, h: float):
		self.w, self.h = w, h
	def area(self) -> float:
		return self.w * self.h

def print_area(shape: Shape) -> None:  # Rectangle KHÔNG cần khai báo "implements Shape"
	print(shape.area())

print_area(Rectangle(10, 5))  # hoạt động — structural typing, giống hệt interface của Go!
```

`Protocol` (PEP 544) hoạt động theo **structural typing** giống hệt interface của Go ([Go Bài 10](../Go/10_interfaces.md)) — không cần khai báo `implements` tường minh, chỉ cần có đủ method. Đây là công cụ quan trọng nhất để làm dependency injection kiểu Pythonic ở [Bài 16](./16_clean_architecture.md).

## 7. Kiểm tra static với `mypy`

```powershell
pip install mypy
mypy your_file.py
```

```python
# example.py
def add(a: int, b: int) -> int:
	return a + b

result: str = add(1, 2)  # mypy báo lỗi: Incompatible types (int vs str)
```

Chạy `mypy example.py` sẽ báo:
```
example.py:4: error: Incompatible types in assignment (expression has type "int", variable has type "str")
```

## Bài tập

1. **Thêm type hint đầy đủ**: lấy code `Rectangle`/`Circle`/`Shape` (Bài 7 hoặc Go Bài 10), thêm type hint đầy đủ cho tham số và giá trị trả về.
2. **`mypy`**: cài `mypy`, cố tình viết 1 lỗi kiểu (gán sai kiểu), chạy `mypy` để xác nhận nó phát hiện được.
3. **`Optional`**: viết hàm `find_by_id(items: list[User], id: int) -> User | None`, xử lý đúng trường hợp không tìm thấy.
4. **`Protocol`**: viết `Protocol` tên `Notifier` với method `notify(message: str) -> None`, viết 2 class implement (không cần kế thừa tường minh) là `EmailNotifier`, `SMSNotifier`, viết hàm nhận `Notifier` dùng chung cho cả 2 — liên hệ [Go Bài 10 bài tập 4](../Go/10_interfaces.md).

## Tiếp theo
→ [Bài 14: Testing với pytest](./14_testing.md)
