# Bài 3: Luồng điều khiển trong ABAP

## Mục tiêu
- `IF/ELSEIF/ELSE`, `CASE`, `DO`, `WHILE`, `LOOP AT`.
- `CHECK`, `CONTINUE`, `EXIT` — kiểm soát luồng trong vòng lặp.

## 1. `IF / ELSEIF / ELSE`

```abap
DATA(lv_score) = 75.

IF lv_score >= 90.
  WRITE / 'Xuất sắc'.
ELSEIF lv_score >= 70.
  WRITE / 'Khá'.
ELSE.
  WRITE / 'Cần cố gắng'.
ENDIF.
```

Mọi khối lệnh ABAP kết thúc bằng từ khóa `END...` tương ứng (`ENDIF`, `ENDDO`, `ENDLOOP`...) — khác Go/Python dùng `{}`/thụt lề.

## 2. `CASE` — tương đương `switch` của Go

```abap
DATA(lv_day) = 3.

CASE lv_day.
  WHEN 1 OR 7.
    WRITE / 'Cuối tuần'.
  WHEN 2 OR 3 OR 4 OR 5 OR 6.
    WRITE / 'Ngày thường'.
  WHEN OTHERS.
    WRITE / 'Không hợp lệ'.
ENDCASE.
```

## 3. `DO` — vòng lặp đếm số lần cố định

```abap
DO 5 TIMES.
  WRITE / sy-index.  " sy-index: biến hệ thống, tự tăng 1, 2, 3... trong mỗi vòng DO
ENDDO.

" DO vô hạn — phải tự EXIT
DO.
  IF sy-index > 3.
    EXIT.
  ENDIF.
  WRITE / sy-index.
ENDDO.
```

## 4. `WHILE`

```abap
DATA(lv_n) = 0.
WHILE lv_n < 5.
  WRITE / lv_n.
  lv_n = lv_n + 1.
ENDWHILE.
```

## 5. `LOOP AT` — duyệt internal table (chi tiết đầy đủ ở [Bài 4](./4_internal_tables.md))

```abap
DATA: lt_numbers TYPE TABLE OF i.
lt_numbers = VALUE #( ( 10 ) ( 20 ) ( 30 ) ).

LOOP AT lt_numbers INTO DATA(lv_number).
  WRITE / lv_number.
ENDLOOP.

" Có index
LOOP AT lt_numbers INTO DATA(lv_num2) FROM 1.
  WRITE: / sy-tabix, lv_num2.  " sy-tabix: index hiện tại trong LOOP AT
ENDLOOP.
```

## 6. `CHECK`, `CONTINUE`, `EXIT` trong vòng lặp

```abap
DO 10 TIMES.
  IF sy-index MOD 2 = 0.
    CONTINUE.  " bỏ qua vòng hiện tại, sang vòng tiếp theo — giống continue của Go
  ENDIF.
  IF sy-index > 7.
    EXIT.       " thoát hẳn vòng lặp — giống break của Go
  ENDIF.
  WRITE / sy-index.
ENDDO.
```

`CHECK` là cách viết cũ hơn, tương đương "nếu điều kiện sai thì `CONTINUE`":

```abap
LOOP AT lt_numbers INTO DATA(lv_val).
  CHECK lv_val > 15.  " nếu lv_val <= 15, tự động nhảy sang vòng tiếp theo (như CONTINUE)
  WRITE / lv_val.
ENDLOOP.
```

**Khuyến khích dùng `IF ... CONTINUE` tường minh thay vì `CHECK`** trong code hiện đại — `CHECK` dễ gây khó đọc khi đặt giữa 1 khối logic dài.

## Ví dụ đầy đủ: FizzBuzz

```abap
REPORT z_fizzbuzz.

DO 100 TIMES.
  DATA(lv_i) = sy-index.

  IF lv_i MOD 15 = 0.
    WRITE / 'FizzBuzz'.
  ELSEIF lv_i MOD 3 = 0.
    WRITE / 'Fizz'.
  ELSEIF lv_i MOD 5 = 0.
    WRITE / 'Buzz'.
  ELSE.
    WRITE / lv_i.
  ENDIF.
ENDDO.
```

## Bài tập

1. **FizzBuzz 1-100**: dùng code mẫu trên làm nền, tự gõ lại bằng `DO...TIMES`.
2. **Số nguyên tố**: viết chương trình kiểm tra và in các số nguyên tố từ 2-100 (dùng `DO`/`WHILE` + `IF`).
3. **`LOOP AT` với `CONTINUE`/`EXIT`**: tạo 1 internal table số, dùng `LOOP AT` in ra chỉ các số > 15, dừng hẳn nếu gặp số > 50.
4. **`CASE` phân loại**: viết chương trình nhận điểm số (0-100), dùng `CASE` (không phải `IF`) để phân loại A/B/C/D/F theo khoảng điểm (gợi ý: dùng `CASE TYPE OF` hoặc chia nhỏ điều kiện bằng biến trung gian vì `CASE` cổ điển của ABAP so sánh giá trị chính xác, không so sánh khoảng).

## Tiếp theo
→ [Bài 4: Internal Table (Bảng nội bộ)](./4_internal_tables.md)
