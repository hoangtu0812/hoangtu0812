# Bài 15: Logistic Regression & Classification Cơ Bản

## Mục tiêu
- Hàm sigmoid, mô hình xác suất cho phân loại.
- Cross-Entropy Loss (đã chứng minh từ MLE ở [Bài 10](./10_statistics_inference.md)).
- Decision boundary.

## 1. Vì sao không dùng Linear Regression cho Classification?

Nếu dùng $\hat{y} = \vec{w}^T\vec{x}+b$ trực tiếp cho nhãn 0/1, giá trị dự đoán có thể ra ngoài khoảng $[0,1]$ — vô nghĩa khi diễn giải là xác suất. Cần 1 hàm "nén" giá trị về $[0,1]$: hàm **sigmoid**.

## 2. Hàm Sigmoid

$$\sigma(z) = \frac{1}{1+e^{-z}}$$

```python
import numpy as np

def sigmoid(z):
	return 1 / (1 + np.exp(-z))

z = np.linspace(-10, 10, 100)
import matplotlib.pyplot as plt
plt.plot(z, sigmoid(z))
plt.axhline(0.5, color='r', linestyle='--')
plt.title("Hàm Sigmoid")
plt.show()
```

Tính chất: $\sigma(z) \to 0$ khi $z \to -\infty$, $\sigma(z) \to 1$ khi $z \to \infty$, $\sigma(0) = 0.5$. Đạo hàm đẹp: $\sigma'(z) = \sigma(z)(1-\sigma(z))$ — dùng trực tiếp trong Backpropagation ([Bài 21](./21_backpropagation.md)).

## 3. Mô hình Logistic Regression

$$P(y=1|\vec{x}) = \sigma(\vec{w}^T\vec{x}+b) = \frac{1}{1+e^{-(\vec{w}^T\vec{x}+b)}}$$

Dự đoán nhãn: $\hat{y} = 1$ nếu $P(y=1|\vec{x}) \geq 0.5$, ngược lại $\hat{y} = 0$.

## 4. Hàm mất mát — Binary Cross-Entropy

Như đã chứng minh chi tiết từ MLE ở [Bài 10 mục 3](./10_statistics_inference.md):

$$L(\vec{w}) = -\frac{1}{m}\sum_{i=1}^m \left[y_i\log(\hat{y_i}) + (1-y_i)\log(1-\hat{y_i})\right]$$

**Vì sao không dùng MSE cho classification?** MSE với sigmoid tạo ra hàm loss **không lồi** (nhiều cực tiểu cục bộ) — Gradient Descent dễ mắc kẹt. Cross-Entropy giữ tính lồi ([Bài 7 mục 3](./7_calculus_optimization.md), [Bài 11 mục 4](./11_discrete_math_optimization.md)), đảm bảo hội tụ về nghiệm tối ưu toàn cục.

## 5. Gradient của Cross-Entropy — đẹp bất ngờ, giống hệt Linear Regression!

$$\nabla_{\vec{w}} L = \frac{1}{m}X^T(\sigma(X\vec{w}) - \vec{y})$$

So sánh với gradient MSE của Linear Regression ([Bài 6 mục 3](./6_calculus_matrix_calculus.md)): $\nabla L = \frac{2}{m}X^T(X\vec{w}-\vec{y})$ — cùng dạng "$X^T \times \text{sai số}$", chỉ khác $X\vec{w}$ được thay bằng $\sigma(X\vec{w})$. Đây không phải trùng hợp: cả 2 đều là hệ quả của việc chọn loss = MLE cho phân phối tương ứng ([Bài 10](./10_statistics_inference.md)).

## 6. Cài đặt from scratch

