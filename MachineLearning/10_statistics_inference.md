# Bài 10: Ước lượng thống kê & Kiểm định giả thuyết

## Mục tiêu
- Suy luận đầy đủ vì sao MLE dẫn tới MSE (Linear Regression) và Cross-Entropy (Logistic Regression) — TỪNG BƯỚC đại số, không chỉ phát biểu kết quả.
- Hiểu MAP và mối liên hệ với Regularization.
- Nắm ý nghĩa p-value/khoảng tin cậy trước khi dùng scipy.

---

## PHẦN A — Ý NGHĨA TOÁN HỌC

### 1. Maximum Likelihood Estimation (MLE) — nguyên lý "tham số nào khiến dữ liệu ĐÃ QUAN SÁT là hợp lý nhất?"

Cho dữ liệu $x_1,...,x_n$ (giả sử độc lập — [Bài 8 mục 4](./8_probability_basics.md)) sinh ra từ 1 phân phối có tham số $\theta$ chưa biết. **Likelihood** $L(\theta)=P(D|\theta)=\prod_iP(x_i|\theta)$ đo "nếu tham số thật sự là $\theta$, xác suất để quan sát được ĐÚNG bộ dữ liệu này là bao nhiêu". MLE chọn $\theta$ làm cho con số này LỚN NHẤT — nói cách khác, "tham số nào khiến dữ liệu ta đã thấy trở nên ít bất ngờ/hợp lý nhất".

**Vì sao lấy log?** Tích của nhiều xác suất nhỏ ($<1$) dễ gây **underflow số học** (kết quả về gần 0 tới mức máy tính làm tròn thành 0 hẳn) khi $n$ lớn. Vì $\log$ là hàm **đơn điệu tăng**, $\theta$ tối đa hóa $L(\theta)$ CŨNG chính là $\theta$ tối đa hóa $\log L(\theta)$ — chuyển tích thành tổng ($\log\prod_i P(x_i|\theta)=\sum_i\log P(x_i|\theta)$), vừa tránh underflow vừa dễ lấy đạo hàm hơn (đạo hàm của tổng dễ hơn đạo hàm của tích — quy tắc tích ở [Bài 5 mục 3](./5_calculus_derivatives.md) phức tạp hơn quy tắc cộng nhiều).

**Ví dụ tính tay — MLE cho Bernoulli:** tung đồng xu $n$ lần, $k$ lần ngửa. $L(p)=p^k(1-p)^{n-k}$, $\log L(p)=k\log p+(n-k)\log(1-p)$. Đạo hàm theo $p$ ([Bài 5 mục 3](./5_calculus_derivatives.md)) và cho $=0$:

$$\frac{d\log L}{dp} = \frac{k}{p}-\frac{n-k}{1-p}=0 \;\Rightarrow\; k(1-p)=(n-k)p \;\Rightarrow\; p_{MLE}=\frac{k}{n}$$

Kết quả trùng khớp trực giác: ước lượng tốt nhất cho xác suất ngửa chính là **tỷ lệ ngửa quan sát được**.

### 2. MLE ⟹ MSE — chứng minh đầy đủ TẠI SAO Linear Regression dùng MSE

Giả định mô hình sinh dữ liệu: $y=\vec{w}^T\vec{x}+\epsilon$, nhiễu $\epsilon\sim N(0,\sigma^2)$ ([Bài 9 mục 6](./9_probability_distributions.md)). Vậy với $\vec{x}$ cho trước, $y$ tuân theo $N(\vec{w}^T\vec{x},\sigma^2)$ — áp dụng công thức PDF Gaussian ([Bài 9 mục 6](./9_probability_distributions.md)):

$$P(y_i|\vec{x_i},\vec{w}) = \frac{1}{\sigma\sqrt{2\pi}}\exp\left(-\frac{(y_i-\vec{w}^T\vec{x_i})^2}{2\sigma^2}\right)$$

**Bước 1 — viết log-likelihood toàn bộ dataset** (dùng tính độc lập, tổng của log — mục 1):

$$\log L(\vec{w}) = \sum_{i=1}^n\log P(y_i|\vec{x_i},\vec{w}) = \sum_{i=1}^n\left[-\log(\sigma\sqrt{2\pi}) - \frac{(y_i-\vec{w}^T\vec{x_i})^2}{2\sigma^2}\right]$$

**Bước 2 — tách phần phụ thuộc $\vec{w}$:** số hạng $-\log(\sigma\sqrt{2\pi})$ không chứa $\vec{w}$ (hằng số theo $\vec{w}$), nên:

$$\log L(\vec{w}) = -\frac{1}{2\sigma^2}\sum_{i=1}^n(y_i-\vec{w}^T\vec{x_i})^2 + \text{const}$$

