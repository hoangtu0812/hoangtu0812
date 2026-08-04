# Bài 18: Decision Tree, Random Forest & Ensemble Methods

## Mục tiêu
- Entropy/Gini Index — tiêu chí chia nhánh.
- Decision Tree, Bagging, Random Forest, Boosting.

## 1. Decision Tree — trực giác

Cây quyết định chia dữ liệu bằng chuỗi câu hỏi if/else dựa trên đặc trưng, mỗi node lá là 1 dự đoán — đây là dạng đặc biệt của đồ thị (cây) đã nhắc ở [Bài 11 mục 2](./11_discrete_math_optimization.md).

```
Diện tích > 100m²?
├── Có → Số phòng > 3?
│         ├── Có → Giá cao
│         └── Không → Giá trung bình
└── Không → Giá thấp
```

## 2. Entropy — đo độ "hỗn loạn" của tập dữ liệu

$$H(S) = -\sum_{i=1}^k p_i \log_2(p_i)$$

với $p_i$ = tỷ lệ mẫu thuộc lớp $i$ trong tập $S$ — đây là **entropy thông tin (Shannon entropy)**, liên hệ trực tiếp phân phối xác suất rời rạc ([Bài 9](./9_probability_distributions.md)).

```python
import numpy as np

def entropy(y):
	_, counts = np.unique(y, return_counts=True)
	probabilities = counts / len(y)
	return -np.sum(probabilities * np.log2(probabilities + 1e-15))

print(entropy(np.array([0,0,0,0])))        # 0 — hoàn toàn thuần khiết, không "hỗn loạn"
print(entropy(np.array([0,0,1,1])))        # 1 — hỗn loạn tối đa (50/50)
print(entropy(np.array([0,0,0,1])))        # ~0.81 — trung gian
```

Entropy = 0 khi tất cả mẫu cùng 1 lớp (thuần khiết); entropy cực đại khi các lớp phân bố đều nhau (bất định lớn nhất).

## 3. Gini Index — lựa chọn thay thế cho Entropy

$$Gini(S) = 1 - \sum_{i=1}^k p_i^2$$

```python
def gini(y):
	_, counts = np.unique(y, return_counts=True)
	probabilities = counts / len(y)
	return 1 - np.sum(probabilities**2)
```

Gini tính toán nhanh hơn Entropy (không cần $\log$), thường cho kết quả cây tương tự — scikit-learn dùng Gini làm mặc định.

## 4. Information Gain — tiêu chí chọn đặc trưng để chia nhánh

$$IG(S, A) = H(S) - \sum_{v \in \text{values}(A)} \frac{|S_v|}{|S|}H(S_v)$$

Thuật toán xây cây (ID3/CART) chọn đặc trưng $A$ có **Information Gain cao nhất** tại mỗi bước — giảm entropy nhiều nhất sau khi chia.

```python
def information_gain(y, y_left, y_right):
	n = len(y)
	weighted_entropy = (len(y_left)/n)*entropy(y_left) + (len(y_right)/n)*entropy(y_right)
	return entropy(y) - weighted_entropy
```

## 5. Decision Tree với scikit-learn

```python
from sklearn.tree import DecisionTreeClassifier, plot_tree
import matplotlib.pyplot as plt

tree = DecisionTreeClassifier(max_depth=3, criterion="gini", random_state=42)
tree.fit(X_train, y_train)

plt.figure(figsize=(12, 8))
plot_tree(tree, feature_names=feature_names, filled=True)
plt.show()

print("Feature importance:", tree.feature_importances_)
```

`max_depth` là hyperparameter điều khiển độ phức tạp — cây quá sâu dễ **overfitting** (liên hệ [Bài 16-17](./16_model_evaluation.md)): mỗi lá chỉ chứa 1-2 mẫu, "học thuộc" dữ liệu train thay vì học pattern tổng quát.

## 6. Bagging (Bootstrap Aggregating) — giảm variance bằng cách kết hợp nhiều model

Ý tưởng: lấy mẫu **có hoàn lại (bootstrap)** từ dataset gốc $B$ lần, huấn luyện $B$ model độc lập, kết quả cuối = trung bình (regression) hoặc vote đa số (classification).

```python
from sklearn.ensemble import BaggingClassifier

bagging = BaggingClassifier(
	estimator=DecisionTreeClassifier(),
	n_estimators=50,
	random_state=42
)
bagging.fit(X_train, y_train)
```

