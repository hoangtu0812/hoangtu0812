# Bài 3: Ma trận & Hệ phương trình tuyến tính

## Mục tiêu
- Phép toán ma trận: cộng, nhân, chuyển vị, nghịch đảo, định thức.
- Giải hệ phương trình tuyến tính.
- Ứng dụng: biểu diễn dataset dạng ma trận, Normal Equation cho Linear Regression.

## 1. Ma trận là gì?

Ma trận là bảng số chữ nhật $m \times n$ ($m$ dòng, $n$ cột):

$$A = \begin{bmatrix} a_{11} & a_{12} & \dots & a_{1n} \\ a_{21} & a_{22} & \dots & a_{2n} \\ \vdots & & \ddots & \vdots \\ a_{m1} & a_{m2} & \dots & a_{mn} \end{bmatrix} \in \mathbb{R}^{m \times n}$$

**Trong ML:** toàn bộ dataset là 1 ma trận **design matrix** $X \in \mathbb{R}^{m \times n}$ — $m$ dòng (số mẫu), $n$ cột (số đặc trưng). Đây là cách dữ liệu được biểu diễn xuyên suốt mọi thuật toán ML.

```python
import numpy as np

# Dataset: 4 căn nhà, mỗi nhà có 2 đặc trưng [diện tích, số phòng]
X = np.array([
	[120, 3],
	[85,  2],
	[200, 4],
	[60,  1],
])
print(X.shape)  # (4, 2) — 4 mẫu, 2 đặc trưng
```

## 2. Phép toán ma trận cơ bản

```python
A = np.array([[1, 2], [3, 4]])
B = np.array([[5, 6], [7, 8]])

print(A + B)       # cộng từng phần tử: [[6 8] [10 12]]
print(A * B)        # nhân TỪNG PHẦN TỬ (element-wise) — KHÔNG phải nhân ma trận!
print(A @ B)         # nhân ma trận THẬT SỰ (matrix multiplication) — dùng @ hoặc np.matmul
print(A.T)            # chuyển vị (transpose) — đổi dòng thành cột
```

**Cạm bẫy phổ biến nhất khi mới học NumPy:** `A * B` là nhân element-wise, `A @ B` mới là phép nhân ma trận chuẩn — nhầm lẫn 2 cái này là nguồn lỗi rất thường gặp khi code ML.

## 3. Phép nhân ma trận — công thức

$$C = AB \quad\text{với}\quad C_{ij} = \sum_{k=1}^n A_{ik}B_{kj}$$

Điều kiện: số cột của $A$ phải bằng số dòng của $B$. Kích thước: $A_{m \times n} \cdot B_{n \times p} = C_{m \times p}$.

```python
A = np.array([[1, 2, 3], [4, 5, 6]])  # 2x3
B = np.array([[7, 8], [9, 10], [11, 12]])  # 3x2
C = A @ B  # kết quả 2x2
print(C)
```

**Ứng dụng ML then chốt:** dự đoán của Linear Regression cho **toàn bộ dataset cùng lúc** (không phải từng mẫu 1 như [Bài 2](./2_linear_algebra_vectors.md)) là:

$$\hat{y} = X\vec{w} + b$$

với $X \in \mathbb{R}^{m \times n}$ (design matrix), $\vec{w} \in \mathbb{R}^n$ (vector trọng số) → $\hat{y} \in \mathbb{R}^m$ (vector dự đoán cho $m$ mẫu) — đây là lý do NumPy vector hóa nhanh hơn hẳn vòng lặp `for` thủ công.

## 4. Ma trận đơn vị (Identity Matrix) & Ma trận nghịch đảo (Inverse)

```python
I = np.eye(3)  # ma trận đơn vị 3x3 — đường chéo = 1, còn lại = 0
print(I)

A = np.array([[4, 7], [2, 6]])
A_inv = np.linalg.inv(A)
print(A_inv)
print(A @ A_inv)  # xấp xỉ ma trận đơn vị I (do sai số floating-point)
```

$A^{-1}A = AA^{-1} = I$ — tương tự số nghịch đảo $a^{-1} \cdot a = 1$ trong số học thông thường. **Không phải ma trận nào cũng có nghịch đảo** — chỉ ma trận vuông và **không suy biến** (định thức khác 0) mới khả nghịch.

## 5. Định thức (Determinant)

```python
A = np.array([[4, 7], [2, 6]])
det = np.linalg.det(A)  # 4*6 - 7*2 = 10
print(det)
```

