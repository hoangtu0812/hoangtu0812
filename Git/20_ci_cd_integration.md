# Bài 20: CI/CD & Git — Kết Nối Thực Tế

## Mục tiêu
- Hiểu CI/CD dựa vào sự kiện Git (push, PR) để trigger tự động thế nào.
- Viết 1 GitHub Actions workflow cơ bản chạy test khi có PR.
- Kết hợp Branch Protection để bắt buộc CI pass trước khi merge.

## 1. CI/CD dựa vào Git event ra sao

CI/CD (Continuous Integration/Continuous Deployment) là quy trình TỰ ĐỘNG chạy test/build/deploy mỗi khi code thay đổi — GitHub Actions (hay Jenkins, GitLab CI...) LẮNG NGHE các sự kiện Git cụ thể (`push`, `pull_request`, tạo `tag` mới — [Bài 11](./11_tags.md)) để tự động kích hoạt workflow tương ứng.

## 2. Cấu trúc GitHub Actions Workflow cơ bản

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4       # tự động clone ĐÚNG commit đã trigger event này

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: "3.12"

      - name: Cài dependency
        run: pip install -r requirements.txt

      - name: Chạy test
        run: pytest --cov

      - name: Chạy linter
        run: flake8 .
```

`on.push`/`on.pull_request` chỉ định CHÍNH XÁC sự kiện Git nào kích hoạt workflow — `pull_request` trigger MỖI LẦN có commit mới push vào branch nguồn của PR đang mở, cho phép CI chạy lại tự động khi tác giả sửa theo góp ý review ([Bài 14 mục 2](./14_pull_requests.md)).

## 3. Workflow đa bước — build, test, deploy tách biệt

```yaml
name: Full Pipeline

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm ci
      - run: npm test

  deploy:
    needs: test              # CHỈ chạy nếu job "test" thành công
    if: github.ref == 'refs/heads/main' && github.event_name == 'push'   # CHỈ deploy khi push thẳng vào main (không phải PR)
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Deploy to production
        run: ./deploy.sh
```

`needs: test` đảm bảo `deploy` KHÔNG BAO GIỜ chạy nếu `test` fail — đây là lớp bảo vệ tự động, không phụ thuộc con người nhớ "chạy test trước khi deploy".

## 4. Kết hợp Branch Protection — biến CI thành RÀO CHẮN bắt buộc

Nhắc lại [Bài 14 mục 7](./14_pull_requests.md): Branch Protection Rule có thể yêu cầu **status check** (chính là kết quả job CI) PASS trước khi nút "Merge" trên GitHub được phép bấm. Kết hợp với workflow ở mục 2, quy trình trở thành:

```
1. Tạo PR
2. GitHub TỰ ĐỘNG chạy CI (workflow "test")
3. Nếu CI FAIL -> nút Merge bị KHÓA, dù đã có approval
4. Tác giả sửa lỗi, push thêm commit -> CI chạy lại
5. CI PASS + có approval -> nút Merge mở khóa
```

Đây là cách BIẾN quy tắc "phải test trước khi merge" từ 1 quy ước bằng lời (dễ bị quên/bỏ qua) thành 1 RÀO CHẮN kỹ thuật không thể lách qua (trừ khi admin cố tình tắt rule).

## 5. Trigger theo Tag — tự động deploy khi có release mới

```yaml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'    # CHỈ trigger khi push tag khớp pattern Semantic Versioning (Bài 11 mục 6)

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build & Publish
        run: |
          npm run build
          npm publish
```

Kết hợp trực tiếp với quy trình tag ([Bài 11](./11_tags.md)): mỗi lần bạn `git tag v1.2.0 && git push origin v1.2.0`, workflow này TỰ ĐỘNG build và publish — không cần thao tác thủ công nào thêm.

## 6. `git log`/CI kết hợp để debug pipeline fail

Khi CI fail, luôn kiểm tra CHÍNH XÁC commit nào đang được test (`github.sha` trong context của GitHub Actions tương ứng đúng 1 commit hash cụ thể) — nếu nghi ngờ 1 job CI fail "không rõ lý do", dùng kỹ thuật tương tự `git bisect` ([Bài 19 mục 1](./19_debugging_tools.md)) trên chính lịch sử CI runs để xác định commit nào bắt đầu gây fail.

## Ví dụ đầy đủ: thiết lập CI cho project mới

```yaml
# .github/workflows/ci.yml
name: CI

on:
  pull_request:
    branches: [main]

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: "npm"

      - name: Install dependencies
        run: npm ci

      - name: Run linter
        run: npm run lint

      - name: Run tests with coverage
        run: npm test -- --coverage

      - name: Upload coverage report
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage/
```

Sau khi commit file này vào `.github/workflows/`, push lên GitHub, mở 1 PR thử — tab "Checks" trên PR sẽ hiện kết quả job `lint-and-test` chạy tự động.

## Bài tập

1. **Viết CI workflow cơ bản**: tạo `.github/workflows/ci.yml` cho 1 project thật (Go/Python/Node.js bạn có sẵn), chạy test tự động khi mở PR.
2. **Thiết lập Branch Protection**: bật yêu cầu status check PASS cho `main`, thử tạo PR có test FAIL — verify nút merge bị khóa; sửa lỗi, verify nút mở khóa lại.
3. **Workflow deploy theo tag**: viết workflow trigger khi push tag `v*.*.*`, chỉ cần in ra "Deploying version $GITHUB_REF" (giả lập, không cần deploy thật) để verify trigger đúng.
4. **Multi-job với `needs`**: viết workflow có 2 job (`test` và `build`), `build` chỉ chạy sau khi `test` pass — verify bằng cách cố tình cho `test` fail, quan sát `build` bị bỏ qua.

## Tổng kết Giai đoạn 3
Bạn đã nắm workflow team (Git Flow/GitHub Flow/Trunk-based), quy trình Pull Request/Code Review, submodule/subtree, Git Hooks, cơ chế bên trong Git (Internals), quy ước commit chuẩn, công cụ debug (bisect/blame/worktree), và kết nối Git với CI/CD. Đây là bộ kỹ năng đầy đủ để làm việc chuyên nghiệp trong bất kỳ team phát triển phần mềm nào. Bước cuối là ghép mọi thứ vào 1 dự án thực hành hoàn chỉnh.

## Tiếp theo
→ [Dự án thực hành: Mô phỏng quy trình làm việc nhóm hoàn chỉnh](./21_capstone_project.md)
