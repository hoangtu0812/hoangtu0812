# Bài 17: Regularization & Bias-Variance Tradeoff

## Mục tiêu
- Ridge (L2), Lasso (L1), Elastic Net.
- Bias-Variance decomposition.
- Kết nối Regularization với MAP ([Bài 10](./10_statistics_inference.md)).

## 1. Bias-Variance Tradeoff — phân tích lỗi model

Sai số tổng quát hóa (generalization error) phân tích thành:

$$\text{Expected Error} = \underbrace{\text{Bias}^2}_{\text{lỗi do model quá đơn giản}} + \underbrace{\text{Variance}}_{\text{lỗi do model quá nhạy với dữ liệu train}} + \underbrace{\text{Irreducible Error}}_{\text{nhiễu vốn có}}$$

- **High Bias (Underfitting)**: model quá đơn giản (vd Linear Regression cho quan hệ phi tuyến phức tạp) — sai số cao dù dùng dữ liệu nào.
- **High Variance (Overfitting)**: model quá phức tạp (vd Polynomial bậc cao — [Bài 16 mục 4](./16_model_evaluation.md)) — dự đoán thay đổi mạnh nếu đổi tập train khác.

```python
import numpy as np
import matplotlib.pyplot as plt

def bias_variance_demo(true_func, n_experiments=100, degree=1, n_samples=20, noise=1.0):
	x_test = np.linspace(0, 10, 50)
	predictions = []

	for _ in range(n_experiments):
		x_train = np.random.uniform(0, 10, n_samples)
		y_train = true_func(x_train) + np.random.randn(n_samples) * noise

		coeffs = np.polyfit(x_train, y_train, degree)
		y_pred = np.polyval(coeffs, x_test)
		predictions.append(y_pred)

	predictions = np.array(predictions)
	mean_pred = predictions.mean(axis=0)
	true_values = true_func(x_test)

	bias_squared = np.mean((mean_pred - true_values)**2)
	variance = np.mean(predictions.var(axis=0))
	return bias_squared, variance

true_func = lambda x: 2*x + 3
for degree in [1, 3, 10]:
	bias_sq, var = bias_variance_demo(true_func, degree=degree)
	print(f"Degree={degree}: Bias²={bias_sq:.3f}, Variance={var:.3f}")
	# degree=1 (đúng bậc): bias thấp, variance thấp
	# degree=10 (quá phức tạp): bias thấp NHƯNG variance CAO — overfitting
```

## 2. Regularization = thêm "phạt" vào hàm loss để giảm variance

$$L_{regularized}(\vec{w}) = L_{original}(\vec{w}) + \lambda \cdot R(\vec{w})$$

$\lambda$ (regularization strength) điều khiển mức độ "phạt" — $\lambda$ lớn → model đơn giản hơn (bias tăng, variance giảm); $\lambda = 0$ → không regularize (như model gốc).

## 3. Ridge Regression (L2 Regularization)

$$L_{Ridge}(\vec{w}) = \frac{1}{m}\sum_{i=1}^m(y_i-\hat{y_i})^2 + \lambda\|\vec{w}\|_2^2$$

Như đã chứng minh ở [Bài 10 mục 4](./10_statistics_inference.md), đây chính là **MAP với prior Gaussian** lên trọng số $\vec{w} \sim N(0, \tau^2)$.

```python
from sklearn.linear_model import Ridge

ridge = Ridge(alpha=1.0)  # alpha = lambda trong công thức
ridge.fit(X_train, y_train)
print(ridge.coef_)  # trọng số sẽ nhỏ hơn Linear Regression thường, không có cái nào = 0 hẳn
```

Gradient của Ridge (cộng thêm vào gradient MSE gốc — [Bài 14 mục 4](./14_linear_regression.md)):

$$\nabla_{\vec{w}} L_{Ridge} = \frac{2}{m}X^T(X\vec{w}-\vec{y}) + 2\lambda\vec{w}$$

## 4. Lasso Regression (L1 Regularization)

$$L_{Lasso}(\vec{w}) = \frac{1}{m}\sum_{i=1}^m(y_i-\hat{y_i})^2 + \lambda\|\vec{w}\|_1$$

```python
from sklearn.linear_model import Lasso

lasso = Lasso(alpha=0.1)
lasso.fit(X_train, y_train)
print(lasso.coef_)  # MỘT SỐ trọng số sẽ = 0 CHÍNH XÁC — Lasso tự động "chọn đặc trưng"
```

