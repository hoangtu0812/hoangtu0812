# Lộ Trình Học SAP Developer Chi Tiết

> Cấu trúc giống [Go/ROADMAP.md](../Go/ROADMAP.md) và [Python/ROADMAP.md](../Python/ROADMAP.md): mỗi bài có file chi tiết riêng, kết thúc bằng dự án capstone. Vì bạn (theo README) đã làm việc với **SAP ABAP, SAP CAP, SAP BTP** tại BSR, lộ trình này đi từ nền tảng ABAP cổ điển → ABAP hiện đại (CDS/RAP) → SAP BTP & CAP, để hệ thống hóa lại và lấp khoảng trống kiến thức, thay vì dạy lại từ số 0.

---

## Giai đoạn 0 — Giới thiệu & Cài đặt

### [Bài 1: Giới thiệu hệ sinh thái SAP & Cài đặt](./1_get_started.md)
- Bức tranh tổng thể: ECC/S4HANA (on-premise, ABAP), SAP BTP (cloud, CAP/Node.js/Java), Fiori/UI5 (frontend), điểm giao nhau giữa chúng.
- Cài đặt: Eclipse + ABAP Development Tools (ADT), SAP BTP Trial account, SAP Business Application Studio (BAS), Node.js + `@sap/cds-dk` (CDS CLI) cho CAP.
- **Bài tập:** Tạo BTP trial account, cài `@sap/cds-dk`, chạy `cds --version`; kết nối ADT tới 1 hệ thống ABAP (hoặc SAP steampunk/ABAP trial).

---

## Giai đoạn 1 — ABAP Cơ Bản (Classic Foundations)

### [Bài 2: Cú pháp ABAP cơ bản](./2_abap_syntax.md)
- `DATA`, kiểu built-in (`i`, `string`, `p`, `d`...), `TYPES`, toán tử, `WRITE`, khai báo inline (`DATA(x) = ...`).
- **Bài tập:** Viết report in "Hello World"; khai báo biến bằng cả cú pháp cũ và inline, so sánh.

### [Bài 3: Luồng điều khiển trong ABAP](./3_abap_control_flow.md)
- `IF/ELSEIF/ELSE`, `CASE`, `DO`, `WHILE`, `LOOP AT`, `CHECK`/`CONTINUE`/`EXIT`.
- **Bài tập:** FizzBuzz bằng ABAP; kiểm tra số nguyên tố.

### [Bài 4: Internal Table (Bảng nội bộ)](./4_internal_tables.md)
- `STANDARD TABLE`, `SORTED TABLE`, `HASHED TABLE`, `LOOP AT`, `READ TABLE` (với `BINARY SEARCH`/key), `APPEND`/`INSERT`/`DELETE`, `TYPES` structure.
- **Bài tập:** Đọc dữ liệu vào internal table, lọc và sắp xếp; so sánh hiệu năng `STANDARD` vs `HASHED` table khi tra cứu nhiều lần.

### [Bài 5: Modularization — Function Module & Method](./5_modularization.md)
- Subroutine (`FORM`/`PERFORM` — legacy), Function Module, Method trong class — khi nào dùng cái nào.
- **Bài tập:** Viết 1 Function Module tính tổng internal table; viết lại bằng static method.

### [Bài 6: ABAP Objects — OOP cơ bản](./6_abap_oo.md)
- `CLASS...DEFINITION/IMPLEMENTATION`, `INTERFACE`, kế thừa (`INHERITING FROM`), visibility (`PUBLIC/PROTECTED/PRIVATE`), constructor.
- **Bài tập:** Class `ZCL_SHAPE` (abstract) với subclass `ZCL_RECTANGLE`, `ZCL_CIRCLE`, method `GET_AREA`.

### [Bài 7: Open SQL — Truy vấn Database](./7_open_sql.md)
- `SELECT` hiện đại (inline declaration, `SELECT ... INTO TABLE @DATA(...)`), JOIN, `WHERE`, tránh SELECT * trong vòng lặp, CDS View dùng trong Open SQL.
- **Bài tập:** Viết SELECT lấy dữ liệu từ 2 bảng chuẩn (vd `SFLIGHT`/`SBOOK`) bằng JOIN, dùng inline declaration.

### [Bài 8: Xử lý lỗi trong ABAP](./8_abap_exceptions.md)
- Class-based exception (`TRY/CATCH/CLEANUP`), tự định nghĩa exception class kế thừa `CX_STATIC_CHECK`/`CX_DYNAMIC_CHECK`.
- **Bài tập:** Viết method raise custom exception khi input không hợp lệ, bắt bằng `TRY/CATCH`.

