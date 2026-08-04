# Lộ Trình Học Python Chi Tiết

> Cấu trúc giống hệt [Go/ROADMAP.md](../Go/ROADMAP.md): mỗi bài có file chi tiết riêng (lý thuyết + code mẫu + bài tập), kết thúc bằng dự án capstone — REST API có phân quyền, lần này dùng FastAPI.

---

## Giai đoạn 0 — Giới thiệu & Cài đặt

### [Bài 1: Giới thiệu & Cài đặt](./1_get_started.md)
- Python là gì, dùng cho việc gì (scripting, data, AI/ML, backend web), so sánh nhanh với Go (typing động vs tĩnh, interpreter vs compiler).
- Cài Python, VS Code + extension Python/Pylance, tạo virtual environment, `pip`.
- **Bài tập:** Viết Hello World, chạy bằng `python`, tạo và kích hoạt venv, cài 1 package bằng `pip`.

---

## Giai đoạn 1 — Python Cơ Bản (Fundamentals)

### [Bài 2: Biến, kiểu dữ liệu, toán tử](./2_variables_types.md)
- Typing động, kiểu cơ bản (int, float, str, bool, None), toán tử số học/so sánh/logic, f-string.
- **Bài tập:** Đổi độ C↔F; in kiểu bằng `type()`; thử `is` vs `==`.

### [Bài 3: Luồng điều khiển](./3_control_flow.md)
- `if/elif/else`, `for`, `while`, `break/continue`, `range()`, list comprehension cơ bản, walrus operator `:=`.
- **Bài tập:** FizzBuzz; số nguyên tố; bảng cửu chương.

### [Bài 4: Hàm (Functions)](./4_functions.md)
- `def`, tham số mặc định, `*args`/`**kwargs`, lambda, closures, docstring.
- **Bài tập:** Hàm chia an toàn; hàm biến đối số; closure counter.

### [Bài 5: List, Tuple, Set, Dict](./5_collections.md)
- Cấu trúc dữ liệu built-in, slicing, mutable vs immutable, comprehension cho list/dict/set.
- **Bài tập:** Loại bỏ trùng lặp; đếm tần suất từ; đảo ngược list.

### [Bài 6: String & Comprehension nâng cao](./6_strings_comprehensions.md)
- Method của string, format, regex cơ bản (`re`), nested comprehension, generator expression.
- **Bài tập:** Đếm từ trong câu; kiểm tra palindrome; parse CSV đơn giản bằng `split`.

### [Bài 7: OOP trong Python](./7_oop.md)
- `class`, `__init__`, kế thừa, đa hình, encapsulation (`_private`/`__name mangling`), dunder methods (`__str__`, `__eq__`), `@property`, dataclass.
- **Bài tập:** `Rectangle`/`Circle` kế thừa `Shape`; dùng `@dataclass` cho `Point`.

### [Bài 8: Xử lý lỗi (Exception Handling)](./8_exceptions.md)
- `try/except/else/finally`, custom exception, `raise ... from`, context manager cơ bản (`with`).
- **Bài tập:** Custom `ValidationError`; hàm dùng `with open()` an toàn.

### [Bài 9: Module, Package & Virtual Environment](./9_modules_packages.md)
- `import`, package (`__init__.py`), `pip`/`requirements.txt`, venv/`poetry`, cấu trúc project chuẩn.
- **Bài tập:** Tách code thành package riêng, viết `requirements.txt`.

---

## Giai đoạn 2 — Python Trung Cấp (Intermediate)

### [Bài 10: Iterator & Generator](./10_iterators_generators.md)
- Protocol `__iter__`/`__next__`, `yield`, generator expression, lazy evaluation.
- **Bài tập:** Generator sinh dãy Fibonacci; iterator tự định nghĩa.

### [Bài 11: Decorator & Context Manager](./11_decorators_context_managers.md)
- Hàm bậc cao, decorator cơ bản/có tham số, `functools.wraps`, context manager (`with`, `contextlib.contextmanager`).
- **Bài tập:** Decorator đo thời gian chạy; decorator retry; context manager tự viết.

### [Bài 12: File I/O & JSON/CSV](./12_file_io_json.md)
- Đọc/ghi file, `json` module, `csv` module, `pathlib`.
- **Bài tập:** Ghi/đọc danh sách user ra JSON; đọc CSV thành list of dict.

### [Bài 13: Type Hints & Static Checking](./13_type_hints.md)
- `typing` module, `Optional`, `Union`, generic (`TypeVar`), `mypy`.
- **Bài tập:** Thêm type hint đầy đủ cho code Bài 7, chạy `mypy` kiểm tra.

### [Bài 14: Testing với pytest](./14_testing.md)
- `pytest`, fixture, parametrize, mock (`unittest.mock`), coverage.
- **Bài tập:** Test table-driven bằng `@pytest.mark.parametrize`; mock 1 dependency.

