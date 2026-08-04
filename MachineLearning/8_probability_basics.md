# Bài 8: Xác suất cơ bản & Định lý Bayes

## Mục tiêu
- Không gian mẫu, biến cố, xác suất có điều kiện, độc lập.
- Định lý Bayes — nền tảng của Naive Bayes Classifier.

## 1. Khái niệm cơ bản

- **Không gian mẫu ($\Omega$)**: tập hợp mọi kết quả có thể xảy ra. Vd: tung xúc xắc, $\Omega = \{1,2,3,4,5,6\}$.
- **Biến cố (Event)**: 1 tập con của $\Omega$. Vd: "ra số chẵn" = $\{2,4,6\}$.
- **Xác suất $P(A)$**: số đo mức độ "chắc chắn" của biến cố $A$, $0 \leq P(A) \leq 1$.

$$P(A) = \frac{\text{số kết quả thuận lợi cho } A}{\text{tổng số kết quả có thể}} \quad \text{(với không gian mẫu hữu hạn, đồng khả năng)}$$

## 2. Quy tắc cộng & quy tắc nhân xác suất

**Quy tắc cộng:**
$$P(A \cup B) = P(A) + P(B) - P(A \cap B)$$

**Quy tắc nhân (cho biến cố độc lập):**
$$P(A \cap B) = P(A)P(B) \quad \text{nếu } A, B \text{ độc lập}$$

```python
import numpy as np

# Mô phỏng: tung xúc xắc 100,000 lần, ước lượng P(số chẵn)
np.random.seed(42)
rolls = np.random.randint(1, 7, size=100_000)
p_even = np.mean(rolls % 2 == 0)
print(p_even)  # xấp xỉ 0.5
```

## 3. Xác suất có điều kiện (Conditional Probability)

Xác suất của $A$ **biết rằng** $B$ đã xảy ra:

$$P(A|B) = \frac{P(A \cap B)}{P(B)} \quad (\text{với } P(B) > 0)$$

**Ví dụ:** trong 100 email, 30 là spam. Trong 30 spam đó, 27 chứa từ "khuyến mãi". Trong 70 email không spam, chỉ 5 chứa từ đó.

$$P(\text{spam} \mid \text{"khuyến mãi"}) = \frac{P(\text{spam} \cap \text{"khuyến mãi"})}{P(\text{"khuyến mãi"})} = \frac{27/100}{32/100} = 0.84$$

```python
p_spam_and_word = 27/100
p_word = (27+5)/100
p_spam_given_word = p_spam_and_word / p_word
print(p_spam_given_word)  # 0.84
```

## 4. Biến cố độc lập (Independence)

$A$ và $B$ **độc lập** nếu biết $B$ xảy ra không thay đổi xác suất của $A$:

$$P(A|B) = P(A) \quad\Leftrightarrow\quad P(A \cap B) = P(A)P(B)$$

**Ứng dụng ML trực tiếp:** thuật toán **Naive Bayes** ([Bài 19](./19_svm_knn_clustering.md)) giả định các đặc trưng (feature) **độc lập có điều kiện** với nhau khi biết nhãn — giả định "ngây thơ" (naive) này giúp tính toán đơn giản hơn nhiều, dù thực tế các đặc trưng thường không hoàn toàn độc lập.

## 5. Định lý Bayes — công thức quan trọng bậc nhất cho ML

$$P(A|B) = \frac{P(B|A)P(A)}{P(B)}$$

Trong ngữ cảnh ML/thống kê, thường viết dưới dạng:

$$\underbrace{P(\theta|D)}_{\text{posterior}} = \frac{\overbrace{P(D|\theta)}^{\text{likelihood}} \cdot \overbrace{P(\theta)}^{\text{prior}}}{\underbrace{P(D)}_{\text{evidence}}}$$

- **Prior $P(\theta)$**: niềm tin ban đầu về tham số $\theta$ (trước khi thấy dữ liệu).
- **Likelihood $P(D|\theta)$**: xác suất quan sát dữ liệu $D$ nếu tham số là $\theta$.
- **Posterior $P(\theta|D)$**: niềm tin **sau khi** đã thấy dữ liệu — đây là cái ta thực sự muốn biết.

