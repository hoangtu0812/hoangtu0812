# Bài 6: Đạo hàm ma trận & Chain Rule (nền tảng Backpropagation)

## Mục tiêu
- Hiểu vì sao cần "gói gọn" hàng triệu đạo hàm riêng thành công thức ma trận.
- Suy luận đầy đủ công thức gradient của MSE — không chỉ chép công thức.
- Hiểu chain rule nhiều lớp bằng trực giác trước khi thấy code.

---

## PHẦN A — Ý NGHĨA TOÁN HỌC

### 1. Vì sao cần đạo hàm ma trận?

Mạng neural có thể có hàng triệu tham số. Viết riêng $\frac{\partial L}{\partial w_1}, \frac{\partial L}{\partial w_2}, ...$ cho từng tham số là bất khả thi để trình bày và lập trình. Đạo hàm ma trận cho phép gói TOÀN BỘ gradient theo hàng triệu tham số vào **1 công thức ma trận duy nhất** — không phải kỹ thuật mới, chỉ là cách viết gọn của những gì đã học ở [Bài 5](./5_calculus_derivatives.md), áp dụng đồng loạt cho mọi thành phần.

### 2. Gradient của hàm vô hướng theo vector — nhắc lại và mở rộng từ Bài 5

Với $f:\mathbb{R}^n\to\mathbb{R}$ (nhận vector, trả về 1 số — đúng dạng hàm loss), $\nabla_{\vec{x}}f$ là vector đạo hàm riêng ([Bài 5 mục 6](./5_calculus_derivatives.md)). Ta suy ra 2 công thức nền tảng hay gặp nhất trong ML bằng cách áp dụng đạo hàm riêng cho TỪNG thành phần rồi gộp lại.

**Công thức 1 — gradient của dạng tuyến tính** $f(\vec{x})=\vec{a}^T\vec{x}=\sum_i a_ix_i$:

Đạo hàm riêng theo $x_j$: $\frac{\partial f}{\partial x_j} = a_j$ (vì trong tổng $\sum_i a_ix_i$, chỉ số hạng $a_jx_j$ chứa $x_j$, đạo hàm của nó theo $x_j$ là $a_j$, các số hạng khác là hằng số theo $x_j$ nên đạo hàm = 0). Gộp lại mọi $j$: $\nabla_{\vec{x}}f = \vec{a}$.

**Công thức 2 — gradient của dạng toàn phương** $f(\vec{x})=\vec{x}^TA\vec{x}=\sum_i\sum_j A_{ij}x_ix_j$ (với $A$ đối xứng):

Đạo hàm riêng theo $x_k$ cần xét mọi số hạng có chứa $x_k$ (cả khi $i=k$ và khi $j=k$), áp dụng quy tắc tích ([Bài 5 mục 3](./5_calculus_derivatives.md)) — kết quả rút gọn: $\frac{\partial f}{\partial x_k} = 2\sum_j A_{kj}x_j = 2(A\vec{x})_k$. Gộp lại: $\nabla_{\vec{x}}f = 2A\vec{x}$.

**Verify bằng ví dụ số nhỏ** ($n=2$): $A=\begin{bmatrix}a&b\\b&d\end{bmatrix}$ (đối xứng), $f(x_1,x_2)=ax_1^2+2bx_1x_2+dx_2^2$. Đạo hàm riêng: $\frac{\partial f}{\partial x_1}=2ax_1+2bx_2$, $\frac{\partial f}{\partial x_2}=2bx_1+2dx_2$ — viết dưới dạng vector đúng bằng $2A\vec{x}$.

### 3. Suy luận đầy đủ Gradient của MSE — TỪNG BƯỚC, không chỉ chép công thức

Hàm loss MSE của Linear Regression ([Bài 14](./14_linear_regression.md)):

$$L(\vec{w}) = \frac{1}{m}\|X\vec{w}-\vec{y}\|^2 = \frac{1}{m}(X\vec{w}-\vec{y})^T(X\vec{w}-\vec{y})$$

**Bước 1 — khai triển** (dùng tính chất $\|\vec{v}\|^2=\vec{v}\cdot\vec{v}$ từ [Bài 2 mục 4](./2_linear_algebra_vectors.md)):

$$L(\vec{w}) = \frac{1}{m}\left[\vec{w}^TX^TX\vec{w} - 2\vec{y}^TX\vec{w} + \vec{y}^T\vec{y}\right]$$

