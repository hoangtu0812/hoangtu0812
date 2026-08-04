# Bài 20: Logging, Config, Deployment

## Mục tiêu
- Structured logging bằng `logging` module.
- Quản lý config bằng Pydantic Settings.
- Dockerize FastAPI, chạy bằng Uvicorn/Gunicorn, graceful shutdown.

## 1. Structured logging với `logging`

```python
import logging
import json

class JSONFormatter(logging.Formatter):
	def format(self, record):
		log_obj = {
			"time": self.formatTime(record),
			"level": record.levelname,
			"message": record.getMessage(),
			"logger": record.name,
		}
		return json.dumps(log_obj, ensure_ascii=False)

logger = logging.getLogger("taskapi")
logger.setLevel(logging.INFO)
handler = logging.StreamHandler()
handler.setFormatter(JSONFormatter())
logger.addHandler(handler)

logger.info("Server đang khởi động")
logger.warning("Connection pool gần đầy")
logger.error("Không thể kết nối database", exc_info=True)  # exc_info=True kèm traceback nếu trong except
```

Trong thực tế, thư viện `structlog` phổ biến hơn cho log có cấu trúc phức tạp, nhưng `logging` built-in đã đủ dùng cho hầu hết dự án — tương đương `log/slog` của Go ([Go Bài 20](../Go/20_logging_config_deploy.md)).

### Middleware log mỗi request (nối tiếp [Bài 18](./18_rest_api.md))

```python
import time
from fastapi import Request

@app.middleware("http")
async def logging_middleware(request: Request, call_next):
	start = time.perf_counter()
	response = await call_next(request)
	logger.info(json.dumps({
		"method": request.method,
		"path": request.url.path,
		"status_code": response.status_code,
		"duration_ms": round((time.perf_counter() - start) * 1000, 2),
	}, ensure_ascii=False))
	return response
```

## 2. Quản lý config với Pydantic Settings

```powershell
pip install pydantic-settings python-dotenv
```

```python
# app/core/config.py
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
	port: int = 8080
	database_url: str
	jwt_secret: str

	class Config:
		env_file = ".env"

settings = Settings()  # tự đọc từ .env + biến môi trường hệ thống, validate kiểu tự động
```

```
# .env
PORT=8080
DATABASE_URL=postgresql://user:pass@localhost:5432/taskdb
JWT_SECRET=doi-secret-nay-trong-production
```

Pydantic Settings **validate kiểu tự động** (vd `port` phải convert được sang `int`, nếu không sẽ raise lỗi ngay khi khởi động) — an toàn hơn đọc `os.getenv()` thủ công, gần tương đương struct `Config` của Go ([Go Bài 20](../Go/20_logging_config_deploy.md)).

**`.env` KHÔNG commit vào git** — thêm vào `.gitignore`, chỉ commit `.env.example`.

## 3. Lifespan — thay thế `@app.on_event` cũ, tương đương graceful shutdown của Go

```python
from contextlib import asynccontextmanager
from fastapi import FastAPI

@asynccontextmanager
async def lifespan(app: FastAPI):
	# Code TRƯỚC yield chạy lúc STARTUP
	logger.info("Server đang khởi động")
	yield
	# Code SAU yield chạy lúc SHUTDOWN — nơi đóng connection pool, dọn dẹp resource
	logger.info("Server đang tắt, dọn dẹp resource...")

app = FastAPI(lifespan=lifespan)
```

Uvicorn tự xử lý tín hiệu `SIGTERM`/`SIGINT` để gọi phần shutdown của `lifespan`, tương đương `srv.Shutdown(ctx)` của Go.

## 4. Dockerize FastAPI

```dockerfile
# Dockerfile
FROM python:3.12-slim AS base

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8080
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
```

Với production, nên dùng **Gunicorn quản lý nhiều Uvicorn worker** để tận dụng đa core (vì GIL — xem [Bài 15](./15_concurrency.md)):

```dockerfile
CMD ["gunicorn", "main:app", "-w", "4", "-k", "uvicorn.workers.UvicornWorker", "-b", "0.0.0.0:8080"]
```

## 5. `docker-compose.yml`

```yaml
version: "3.9"
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://user:pass@db:5432/taskdb
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: taskdb
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

Chạy: `docker-compose up --build`

## Bài tập

1. **Structured logging**: thay `print()` trong 1 project nhỏ đã viết trước đó bằng `logging` + middleware log request.
2. **Pydantic Settings**: viết `Settings` class, tạo `.env`, đảm bảo `.env` nằm trong `.gitignore`.
3. **`lifespan`**: thêm log "startup"/"shutdown" vào `lifespan`, chạy `uvicorn` và `Ctrl+C` để quan sát log shutdown xuất hiện.
4. **Dockerize**: viết `Dockerfile`, build image, chạy `docker run -p 8080:8080 <image>` và gọi thử `/docs`.

## Tổng kết Giai đoạn 3
Bạn đã có đủ kiến thức để xây REST API hoàn chỉnh bằng FastAPI: kiến trúc phân tầng, database, routing/middleware, auth/RBAC, logging/config/deploy — cùng bộ kỹ năng với track Go, khác framework/ngôn ngữ.

## Tiếp theo
→ [Dự án Capstone: FastAPI Task API có phân quyền](./21_capstone_project.md)
