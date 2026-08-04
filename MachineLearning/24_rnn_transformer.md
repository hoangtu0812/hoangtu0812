# Bài 24: RNN/LSTM & Giới Thiệu Attention/Transformer

## Mục tiêu
- Hiểu vì sao dữ liệu tuần tự (chuỗi thời gian, văn bản) cần kiến trúc khác CNN/MLP.
- Nắm vanishing gradient trong RNN — hệ quả TRỰC TIẾP của chain rule đã học ([Bài 6](./6_calculus_matrix_calculus.md)), và cách LSTM giải quyết.
- Hiểu ý tưởng cốt lõi của Attention trước khi xem kiến trúc Transformer.

---

## PHẦN A — Ý NGHĨA TOÁN HỌC

### 1. Vì sao MLP/CNN không phù hợp cho dữ liệu tuần tự?

MLP ([Bài 20](./20_neural_networks.md)) yêu cầu input kích thước CỐ ĐỊNH — không xử lý được câu văn có độ dài thay đổi. Quan trọng hơn: MLP/CNN không có khái niệm "thứ tự" — nếu hoán đổi vị trí 2 từ trong câu, MLP xử lý y hệt (không phân biệt được), trong khi thứ tự từ MANG THÔNG TIN QUAN TRỌNG ("chó cắn người" khác "người cắn chó"). Cần 1 kiến trúc có "bộ nhớ" giữ lại thông tin từ các bước TRƯỚC khi xử lý bước HIỆN TẠI.

### 2. Recurrent Neural Network (RNN) — chia sẻ trọng số THEO THỜI GIAN

$$\vec{h}_t = \phi(W_{hh}\vec{h}_{t-1} + W_{xh}\vec{x}_t + \vec{b})$$

**Ý nghĩa:** $\vec{h}_t$ (hidden state) đóng vai trò "bộ nhớ" — tóm tắt thông tin từ MỌI bước trước đó ($\vec{x}_1,...,\vec{x}_t$) thành 1 vector kích thước cố định. Cùng bộ trọng số $W_{hh},W_{xh}$ được dùng LẶP LẠI ở MỌI bước thời gian — đây là "parameter sharing theo thời gian", cùng ý tưởng chia sẻ tham số theo không gian của CNN ([Bài 23 mục 3](./23_cnn.md)), giờ áp dụng theo chiều thời gian: giả định "quy luật xử lý bước $t$ giống hệt quy luật xử lý bước $t-1,t-2,...$" (translation invariance theo thời gian).

### 3. Vanishing/Exploding Gradient trong RNN — hệ quả TOÁN HỌC trực tiếp từ Chain Rule

Huấn luyện RNN dùng **Backpropagation Through Time (BPTT)** — về bản chất chỉ là Backpropagation thường ([Bài 21](./21_backpropagation.md)) áp dụng lên mạng "trải phẳng theo thời gian" (coi mỗi bước thời gian như 1 lớp riêng biệt, dùng CHUNG trọng số).

Xét ảnh hưởng của $\vec{h}_1$ (bước đầu) lên loss tại bước $T$ (bước cuối, T rất lớn với chuỗi dài) — áp dụng chain rule liên tiếp qua $T-1$ bước ([Bài 6 mục 5](./6_calculus_matrix_calculus.md)):