**Bước 2 — áp dụng 2 công thức vừa suy ra ở mục 2 cho từng số hạng:**
- $\vec{w}^TX^TX\vec{w}$ là dạng toàn phương với $A=X^TX$ (đối xứng) → gradient $=2X^TX\vec{w}$.
- $-2\vec{y}^TX\vec{w}$ là dạng tuyến tính theo $\vec{w}$ với $\vec{a}=-2X^T\vec{y}$ → gradient $=-2X^T\vec{y}$.
- $\vec{y}^T\vec{y}$ không chứa $\vec{w}$ → gradient $=\vec{0}$.

**Bước 3 — gộp lại:**

$$\nabla_{\vec{w}}L = \frac{1}{m}\left[2X^TX\vec{w}-2X^T\vec{y}\right] = \frac{2}{m}X^T(X\vec{w}-\vec{y})$$

Đây chính xác là công thức đã dùng trực tiếp ở [Bài 14](./14_linear_regression.md) — giờ bạn đã thấy nó đến từ đâu, từng bước, không phải chỉ "được cho sẵn".

### 4. Jacobian — khi hàm trả về CẢ MỘT VECTOR, không chỉ 1 số

Với $f:\mathbb{R}^n\to\mathbb{R}^m$ (nhận vector, trả về vector), mỗi thành phần đầu ra $f_i$ có gradient riêng theo $\vec{x}$ — gộp $m$ gradient đó (mỗi cái là 1 dòng) thành ma trận Jacobian $m\times n$:

$$J_{ij} = \frac{\partial f_i}{\partial x_j}$$

**Trực giác:** Jacobian là "phiên bản ma trận" của đạo hàm — trong khi $f'(x)$ (1 số) mô tả hàm 1 chiều thay đổi ra sao, $J$ (1 ma trận) mô tả CẢ vector đầu ra thay đổi ra sao khi CẢ vector đầu vào thay đổi, tại 1 điểm cụ thể. Mỗi lớp (layer) của mạng neural là 1 hàm nhận vector, trả vector — Jacobian của nó chính là "viên gạch" mà Backpropagation ghép nối qua nhiều lớp.

### 5. Chain Rule qua nhiều lớp — TỪNG BƯỚC suy luận, chuẩn bị cho Backpropagation

Xét mạng neural đơn giản nhất có thể: 1 input $x$, 1 trọng số lớp 1 ($w_1$), hàm kích hoạt $\sigma$, 1 trọng số lớp 2 ($w_2$), rồi tính loss:

$$z = w_1x \quad\to\quad a=\sigma(z) \quad\to\quad \hat{y}=w_2a \quad\to\quad L=(\hat{y}-y)^2$$

Đây là **hàm hợp 4 tầng** — muốn biết $\frac{\partial L}{\partial w_1}$ (loss thay đổi thế nào nếu chỉnh $w_1$, tham số Ở TẦNG ĐẦU TIÊN, xa nhất so với $L$), áp dụng chain rule liên tiếp ([Bài 5 mục 4](./5_calculus_derivatives.md)) qua TỪNG mắt xích trung gian:

$$\frac{\partial L}{\partial w_1} = \underbrace{\frac{\partial L}{\partial\hat{y}}}_{\text{bước 4}\to3} \cdot \underbrace{\frac{\partial\hat{y}}{\partial a}}_{3\to2} \cdot \underbrace{\frac{\partial a}{\partial z}}_{2\to1} \cdot \underbrace{\frac{\partial z}{\partial w_1}}_{1\to w_1}$$

**Ý nghĩa trực giác của từng mắt xích** (đúng như phép "tỷ giá nối tiếp" ở [Bài 5 mục 4](./5_calculus_derivatives.md)): "$w_1$ thay đổi 1 chút làm $z$ thay đổi theo tỷ lệ $\partial z/\partial w_1$; $z$ thay đổi làm $a$ thay đổi theo tỷ lệ $\partial a/\partial z$; $a$ thay đổi làm $\hat{y}$ thay đổi theo tỷ lệ $\partial\hat{y}/\partial a$; $\hat{y}$ thay đổi làm $L$ thay đổi theo tỷ lệ $\partial L/\partial\hat{y}$" — nhân tất cả các tỷ lệ liên tiếp lại, ta được ảnh hưởng TỔNG HỢP của $w_1$ lên $L$, dù chúng cách nhau qua 3 tầng trung gian.

