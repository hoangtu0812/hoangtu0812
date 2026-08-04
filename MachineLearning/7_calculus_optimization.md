# Bài 7: Tối ưu hóa không ràng buộc & Taylor Series

## Mục tiêu
- Hiểu vì sao $\nabla f=0$ chỉ là điều kiện CẦN, không ĐỦ để tìm cực tiểu.
- Suy luận Taylor Series từ ý tưởng "xấp xỉ hàm bằng đa thức".
- Chứng minh trực giác vì sao Gradient Descent hội tụ — TRƯỚC KHI chạy bất kỳ dòng code nào.

---

## PHẦN A — Ý NGHĨA TOÁN HỌC

### 1. Điểm cực trị — điều kiện cần bậc 1

Tại điểm cực tiểu/cực đại của $f$, tiếp tuyến (đồ thị đạo hàm — [Bài 5 mục 1](./5_calculus_derivatives.md)) phải **nằm ngang** — nếu còn dốc theo hướng nào, ta luôn có thể di chuyển theo hướng đó để giảm/tăng $f$ thêm nữa, nghĩa là điểm đó CHƯA phải cực trị. Vậy điều kiện cần: $f'(x^*)=0$ (hàm nhiều biến: $\nabla f(x^*)=\vec{0}$, [Bài 5 mục 6](./5_calculus_derivatives.md)).

**Nhưng $\nabla f=0$ chưa đủ** — 3 khả năng có thể xảy ra tại điểm đó:
- **Cực tiểu**: đáy 1 "thung lũng" — mọi hướng xung quanh đều đi lên.
- **Cực đại**: đỉnh 1 "ngọn đồi" — mọi hướng xung quanh đều đi xuống.
- **Điểm yên ngựa (saddle point)**: giống yên ngựa thật — đi lên theo hướng này, đi xuống theo hướng khác, dù tại chính điểm đó "phẳng" theo mọi hướng tức thời.

Cần thêm thông tin bậc 2 (độ **cong** của hàm) để phân biệt 3 trường hợp này.

### 2. Hessian Matrix — "đạo hàm của đạo hàm", đo độ cong

Đạo hàm bậc 2 $f''(x)$ (hàm 1 biến) đo tốc độ thay đổi CỦA CHÍNH đạo hàm — nếu $f''(x)>0$, độ dốc đang tăng dần (hàm "cong lên", giống đáy chén) → đó là dấu hiệu của cực tiểu. Nếu $f''(x)<0$, độ dốc đang giảm dần (hàm "cong xuống", giống đỉnh đồi) → dấu hiệu cực đại.

Với hàm nhiều biến, gộp mọi đạo hàm riêng bậc 2 thành ma trận Hessian:

$$H = \begin{bmatrix}\partial^2f/\partial x_1^2 & \partial^2f/\partial x_1\partial x_2\\ \partial^2f/\partial x_2\partial x_1 & \partial^2f/\partial x_2^2\end{bmatrix}$$

**Điều kiện đủ bậc 2** (tại điểm đã có $\nabla f=0$):
- $H$ **xác định dương** (mọi trị riêng $>0$ — [Bài 4 mục 1](./4_linear_algebra_eigen_svd.md)) → hàm "cong lên" theo MỌI hướng → **cực tiểu**.
- $H$ **xác định âm** (mọi trị riêng $<0$) → cong xuống mọi hướng → **cực đại**.
- $H$ có cả trị riêng dương lẫn âm → cong lên theo hướng này, cong xuống theo hướng khác → **saddle point**.

Đây là lý do vì sao [Bài 4](./4_linear_algebra_eigen_svd.md) (trị riêng) lại cần thiết ở đây — dấu của trị riêng Hessian chính là "công cụ phân loại" điểm dừng.

### 3. Hàm lồi — trường hợp "trong mơ" để tối ưu

Hàm $f$ là **lồi** nếu đoạn thẳng nối 2 điểm bất kỳ trên đồ thị luôn nằm **phía trên hoặc trùng** với chính đồ thị đó:

$$f(\lambda x_1+(1-\lambda)x_2) \leq \lambda f(x_1)+(1-\lambda)f(x_2), \quad\forall\lambda\in[0,1]$$

Về mặt Hessian: $f$ lồi $\Leftrightarrow$ $H\succeq0$ (xác định bán dương) tại MỌI điểm — tức hàm luôn "cong lên hoặc phẳng", không bao giờ "cong xuống" ở đâu cả.

