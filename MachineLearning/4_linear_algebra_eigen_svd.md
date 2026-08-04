# Bài 4: Trị riêng, Vector riêng & SVD

## Mục tiêu
- Eigenvalue/Eigenvector, ý nghĩa hình học.
- Singular Value Decomposition (SVD).
- Ứng dụng trực tiếp: PCA (giảm chiều dữ liệu).

## 1. Trị riêng & Vector riêng (Eigenvalue/Eigenvector)

Với ma trận vuông $A$, vector $\vec{v} \neq 0$ là **vector riêng** với **trị riêng** $\lambda$ nếu:

$$A\vec{v} = \lambda\vec{v}$$

Ý nghĩa hình học: khi ma trận $A$ "biến đổi" không gian (xoay, kéo dãn...), hầu hết vector đổi cả hướng lẫn độ dài — nhưng vector riêng $\vec{v}$ chỉ bị **kéo dãn theo đúng hướng cũ** (theo hệ số $\lambda$), không đổi hướng.

```python
import numpy as np

A = np.array([[2, 0], [0, 3]])
eigenvalues, eigenvectors = np.linalg.eig(A)
print("Trị riêng:", eigenvalues)      # [2. 3.]
print("Vector riêng:\n", eigenvectors)  # các cột là vector riêng tương ứng

# Kiểm chứng: A @ v = lambda * v
v = eigenvectors[:, 0]
lam = eigenvalues[0]
print(A @ v)          # phải xấp xỉ lam * v
print(lam * v)
```

## 2. Chéo hóa ma trận (Eigendecomposition)

Nếu $A$ có đủ $n$ vector riêng độc lập tuyến tính, có thể phân tích:

$$A = V\Lambda V^{-1}$$

trong đó $V$ là ma trận có các cột là vector riêng, $\Lambda$ là ma trận đường chéo chứa trị riêng tương ứng. Phân tích này giúp tính $A^k$ (lũy thừa ma trận) rất nhanh: $A^k = V\Lambda^k V^{-1}$ — hữu ích khi phân tích hành vi dài hạn của hệ động lực (dùng trong 1 số thuật toán RL/Markov Chain).

## 3. Ma trận đối xứng & xác định dương (Positive Semi-Definite)

Ma trận **hiệp phương sai (covariance matrix)** — nền tảng của PCA — luôn là ma trận **đối xứng và xác định dương (PSD)**:

$$\vec{x}^T A \vec{x} \geq 0 \quad \forall \vec{x}$$

Tính chất quan trọng: mọi ma trận đối xứng đều có **trị riêng thực** và **vector riêng trực giao** (vuông góc với nhau) — đây chính là lý do PCA hoạt động đúng đắn về mặt toán học.

## 4. Singular Value Decomposition (SVD) — tổng quát hơn eigendecomposition

SVD phân tích **bất kỳ ma trận nào** (kể cả không vuông) thành:

$$A = U\Sigma V^T$$

- $U \in \mathbb{R}^{m \times m}$: ma trận trực giao (cột là "left singular vectors").
- $\Sigma \in \mathbb{R}^{m \times n}$: ma trận đường chéo chứa **singular values** $\sigma_1 \geq \sigma_2 \geq ... \geq 0$ (luôn không âm, sắp giảm dần).
- $V \in \mathbb{R}^{n \times n}$: ma trận trực giao (cột là "right singular vectors").

```python
A = np.array([[3, 1], [1, 3], [0, 2]])  # ma trận 3x2, KHÔNG vuông
U, S, Vt = np.linalg.svd(A)
print("U shape:", U.shape)     # (3, 3)
print("Singular values:", S)    # [4.47, 2.83]
print("Vt shape:", Vt.shape)     # (2, 2)

# Tái tạo lại A từ U, S, Vt
Sigma = np.zeros((3, 2))
np.fill_diagonal(Sigma, S)
A_reconstructed = U @ Sigma @ Vt
print(np.allclose(A, A_reconstructed))  # True
```

## 5. SVD dùng để nén dữ liệu (Low-rank Approximation)

Chỉ giữ lại $k$ singular value **lớn nhất** cho ta xấp xỉ tốt nhất của $A$ với hạng (rank) $k$ — đây là nguyên lý nén ảnh, nén dữ liệu:

```python
def low_rank_approx(A, k):
	U, S, Vt = np.linalg.svd(A, full_matrices=False)
	return U[:, :k] @ np.diag(S[:k]) @ Vt[:k, :]

A = np.random.rand(100, 50)  # ma trận lớn giả lập (vd 1 ảnh grayscale)
A_compressed = low_rank_approx(A, k=10)  # chỉ giữ 10 "thành phần" quan trọng nhất
print(f"Sai số nén: {np.linalg.norm(A - A_compressed):.4f}")
```

