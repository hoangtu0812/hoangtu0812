# Bài 9: Debugging & ABAP Unit Testing

## Mục tiêu
- Dùng ABAP Debugger (breakpoint, watchpoint) thay vì `WRITE` để kiểm tra giá trị.
- Viết ABAP Unit test (`CLASS ... FOR TESTING`), tương đương `go test`/`pytest`.

## 1. ABAP Debugger — breakpoint

```abap
DATA(lv_total) = 0.
LOOP AT lt_numbers INTO DATA(lv_num).
  BREAK-POINT.  " chương trình dừng TẠI ĐÂY khi chạy trong môi trường dev (KHÔNG dùng trong code production)
  lv_total = lv_total + lv_num.
ENDLOOP.
```

Cách tốt hơn `BREAK-POINT` cứng trong code: đặt **breakpoint từ ADT** (double-click lề trái dòng code) — không cần sửa code, không rủi ro quên xóa `BREAK-POINT` trước khi lên production.

Trong ADT Debugger:
- **F5**: step into (nhảy vào bên trong method được gọi).
- **F6**: step over (chạy qua dòng hiện tại, không nhảy vào method con).
- **F7**: step return (chạy tới hết method hiện tại rồi dừng).
- **F8**: continue (chạy tiếp tới breakpoint kế tiếp).
- **Watchpoint**: dừng chương trình khi 1 biến cụ thể **thay đổi giá trị** — hữu ích khi không biết chính xác đoạn code nào làm sai giá trị.

## 2. Debugger variable view — thay thế hoàn toàn `WRITE` để debug

Thay vì:
```abap
WRITE: / 'lv_total hiện tại là:', lv_total.  " phải sửa code, chạy lại, xóa lại sau khi xong
```

Dùng ADT Debugger: đặt breakpoint, chạy chương trình, xem trực tiếp giá trị mọi biến trong panel "Variables" — không cần sửa code, có thể **sửa giá trị biến ngay lúc debug** để thử các nhánh logic khác nhau.

## 3. ABAP Unit — viết test tự động

```abap
CLASS zcl_calculator DEFINITION.
  PUBLIC SECTION.
    CLASS-METHODS divide
      IMPORTING iv_a           TYPE p
                iv_b           TYPE p
      RETURNING VALUE(rv_result) TYPE p
      RAISING   cx_sy_zerodivide.
ENDCLASS.

CLASS zcl_calculator IMPLEMENTATION.
  METHOD divide.
    IF iv_b = 0.
      RAISE EXCEPTION TYPE cx_sy_zerodivide.
    ENDIF.
    rv_result = iv_a / iv_b.
  ENDMETHOD.
ENDCLASS.
```

```abap
CLASS ltc_calculator DEFINITION FOR TESTING
  DURATION SHORT
  RISK LEVEL HARMLESS.
  PRIVATE SECTION.
    METHODS: test_divide_normal FOR TESTING,
             test_divide_by_zero FOR TESTING.
ENDCLASS.

CLASS ltc_calculator IMPLEMENTATION.
  METHOD test_divide_normal.
    DATA(lv_result) = zcl_calculator=>divide( iv_a = 10 iv_b = 2 ).
    cl_abap_unit_assert=>assert_equals(
      act = lv_result
      exp = 5
      msg = 'Kết quả chia phải bằng 5' ).
  ENDMETHOD.

  METHOD test_divide_by_zero.
    TRY.
        zcl_calculator=>divide( iv_a = 10 iv_b = 0 ).
        cl_abap_unit_assert=>fail( 'Phải raise cx_sy_zerodivide' ).
      CATCH cx_sy_zerodivide.
        " test pass — đúng như mong đợi
    ENDTRY.
  ENDMETHOD.
ENDCLASS.
```

Chạy: click phải trong ADT → **Run As → ABAP Unit Test**, hoặc `Ctrl+Shift+F10`.

