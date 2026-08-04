# Lộ Trình Học Machine Learning Chi Tiết (kèm Toán nền tảng)

> Cấu trúc giống [Go/ROADMAP.md](../Go/ROADMAP.md), [Python/ROADMAP.md](../Python/ROADMAP.md), [SAP/ROADMAP.md](../SAP/ROADMAP.md): mỗi bài có file chi tiết riêng (lý thuyết + công thức + code mẫu + bài tập). Khác biệt lớn nhất của track này: ML đòi hỏi nền toán vững — nên lộ trình dành **hẳn Phần I** để học/ôn 4 mảng toán cốt lõi (Đại số tuyến tính, Giải tích, Xác suất & Thống kê, Toán rời rạc & Tối ưu hóa) TRƯỚC khi vào ML, và mỗi công thức toán đều được nối thẳng tới chỗ nó xuất hiện trong ML/Deep Learning.
>
> **Yêu cầu tiên quyết:** đã quen lập trình (dùng [Python/ROADMAP.md](../Python/ROADMAP.md) làm nền, đặc biệt [Python Bài 5](../Python/5_collections.md) và [Python Bài 17](../Python/17_database.md) về xử lý dữ liệu).

---

## Giai đoạn 0 — Giới thiệu & Cài đặt

### [Bài 1: Giới thiệu Machine Learning & Cài đặt môi trường](./1_get_started.md)
- ML là gì, 3 loại bài toán (Supervised/Unsupervised/Reinforcement Learning), ML fit vào đâu trong AI/Data Science.
- Cài Python cho ML: Anaconda/Miniconda, Jupyter Notebook, `numpy`/`pandas`/`matplotlib`/`scikit-learn`.
- **Bài tập:** Cài môi trường, chạy thử Jupyter Notebook, load 1 dataset mẫu (`sklearn.datasets`).

---

## Phần I — Toán Nền Tảng Cho Machine Learning

> Mỗi bài toán đều có mục "Ứng dụng trong ML" nối trực tiếp công thức tới thuật toán cụ thể ở Phần II/III — không học toán suông.

### A. Đại số tuyến tính (Linear Algebra)

#### [Bài 2: Vector & Không gian vector](./2_linear_algebra_vectors.md)
- Vector, phép cộng/nhân vô hướng, tích vô hướng (dot product), norm, khoảng cách, góc giữa 2 vector, vector độc lập tuyến tính.
- **Ứng dụng ML:** biểu diễn dữ liệu (feature vector), cosine similarity (recommendation), embedding.

#### [Bài 3: Ma trận & Hệ phương trình tuyến tính](./3_linear_algebra_matrices.md)
- Phép toán ma trận (cộng, nhân, chuyển vị), ma trận nghịch đảo, định thức, hạng ma trận, giải hệ phương trình tuyến tính (Gaussian elimination).
- **Ứng dụng ML:** Linear Regression dạng ma trận (Normal Equation), biểu diễn dataset dạng ma trận (design matrix).

#### [Bài 4: Trị riêng, Vector riêng & SVD](./4_linear_algebra_eigen_svd.md)
- Eigenvalue/Eigenvector, chéo hóa ma trận, Singular Value Decomposition (SVD), ma trận xác định dương (positive semi-definite).
- **Ứng dụng ML:** PCA (Principal Component Analysis), giảm chiều dữ liệu, nén dữ liệu, recommendation system (matrix factorization).

### B. Giải tích (Calculus)

#### [Bài 5: Đạo hàm & Gradient](./5_calculus_derivatives.md)
- Đạo hàm 1 biến, quy tắc đạo hàm (chain rule, product rule), đạo hàm riêng (partial derivative), gradient của hàm nhiều biến.
- **Ứng dụng ML:** cách 1 model "học" — gradient chỉ hướng tăng nhanh nhất của hàm loss.

#### [Bài 6: Đạo hàm ma trận & Chain Rule (nền tảng Backpropagation)](./6_calculus_matrix_calculus.md)
- Đạo hàm theo vector/ma trận, Jacobian, chain rule nhiều lớp.
- **Ứng dụng ML:** đây CHÍNH LÀ toán học đứng sau thuật toán Backpropagation ([Bài 21](./21_neural_networks.md)).

