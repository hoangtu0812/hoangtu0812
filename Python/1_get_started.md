# Bài 1: Giới thiệu & Cài đặt

## Mục tiêu
- Hiểu Python là gì, dùng cho việc gì, so sánh nhanh với Go (bạn đã học ở [Go/ROADMAP.md](../Go/ROADMAP.md)).
- Cài Python, VS Code + Pylance, tạo virtual environment, dùng `pip`.

## 1. Python là gì?

Python là ngôn ngữ **thông dịch (interpreted)**, **typing động (dynamically typed)**, cú pháp gần với ngôn ngữ tự nhiên, hệ sinh thái cực mạnh cho: scripting/tự động hóa, data science/AI/ML, backend web (Django/FastAPI/Flask).

### So sánh nhanh với Go

| | Go | Python |
|---|---|---|
| Typing | Tĩnh, kiểm tra lúc compile | Động, kiểm tra lúc chạy (có thể thêm type hint — [Bài 13](./13_type_hints.md)) |
| Thực thi | Biên dịch ra binary native | Thông dịch bởi CPython (chậm hơn nhưng linh hoạt hơn) |
| Concurrency | goroutine/channel dựng sẵn | GIL giới hạn thread song song thật sự — dùng `asyncio`/`multiprocessing` ([Bài 15](./15_concurrency.md)) |
| Error handling | trả `error` làm giá trị | `try/except` (exception) |
| Tốc độ dev | chậm hơn 1 chút, ít "magic" | rất nhanh để viết, nhiều thư viện sẵn |

## 2. Cài đặt

1. Tải Python tại https://www.python.org/downloads/ (khuyến khích bản mới nhất 3.12+).
2. Kiểm tra cài đặt:

```powershell
python --version
pip --version
```

3. Cài extension **Python** và **Pylance** cho VS Code (từ Microsoft).

## 3. Virtual environment — bắt buộc phải biết

Khác Go (mỗi project tự quản lý dependency qua `go.mod`), Python mặc định cài package **toàn hệ thống** — dễ xung đột version giữa các project. Virtual environment (venv) tạo 1 "hộp cách ly" riêng cho từng project.

```powershell
# Tạo venv trong thư mục project
python -m venv venv

# Kích hoạt (Windows PowerShell)
.\venv\Scripts\Activate.ps1

# Kích hoạt (Windows cmd)
.\venv\Scripts\activate.bat

# Sau khi kích hoạt, prompt sẽ hiện (venv) ở đầu dòng
# Cài package — chỉ ảnh hưởng venv này, không ảnh hưởng hệ thống
pip install requests

# Lưu danh sách dependency (tương đương go.mod)
pip freeze > requirements.txt

# Cài lại từ requirements.txt ở máy khác
pip install -r requirements.txt

# Thoát venv
deactivate
```

**Quy tắc:** LUÔN tạo venv riêng cho mỗi project Python, không bao giờ `pip install` trực tiếp vào Python hệ thống.

## 4. Hello World

```python
# hello.py
print("Hello, World!")
```

Chạy: `python hello.py`

Không cần `package main`/`func main()` như Go — Python thực thi tuần tự từ trên xuống, script nào cũng chạy được trực tiếp.

## 5. Chương trình có `if __name__ == "__main__":`

```python
def greet(name: str) -> str:
	return f"Xin chào, {name}!"

if __name__ == "__main__":
	print(greet("Ben"))
```

`if __name__ == "__main__":` đảm bảo code trong khối này chỉ chạy khi file được chạy trực tiếp (`python hello.py`), KHÔNG chạy khi file bị import từ file khác — idiom bắt buộc phải biết, tương đương "entry point" trong Go nhưng linh hoạt hơn vì 1 file Python vừa có thể là script vừa có thể là module để import.

## Bài tập

1. **Hello World**: viết `hello.py` in "Hello, World!", chạy bằng `python hello.py`.
2. **Venv**: tạo venv trong 1 thư mục mới, kích hoạt, cài package `requests` bằng `pip install requests`, chạy `pip freeze` để xem nó xuất hiện trong danh sách.
3. **`if __name__ == "__main__"`**: viết 2 file `greet.py` (chứa hàm `greet`) và `main.py` (import `greet` từ `greet.py` và gọi nó) — thêm `print("Đây là greet.py")` ở top-level của `greet.py` (ngoài hàm) để quan sát nó chạy khi nào.

## Tiếp theo
→ [Bài 2: Biến, kiểu dữ liệu, toán tử](./2_variables_types.md)
