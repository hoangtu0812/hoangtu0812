# Dự Án Capstone: Pipeline Machine Learning Hoàn Chỉnh

## Mục tiêu
Ghép toàn bộ kiến thức từ Bài 1 → Bài 24 thành 1 dự án hoàn chỉnh: từ EDA, tiền xử lý, huấn luyện NHIỀU model (bao gồm cả cài from-scratch để chứng minh hiểu toán), đánh giá, tuning, tới đóng gói kết quả — áp dụng toán Phần I xuyên suốt, đúng tinh thần "hiểu từ gốc rễ" của cả lộ trình.

**Đề bài:** Dự đoán khả năng mắc bệnh tim (Heart Disease Prediction) — bài toán **phân loại nhị phân** kinh điển, dataset công khai (UCI Heart Disease / Kaggle), đủ nhỏ để chạy nhanh nhưng đủ thực tế để áp dụng toàn bộ quy trình chuẩn.

## Cấu trúc thư mục

```
heart-disease-ml/
├── data/
│   └── heart.csv
├── notebooks/
│   ├── 01_eda.ipynb
│   ├── 02_preprocessing.ipynb
│   ├── 03_classical_models.ipynb
│   ├── 04_neural_network_from_scratch.ipynb
│   └── 05_neural_network_pytorch.ipynb
├── src/
│   ├── preprocessing.py
│   ├── models/
│   │   ├── logistic_regression_scratch.py   # from Bài 15
│   │   ├── neural_network_scratch.py         # from Bài 21
│   │   └── train_classical.py                  # scikit-learn models
│   └── evaluate.py
├── report/
│   └── final_report.md
├── requirements.txt
└── README.md
```

## Bước 1: EDA — áp dụng [Bài 13](./13_data_visualization_eda.md) đầy đủ

```python
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

df = pd.read_csv("data/heart.csv")

print(df.shape, df.info(), df.describe())
print(df.isnull().sum())
print(df['target'].value_counts(normalize=True))  # kiểm tra mất cân bằng lớp — liên hệ Bài 16 mục 6

sns.heatmap(df.corr(), annot=True, cmap="coolwarm")
plt.show()

sns.pairplot(df, hue="target", vars=["age", "chol", "thalach", "oldpeak"])
plt.show()
```

**Câu hỏi cần trả lời qua EDA** (liên hệ trực tiếp Phần I): đặc trưng nào tương quan mạnh với `target` (dùng ma trận hiệp phương sai — [Bài 9 mục 4](./9_probability_distributions.md))? Có multicollinearity giữa các đặc trưng không ([Bài 3 mục 7](./3_linear_algebra_matrices.md))? Phân phối mỗi đặc trưng có gần Normal không, có outlier không ([Bài 9 mục 6](./9_probability_distributions.md), [Bài 13 mục 5](./13_data_visualization_eda.md))?

## Bước 2: Tiền xử lý — áp dụng [Bài 12](./12_numpy_pandas.md)

```python
# src/preprocessing.py
import numpy as np
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler

def load_and_preprocess(path: str, test_size=0.2, random_state=42):
	df = pd.read_csv(path)
	df = df.drop_duplicates().dropna()

	X = df.drop(columns=["target"]).to_numpy()
	y = df["target"].to_numpy()

	X_train, X_test, y_train, y_test = train_test_split(
		X, y, test_size=test_size, random_state=random_state, stratify=y  # stratify: giữ tỷ lệ lớp — Bài 16
	)

	scaler = StandardScaler()  # chuẩn hóa — Bài 9 mục 3, bắt buộc cho Gradient Descent hội tụ tốt
	X_train = scaler.fit_transform(X_train)
	X_test = scaler.transform(X_test)  # CHỈ fit trên train, tránh data leakage

	return X_train, X_test, y_train, y_test, scaler
```

## Bước 3: Model cổ điển — áp dụng [Bài 14-19](./14_linear_regression.md)