### [Bài 9: Debugging & ABAP Unit Testing](./9_abap_debug_testing.md)
- ABAP Debugger (breakpoint, watchpoint), ABAP Unit (`CLASS ... FOR TESTING`), test double đơn giản.
- **Bài tập:** Viết ABAP Unit test cho class Bài 6.

---

## Giai đoạn 2 — ABAP Hiện Đại: CDS & RAP

### [Bài 10: CDS View cơ bản](./10_cds_basics.md)
- Core Data Services: `DEFINE VIEW ENTITY`, annotation cơ bản (`@AccessControl`, `@UI`), field alias, tại sao CDS thay thế dần view cổ điển.
- **Bài tập:** Viết CDS view đơn giản trên 1 bảng chuẩn, thêm annotation `@UI.lineItem`.

### [Bài 11: CDS nâng cao](./11_cds_advanced.md)
- Association (`ASSOCIATION TO`), `path expression`, extend view, analytics annotation (`@Analytics.query`), virtual element.
- **Bài tập:** Viết CDS view có association tới 1 view khác, expose field qua path expression.

### [Bài 12: RAP — Giới thiệu (RESTful ABAP Programming Model)](./12_rap_intro.md)
- Kiến trúc RAP: Behavior Definition, Behavior Implementation, Business Object, managed vs unmanaged.
- **Bài tập:** Tạo Business Object đơn giản (CRUD) cho 1 custom table bằng RAP managed scenario.

### [Bài 13: RAP — Xây dựng Transactional App](./13_rap_transactional.md)
- Draft handling, validation (`VALIDATION ... FOR`), determination (`DETERMINATION ... FOR`), action tự định nghĩa.
- **Bài tập:** Thêm 1 validation (vd không cho save nếu field trống) và 1 action (vd "Approve") vào Business Object Bài 12.

### [Bài 14: OData Services](./14_odata_services.md)
- Expose CDS/RAP Business Object qua OData V2/V4, Service Definition + Service Binding, test qua SAP Gateway Client.
- **Bài tập:** Publish Business Object Bài 13 thành OData V4 service, test bằng Gateway Client hoặc Postman.

### [Bài 15: Authorization trong ABAP](./15_abap_authorization.md)
- Authorization Object, PFCG Role, `AUTHORITY-CHECK` trong code cổ điển, Access Control (DCL) trong CDS (`@AccessControl` + `DEFINE ROLE`).
- **Bài tập:** Viết CDS Access Control (DCL) giới hạn dữ liệu theo 1 field (vd chỉ thấy dữ liệu của phòng ban mình).

---

## Giai đoạn 3 — SAP BTP & CAP (Cloud-Native Development)

### [Bài 16: Giới thiệu SAP BTP](./16_btp_intro.md)
- Kiến trúc BTP: Cloud Foundry vs Kyma, Subaccount/Space, các service chính (HANA Cloud, XSUAA, Destination, Application Router).
- **Bài tập:** Tạo subaccount trial, khởi tạo 1 HANA Cloud instance (hoặc SQLite cho dev), khám phá BTP Cockpit.

### [Bài 17: CAP cơ bản (Cloud Application Programming Model)](./17_cap_basics.md)
- `cds init`, CDS trong CAP (data model + service definition), `cds watch`, so sánh CDS trong CAP vs CDS trong ABAP (Bài 10).
- **Bài tập:** `cds init taskapp`, định nghĩa entity `Tasks`, service `TaskService`, chạy `cds watch` và test qua `/tasks`.

### [Bài 18: CAP nâng cao — Custom Logic](./18_cap_advanced.md)
- Event handler (`srv.on`, `srv.before`, `srv.after`) bằng Node.js hoặc Java, kết nối HANA Cloud/PostgreSQL, `cds-dbm`/deploy database.
- **Bài tập:** Viết custom handler validate dữ liệu trước khi CREATE, viết handler tính toán field trước khi trả response.

### [Bài 19: Authentication & Authorization trên BTP (XSUAA)](./19_btp_auth.md)
- XSUAA service, `xs-security.json` (scope, role template), annotation `@requires`/`@restrict` trong CDS service, test với mock user (`cds bind`/`.cdsrc.json` mock auth).
- **Bài tập:** Thêm `@requires: 'authenticated-user'` cho service, `@restrict` chỉ cho phép role `admin` xóa task, test bằng mock user trong CAP.

