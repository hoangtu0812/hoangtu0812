# Bài 8: Xử lý lỗi trong ABAP

## Mục tiêu
- Class-based exception: `TRY/CATCH/CLEANUP`.
- Tự định nghĩa exception class.
- So sánh với error handling của Go/Python.

## 1. `TRY / CATCH` cơ bản

```abap
TRY.
    DATA(lv_result) = 10 / 0.
  CATCH cx_sy_zerodivide INTO DATA(lx_error).
    WRITE: / 'Lỗi:', lx_error->get_text( ).
ENDTRY.
```

`CATCH ... INTO` bắt exception object, `get_text( )` lấy thông báo lỗi — tương đương bắt exception bằng `except Exception as e:` của Python ([Python Bài 8](../Python/8_exceptions.md)).

## 2. Bắt nhiều loại exception

```abap
TRY.
    " logic có thể raise nhiều loại lỗi khác nhau
    DATA(lv_result) = 10 / 0.
  CATCH cx_sy_zerodivide INTO DATA(lx_zero).
    WRITE / 'Chia cho 0'.
  CATCH cx_sy_conversion_error INTO DATA(lx_conv).
    WRITE / 'Lỗi convert dữ liệu'.
  CATCH cx_root INTO DATA(lx_root).  " bắt MỌI exception khác — đặt SAU CÙNG, giống except Exception: của Python
    WRITE: / 'Lỗi không xác định:', lx_root->get_text( ).
ENDTRY.
```

`cx_root` là base class của toàn bộ exception hierarchy trong ABAP — tương đương `Exception` của Python hoặc `error` interface của Go.

## 3. Tự định nghĩa exception class

```abap
CLASS zcx_validation_error DEFINITION INHERITING FROM cx_static_check.
  PUBLIC SECTION.
    DATA: mv_field   TYPE string,
          mv_message TYPE string.

    METHODS: constructor
      IMPORTING iv_field   TYPE string
                iv_message TYPE string
                previous   LIKE previous OPTIONAL.
ENDCLASS.

CLASS zcx_validation_error IMPLEMENTATION.
  METHOD constructor.
    super->constructor( previous = previous ).
    mv_field = iv_field.
    mv_message = iv_message.
  ENDMETHOD.
ENDCLASS.
```

Dùng:

```abap
CLASS zcl_validator DEFINITION.
  PUBLIC SECTION.
    CLASS-METHODS validate_age
      IMPORTING iv_age TYPE i
      RAISING   zcx_validation_error.
ENDCLASS.

CLASS zcl_validator IMPLEMENTATION.
  METHOD validate_age.
    IF iv_age < 0 OR iv_age > 150.
      RAISE EXCEPTION TYPE zcx_validation_error
        EXPORTING
          iv_field   = 'age'
          iv_message = 'phải nằm trong khoảng 0-150'.
    ENDIF.
  ENDMETHOD.
ENDCLASS.

TRY.
    zcl_validator=>validate_age( -5 ).
  CATCH zcx_validation_error INTO DATA(lx_val).
    WRITE: / 'Lỗi field', lx_val->mv_field, ':', lx_val->mv_message.
ENDTRY.
```

Tương đương custom error `ValidationError` của Go ([Go Bài 8](../Go/8_error_handling.md)) hoặc Python ([Python Bài 8](../Python/8_exceptions.md)) — nhưng ABAP yêu cầu định nghĩa **class riêng kế thừa từ `cx_static_check`/`cx_dynamic_check`/`cx_no_check`**.

## 4. `CX_STATIC_CHECK` vs `CX_DYNAMIC_CHECK` vs `CX_NO_CHECK`

| Base class | Ý nghĩa | Tương đương |
|---|---|---|
| `CX_STATIC_CHECK` | Compiler BẮT BUỘC method gọi phải khai báo `RAISING` hoặc `TRY/CATCH` | `error` phải kiểm tra tường minh (Go), checked exception |
| `CX_DYNAMIC_CHECK` | Kiểm tra lúc runtime, không bắt buộc lúc compile | ít dùng hơn |
| `CX_NO_CHECK` | Không bắt buộc kiểm tra gì cả | tương tự unchecked/runtime exception |

**Khuyến khích dùng `CX_STATIC_CHECK`** cho exception nghiệp vụ — buộc lập trình viên gọi hàm phải xử lý lỗi tường minh, giống triết lý "luôn kiểm tra `err != nil`" của Go.

## 5. `CLEANUP` — tương đương `finally` của Python / `defer` của Go

```abap
TRY.
    " mở resource, xử lý có thể lỗi
    DATA(lv_result) = 10 / 0.
  CATCH cx_sy_zerodivide.
    WRITE / 'Xử lý lỗi'.
  CLEANUP.
    " LUÔN chạy khi có exception xảy ra và KHÔNG bị catch ở khối CATCH tương ứng
    " Thường dùng để rollback transaction dở dang
    WRITE / 'Dọn dẹp'.
ENDTRY.
```

## Ví dụ đầy đủ

```abap
CLASS zcx_validation_error DEFINITION INHERITING FROM cx_static_check.
  PUBLIC SECTION.
    DATA: mv_field TYPE string, mv_message TYPE string.
    METHODS constructor IMPORTING iv_field TYPE string iv_message TYPE string.
ENDCLASS.

CLASS zcx_validation_error IMPLEMENTATION.
  METHOD constructor.
    super->constructor( ).
    mv_field = iv_field.
    mv_message = iv_message.
  ENDMETHOD.
ENDCLASS.

CLASS zcl_order_validator DEFINITION.
  PUBLIC SECTION.
    CLASS-METHODS validate_quantity
      IMPORTING iv_qty TYPE i
      RAISING   zcx_validation_error.
ENDCLASS.

CLASS zcl_order_validator IMPLEMENTATION.
  METHOD validate_quantity.
    IF iv_qty <= 0.
      RAISE EXCEPTION TYPE zcx_validation_error
        EXPORTING iv_field = 'quantity' iv_message = 'phải lớn hơn 0'.
    ENDIF.
  ENDMETHOD.
ENDCLASS.

START-OF-SELECTION.
  TRY.
      zcl_order_validator=>validate_quantity( -5 ).
    CATCH zcx_validation_error INTO DATA(lx_error).
      WRITE: / 'Lỗi validate:', lx_error->mv_field, '-', lx_error->mv_message.
  ENDTRY.
```

## Bài tập

1. **`TRY/CATCH` cơ bản**: viết code cố tình chia cho 0, bắt `cx_sy_zerodivide`, in thông báo lỗi.
2. **Custom exception**: viết `ZCX_VALIDATION_ERROR` như ví dụ, viết method `validate_age` raise lỗi khi không hợp lệ.
3. **Nhiều `CATCH`**: viết code có thể phát sinh 2 loại lỗi khác nhau, bắt riêng từng loại + 1 `CATCH cx_root` cuối cùng.
4. **So sánh với Go/Python**: viết ghi chú so sánh `RAISING` (khai báo bắt buộc exception có thể raise) với `error` return của Go và exception ngầm định của Python — điều gì ABAP làm chặt chẽ hơn, điều gì Go làm rõ ràng hơn.

## Tiếp theo
→ [Bài 9: Debugging & ABAP Unit Testing](./9_abap_debug_testing.md)