#### [Bài 7: Tối ưu hóa không ràng buộc & Taylor Series](./7_calculus_optimization.md)
- Điểm cực trị, Hessian matrix, điều kiện tối ưu bậc 1/bậc 2, khai triển Taylor, ý nghĩa hình học của gradient descent.
- **Ứng dụng ML:** vì sao Gradient Descent hội tụ, learning rate ảnh hưởng thế nào.

### C. Xác suất & Thống kê (Probability & Statistics)

#### [Bài 8: Xác suất cơ bản & Định lý Bayes](./8_probability_basics.md)
- Không gian mẫu, biến cố, xác suất có điều kiện, độc lập, quy tắc nhân/cộng xác suất, định lý Bayes.
- **Ứng dụng ML:** Naive Bayes Classifier, suy luận Bayesian, spam filter.

#### [Bài 9: Biến ngẫu nhiên & Phân phối xác suất](./9_probability_distributions.md)
- Biến ngẫu nhiên rời rạc/liên tục, hàm mật độ (PDF), hàm phân phối tích lũy (CDF), kỳ vọng, phương sai, hiệp phương sai.
- Các phân phối quan trọng: Bernoulli, Binomial, Poisson, Uniform, **Normal/Gaussian**.
- **Ứng dụng ML:** giả định phân phối trong Linear Regression, khởi tạo trọng số mạng neural, Gaussian Naive Bayes.

#### [Bài 10: Ước lượng thống kê & Kiểm định giả thuyết](./10_statistics_inference.md)
- Maximum Likelihood Estimation (MLE), Maximum A Posteriori (MAP), khoảng tin cậy, kiểm định giả thuyết (hypothesis testing), p-value.
- **Ứng dụng ML:** vì sao hàm loss Cross-Entropy = MLE cho phân phối Bernoulli/Categorical, A/B testing đánh giá model.

### D. Toán rời rạc & Tối ưu hóa (Discrete Math & Optimization)

#### [Bài 11: Toán rời rạc, Độ phức tạp thuật toán & Tối ưu hóa lồi](./11_discrete_math_optimization.md)
- Tổ hợp/hoán vị cơ bản, lý thuyết đồ thị sơ lược (dùng trong decision tree, graph neural network), Big-O.
- Hàm lồi (convex function), tối ưu hóa lồi, **Gradient Descent** và các biến thể (Batch/Stochastic/Mini-batch).
- **Ứng dụng ML:** đây là thuật toán tối ưu dùng trong HẦU HẾT model ML/DL — nền tảng bắt buộc trước khi học Neural Network.

---

## Phần II — Machine Learning Cơ Bản

### [Bài 12: NumPy & Pandas cho Data Science](./12_numpy_pandas.md)
- NumPy array, vector hóa phép toán (liên hệ [Bài 2-3](./2_linear_algebra_vectors.md)), Pandas DataFrame, xử lý dữ liệu thiếu/trùng lặp.

### [Bài 13: Trực quan hóa dữ liệu & EDA (Exploratory Data Analysis)](./13_data_visualization_eda.md)
- Matplotlib/Seaborn, phân tích thống kê mô tả (liên hệ [Bài 9](./9_probability_distributions.md)), phát hiện outlier, correlation.

### [Bài 14: Linear Regression — Từ toán tới code](./14_linear_regression.md)
- Mô hình hồi quy tuyến tính, hàm mất mát MSE, Normal Equation (liên hệ [Bài 3](./3_linear_algebra_matrices.md)), Gradient Descent (liên hệ [Bài 11](./11_discrete_math_optimization.md)).

### [Bài 15: Logistic Regression & Classification cơ bản](./15_logistic_regression.md)
- Hàm sigmoid, Cross-Entropy Loss (liên hệ MLE — [Bài 10](./10_statistics_inference.md)), decision boundary.

### [Bài 16: Đánh giá mô hình (Model Evaluation)](./16_model_evaluation.md)
- Train/test/validation split, k-fold cross-validation, Confusion Matrix, Precision/Recall/F1, ROC-AUC, Overfitting vs Underfitting.

### [Bài 17: Regularization & Bias-Variance Tradeoff](./17_regularization.md)
- Ridge (L2), Lasso (L1), Elastic Net, Bias-Variance decomposition, vì sao regularization hoạt động (liên hệ MAP — [Bài 10](./10_statistics_inference.md)).

