# Bài 1: Giới thiệu Machine Learning & Cài đặt môi trường

## Mục tiêu
- Hiểu Machine Learning là gì, 3 loại bài toán chính.
- Cài đặt môi trường: Anaconda, Jupyter Notebook, các thư viện cốt lõi.

## 1. Machine Learning là gì?

Machine Learning (ML) là nhánh của AI, nơi máy tính **học quy luật từ dữ liệu** thay vì được lập trình luật cứng thủ công. Thay vì viết `if/else` cho mọi trường hợp (như code truyền thống ở [Go](../Go/ROADMAP.md)/[Python](../Python/ROADMAP.md)), bạn cung cấp dữ liệu + 1 hàm mục tiêu (loss function), thuật toán tự điều chỉnh tham số để tối ưu hàm đó.

```
Lập trình truyền thống:  Dữ liệu + Luật (code)  →  Kết quả
Machine Learning:         Dữ liệu + Kết quả       →  Luật (model)
```

## 2. 3 loại bài toán ML chính

| Loại | Dữ liệu | Mục tiêu | Ví dụ | Bài học |
|---|---|---|---|---|
| **Supervised Learning** | Có nhãn (input, output đúng) | Học ánh xạ input → output | Dự đoán giá nhà, phân loại email spam | [Bài 14-19](./14_linear_regression.md) |
| **Unsupervised Learning** | Không có nhãn | Tìm cấu trúc ẩn trong dữ liệu | Phân cụm khách hàng, giảm chiều dữ liệu | [Bài 19](./19_svm_knn_clustering.md) |
| **Reinforcement Learning** | Agent tương tác với môi trường | Tối đa hóa phần thưởng tích lũy | Chơi game, robot | Ngoài phạm vi lộ trình này |

Supervised Learning lại chia 2 nhóm: **Regression** (dự đoán giá trị liên tục, vd giá nhà) và **Classification** (dự đoán nhãn rời rạc, vd spam/không spam).

## 3. Vì sao cần học toán trước khi học ML?

```
Toán nền tảng (Phần I)          →  ML dùng toán đó (Phần II-III)
Đại số tuyến tính (vector/ma trận)  →  Biểu diễn dữ liệu, phép tính trong model
Giải tích (đạo hàm/gradient)         →  Thuật toán "học" (Gradient Descent, Backprop)
Xác suất & Thống kê                   →  Hàm loss, đánh giá độ tin cậy, giả định model
Toán rời rạc & Tối ưu hóa              →  Tối ưu hàm loss, độ phức tạp thuật toán
```

Nếu bỏ qua toán, bạn vẫn gọi được `model.fit(X, y)` của scikit-learn — nhưng sẽ không biết **tại sao** model không hội tụ, **tại sao** cần regularization, hay **cách** đọc 1 paper ML mới. Lộ trình này dành hẳn Phần I cho toán trước khi chạm vào code ML thật.

## 4. Cài đặt môi trường

### Cách 1: Anaconda (khuyến khích cho người mới)

Tải Anaconda tại https://www.anaconda.com/download — cài sẵn Python + hầu hết thư viện Data Science + Jupyter.

```powershell
conda create -n ml python=3.11
conda activate ml
conda install numpy pandas matplotlib seaborn scikit-learn jupyter
```

### Cách 2: venv thuần (nếu đã quen từ [Python Bài 1](../Python/1_get_started.md))

```powershell
python -m venv ml-env
.\ml-env\Scripts\Activate.ps1
pip install numpy pandas matplotlib seaborn scikit-learn jupyter notebook
```

### Cho Deep Learning (dùng từ Phần III — [Bài 20](./20_neural_networks.md) trở đi)

```powershell
pip install torch torchvision  # PyTorch — framework khuyến khích cho lộ trình này
```

## 5. Jupyter Notebook — công cụ chuẩn cho ML

```powershell
jupyter notebook
```

Khác với viết script `.py` chạy 1 lần (như ở [Go](../Go/ROADMAP.md)/[Python](../Python/ROADMAP.md)), Jupyter Notebook cho phép chạy từng "cell" độc lập, giữ lại state giữa các lần chạy — cực kỳ phù hợp cho việc khám phá dữ liệu (EDA — [Bài 13](./13_data_visualization_eda.md)) và thử nghiệm model lặp đi lặp lại.

```python
# Cell 1
import numpy as np
import pandas as pd
print(np.__version__, pd.__version__)
```

```python
# Cell 2 — chạy độc lập, dùng lại biến từ Cell 1
arr = np.array([1, 2, 3])
print(arr.mean())
```

## 6. Load thử 1 dataset mẫu

```python
from sklearn.datasets import load_iris
import pandas as pd

data = load_iris()
df = pd.DataFrame(data.data, columns=data.feature_names)
df['target'] = data.target

print(df.head())
print(df.describe())
```

Dataset Iris (phân loại 3 loài hoa dựa trên 4 đặc trưng) sẽ được dùng làm ví dụ xuyên suốt ở Phần II.

## Bài tập

1. **Cài môi trường**: cài Anaconda hoặc venv + thư viện cốt lõi, xác nhận `import numpy, pandas, matplotlib, sklearn` không lỗi.
2. **Jupyter cơ bản**: mở Jupyter Notebook, tạo 3 cell: (1) import thư viện, (2) tạo 1 mảng NumPy, (3) tính mean/std của mảng đó.
3. **Load dataset**: load `load_iris()` (hoặc `load_diabetes()`), in `df.head()`, `df.describe()`, `df.info()` — quan sát ý nghĩa từng cột.
4. **Phân loại bài toán**: với 3 bài toán sau, xác định loại ML phù hợp (Regression/Classification/Clustering) và giải thích tại sao: (a) dự đoán nhiệt độ ngày mai, (b) phân loại ảnh chó/mèo, (c) nhóm khách hàng theo hành vi mua sắm không có nhãn trước.

## Tiếp theo
→ [Bài 2: Vector & Không gian vector](./2_linear_algebra_vectors.md)
