# Bài 9: Module, Package & Virtual Environment

## Mục tiêu
- `import` — module vs package.
- Cấu trúc package chuẩn với `__init__.py`.
- `pip`, `requirements.txt`, `venv`/`poetry` — tương đương `go.mod` của Go.

## 1. Module — mỗi file `.py` là 1 module

```python
# shapes.py
def area_rectangle(width: float, height: float) -> float:
	return width * height

PI = 3.14159
```

```python
# main.py
import shapes

print(shapes.area_rectangle(10, 5))
print(shapes.PI)

# Hoặc import cụ thể tên cần dùng
from shapes import area_rectangle, PI
print(area_rectangle(10, 5))

# Đặt alias (thường dùng với thư viện tên dài, vd pandas as pd)
import shapes as sh
print(sh.PI)
```

## 2. Package — thư mục chứa nhiều module + `__init__.py`

```
myproject/
├── main.py
└── shapes/
    ├── __init__.py       # đánh dấu shapes/ là 1 package, có thể export sẵn 1 số tên
    ├── rectangle.py
    └── circle.py
```

```python
# shapes/rectangle.py
class Rectangle:
	def __init__(self, width: float, height: float):
		self.width = width
		self.height = height

	def area(self) -> float:
		return self.width * self.height
```

```python
# shapes/circle.py
import math

class Circle:
	def __init__(self, radius: float):
		self.radius = radius

	def area(self) -> float:
		return math.pi * self.radius ** 2
```

```python
# shapes/__init__.py — cho phép import gọn hơn từ package
from .rectangle import Rectangle
from .circle import Circle
```

```python
# main.py
from shapes import Rectangle, Circle  # nhờ __init__.py, không cần shapes.rectangle.Rectangle

r = Rectangle(10, 5)
c = Circle(3)
print(r.area(), c.area())
```

Từ Python 3.3+, `__init__.py` không còn **bắt buộc** để 1 thư mục được coi là package (namespace package), nhưng vẫn nên có để kiểm soát rõ ràng những gì package export ra ngoài — tương đương triết lý "chỉ export cái cần thiết" ở [Go Bài 9](../Go/9_packages_modules.md).

## 3. Absolute import vs relative import

```python
# Absolute import — luôn từ gốc project, RÕ RÀNG, khuyến khích dùng
from shapes.rectangle import Rectangle

# Relative import — CHỈ dùng bên trong package, dựa vào vị trí file hiện tại
# (trong shapes/__init__.py)
from .rectangle import Rectangle   # . = package hiện tại
from ..utils import helper          # .. = package cha
```

## 4. `pip` và `requirements.txt`

```powershell
# Cài 1 package cụ thể version
pip install requests==2.31.0

# Cài từ file requirements.txt (tương đương go.sum + go mod download)
pip install -r requirements.txt

# Xuất danh sách package hiện có trong venv ra file
pip freeze > requirements.txt

# Gỡ 1 package
pip uninstall requests
```

`requirements.txt` mẫu:
```
requests==2.31.0
fastapi==0.110.0
sqlalchemy>=2.0,<3.0
```

## 5. Poetry — quản lý dependency hiện đại hơn (tùy chọn, khuyến khích cho project lớn)

```powershell
pip install poetry
poetry init          # tạo pyproject.toml, tương đương go.mod
poetry add requests  # thêm dependency, tự cập nhật pyproject.toml + poetry.lock (giống go.sum)
poetry install        # cài từ pyproject.toml + lock file
poetry shell           # kích hoạt venv do poetry tự quản lý
```

`poetry` giải quyết vấn đề `pip freeze` hay gặp: `requirements.txt` không phân biệt được "dependency trực tiếp bạn cài" với "dependency của dependency" — `pyproject.toml` + `poetry.lock` làm rõ điều này, tương tự `go.mod`/`go.sum`.

## 6. Cấu trúc project chuẩn (tham khảo cho các bài sau)

```
myproject/
├── venv/                  # KHÔNG commit vào git
├── src/ hoặc app/
│   ├── __init__.py
│   ├── main.py
│   └── shapes/
│       ├── __init__.py
│       ├── rectangle.py
│       └── circle.py
├── tests/
│   └── test_shapes.py
├── requirements.txt
├── .gitignore              # PHẢI có venv/, __pycache__/, .env
└── README.md
```

## Bài tập

1. **Tách package `shapes`**: lấy `Rectangle`/`Circle` (Bài 7) ra thành package riêng như ví dụ trên, import và dùng từ `main.py`.
2. **`requirements.txt`**: tạo venv mới, cài `requests` và `python-dotenv`, chạy `pip freeze > requirements.txt`, kiểm tra nội dung file.
3. **Absolute vs relative import**: thử cả 2 cách import trong package `shapes`, quan sát lỗi khi chạy relative import trực tiếp như 1 script độc lập (`python shapes/rectangle.py`) thay vì qua `main.py` — giải thích tại sao.
4. **`.gitignore`**: viết file `.gitignore` chuẩn cho 1 project Python (venv, `__pycache__`, `.env`, file build).

## Tổng kết Giai đoạn 1
Bạn đã nắm nền tảng Python: biến/kiểu, control flow, hàm, collections, string, OOP, exception, module/package. Giai đoạn 2 sẽ đi vào các khái niệm khiến Python mạnh cho backend: iterator/generator, decorator, testing, và concurrency.

## Tiếp theo
→ [Bài 10: Iterator & Generator](./10_iterators_generators.md)
