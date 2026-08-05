# Bài 12: Sửa Lại Lịch Sử — rebase -i, amend, reflog

## Mục tiêu
- Dọn dẹp lịch sử commit lộn xộn TRƯỚC KHI chia sẻ, bằng `rebase -i` và `amend`.
- Nắm `reflog` — "lưới an toàn cuối cùng" của Git, gần như luôn cứu được khi "lỡ tay".

## 1. `git commit --amend` — sửa lại commit GẦN NHẤT

```powershell
git commit --amend -m "Message mới, chính xác hơn"    # sửa MESSAGE của commit gần nhất

# Sửa cả NỘI DUNG commit gần nhất (quên add 1 file)
git add file-quen-add.txt
git commit --amend --no-edit    # giữ nguyên message cũ, chỉ thêm file vào commit gần nhất
```

**Bản chất:** `amend` không "sửa" commit cũ tại chỗ — nó tạo 1 commit HOÀN TOÀN MỚI (nội dung + message mới), rồi di chuyển con trỏ branch trỏ tới commit mới đó thay vì commit cũ (commit cũ vẫn tồn tại trong object database, chỉ không còn branch nào trỏ tới — sẽ "biến mất" dần khi Git dọn rác, nhưng vẫn có thể tìm lại qua `reflog` mục 4 trong lúc chưa bị dọn). Vì tạo commit MỚI, **quy tắc vàng của rebase** ([Bài 7 mục 3](./7_merge_vs_rebase.md)) áp dụng y hệt: chỉ `amend` commit CHƯA push.

## 2. `git rebase -i` — Interactive Rebase, công cụ dọn dẹp lịch sử mạnh nhất

```powershell
git rebase -i HEAD~4    # mở editor cho 4 commit gần nhất
```

Editor hiện danh sách commit (từ CŨ tới MỚI, ngược với `git log`), mỗi dòng có thể đổi hành động:

```
pick a1b2c3d Add: tính năng đăng nhập
pick e4f5g6h Fix: sửa lỗi chính tả
pick h7i8j9k Fix: sửa lỗi chính tả (lần 2, thật ra vẫn còn lỗi)
pick k1l2m3n Add: viết test cho đăng nhập
```

**Các hành động khả dụng:**

| Lệnh | Ý nghĩa |
|---|---|
| `pick` | giữ nguyên commit |
| `reword` | giữ nội dung, cho phép SỬA message |
| `edit` | dừng lại TẠI commit này để sửa nội dung (dùng `git commit --amend`, rồi `git rebase --continue`) |
| `squash` | GỘP commit này vào commit NGAY TRÊN nó, gộp cả message (cho chỉnh sửa) |
| `fixup` | giống squash, nhưng BỎ LUÔN message của commit này (dùng message của commit trên) |
| `drop` | XÓA HẲN commit này |

**Ví dụ dọn dẹp** — gộp 2 commit "Fix: sửa lỗi chính tả" thành 1, vì commit thứ 2 thực chất chỉ là "sửa cho đúng cái đã sửa dở lần 1":

```
pick a1b2c3d Add: tính năng đăng nhập
pick e4f5g6h Fix: sửa lỗi chính tả
fixup h7i8j9k Fix: sửa lỗi chính tả (lần 2, thật ra vẫn còn lỗi)
pick k1l2m3n Add: viết test cho đăng nhập
```

Lưu và đóng editor — Git tự động thực hiện, kết quả còn lại 3 commit thay vì 4, lịch sử SẠCH hơn nhiều để review sau này.

## 3. Sắp xếp LẠI thứ tự commit bằng rebase -i

Chỉ cần đổi THỨ TỰ các dòng trong editor — Git sẽ áp dụng theo đúng thứ tự MỚI. **Cẩn thận:** nếu các commit phụ thuộc lẫn nhau (commit sau sửa file mà commit trước tạo ra), đổi thứ tự có thể gây conflict khi rebase — Git sẽ dừng lại đúng như conflict thường ([Bài 8](./8_conflicts.md)), xử lý bằng `git rebase --continue` sau khi giải quyết.

## 4. `git reflog` — "lưới an toàn" ghi lại MỌI nơi HEAD từng trỏ tới

Đây là công cụ QUAN TRỌNG NHẤT của bài, giải quyết nỗi sợ lớn nhất khi mới học Git: "lỡ làm mất commit thì sao?"

```powershell
git reflog    # liệt kê MỌI thao tác đã di chuyển HEAD: commit, checkout, reset, rebase, merge...
```

