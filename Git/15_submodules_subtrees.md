# Bài 15: Submodule & Subtree

## Mục tiêu
- Hiểu vấn đề "repo lồng nhau" mà submodule/subtree giải quyết.
- So sánh 2 giải pháp, biết ưu nhược điểm để chọn đúng.

## 1. Vấn đề: khi nào cần "repo bên trong repo"?

Giả sử project của bạn dùng chung 1 thư viện nội bộ (vd 1 bộ UI component) mà TEAM KHÁC đang phát triển trong 1 repo RIÊNG — bạn cần "nhúng" repo đó vào project của mình, nhưng vẫn muốn nó giữ được lịch sử Git RIÊNG, có thể cập nhật độc lập khi team kia release bản mới. Copy-paste code thủ công mất khả năng theo dõi version/cập nhật dễ dàng — đây là lúc cần submodule hoặc subtree.

## 2. Git Submodule — repo con vẫn là repo ĐỘC LẬP, chỉ lưu "tham chiếu"

```powershell
git submodule add https://github.com/team/shared-ui.git libs/shared-ui
git commit -m "Add: submodule shared-ui"
```

**Cơ chế:** repo cha KHÔNG lưu code của `shared-ui` trực tiếp — nó chỉ lưu 1 file `.gitmodules` (chứa URL + đường dẫn) và 1 con trỏ tới **1 commit CỤ THỂ** của repo `shared-ui`. Thư mục `libs/shared-ui` về bản chất là 1 Git repo RIÊNG BIỆT, hoàn toàn tách khỏi lịch sử của repo cha.

```powershell
# Clone repo có submodule — PHẢI thêm cờ, nếu không thư mục submodule sẽ TRỐNG
git clone --recurse-submodules https://github.com/user/main-project.git

# Nếu đã clone rồi mới nhớ ra:
git submodule init
git submodule update

# Cập nhật submodule lên commit MỚI NHẤT của repo con
git submodule update --remote libs/shared-ui
git add libs/shared-ui
git commit -m "Update: shared-ui lên version mới"
```

**Ưu điểm:** lịch sử 2 repo hoàn toàn TÁCH BIỆT, rõ ràng ai chịu trách nhiệm gì; repo con có thể dùng chung cho NHIỀU repo cha khác nhau. **Nhược điểm:** trải nghiệm dùng khá rắc rối cho người mới (dễ quên `--recurse-submodules`, dễ quên `submodule update` sau khi pull, submodule "trỏ" tới đúng 1 commit cụ thể — nếu quên commit lại tham chiếu sau khi update, người khác vẫn thấy version CŨ).

## 3. Git Subtree — nhúng TOÀN BỘ code repo con TRỰC TIẾP vào lịch sử repo cha

```powershell
git subtree add --prefix=libs/shared-ui https://github.com/team/shared-ui.git main --squash
```

**Cơ chế:** khác hẳn submodule, subtree COPY toàn bộ code (và tùy chọn, cả lịch sử) của repo con vào ĐÚNG lịch sử của repo cha — sau lệnh này, `libs/shared-ui` là thư mục BÌNH THƯỜNG như mọi thư mục khác trong repo cha, KHÔNG phải 1 Git repo riêng. `--squash` gộp toàn bộ lịch sử repo con thành 1 commit duy nhất khi thêm vào (tùy chọn, tránh làm phình lịch sử repo cha).

```powershell
# Cập nhật code mới nhất từ repo con vào subtree
git subtree pull --prefix=libs/shared-ui https://github.com/team/shared-ui.git main --squash

# Đẩy thay đổi TỪ subtree NGƯỢC LẠI repo con gốc (nếu bạn có sửa trực tiếp trong thư mục đó)
git subtree push --prefix=libs/shared-ui https://github.com/team/shared-ui.git main
```

**Ưu điểm:** người clone repo cha KHÔNG CẦN làm gì đặc biệt (`git clone` bình thường là đủ, không cần cờ `--recurse-submodules`) — trải nghiệm đơn giản hơn nhiều cho người dùng cuối. **Nhược điểm:** lịch sử repo cha "phình to" hơn (chứa cả code của repo con), thao tác đồng bộ 2 chiều (`subtree push`) phức tạp hơn submodule.

## 4. So sánh nhanh — chọn cái nào?

| | Submodule | Subtree |
|---|---|---|
| Code lưu ở đâu | Repo riêng, chỉ lưu tham chiếu | Copy trực tiếp vào lịch sử repo cha |
| Trải nghiệm clone | Cần `--recurse-submodules` | `git clone` bình thường là đủ |
| Kích thước repo cha | Nhỏ (chỉ lưu con trỏ) | Lớn hơn (chứa cả code con) |
| Phù hợp | Thư viện dùng chung nhiều repo, muốn tách biệt rõ trách nhiệm | Muốn đơn giản hóa trải nghiệm người dùng cuối, ít cần đồng bộ ngược |
| Độ phổ biến thực tế | Phổ biến hơn cho thư viện lớn, nhiều team | Ít phổ biến hơn, nhưng đơn giản hơn cho use case nhỏ |

## 5. Khi nào KHÔNG cần cả 2 — cân nhắc package manager

Với hầu hết trường hợp thư viện code dùng lại (không phải cần đồng bộ SOURCE CODE trực tiếp), **package manager** (npm, pip, Go modules — [Go Bài 9](../Go/9_packages_modules.md), [Python Bài 9](../Python/9_modules_packages.md)) thường là lựa chọn TỐT HƠN submodule/subtree — publish thư viện thành package có version rõ ràng, quản lý dependency chuẩn hóa, không cần biết gì về cơ chế Git bên dưới. Chỉ nên cân nhắc submodule/subtree khi thực sự cần theo dõi/sửa SOURCE CODE của thư viện đó cùng lúc với code chính, không chỉ dùng nó như 1 dependency đóng gói sẵn.

## Bài tập

1. **Thêm submodule**: tạo 2 repo test (1 "chính", 1 "thư viện"), thêm thư viện làm submodule của repo chính, thử clone lại (thiếu cờ trước, rồi đúng cờ) để thấy khác biệt.
2. **Cập nhật submodule**: sửa code ở repo "thư viện", `submodule update --remote` ở repo chính, verify lịch sử ghi nhận đúng.
3. **Thêm subtree**: lặp lại tình huống tương tự nhưng dùng `git subtree add`, so sánh cấu trúc `.git` và trải nghiệm clone giữa 2 cách.
4. **Ra quyết định**: với project thực tế bạn đang làm (nếu có), đánh giá xem có tình huống nào phù hợp dùng submodule/subtree không, hay package manager đã đủ — viết ra lý do quyết định.

## Tiếp theo
→ [Bài 16: Git Hooks](./16_git_hooks.md)