### [Bài 20: Fiori Elements & Deployment](./20_fiori_deployment.md)
- Sinh Fiori Elements app từ OData service (List Report/Object Page), MTA (`mta.yaml`), `cf deploy`/`cds deploy`.
- **Bài tập:** Sinh Fiori Elements app cho `TaskService` (Bài 17), chạy thử local; viết `mta.yaml` cơ bản.

---

## Giai đoạn 4 — Dự Án: CAP Application Có Phân Quyền (Capstone)

> **File chi tiết đầy đủ: [21_capstone_project.md](./21_capstone_project.md)**

**Đề bài:** Xây dựng CAP application quản lý "Task" — cùng ngữ cảnh nghiệp vụ với dự án Go/Python để so sánh trực tiếp 3 stack (Go REST API, Python FastAPI, SAP CAP) trên cùng 1 bài toán: `user` chỉ CRUD task của mình, `admin` quản lý mọi task và user, có Fiori Elements UI, phân quyền qua XSUAA + CDS `@restrict`, deploy thử lên BTP trial.

### Cấu trúc thư mục đề xuất
```
taskapp/
├── db/
│   └── schema.cds          # data model: Users, Tasks
├── srv/
│   ├── task-service.cds    # service definition + @requires/@restrict
│   └── task-service.js     # custom handler (ownership check, validation)
├── app/
│   └── taskapp/             # Fiori Elements app (sinh từ service)
├── xs-security.json         # định nghĩa scope/role cho XSUAA
├── mta.yaml                 # deployment descriptor
├── package.json
└── .cdsrc.json               # mock auth cho dev local
```

### Các bước triển khai
1. `cds init taskapp`, định nghĩa data model `db/schema.cds` (`Users`, `Tasks` với `owner` association).
2. Định nghĩa `TaskService` trong `srv/task-service.cds`, expose `Tasks` entity.
3. Thêm `@requires: 'authenticated-user'` cho service, `@restrict` cho action xóa (chỉ `admin`).
4. Viết custom handler `srv/task-service.js`: `before CREATE` gán `owner` = user hiện tại; `before UPDATE/DELETE` kiểm tra ownership (user thường chỉ sửa task của mình, trừ khi có role admin).
5. Cấu hình mock user trong `.cdsrc.json` để test local với các role khác nhau (`admin`, `user`).
6. Sinh Fiori Elements List Report app từ `TaskService`.
7. Viết `xs-security.json` định nghĩa scope `Admin`, `User`, role template tương ứng.
8. Deploy database (HANA Cloud hoặc SQLite cho dev), test CRUD qua `/tasks`.
9. Viết test cơ bản (`cds test` hoặc Jest) cho handler ownership.
10. Đóng gói `mta.yaml`, deploy thử lên BTP trial bằng `cf deploy`/`mbt build`.

### Tiêu chí hoàn thành
- [ ] Service yêu cầu xác thực, user chưa đăng nhập bị từ chối (401).
- [ ] `user` không sửa/xóa được task của người khác (test bằng mock user khác role/id).
- [ ] `admin` xem/sửa/xóa được mọi task.
- [ ] Fiori Elements app hiển thị đúng danh sách task theo quyền.
- [ ] Deploy thử thành công lên BTP trial (hoặc chạy tốt bằng `cds watch` nếu chưa có quyền deploy).

## Gợi ý cách học
- Vì đã có nền ABAP/CAP/BTP, hãy tập trung nhiều thời gian hơn cho Giai đoạn 2-4 (RAP, CDS nâng cao, CAP, XSUAA) — đây là phần dùng trong dự án thực tế nhiều hơn ABAP cổ điển.
- So sánh xuyên suốt với 2 track kia: CDS `@restrict` (SAP) ↔ middleware RBAC (Go/Python), Behavior Definition (RAP) ↔ service layer (Go/Python), XSUAA scope ↔ JWT role claim.
- Ưu tiên làm capstone trên môi trường BTP trial (miễn phí) để có trải nghiệm gần nhất với công việc thực tế.

## Tài liệu tham khảo
- SAP CAP: https://cap.cloud.sap/docs/
- ABAP RAP: https://help.sap.com/docs/abap-cloud
- SAP BTP: https://help.sap.com/docs/btp
- ABAP Keyword Documentation (trong ADT hoặc SE38 → F1 trên từng lệnh)
