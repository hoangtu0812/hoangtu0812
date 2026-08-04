# Bài 7: Tối ưu hóa không ràng buộc & Taylor Series

## Mục tiêu
- Điểm cực trị, Hessian matrix, điều kiện tối ưu bậc 1/bậc 2.
- Khai triển Taylor.
- Hiểu vì sao Gradient Descent hội tụ, learning rate ảnh hưởng thế nào.

## 1. Điểm cực trị (Critical Point)

Điểm $x^*$ là **điểm cực trị** của $f$ nếu $f'(x^*) = 0$ (với hàm nhiều biến: $\nabla f(x^*) = \vec{0}$) — đây là **điều kiện cần bậc 1 (First-Order Condition)**.

Điểm cực trị có thể là:
- **Cực tiểu (minimum)**: giá trị nhỏ nhất cục bộ.
- **Cực đại (maximum)**: giá trị lớn nhất cục bộ.
- **Điểm yên ngựa (saddle point)**: không phải cực đại lẫn cực tiểu — tăng theo hướng này, giảm theo hướng khác.

**Trong ML:** huấn luyện model = tìm $\vec{w}^*$ sao cho hàm loss $L(\vec{w})$ đạt cực tiểu — đây chính là bài toán tối ưu hóa không ràng buộc.

## 2. Hessian Matrix — điều kiện bậc 2, phân biệt cực tiểu/cực đại/saddle point

Hessian là ma trận đạo hàm riêng bậc 2:

$$H = \begin{bmatrix} \partial^2f/\partial x_1^2 & \partial^2f/\partial x_1\partial x_2 \\ \partial^2f/\partial x_2\partial x_1 & \partial^2f/\partial x_2^2 \end{bmatrix}$$

Tại điểm cực trị $x^*$ (đã có $\nabla f(x^*) = 0$):
- $H$ **xác định dương** (positive definite, mọi trị riêng > 0) → $x^*$ là **cực tiểu**.
- $H$ **xác định âm** (mọi trị riêng < 0) → $x^*$ là **cực đại**.
- $H$ có cả trị riêng dương và âm → $x^*$ là **saddle point**.

```python
import numpy as np

def hessian_numerical(f, x, h=1e-4):
	n = len(x)
	H = np.zeros((n, n))
	for i in range(n):
		for j in range(n):
			x_pp = x.copy(); x_pp[i] += h; x_pp[j] += h
			x_pm = x.copy(); x_pm[i] += h; x_pm[j] -= h
			x_mp = x.copy(); x_mp[i] -= h; x_mp[j] += h
			x_mm = x.copy(); x_mm[i] -= h; x_mm[j] -= h
			H[i,j] = (f(x_pp) - f(x_pm) - f(x_mp) + f(x_mm)) / (4*h*h)
	return H

def f(x):
	return x[0]**2 + x[1]**2  # hàm bát úp — cực tiểu tại (0,0)

H = hessian_numerical(f, np.array([1.0, 1.0]))
print(H)  # xấp xỉ [[2, 0], [0, 2]] — cả 2 trị riêng đều dương -> cực tiểu
```

**Ứng dụng ML quan trọng:** saddle point là lý do chính khiến việc huấn luyện mạng neural sâu bị "chững lại" tạm thời (gradient gần 0 nhưng chưa phải cực tiểu thật) — các optimizer nâng cao như Adam ([Bài 22](./22_optimizers_dl_regularization.md)) được thiết kế để vượt qua vấn đề này tốt hơn Gradient Descent thuần.

## 3. Hàm lồi (Convex Function) — trường hợp "dễ" nhất để tối ưu

Hàm $f$ là **lồi** nếu Hessian của nó luôn xác định dương (hoặc bán dương) tại MỌI điểm — khi đó **mọi cực tiểu cục bộ đều là cực tiểu toàn cục**, không có nguy cơ "mắc kẹt" ở cực tiểu địa phương tệ.

$$f(\lambda x_1 + (1-\lambda)x_2) \leq \lambda f(x_1) + (1-\lambda)f(x_2), \quad \forall \lambda \in [0,1]$$

**Ứng dụng ML:** hàm loss MSE của Linear Regression ([Bài 14](./14_linear_regression.md)) và Cross-Entropy của Logistic Regression ([Bài 15](./15_logistic_regression.md)) đều **lồi** — Gradient Descent được đảm bảo hội tụ về cực tiểu toàn cục. Ngược lại, hàm loss của mạng neural sâu **KHÔNG lồi** — đây là lý do Deep Learning khó tối ưu hơn nhiều, cần các kỹ thuật nâng cao ([Bài 22](./22_optimizers_dl_regularization.md)). Chi tiết về tối ưu hóa lồi ở [Bài 11](./11_discrete_math_optimization.md).

## 4. Khai triển Taylor — xấp xỉ hàm bằng đa thức

