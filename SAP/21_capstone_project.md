# Dự án Capstone: CAP Application Có Phân Quyền

## Mục tiêu
Ghép toàn bộ kiến thức từ Bài 1 → Bài 20 thành 1 ứng dụng CAP hoàn chỉnh — **cùng ngữ cảnh nghiệp vụ với [Go/21_capstone_project.md](../Go/21_capstone_project.md) và [Python/21_capstone_project.md](../Python/21_capstone_project.md)**: quản lý "Task", `user` chỉ CRUD task của mình, `admin` quản lý mọi task, có UI Fiori Elements, phân quyền qua XSUAA + CDS `@restrict`.

**Kiến thức áp dụng:** [Bài 17](./17_cap_basics.md) (CAP cơ bản), [Bài 18](./18_cap_advanced.md) (custom logic), [Bài 19](./19_btp_auth.md) (XSUAA auth), [Bài 20](./20_fiori_deployment.md) (Fiori/deploy).

## Yêu cầu chức năng
Giống 2 bản Go/Python: `user` tạo/xem/sửa/xóa task của chính mình; `admin` xem/sửa/xóa mọi task; phân quyền qua scope `Admin`/`User`.

## Cấu trúc thư mục

```
taskapp/
├── db/
│   ├── schema.cds            # entity Users, Tasks
│   └── data/                  # CSV seed data cho dev
├── srv/
│   ├── task-service.cds       # service definition + @requires/@restrict
│   └── task-service.js         # custom handler: ownership check, validation
├── app/
│   └── taskapp/                # Fiori Elements app (sinh từ service)
├── xs-security.json            # scope/role cho XSUAA
├── mta.yaml                    # deployment descriptor
├── package.json
├── .cdsrc.json                  # mock auth cho dev local
└── test/
    └── task-service.test.js
```

## Bước 1: Bootstrap

```powershell
cds init taskapp
cd taskapp
npm install
npm install @cap-js/sqlite --save-dev
```

## Bước 2: Data model — `db/schema.cds`

```cds
namespace taskapp;

using { cuid, managed } from '@sap/cds/common';

entity Tasks : cuid, managed {
  title : String(255) not null;
  done  : Boolean default false;
  owner : String(255) not null;  // user ID của người tạo task
}
```

`cuid` (tự sinh UUID) và `managed` (tự thêm `createdAt`/`modifiedAt`) là aspect chuẩn của CAP — tương đương `id SERIAL PRIMARY KEY` + `created_at TIMESTAMP` viết tay ở Go/Python ([Go Bài 17](../Go/17_database.md), [Python Bài 17](../Python/17_database.md)), nhưng chỉ cần 1 dòng `using`.

## Bước 3: Service definition — `srv/task-service.cds`

```cds
using taskapp from '../db/schema';

@requires: 'authenticated-user'
service TaskService {

  @restrict: [
    { grant: ['READ', 'CREATE'], to: 'User' },
    { grant: ['READ', 'CREATE', 'UPDATE', 'DELETE'], to: 'Admin' }
  ]
  entity Tasks as projection on taskapp.Tasks actions {
    action markComplete() returns Tasks;
  };

  @restrict: [{ grant: '*', to: 'Admin' }]
  @readonly
  entity Users as projection on taskapp.Users;  // giả sử có entity Users riêng để admin xem danh sách
}
```

## Bước 4: Custom logic — `srv/task-service.js`

```javascript
const cds = require('@sap/cds');

module.exports = cds.service.impl(async function (srv) {
  const { Tasks } = this.entities;

  // Validation + gán owner tự động khi tạo mới
  srv.before('CREATE', Tasks, (req) => {
    if (!req.data.title?.trim()) {
      req.error(400, 'Title không được để trống');
    }
    req.data.owner = req.user.id;  // owner LUÔN là user hiện tại, KHÔNG cho client tự set
    req.data.done = false;
  });

  // Ownership check cho UPDATE/DELETE — logic CỐT LÕI của phân quyền, giống hệt Go/Python
  srv.before(['UPDATE', 'DELETE'], Tasks, async (req) => {
    if (req.user.is('Admin')) return;  // admin luôn được phép

    const task = await SELECT.one.from(Tasks).where({ ID: req.params[0].ID });
    if (!task) {
      req.error(404, 'Không tìm thấy task');
      return;
    }
    if (task.owner !== req.user.id) {
      req.reject(403, 'Bạn không có quyền sửa/xóa task này');
    }
  });

  // READ: user thường chỉ thấy task của mình, admin thấy tất cả
  srv.before('READ', Tasks, (req) => {
    if (!req.user.is('Admin')) {
      req.query.where({ owner: req.user.id });
    }
  });

  srv.on('markComplete', Tasks, async (req) => {
    const taskId = req.params[0].ID;
    await UPDATE(Tasks, taskId).set({ done: true });
    return SELECT.one.from(Tasks).where({ ID: taskId });
  });
});
```

Đối chiếu trực tiếp với Go ([Go Bài 21](../Go/21_capstone_project.md)) và Python ([Python Bài 21](../Python/21_capstone_project.md)): `srv.before(['UPDATE','DELETE'], ...)` ở đây làm CHÍNH XÁC công việc của `TaskService.UpdateTask`/`DeleteTask` (kiểm tra `requesterRole != admin && task.OwnerID != requesterID`) — chỉ khác là CAP gắn logic vào **event của entity** thay vì method của service class tự viết.

