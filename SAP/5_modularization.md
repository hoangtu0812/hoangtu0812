# Bài 5: Modularization — Function Module & Method

## Mục tiêu
- Phân biệt 3 cách "đóng gói" code trong ABAP: Subroutine (legacy), Function Module, Method.
- Biết khi nào dùng cái nào trong code hiện đại.

## 1. Subroutine (`FORM`/`PERFORM`) — legacy, chỉ nên NHẬN BIẾT khi đọc code cũ

```abap
FORM calculate_total USING iv_price TYPE p iv_qty TYPE i
                      CHANGING cv_total TYPE p.
  cv_total = iv_price * iv_qty.
ENDFORM.

START-OF-SELECTION.
  DATA lv_total TYPE p DECIMALS 2.
  PERFORM calculate_total USING '100.00' 3 CHANGING lv_total.
  WRITE / lv_total.
```

**Không dùng `FORM`/`PERFORM` trong code mới** — SAP khuyến nghị dùng Method (ABAP Objects) từ lâu. Bạn chỉ cần nhận diện được cú pháp này khi đọc code cũ trong hệ thống công ty.

## 2. Function Module — đóng gói logic dùng chung, có thể gọi từ xa (RFC)

```abap
FUNCTION z_calculate_total.
*"----------------------------------------------------------------------
*"*"Local interface:
*"  IMPORTING
*"     VALUE(IV_PRICE) TYPE  P
*"     VALUE(IV_QTY) TYPE  I
*"  EXPORTING
*"     VALUE(EV_TOTAL) TYPE  P
*"  EXCEPTIONS
*"      INVALID_QUANTITY
*"----------------------------------------------------------------------
  IF iv_qty <= 0.
    RAISE invalid_quantity.
  ENDIF.
  ev_total = iv_price * iv_qty.
ENDFUNCTION.
```

Gọi Function Module:

```abap
DATA lv_total TYPE p DECIMALS 2.

CALL FUNCTION 'Z_CALCULATE_TOTAL'
  EXPORTING
    iv_price = '100.00'
    iv_qty   = 3
  IMPORTING
    ev_total = lv_total
  EXCEPTIONS
    invalid_quantity = 1
    OTHERS            = 2.

IF sy-subrc <> 0.
  WRITE / 'Lỗi khi tính tổng'.
ENDIF.
```

Function Module vẫn được dùng rộng rãi trong hệ thống ABAP cổ điển, đặc biệt cho **RFC** (Remote Function Call — gọi từ hệ thống khác) và **BAPI** (Function Module chuẩn hóa để thao tác business object). Bạn sẽ gặp rất nhiều trong code SAP hiện có tại công ty.

## 3. Method (ABAP Objects) — cách hiện đại, khuyến khích dùng cho code mới

```abap
CLASS zcl_calculator DEFINITION.
  PUBLIC SECTION.
    CLASS-METHODS: calculate_total
      IMPORTING iv_price       TYPE p
                iv_qty         TYPE i
      RETURNING VALUE(rv_total) TYPE p
      RAISING   cx_sy_arithmetic_error.
ENDCLASS.

CLASS zcl_calculator IMPLEMENTATION.
  METHOD calculate_total.
    IF iv_qty <= 0.
      RAISE EXCEPTION TYPE cx_sy_arithmetic_error.
    ENDIF.
    rv_total = iv_price * iv_qty.
  ENDMETHOD.
ENDCLASS.
```

Gọi static method:

```abap
DATA(lv_total) = zcl_calculator=>calculate_total( iv_price = '100.00' iv_qty = 3 ).
```

`RETURNING VALUE(...)` cho phép gọi method **inline trong biểu thức** (như hàm bình thường của Go/Python) — khác `EXPORTING`/`IMPORTING` của Function Module luôn cần statement riêng (`CALL FUNCTION`).

## 4. So sánh 3 cách — khi nào dùng gì

| | Subroutine | Function Module | Method |
|---|---|---|---|
| Trạng thái | Legacy, KHÔNG dùng cho code mới | Vẫn dùng cho RFC/BAPI | **Mặc định cho code mới** |
| Gọi từ xa (RFC) | Không | Có | Có (qua wrapper) |
| Return inline trong biểu thức | Không | Không | Có (`RETURNING VALUE`) |
| Exception handling | `sy-subrc` thủ công | `EXCEPTIONS` + `sy-subrc` | Class-based exception (`TRY/CATCH` — xem [Bài 8](./8_abap_exceptions.md)) |
| Tương đương Go/Python | — | tương tự hàm public trong 1 package | tương đương method/function hiện đại |

## Ví dụ đầy đủ: viết lại 1 Function Module bằng Method

```abap
" Function Module tính tổng internal table
FUNCTION z_sum_table.
*"  IMPORTING VALUE(IT_NUMBERS) TYPE  TY_NUMBER_TABLE
*"  EXPORTING VALUE(EV_SUM) TYPE  I
  DATA(lv_sum) = 0.
  LOOP AT it_numbers INTO DATA(lv_num).
    lv_sum = lv_sum + lv_num.
  ENDLOOP.
  ev_sum = lv_sum.
ENDFUNCTION.
```

```abap
" Viết lại bằng static method — gọn hơn, gọi inline được
CLASS zcl_math_utils DEFINITION.
  PUBLIC SECTION.
    TYPES ty_number_table TYPE STANDARD TABLE OF i.
    CLASS-METHODS sum_table
      IMPORTING it_numbers      TYPE ty_number_table
      RETURNING VALUE(rv_sum)   TYPE i.
ENDCLASS.

CLASS zcl_math_utils IMPLEMENTATION.
  METHOD sum_table.
    LOOP AT it_numbers INTO DATA(lv_num).
      rv_sum = rv_sum + lv_num.
    ENDLOOP.
  ENDMETHOD.
ENDCLASS.

" Gọi:
DATA(lt_numbers) = VALUE zcl_math_utils=>ty_number_table( ( 1 ) ( 2 ) ( 3 ) ).
DATA(lv_total) = zcl_math_utils=>sum_table( lt_numbers ).
WRITE / lv_total.  " 6
```

## Bài tập

1. **Function Module**: viết Function Module `Z_CALCULATE_TOTAL` như ví dụ, gọi thử bằng `CALL FUNCTION`.
2. **Viết lại bằng Method**: chuyển Function Module trên thành static method của 1 class, gọi bằng cú pháp inline `RETURNING VALUE`.
3. **Đọc code cũ**: tìm 1 `FORM`/`PERFORM` trong hệ thống công ty (nếu có quyền truy cập), giải thích logic, và viết lại phần đó bằng method (chỉ để luyện tập, không cần deploy).
4. **So sánh error handling**: viết cùng 1 logic (chia 2 số) bằng cả `EXCEPTIONS` của Function Module và class-based exception của Method — so sánh độ rõ ràng.

## Tiếp theo
→ [Bài 6: ABAP Objects — OOP cơ bản](./6_abap_oo.md)
