# Bài 1: Giới thiệu hệ sinh thái SAP & Cài đặt

## Mục tiêu
- Nắm bức tranh tổng thể hệ sinh thái SAP: ECC/S4HANA (on-premise), SAP BTP (cloud), Fiori/UI5 (frontend), và cách chúng giao nhau.
- Cài đặt môi trường: Eclipse + ADT, SAP BTP trial, Business Application Studio, `@sap/cds-dk`.

## 1. Bức tranh tổng thể

```
┌─────────────────────────┐        ┌──────────────────────────┐
│  ECC / S/4HANA (on-prem) │        │      SAP BTP (cloud)      │
│  - Ngôn ngữ: ABAP         │        │  - CAP (Node.js/Java)     │
│  - Data model: DDIC,       │◄──────►│  - HANA Cloud             │
│    CDS View                │  OData │  - XSUAA (auth)            │
│  - RAP (Business Object)   │  RFC   │  - Application Router      │
└─────────────────────────┘        └──────────────────────────┘
                    ▲                          ▲
                    └───────────┬──────────────┘
                                │
                      ┌───────────────────┐
                      │  Fiori / UI5 (FE)  │
                      │  Fiori Elements     │
                      └───────────────────┘
```

- **ECC/S4HANA**: hệ thống ERP lõi, logic nghiệp vụ viết bằng **ABAP**. Đây là nơi bạn đã và đang làm việc tại BSR.
- **SAP BTP (Business Technology Platform)**: nền tảng cloud, nơi viết extension/side-by-side application bằng **CAP** (Node.js hoặc Java), tận dụng HANA Cloud, các dịch vụ xác thực/tích hợp có sẵn.
- **Fiori/UI5**: framework frontend chuẩn của SAP, chạy trên cả 2 nền (ABAP RAP hoặc CAP) qua OData.

Lộ trình này đi theo đúng thứ tự: **ABAP cổ điển → ABAP hiện đại (CDS/RAP) → BTP & CAP**, để bạn hệ thống hóa từ nền tảng đã có tới phần cloud-native.

## 2. So sánh nhanh với Go/Python (bạn đã học ở 2 track kia)

| | ABAP (classic) | CAP (BTP) | Go/Python |
|---|---|---|---|
| Chạy ở đâu | ABAP Application Server | Node.js/Java runtime trên Cloud Foundry/Kyma | bất kỳ đâu |
| Data model | DDIC table, CDS View | CDS trong CAP (`.cds` file) | struct/dataclass + ORM |
| Expose API | OData Service (từ CDS/RAP) | OData/REST tự sinh từ CDS service | viết route thủ công (Gin/FastAPI) |
| Phân quyền | PFCG Role + Authorization Object | XSUAA scope + `@requires`/`@restrict` | middleware RBAC tự viết |
| Đặc biệt | Rất tích hợp sẵn với ERP data | CDS tự sinh OData — ít code hơn Go/Python cho CRUD cơ bản | tự do, không ràng buộc ERP |

## 3. Cài đặt môi trường

### Cho ABAP (Giai đoạn 1-2)
1. **SAP GUI** (nếu công ty dùng on-premise) hoặc quyền truy cập hệ thống ABAP thật.
2. **Eclipse + ABAP Development Tools (ADT)**: tải Eclipse, cài plugin ADT từ Update Site của SAP — đây là IDE hiện đại thay thế SE38/SE80 cổ điển, hỗ trợ tốt cho CDS/RAP.
3. Nếu không có hệ thống công ty để thực hành: đăng ký **SAP BTP ABAP Environment trial** (ABAP steampunk) — miễn phí, đủ để luyện Bài 2-15.

### Cho BTP/CAP (Giai đoạn 3)
1. Đăng ký **SAP BTP Trial account**: https://cockpit.hanatrial.ondemand.com
2. Cài **Node.js** (LTS) + CDS CLI:

```powershell
npm install -g @sap/cds-dk
cds --version
```

3. **SAP Business Application Studio (BAS)**: IDE cloud chạy ngay trong trình duyệt, tích hợp sẵn CDS/CAP/Fiori tooling — tạo dev space "Full Stack Cloud Application" trong BTP Cockpit.
4. (Tùy chọn) VS Code + extension "SAP CDS Language Support" nếu muốn code local thay vì BAS.

## Bài tập

1. **Khám phá ADT**: cài Eclipse + ADT, kết nối tới hệ thống ABAP hiện có ở công ty (hoặc trial), mở 1 chương trình mẫu (vd `DEMO_SELECT`), đọc code.
2. **Tạo BTP Trial**: đăng ký tài khoản trial, khám phá BTP Cockpit — tìm subaccount, space, các service khả dụng.
3. **Cài `@sap/cds-dk`**: cài Node.js + `@sap/cds-dk`, chạy `cds --version` để xác nhận.
4. **So sánh**: viết 3-5 dòng ghi chú (bằng lời) so sánh trải nghiệm dev ABAP cổ điển (SAP GUI) với CAP (BAS/VS Code) — điều gì bạn thấy quen thuộc từ Go/Python, điều gì hoàn toàn mới.

## Tiếp theo
→ [Bài 2: Cú pháp ABAP cơ bản](./2_abap_syntax.md)