**Hệ quả cực kỳ quan trọng:** với hàm lồi, KHÔNG THỂ có saddle point hay "cực tiểu giả" — **mọi điểm có $\nabla f=0$ đều là cực tiểu TOÀN CỤC**. Đây là lý do Linear Regression (MSE) và Logistic Regression (Cross-Entropy) "dễ" tối ưu: hàm loss của chúng lồi ([Bài 6 mục 3](./6_calculus_matrix_calculus.md) đã chỉ ra $\nabla^2(\vec{x}^TA\vec{x})=2A\succeq0$ khi $A=X^TX$, vì $\vec{v}^TX^TX\vec{v}=\|X\vec{v}\|^2\geq0$ luôn đúng). Ngược lại, mạng neural sâu có hàm loss **không lồi** — có thể có vô số saddle point/cực tiểu cục bộ, đây là lý do Deep Learning cần các kỹ thuật tối ưu tinh vi hơn ([Bài 22](./22_optimizers_dl_regularization.md)).

### 4. Khai triển Taylor — xấp xỉ hàm phức tạp bằng đa thức đơn giản

**Ý tưởng xuất phát:** quanh 1 điểm $x_0$, hàm $f$ bất kỳ (dù phức tạp) có thể xấp xỉ tốt bằng 1 đa thức, nếu ta biết đủ thông tin về đạo hàm tại $x_0$. Bậc càng cao, xấp xỉ càng chính xác (nhưng chỉ **cục bộ**, gần $x_0$).

**Bậc 0** (chỉ dùng giá trị hàm): $f(x)\approx f(x_0)$ — xấp xỉ tệ nhất, chỉ đúng khi $x=x_0$.

**Bậc 1** (thêm thông tin độ dốc — chính là phương trình tiếp tuyến ở [Bài 5 mục 1](./5_calculus_derivatives.md)): $f(x)\approx f(x_0)+f'(x_0)(x-x_0)$.

**Bậc 2** (thêm thông tin độ cong — mục 2): $f(x)\approx f(x_0)+f'(x_0)(x-x_0)+\frac{1}{2}f''(x_0)(x-x_0)^2$.

Với hàm nhiều biến (dùng gradient & Hessian):

$$f(\vec{x})\approx f(\vec{x_0})+\nabla f(\vec{x_0})^T(\vec{x}-\vec{x_0})+\frac{1}{2}(\vec{x}-\vec{x_0})^TH(\vec{x}-\vec{x_0})$$

**Ví dụ tính tay** (bậc 1, kiểm chứng trực giác "gần đúng"): $f(x)=e^x$, $x_0=0$. $f(0)=1$, $f'(0)=e^0=1$. Xấp xỉ bậc 1: $f(x)\approx 1+x$. Với $x=0.1$: xấp xỉ $=1.1$, giá trị thật $e^{0.1}\approx1.105$ — sai số rất nhỏ vì $x$ gần $x_0$. Với $x=2$: xấp xỉ $=3$, giá trị thật $e^2\approx7.389$ — sai số LỚN vì $x$ đã xa $x_0$, đúng như đã cảnh báo "chỉ chính xác cục bộ".

### 5. Vì sao Gradient Descent hội tụ — chứng minh trực giác bằng Taylor bậc 1

Đây là kết quả quan trọng nhất của bài — hãy đi chậm từng bước.

Xét bước cập nhật Gradient Descent: $\vec{w}_{t+1}=\vec{w}_t-\eta\nabla L(\vec{w}_t)$ (η là learning rate). Áp dụng khai triển Taylor bậc 1 (mục 4) cho $L$ tại điểm $\vec{w}_{t+1}$, khai triển quanh $\vec{w}_t$:

$$L(\vec{w}_{t+1}) \approx L(\vec{w}_t) + \nabla L(\vec{w}_t)^T(\vec{w}_{t+1}-\vec{w}_t)$$

Thay $\vec{w}_{t+1}-\vec{w}_t = -\eta\nabla L(\vec{w}_t)$ (đúng theo công thức cập nhật):

$$L(\vec{w}_{t+1}) \approx L(\vec{w}_t) - \eta\nabla L(\vec{w}_t)^T\nabla L(\vec{w}_t) = L(\vec{w}_t) - \eta\|\nabla L(\vec{w}_t)\|^2$$

Vì $\|\nabla L(\vec{w}_t)\|^2\geq0$ luôn đúng (norm bình phương không bao giờ âm — [Bài 2 mục 4](./2_linear_algebra_vectors.md)), và $\eta>0$, số hạng $-\eta\|\nabla L\|^2 \leq 0$. Vậy: **$L(\vec{w}_{t+1}) \leq L(\vec{w}_t)$** — mỗi bước cập nhật đều làm giảm (hoặc giữ nguyên) giá trị loss, MIỄN LÀ $\eta$ đủ nhỏ để xấp xỉ Taylor bậc 1 còn chính xác.

![Gradient Descent đi ngược chiều gradient tới cực tiểu](./images/gradient_descent_contour.svg)

