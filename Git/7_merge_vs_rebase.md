# Bài 7: Merge vs Rebase

## Mục tiêu
- Hiểu sâu 2 cách tích hợp thay đổi giữa các branch, khác biệt về LỊCH SỬ tạo ra.
- Biết khi nào nên dùng cái nào — đây là quyết định TRIẾT LÝ, không có đáp án "đúng tuyệt đối".

## 1. Nhắc lại Merge — giữ nguyên lịch sử THẬT

Đã giới thiệu ở [Bài 4 mục 4](./4_branching.md): merge tạo 1 commit mới có **2 parent**, giữ nguyên dấu vết rằng 2 nhánh đã từng tách ra và giờ hợp lại.

```
main:     A---B-------E---M     (M = merge commit, có 2 parent: E và D)
                \     /   /
feature:         C---D---'
```

**Ưu điểm:** lịch sử phản ánh CHÍNH XÁC những gì đã xảy ra (branch nào, khi nào tách, khi nào hợp) — không có rủi ro mất mát thông tin. **Nhược điểm:** với nhiều branch/merge, `git log --graph` có thể trông rối rắm, khó theo dõi tuyến tính.

## 2. Rebase — "viết lại" lịch sử để trông THẲNG HÀNG

```powershell
git switch feature
git rebase main
```

```
Trước rebase:
main:     A---B---E
                \
feature:         C---D

Sau "git rebase main" (đang đứng trên feature):
main:     A---B---E
                    \
feature:             C'---D'   (C', D' là commit MỚI, nội dung giống C,D nhưng lấy E làm nền, hash khác hoàn toàn)
```

**Cơ chế:** Git tạm "gỡ" các commit riêng của `feature` (C, D) ra, di chuyển điểm nhánh `feature` tới ngay sau `E` (đầu mới nhất của `main`), rồi "áp" lại TỪNG commit đó lên trên (theo đúng thứ tự) — mỗi commit áp lại thực chất là 1 commit MỚI (nội dung thay đổi giống hệt nhưng có parent khác, nên **hash khác hoàn toàn** — nhắc lại nội dung + hash gắn chặt, xem [Bài 17](./17_git_internals.md)).

**Ưu điểm:** lịch sử SẠCH, thẳng hàng (linear) — dễ đọc `git log --oneline` như 1 câu chuyện liền mạch, không có merge commit "nhiễu". **Nhược điểm:** VIẾT LẠI lịch sử — commit cũ (C, D) bị "thay thế" bằng commit mới (C', D') có hash khác. Đây là lý do dẫn tới **quy tắc vàng của rebase**.

## 3. Quy tắc vàng: KHÔNG BAO GIỜ rebase commit ĐÃ PUSH lên remote mà người khác có thể đã dựa vào

Nếu bạn rebase commit đã push, và người khác đã `pull` những commit CŨ (C, D) về máy họ, sau khi bạn force-push commit MỚI (C', D'), lịch sử của họ và của bạn hoàn toàn "lệch pha" — dẫn tới tình huống rất khó gỡ (commit trùng lặp, conflict hàng loạt). Quy tắc an toàn: **chỉ rebase commit CÒN Ở LOCAL, chưa từng push, hoặc push lên 1 branch riêng CHỈ MÌNH BẠN dùng** (feature branch cá nhân).

## 4. So sánh trực tiếp qua ví dụ

```powershell
# Cùng 1 tình huống: main có thêm 1 commit trong lúc bạn làm feature branch

git init compare-demo; cd compare-demo
echo "1" > f.txt; git add .; git commit -m "commit 1"
git switch -c feature
echo "2" > f.txt; git add .; git commit -m "commit 2 (feature)"
git switch main
echo "1b" > f.txt; git add .; git commit -m "commit 1b (main)"

# Cách A: merge
git switch feature
git switch -c feature-merge-test
git merge main
git log --oneline --graph      # thấy merge commit, cấu trúc "hình thoi"

# Cách B: rebase (thử trên nhánh khác để so sánh, không ảnh hưởng nhánh gốc)
git switch feature
git switch -c feature-rebase-test
git rebase main
git log --oneline --graph      # thấy lịch sử THẲNG HÀNG, không có merge commit
```

## 5. `git pull --rebase` — áp dụng rebase khi đồng bộ với remote

```powershell
git pull --rebase origin main
```

Thay vì `fetch + merge` (mặc định — [Bài 5 mục 6](./5_remote.md)), lệnh này làm `fetch + rebase` — các commit local CHƯA push được "đặt lên trên" các commit mới từ remote, giữ lịch sử thẳng hàng thay vì tạo merge commit mỗi lần đồng bộ. Nhiều team đặt `pull.rebase=true` làm mặc định để tránh merge commit "rác" xuất hiện liên tục.

## 6. Interactive Rebase — xem trước, chi tiết đầy đủ ở Bài 12

`git rebase -i` không chỉ dùng để tích hợp branch khác, mà còn để **sửa lại lịch sử commit của chính mình** trước khi chia sẻ (gộp, xóa, sửa message nhiều commit cùng lúc) — chi tiết đầy đủ ở [Bài 12](./12_rewriting_history.md).

## 7. Bảng quyết định nhanh: khi nào dùng gì

| Tình huống | Nên dùng |
|---|---|
| Merge feature branch đã hoàn thành vào `main` (qua Pull Request) | **Merge** — giữ dấu vết lịch sử feature đã tồn tại, đúng chuẩn hầu hết team dùng GitHub Flow ([Bài 13](./13_workflows.md)) |
| Cập nhật feature branch CỦA RIÊNG BẠN với thay đổi mới nhất từ `main` (trước khi tạo PR) | **Rebase** — giữ lịch sử sạch, dễ review hơn khi PR chỉ hiện đúng commit của bạn |
| Dọn dẹp lịch sử commit lộn xộn TRƯỚC KHI push lần đầu | **Interactive Rebase** ([Bài 12](./12_rewriting_history.md)) |
| Commit ĐÃ PUSH, người khác có thể đã pull | **KHÔNG rebase** — chỉ merge, hoặc revert nếu cần hoàn tác ([Bài 10](./10_undo_changes.md)) |

## Bài tập

1. **So sánh trực quan**: thực hiện đúng "Ví dụ" mục 4, so sánh 2 kết quả `git log --graph` giữa nhánh merge và nhánh rebase, mô tả bằng lời sự khác biệt.
2. **Giải quyết conflict trong rebase**: tạo tình huống 2 branch sửa CÙNG dòng, thử rebase — quan sát conflict xuất hiện TỪNG COMMIT một (khác cách merge chỉ conflict 1 lần), dùng `git rebase --continue`/`--abort` để xử lý.
3. **`pull --rebase`**: thiết lập `git config pull.rebase true` cho 1 repo test, thử pull khi có commit mới trên remote, verify lịch sử thẳng hàng.
4. **Thảo luận (viết ra suy nghĩ)**: với 1 team 5 người làm chung 1 project, đề xuất quy tắc dùng merge/rebase cho từng tình huống (feature branch cá nhân, merge vào main, hotfix khẩn cấp) — dựa trên bảng quyết định mục 7.

## Tiếp theo
→ [Bài 8: Giải quyết Conflict](./8_conflicts.md)
