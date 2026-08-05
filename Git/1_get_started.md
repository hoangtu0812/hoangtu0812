# Bài 1: Giới thiệu Git & Cài đặt

## Mục tiêu
- Hiểu Git là gì, vì sao cần version control.
- Phân biệt Git (công cụ) và GitHub/GitLab/Bitbucket (dịch vụ hosting).
- Cài đặt và cấu hình Git.

## 1. Version Control là gì, và tại sao cần nó?

Version Control System (VCS) ghi lại lịch sử thay đổi của code theo thời gian — cho phép: quay lại phiên bản cũ khi có bug, biết AI đã sửa gì và tại sao, nhiều người cùng làm việc trên 1 codebase mà không ghi đè lẫn nhau. Không có VCS, làm việc nhóm thường dẫn tới việc gửi file qua lại (`project_final_v2_thật.zip`) — không thể theo dõi lịch sử, dễ mất code, không thể hợp nhất thay đổi của nhiều người an toàn.

## 2. Git khác gì các VCS cũ (SVN, CVS)?

Git là hệ thống **phân tán (distributed)** — mỗi máy clone repo có TOÀN BỘ lịch sử, không chỉ 1 bản snapshot hiện tại như VCS tập trung (centralized) kiểu SVN. Điều này nghĩa là:
- Có thể commit, xem lịch sử, tạo branch hoàn toàn **offline** — không cần kết nối server.
- Nếu server hỏng, bất kỳ máy nào đã clone cũng có bản sao đầy đủ để khôi phục.
- Branching trong Git cực nhẹ và nhanh (chi tiết ở [Bài 4](./4_branching.md), [Bài 17](./17_git_internals.md)) — khác SVN, nơi tạo branch tốn kém và chậm.

## 3. Git vs GitHub/GitLab/Bitbucket — phân biệt quan trọng

**Git** là phần mềm quản lý phiên bản, chạy hoàn toàn local trên máy bạn — không bắt buộc cần internet hay tài khoản nào.

**GitHub/GitLab/Bitbucket** là các **dịch vụ hosting** repository Git trên cloud, cộng thêm tính năng: giao diện web, Pull Request/Merge Request, Issue tracking, CI/CD (GitHub Actions/GitLab CI), quản lý quyền truy cập team. Bạn có thể dùng Git mà KHÔNG CẦN GitHub (chỉ làm việc local, hoặc tự host server Git riêng) — nhưng hầu hết công việc nhóm thực tế đều kết hợp cả 2.

## 4. Cài đặt Git

```powershell
# Windows: tải tại https://git-scm.com/download/win, hoặc dùng winget
winget install --id Git.Git -e --source winget

# Kiểm tra cài đặt
git --version
```

## 5. Cấu hình lần đầu — BẮT BUỘC trước khi commit

```powershell
git config --global user.name "Ten Cua Ban"
git config --global user.email "email@example.com"

# Xem lại cấu hình
git config --list

# Đặt editor mặc định (dùng khi Git cần mở editor, vd viết commit message dài, rebase -i)
git config --global core.editor "code --wait"   # VS Code
```

Thông tin `user.name`/`user.email` được **gắn vào MỌI commit** bạn tạo ra — đây là "chữ ký" xác định ai đã thực hiện thay đổi, hiển thị trong `git log` và trên GitHub.

## 6. SSH Key — xác thực với GitHub không cần nhập mật khẩu mỗi lần

```powershell
ssh-keygen -t ed25519 -C "email@example.com"
# Nhấn Enter để chấp nhận đường dẫn mặc định, có thể đặt passphrase hoặc để trống

# Copy public key để thêm vào GitHub (Settings -> SSH and GPG keys)
cat ~/.ssh/id_ed25519.pub
```

Sau khi thêm public key vào GitHub, test kết nối:

```powershell
ssh -T git@github.com
```

## 7. 3 lệnh sẽ dùng nhiều nhất ngay từ đầu (xem trước, chi tiết ở Bài 3)

```powershell
git status    # xem trạng thái hiện tại của repo
git add       # đưa thay đổi vào staging area
git commit     # lưu snapshot thay đổi vào lịch sử
```

## Bài tập

1. **Cài đặt & cấu hình**: cài Git, chạy `git config --global user.name`/`user.email`, verify bằng `git config --list`.
2. **Tạo SSH key**: tạo SSH key, thêm public key vào tài khoản GitHub, test bằng `ssh -T git@github.com`.
3. **Phân biệt Git vs GitHub**: viết 3-4 dòng giải thích bằng lời (không tra cứu) sự khác biệt Git và GitHub cho 1 người mới hoàn toàn — luyện khả năng diễn đạt lại kiến thức.
4. **Khám phá `git --help`**: chạy `git help` và `git help commit`, đọc qua cấu trúc tài liệu — đây là nguồn tra cứu chính thức nên quen thuộc ngay từ đầu.

## Tiếp theo
→ [Bài 2: Khái niệm cốt lõi — 3 vùng làm việc](./2_core_concepts.md)
