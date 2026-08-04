# Bài 11: Toán rời rạc, Độ phức tạp thuật toán & Tối ưu hóa lồi

## Mục tiêu
- Tổ hợp/hoán vị cơ bản, lý thuyết đồ thị sơ lược, Big-O.
- Hàm lồi & tối ưu hóa lồi.
- **Gradient Descent** và các biến thể — thuật toán tối ưu dùng trong hầu hết ML/DL.

## 1. Tổ hợp & Hoán vị — nền tảng đếm xác suất

**Hoán vị (Permutation)**: số cách sắp xếp $n$ phần tử: $n! = n(n-1)(n-2)\cdots 1$

**Tổ hợp (Combination)**: số cách chọn $k$ phần tử từ $n$ (không quan tâm thứ tự):

$$\binom{n}{k} = \frac{n!}{k!(n-k)!}$$

```python
from math import comb, perm

print(perm(5, 3))   # 60 — số cách sắp xếp 3 trong 5 phần tử có thứ tự
print(comb(5, 3))    # 10 — số cách chọn 3 trong 5 phần tử không thứ tự
```

**Ứng dụng ML:** công thức Binomial ([Bài 9](./9_probability_distributions.md)) dùng trực tiếp $\binom{n}{k}$; k-fold cross-validation ([Bài 16](./16_model_evaluation.md)) liên quan tới cách chia tổ hợp dữ liệu; feature selection đôi khi cần xét $\binom{n}{k}$ tổ hợp đặc trưng.

## 2. Lý thuyết đồ thị sơ lược

Đồ thị $G = (V, E)$ gồm tập đỉnh $V$ và tập cạnh $E$. Trong ML:
- **Decision Tree/Random Forest** ([Bài 18](./18_trees_ensembles.md)) là cấu trúc cây (1 dạng đặc biệt của đồ thị) — mỗi node là 1 quyết định rẽ nhánh.
- **Graph Neural Network** (ngoài phạm vi lộ trình cơ bản này) áp dụng deep learning trực tiếp lên dữ liệu dạng đồ thị (mạng xã hội, phân tử hóa học).
- **Neural Network** bản thân cũng là 1 đồ thị có hướng (computational graph) — mỗi node là 1 phép toán, cạnh là luồng dữ liệu, chính là cấu trúc mà PyTorch/TensorFlow dùng để tự động tính đạo hàm (autograd).

## 3. Độ phức tạp thuật toán (Big-O) — áp dụng cho ML

Nhắc lại khái niệm Big-O (nếu đã quen từ lập trình competitive/[Go](../Go/ROADMAP.md)/[Python](../Python/ROADMAP.md)) trong ngữ cảnh ML:

| Thuật toán | Độ phức tạp (huấn luyện) | Ghi chú |
|---|---|---|
| Linear Regression (Normal Equation) | $O(n^3 + n^2m)$ | $n$=số đặc trưng, $m$=số mẫu — nghịch đảo ma trận $O(n^3)$ tốn kém khi $n$ lớn |
| Linear Regression (Gradient Descent) | $O(nm)$ mỗi bước | Rẻ hơn Normal Equation khi $n$ lớn — lý do dùng GD cho dataset nhiều đặc trưng |
| k-NN (dự đoán) | $O(m)$ mỗi truy vấn | Không có "huấn luyện" thật sự, nhưng dự đoán chậm với dataset lớn |
| k-Means | $O(kmi)$ | $k$=số cụm, $i$=số vòng lặp |
| Decision Tree | $O(mn\log m)$ | Xây cây |

Đây chính là lý do **Normal Equation không dùng được cho Deep Learning** (hàng triệu tham số → $n^3$ là bất khả thi), buộc phải dùng Gradient Descent.

## 4. Hàm lồi & Tập lồi (Convex Set/Function) — nhắc lại từ Bài 7, mở rộng

Tập $S$ là **lồi** nếu với mọi $x_1, x_2 \in S$, đoạn thẳng nối chúng cũng nằm trong $S$:

$$\lambda x_1 + (1-\lambda)x_2 \in S, \quad \forall \lambda \in [0,1]$$

Hàm $f$ lồi $\Leftrightarrow$ tập "trên đồ thị" (epigraph) của nó là tập lồi $\Leftrightarrow$ Hessian $\succeq 0$ mọi nơi ([Bài 7 mục 3](./7_calculus_optimization.md)).

**Tính chất cực kỳ quan trọng:** với hàm lồi, **mọi cực tiểu cục bộ = cực tiểu toàn cục**. Đây là lý do Linear/Logistic Regression "dễ" tối ưu — Gradient Descent CHẮC CHẮN tìm được nghiệm tốt nhất (với learning rate phù hợp), không lo bị "mắc kẹt".

## 5. Gradient Descent — thuật toán trung tâm của ML/DL

Công thức cập nhật (đã giới thiệu ở [Bài 5](./5_calculus_derivatives.md), [Bài 7](./7_calculus_optimization.md)):

$$\vec{w}_{t+1} = \vec{w}_t - \eta \nabla L(\vec{w}_t)$$

```python
import numpy as np

def gradient_descent(grad_f, w_init, learning_rate=0.01, n_iters=1000, tol=1e-6):
	w = w_init.copy()
	history = [w.copy()]
	for i in range(n_iters):
		grad = grad_f(w)
		w = w - learning_rate * grad
		history.append(w.copy())
		if np.linalg.norm(grad) < tol:  # điều kiện dừng sớm khi gradient đủ nhỏ
			print(f"Hội tụ sau {i+1} bước")
			break
	return w, history
```

