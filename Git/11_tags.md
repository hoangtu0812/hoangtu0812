# Bài 11: Tags & Releases

## Mục tiêu
- Hiểu tag khác branch thế nào, khi nào dùng tag.
- Phân biệt lightweight tag vs annotated tag.
- Liên hệ Semantic Versioning.

## 1. Tag là gì, khác branch thế nào?

Tag là 1 con trỏ trỏ tới **1 commit CỤ THỂ**, giống branch — nhưng khác biệt cốt lõi: **tag KHÔNG di chuyển**. Branch tự động "tiến" theo mỗi commit mới ([Bài 4 mục 1](./4_branching.md)), còn tag gắn CỐ ĐỊNH vào đúng 1 điểm trong lịch sử, dùng để đánh dấu các mốc quan trọng — điển hình nhất là **phiên bản phát hành (release)**.

## 2. Lightweight Tag vs Annotated Tag

```powershell
# Lightweight — chỉ là 1 con trỏ đơn giản, không kèm metadata
git tag v1.0.0-lw

# Annotated — có kèm message, tác giả, ngày tháng, có thể ký (sign) — KHUYẾN KHÍCH cho release thật
git tag -a v1.0.0 -m "Phiên bản đầu tiên phát hành production"
```

Annotated tag là 1 **object riêng** trong Git database (chi tiết [Bài 17](./17_git_internals.md)), lưu đầy đủ thông tin ai tạo tag, khi nào, message gì — lightweight tag chỉ là 1 con trỏ trần, không lưu metadata gì thêm. **Với release chính thức, luôn dùng annotated tag** (`-a`) để có đầy đủ thông tin truy vết sau này.

## 3. Xem, xóa tag

```powershell
git tag                       # liệt kê mọi tag
git tag -l "v1.*"                # lọc theo pattern
git show v1.0.0                    # xem chi tiết tag (message + commit nó trỏ tới)

git tag -d v1.0.0-lw                # xóa tag LOCAL
git push origin --delete v1.0.0       # xóa tag trên REMOTE (xóa local không tự động xóa remote)
```

## 4. Tag cho 1 commit TRONG QUÁ KHỨ (không phải commit hiện tại)

```powershell
git log --oneline               # tìm hash commit muốn đánh tag
git tag -a v0.9.0 <commit-hash> -m "Phiên bản beta"
```

## 5. Push tag lên remote

```powershell
git push origin v1.0.0             # push 1 tag cụ thể
git push origin --tags               # push MỌI tag local chưa có trên remote (cẩn thận nếu có tag thử nghiệm)
```

**Lưu ý:** `git push` (không chỉ định gì thêm) **KHÔNG tự động push tag** — phải push tường minh, tránh vô tình đẩy tag thử nghiệm lên remote chung.

## 6. Semantic Versioning (SemVer) — quy ước đặt tên version chuẩn

```
v<MAJOR>.<MINOR>.<PATCH>
```

- **MAJOR** tăng khi có breaking change (thay đổi không tương thích ngược).
- **MINOR** tăng khi thêm tính năng mới nhưng vẫn tương thích ngược.
- **PATCH** tăng khi chỉ fix bug, không thêm tính năng.

Ví dụ: `v2.3.1` — major version 2, đã có 3 lần thêm tính năng kể từ v2.0.0, và 1 lần fix bug kể từ v2.3.0. Tuân theo SemVer giúp người dùng thư viện của bạn (hoặc chính team bạn) biết NGAY từ số phiên bản: update lên bản mới có an toàn không (chỉ PATCH/MINOR thường an toàn, MAJOR cần đọc changelog cẩn thận).

## 7. Checkout tại 1 tag — xem code TẠI thời điểm release đó

```powershell
git checkout v1.0.0    # HEAD chuyển tới ĐÚNG commit mà tag trỏ tới — ở trạng thái "detached HEAD"
```

**"Detached HEAD"** nghĩa là `HEAD` đang trỏ TRỰC TIẾP vào 1 commit, KHÔNG thông qua branch nào — nếu bạn commit thêm ở trạng thái này, commit mới đó KHÔNG thuộc branch nào cả, dễ bị "mất dấu" (chỉ còn tìm lại qua `reflog` — [Bài 12](./12_rewriting_history.md)) nếu chuyển sang branch khác mà quên tạo branch mới trước. Nếu muốn SỬA code tại điểm tag đó, nên tạo branch mới ngay:

```powershell
git checkout -b hotfix-v1.0.1 v1.0.0    # tạo branch MỚI từ đúng điểm tag, an toàn để commit tiếp
```

## 8. GitHub Releases — lớp bọc thêm tính năng trên tag

Trên GitHub, "Releases" (Repo → Releases → Draft a new release) gắn với 1 tag, cho phép thêm: changelog dạng văn bản đẹp, đính kèm file build sẵn (binary, installer), tự động sinh "Release Notes" từ Pull Request đã merge kể từ release trước. Về mặt Git thuần, Release chỉ là tag + metadata bổ sung do GitHub quản lý.

## Ví dụ đầy đủ

```powershell
git init tags-demo; cd tags-demo
echo "v1" > app.txt; git add .; git commit -m "Initial release"
git tag -a v1.0.0 -m "Bản phát hành đầu tiên"

echo "v1.1" > app.txt; git add .; git commit -m "Thêm tính năng X"
git tag -a v1.1.0 -m "Thêm tính năng X, tương thích ngược"

echo "v1.1.1" > app.txt; git add .; git commit -m "Fix bug Y"
git tag -a v1.1.1 -m "Fix bug Y"

git tag    # liệt kê: v1.0.0, v1.1.0, v1.1.1
git log --oneline --decorate    # thấy tag hiển thị kèm commit tương ứng

git checkout -b hotfix v1.0.0    # cần fix khẩn cấp trên bản v1.0.0 cũ, không phải bản mới nhất
```

## Bài tập

1. **Tạo annotated tag**: tạo 3 commit, gắn tag theo đúng Semantic Versioning cho từng mốc (vd feature mới → tăng MINOR, fix bug → tăng PATCH).
2. **Push & xóa tag**: push tag lên GitHub, sau đó thử xóa cả local và remote, verify bằng `git ls-remote --tags origin`.
3. **Checkout tại tag & tạo hotfix branch**: checkout về 1 tag cũ, verify trạng thái "detached HEAD" (Git sẽ cảnh báo), tạo branch mới từ đó để sửa an toàn.
4. **Semantic Versioning thực hành**: cho 1 danh sách thay đổi giả định (vd "thêm API mới", "đổi tên field trong response — breaking", "sửa lỗi chính tả"), xác định version tiếp theo nên là gì nếu bản hiện tại là `v2.4.1`.

## Tiếp theo
→ [Bài 12: Sửa lại lịch sử — rebase -i, amend, reflog](./12_rewriting_history.md)