**Bước 3 — kết luận:** tối đa hóa $\log L(\vec{w})$ (vì $-\frac{1}{2\sigma^2}<0$) **tương đương** tối thiểu hóa $\sum_i(y_i-\vec{w}^T\vec{x_i})^2$ — **CHÍNH XÁC** là hàm loss MSE (chỉ khác hệ số $\frac1m$ không ảnh hưởng vị trí cực tiểu). Vậy: **MSE không phải lựa chọn tùy tiện** — nó là hệ quả TẤT YẾU của việc giả định nhiễu tuân theo phân phối Gaussian. Nếu nhiễu thực tế KHÔNG Gaussian (vd có nhiều outlier — nhiễu đuôi dày), MSE không còn là lựa chọn "tối ưu về mặt lý thuyết xác suất" nữa — đây là lý do đôi khi dùng MAE (ít nhạy outlier hơn — [Bài 16 mục 8](./16_model_evaluation.md)) thay thế.

### 3. MLE ⟹ Cross-Entropy — chứng minh đầy đủ TẠI SAO Logistic Regression dùng Cross-Entropy

Với bài toán phân loại nhị phân, $y\sim\text{Bernoulli}(p)$ với $p=\sigma(\vec{w}^T\vec{x})$ ([Bài 9 mục 5](./9_probability_distributions.md), [Bài 15](./15_logistic_regression.md)). Viết gọn PMF Bernoulli cho cả 2 trường hợp $y=0,1$ bằng 1 công thức: $P(y_i|\vec{x_i}) = p_i^{y_i}(1-p_i)^{1-y_i}$ (thay $y_i=1$ vào ra $p_i$, thay $y_i=0$ vào ra $1-p_i$ — đúng như định nghĩa Bernoulli).

**Log-likelihood toàn bộ dataset:**

$$\log L(\vec{w}) = \sum_{i=1}^n\left[y_i\log p_i+(1-y_i)\log(1-p_i)\right]$$

**Kết luận:** tối đa hóa biểu thức này tương đương tối thiểu hóa $-\sum_i[y_i\log p_i+(1-y_i)\log(1-p_i)]$ — **CHÍNH XÁC** là Binary Cross-Entropy Loss. Cùng logic suy luận như MSE (mục 2), chỉ khác phân phối giả định (Bernoulli thay vì Gaussian) — đây là hình mẫu chung: **chọn hàm loss = chọn giả định phân phối cho dữ liệu, rồi lấy MLE**.

### 4. Maximum A Posteriori (MAP) — MLE có thêm "niềm tin trước" (prior)

Nhắc lại định lý Bayes ([Bài 8 mục 5](./8_probability_basics.md)): $P(\theta|D)\propto P(D|\theta)P(\theta)$. MAP tối đa hóa **posterior** thay vì chỉ likelihood:

$$\theta_{MAP} = \arg\max_\theta\log P(D|\theta) + \log P(\theta)$$

So với MLE, MAP có thêm số hạng $\log P(\theta)$ — "phạt" các giá trị $\theta$ mà ta tin là ít khả năng xảy ra từ trước (prior).

**Suy luận đầy đủ: MAP với prior Gaussian ⟹ Ridge Regression.** Giả sử tin rằng trọng số nên nhỏ, đặt prior $\vec{w}\sim N(\vec{0},\tau^2I)$ — áp dụng công thức PDF Gaussian đa biến (tương tự mục 2, mở rộng cho vector): $\log P(\vec{w}) = -\frac{1}{2\tau^2}\|\vec{w}\|^2+\text{const}$. Cộng vào log-likelihood MSE (mục 2):

$$\log P(\vec{w}|D) = -\frac{1}{2\sigma^2}\sum_i(y_i-\vec{w}^T\vec{x_i})^2 - \frac{1}{2\tau^2}\|\vec{w}\|^2 + \text{const}$$

Tối đa hóa biểu thức này tương đương tối thiểu hóa $\sum_i(y_i-\vec{w}^T\vec{x_i})^2 + \lambda\|\vec{w}\|^2$ (với $\lambda=\sigma^2/\tau^2$) — **CHÍNH XÁC là hàm loss của Ridge Regression** ([Bài 17](./17_regularization.md))! Vậy regularization L2 không phải "trick kỹ thuật" — nó là **MAP với niềm tin trước rằng trọng số nên nhỏ**. Tương tự, prior Laplace (thay vì Gaussian) lên $\vec{w}$ dẫn tới Lasso Regression (L1) — chứng minh đầy đủ dành làm bài tập.

### 5. Khoảng tin cậy — định lượng độ KHÔNG chắc chắn của ước lượng

Ước lượng $\bar{x}$ (trung bình mẫu) chỉ là 1 con số — nhưng nếu lấy mẫu khác, $\bar{x}$ sẽ khác đi chút ít. Khoảng tin cậy 95% $\bar{x}\pm1.96\cdot\frac{s}{\sqrt{n}}$ định lượng "vùng dao động hợp lý" đó — hệ số $\frac{s}{\sqrt{n}}$ (sai số chuẩn — standard error) GIẢM khi $n$ tăng, phản ánh đúng trực giác: **mẫu càng lớn, ước lượng càng đáng tin cậy** (chính là hệ quả của $\text{Var}(\bar{X})=\text{Var}(X)/n$ khi các mẫu độc lập — [Bài 9 mục 3](./9_probability_distributions.md)).

