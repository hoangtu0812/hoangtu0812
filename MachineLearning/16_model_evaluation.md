# Bài 16: Đánh Giá Mô Hình (Model Evaluation)

## Mục tiêu
- Train/test/validation split, k-fold cross-validation.
- Confusion Matrix, Precision/Recall/F1, ROC-AUC.
- Overfitting vs Underfitting.

## 1. Vì sao không đánh giá model trên chính dữ liệu huấn luyện?

Model có thể "học thuộc lòng" dữ liệu huấn luyện (overfitting) mà không tổng quát hóa được cho dữ liệu mới — accuracy trên tập train cao KHÔNG đảm bảo model tốt trong thực tế. Cần dữ liệu **chưa từng thấy** để đánh giá khách quan.

## 2. Train/Test/Validation Split

```python
from sklearn.model_selection import train_test_split

X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Với bài toán cần tuning hyperparameter, tách thêm validation set
X_train, X_val, y_train, y_val = train_test_split(X_train, y_train, test_size=0.25, random_state=42)
# Kết quả: 60% train, 20% validation, 20% test
```

- **Train set**: dùng để huấn luyện model.
- **Validation set**: dùng để chọn hyperparameter (learning rate, độ sâu cây...), KHÔNG dùng để huấn luyện.
- **Test set**: chỉ dùng 1 LẦN DUY NHẤT ở cuối cùng để đánh giá khách quan — không được "nhìn" test set trong lúc tuning, nếu không kết quả sẽ bị lạc quan giả tạo (data leakage).

## 3. k-Fold Cross-Validation — dùng dữ liệu hiệu quả hơn

Chia dữ liệu thành $k$ phần bằng nhau, huấn luyện $k$ lần, mỗi lần dùng 1 phần làm validation, $k-1$ phần còn lại làm train. Kết quả cuối = trung bình $k$ lần đánh giá — liên hệ tổ hợp cách chia dữ liệu ([Bài 11 mục 1](./11_discrete_math_optimization.md)).

```python
from sklearn.model_selection import cross_val_score
from sklearn.linear_model import LogisticRegression

model = LogisticRegression()
scores = cross_val_score(model, X, y, cv=5, scoring="accuracy")
print(f"Accuracy mỗi fold: {scores}")
print(f"Trung bình: {scores.mean():.4f} ± {scores.std():.4f}")
```

`scores.std()` (độ lệch chuẩn — [Bài 9 mục 4](./9_probability_distributions.md)) cho biết độ ổn định của model — std lớn nghĩa là hiệu năng dao động mạnh giữa các fold, đáng lo ngại.

## 4. Overfitting vs Underfitting — liên hệ Bias-Variance ([Bài 17](./17_regularization.md))

```
Underfitting: model quá đơn giản, train_error CAO, test_error CAO — model chưa "học" đủ pattern.
Đúng mức:      train_error thấp, test_error THẤP và GẦN train_error.
Overfitting:   model quá phức tạp, train_error RẤT thấp, test_error CAO — model "học thuộc" nhiễu.
```

```python
from sklearn.preprocessing import PolynomialFeatures
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error

for degree in [1, 4, 15]:
	poly = PolynomialFeatures(degree=degree)
	X_train_poly = poly.fit_transform(X_train)
	X_test_poly = poly.transform(X_test)

	model = LinearRegression().fit(X_train_poly, y_train)
	train_mse = mean_squared_error(y_train, model.predict(X_train_poly))
	test_mse = mean_squared_error(y_test, model.predict(X_test_poly))
	print(f"Degree={degree}: train_mse={train_mse:.3f}, test_mse={test_mse:.3f}")
	# degree=1: cả 2 cao (underfitting) | degree vừa: cả 2 thấp | degree=15: train thấp, test CAO (overfitting)
```

## 5. Confusion Matrix — nền tảng của mọi metric phân loại

|  | Dự đoán: Positive | Dự đoán: Negative |
|---|---|---|
| **Thật: Positive** | True Positive (TP) | False Negative (FN) |
| **Thật: Negative** | False Positive (FP) | True Negative (TN) |

```python
from sklearn.metrics import confusion_matrix, ConfusionMatrixDisplay
import matplotlib.pyplot as plt

cm = confusion_matrix(y_test, y_pred)
ConfusionMatrixDisplay(cm).plot()
plt.show()
```

