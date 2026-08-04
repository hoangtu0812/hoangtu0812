# Bài 4: Internal Table (Bảng nội bộ)

## Mục tiêu
- Hiểu Internal Table — cấu trúc dữ liệu quan trọng nhất trong ABAP, tương đương slice của Go / list của Python.
- Phân biệt `STANDARD`, `SORTED`, `HASHED` table.
- `LOOP AT`, `READ TABLE`, `APPEND`/`INSERT`/`DELETE`.

## 1. Vì sao Internal Table quan trọng bậc nhất trong ABAP?

Mọi kết quả `SELECT` từ database, mọi danh sách dữ liệu xử lý trong ABAP đều nằm trong internal table. Nắm vững nó quan trọng ngang việc nắm vững slice trong Go ([Go Bài 5](../Go/5_slices_maps.md)) hoặc list trong Python ([Python Bài 5](../Python/5_collections.md)).

## 2. Khai báo structure trước

```abap
TYPES: BEGIN OF ty_user,
         id    TYPE i,
         name  TYPE string,
         email TYPE string,
       END OF ty_user.
```

`TYPES ... BEGIN OF ... END OF` tương đương định nghĩa `struct` của Go hoặc `dataclass` của Python — mỗi dòng của internal table sẽ có cấu trúc này.

## 3. 3 loại Internal Table

| Loại | Đặc điểm | Khi nào dùng | Tương đương |
|---|---|---|---|
| `STANDARD TABLE` | Không có key, truy cập tuần tự | Danh sách đơn giản, thứ tự quan trọng | `[]T` của Go |
| `SORTED TABLE` | Luôn giữ thứ tự sắp xếp theo key | Cần dữ liệu luôn sắp xếp, tìm kiếm nhanh (binary search tự động) | tương tự `SORTED TABLE` không có tương đương trực tiếp |
| `HASHED TABLE` | Truy cập bằng hash key, O(1) | Bảng lớn, tra cứu theo key liên tục, KHÔNG cần duyệt tuần tự | `map[K]V` của Go |

```abap
DATA: lt_standard TYPE STANDARD TABLE OF ty_user,
      lt_sorted    TYPE SORTED TABLE OF ty_user WITH UNIQUE KEY id,
      lt_hashed    TYPE HASHED TABLE OF ty_user WITH UNIQUE KEY id.
```

## 4. Thêm dữ liệu — `APPEND`, `INSERT`, `VALUE #(...)`

```abap
DATA: lt_users TYPE STANDARD TABLE OF ty_user,
      ls_user  TYPE ty_user.

" Cách 1: gán field rồi APPEND
ls_user-id = 1.
ls_user-name = 'Ben'.
ls_user-email = 'ben@example.com'.
APPEND ls_user TO lt_users.

" Cách 2: APPEND VALUE #(...) — gọn hơn, 1 dòng
APPEND VALUE #( id = 2 name = 'Alice' email = 'alice@example.com' ) TO lt_users.

" Cách 3: khởi tạo cả table 1 lần bằng VALUE #(...) — tương đương slice literal của Go
DATA(lt_users2) = VALUE ty_user_table(
  ( id = 1 name = 'Ben' email = 'ben@example.com' )
  ( id = 2 name = 'Alice' email = 'alice@example.com' )
).
```

## 5. `LOOP AT` — duyệt table (nối tiếp [Bài 3](./3_abap_control_flow.md))

```abap
LOOP AT lt_users INTO DATA(ls_current).
  WRITE: / ls_current-id, ls_current-name, ls_current-email.
ENDLOOP.

" Duyệt kèm điều kiện lọc (Open SQL-style filter ngay trong LOOP AT — hiện đại, hiệu quả)
LOOP AT lt_users INTO DATA(ls_filtered) WHERE id > 1.
  WRITE / ls_filtered-name.
ENDLOOP.
```

## 6. `READ TABLE` — tìm 1 dòng cụ thể (tương đương lookup map hoặc tìm trong slice)

```abap
" Đọc theo index (chậm với STANDARD TABLE lớn nếu dùng sai cách)
READ TABLE lt_users INTO DATA(ls_by_index) INDEX 1.

" Đọc theo key — RẤT quan trọng dùng WITH KEY để ABAP tối ưu tìm kiếm
READ TABLE lt_users INTO DATA(ls_by_key) WITH KEY id = 2.
IF sy-subrc = 0.  " sy-subrc = 0 nghĩa là tìm thấy — kiểm tra BẮT BUỘC sau mỗi READ TABLE
  WRITE / ls_by_key-name.
ELSE.
  WRITE / 'Không tìm thấy'.
ENDIF.

" Với SORTED/HASHED table + khai báo key đúng, ABAP tự dùng binary search/hash — nhanh hơn nhiều so với STANDARD TABLE khi bảng lớn
READ TABLE lt_sorted INTO DATA(ls_sorted) WITH KEY id = 2 BINARY SEARCH.
```

`sy-subrc` là biến hệ thống chứa mã trả về của lệnh vừa thực thi — `0` nghĩa là thành công, tương đương kiểm tra `err != nil` của Go hoặc `if not found:` của Python, nhưng dựa vào **biến toàn cục ngầm định** thay vì giá trị trả về tường minh — đặc thù cần làm quen của ABAP.

## 7. `DELETE`, `MODIFY`

```abap
DELETE lt_users WHERE id = 1.

READ TABLE lt_users INTO DATA(ls_to_modify) WITH KEY id = 2.
ls_to_modify-name = 'Alice Updated'.
MODIFY lt_users FROM ls_to_modify INDEX sy-tabix.  " sy-tabix có được từ READ TABLE ngay trước đó
```

## Ví dụ đầy đủ

```abap
REPORT z_internal_table_demo.

TYPES: BEGIN OF ty_product,
         id    TYPE i,
         name  TYPE string,
         price TYPE p DECIMALS 2,
       END OF ty_product.

DATA(lt_products) = VALUE ty_product_table(
  ( id = 1 name = 'Bàn phím' price = '250000.00' )
  ( id = 2 name = 'Chuột'    price = '150000.00' )
  ( id = 3 name = 'Màn hình' price = '3500000.00' )
).

" In các sản phẩm > 200,000
LOOP AT lt_products INTO DATA(ls_product) WHERE price > '200000.00'.
  WRITE: / ls_product-name, ls_product-price.
ENDLOOP.

" Tìm 1 sản phẩm cụ thể
READ TABLE lt_products INTO DATA(ls_found) WITH KEY id = 2.
IF sy-subrc = 0.
  WRITE: / 'Tìm thấy:', ls_found-name.
ENDIF.
```

## Bài tập

1. **Tạo & duyệt table**: định nghĩa structure + internal table cho "sản phẩm" (id, tên, giá), thêm 5 dòng, duyệt và in bằng `LOOP AT`.
2. **Lọc & sắp xếp**: dùng `LOOP AT ... WHERE` để lọc, dùng `SORT lt_table BY field` để sắp xếp.
3. **`READ TABLE`**: tìm 1 dòng theo key, kiểm tra `sy-subrc`, xử lý cả 2 trường hợp tìm thấy/không.
4. **Nâng cao — so sánh hiệu năng**: tạo 1 `STANDARD TABLE` và 1 `HASHED TABLE` cùng 10.000 dòng, viết vòng lặp tra cứu 1.000 lần bằng `READ TABLE ... WITH KEY`, so sánh (bằng `GET RUN TIME` hoặc quan sát log) tốc độ giữa 2 loại table.

## Tiếp theo
→ [Bài 5: Modularization — Function Module & Method](./5_modularization.md)
