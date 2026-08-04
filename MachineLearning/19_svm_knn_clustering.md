# Bài 19: SVM, k-NN, Naive Bayes & Unsupervised Learning

## Mục tiêu
- Support Vector Machine (margin, kernel trick), k-Nearest Neighbors, Naive Bayes.
- k-Means Clustering, PCA thực hành.

## 1. Support Vector Machine (SVM) — tối đa hóa margin

SVM tìm siêu phẳng (hyperplane) phân tách 2 lớp sao cho **khoảng cách (margin)** từ siêu phẳng tới điểm gần nhất mỗi lớp là **lớn nhất**:

$$\max_{\vec{w},b} \frac{2}{\|\vec{w}\|} \quad\text{sao cho}\quad y_i(\vec{w}^T\vec{x_i}+b) \geq 1, \, \forall i$$

Các điểm nằm đúng trên biên margin gọi là **support vectors** — chỉ chúng quyết định vị trí siêu phẳng, các điểm khác không ảnh hưởng.

```python
from sklearn.svm import SVC
from sklearn.datasets import make_classification
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler

X, y = make_classification(n_samples=200, n_features=2, n_redundant=0, n_clusters_per_class=1, random_state=42)
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)
X_train, X_test = StandardScaler().fit_transform(X_train), StandardScaler().fit_transform(X_test)

svm = SVC(kernel="linear", C=1.0)
svm.fit(X_train, y_train)
print("Support vectors:", svm.support_vectors_.shape[0])
```

## 2. Soft Margin — cho phép sai số khi dữ liệu không tách rời tuyến tính hoàn hảo

$$\min_{\vec{w},b,\xi} \frac{1}{2}\|\vec{w}\|^2 + C\sum_i \xi_i$$

$C$ (regularization — liên hệ [Bài 17](./17_regularization.md)) điều khiển đánh đổi giữa margin rộng và số điểm bị phân loại sai: $C$ nhỏ → margin rộng, chấp nhận nhiều lỗi (giảm variance); $C$ lớn → cố phân loại đúng mọi điểm, margin hẹp (dễ overfitting).

## 3. Kernel Trick — xử lý dữ liệu KHÔNG tách rời tuyến tính

Khi dữ liệu không tách được bằng đường thẳng trong không gian gốc, kernel "ánh xạ ngầm" dữ liệu lên không gian nhiều chiều hơn (nơi nó CÓ THỂ tách tuyến tính) mà **không cần tính tường minh phép ánh xạ** — chỉ cần tính hàm kernel $K(\vec{x_i}, \vec{x_j})$ thay cho dot product ([Bài 2](./2_linear_algebra_vectors.md)):

```python
svm_rbf = SVC(kernel="rbf", C=1.0, gamma="scale")  # RBF (Gaussian) kernel — phổ biến nhất
svm_rbf.fit(X_train, y_train)
```

Kernel RBF phổ biến nhất: $K(\vec{x_i},\vec{x_j}) = e^{-\gamma\|\vec{x_i}-\vec{x_j}\|^2}$ — chính là hàm mật độ Gaussian ([Bài 9 mục 6](./9_probability_distributions.md)) áp dụng lên khoảng cách Euclidean ([Bài 2 mục 5](./2_linear_algebra_vectors.md)).

## 4. k-Nearest Neighbors (k-NN) — thuật toán đơn giản nhất, dựa hoàn toàn vào khoảng cách

```python
from sklearn.neighbors import KNeighborsClassifier

knn = KNeighborsClassifier(n_neighbors=5)
knn.fit(X_train, y_train)
print("Accuracy:", knn.score(X_test, y_test))
```

k-NN **không "học" tham số** — dự đoán bằng cách tìm $k$ điểm gần nhất (theo Euclidean distance — [Bài 2 mục 5](./2_linear_algebra_vectors.md)) trong tập train, vote đa số (classification) hoặc lấy trung bình (regression). Chi phí dự đoán $O(m)$ mỗi truy vấn ([Bài 11 mục 3](./11_discrete_math_optimization.md)) — chậm với dataset lớn. **Luôn chuẩn hóa dữ liệu trước** ([Bài 9 mục 4](./9_probability_distributions.md)) vì k-NN cực kỳ nhạy với thang đo đặc trưng.

## 5. Naive Bayes — áp dụng trực tiếp Định lý Bayes

Đã giới thiệu công thức ở [Bài 8 mục 7](./8_probability_basics.md):

$$P(y|\vec{x}) \propto P(y)\prod_{i=1}^n P(x_i|y)$$

```python
from sklearn.naive_bayes import GaussianNB

gnb = GaussianNB()  # giả định P(x_i|y) tuân theo phân phối Gaussian
gnb.fit(X_train, y_train)
print("Accuracy:", gnb.score(X_test, y_test))
```

`GaussianNB` giả định mỗi đặc trưng tuân theo phân phối chuẩn ([Bài 9 mục 6](./9_probability_distributions.md)) với mỗi lớp — ước lượng $\mu, \sigma$ bằng MLE ([Bài 10 mục 1](./10_statistics_inference.md)) từ dữ liệu train. `MultinomialNB` dùng cho dữ liệu đếm (vd số lần xuất hiện từ trong văn bản — bài toán spam filter kinh điển).

## 6. k-Means Clustering — Unsupervised Learning