## Bước 5: `xs-security.json`

```json
{
  "xsappname": "taskapp",
  "tenant-mode": "dedicated",
  "scopes": [
    { "name": "$XSAPPNAME.Admin", "description": "Quản trị mọi task" },
    { "name": "$XSAPPNAME.User", "description": "Quản lý task cá nhân" }
  ],
  "role-templates": [
    { "name": "Admin", "scope-references": ["$XSAPPNAME.Admin"] },
    { "name": "User", "scope-references": ["$XSAPPNAME.User"] }
  ]
}
```

## Bước 6: Mock user cho dev local — `.cdsrc.json`

```json
{
  "requires": {
    "auth": {
      "kind": "mocked",
      "users": {
        "ben": { "password": "pass", "roles": ["User"] },
        "alice": { "password": "pass", "roles": ["User"] },
        "admin": { "password": "pass", "roles": ["Admin"] }
      }
    }
  }
}
```

## Bước 7: Test

```javascript
// test/task-service.test.js
const cds = require('@sap/cds');
const { GET, POST, PATCH, DELETE, expect } = cds.test(__dirname + '/..');

describe('TaskService ownership', () => {
  const benAuth = { auth: { username: 'ben', password: 'pass' } };
  const aliceAuth = { auth: { username: 'alice', password: 'pass' } };
  const adminAuth = { auth: { username: 'admin', password: 'pass' } };

  let taskId;

  it('ben tạo task mới thành công', async () => {
    const res = await POST('/odata/v4/task/Tasks', { title: 'Task của Ben' }, benAuth);
    expect(res.status).to.equal(201);
    taskId = res.data.ID;
  });

  it('alice KHÔNG sửa được task của ben', async () => {
    const res = await PATCH(`/odata/v4/task/Tasks(${taskId})`, { title: 'Cố sửa' }, aliceAuth)
      .catch(e => e.response);
    expect(res.status).to.equal(403);
  });

  it('admin sửa được task của ben', async () => {
    const res = await PATCH(`/odata/v4/task/Tasks(${taskId})`, { title: 'Sửa bởi admin' }, adminAuth);
    expect(res.status).to.equal(200);
  });
});
```

Chạy: `npm test` (CAP tích hợp sẵn `cds.test` dựa trên Jest/Mocha).

## Bước 8: Fiori Elements UI

```powershell
cds add fiori
```

Sinh `app/taskapp/annotations.cds`, cấu hình `UI.LineItem` hiển thị `title`, `done`, `owner` — xem chi tiết ở [Bài 20](./20_fiori_deployment.md).

## Bước 9: `mta.yaml` & Deploy

Dùng mẫu đầy đủ ở [Bài 20 mục 4](./20_fiori_deployment.md), điều chỉnh tên module thành `taskapp-srv`, `taskapp-db-deployer`, `taskapp-app-content`.

```powershell
mbt build
cf deploy mta_archives/taskapp_1.0.0.mtar
```

## Checklist hoàn thành

- [ ] Service yêu cầu xác thực (`@requires: 'authenticated-user'`) — request không có auth bị từ chối.
- [ ] `user` không sửa/xóa được task của người khác (test bằng 2 mock user `ben`/`alice`).
- [ ] `admin` (mock user `admin`) xem/sửa/xóa được mọi task.
- [ ] `owner` field không thể bị client tự set khi tạo task (luôn gán từ `req.user.id`).
- [ ] Fiori Elements app hiển thị đúng danh sách task theo quyền của user đăng nhập.
- [ ] `npm test` chạy pass với các case ownership ở Bước 7.
- [ ] Deploy thử thành công lên BTP trial (hoặc chạy tốt bằng `cds watch` với mock auth nếu chưa có quyền deploy).

## So sánh 3 stack

Sau khi hoàn thành cả 3 capstone (Go/Python/SAP), điểm khác biệt cốt lõi:

| | Go | Python (FastAPI) | SAP (CAP) |
|---|---|---|---|
| Boilerplate CRUD | Viết tay đầy đủ (handler/service/repo) | Viết tay, ít hơn Go nhờ Pydantic/FastAPI | **Gần như 0** — CDS tự sinh CRUD + OData |
| Ownership check | Code tường minh trong `service.UpdateTask` | Code tường minh trong `TaskService.update_task` | Code trong `srv.before(['UPDATE','DELETE'])` — cùng ý tưởng |
| Auth server | Tự viết JWT hoàn toàn | Tự viết JWT hoàn toàn | **XSUAA managed** — không tự viết |
| UI | Không có (chỉ API) | Không có (chỉ API, có Swagger docs) | **Fiori Elements tự sinh** từ annotation |
| Tốc độ phát triển CRUD cơ bản | Chậm nhất (nhiều boilerplate) | Trung bình | Nhanh nhất nếu đúng use case ERP |
| Độ tự do tùy biến | Cao nhất | Cao | Thấp hơn — phải theo "khuôn" CAP/RAP |

---
Hoàn thành cả 3 track (Go, Python, SAP) nghĩa là bạn có góc nhìn đối chiếu thực tế giữa: ngôn ngữ biên dịch hiệu năng cao (Go), ngôn ngữ động linh hoạt cho backend nhanh (Python), và nền tảng low-code/declarative chuyên biệt cho ERP (SAP CAP/RAP) — nền tảng vững để chọn đúng công cụ cho từng bài toán trong công việc thực tế tại BSR.