### [Bài 18: Decision Tree, Random Forest & Ensemble Methods](./18_trees_ensembles.md)
- Entropy/Gini Index (liên hệ [Bài 9](./9_probability_distributions.md)), Decision Tree, Bagging, Random Forest, Boosting (Gradient Boosting/XGBoost sơ lược).

### [Bài 19: SVM, k-NN, Naive Bayes & Unsupervised Learning (Clustering, PCA)](./19_svm_knn_clustering.md)
- Support Vector Machine (margin, kernel trick), k-Nearest Neighbors, Naive Bayes (liên hệ [Bài 8](./8_probability_basics.md)), k-Means Clustering, PCA thực hành (liên hệ [Bài 4](./4_linear_algebra_eigen_svd.md)).

---

## Phần III — Deep Learning

### [Bài 20: Neural Network cơ bản — Perceptron & MLP](./20_neural_networks.md)
- Perceptron, Multi-Layer Perceptron (MLP), hàm kích hoạt (Sigmoid/ReLU/Tanh/Softmax), forward propagation.

### [Bài 21: Backpropagation chi tiết (Toán + Code from scratch)](./21_backpropagation.md)
- Áp dụng trực tiếp chain rule ([Bài 6](./6_calculus_matrix_calculus.md)) để tính gradient qua nhiều lớp, cài đặt neural network from scratch bằng NumPy.

### [Bài 22: Optimizer nâng cao & Regularization trong Deep Learning](./22_optimizers_dl_regularization.md)
- SGD với Momentum, RMSProp, Adam, Learning Rate Scheduling, Dropout, Batch Normalization, Early Stopping.

### [Bài 23: CNN (Convolutional Neural Networks)](./23_cnn.md)
- Phép tích chập (convolution), pooling, kiến trúc CNN cơ bản, ứng dụng xử lý ảnh, transfer learning sơ lược.

### [Bài 24: RNN/LSTM & Giới thiệu Attention/Transformer](./24_rnn_transformer.md)
- Recurrent Neural Network, vanishing gradient, LSTM/GRU, cơ chế Attention, tổng quan kiến trúc Transformer.

---

## Phần IV — Dự Án Capstone

> **File chi tiết đầy đủ: [25_capstone_project.md](./25_capstone_project.md)**

**Đề bài:** Xây dựng pipeline ML hoàn chỉnh — từ EDA, tiền xử lý, huấn luyện nhiều model (Linear/Logistic Regression, Random Forest, Neural Network from scratch, rồi so sánh với PyTorch/scikit-learn), đánh giá, tuning, và đóng gói kết quả thành báo cáo + demo dự đoán — áp dụng toán từ Phần I xuyên suốt từng bước, đúng tinh thần "hiểu từ gốc rễ" thay vì chỉ gọi API.

## Gợi ý cách học
- **Đừng bỏ qua Phần I dù thấy "chỉ cần gọi `model.fit()`"** — hiểu toán giúp bạn debug được vì sao model không hội tụ, chọn đúng regularization, đọc được paper mới. Đây là ranh giới giữa "biết dùng thư viện" và "hiểu ML thật sự".
- Với mỗi công thức toán, tự tay làm ít nhất 1 ví dụ số nhỏ bằng tay (giấy bút) trước khi code — trực giác hình thành từ đó, không phải từ đọc code.
- Sau khi xong Phần I, nên implement **ít nhất 1 thuật toán ML from scratch chỉ bằng NumPy** (không dùng scikit-learn) — Linear Regression hoặc Neural Network — để chắc chắn hiểu toán đã "ngấm" thành code.
- Toán rời rạc/Big-O ([Bài 11](./11_discrete_math_optimization.md)) liên hệ trực tiếp kỹ năng thuật toán bạn đã có từ [Go](../Go/ROADMAP.md)/[Python](../Python/ROADMAP.md) — tận dụng nền đó thay vì học lại từ 0.

## Tài liệu tham khảo
- Mathematics for Machine Learning (Deisenroth, Faisal, Ong) — sách miễn phí: https://mml-book.github.io/
- 3Blue1Brown — Essence of Linear Algebra / Neural Networks (YouTube, trực giác hình học rất tốt)
- StatQuest with Josh Starmer (YouTube) — giải thích thống kê/ML trực quan
- scikit-learn docs: https://scikit-learn.org/stable/
- PyTorch docs: https://pytorch.org/docs/
- Deep Learning Book (Goodfellow, Bengio, Courville) — https://www.deeplearningbook.org/
