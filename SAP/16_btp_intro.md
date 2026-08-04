# Bài 16: Giới thiệu SAP BTP

## Mục tiêu
- Kiến trúc BTP: Cloud Foundry vs Kyma, Subaccount/Space.
- Các service chính: HANA Cloud, XSUAA, Destination, Application Router.

## 1. SAP BTP là gì?

SAP Business Technology Platform là nền tảng cloud (PaaS) của SAP, nơi bạn xây dựng **side-by-side extension** (mở rộng chức năng ERP mà không sửa trực tiếp core system) hoặc **ứng dụng độc lập** dùng CAP/Java/Node.js, tận dụng các managed service (database, auth, integration) thay vì tự vận hành hạ tầng — tương đương AWS/Azure/GCP nhưng tích hợp sâu với hệ sinh thái SAP.

## 2. Cloud Foundry vs Kyma — 2 runtime chính

| | Cloud Foundry | Kyma |
|---|---|---|
| Nền tảng | PaaS truyền thống, `cf push` deploy | Kubernetes-based, container-native |
| Phù hợp | CAP app, Java/Node.js truyền thống | Microservice phức tạp, cần Kubernetes trực tiếp |
| Độ phức tạp | Thấp hơn, khuyến khích cho người mới | Cao hơn, cần hiểu Kubernetes |

Lộ trình này (và capstone) dùng **Cloud Foundry** — phù hợp nhất cho CAP application ([Bài 17-18](./17_cap_basics.md)).

## 3. Subaccount & Space

```
Global Account (tổ chức)
  └── Subaccount (môi trường, vd "Dev", "Prod")
        └── Space (Cloud Foundry, nhóm ứng dụng + service instance)
              ├── App: taskapp-srv
              ├── App: taskapp-approuter
              └── Service instances: HANA Cloud, XSUAA, Destination
```

- **Subaccount**: đơn vị quản lý entitlement (quota dùng service gì, bao nhiêu), thường tách theo môi trường (dev/test/prod) hoặc theo team.
- **Space**: nhóm logic bên trong subaccount, chứa app + service instance đang chạy — tương đương "namespace" nếu so với Kubernetes.

## 4. Các service chính hay dùng

| Service | Vai trò | Liên hệ Go/Python |
|---|---|---|
| **HANA Cloud** | Database chính (hoặc dùng SQLite cho dev local) | tương đương managed PostgreSQL ([Go Bài 17](../Go/17_database.md), [Python Bài 17](../Python/17_database.md)) |
| **XSUAA** | Authorization server (OAuth2), quản lý scope/role | tương đương service tự implement JWT + RBAC ([Go Bài 19](../Go/19_auth.md), [Python Bài 19](../Python/19_auth.md)) — nhưng BTP cung cấp sẵn |
| **Destination** | Quản lý kết nối tới hệ thống ngoài (vd gọi ngược về S/4HANA on-premise) | tương đương config biến môi trường cho URL external API |
| **Application Router (approuter)** | Reverse proxy phía trước app, xử lý routing + tích hợp XSUAA login | tương đương API Gateway/nginx reverse proxy |
| **Cloud Logging/Application Logging** | Log tập trung | tương đương ELK/Loki mà `log/slog` (Go) hoặc `logging` (Python) ghi vào |

## 5. BTP Cockpit — giao diện quản lý

Đăng nhập https://cockpit.hanatrial.ondemand.com (trial) hoặc cockpit công ty, bạn sẽ thấy:
- **Entitlements**: dịch vụ nào được phép dùng, quota bao nhiêu.
- **Instances and Subscriptions**: service đã tạo (vd HANA Cloud instance, XSUAA instance).
- **Spaces**: nơi app CAP của bạn chạy sau khi `cf push`.

## 6. Luồng 1 request đi qua BTP (tổng quan, chi tiết ở các bài sau)

```
Browser (Fiori app)
   │
   ▼
Application Router (approuter) ── xác thực qua XSUAA (OAuth2 login) ──► XSUAA
   │ (đã có JWT token hợp lệ)
   ▼
CAP Service (Node.js) ── kiểm tra scope trong JWT (@requires/@restrict — Bài 19) ──► trả data
   │
   ▼
HANA Cloud (hoặc SQLite lúc dev)
```

Luồng này tương đương: client → middleware auth (Go/Gin hoặc FastAPI `Depends`) → service → database — chỉ khác XSUAA + approuter là service **có sẵn**, không cần tự viết.

## Bài tập

1. **Khám phá Cockpit**: đăng nhập BTP Trial (hoặc subaccount công ty), tìm mục Entitlements, Instances, Spaces — ghi chú lại các service đang khả dụng.
2. **Tạo HANA Cloud instance** (hoặc xác nhận dùng SQLite cho dev nếu trial không đủ quota): tạo thử 1 instance database.
3. **So sánh kiến trúc**: vẽ (bằng lời hoặc sơ đồ ASCII) sơ đồ luồng request cho app CAP của bạn, đối chiếu với sơ đồ middleware/service/DB của Go hoặc Python bạn đã học.
4. **Cloud Foundry CLI**: cài `cf` CLI, đăng nhập vào subaccount trial (`cf login`), chạy `cf target` để xác nhận đã kết nối đúng org/space.

## Tiếp theo
→ [Bài 17: CAP cơ bản](./17_cap_basics.md)
