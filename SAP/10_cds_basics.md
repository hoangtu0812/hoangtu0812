# Bài 10: CDS View cơ bản

## Mục tiêu
- Hiểu Core Data Services (CDS) — nền tảng data model hiện đại của SAP.
- Viết `DEFINE VIEW ENTITY`, dùng annotation cơ bản.
- Vì sao CDS thay thế dần view cổ điển (SE11 database view, `INFOSET`).

## 1. CDS View là gì?

CDS View là view được định nghĩa bằng ngôn ngữ khai báo (DDL — Data Definition Language) ngay trong ADT, biên dịch xuống database view thật sự. Khác view cổ điển (SE11), CDS hỗ trợ **annotation** phong phú (metadata cho UI, phân tích, phân quyền), **association** (quan hệ giữa các view), và là nền tảng bắt buộc cho RAP ([Bài 12](./12_rap_intro.md)).

So sánh nhanh: CDS View giống việc định nghĩa 1 SQL `VIEW` có kèm theo rất nhiều "annotation" (metadata) mà Go/Python không có khái niệm tương đương trực tiếp — gần nhất có thể ví CDS annotation như OpenAPI spec tự sinh từ code (`swaggo` ở Go, tương tự Pydantic model của FastAPI ở Python) nhưng ở tầng data model.

## 2. Cú pháp cơ bản

```abap
@AbapCatalog.sqlViewName: 'ZVFLIGHTOVW'
@AccessControl.authorizationCheck: #CHECK
@EndUserText.label: 'Flight Overview'
define view entity ZC_FlightOverview
  as select from sflight
{
  key carrid,
  key connid,
      fldate,
      price,
      currency,
      seatsmax,
      seatsocc
}
```

- `@AbapCatalog.sqlViewName`: tên view thật trong database (tối đa 16 ký tự, quy ước cũ).
- `define view entity`: cú pháp CDS hiện đại (khác `define view` cũ hơn, không hỗ trợ đầy đủ tính năng RAP).
- `key`: đánh dấu field là khóa chính của view — BẮT BUỘC phải có ít nhất 1 field `key`.

## 3. Dùng CDS View trong Open SQL (liên hệ [Bài 7](./7_open_sql.md))

```abap
SELECT carrid, connid, price
  FROM zc_flightoverview
  INTO TABLE @DATA(lt_flights)
  WHERE carrid = 'LH'.
```

CDS View dùng được y hệt 1 bảng thông thường trong `SELECT` — điểm mạnh lớn: logic join/tính toán phức tạp được "đóng gói" 1 lần trong CDS, mọi nơi dùng lại chỉ cần `SELECT` đơn giản.

## 4. Annotation `@UI` — hiển thị trên Fiori Elements (liên hệ [Bài 20](./20_fiori_deployment.md))

```abap
@UI: {
  headerInfo: { typeName: 'Flight', typeNamePlural: 'Flights' }
}
define view entity ZC_FlightOverview
  as select from sflight
{
  key carrid,
  key connid,

  @UI.lineItem: [{ position: 10 }]
  @UI.identification: [{ position: 10 }]
  fldate,

  @UI.lineItem: [{ position: 20 }]
  price,

  currency
}
```

`@UI.lineItem` đánh dấu field nào hiển thị trong danh sách (List Report) khi CDS View được expose thành Fiori Elements app — bạn KHÔNG cần viết code frontend riêng, chỉ cần annotation đúng.

## 5. Field alias, tính toán trong CDS

```abap
define view entity ZC_FlightOverview
  as select from sflight
{
  key carrid,
  key connid,
      fldate as flight_date,              // alias — đổi tên field khi expose ra ngoài
      price * 1.1 as price_with_tax,        // biểu thức tính toán ngay trong CDS
      case
        when seatsocc >= seatsmax then 'FULL'
        else 'AVAILABLE'
      end as availability_status            // logic điều kiện — tương tự CASE của SQL chuẩn
}
```

## 6. Vì sao CDS thay thế view cổ điển

| | View cổ điển (SE11) | CDS View |
|---|---|---|
| Association/quan hệ | Không hỗ trợ trực tiếp | Có (`ASSOCIATION TO` — xem [Bài 11](./11_cds_advanced.md)) |
| Annotation UI/Auth | Không có | Có (`@UI`, `@AccessControl`...) |
| Nền tảng cho RAP | Không | **Bắt buộc** — mọi Business Object RAP đều dựa trên CDS |
| Code pushdown (tính toán tại DB, không kéo dữ liệu về ABAP) | Hạn chế | Tối ưu, đặc biệt trên HANA |

## Ví dụ đầy đủ

```abap
@AbapCatalog.sqlViewName: 'ZVPRODOVERVIEW'
@AccessControl.authorizationCheck: #CHECK
@EndUserText.label: 'Product Overview'
@UI: { headerInfo: { typeName: 'Product', typeNamePlural: 'Products' } }
define view entity ZC_ProductOverview
  as select from zproduct
{
  key product_id,

  @UI.lineItem: [{ position: 10 }]
  @UI.identification: [{ position: 10 }]
  product_name,

  @UI.lineItem: [{ position: 20 }]
  price,

  @UI.lineItem: [{ position: 30 }]
  case
    when stock_qty = 0 then 'OUT_OF_STOCK'
    when stock_qty < 10 then 'LOW_STOCK'
    else 'IN_STOCK'
  end as stock_status
}
```

## Bài tập

1. **CDS View đơn giản**: viết CDS view trên 1 bảng chuẩn (`SFLIGHT` hoặc tương tự), chọn ra 5-6 field, kiểm tra bằng Data Preview trong ADT.
2. **Thêm annotation `@UI`**: thêm `@UI.lineItem` cho các field muốn hiển thị dạng danh sách.
3. **Field tính toán**: thêm 1 field tính toán (vd giá đã bao gồm thuế) và 1 field `CASE WHEN` (vd trạng thái dựa trên số lượng).
4. **Dùng trong Open SQL**: viết 1 report `SELECT` từ CDS view vừa tạo, so sánh với việc `SELECT` trực tiếp từ bảng gốc.

## Tiếp theo
→ [Bài 11: CDS nâng cao](./11_cds_advanced.md)
