# Bài 19: Debug Với Git — bisect, blame, worktree

## Mục tiêu
- Dùng `git bisect` tìm CHÍNH XÁC commit gây ra bug bằng binary search, không cần đoán mò.
- Dùng `git blame` truy vết ai/khi nào/tại sao 1 dòng code thay đổi.
- Dùng `git worktree` làm việc song song trên nhiều branch mà không cần stash/switch liên tục.

## 1. `git bisect` — tìm commit gây lỗi bằng binary search

**Tình huống:** code chạy tốt cách đây 2 tuần (100 commit trước), giờ phát hiện bug — nhưng không biết CHÍNH XÁC commit nào gây ra. Kiểm tra thủ công từng commit trong 100 commit là bất khả thi; `git bisect` áp dụng **binary search** ([Go Bài 3](../Go/3_control_flow.md) nếu bạn đã học tư duy thuật toán) để tìm ra chỉ trong $\log_2(100)\approx7$ lần thử.

```powershell
git bisect start
git bisect bad                    # đánh dấu commit HIỆN TẠI là "có bug"
git bisect good v1.0.0               # đánh dấu 1 commit CŨ đã biết chắc "không có bug" (vd tag release trước)

# Git tự động checkout tới commit Ở GIỮA khoảng good...bad
# Bạn TEST code tại điểm đó, rồi báo kết quả:
git bisect good    # nếu commit này KHÔNG có bug
git bisect bad       # nếu commit này CÓ bug

# Git tiếp tục thu hẹp khoảng, lặp lại tới khi tìm ra ĐÚNG commit đầu tiên gây bug
# Cuối cùng in ra: "a1b2c3d is the first bad commit"

git bisect reset    # kết thúc, quay lại branch/commit ban đầu trước khi bisect
```

## 2. `git bisect run` — tự động hóa hoàn toàn nếu có script test

Nếu có 1 script/lệnh trả về đúng exit code (0 = pass, khác 0 = fail), có thể để Git TỰ CHẠY toàn bộ quá trình, không cần bạn test thủ công từng bước:

```powershell
git bisect start HEAD v1.0.0
git bisect run npm test    # hoặc: python -m pytest, go test ./..., tùy stack

# Git tự động checkout từng commit, chạy lệnh, đọc exit code, tự thu hẹp khoảng
# In thẳng ra commit gây lỗi mà KHÔNG CẦN bạn can thiệp
```

Đây là cách dùng MẠNH MẼ NHẤT của bisect — biến việc "tìm bug trong 100 commit" từ vài giờ thủ công thành vài phút tự động.

## 3. `git blame` — ai đã viết dòng này, khi nào, commit nào

```powershell
git blame file.py                       # mỗi dòng kèm hash commit + tác giả + ngày tạo ra nó
git blame -L 10,20 file.py                 # chỉ xem dòng 10-20
git blame -w file.py                        # bỏ qua thay đổi CHỈ về khoảng trắng (whitespace) — tránh "nhiễu" do format lại
```

**Dùng để làm gì:** khi thấy 1 đoạn code "kỳ lạ" không hiểu tại sao viết vậy, `blame` chỉ ra ĐÚNG commit đã thêm nó — xem commit message + PR liên quan (nếu tuân theo quy ước ở [Bài 18](./18_team_conventions.md)) thường giải thích rõ lý do, tránh việc vô tình "sửa lại cho đẹp" rồi phá vỡ 1 fix quan trọng đã có chủ đích từ trước.

**"Blame" xuyên qua refactor** — nếu dòng code bị di chuyển qua nhiều commit refactor, dùng `git log --follow -p -- file.py` để xem toàn bộ lịch sử file kể cả qua các lần đổi tên/di chuyển.

## 4. `git worktree` — làm việc trên NHIỀU branch CÙNG LÚC, không cần switch

