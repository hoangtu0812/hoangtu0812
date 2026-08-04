# Bài 21: Backpropagation Chi Tiết (Toán + Code From Scratch)

## Mục tiêu
- Suy luận ĐẦY ĐỦ, từng bước, công thức Backpropagation cho mạng nhiều lớp — đây là bài quan trọng nhất về mặt toán học của cả lộ trình.
- Chỉ sau khi hiểu trọn vẹn suy luận, mới cài đặt Neural Network from scratch bằng NumPy (không dùng PyTorch/TensorFlow) để chứng minh bạn hiểu đúng cơ chế.

---

## PHẦN A — Ý NGHĨA TOÁN HỌC

### 1. Bài toán: cần tính gradient của loss theo HÀNG TRIỆU tham số

Sau forward propagation ([Bài 20 mục 4](./20_neural_networks.md)), ta có $\hat{y}$ và loss $L(\hat{y},y)$. Để huấn luyện bằng Gradient Descent ([Bài 11 mục 5](./11_discrete_math_optimization.md)), cần $\frac{\partial L}{\partial W^{[l]}}$ cho MỌI lớp $l$ — kể cả các lớp ĐẦU TIÊN, cách xa $L$ qua rất nhiều phép biến đổi trung gian. Tính trực tiếp bằng định nghĩa (numerical differentiation — [Bài 5](./5_calculus_derivatives.md)) cho từng tham số riêng lẻ sẽ tốn $O(\text{số tham số})$ lần forward pass — hoàn toàn bất khả thi với mạng triệu tham số. Backpropagation là thuật toán tính CHÍNH XÁC mọi gradient đó chỉ với 1 lần forward + 1 lần "lan truyền ngược" — nhờ áp dụng có hệ thống Chain Rule đã học ở [Bài 6 mục 5](./6_calculus_matrix_calculus.md).

### 2. Thiết lập ký hiệu cho mạng $L$ lớp

$$\vec{z}^{[l]} = W^{[l]}\vec{a}^{[l-1]}+\vec{b}^{[l]}, \qquad \vec{a}^{[l]}=\phi(\vec{z}^{[l]}) \qquad (\vec{a}^{[0]}=\vec{x})$$

Mục tiêu: tính $\frac{\partial L}{\partial W^{[l]}}$ và $\frac{\partial L}{\partial \vec{b}^{[l]}}$ cho MỌI $l=1,...,L$.

### 3. Định nghĩa "lỗi cục bộ" $\delta^{[l]} = \frac{\partial L}{\partial\vec{z}^{[l]}}$ — chìa khóa của Backpropagation

Thay vì tính trực tiếp $\frac{\partial L}{\partial W^{[l]}}$ (phức tạp vì $W^{[l]}$ ảnh hưởng tới $L$ qua CHUỖI DÀI các lớp phía sau), ta chia bài toán thành 2 bước dễ hơn: (1) tính $\delta^{[l]}$ — "$L$ nhạy cảm thế nào với $\vec{z}^{[l]}$" — TỪ LỚP CUỐI LÙI VỀ LỚP ĐẦU; (2) từ $\delta^{[l]}$ suy ra $\frac{\partial L}{\partial W^{[l]}}$ dễ dàng (bước cuối, mục 6).

### 4. Bước khởi đầu — tính $\delta$ ở LỚP CUỐI (dễ nhất, vì gần $L$ nhất)

$$\delta^{[L]} = \frac{\partial L}{\partial\vec{z}^{[L]}} = \frac{\partial L}{\partial\vec{a}^{[L]}}\odot\phi'(\vec{z}^{[L]})$$

($\odot$ = nhân từng phần tử — element-wise, vì mỗi $z_i^{[L]}$ chỉ ảnh hưởng tới $a_i^{[L]}$ tương ứng qua $\phi$, không ảnh hưởng chéo). Đây là chain rule 1 bước: $\vec{z}^{[L]}\to\vec{a}^{[L]}\to L$, giống hệt ví dụ 1 neuron ở [Bài 6 mục 5](./6_calculus_matrix_calculus.md).

