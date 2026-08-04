# Bài 13: RAP — Xây dựng Transactional App

## Mục tiêu
- Draft handling.
- Validation (`VALIDATION ... FOR`), Determination (`DETERMINATION ... FOR`).
- Action tự định nghĩa.

## 1. Draft handling — "nháp" trước khi save thật sự

Draft cho phép user chỉnh sửa dữ liệu qua nhiều bước (nhiều màn hình Fiori Elements) mà **chưa commit** vào database chính thức — giống trạng thái "chưa submit form" trong web app thông thường, nhưng SAP chuẩn hóa nó thành 1 cơ chế framework hỗ trợ sẵn.

```abap
managed implementation in class zbp_i_task unique;
strict(2);
with draft;   // bật draft handling cho Business Object này

define behavior for ZI_Task alias Task
persistent table ztask
draft table ztask_d      // bảng draft riêng, RAP tự tạo/quản lý
lock master
authorization master ( instance )
etag master LastChangedAt
{
  create;
  update;
  delete;

  draft action Activate;    // chuyển từ draft -> dữ liệu chính thức
  draft action Discard;     // hủy bỏ draft
  draft action Edit;         // bắt đầu chỉnh sửa (tạo draft từ dữ liệu đã active)
  draft action Resume;
  draft determine action Prepare;

  field ( mandatory ) Title;
}
```

Bảng draft (`ztask_d`) cần tạo riêng, cấu trúc gần giống bảng chính nhưng thêm field kỹ thuật quản lý draft (`%pid`, `%lock`...) — thường sinh tự động qua template trong ADT khi bật `with draft`.

## 2. Validation — kiểm tra dữ liệu hợp lệ trước khi lưu

```abap
managed implementation in class zbp_i_task unique;
strict(2);

define behavior for ZI_Task alias Task
persistent table ztask
lock master
authorization master ( instance )
{
  create;
  update;
  delete;

  field ( mandatory ) Title;

  validation validateTitle on save { create; update; }
}
```

```abap
CLASS lhc_task DEFINITION INHERITING FROM cl_abap_behavior_handler.
  PRIVATE SECTION.
    METHODS validateTitle FOR VALIDATE ON SAVE
      IMPORTING keys FOR Task~validateTitle.
ENDCLASS.

CLASS lhc_task IMPLEMENTATION.
  METHOD validateTitle.
    READ ENTITIES OF ZI_Task IN LOCAL MODE
      ENTITY Task
      FIELDS ( Title )
      WITH CORRESPONDING #( keys )
      RESULT DATA(lt_tasks).

    LOOP AT lt_tasks INTO DATA(ls_task).
      IF ls_task-Title IS INITIAL.
        APPEND VALUE #( %tky = ls_task-%tky ) TO failed-task.
        APPEND VALUE #(
          %tky = ls_task-%tky
          %msg = new_message( id       = 'ZTASK_MSG'
                               number   = '001'
                               severity = if_abap_behv_message=>severity-error )
        ) TO reported-task.
      ENDIF.
    ENDLOOP.
  ENDMETHOD.
ENDCLASS.
```

Validation chạy **trước khi save**, có thể thêm entry vào `failed`/`reported` để ngăn không cho lưu và hiển thị lỗi rõ ràng cho user trên Fiori UI — tương đương validation ở tầng service của Go/Python ([Go Bài 19](../Go/19_auth.md), [Python Bài 19](../Python/19_auth.md)) nhưng tích hợp sẵn cơ chế báo lỗi UI.

## 3. Determination — tự động tính/gán giá trị field

```abap
define behavior for ZI_Task alias Task
persistent table ztask
{
  create;
  update;
  delete;

  determination setDefaultStatus on modify { create; }
}
```

```abap
CLASS lhc_task DEFINITION INHERITING FROM cl_abap_behavior_handler.
  PRIVATE SECTION.
    METHODS setDefaultStatus FOR DETERMINE ON MODIFY
      IMPORTING keys FOR Task~setDefaultStatus.
ENDCLASS.

CLASS lhc_task IMPLEMENTATION.
  METHOD setDefaultStatus.
    MODIFY ENTITIES OF ZI_Task IN LOCAL MODE
      ENTITY Task
      UPDATE FIELDS ( Done )
      WITH VALUE #( FOR key IN keys ( %tky = key-%tky Done = abap_false ) ).
  ENDMETHOD.
ENDCLASS.
```

Determination tương đương `before CREATE` handler trong CAP ([Bài 18](./18_cap_advanced.md)) hoặc gán default value trong constructor của Go/Python.

## 4. Action tùy chỉnh — nghiệp vụ ngoài CRUD chuẩn

```abap
define behavior for ZI_Task alias Task
persistent table ztask
{
  create;
  update;
  delete;

  action markComplete result [1] $self;
}
```

```abap
CLASS lhc_task DEFINITION INHERITING FROM cl_abap_behavior_handler.
  PRIVATE SECTION.
    METHODS markComplete FOR MODIFY
      IMPORTING keys FOR ACTION Task~markComplete RESULT result.
ENDCLASS.

CLASS lhc_task IMPLEMENTATION.
  METHOD markComplete.
    MODIFY ENTITIES OF ZI_Task IN LOCAL MODE
      ENTITY Task
      UPDATE FIELDS ( Done )
      WITH VALUE #( FOR key IN keys ( %tky = key-%tky Done = abap_true ) )
      FAILED failed
      REPORTED reported.

    READ ENTITIES OF ZI_Task IN LOCAL MODE
      ENTITY Task
      ALL FIELDS WITH CORRESPONDING #( keys )
      RESULT DATA(lt_result).

    result = VALUE #( FOR ls_res IN lt_result ( %tky = ls_res-%tky %param = ls_res ) ).
  ENDMETHOD.
ENDCLASS.
```

Action là cách RAP hỗ trợ nghiệp vụ **không chỉ là CRUD đơn thuần** — tương đương endpoint tùy chỉnh như `POST /tasks/{id}/complete` mà bạn tự viết route riêng trong Go/Python ([Go Bài 18](../Go/18_rest_api.md)) thay vì dùng `PUT` chuẩn.

## Bài tập

1. **Validation**: thêm validation cho Business Object `Task` (Bài 12) — không cho save nếu `Title` rỗng, hiển thị message lỗi rõ ràng.
2. **Determination**: thêm determination tự động gán `Done = false` khi tạo task mới.
3. **Action tùy chỉnh**: thêm action `markComplete` như ví dụ, test qua Preview trong ADT.
4. **Draft (nâng cao)**: bật `with draft` cho Business Object, quan sát sự khác biệt trải nghiệm khi test qua Fiori Elements Preview (dữ liệu chỉ lưu thật khi "Save", có thể "Cancel" giữa chừng).

## Tiếp theo
→ [Bài 14: OData Services](./14_odata_services.md)
