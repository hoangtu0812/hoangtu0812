# Bài 13: Git Workflow Cho Team

## Mục tiêu
- Hiểu 3 workflow phổ biến nhất: Git Flow, GitHub Flow, Trunk-based Development.
- Biết chọn workflow phù hợp với quy mô team và tần suất release.

## 1. Vì sao cần "workflow" — không chỉ biết lệnh Git là đủ

Biết `merge`, `rebase`, `branch` ([Bài 4](./4_branching.md), [Bài 7](./7_merge_vs_rebase.md)) là điều kiện CẦN, nhưng team cần thêm 1 **quy ước chung**: đặt tên branch thế nào, khi nào tạo branch mới, khi nào merge vào đâu, ai review trước khi merge — nếu không, mỗi người làm 1 kiểu, lịch sử hỗn loạn, dễ xung đột quy trình dù không xung đột code.

## 2. Git Flow — quy trình đầy đủ, nhiều loại branch, phù hợp release theo chu kỳ

```
main        ────────────●─────────────●────────  (chỉ chứa code ĐÃ RELEASE)
                          \           /
release/1.2                ●───●───●              (chuẩn bị release, chỉ fix bug nhỏ)
                           /         \
develop     ──●───●───●──●───────────●───●──────  (nhánh tích hợp chính, code đang phát triển)
               \       \
feature/login   ●───●   \
                          feature/payment ●───●
```

- **`main`**: LUÔN ở trạng thái production-ready, mỗi commit tương ứng 1 bản đã release.
- **`develop`**: nhánh tích hợp, nơi các feature branch merge vào trước khi release.
- **`feature/*`**: mỗi tính năng 1 branch riêng, tách từ `develop`, merge lại vào `develop`.
- **`release/*`**: khi `develop` đủ tính năng cho 1 bản release, tạo branch riêng để ổn định (chỉ fix bug, không thêm tính năng mới), rồi merge vào CẢ `main` và `develop`.
- **`hotfix/*`**: fix khẩn cấp trực tiếp từ `main` (bug production), merge lại vào CẢ `main` và `develop`.

**Phù hợp:** team lớn, sản phẩm release theo lịch cố định (vd mỗi tháng 1 bản), cần hỗ trợ nhiều phiên bản cùng lúc (khách hàng dùng bản cũ vẫn cần fix). **Nhược điểm:** phức tạp, nhiều loại branch dễ gây nhầm lẫn cho team nhỏ hoặc release liên tục.

## 3. GitHub Flow — đơn giản hơn nhiều, phù hợp deploy liên tục

```
main   ──●───────●───────────●──────────  (LUÔN deploy được, mọi merge vào đây có thể lên production ngay)
          \       \           \
feature-a  ●──●    \           feature-c ●──●
                     feature-b ●───●───●
```

Chỉ có **1 loại branch dài hạn**: `main`. Mọi tính năng đi qua: tạo branch từ `main` → code → mở Pull Request → review → merge thẳng vào `main` → deploy ngay (thường tự động qua CI/CD — [Bài 20](./20_ci_cd_integration.md)).

**Phù hợp:** team nhỏ-vừa, sản phẩm web/SaaS deploy liên tục (nhiều lần/ngày), không cần hỗ trợ song song nhiều phiên bản cũ. Đây là workflow phổ biến NHẤT cho startup và team hiện đại, đơn giản, ít overhead.

## 4. Trunk-based Development — cực kỳ đơn giản, branch sống RẤT ngắn

Tương tự GitHub Flow nhưng nhấn mạnh: branch (nếu có) chỉ tồn tại **vài giờ tới 1-2 ngày** trước khi merge lại vào `main` (trunk) — khuyến khích commit thẳng vào `main` nếu thay đổi đủ nhỏ, dùng **feature flag** (bật/tắt tính năng bằng cấu hình runtime, không phải bằng branch) để "ẩn" tính năng chưa hoàn thiện dù code đã nằm trên `main`.

**Phù hợp:** team có kỷ luật CI/CD rất mạnh (test tự động đầy đủ, deploy liên tục nhiều lần/ngày), thường thấy ở các công ty công nghệ lớn (Google, Facebook) — giảm thiểu tối đa "merge hell" vì branch không bao giờ sống đủ lâu để phân kỳ nhiều với `main`.

## 5. Bảng so sánh nhanh

| | Git Flow | GitHub Flow | Trunk-based |
|---|---|---|---|
| Số loại branch dài hạn | Nhiều (main, develop, release, hotfix) | 1 (main) | 1 (main/trunk) |
| Độ phức tạp | Cao | Thấp | Rất thấp |
| Tần suất release | Theo chu kỳ (vd hàng tháng) | Liên tục | Liên tục, rất thường xuyên |
| Hỗ trợ nhiều version song song | Tốt | Kém | Kém |
| Yêu cầu CI/CD mạnh | Không bắt buộc | Khuyến khích | Bắt buộc |
| Phù hợp | Team lớn, sản phẩm đóng gói/enterprise | Team nhỏ-vừa, web/SaaS | Team lớn với văn hóa kỹ thuật rất mạnh |

## 6. Quy ước đặt tên branch — dù chọn workflow nào

```
feature/ten-tinh-nang       # tính năng mới
fix/mo-ta-bug                  # sửa bug
hotfix/mo-ta-khan-cap            # fix khẩn cấp production
chore/cap-nhat-dependency          # việc vặt, không phải tính năng/bug
docs/cap-nhat-readme                 # chỉ sửa tài liệu
```

Tên branch RÕ RÀNG giúp mọi người (kể cả không phải người tạo) hiểu ngay branch đó làm gì khi nhìn `git branch -a`.

## Ví dụ đầy đủ: GitHub Flow thực hành

```powershell
git switch main
git pull origin main               # LUÔN đồng bộ main mới nhất trước khi tạo branch mới

git switch -c feature/them-nut-dang-xuat
# ... code ...
git add .; git commit -m "Add: nút đăng xuất ở header"
git push -u origin feature/them-nut-dang-xuat

# Trên GitHub: mở Pull Request từ feature/them-nut-dang-xuat vào main
# Team review, CI chạy test tự động (Bài 20)
# Sau khi approve: merge vào main (thường qua nút "Squash and merge" trên GitHub)

git switch main
git pull origin main                # đồng bộ lại, giờ đã có tính năng vừa merge
git branch -d feature/them-nut-dang-xuat   # xóa branch local đã merge xong
```

## Bài tập

1. **Mô phỏng GitHub Flow**: thực hành đúng "Ví dụ đầy đủ" cho 2 tính năng khác nhau, mỗi tính năng đi qua đủ quy trình branch → PR (dùng GitHub thật) → merge.
2. **So sánh workflow**: với 3 tình huống team khác nhau (team 3 người làm app nội bộ deploy hàng tuần; team 50 người làm phần mềm đóng gói bán cho doanh nghiệp, release mỗi quý; team 200 người tại 1 công ty công nghệ lớn deploy hàng chục lần/ngày) — đề xuất workflow phù hợp nhất cho mỗi trường hợp, giải thích lý do.
3. **Vẽ sơ đồ Git Flow**: tự vẽ lại sơ đồ Git Flow ở mục 2 (không nhìn lại bài), chỉ rõ branch nào merge vào branch nào.
4. **Quy ước đặt tên branch**: viết quy ước đặt tên branch cho 1 project giả định của bạn, áp dụng nó khi tạo 3-4 branch thử nghiệm.

## Tiếp theo
→ [Bài 14: Pull Request & Code Review](./14_pull_requests.md)
