# Bài 7: Open SQL — Truy vấn Database

## Mục tiêu
- Viết `SELECT` hiện đại với inline declaration.
- JOIN, `WHERE`, tránh anti-pattern `SELECT *` trong vòng lặp.
- Liên hệ CDS View dùng trong Open SQL.

## 1. `SELECT` cơ bản — cú pháp hiện đại (inline declaration)

```abap
" Lấy 1 dòng
SELECT SINGLE carrid, connid, fldate
  FROM sflight
  INTO @DATA(ls_flight)
  WHERE carrid = 'LH' AND connid = '0400'.

IF sy-subrc = 0.
  WRITE: / ls_flight-carrid, ls_flight-connid.
ENDIF.

" Lấy nhiều dòng vào internal table
SELECT carrid, connid, fldate, price
  FROM sflight
  INTO TABLE @DATA(lt_flights)
  WHERE carrid = 'LH'.
```

So với `database/sql` của Go ([Go Bài 17](../Go/17_database.md)) hoặc SQLAlchemy của Python ([Python Bài 17](../Python/17_database.md)): Open SQL là **DSL tích hợp sẵn trong ngôn ngữ ABAP** — không cần driver riêng, câu SQL được kiểm tra cú pháp ngay lúc biên dịch (syntax check), và tự động tối ưu cho database bên dưới (HANA hoặc DB khác).

## 2. So sánh cú pháp cũ vs mới

```abap
" Cách CŨ (trước ABAP 7.40) — vẫn gặp nhiều trong code cũ tại công ty
DATA: lt_flights TYPE TABLE OF sflight,
      ls_flight  TYPE sflight.

SELECT * FROM sflight INTO TABLE lt_flights WHERE carrid = 'LH'.

" Cách MỚI (khuyến khích) — inline declaration + chỉ định field cần thiết
SELECT carrid, connid, fldate, price
  FROM sflight
  INTO TABLE @DATA(lt_flights_new)
  WHERE carrid = 'LH'.
```

**Luôn liệt kê field cần thiết thay vì `SELECT *`** — giảm băng thông, tránh lấy dư dữ liệu không dùng, và làm rõ code cần field gì (nguyên tắc giống việc chỉ export cái cần thiết ở [Go Bài 9](../Go/9_packages_modules.md)).

## 3. JOIN

```abap
SELECT sflight~carrid, sflight~connid, sflight~fldate,
       scarr~carrname
  FROM sflight
  INNER JOIN scarr ON sflight~carrid = scarr~carrid
  INTO TABLE @DATA(lt_flights_with_carrier)
  WHERE sflight~carrid = 'LH'.

LOOP AT lt_flights_with_carrier INTO DATA(ls_row).
  WRITE: / ls_row-carrname, ls_row-connid.
ENDLOOP.
```

`LEFT OUTER JOIN` dùng khi muốn giữ lại dòng bên trái dù không có dòng khớp bên phải — cú pháp tương tự SQL chuẩn.

## 4. Anti-pattern nghiêm trọng: `SELECT` bên trong `LOOP`

```abap
" SAI — N+1 query problem, mỗi lần lặp lại query database 1 lần — CỰC KỲ CHẬM với bảng lớn
LOOP AT lt_flights INTO DATA(ls_flight).
  SELECT SINGLE carrname FROM scarr INTO @DATA(lv_carrname) WHERE carrid = @ls_flight-carrid.
  WRITE: / ls_flight-connid, lv_carrname.
ENDLOOP.

" ĐÚNG — 1 JOIN duy nhất, hoặc dùng FOR ALL ENTRIES nếu bắt buộc phải tách 2 query
SELECT sflight~connid, scarr~carrname
  FROM sflight
  INNER JOIN scarr ON sflight~carrid = scarr~carrid
  INTO TABLE @DATA(lt_result)
  WHERE sflight~carrid = 'LH'.
```

Đây tương đương lỗi N+1 query trong Go/Python (gọi `db.QueryRowContext` lặp lại trong `for` thay vì 1 `JOIN` — xem [Go Bài 17](../Go/17_database.md)) — nhưng trong ABAP, hậu quả hiệu năng còn nghiêm trọng hơn vì có thể chạy trên dữ liệu ERP hàng triệu dòng.

## 5. `FOR ALL ENTRIES` — khi buộc phải tách query (bảng ở 2 nguồn khác nhau)

```abap
IF lt_flights IS NOT INITIAL.  " BẮT BUỘC kiểm tra rỗng trước FOR ALL ENTRIES — nếu rỗng sẽ lấy TOÀN BỘ bảng!
  SELECT carrid, carrname
    FROM scarr
    INTO TABLE @DATA(lt_carriers)
    FOR ALL ENTRIES IN @lt_flights
    WHERE carrid = @lt_flights-carrid.
ENDIF.
```

**Cạm bẫy kinh điển:** nếu `lt_flights` rỗng mà không kiểm tra trước, `FOR ALL ENTRIES` sẽ bỏ qua điều kiện `WHERE` và lấy **TOÀN BỘ bảng** — luôn nhớ `IF ... IS NOT INITIAL` trước khi dùng.

## 6. Liên hệ CDS View trong Open SQL (xem đầy đủ ở [Bài 10](./10_cds_basics.md))

```abap
" CDS View cũng có thể SELECT trực tiếp như 1 bảng thông thường
SELECT carrid, connid, total_price
  FROM zc_flight_overview  " tên CDS View, không phải table thật
  INTO TABLE @DATA(lt_overview)
  WHERE carrid = 'LH'.
```

## Ví dụ đầy đủ

```abap
REPORT z_flight_report.

SELECT sflight~carrid, sflight~connid, sflight~fldate, sflight~price,
       scarr~carrname
  FROM sflight
  INNER JOIN scarr ON sflight~carrid = scarr~carrid
  INTO TABLE @DATA(lt_flights)
  WHERE sflight~carrid = 'LH'
  ORDER BY sflight~fldate.

LOOP AT lt_flights INTO DATA(ls_flight).
  WRITE: / ls_flight-carrname, ls_flight-connid, ls_flight-fldate, ls_flight-price.
ENDLOOP.
```

## Bài tập

1. **`SELECT` inline**: viết `SELECT` lấy dữ liệu từ `SFLIGHT` (hoặc bảng tương tự) vào internal table bằng inline declaration, chỉ định rõ field cần lấy.
2. **JOIN**: viết `SELECT` JOIN `SFLIGHT` với `SCARR`, lấy tên hãng bay kèm thông tin chuyến bay.
3. **Sửa lỗi N+1**: viết 1 đoạn code cố tình có `SELECT` bên trong `LOOP` (anti-pattern), sau đó viết lại bằng `JOIN` hoặc `FOR ALL ENTRIES`, so sánh.
4. **`FOR ALL ENTRIES` an toàn**: viết code dùng `FOR ALL ENTRIES` có kiểm tra `IS NOT INITIAL` trước — giải thích bằng comment điều gì xảy ra nếu bỏ qua bước kiểm tra này.

## Tiếp theo
→ [Bài 8: Xử lý lỗi trong ABAP](./8_abap_exceptions.md)