**Trường hợp đặc biệt CỰC KỲ đẹp — Softmax + Cross-Entropy** (tổ hợp thường dùng nhất ở lớp cuối bài toán phân loại): dù $\frac{\partial L}{\partial\vec{a}^{[L]}}$ và $\phi'(\vec{z}^{[L]})$ riêng lẻ đều phức tạp (softmax có đạo hàm dạng ma trận Jacobian đầy đủ — [Bài 6 mục 4](./6_calculus_matrix_calculus.md), không phải element-wise đơn giản), khi NHÂN LẠI theo đúng chain rule, mọi thứ rút gọn ngoạn mục thành:

$$\delta^{[L]} = \vec{a}^{[L]} - \vec{y}$$

(với $\vec{y}$ là one-hot vector của nhãn thật) — đơn giản đến bất ngờ: "lỗi cục bộ ở lớp cuối = dự đoán trừ nhãn thật". Đây là lý do tổ hợp Softmax+Cross-Entropy được dùng gần như mặc định — vừa có ý nghĩa xác suất đẹp ([Bài 10 mục 3](./10_statistics_inference.md)), vừa cho công thức gradient cực gọn.

### 5. Bước lan truyền ngược — suy ra $\delta^{[l-1]}$ TỪ $\delta^{[l]}$ (trái tim của thuật toán)

Đây là bước chain rule QUAN TRỌNG NHẤT: $\vec{z}^{[l-1]}$ ảnh hưởng tới $L$ CHỈ THÔNG QUA $\vec{a}^{[l-1]}$, rồi qua $\vec{z}^{[l]}$ (vì $\vec{z}^{[l]}=W^{[l]}\vec{a}^{[l-1]}+\vec{b}^{[l]}$), rồi mới tới $L$. Áp dụng chain rule qua chuỗi $\vec{z}^{[l-1]}\to\vec{a}^{[l-1]}\to\vec{z}^{[l]}\to L$:

$$\delta^{[l-1]} = \frac{\partial L}{\partial\vec{z}^{[l-1]}} = \left(\frac{\partial\vec{z}^{[l]}}{\partial\vec{a}^{[l-1]}}\right)^T\frac{\partial L}{\partial\vec{z}^{[l]}}\odot\phi'(\vec{z}^{[l-1]})$$

Vì $\vec{z}^{[l]}=W^{[l]}\vec{a}^{[l-1]}+\vec{b}^{[l]}$ là hàm TUYẾN TÍNH theo $\vec{a}^{[l-1]}$, Jacobian $\frac{\partial\vec{z}^{[l]}}{\partial\vec{a}^{[l-1]}}=W^{[l]}$ (đúng theo công thức gradient dạng tuyến tính đã suy ra ở [Bài 6 mục 2](./6_calculus_matrix_calculus.md), áp dụng cho từng dòng của $W^{[l]}$). Vậy:

$$\delta^{[l-1]} = \left((W^{[l]})^T\delta^{[l]}\right)\odot\phi'(\vec{z}^{[l-1]})$$

**Đây chính là công thức "lan truyền NGƯỢC":** để có lỗi cục bộ ở lớp $l-1$, ta lấy lỗi cục bộ ở lớp $l$ (đã biết từ bước trước), "đẩy ngược" qua $(W^{[l]})^T$ (CHÚ Ý: dùng **chuyển vị** của đúng ma trận trọng số đã dùng lúc forward — [Bài 3 mục 2](./3_linear_algebra_matrices.md), phép nhân với $W^T$ ở đây thực hiện vai trò "đảo ngược" hướng truyền dữ liệu so với lúc forward dùng $W$), rồi nhân với $\phi'(\vec{z}^{[l-1]})$ (mức độ "mở cổng" của hàm kích hoạt tại lớp đó — nếu $\phi'\approx0$, như vùng bão hòa của sigmoid ở [Bài 20 mục 3](./20_neural_networks.md), lỗi gần như KHÔNG truyền qua được — đây chính là hiện tượng vanishing gradient nhìn từ góc độ công thức).

