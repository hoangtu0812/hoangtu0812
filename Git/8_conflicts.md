# Bài 8: Giải Quyết Conflict

## Mục tiêu
- Đọc hiểu conflict marker, giải quyết conflict thủ công một cách tự tin.
- Phân biệt trải nghiệm conflict trong merge vs trong rebase.

## 1. Conflict xảy ra khi nào?

Git tự động kết hợp thay đổi từ 2 nguồn khác nhau ([Bài 4 mục 4](./4_branching.md)) — nhưng nếu CẢ 2 nguồn sửa **cùng 1 vùng dòng** của cùng 1 file theo cách KHÁC NHAU, Git không thể tự quyết định giữ bản nào — nó dừng lại, đánh dấu vùng xung đột, để BẠN quyết định.

## 2. Đọc Conflict Marker

```
<<<<<<< HEAD
Đây là nội dung ở branch hiện tại (nơi bạn đang đứng, thường là main)
=======
Đây là nội dung ở branch đang merge vào (feature)
>>>>>>> feature
```

- `<<<<<<< HEAD` tới `=======`: nội dung từ branch bạn ĐANG ĐỨNG.
- `=======` tới `>>>>>>> feature`: nội dung từ branch đang được merge VÀO.

**Nhiệm vụ của bạn:** sửa lại đoạn này thành nội dung ĐÚNG bạn muốn (có thể giữ 1 trong 2, kết hợp cả 2, hoặc viết lại hoàn toàn khác), rồi **xóa hết các dòng marker** (`<<<<<<<`, `=======`, `>>>>>>>`).

## 3. Quy trình giải quyết conflict khi merge

```powershell
git merge feature
# CONFLICT (content): Merge conflict in file.txt
# Automatic merge failed; fix conflicts and then commit the result.

git status    # liệt kê rõ file nào đang conflict ("both modified")

# Mở file.txt, sửa thủ công, xóa marker
git add file.txt          # đánh dấu ĐÃ giải quyết xong file này
git status                  # verify không còn file nào "Unmerged paths"

git commit                    # HOÀN THÀNH merge — Git tự điền sẵn message "Merge branch 'feature'"
```

**Nếu muốn HỦY merge giữa chừng** (quyết định chưa sẵn sàng xử lý conflict lúc này):

```powershell
git merge --abort    # quay lại trạng thái TRƯỚC khi merge, như chưa từng chạy lệnh merge
```

## 4. Quy trình giải quyết conflict khi rebase — KHÁC merge

Rebase áp lại TỪNG commit một ([Bài 7 mục 2](./7_merge_vs_rebase.md)) — nếu nhiều commit đều gây conflict, bạn phải giải quyết **LẦN LƯỢT từng commit**, không phải 1 lần duy nhất như merge:

```powershell
git rebase main
# CONFLICT (content): Merge conflict in file.txt
# (sửa file.txt, xóa marker)

git add file.txt
git rebase --continue    # tiếp tục áp commit TIẾP THEO — có thể lại conflict nữa nếu có nhiều commit

# Nếu muốn bỏ dở giữa chừng
git rebase --abort         # quay lại TRẠNG THÁI TRƯỚC khi bắt đầu rebase
```

**Điểm khác biệt quan trọng cần nhớ:** merge conflict giải quyết 1 LẦN cho TOÀN BỘ thay đổi giữa 2 nhánh; rebase conflict có thể lặp lại NHIỀU LẦN (mỗi lần cho 1 commit) nếu các commit riêng lẻ đều đụng tới vùng xung đột.

## 5. Công cụ hỗ trợ merge trực quan (Merge Tool)

```powershell
git config --global merge.tool vscode
git config --global mergetool.vscode.cmd 'code --wait $MERGED'

git mergetool    # mở VS Code (hoặc tool đã cấu hình) hiển thị conflict trực quan 3 cột
```

VS Code có sẵn giao diện xử lý conflict rất trực quan (nút "Accept Current Change"/"Accept Incoming Change"/"Accept Both Changes") khi mở file có conflict marker — không cần cấu hình `mergetool` riêng nếu chỉ dùng VS Code, chỉ cần mở file conflict trực tiếp.

## 6. Conflict trong file KHÔNG PHẢI code (JSON, YAML, binary)

Với file cấu hình (JSON/YAML), conflict marker vẫn chèn vào y hệt — cần cẩn thận giữ cú pháp hợp lệ sau khi xóa marker (dễ để sót dấu ngoặc/thụt lề sai). Với file **binary** (ảnh, file nén), Git KHÔNG THỂ hiển thị "diff" dễ đọc — thường phải chọn giữ HẲN 1 trong 2 phiên bản:

```powershell
git checkout --ours path/to/image.png     # giữ bản của branch hiện tại (HEAD)
git checkout --theirs path/to/image.png     # giữ bản đang merge vào
git add path/to/image.png
```

**Lưu ý:** trong lúc REBASE, ý nghĩa `--ours`/`--theirs` bị **ĐẢO NGƯỢC** so với merge (vì Git coi commit đang được "áp lại" là "theirs") — luôn kiểm tra kỹ bằng `git status`/xem nội dung trước khi quyết định, đừng chỉ nhớ máy móc.

## Ví dụ đầy đủ

```powershell
git init conflict-demo; cd conflict-demo
echo "Xin chao" > greeting.txt
git add .; git commit -m "Initial"

git switch -c feature
echo "Xin chao cac ban" > greeting.txt
git add .; git commit -m "Sua greeting (feature)"

git switch main
echo "Xin chao moi nguoi" > greeting.txt
git add .; git commit -m "Sua greeting (main)"

git merge feature
# CONFLICT — mở greeting.txt, sẽ thấy:
# <<<<<<< HEAD
# Xin chao moi nguoi
# =======
# Xin chao cac ban
# >>>>>>> feature

# Quyết định giữ: "Xin chao moi nguoi va cac ban"
# Xóa hết marker, lưu file

git add greeting.txt
git commit
git log --oneline --graph
```

## Bài tập

1. **Giải quyết conflict merge cơ bản**: làm theo đúng "Ví dụ đầy đủ", tự tay giải quyết conflict.
2. **`merge --abort`**: tạo lại tình huống conflict, thử `git merge --abort`, verify repo quay về trạng thái trước merge.
3. **Conflict trong rebase — lặp lại nhiều lần**: tạo 3 commit trên `feature` đều sửa cùng 1 vùng dòng, rebase lên `main` cũng đã sửa vùng đó — quan sát và giải quyết conflict LẦN LƯỢT cho từng commit bằng `rebase --continue`.
4. **Dùng merge tool**: cấu hình `git mergetool` với VS Code, thử giải quyết 1 conflict bằng giao diện trực quan thay vì sửa tay text.

## Tiếp theo
→ [Bài 9: Stash & Cherry-pick](./9_stash_cherry_pick.md)
