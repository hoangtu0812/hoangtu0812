# Bài 12: NumPy & Pandas cho Data Science

## Mục tiêu
- NumPy array, vector hóa (vectorization) — cài đặt thực tế cho toán ở Phần I.
- Pandas DataFrame, xử lý dữ liệu thiếu/trùng lặp.

## 1. NumPy array — cài đặt thực tế của vector/ma trận ([Bài 2-3](./2_linear_algebra_vectors.md))

```python
import numpy as np

arr1d = np.array([1, 2, 3])           # vector
arr2d = np.array([[1, 2], [3, 4]])     # ma trận
zeros = np.zeros((3, 4))                # ma trận 0, shape (3,4)
ones = np.ones((2, 2))
identity = np.eye(3)                     # ma trận đơn vị
random_arr = np.random.randn(3, 3)       # ma trận ngẫu nhiên từ N(0,1)

print(arr2d.shape, arr2d.dtype, arr2d.ndim)
```

## 2. Vector hóa (Vectorization) — vì sao NumPy nhanh hơn vòng lặp Python

```python
import time

n = 1_000_000
a = list(range(n))
b = list(range(n))

start = time.perf_counter()
c = [a[i] + b[i] for i in range(n)]  # vòng lặp Python thuần — CHẬM
print(f"Python loop: {time.perf_counter()-start:.4f}s")

a_np, b_np = np.array(a), np.array(b)
start = time.perf_counter()
c_np = a_np + b_np  # vector hóa — NHANH HƠN HÀNG CHỤC LẦN
print(f"NumPy vectorized: {time.perf_counter()-start:.4f}s")
```

NumPy thực thi phép toán bằng code C biên dịch sẵn, xử lý cả mảng cùng lúc thay vì lặp từng phần tử qua interpreter Python — đây là lý do TUYỆT ĐỐI không viết vòng lặp `for` thủ công khi làm việc với dữ liệu số lớn trong ML.

## 3. Indexing, Slicing, Broadcasting

```python
arr = np.array([[1, 2, 3], [4, 5, 6], [7, 8, 9]])

print(arr[0, :])       # dòng đầu: [1 2 3]
print(arr[:, 1])         # cột thứ 2: [2 5 8]
print(arr[arr > 5])       # boolean indexing: [6 7 8 9]
print(arr[0:2, 1:3])       # slicing 2D: [[2 3] [5 6]]

# Broadcasting — tự động "mở rộng" mảng nhỏ hơn để khớp phép toán
vec = np.array([10, 20, 30])
print(arr + vec)  # cộng vec vào MỖI DÒNG của arr, không cần vòng lặp
```

Broadcasting là cơ chế đứng sau rất nhiều đoạn code ML "gọn" — vd chuẩn hóa dữ liệu ($X - \mu$) cho cả dataset chỉ trong 1 dòng.

## 4. Các hàm thống kê NumPy — cài đặt của Bài 9

```python
data = np.array([[1, 2], [3, 4], [5, 6]])

print(data.mean(axis=0))   # mean theo CỘT: [3. 4.]
print(data.mean(axis=1))   # mean theo DÒNG: [1.5 3.5 5.5]
print(data.std(axis=0))     # độ lệch chuẩn theo cột
print(np.cov(data.T))        # ma trận hiệp phương sai (Bài 9 mục 5)
```

`axis=0` = "thu gọn theo chiều dòng" (tính cho mỗi cột), `axis=1` = "thu gọn theo chiều cột" (tính cho mỗi dòng) — nguồn nhầm lẫn phổ biến nhất khi mới học NumPy, luôn kiểm tra bằng ví dụ nhỏ nếu không chắc.

## 5. Pandas DataFrame — bảng dữ liệu có nhãn

```python
import pandas as pd

df = pd.DataFrame({
	"name": ["Ben", "Alice", "Carol"],
	"age": [25, 30, 28],
	"salary": [50000, 60000, 55000],
})

print(df.head())
print(df.dtypes)
print(df["age"])              # chọn 1 cột — trả về Series
print(df[["name", "age"]])      # chọn nhiều cột — trả về DataFrame
print(df[df["age"] > 26])        # lọc theo điều kiện
```

