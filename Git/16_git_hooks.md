# Bài 16: Git Hooks

## Mục tiêu
- Hiểu Git Hooks là gì, tự động hóa kiểm tra chất lượng TRƯỚC KHI commit/push.
- Viết 1 hook thực tế bằng script shell, và biết công cụ quản lý hook hiện đại (Husky/pre-commit).

## 1. Git Hook là gì?

Hook là **script tự động chạy** khi 1 sự kiện Git cụ thể xảy ra (trước/sau commit, trước push, sau merge...) — nằm trong `.git/hooks/`, mỗi hook là 1 file thực thi (không cần đuôi mở rộng) đặt đúng tên sự kiện.

```powershell
ls .git/hooks/    # Git tự tạo sẵn các file .sample — đổi tên (bỏ .sample) và chỉnh sửa để kích hoạt
```

## 2. Các hook phổ biến nhất

| Hook | Chạy khi nào | Dùng để |
|---|---|---|
| `pre-commit` | TRƯỚC khi tạo commit (sau khi gõ `git commit`, trước khi lưu) | Chạy linter, formatter, test nhanh — HỦY commit nếu script trả về lỗi (exit code khác 0) |
| `commit-msg` | Sau khi viết message, TRƯỚC khi commit hoàn tất | Validate format message (vd bắt buộc theo Conventional Commits — [Bài 18](./18_team_conventions.md)) |
| `pre-push` | TRƯỚC khi push lên remote | Chạy test đầy đủ hơn, chặn push nếu test fail |
| `post-merge` | SAU khi merge hoàn tất | Tự động cài lại dependency nếu `package.json`/`requirements.txt` thay đổi |

## 3. Viết `pre-commit` hook cơ bản — chặn commit có debug code sót lại

```bash
#!/bin/sh
# .git/hooks/pre-commit — PHẢI có quyền thực thi (chmod +x trên Linux/Mac; Windows Git Bash tự xử lý)

if git diff --cached | grep -E "console\.log|debugger" > /dev/null; then
  echo "❌ Phát hiện console.log hoặc debugger trong code sắp commit — vui lòng xóa trước khi commit."
  exit 1   # exit code khác 0 -> HỦY commit
fi

exit 0
```

```powershell
chmod +x .git/hooks/pre-commit    # trên Windows Git Bash, thường không bắt buộc nhưng nên thêm cho chắc
```

Thử commit code có `console.log` — Git sẽ HIỆN thông báo lỗi và TỪ CHỐI tạo commit, buộc bạn xóa dòng debug trước.

## 4. Viết `commit-msg` hook — validate format message

```bash
#!/bin/sh
# .git/hooks/commit-msg
# $1 = đường dẫn tới file TẠM chứa message vừa gõ

commit_msg=$(cat "$1")
pattern="^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?: .{1,50}"

if ! echo "$commit_msg" | grep -qE "$pattern"; then
  echo "❌ Commit message phải theo format: <type>: <mô tả ngắn>"
  echo "   Ví dụ: feat: thêm chức năng đăng nhập"
  exit 1
fi
```

## 5. Vấn đề của hook trong `.git/hooks/` — KHÔNG được commit/chia sẻ tự động

`.git/hooks/` nằm TRONG thư mục `.git/` — theo định nghĩa, nó **KHÔNG BAO GIỜ** được Git track/commit (giống mọi thứ khác trong `.git/`). Nghĩa là hook bạn viết chỉ tồn tại TRÊN MÁY BẠN — đồng nghiệp clone repo về sẽ KHÔNG có hook này, trừ khi họ tự copy tay. Đây là lý do cần công cụ quản lý hook chia sẻ được qua Git.

## 6. Husky (cho project Node.js) — quản lý hook CHIA SẺ được qua Git

```powershell
npm install --save-dev husky
npx husky init
echo "npm test" > .husky/pre-commit
```

Husky lưu hook trong thư mục `.husky/` (KHÔNG phải `.git/hooks/`) — thư mục này BÌNH THƯỜNG, được Git track/commit như mọi file khác, nên khi đồng nghiệp `npm install`, hook TỰ ĐỘNG được kích hoạt cho họ — giải quyết đúng vấn đề ở mục 5.

## 7. `pre-commit` framework (đa ngôn ngữ, phổ biến cho Python/đa dạng stack)

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
  - repo: https://github.com/psf/black
    rev: 23.12.0
    hooks:
      - id: black
```

```powershell
pip install pre-commit
pre-commit install    # kích hoạt hook cho repo hiện tại, đọc cấu hình từ .pre-commit-config.yaml
```

File `.pre-commit-config.yaml` được commit vào repo (giống cách Husky làm) — mọi người trong team chỉ cần `pre-commit install` 1 lần sau khi clone để kích hoạt cùng bộ hook.

## 8. Hook local vs kiểm tra trên CI — vì sao cần CẢ 2

Hook local ([Bài 16](./16_git_hooks.md) này) giúp phát hiện lỗi SỚM, ngay trên máy trước khi commit/push — tiết kiệm thời gian chờ CI. Nhưng hook local **có thể bị bỏ qua** (`git commit --no-verify` bỏ qua mọi hook, hoặc đơn giản là đồng nghiệp chưa cài hook) — vì vậy KHÔNG được coi hook local là lớp bảo vệ DUY NHẤT; luôn có thêm kiểm tra tương tự chạy trên CI ([Bài 20](./20_ci_cd_integration.md)) như "lưới an toàn cuối cùng" trước khi merge vào `main`.

## Bài tập

1. **Viết `pre-commit` chặn debug code**: làm theo mục 3, thử commit code có `console.log`/`fmt.Println` debug, verify bị chặn; xóa dòng đó, verify commit thành công.
2. **Viết `commit-msg` validate format**: làm theo mục 4, thử commit với message KHÔNG đúng format, verify bị từ chối; sửa đúng format, verify thành công.
3. **Thiết lập Husky (nếu project Node.js) hoặc `pre-commit` (Python)**: chọn 1 trong 2 phù hợp với stack bạn quen, thiết lập hook chạy test/lint tự động, verify hook được commit vào repo (không nằm trong `.git/`).
4. **Bypass hook có chủ đích**: thử `git commit --no-verify` để bỏ qua hook đã thiết lập — giải thích bằng lời khi nào việc này CÓ THỂ chấp nhận được (vd fix khẩn cấp) và tại sao vẫn cần kiểm tra lại trên CI sau đó.

## Tiếp theo
→ [Bài 17: Git Internals — Git hoạt động thế nào bên trong](./17_git_internals.md)
