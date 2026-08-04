# Bài 23: CNN (Convolutional Neural Networks)

## Mục tiêu
- Hiểu phép tích chập (convolution) từ ý nghĩa toán học, không chỉ "công cụ xử lý ảnh".
- Nắm vì sao CNN hiệu quả hơn MLP thuần cho dữ liệu ảnh — lý luận bằng số lượng tham số.

---

## PHẦN A — Ý NGHĨA TOÁN HỌC

### 1. Vì sao KHÔNG dùng MLP thuần cho ảnh? — lý luận bằng con số

Ảnh $224\times224\times3$ (RGB) có $224\times224\times3=150{,}528$ giá trị đầu vào. Nếu lớp ẩn đầu tiên của MLP ([Bài 20](./20_neural_networks.md)) có 1000 neuron, số tham số CHỈ RIÊNG lớp này đã là $150{,}528\times1000\approx150$ triệu — cực kỳ tốn bộ nhớ, dễ overfitting ([Bài 16-17](./16_model_evaluation.md), model quá nhiều tham số so với dữ liệu), và (quan trọng nhất) **bỏ lỡ hoàn toàn cấu trúc không gian** của ảnh: MLP coi mỗi pixel là 1 đặc trưng độc lập, không "biết" rằng pixel liền kề thường liên quan chặt chẽ với nhau (1 cạnh, 1 góc, 1 vùng màu đồng nhất).

### 2. Phép tích chập (Convolution) — trượt 1 "bộ lọc" nhỏ qua toàn ảnh

Convolution 2D với kernel (bộ lọc) $K$ kích thước $k\times k$, áp dụng lên ảnh $I$:

$$(I*K)_{i,j} = \sum_{m=0}^{k-1}\sum_{n=0}^{k-1} I_{i+m,j+n}\cdot K_{m,n}$$

**Ý nghĩa:** mỗi giá trị đầu ra là 1 **dot product** ([Bài 2 mục 3](./2_linear_algebra_vectors.md)) giữa kernel và 1 "cửa sổ" $k\times k$ của ảnh gốc — sau đó kernel TRƯỢT sang vị trí kế tiếp và lặp lại. Về bản chất, mỗi giá trị đầu ra đo "mức độ giống nhau" giữa vùng ảnh cục bộ đó và pattern mà kernel biểu diễn (giống ý nghĩa dot product đo mức "cùng hướng" đã học ở [Bài 2](./2_linear_algebra_vectors.md), giờ áp dụng cho pattern không gian 2D thay vì vector đặc trưng).

**Ví dụ trực giác:** 1 kernel $3\times3$ dạng $\begin{bmatrix}-1&0&1\\-1&0&1\\-1&0&1\end{bmatrix}$ tính "chênh lệch độ sáng trái-phải" tại mỗi vị trí — kết quả tích chập sẽ có giá trị LỚN (dương hoặc âm) đúng tại các cạnh dọc trong ảnh (nơi độ sáng đổi đột ngột theo chiều ngang), và gần 0 ở vùng đồng nhất — đây chính là 1 **bộ phát hiện cạnh (edge detector)** cổ điển, minh họa CNN "học được" cách phát hiện đặc trưng từ đâu.

### 3. 2 tính chất then chốt khiến CNN hiệu quả hơn MLP cho ảnh

**Chia sẻ tham số (Parameter Sharing):** CÙNG 1 kernel được dùng lặp lại ở MỌI vị trí trên ảnh — thay vì mỗi vị trí có 1 bộ trọng số riêng (như MLP), toàn bộ ảnh chỉ cần $k\times k$ tham số cho 1 kernel. Giả định ẩn đằng sau: "nếu 1 pattern (vd cạnh dọc) hữu ích để phát hiện ở góc trên-trái ảnh, nó CŨNG hữu ích để phát hiện ở bất kỳ đâu khác trên ảnh" — giả định hợp lý cho hầu hết bài toán thị giác máy tính (**translation invariance**).

