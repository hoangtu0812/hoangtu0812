# Bài 9: Stash & Cherry-pick

## Mục tiêu
- Dùng `git stash` để tạm "cất" thay đổi dở dang mà không cần commit.
- Dùng `git cherry-pick` để lấy 1 commit cụ thể từ branch khác.

## 1. `git stash` — cất thay đổi dở dang, quay lại trạng thái sạch tạm thời

**Tình huống điển hình:** đang sửa dở `feature.py`, sếp/đồng nghiệp báo có bug khẩn cấp cần fix ngay trên `main` — nhưng thay đổi hiện tại CHƯA sẵn sàng để commit. `git stash` giải quyết đúng vấn đề này:

```powershell
git stash                       # cất TOÀN BỘ thay đổi (đã track) chưa commit, trả Working Directory về sạch
git stash push -m "Đang làm dở tính năng lọc"   # cất kèm message mô tả — nên dùng để dễ nhớ sau này

git switch main
# ... fix bug khẩn cấp, commit, push ...

git switch feature-branch-cu
git stash list                    # xem danh sách các stash đã cất (có thể nhiều)
git stash pop                       # LẤY LẠI stash gần nhất, ĐỒNG THỜI xóa nó khỏi danh sách stash
git stash apply                      # LẤY LẠI nhưng GIỮ NGUYÊN trong danh sách stash (dùng khi muốn áp dụng ở NHIỀU branch)
```

## 2. Các thao tác khác với stash

```powershell
git stash list                     # liệt kê mọi stash, dạng stash@{0}, stash@{1}...
git stash show stash@{0}             # xem tóm tắt thay đổi trong 1 stash cụ thể
git stash show -p stash@{0}           # xem ĐẦY ĐỦ diff trong stash đó

git stash drop stash@{0}                # xóa 1 stash cụ thể (không áp dụng lại)
git stash clear                            # xóa TOÀN BỘ stash — cẩn trọng, không thể hoàn tác

git stash -u                                # bao gồm CẢ file untracked (mặc định stash BỎ QUA file chưa track)
git stash push -- file1.txt file2.txt         # chỉ stash 1 số file cụ thể, không phải toàn bộ
```

## 3. Stash hoạt động thế nào về bản chất

Stash thực chất tạo ra các **commit đặc biệt** (không nằm trên branch nào, được tham chiếu qua 1 ref riêng `refs/stash`) lưu trạng thái Working Directory + Staging Area tại thời điểm đó — đây là lý do stash tồn tại được dù bạn chuyển qua chuyển lại nhiều branch, và vì sao nó "biến mất" nếu bạn xóa `.git` hoàn toàn (nó vẫn là dữ liệu Git bình thường, chỉ không gắn với branch cụ thể).

## 4. `git cherry-pick` — lấy CHÍNH XÁC 1 commit từ branch khác

**Tình huống điển hình:** 1 commit fix bug quan trọng đã nằm trên branch `feature-x`, nhưng bạn cần áp dụng NGAY commit đó vào `main` (hoặc branch release khác) mà KHÔNG merge toàn bộ `feature-x` (vì nó còn nhiều thay đổi khác chưa sẵn sàng).

```powershell
git log feature-x --oneline     # tìm hash của commit cần lấy, vd a1b2c3d

git switch main
git cherry-pick a1b2c3d           # áp dụng ĐÚNG thay đổi của commit đó lên main, tạo commit MỚI (hash khác)
```

Cherry-pick về bản chất tương tự 1 bước của rebase ([Bài 7 mục 2](./7_merge_vs_rebase.md)) — lấy nội dung thay đổi (diff) của 1 commit, áp lên nền hiện tại, tạo commit mới với parent khác.

## 5. Cherry-pick nhiều commit & xử lý conflict

```powershell
git cherry-pick a1b2c3d e4f5g6h        # nhiều commit cùng lúc, theo đúng thứ tự liệt kê
git cherry-pick a1b2c3d^..e4f5g6h        # 1 dải commit liên tiếp (LƯU Ý: commit ĐẦU dùng ^ để bao gồm chính nó)

# Nếu cherry-pick gây conflict — xử lý TƯƠNG TỰ merge/rebase (Bài 8)
git cherry-pick a1b2c3d
# (giải quyết conflict thủ công)
git add file.txt
git cherry-pick --continue

git cherry-pick --abort    # hủy bỏ nếu cần
```

## 6. Khi nào dùng cherry-pick — và khi nào KHÔNG nên

**Nên dùng:** áp dụng hotfix khẩn cấp từ 1 branch sang branch release khác; khôi phục 1 commit đã lỡ xóa nhầm (tìm lại qua `git reflog` — [Bài 12](./12_rewriting_history.md)) vào đúng branch cần.

**Không nên lạm dụng:** cherry-pick nhiều commit liên tục thay vì merge cả branch tạo ra lịch sử "trùng lặp nội dung nhưng khác hash" — dễ gây nhầm lẫn khi merge branch gốc sau này (Git có thể coi đó là thay đổi khác, gây conflict không cần thiết). Nếu cần lấy TOÀN BỘ 1 branch, hãy merge/rebase ([Bài 7](./7_merge_vs_rebase.md)) thay vì cherry-pick từng commit một.

## Ví dụ đầy đủ

```powershell
git init stash-cherry-demo; cd stash-cherry-demo
echo "v1" > app.py; git add .; git commit -m "Initial"

git switch -c feature
echo "dang lam dang do" >> app.py    # KHÔNG commit

# Bug khẩn cấp xuất hiện!
git stash push -m "feature dang do"
git switch main
echo "hotfix" > hotfix.py
git add .; git commit -m "Hotfix khẩn cấp"

# Quay lại làm tiếp feature
git switch feature
git stash pop
cat app.py    # thấy lại "dang lam dang do"

# Giờ muốn áp dụng hotfix đó sang 1 branch release khác
git switch -c release
git log main --oneline    # tìm hash commit "Hotfix khẩn cấp"
git cherry-pick <hash-cua-hotfix>
```

## Bài tập

1. **Stash cơ bản**: làm theo đúng "Ví dụ đầy đủ" phần stash, verify `stash pop` khôi phục đúng thay đổi dở dang.
2. **Stash nhiều lần**: tạo 2-3 stash khác nhau (dùng message rõ ràng), dùng `git stash list`/`git stash show` để phân biệt, áp dụng đúng cái cần bằng `git stash apply stash@{n}`.
3. **Cherry-pick 1 commit**: làm theo "Ví dụ đầy đủ" phần cherry-pick, verify commit trên `release` có CÙNG nội dung nhưng KHÁC hash so với commit gốc trên `main`.
4. **Cherry-pick gây conflict**: tạo tình huống cherry-pick 1 commit mà vùng nó sửa đã bị thay đổi khác ở nhánh đích — giải quyết conflict, hoàn thành bằng `--continue`.

## Tiếp theo
→ [Bài 10: Undo Changes — reset, revert, restore](./10_undo_changes.md)
