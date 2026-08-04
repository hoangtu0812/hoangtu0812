# Bài 6: Đạo hàm ma trận & Chain Rule (nền tảng Backpropagation)

## Mục tiêu
- Đạo hàm theo vector/ma trận (matrix calculus).
- Jacobian.
- Áp dụng chain rule qua nhiều lớp — CHÍNH LÀ toán học của Backpropagation.

## 1. Vì sao cần đạo hàm ma trận?

Trong ML/Deep Learning, hàm loss $L$ phụ thuộc vào **hàng nghìn/triệu tham số** (trọng số của mạng neural), không chỉ 1-2 biến như [Bài 5](./5_calculus_derivatives.md). Đạo hàm ma trận cho phép biểu diễn **toàn bộ gradient** của $L$ theo mọi tham số bằng 1 công thức ma trận gọn gàng, thay vì viết hàng triệu đạo hàm riêng lẻ.

## 2. Gradient của hàm vô hướng theo vector

Với $f: \mathbb{R}^n \to \mathbb{R}$ (hàm nhận vector, trả về 1 số — đúng dạng của hàm loss), gradient là vector đạo hàm riêng ([Bài 5 mục 5](./5_calculus_derivatives.md)):

$$\nabla_{\vec{x}} f = \begin{bmatrix} \partial f/\partial x_1 \\ \vdots \\ \partial f/\partial x_n \end{bmatrix}$$

### Công thức hay dùng nhất trong ML: gradient của dạng tuyến tính và dạng toàn phương

$$f(\vec{x}) = \vec{a}^T\vec{x} \quad\Rightarrow\quad \nabla_{\vec{x}} f = \vec{a}$$

$$f(\vec{x}) = \vec{x}^TA\vec{x} \quad\Rightarrow\quad \nabla_{\vec{x}} f = (A + A^T)\vec{x} \quad (\text{= } 2A\vec{x} \text{ nếu } A \text{ đối xứng})$$

```python
import numpy as np

def numerical_gradient(f, x, h=1e-5):
	grad = np.zeros_like(x, dtype=float)
	for i in range(len(x)):
		x_plus = x.copy(); x_plus[i] += h
		x_minus = x.copy(); x_minus[i] -= h
		grad[i] = (f(x_plus) - f(x_minus)) / (2 * h)
	return grad

a = np.array([2.0, 3.0])
f = lambda x: a @ x
x = np.array([1.0, 1.0])
print(numerical_gradient(f, x))  # xấp xỉ [2, 3] = a — khớp công thức
```

## 3. Gradient của hàm loss MSE — ví dụ áp dụng trực tiếp trong Linear Regression

Hàm loss MSE (Mean Squared Error, chi tiết ở [Bài 14](./14_linear_regression.md)):

$$L(\vec{w}) = \frac{1}{m}\|X\vec{w} - \vec{y}\|^2 = \frac{1}{m}(X\vec{w}-\vec{y})^T(X\vec{w}-\vec{y})$$

Áp dụng đạo hàm ma trận:

$$\nabla_{\vec{w}} L = \frac{2}{m}X^T(X\vec{w} - \vec{y})$$

Đây chính xác là công thức gradient dùng trong Gradient Descent cho Linear Regression — xem cách nó được dùng thực tế ở [Bài 14](./14_linear_regression.md).

```python
def mse_loss(w, X, y):
	m = len(y)
	return (1/m) * np.sum((X @ w - y)**2)

def mse_gradient(w, X, y):
	m = len(y)
	return (2/m) * X.T @ (X @ w - y)

X = np.array([[1, 1], [1, 2], [1, 3]], dtype=float)  # cột đầu = 1 cho bias
y = np.array([2, 3, 5], dtype=float)
w = np.array([0.0, 0.0])

print(mse_gradient(w, X, y))
print(numerical_gradient(lambda w: mse_loss(w, X, y), w))  # phải xấp xỉ bằng nhau
```

## 4. Jacobian — tổng quát hóa gradient cho hàm trả về vector

Với $f: \mathbb{R}^n \to \mathbb{R}^m$ (hàm nhận vector, trả về **vector**), Jacobian là ma trận $m \times n$ chứa mọi đạo hàm riêng:

