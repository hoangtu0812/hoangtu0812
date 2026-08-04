# Bài 13: Trực quan hóa dữ liệu & EDA (Exploratory Data Analysis)

## Mục tiêu
- Matplotlib/Seaborn cho các loại biểu đồ cơ bản.
- Phân tích thống kê mô tả (nối tiếp [Bài 9](./9_probability_distributions.md)).
- Phát hiện outlier, phân tích correlation.

## 1. Vì sao EDA là bước KHÔNG THỂ bỏ qua

Trước khi chọn/huấn luyện bất kỳ model nào, phải hiểu dữ liệu: phân phối ra sao, có outlier không, các đặc trưng tương quan thế nào. Bỏ qua EDA thường dẫn tới chọn sai model, bỏ sót lỗi dữ liệu, hoặc hiểu sai kết quả sau này.

## 2. Matplotlib — biểu đồ cơ bản

```python
import matplotlib.pyplot as plt
import numpy as np

x = np.linspace(0, 10, 100)
y = np.sin(x)

plt.figure(figsize=(8, 4))
plt.plot(x, y, label="sin(x)")
plt.xlabel("x")
plt.ylabel("y")
plt.title("Đồ thị sin")
plt.legend()
plt.grid(True)
plt.show()
```

### Histogram — quan sát phân phối (liên hệ [Bài 9](./9_probability_distributions.md))

```python
data = np.random.normal(50, 10, 1000)
plt.hist(data, bins=30, edgecolor="black")
plt.title("Phân phối dữ liệu")
plt.show()
```

### Scatter plot — quan hệ giữa 2 biến

```python
x = np.random.rand(100) * 10
y = 2 * x + np.random.randn(100) * 2  # quan hệ tuyến tính + nhiễu
plt.scatter(x, y, alpha=0.6)
plt.xlabel("x"); plt.ylabel("y")
plt.title("Scatter plot — quan hệ tuyến tính")
plt.show()
```

### Boxplot — phát hiện outlier

```python
data_with_outliers = np.concatenate([np.random.normal(50, 5, 100), [100, 105, -10]])
plt.boxplot(data_with_outliers)
plt.title("Boxplot — điểm ngoài whisker là outlier")
plt.show()
```

Boxplot hiển thị: đường giữa hộp = median, mép hộp = quartile 1/3 (Q1/Q3), "whisker" mở rộng tới $Q1 - 1.5 \times IQR$ và $Q3 + 1.5 \times IQR$ ($IQR = Q3-Q1$) — điểm nằm ngoài khoảng này được coi là **outlier** theo quy ước thống kê thông dụng.

## 3. Seaborn — biểu đồ thống kê nâng cao, dựng trên Matplotlib

```python
import seaborn as sns
import pandas as pd

df = sns.load_dataset("iris")  # dataset mẫu có sẵn trong seaborn

# Pairplot — ma trận scatter plot cho MỌI cặp biến, tô màu theo nhãn
sns.pairplot(df, hue="species")
plt.show()

# Heatmap — trực quan hóa ma trận tương quan
correlation_matrix = df.select_dtypes(include=[np.number]).corr()
sns.heatmap(correlation_matrix, annot=True, cmap="coolwarm")
plt.title("Ma trận tương quan")
plt.show()
```

Ma trận tương quan chính là chuẩn hóa của ma trận hiệp phương sai ([Bài 9 mục 5](./9_probability_distributions.md)) — heatmap giúp nhìn nhanh cặp đặc trưng nào tương quan mạnh (có thể gây multicollinearity — [Bài 3 mục 7](./3_linear_algebra_matrices.md)).

### Boxplot theo nhóm — so sánh phân phối giữa các category

```python
sns.boxplot(x="species", y="petal_length", data=df)
plt.title("Phân phối petal_length theo loài")
plt.show()
```

## 4. Thống kê mô tả — nối trực tiếp lý thuyết Bài 9

```python
print(df.describe())         # count, mean, std, min, 25%, 50%, 75%, max cho mỗi cột số
print(df["petal_length"].skew())   # độ lệch (skewness)
print(df["petal_length"].kurt())    # độ nhọn (kurtosis)
```

## 5. Phát hiện outlier bằng phương pháp thống kê (bổ sung cho boxplot)

### Phương pháp IQR

```python
def detect_outliers_iqr(series):
	Q1, Q3 = series.quantile(0.25), series.quantile(0.75)
	IQR = Q3 - Q1
	lower, upper = Q1 - 1.5*IQR, Q3 + 1.5*IQR
	return series[(series < lower) | (series > upper)]

outliers = detect_outliers_iqr(df["sepal_width"])
print(outliers)
```

### Phương pháp Z-score (dựa trên độ lệch chuẩn — liên hệ Normal distribution Bài 9)

```python
def detect_outliers_zscore(series, threshold=3):
	z_scores = (series - series.mean()) / series.std()
	return series[abs(z_scores) > threshold]
```

Z-score đo "1 điểm dữ liệu cách bao nhiêu độ lệch chuẩn so với trung bình" — với phân phối chuẩn, $|z| > 3$ chỉ xảy ra với xác suất < 0.3%, nên thường dùng làm ngưỡng phát hiện outlier.

## 6. Correlation — cẩn thận: correlation KHÔNG phải nhân quả

```python
print(df.select_dtypes(include=[np.number]).corr())
```

Hệ số tương quan cao giữa 2 biến ($|\rho|$ gần 1) chỉ cho biết có quan hệ **tuyến tính**, KHÔNG chứng minh biến này **gây ra** biến kia — luôn nhắc bản thân điều này khi diễn giải kết quả EDA, tránh kết luận sai lầm trong phân tích thực tế.

## Ví dụ đầy đủ: quy trình EDA hoàn chỉnh

```python
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns

def full_eda_report(df: pd.DataFrame, target_col: str):
	print("=== Kích thước & kiểu dữ liệu ===")
	print(df.shape)
	print(df.dtypes)

	print("\n=== Giá trị thiếu ===")
	print(df.isnull().sum())

	print("\n=== Thống kê mô tả ===")
	print(df.describe())

	print("\n=== Ma trận tương quan ===")
	numeric_df = df.select_dtypes(include=[np.number])
	corr = numeric_df.corr()
	sns.heatmap(corr, annot=True, cmap="coolwarm")
	plt.title("Correlation Matrix")
	plt.show()

	print(f"\n=== Tương quan với target '{target_col}' ===")
	print(corr[target_col].sort_values(ascending=False))

if __name__ == "__main__":
	df = sns.load_dataset("iris")
	df["species_code"] = df["species"].astype("category").cat.codes
	full_eda_report(df, "species_code")
```

## Bài tập

1. **Histogram & Boxplot**: load dataset Iris (hoặc dataset khác), vẽ histogram và boxplot cho từng cột số, nhận xét về phân phối và outlier.
2. **Pairplot & Heatmap**: dùng `sns.pairplot` và `sns.heatmap` cho dataset đó, xác định cặp đặc trưng nào tương quan mạnh nhất.
3. **Phát hiện outlier**: viết cả 2 hàm `detect_outliers_iqr` và `detect_outliers_zscore`, so sánh kết quả trên cùng 1 cột dữ liệu — chúng có luôn cho kết quả giống nhau không?
4. **EDA report hoàn chỉnh**: dùng code mẫu trên làm nền, tự viết lại `full_eda_report`, chạy trên 1 dataset thật (từ Kaggle hoặc `sklearn.datasets`) trước khi bước vào [Bài 14](./14_linear_regression.md).

## Tiếp theo
→ [Bài 14: Linear Regression — Từ toán tới code](./14_linear_regression.md)
