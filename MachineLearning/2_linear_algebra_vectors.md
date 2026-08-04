# Bài 2: Vector & Không gian vector

## Mục tiêu
- Hiểu vector, các phép toán cơ bản, tích vô hướng, norm.
- Cosine similarity, góc giữa 2 vector.
- Kết nối trực tiếp tới cách ML biểu diễn dữ liệu.

## 1. Vector là gì?

Vector là 1 danh sách số có thứ tự, biểu diễn 1 điểm/hướng trong không gian nhiều chiều:

$$\vec{v} = \begin{bmatrix} v_1 \\ v_2 \\ \vdots \\ v_n \end{bmatrix} \in \mathbb{R}^n$$

**Trong ML:** mỗi mẫu dữ liệu (1 dòng dataset) là 1 vector — vd 1 căn nhà với 3 đặc trưng (diện tích, số phòng, tuổi nhà) là vector trong $\mathbb{R}^3$: $\vec{x} = [120, 3, 5]$.

```python
import numpy as np

v = np.array([120, 3, 5])  # 1 mẫu dữ liệu: [diện tích, số phòng, tuổi nhà]
print(v.shape)  # (3,) — vector 3 chiều
```

## 2. Phép toán vector cơ bản

```python
u = np.array([1, 2, 3])
v = np.array([4, 5, 6])

print(u + v)       # [5 7 9] — cộng từng phần tử
print(u - v)       # [-3 -3 -3]
print(2 * u)        # [2 4 6] — nhân vô hướng (scalar multiplication)
```

Về mặt toán học:
$$\vec{u} + \vec{v} = \begin{bmatrix} u_1+v_1 \\ u_2+v_2 \\ u_3+v_3 \end{bmatrix}, \quad c\vec{u} = \begin{bmatrix} cu_1 \\ cu_2 \\ cu_3 \end{bmatrix}$$

## 3. Tích vô hướng (Dot Product) — phép toán quan trọng nhất trong ML

$$\vec{u} \cdot \vec{v} = \sum_{i=1}^n u_i v_i = u_1v_1 + u_2v_2 + \dots + u_nv_n$$

```python
u = np.array([1, 2, 3])
v = np.array([4, 5, 6])

dot = np.dot(u, v)  # 1*4 + 2*5 + 3*6 = 32
print(dot)
```

**Ứng dụng ML trực tiếp nhất:** dự đoán của Linear Regression ([Bài 14](./14_linear_regression.md)) chính là tích vô hướng giữa vector trọng số $\vec{w}$ và vector đặc trưng $\vec{x}$:

$$\hat{y} = \vec{w} \cdot \vec{x} + b = w_1x_1 + w_2x_2 + \dots + w_nx_n + b$$

Đây cũng chính là phép tính lõi bên trong **mỗi neuron** của mạng neural ([Bài 20](./20_neural_networks.md)).

## 4. Norm — "độ dài" của vector

Norm phổ biến nhất: **Euclidean norm (L2 norm)**:

$$\|\vec{v}\|_2 = \sqrt{v_1^2 + v_2^2 + \dots + v_n^2} = \sqrt{\vec{v} \cdot \vec{v}}$$

```python
v = np.array([3, 4])
norm_l2 = np.linalg.norm(v)  # sqrt(3^2 + 4^2) = sqrt(25) = 5
print(norm_l2)

norm_l1 = np.linalg.norm(v, ord=1)  # |3| + |4| = 7  — L1 norm
print(norm_l1)
```

**Ứng dụng ML:** L2 norm dùng trong Ridge Regression, L1 norm dùng trong Lasso Regression ([Bài 17](./17_regularization.md)) để "phạt" trọng số lớn — chính là norm của vector trọng số $\vec{w}$.

## 5. Khoảng cách giữa 2 vector (Euclidean Distance)

$$d(\vec{u}, \vec{v}) = \|\vec{u} - \vec{v}\|_2 = \sqrt{\sum_{i=1}^n (u_i - v_i)^2}$$

