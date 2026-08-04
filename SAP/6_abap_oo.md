# Bài 6: ABAP Objects — OOP cơ bản

## Mục tiêu
- `CLASS...DEFINITION/IMPLEMENTATION`, `INTERFACE`.
- Kế thừa (`INHERITING FROM`), visibility, constructor.

## 1. Class cơ bản

```abap
CLASS zcl_rectangle DEFINITION.
  PUBLIC SECTION.
    METHODS: constructor
      IMPORTING iv_width  TYPE p
                iv_height TYPE p,
      get_area
        RETURNING VALUE(rv_area) TYPE p,
      get_perimeter
        RETURNING VALUE(rv_perimeter) TYPE p.

  PRIVATE SECTION.
    DATA: mv_width  TYPE p,
          mv_height TYPE p.
ENDCLASS.

CLASS zcl_rectangle IMPLEMENTATION.
  METHOD constructor.
    mv_width = iv_width.
    mv_height = iv_height.
  ENDMETHOD.

  METHOD get_area.
    rv_area = mv_width * mv_height.
  ENDMETHOD.

  METHOD get_perimeter.
    rv_perimeter = 2 * ( mv_width + mv_height ).
  ENDMETHOD.
ENDCLASS.
```

Gọi:

```abap
DATA(lo_rect) = NEW zcl_rectangle( iv_width = 10 iv_height = 5 ).
WRITE: / lo_rect->get_area( ), lo_rect->get_perimeter( ).
```

Quy ước đặt tên phổ biến: `mv_` cho instance attribute (member variable), `iv_`/`ev_`/`cv_`/`rv_` cho tham số IMPORTING/EXPORTING/CHANGING/RETURNING, `lo_` cho local object reference, `lt_`/`ls_` cho local table/structure.

## 2. Visibility — `PUBLIC/PROTECTED/PRIVATE SECTION`

```
PUBLIC SECTION     — truy cập được từ MỌI nơi (tương đương exported của Go, public của Python convention)
PROTECTED SECTION  — chỉ class này và class con truy cập được
PRIVATE SECTION    — chỉ class này truy cập được
```

Tương đương `PUBLIC`/`PRIVATE`/`_protected` (convention) của Python ([Python Bài 7](../Python/7_oop.md)), nhưng ABAP **enforce thật sự ở compile-time** — khác Python chỉ là quy ước tên biến.

## 3. Interface

```abap
INTERFACE zif_shape.
  METHODS get_area RETURNING VALUE(rv_area) TYPE p.
ENDINTERFACE.

CLASS zcl_rectangle DEFINITION.
  PUBLIC SECTION.
    INTERFACES zif_shape.  " class CAM KẾT implement mọi method của interface
    METHODS constructor IMPORTING iv_width TYPE p iv_height TYPE p.
  PRIVATE SECTION.
    DATA: mv_width TYPE p, mv_height TYPE p.
ENDCLASS.

CLASS zcl_rectangle IMPLEMENTATION.
  METHOD constructor.
    mv_width = iv_width.
    mv_height = iv_height.
  ENDMETHOD.

  METHOD zif_shape~get_area.  " tên method phải có tiền tố interface~
    rv_area = mv_width * mv_height.
  ENDMETHOD.
ENDCLASS.
```

**Khác Go/Python:** ABAP `INTERFACE` giống Go interface ở chỗ định nghĩa hợp đồng method, nhưng **phải khai báo tường minh** `INTERFACES zif_shape` trong class (giống Java) — không phải structural typing ngầm định như Go's interface ([Go Bài 10](../Go/10_interfaces.md)) hay Python's `Protocol` ([Python Bài 13](../Python/13_type_hints.md)).

## 4. Kế thừa — `INHERITING FROM`

