# Bài 8: Xác suất cơ bản & Định lý Bayes

## Mục tiêu
- Hiểu xác suất từ định nghĩa gốc, không chỉ công thức.
- Suy luận đầy đủ định lý Bayes, hiểu TẠI SAO nó phản trực giác trong nhiều trường hợp.
- Chỉ sau khi hiểu thấu, mới verify bằng mô phỏng code.

---

## PHẦN A — Ý NGHĨA TOÁN HỌC

### 1. Xác suất là gì? — từ đếm cổ điển tới niềm tin có cập nhật

**Định nghĩa cổ điển** (khi mọi kết quả đồng khả năng): $P(A) = \frac{\text{số kết quả thuận lợi cho }A}{\text{tổng số kết quả có thể}}$. Vd tung xúc xắc, $\Omega=\{1,...,6\}$, $P(\text{số chẵn})=P(\{2,4,6\})=3/6=0.5$.

**Định nghĩa tần suất** (khi lặp lại thí nghiệm nhiều lần): $P(A) \approx \frac{\text{số lần }A\text{ xảy ra}}{\text{tổng số lần thử}}$ khi số lần thử $\to\infty$ — đây là nền tảng để dùng mô phỏng (simulation) verify lý thuyết xác suất, sẽ dùng ở Phần B.

**Ý nghĩa sâu hơn cho ML — xác suất như "mức độ tin tưởng":** trong ML, ta thường nói "xác suất model dự đoán đúng là 80%" dù sự kiện đó không "lặp lại vô hạn lần" theo nghĩa cổ điển — đây là cách diễn giải **Bayesian**: xác suất đo mức độ tin tưởng/không chắc chắn về 1 mệnh đề, có thể **cập nhật** khi có thêm bằng chứng (chính là nội dung định lý Bayes ở mục 5).

### 2. Quy tắc cộng — trực giác từ việc "đếm bị trùng"

$$P(A\cup B) = P(A)+P(B)-P(A\cap B)$$

**Vì sao trừ đi $P(A\cap B)$?** Nếu cộng trực tiếp $P(A)+P(B)$, các kết quả vừa thuộc $A$ vừa thuộc $B$ (phần giao) bị **đếm 2 lần** — phải trừ lại 1 lần để đếm đúng. Hình dung 2 vòng tròn Venn chồng lên nhau: diện tích hợp = diện tích 2 vòng tròn cộng lại, trừ đi phần chồng lấn (đã bị tính 2 lần).

### 3. Xác suất có điều kiện — "thu hẹp không gian mẫu"

$$P(A|B) = \frac{P(A\cap B)}{P(B)}$$

**Ý nghĩa:** một khi biết $B$ đã xảy ra, "thế giới có thể" thu hẹp lại chỉ còn các kết quả nằm trong $B$ — ta hỏi tiếp trong "thế giới thu hẹp" đó, tỷ lệ nào thuộc $A$. Đây là lý do công thức chia cho $P(B)$: "chuẩn hóa lại" không gian mẫu đã thu hẹp về tổng xác suất 1.

**Ví dụ tính tay:** trong 100 email, 30 spam (27 trong đó chứa từ "khuyến mãi"), 70 không spam (5 chứa từ đó). Biết 1 email chứa "khuyến mãi" (thu hẹp không gian mẫu về $32$ email chứa từ này), tỷ lệ nào trong 32 email đó là spam?

$$P(\text{spam}\mid\text{"khuyến mãi"}) = \frac{27/100}{32/100} = \frac{27}{32} = 0.84$$

### 4. Biến cố độc lập — "biết thêm cũng không đổi niềm tin"

$A,B$ độc lập $\Leftrightarrow P(A|B)=P(A) \Leftrightarrow P(A\cap B)=P(A)P(B)$ — biết $B$ xảy ra KHÔNG cung cấp thêm thông tin gì về khả năng $A$ xảy ra.

