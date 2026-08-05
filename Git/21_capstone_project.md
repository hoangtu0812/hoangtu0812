# Dự Án Thực Hành: Mô Phỏng Quy Trình Làm Việc Nhóm Hoàn Chỉnh

## Mục tiêu
Ghép toàn bộ kiến thức từ Bài 1 → Bài 20 thành 1 quy trình thực hành đầy đủ — từ khởi tạo repo, thiết lập bảo vệ, làm việc qua feature branch + PR + review, xử lý conflict thực tế, release có tag, tới CI tự động. Nếu có thể, rủ 1 người bạn/đồng nghiệp cùng làm để trải nghiệm ĐÚNG tình huống nhiều người — nếu làm 1 mình, dùng 2 tài khoản GitHub hoặc 2 thư mục clone để mô phỏng "người thứ 2".

## Đề bài
Xây dựng 1 project nhỏ (vd Todo API bằng ngôn ngữ bạn đã học — [Go](../Go/ROADMAP.md)/[Python](../Python/ROADMAP.md)) theo ĐÚNG quy trình GitHub Flow ([Bài 13 mục 3](./13_workflows.md)), với ít nhất 2 "người" cùng đóng góp.

## Bước 1: Khởi tạo repo & thiết lập ban đầu

```powershell
mkdir team-workflow-demo; cd team-workflow-demo
git init
echo "# Team Workflow Demo" > README.md
git add README.md
git commit -m "docs: khởi tạo project"

git remote add origin git@github.com:username/team-workflow-demo.git
git push -u origin main
```

Viết `.gitignore` phù hợp ([Bài 6](./6_gitignore.md)) cho ngôn ngữ bạn chọn.

## Bước 2: Thiết lập Branch Protection cho `main`

Trên GitHub: Settings → Branches → Add rule cho `main` — bật "Require a pull request before merging", "Require status checks to pass" (sẽ có tác dụng sau khi thêm CI ở Bước 6).

## Bước 3: Thiết lập CI cơ bản ([Bài 20](./20_ci_cd_integration.md))

```yaml
# .github/workflows/ci.yml
name: CI
on:
  pull_request:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Chạy test
        run: echo "Chạy test ở đây (thay bằng lệnh test thật của project bạn)"
```

Commit và push trực tiếp (lần cuối được phép push thẳng, vì `main` rule cần code ĐÃ có sẵn workflow này trước khi áp dụng nghiêm ngặt).

## Bước 4: Thiết lập Git Hook local ([Bài 16](./16_git_hooks.md))

Thêm `pre-commit` hook (qua Husky/`pre-commit` framework) chặn debug code hoặc chạy linter nhanh trước mỗi commit.

## Bước 5: "Người A" phát triển tính năng 1 — quy trình GitHub Flow đầy đủ

```powershell
git switch -c feature/them-endpoint-tao-task
# ... code endpoint POST /tasks ...
git add .
git commit -m "feat(tasks): thêm endpoint tạo task mới"    # theo Conventional Commits — Bài 18
git push -u origin feature/them-endpoint-tao-task
```

Mở PR trên GitHub với description đầy đủ ([Bài 14 mục 3](./14_pull_requests.md)). CI tự động chạy ([Bài 20](./20_ci_cd_integration.md)).

## Bước 6: "Người B" review PR của Người A

Để lại ít nhất 1 comment dạng "suggestion", 1 câu hỏi ([Bài 14 mục 5](./14_pull_requests.md)). Người A phản hồi/sửa, push thêm commit — verify PR tự động cập nhật, CI chạy lại.

## Bước 7: Tạo CONFLICT thực tế giữa 2 người

Trong lúc PR của Người A đang mở, "Người B" tạo branch riêng SỬA CÙNG 1 VÙNG code (vd cùng file cấu hình route), merge branch của mình vào `main` TRƯỚC:

```powershell
# Người B
git switch main; git pull
git switch -c feature/cau-hinh-route
# ... sửa route.py ở ĐÚNG vùng Người A cũng đang sửa ...
git commit -m "feat(routes): tái cấu trúc đăng ký route"
git push -u origin feature/cau-hinh-route
# Mở PR, review nhanh, merge vào main (squash merge — Bài 14 mục 4)
```