Định thức $= 0$ nghĩa là ma trận **suy biến (singular)** — không có nghịch đảo, các dòng/cột phụ thuộc tuyến tính lẫn nhau ([Bài 2 mục 7](./2_linear_algebra_vectors.md)). **Ứng dụng ML:** nếu design matrix có 2 cột đặc trưng phụ thuộc tuyến tính, Normal Equation (mục 7) sẽ thất bại vì cần nghịch đảo ma trận suy biến.

## 6. Giải hệ phương trình tuyến tính

Hệ phương trình $A\vec{x} = \vec{b}$:

$$\begin{cases} 2x_1 + 3x_2 = 8 \\ x_1 - x_2 = 1 \end{cases} \quad\Leftrightarrow\quad \begin{bmatrix} 2 & 3 \\ 1 & -1 \end{bmatrix}\begin{bmatrix} x_1 \\ x_2 \end{bmatrix} = \begin{bmatrix} 8 \\ 1 \end{bmatrix}$$

```python
A = np.array([[2, 3], [1, -1]])
b = np.array([8, 1])

x = np.linalg.solve(A, b)  # cách ĐÚNG — nhanh và ổn định số học hơn tính A_inv @ b
print(x)  # [2.2, 1.2]
```

**Luôn dùng `np.linalg.solve()` thay vì `np.linalg.inv(A) @ b`** — về mặt toán học cho cùng kết quả, nhưng `solve()` dùng thuật toán số học ổn định hơn (không tính nghịch đảo tường minh), quan trọng khi ma trận lớn trong ML thực tế.

## 7. Normal Equation — nghiệm đóng (closed-form) của Linear Regression

Thay vì dùng Gradient Descent ([Bài 11](./11_discrete_math_optimization.md)) để tìm $\vec{w}$ tối ưu cho Linear Regression, có thể giải trực tiếp bằng đại số tuyến tính:

$$\vec{w}^* = (X^TX)^{-1}X^T\vec{y}$$

```python
X = np.array([[120, 3], [85, 2], [200, 4], [60, 1]])
X_with_bias = np.hstack([np.ones((X.shape[0], 1)), X])  # thêm cột 1 cho hệ số bias b
y = np.array([500, 350, 800, 250])  # giá nhà (đơn vị triệu)

w = np.linalg.inv(X_with_bias.T @ X_with_bias) @ X_with_bias.T @ y
print(w)  # [bias, w1, w2]
```

Chi tiết đầy đủ về Linear Regression (bao gồm cả Gradient Descent để so sánh) ở [Bài 14](./14_linear_regression.md).

## Ví dụ đầy đủ

```python
import numpy as np

def solve_linear_regression_normal_eq(X, y):
	"""Giải Linear Regression bằng Normal Equation."""
	X_b = np.hstack([np.ones((X.shape[0], 1)), X])
	w = np.linalg.inv(X_b.T @ X_b) @ X_b.T @ y
	return w  # w[0] = bias, w[1:] = trọng số các đặc trưng

def predict(X, w):
	X_b = np.hstack([np.ones((X.shape[0], 1)), X])
	return X_b @ w

if __name__ == "__main__":
	X = np.array([[1], [2], [3], [4]], dtype=float)
	y = np.array([3, 5, 7, 9], dtype=float)  # y = 2x + 1

	w = solve_linear_regression_normal_eq(X, y)
	print("w:", w)  # xấp xỉ [1, 2] — bias=1, weight=2

	predictions = predict(np.array([[5], [6]]), w)
	print("Dự đoán:", predictions)  # xấp xỉ [11, 13]
```

## Bài tập

1. **Phép toán ma trận**: tự tính bằng tay phép nhân 2 ma trận 2x2, kiểm tra lại bằng `@`; phân biệt rõ `A * B` vs `A @ B` bằng ví dụ cụ thể.
2. **Định thức & nghịch đảo**: tính định thức và ma trận nghịch đảo của 1 ma trận 2x2 bằng tay, verify bằng `np.linalg.det`/`np.linalg.inv`.
3. **Giải hệ phương trình**: giải hệ 3 phương trình 3 ẩn bằng `np.linalg.solve`, so sánh với giải bằng tay (phương pháp thế hoặc khử Gauss).
4. **Normal Equation**: dùng code mẫu trên làm nền, tự viết lại `solve_linear_regression_normal_eq`, thử với 1 dataset tự tạo có quan hệ tuyến tính rõ ràng ($y = 3x - 2$ + nhiễu nhỏ).

## Tiếp theo
→ [Bài 4: Trị riêng, Vector riêng & SVD](./4_linear_algebra_eigen_svd.md)
