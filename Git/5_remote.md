# Bài 5: Làm Việc Với Remote

## Mục tiêu
- Hiểu remote repository là gì, cách repo local và remote đồng bộ với nhau.
- Thành thạo `clone`, `push`, `pull`, `fetch`, và khái niệm upstream tracking.

## 1. Remote là gì?

Remote là 1 phiên bản của repo được lưu ở nơi khác (thường là GitHub/GitLab, nhưng cũng có thể là 1 máy khác trong mạng nội bộ, hoặc 1 thư mục khác trên cùng máy). Vì Git là hệ thống **phân tán** ([Bài 1 mục 2](./1_get_started.md)), remote không phải "server chủ" theo nghĩa bắt buộc — nó chỉ là 1 bản sao khác mà bạn đồng bộ dữ liệu qua lại, dù trong thực tế thường có 1 remote được coi là "nguồn chân lý chung" của team (thường đặt tên `origin`).

## 2. `git clone` — tạo bản sao local từ remote

```powershell
git clone https://github.com/user/repo.git
git clone git@github.com:user/repo.git my-folder-name   # dùng SSH, đổi tên thư mục
```

`clone` tự động: tải toàn bộ lịch sử, tạo remote tên `origin` trỏ về URL gốc, checkout sẵn branch mặc định (thường `main`).

## 3. `git remote` — quản lý danh sách remote

```powershell
git remote -v                                  # xem remote hiện có (và URL)
git remote add origin git@github.com:user/repo.git   # thêm remote mới (khi bắt đầu từ git init, chưa có origin)
git remote remove origin                          # xóa remote
git remote rename origin upstream                   # đổi tên
```

## 4. `git push` — đẩy commit local lên remote

```powershell
git push origin main                 # đẩy branch main lên remote "origin"
git push -u origin main                # -u (--set-upstream): ghi nhớ liên kết, lần sau chỉ cần "git push"
git push                                 # sau khi đã -u, không cần chỉ định lại remote/branch

git push origin feature-login            # đẩy 1 branch mới lên remote lần đầu
```

**"Upstream tracking"** là liên kết giữa branch local và branch remote tương ứng (vd `main` local theo dõi `origin/main`) — sau khi thiết lập (`-u`), Git biết CHÍNH XÁC push/pull đi đâu mà không cần gõ lại tên remote/branch mỗi lần.

## 5. `git fetch` vs `git pull` — khác biệt QUAN TRỌNG cần phân biệt rõ

```powershell
git fetch origin      # CHỈ tải về dữ liệu mới từ remote, KHÔNG tự động merge vào branch hiện tại
git pull origin main    # = git fetch + git merge (hoặc rebase nếu cấu hình) NGAY LẬP TỨC
```

`git fetch` AN TOÀN hơn — tải dữ liệu mới về nhưng để bạn TỰ QUYẾT ĐỊNH khi nào và cách nào (merge/rebase) tích hợp nó vào công việc hiện tại, có thời gian xem trước (`git log origin/main`) những gì đã thay đổi. `git pull` tiện lợi hơn (1 lệnh) nhưng "vội vàng" hơn — với người mới, nên ưu tiên thói quen `fetch` rồi xem xét, thay vì `pull` ngay theo phản xạ, để tránh merge/rebase bất ngờ khi chưa hiểu rõ thay đổi từ remote.

```powershell
# Workflow AN TOÀN hơn thay vì pull ngay
git fetch origin
git log HEAD..origin/main --oneline    # xem TRƯỚC những commit mới trên remote mà mình chưa có
git merge origin/main                    # rồi mới quyết định merge (hoặc rebase)
```

## 6. Cấu hình pull mặc định — merge hay rebase?

```powershell
git config --global pull.rebase false   # mặc định của Git: pull = fetch + merge
git config --global pull.rebase true      # pull = fetch + rebase (lịch sử thẳng hàng hơn — Bài 7)
```

Nên thống nhất lựa chọn này TRONG TOÀN TEAM để lịch sử nhất quán — tùy chọn tùy theo Git Workflow team chọn ([Bài 13](./13_workflows.md)).

## 7. Xử lý khi push bị từ chối (remote có commit mới mà local chưa có)

```
! [rejected] main -> main (fetch first)
error: failed to push some refs
```

Nghĩa là remote đã có commit mới hơn (do người khác push trước) mà local chưa đồng bộ. Giải pháp:

```powershell
git pull origin main       # đồng bộ trước (merge hoặc rebase tùy cấu hình)
# giải quyết conflict nếu có (Bài 8)
git push origin main         # rồi mới push lại
```

**KHÔNG BAO GIỜ** dùng `git push --force` để "ép" đẩy đè lên remote trừ khi thực sự hiểu hậu quả (xóa mất commit của người khác) — chi tiết về force push an toàn hơn (`--force-with-lease`) ở [Bài 12](./12_rewriting_history.md).

## Ví dụ đầy đủ: từ `git init` tới remote GitHub

```powershell
# Trên GitHub: tạo repo rỗng tên "my-project" (KHÔNG khởi tạo README để tránh conflict ngay từ đầu)

git init my-project; cd my-project
echo "# My Project" > README.md
git add README.md; git commit -m "Initial commit"

git remote add origin git@github.com:username/my-project.git
git push -u origin main

# Từ máy khác (hoặc thư mục khác), lấy code về:
git clone git@github.com:username/my-project.git
```

## Bài tập

1. **Tạo repo & đẩy lên GitHub**: làm theo đúng "Ví dụ đầy đủ", tạo 1 repo thật trên GitHub.
2. **Clone về vị trí khác**: clone repo vừa tạo vào 1 thư mục khác (mô phỏng "máy khác"), thử sửa file, commit, push từ đó.
3. **`fetch` rồi mới `merge`**: từ bản clone thứ 2, tạo 1 commit mới, push. Quay lại bản clone gốc, thực hành đúng quy trình `git fetch` + `git log HEAD..origin/main` + `git merge` thay vì `pull` ngay.
4. **Giả lập & xử lý push bị từ chối**: từ 2 bản clone, cả 2 cùng sửa 1 file KHÁC NHAU rồi cùng cố push — quan sát lỗi rejected ở bản push sau, xử lý đúng quy trình mục 7.

## Tiếp theo
→ [Bài 6: .gitignore & Quản lý file](./6_gitignore.md)