Khai triển Taylor bậc 2 quanh điểm $x_0$:

$$f(x) \approx f(x_0) + f'(x_0)(x-x_0) + \frac{1}{2}f''(x_0)(x-x_0)^2$$

Với hàm nhiều biến (dạng ma trận, dùng gradient & Hessian):

$$f(\vec{x}) \approx f(\vec{x_0}) + \nabla f(\vec{x_0})^T(\vec{x}-\vec{x_0}) + \frac{1}{2}(\vec{x}-\vec{x_0})^TH(\vec{x}-\vec{x_0})$$

```python
def f(x):
	return np.exp(x)

def taylor_approx_2nd_order(f, x0, x, h=1e-5):
	f_x0 = f(x0)
	f_prime = (f(x0+h) - f(x0-h)) / (2*h)
	f_double_prime = (f(x0+h) - 2*f(x0) + f(x0-h)) / (h**2)
	return f_x0 + f_prime * (x - x0) + 0.5 * f_double_prime * (x - x0)**2

x0 = 0
for x in [0.1, 0.5, 1.0, 2.0]:
	approx = taylor_approx_2nd_order(f, x0, x)
	real = f(x)
	print(f"x={x}: xấp xỉ={approx:.4f}, thật={real:.4f}, sai số={abs(approx-real):.4f}")
```

Sai số tăng dần khi $x$ xa $x_0$ — xấp xỉ Taylor chỉ chính xác **cục bộ** (gần điểm khai triển).

## 5. Vì sao Gradient Descent hoạt động — giải thích bằng Taylor

Khai triển Taylor bậc 1 quanh điểm hiện tại $\vec{w_t}$:

$$L(\vec{w_t} - \eta\nabla L(\vec{w_t})) \approx L(\vec{w_t}) - \eta\|\nabla L(\vec{w_t})\|^2$$

Vì $\|\nabla L(\vec{w_t})\|^2 \geq 0$, nên với $\eta$ (learning rate) đủ nhỏ, giá trị hàm loss **luôn giảm** sau mỗi bước cập nhật $\vec{w_{t+1}} = \vec{w_t} - \eta\nabla L(\vec{w_t})$ — đây chính là chứng minh trực giác vì sao Gradient Descent hội tụ, và tại sao chọn $\eta$ quá lớn có thể khiến xấp xỉ Taylor bậc 1 không còn chính xác (loss có thể TĂNG thay vì giảm — hiện tượng "divergence").

```python
def demonstrate_learning_rate_effect(f, grad_f, start, learning_rates, steps=10):
	for lr in learning_rates:
		w = start
		for _ in range(steps):
			w = w - lr * grad_f(w)
		print(f"learning_rate={lr}: w cuối = {w:.4f}, loss = {f(w):.4f}")

f = lambda w: (w - 5)**2
grad_f = lambda w: 2*(w - 5)

demonstrate_learning_rate_effect(f, grad_f, start=0.0, learning_rates=[0.01, 0.1, 0.5, 1.0, 1.1])
# lr quá nhỏ: hội tụ chậm. lr vừa: hội tụ tốt. lr quá lớn (>=1.0 với hàm này): DAO ĐỘNG hoặc PHÂN KỲ
```

## Bài tập

1. **Hessian & phân loại điểm cực trị**: tính Hessian bằng tay cho $f(x,y) = x^2 - y^2$ (hàm saddle kinh điển) tại điểm $(0,0)$, xác định đây là cực tiểu/cực đại/saddle point; verify bằng `hessian_numerical`.
2. **Kiểm tra hàm lồi**: viết code kiểm tra Hessian có xác định dương tại nhiều điểm ngẫu nhiên hay không (dùng `np.linalg.eigvals(H) > 0`) cho hàm MSE loss ([Bài 6 mục 3](./6_calculus_matrix_calculus.md)) — xác nhận nó luôn lồi.
3. **Taylor approximation**: dùng code mẫu mục 4, thử với hàm $\sin(x)$ quanh $x_0 = 0$, vẽ đồ thị so sánh hàm thật với xấp xỉ Taylor bậc 2 trên khoảng $[-2, 2]$ (liên hệ [Bài 13](./13_data_visualization_eda.md)).
4. **Learning rate**: dùng code mẫu mục 5, thử nhiều learning rate khác nhau, quan sát và giải thích hiện tượng hội tụ chậm/vừa/dao động/phân kỳ.

## Tổng kết phần Giải tích
Bạn đã nắm đạo hàm, gradient, đạo hàm ma trận, và lý do toán học vì sao Gradient Descent hoạt động. Tiếp theo là Xác suất & Thống kê — nền tảng để hiểu hàm loss đến từ đâu và cách đánh giá độ tin cậy của model.

## Tiếp theo
→ [Bài 8: Xác suất cơ bản & Định lý Bayes](./8_probability_basics.md)
