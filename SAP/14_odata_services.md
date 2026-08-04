# Bài 14: OData Services

## Mục tiêu
- Expose CDS/RAP Business Object qua OData V2/V4.
- Service Definition + Service Binding.
- Test qua SAP Gateway Client / Postman.

## 1. OData là gì, và vì sao SAP dùng nó làm chuẩn

OData (Open Data Protocol) là chuẩn REST API mở rộng, có metadata mô tả entity/field/quan hệ (`$metadata`), hỗ trợ query linh hoạt qua URL (`$filter`, `$select`, `$expand`, `$top`...) mà không cần viết endpoint riêng cho từng nhu cầu — khác REST API tự viết tay của Go/Python ([Go Bài 18](../Go/18_rest_api.md), [Python Bài 18](../Python/18_rest_api.md)) nơi bạn tự định nghĩa từng route.

## 2. Service Definition — chọn CDS View nào để expose

```abap
@EndUserText.label: 'Task Service Definition'
define service ZUI_TASK
{
  expose ZC_Task as Task;
}
```

Service Definition liệt kê các Projection View (đã có Behavior Definition — [Bài 12](./12_rap_intro.md)) sẽ được expose ra OData — tương đương việc chọn router nào `app.include_router()` trong FastAPI ([Python Bài 18](../Python/18_rest_api.md)).

## 3. Service Binding — chọn protocol (OData V2/V4) và publish

Tạo Service Binding trong ADT (New → Service Binding), chọn:
- **Binding Type**: `OData V4 - UI` (khuyến khích cho app mới, đặc biệt Fiori Elements) hoặc `OData V2 - UI` (tương thích ngược, nhiều hệ thống cũ vẫn dùng).
- **Service Definition**: trỏ tới `ZUI_TASK` vừa tạo.

Sau khi tạo, click **Publish** để service sẵn sàng phục vụ request — tương đương `r.Run(":8080")` của Gin ([Go Bài 18](../Go/18_rest_api.md)) hoặc `uvicorn main:app` của FastAPI, nhưng ở đây là "publish" trong hệ thống ABAP thay vì chạy process riêng.

## 4. Test qua Preview / Gateway Client

Trong ADT, click phải vào Service Binding → **Preview** để mở Fiori Elements UI tự sinh (test nhanh CRUD qua giao diện — không cần viết code frontend).

Hoặc test thô qua URL trực tiếp (Gateway Client / Postman):

```
GET /sap/opu/odata4/sap/zui_task/srvd/sap/zui_task/0001/Task
GET /sap/opu/odata4/sap/zui_task/srvd/sap/zui_task/0001/Task(TaskId='...')
POST /sap/opu/odata4/sap/zui_task/srvd/sap/zui_task/0001/Task
PATCH /sap/opu/odata4/sap/zui_task/srvd/sap/zui_task/0001/Task(TaskId='...')
DELETE /sap/opu/odata4/sap/zui_task/srvd/sap/zui_task/0001/Task(TaskId='...')
```

## 5. Query option — sức mạnh built-in của OData

```
GET /Task?$filter=Done eq false
GET /Task?$select=TaskId,Title
GET /Task?$orderby=LastChangedAt desc
GET /Task?$top=10&$skip=20
GET /Task?$expand=_Owner
```

So sánh: những gì bạn phải tự viết bằng tay trong Go (`?done=false` + parse query param thủ công — [Go Bài 18](../Go/18_rest_api.md)) hoặc FastAPI (Pydantic query param), OData **tự động hỗ trợ sẵn** — chỉ cần expose đúng field và association trong CDS View.

## 6. `$metadata` — API self-describing

```
GET /sap/opu/odata4/sap/zui_task/srvd/sap/zui_task/0001/$metadata
```

Trả về XML mô tả toàn bộ entity, field, kiểu dữ liệu, action — tương đương OpenAPI/Swagger spec tự sinh (`/docs` của FastAPI — [Python Bài 18](../Python/18_rest_api.md), hoặc `swaggo` của Go) nhưng là chuẩn OData, được Fiori Elements dùng để tự sinh UI hoàn toàn không cần code frontend thủ công.

## 7. So sánh nhanh: RAP+OData vs tự viết REST API (Go/Python)

| | RAP + OData | Go/Python tự viết REST |
|---|---|---|
| Boilerplate CRUD | Gần như 0 — managed scenario tự sinh | Phải viết route + handler cho từng entity |
| Query linh hoạt ($filter, $expand...) | Có sẵn | Tự implement (hoặc dùng thư viện) |
| UI tự sinh | Có (Fiori Elements — [Bài 20](./20_fiori_deployment.md)) | Không — phải tự viết frontend |
| Tùy biến logic phức tạp | Cần học đúng "khuôn" RAP (validation/determination/action) | Tự do hoàn toàn |
| Phù hợp | Ứng dụng nghiệp vụ chuẩn, tích hợp ERP | Ứng dụng tự do, microservice độc lập |

## Bài tập

1. **Service Definition + Binding**: expose Business Object `Task` (Bài 12-13) qua Service Definition + Service Binding OData V4.
2. **Preview**: mở Preview trong ADT, thực hiện CRUD qua giao diện Fiori Elements tự sinh.
3. **Test qua Postman**: gọi trực tiếp OData endpoint bằng Postman/curl (GET list, GET single, POST tạo mới), so sánh format response với REST API bạn từng viết ở Go/Python.
4. **Query option**: thử `$filter`, `$select`, `$top` trên URL, quan sát kết quả trả về khác nhau thế nào mà không cần sửa code service.

## Tiếp theo
→ [Bài 15: Authorization trong ABAP](./15_abap_authorization.md)
