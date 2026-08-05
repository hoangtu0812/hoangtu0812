# Bài 17: Git Internals — Git Hoạt Động Thế Nào Bên Trong

## Mục tiêu
- Hiểu Object Database (blob/tree/commit/tag) — nền tảng giải thích MỌI lệnh Git đã học trước đó.
- Tự tay khám phá object thật bằng `git cat-file`/`git hash-object`.

## 1. Vì sao học Internals — không phải chỉ để biết cho vui

Toàn bộ những "khó hiểu" khi học Git (tại sao commit có hash dài ngoằng, tại sao branch "nhẹ" mà SVN branch lại nặng, tại sao rebase tạo hash MỚI dù nội dung "giống hệt") đều có câu trả lời RÕ RÀNG khi hiểu cơ chế bên trong — phần này biến các lệnh Git từ "học thuộc" thành "suy luận được".

## 2. `.git/` — toàn bộ "bộ não" của repo nằm ở đây

```
.git/
├── HEAD              # con trỏ hiện tại đang trỏ tới branch nào (hoặc commit nào nếu detached)
├── config              # cấu hình riêng của repo này (user, remote...)
├── objects/              # OBJECT DATABASE — nơi lưu MỌI dữ liệu thật (blob, tree, commit, tag)
├── refs/
│   ├── heads/             # mỗi file = 1 branch, nội dung = hash commit nó trỏ tới
│   └── tags/                # mỗi file = 1 tag
└── logs/                      # dữ liệu reflog (Bài 12 mục 4)
```

## 3. 4 loại Object — nền tảng của MỌI thứ Git làm

**Blob (Binary Large Object)** — lưu NỘI DUNG THÔ của 1 file (không lưu tên file, không lưu quyền — chỉ nội dung). 2 file có nội dung GIỐNG HỆT nhau (dù tên khác nhau, ở thư mục khác nhau) chia sẻ CHUNG 1 blob — đây là cách Git tiết kiệm dung lượng.

**Tree** — tương đương 1 thư mục: danh sách các entry, mỗi entry trỏ tới 1 blob (file) hoặc 1 tree khác (thư mục con), kèm tên và quyền truy cập.

**Commit** — trỏ tới 1 tree (snapshot TOÀN BỘ thư mục gốc tại thời điểm đó), kèm: parent commit (1 với commit thường, 2 với merge commit, 0 với commit đầu tiên), tác giả, thời gian, message.

**Tag (annotated)** — trỏ tới 1 commit, kèm metadata (mục [Bài 11 mục 2](./11_tags.md)).

## 4. SHA-1 Hash — "địa chỉ" của mọi object, tính từ chính NỘI DUNG

Mỗi object được đặt tên bằng hash SHA-1 (40 ký tự hex) tính từ CHÍNH NỘI DUNG của nó — đây gọi là **content-addressable storage**. Hệ quả quan trọng:

- 2 object có nội dung GIỐNG HỆT nhau LUÔN có CÙNG hash — đây là lý do 2 file giống nhau chia sẻ chung 1 blob (mục 3).
- CHỈ CẦN đổi 1 bit nội dung, hash thay đổi HOÀN TOÀN — đây là lý do commit sau `rebase`/`amend` ([Bài 7](./7_merge_vs_rebase.md), [Bài 12](./12_rewriting_history.md)) có hash khác dù nội dung "gần giống" bản gốc: parent của nó đã đổi, mà commit object LƯU CẢ parent trong nội dung được hash, nên chỉ cần đổi parent, hash toàn bộ commit đổi theo.
- Có thể VERIFY tính toàn vẹn dữ liệu — nếu ai đó cố tình sửa nội dung 1 object mà không cập nhật lại mọi hash liên quan, Git phát hiện ngay lập tức vì hash không khớp.

## 5. Khám phá object THẬT bằng tay

```powershell
git init internals-demo; cd internals-demo
echo "Hello Git" > file.txt
git add file.txt

# Tìm hash của blob vừa tạo (nằm trong staging area)
git ls-files -s    # in ra mode, hash (SHA-1), tên file

# Xem NỘI DUNG THẬT của object (blob) qua hash đó
git cat-file -p <hash-vua-thay>    # in ra "Hello Git" — đúng là nội dung file

# Xem KIỂU của object
git cat-file -t <hash-vua-thay>    # in ra "blob"

git commit -m "commit dau tien"
git cat-file -p HEAD                 # xem NỘI DUNG object commit — thấy tree, author, message
git cat-file -p HEAD^{tree}            # xem tree mà commit đó trỏ tới — thấy danh sách file + hash blob
```

