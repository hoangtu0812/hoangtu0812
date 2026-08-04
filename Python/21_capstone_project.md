# Dự án Capstone: FastAPI Task API Có Phân Quyền

## Mục tiêu
Ghép toàn bộ kiến thức từ Bài 1 → Bài 20 thành 1 dự án hoàn chỉnh — **cùng ngữ cảnh nghiệp vụ với dự án Go** ([Go/21_capstone_project.md](../Go/21_capstone_project.md)) để bạn so sánh trực tiếp 2 stack trên cùng 1 bài toán.

**Kiến thức áp dụng:** [Bài 16](./16_clean_architecture.md) (kiến trúc), [Bài 17](./17_database.md) (database), [Bài 18](./18_rest_api.md) (REST API), [Bài 19](./19_auth.md) (auth/RBAC), [Bài 20](./20_logging_config_deploy.md) (logging/deploy).

## Yêu cầu chức năng
Giống hệt bản Go: `user` đăng ký/đăng nhập, CRUD task của chính mình; `admin` quản lý mọi task và user; mật khẩu hash bcrypt, JWT xác thực.

## Cấu trúc thư mục

```
taskapi/
├── app/
│   ├── main.py
│   ├── core/
│   │   ├── config.py         # Pydantic Settings
│   │   ├── database.py        # SQLAlchemy engine/session
│   │   └── security.py        # hash password, JWT
│   ├── domain/
│   │   ├── user.py            # dataclass User + Protocol UserRepository
│   │   └── task.py            # dataclass Task + Protocol TaskRepository
│   ├── models.py               # SQLAlchemy ORM model (UserModel, TaskModel)
│   ├── repositories/
│   │   ├── user_sqlalchemy.py
│   │   └── task_sqlalchemy.py
│   ├── services/
│   │   ├── auth_service.py
│   │   └── task_service.py
│   ├── api/
│   │   ├── deps.py             # get_current_user, require_role
│   │   ├── auth_router.py
│   │   ├── task_router.py
│   │   └── admin_router.py
│   └── schemas/                # Pydantic request/response model
│       ├── user_schema.py
│       └── task_schema.py
├── tests/
├── alembic/
├── .env.example
├── requirements.txt
├── Dockerfile
└── docker-compose.yml
```

## Bước 1: Bootstrap

```powershell
mkdir taskapi; cd taskapi
python -m venv venv
.\venv\Scripts\Activate.ps1
pip install fastapi "uvicorn[standard]" sqlalchemy psycopg2-binary alembic "passlib[bcrypt]" "python-jose[cryptography]" pydantic-settings python-dotenv pytest httpx
```

## Bước 2: Domain (`app/domain/`)

```python
# app/domain/user.py
from dataclasses import dataclass
from typing import Protocol

@dataclass
class User:
	id: int
	name: str
	email: str
	password_hash: str
	role: str = "user"

class UserRepository(Protocol):
	def create(self, user: User) -> User: ...
	def find_by_id(self, user_id: int) -> User | None: ...
	def find_by_email(self, email: str) -> User | None: ...
	def list_all(self) -> list[User]: ...
```

```python
# app/domain/task.py
from dataclasses import dataclass
from typing import Protocol

@dataclass
class Task:
	id: int
	owner_id: int
	title: str
	done: bool = False

class TaskRepository(Protocol):
	def create(self, task: Task) -> Task: ...
	def find_by_id(self, task_id: int) -> Task | None: ...
	def list_by_owner(self, owner_id: int) -> list[Task]: ...
	def list_all(self) -> list[Task]: ...
	def update(self, task: Task) -> Task: ...
	def delete(self, task_id: int) -> None: ...
```

## Bước 3: `core/security.py`

```python
from datetime import datetime, timedelta, timezone
from passlib.context import CryptContext
from jose import jwt, JWTError
from app.core.config import settings

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

def hash_password(password: str) -> str:
	return pwd_context.hash(password)

def verify_password(plain: str, hashed: str) -> bool:
	return pwd_context.verify(plain, hashed)

def create_access_token(user_id: int, role: str) -> str:
	payload = {
		"sub": str(user_id),
		"role": role,
		"exp": datetime.now(timezone.utc) + timedelta(hours=24),
	}
	return jwt.encode(payload, settings.jwt_secret, algorithm="HS256")

def verify_access_token(token: str) -> dict:
	try:
		return jwt.decode(token, settings.jwt_secret, algorithms=["HS256"])
	except JWTError as e:
		raise ValueError("token không hợp lệ") from e
```

## Bước 4: Repository (SQLAlchemy)