**Kết nối cục bộ (Local Connectivity):** mỗi giá trị đầu ra chỉ phụ thuộc vào 1 VÙNG NHỎ ($k\times k$) của input, không phải TOÀN BỘ ảnh như MLP (mỗi neuron MLP kết nối với MỌI pixel input). Điều này khớp với trực giác: để nhận biết 1 cạnh/góc, chỉ cần nhìn vùng lân cận cục bộ, không cần "nhìn" toàn bộ ảnh cùng lúc.

**Hệ quả về số tham số:** 1 lớp Conv với 64 kernel $3\times3\times3$ (áp dụng cho ảnh RGB) chỉ có $64\times(3\times3\times3+1)=1{,}792$ tham số — SO VỚI hàng trăm triệu của MLP ở mục 1 — dù vẫn xử lý được ảnh kích thước đầy đủ. Đây là lý do TOÁN HỌC (giảm số tham số nhờ 2 tính chất trên) khiến CNN vừa hiệu quả hơn về tính toán, vừa giảm nguy cơ overfitting ([Bài 17](./17_regularization.md)).

### 4. Padding & Stride — kiểm soát kích thước đầu ra

Không có padding, ảnh $n\times n$ qua kernel $k\times k$ cho đầu ra $(n-k+1)\times(n-k+1)$ — nhỏ dần sau mỗi lớp (vì kernel không "vừa" tại các pixel biên). **Padding** (thêm viền 0 quanh ảnh) giữ nguyên kích thước nếu cần. **Stride** $s$ (bước nhảy của kernel mỗi lần trượt, thay vì luôn là 1) cho công thức tổng quát kích thước đầu ra:

$$\text{output size} = \left\lfloor\frac{n+2p-k}{s}\right\rfloor + 1$$

($p$=padding). Stride $>1$ vừa giảm kích thước đầu ra (giảm tính toán ở lớp sau) vừa tăng "trường nhìn" hiệu dụng của các lớp sâu hơn.

### 5. Pooling — giảm chiều có chủ đích, giữ lại thông tin quan trọng nhất

**Max Pooling** lấy giá trị LỚN NHẤT trong mỗi vùng $k\times k$ (không chồng lấp, thường $k=2,s=2$) — giữ lại "đặc trưng mạnh nhất" trong vùng, loại bỏ chi tiết nhiễu/dư thừa. Về ý nghĩa, đây là 1 dạng ép buộc **bất biến nhỏ với dịch chuyển** (small translation invariance): nếu 1 đặc trưng dịch chuyển vài pixel trong vùng pooling, kết quả max không đổi — tăng độ "bền" của biểu diễn trước những biến dạng nhỏ vô hại của ảnh đầu vào.

### 6. Kiến trúc CNN điển hình — chuỗi lớp lặp lại

$$\text{Input} \to [\text{Conv}\to\text{ReLU}\to\text{Pooling}]\times N \to \text{Flatten} \to \text{Fully Connected (MLP)} \to \text{Softmax}$$

Các lớp Conv đầu học đặc trưng ĐƠN GIẢN (cạnh, góc, màu — gần input), các lớp Conv sâu hơn kết hợp đặc trưng đơn giản thành đặc trưng PHỨC TẠP hơn (hình dạng, texture, rồi tới bộ phận đối tượng) — đây là hệ quả tự nhiên của việc XẾP CHỒNG nhiều lớp phi tuyến, cùng nguyên lý đã chứng minh ở [Bài 20 mục 2](./20_neural_networks.md) (không có phi tuyến, xếp chồng vô nghĩa) áp dụng cho không gian đặc trưng ảnh.

---

## PHẦN B — Cài đặt & Minh họa bằng code

