# Bài 22: Optimizer Nâng Cao & Regularization Trong Deep Learning

## Mục tiêu
- Hiểu vì sao Gradient Descent thuần ([Bài 11](./11_discrete_math_optimization.md)) chưa đủ tốt cho Deep Learning — và ý nghĩa toán học của từng cải tiến.
- Nắm Dropout, Batch Normalization, Early Stopping — không chỉ "dùng vì nó hiệu quả" mà hiểu cơ chế.

---

## PHẦN A — Ý NGHĨA TOÁN HỌC CỦA CÁC OPTIMIZER

### 1. Vấn đề của SGD/Mini-batch GD thuần trên bề mặt loss không lồi

Hàm loss của mạng neural **không lồi** ([Bài 7 mục 3](./7_calculus_optimization.md), [Bài 11 mục 4](./11_discrete_math_optimization.md)) — bề mặt của nó gồ ghề, có nhiều saddle point ([Bài 7 mục 1-2](./7_calculus_optimization.md)) và "thung lũng hẹp, dốc đứng 2 bên nhưng thoải dọc đáy" (ravine). Gradient Descent thuần với learning rate cố định gặp 2 vấn đề: (1) dao động qua lại ở thung lũng hẹp thay vì tiến nhanh dọc đáy, (2) "mắc kẹt" hoặc di chuyển cực chậm ở gần saddle point (nơi gradient gần 0 theo nhiều hướng).

### 2. Momentum — mượn ý tưởng "quán tính" từ vật lý

$$\vec{v}_t = \beta\vec{v}_{t-1} + (1-\beta)\nabla L(\vec{w}_t), \qquad \vec{w}_{t+1} = \vec{w}_t - \eta\vec{v}_t$$

**Ý nghĩa toán học:** $\vec{v}_t$ là trung bình động (exponential moving average) của các gradient GẦN ĐÂY, không chỉ gradient TẠI THỜI ĐIỂM HIỆN TẠI. Khai triển đệ quy: $\vec{v}_t = (1-\beta)\sum_{k=0}^t\beta^{t-k}\nabla L(\vec{w}_k)$ — gradient càng CŨ đóng góp trọng số càng nhỏ (giảm theo cấp số nhân $\beta^{t-k}$), gradient GẦN ĐÂY đóng góp nhiều nhất.

**Vì sao giúp ích:** nếu gradient liên tục "dao động qua lại" theo 1 hướng (vd 2 thành thung lũng hẹp) nhưng NHẤT QUÁN tiến theo hướng khác (dọc đáy thung lũng), trung bình động sẽ "triệt tiêu" phần dao động qua lại (cộng dồn các dấu +/- xen kẽ gần bằng 0) trong khi CỘNG DỒN phần nhất quán (cùng dấu liên tục) — kết quả: di chuyển mượt hơn, nhanh hơn dọc hướng đúng, giảm dao động ngang. Đây chính là "quán tính" — giống 1 quả bóng lăn xuống thung lũng, tích lũy đà theo hướng dốc chính, không đổi hướng đột ngột mỗi khi gặp gợn nhỏ.

### 3. RMSProp — tự động điều chỉnh learning rate theo TỪNG tham số

$$s_t = \beta s_{t-1} + (1-\beta)(\nabla L(\vec{w}_t))^2, \qquad \vec{w}_{t+1} = \vec{w}_t - \frac{\eta}{\sqrt{s_t}+\epsilon}\nabla L(\vec{w}_t)$$

**Ý nghĩa:** $s_t$ là trung bình động của **bình phương** gradient (ước lượng độ lớn/mức "biến động" của gradient theo từng tham số riêng — liên hệ ý nghĩa phương sai ở [Bài 9 mục 3](./9_probability_distributions.md), $E[X^2]$ là 1 phần của công thức phương sai). Chia learning rate cho $\sqrt{s_t}$: tham số nào có gradient LỚN/biến động MẠNH liên tục sẽ được "hãm bớt" tốc độ cập nhật (chia cho số lớn hơn); tham số nào có gradient nhỏ/ổn định được "tăng tốc" tương đối (chia cho số nhỏ hơn). Điều này giải quyết vấn đề: các tham số khác nhau trong mạng neural thường cần tốc độ học khác nhau, mà 1 learning rate cố định cho TẤT CẢ tham số ([Bài 11](./11_discrete_math_optimization.md)) không tối ưu.

### 4. Adam — kết hợp Momentum + RMSProp, optimizer mặc định phổ biến nhất

$$\vec{m}_t = \beta_1\vec{m}_{t-1}+(1-\beta_1)\nabla L, \qquad \vec{s}_t = \beta_2\vec{s}_{t-1}+(1-\beta_2)(\nabla L)^2$$

