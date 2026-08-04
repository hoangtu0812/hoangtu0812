# Bài 18: REST API với FastAPI

## Mục tiêu
- Viết REST API bằng FastAPI: routing, Pydantic model, dependency injection, middleware.
- So sánh nhanh với Gin của Go ([Go Bài 18](../Go/18_rest_api.md)).

## 1. Cài đặt & app tối thiểu

```powershell
pip install fastapi "uvicorn[standard]"
```

```python
# main.py
from fastapi import FastAPI

app = FastAPI(title="Todo API")

@app.get("/")
def root():
	return {"message": "Hello, FastAPI!"}
```

Chạy: `uvicorn main:app --reload` — FastAPI **tự sinh Swagger UI** tại `/docs` (khác Go/Gin cần thêm thư viện `swaggo` — [Go Bài 18](../Go/18_rest_api.md)).

## 2. Pydantic model — tương đương struct + json tag + validator của Go

```python
from pydantic import BaseModel, Field

class TodoCreate(BaseModel):
	title: str = Field(min_length=1, max_length=255)
	done: bool = False

class TodoResponse(BaseModel):
	id: int
	title: str
	done: bool

	model_config = {"from_attributes": True}  # cho phép tạo từ SQLAlchemy model (ORM object)
```

Pydantic tự động **validate** request body theo type hint + `Field` constraint — tương đương `binding:"required,min=..."` của Gin ([Go Bài 18](../Go/18_rest_api.md)), nhưng tích hợp sẵn, không cần thư viện riêng.

## 3. CRUD API đầy đủ

```python
from fastapi import FastAPI, HTTPException, status

app = FastAPI()

todos: dict[int, dict] = {1: {"id": 1, "text": "Học FastAPI", "done": False}}
next_id = 2

@app.get("/todos", response_model=list[TodoResponse])
def list_todos():
	return list(todos.values())

@app.post("/todos", response_model=TodoResponse, status_code=status.HTTP_201_CREATED)
def create_todo(payload: TodoCreate):
	global next_id
	todo = {"id": next_id, "title": payload.title, "done": payload.done}
	todos[next_id] = todo
	next_id += 1
	return todo

@app.get("/todos/{todo_id}", response_model=TodoResponse)
def get_todo(todo_id: int):
	todo = todos.get(todo_id)
	if todo is None:
		raise HTTPException(status_code=404, detail="Không tìm thấy todo")
	return todo

@app.put("/todos/{todo_id}", response_model=TodoResponse)
def update_todo(todo_id: int, payload: TodoCreate):
	if todo_id not in todos:
		raise HTTPException(status_code=404, detail="Không tìm thấy todo")
	todos[todo_id].update(payload.model_dump())
	return todos[todo_id]

@app.delete("/todos/{todo_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_todo(todo_id: int):
	if todo_id not in todos:
		raise HTTPException(status_code=404, detail="Không tìm thấy todo")
	del todos[todo_id]
```

`todo_id: int` trong path parameter tự động parse + validate kiểu (nếu client gửi `/todos/abc`, FastAPI tự trả `422 Unprocessable Entity` — không cần `strconv.Atoi` thủ công như Go).

## 4. Dependency Injection với `Depends()`

```python
from fastapi import Depends
from sqlalchemy.orm import Session
from app.core.database import get_session

@app.get("/todos/{todo_id}")
def get_todo(todo_id: int, session: Session = Depends(get_session)):
	# session được FastAPI tự "tiêm" vào, tự đóng sau khi request xử lý xong
	...
```

`Depends()` là cơ chế trung tâm của FastAPI — dùng cho cả kết nối DB, xác thực (xem [Bài 19](./19_auth.md)), và validate logic dùng chung.

## 5. Middleware

```python
import time
from fastapi import Request

@app.middleware("http")
async def logging_middleware(request: Request, call_next):
	start = time.perf_counter()
	response = await call_next(request)  # gọi tiếp handler/middleware kế tiếp trong chuỗi
	elapsed = time.perf_counter() - start
	print(f"{request.method} {request.url.path} - {response.status_code} - {elapsed:.4f}s")
	return response
```

## 6. Router riêng biệt (`APIRouter`) — chia nhỏ theo domain, giống router riêng của Go

```python
# app/api/todo_router.py
from fastapi import APIRouter

router = APIRouter(prefix="/todos", tags=["todos"])

@router.get("")
def list_todos():
	...

@router.post("")
def create_todo(payload: TodoCreate):
	...
```

```python
# main.py
from app.api.todo_router import router as todo_router

app.include_router(todo_router)
```

## 7. Response format chuẩn hóa & exception handler toàn cục

```python
from fastapi import Request
from fastapi.responses import JSONResponse

@app.exception_handler(ValueError)
async def value_error_handler(request: Request, exc: ValueError):
	return JSONResponse(status_code=400, content={"success": False, "error": str(exc)})
```

## Bài tập

1. **API Todo đầy đủ**: viết `GET/POST /todos`, `GET/PUT/DELETE /todos/{id}` theo code mẫu, test qua `/docs` (Swagger UI tự sinh).
2. **Pydantic validation**: thêm `Field` constraint cho `TodoCreate` (vd `title` không được rỗng), thử gửi request sai để xem FastAPI tự trả lỗi 422 chi tiết thế nào.
3. **Middleware logging**: viết middleware log method/path/status/thời gian xử lý cho mỗi request.
4. **Nâng cao**: tách route Todo ra `APIRouter` riêng trong package `app/api/`, dùng `app.include_router()`, viết exception handler toàn cục cho `ValueError` như ví dụ.

## Tiếp theo
→ [Bài 19: Authentication & Authorization](./19_auth.md)