**Vấn đề:** đang code dở trên `feature-a` (chưa muốn commit/stash), cần gấp kiểm tra code trên `main` (vd để so sánh hành vi cũ) — `git switch` sẽ đổi TOÀN BỘ Working Directory, buộc phải stash trước ([Bài 9](./9_stash_cherry_pick.md)).

```powershell
git worktree add ../main-check main    # tạo 1 THƯ MỤC MỚI, checkout branch "main" vào đó — SONG SONG với thư mục hiện tại
```

Giờ bạn có 2 thư mục làm việc, CÙNG 1 repo (chia sẻ chung `.git/`), mỗi thư mục đứng ở 1 branch khác nhau — không cần stash, không cần switch, có thể mở cả 2 trong 2 cửa sổ editor song song.

```powershell
git worktree list         # xem mọi worktree đang có
git worktree remove ../main-check    # xong việc, xóa bỏ (không ảnh hưởng repo gốc)
```

**Ứng dụng thực tế phổ biến nhất:** chạy test/build trên 1 branch trong lúc vẫn code tiếp ở branch khác; so sánh hành vi giữa 2 version mà không mất thời gian chuyển đổi qua lại liên tục.

## 5. Kết hợp bisect + blame trong quy trình debug thực tế

```
1. Phát hiện bug -> git bisect (tự động nếu có test) -> tìm ra commit gây lỗi
2. git show <commit-gay-loi>    -> xem CHÍNH XÁC thay đổi gì đã gây ra vấn đề
3. git blame vùng code liên quan -> hiểu bối cảnh rộng hơn (đã sửa gì trước đó, tại sao)
4. Quyết định: git revert commit đó (Bài 10), hoặc fix thêm dựa trên hiểu biết đầy đủ
```

## Ví dụ đầy đủ

```powershell
git init bisect-demo; cd bisect-demo
echo "def add(a, b): return a + b" > math.py
git add .; git commit -m "commit 1: hàm add đúng"

echo "def add(a, b): return a - b" > math.py    # cố tình tạo bug ở đây
git add .; git commit -m "commit 2: refactor (lỡ gây bug)"

echo "def add(a, b): return a - b  # thêm comment" > math.py
git add .; git commit -m "commit 3: thêm comment"

echo "def add(a, b): return a - b" > math.py
echo "def multiply(a, b): return a * b" >> math.py
git add .; git commit -m "commit 4: thêm hàm multiply"

# Giờ giả sử bạn PHÁT HIỆN add() bị sai, không nhớ từ commit nào
git bisect start
git bisect bad HEAD
git bisect good HEAD~3    # commit 1, biết chắc đúng

# Git checkout tới commit giữa (commit 2 hoặc 3) -> bạn kiểm tra math.py
cat math.py    # nếu thấy "a - b" -> bad; nếu "a + b" -> good
git bisect bad    # (ví dụ, nếu đang ở commit có bug)
# Git tiếp tục thu hẹp, cuối cùng báo "commit 2 is the first bad commit"

git bisect reset
git blame math.py    # xác nhận đúng commit 2 đã sửa dòng này
```

## Bài tập

1. **`git bisect` thủ công**: làm theo đúng "Ví dụ đầy đủ", tự tay chạy qua toàn bộ quy trình bisect tìm ra "commit 2" là nguồn gốc bug.
2. **`git bisect run` tự động**: viết 1 script Python/shell đơn giản kiểm tra `add(2,3) == 5`, dùng `git bisect run` để tự động tìm bug thay vì kiểm tra thủ công.
3. **`git blame`**: trên 1 file bất kỳ có lịch sử vài commit, dùng `git blame` xác định dòng cụ thể do commit nào tạo ra, xem message commit đó bằng `git show`.
4. **`git worktree`**: tạo 1 worktree phụ trỏ tới branch khác, thử sửa file ở CẢ 2 thư mục cùng lúc, verify chúng độc lập nhưng chia sẻ chung lịch sử Git (`git log` ở worktree phụ thấy đủ lịch sử).

## Tiếp theo
→ [Bài 20: CI/CD & Git — kết nối thực tế](./20_ci_cd_integration.md)