$$J = \begin{bmatrix} \partial f_1/\partial x_1 & \dots & \partial f_1/\partial x_n \\ \vdots & \ddots & \vdots \\ \partial f_m/\partial x_1 & \dots & \partial f_m/\partial x_n \end{bmatrix}$$

**Ứng dụng ML:** mỗi lớp (layer) của mạng neural là 1 hàm nhận vector đầu vào, trả về vector đầu ra — Jacobian mô tả cách output của lớp đó thay đổi theo input, chính là "viên gạch" mà Backpropagation ghép lại qua nhiều lớp.

## 5. Chain Rule cho nhiều lớp — chuẩn bị trực tiếp cho Backpropagation

Giả sử mạng neural đơn giản 2 lớp: $z = W_1\vec{x}$, $a = \sigma(z)$ (hàm kích hoạt), $\hat{y} = W_2 a$, $L = \text{loss}(\hat{y}, y)$.

Chain Rule cho ta:

$$\frac{\partial L}{\partial W_1} = \frac{\partial L}{\partial \hat{y}} \cdot \frac{\partial \hat{y}}{\partial a} \cdot \frac{\partial a}{\partial z} \cdot \frac{\partial z}{\partial W_1}$$

Đây **chính xác** là cách Backpropagation hoạt động: tính gradient của loss theo output trước, rồi "lan truyền ngược" qua từng lớp bằng cách nhân liên tiếp các đạo hàm riêng — mỗi bước chỉ cần biết đạo hàm cục bộ của lớp đó (local gradient). Chi tiết đầy đủ với code from scratch ở [Bài 21](./21_backpropagation.md).

```python
# Minh họa numerically cho mạng 2 lớp cực đơn giản (1 neuron mỗi lớp)
def sigmoid(z):
	return 1 / (1 + np.exp(-z))

def sigmoid_derivative(z):
	s = sigmoid(z)
	return s * (1 - s)

def forward(x, w1, w2):
	z = w1 * x
	a = sigmoid(z)
	y_hat = w2 * a
	return z, a, y_hat

def backward(x, y, w1, w2):
	z, a, y_hat = forward(x, w1, w2)
	loss = (y_hat - y)**2

	dL_dyhat = 2 * (y_hat - y)         # đạo hàm loss theo output
	dyhat_da = w2                        # đạo hàm output theo a
	da_dz = sigmoid_derivative(z)         # đạo hàm sigmoid
	dz_dw1 = x                             # đạo hàm z theo w1

	# Chain rule: nhân liên tiếp các đạo hàm cục bộ
	dL_dw1 = dL_dyhat * dyhat_da * da_dz * dz_dw1
	dL_dw2 = dL_dyhat * a

	return dL_dw1, dL_dw2, loss

grad_w1, grad_w2, loss = backward(x=1.0, y=1.0, w1=0.5, w2=0.8)
print(f"Gradient w1: {grad_w1:.4f}, Gradient w2: {grad_w2:.4f}, Loss: {loss:.4f}")
```

## Bài tập

1. **Gradient dạng tuyến tính & toàn phương**: verify bằng code công thức $\nabla_{\vec{x}} (\vec{a}^T\vec{x}) = \vec{a}$ và $\nabla_{\vec{x}} (\vec{x}^TA\vec{x}) = 2A\vec{x}$ (với $A$ đối xứng) bằng `numerical_gradient`.
2. **Gradient MSE**: dùng code mẫu mục 3, verify `mse_gradient` khớp với `numerical_gradient` trên cùng dataset, thử với dataset lớn hơn (5-10 mẫu, 2-3 đặc trưng).
3. **Chain rule qua 2 lớp**: dùng code mẫu mục 5, tính tay đạo hàm $\frac{\partial L}{\partial w_1}$ cho 1 bộ giá trị cụ thể ($x=2, y=0, w_1=0.3, w_2=0.5$), verify bằng code.
4. **Nâng cao**: mở rộng ví dụ mục 5 thành mạng 2 neuron ở lớp ẩn (thay vì 1), viết lại `forward`/`backward` — đây chính là bước đệm trực tiếp cho [Bài 21](./21_backpropagation.md).

## Tiếp theo
→ [Bài 7: Tối ưu hóa không ràng buộc & Taylor Series](./7_calculus_optimization.md)
