# Lộ Trình Học Git Chi Tiết

> Cấu trúc giống [Go/ROADMAP.md](../Go/ROADMAP.md), [Python/ROADMAP.md](../Python/ROADMAP.md), [SAP/ROADMAP.md](../SAP/ROADMAP.md), [MachineLearning/ROADMAP.md](../MachineLearning/ROADMAP.md): mỗi bài có file chi tiết riêng (lý thuyết + ý nghĩa cơ chế bên trong + lệnh thực hành + bài tập), kết thúc bằng dự án thực hành mô phỏng quy trình làm việc nhóm hoàn chỉnh.

---

## Giai đoạn 0 — Giới thiệu & Cài đặt

### [Bài 1: Giới thiệu Git & Cài đặt](./1_get_started.md)
- Git là gì, vì sao cần version control, phân biệt Git (công cụ) và GitHub/GitLab/Bitbucket (dịch vụ hosting).
- Cài đặt Git, cấu hình `user.name`/`user.email`, SSH key cho GitHub.
- **Bài tập:** Cài Git, cấu hình, tạo SSH key và thêm vào GitHub.

---

## Giai đoạn 1 — Git Cơ Bản (Fundamentals)

### [Bài 2: Khái niệm cốt lõi — 3 vùng làm việc](./2_core_concepts.md)
- Working Directory, Staging Area (Index), Repository (`.git`), snapshot vs diff — vì sao Git lưu trạng thái khác SVN/CVS.
- **Bài tập:** Vẽ sơ đồ 3 vùng, mô tả file di chuyển qua chúng thế nào.

### [Bài 3: Lệnh cơ bản — init, add, commit, log](./3_basic_commands.md)
- `git init`, `git add`, `git commit`, `git status`, `git log`, `git diff`.
- **Bài tập:** Tạo repo, thực hiện 5 commit có ý nghĩa, xem lại lịch sử.

### [Bài 4: Branching cơ bản](./4_branching.md)
- `git branch`, `git switch`/`git checkout`, `git merge` (fast-forward vs 3-way merge).
- **Bài tập:** Tạo branch mới, commit trên đó, merge về `main`.

### [Bài 5: Làm việc với Remote](./5_remote.md)
- `git remote`, `git push`, `git pull`, `git fetch`, `git clone`, upstream tracking.
- **Bài tập:** Đẩy repo local lên GitHub, clone về máy khác (hoặc thư mục khác), thử pull/push.

### [Bài 6: .gitignore & Quản lý file](./6_gitignore.md)
- `.gitignore` pattern, `git rm --cached`, xử lý file đã lỡ commit nhạy cảm.
- **Bài tập:** Viết `.gitignore` cho 1 project Python/Node, gỡ file đã lỡ track.

---

## Giai đoạn 2 — Git Trung Cấp (Intermediate)

### [Bài 7: Merge vs Rebase](./7_merge_vs_rebase.md)
- So sánh lịch sử sau merge vs rebase, khi nào dùng cái nào, `git rebase` cơ bản.
- **Bài tập:** Tạo cùng 1 tình huống 2 branch phân kỳ, thử cả merge và rebase, so sánh log.

### [Bài 8: Giải quyết Conflict](./8_conflicts.md)
- Đọc conflict marker, công cụ merge tool, conflict trong merge vs trong rebase.
- **Bài tập:** Cố tình tạo conflict, giải quyết bằng tay và bằng merge tool.

### [Bài 9: Stash & Cherry-pick](./9_stash_cherry_pick.md)
- `git stash` (save/pop/list/drop), `git cherry-pick`.
- **Bài tập:** Stash thay đổi dở dang để chuyển branch khẩn cấp, cherry-pick 1 commit từ branch khác.

### [Bài 10: Undo Changes — reset, revert, restore](./10_undo_changes.md)
- `git reset` (soft/mixed/hard), `git revert`, `git restore`, `git checkout -- file`.
- **Bài tập:** Thử cả 3 loại reset, quan sát khác biệt; dùng revert cho commit đã push.

### [Bài 11: Tags & Releases](./11_tags.md)
- Lightweight tag vs annotated tag, `git tag`, liên hệ semantic versioning.
- **Bài tập:** Tạo tag cho 1 "release", push tag lên remote.