Lặp lại công thức này từ $l=L$ lùi về $l=1$, ta có $\delta^{[l]}$ cho MỌI lớp — đây là lý do gọi là "back"-propagation: tính toán đi NGƯỢC chiều so với forward.

### 6. Bước cuối — từ $\delta^{[l]}$ suy ra gradient của $W^{[l]}, \vec{b}^{[l]}$

Vì $\vec{z}^{[l]}=W^{[l]}\vec{a}^{[l-1]}+\vec{b}^{[l]}$, áp dụng chain rule 1 bước cuối cùng:

$$\frac{\partial L}{\partial W^{[l]}} = \delta^{[l]}(\vec{a}^{[l-1]})^T, \qquad \frac{\partial L}{\partial\vec{b}^{[l]}} = \delta^{[l]}$$

(công thức $\frac{\partial L}{\partial W^{[l]}}=\delta^{[l]}(\vec{a}^{[l-1]})^T$ là tích ngoài — outer product — của 2 vector, cho ra đúng ma trận cùng kích thước với $W^{[l]}$; có thể kiểm chứng bằng cách viết ra từng phần tử $\frac{\partial L}{\partial W^{[l]}_{ij}}=\delta_i^{[l]}a_j^{[l-1]}$, đúng theo quy tắc đạo hàm dạng tuyến tính ở [Bài 6 mục 2](./6_calculus_matrix_calculus.md) áp dụng riêng cho từng phần tử ma trận).

**Tóm tắt toàn bộ thuật toán (4 bước):**
1. Forward pass: tính $\vec{z}^{[l]},\vec{a}^{[l]}$ mọi lớp, lưu lại để dùng sau.
2. Tính $\delta^{[L]}$ ở lớp cuối (mục 4).
3. Lan truyền ngược $\delta^{[l-1]}$ từ $\delta^{[l]}$, với $l=L,...,2$ (mục 5).
4. Từ mỗi $\delta^{[l]}$, tính gradient $W^{[l]},\vec{b}^{[l]}$ (mục 6), cập nhật bằng Gradient Descent ([Bài 11](./11_discrete_math_optimization.md)).

---

## PHẦN B — Cài đặt Neural Network From Scratch (NumPy thuần)

