# Bài 17: Làm việc với Database

## Mục tiêu
- Dùng SQLAlchemy (ORM phổ biến nhất Python) để định nghĩa model và CRUD.
- Quản lý migration với Alembic.
- Kết nối PostgreSQL.

## 1. Cài đặt

```powershell
pip install sqlalchemy psycopg2-binary alembic
```

## 2. Định nghĩa model (SQLAlchemy 2.0 style — dùng `Mapped`/`mapped_column`)

```python
# app/core/database.py
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, DeclarativeBase

DATABASE_URL = "postgresql://user:pass@localhost:5432/taskdb"

engine = create_engine(DATABASE_URL, pool_size=10, max_overflow=20)
SessionLocal = sessionmaker(bind=engine, autoflush=False)

class Base(DeclarativeBase):
	pass

def get_session():
	session = SessionLocal()
	try:
		yield session   # dùng với Depends() của FastAPI — tự đóng session sau khi xử lý xong request
	finally:
		session.close()
```

```python
# app/models.py
from sqlalchemy import String, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.core.database import Base

class UserModel(Base):
	__tablename__ = "users"

	id: Mapped[int] = mapped_column(primary_key=True)
	name: Mapped[str] = mapped_column(String(255))
	email: Mapped[str] = mapped_column(String(255), unique=True)
	password_hash: Mapped[str] = mapped_column(String(255))
	role: Mapped[str] = mapped_column(String(20), default="user")

	tasks: Mapped[list["TaskModel"]] = relationship(back_populates="owner")

class TaskModel(Base):
	__tablename__ = "tasks"

	id: Mapped[int] = mapped_column(primary_key=True)
	owner_id: Mapped[int] = mapped_column(ForeignKey("users.id"))
	title: Mapped[str] = mapped_column(String(255))
	done: Mapped[bool] = mapped_column(default=False)

	owner: Mapped["UserModel"] = relationship(back_populates="tasks")
```

`Mapped[int]`/`mapped_column` dùng type hint ([Bài 13](./13_type_hints.md)) để định nghĩa schema — tương tự struct tag `db:"..."` nếu Go dùng ORM (GORM) thay vì `database/sql` thuần ([Go Bài 17](../Go/17_database.md)).

## 3. CRUD cơ bản

```python
from sqlalchemy.orm import Session
from app.models import UserModel

def create_user(session: Session, name: str, email: str, password_hash: str) -> UserModel:
	user = UserModel(name=name, email=email, password_hash=password_hash)
	session.add(user)
	session.commit()
	session.refresh(user)  # lấy lại id vừa được DB sinh ra
	return user

def get_user_by_id(session: Session, user_id: int) -> UserModel | None:
	return session.get(UserModel, user_id)

def list_users(session: Session) -> list[UserModel]:
	return session.query(UserModel).all()

def update_user_email(session: Session, user_id: int, email: str) -> None:
	user = session.get(UserModel, user_id)
	if user is None:
		raise ValueError("không tìm thấy user")
	user.email = email
	session.commit()

def delete_user(session: Session, user_id: int) -> None:
	user = session.get(UserModel, user_id)
	if user:
		session.delete(user)
		session.commit()
```

## 4. Query nâng cao — JOIN, filter

```python
from sqlalchemy import select

def get_tasks_by_owner(session: Session, owner_id: int) -> list["TaskModel"]:
	stmt = select(TaskModel).where(TaskModel.owner_id == owner_id)
	return list(session.scalars(stmt))

def get_undone_tasks_with_owner_name(session: Session):
	stmt = (
		select(TaskModel.title, UserModel.name)
		.join(UserModel, TaskModel.owner_id == UserModel.id)
		.where(TaskModel.done == False)
	)
	return session.execute(stmt).all()
```

**Luôn dùng query builder của SQLAlchemy (không nối chuỗi SQL thủ công)** — tự động tham số hóa, tránh SQL Injection, giống nguyên tắc `$1, $2` của Go ([Go Bài 17](../Go/17_database.md)).

## 5. Transaction

```python
def transfer_balance(session: Session, from_id: int, to_id: int, amount: float) -> None:
	try:
		from_account = session.get(AccountModel, from_id)
		to_account = session.get(AccountModel, to_id)
		from_account.balance -= amount
		to_account.balance += amount
		session.commit()  # commit cả 2 thay đổi cùng lúc
	except Exception:
		session.rollback()  # rollback nếu có lỗi giữa chừng
		raise
```

## 6. Migration với Alembic

```powershell
alembic init alembic
```

Cấu hình `alembic.ini` trỏ `sqlalchemy.url` tới database, rồi:

```powershell
alembic revision --autogenerate -m "create users and tasks table"
alembic upgrade head
```

File migration tự sinh (`alembic/versions/xxxx_create_users_and_tasks_table.py`):

```python
def upgrade():
	op.create_table(
		"users",
		sa.Column("id", sa.Integer, primary_key=True),
		sa.Column("name", sa.String(255), nullable=False),
		sa.Column("email", sa.String(255), unique=True, nullable=False),
		sa.Column("password_hash", sa.String(255), nullable=False),
		sa.Column("role", sa.String(20), server_default="user"),
	)

def downgrade():
	op.drop_table("users")
```

## Bài tập

1. **Kết nối DB**: cài PostgreSQL (hoặc Docker), cấu hình `DATABASE_URL`, kiểm tra kết nối bằng `engine.connect()`.
2. **CRUD `users`**: định nghĩa `UserModel`, viết đủ 4 hàm CRUD như ví dụ.
3. **Migration**: dùng Alembic tạo migration cho bảng `users` và `tasks`, chạy `alembic upgrade head`.
4. **Nâng cao**: viết `SqlAlchemyUserRepository` implement `Protocol UserRepository` từ [Bài 16](./16_clean_architecture.md), viết test dùng SQLite in-memory (`create_engine("sqlite:///:memory:")`) thay vì Postgres thật để test nhanh.

## Tiếp theo
→ [Bài 18: REST API với FastAPI](./18_rest_api.md)