```python
class LogisticRegression:
	def __init__(self, learning_rate=0.1, n_iters=1000):
		self.lr = learning_rate
		self.n_iters = n_iters
		self.loss_history = []

	def fit(self, X, y):
		m, n = X.shape
		X_b = np.hstack([np.ones((m, 1)), X])
		self.w = np.zeros(n + 1)

		for _ in range(self.n_iters):
			z = X_b @ self.w
			y_pred = sigmoid(z)
			error = y_pred - y
			grad = (1/m) * X_b.T @ error
			self.w -= self.lr * grad

			eps = 1e-15
			y_pred_clip = np.clip(y_pred, eps, 1-eps)
			loss = -np.mean(y*np.log(y_pred_clip) + (1-y)*np.log(1-y_pred_clip))
			self.loss_history.append(loss)
		return self

	def predict_proba(self, X):
		X_b = np.hstack([np.ones((X.shape[0], 1)), X])
		return sigmoid(X_b @ self.w)

	def predict(self, X, threshold=0.5):
		return (self.predict_proba(X) >= threshold).astype(int)
```

## 7. Decision Boundary — ý nghĩa hình học

$\hat{y} = 1$ khi $\vec{w}^T\vec{x}+b \geq 0$ — biên quyết định (decision boundary) chính là siêu phẳng (hyperplane) $\vec{w}^T\vec{x}+b = 0$, tương tự siêu phẳng trong SVM ([Bài 19](./19_svm_knn_clustering.md)).

```python
def plot_decision_boundary(model, X, y):
	x1_min, x1_max = X[:,0].min()-1, X[:,0].max()+1
	x2_min, x2_max = X[:,1].min()-1, X[:,1].max()+1
	xx1, xx2 = np.meshgrid(np.linspace(x1_min, x1_max, 200), np.linspace(x2_min, x2_max, 200))
	grid = np.c_[xx1.ravel(), xx2.ravel()]
	preds = model.predict(grid).reshape(xx1.shape)

	plt.contourf(xx1, xx2, preds, alpha=0.3)
	plt.scatter(X[:,0], X[:,1], c=y, edgecolors='k')
	plt.title("Decision Boundary")
	plt.show()
```

## 8. Multi-class Classification — Softmax Regression

Với > 2 lớp, mở rộng sigmoid thành **softmax**:

$$P(y=k|\vec{x}) = \frac{e^{z_k}}{\sum_{j=1}^K e^{z_j}}, \quad z_k = \vec{w_k}^T\vec{x}+b_k$$

Softmax "chuẩn hóa" $K$ giá trị thành phân phối xác suất (tổng = 1) — dùng ở lớp cuối của hầu hết mạng neural phân loại nhiều lớp ([Bài 20](./20_neural_networks.md)).

```python
def softmax(z):
	exp_z = np.exp(z - np.max(z, axis=1, keepdims=True))  # trừ max để ổn định số học
	return exp_z / np.sum(exp_z, axis=1, keepdims=True)
```

## Ví dụ đầy đủ

```python
import numpy as np
from sklearn.datasets import make_classification
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler

X, y = make_classification(n_samples=300, n_features=2, n_redundant=0, n_clusters_per_class=1, random_state=42)
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

scaler = StandardScaler()
X_train = scaler.fit_transform(X_train)
X_test = scaler.transform(X_test)

model = LogisticRegression(learning_rate=0.5, n_iters=500).fit(X_train, y_train)
y_pred = model.predict(X_test)

accuracy = np.mean(y_pred == y_test)
print(f"Accuracy: {accuracy:.2%}")
```

## Bài tập

1. **Sigmoid & đạo hàm**: verify công thức $\sigma'(z) = \sigma(z)(1-\sigma(z))$ bằng `numerical_derivative` ([Bài 5](./5_calculus_derivatives.md)).
2. **Cài đặt Logistic Regression**: dùng code mẫu trên làm nền, tự viết lại, huấn luyện trên dataset phân loại nhị phân tự tạo (`make_classification`), đạt accuracy > 85%.
3. **Vẽ decision boundary**: dùng code mẫu mục 7, trực quan hóa ranh giới phân loại model vừa huấn luyện.
4. **So sánh MSE vs Cross-Entropy**: thử thay hàm loss trong `fit()` bằng MSE thay vì Cross-Entropy, quan sát quá trình hội tụ (`loss_history`) khác biệt thế nào — liên hệ lý do đã giải thích ở mục 4.

## Tiếp theo
→ [Bài 16: Đánh giá mô hình (Model Evaluation)](./16_model_evaluation.md)
