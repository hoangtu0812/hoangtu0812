# Bài 17: CAP Cơ Bản (Cloud Application Programming Model)

## Mục tiêu
- `cds init`, định nghĩa data model + service definition bằng CDS trong CAP.
- `cds watch`, so sánh CDS trong CAP với CDS trong ABAP ([Bài 10](./10_cds_basics.md)).

## 1. CAP là gì?

CAP (Cloud Application Programming Model) là framework của SAP để xây dựng backend service trên BTP, dùng **CDS làm data model + service definition**, chạy trên Node.js hoặc Java runtime. CAP tự động sinh OData/REST API từ CDS — cùng triết lý "khai báo nhiều, code ít" như RAP ([Bài 12](./12_rap_intro.md)), nhưng chạy cloud-native thay vì trong ABAP application server.

So với Go/Python: CAP giống việc kết hợp Pydantic model + SQLAlchemy model + FastAPI router **tự động sinh ra chỉ từ 1 file CDS khai báo**, thay vì viết tay từng phần như [Python Bài 16-18](../Python/16_clean_architecture.md).

## 2. `cds init` — khởi tạo project

```powershell
cds init taskapp
cd taskapp
npm install
```

Cấu trúc sinh ra:

```
taskapp/
├── db/            # data model (.cds)
├── srv/            # service definition (.cds) + custom logic (.js)
├── app/             # UI (Fiori) — thêm sau ở Bài 20
├── package.json
└── .cdsrc.json      # config CAP (thêm sau)
```

## 3. Data model — `db/schema.cds`

```cds
namespace taskapp;

entity Tasks {
  key ID    : UUID;
      title : String(255) not null;
      done  : Boolean default false;
      owner : String(255);
}
```

So sánh với CDS trong ABAP ([Bài 10](./10_cds_basics.md)): cú pháp **gần giống nhau** (cả 2 đều dùng CDS DDL), nhưng CDS trong CAP **định nghĩa bảng mới** (giống `CREATE TABLE`), còn CDS trong ABAP (`define view entity`) thường **định nghĩa view trên bảng đã có sẵn**. Đây là khác biệt quan trọng nhất khi chuyển đổi tư duy giữa 2 track.

## 4. Service Definition — `srv/task-service.cds`

```cds
using taskapp from '../db/schema';

service TaskService {
  entity Tasks as projection on taskapp.Tasks;
}
```

Chỉ với 3 dòng này, CAP đã tự sinh đầy đủ OData endpoint CRUD cho `Tasks` — tương đương việc bạn phải tự viết cả 5 route (`GET/POST/PUT/DELETE`) trong FastAPI/Gin ([Go Bài 18](../Go/18_rest_api.md), [Python Bài 18](../Python/18_rest_api.md)).

## 5. `cds watch` — chạy dev server

```powershell
cds watch
```

Mặc định dùng SQLite in-memory cho dev (không cần cấu hình HANA Cloud ngay) — tự động tạo bảng từ `schema.cds`, tự reload khi sửa code. Truy cập:

```
http://localhost:4004/odata/v4/task/Tasks
```

## 6. Thêm dữ liệu mẫu — CSV seed data

```
db/data/taskapp-Tasks.csv
```

```csv
ID,title,done,owner
1a2b3c4d-...,Học CAP,false,ben
```

CAP tự động load file CSV này vào SQLite khi `cds watch` khởi động — tiện cho việc test nhanh mà không cần viết migration/seed script riêng như Alembic ([Python Bài 17](../Python/17_database.md)).

## 7. So sánh nhanh CDS trong CAP vs CDS trong ABAP

| | CDS trong ABAP (RAP) | CDS trong CAP |
|---|---|---|
| Chạy ở đâu | ABAP Application Server | Node.js/Java trên BTP |
| Vai trò chính | View trên bảng có sẵn, nền tảng cho RAP | Định nghĩa bảng MỚI + service |
| Sinh OData | Cần Service Definition + Binding riêng ([Bài 14](./14_odata_services.md)) | Tự động từ Service Definition |
| Custom logic | Behavior Implementation (ABAP class) | Event handler (JavaScript/TypeScript hoặc Java — [Bài 18](./18_cap_advanced.md)) |
| Database | HANA/AnyDB có sẵn của hệ thống ERP | HANA Cloud hoặc SQLite (dev) |

## Ví dụ đầy đủ

```cds
// db/schema.cds
namespace taskapp;

entity Tasks {
  key ID    : UUID;
      title : String(255) not null;
      done  : Boolean default false;
      owner : String(255);
}
```

```cds
// srv/task-service.cds
using taskapp from '../db/schema';

service TaskService {
  entity Tasks as projection on taskapp.Tasks;
}
```

Chạy `cds watch`, test:
```
GET  http://localhost:4004/odata/v4/task/Tasks
POST http://localhost:4004/odata/v4/task/Tasks
     Body: { "title": "Học CAP", "owner": "ben" }
```

## Bài tập

1. **`cds init`**: tạo project `taskapp`, định nghĩa entity `Tasks` như ví dụ, chạy `cds watch`.
2. **Service Definition**: expose `Tasks` qua `TaskService`, test CRUD qua URL (Postman hoặc trình duyệt cho GET).
3. **Seed data**: thêm file CSV mẫu vào `db/data/`, xác nhận dữ liệu load khi `cds watch` khởi động lại.
4. **So sánh cấu trúc**: viết ghi chú so sánh `entity Tasks { ... }` (CAP) với `TypeModel`/`dataclass` + SQLAlchemy model (Python — [Python Bài 17](../Python/17_database.md)) hoặc `struct` + GORM tag (Go — [Go Bài 17](../Go/17_database.md)).

## Tiếp theo
→ [Bài 18: CAP nâng cao — Custom Logic](./18_cap_advanced.md)
