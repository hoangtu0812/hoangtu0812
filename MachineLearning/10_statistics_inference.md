# Bài 10: Ước lượng thống kê & Kiểm định giả thuyết

## Mục tiêu
- Maximum Likelihood Estimation (MLE), Maximum A Posteriori (MAP).
- Khoảng tin cậy, kiểm định giả thuyết, p-value.
- Hiểu vì sao hàm loss MSE/Cross-Entropy đến từ đâu.

## 1. Maximum Likelihood Estimation (MLE) — nguyên lý nền tảng của hầu hết ML

MLE tìm tham số $\theta$ khiến dữ liệu **quan sát được có xác suất cao nhất**:

$$\theta_{MLE} = \arg\max_\theta P(D|\theta) = \arg\max_\theta \prod_{i=1}^n P(x_i|\theta)$$

Thường lấy log để biến tích thành tổng (dễ tối ưu hơn, tránh underflow số học):

$$\theta_{MLE} = \arg\max_\theta \sum_{i=1}^n \log P(x_i|\theta)$$

### Ví dụ: MLE cho phân phối Normal

Cho dữ liệu $x_1, ..., x_n \sim N(\mu, \sigma^2)$, giải $\frac{\partial}{\partial \mu}\log L = 0$ cho ra:

$$\mu_{MLE} = \frac{1}{n}\sum_{i=1}^n x_i \quad (\text{chính là trung bình mẫu!})$$

```python
import numpy as np

np.random.seed(42)
true_mu, true_sigma = 5, 2
data = np.random.normal(true_mu, true_sigma, size=1000)

mu_mle = np.mean(data)   # MLE của mean = trung bình mẫu
sigma_mle = np.std(data)  # MLE của std = độ lệch chuẩn mẫu
print(f"mu ước lượng: {mu_mle:.3f} (thật: {true_mu})")
print(f"sigma ước lượng: {sigma_mle:.3f} (thật: {true_sigma})")
```

## 2. MLE là nguồn gốc của hàm loss MSE — kết nối trực tiếp tới Linear Regression

Giả định Linear Regression: $y = \vec{w}^T\vec{x} + \epsilon$, với nhiễu $\epsilon \sim N(0, \sigma^2)$. Khi đó $y|\vec{x} \sim N(\vec{w}^T\vec{x}, \sigma^2)$.

Log-likelihood của toàn bộ dataset:

$$\log L(\vec{w}) = \sum_{i=1}^n \log P(y_i|\vec{x_i}, \vec{w}) = -\frac{1}{2\sigma^2}\sum_{i=1}^n (y_i - \vec{w}^T\vec{x_i})^2 + \text{const}$$

**Tối đa hóa log-likelihood $\Leftrightarrow$ tối thiểu hóa $\sum (y_i - \hat{y_i})^2$** — chính xác là hàm loss MSE! Đây là lý do sâu xa vì sao Linear Regression dùng MSE: nó **chính là MLE dưới giả định nhiễu Gaussian**, không phải lựa chọn tùy tiện. Chi tiết thực hành ở [Bài 14](./14_linear_regression.md).

## 3. MLE là nguồn gốc của Cross-Entropy Loss — kết nối tới Logistic Regression

Với bài toán phân loại nhị phân, $y \sim \text{Bernoulli}(p)$ với $p = \sigma(\vec{w}^T\vec{x})$ (hàm sigmoid). Log-likelihood:

$$\log L(\vec{w}) = \sum_{i=1}^n \left[y_i\log(p_i) + (1-y_i)\log(1-p_i)\right]$$

**Tối đa hóa log-likelihood $\Leftrightarrow$ tối thiểu hóa $-\sum[y_i\log p_i + (1-y_i)\log(1-p_i)]$** — chính là hàm **Binary Cross-Entropy Loss**, dùng trong Logistic Regression ([Bài 15](./15_logistic_regression.md)) và toàn bộ mạng neural cho bài toán phân loại.

## 4. Maximum A Posteriori (MAP) — MLE + Prior (kết nối Bayes)

MAP tìm $\theta$ tối đa hóa **posterior** thay vì likelihood (liên hệ định lý Bayes — [Bài 8](./8_probability_basics.md)):

$$\theta_{MAP} = \arg\max_\theta P(\theta|D) = \arg\max_\theta P(D|\theta)P(\theta)$$

Nếu chọn prior $P(\vec{w}) \sim N(0, \tau^2)$ (tin rằng trọng số nên nhỏ), MAP cho Linear Regression trở thành:

$$\vec{w}_{MAP} = \arg\min_{\vec{w}} \underbrace{\sum_i(y_i - \vec{w}^T\vec{x_i})^2}_{\text{MLE term (MSE)}} + \underbrace{\lambda\|\vec{w}\|^2}_{\text{prior term}}$$