```python
import numpy as np

# Mục 2: cài đặt convolution 2D from scratch — verify ý nghĩa "trượt kernel"
def conv2d(image, kernel, stride=1, padding=0):
	if padding > 0:
		image = np.pad(image, padding, mode='constant')
	h, w = image.shape
	kh, kw = kernel.shape
	out_h = (h - kh) // stride + 1
	out_w = (w - kw) // stride + 1
	output = np.zeros((out_h, out_w))

	for i in range(out_h):
		for j in range(out_w):
			region = image[i*stride:i*stride+kh, j*stride:j*stride+kw]
			output[i, j] = np.sum(region * kernel)  # dot product 2D — mục 2
	return output

# Verify bằng edge detector — mục 2 ví dụ trực giác
image = np.array([
	[10, 10, 10, 0, 0, 0],
	[10, 10, 10, 0, 0, 0],
	[10, 10, 10, 0, 0, 0],
	[10, 10, 10, 0, 0, 0],
])
vertical_edge_kernel = np.array([[-1, 0, 1], [-1, 0, 1], [-1, 0, 1]])
result = conv2d(image, vertical_edge_kernel)
print(result)  # giá trị lớn (âm/dương) đúng tại cột chuyển từ sáng (10) sang tối (0)

# Mục 4: verify công thức kích thước output
def output_size(n, k, s=1, p=0):
	return (n + 2*p - k) // s + 1

print(output_size(28, 3, s=1, p=0))  # ảnh 28x28, kernel 3x3, không padding -> 26
print(output_size(28, 3, s=1, p=1))  # với padding=1 -> vẫn 28 (giữ nguyên kích thước)

# Mục 5: max pooling
def max_pool2d(feature_map, size=2, stride=2):
	h, w = feature_map.shape
	out_h, out_w = h // stride, w // stride
	output = np.zeros((out_h, out_w))
	for i in range(out_h):
		for j in range(out_w):
			region = feature_map[i*stride:i*stride+size, j*stride:j*stride+size]
			output[i, j] = np.max(region)
	return output

print(max_pool2d(result))
```

### CNN thực tế với PyTorch (dùng framework thay vì from-scratch cho model thật)

```python
import torch
import torch.nn as nn

class SimpleCNN(nn.Module):
	def __init__(self, num_classes=10):
		super().__init__()
		self.conv1 = nn.Conv2d(3, 32, kernel_size=3, padding=1)  # input 3 kênh (RGB) -> 32 kernel
		self.conv2 = nn.Conv2d(32, 64, kernel_size=3, padding=1)
		self.pool = nn.MaxPool2d(2, 2)
		self.fc = nn.Linear(64 * 8 * 8, num_classes)  # giả sử ảnh 32x32 sau 2 lần pool còn 8x8

	def forward(self, x):
		x = self.pool(torch.relu(self.conv1(x)))   # Conv -> ReLU -> Pool, đúng kiến trúc mục 6
		x = self.pool(torch.relu(self.conv2(x)))
		x = x.view(x.size(0), -1)                     # flatten
		return self.fc(x)                               # fully connected -> logits, đưa vào softmax ngoài

model = SimpleCNN(num_classes=10)
print(sum(p.numel() for p in model.parameters()))  # đếm tổng tham số — so sánh với MLP thuần ở mục 1
```

## Bài tập

1. **Cài `conv2d` from scratch**: dùng code mẫu trên, thử với kernel phát hiện cạnh NGANG (chuyển vị kernel dọc ở ví dụ), verify nó phát hiện đúng cạnh ngang thay vì dọc.
2. **Verify công thức kích thước output**: tính tay kích thước output cho ảnh $32\times32$, kernel $5\times5$, stride $2$, padding $2$; verify bằng `output_size`.
3. **So sánh số tham số CNN vs MLP**: tính tay số tham số của 1 lớp Conv (64 kernel $3\times3\times3$) và 1 lớp MLP tương đương xử lý cùng kích thước ảnh — verify chênh lệch lớn như lý luận ở mục 3.
4. **Huấn luyện `SimpleCNN` trên CIFAR-10**: dùng `torchvision.datasets.CIFAR10`, huấn luyện vài epoch, đo accuracy — so sánh với thử dùng MLP thuần cho cùng dataset (liên hệ [Bài 20](./20_neural_networks.md)).

## Tiếp theo
→ [Bài 24: RNN/LSTM & Giới thiệu Attention/Transformer](./24_rnn_transformer.md)