**Vì sao Bagging giảm variance?** Trung bình của nhiều biến ngẫu nhiên (các model, mỗi cái huấn luyện trên 1 bootstrap sample khác nhau) có phương sai thấp hơn từng biến riêng lẻ — hệ quả trực tiếp của tính chất $\text{Var}(\bar{X}) = \text{Var}(X)/n$ khi các biến độc lập (liên hệ [Bài 9](./9_probability_distributions.md)).

## 7. Random Forest — Bagging + chọn ngẫu nhiên đặc trưng

```python
from sklearn.ensemble import RandomForestClassifier

rf = RandomForestClassifier(n_estimators=100, max_depth=5, random_state=42)
rf.fit(X_train, y_train)

print("Accuracy:", rf.score(X_test, y_test))
print("Feature importance:", rf.feature_importances_)
```

Ngoài bootstrap sample (Bagging), mỗi cây trong Random Forest còn chỉ xét **1 tập con ngẫu nhiên các đặc trưng** tại mỗi lần chia nhánh — giúp các cây "đa dạng hóa" hơn nữa (giảm tương quan giữa các cây), giảm variance mạnh hơn Bagging thuần.

## 8. Boosting — xây model TUẦN TỰ, mỗi model sửa lỗi của model trước

Khác Bagging (các model độc lập, song song), Boosting huấn luyện **tuần tự**: model sau tập trung vào các mẫu mà model trước dự đoán SAI.

```python
from sklearn.ensemble import GradientBoostingClassifier

gb = GradientBoostingClassifier(n_estimators=100, learning_rate=0.1, max_depth=3, random_state=42)
gb.fit(X_train, y_train)
```

Gradient Boosting áp dụng **Gradient Descent trong không gian hàm số** ([Bài 11](./11_discrete_math_optimization.md)) — mỗi cây mới được huấn luyện để xấp xỉ **gradient âm của hàm loss** so với dự đoán hiện tại, tương tự cách Gradient Descent cập nhật tham số ngược chiều gradient. XGBoost/LightGBM là các thư viện tối ưu hóa hiệu năng cao của ý tưởng này, rất phổ biến trong thi Kaggle và production thực tế.

## 9. So sánh Bagging vs Boosting

| | Bagging (Random Forest) | Boosting (Gradient Boosting) |
|---|---|---|
| Huấn luyện | Song song, độc lập | Tuần tự, model sau phụ thuộc model trước |
| Mục tiêu chính | Giảm **variance** | Giảm **bias** (và variance) |
| Dễ overfit? | Ít hơn (nhờ đa dạng hóa) | Dễ overfit hơn nếu quá nhiều estimator/learning rate cao |
| Tốc độ huấn luyện | Nhanh hơn (song song hóa được) | Chậm hơn (tuần tự) |

## Ví dụ đầy đủ

```python
from sklearn.datasets import load_breast_cancer
from sklearn.model_selection import train_test_split, cross_val_score
from sklearn.tree import DecisionTreeClassifier
from sklearn.ensemble import RandomForestClassifier, GradientBoostingClassifier

data = load_breast_cancer()
X_train, X_test, y_train, y_test = train_test_split(data.data, data.target, test_size=0.2, random_state=42)

models = {
	"Decision Tree": DecisionTreeClassifier(max_depth=5, random_state=42),
	"Random Forest": RandomForestClassifier(n_estimators=100, max_depth=5, random_state=42),
	"Gradient Boosting": GradientBoostingClassifier(n_estimators=100, random_state=42),
}

for name, model in models.items():
	scores = cross_val_score(model, X_train, y_train, cv=5)
	model.fit(X_train, y_train)
	test_acc = model.score(X_test, y_test)
	print(f"{name}: CV accuracy={scores.mean():.4f}, Test accuracy={test_acc:.4f}")
```

## Bài tập

1. **Entropy/Gini bằng tay**: tính tay entropy và Gini cho tập $\{0,0,0,1,1\}$, verify bằng code.
2. **Decision Tree & overfitting**: huấn luyện với `max_depth` = 2, 5, None (không giới hạn), so sánh train/test accuracy — quan sát overfitting khi cây quá sâu.
3. **Random Forest**: huấn luyện trên dataset thật (`load_breast_cancer` hoặc tương tự), in `feature_importances_`, xác định 3 đặc trưng quan trọng nhất.
4. **So sánh 3 model**: dùng code mẫu trên làm nền, tự viết lại, so sánh Decision Tree đơn lẻ vs Random Forest vs Gradient Boosting trên cùng dataset — nhận xét về accuracy và độ ổn định (std của cross-validation).

## Tiếp theo
→ [Bài 19: SVM, k-NN, Naive Bayes & Unsupervised Learning](./19_svm_knn_clustering.md)
