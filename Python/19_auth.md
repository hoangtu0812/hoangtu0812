# Bài 19: Authentication & Authorization

## Mục tiêu
- Hash password bằng `passlib`/`bcrypt`.
- Tạo/verify JWT bằng `python-jose`.
- Dependency xác thực + phân quyền RBAC trong FastAPI.

## 1. Hash password

```powershell
pip install "passlib[bcrypt]"
```

```python
from passlib.context import CryptContext

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

def hash_password(password: str) -> str:
	return pwd_context.hash(password)

def verify_password(plain: str, hashed: str) -> bool:
	return pwd_context.verify(plain, hashed)
```

## 2. Tạo & verify JWT

```powershell
pip install "python-jose[cryptography]"
```

```python
from datetime import datetime, timedelta, timezone
from jose import jwt, JWTError

JWT_SECRET = "doc-tu-env-KHONG-hardcode"
JWT_ALGORITHM = "HS256"

def create_access_token(user_id: int, role: str) -> str:
	payload = {
		"sub": str(user_id),
		"role": role,
		"exp": datetime.now(timezone.utc) + timedelta(hours=24),
	}
	return jwt.encode(payload, JWT_SECRET, algorithm=JWT_ALGORITHM)

def verify_access_token(token: str) -> dict:
	try:
		return jwt.decode(token, JWT_SECRET, algorithms=[JWT_ALGORITHM])
	except JWTError as e:
		raise ValueError("token không hợp lệ hoặc đã hết hạn") from e
```

## 3. Flow đăng ký / đăng nhập

```python
from app.domain.user import User, UserRepository

class AuthService:
	def __init__(self, repo: UserRepository):
		self.repo = repo

	def register(self, name: str, email: str, password: str) -> User:
		if self.repo.find_by_email(email) is not None:
			raise ValueError("email đã tồn tại")

		user = User(id=0, name=name, email=email, role="user")
		user.password_hash = hash_password(password)  # giả sử User có field password_hash
		return self.repo.create(user)

	def login(self, email: str, password: str) -> str:
		user = self.repo.find_by_email(email)
		if user is None or not verify_password(password, user.password_hash):
			raise ValueError("email hoặc mật khẩu không đúng")  # KHÔNG tiết lộ cái nào sai
		return create_access_token(user.id, user.role)
```

## 4. Dependency xác thực trong FastAPI — `OAuth2PasswordBearer`

```python
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="auth/login")  # chỉ để Swagger UI hiện nút "Authorize"

def get_current_user(token: str = Depends(oauth2_scheme)) -> dict:
	try:
		payload = verify_access_token(token)
	except ValueError:
		raise HTTPException(
			status_code=status.HTTP_401_UNAUTHORIZED,
			detail="Token không hợp lệ",
			headers={"WWW-Authenticate": "Bearer"},
		)
	return {"user_id": int(payload["sub"]), "role": payload["role"]}

@app.get("/profile")
def get_profile(current_user: dict = Depends(get_current_user)):
	return current_user
```

`Depends(get_current_user)` tương đương middleware xác thực của Gin ([Go Bài 19](../Go/19_auth.md)) — nhưng ở FastAPI, xác thực là 1 **dependency** gắn trực tiếp vào từng route (hoặc router) thay vì middleware toàn cục, cho phép áp dụng linh hoạt route nào cần, route nào không.

## 5. Dependency phân quyền (RBAC)

```python
def require_role(*roles: str):
	def checker(current_user: dict = Depends(get_current_user)) -> dict:
		if current_user["role"] not in roles:
			raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Không đủ quyền truy cập")
		return current_user
	return checker

@app.get("/admin/users")
def list_all_users(current_user: dict = Depends(require_role("admin"))):
	...
```

`require_role("admin")` là 1 "dependency factory" — trả về dependency mới mỗi lần gọi, tương tự decorator có tham số ([Bài 11](./11_decorators_context_managers.md)) và middleware `RequireRole(...)` của Gin ([Go Bài 19](../Go/19_auth.md)).

## 6. Phân quyền theo ownership — logic nằm ở tầng service

```python
class TaskService:
	def __init__(self, repo):
		self.repo = repo

	def update_task(self, requester_id: int, requester_role: str, task_id: int, title: str, done: bool):
		task = self.repo.find_by_id(task_id)
		if task is None:
			raise ValueError("không tìm thấy task")

		if requester_role != "admin" and task.owner_id != requester_id:
			raise PermissionError("bạn không có quyền sửa task này")

		task.title = title
		task.done = done
		return self.repo.update(task)
```

```python
@app.put("/tasks/{task_id}")
def update_task(task_id: int, payload: TaskUpdate, current_user: dict = Depends(get_current_user)):
	try:
		return task_service.update_task(
			current_user["user_id"], current_user["role"], task_id, payload.title, payload.done
		)
	except PermissionError as e:
		raise HTTPException(status_code=403, detail=str(e))
	except ValueError as e:
		raise HTTPException(status_code=404, detail=str(e))
```

Đúng theo nguyên tắc ở [Go Bài 19](../Go/19_auth.md): logic ownership nằm ở **service**, dependency `require_role` chỉ kiểm tra role chung chung.

## Bài tập

1. **Hash & verify password**: viết `hash_password`/`verify_password` bằng `passlib`, test với password đúng/sai.
2. **Tạo & verify JWT**: viết `create_access_token`/`verify_access_token`, test với token hết hạn (đặt `exp` trong quá khứ).
3. **Dependency auth + RBAC**: viết `get_current_user` và `require_role`, tạo 3 route (public, cần login, cần admin), test qua `/docs` bằng nút "Authorize".
4. **Nâng cao — ownership check**: viết `update_task` như ví dụ, viết test (dùng mock repository) cho 3 case: admin sửa task người khác, user sửa task của mình, user sửa task người khác bị từ chối — giống bài tập ở [Go Bài 19](../Go/19_auth.md).

## Tiếp theo
→ [Bài 20: Logging, Config, Deployment](./20_logging_config_deploy.md)
