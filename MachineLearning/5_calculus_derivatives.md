# Bài 5: Đạo hàm & Gradient

## Mục tiêu
- Đạo hàm 1 biến, các quy tắc đạo hàm cơ bản.
- Đạo hàm riêng, gradient của hàm nhiều biến.
- Hiểu vì sao gradient là "kim chỉ nam" để model học.

## 1. Đạo hàm là gì? — trực giác

Đạo hàm $f'(x)$ đo **tốc độ thay đổi tức thời** của hàm $f$ tại điểm $x$ — về mặt hình học, là **độ dốc (slope)** của tiếp tuyến tại điểm đó:

$$f'(x) = \lim_{h \to 0} \frac{f(x+h) - f(x)}{h}$$

```python
import numpy as np

def f(x):
	return x**2

def numerical_derivative(f, x, h=1e-5):
	"""Xấp xỉ đạo hàm bằng phương pháp sai phân hữu hạn (finite difference)."""
	return (f(x + h) - f(x - h)) / (2 * h)

print(numerical_derivative(f, 3))  # xấp xỉ 6 (đạo hàm thật của x^2 tại x=3 là 2x=6)
```

**Trong ML:** đạo hàm cho biết nếu tăng 1 tham số của model lên 1 chút, hàm loss thay đổi thế nào — đây là thông tin cốt lõi để biết **nên chỉnh tham số theo hướng nào**.

## 2. Quy tắc đạo hàm cơ bản

| Hàm | Đạo hàm |
|---|---|
| $f(x) = c$ (hằng số) | $f'(x) = 0$ |
| $f(x) = x^n$ | $f'(x) = nx^{n-1}$ |
| $f(x) = e^x$ | $f'(x) = e^x$ |
| $f(x) = \ln x$ | $f'(x) = \frac{1}{x}$ |
| $f(x) = \sin x$ | $f'(x) = \cos x$ |

**Quy tắc cộng:** $(f+g)' = f' + g'$
**Quy tắc tích (Product rule):** $(fg)' = f'g + fg'$
**Quy tắc thương (Quotient rule):** $\left(\frac{f}{g}\right)' = \frac{f'g - fg'}{g^2}$

## 3. Chain Rule — quy tắc quan trọng nhất cho ML

$$\text{Nếu } y = f(g(x)) \text{ thì } \frac{dy}{dx} = f'(g(x)) \cdot g'(x)$$

Cách viết dùng ký hiệu Leibniz (dễ nhớ hơn khi tính chuỗi dài):

$$\frac{dy}{dx} = \frac{dy}{du} \cdot \frac{du}{dx} \quad \text{với } u = g(x)$$

**Ví dụ:** $y = (3x+1)^2$. Đặt $u = 3x+1$, $y = u^2$.
$$\frac{dy}{du} = 2u, \quad \frac{du}{dx} = 3 \quad\Rightarrow\quad \frac{dy}{dx} = 2u \cdot 3 = 6(3x+1)$$

```python
def y(x):
	return (3*x + 1)**2

print(numerical_derivative(y, 2))  # xấp xỉ 6*(3*2+1) = 42
```

**Ứng dụng ML then chốt:** Chain Rule là công thức toán học đứng sau thuật toán **Backpropagation** ([Bài 21](./21_backpropagation.md)) — mạng neural là hàm hợp của rất nhiều lớp, và Chain Rule cho phép tính gradient qua từng lớp một cách có hệ thống.

## 4. Đạo hàm riêng (Partial Derivative) — cho hàm nhiều biến

Với hàm $f(x, y)$, đạo hàm riêng theo $x$ (ký hiệu $\frac{\partial f}{\partial x}$) là đạo hàm khi **coi $y$ là hằng số**:

$$f(x, y) = x^2y + 3y^2 \quad\Rightarrow\quad \frac{\partial f}{\partial x} = 2xy, \quad \frac{\partial f}{\partial y} = x^2 + 6y$$

