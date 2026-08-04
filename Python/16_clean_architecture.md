# Bài 16: Kiến trúc project & Clean Architecture

## Mục tiêu
- Áp dụng mô hình phân tầng router → service → repository (giống [Go Bài 16](../Go/16_clean_architecture.md)) vào Python.
- Dùng `Protocol` để dependency injection, không cần framework DI phức tạp.

## 1. Mô hình 3 tầng — giống hệt bên Go, chỉ khác tên gọi tầng trên cùng

```
Router (FastAPI APIRouter)  →  Service (business logic)  →  Repository (data access)
        ↓                             ↓                            ↓
   nhận HTTP request             xử lý nghiệp vụ               CRUD với DB
   gọi service                    validate, tính toán          KHÔNG biết gì về HTTP
   trả response
```

## 2. Domain layer — model + Protocol (tương đương `domain.UserRepository` của Go)

```python
# app/domain/user.py
from dataclasses import dataclass
from typing import Protocol

@dataclass
class User:
	id: int
	name: str
	email: str
	role: str  # "admin" hoặc "user"

class UserRepository(Protocol):
	def find_by_id(self, user_id: int) -> User | None: ...
	def find_by_email(self, email: str) -> User | None: ...
	def create(self, user: User) -> User: ...
```

`Protocol` (xem [Bài 13](./13_type_hints.md)) đóng đúng vai trò `interface` của Go — Service sẽ phụ thuộc vào `UserRepository` (Protocol), không phụ thuộc trực tiếp implementation cụ thể.

## 3. Repository layer — implement Protocol

```python
# app/repositories/user_sqlalchemy.py
from sqlalchemy.orm import Session
from app.domain.user import User, UserRepository

class SqlAlchemyUserRepository:  # KHÔNG cần khai báo "implements UserRepository" — Protocol tự khớp
	def __init__(self, session: Session):
		self.session = session

	def find_by_id(self, user_id: int) -> User | None:
		row = self.session.get(UserModel, user_id)  # UserModel: SQLAlchemy model — xem Bài 17
		if row is None:
			return None
		return User(id=row.id, name=row.name, email=row.email, role=row.role)

	def find_by_email(self, email: str) -> User | None:
		row = self.session.query(UserModel).filter_by(email=email).first()
		return User(id=row.id, name=row.name, email=row.email, role=row.role) if row else None

	def create(self, user: User) -> User:
		row = UserModel(name=user.name, email=user.email, role=user.role)
		self.session.add(row)
		self.session.commit()
		user.id = row.id
		return user
```

## 4. Service layer — business logic, phụ thuộc Protocol

```python
# app/services/user_service.py
from app.domain.user import User, UserRepository

class UserService:
	def __init__(self, repo: UserRepository):  # dependency injection qua constructor
		self.repo = repo

	def get_user(self, user_id: int) -> User:
		user = self.repo.find_by_id(user_id)
		if user is None:
			raise ValueError(f"Không tìm thấy user {user_id}")
		return user
```

## 5. Router layer (FastAPI) — chi tiết đầy đủ ở [Bài 18](./18_rest_api.md)

```python
# app/api/user_router.py
from fastapi import APIRouter, Depends, HTTPException
from app.services.user_service import UserService

router = APIRouter()

def get_user_service() -> UserService:
	# nơi DUY NHẤT biết implementation cụ thể (SQLAlchemy) — tương đương main.go của Go
	from app.repositories.user_sqlalchemy import SqlAlchemyUserRepository
	from app.core.database import get_session
	return UserService(SqlAlchemyUserRepository(next(get_session())))

@router.get("/users/{user_id}")
def get_user(user_id: int, service: UserService = Depends(get_user_service)):
	try:
		return service.get_user(user_id)
	except ValueError as e:
		raise HTTPException(status_code=404, detail=str(e))
```

`Depends()` là hệ thống dependency injection built-in của FastAPI — tương đương việc "tiêm" `UserRepository` cụ thể vào `UserService` trong `main.go` của Go.

## 6. Test service mà không cần DB thật (nối tiếp [Bài 14](./14_testing.md))

```python
# tests/test_user_service.py
from unittest.mock import Mock
from app.domain.user import User, UserRepository
from app.services.user_service import UserService

def test_get_user_found():
	mock_repo = Mock(spec=UserRepository)
	mock_repo.find_by_id.return_value = User(id=1, name="Ben", email="ben@x.com", role="user")

	service = UserService(mock_repo)
	user = service.get_user(1)

	assert user.name == "Ben"

def test_get_user_not_found():
	mock_repo = Mock(spec=UserRepository)
	mock_repo.find_by_id.return_value = None

	service = UserService(mock_repo)
	try:
		service.get_user(999)
		assert False, "phải raise ValueError"
	except ValueError:
		pass
```

## 7. Cấu trúc thư mục đầy đủ (dùng cho capstone — [Bài 21](./21_capstone_project.md))

```
app/
├── main.py
├── core/              # config, security (hash, jwt)
├── domain/             # models + Protocol repository
├── api/                # routers — tương đương handler của Go
├── services/           # business logic
├── repositories/        # SQLAlchemy implementation
└── middleware hoặc dependencies.py  # auth dependency, RBAC dependency
```

## Bài tập

1. **Tách 3 tầng cho `Shape`**: thiết kế `domain` (dataclass + Protocol `ShapeRepository` giả lập lưu trữ), `service` (tính tổng diện tích), `router`/CLI giả lập gọi service.
2. **Mock repository test**: viết `InMemoryUserRepository` (dùng `dict` thay database thật) implement cùng `Protocol`, dùng nó test `UserService` không cần DB.
3. **So sánh với Go**: viết 1 đoạn ghi chú (comment hoặc file riêng) đối chiếu: `Protocol` ↔ `interface`, `Depends()` ↔ constructor injection thủ công trong `main.go`.

## Tiếp theo
→ [Bài 17: Làm việc với Database](./17_database.md)
