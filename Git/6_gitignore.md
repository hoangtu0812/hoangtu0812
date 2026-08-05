# Bài 6: .gitignore & Quản Lý File

## Mục tiêu
- Viết `.gitignore` đúng cách, hiểu pattern matching.
- Xử lý tình huống đã lỡ commit file không nên track (kể cả file nhạy cảm).

## 1. `.gitignore` là gì và tại sao cần

File `.gitignore` liệt kê pattern các file/thư mục Git nên **bỏ qua**, không hiện trong `git status`, không thể `git add` (trừ khi ép buộc). Cần thiết cho: file build tự sinh (`node_modules/`, `__pycache__/`, `*.exe`), file cấu hình cá nhân/máy (`.env`, `.vscode/settings.json`), file lớn/binary không nên nằm trong lịch sử Git (dataset, model file).

## 2. Cú pháp pattern cơ bản

```gitignore
# Comment bắt đầu bằng #

*.log                 # bỏ qua MỌI file có đuôi .log, ở BẤT KỲ thư mục nào
/build/                # chỉ bỏ qua thư mục "build" ở GỐC repo (dấu / đầu = chỉ định vị trí tuyệt đối từ root)
temp/                  # bỏ qua thư mục "temp" ở BẤT KỲ đâu trong repo
!important.log          # NGOẠI LỆ — dấu ! phủ định, GIỮ LẠI file này dù khớp pattern *.log ở trên
**/logs                  # bỏ qua thư mục "logs" ở BẤT KỲ độ sâu nào (** = khớp nhiều cấp thư mục)
config/*.local.yaml        # bỏ qua file khớp pattern CHỈ trong thư mục config/ (không đệ quy con)
```

**Thứ tự quan trọng:** pattern sau có thể ghi đè pattern trước (đặc biệt với `!`) — Git đọc từ trên xuống, áp dụng rule cuối cùng khớp.

## 3. Mẫu `.gitignore` thực tế theo loại project

```gitignore
# Python
venv/
__pycache__/
*.pyc
.env

# Node.js
node_modules/
dist/
.env

# Go
*.exe
*.test

# IDE
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db
```

GitHub cung cấp template `.gitignore` chuẩn cho hầu hết ngôn ngữ tại https://github.com/github/gitignore — nên bắt đầu từ đó thay vì viết từ đầu.

## 4. `.gitignore` CHỈ áp dụng cho file CHƯA từng track

Đây là hiểu lầm phổ biến nhất: nếu file đã LỠ được `git add`/`commit` TRƯỚC KHI thêm vào `.gitignore`, Git vẫn tiếp tục theo dõi nó — `.gitignore` chỉ ngăn file MỚI, CHƯA từng track.

```powershell
# Gỡ 1 file đã lỡ track ra khỏi Git (nhưng GIỮ LẠI file thật trên đĩa)
git rm --cached config.env
echo "config.env" >> .gitignore
git add .gitignore
git commit -m "Ngừng track config.env, thêm vào .gitignore"

# Gỡ CẢ 1 thư mục đã lỡ track
git rm -r --cached node_modules/
```

`--cached` là phần MẤU CHỐT — không có nó, `git rm` sẽ XÓA LUÔN file thật khỏi đĩa, không chỉ ngừng track.

## 5. Xử lý khi đã lỡ COMMIT file nhạy cảm (secret, API key)

Nếu file nhạy cảm đã push lên remote, chỉ `git rm --cached` KHÔNG ĐỦ — nó vẫn còn nguyên trong LỊCH SỬ (những commit cũ vẫn chứa nó, ai cũng có thể xem lại bằng `git log -p`). Cần "sửa lại lịch sử":

```powershell
# Cách nhanh cho repo cá nhân, chưa share nhiều: dùng git filter-repo (khuyến khích hơn filter-branch cũ)
pip install git-filter-repo
git filter-repo --path config.env --invert-paths   # xóa file này khỏi TOÀN BỘ lịch sử

# Sau đó BẮT BUỘC force push (vì lịch sử đã đổi hoàn toàn) và thông báo cho team
git push origin --force --all
```

**Quan trọng nhất:** nếu secret đã bị lộ (API key, password), việc xóa khỏi lịch sử Git là CHƯA ĐỦ — phải **thu hồi/đổi lại (rotate) secret đó ngay lập tức**, vì nó có thể đã bị crawl/cache lại ở đâu đó (fork của người khác, cache của GitHub) trước khi bạn kịp xóa.

## 6. `.gitignore` lồng nhau & global gitignore

```powershell
# .gitignore riêng cho máy cá nhân, KHÔNG commit vào repo (áp dụng cho MỌI repo trên máy)
git config --global core.excludesfile ~/.gitignore_global
```

Hữu ích cho file đặc thù máy bạn (vd `.DS_Store` trên Mac) mà không muốn buộc cả team phải thêm vào `.gitignore` chung của mỗi project.

## Bài tập

1. **Viết `.gitignore`**: tạo project Python hoặc Node.js giả lập, viết `.gitignore` phù hợp (dùng template từ GitHub làm tham khảo), verify bằng cách tạo file rác thuộc các pattern đó và `git status` không thấy chúng.
2. **Gỡ file đã lỡ track**: cố tình commit 1 file (vd `secret.txt`) TRƯỚC KHI thêm vào `.gitignore`, sau đó thực hành đúng quy trình mục 4 để gỡ nó ra.
3. **Global gitignore**: thiết lập global gitignore cho file đặc thù hệ điều hành của bạn.
4. **Nâng cao (chỉ thực hành trên repo test, KHÔNG làm trên repo thật)**: cài `git-filter-repo`, thử xóa 1 file khỏi TOÀN BỘ lịch sử của 1 repo test, verify bằng `git log --all -- <file>` không còn thấy nó.

## Tổng kết Giai đoạn 1
Bạn đã nắm vòng lặp làm việc cơ bản: 3 vùng làm việc, commit, branch, remote, và quản lý file qua `.gitignore`. Giai đoạn 2 sẽ đi sâu vào các tình huống thực tế phức tạp hơn: merge/rebase, conflict, và undo thay đổi.

## Tiếp theo
→ [Bài 7: Merge vs Rebase](./7_merge_vs_rebase.md)