## 6. Ví dụ kinh điển: xét nghiệm y tế (minh họa trực giác Bayes)

Bệnh hiếm gặp: $P(\text{bệnh}) = 0.001$. Xét nghiệm chính xác 99%: $P(\text{dương tính}|\text{bệnh}) = 0.99$, $P(\text{dương tính}|\text{không bệnh}) = 0.01$ (dương tính giả). Hỏi: nếu xét nghiệm dương tính, xác suất THẬT SỰ mắc bệnh là bao nhiêu?

```python
p_disease = 0.001
p_positive_given_disease = 0.99
p_positive_given_no_disease = 0.01
p_no_disease = 1 - p_disease

p_positive = (p_positive_given_disease * p_disease +
              p_positive_given_no_disease * p_no_disease)

p_disease_given_positive = (p_positive_given_disease * p_disease) / p_positive
print(p_disease_given_positive)  # chỉ khoảng 0.09 (9%) — RẤT PHẢN TRỰC GIÁC!
```

Dù test "chính xác 99%", xác suất thật sự mắc bệnh khi dương tính chỉ ~9% — vì bệnh quá hiếm (prior thấp), số dương tính giả (từ nhóm khỏe mạnh đông hơn nhiều) áp đảo số dương tính thật. Đây là bài học quan trọng: **luôn cân nhắc prior**, không chỉ nhìn độ chính xác của "test"/model.

## 7. Naive Bayes Classifier — ứng dụng Bayes trực tiếp trong ML

Để phân loại email spam/không-spam dựa trên các từ xuất hiện $w_1, ..., w_n$:

$$P(\text{spam} \mid w_1,...,w_n) \propto P(\text{spam}) \prod_{i=1}^n P(w_i|\text{spam})$$

(nhờ giả định độc lập ở mục 4, $P(w_1,...,w_n|\text{spam}) \approx \prod P(w_i|\text{spam})$ thay vì phải tính xác suất kết hợp phức tạp)

```python
def naive_bayes_spam_score(words, word_probs_spam, p_spam):
	"""Tính P(spam | words) không chuẩn hóa (bỏ qua evidence P(words))."""
	score = np.log(p_spam)  # dùng log để tránh underflow khi nhân nhiều xác suất nhỏ
	for word in words:
		score += np.log(word_probs_spam.get(word, 1e-6))
	return score
```

Chi tiết cài đặt đầy đủ Naive Bayes Classifier ở [Bài 19](./19_svm_knn_clustering.md).

## Bài tập

1. **Xác suất cơ bản**: mô phỏng tung 2 xúc xắc 100,000 lần bằng NumPy, ước lượng $P(\text{tổng} = 7)$, so sánh với tính bằng tay (đếm số cách ra tổng 7 / 36 khả năng).
2. **Xác suất có điều kiện**: dùng ví dụ email spam ở mục 3, tự tạo 1 bảng số liệu khác, tính $P(\text{spam}|\text{từ khóa})$.
3. **Định lý Bayes — xét nghiệm y tế**: dùng code mẫu mục 6, thử thay đổi $P(\text{bệnh})$ thành 0.01 (bệnh phổ biến hơn) và 0.3, quan sát $P(\text{bệnh}|\text{dương tính})$ thay đổi thế nào — giải thích bằng lời ý nghĩa của prior.
4. **Nâng cao**: viết hàm mô phỏng Monte Carlo (random simulation) để verify kết quả bài tập 3 bằng cách sinh ngẫu nhiên 1 triệu người theo đúng tỷ lệ bệnh/xét nghiệm, đếm trực tiếp thay vì dùng công thức Bayes — 2 kết quả phải xấp xỉ nhau.

## Tiếp theo
→ [Bài 9: Biến ngẫu nhiên & Phân phối xác suất](./9_probability_distributions.md)
