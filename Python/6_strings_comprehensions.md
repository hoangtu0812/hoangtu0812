# Bài 6: String & Comprehension nâng cao

## Mục tiêu
- Thành thạo các method xử lý chuỗi hay dùng nhất.
- Regex cơ bản với `re`.
- Nested comprehension, generator expression.

## 1. String method thường dùng

```python
s = "  Hello, World!  "

print(s.strip())            # "Hello, World!" — bỏ khoảng trắng 2 đầu
print(s.lower())            # "  hello, world!  "
print(s.upper())            # "  HELLO, WORLD!  "
print(s.strip().split(", ")) # ['Hello', 'World!']
print("-".join(["a", "b", "c"]))  # "a-b-c"
print(s.strip().replace("World", "Python"))  # "Hello, Python!"
print(s.strip().startswith("Hello"))  # True
print("42".isdigit())        # True
print(s.strip().find("World"))  # 7 — trả -1 nếu không tìm thấy
```

## 2. f-string nâng cao (nối tiếp [Bài 2](./2_variables_types.md))

```python
value = 1234567.891
print(f"{value:,.2f}")     # "1,234,567.89" — dấu phẩy phân cách hàng nghìn
print(f"{value:>15,.2f}")  # căn phải, độ rộng 15
print(f"{'text':^20}")     # căn giữa trong độ rộng 20

name = "ben"
print(f"{name!r}")  # 'ben' — !r gọi repr() thay vì str()
```

## 3. Regex cơ bản với `re`

```python
import re

text = "Liên hệ: ben@example.com hoặc alice@test.org"

emails = re.findall(r"[\w.+-]+@[\w-]+\.[\w.-]+", text)
print(emails)  # ['ben@example.com', 'alice@test.org']

match = re.search(r"(\d{3})-(\d{3})-(\d{4})", "Gọi 090-123-4567")
if match:
	print(match.group())    # 090-123-4567
	print(match.groups())    # ('090', '123', '4567')

cleaned = re.sub(r"\s+", " ", "quá   nhiều    khoảng   trắng")
print(cleaned)  # "quá nhiều khoảng trắng"
```

## 4. Nested comprehension

```python
# Ma trận 3x3
matrix = [[i * 3 + j for j in range(3)] for i in range(3)]
print(matrix)  # [[0, 1, 2], [3, 4, 5], [6, 7, 8]]

# "Flatten" ma trận thành 1 list phẳng
flat = [num for row in matrix for num in row]
print(flat)  # [0, 1, 2, 3, 4, 5, 6, 7, 8]
```

Quy tắc đọc nested comprehension: đọc **từ trái sang phải giống viết vòng lặp lồng nhau bình thường** — `for row in matrix` là vòng ngoài, `for num in row` là vòng trong.

## 5. Generator expression — lazy, tiết kiệm bộ nhớ

```python
# List comprehension: tạo TOÀN BỘ list trong bộ nhớ ngay lập tức
squares_list = [x**2 for x in range(1_000_000)]

# Generator expression (chỉ khác cú pháp: () thay vì []): sinh giá trị LAZY, từng cái một khi cần
squares_gen = (x**2 for x in range(1_000_000))

print(sum(squares_gen))  # vẫn hoạt động, nhưng KHÔNG tốn bộ nhớ lưu 1 triệu số cùng lúc
```

Dùng generator expression khi chỉ cần **duyệt qua 1 lần** (vd truyền vào `sum()`, `max()`, `any()`), dùng list comprehension khi cần **truy cập nhiều lần** hoặc cần index/slicing. Xem thêm generator đầy đủ ở [Bài 10](./10_iterators_generators.md).

## Ví dụ đầy đủ

```python
import re

def count_words(sentence: str) -> int:
	return len(sentence.split())

def is_palindrome(s: str) -> bool:
	cleaned = re.sub(r"[^a-z0-9]", "", s.lower())
	return cleaned == cleaned[::-1]

def parse_csv_line(line: str) -> list[str]:
	return [field.strip() for field in line.split(",")]

if __name__ == "__main__":
	print(count_words("Học Python thật là thú vị"))
	print(is_palindrome("A man a plan a canal Panama"))  # True
	print(parse_csv_line("Ben, 25, Quang Ngai"))
```

## Bài tập

1. **Đếm từ trong câu**: viết `count_words(sentence: str) -> int`.
2. **Palindrome**: viết `is_palindrome(s: str) -> bool` bỏ qua khoảng trắng/hoa-thường/dấu câu (dùng `re.sub`).
3. **Parse CSV đơn giản**: viết hàm nhận 1 dòng CSV (chuỗi phân cách bởi dấu phẩy), trả về list các field đã `strip()`.
4. **Nâng cao**: viết nested list comprehension chuyển ma trận 2 chiều thành **ma trận chuyển vị (transpose)** — không dùng thư viện ngoài (gợi ý: `[[row[i] for row in matrix] for i in range(len(matrix[0]))]`).

## Tiếp theo
→ [Bài 7: OOP trong Python](./7_oop.md)