Output mẫu:
```
a1b2c3d HEAD@{0}: reset: moving to HEAD~1
e4f5g6h HEAD@{1}: commit: Fix bug quan trọng
h7i8j9k HEAD@{2}: checkout: moving from feature to main
```

**Ý nghĩa:** dù bạn vừa `git reset --hard` "xóa mất" commit `e4f5g6h`, commit đó **CHƯA hề bị xóa khỏi object database** ([Bài 17](./17_git_internals.md)) — nó chỉ không còn branch nào trỏ tới. `reflog` vẫn ghi nhớ nó tồn tại, cho phép khôi phục:

```powershell
git reset --hard HEAD@{1}    # hoặc: git reset --hard e4f5g6h
# khôi phục CHÍNH XÁC về trạng thái tại thời điểm đó
```

**Giới hạn của reflog:** chỉ lưu trên MÁY LOCAL (không đồng bộ qua remote — reflog của bạn khác reflog người khác), và có thời hạn mặc định (~90 ngày cho commit "unreachable", ~30 ngày cho reflog entry khác) trước khi Git dọn rác thật sự (`git gc`). Trong hầu hết tình huống thực tế ("lỡ tay" gần đây), reflog gần như LUÔN cứu được.

## 5. `git push --force-with-lease` — force push AN TOÀN hơn `--force`

Sau khi rebase/amend commit ĐÃ PUSH (chỉ nên làm trên branch riêng của bạn — [Bài 7 mục 3](./7_merge_vs_rebase.md)), cần force push để cập nhật remote:

```powershell
git push --force origin feature-branch              # NGUY HIỂM — ghi đè remote KỂ CẢ nếu có commit mới của người khác mà bạn chưa biết
git push --force-with-lease origin feature-branch      # AN TOÀN HƠN — Git kiểm tra remote CHƯA bị ai khác thay đổi kể từ lần bạn fetch gần nhất, TỪ CHỐI nếu có
```

**Luôn ưu tiên `--force-with-lease`** — nó bảo vệ bạn khỏi vô tình ghi đè lên commit của đồng nghiệp mà bạn chưa kịp thấy (vd họ vừa push thêm gì đó lên CHÍNH branch bạn đang force push, trong lúc bạn đang rebase).

## Ví dụ đầy đủ

```powershell
git init rewrite-demo; cd rewrite-demo
echo "1" > f.txt; git add .; git commit -m "commit 1"
echo "2" > f.txt; git add .; git commit -m "typo fix"
echo "3" > f.txt; git add .; git commit -m "typo fix lai"
echo "4" > f.txt; git add .; git commit -m "Add tinh nang moi"

git rebase -i HEAD~4
# đổi dòng 2 "typo fix lai" thành "fixup", lưu lại
git log --oneline    # giờ chỉ còn 3 commit

# Giả lập "lỡ tay" reset --hard
git reset --hard HEAD~2
git log --oneline    # commit gần nhất "biến mất"

# CỨU LẠI bằng reflog
git reflog
git reset --hard HEAD@{1}    # khôi phục về TRƯỚC lúc reset --hard
git log --oneline    # commit đã quay lại đầy đủ
```

## Bài tập

1. **`amend`**: commit thiếu 1 file, dùng `--amend` để bổ sung mà không tạo commit riêng.
2. **`rebase -i` với squash/fixup**: tạo 4-5 commit "vụn vặt" (như ví dụ), dọn thành 2-3 commit có ý nghĩa rõ ràng.
3. **Cứu commit đã mất bằng reflog**: cố tình `git reset --hard` lùi lại vài commit, dùng `reflog` để khôi phục về đúng trạng thái trước đó — làm ĐÚNG theo "Ví dụ đầy đủ" để tự tay trải nghiệm cảm giác "cứu được".
4. **`--force-with-lease`**: mô phỏng tình huống 2 người cùng làm trên 1 branch remote — thử force push bằng cả `--force` và `--force-with-lease` khi remote đã có thay đổi mới, so sánh hành vi khác nhau.

## Tổng kết Giai đoạn 2
Bạn đã nắm merge/rebase, giải quyết conflict, stash/cherry-pick, và các cách hoàn tác/sửa lịch sử — đặc biệt `reflog`, "lưới an toàn" giúp bạn tự tin thử nghiệm Git mà không sợ mất dữ liệu. Giai đoạn 3 sẽ đi vào quy trình làm việc nhóm và các công cụ nâng cao.

## Tiếp theo
→ [Bài 13: Git Workflow cho team](./13_workflows.md)