## 6. Tự tạo 1 blob bằng `hash-object` — hiểu content-addressable từ gốc

```powershell
echo "Noi dung test" | git hash-object --stdin           # chỉ TÍNH hash, KHÔNG lưu vào object database
echo "Noi dung test" | git hash-object --stdin -w          # tính hash VÀ lưu thật vào .git/objects/

git cat-file -p <hash-vua-in-ra>    # verify lấy lại đúng "Noi dung test"
```

Thử chạy lại CHÍNH XÁC 2 dòng trên 1 lần nữa — hash in ra sẽ **GIỐNG HỆT** lần trước, vì cùng nội dung → cùng hash (mục 4) — dù bạn "tạo" object 2 lần, Git chỉ lưu 1 bản duy nhất (ghi đè lên chính nó, không tốn thêm dung lượng).

## 7. Branch chỉ là 1 FILE TRẦN chứa 1 hash — giải thích vì sao branch "nhẹ"

```powershell
cat .git/refs/heads/main    # nội dung CHỈ LÀ 1 dòng: hash của commit mà branch "main" đang trỏ tới
```

Tạo branch mới = tạo 1 FILE MỚI trong `.git/refs/heads/` chứa hash — nhanh gần như tức thời (ghi 1 file text vài chục byte), KHÁC HẲN việc copy toàn bộ code (như branching kiểu SVN cũ) — đây là lý do "branching cực nhẹ" đã nhắc ở [Bài 4 mục 1](./4_branching.md) giờ có lời giải thích cụ thể.

```powershell
git branch new-branch
cat .git/refs/heads/new-branch    # thấy CHÍNH XÁC cùng hash với main (vì mới tách, chưa commit gì thêm)
```

## 8. `HEAD` — 1 con trỏ TỚI con trỏ khác (thường)

```powershell
cat .git/HEAD    # thường thấy: ref: refs/heads/main
```

`HEAD` thường KHÔNG trỏ trực tiếp tới 1 commit — nó trỏ tới 1 branch (refs/heads/main), rồi branch đó MỚI trỏ tới commit. Đây là lý do khi bạn commit, chỉ cần cập nhật 1 chỗ (file `refs/heads/main`) là branch VÀ HEAD đều "tự động" trỏ đúng commit mới — trong trạng thái "detached HEAD" ([Bài 11 mục 7](./11_tags.md)), `HEAD` trỏ TRỰC TIẾP vào hash commit, bỏ qua bước trung gian branch.

## 9. `git gc` — dọn rác object không còn ai tham chiếu

Object "unreachable" (không branch/tag/reflog nào trỏ tới, vd sau `reset --hard` — [Bài 10](./10_undo_changes.md)) không bị xóa NGAY, chỉ chờ `git gc` (chạy tự động định kỳ, hoặc thủ công) dọn dẹp sau 1 khoảng thời gian — đây là lý do `reflog` ([Bài 12 mục 4](./12_rewriting_history.md)) vẫn cứu được dữ liệu "đã mất" trong phần lớn trường hợp thực tế.

## Bài tập

1. **Khám phá blob/tree/commit thật**: làm theo đúng mục 5, tự tay xem nội dung blob, tree, commit của 1 repo bạn tạo.
2. **Verify content-addressable**: tạo 2 file KHÁC TÊN nhưng NỘI DUNG GIỐNG HỆT nhau, `git add` cả 2, dùng `git ls-files -s` verify chúng có CÙNG hash blob.
3. **Đọc trực tiếp `refs/heads/`**: tạo vài branch, `cat` từng file trong `.git/refs/heads/`, so sánh với output của `git log --oneline` để verify hash khớp nhau.
4. **Giải thích lại bằng lời** (không nhìn bài): tại sao `git commit --amend` tạo ra commit có hash HOÀN TOÀN khác, dù nội dung file gần như không đổi? (Gợi ý: commit object lưu CẢ thời gian, và nếu sửa message/file, nội dung object thay đổi → hash thay đổi theo mục 4).

## Tiếp theo
→ [Bài 18: Quy ước làm việc nhóm](./18_team_conventions.md)