```abap
CLASS zcl_shape DEFINITION ABSTRACT.
  PUBLIC SECTION.
    METHODS get_area ABSTRACT RETURNING VALUE(rv_area) TYPE p.
ENDCLASS.

CLASS zcl_circle DEFINITION INHERITING FROM zcl_shape.
  PUBLIC SECTION.
    METHODS: constructor IMPORTING iv_radius TYPE p,
             get_area REDEFINITION.
  PRIVATE SECTION.
    DATA mv_radius TYPE p.
ENDCLASS.

CLASS zcl_circle IMPLEMENTATION.
  METHOD constructor.
    mv_radius = iv_radius.
  ENDMETHOD.

  METHOD get_area.
    rv_area = '3.14159' * mv_radius * mv_radius.
  ENDMETHOD.
ENDCLASS.
```

`ABSTRACT` (class hoặc method) tương đương định nghĩa interface tối thiểu bắt buộc subclass phải implement — giống base class trừu tượng của Python ([Python Bài 7](../Python/7_oop.md)).

## 5. Đa hình

```abap
DATA: lo_shape TYPE REF TO zcl_shape.
lo_shape = NEW zcl_circle( iv_radius = 5 ).
WRITE / lo_shape->get_area( ).  " gọi đúng implementation của Circle nhờ REDEFINITION
```

## Ví dụ đầy đủ

```abap
CLASS zcl_shape DEFINITION ABSTRACT.
  PUBLIC SECTION.
    METHODS get_area ABSTRACT RETURNING VALUE(rv_area) TYPE p.
ENDCLASS.

CLASS zcl_rectangle DEFINITION INHERITING FROM zcl_shape.
  PUBLIC SECTION.
    METHODS: constructor IMPORTING iv_width TYPE p iv_height TYPE p,
             get_area REDEFINITION.
  PRIVATE SECTION.
    DATA: mv_width TYPE p, mv_height TYPE p.
ENDCLASS.

CLASS zcl_rectangle IMPLEMENTATION.
  METHOD constructor.
    mv_width = iv_width.
    mv_height = iv_height.
  ENDMETHOD.
  METHOD get_area.
    rv_area = mv_width * mv_height.
  ENDMETHOD.
ENDCLASS.

CLASS zcl_circle DEFINITION INHERITING FROM zcl_shape.
  PUBLIC SECTION.
    METHODS: constructor IMPORTING iv_radius TYPE p,
             get_area REDEFINITION.
  PRIVATE SECTION.
    DATA mv_radius TYPE p.
ENDCLASS.

CLASS zcl_circle IMPLEMENTATION.
  METHOD constructor.
    mv_radius = iv_radius.
  ENDMETHOD.
  METHOD get_area.
    rv_area = '3.14159' * mv_radius * mv_radius.
  ENDMETHOD.
ENDCLASS.

START-OF-SELECTION.
  DATA(lt_shapes) = VALUE zcl_shape_table( ( NEW zcl_rectangle( iv_width = 10 iv_height = 5 ) )
                                            ( NEW zcl_circle( iv_radius = 3 ) ) ).
  LOOP AT lt_shapes INTO DATA(lo_shape).
    WRITE / lo_shape->get_area( ).
  ENDLOOP.
```

## Bài tập

1. **`ZCL_RECTANGLE`**: viết class như ví dụ đầu bài, với `get_area`/`get_perimeter`.
2. **`ZIF_SHAPE` + implement**: viết interface `ZIF_SHAPE`, cho `ZCL_RECTANGLE` và `ZCL_CIRCLE` cùng implement.
3. **Kế thừa `ZCL_SHAPE` abstract**: viết class abstract `ZCL_SHAPE` với method abstract `get_area`, cho `ZCL_RECTANGLE`, `ZCL_CIRCLE` kế thừa và `REDEFINITION`.
4. **Nâng cao**: viết `constructor` cho `ZCL_RECTANGLE` raise exception (xem trước ở [Bài 8](./8_abap_exceptions.md)) nếu `iv_width`/`iv_height` <= 0.

## Tiếp theo
→ [Bài 7: Open SQL — Truy vấn Database](./7_open_sql.md)