**Ứng dụng ML quan trọng — Naive Bayes:** giả định các đặc trưng **độc lập có điều kiện** khi biết nhãn, giúp $P(x_1,...,x_n|y) \approx \prod_iP(x_i|y)$ thay vì phải ước lượng 1 phân phối kết hợp cực kỳ phức tạp (số tham số cần ước lượng tăng theo cấp số nhân với số đặc trưng nếu không có giả định này) — "ngây thơ" (naive) vì thực tế các đặc trưng hiếm khi độc lập hoàn toàn, nhưng giả định này vẫn cho kết quả thực dụng tốt trong nhiều bài toán.

### 5. Định lý Bayes — suy luận đầy đủ từ định nghĩa xác suất có điều kiện

Từ mục 3, ta có 2 cách viết $P(A\cap B)$:

$$P(A\cap B) = P(A|B)P(B) = P(B|A)P(A)$$

(2 cách viết này đúng vì $P(A\cap B)$ đối xứng — giao của $A$ và $B$ không phụ thuộc thứ tự). Từ đẳng thức $P(A|B)P(B)=P(B|A)P(A)$, chia cả 2 vế cho $P(B)$:

$$P(A|B) = \frac{P(B|A)P(A)}{P(B)}$$

Đây là **định lý Bayes** — về bản chất chỉ là "đảo chiều" xác suất có điều kiện, cho phép tính $P(A|B)$ từ $P(B|A)$ (thường dễ ước lượng hơn) và ngược lại.

Viết dưới dạng dùng trong ML (ước lượng tham số $\theta$ từ dữ liệu $D$):

$$\underbrace{P(\theta|D)}_{\text{posterior — niềm tin SAU khi thấy data}} = \frac{\overbrace{P(D|\theta)}^{\text{likelihood}}\cdot\overbrace{P(\theta)}^{\text{prior — niềm tin TRƯỚC khi thấy data}}}{\underbrace{P(D)}_{\text{evidence, hằng số chuẩn hóa}}}$$

### 6. Ví dụ kinh điển — TẠI SAO trực giác thường SAI khi bỏ qua prior

Bệnh hiếm: $P(\text{bệnh})=0.001$. Test chính xác 99%: $P(\text{dương}|\text{bệnh})=0.99$, $P(\text{dương}|\text{không bệnh})=0.01$. Hỏi: dương tính thì xác suất THẬT SỰ mắc bệnh?

**Suy luận từng bước bằng công thức xác suất toàn phần** (mở rộng mẫu số $P(B)$ trong định lý Bayes thành tổng theo mọi trường hợp có thể của $A$):

$$P(\text{dương}) = P(\text{dương}|\text{bệnh})P(\text{bệnh}) + P(\text{dương}|\text{không bệnh})P(\text{không bệnh})$$
$$= 0.99\times0.001 + 0.01\times0.999 = 0.00099+0.00999=0.01098$$

Áp dụng Bayes:

$$P(\text{bệnh}|\text{dương}) = \frac{0.99\times0.001}{0.01098} \approx 0.0902 \; (9\%)$$

**Vì sao chỉ 9% dù test "chính xác 99%"?** Vì bệnh quá HIẾM (prior nhỏ) — trong 1000 người, chỉ ~1 người thực sự bệnh (và gần như chắc chắn có kết quả dương), nhưng trong 999 người khỏe mạnh còn lại, vẫn có ~10 người dương tính GIẢ (1% của 999). Vậy trong tổng ~11 người dương tính, chỉ 1 người thực sự bệnh — số dương tính giả (từ nhóm đông đảo người khỏe) áp đảo số dương tính thật (từ nhóm ít ỏi người bệnh). Bài học: **luôn cân nhắc prior**, không chỉ nhìn độ chính xác của "test"/model — sai lầm bỏ qua prior này rất phổ biến khi diễn giải kết quả model phân loại trong thực tế (vd hệ thống phát hiện gian lận, hiếm khi xảy ra).

### 7. Naive Bayes Classifier — áp dụng trực tiếp mọi thứ vừa học