**Khác biệt cốt lõi giữa L1 và L2:** L1 norm có "góc nhọn" tại 0 (không khả vi tại đó — hình học của $\|\vec{w}\|_1$ là hình thoi trong không gian tham số), khiến nghiệm tối ưu thường rơi đúng vào các trục tọa độ → nhiều trọng số bị đẩy về **CHÍNH XÁC 0**. L2 norm (hình tròn/elip) không có góc nhọn, nghiệm chỉ bị "co nhỏ" chứ hiếm khi về đúng 0.

$\Rightarrow$ **Lasso hữu ích cho Feature Selection tự động** (loại bỏ đặc trưng không quan trọng), **Ridge hữu ích khi muốn giữ lại mọi đặc trưng nhưng giảm ảnh hưởng của đặc trưng nhiễu**.

## 5. Elastic Net — kết hợp cả 2

$$L_{ElasticNet} = L_{original} + \lambda_1\|\vec{w}\|_1 + \lambda_2\|\vec{w}\|_2^2$$

```python
from sklearn.linear_model import ElasticNet

elastic = ElasticNet(alpha=0.1, l1_ratio=0.5)  # l1_ratio: tỷ lệ L1 so với L2
elastic.fit(X_train, y_train)
```

## 6. Chọn $\lambda$ tối ưu bằng Cross-Validation ([Bài 16 mục 3](./16_model_evaluation.md))

```python
from sklearn.linear_model import RidgeCV

alphas = [0.001, 0.01, 0.1, 1, 10, 100]
ridge_cv = RidgeCV(alphas=alphas, cv=5)
ridge_cv.fit(X_train, y_train)
print("Best alpha:", ridge_cv.alpha_)
```

## 7. Trực quan hóa ảnh hưởng của $\lambda$

```python
import numpy as np
import matplotlib.pyplot as plt

alphas = np.logspace(-3, 3, 50)
coefs = []
for alpha in alphas:
	ridge = Ridge(alpha=alpha).fit(X_train, y_train)
	coefs.append(ridge.coef_)

plt.plot(alphas, coefs)
plt.xscale("log")
plt.xlabel("alpha (lambda)"); plt.ylabel("Trọng số")
plt.title("Ridge: trọng số co dần về 0 khi alpha tăng")
plt.show()
```

## Ví dụ đầy đủ: so sánh cả 3 trên dataset dễ overfitting

```python
import numpy as np
from sklearn.linear_model import LinearRegression, Ridge, Lasso
from sklearn.preprocessing import PolynomialFeatures
from sklearn.metrics import mean_squared_error

np.random.seed(42)
X = np.random.uniform(0, 10, (30, 1))
y = 2*X.flatten() + 3 + np.random.randn(30) * 3

poly = PolynomialFeatures(degree=10)
X_poly = poly.fit_transform(X)

models = {
	"Linear (không regularize)": LinearRegression(),
	"Ridge (L2)": Ridge(alpha=10),
	"Lasso (L1)": Lasso(alpha=0.5),
}

for name, model in models.items():
	model.fit(X_poly, y)
	mse = mean_squared_error(y, model.predict(X_poly))
	n_nonzero = np.sum(np.abs(model.coef_) > 1e-4) if hasattr(model, 'coef_') else "-"
	print(f"{name}: train_mse={mse:.3f}, số trọng số khác 0={n_nonzero}")
```

## Bài tập

1. **Bias-Variance demo**: dùng code mẫu mục 1, thử với `degree` = 1, 3, 10, quan sát bias/variance thay đổi thế nào, vẽ đồ thị minh họa.
2. **Ridge vs Lasso**: huấn luyện cả 2 trên cùng dataset có nhiều đặc trưng nhiễu (không liên quan tới target), so sánh số trọng số bị đưa về 0.
3. **Chọn alpha bằng CV**: dùng `RidgeCV`/`LassoCV`, tìm alpha tối ưu, giải thích tại sao alpha quá lớn/quá nhỏ đều không tốt.
4. **Nâng cao**: viết lại Gradient Descent cho Ridge Regression from scratch (dựa trên `LinearRegressionGD` ở [Bài 14](./14_linear_regression.md), thêm term $2\lambda\vec{w}$ vào gradient), verify với `sklearn.linear_model.Ridge`.

## Tiếp theo
→ [Bài 18: Decision Tree, Random Forest & Ensemble Methods](./18_trees_ensembles.md)
