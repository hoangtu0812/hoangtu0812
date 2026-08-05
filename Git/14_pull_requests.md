# Bài 14: Pull Request & Code Review

## Mục tiêu
- Hiểu quy trình Pull Request (PR) từ đầu tới cuối.
- Biết chọn đúng kiểu merge trên GitHub (Merge commit / Squash / Rebase).
- Viết PR description và review comment hiệu quả.

## 1. Pull Request là gì?

PR là 1 khái niệm của **GitHub** (GitLab gọi là "Merge Request" — cùng ý tưởng), KHÔNG phải lệnh Git thuần — nó là lớp quy trình được xây trên nền Git: đề xuất merge 1 branch vào branch khác, kèm không gian để review, comment, chạy CI tự động ([Bài 20](./20_ci_cd_integration.md)) TRƯỚC KHI merge thật sự xảy ra.

## 2. Quy trình PR đầy đủ

```
1. Tạo branch từ main, code, push lên remote (Bài 13 mục 6)
2. Mở PR trên GitHub: chọn branch nguồn -> branch đích, viết description
3. CI tự động chạy (test, lint, build)
4. Reviewer đọc code, để lại comment/suggestion
5. Tác giả sửa theo góp ý, push thêm commit (PR TỰ ĐỘNG cập nhật, không cần tạo PR mới)
6. Reviewer approve
7. Merge vào branch đích (chọn 1 trong 3 kiểu — mục 4)
8. Xóa branch đã merge (thường có nút tự động trên GitHub)
```

## 3. Viết PR Description tốt

```markdown
## Tóm tắt
Thêm chức năng đăng xuất, tự động clear session khi user click nút "Đăng xuất".

## Thay đổi chính
- Thêm nút "Đăng xuất" ở header
- Thêm API endpoint POST /auth/logout
- Clear localStorage token khi logout thành công

## Cách test
1. Đăng nhập bằng tài khoản test
2. Click "Đăng xuất" ở header
3. Verify redirect về trang login, token đã bị xóa

## Ảnh chụp màn hình (nếu có thay đổi UI)
[đính kèm ảnh]
```

Description tốt giúp reviewer hiểu NGỮ CẢNH (tại sao cần thay đổi này) mà không phải tự đọc toàn bộ diff để đoán — đặc biệt quan trọng khi reviewer không tham gia từ đầu.

## 4. 3 kiểu merge trên GitHub — khác biệt về LỊCH SỬ tạo ra

**Create a merge commit** — giống `git merge` thường ([Bài 4 mục 4](./4_branching.md)): giữ nguyên MỌI commit của branch, thêm 1 merge commit. Lịch sử đầy đủ nhất nhưng có thể "nhiễu" nếu branch có nhiều commit vụn vặt.

**Squash and merge** — gộp TOÀN BỘ commit của PR thành **1 commit DUY NHẤT** trên branch đích, dùng title/description của PR làm message. Lịch sử `main` SẠCH NHẤT — mỗi PR = đúng 1 commit, dễ theo dõi, dễ revert nguyên 1 tính năng nếu cần. **Đây là lựa chọn phổ biến nhất** cho GitHub Flow ([Bài 13 mục 3](./13_workflows.md)), đặc biệt khi commit trong quá trình code không cần giữ chi tiết (đã có PR description đầy đủ thay thế).

**Rebase and merge** — áp dụng rebase ([Bài 7 mục 2](./7_merge_vs_rebase.md)) từng commit của branch lên đầu branch đích, KHÔNG tạo merge commit, nhưng VẪN GIỮ từng commit riêng lẻ (khác squash gộp thành 1). Phù hợp khi từng commit trong PR đã được viết CẨN THẬN, có ý nghĩa riêng đáng giữ lại.

| Kiểu merge | Số commit thêm vào main | Khi nào dùng |
|---|---|---|
| Merge commit | Toàn bộ + 1 merge commit | Muốn giữ dấu vết đầy đủ nhất, ít quan tâm lịch sử gọn |
| Squash and merge | 1 | **Mặc định khuyến khích** — lịch sử main sạch, mỗi PR = 1 dòng log |
| Rebase and merge | Toàn bộ, không merge commit | Từng commit trong PR đã chỉn chu, đáng giữ riêng |

## 5. Code Review — góc nhìn người REVIEW

Review tập trung vào: logic có đúng không, có edge case nào bị bỏ sót, code có dễ đọc/bảo trì không, có vi phạm convention team không, có thiếu test không — KHÔNG nên chỉ soi lỗi format (nên để linter/formatter tự động làm — liên hệ [Bài 16: Git Hooks](./16_git_hooks.md)).

```markdown
<!-- Comment kiểu "suggestion" trên GitHub — reviewer có thể đề xuất sửa TRỰC TIẾP -->
```suggestion
if user is None:
    raise ValueError("User không tồn tại")
```

<!-- Comment dạng câu hỏi, không áp đặt -->
Hàm này có xử lý trường hợp `list` rỗng không? Mình thấy chưa thấy check ở đây.
```

**Nguyên tắc phản hồi tốt:** đặt câu hỏi thay vì ra lệnh khi không chắc chắn ý đồ tác giả; khen ngợi phần code tốt, không chỉ nêu vấn đề; phân biệt rõ "phải sửa" (blocking) và "gợi ý, tùy chọn" (nitpick) để tác giả biết ưu tiên gì trước.

## 6. Góc nhìn người ĐƯỢC review

Không nên coi comment là "chỉ trích cá nhân" — mục tiêu chung là code tốt hơn. Trả lời rõ ràng khi KHÔNG đồng ý với góp ý (giải thích lý do, không chỉ im lặng bỏ qua); resolve conversation sau khi đã sửa hoặc đã giải thích thỏa đáng; nếu PR quá lớn khó review, cân nhắc tách nhỏ thành nhiều PR độc lập.

## 7. Branch Protection Rule — bắt buộc quy trình qua PR

Trên GitHub (Settings → Branches → Add rule), có thể thiết lập cho `main`:
- Yêu cầu ít nhất 1 (hoặc N) approval trước khi merge.
- Yêu cầu CI pass (status check) trước khi merge được phép bấm.
- Không cho phép push trực tiếp vào `main` (bắt buộc qua PR, kể cả admin nếu muốn).

Đây là cách "ép" quy trình PR ([Bài 13](./13_workflows.md)) được tuân thủ THẬT SỰ, không chỉ là quy ước bằng lời — chi tiết kết hợp với CI ở [Bài 20](./20_ci_cd_integration.md).

## Bài tập

1. **Tạo PR thật trên GitHub**: tạo 1 branch, push, mở PR với description đầy đủ theo mẫu ở mục 3.
2. **Thử cả 3 kiểu merge**: tạo 3 PR nhỏ tương tự nhau, merge mỗi cái bằng 1 kiểu khác nhau (Merge commit/Squash/Rebase), so sánh `git log --graph` của `main` sau đó.
3. **Review PR của chính mình**: tự đóng vai reviewer cho PR của bạn (hoặc nhờ đồng nghiệp/bạn học), để lại ít nhất 3 comment (1 "suggestion" code, 1 câu hỏi, 1 lời khen) theo đúng tinh thần mục 5.
4. **Thiết lập Branch Protection**: bật rule yêu cầu PR + status check cho `main` trên 1 repo test, thử push trực tiếp vào `main` để verify bị chặn.

## Tiếp theo
→ [Bài 15: Submodule & Subtree](./15_submodules_subtrees.md)