### [Bài 12: Sửa lại lịch sử — rebase -i, amend, reflog](./12_rewriting_history.md)
- `git commit --amend`, `git rebase -i` (squash/reword/drop), `git reflog` (lưới an toàn cuối cùng).
- **Bài tập:** Gộp 3 commit thành 1 bằng squash, sửa message commit cũ, dùng reflog khôi phục sau khi "lỡ tay" reset --hard.

---

## Giai đoạn 3 — Git Nâng Cao (Advanced)

### [Bài 13: Git Workflow cho team](./13_workflows.md)
- Git Flow, GitHub Flow, Trunk-based Development — so sánh, khi nào dùng loại nào.
- **Bài tập:** Mô phỏng 1 tính năng đi qua đúng quy trình GitHub Flow (branch → PR → merge).

### [Bài 14: Pull Request & Code Review](./14_pull_requests.md)
- Quy trình PR, review comment, resolve conversation, squash merge vs merge commit vs rebase merge trên GitHub.
- **Bài tập:** Tạo PR thật trên GitHub (dùng 2 branch), tự review, thử cả 3 kiểu merge.

### [Bài 15: Submodule & Subtree](./15_submodules_subtrees.md)
- Khi nào cần quản lý repo lồng nhau, `git submodule`, `git subtree` — ưu nhược điểm.
- **Bài tập:** Thêm 1 submodule vào project, clone lại với `--recurse-submodules`.

### [Bài 16: Git Hooks](./16_git_hooks.md)
- `pre-commit`, `commit-msg`, `pre-push` — tự động hóa kiểm tra chất lượng code.
- **Bài tập:** Viết `pre-commit` hook chặn commit nếu có `console.log`/`fmt.Println` debug sót lại.

### [Bài 17: Git Internals — Git hoạt động thế nào bên trong](./17_git_internals.md)
- Object database (blob/tree/commit/tag), SHA-1 hash, `.git` directory structure, refs là gì.
- **Bài tập:** Dùng `git cat-file`/`git hash-object` khám phá object của 1 commit thật.

### [Bài 18: Quy ước làm việc nhóm](./18_team_conventions.md)
- Conventional Commits, chiến lược đặt tên branch, code review checklist.
- **Bài tập:** Viết lại lịch sử commit lộn xộn theo chuẩn Conventional Commits.

### [Bài 19: Debug với Git — bisect, blame, worktree](./19_debugging_tools.md)
- `git bisect` tìm commit gây lỗi, `git blame` truy vết thay đổi, `git worktree` làm việc song song nhiều branch.
- **Bài tập:** Dùng `git bisect` tìm commit gây lỗi trong 1 chuỗi commit giả lập.

### [Bài 20: CI/CD & Git — kết nối thực tế](./20_ci_cd_integration.md)
- GitHub Actions trigger theo push/PR, branch protection rule, tự động chạy test trước khi merge.
- **Bài tập:** Viết 1 GitHub Actions workflow chạy test khi có PR.

---

## Giai đoạn 4 — Dự Án Thực Hành (Capstone)

> **File chi tiết đầy đủ: [21_capstone_project.md](./21_capstone_project.md)**

**Đề bài:** Mô phỏng đầy đủ quy trình làm việc nhóm trên 1 project nhỏ — từ khởi tạo repo, thiết lập branch protection, làm việc qua feature branch + PR + code review, xử lý conflict thực tế, release có tag, tới thiết lập CI chạy test tự động.

## Gợi ý cách học
- Git học tốt nhất bằng cách LÀM, không phải đọc — mỗi bài nên thực hành trực tiếp trên 1 repo thử (tạo trong [scratchpad hoặc 1 thư mục ngoài project chính]) trước khi áp dụng vào project thật.
- Hiểu Git Internals ([Bài 17](./17_git_internals.md)) sớm giúp các lệnh "khó nhớ" như rebase/reset trở nên trực giác — vì bản chất chỉ là di chuyển con trỏ (ref) trỏ tới các object bất biến.
- Không sợ "làm hỏng" khi luyện tập — `git reflog` ([Bài 12](./12_rewriting_history.md)) gần như luôn cứu được, đây là lý do nên thử nghiệm mạnh dạn trên repo test.

## Tài liệu tham khảo
- Pro Git book (miễn phí): https://git-scm.com/book
- Git documentation: https://git-scm.com/docs
- Learn Git Branching (trực quan, tương tác): https://learngitbranching.js.org/
- Conventional Commits: https://www.conventionalcommits.org/