## 6. Precision, Recall, F1-score

$$\text{Precision} = \frac{TP}{TP+FP} \quad\text{(trong số dự đoán Positive, bao nhiêu % đúng)}$$
$$\text{Recall} = \frac{TP}{TP+FN} \quad\text{(trong số Positive thật, bắt được bao nhiêu %)}$$
$$F_1 = 2 \cdot \frac{\text{Precision} \cdot \text{Recall}}{\text{Precision}+\text{Recall}} \quad\text{(trung bình điều hòa, cân bằng cả 2)}$$

```python
from sklearn.metrics import precision_score, recall_score, f1_score, classification_report

print(classification_report(y_test, y_pred))
```

**Khi nào ưu tiên Precision, khi nào ưu tiên Recall?**
- **Ưu tiên Recall** khi bỏ sót Positive rất tốn kém (vd chẩn đoán ung thư — thà báo nhầm còn hơn bỏ sót ca bệnh thật).
- **Ưu tiên Precision** khi báo động giả rất tốn kém (vd hệ thống chặn giao dịch gian lận — chặn nhầm giao dịch hợp lệ gây phiền cho khách hàng).

**Vì sao không chỉ dùng Accuracy?** Với dataset mất cân bằng (vd 99% negative, 1% positive), model dự đoán TOÀN BỘ là "negative" vẫn đạt accuracy 99% nhưng VÔ DỤNG — Precision/Recall/F1 phản ánh đúng chất lượng hơn nhiều trong trường hợp này.

## 7. ROC Curve & AUC

ROC (Receiver Operating Characteristic) vẽ **True Positive Rate** (= Recall) theo **False Positive Rate** ($FPR = \frac{FP}{FP+TN}$) khi thay đổi threshold phân loại từ 0 đến 1.

```python
from sklearn.metrics import roc_curve, roc_auc_score

y_proba = model.predict_proba(X_test)[:, 1]
fpr, tpr, thresholds = roc_curve(y_test, y_proba)
auc = roc_auc_score(y_test, y_proba)

plt.plot(fpr, tpr, label=f"AUC = {auc:.3f}")
plt.plot([0,1], [0,1], 'k--', label="Random guess")
plt.xlabel("False Positive Rate"); plt.ylabel("True Positive Rate")
plt.legend()
plt.show()
```

AUC (Area Under Curve) $\in [0,1]$: AUC=1 là model hoàn hảo, AUC=0.5 tương đương đoán ngẫu nhiên. AUC không phụ thuộc threshold cụ thể — hữu ích khi so sánh model tổng quát, không chỉ tại 1 ngưỡng 0.5.

## 8. Metric cho Regression (nối tiếp [Bài 14](./14_linear_regression.md))

```python
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score

mse = mean_squared_error(y_test, y_pred)
mae = mean_absolute_error(y_test, y_pred)  # ít nhạy với outlier hơn MSE
r2 = r2_score(y_test, y_pred)               # tỷ lệ phương sai giải thích được, [0,1] càng gần 1 càng tốt
```

$R^2 = 1 - \frac{\sum(y_i-\hat{y_i})^2}{\sum(y_i-\bar{y})^2}$ — so sánh model với baseline "chỉ đoán trung bình" — liên hệ khái niệm phương sai ([Bài 9](./9_probability_distributions.md)).

## Bài tập

1. **Train/test split & cross-validation**: huấn luyện Logistic Regression ([Bài 15](./15_logistic_regression.md)) trên dataset của bạn, so sánh accuracy từ 1 lần split vs trung bình 5-fold cross-validation.
2. **Overfitting demo**: dùng code mẫu mục 4, thử với nhiều `degree` khác nhau, vẽ đồ thị train_mse/test_mse theo degree — xác định "điểm ngọt" (sweet spot).
3. **Confusion Matrix & metrics**: tính tay Precision/Recall/F1 từ 1 confusion matrix cụ thể (tự cho số TP/FP/FN/TN), verify bằng `classification_report`.
4. **ROC-AUC trên dataset mất cân bằng**: tạo dataset mất cân bằng (`make_classification(weights=[0.95, 0.05])`), so sánh Accuracy vs F1 vs AUC — giải thích tại sao Accuracy gây hiểu lầm trong trường hợp này.

## Tiếp theo
→ [Bài 17: Regularization & Bias-Variance Tradeoff](./17_regularization.md)
