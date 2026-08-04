# Bài 9: Biến ngẫu nhiên & Phân phối xác suất

## Mục tiêu
- Biến ngẫu nhiên rời rạc/liên tục, PDF, CDF.
- Kỳ vọng, phương sai, hiệp phương sai.
- Các phân phối quan trọng: Bernoulli, Binomial, Poisson, Uniform, Normal/Gaussian.

## 1. Biến ngẫu nhiên (Random Variable)

Biến ngẫu nhiên $X$ ánh xạ mỗi kết quả trong không gian mẫu thành 1 số. Có 2 loại:
- **Rời rạc (discrete)**: nhận giá trị đếm được (vd số lần tung được mặt ngửa).
- **Liên tục (continuous)**: nhận giá trị trong 1 khoảng số thực (vd chiều cao, thời gian).

## 2. PMF, PDF, CDF

**Với biến rời rạc** — Probability Mass Function (PMF): $P(X = x)$ cho từng giá trị cụ thể.

**Với biến liên tục** — Probability Density Function (PDF) $f(x)$: xác suất nằm trong 1 khoảng, không phải tại 1 điểm (vì $P(X=x) = 0$ với biến liên tục):

$$P(a \leq X \leq b) = \int_a^b f(x)\,dx$$

**Cumulative Distribution Function (CDF)** — dùng cho cả 2 loại: $F(x) = P(X \leq x)$.

```python
import numpy as np
from scipy import stats

# Biến liên tục: phân phối chuẩn N(0, 1)
x = 1.5
pdf_value = stats.norm.pdf(x, loc=0, scale=1)  # giá trị "mật độ" tại x=1.5
cdf_value = stats.norm.cdf(x, loc=0, scale=1)  # P(X <= 1.5)
print(f"PDF: {pdf_value:.4f}, CDF: {cdf_value:.4f}")
```

## 3. Kỳ vọng (Expectation / Mean)

$$E[X] = \sum_x x \cdot P(X=x) \quad \text{(rời rạc)} \qquad E[X] = \int x f(x)\,dx \quad \text{(liên tục)}$$

Kỳ vọng là "giá trị trung bình lý thuyết" nếu lặp lại thí nghiệm vô hạn lần.

```python
# Kỳ vọng khi tung xúc xắc
values = np.array([1, 2, 3, 4, 5, 6])
probs = np.array([1/6]*6)
expected_value = np.sum(values * probs)
print(expected_value)  # 3.5
```

Tính chất quan trọng dùng nhiều trong ML: **tuyến tính của kỳ vọng** $E[aX + bY] = aE[X] + bE[Y]$ — đúng dù $X, Y$ có độc lập hay không.

## 4. Phương sai (Variance) & Độ lệch chuẩn (Standard Deviation)

$$\text{Var}(X) = E[(X-E[X])^2] = E[X^2] - (E[X])^2, \qquad \sigma = \sqrt{\text{Var}(X)}$$

```python
data = np.array([2, 4, 4, 4, 5, 5, 7, 9])
variance = np.var(data)         # phương sai
std_dev = np.std(data)           # độ lệch chuẩn
print(variance, std_dev)
```

**Ứng dụng ML:** phương sai đo độ "phân tán" của dữ liệu — nền tảng của **chuẩn hóa dữ liệu (standardization)**: $z = \frac{x - \mu}{\sigma}$, bước tiền xử lý gần như bắt buộc trước khi huấn luyện nhiều model (đặc biệt SVM, k-NN, Neural Network — [Bài 19](./19_svm_knn_clustering.md), [Bài 20](./20_neural_networks.md)).

## 5. Hiệp phương sai (Covariance) & Ma trận hiệp phương sai

$$\text{Cov}(X,Y) = E[(X-E[X])(Y-E[Y])]$$

- $\text{Cov}(X,Y) > 0$: $X, Y$ có xu hướng tăng/giảm cùng nhau.
- $\text{Cov}(X,Y) < 0$: $X, Y$ có xu hướng ngược chiều.
- $\text{Cov}(X,Y) = 0$: không có quan hệ tuyến tính (KHÔNG có nghĩa là độc lập hoàn toàn).

**Hệ số tương quan (correlation)** chuẩn hóa covariance về $[-1, 1]$: $\rho = \frac{\text{Cov}(X,Y)}{\sigma_X\sigma_Y}$.

```python
X = np.array([1, 2, 3, 4, 5])
Y = np.array([2, 4, 5, 4, 5])
cov_matrix = np.cov(X, Y)  # ma trận 2x2: đường chéo là Var(X), Var(Y); off-diagonal là Cov(X,Y)
print(cov_matrix)

correlation = np.corrcoef(X, Y)[0, 1]
print(correlation)
```