**Vì sao $\eta$ quá lớn lại nguy hiểm?** Nếu $\eta$ lớn, bước nhảy $\vec{w}_{t+1}-\vec{w}_t$ lớn, điểm $\vec{w}_{t+1}$ đã "đi xa" khỏi $\vec{w}_t$ đến mức xấp xỉ Taylor bậc 1 (chỉ đúng CỤC BỘ, như đã thấy ở mục 4 ví dụ $x=2$) không còn phản ánh đúng hàm $L$ thật nữa — loss có thể **TĂNG** thay vì giảm, dẫn tới dao động hoặc phân kỳ (divergence).

---

## PHẦN B — Cài đặt & Minh họa bằng code

```python
import numpy as np

# Mục 2: Hessian numerical — verify phân loại điểm cực trị
def hessian_numerical(f, x, h=1e-4):
	n = len(x)
	H = np.zeros((n, n))
	for i in range(n):
		for j in range(n):
			x_pp = x.copy(); x_pp[i] += h; x_pp[j] += h
			x_pm = x.copy(); x_pm[i] += h; x_pm[j] -= h
			x_mp = x.copy(); x_mp[i] -= h; x_mp[j] += h
			x_mm = x.copy(); x_mm[i] -= h; x_mm[j] -= h
			H[i,j] = (f(x_pp) - f(x_pm) - f(x_mp) + f(x_mm)) / (4*h*h)
	return H

f_bowl = lambda x: x[0]**2 + x[1]**2  # cực tiểu tại (0,0)
H = hessian_numerical(f_bowl, np.array([1.0, 1.0]))
print(H)                                    # xấp xỉ [[2,0],[0,2]]
print(np.linalg.eigvals(H))                  # cả 2 trị riêng dương -> cực tiểu, khớp lý thuyết mục 2

f_saddle = lambda x: x[0]**2 - x[1]**2  # saddle point tại (0,0)
H_saddle = hessian_numerical(f_saddle, np.array([0.0, 0.0]))
print(np.linalg.eigvals(H_saddle))           # 1 dương, 1 âm -> saddle point

# Mục 4: verify Taylor xấp xỉ bậc 2
def taylor_approx_2nd_order(f, x0, x, h=1e-5):
	f_x0 = f(x0)
	f_prime = (f(x0+h) - f(x0-h)) / (2*h)
	f_double_prime = (f(x0+h) - 2*f(x0) + f(x0-h)) / (h**2)
	return f_x0 + f_prime*(x-x0) + 0.5*f_double_prime*(x-x0)**2

f_exp = np.exp
for x in [0.1, 0.5, 1.0, 2.0]:
	approx = taylor_approx_2nd_order(f_exp, 0, x)
	print(f"x={x}: xấp xỉ={approx:.4f}, thật={f_exp(x):.4f}, sai số={abs(approx-f_exp(x)):.4f}")
	# sai số tăng dần khi x xa điểm khai triển x0=0 — đúng như lý luận ở mục 4

# Mục 5: verify ảnh hưởng của learning rate
def demonstrate_learning_rate_effect(f, grad_f, start, learning_rates, steps=10):
	for lr in learning_rates:
		w = start
		for _ in range(steps):
			w = w - lr * grad_f(w)
		print(f"lr={lr}: w cuối = {w:.4f}, loss = {f(w):.4f}")

f_1d = lambda w: (w - 5)**2
grad_1d = lambda w: 2*(w - 5)
demonstrate_learning_rate_effect(f_1d, grad_1d, start=0.0, learning_rates=[0.01, 0.1, 0.5, 1.0, 1.1])
# lr nhỏ: hội tụ chậm. lr vừa: hội tụ tốt. lr=1.0: dao động ổn định. lr=1.1: PHÂN KỲ — đúng cảnh báo ở mục 5
```

## Bài tập

1. **Hessian & phân loại điểm**: tính tay Hessian của $f(x,y)=x^2-y^2$ tại $(0,0)$, xác định dấu 2 trị riêng bằng tay, verify bằng `hessian_numerical`.
2. **Kiểm tra tính lồi của MSE**: dùng `np.linalg.eigvals`, verify Hessian của hàm MSE ([Bài 6 mục 3](./6_calculus_matrix_calculus.md)) luôn có trị riêng $\geq0$ tại nhiều điểm ngẫu nhiên — xác nhận nó lồi.
3. **Taylor approximation**: thử với hàm $\sin(x)$ quanh $x_0=0$, vẽ đồ thị so sánh hàm thật với xấp xỉ Taylor bậc 2 trên $[-2,2]$ (liên hệ [Bài 13](./13_data_visualization_eda.md)) — quan sát trực quan vùng nào xấp xỉ tốt/tệ.
4. **Learning rate**: chạy `demonstrate_learning_rate_effect` với nhiều mức learning rate, tự giải thích bằng lời (dựa vào lập luận Taylor ở mục 5) tại sao $\eta$ quá lớn gây phân kỳ.

## Tiếp theo
→ [Bài 8: Xác suất cơ bản & Định lý Bayes](./8_probability_basics.md)