Giờ PR của Người A bị "outdated" so với `main` — GitHub sẽ báo conflict. Người A xử lý:

```powershell
git switch feature/them-endpoint-tao-task
git fetch origin
git rebase origin/main    # hoặc merge, tùy quy ước team (Bài 7 mục 7)
# giải quyết conflict thủ công (Bài 8)
git add route.py
git rebase --continue
git push --force-with-lease origin feature/them-endpoint-tao-task   # Bài 12 mục 5, vì đã rebase commit đã push
```

## Bước 8: Hoàn thành merge, dọn dẹp

```powershell
# Sau khi PR được approve, merge trên GitHub (squash and merge — Bài 14 mục 4)
git switch main
git pull origin main
git branch -d feature/them-endpoint-tao-task
```

## Bước 9: Release có tag ([Bài 11](./11_tags.md))

```powershell
git tag -a v1.0.0 -m "Phiên bản đầu tiên: tạo và liệt kê task"
git push origin v1.0.0
```

Nếu đã thiết lập workflow trigger theo tag ([Bài 20 mục 5](./20_ci_cd_integration.md)), verify nó tự động chạy.

## Bước 10: Xử lý hotfix khẩn cấp

Giả lập phát hiện bug nghiêm trọng trên `main` NGAY SAU khi release:

```powershell
git switch main; git pull
git switch -c hotfix/sua-loi-crash-khi-tao-task
# ... fix ...
git commit -m "fix(tasks): sửa lỗi crash khi title rỗng"
git push -u origin hotfix/sua-loi-crash-khi-tao-task
# Mở PR, review NHANH (ưu tiên cao), merge ngay sau khi CI pass
```

```powershell
git switch main; git pull
git tag -a v1.0.1 -m "Hotfix: sửa lỗi crash khi tạo task với title rỗng"
git push origin v1.0.1
```

## Bước 11 (nếu "lỡ tay"): thực hành cứu nguy bằng reflog

Cố tình (trên branch TEST riêng, không phải `main`) chạy `git reset --hard` lùi lại vài commit, rồi dùng `git reflog` ([Bài 12 mục 4](./12_rewriting_history.md)) khôi phục lại — để chắc chắn bạn tự tin với "lưới an toàn" này trước khi kết thúc dự án.

## Checklist hoàn thành

- [ ] Repo có `.gitignore` phù hợp, Branch Protection bật cho `main`.
- [ ] CI chạy tự động trên mọi PR ([Bài 20](./20_ci_cd_integration.md)).
- [ ] Git Hook local hoạt động (chặn được ít nhất 1 loại lỗi trước commit).
- [ ] Có ít nhất 2 PR hoàn chỉnh với description, review comment, CI pass.
- [ ] Đã xử lý ít nhất 1 conflict THỰC TẾ (không phải giả lập đơn giản) giữa 2 nhánh.
- [ ] Có ít nhất 1 release với annotated tag theo Semantic Versioning.
- [ ] Đã xử lý 1 tình huống hotfix khẩn cấp đúng quy trình.
- [ ] Lịch sử commit trên `main` tuân theo Conventional Commits ([Bài 18](./18_team_conventions.md)), đọc `git log --oneline` thấy RÕ RÀNG mạch câu chuyện phát triển.
- [ ] Tự tin dùng `reflog` để cứu nguy khi "lỡ tay" thao tác nguy hiểm.

## Mở rộng (tùy chọn)
- Thử áp dụng Git Flow ([Bài 13 mục 2](./13_workflows.md)) thay vì GitHub Flow cho cùng project, so sánh trải nghiệm.
- Thêm submodule cho 1 thư viện dùng chung ([Bài 15](./15_submodules_subtrees.md)).
- Viết workflow CI đầy đủ hơn: lint + test + build + deploy staging tự động khi merge vào `main`.

---
Hoàn thành dự án này nghĩa là bạn đã trải nghiệm ĐẦY ĐỦ vòng đời làm việc thực tế với Git trong môi trường nhóm — không chỉ biết lệnh, mà biết ÁP DỤNG đúng lúc, đúng quy trình, và quan trọng nhất: không còn sợ "làm hỏng" vì đã hiểu rõ cách cứu nguy khi cần.