$$\frac{\partial\vec{h}_T}{\partial\vec{h}_1} = \prod_{t=2}^{T}\frac{\partial\vec{h}_t}{\partial\vec{h}_{t-1}} \approx \prod_{t=2}^{T}W_{hh}^T\text{diag}(\phi'(\cdot))$$

**Vấn đề:** đây là TÍCH của $T-1$ ma trận Jacobian gần giống nhau (vì CÙNG $W_{hh}$ mọi bước — mục 2). Nếu "độ lớn" của $W_{hh}$ (đo bằng trị riêng lớn nhất — [Bài 4 mục 1](./4_linear_algebra_eigen_svd.md)) nhỏ hơn 1, tích $T-1$ số nhỏ hơn 1 liên tiếp tiến về 0 RẤT NHANH (theo cấp số nhân — **vanishing gradient**, giống hiện tượng đã thấy với sigmoid ở [Bài 20 mục 3](./20_neural_networks.md), nhưng ở đây trầm trọng hơn vì lặp lại CÙNG 1 ma trận hàng trăm/nghìn lần thay vì vài chục lớp). Ngược lại, nếu độ lớn $>1$, tích bùng nổ theo cấp số nhân (**exploding gradient**). Đây là lý do RNN "vanilla" học RẤT KÉM các phụ thuộc XA (long-range dependency) — thông tin từ đầu chuỗi gần như "biến mất" khi lan truyền gradient qua chuỗi dài.

### 4. LSTM — thiết kế "đường tắt" (Cell State) để gradient không bị "co lại" theo cấp số nhân

Ý tưởng cốt lõi: thay vì gradient phải đi qua PHÉP NHÂN ma trận LẶP LẠI (nguồn gốc vanishing ở mục 3), LSTM thêm 1 đường truyền phụ — **cell state** $\vec{c}_t$ — cập nhật chủ yếu bằng PHÉP CỘNG:

$$\vec{c}_t = \vec{f}_t\odot\vec{c}_{t-1} + \vec{i}_t\odot\tilde{\vec{c}}_t$$

với $\vec{f}_t$ (**forget gate**), $\vec{i}_t$ (**input gate**) $\in(0,1)$ (tính bằng sigmoid — [Bài 20 mục 3](./20_neural_networks.md), đóng vai trò "van điều tiết" bao nhiêu % thông tin cũ được GIỮ và bao nhiêu % thông tin mới được THÊM VÀO).

**Vì sao phép CỘNG giải quyết vanishing gradient?** Đạo hàm của $\vec{c}_t$ theo $\vec{c}_{t-1}$ chủ yếu là $\vec{f}_t$ (từ số hạng $\vec{f}_t\odot\vec{c}_{t-1}$, đạo hàm của phép nhân element-wise theo 1 trong 2 thừa số) — nếu forget gate $\vec{f}_t\approx1$ (mạng "quyết định" giữ lại thông tin), gradient truyền qua gần như NGUYÊN VẸN (nhân với số gần 1, không co lại theo cấp số nhân như nhân lặp lại với $W_{hh}$ ở mục 3). Đây là "đường cao tốc gradient" (gradient highway) — cho phép LSTM học được phụ thuộc XA hơn nhiều so với RNN thuần.

### 5. Ý tưởng cốt lõi của Attention — "nhìn lại TOÀN BỘ chuỗi, không chỉ hidden state cuối"

RNN/LSTM nén TOÀN BỘ lịch sử vào 1 vector $\vec{h}_t$ kích thước cố định — với chuỗi rất dài, đây là **nút thắt cổ chai thông tin** (information bottleneck), dù có cell state cũng khó giữ trọn vẹn thông tin chi tiết từ rất xa. Attention giải quyết bằng cách: khi xử lý 1 vị trí, cho phép mô hình **"nhìn lại" TRỰC TIẾP mọi vị trí khác** trong chuỗi (không phải chỉ qua trung gian hidden state), và tính 1 trọng số "chú ý" cho biết vị trí nào quan trọng nhất với vị trí hiện tại.

**Công thức Scaled Dot-Product Attention** — dùng CHÍNH XÁC dot product đã học ở [Bài 2 mục 3](./2_linear_algebra_vectors.md) để đo "mức độ liên quan":

$$\text{Attention}(Q,K,V) = \text{softmax}\left(\frac{QK^T}{\sqrt{d_k}}\right)V$$

- $Q$ (Query — "tôi đang tìm thông tin gì"), $K$ (Key — "mỗi vị trí khác đại diện cho thông tin gì"), $V$ (Value — "nội dung thực sự của mỗi vị trí") — đều là các phép biến đổi tuyến tính ([Bài 3 mục 1](./3_linear_algebra_matrices.md)) của input.
- $QK^T$: dot product giữa Query và MỌI Key — đo độ "khớp" giữa vị trí đang xét và mọi vị trí khác (đúng ý nghĩa dot product ở [Bài 2 mục 3](./2_linear_algebra_vectors.md): điểm càng cao, 2 vector càng "cùng hướng"/liên quan).
- Chia $\sqrt{d_k}$: chuẩn hóa để giá trị dot product không "nổ" quá lớn khi số chiều $d_k$ tăng (dot product của 2 vector ngẫu nhiên có phương sai tỷ lệ thuận với số chiều — liên hệ [Bài 9 mục 2](./9_probability_distributions.md), tính tổng phương sai của tổng các biến độc lập), tránh đẩy softmax vào vùng bão hòa gradient gần 0 ([Bài 20 mục 3](./20_neural_networks.md)).
- `softmax` ([Bài 15 mục 8](./15_logistic_regression.md)): chuẩn hóa điểm khớp thành trọng số "chú ý" có tổng $=1$ — 1 phân phối xác suất trên các vị trí.
- Nhân với $V$: kết quả cuối là **trung bình có trọng số** (weighted average — mở rộng ý tưởng kỳ vọng ở [Bài 9 mục 2](./9_probability_distributions.md)) của mọi Value, trọng số cao cho vị trí "liên quan nhiều" tới Query hiện tại.

### 6. Transformer — bỏ hẳn RNN, chỉ dùng Attention

Kiến trúc Transformer (Vaswani et al., 2017) thay thế HOÀN TOÀN cấu trúc tuần tự của RNN bằng **Self-Attention** (Q,K,V đều từ CÙNG 1 chuỗi input) xếp chồng nhiều lớp, cộng thêm **Positional Encoding** (vì Attention tự nó không phân biệt thứ tự — không giống RNN xử lý tuần tự tự nhiên có thứ tự, cần "tiêm" thêm thông tin vị trí vào vector input). Ưu điểm lớn nhất: mọi vị trí trong chuỗi có thể "nhìn" trực tiếp mọi vị trí khác CÙNG LÚC (tính toán song song được hoàn toàn, không phải tuần tự từng bước thời gian như RNN ở mục 2) — vừa giải quyết triệt để vanishing gradient (không còn tích lặp lại ma trận qua nhiều bước như mục 3), vừa tận dụng được phần cứng song song (GPU) hiệu quả hơn nhiều, là nền tảng của các mô hình ngôn ngữ lớn hiện đại.

---

## PHẦN B — Cài đặt & Minh họa bằng code

```python
import numpy as np

def sigmoid(z): return 1 / (1 + np.exp(-z))
def tanh(z): return np.tanh(z)

# Mục 2: RNN cell đơn giản — verify công thức forward
def rnn_cell_forward(x_t, h_prev, W_hh, W_xh, b):
	return tanh(W_hh @ h_prev + W_xh @ x_t + b)

# Mục 3: minh họa vanishing gradient bằng số — verify hiện tượng thay vì chỉ khẳng định
def demonstrate_vanishing_gradient(W_hh_scale, T=50):
	np.random.seed(42)
	W_hh = np.random.randn(4, 4) * W_hh_scale
	grad = np.eye(4)
	norms = []
	for t in range(T):
		grad = W_hh.T @ grad  # tích lặp lại như công thức mục 3
		norms.append(np.linalg.norm(grad))
	return norms

norms_small = demonstrate_vanishing_gradient(W_hh_scale=0.3)  # trị riêng nhỏ -> vanish
norms_large = demonstrate_vanishing_gradient(W_hh_scale=1.5)  # trị riêng lớn -> explode
print("Gradient norm sau 50 bước (scale nhỏ):", norms_small[-1])   # gần 0
print("Gradient norm sau 50 bước (scale lớn):", norms_large[-1])   # rất lớn hoặc inf

# Mục 5: Scaled Dot-Product Attention from scratch
def softmax(z):
	exp_z = np.exp(z - np.max(z, axis=-1, keepdims=True))
	return exp_z / np.sum(exp_z, axis=-1, keepdims=True)

def scaled_dot_product_attention(Q, K, V):
	d_k = Q.shape[-1]
	scores = Q @ K.T / np.sqrt(d_k)   # dot product từng cặp Query-Key, chuẩn hóa theo mục 5
	weights = softmax(scores)           # chuẩn hóa thành trọng số "chú ý"
	output = weights @ V                 # trung bình có trọng số của Value
	return output, weights

np.random.seed(42)
seq_len, d_model = 4, 8
Q = np.random.randn(seq_len, d_model)
K = np.random.randn(seq_len, d_model)
V = np.random.randn(seq_len, d_model)

output, attn_weights = scaled_dot_product_attention(Q, K, V)
print("Attention weights (mỗi dòng tổng = 1):\n", attn_weights.sum(axis=1))  # verify [1,1,1,1]
```

### LSTM thực tế với PyTorch

```python
import torch
import torch.nn as nn

class SimpleLSTMClassifier(nn.Module):
	def __init__(self, vocab_size, embed_dim, hidden_dim, num_classes):
		super().__init__()
		self.embedding = nn.Embedding(vocab_size, embed_dim)
		self.lstm = nn.LSTM(embed_dim, hidden_dim, batch_first=True)  # tự cài forget/input/output gate (mục 4)
		self.fc = nn.Linear(hidden_dim, num_classes)

	def forward(self, x):
		embedded = self.embedding(x)
		_, (h_n, c_n) = self.lstm(embedded)   # h_n: hidden state cuối, c_n: cell state cuối (mục 4)
		return self.fc(h_n.squeeze(0))
```

## Bài tập

1. **RNN cell forward**: cài `rnn_cell_forward` như trên, chạy qua 1 chuỗi 5 bước, in `h_t` mỗi bước.
2. **Verify vanishing/exploding gradient bằng số**: chạy `demonstrate_vanishing_gradient` với nhiều `W_hh_scale` (0.1, 0.5, 0.9, 1.1, 2.0), vẽ đồ thị `norm` theo số bước $t$ (liên hệ [Bài 13](./13_data_visualization_eda.md)) — xác định ngưỡng scale nào bắt đầu gây vanish/explode rõ rệt.
3. **Attention weights**: dùng `scaled_dot_product_attention`, thử với $Q=K$ (self-attention), quan sát vị trí nào "chú ý" tới chính nó nhiều nhất — giải thích bằng lời dựa vào ý nghĩa dot product.
4. **So sánh RNN vs LSTM trên chuỗi dài**: dùng PyTorch, huấn luyện `nn.RNN` và `nn.LSTM` trên bài toán cần nhớ thông tin từ ĐẦU chuỗi dài (vd "chuỗi có bắt đầu bằng token X hay không", chuỗi dài 100+ bước) — so sánh accuracy, minh chứng thực nghiệm cho lý luận vanishing gradient ở mục 3-4.

## Tổng kết Phần III — Deep Learning
Bạn đã đi từ neuron đơn lẻ (Perceptron) tới Backpropagation đầy đủ, các optimizer nâng cao, CNN cho dữ liệu ảnh, và RNN/LSTM/Attention cho dữ liệu tuần tự — mọi công thức đều được suy luận từ nền tảng Toán ở Phần I. Bước cuối cùng là ghép TẤT CẢ kiến thức thành 1 dự án hoàn chỉnh.

## Tiếp theo
→ [Dự án Capstone: Pipeline Machine Learning hoàn chỉnh](./25_capstone_project.md)