Cộng thêm bước "hiệu chỉnh thiên lệch" (bias correction) — vì $\vec{m}_0=\vec{s}_0=\vec{0}$, những bước đầu $\vec{m}_t,\vec{s}_t$ bị lệch thấp so với giá trị trung bình thật (do khởi tạo từ 0), nên chia lại:

$$\hat{\vec{m}}_t = \frac{\vec{m}_t}{1-\beta_1^t}, \qquad \hat{\vec{s}}_t = \frac{\vec{s}_t}{1-\beta_2^t}, \qquad \vec{w}_{t+1}=\vec{w}_t-\frac{\eta}{\sqrt{\hat{\vec{s}}_t}+\epsilon}\hat{\vec{m}}_t$$

Adam = "Momentum cho HƯỚNG di chuyển" ($\hat{\vec{m}}_t$) + "RMSProp cho TỐC ĐỘ theo từng tham số" ($\hat{\vec{s}}_t$) — kết hợp cả 2 lợi ích ở mục 2-3, là lý do nó thường là lựa chọn optimizer MẶC ĐỊNH đầu tiên khi thử nghiệm 1 kiến trúc mạng mới.

### 5. Learning Rate Scheduling — vì sao learning rate nên GIẢM DẦN theo thời gian

Nhắc lại [Bài 7 mục 5](./7_calculus_optimization.md): learning rate lớn giúp hội tụ NHANH lúc đầu (khi còn xa cực tiểu, xấp xỉ Taylor bậc 1 vẫn khá chính xác trên khoảng cách lớn), nhưng GẦN cực tiểu, bước nhảy lớn dễ "nhảy qua" đáy rồi dao động qua lại mãi không ổn định chính xác. Learning rate scheduling (vd giảm dần theo epoch, hoặc giảm khi loss ngừng cải thiện) tận dụng learning rate lớn để tiến nhanh giai đoạn đầu, rồi thu nhỏ dần để "định vị chính xác" đáy cực tiểu ở giai đoạn cuối.

---

## PHẦN B — Ý NGHĨA & CÀI ĐẶT REGULARIZATION CHO DEEP LEARNING

### 6. Dropout — "buộc" mạng không phụ thuộc quá mức vào 1 neuron cụ thể

Trong lúc huấn luyện, MỖI bước, ngẫu nhiên "tắt" (đặt = 0) 1 tỷ lệ $p$ neuron ở mỗi lớp — mạng buộc phải học biểu diễn KHÔNG phụ thuộc vào bất kỳ tổ hợp neuron cụ thể nào (vì tổ hợp đó có thể bị tắt bất cứ lúc nào), giống việc huấn luyện đồng thời NHIỀU mạng con (kiến trúc khác nhau do neuron bị tắt khác nhau mỗi bước) rồi lấy trung bình — về bản chất là 1 dạng **ensemble ngầm** ([Bài 18 mục 6](./18_trees_ensembles.md), cùng logic Bagging: kết hợp nhiều model giảm variance).

```python
def dropout_forward(a, keep_prob=0.8, training=True):
	if not training:
		return a  # KHÔNG dropout lúc test/inference — dùng toàn bộ mạng
	mask = (np.random.rand(*a.shape) < keep_prob) / keep_prob  # chia lại để giữ kỳ vọng không đổi
	return a * mask
```

Chia cho `keep_prob` ("inverted dropout") đảm bảo $E[\text{output sau dropout}]$ = output gốc ([Bài 9 mục 2](./9_probability_distributions.md), tính tuyến tính kỳ vọng) — để không cần điều chỉnh gì thêm lúc test.

### 7. Batch Normalization — chuẩn hóa dữ liệu GIỮA các lớp, không chỉ ở input

Nhắc lại [Bài 9 mục 3](./9_probability_distributions.md): chuẩn hóa input giúp Gradient Descent ổn định hơn. Batch Norm áp dụng Ý TƯỞNG TƯƠNG TỰ nhưng cho đầu ra của MỌI lớp ẩn, không chỉ input ban đầu — vì phân phối của $\vec{z}^{[l]}$ thay đổi liên tục trong lúc huấn luyện (do các lớp trước cũng đang thay đổi trọng số — hiện tượng gọi là "internal covariate shift"), gây khó khăn cho việc học ổn định ở các lớp sau:

$$\hat{z} = \frac{z-\mu_{batch}}{\sqrt{\sigma_{batch}^2+\epsilon}}, \qquad y = \gamma\hat{z}+\beta$$

$\mu_{batch},\sigma_{batch}$ tính TRÊN TỪNG MINI-BATCH ([Bài 11 mục 5](./11_discrete_math_optimization.md)) lúc huấn luyện; $\gamma,\beta$ là tham số HỌC ĐƯỢC, cho phép mạng tự "hoàn tác" việc chuẩn hóa nếu thực sự cần thiết (linh hoạt hơn ép cứng mean=0,std=1).