## 6. 3 biến thể của Gradient Descent — khác biệt ở LƯỢNG DỮ LIỆU dùng mỗi bước

### Batch Gradient Descent — dùng TOÀN BỘ dataset mỗi bước cập nhật

```python
def batch_gradient_descent(X, y, w_init, lr=0.01, n_iters=100):
	w = w_init.copy()
	m = len(y)
	for _ in range(n_iters):
		grad = (2/m) * X.T @ (X @ w - y)  # dùng TOÀN BỘ m mẫu
		w = w - lr * grad
	return w
```
Ổn định, hội tụ mượt, nhưng **chậm** với dataset lớn (mỗi bước phải duyệt hết dữ liệu).

### Stochastic Gradient Descent (SGD) — dùng 1 mẫu ngẫu nhiên mỗi bước

```python
def stochastic_gradient_descent(X, y, w_init, lr=0.01, n_epochs=50):
	w = w_init.copy()
	m = len(y)
	for epoch in range(n_epochs):
		indices = np.random.permutation(m)
		for i in indices:
			xi, yi = X[i:i+1], y[i:i+1]
			grad = 2 * xi.T @ (xi @ w - yi)  # CHỈ dùng 1 mẫu
			w = w - lr * grad.flatten()
	return w
```
Nhanh hơn nhiều mỗi bước, nhưng đường hội tụ "nhiễu" (dao động) vì mỗi mẫu chỉ là ước lượng thô của gradient thật.

### Mini-batch Gradient Descent — cân bằng giữa 2 cách trên (dùng phổ biến NHẤT trong Deep Learning)

```python
def minibatch_gradient_descent(X, y, w_init, lr=0.01, batch_size=32, n_epochs=50):
	w = w_init.copy()
	m = len(y)
	for epoch in range(n_epochs):
		indices = np.random.permutation(m)
		for start in range(0, m, batch_size):
			batch_idx = indices[start:start+batch_size]
			X_batch, y_batch = X[batch_idx], y[batch_idx]
			grad = (2/len(batch_idx)) * X_batch.T @ (X_batch @ w - y_batch)
			w = w - lr * grad
	return w
```

`batch_size` (thường 32, 64, 128, 256) là hyperparameter cân bằng tốc độ và độ ổn định — đây chính là biến thể dùng khi huấn luyện mạng neural ([Bài 20-22](./20_neural_networks.md)).

## 7. So sánh trực quan 3 biến thể

| | Batch GD | SGD | Mini-batch GD |
|---|---|---|---|
| Dữ liệu/bước | Toàn bộ $m$ mẫu | 1 mẫu | $b$ mẫu ($1 < b < m$) |
| Tốc độ/bước | Chậm | Rất nhanh | Vừa phải |
| Độ ổn định hội tụ | Mượt | Nhiễu, dao động | Khá mượt |
| Dùng thực tế trong DL | Hiếm (dataset lớn không khả thi) | Hiếm khi dùng thuần túy | **Chuẩn công nghiệp** |

## Ví dụ đầy đủ: so sánh 3 biến thể trên cùng bài toán

```python
import numpy as np

np.random.seed(42)
m = 200
X = np.hstack([np.ones((m,1)), np.random.randn(m,1)*2])
true_w = np.array([3, 5])
y = X @ true_w + np.random.randn(m) * 0.5

w_init = np.zeros(2)
w_batch = batch_gradient_descent(X, y, w_init, lr=0.05, n_iters=200)
w_sgd = stochastic_gradient_descent(X, y, w_init, lr=0.01, n_epochs=20)
w_minibatch = minibatch_gradient_descent(X, y, w_init, lr=0.05, batch_size=32, n_epochs=50)

print("True w:", true_w)
print("Batch GD:", w_batch)
print("SGD:", w_sgd)
print("Mini-batch GD:", w_minibatch)
```

## Bài tập

1. **Tổ hợp/hoán vị**: tính tay số cách chia 20 mẫu dữ liệu thành 5-fold cross-validation (mỗi fold 4 mẫu) — liên hệ [Bài 16](./16_model_evaluation.md).
2. **Kiểm tra hàm lồi**: viết hàm kiểm tra Hessian dương tại nhiều điểm ngẫu nhiên cho 3 hàm: $f(x)=x^2$ (lồi), $f(x)=x^4-x^2$ (không lồi toàn cục), $f(x,y)=x^2+y^2$ (lồi).
3. **Implement 3 biến thể GD**: dùng code mẫu trên làm nền, tự viết lại cả 3, chạy trên cùng dataset, vẽ đồ thị loss theo epoch (liên hệ [Bài 13](./13_data_visualization_eda.md)) để so sánh trực quan tốc độ hội tụ và độ nhiễu.
4. **Nâng cao**: thử nhiều `batch_size` khác nhau (8, 32, 128) cho Mini-batch GD, quan sát ảnh hưởng tới tốc độ hội tụ và độ dao động của loss.

## Tổng kết Phần I — Toán Nền Tảng
Bạn đã hoàn thành 4 mảng toán cốt lõi: Đại số tuyến tính (biểu diễn dữ liệu), Giải tích (cách model học qua gradient), Xác suất & Thống kê (nguồn gốc hàm loss), Toán rời rạc & Tối ưu hóa (thuật toán tối ưu). Từ đây, mọi công thức ở Phần II-III sẽ được giải thích bằng cách trỏ ngược lại các bài này.

## Tiếp theo
→ [Bài 12: NumPy & Pandas cho Data Science](./12_numpy_pandas.md)
