# Bài 10: Undo Changes — reset, revert, restore

## Mục tiêu
- Phân biệt RÕ RÀNG 4 lệnh "hoàn tác" dễ nhầm lẫn nhất của Git: `reset` (soft/mixed/hard), `revert`, `restore`, `checkout -- file`.
- Biết chọn đúng lệnh cho đúng tình huống — đây là kỹ năng quan trọng nhất để "yên tâm" dùng Git.

## 1. Bảng tổng quan — chọn đúng lệnh cho đúng nhu cầu

| Muốn làm gì | Dùng lệnh |
|---|---|
| Bỏ thay đổi CHƯA staged trong 1 file (về lại như commit gần nhất) | `git restore <file>` |
| Bỏ staging (đưa file từ Staging Area về lại Modified, KHÔNG mất thay đổi) | `git restore --staged <file>` |
| Di chuyển HEAD/branch về commit cũ, GIỮ thay đổi ở Staging Area | `git reset --soft <commit>` |
| Di chuyển HEAD/branch về commit cũ, GIỮ thay đổi ở Working Directory (bỏ staging) | `git reset --mixed <commit>` (mặc định) |
| Di chuyển HEAD/branch về commit cũ, XÓA SẠCH mọi thay đổi | `git reset --hard <commit>` |
| Hoàn tác 1 commit ĐÃ PUSH mà KHÔNG xóa lịch sử | `git revert <commit>` |

## 2. `git restore` — hoàn tác ở cấp độ Working Directory/Staging Area (không di chuyển commit)

```powershell
git restore file.txt              # bỏ thay đổi CHƯA staged, khôi phục về bản ở commit gần nhất
git restore --staged file.txt       # bỏ staging (unstage), file vẫn GIỮ thay đổi, chỉ về lại "Modified"
git restore --source=HEAD~2 file.txt  # khôi phục file về ĐÚNG trạng thái ở 2 commit trước
```

`restore` chỉ thao tác trên **1 file cụ thể**, KHÔNG di chuyển `HEAD`/branch — an toàn hơn `reset` cho các tình huống chỉ cần sửa 1 file.

## 3. `git reset` — di chuyển con trỏ branch, 3 mức độ ảnh hưởng

Nhắc lại `HEAD` ([Bài 2 mục 5](./2_core_concepts.md)) là con trỏ trỏ commit hiện tại. `git reset <commit>` di chuyển `HEAD` (và branch đang đứng) TỚI 1 commit khác — khác nhau ở việc XỬ LÝ những thay đổi đã có SAU commit đó thế nào:

```
Trước reset:  A---B---C  (HEAD tại C)
git reset B:  A---B      (HEAD tại B — thay đổi của C đi đâu tùy loại reset)
```

**`--soft`**: di chuyển `HEAD`, nhưng GIỮ NGUYÊN mọi thay đổi ở **Staging Area** — như thể bạn chưa từng commit C, nhưng đã `git add` sẵn. Dùng khi muốn "gộp lại" nhiều commit thành 1 (soft reset về commit trước đó, rồi commit lại 1 lần).

```powershell
git reset --soft HEAD~1    # "undo" commit gần nhất, giữ thay đổi ở staging, sẵn sàng commit lại
```

**`--mixed`** (mặc định nếu không chỉ định gì): di chuyển `HEAD`, đưa thay đổi về **Working Directory** (bỏ staging). Dùng khi muốn "commit lại từ đầu, chia nhỏ khác đi".

```powershell
git reset HEAD~1    # tương đương --mixed, thay đổi của commit gần nhất giờ ở Working Directory, CHƯA staged
```

**`--hard`**: di chuyển `HEAD`, **XÓA SẠCH** mọi thay đổi (cả staged lẫn working directory) — CỰC KỲ NGUY HIỂM nếu chưa commit/backup, vì thay đổi bị mất KHÔNG THỂ khôi phục qua thao tác Git thông thường (chỉ có thể cứu qua `reflog` — [Bài 12](./12_rewriting_history.md), và chỉ khi thay đổi đó ĐÃ từng được commit ở đâu đó).

```powershell
git reset --hard HEAD~1    # XÓA VĨNH VIỄN commit gần nhất VÀ mọi thay đổi liên quan — cẩn trọng tuyệt đối
```

## 4. `reset` trên commit ĐÃ PUSH — vì sao nguy hiểm