### 8. Early Stopping — dừng huấn luyện dựa trên validation loss ([Bài 16 mục 2](./16_model_evaluation.md))

Theo dõi loss trên tập validation mỗi epoch — dừng huấn luyện khi validation loss NGỪNG cải thiện (hoặc bắt đầu tăng) dù train loss vẫn giảm — đây chính là dấu hiệu overfitting bắt đầu xảy ra ([Bài 16 mục 4](./16_model_evaluation.md), [Bài 17 mục 1](./17_regularization.md)). Về bản chất, Early Stopping là 1 dạng regularization "ngầm": giới hạn SỐ BƯỚC cập nhật tham số, tương tự cách L2 regularization ([Bài 17 mục 3](./17_regularization.md)) giới hạn ĐỘ LỚN tham số — cả 2 đều giảm "không gian giả thuyết hiệu dụng" mà model có thể khai thác để khớp nhiễu.

---

## PHẦN C — Cài đặt & Minh họa bằng code

```python
import numpy as np

class AdamOptimizer:
	def __init__(self, params_shapes, lr=0.001, beta1=0.9, beta2=0.999, eps=1e-8):
		self.lr, self.beta1, self.beta2, self.eps = lr, beta1, beta2, eps
		self.m = [np.zeros(s) for s in params_shapes]
		self.s = [np.zeros(s) for s in params_shapes]
		self.t = 0

	def step(self, params, grads):
		self.t += 1
		for i in range(len(params)):
			self.m[i] = self.beta1*self.m[i] + (1-self.beta1)*grads[i]
			self.s[i] = self.beta2*self.s[i] + (1-self.beta2)*(grads[i]**2)

			m_hat = self.m[i] / (1 - self.beta1**self.t)   # hiệu chỉnh thiên lệch — mục 4
			s_hat = self.s[i] / (1 - self.beta2**self.t)

			params[i] -= self.lr * m_hat / (np.sqrt(s_hat) + self.eps)
		return params

# So sánh trực quan: SGD thuần vs Momentum vs Adam trên hàm loss có "ravine" (thung lũng hẹp)
def compare_optimizers(f, grad_f, start, n_steps=100):
	# hàm ravine: dốc mạnh theo x, thoải theo y — mô phỏng "thung lũng hẹp" mục 1
	results = {}

	w = start.copy()  # SGD thuần
	path_sgd = [w.copy()]
	for _ in range(n_steps):
		w = w - 0.01 * grad_f(w)
		path_sgd.append(w.copy())
	results['SGD'] = np.array(path_sgd)

	w = start.copy(); v = np.zeros_like(w)  # Momentum
	path_mom = [w.copy()]
	for _ in range(n_steps):
		v = 0.9*v + 0.1*grad_f(w)
		w = w - 0.01 * v
		path_mom.append(w.copy())
	results['Momentum'] = np.array(path_mom)

	return results

f_ravine = lambda w: 10*w[0]**2 + w[1]**2  # dốc mạnh theo x (hệ số 10), thoải theo y
grad_ravine = lambda w: np.array([20*w[0], 2*w[1]])

results = compare_optimizers(f_ravine, grad_ravine, np.array([1.0, 1.0]))
print("SGD bước cuối:", results['SGD'][-1])
print("Momentum bước cuối:", results['Momentum'][-1])  # thường gần (0,0) hơn với cùng số bước
```

## Bài tập

1. **Cài đặt Momentum & RMSProp**: viết class tương tự `AdamOptimizer` nhưng chỉ Momentum, và chỉ RMSProp — so sánh cả 3 (SGD/Momentum/RMSProp/Adam) trên cùng hàm "ravine" ở ví dụ cuối bài, vẽ quỹ đạo hội tụ (liên hệ [Bài 13](./13_data_visualization_eda.md)).
2. **Verify Dropout giữ kỳ vọng không đổi**: chạy `dropout_forward` nhiều lần (vd 10,000 lần) trên cùng input, tính trung bình output — verify nó xấp xỉ input gốc, đúng lý luận ở mục 6.
3. **Tích hợp Adam vào `NeuralNetworkFromScratch`** ([Bài 21](./21_backpropagation.md)): thay bước `update()` đơn giản bằng `AdamOptimizer`, so sánh tốc độ hội tụ (loss theo epoch) với Gradient Descent thuần.
4. **Early Stopping**: huấn luyện 1 mạng neural trên dataset dễ overfit (ít mẫu, mạng lớn), theo dõi train loss và validation loss mỗi epoch, tự cài logic dừng khi validation loss không cải thiện sau $k$ epoch liên tiếp.

## Tiếp theo
→ [Bài 23: CNN (Convolutional Neural Networks)](./23_cnn.md)
