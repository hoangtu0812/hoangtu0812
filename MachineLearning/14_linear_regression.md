# Bài 14: Linear Regression — Từ Toán Tới Code

## Mục tiêu
- Ghép toàn bộ toán Phần I lại thành thuật toán ML đầu tiên: Linear Regression.
- Cài đặt bằng cả Normal Equation và Gradient Descent, rồi so sánh với scikit-learn.

## 1. Mô hình toán học

$$\hat{y} = \vec{w}^T\vec{x} + b = w_1x_1 + w_2x_2 + \dots + w_nx_n + b$$

Đây chính là dot product ([Bài 2](./2_linear_algebra_vectors.md)) — với toàn bộ dataset viết dưới dạng ma trận ([Bài 3](./3_linear_algebra_matrices.md)):

$$\hat{\vec{y}} = X\vec{w} + b$$

## 2. Hàm mất mát (Loss Function) — MSE

$$L(\vec{w}, b) = \frac{1}{m}\sum_{i=1}^m (y_i - \hat{y_i})^2$$

Như đã chứng minh ở [Bài 10 mục 2](./10_statistics_inference.md), MSE chính là **MLE dưới giả định nhiễu Gaussian** — không phải lựa chọn tùy tiện.

## 3. Cách 1: Normal Equation (nghiệm đóng)

$$\vec{w}^* = (X^TX)^{-1}X^T\vec{y}$$

(đã chi tiết ở [Bài 3 mục 7](./3_linear_algebra_matrices.md))

```python
import numpy as np

class LinearRegressionNormalEq:
	def fit(self, X, y):
		X_b = np.hstack([np.ones((X.shape[0], 1)), X])
		self.w = np.linalg.inv(X_b.T @ X_b) @ X_b.T @ y
		return self

	def predict(self, X):
		X_b = np.hstack([np.ones((X.shape[0], 1)), X])
		return X_b @ self.w
```

**Ưu điểm:** nghiệm chính xác, không cần chọn learning rate. **Nhược điểm:** $O(n^3)$ do nghịch đảo ma trận — không khả thi khi $n$ (số đặc trưng) lớn.

## 4. Cách 2: Gradient Descent

Gradient của MSE (đã tính ở [Bài 6 mục 3](./6_calculus_matrix_calculus.md)):

$$\nabla_{\vec{w}} L = \frac{2}{m}X^T(X\vec{w} - \vec{y})$$

```python
class LinearRegressionGD:
	def __init__(self, learning_rate=0.01, n_iters=1000):
		self.lr = learning_rate
		self.n_iters = n_iters
		self.loss_history = []

	def fit(self, X, y):
		m, n = X.shape
		X_b = np.hstack([np.ones((m, 1)), X])
		self.w = np.zeros(n + 1)

		for _ in range(self.n_iters):
			y_pred = X_b @ self.w
			error = y_pred - y
			grad = (2/m) * X_b.T @ error
			self.w -= self.lr * grad

			loss = np.mean(error**2)
			self.loss_history.append(loss)
		return self

	def predict(self, X):
		X_b = np.hstack([np.ones((X.shape[0], 1)), X])
		return X_b @ self.w
```

## 5. Chuẩn hóa dữ liệu — BẮT BUỘC khi dùng Gradient Descent

Nếu các đặc trưng có thang đo khác nhau nhiều (vd diện tích ~100, số phòng ~3), Gradient Descent hội tụ RẤT chậm hoặc dao động (liên hệ [Bài 7 mục 5](./7_calculus_optimization.md) về ảnh hưởng learning rate). Luôn chuẩn hóa trước ([Bài 9 mục 4](./9_probability_distributions.md), [Bài 12](./12_numpy_pandas.md)):

```python
def standardize(X):
	mean, std = X.mean(axis=0), X.std(axis=0)
	return (X - mean) / std, mean, std
```

## 6. So sánh với scikit-learn — verify implementation của bạn đúng