Để phân loại email dựa trên các từ $w_1,...,w_n$ xuất hiện, áp dụng Bayes (mục 5) và giả định độc lập (mục 4):

$$P(\text{spam}|w_1,...,w_n) = \frac{P(w_1,...,w_n|\text{spam})P(\text{spam})}{P(w_1,...,w_n)} \propto P(\text{spam})\prod_{i=1}^nP(w_i|\text{spam})$$

Ký hiệu $\propto$ (tỷ lệ thuận) vì mẫu số $P(w_1,...,w_n)$ giống nhau cho mọi lớp cần so sánh (spam hay không spam) — không ảnh hưởng tới việc lớp nào có xác suất CAO HƠN, nên có thể bỏ qua khi chỉ cần so sánh/phân loại. Cài đặt đầy đủ ở [Bài 19](./19_svm_knn_clustering.md).

---

## PHẦN B — Cài đặt & Minh họa bằng code

```python
import numpy as np

# Mục 1: verify định nghĩa tần suất bằng mô phỏng
np.random.seed(42)
rolls = np.random.randint(1, 7, size=100_000)
print(np.mean(rolls % 2 == 0))  # xấp xỉ 0.5, khớp định nghĩa cổ điển

# Mục 3: verify xác suất có điều kiện
p_spam_and_word = 27/100
p_word = 32/100
print(p_spam_and_word / p_word)  # 0.84, khớp tính tay

# Mục 6: verify bài toán xét nghiệm y tế bằng công thức
def bayes_medical_test(p_disease, p_pos_given_disease, p_pos_given_no_disease):
	p_no_disease = 1 - p_disease
	p_positive = (p_pos_given_disease*p_disease + p_pos_given_no_disease*p_no_disease)
	return (p_pos_given_disease * p_disease) / p_positive

print(bayes_medical_test(0.001, 0.99, 0.01))  # xấp xỉ 0.09 — khớp tính tay

# Verify LẦN NỮA bằng mô phỏng Monte Carlo — cách khác hoàn toàn để tới cùng đáp số
def monte_carlo_verify(p_disease, p_pos_given_disease, p_pos_given_no_disease, n=1_000_000):
	has_disease = np.random.rand(n) < p_disease
	positive = np.where(
		has_disease,
		np.random.rand(n) < p_pos_given_disease,
		np.random.rand(n) < p_pos_given_no_disease,
	)
	return has_disease[positive].mean()  # trong số dương tính, tỷ lệ thực sự bệnh

print(monte_carlo_verify(0.001, 0.99, 0.01))  # phải xấp xỉ 0.09, khớp công thức Bayes
```

## Bài tập

1. **Xác suất cổ điển vs mô phỏng**: mô phỏng tung 2 xúc xắc 100,000 lần, ước lượng $P(\text{tổng}=7)$ bằng tần suất, verify bằng tính tay (đếm số cách/36).
2. **Xác suất có điều kiện tự đặt**: tạo 1 bảng số liệu khác (vd chẩn đoán bệnh dựa trên triệu chứng), tính $P(\text{bệnh}|\text{triệu chứng})$ bằng tay trước, verify bằng code.
3. **Định lý Bayes — thử đổi prior**: dùng `bayes_medical_test`, thử $P(\text{bệnh})=0.01$ và $0.3$, giải thích bằng lời (dựa vào lập luận mục 6) tại sao $P(\text{bệnh}|\text{dương})$ tăng mạnh khi bệnh phổ biến hơn.
4. **Monte Carlo verify**: dùng `monte_carlo_verify`, xác nhận nó cho kết quả gần khớp `bayes_medical_test` — giải thích bằng lời tại sao 2 cách tiếp cận hoàn toàn khác nhau (công thức đóng vs mô phỏng ngẫu nhiên) lại cho cùng đáp số.

## Tiếp theo
→ [Bài 9: Biến ngẫu nhiên & Phân phối xác suất](./9_probability_distributions.md)
