# Bài 15: Authorization trong ABAP

## Mục tiêu
- Authorization Object, PFCG Role, `AUTHORITY-CHECK`.
- Access Control (DCL) trong CDS — `@AccessControl` + `DEFINE ROLE`.
- Đây là phần tương đương RBAC của Go/Python ([Go Bài 19](../Go/19_auth.md), [Python Bài 19](../Python/19_auth.md)).

## 1. Authorization Object — đơn vị phân quyền cơ bản trong ABAP

Authorization Object định nghĩa **field nào** dùng để kiểm soát quyền (vd `ACTVT` = loại hành động: 01-Create, 02-Change, 03-Display..., cộng thêm field nghiệp vụ như `BUKRS` = mã công ty). Đây là khái niệm cấp thấp hơn nhiều so với "role" đơn giản (`admin`/`user`) của Go/Python — ABAP cho phép phân quyền **rất chi tiết theo giá trị dữ liệu**.

## 2. `AUTHORITY-CHECK` — kiểm tra quyền trong code cổ điển

```abap
AUTHORITY-CHECK OBJECT 'Z_TASK'
  ID 'ACTVT' FIELD '02'      " '02' = Change
  ID 'OWNER' FIELD lv_owner_id.

IF sy-subrc <> 0.
  MESSAGE 'Bạn không có quyền sửa task này' TYPE 'E'.
ENDIF.
```

`sy-subrc = 0` sau `AUTHORITY-CHECK` nghĩa là user hiện tại **có** quyền — tương đương middleware `RequireRole` của Go/Python nhưng kiểm tra tại **bất kỳ điểm nào trong code** thay vì chỉ ở tầng router/middleware.

## 3. PFCG Role — gói các Authorization Object lại, gán cho user

PFCG (transaction quản lý role trong SAP GUI) cho phép admin:
1. Tạo 1 Role (vd `Z_TASK_ADMIN`, `Z_TASK_USER`).
2. Gắn Authorization Object vào role, chỉ định giá trị field cụ thể (vd `ACTVT = 01,02,03,06` cho admin; `ACTVT = 02,03` và `OWNER = &user_id` cho user thường).
3. Gán role cho user (transaction `SU01`).

Đây tương đương việc gán `role: "admin"` vào JWT claim của Go/Python ([Go Bài 19](../Go/19_auth.md)), nhưng PFCG Role được quản lý **hoàn toàn qua giao diện, không cần deploy code** — thay đổi quyền không cần release mới.

## 4. Access Control (DCL) trong CDS — cách hiện đại, tích hợp với RAP

```abap
@EndUserText.label: 'Authorization for Task'
@MappingRole: true
define role ZTASK_DCL {
  grant select on ZI_Task
  where OwnerId = aspect pfcg_auth( 'Z_TASK', 'OWNER', actvt = '03' );
}
```

Access Control (DCL — Data Control Language) gắn trực tiếp vào CDS View qua annotation `@AccessControl.authorizationCheck: #CHECK` (đã thấy ở [Bài 10](./10_cds_basics.md)) — RAP tự động áp dụng filter phân quyền **ngay tại tầng data**, không cần viết `AUTHORITY-CHECK` thủ công trong Behavior Implementation.

### Ví dụ: chỉ cho xem task của chính mình (trừ khi có quyền admin xem tất cả)

```abap
define role ZTASK_DCL {
  grant select, insert, update, delete on ZI_Task
  where ( OwnerId = $session.user ) or aspect pfcg_auth( 'Z_TASK_ADMIN', 'ACTVT', actvt = '*' );
}
```

`$session.user` lấy user hiện tại đang đăng nhập — CDS tự động thêm điều kiện `WHERE OwnerId = <current_user>` vào MỌI query trên `ZI_Task`, trừ khi user đó có Authorization Object `Z_TASK_ADMIN`. Đây chính xác là logic ownership check bạn đã viết thủ công ở tầng service trong Go/Python ([Go Bài 19 mục 7](../Go/19_auth.md), [Python Bài 19 mục "Phân quyền theo ownership"](../Python/19_auth.md)) — nhưng RAP cho phép khai báo nó **ở tầng data model**, áp dụng tự động cho MỌI entry point (OData, RAP action, report...) mà không cần lặp lại logic.

## 5. So sánh trực tiếp với Go/Python RBAC

| | ABAP (PFCG + DCL) | Go/Python (JWT + middleware) |
|---|---|---|
| Định nghĩa role | PFCG (giao diện, không cần code) | Code (`role: "admin"` trong DB/JWT) |
| Kiểm tra quyền request-level | `AUTHORITY-CHECK` trong code, hoặc DCL tự động | Middleware `RequireRole()` |
| Kiểm tra quyền row-level (ownership) | DCL với `$session.user` — tự động ở tầng data | Logic thủ công trong service layer |
| Thay đổi quyền không cần deploy | Có (PFCG) | Không (phải sửa code hoặc DB role mapping) |
| Độ chi tiết | Rất cao (theo field, theo giá trị) | Thường đơn giản hơn (role-based, đôi khi thêm ownership) |

## Bài tập

1. **`AUTHORITY-CHECK` cơ bản**: viết code dùng `AUTHORITY-CHECK` với 1 Authorization Object chuẩn có sẵn (vd `S_TABU_DIS`), kiểm tra `sy-subrc`.
2. **Tìm hiểu PFCG**: nếu có quyền truy cập hệ thống công ty, mở transaction PFCG, xem 1 role thật đang được dùng — liệt kê các Authorization Object trong đó.
3. **DCL cho Business Object `Task`**: viết Access Control (DCL) cho `ZI_Task` (Bài 12), giới hạn user chỉ thấy task của chính mình (dùng `$session.user`), test bằng cách đăng nhập 2 user khác nhau qua Preview.
4. **So sánh**: viết ghi chú so sánh cách ABAP xử lý authorization (khai báo, tách biệt code nghiệp vụ) với middleware RBAC bạn viết ở Go/Python — ưu/nhược điểm mỗi cách.

## Tổng kết Giai đoạn 2
Bạn đã nắm phần ABAP hiện đại: CDS View, RAP (Business Object, validation/determination/action), OData Service, và Authorization. Đây là bộ kỹ năng dùng nhiều nhất trong các dự án S/4HANA hiện tại. Giai đoạn 3 sẽ mở rộng sang SAP BTP & CAP — phần cloud-native.

## Tiếp theo
→ [Bài 16: Giới thiệu SAP BTP](./16_btp_intro.md)