```python
# src/models/train_classical.py
from sklearn.linear_model import LogisticRegression
from sklearn.ensemble import RandomForestClassifier, GradientBoostingClassifier
from sklearn.svm import SVC
from sklearn.model_selection import cross_val_score, GridSearchCV

def train_and_compare(X_train, y_train, X_test, y_test):
	models = {
		"Logistic Regression": LogisticRegression(),
		"Random Forest": RandomForestClassifier(n_estimators=100, random_state=42),
		"Gradient Boosting": GradientBoostingClassifier(random_state=42),
		"SVM (RBF)": SVC(kernel="rbf", probability=True),
	}

	results = {}
	for name, model in models.items():
		cv_scores = cross_val_score(model, X_train, y_train, cv=5, scoring="f1")  # Bài 16
		model.fit(X_train, y_train)
		test_acc = model.score(X_test, y_test)
		results[name] = {"cv_f1_mean": cv_scores.mean(), "cv_f1_std": cv_scores.std(), "test_acc": test_acc}
		print(f"{name}: CV F1={cv_scores.mean():.3f}±{cv_scores.std():.3f}, Test Acc={test_acc:.3f}")

	return results

def tune_random_forest(X_train, y_train):
	"""Tuning hyperparameter bằng GridSearchCV — Bài 16, Bài 17."""
	param_grid = {"n_estimators": [50, 100, 200], "max_depth": [3, 5, 10, None]}
	grid = GridSearchCV(RandomForestClassifier(random_state=42), param_grid, cv=5, scoring="f1")
	grid.fit(X_train, y_train)
	print("Best params:", grid.best_params_)
	return grid.best_estimator_
```

## Bước 4: Logistic Regression from scratch — verify hiểu [Bài 15](./15_logistic_regression.md)

```python
# src/models/logistic_regression_scratch.py
import numpy as np

class LogisticRegressionScratch:
	def __init__(self, lr=0.1, n_iters=1000, l2_lambda=0.01):
		self.lr, self.n_iters, self.l2_lambda = lr, n_iters, l2_lambda

	def fit(self, X, y):
		m, n = X.shape
		X_b = np.hstack([np.ones((m,1)), X])
		self.w = np.zeros(n+1)

		for _ in range(self.n_iters):
			z = X_b @ self.w
			y_pred = 1 / (1 + np.exp(-z))
			grad = (1/m) * X_b.T @ (y_pred - y) + 2*self.l2_lambda*self.w  # + Ridge — Bài 17
			self.w -= self.lr * grad
		return self

	def predict(self, X):
		X_b = np.hstack([np.ones((X.shape[0],1)), X])
		return (1/(1+np.exp(-(X_b @ self.w))) >= 0.5).astype(int)
```

**Kiểm tra bắt buộc:** so sánh accuracy của `LogisticRegressionScratch` với `sklearn.linear_model.LogisticRegression` trên CÙNG dataset — phải xấp xỉ nhau (chênh lệch nhỏ do khác cách tối ưu/regularization mặc định). Nếu chênh lệch lớn, đó là dấu hiệu implementation from-scratch có bug — quay lại verify từng công thức ở [Bài 15](./15_logistic_regression.md).

## Bước 5: Neural Network — from scratch VÀ PyTorch — verify [Bài 20-22](./20_neural_networks.md)

```python
# So sánh: NeuralNetworkFromScratch (Bài 21) vs PyTorch
from src.models.neural_network_scratch import NeuralNetworkFromScratch
import torch
import torch.nn as nn

# From scratch
nn_scratch = NeuralNetworkFromScratch([X_train.shape[1], 16, 8, 2], lr=0.05)
for epoch in range(200):
	total_loss = 0
	for i in range(len(X_train)):
		y_onehot = np.eye(2)[y_train[i]]
		total_loss += nn_scratch.train_step(X_train[i], y_onehot)
	if epoch % 50 == 0:
		print(f"Epoch {epoch}: loss={total_loss/len(X_train):.4f}")

# PyTorch — kiến trúc TƯƠNG ĐƯƠNG, để so sánh tốc độ hội tụ & kết quả
class TorchMLP(nn.Module):
	def __init__(self, input_dim):
		super().__init__()
		self.net = nn.Sequential(
			nn.Linear(input_dim, 16), nn.ReLU(),
			nn.Linear(16, 8), nn.ReLU(),
			nn.Linear(8, 2),
		)
	def forward(self, x):
		return self.net(x)

model = TorchMLP(X_train.shape[1])
optimizer = torch.optim.Adam(model.parameters(), lr=0.01)  # Bài 22
criterion = nn.CrossEntropyLoss()

X_train_t = torch.FloatTensor(X_train)
y_train_t = torch.LongTensor(y_train)

for epoch in range(200):
	optimizer.zero_grad()
	outputs = model(X_train_t)
	loss = criterion(outputs, y_train_t)
	loss.backward()   # PyTorch tự động làm chính xác những gì Bài 21 cài thủ công
	optimizer.step()
	if epoch % 50 == 0:
		print(f"Epoch {epoch}: loss={loss.item():.4f}")
```

**Ý nghĩa của bước so sánh này:** verify rằng những gì bạn cài THỦ CÔNG ở [Bài 21](./21_backpropagation.md) (Backpropagation) cho kết quả tương đương những gì PyTorch làm TỰ ĐỘNG (autograd) — đây là bằng chứng mạnh mẽ nhất rằng bạn thực sự HIỂU cơ chế bên trong, không chỉ biết "gọi API".