```python
from sklearn.linear_model import LinearRegression as SklearnLR
from sklearn.preprocessing import StandardScaler

X = np.array([[120, 3], [85, 2], [200, 4], [60, 1], [150, 3]], dtype=float)
y = np.array([500, 350, 800, 250, 600], dtype=float)

scaler = StandardScaler()
X_scaled = scaler.fit_transform(X)

# Model tự viết
my_model = LinearRegressionGD(learning_rate=0.1, n_iters=1000).fit(X_scaled, y)
print("My model weights:", my_model.w)

# scikit-learn
sklearn_model = SklearnLR().fit(X_scaled, y)
print("Sklearn coef:", sklearn_model.coef_, "intercept:", sklearn_model.intercept_)
```

2 kết quả phải xấp xỉ nhau — đây là bước "unit test" quan trọng khi tự cài thuật toán ML from scratch (liên hệ tinh thần testing ở [Go Bài 14](../Go/14_testing.md)/[Python Bài 14](../Python/14_testing.md)).

## 7. Multicollinearity — vấn đề thực tế liên hệ Bài 3

Nếu 2 đặc trưng phụ thuộc tuyến tính mạnh ([Bài 3 mục 7](./3_linear_algebra_matrices.md)), $X^TX$ gần suy biến, Normal Equation cho kết quả không ổn định (trọng số rất lớn/nhỏ bất thường). Kiểm tra bằng ma trận tương quan ([Bài 13](./13_data_visualization_eda.md)) trước khi huấn luyện, hoặc dùng Ridge Regression ([Bài 17](./17_regularization.md)) để giảm ảnh hưởng.

## 8. Polynomial Regression — mở rộng Linear Regression cho quan hệ phi tuyến

```python
from sklearn.preprocessing import PolynomialFeatures

poly = PolynomialFeatures(degree=2)
X_poly = poly.fit_transform(X)  # thêm các cột x^2, x1*x2...

model = SklearnLR().fit(X_poly, y)
```

Vẫn là "Linear" Regression vì **tuyến tính theo tham số $\vec{w}$** — chỉ đặc trưng đầu vào được biến đổi phi tuyến. Cẩn thận: degree cao dễ gây overfitting (chi tiết ở [Bài 16-17](./16_model_evaluation.md)).

## Ví dụ đầy đủ

```python
import numpy as np
import matplotlib.pyplot as plt

def generate_data(n=100, noise=1.0, seed=42):
	np.random.seed(seed)
	X = np.random.rand(n, 1) * 10
	y = 3 * X.flatten() + 5 + np.random.randn(n) * noise
	return X, y

if __name__ == "__main__":
	X, y = generate_data()

	model = LinearRegressionGD(learning_rate=0.01, n_iters=500)
	model.fit(X, y)

	plt.plot(model.loss_history)
	plt.xlabel("Iteration"); plt.ylabel("MSE Loss")
	plt.title("Quá trình hội tụ của Gradient Descent")
	plt.show()

	print("Weights (bias, w1):", model.w)  # xấp xỉ [5, 3]

	X_test = np.array([[2], [8]])
	print("Dự đoán:", model.predict(X_test))
```

## Bài tập

1. **Cài đặt cả 2 cách**: viết `LinearRegressionNormalEq` và `LinearRegressionGD` như trên, verify chúng cho kết quả gần giống nhau trên cùng dataset.
2. **Chuẩn hóa & learning rate**: thử huấn luyện KHÔNG chuẩn hóa dữ liệu (đặc trưng thang đo rất khác nhau), quan sát Gradient Descent hội tụ chậm/dao động; sau đó chuẩn hóa và so sánh.
3. **So sánh với scikit-learn**: verify model tự viết cho kết quả gần giống `sklearn.linear_model.LinearRegression`.
4. **Polynomial Regression**: dùng `PolynomialFeatures` fit dữ liệu có quan hệ phi tuyến (vd $y = x^2 + \text{nhiễu}$), so sánh với Linear Regression bậc 1 — quan sát bậc nào fit tốt hơn.

## Tiếp theo
→ [Bài 15: Logistic Regression & Classification cơ bản](./15_logistic_regression.md)
