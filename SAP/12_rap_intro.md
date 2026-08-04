# Bài 12: RAP — Giới thiệu (RESTful ABAP Programming Model)

## Mục tiêu
- Hiểu kiến trúc RAP: Behavior Definition, Behavior Implementation, Business Object.
- Phân biệt managed vs unmanaged scenario.
- Tạo Business Object CRUD đơn giản bằng managed scenario.

## 1. RAP là gì, và vì sao quan trọng

RAP là mô hình lập trình chuẩn của SAP cho **transactional app** trên nền ABAP (cả on-premise S/4HANA và ABAP Cloud/steampunk trên BTP). RAP thay thế các cách làm cũ (BOPF, hoặc code CRUD thủ công) bằng 1 kiến trúc thống nhất, tự động sinh OData service, hỗ trợ draft, validation, action... với rất ít code viết tay.

So sánh với 2 track kia: RAP giống việc kết hợp **domain model + service layer + auto-generated API** trong 1 công cụ — tương tự cách CAP ([Bài 17](./17_cap_basics.md)) tự sinh OData từ CDS, nhưng RAP chạy trên ABAP thay vì Node.js/Java.

## 2. Kiến trúc RAP — 3 thành phần chính

```
CDS View (Interface View + Projection View)   ← data model, dựa trên Bài 10-11
        │
Behavior Definition (.bdef)                     ← khai báo hành vi: CREATE/UPDATE/DELETE nào được phép,
        │                                          validation nào, action nào
Behavior Implementation (class ABAP)             ← code xử lý logic thật (validation, determination, action)
        │
Service Definition + Service Binding             ← expose ra OData (chi tiết ở Bài 14)
```

## 3. Managed vs Unmanaged scenario

| | Managed | Unmanaged |
|---|---|---|
| CRUD cơ bản | RAP framework tự sinh (Create/Update/Delete built-in) | Bạn tự viết code CRUD hoàn toàn |
| Khi nào dùng | **Mặc định, khuyến khích** cho Business Object mới | Khi logic CRUD quá đặc thù, không map trực tiếp vào 1 bảng |
| Tốc độ phát triển | Nhanh — ít code | Chậm hơn, kiểm soát toàn bộ |

Bài này tập trung **managed scenario** — trường hợp dùng phổ biến nhất.

## 4. CDS View cho RAP — Interface View + Projection View

```abap
" Interface View — data model "gốc", KHÔNG expose trực tiếp ra ngoài
@AbapCatalog.sqlViewName: 'ZVTASKI'
@AccessControl.authorizationCheck: #CHECK
define root view entity ZI_Task
  as select from ztask
{
  key task_id as TaskId,
      owner_id as OwnerId,
      title    as Title,
      done     as Done,
      @Semantics.systemDateTime.lastChangedAt: true
      last_changed_at as LastChangedAt
}
```

```abap
" Projection View — lớp "mặt tiền" expose ra OData, thường 1-1 với Interface View
@AbapCatalog.sqlViewName: 'ZVTASKC'
@AccessControl.authorizationCheck: #CHECK
define root view entity ZC_Task
  provider contract transactional_query
  as projection on ZI_Task
{
  key TaskId,
      OwnerId,
      Title,
      Done,
      LastChangedAt
}
```

## 5. Behavior Definition — khai báo hành vi cho Interface View

```abap
managed implementation in class zbp_i_task unique;
strict(2);

define behavior for ZI_Task alias Task
persistent table ztask
lock master
authorization master ( instance )
etag master LastChangedAt
{
  create;
  update;
  delete;

  field ( readonly ) TaskId;
  field ( mandatory ) Title;

  mapping for ztask
  {
    TaskId = task_id;
    OwnerId = owner_id;
    Title = title;
    Done = done;
    LastChangedAt = last_changed_at;
  }
}
```

- `managed implementation in class ...`: khai báo class xử lý logic bổ sung (validation/determination — [Bài 13](./13_rap_transactional.md)).
- `create; update; delete;`: bật các thao tác CRUD được phép — RAP tự sinh implementation cho các thao tác này.
- `field ( readonly )`: field không được sửa qua API (vd khóa chính tự sinh).
- `field ( mandatory )`: field bắt buộc khi tạo mới.
- `mapping for ztask`: ánh xạ field CDS ↔ field bảng database.

## 6. Behavior Definition cho Projection View (expose)

```abap
projection;
strict(2);

define behavior for ZC_Task alias Task
{
  use create;
  use update;
  use delete;

  use association _Owner;
}
```

## 7. Business Object — Interface View + Behavior Definition + Behavior Implementation ghép lại

Không có 1 từ khóa "Business Object" cụ thể — thuật ngữ này chỉ **toàn bộ tổ hợp** CDS + Behavior Definition + Behavior Implementation cho 1 thực thể nghiệp vụ (`Task` trong ví dụ trên).

## Ví dụ tổng thể — luồng tạo Business Object `Task` (managed, CRUD cơ bản)

1. Tạo bảng database `ZTASK` (SE11 hoặc ADT Data Definition).
2. Tạo Interface View `ZI_Task` (mục 4).
3. Tạo Behavior Definition cho `ZI_Task` (mục 5) — bật `create/update/delete`.
4. Tạo Projection View `ZC_Task` (mục 4).
5. Tạo Behavior Definition cho `ZC_Task` (mục 6) — `use create/update/delete`.
6. Tạo Service Definition + Service Binding (chi tiết [Bài 14](./14_odata_services.md)) để test qua Gateway Client/Postman.

## Bài tập

1. **Bảng + Interface View**: tạo bảng database đơn giản (vd `ZTASK` với `task_id`, `title`, `done`), viết Interface View `ZI_Task`.
2. **Behavior Definition managed**: viết Behavior Definition bật `create/update/delete` như ví dụ.
3. **Projection View**: viết Projection View `ZC_Task` + Behavior Definition tương ứng dùng `use create/update/delete`.
4. **Test qua Preview**: dùng chức năng "Preview" trong ADT (click phải trên Service Binding sau khi tạo — xem trước ở [Bài 14](./14_odata_services.md)) để test CRUD trực tiếp mà chưa cần viết validation/action gì thêm.

## Tiếp theo
→ [Bài 13: RAP — Xây dựng Transactional App](./13_rap_transactional.md)