## Bước 6: Đánh giá tổng hợp — áp dụng [Bài 16](./16_model_evaluation.md) đầy đủ

```python
# src/evaluate.py
from sklearn.metrics import classification_report, roc_auc_score, roc_curve
import matplotlib.pyplot as plt

def full_evaluation_report(models: dict, X_test, y_test):
	plt.figure(figsize=(8,6))
	for name, model in models.items():
		y_pred = model.predict(X_test)
		y_proba = model.predict_proba(X_test)[:,1] if hasattr(model, "predict_proba") else y_pred

		print(f"\n=== {name} ===")
		print(classification_report(y_test, y_pred))

		fpr, tpr, _ = roc_curve(y_test, y_proba)
		auc = roc_auc_score(y_test, y_proba)
		plt.plot(fpr, tpr, label=f"{name} (AUC={auc:.3f})")

	plt.plot([0,1],[0,1],'k--')
	plt.xlabel("FPR"); plt.ylabel("TPR"); plt.legend()
	plt.title("So sánh ROC Curve giữa các model")
	plt.savefig("report/roc_comparison.png")
	plt.show()
```

## Bước 7: Đóng gói báo cáo

`report/final_report.md` nên có:
1. **Tóm tắt bài toán & dataset** (số mẫu, số đặc trưng, tỷ lệ lớp).
2. **Phát hiện chính từ EDA** (đặc trưng quan trọng nhất, tương quan, outlier).
3. **Bảng so sánh MỌI model** (accuracy, precision, recall, F1, AUC — [Bài 16](./16_model_evaluation.md)) kèm ROC curve.
4. **Phân tích bias-variance** cho model tốt nhất ([Bài 17](./17_regularization.md)) — model có đang overfit/underfit không, đã tuning $\lambda$/`max_depth`/dropout thế nào.
5. **So sánh Neural Network from-scratch vs PyTorch** — verify kết quả khớp nhau, chứng minh hiểu đúng Backpropagation.
6. **Kết luận & hướng cải thiện** (thêm dữ liệu, feature engineering, thử kiến trúc DL phức tạp hơn nếu phù hợp bài toán).

## Checklist hoàn thành

- [ ] EDA đầy đủ: phân phối, correlation, outlier, kiểm tra mất cân bằng lớp.
- [ ] Tiền xử lý đúng thứ tự: split TRƯỚC, fit scaler CHỈ trên train (tránh data leakage).
- [ ] Huấn luyện tối thiểu 4 model cổ điển (Logistic Regression, Random Forest, Gradient Boosting, SVM), so sánh bằng cross-validation.
- [ ] Cài đặt Logistic Regression from scratch, verify khớp scikit-learn.
- [ ] Cài đặt Neural Network from scratch ([Bài 21](./21_backpropagation.md)), verify gradient checking, verify khớp PyTorch.
- [ ] Tuning hyperparameter cho model tốt nhất bằng GridSearchCV/cross-validation.
- [ ] Báo cáo đầy đủ metric (không chỉ Accuracy — đặc biệt quan trọng nếu dataset mất cân bằng).
- [ ] `final_report.md` viết rõ ràng, có thể đưa vào portfolio/CV.

## Mở rộng (tùy chọn, sau khi hoàn thành checklist)
- Thử SHAP/permutation importance để giải thích model (interpretability).
- Deploy model tốt nhất thành API nhỏ bằng FastAPI ([Python Bài 18](../Python/18_rest_api.md)) — kết nối trực tiếp track Python đã học.
- Thử kiến trúc CNN ([Bài 23](./23_cnn.md)) trên 1 bài toán ảnh khác (MNIST/CIFAR-10) để có kinh nghiệm với dữ liệu ảnh thật.
- Viết unit test cho pipeline tiền xử lý + model bằng `pytest` ([Python Bài 14](../Python/14_testing.md)).

---
Hoàn thành bài này nghĩa là bạn đã đi trọn vẹn từ **vector/ma trận cơ bản** ([Bài 2](./2_linear_algebra_vectors.md)) tới **huấn luyện và đánh giá 1 hệ thống Machine Learning hoàn chỉnh**, hiểu rõ TOÁN đứng sau mỗi công thức thay vì chỉ gọi API — nền tảng vững chắc để tiếp tục học sâu hơn về bất kỳ hướng chuyên biệt nào của AI/ML (Computer Vision, NLP, Reinforcement Learning...) trong tương lai.
