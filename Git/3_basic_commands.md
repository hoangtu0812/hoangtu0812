# Bài 3: Lệnh Cơ Bản — init, add, commit, log

## Mục tiêu
- Thành thạo vòng lặp làm việc hàng ngày: `status → add → commit → log`.
- Viết commit message có ý nghĩa, hiểu `git diff`/`git log` với các tùy chọn hữu ích.

## 1. `git init` — khởi tạo repository

```powershell
git init my-project
cd my-project
# hoặc: cd thư-mục-đã-có, rồi git init
```

Tạo thư mục ẩn `.git/` — đây là TOÀN BỘ dữ liệu Git cần (object database, refs, config — chi tiết [Bài 17](./17_git_internals.md)). Xóa thư mục `.git/` = xóa sạch lịch sử Git, nhưng KHÔNG ảnh hưởng file thật trong Working Directory.

## 2. `git add` — đưa thay đổi vào Staging Area

```powershell
git add file.txt          # 1 file cụ thể
git add file1.txt file2.txt  # nhiều file
git add .                   # mọi thay đổi trong thư mục hiện tại (và con)
git add -A                   # mọi thay đổi trong TOÀN BỘ repo (kể cả ngoài thư mục hiện tại)
git add -p                    # chọn từng "hunk" (đoạn thay đổi) để stage — rất hữu ích khi cần tách commit rõ ràng
```

**Thận trọng với `git add .`**: dễ vô tình stage cả file không nên commit (secret, file build, file rác) — luôn `git status` kiểm tra TRƯỚC khi `git add .`, và có `.gitignore` phù hợp ([Bài 6](./6_gitignore.md)).

## 3. `git commit` — lưu snapshot vào lịch sử

```powershell
git commit -m "Thêm chức năng đăng nhập"
git commit -am "Fix bug + update docs"   # -a: tự động stage mọi file ĐÃ TRACKED bị sửa (KHÔNG stage file mới/untracked)
git commit --amend                         # sửa lại commit GẦN NHẤT (chi tiết Bài 12 — chỉ dùng khi commit CHƯA push)
```

### Viết commit message tốt

```
<Loại>: Tóm tắt ngắn gọn, dưới 50 ký tự, thì mệnh lệnh (Add, không phải Added)

Giải thích chi tiết hơn nếu cần: TẠI SAO thay đổi này cần thiết,
không phải MÔ TẢ code đã thay đổi gì (code đã tự nói lên điều đó).
Wrap dòng ở khoảng 72 ký tự.
```

Ví dụ tốt: `Fix: sửa lỗi timeout khi gọi API thanh toán` — ngắn gọn, rõ HÀNH ĐỘNG và MỤC ĐÍCH. Ví dụ tệ: `update` hoặc `fix bug` (không ai biết sau này commit này làm gì nếu chỉ đọc log). Chi tiết chuẩn hóa quy ước message ở [Bài 18](./18_team_conventions.md) (Conventional Commits).

## 4. `git status` & `git diff` — luôn kiểm tra TRƯỚC khi commit

```powershell
git status                 # tổng quan trạng thái
git status -s                # dạng ngắn gọn (short format), quen thuộc khi dùng nhiều

git diff                     # thay đổi CHƯA staged (Working Directory vs Staging Area)
git diff --staged             # thay đổi ĐÃ staged (Staging Area vs commit gần nhất)
git diff HEAD                  # tổng hợp CẢ 2 (Working Directory vs commit gần nhất)
git diff commit1 commit2        # so sánh 2 commit bất kỳ
```

## 5. `git log` — xem lịch sử

```powershell
git log                      # đầy đủ: hash, author, date, message
git log --oneline             # rút gọn mỗi commit 1 dòng — dùng THƯỜNG XUYÊN NHẤT
git log --oneline --graph --all   # hiện đồ thị branch trực quan — cực hữu ích khi có nhiều branch
git log -p                     # hiện luôn diff của từng commit
git log --author="Ben"           # lọc theo tác giả
git log --since="2 weeks ago"     # lọc theo thời gian
git log -- path/to/file.txt        # chỉ xem lịch sử của 1 file cụ thể
git log --stat                      # thống kê số dòng thêm/xóa mỗi commit
```

`git log --oneline --graph --all` nên trở thành phản xạ — cho cái nhìn nhanh về "bức tranh tổng thể" của repo, đặc biệt hữu ích trước khi merge/rebase ([Bài 7](./7_merge_vs_rebase.md)).

## 6. `git show` — xem chi tiết 1 commit cụ thể

```powershell
git show <commit-hash>       # xem đầy đủ thông tin + diff của 1 commit
git show HEAD                  # commit hiện tại
git show HEAD~2                 # commit 2 bước trước HEAD
```

## Ví dụ đầy đủ: workflow hàng ngày

```powershell
git status                    # 1. Xem trạng thái hiện tại
git diff                        # 2. Xem lại CHÍNH XÁC những gì đã sửa
git add feature.py               # 3. Stage phần muốn commit
git diff --staged                  # 4. Kiểm tra lại lần cuối trước khi commit
git commit -m "Add: tính năng lọc dữ liệu theo ngày"  # 5. Commit
git log --oneline -5                # 6. Xác nhận commit đã vào lịch sử
```

## Bài tập

1. **Vòng lặp cơ bản**: tạo repo mới, thực hiện ĐÚNG quy trình 6 bước ở "Ví dụ đầy đủ" cho 3 commit khác nhau.
2. **`git add -p`**: sửa 1 file ở 3 vị trí khác nhau (không liên quan nhau), dùng `git add -p` tách thành 2-3 commit riêng biệt, có message rõ ràng cho từng commit.
3. **Đọc `git log`**: thử mọi tùy chọn ở mục 5 trên 1 repo có sẵn lịch sử (vd clone 1 project open-source nhỏ), so sánh lượng thông tin hiển thị.
4. **Viết commit message chuẩn**: tự review lại 5 commit gần nhất bạn từng viết (trong bất kỳ project nào), đánh giá theo tiêu chí ở mục 3 — commit nào tốt, commit nào cần cải thiện, viết lại chúng cho tốt hơn.

## Tiếp theo
→ [Bài 4: Branching cơ bản](./4_branching.md)