Ma trận hiệp phương sai chính là ma trận đối xứng xác định dương dùng trong PCA ([Bài 4](./4_linear_algebra_eigen_svd.md)).

## 6. Các phân phối quan trọng cho ML

### Bernoulli — 1 lần thử nhị phân (thành công/thất bại)
$$P(X=1) = p, \quad P(X=0) = 1-p$$
**Ứng dụng:** mô hình hóa nhãn nhị phân trong Logistic Regression ([Bài 15](./15_logistic_regression.md)).

### Binomial — số lần thành công trong $n$ lần thử Bernoulli độc lập
$$P(X=k) = \binom{n}{k}p^k(1-p)^{n-k}$$

### Poisson — số sự kiện xảy ra trong 1 khoảng thời gian/không gian cố định
$$P(X=k) = \frac{\lambda^k e^{-\lambda}}{k!}$$
**Ứng dụng:** mô hình hóa số lượng request/giây, số khách hàng đến cửa hàng.

### Uniform — mọi giá trị trong khoảng có xác suất bằng nhau
$$f(x) = \frac{1}{b-a}, \quad x \in [a,b]$$
**Ứng dụng:** khởi tạo trọng số ngẫu nhiên trong mạng neural.

### Normal/Gaussian — phân phối quan trọng NHẤT trong ML
$$f(x) = \frac{1}{\sigma\sqrt{2\pi}}e^{-\frac{(x-\mu)^2}{2\sigma^2}}$$

```python
import matplotlib.pyplot as plt

x = np.linspace(-4, 4, 200)
for sigma in [0.5, 1, 2]:
	plt.plot(x, stats.norm.pdf(x, loc=0, scale=sigma), label=f"σ={sigma}")
plt.legend()
plt.title("Phân phối chuẩn với các độ lệch chuẩn khác nhau")
plt.show()
```

**Vì sao Normal quan trọng bậc nhất:**
1. **Central Limit Theorem**: tổng/trung bình của nhiều biến ngẫu nhiên độc lập (bất kể phân phối gốc) tiến về phân phối chuẩn khi $n$ đủ lớn — lý do nhiễu đo lường trong thực tế thường được giả định là Gaussian.
2. Linear Regression ([Bài 14](./14_linear_regression.md)) giả định nhiễu $\epsilon \sim N(0, \sigma^2)$ — hàm loss MSE chính là kết quả của MLE dưới giả định này ([Bài 10](./10_statistics_inference.md)).
3. Khởi tạo trọng số mạng neural thường dùng phân phối chuẩn (Xavier/He initialization — [Bài 20](./20_neural_networks.md)).

## Ví dụ đầy đủ

```python
import numpy as np
from scipy import stats

def describe_distribution(data):
	return {
		"mean": np.mean(data),
		"variance": np.var(data),
		"std": np.std(data),
		"skewness": stats.skew(data),   # độ lệch (bất đối xứng)
		"kurtosis": stats.kurtosis(data),  # độ "nhọn" so với phân phối chuẩn
	}

np.random.seed(42)
normal_data = np.random.normal(loc=50, scale=10, size=1000)
uniform_data = np.random.uniform(low=0, high=100, size=1000)

print("Normal:", describe_distribution(normal_data))
print("Uniform:", describe_distribution(uniform_data))
```

## Bài tập

1. **Kỳ vọng & phương sai**: tính bằng tay kỳ vọng và phương sai của biến ngẫu nhiên $X$ với PMF $P(X=1)=0.2, P(X=2)=0.5, P(X=3)=0.3$, verify bằng code.
2. **Covariance & correlation**: tạo 2 mảng dữ liệu có tương quan dương mạnh, tương quan âm, và không tương quan; tính và so sánh hệ số tương quan của cả 3 cặp.
3. **So sánh phân phối**: mô phỏng (dùng `np.random`) 10,000 mẫu từ Binomial, Poisson, Normal; vẽ histogram của cả 3 (liên hệ [Bài 13](./13_data_visualization_eda.md)), so sánh hình dạng.
4. **Central Limit Theorem**: mô phỏng — lấy trung bình của 30 mẫu từ phân phối Uniform, lặp lại 1000 lần, vẽ histogram của các giá trị trung bình đó — quan sát nó xấp xỉ phân phối chuẩn dù dữ liệu gốc là Uniform.

## Tiếp theo
→ [Bài 10: Ước lượng thống kê & Kiểm định giả thuyết](./10_statistics_inference.md)
