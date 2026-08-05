# Bài 4: Branching Cơ Bản

## Mục tiêu
- Hiểu branch là gì về bản chất (chỉ là 1 con trỏ nhẹ, không phải bản sao toàn bộ code).
- Thành thạo tạo, chuyển, xóa branch, và merge cơ bản (fast-forward vs 3-way merge).

## 1. Branch là gì — bản chất cực kỳ đơn giản

Một branch chỉ là **1 con trỏ di động (movable pointer)** trỏ tới 1 commit cụ thể — KHÔNG phải bản sao toàn bộ code như nhiều người hình dung ban đầu. Khi bạn commit trên 1 branch, Git tạo commit mới rồi di chuyển con trỏ branch đó tới commit mới — cực kỳ nhẹ và nhanh (ngược lại hoàn toàn với branching tốn kém của SVN). Chi tiết cơ chế đầy đủ ở [Bài 17: Git Internals](./17_git_internals.md).

```
main:     A---B---C
                    \
feature:             D---E   (con trỏ "feature" trỏ tới E, "main" vẫn trỏ tới C)
```

## 2. Tạo & chuyển branch

```powershell
git branch feature-login        # TẠO branch mới, nhưng KHÔNG tự chuyển sang nó
git switch feature-login          # CHUYỂN sang branch đó (lệnh hiện đại, khuyến khích)
git checkout feature-login         # cách cũ hơn, vẫn phổ biến, làm cùng việc

git switch -c feature-login          # tạo VÀ chuyển sang branch mới cùng lúc (phổ biến nhất)
git checkout -b feature-login          # tương đương, cách cũ

git branch                              # liệt kê mọi branch local, dấu * chỉ branch hiện tại
git branch -a                             # liệt kê CẢ branch remote-tracking
```

**Vì sao có cả `switch`/`checkout`?** `checkout` lịch sử làm QUÁ NHIỀU việc khác nhau (chuyển branch, khôi phục file, chuyển tới commit cụ thể — dễ gây nhầm lẫn). Git giới thiệu `switch` (chỉ chuyển branch) và `restore` (chỉ khôi phục file — [Bài 10](./10_undo_changes.md)) để tách rõ trách nhiệm — nên ưu tiên dùng 2 lệnh mới này cho code rõ ràng, dù `checkout` vẫn hoạt động và còn phổ biến trong tài liệu/tutorial cũ.

## 3. Xóa branch

```powershell
git branch -d feature-login    # xóa AN TOÀN — Git từ chối nếu branch chưa merge (tránh mất commit)
git branch -D feature-login     # xóa CƯỠNG BỨC — dùng khi chắc chắn muốn bỏ, kể cả chưa merge
```

## 4. Merge — 2 cơ chế khác nhau hoàn toàn

### Fast-forward merge — khi branch đích KHÔNG có commit mới nào kể từ lúc tách nhánh

```
main:     A---B
                \
feature:         C---D

# sau "git switch main; git merge feature"
main:     A---B---C---D   (chỉ đơn giản DI CHUYỂN con trỏ main tới D, không tạo commit mới)
```

```powershell
git switch main
git merge feature-login
```

### 3-way merge — khi CẢ 2 branch đều có commit mới kể từ lúc tách nhánh

```
main:     A---B-------E
                \     /
feature:         C---D

# Git tạo 1 COMMIT MERGE MỚI (M) có 2 parent, kết hợp thay đổi từ cả 2 nhánh
main:     A---B-------E---M
                \     /   /
feature:         C---D---'
```

Git tìm **common ancestor** (tổ tiên chung — ở đây là B) của 2 branch, so sánh cả 2 nhánh với tổ tiên đó, rồi tự động kết hợp các thay đổi KHÔNG xung đột (dùng thuật toán 3-way merge — so sánh 3 phiên bản: ancestor, branch 1, branch 2). Nếu 2 branch sửa CÙNG 1 dòng theo cách khác nhau, xảy ra **conflict** — cần giải quyết thủ công, chi tiết ở [Bài 8](./8_conflicts.md).

## 5. So sánh Merge vs Rebase — xem trước, chi tiết đầy đủ ở Bài 7

`git merge` giữ nguyên lịch sử THẬT (kể cả các nhánh rẽ), tạo merge commit. `git rebase` "viết lại" lịch sử để trông như thẳng hàng (linear) — cả 2 đều hợp lệ, khác nhau về TRIẾT LÝ quản lý lịch sử, sẽ so sánh kỹ ở [Bài 7](./7_merge_vs_rebase.md).

## Ví dụ đầy đủ

```powershell
git init branching-demo; cd branching-demo
echo "line 1" > file.txt
git add file.txt; git commit -m "Initial commit"

git switch -c feature-a
echo "line 2 from feature-a" >> file.txt
git add file.txt; git commit -m "Add: dòng từ feature-a"

git switch main
echo "line 2 from main" >> file.txt
git add file.txt; git commit -m "Add: dòng từ main"

git merge feature-a   # 2 branch đều có commit mới -> xảy ra CONFLICT (cùng sửa dòng 2)
# giải quyết conflict thủ công (Bài 8), rồi:
git add file.txt
git commit -m "Merge: kết hợp feature-a vào main"

git log --oneline --graph --all   # xem hình dạng lịch sử sau merge
```

## Bài tập

1. **Fast-forward merge**: tạo branch mới, commit vài lần, quay lại `main` (KHÔNG commit gì thêm trên `main`), merge — verify bằng `git log --graph` rằng đây là fast-forward (không có merge commit).
2. **3-way merge KHÔNG conflict**: tạo 2 branch từ cùng điểm, mỗi branch sửa 1 FILE KHÁC NHAU, merge — quan sát Git tự động kết hợp không cần can thiệp.
3. **3-way merge CÓ conflict**: làm lại ví dụ đầy đủ ở cuối bài, tự giải quyết conflict, hoàn thành merge.
4. **Xóa branch an toàn vs cưỡng bức**: tạo 1 branch, commit, KHÔNG merge, thử `git branch -d` (sẽ bị từ chối), giải thích tại sao, rồi dùng `-D` để xóa cưỡng bức.

## Tiếp theo
→ [Bài 5: Làm việc với Remote](./5_remote.md)
