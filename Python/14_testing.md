# Bài 14: Testing với pytest

## Mục tiêu
- Viết test bằng `pytest` (thư viện test phổ biến nhất, hơn hẳn `unittest` built-in).
- Fixture, parametrize (tương đương table-driven test của Go).
- Mock dependency bằng `unittest.mock`.

## 1. Cài đặt & quy ước file test

```powershell
pip install pytest
```

File test có tiền tố `test_` hoặc hậu tố `_test`, hàm test có tiền tố `test_`:

```
myproject/
├── shapes/
│   └── rectangle.py
└── tests/
    └── test_rectangle.py
```

## 2. Test cơ bản

```python
# shapes/rectangle.py
def divide(a: float, b: float) -> float:
	if b == 0:
		raise ValueError("không thể chia cho 0")
	return a / b
```

```python
# tests/test_rectangle.py
import pytest
from shapes.rectangle import divide

def test_divide_normal():
	assert divide(10, 2) == 5

def test_divide_by_zero_raises():
	with pytest.raises(ValueError, match="không thể chia"):
		divide(10, 0)
```

Chạy: `pytest` hoặc `pytest -v` (verbose, in tên từng test) hoặc `pytest tests/test_rectangle.py::test_divide_normal` (chạy 1 test cụ thể).

`assert` trần là đủ — `pytest` tự phân tích biểu thức và in thông báo lỗi chi tiết khi fail (khác `unittest` cần `self.assertEqual(...)`).

## 3. Parametrize — tương đương table-driven test của Go

```python
import pytest
from shapes.rectangle import divide

@pytest.mark.parametrize("a, b, expected", [
	(10, 2, 5),
	(9, 3, 3),
	(-10, 2, -5),
	(0, 5, 0),
])
def test_divide_various(a, b, expected):
	assert divide(a, b) == expected

@pytest.mark.parametrize("b", [0])
def test_divide_by_zero(b):
	with pytest.raises(ValueError):
		divide(10, b)
```

`pytest -v` sẽ in ra từng case riêng biệt (`test_divide_various[10-2-5]`...) — tương đương subtest `t.Run()` của Go ([Go Bài 14](../Go/14_testing.md)).

## 4. Fixture — setup/teardown tái sử dụng

```python
import pytest

@pytest.fixture
def sample_users():
	# "Arrange" — chuẩn bị dữ liệu dùng chung cho nhiều test
	return [
		{"id": 1, "name": "Ben"},
		{"id": 2, "name": "Alice"},
	]

def test_find_user(sample_users):
	# pytest tự động "tiêm" fixture vào bằng cách khớp TÊN THAM SỐ
	user = next(u for u in sample_users if u["id"] == 1)
	assert user["name"] == "Ben"

@pytest.fixture
def db_connection():
	conn = create_fake_connection()
	yield conn          # giá trị trả cho test — code TRƯỚC yield là setup
	conn.close()          # code SAU yield là teardown, LUÔN chạy dù test pass/fail
```

## 5. Mock dependency

```python
from unittest.mock import Mock, patch

class EmailService:
	def send(self, to: str, message: str) -> bool:
		# giả sử đây gọi API email thật, tốn thời gian/tiền
		...

class NotificationService:
	def __init__(self, email_service: EmailService):
		self.email_service = email_service

	def notify(self, user_email: str) -> bool:
		return self.email_service.send(user_email, "Xin chào!")

def test_notify_calls_email_service():
	mock_email = Mock(spec=EmailService)  # spec= giới hạn mock chỉ có method của EmailService thật
	mock_email.send.return_value = True

	service = NotificationService(mock_email)
	result = service.notify("ben@example.com")

	assert result is True
	mock_email.send.assert_called_once_with("ben@example.com", "Xin chào!")
```

`Mock`/`patch` tương đương việc viết mock implementation cho interface trong Go ([Go Bài 14](../Go/14_testing.md)) — nhưng Python cho phép mock "runtime" không cần định nghĩa interface tường minh trước (tuy vậy nên dùng `Protocol` — [Bài 13](./13_type_hints.md) — để giữ rõ ràng).

## 6. Test coverage

```powershell
pip install pytest-cov
pytest --cov=shapes --cov-report=term-missing
```

## Ví dụ đầy đủ

```python
# tests/test_user_service.py
import pytest
from unittest.mock import Mock

class UserNotFoundError(Exception):
	pass

class UserRepository:  # Protocol thực tế nên dùng typing.Protocol — xem Bài 13
	def find_by_id(self, user_id: int) -> dict: ...

class UserService:
	def __init__(self, repo: UserRepository):
		self.repo = repo

	def get_user_name(self, user_id: int) -> str:
		user = self.repo.find_by_id(user_id)
		if user is None:
			raise UserNotFoundError(f"Không tìm thấy user {user_id}")
		return user["name"]

@pytest.fixture
def mock_repo():
	return Mock(spec=UserRepository)

def test_get_user_name_success(mock_repo):
	mock_repo.find_by_id.return_value = {"id": 1, "name": "Ben"}
	service = UserService(mock_repo)

	assert service.get_user_name(1) == "Ben"

def test_get_user_name_not_found(mock_repo):
	mock_repo.find_by_id.return_value = None
	service = UserService(mock_repo)

	with pytest.raises(UserNotFoundError):
		service.get_user_name(999)
```

## Bài tập

1. **Parametrize cho `divide`**: viết test bao phủ ít nhất 4 case bằng `@pytest.mark.parametrize`.
2. **Fixture**: viết fixture cung cấp danh sách user mẫu, dùng trong ít nhất 2 test khác nhau.
3. **Mock dependency**: viết `UserService` phụ thuộc `UserRepository` (Protocol), viết test dùng `Mock` để kiểm tra logic mà không cần database thật — mirror bài tập tương tự ở [Go Bài 14](../Go/14_testing.md).
4. **Coverage**: cài `pytest-cov`, chạy để đạt >80% coverage cho package bạn đã viết ở các bài trước.

## Tiếp theo
→ [Bài 15: Concurrency (threading, multiprocessing, asyncio)](./15_concurrency.md)
