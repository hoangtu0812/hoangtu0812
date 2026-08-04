# Bài 12: File I/O & JSON/CSV

## Mục tiêu
- Đọc/ghi file bằng `open()`/`pathlib`.
- `json` module — tương đương `encoding/json` của Go.
- `csv` module.

## 1. Đọc/ghi file cơ bản

```python
# Ghi file
with open("output.txt", "w", encoding="utf-8") as f:
	f.write("Xin chào Python!\n")
	f.write("Dòng thứ hai.\n")

# Đọc toàn bộ file
with open("output.txt", "r", encoding="utf-8") as f:
	content = f.read()
	print(content)

# Đọc theo từng dòng — tiết kiệm bộ nhớ với file lớn (liên hệ generator ở Bài 10)
with open("output.txt", "r", encoding="utf-8") as f:
	for line in f:
		print(line.strip())

# Đọc tất cả dòng thành list
with open("output.txt", "r", encoding="utf-8") as f:
	lines = f.readlines()

# Ghi thêm vào cuối file (append) thay vì ghi đè
with open("output.txt", "a", encoding="utf-8") as f:
	f.write("Dòng được thêm vào.\n")
```

**Luôn chỉ định `encoding="utf-8"` tường minh** — mặc định phụ thuộc hệ điều hành, dễ gây lỗi khi làm việc với tiếng Việt có dấu.

## 2. `pathlib` — cách hiện đại hơn thao tác đường dẫn (khuyến khích thay `os.path`)

```python
from pathlib import Path

p = Path("data") / "users.json"  # nối đường dẫn bằng toán tử / — hoạt động cross-platform
print(p.exists())
print(p.parent)     # thư mục cha
print(p.suffix)      # ".json"

p.parent.mkdir(parents=True, exist_ok=True)  # tạo thư mục nếu chưa có

# pathlib cũng hỗ trợ đọc/ghi trực tiếp, không cần open()
p.write_text("hello", encoding="utf-8")
content = p.read_text(encoding="utf-8")
```

## 3. `json` module

```python
import json
from dataclasses import dataclass, asdict

@dataclass
class User:
	id: int
	name: str
	email: str

user = User(id=1, name="Ben", email="ben@example.com")

# dict -> JSON string
json_str = json.dumps(asdict(user), ensure_ascii=False, indent=2)
print(json_str)
# {
#   "id": 1,
#   "name": "Ben",
#   "email": "ben@example.com"
# }

# JSON string -> dict
data = json.loads(json_str)
print(data["name"])  # Ben

# Ghi trực tiếp ra file
with open("user.json", "w", encoding="utf-8") as f:
	json.dump(asdict(user), f, ensure_ascii=False, indent=2)

# Đọc trực tiếp từ file
with open("user.json", "r", encoding="utf-8") as f:
	loaded = json.load(f)
```

`ensure_ascii=False` giữ nguyên ký tự Unicode (tiếng Việt) thay vì escape thành `\uXXXX` — nên luôn bật khi làm việc với dữ liệu tiếng Việt.

### Convert dict → object (ngược lại với `dataclasses.asdict`)

```python
data = {"id": 2, "name": "Alice", "email": "alice@example.com"}
user = User(**data)  # unpack dict thành keyword argument — khớp với field của dataclass
print(user)
```

Ở [Bài 18](./18_rest_api.md), FastAPI + Pydantic sẽ tự động hóa việc này (tương đương struct tag `json:"..."` của Go).

## 4. `csv` module

```python
import csv

# Ghi CSV
with open("users.csv", "w", newline="", encoding="utf-8") as f:
	writer = csv.writer(f)
	writer.writerow(["id", "name", "email"])
	writer.writerow([1, "Ben", "ben@example.com"])
	writer.writerow([2, "Alice", "alice@example.com"])

# Đọc CSV — trả về list các dòng (mỗi dòng là list các field)
with open("users.csv", "r", encoding="utf-8") as f:
	reader = csv.reader(f)
	header = next(reader)  # đọc dòng tiêu đề riêng
	for row in reader:
		print(row)

# DictReader — đọc mỗi dòng thành dict (key = tên cột) — TIỆN HƠN NHIỀU trong thực tế
with open("users.csv", "r", encoding="utf-8") as f:
	reader = csv.DictReader(f)
	for row in reader:
		print(row["name"], row["email"])
```

## Ví dụ đầy đủ

```python
import json
from dataclasses import dataclass, asdict
from pathlib import Path

@dataclass
class User:
	id: int
	name: str
	email: str

def save_users(users: list[User], path: str) -> None:
	data = [asdict(u) for u in users]
	Path(path).write_text(json.dumps(data, ensure_ascii=False, indent=2), encoding="utf-8")

def load_users(path: str) -> list[User]:
	data = json.loads(Path(path).read_text(encoding="utf-8"))
	return [User(**item) for item in data]

if __name__ == "__main__":
	users = [
		User(1, "Ben", "ben@example.com"),
		User(2, "Alice", "alice@example.com"),
	]
	save_users(users, "users.json")
	print("Đã lưu users.json")

	loaded = load_users("users.json")
	print(f"Đã đọc lại {len(loaded)} user")
	for u in loaded:
		print(u)
```

## Bài tập

1. **Ghi/đọc JSON**: dùng code mẫu trên làm nền, tự viết lại `save_users`/`load_users`.
2. **Đọc file theo dòng**: tạo file `.txt` nhiều dòng, đếm số dòng và số từ tổng cộng.
3. **CSV → list of dict**: viết hàm đọc file CSV bằng `csv.DictReader`, trả về `list[dict]`.
4. **Nâng cao**: viết hàm convert 1 file CSV `users.csv` sang `users.json` (đọc bằng `DictReader`, ghi bằng `json.dump`) — mô phỏng 1 tác vụ ETL nhỏ.

## Tiếp theo
→ [Bài 13: Type Hints & Static Checking](./13_type_hints.md)