## 6. Ứng dụng quan trọng nhất: PCA (Principal Component Analysis)

PCA tìm các hướng (principal components) mà dữ liệu **biến thiên nhiều nhất** — đây chính là các **vector riêng của ma trận hiệp phương sai** dữ liệu, sắp theo trị riêng giảm dần.

```python
def pca_manual(X, n_components):
	# Bước 1: chuẩn hóa dữ liệu về mean = 0
	X_centered = X - X.mean(axis=0)

	# Bước 2: tính ma trận hiệp phương sai
	cov_matrix = np.cov(X_centered.T)

	# Bước 3: tìm trị riêng & vector riêng
	eigenvalues, eigenvectors = np.linalg.eig(cov_matrix)

	# Bước 4: sắp xếp theo trị riêng giảm dần, lấy k vector riêng đầu
	idx = np.argsort(eigenvalues)[::-1]
	top_k_vectors = eigenvectors[:, idx[:n_components]]

	# Bước 5: chiếu dữ liệu lên không gian mới (giảm chiều)
	X_reduced = X_centered @ top_k_vectors
	return X_reduced, eigenvalues[idx[:n_components]]

from sklearn.datasets import load_iris
data = load_iris()
X_reduced, explained = pca_manual(data.data, n_components=2)
print(X_reduced.shape)  # (150, 2) — từ 4 chiều xuống còn 2 chiều
```

Trị riêng lớn = phương sai (thông tin) mà thành phần đó "giữ lại" — nên ta luôn sắp xếp và chọn trị riêng lớn nhất. **Tỷ lệ phương sai giải thích được** (explained variance ratio) $= \frac{\lambda_i}{\sum \lambda_j}$ cho biết mất bao nhiêu thông tin khi giảm chiều — đây là con số quan trọng khi quyết định giữ bao nhiêu chiều.

Chi tiết dùng PCA thực hành với scikit-learn ở [Bài 19](./19_svm_knn_clustering.md).

## Ví dụ đầy đủ

```python
import numpy as np

def pca_manual(X, n_components):
	X_centered = X - X.mean(axis=0)
	cov_matrix = np.cov(X_centered.T)
	eigenvalues, eigenvectors = np.linalg.eig(cov_matrix)
	idx = np.argsort(eigenvalues)[::-1]
	eigenvalues, eigenvectors = eigenvalues[idx], eigenvectors[:, idx]

	top_k = eigenvectors[:, :n_components]
	X_reduced = X_centered @ top_k

	explained_ratio = eigenvalues[:n_components] / eigenvalues.sum()
	return X_reduced, explained_ratio

if __name__ == "__main__":
	np.random.seed(42)
	X = np.random.multivariate_normal(mean=[0, 0], cov=[[3, 2], [2, 2]], size=200)

	X_reduced, ratio = pca_manual(X, n_components=1)
	print("Tỷ lệ phương sai giữ lại:", ratio)  # gần 1.0 vì dữ liệu 2D có tương quan mạnh
```

## Bài tập

1. **Trị riêng/vector riêng bằng tay**: tính trị riêng của ma trận $\begin{bmatrix} 4 & 1 \\ 2 & 3 \end{bmatrix}$ bằng tay (giải phương trình đặc trưng $\det(A - \lambda I) = 0$), verify bằng `np.linalg.eig`.
2. **SVD & nén dữ liệu**: tạo 1 ma trận ngẫu nhiên 50x50, thử nén với $k = 5, 10, 20$, vẽ đồ thị sai số nén giảm dần theo $k$ (liên hệ [Bài 13](./13_data_visualization_eda.md) cho phần vẽ).
3. **PCA thủ công**: dùng code mẫu trên làm nền, tự viết lại `pca_manual`, áp dụng cho dataset Iris, giảm từ 4 chiều xuống 2 chiều, in tỷ lệ phương sai giữ lại.
4. **Nâng cao**: so sánh kết quả PCA thủ công (mục 6) với `sklearn.decomposition.PCA` — verify 2 kết quả gần giống nhau (có thể lệch dấu do vector riêng không duy nhất về hướng).

## Tổng kết phần Đại số tuyến tính
Bạn đã nắm vector, ma trận, trị riêng/SVD — đủ nền để hiểu cách dữ liệu được biểu diễn và biến đổi trong mọi thuật toán ML. Tiếp theo là Giải tích — công cụ giúp model "học" từ dữ liệu.

## Tiếp theo
→ [Bài 5: Đạo hàm & Gradient](./5_calculus_derivatives.md)
