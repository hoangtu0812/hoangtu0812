# Bài 18: CAP Nâng Cao — Custom Logic

## Mục tiêu
- Event handler (`srv.on`, `srv.before`, `srv.after`) bằng Node.js.
- Kết nối HANA Cloud/PostgreSQL thay vì SQLite dev.
- Validate và tính toán field trước khi trả response.

## 1. Event handler — nơi viết logic nghiệp vụ tùy chỉnh

```javascript
// srv/task-service.js
const cds = require('@sap/cds');

module.exports = cds.service.impl(async function (srv) {
  const { Tasks } = this.entities;

  // before CREATE — chạy TRƯỚC khi ghi vào DB, giống validation ở RAP (Bài 13) hoặc middleware validate ở Go/Python
  srv.before('CREATE', Tasks, async (req) => {
    if (!req.data.title || req.data.title.trim() === '') {
      req.error(400, 'Title không được để trống');
    }
  });

  // after READ — biến đổi dữ liệu SAU khi đọc, TRƯỚC khi trả về client
  srv.after('READ', Tasks, (tasks) => {
    const list = Array.isArray(tasks) ? tasks : [tasks];
    list.forEach((task) => {
      task.statusLabel = task.done ? 'Hoàn thành' : 'Đang xử lý';  // field tính toán, giống virtual element (SAP Bài 11) hoặc @property (Python Bài 7)
    });
  });

  // on CREATE — thay thế HOÀN TOÀN logic mặc định (ít dùng hơn before/after)
  srv.on('CREATE', Tasks, async (req) => {
    const task = { ...req.data, done: false };
    const result = await INSERT.into(Tasks).entries(task);
    return task;
  });
});
```

`before`/`after`/`on` tương đương middleware chain của Go ([Go Bài 18](../Go/18_rest_api.md)) hoặc dependency của FastAPI ([Python Bài 18](../Python/18_rest_api.md)) — nhưng gắn trực tiếp vào **entity + CRUD event** thay vì route path, nên logic tự động áp dụng cho MỌI cách client gọi tới entity đó (kể cả qua `$batch`, deep insert...).

## 2. Truy vấn database trong handler — CQL (Core Query Language)

```javascript
srv.on('READ', Tasks, async (req) => {
  return await SELECT.from(Tasks).where({ done: false });
});

srv.after('CREATE', Tasks, async (task, req) => {
  // Query liên quan sau khi tạo — vd cập nhật 1 bảng thống kê khác
  const { Stats } = this.entities;
  await UPDATE(Stats).set('taskCount += 1').where({ owner: task.owner });
});
```

CQL (`SELECT.from()`, `INSERT.into()`, `UPDATE()`) là DSL query builder của CAP — tương đương `select()` của SQLAlchemy ([Python Bài 17](../Python/17_database.md)) hoặc query builder Go — tự động tham số hóa, an toàn khỏi SQL Injection.

## 3. Custom action/function (nghiệp vụ ngoài CRUD chuẩn)

```cds
// srv/task-service.cds
service TaskService {
  entity Tasks as projection on taskapp.Tasks actions {
    action markComplete() returns Tasks;
  };
}
```

```javascript
// srv/task-service.js
srv.on('markComplete', Tasks, async (req) => {
  const taskId = req.params[0].ID;
  await UPDATE(Tasks, taskId).set({ done: true });
  return await SELECT.one.from(Tasks).where({ ID: taskId });
});
```

Tương đương action tùy chỉnh trong RAP ([Bài 13](./13_rap_transactional.md)) hoặc endpoint tùy chỉnh `POST /tasks/{id}/complete` tự viết ở Go/Python.

## 4. Kết nối HANA Cloud / PostgreSQL thay SQLite

```powershell
npm install @cap-js/postgres
```

```json
// package.json
{
  "cds": {
    "requires": {
      "db": {
        "kind": "postgres",
        "credentials": {
          "host": "localhost",
          "port": 5432,
          "database": "taskdb",
          "user": "user",
          "password": "pass"
        }
      }
    }
  }
}
```

```powershell
cds deploy --to postgres   # tạo bảng theo schema.cds trên Postgres thật
```

Với HANA Cloud trên BTP thực tế, thường dùng `cds bind` để tự động lấy credential từ service instance đã tạo trong Cockpit — không hardcode thông tin kết nối như trên (chỉ dùng cho local dev).

## 5. Draft, `requires` cho CAP (liên hệ [Bài 19](./19_btp_auth.md))

```cds
service TaskService {
  @odata.draft.enabled
  entity Tasks as projection on taskapp.Tasks;
}
```

`@odata.draft.enabled` bật draft handling cho CAP — tương đương `with draft` của RAP ([Bài 13](./13_rap_transactional.md)), CAP tự sinh bảng draft và logic Activate/Discard, dùng cho Fiori Elements Object Page phức tạp.

## Ví dụ đầy đủ

```javascript
// srv/task-service.js
const cds = require('@sap/cds');

module.exports = cds.service.impl(async function (srv) {
  const { Tasks } = this.entities;

  srv.before('CREATE', Tasks, (req) => {
    if (!req.data.title?.trim()) {
      req.error(400, 'Title không được để trống');
    }
    req.data.done = false;  // determination — luôn khởi tạo done = false khi tạo mới
  });

  srv.after('READ', Tasks, (tasks) => {
    const list = Array.isArray(tasks) ? tasks : [tasks];
    list.forEach((t) => { t.statusLabel = t.done ? 'Hoàn thành' : 'Đang xử lý'; });
  });

  srv.on('markComplete', Tasks, async (req) => {
    const taskId = req.params[0].ID;
    await UPDATE(Tasks, taskId).set({ done: true });
    return SELECT.one.from(Tasks).where({ ID: taskId });
  });
});
```

## Bài tập

1. **`before CREATE` validate**: thêm handler validate `title` không rỗng cho entity `Tasks` (Bài 17), test bằng POST thiếu field.
2. **`after READ` field tính toán**: thêm field `statusLabel` tính từ `done`, xác nhận nó xuất hiện trong response GET.
3. **Custom action**: thêm action `markComplete` như ví dụ, gọi qua Postman (`POST /Tasks(ID)/TaskService.markComplete`).
4. **Kết nối Postgres**: cấu hình `@cap-js/postgres`, chạy `cds deploy --to postgres`, xác nhận bảng `Tasks` được tạo trên Postgres thật thay vì SQLite.

## Tiếp theo
→ [Bài 19: Authentication & Authorization trên BTP (XSUAA)](./19_btp_auth.md)