**Đây chính xác là Ridge Regression (L2 regularization)!** — regularization không phải "trick kỹ thuật" tùy tiện, mà là **MAP với prior Gaussian lên trọng số**. Chi tiết ở [Bài 17](./17_regularization.md).

## 5. Ước lượng khoảng (Confidence Interval)

Khoảng tin cậy 95% cho trung bình:

$$\bar{x} \pm 1.96 \cdot \frac{s}{\sqrt{n}}$$

```python
from scipy import stats

data = np.random.normal(50, 10, size=100)
mean = np.mean(data)
sem = stats.sem(data)  # standard error of the mean

ci = stats.t.interval(confidence=0.95, df=len(data)-1, loc=mean, scale=sem)
print(f"Khoảng tin cậy 95%: {ci}")
```

Ý nghĩa: nếu lặp lại thí nghiệm nhiều lần, ~95% khoảng tính được sẽ chứa giá trị trung bình thật của tổng thể.

## 6. Kiểm định giả thuyết (Hypothesis Testing) & p-value

Quy trình: đặt **giả thuyết không (null hypothesis) $H_0$**, tính xác suất quan sát được kết quả (hoặc cực đoan hơn) **nếu $H_0$ đúng** — đó chính là **p-value**. p-value nhỏ (thường < 0.05) → bác bỏ $H_0$.

```python
from scipy import stats

# So sánh 2 nhóm: model A vs model B (vd accuracy trên nhiều lần chạy)
model_a_scores = [0.85, 0.87, 0.86, 0.84, 0.88]
model_b_scores = [0.80, 0.82, 0.81, 0.79, 0.83]

t_stat, p_value = stats.ttest_ind(model_a_scores, model_b_scores)
print(f"p-value: {p_value:.4f}")

if p_value < 0.05:
	print("Khác biệt có ý nghĩa thống kê — model A thực sự tốt hơn B")
else:
	print("Chưa đủ bằng chứng khác biệt có ý nghĩa")
```

**Ứng dụng ML:** so sánh 2 model (A/B testing), xác định 1 đặc trưng có thực sự ảnh hưởng đến biến mục tiêu hay không, đánh giá cải tiến model có ý nghĩa thống kê hay chỉ là nhiễu ngẫu nhiên.

## Ví dụ đầy đủ: kết nối MLE với hàm loss

```python
import numpy as np

def mse_loss_from_mle(y_true, y_pred):
	"""MSE = -log-likelihood (bỏ hằng số) dưới giả định nhiễu Gaussian."""
	return np.mean((y_true - y_pred)**2)

def bce_loss_from_mle(y_true, p_pred, eps=1e-15):
	"""Binary Cross-Entropy = -log-likelihood dưới giả định Bernoulli."""
	p_pred = np.clip(p_pred, eps, 1-eps)  # tránh log(0)
	return -np.mean(y_true * np.log(p_pred) + (1-y_true) * np.log(1-p_pred))

y_true_reg = np.array([3.0, 5.0, 2.5])
y_pred_reg = np.array([2.8, 5.2, 2.3])
print("MSE:", mse_loss_from_mle(y_true_reg, y_pred_reg))

y_true_clf = np.array([1, 0, 1, 1])
p_pred_clf = np.array([0.9, 0.1, 0.8, 0.6])
print("BCE:", bce_loss_from_mle(y_true_clf, p_pred_clf))
```

## Bài tập

1. **MLE cho Bernoulli**: cho dữ liệu tung đồng xu 100 lần, 62 lần ngửa, tính $p_{MLE}$ bằng công thức (đạo hàm log-likelihood, giải phương trình = 0), so sánh với trực giác (tỷ lệ ngửa).
2. **Kết nối MLE-MSE**: dùng code mẫu mục 2, verify bằng tay rằng tối thiểu hóa MSE = tối đa hóa log-likelihood Gaussian cho 1 bộ dữ liệu nhỏ (3-5 điểm).
3. **Khoảng tin cậy**: tính khoảng tin cậy 95% cho accuracy của 1 model (giả lập bằng list điểm số), giải thích ý nghĩa của khoảng đó.
4. **Kiểm định giả thuyết**: dùng code mẫu mục 6, thử với 2 nhóm dữ liệu có sự khác biệt rất nhỏ (gần như giống hệt nhau) — quan sát p-value lớn, giải thích tại sao KHÔNG bác bỏ được $H_0$ trong trường hợp này.

## Tổng kết phần Xác suất & Thống kê
Bạn đã nắm xác suất, phân phối, MLE/MAP — đây là lý do toán học đằng sau MỌI hàm loss trong ML (MSE, Cross-Entropy) và regularization. Phần cuối của toán nền tảng là Toán rời rạc & Tối ưu hóa — công cụ thực sự "giải" bài toán tối ưu hàm loss đó.

## Tiếp theo
→ [Bài 11: Toán rời rạc, Độ phức tạp thuật toán & Tối ưu hóa lồi](./11_discrete_math_optimization.md)
