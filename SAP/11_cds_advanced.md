# Bài 11: CDS Nâng Cao

## Mục tiêu
- Association (`ASSOCIATION TO`), path expression.
- Extend view, analytics annotation.
- Virtual element.

## 1. Association — quan hệ giữa các CDS View (tương đương JOIN nhưng khai báo, lazy)

```abap
define view entity ZC_FlightOverview
  as select from sflight
  association [1..1] to zc_carrier as _Carrier on $projection.carrid = _Carrier.carrid
{
  key carrid,
  key connid,
      fldate,
      price,

      _Carrier  // "expose" association ra ngoài để nơi dùng có thể "path expression" tới
}
```

`association [1..1] to ...` khai báo quan hệ nhưng **không tự động JOIN ngay** — chỉ khi consumer thực sự dùng tới (`_Carrier.carrname` trong path expression hoặc `SELECT`) thì SQL JOIN mới được sinh ra (lazy). Đây khác JOIN tường minh ở [Bài 7](./7_open_sql.md), tối ưu hơn khi không phải lúc nào cũng cần dữ liệu liên kết.

## 2. Path expression — dùng association để lấy field từ view liên kết

```abap
define view entity ZC_FlightWithCarrierName
  as select from ZC_FlightOverview as Flight
{
  key Flight.carrid,
  key Flight.connid,
      Flight.fldate,
      Flight._Carrier.carrname as carrier_name  // path expression — "đi qua" association để lấy field
}
```

Path expression tương đương truy cập field qua quan hệ (`user.Profile.Email`) trong ORM của Go/Python (GORM `Preload`, SQLAlchemy `relationship` — [Go Bài 17](../Go/17_database.md), [Python Bài 17](../Python/17_database.md)), nhưng được khai báo tường minh ngay trong data model thay vì code truy vấn.

## 3. Cardinality của Association

```
[1..1]  — chính xác 1 dòng liên kết (giống foreign key NOT NULL trỏ tới bảng cha)
[0..1]  — 0 hoặc 1 dòng (optional)
[0..*]  — 0 hoặc nhiều dòng (1-to-many)
[1..*]  — ít nhất 1, có thể nhiều
```

```abap
define view entity ZC_CarrierWithFlights
  as select from zc_carrier as Carrier
  association [0..*] to ZC_FlightOverview as _Flights on $projection.carrid = _Flights.carrid
{
  key Carrier.carrid,
      Carrier.carrname,
      _Flights  // expose để consumer truy vấn danh sách chuyến bay của hãng này
}
```

## 4. Extend View — mở rộng CDS view có sẵn KHÔNG sửa code gốc

```abap
extend view entity ZC_FlightOverview with
{
  extension.zz_custom_note  // thêm field từ bảng extension (append structure) mà không đụng vào view gốc
}
```

`extend view` tương đương composition/embedding thay vì sửa trực tiếp code người khác — an toàn khi cần custom trên chuẩn SAP mà không phá vỡ khả năng update hệ thống sau này (tương tự triết lý "mở để mở rộng, đóng để sửa đổi" — Open/Closed Principle).

## 5. Analytics annotation — cho báo cáo/dashboard

```abap
@Analytics.query: true
define view entity ZC_SalesAnalysis
  as select from zsalesorder
{
  key product_id,

  @Aggregation.default: #SUM
  quantity,

  @Aggregation.default: #SUM
  total_amount,

  @Analytics.dimension: true
  region
}
```

`@Aggregation.default: #SUM` cho biết field này nên được **tổng hợp** (sum) khi dùng trong báo cáo/pivot table — tương đương định nghĩa measure trong BI tool, giúp Fiori Analytics App tự động biết cách tổng hợp dữ liệu.

## 6. Virtual Element — field tính toán "ảo", KHÔNG lưu trong DB, tính lúc runtime bằng code ABAP

```abap
define view entity ZC_ProductOverview
  as select from zproduct
{
  key product_id,
      product_name,
      price,

  virtual discount_price : abap.dec(10,2)  // khai báo virtual element
}
```

```abap
" Class xử lý logic tính virtual element
CLASS zcl_product_discount_calc DEFINITION.
  PUBLIC SECTION.
    INTERFACES if_sadl_exit_calc_element_read.
ENDCLASS.

CLASS zcl_product_discount_calc IMPLEMENTATION.
  METHOD if_sadl_exit_calc_element_read~calculate.
    " logic tính discount_price dựa trên price, gán vào kết quả trả về
  ENDMETHOD.
ENDCLASS.
```

Virtual element hữu ích khi logic tính toán quá phức tạp để viết bằng biểu thức CDS thuần túy (mục 5 ở [Bài 10](./10_cds_basics.md)), cần code ABAP thật sự — nhưng nên hạn chế dùng vì mất đi lợi ích "code pushdown" (tính toán tại DB) của CDS.

## Ví dụ đầy đủ

```abap
define view entity ZC_OrderWithCustomer
  as select from zorder as Order
  association [1..1] to zcustomer as _Customer on $projection.customer_id = _Customer.customer_id
{
  key Order.order_id,
      Order.customer_id,
      Order.order_date,
      Order.total_amount,

      _Customer.customer_name,   // path expression lấy tên khách hàng
      _Customer.email as customer_email,

      _Customer  // vẫn expose association gốc để consumer có thể query sâu hơn nếu cần
}
```

## Bài tập

1. **Association + path expression**: viết 2 CDS view có quan hệ với nhau (vd `Order` và `Customer`), dùng association + path expression để lấy tên khách hàng ngay trong view `Order`.
2. **Cardinality**: thử cả `[1..1]` và `[0..*]`, giải thích khi nào dùng loại nào.
3. **Extend view**: thêm 1 field mới vào CDS view có sẵn bằng `extend view`, không sửa code gốc.
4. **Nâng cao — analytics**: viết CDS view có `@Analytics.query: true` với ít nhất 1 field `@Aggregation.default: #SUM` và 1 field `@Analytics.dimension: true`, giải thích ý nghĩa khi dùng trong Fiori Analytics.

## Tiếp theo
→ [Bài 12: RAP — Giới thiệu](./12_rap_intro.md)