```python
# app/repositories/task_sqlalchemy.py
from sqlalchemy.orm import Session
from app.domain.task import Task
from app.models import TaskModel

class SqlAlchemyTaskRepository:
	def __init__(self, session: Session):
		self.session = session

	def _to_domain(self, row: TaskModel) -> Task:
		return Task(id=row.id, owner_id=row.owner_id, title=row.title, done=row.done)

	def create(self, task: Task) -> Task:
		row = TaskModel(owner_id=task.owner_id, title=task.title, done=task.done)
		self.session.add(row)
		self.session.commit()
		self.session.refresh(row)
		task.id = row.id
		return task

	def find_by_id(self, task_id: int) -> Task | None:
		row = self.session.get(TaskModel, task_id)
		return self._to_domain(row) if row else None

	def list_by_owner(self, owner_id: int) -> list[Task]:
		rows = self.session.query(TaskModel).filter_by(owner_id=owner_id).all()
		return [self._to_domain(r) for r in rows]

	def list_all(self) -> list[Task]:
		return [self._to_domain(r) for r in self.session.query(TaskModel).all()]

	def update(self, task: Task) -> Task:
		row = self.session.get(TaskModel, task.id)
		row.title = task.title
		row.done = task.done
		self.session.commit()
		return task

	def delete(self, task_id: int) -> None:
		row = self.session.get(TaskModel, task_id)
		if row:
			self.session.delete(row)
			self.session.commit()
```

`user_sqlalchemy.py` viết tương tự (xem mẫu ở [Bài 16](./16_clean_architecture.md) và [Bài 17](./17_database.md)).

## Bước 5: Service — chứa logic phân quyền ownership

```python
# app/services/task_service.py
from app.domain.task import Task, TaskRepository

class ForbiddenError(Exception):
	pass

class TaskService:
	def __init__(self, repo: TaskRepository):
		self.repo = repo

	def create_task(self, owner_id: int, title: str) -> Task:
		return self.repo.create(Task(id=0, owner_id=owner_id, title=title))

	def list_tasks(self, requester_id: int, requester_role: str) -> list[Task]:
		if requester_role == "admin":
			return self.repo.list_all()
		return self.repo.list_by_owner(requester_id)

	def update_task(self, requester_id: int, requester_role: str, task_id: int, title: str, done: bool) -> Task:
		task = self.repo.find_by_id(task_id)
		if task is None:
			raise ValueError("không tìm thấy task")

		if requester_role != "admin" and task.owner_id != requester_id:
			raise ForbiddenError("bạn không có quyền sửa task này")

		task.title = title
		task.done = done
		return self.repo.update(task)

	def delete_task(self, requester_id: int, requester_role: str, task_id: int) -> None:
		task = self.repo.find_by_id(task_id)
		if task is None:
			raise ValueError("không tìm thấy task")
		if requester_role != "admin" and task.owner_id != requester_id:
			raise ForbiddenError("bạn không có quyền xóa task này")
		self.repo.delete(task_id)
```

```python
# app/services/auth_service.py
from app.domain.user import User, UserRepository
from app.core.security import hash_password, verify_password, create_access_token

class AuthService:
	def __init__(self, repo: UserRepository):
		self.repo = repo

	def register(self, name: str, email: str, password: str) -> User:
		if self.repo.find_by_email(email):
			raise ValueError("email đã tồn tại")
		user = User(id=0, name=name, email=email, password_hash=hash_password(password), role="user")
		return self.repo.create(user)

	def login(self, email: str, password: str) -> str:
		user = self.repo.find_by_email(email)
		if user is None or not verify_password(password, user.password_hash):
			raise ValueError("email hoặc mật khẩu không đúng")
		return create_access_token(user.id, user.role)
```

## Bước 6-7: Dependencies xác thực + phân quyền

```python
# app/api/deps.py
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from sqlalchemy.orm import Session
from app.core.database import get_session
from app.core.security import verify_access_token
from app.repositories.task_sqlalchemy import SqlAlchemyTaskRepository
from app.services.task_service import TaskService

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="auth/login")

def get_current_user(token: str = Depends(oauth2_scheme)) -> dict:
	try:
		payload = verify_access_token(token)
	except ValueError:
		raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Token không hợp lệ")
	return {"user_id": int(payload["sub"]), "role": payload["role"]}

def require_role(*roles: str):
	def checker(current_user: dict = Depends(get_current_user)) -> dict:
		if current_user["role"] not in roles:
			raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Không đủ quyền")
		return current_user
	return checker

def get_task_service(session: Session = Depends(get_session)) -> TaskService:
	return TaskService(SqlAlchemyTaskRepository(session))
```

## Bước 8: Router