## 4. `LTC_...` — local test class convention

`LTC_` (Local Test Class) là tiền tố quy ước cho class test, đặt trong cùng include hoặc `.testclasses.abap` riêng — tương đương file `_test.go` của Go hoặc `test_*.py` của Python ([Go Bài 14](../Go/14_testing.md), [Python Bài 14](../Python/14_testing.md)).

`cl_abap_unit_assert=>` cung cấp các method assertion:

```abap
cl_abap_unit_assert=>assert_equals( act = ... exp = ... ).
cl_abap_unit_assert=>assert_true( act = ... ).
cl_abap_unit_assert=>assert_false( act = ... ).
cl_abap_unit_assert=>assert_bound( act = ... ).      " kiểm tra object reference không rỗng
cl_abap_unit_assert=>fail( msg = '...' ).             " chủ động fail test
```

## 5. Test double đơn giản — dependency injection thủ công trong ABAP

```abap
INTERFACE zif_user_repository.
  METHODS find_by_id IMPORTING iv_id TYPE i RETURNING VALUE(rs_user) TYPE ty_user.
ENDINTERFACE.

CLASS zcl_user_service DEFINITION.
  PUBLIC SECTION.
    METHODS: constructor IMPORTING io_repo TYPE REF TO zif_user_repository,
             get_user_name IMPORTING iv_id TYPE i RETURNING VALUE(rv_name) TYPE string.
  PRIVATE SECTION.
    DATA mo_repo TYPE REF TO zif_user_repository.
ENDCLASS.

CLASS zcl_user_service IMPLEMENTATION.
  METHOD constructor.
    mo_repo = io_repo.
  ENDMETHOD.
  METHOD get_user_name.
    DATA(ls_user) = mo_repo->find_by_id( iv_id ).
    rv_name = ls_user-name.
  ENDMETHOD.
ENDCLASS.

" Test double implement CÙNG interface, không cần database thật
CLASS ltd_fake_user_repo DEFINITION FOR TESTING.
  PUBLIC SECTION.
    INTERFACES zif_user_repository.
ENDCLASS.

CLASS ltd_fake_user_repo IMPLEMENTATION.
  METHOD zif_user_repository~find_by_id.
    rs_user-id = 1.
    rs_user-name = 'Ben (fake)'.
  ENDMETHOD.
ENDCLASS.
```

Đây chính là ý tưởng mock repository ở [Go Bài 14](../Go/14_testing.md) và [Python Bài 14](../Python/14_testing.md) — nhờ `INTERFACE`, service phụ thuộc vào hợp đồng chứ không phụ thuộc implementation thật, cho phép test không cần kết nối database.

## Bài tập

1. **Breakpoint trong ADT**: đặt breakpoint trong 1 method bất kỳ (không dùng `BREAK-POINT` cứng), chạy debug, dùng F5/F6/F8 để quan sát luồng thực thi và giá trị biến.
2. **ABAP Unit cho `ZCL_CALCULATOR`**: viết `LTC_CALCULATOR` như ví dụ, test cả case thành công và case raise exception.
3. **Watchpoint**: đặt watchpoint theo dõi 1 biến trong vòng lặp, chạy debug để xem chương trình tự dừng khi giá trị biến đó thay đổi.
4. **Test double**: viết interface `ZIF_USER_REPOSITORY` + fake implementation cho test, viết `ZCL_USER_SERVICE` phụ thuộc interface, test bằng fake repo.

## Tổng kết Giai đoạn 1
Bạn đã ôn lại/hệ thống hóa nền tảng ABAP cổ điển: cú pháp, control flow, internal table, modularization, OOP, Open SQL, exception, debugging/testing. Giai đoạn 2 sẽ đi vào phần ABAP hiện đại dùng nhiều trong dự án thực tế: CDS View và RAP.

## Tiếp theo
→ [Bài 10: CDS View cơ bản](./10_cds_basics.md)
