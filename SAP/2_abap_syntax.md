# Bài 2: Cú pháp ABAP cơ bản

## Mục tiêu
- Khai báo biến với `DATA`/`TYPES`, dùng kiểu built-in.
- Toán tử, `WRITE`, khai báo inline (`DATA(x) = ...`).

## 1. Khai báo biến — `DATA`

```abap
REPORT z_hello_world.

DATA: lv_name  TYPE string,
      lv_age   TYPE i,
      lv_price TYPE p DECIMALS 2,
      lv_flag  TYPE abap_bool.

lv_name = 'Ben'.
lv_age = 25.
lv_price = '19.99'.
lv_flag = abap_true.

WRITE: / lv_name, lv_age, lv_price.
```

So với Go/Python: `DATA` giống `var` của Go (khai báo tường minh kiểu), nhưng ABAP cổ điển **không có type inference** — luôn phải khai báo `TYPE` tường minh (khác `:=` của Go hay gán trực tiếp của Python).

## 2. Kiểu built-in thường dùng

| Kiểu ABAP | Ý nghĩa | Tương đương Go/Python |
|---|---|---|
| `i` | Số nguyên 4 byte | `int` |
| `string` | Chuỗi ký tự độ dài động | `string`/`str` |
| `c LENGTH n` | Chuỗi ký tự độ dài CỐ ĐỊNH n | không có tương đương trực tiếp — đặc thù ABAP |
| `p DECIMALS n` | Số thập phân đóng gói (packed) — dùng cho tiền tệ | `decimal.Decimal`/tương đương |
| `d` | Ngày (YYYYMMDD) | `time.Time`/`datetime.date` |
| `t` | Giờ (HHMMSS) | `time.Time`/`datetime.time` |
| `abap_bool` | Boolean (`abap_true`/`abap_false`, thực chất là `c LENGTH 1`) | `bool` |

**Lưu ý quan trọng:** `p DECIMALS n` (packed number) LUÔN được dùng cho giá trị tiền tệ/số lượng trong ABAP — KHÔNG dùng `f` (floating point) vì sai số làm tròn có thể gây lỗi nghiêm trọng trong tính toán tài chính.

## 3. Khai báo inline — cú pháp hiện đại (từ ABAP 7.40+)

```abap
" Cách cổ điển — phải khai báo TRƯỚC khi dùng
DATA lv_count TYPE i.
lv_count = 10.

" Cách hiện đại — khai báo NGAY TẠI ĐIỂM DÙNG, kiểu tự suy luận từ ngữ cảnh
DATA(lv_count2) = 10.

" Rất phổ biến trong SELECT hiện đại (xem Bài 7)
SELECT SINGLE * FROM sflight INTO @DATA(ls_flight) WHERE carrid = 'LH'.
```

Cú pháp inline (`DATA(...)`) gần giống `:=` của Go — nhưng khác biệt: ABAP vẫn **suy luận kiểu tĩnh tại compile-time**, không phải dynamic typing như Python.

## 4. Toán tử

```abap
DATA: lv_a TYPE i VALUE 10,
      lv_b TYPE i VALUE 3.

DATA(lv_sum)  = lv_a + lv_b.   " 13
DATA(lv_diff) = lv_a - lv_b.   " 7
DATA(lv_prod) = lv_a * lv_b.   " 30
DATA(lv_div)  = lv_a / lv_b.   " 3.333... (nếu kiểu là i thì làm tròn xuống 3)
DATA(lv_mod)  = lv_a MOD lv_b. " 1

IF lv_a > lv_b AND lv_b > 0.
  WRITE / 'Điều kiện đúng'.
ENDIF.
```

## 5. Xử lý chuỗi cơ bản

```abap
DATA(lv_greeting) = |Xin chào, { lv_name }!|.  " string template — tương đương f-string của Python
WRITE / lv_greeting.

DATA(lv_upper) = to_upper( lv_name ).
DATA(lv_len)   = strlen( lv_name ).
CONCATENATE 'Hello' 'World' INTO DATA(lv_concat) SEPARATED BY space.
```

String template `|...{ }...|` (từ ABAP 7.40+) tương đương f-string của Python — nên dùng thay cho `CONCATENATE` cổ điển khi có thể.

## 6. `WRITE` — output cơ bản (chỉ dùng để debug/report đơn giản)

```abap
WRITE: / 'Hello, World!'.
WRITE: / 'Tên:', lv_name, 'Tuổi:', lv_age.
```

Trong ABAP hiện đại (RAP/CAP), `WRITE` gần như không dùng trong logic nghiệp vụ thật — chỉ hữu ích cho report cổ điển hoặc debug nhanh. Ở [Bài 9](./9_abap_debug_testing.md) bạn sẽ dùng ABAP Debugger thay vì `WRITE` để kiểm tra giá trị biến.

## Ví dụ đầy đủ

```abap
REPORT z_hello_world.

DATA: lv_name TYPE string VALUE 'Ben',
      lv_age  TYPE i      VALUE 25.

DATA(lv_greeting) = |Xin chào, { lv_name }! Bạn { lv_age } tuổi.|.
WRITE / lv_greeting.

DATA(lv_next_year_age) = lv_age + 1.
WRITE: / 'Năm sau bạn', lv_next_year_age, 'tuổi'.
```

## Bài tập

1. **Hello World**: viết report in "Hello, World!" bằng cả `WRITE` và string template `|...|`.
2. **Khai báo biến**: khai báo biến bằng cả `DATA` cổ điển và `DATA(...)` inline, so sánh code, giải thích khi nào nên dùng cách nào.
3. **Toán tử**: viết chương trình tính và in tổng, hiệu, tích, thương, dư của 2 số.
4. **String template**: viết chương trình in ra câu chào có tên và tuổi bằng `|...{ }...|`, so sánh với `CONCATENATE`.

## Tiếp theo
→ [Bài 3: Luồng điều khiển trong ABAP](./3_abap_control_flow.md)
