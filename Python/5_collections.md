# Bài 5: List, Tuple, Set, Dict

## Mục tiêu
- Phân biệt 4 cấu trúc dữ liệu built-in quan trọng nhất của Python.
- Slicing, mutable vs immutable.
- Comprehension cho list/dict/set.

## 1. List — mutable, có thứ tự (tương đương slice của Go)

```python
fruits = ["táo", "cam", "chuối"]

fruits.append("xoài")       # thêm cuối — giống append() của Go
fruits.insert(0, "nho")      # chèn tại vị trí
fruits.remove("cam")         # xóa theo GIÁ TRỊ (khác Go, không có sẵn built-in)
last = fruits.pop()          # xóa và trả về phần tử cuối

print(fruits, last)
print(len(fruits))
```

### Slicing — cú pháp mạnh hơn Go rất nhiều

```python
nums = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

print(nums[2:5])    # [2, 3, 4] — giống Go: from inclusive, to exclusive
print(nums[:3])      # [0, 1, 2] — bỏ trống start = từ đầu
print(nums[7:])      # [7, 8, 9] — bỏ trống stop = tới cuối
print(nums[-3:])     # [7, 8, 9] — index âm = đếm từ cuối (Go KHÔNG hỗ trợ)
print(nums[::2])     # [0, 2, 4, 6, 8] — step = 2
print(nums[::-1])    # đảo ngược toàn bộ list — idiom cực hay dùng
```

## 2. Tuple — immutable, có thứ tự

```python
point = (3, 4)
x, y = point   # unpacking — rất phổ biến

# Tuple KHÔNG sửa được sau khi tạo
# point[0] = 5  # TypeError: 'tuple' object does not support item assignment
```

Dùng tuple khi dữ liệu **cố định, không đổi** (tọa độ, RGB color, key của dict cần nhiều giá trị) — tương tự triết lý dùng `const`/struct nhỏ bất biến trong Go.

## 3. Set — không trùng lặp, không thứ tự (tương đương `map[T]bool` của Go)

```python
unique_nums = {1, 2, 3, 2, 1}
print(unique_nums)  # {1, 2, 3}

a = {1, 2, 3}
b = {2, 3, 4}
print(a | b)   # hợp: {1, 2, 3, 4}
print(a & b)   # giao: {2, 3}
print(a - b)   # hiệu: {1}

# Kiểm tra tồn tại — O(1), rất nhanh, giống lookup map của Go
print(3 in a)  # True
```

## 4. Dict — key-value (tương đương map của Go)

```python
scores = {"Alice": 90, "Bob": 85}

scores["Carol"] = 95      # thêm/sửa
print(scores.get("Dave"))          # None — KHÔNG lỗi nếu key không tồn tại (khác truy cập bằng [])
print(scores.get("Dave", 0))       # 0 — giá trị mặc định nếu không tồn tại

if "Alice" in scores:              # kiểm tra key tồn tại — tương đương "comma ok" của Go
	print(scores["Alice"])

del scores["Bob"]                   # xóa key

for name, score in scores.items():  # duyệt cả key và value
	print(name, score)
```

**Từ Python 3.7+**, dict giữ **đúng thứ tự chèn** (insertion order) — khác Go, nơi thứ tự duyệt map luôn ngẫu nhiên.

## 5. Comprehension cho cả 3 loại

```python
# List comprehension (đã học ở Bài 3)
squares = [x**2 for x in range(10)]

# Set comprehension
unique_lengths = {len(word) for word in ["a", "bb", "ccc", "dd"]}

# Dict comprehension
word_lengths = {word: len(word) for word in ["python", "go", "sap"]}
print(word_lengths)  # {'python': 6, 'go': 2, 'sap': 3}
```

## Ví dụ đầy đủ: Đếm tần suất từ

```python
def word_frequency(sentence: str) -> dict[str, int]:
	words = sentence.lower().split()
	freq: dict[str, int] = {}
	for word in words:
		freq[word] = freq.get(word, 0) + 1
	return freq

if __name__ == "__main__":
	sentence = "the quick brown fox jumps over the lazy dog the fox runs"
	freq = word_frequency(sentence)
	for word, count in sorted(freq.items()):
		print(f"{word}: {count}")
```

## Bài tập

1. **Loại bỏ trùng lặp**: viết hàm `remove_duplicates(items: list) -> list` giữ nguyên thứ tự xuất hiện đầu tiên (gợi ý: dùng `dict.fromkeys()` hoặc set + kiểm tra thứ tự).
2. **Đếm tần suất từ**: dùng code mẫu trên làm nền, tự viết lại.
3. **Đảo ngược list**: thử cả 2 cách — dùng slicing `[::-1]` và dùng vòng lặp thủ công (2 con trỏ, giống Bài 5 của Go).
4. **Nâng cao**: viết hàm nhận 2 list, trả về phần tử chung (dùng `set` intersection), phần tử chỉ có ở list 1 (set difference).

## Tiếp theo
→ [Bài 6: String & Comprehension nâng cao](./6_strings_comprehensions.md)
