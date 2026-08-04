# Bài 7: OOP trong Python

## Mục tiêu
- `class`, `__init__`, kế thừa, đa hình.
- Encapsulation theo convention (Python không có `private` thật sự).
- Dunder method, `@property`, `@dataclass`.

## 1. Class cơ bản

```python
class Rectangle:
	def __init__(self, width: float, height: float):
		self.width = width
		self.height = height

	def area(self) -> float:
		return self.width * self.height

	def perimeter(self) -> float:
		return 2 * (self.width + self.height)

rect = Rectangle(10, 5)
print(rect.area())       # 50
print(rect.perimeter())  # 30
```

`self` tương đương receiver của method trong Go — luôn là tham số đầu tiên, Python **không tự động ẩn nó** như Go (`self` phải viết tường minh trong mọi method).

## 2. Kế thừa

```python
class Animal:
	def __init__(self, name: str):
		self.name = name

	def describe(self) -> str:
		return f"{self.name} là một loài động vật"

class Dog(Animal):
	def __init__(self, name: str, breed: str):
		super().__init__(name)  # gọi __init__ của class cha — bắt buộc nếu override __init__
		self.breed = breed

	def describe(self) -> str:  # override — Python không cần từ khóa gì đặc biệt
		return f"{self.name} là chó giống {self.breed}"

class Cat(Animal):
	pass  # không override gì, dùng nguyên describe() của Animal

dog = Dog("Rex", "Poodle")
cat = Cat("Mimi")
print(dog.describe())  # Rex là chó giống Poodle
print(cat.describe())  # Mimi là một loài động vật
```

**So sánh với Go:** Go dùng **composition** (embedding struct — [Go Bài 6](../Go/6_structs_methods.md)), Python dùng **kế thừa thật sự** (`class Dog(Animal)`). Cả 2 đều đạt được tái sử dụng code, nhưng triết lý khác nhau — Python OOP gần Java/C# hơn.

## 3. Đa hình (Polymorphism)

```python
def print_description(animal: Animal) -> None:
	print(animal.describe())  # gọi đúng method của SUBCLASS thực tế, không phải Animal

animals = [Dog("Rex", "Poodle"), Cat("Mimi")]
for a in animals:
	print_description(a)  # tự động gọi đúng describe() của Dog hoặc Cat
```

## 4. Encapsulation — convention, không phải rào chắn cứng

Python **không có** `private`/`public`/`protected` thật sự như Java. Thay vào đó dùng convention đặt tên:

```python
class BankAccount:
	def __init__(self, balance: float):
		self._balance = balance     # 1 gạch dưới: "protected theo quy ước" — vẫn truy cập được, nhưng ngầm hiểu là không nên
		self.__pin = "1234"          # 2 gạch dưới: "name mangling" — Python đổi tên thành _BankAccount__pin để khó truy cập nhầm

	def deposit(self, amount: float) -> None:
		self._balance += amount

	def get_balance(self) -> float:
		return self._balance

account = BankAccount(100)
print(account._balance)      # vẫn truy cập được — Python tin tưởng lập trình viên tự giác
# print(account.__pin)        # AttributeError — do name mangling
print(account._BankAccount__pin)  # vẫn lách được nếu cố tình — nhưng KHÔNG NÊN làm vậy
```

## 5. `@property` — getter/setter kiểu Pythonic

```python
class Circle:
	def __init__(self, radius: float):
		self._radius = radius

	@property
	def radius(self) -> float:
		return self._radius

	@radius.setter
	def radius(self, value: float) -> None:
		if value <= 0:
			raise ValueError("Bán kính phải dương")
		self._radius = value

	@property
	def area(self) -> float:  # "computed property" — gọi như attribute, không cần ()
		import math
		return math.pi * self._radius ** 2

c = Circle(5)
print(c.radius)  # 5 — gọi như attribute, không phải c.radius()
print(c.area)     # 78.53...
c.radius = 10     # dùng setter, tự động validate
```

## 6. Dunder methods (magic methods)

```python
class Point:
	def __init__(self, x: int, y: int):
		self.x = x
		self.y = y

	def __str__(self) -> str:       # dùng khi print(point) — giống fmt.Stringer của Go
		return f"Point({self.x}, {self.y})"

	def __eq__(self, other) -> bool: # dùng khi so sánh point1 == point2
		return self.x == other.x and self.y == other.y

	def __add__(self, other):        # dùng khi point1 + point2
		return Point(self.x + other.x, self.y + other.y)

p1 = Point(1, 2)
p2 = Point(3, 4)
print(p1)             # Point(1, 2) — nhờ __str__
print(p1 == Point(1, 2))  # True — nhờ __eq__
print(p1 + p2)          # Point(4, 6) — nhờ __add__
```

## 7. `@dataclass` — giảm boilerplate cho class chỉ chứa dữ liệu

```python
from dataclasses import dataclass

@dataclass
class Point:
	x: int
	y: int

# Tự động sinh __init__, __repr__, __eq__ — KHÔNG cần viết tay
p1 = Point(1, 2)
p2 = Point(1, 2)
print(p1)          # Point(x=1, y=2) — __repr__ tự sinh
print(p1 == p2)     # True — __eq__ tự sinh
```

`@dataclass` gần giống struct trần của Go (chỉ có field, method thêm sau) — dùng khi class chủ yếu là "túi đựng dữ liệu".

## Bài tập

1. **`Shape` kế thừa**: viết `Shape` (base), `Rectangle`, `Circle` kế thừa, mỗi lớp implement `area()` riêng, viết hàm `print_area(shape: Shape)` dùng đa hình.
2. **`@dataclass` cho `Point`**: viết `Point` bằng `@dataclass`, so sánh với cách viết class thủ công (mục 6) — đếm số dòng code tiết kiệm được.
3. **`@property` validate**: viết `BankAccount` với `@property balance`, setter raise `ValueError` nếu gán số âm.
4. **Nâng cao**: viết class `Vector` với `__add__`, `__sub__`, `__mul__` (nhân vô hướng), `__str__`, `__eq__` — mô phỏng phép toán vector 2D.

## Tiếp theo
→ [Bài 8: Xử lý lỗi (Exception Handling)](./8_exceptions.md)