### [Bài 15: Concurrency (threading, multiprocessing, asyncio)](./15_concurrency.md)
- GIL là gì, khi nào dùng threading/multiprocessing/asyncio, `async def`/`await`.
- **Bài tập:** Gọi song song nhiều "API" giả lập bằng `asyncio.gather`.

---

## Giai đoạn 3 — Python Nâng Cao (Advanced)

### [Bài 16: Kiến trúc project & Clean Architecture](./16_clean_architecture.md)
- Tách layer (router → service → repository), dependency injection qua constructor, Protocol (interface kiểu Python).

### [Bài 17: Làm việc với Database](./17_database.md)
- `SQLAlchemy` (Core + ORM), `Alembic` migration, kết nối PostgreSQL.
- **Bài tập:** CRUD bảng `users` bằng SQLAlchemy ORM.

### [Bài 18: REST API với FastAPI](./18_rest_api.md)
- Routing, Pydantic model (request/response validation), dependency injection của FastAPI, middleware.
- **Bài tập:** API `/todos` CRUD đầy đủ bằng FastAPI.

### [Bài 19: Authentication & Authorization](./19_auth.md)
- Hash password (`passlib`/`bcrypt`), JWT (`python-jose`), OAuth2PasswordBearer, RBAC.

### [Bài 20: Logging, Config, Deployment](./20_logging_config_deploy.md)
- `logging` module, Pydantic Settings (`.env`), Dockerize FastAPI, Uvicorn/Gunicorn, graceful shutdown.

---

## Giai đoạn 4 — Dự Án: REST API có Phân Quyền (Capstone)

> **File chi tiết đầy đủ: [21_capstone_project.md](./21_capstone_project.md)**

**Đề bài:** Xây dựng API quản lý "Task" bằng FastAPI với 2 role `admin`/`user`, cùng ngữ cảnh với dự án Go để dễ so sánh 2 ngôn ngữ.

### Cấu trúc thư mục đề xuất
```
taskapi/
├── app/
│   ├── main.py
│   ├── core/            # config, security (hash, jwt)
│   ├── domain/           # models + interface (Protocol) repository
│   ├── api/              # routers (FastAPI APIRouter) — tương đương handler
│   ├── services/         # business logic
│   ├── repositories/     # SQLAlchemy implementation
│   └── middleware/       # auth dependency, RBAC dependency
├── tests/
├── alembic/               # migration
├── .env.example
├── requirements.txt
├── Dockerfile
└── docker-compose.yml
```

### Các bước triển khai
1. Bootstrap project (venv, FastAPI, SQLAlchemy, Alembic).
2. Domain models (`User`, `Task`) + Protocol repository.
3. Database layer + migration bảng `users`, `tasks`.
4. Repository layer (SQLAlchemy implementation).
5. Service layer (đăng ký/đăng nhập, CRUD task với check ownership).
6. Dependency xác thực (`OAuth2PasswordBearer` + verify JWT).
7. Dependency phân quyền (`require_role("admin")`).
8. Router/API: `/auth/register`, `/auth/login`, `/tasks`, `/admin/users`.
9. Validation bằng Pydantic, response format chuẩn hóa.
10. Testing (pytest + `TestClient`, mock repository).
11. Logging & graceful shutdown (Uvicorn lifespan).
12. Dockerize + docker-compose kèm Postgres.
13. Test thủ công bằng `curl`/Swagger UI (`/docs` có sẵn nhờ FastAPI).

### Tiêu chí hoàn thành
- [ ] Đăng ký/đăng nhập trả JWT hợp lệ, password hash bằng bcrypt.
- [ ] `user` không sửa/xóa được task người khác, `admin` thao tác được mọi task.
- [ ] `pytest` chạy pass, có test cho logic ownership.
- [ ] `docker-compose up` chạy được, `/docs` (Swagger) truy cập được.

## Gợi ý cách học
- Vì bạn đã học Go trước, hãy chủ động so sánh: `interface` (Go) ↔ `Protocol`/ABC (Python), `goroutine` ↔ `asyncio`/`threading`, `struct tag` ↔ Pydantic field, `go mod` ↔ `requirements.txt`/`poetry`. Việc đối chiếu giúp nhớ nhanh hơn học từ đầu.
- Giai đoạn 4 nên tái sử dụng đúng schema DB (`users`, `tasks`) như bên Go để so sánh trực tiếp 2 stack.

## Tài liệu tham khảo
- Python docs: https://docs.python.org/3/
- Real Python: https://realpython.com/
- FastAPI docs: https://fastapi.tiangolo.com/
- SQLAlchemy docs: https://docs.sqlalchemy.org/