```python
def f(x, y):
	return x**2 * y + 3 * y**2

def partial_derivative_x(f, x, y, h=1e-5):
	return (f(x + h, y) - f(x - h, y)) / (2 * h)

def partial_derivative_y(f, x, y, h=1e-5):
	return (f(x, y + h) - f(x, y - h)) / (2 * h)

print(partial_derivative_x(f, 2, 3))  # xấp xỉ 2*2*3 = 12
print(partial_derivative_y(f, 2, 3))  # xấp xỉ 4 + 18 = 22
```

## 5. Gradient — vector các đạo hàm riêng

Với hàm $f(x_1, ..., x_n)$, gradient là vector:

$$\nabla f = \begin{bmatrix} \frac{\partial f}{\partial x_1} \\ \vdots \\ \frac{\partial f}{\partial x_n} \end{bmatrix}$$

**Ý nghĩa hình học quan trọng nhất:** $\nabla f$ chỉ **hướng tăng nhanh nhất** của hàm $f$ tại điểm đang xét. Ngược lại, $-\nabla f$ chỉ hướng **giảm nhanh nhất**.

```python
def gradient(f, point, h=1e-5):
	"""Tính gradient bằng numerical differentiation cho hàm nhiều biến."""
	grad = np.zeros_like(point, dtype=float)
	for i in range(len(point)):
		point_plus = point.copy(); point_plus[i] += h
		point_minus = point.copy(); point_minus[i] -= h
		grad[i] = (f(point_plus) - f(point_minus)) / (2 * h)
	return grad

def f(p):
	x, y = p
	return x**2 + y**2  # hàm "bát úp" — cực tiểu tại (0,0)

point = np.array([3.0, 4.0])
print(gradient(f, point))  # xấp xỉ [6, 8] = [2x, 2y] tại (3,4)
```

**Đây chính là nền tảng của Gradient Descent** ([Bài 11](./11_discrete_math_optimization.md)): để giảm hàm loss, ta di chuyển tham số theo hướng $-\nabla f$ (ngược chiều gradient).

## Ví dụ đầy đủ: minh họa "học" bằng gradient

```python
import numpy as np

def loss(w):
	"""Hàm loss giả lập — cực tiểu tại w=5."""
	return (w - 5)**2

def gradient_1d(f, w, h=1e-5):
	return (f(w + h) - f(w - h)) / (2 * h)

def gradient_descent_demo(start_w, learning_rate=0.1, steps=20):
	w = start_w
	history = [w]
	for _ in range(steps):
		grad = gradient_1d(loss, w)
		w = w - learning_rate * grad  # di chuyển NGƯỢC chiều gradient
		history.append(w)
	return history

history = gradient_descent_demo(start_w=0.0)
print("w hội tụ về:", history[-1])  # tiến dần về 5
for i, w in enumerate(history[:5]):
	print(f"Bước {i}: w = {w:.4f}, loss = {loss(w):.4f}")
```

## Bài tập

1. **Đạo hàm bằng tay**: tính đạo hàm bằng tay của $f(x) = 3x^2 - 5x + 2$, $g(x) = e^{2x}$, $h(x) = \ln(x^2+1)$ (dùng chain rule cho $h$); verify bằng `numerical_derivative`.
2. **Đạo hàm riêng**: cho $f(x,y) = x^2y^3 + \sin(x)$, tính $\frac{\partial f}{\partial x}$ và $\frac{\partial f}{\partial y}$ bằng tay, verify bằng code.
3. **Gradient & hướng giảm nhanh nhất**: với $f(x,y) = (x-2)^2 + (y-3)^2$, tính gradient tại điểm $(0,0)$, giải thích tại sao di chuyển theo $-\nabla f$ đưa điểm gần lại cực tiểu $(2,3)$.
4. **Nâng cao**: mở rộng `gradient_descent_demo` cho hàm 2 biến $f(x,y) = (x-2)^2 + (y-3)^2$, chạy nhiều bước, in ra quỹ đạo hội tụ về $(2,3)$.

## Tiếp theo
→ [Bài 6: Đạo hàm ma trận & Chain Rule (nền tảng Backpropagation)](./6_calculus_matrix_calculus.md)