Chia dữ liệu thành $k$ cụm, mỗi điểm thuộc cụm có **centroid (tâm)** gần nhất, lặp lại tối ưu:

$$\min_{\{C_1,...,C_k\}} \sum_{j=1}^k \sum_{\vec{x} \in C_j} \|\vec{x} - \vec{\mu_j}\|^2$$

```python
from sklearn.cluster import KMeans
import matplotlib.pyplot as plt

kmeans = KMeans(n_clusters=3, random_state=42, n_init=10)
labels = kmeans.fit_predict(X)

plt.scatter(X[:,0], X[:,1], c=labels, cmap="viridis")
plt.scatter(kmeans.cluster_centers_[:,0], kmeans.cluster_centers_[:,1], marker="X", s=200, c="red")
plt.title("K-Means Clustering")
plt.show()
```

Thuật toán lặp 2 bước tới hội tụ: (1) gán mỗi điểm vào centroid gần nhất (Euclidean distance — [Bài 2](./2_linear_algebra_vectors.md)), (2) cập nhật centroid = trung bình các điểm trong cụm ([Bài 9 mục 3](./9_probability_distributions.md)) — độ phức tạp $O(kmi)$ đã nhắc ở [Bài 11 mục 3](./11_discrete_math_optimization.md).

### Chọn $k$ bằng Elbow Method

```python
inertias = []
for k in range(1, 10):
	km = KMeans(n_clusters=k, random_state=42, n_init=10).fit(X)
	inertias.append(km.inertia_)  # tổng khoảng cách bình phương tới centroid

plt.plot(range(1, 10), inertias, marker="o")
plt.xlabel("k"); plt.ylabel("Inertia")
plt.title("Elbow Method")
plt.show()
```

Tìm điểm "khuỷu tay" (elbow) — nơi inertia ngừng giảm mạnh — làm $k$ hợp lý.

## 7. PCA thực hành với scikit-learn (nối tiếp [Bài 4](./4_linear_algebra_eigen_svd.md))

```python
from sklearn.decomposition import PCA
from sklearn.datasets import load_iris

data = load_iris()
pca = PCA(n_components=2)
X_reduced = pca.fit_transform(data.data)

print("Tỷ lệ phương sai giải thích:", pca.explained_variance_ratio_)

plt.scatter(X_reduced[:,0], X_reduced[:,1], c=data.target, cmap="viridis")
plt.xlabel("PC1"); plt.ylabel("PC2")
plt.title("Iris dataset sau PCA (4D -> 2D)")
plt.show()
```

Kết quả này chính là những gì bạn đã tự cài đặt bằng tay ở [Bài 4 mục 6](./4_linear_algebra_eigen_svd.md) — `pca.components_` chứa các vector riêng, `pca.explained_variance_ratio_` chứa tỷ lệ trị riêng đã tính thủ công.

## Ví dụ đầy đủ: pipeline PCA + Clustering

```python
from sklearn.datasets import load_iris
from sklearn.decomposition import PCA
from sklearn.cluster import KMeans
from sklearn.preprocessing import StandardScaler

data = load_iris()
X_scaled = StandardScaler().fit_transform(data.data)

X_pca = PCA(n_components=2).fit_transform(X_scaled)

kmeans = KMeans(n_clusters=3, random_state=42, n_init=10)
cluster_labels = kmeans.fit_predict(X_pca)

# So sánh cluster tự tìm được với nhãn thật (chỉ để kiểm tra, KHÔNG dùng nhãn khi cluster)
from sklearn.metrics import adjusted_rand_score
print("Adjusted Rand Index:", adjusted_rand_score(data.target, cluster_labels))
```

## Bài tập

1. **SVM với các kernel khác nhau**: huấn luyện SVM với kernel `linear`, `rbf`, `poly` trên dataset không tách tuyến tính (`make_moons`), so sánh decision boundary (dùng hàm vẽ ở [Bài 15 mục 7](./15_logistic_regression.md)).
2. **k-NN & ảnh hưởng của k**: thử `n_neighbors` = 1, 5, 20, quan sát overfitting/underfitting (liên hệ [Bài 16 mục 4](./16_model_evaluation.md)) — k nhỏ dễ overfit, k lớn dễ underfit.
3. **Naive Bayes cho text classification**: dùng `MultinomialNB` + `CountVectorizer` (scikit-learn) phân loại 1 tập văn bản đơn giản thành 2 chủ đề.
4. **k-Means + Elbow method**: áp dụng cho dataset tự tạo (`make_blobs` với 4 cụm biết trước), dùng Elbow Method tìm lại đúng $k=4$, so sánh cluster tìm được với nhãn thật.

## Tổng kết Phần II — Machine Learning Cơ Bản
Bạn đã nắm toàn bộ thuật toán ML cổ điển: Linear/Logistic Regression, đánh giá mô hình, regularization, Tree-based methods, SVM/k-NN/Naive Bayes, và Unsupervised Learning (k-Means, PCA) — mỗi thuật toán đều được giải thích từ toán nền tảng ở Phần I. Phần III sẽ mở rộng sang Deep Learning — nơi Chain Rule ([Bài 6](./6_calculus_matrix_calculus.md)) trở thành công cụ trung tâm.

## Tiếp theo
→ [Bài 20: Neural Network cơ bản — Perceptron & MLP](./20_neural_networks.md)