```python
import numpy as np

def relu(z): return np.maximum(0, z)
def relu_derivative(z): return (z > 0).astype(float)  # đúng theo mục 3: ReLU'(z) = 1 nếu z>0, else 0

def softmax(z):
	exp_z = np.exp(z - np.max(z))
	return exp_z / np.sum(exp_z)

class NeuralNetworkFromScratch:
	def __init__(self, layer_sizes, lr=0.01):
		self.lr = lr
		self.L = len(layer_sizes) - 1
		self.W, self.b = [], []
		for i in range(self.L):
			self.W.append(np.random.randn(layer_sizes[i+1], layer_sizes[i]) * np.sqrt(2/layer_sizes[i]))
			self.b.append(np.zeros(layer_sizes[i+1]))

	def forward(self, x):
		"""Bước 1 (Phần A mục 6, tóm tắt): lưu lại z, a mọi lớp để dùng cho backward."""
		self.a = [x]
		self.z = []
		for l in range(self.L - 1):
			z = self.W[l] @ self.a[-1] + self.b[l]
			self.z.append(z)
			self.a.append(relu(z))
		z_out = self.W[-1] @ self.a[-1] + self.b[-1]
		self.z.append(z_out)
		self.a.append(softmax(z_out))
		return self.a[-1]

	def backward(self, y_true_onehot):
		"""Bước 2-4 (Phần A mục 4-6): tính delta ngược từ lớp cuối, suy ra gradient."""
		grads_W, grads_b = [None]*self.L, [None]*self.L

		# Bước 2 (mục 4): delta lớp cuối — công thức đẹp của Softmax+CrossEntropy
		delta = self.a[-1] - y_true_onehot

		for l in reversed(range(self.L)):
			# Bước 4 (mục 6): gradient của W[l], b[l] từ delta hiện tại
			grads_W[l] = np.outer(delta, self.a[l])   # a[l] chính là a^[l-1] theo ký hiệu 1-indexed ở Phần A
			grads_b[l] = delta

			if l > 0:
				# Bước 3 (mục 5): lan truyền delta ngược về lớp trước
				delta = (self.W[l].T @ delta) * relu_derivative(self.z[l-1])

		return grads_W, grads_b

	def update(self, grads_W, grads_b):
		for l in range(self.L):
			self.W[l] -= self.lr * grads_W[l]
			self.b[l] -= self.lr * grads_b[l]

	def train_step(self, x, y_true_onehot):
		y_pred = self.forward(x)
		loss = -np.sum(y_true_onehot * np.log(y_pred + 1e-15))  # Cross-Entropy — Bài 10 mục 3
		grads_W, grads_b = self.backward(y_true_onehot)
		self.update(grads_W, grads_b)
		return loss

if __name__ == "__main__":
	np.random.seed(42)
	nn = NeuralNetworkFromScratch([4, 8, 3], lr=0.1)  # input 4 chiều, 1 lớp ẩn, output 3 lớp

	# Verify bằng gradient checking — so sánh backprop với numerical differentiation (Bài 5)
	x = np.random.randn(4)
	y_onehot = np.array([0, 1, 0])

	nn.forward(x)
	grads_W, _ = nn.backward(y_onehot)

	def loss_fn(W0):
		old = nn.W[0].copy()
		nn.W[0] = W0
		y_pred = nn.forward(x)
		loss = -np.sum(y_onehot * np.log(y_pred + 1e-15))
		nn.W[0] = old
		return loss

	# So sánh 1 phần tử gradient tính bằng backprop vs numerical differentiation
	i, j = 0, 0
	h = 1e-5
	W_plus = nn.W[0].copy(); W_plus[i,j] += h
	W_minus = nn.W[0].copy(); W_minus[i,j] -= h
	numerical_grad = (loss_fn(W_plus) - loss_fn(W_minus)) / (2*h)

	print(f"Backprop gradient: {grads_W[0][i,j]:.6f}")
	print(f"Numerical gradient: {numerical_grad:.6f}")  # phải xấp xỉ bằng nhau — verify backprop ĐÚNG
```

**Gradient Checking** (kỹ thuật ở cuối ví dụ) là bước KHÔNG THỂ THIẾU khi tự cài Backpropagation — so sánh gradient tính bằng công thức (nhanh, dễ sai khi code) với gradient tính bằng định nghĩa numerical differentiation ([Bài 5 mục 1](./5_calculus_derivatives.md), chậm nhưng gần như luôn đúng) — nếu 2 kết quả khớp nhau, ta có bằng chứng mạnh rằng cài đặt Backpropagation của mình chính xác.

## Bài tập

1. **Tự suy luận lại công thức $\delta^{[l-1]}$ từ $\delta^{[l]}$**: che phần "Bước lan truyền ngược" (mục 5), tự làm lại phép suy luận chain rule, so sánh với công thức trong bài.
2. **Gradient Checking đầy đủ**: mở rộng ví dụ cuối bài, verify gradient checking cho MỌI phần tử của $W^{[0]}$ và $W^{[1]}$ (không chỉ 1 phần tử), báo cáo sai số lớn nhất.
3. **Huấn luyện thật**: dùng `NeuralNetworkFromScratch`, huấn luyện trên dataset Iris (3 lớp, 4 đặc trưng — dùng `train_step` lặp lại nhiều epoch), đo accuracy trên tập test.
4. **Thí nghiệm vanishing gradient**: thay `relu`/`relu_derivative` bằng `sigmoid`/đạo hàm sigmoid trong toàn bộ mạng, xây mạng SÂU (5-6 lớp), quan sát độ lớn của `grads_W[0]` (lớp đầu tiên) so với `grads_W[-1]` (lớp cuối) — verify bằng số hiện tượng vanishing gradient đã lý luận ở mục 5.

## Tiếp theo
→ [Bài 22: Optimizer nâng cao & Regularization trong Deep Learning](./22_optimizers_dl_regularization.md)
