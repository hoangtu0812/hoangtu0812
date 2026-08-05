# Bài 2: Khái niệm cốt lõi — 3 vùng làm việc

## Mục tiêu
- Hiểu Working Directory, Staging Area, Repository — mô hình tinh thần quan trọng nhất để dùng Git thành thạo.
- Hiểu Git lưu **snapshot** (ảnh chụp toàn bộ), không phải diff — khác biệt cốt lõi với nhiều người hình dung sai.

## 1. 3 vùng (3 trees) của Git

```
Working Directory  --git add-->  Staging Area (Index)  --git commit-->  Repository (.git)
   (file bạn                        (bản nháp                            (lịch sử đã
    đang sửa)                    chuẩn bị commit)                          lưu vĩnh viễn)
```

- **Working Directory**: thư mục thật trên máy bạn, chứa file bạn đang sửa trực tiếp.
- **Staging Area (Index)**: "phòng chờ" — nơi bạn CHỌN chính xác thay đổi nào sẽ đi vào commit tiếp theo. Đây là khái niệm KHÔNG có trong nhiều VCS khác — cho phép commit CHỌN LỌC 1 phần thay đổi, không bắt buộc commit toàn bộ file đã sửa cùng lúc.
- **Repository**: nơi lưu lịch sử vĩnh viễn dưới dạng các **commit** — mỗi commit là 1 snapshot đầy đủ tại thời điểm đó.

## 2. Vì sao cần Staging Area — ví dụ thực tế

Giả sử bạn sửa 2 file: `login.py` (fix bug thật sự cần commit ngay) và `notes.txt` (ghi chú cá nhân, chưa muốn commit). Staging Area cho phép:

```powershell
git add login.py        # chỉ đưa login.py vào staging
git commit -m "Fix bug đăng nhập"   # commit CHỈ chứa login.py, notes.txt vẫn ở Working Directory
```

Thậm chí có thể stage CHỈ 1 PHẦN thay đổi trong CÙNG 1 file (`git add -p`) — hữu ích khi bạn sửa nhiều thứ không liên quan trong cùng file nhưng muốn tách thành các commit rõ ràng, có ý nghĩa riêng biệt.

## 3. Git lưu Snapshot, không lưu Diff — khác biệt quan trọng

Nhiều người lầm tưởng Git lưu "diff" (chỉ phần thay đổi) giữa các version, giống 1 số VCS cũ. Thực tế: **mỗi commit là 1 bản chụp TOÀN BỘ trạng thái file tại thời điểm đó** (không phải chỉ phần đã đổi). Để tiết kiệm dung lượng, nếu file KHÔNG đổi giữa 2 commit, Git không lưu lại bản sao mới — chỉ trỏ tới đúng file đã lưu ở commit trước (nhờ cơ chế content-addressable storage dùng SHA-1 hash — chi tiết đầy đủ ở [Bài 17: Git Internals](./17_git_internals.md)).

**Ý nghĩa thực tế:** `git checkout <commit>` không phải "áp dụng ngược 1 loạt diff" — nó chỉ đơn giản là "khôi phục đúng bản snapshot đã lưu", nên cực kỳ nhanh dù lịch sử dài hàng nghìn commit.

## 4. Trạng thái của file trong Working Directory

```
Untracked  →  (git add)  →  Staged  →  (git commit)  →  Committed (unmodified)
                                                              │
                                                     (sửa file) ↓
                                                          Modified
                                                              │
                                                       (git add) ↓
                                                          Staged (lại)
```

- **Untracked**: file mới, Git chưa "biết" tới nó.
- **Staged**: đã `git add`, sẵn sàng cho commit tiếp theo.
- **Modified**: file đã từng commit, giờ bị sửa thêm nhưng CHƯA `git add` lại.
- **Committed/Unmodified**: khớp hoàn toàn với commit gần nhất.

## 5. HEAD — con trỏ "bạn đang ở đâu"

`HEAD` là con trỏ trỏ tới commit hiện tại bạn đang đứng (thường là commit mới nhất của branch đang checkout). Khi commit mới, `HEAD` (và branch nó trỏ theo) tự động di chuyển tới commit mới đó. Hiểu rõ `HEAD` là chìa khóa để hiểu các lệnh phức tạp hơn sau này (`reset`, `checkout`, `rebase` — [Bài 10](./10_undo_changes.md), [Bài 12](./12_rewriting_history.md)).

## Minh họa bằng lệnh thực hành

```powershell
mkdir demo-git; cd demo-git
git init

echo "Hello" > file.txt
git status              # file.txt hiện là "Untracked"

git add file.txt
git status              # giờ hiện "Changes to be committed" — đã ở Staging Area

git commit -m "Thêm file.txt"
git status              # "nothing to commit, working tree clean" — Working Directory khớp Repository

echo "Hello World" >> file.txt
git status              # "Changes not staged for commit" — Modified, chưa staged
git diff                # xem CHÍNH XÁC phần đã thay đổi (chưa staged)

git add file.txt
git diff                # không còn gì hiện — vì diff mặc định so Working Directory với Staging
git diff --staged        # xem phần đã staged, so với commit gần nhất
```

## Bài tập

1. **Theo dõi trạng thái file**: thực hiện đúng chuỗi lệnh ở mục minh họa, sau MỖI lệnh chạy `git status`, ghi chú lại trạng thái file thay đổi ra sao.
2. **Stage 1 phần thay đổi**: sửa 1 file ở 2 chỗ khác nhau, dùng `git add -p` để chỉ stage 1 trong 2 thay đổi đó, verify bằng `git diff --staged`.
3. **Vẽ sơ đồ**: tự vẽ lại sơ đồ 3 vùng (mục 1) bằng lời/hình ảnh của riêng bạn, không nhìn lại bài — kiểm tra đã nội tâm hóa mô hình này chưa.
4. **`git diff` vs `git diff --staged`**: giải thích bằng lời sự khác biệt giữa 2 lệnh này, dựa trên mô hình 3 vùng.

## Tiếp theo
→ [Bài 3: Lệnh cơ bản — init, add, commit, log](./3_basic_commands.md)