`reset` (đặc biệt `--hard`) thay đổi vị trí con trỏ branch — nếu bạn đã push commit đó lên remote và giờ `reset` rồi push lại, bạn cần **force push** để "ép" remote khớp lại — điều này XÓA MẤT commit đó trên remote, ảnh hưởng tới bất kỳ ai đã pull nó về (giống vấn đề rebase đã nêu ở [Bài 7 mục 3](./7_merge_vs_rebase.md)). **Quy tắc:** chỉ `reset` thoải mái trên commit CÒN Ở LOCAL, chưa push.

## 5. `git revert` — cách AN TOÀN để hoàn tác commit ĐÃ PUSH/ĐÃ CHIA SẺ

```powershell
git revert <commit-hash>
```

Thay vì XÓA commit cũ (như `reset`), `revert` tạo 1 commit **MỚI** có nội dung NGƯỢC LẠI hoàn toàn với commit cần hoàn tác — lịch sử VẪN GIỮ NGUYÊN commit gốc (không xóa), chỉ thêm 1 commit mới "phủ định" nó. Đây là lý do `revert` AN TOÀN để dùng trên commit đã push/share — không viết lại lịch sử, không cần force push, không ảnh hưởng tới người khác đã pull.

```powershell
git revert HEAD              # revert commit gần nhất
git revert <hash> --no-edit    # revert mà không mở editor sửa message (dùng message mặc định)
git revert <hash1> <hash2>       # revert nhiều commit (theo thứ tự NGƯỢC, từ mới nhất trước)
```

**So sánh trực quan:**

```
reset --hard:  A---B---C   →   A---B          (C biến mất khỏi lịch sử)
revert:        A---B---C   →   A---B---C---D   (D = "anti-C", lịch sử VẪN CÒN C, chỉ thêm D phủ định nó)
```

## 6. `git checkout -- <file>` — cách CŨ, tương đương `restore` (biết để đọc hiểu code/tutorial cũ)

```powershell
git checkout -- file.txt    # tương đương git restore file.txt (cách cũ, vẫn hoạt động)
```

Như đã nói ở [Bài 4 mục 2](./4_branching.md), `checkout` làm quá nhiều việc khác nhau — `restore`/`switch` là cách hiện đại, rõ ràng hơn, nên ưu tiên dùng cho code MỚI.

## Ví dụ đầy đủ

```powershell
git init undo-demo; cd undo-demo
echo "v1" > f.txt; git add .; git commit -m "commit 1"
echo "v2" > f.txt; git add .; git commit -m "commit 2"
echo "v3" > f.txt; git add .; git commit -m "commit 3"

# Thử soft reset — undo commit 3, giữ thay đổi ở staging
git reset --soft HEAD~1
git status    # thấy f.txt đã "Changes to be committed" (staged)
git commit -m "commit 3 (viết lại message)"

# Thử revert — cách AN TOÀN cho commit đã "push" (giả lập)
git revert HEAD --no-edit
cat f.txt      # nội dung quay về "v2" (đã push commit "v3 viết lại" bị phủ định)
git log --oneline    # vẫn thấy ĐẦY ĐỦ lịch sử, kể cả commit đã bị revert
```

## Bài tập

1. **3 loại reset**: tạo 3 commit, thử LẦN LƯỢT `--soft`, `--mixed`, `--hard` (mỗi lần reset xong, dùng `git reset --hard <hash-gốc>` để quay lại trạng thái ban đầu trước khi thử loại tiếp theo), quan sát khác biệt bằng `git status`.
2. **`restore` vs `reset`**: sửa 1 file, thử `git restore` để bỏ thay đổi — so sánh với việc dùng `git reset --hard HEAD` cho cùng mục đích, giải thích tại sao `restore` phù hợp hơn khi chỉ cần xử lý 1 file.
3. **`revert` cho commit "đã push"**: giả lập 1 commit "đã push" (chỉ cần tưởng tượng), dùng `revert` để hoàn tác, verify lịch sử vẫn còn nguyên commit gốc.
4. **Tình huống thực tế**: bạn vừa `git reset --hard` nhầm, xóa mất 1 commit quan trọng CHƯA push. Trước khi đọc [Bài 12](./12_rewriting_history.md), thử đoán xem có cách nào cứu lại không (gợi ý: Git hiếm khi XÓA THẬT SỰ dữ liệu ngay lập tức).

## Tiếp theo
→ [Bài 11: Tags & Releases](./11_tags.md)