### 6. Kiểm định giả thuyết & p-value — logic "phản chứng"

Đặt giả thuyết không $H_0$ (thường là "không có khác biệt/hiệu ứng"). Tính **p-value** = xác suất quan sát được kết quả (hoặc cực đoan hơn) NẾU $H_0$ đúng. p-value NHỎ nghĩa là: "nếu $H_0$ thật sự đúng, kết quả ta thấy là RẤT hiếm gặp" — đủ hiếm để nghi ngờ $H_0$ và bác bỏ nó (thường dùng ngưỡng 0.05, mang tính quy ước, không phải chân lý toán học tuyệt đối).

**Cẩn thận diễn giải sai phổ biến:** p-value KHÔNG PHẢI "xác suất $H_0$ đúng" — nó là xác suất của DỮ LIỆU giả sử $H_0$ đúng, một chiều suy luận hoàn toàn khác (giống sự khác biệt $P(A|B)$ và $P(B|A)$ đã nhấn mạnh ở [Bài 8 mục 5](./8_probability_basics.md)).

---

## PHẦN B — Cài đặt & Minh họa bằng code

```python
import numpy as np
from scipy import stats

# Mục 1: verify MLE cho Bernoulli/Normal
np.random.seed(42)
data = np.random.normal(5, 2, size=1000)
print("mu_MLE:", np.mean(data), "sigma_MLE:", np.std(data))  # phải xấp xỉ 5, 2

# Mục 2-3: verify MLE = MSE / Cross-Entropy trực tiếp bằng công thức
def mse_loss(y_true, y_pred):
	return np.mean((y_true - y_pred)**2)

def bce_loss(y_true, p_pred, eps=1e-15):
	p_pred = np.clip(p_pred, eps, 1-eps)
	return -np.mean(y_true*np.log(p_pred) + (1-y_true)*np.log(1-p_pred))

y_true_reg = np.array([3.0, 5.0, 2.5]); y_pred_reg = np.array([2.8, 5.2, 2.3])
print("MSE:", mse_loss(y_true_reg, y_pred_reg))

y_true_clf = np.array([1,0,1,1]); p_pred_clf = np.array([0.9,0.1,0.8,0.6])
print("BCE:", bce_loss(y_true_clf, p_pred_clf))

# Mục 5: khoảng tin cậy
data2 = np.random.normal(50, 10, size=100)
mean, sem = np.mean(data2), stats.sem(data2)
ci = stats.t.interval(confidence=0.95, df=len(data2)-1, loc=mean, scale=sem)
print("Khoảng tin cậy 95%:", ci)

# Mục 6: kiểm định giả thuyết
model_a = [0.85, 0.87, 0.86, 0.84, 0.88]
model_b = [0.80, 0.82, 0.81, 0.79, 0.83]
t_stat, p_value = stats.ttest_ind(model_a, model_b)
print("p-value:", p_value)
print("Bác bỏ H0 (khác biệt có ý nghĩa)" if p_value < 0.05 else "Không đủ bằng chứng")
```

## Bài tập

1. **MLE Bernoulli bằng tay**: cho 100 lần tung, 62 lần ngửa, tự giải phương trình $\frac{d\log L}{dp}=0$ như ở mục 1, verify $p_{MLE}=0.62$.
2. **Tự suy luận lại MLE⟹MSE**: che phần "Bước 1-3" ở mục 2, tự làm lại phép suy luận, so khớp kết quả.
3. **MAP⟹Lasso (nâng cao)**: tự suy luận (tương tự mục 4, nhưng với prior Laplace $P(w)\propto e^{-|w|/b}$) để chứng minh MAP với prior Laplace dẫn tới loss có số hạng $\lambda\|\vec{w}\|_1$ — đây chính là Lasso.
4. **Kiểm định giả thuyết**: thử code mẫu mục 6 với 2 nhóm CÓ khác biệt rõ rệt và 2 nhóm gần như GIỐNG NHAU — so sánh p-value, giải thích bằng lời ý nghĩa từng trường hợp.

## Tổng kết phần Xác suất & Thống kê
Bạn đã suy luận ĐẦY ĐỦ (không chỉ chép công thức) vì sao MSE và Cross-Entropy là lựa chọn "đúng" về mặt lý thuyết xác suất, và vì sao Regularization chính là MAP. Phần cuối của toán nền tảng là Toán rời rạc & Tối ưu hóa — công cụ THỰC SỰ "giải" bài toán tối ưu hàm loss đó.

## Tiếp theo
→ [Bài 11: Toán rời rạc, Độ phức tạp thuật toán & Tối ưu hóa lồi](./11_discrete_math_optimization.md)
