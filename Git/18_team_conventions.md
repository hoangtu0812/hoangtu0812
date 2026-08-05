# Bài 18: Quy Ước Làm Việc Nhóm

## Mục tiêu
- Áp dụng Conventional Commits để lịch sử commit dễ đọc, tự động sinh changelog được.
- Xây dựng checklist review team có thể dùng ngay.

## 1. Vì sao cần quy ước chung, không chỉ "mỗi người viết message theo ý mình"

Đã thấy ở [Bài 3 mục 3](./3_basic_commands.md): commit message tốt quan trọng thế nào. Với 1 người, "tốt" có thể tùy cảm tính — nhưng với TEAM, cần 1 **quy ước THỐNG NHẤT** để: mọi người đọc log hiểu ngay không cần đoán, công cụ có thể TỰ ĐỘNG xử lý (sinh changelog, xác định version tiếp theo theo SemVer — [Bài 11 mục 6](./11_tags.md)), tìm kiếm lịch sử theo loại thay đổi (`git log --grep="^fix"`).

## 2. Conventional Commits — chuẩn phổ biến nhất hiện nay

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Các `type` chuẩn:**

| Type | Ý nghĩa | Ảnh hưởng SemVer |
|---|---|---|
| `feat` | Thêm tính năng mới | MINOR |
| `fix` | Sửa bug | PATCH |
| `docs` | Chỉ thay đổi tài liệu | Không ảnh hưởng |
| `style` | Format code, không đổi logic (dấu cách, dấu chấm phẩy) | Không ảnh hưởng |
| `refactor` | Tái cấu trúc code, không thêm tính năng/sửa bug | Không ảnh hưởng |
| `test` | Thêm/sửa test | Không ảnh hưởng |
| `chore` | Việc vặt (update dependency, config build) | Không ảnh hưởng |
| `perf` | Cải thiện hiệu năng | PATCH |

**Ví dụ:**

```
feat(auth): thêm đăng nhập bằng Google OAuth

Cho phép user đăng nhập bằng tài khoản Google thay vì chỉ email/password.
Yêu cầu cấu hình GOOGLE_CLIENT_ID trong .env.

Closes #142
```

```
fix(payment): sửa lỗi timeout khi gọi API thanh toán

BREAKING CHANGE: đổi tên field "amount" thành "totalAmount" trong response
```

`BREAKING CHANGE:` trong footer đánh dấu thay đổi KHÔNG tương thích ngược — công cụ tự động (như `semantic-release`) sẽ hiểu đây là lý do tăng **MAJOR** version ([Bài 11 mục 6](./11_tags.md)), dù type là `fix`.

## 3. Lợi ích thực tế — tự động sinh Changelog

```powershell
npm install -g conventional-changelog-cli
conventional-changelog -p angular -i CHANGELOG.md -s
```

Công cụ đọc TOÀN BỘ lịch sử commit theo đúng format Conventional Commits, TỰ ĐỘNG nhóm lại thành changelog có cấu trúc:

```markdown
## v1.2.0
### Features
- **auth:** thêm đăng nhập bằng Google OAuth

### Bug Fixes
- **payment:** sửa lỗi timeout khi gọi API thanh toán
```

Không có quy ước chung, việc này phải làm THỦ CÔNG, dễ sót, tốn thời gian mỗi lần release.

## 4. Quy ước đặt tên branch (nhắc lại, mở rộng từ Bài 13)

```
feature/JIRA-123-them-dang-nhap-google
fix/JIRA-456-loi-timeout-thanh-toan
```

Nhiều team gắn mã ticket (Jira/Linear/GitHub Issue) vào tên branch — giúp liên kết TRỰC TIẾP giữa code thay đổi và yêu cầu/bug ban đầu, dễ truy vết "tại sao thay đổi này tồn tại" khi xem lại sau nhiều tháng.

## 5. Checklist Code Review chuẩn — dùng lại mỗi lần review

```markdown
## Chức năng
- [ ] Code giải quyết ĐÚNG vấn đề nêu trong PR description
- [ ] Đã xử lý các edge case hợp lý (input rỗng, null, giá trị biên)

## Chất lượng code
- [ ] Tên biến/hàm rõ ràng, không cần đoán
- [ ] Không có code trùng lặp có thể tách hàm dùng chung
- [ ] Không có debug code sót lại (console.log, print thừa)

## Test & An toàn
- [ ] Có test cho logic mới/đã sửa
- [ ] Không có thông tin nhạy cảm (API key, password) trong code/commit
- [ ] Không phá vỡ test cũ (CI pass — Bài 20)

## Tài liệu
- [ ] README/docs cập nhật nếu có thay đổi API công khai
```

## 6. Quy ước xử lý PR bị "stale" (để lâu không ai review/merge)

Team nên thống nhất: PR mở quá bao lâu (vd 3 ngày) không có phản hồi thì cần nhắc; PR có conflict với `main` thì tác giả chủ động rebase/merge cập nhật TRƯỚC khi nhờ review lại — tránh reviewer phải tự xử lý conflict thay tác giả.

## Ví dụ đầy đủ: dọn dẹp lịch sử commit lộn xộn theo Conventional Commits

```powershell
git log --oneline
# a1b2c3d fix stuff
# e4f5g6h wip
# h7i8j9k more changes
# k1l2m3n asdf

# Dùng interactive rebase (Bài 12) để reword lại từng commit
git rebase -i HEAD~4
# đổi "pick" thành "reword" cho mọi dòng, sửa lại từng message:
# feat(search): thêm bộ lọc tìm kiếm theo ngày
# fix(search): sửa lỗi filter không áp dụng khi đổi trang
# test(search): thêm test cho bộ lọc theo ngày
# docs(search): cập nhật README mô tả API filter mới
```

## Bài tập

1. **Viết commit theo Conventional Commits**: thực hiện 5 commit cho 1 tính năng giả định, mỗi commit dùng đúng `type` phù hợp.
2. **Sinh changelog tự động**: cài `conventional-changelog-cli`, chạy trên repo vừa tạo, xem changelog được sinh ra.
3. **Dọn dẹp lịch sử cũ**: lấy 1 repo có lịch sử commit "lộn xộn" (tự tạo hoặc dùng repo cũ của bạn), dùng `rebase -i` với `reword` để viết lại theo Conventional Commits.
4. **Áp dụng checklist review**: dùng checklist ở mục 5, tự review 1 PR/commit của chính bạn, đánh dấu mục nào đạt/chưa đạt.

## Tiếp theo
→ [Bài 19: Debug với Git — bisect, blame, worktree](./19_debugging_tools.md)
