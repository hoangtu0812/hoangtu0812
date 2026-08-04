# Bài 20: Fiori Elements & Deployment

## Mục tiêu
- Sinh Fiori Elements app từ OData service (List Report/Object Page).
- Viết `mta.yaml`, deploy bằng `cf deploy`/`cds deploy`.

## 1. Fiori Elements là gì?

Fiori Elements sinh UI **hoàn toàn từ annotation của OData service** — không viết HTML/CSS/JS thủ công cho các màn hình CRUD chuẩn (danh sách, chi tiết, form sửa). Tương đương việc bạn KHÔNG cần viết React/Vue component cho CRUD cơ bản — chỉ cần khai báo đúng field/annotation, UI tự sinh.

## 2. Sinh Fiori Elements app từ CAP service ([Bài 17](./17_cap_basics.md))

```powershell
cds add fiori
```

Hoặc dùng **SAP Fiori tools** extension (VS Code/BAS) → "Open Application Generator" → chọn template **List Report Page** → trỏ tới OData service (`TaskService`) → chọn entity `Tasks`.

Kết quả sinh ra thư mục `app/taskapp/` với file `annotations.cds` chứa UI annotation:

```cds
// app/taskapp/annotations.cds
using TaskService as service from '../../srv/task-service';

annotate service.Tasks with @(
  UI.LineItem: [
    { Value: title, Label: 'Tiêu đề' },
    { Value: done, Label: 'Hoàn thành' },
    { Value: statusLabel, Label: 'Trạng thái' }
  ],
  UI.HeaderInfo: {
    TypeName: 'Task',
    TypeNamePlural: 'Tasks',
    Title: { Value: title }
  }
);
```

`UI.LineItem` tương đương `@UI.lineItem` bên ABAP CDS ([SAP Bài 10](./10_cds_basics.md)) — cùng khái niệm annotation UI, khác cú pháp một chút do 2 công cụ khác nhau (ADT vs CAP tooling), nhưng ý tưởng giống hệt: khai báo field nào hiển thị ở đâu, thay vì code UI thủ công.

## 3. Chạy thử local

```powershell
cds watch
```

Truy cập `http://localhost:4004/taskapp/webapp/index.html` (đường dẫn cụ thể tùy cấu hình `app/` sinh ra) — Fiori Elements List Report app hiển thị đầy đủ CRUD dựa trên `TaskService`, tự động tôn trọng `@restrict` đã cấu hình ở [Bài 19](./19_btp_auth.md) (user thường không thấy nút xóa nếu không có quyền).

## 4. `mta.yaml` — Multi-Target Application descriptor

```yaml
_schema-version: "3.1"
ID: taskapp
version: 1.0.0

modules:
  - name: taskapp-srv
    type: nodejs
    path: gen/srv
    requires:
      - name: taskapp-db
      - name: taskapp-auth
    provides:
      - name: srv-api
        properties:
          srv-url: ${default-url}

  - name: taskapp-db-deployer
    type: hdb
    path: gen/db
    requires:
      - name: taskapp-db

  - name: taskapp-app-content
    type: com.sap.application.content
    path: app
    requires:
      - name: taskapp-html5-repo-host

resources:
  - name: taskapp-db
    type: com.sap.xs.hdi-container
  - name: taskapp-auth
    type: org.cloudfoundry.managed-service
    parameters:
      service: xsuaa
      service-plan: application
      path: ./xs-security.json
  - name: taskapp-html5-repo-host
    type: org.cloudfoundry.managed-service
    parameters:
      service: html5-apps-repo
      service-plan: app-host
```

`mta.yaml` đóng gói toàn bộ ứng dụng (service backend + UI + database + auth) thành 1 đơn vị deploy — tương đương `docker-compose.yml` của Go/Python ([Go Bài 20](../Go/20_logging_config_deploy.md), [Python Bài 20](../Python/20_logging_config_deploy.md)) nhưng dành riêng cho hệ sinh thái BTP Cloud Foundry, khai báo cả service managed (XSUAA, HANA) chứ không chỉ container.

## 5. Build & Deploy

```powershell
npm install -g mbt   # Multi-Target Application Build Tool
mbt build              # đóng gói thành file .mtar

cf login
cf deploy mta_archives/taskapp_1.0.0.mtar
```

`mbt build` tương đương `docker build` (đóng gói); `cf deploy` tương đương `docker-compose up`/`docker push` + deploy — nhưng target là Cloud Foundry thay vì container registry.

## 6. Kiểm tra sau deploy

Trong BTP Cockpit → Space đã deploy → **Applications**: xem `taskapp-srv` (backend), `taskapp-app-content` đã chạy (status "Started"). Mở app qua **HTML5 Applications** hoặc link được cấp trong Cockpit — lúc này request thật sự đi qua XSUAA (không còn mock user như [Bài 19](./19_btp_auth.md)).

## Ví dụ tổng thể — luồng hoàn chỉnh

```
1. cds add fiori                        → sinh app/taskapp
2. Chỉnh annotations.cds                → khai báo UI.LineItem, HeaderInfo
3. cds watch                             → test local, mock auth
4. Viết mta.yaml                          → đóng gói srv + db + auth + app
5. mbt build && cf deploy ...mtar         → deploy thật lên BTP
6. Mở app qua Cockpit, đăng nhập XSUAA    → test với user/role thật
```

## Bài tập

1. **Sinh Fiori Elements app**: từ `TaskService` (Bài 17-19), sinh app List Report bằng `cds add fiori` hoặc Fiori tools, chạy thử qua `cds watch`.
2. **Tùy chỉnh annotation**: sửa `annotations.cds` để đổi label cột, thêm field `statusLabel` vào `UI.LineItem`.
3. **Viết `mta.yaml`**: viết file cơ bản như ví dụ cho project `taskapp`.
4. **Deploy thử**: nếu có quyền BTP trial, chạy `mbt build` + `cf deploy`, xác nhận app chạy trên Cockpit; nếu không đủ quyền/quota, mô tả lại từng bước bằng lời để hiểu quy trình.

## Tổng kết Giai đoạn 3
Bạn đã hoàn thành phần SAP BTP & CAP: giới thiệu BTP, CAP cơ bản/nâng cao, XSUAA authorization, và Fiori Elements/deployment. Bước cuối là ghép tất cả lại thành 1 ứng dụng capstone hoàn chỉnh, cùng bài toán nghiệp vụ với 2 track Go/Python.

## Tiếp theo
→ [Dự án Capstone: CAP Application có phân quyền](./21_capstone_project.md)