```python
# app/api/task_router.py
from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel, Field
from app.api.deps import get_current_user, get_task_service
from app.services.task_service import TaskService, ForbiddenError

router = APIRouter(prefix="/tasks", tags=["tasks"])

class TaskCreate(BaseModel):
	title: str = Field(min_length=1, max_length=255)

class TaskUpdate(BaseModel):
	title: str = Field(min_length=1, max_length=255)
	done: bool

@router.post("")
def create_task(payload: TaskCreate, current_user: dict = Depends(get_current_user), service: TaskService = Depends(get_task_service)):
	task = service.create_task(current_user["user_id"], payload.title)
	return task

@router.get("")
def list_tasks(current_user: dict = Depends(get_current_user), service: TaskService = Depends(get_task_service)):
	return service.list_tasks(current_user["user_id"], current_user["role"])

@router.put("/{task_id}")
def update_task(task_id: int, payload: TaskUpdate, current_user: dict = Depends(get_current_user), service: TaskService = Depends(get_task_service)):
	try:
		return service.update_task(current_user["user_id"], current_user["role"], task_id, payload.title, payload.done)
	except ForbiddenError as e:
		raise HTTPException(status_code=403, detail=str(e))
	except ValueError as e:
		raise HTTPException(status_code=404, detail=str(e))

@router.delete("/{task_id}", status_code=204)
def delete_task(task_id: int, current_user: dict = Depends(get_current_user), service: TaskService = Depends(get_task_service)):
	try:
		service.delete_task(current_user["user_id"], current_user["role"], task_id)
	except ForbiddenError as e:
		raise HTTPException(status_code=403, detail=str(e))
	except ValueError as e:
		raise HTTPException(status_code=404, detail=str(e))
```

## Bước 9: `main.py`

```python
from contextlib import asynccontextmanager
from fastapi import FastAPI
from app.api import auth_router, task_router, admin_router

@asynccontextmanager
async def lifespan(app: FastAPI):
	print("Server đang khởi động")
	yield
	print("Server đang tắt")

app = FastAPI(title="Task API", lifespan=lifespan)
app.include_router(auth_router.router)
app.include_router(task_router.router)
app.include_router(admin_router.router)
```

## Bước 10: Test service (không cần DB)

```python
# tests/test_task_service.py
from unittest.mock import Mock
import pytest
from app.domain.task import Task, TaskRepository
from app.services.task_service import TaskService, ForbiddenError

@pytest.fixture
def mock_repo():
	repo = Mock(spec=TaskRepository)
	repo.find_by_id.return_value = Task(id=1, owner_id=100, title="Task của user 100")
	return repo

def test_admin_can_update_others_task(mock_repo):
	service = TaskService(mock_repo)
	service.update_task(999, "admin", 1, "Sửa bởi admin", True)  # không raise

def test_owner_can_update_own_task(mock_repo):
	service = TaskService(mock_repo)
	service.update_task(100, "user", 1, "Sửa bởi chủ", True)  # không raise

def test_other_user_cannot_update_task(mock_repo):
	service = TaskService(mock_repo)
	with pytest.raises(ForbiddenError):
		service.update_task(200, "user", 1, "Cố sửa", True)
```

## Bước 11-12: Dockerize (xem mẫu đầy đủ ở [Bài 20](./20_logging_config_deploy.md))

## Bước 13: Test thủ công

```powershell
# Đăng ký
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d '{\"name\":\"Ben\",\"email\":\"ben@test.com\",\"password\":\"12345678\"}'

# Đăng nhập
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{\"email\":\"ben@test.com\",\"password\":\"12345678\"}'

# Hoặc đơn giản hơn: mở http://localhost:8080/docs và test qua Swagger UI
```

## Checklist hoàn thành
- [ ] `POST /auth/register`, `POST /auth/login` hoạt động, password hash bằng bcrypt.
- [ ] `user` không sửa/xóa được task người khác; `admin` thao tác được mọi task.
- [ ] `pytest` chạy pass, có test ownership như Bước 10.
- [ ] `docker-compose up --build` chạy được, `/docs` truy cập được.
- [ ] README mô tả kiến trúc, hướng dẫn chạy.

## So sánh với bản Go
Sau khi hoàn thành, hãy đối chiếu trực tiếp với [Go/21_capstone_project.md](../Go/21_capstone_project.md):
- `Protocol` (Python) ↔ `interface` (Go) — cùng ý tưởng, Go kiểm tra lúc compile, Python chỉ kiểm tra nếu chạy `mypy`.
- `Depends()` (FastAPI) ↔ constructor injection thủ công trong `main.go`.
- Pydantic model ↔ struct + json tag + `validator`.
- `pytest` + `Mock` ↔ `go test` + mock implementation viết tay.
- Tốc độ phát triển (Python nhanh hơn nhờ ít boilerplate) đổi lấy hiệu năng runtime (Go nhanh hơn, không cần GIL).

---
Hoàn thành bài này nghĩa là bạn có 2 stack backend độc lập (Go + Python) cùng giải 1 bài toán — tài sản rất mạnh khi phỏng vấn hoặc chọn công nghệ cho dự án thực tế.