**Đây chính xác là cách Backpropagation hoạt động**: tính $\frac{\partial L}{\partial\hat{y}}$ trước (dễ nhất, ở "đầu ra"), rồi "lan truyền ngược" từng bước về các tầng trước, mỗi bước chỉ cần nhân thêm 1 đạo hàm cục bộ (local gradient) của tầng đó — không cần tính lại từ đầu cho mỗi tham số. Chi tiết đầy đủ với code from scratch ở [Bài 21](./21_backpropagation.md).

---

## PHẦN B — Cài đặt & Minh họa bằng code

```python
import numpy as np

def numerical_gradient(f, x, h=1e-5):
	grad = np.zeros_like(x, dtype=float)
	for i in range(len(x)):
		x_plus = x.copy(); x_plus[i] += h
		x_minus = x.copy(); x_minus[i] -= h
		grad[i] = (f(x_plus) - f(x_minus)) / (2 * h)
	return grad

# Mục 2, Công thức 1: verify gradient của dạng tuyến tính
a = np.array([2.0, 3.0])
f_linear = lambda x: a @ x
x0 = np.array([1.0, 1.0])
print(numerical_gradient(f_linear, x0))  # xấp xỉ [2, 3] = a — khớp

# Mục 2, Công thức 2: verify gradient của dạng toàn phương
A = np.array([[2.0, 1.0], [1.0, 2.0]])  # đối xứng
f_quad = lambda x: x @ A @ x
print(numerical_gradient(f_quad, x0))       # gradient numerical
print(2 * A @ x0)                             # công thức 2A*x — phải khớp

# Mục 3: verify gradient MSE
def mse_loss(w, X, y):
	m = len(y)
	return (1/m) * np.sum((X @ w - y)**2)

def mse_gradient(w, X, y):
	m = len(y)
	return (2/m) * X.T @ (X @ w - y)

X = np.array([[1, 1], [1, 2], [1, 3]], dtype=float)
y = np.array([2, 3, 5], dtype=float)
w0 = np.array([0.0, 0.0])
print(mse_gradient(w0, X, y))
print(numerical_gradient(lambda w: mse_loss(w, X, y), w0))  # phải xấp xỉ bằng nhau

# Mục 5: chain rule qua mạng 2 lớp — verify từng mắt xích
def sigmoid(z): return 1 / (1 + np.exp(-z))
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

	dL_dyhat = 2 * (y_hat - y)
	dyhat_da = w2
	da_dz = sigmoid_derivative(z)
	dz_dw1 = x

	dL_dw1 = dL_dyhat * dyhat_da * da_dz * dz_dw1  # nhân liên tiếp các mắt xích
	dL_dw2 = dL_dyhat * a
	return dL_dw1, dL_dw2, loss

grad_w1, grad_w2, loss = backward(x=1.0, y=1.0, w1=0.5, w2=0.8)
print(f"Gradient w1: {grad_w1:.4f}, Gradient w2: {grad_w2:.4f}, Loss: {loss:.4f}")
```

## Bài tập

1. **Verify 2 công thức nền tảng**: dùng `numerical_gradient`, verify công thức gradient dạng tuyến tính và toàn phương với 2 bộ số $\vec{a}$/$A$ khác (tự chọn), so khớp với suy luận tay ở mục 2.
2. **Tự suy luận lại gradient MSE**: che phần "Bước 1-3" ở mục 3, tự làm lại phép khai triển và áp dụng công thức — so sánh với lời giải trong bài.
3. **Chain rule bằng tay**: với $x=2,y=0,w_1=0.3,w_2=0.5$ ở mục 5, tính tay từng mắt xích ($\partial z/\partial w_1$, $\partial a/\partial z$...) rồi nhân lại, verify bằng `backward()`.
4. **Mở rộng mạng 2 neuron ở lớp ẩn**: thay vì 1 neuron lớp ẩn như ví dụ, viết lại `forward`/`backward` cho 2 neuron — đây là bước đệm trực tiếp cho [Bài 21](./21_backpropagation.md).

## Tiếp theo
→ [Bài 7: Tối ưu hóa không ràng buộc & Taylor Series](./7_calculus_optimization.md)