```python
u = np.array([1, 2])
v = np.array([4, 6])
distance = np.linalg.norm(u - v)  # sqrt((1-4)^2 + (2-6)^2) = sqrt(9+16) = 5
print(distance)
```

**Ứng dụng ML:** k-Nearest Neighbors ([Bài 19](./19_svm_knn_clustering.md)) và k-Means Clustering ([Bài 19](./19_svm_knn_clustering.md)) dựa hoàn toàn vào khoảng cách này để tìm "hàng xóm gần nhất".

## 6. Góc giữa 2 vector & Cosine Similarity

Liên hệ giữa dot product, norm, và góc $\theta$ giữa 2 vector:

$$\vec{u} \cdot \vec{v} = \|\vec{u}\|\|\vec{v}\|\cos\theta \quad\Rightarrow\quad \cos\theta = \frac{\vec{u} \cdot \vec{v}}{\|\vec{u}\|\|\vec{v}\|}$$

```python
def cosine_similarity(u, v):
	return np.dot(u, v) / (np.linalg.norm(u) * np.linalg.norm(v))

u = np.array([1, 0])
v = np.array([1, 1])
print(cosine_similarity(u, v))  # cos(45°) ≈ 0.707
```

**Ứng dụng ML rất phổ biến:** đo độ giống nhau giữa 2 văn bản (biểu diễn dạng vector embedding), recommendation system (độ giống giữa 2 người dùng/sản phẩm), NLP (word embeddings — word2vec).

## 7. Vector độc lập tuyến tính (Linear Independence) — trực giác

Tập vector $\{\vec{v_1}, ..., \vec{v_k}\}$ **độc lập tuyến tính** nếu không vector nào biểu diễn được bằng tổ hợp tuyến tính của các vector còn lại. Ý nghĩa thực tế trong ML: nếu 2 cột đặc trưng trong dataset **phụ thuộc tuyến tính** lẫn nhau (vd "diện tích m²" và "diện tích ft²" — chỉ khác hệ số nhân), chúng mang cùng 1 thông tin — gây vấn đề **multicollinearity** trong Linear Regression ([Bài 14](./14_linear_regression.md)).

## Ví dụ đầy đủ

```python
import numpy as np

def euclidean_distance(u, v):
	return np.linalg.norm(u - v)

def cosine_similarity(u, v):
	return np.dot(u, v) / (np.linalg.norm(u) * np.linalg.norm(v))

# 3 "khách hàng" biểu diễn bằng vector [số lần mua, chi tiêu trung bình]
customer_a = np.array([10, 500])
customer_b = np.array([12, 480])   # tương tự customer_a
customer_c = np.array([2, 5000])    # khác hẳn — ít lần mua nhưng chi rất nhiều

print("Khoảng cách A-B:", euclidean_distance(customer_a, customer_b))
print("Khoảng cách A-C:", euclidean_distance(customer_a, customer_c))
print("Cosine similarity A-B:", cosine_similarity(customer_a, customer_b))
```

## Bài tập

1. **Phép toán cơ bản**: tính bằng tay (giấy bút) rồi kiểm tra bằng NumPy: cộng, trừ, nhân vô hướng cho 2 vector 3 chiều tự chọn.
2. **Dot product & dự đoán tuyến tính**: cho $\vec{w} = [2, -1, 0.5]$, $\vec{x} = [3, 4, 2]$, $b = 1$, tính $\hat{y} = \vec{w} \cdot \vec{x} + b$ bằng tay, sau đó verify bằng `np.dot`.
3. **Norm & khoảng cách**: tính L1 norm, L2 norm của vector $[3, -4, 12]$; tính khoảng cách Euclidean giữa 2 điểm dữ liệu tự chọn.
4. **Cosine similarity ứng dụng**: viết hàm `most_similar(target, candidates)` trả về vector trong `candidates` có cosine similarity cao nhất với `target` — mô phỏng bài toán recommendation đơn giản.

## Tiếp theo
→ [Bài 3: Ma trận & Hệ phương trình tuyến tính](./3_linear_algebra_matrices.md)