## 6. Đọc dữ liệu thực tế

```python
df = pd.read_csv("data.csv")
# df = pd.read_excel("data.xlsx")
# df = pd.read_json("data.json")

print(df.shape)      # (số dòng, số cột)
print(df.info())       # kiểu dữ liệu + số giá trị non-null mỗi cột
print(df.describe())    # thống kê mô tả: mean, std, min, max, quartile
```

## 7. Xử lý dữ liệu thiếu (Missing Data)

```python
df = pd.DataFrame({
	"a": [1, 2, None, 4],
	"b": [5, None, 7, 8],
})

print(df.isnull().sum())       # đếm giá trị thiếu mỗi cột

df_dropped = df.dropna()         # xóa dòng có giá trị thiếu
df_filled_mean = df.fillna(df.mean())  # điền bằng trung bình cột
df_filled_ffill = df.fillna(method="ffill")  # điền bằng giá trị liền trước (forward fill)
```

**Cách xử lý missing data ảnh hưởng trực tiếp tới kết quả model** — điền sai (vd điền 0 cho cột tuổi) có thể tạo ra pattern giả không tồn tại trong dữ liệu thật. Chi tiết chiến lược EDA đầy đủ ở [Bài 13](./13_data_visualization_eda.md).

## 8. Xử lý dữ liệu trùng lặp & Group By

```python
df = df.drop_duplicates()

# GroupBy — tương đương SQL GROUP BY (liên hệ Go/SQLAlchemy nếu đã quen từ Python Bài 17)
sales = pd.DataFrame({
	"region": ["North", "South", "North", "South"],
	"amount": [100, 150, 200, 120],
})
grouped = sales.groupby("region")["amount"].sum()
print(grouped)
```

## 9. Chuyển đổi giữa Pandas ↔ NumPy — kết nối trực tiếp tới ML pipeline

```python
X = df[["age", "salary"]].to_numpy()  # DataFrame -> NumPy array để đưa vào model
```

Hầu hết thư viện ML (scikit-learn, PyTorch) nhận đầu vào là NumPy array — Pandas dùng để khám phá/tiền xử lý dữ liệu, NumPy dùng để tính toán/huấn luyện model.

## Ví dụ đầy đủ

```python
import numpy as np
import pandas as pd

def preprocess_dataset(df: pd.DataFrame, feature_cols: list, target_col: str):
	df = df.drop_duplicates()
	df = df.fillna(df[feature_cols].mean())

	X = df[feature_cols].to_numpy()
	y = df[target_col].to_numpy()

	# Chuẩn hóa (standardization) — liên hệ Bài 9 mục 4
	X_mean, X_std = X.mean(axis=0), X.std(axis=0)
	X_normalized = (X - X_mean) / X_std

	return X_normalized, y

if __name__ == "__main__":
	df = pd.DataFrame({
		"area": [120, 85, 200, 60, np.nan],
		"rooms": [3, 2, 4, 1, 2],
		"price": [500, 350, 800, 250, 300],
	})
	X, y = preprocess_dataset(df, ["area", "rooms"], "price")
	print(X)
	print(y)
```

## Bài tập

1. **Vector hóa**: so sánh thời gian tính tổng bình phương 1 triệu số bằng vòng lặp Python thuần vs `np.sum(arr**2)`.
2. **Broadcasting & chuẩn hóa**: tạo 1 ma trận dữ liệu ngẫu nhiên (100 mẫu, 5 đặc trưng), chuẩn hóa về mean=0, std=1 bằng broadcasting (không dùng vòng lặp).
3. **Pandas EDA cơ bản**: load 1 dataset CSV thật (hoặc `sklearn.datasets`), dùng `.info()`, `.describe()`, `.isnull().sum()` để hiểu dữ liệu trước khi xử lý.
4. **Pipeline tiền xử lý**: dùng code mẫu trên làm nền, tự viết lại `preprocess_dataset`, áp dụng cho 1 dataset có giá trị thiếu và trùng lặp cố ý chèn vào để test.

## Tiếp theo
→ [Bài 13: Trực quan hóa dữ liệu & EDA](./13_data_visualization_eda.md)
