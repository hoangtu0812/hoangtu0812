# Bài 19: Authentication & Authorization trên BTP (XSUAA)

## Mục tiêu
- XSUAA service, `xs-security.json` (scope, role template).
- Annotation `@requires`/`@restrict` trong CDS service.
- Test với mock user (không cần deploy XSUAA thật để dev local).

## 1. XSUAA là gì?

XSUAA (Cross-Consumption User Account and Authentication) là OAuth2 Authorization Server managed service của BTP. Thay vì tự viết `create_access_token`/`verify_access_token` như Go/Python ([Go Bài 19](../Go/19_auth.md), [Python Bài 19](../Python/19_auth.md)), CAP tích hợp sẵn với XSUAA — bạn chỉ cần **khai báo scope/role** cần thiết, XSUAA lo phần cấp phát/verify JWT.

## 2. `xs-security.json` — định nghĩa scope & role template

```json
{
  "xsappname": "taskapp",
  "tenant-mode": "dedicated",
  "scopes": [
    { "name": "$XSAPPNAME.Admin", "description": "Quản trị toàn bộ task" },
    { "name": "$XSAPPNAME.User", "description": "Quản lý task của chính mình" }
  ],
  "role-templates": [
    { "name": "Admin", "description": "Admin role", "scope-references": ["$XSAPPNAME.Admin"] },
    { "name": "User", "description": "User role", "scope-references": ["$XSAPPNAME.User"] }
  ]
}
```

`scopes` tương đương định nghĩa permission (`admin`, `user`) trong Go/Python; `role-templates` là "khuôn" để admin BTP tạo Role Collection thật và gán cho user cụ thể (qua Cockpit) — tương tự PFCG Role của ABAP ([SAP Bài 15](./15_abap_authorization.md)) nhưng dành cho môi trường cloud.

## 3. `@requires` — yêu cầu đã xác thực

```cds
// srv/task-service.cds
using taskapp from '../db/schema';

@requires: 'authenticated-user'
service TaskService {
  entity Tasks as projection on taskapp.Tasks;
}
```

`@requires: 'authenticated-user'` tương đương middleware `Auth()` của Go/Gin hoặc dependency `get_current_user` của FastAPI — mọi request tới service này bắt buộc có JWT hợp lệ.

## 4. `@restrict` — phân quyền theo scope/role (RBAC)

```cds
service TaskService {
  @restrict: [
    { grant: ['READ', 'CREATE'], to: 'User' },
    { grant: '*', to: 'Admin' }
  ]
  entity Tasks as projection on taskapp.Tasks;
}
```

`@restrict` tương đương middleware `RequireRole("admin")` của Go ([Go Bài 19](../Go/19_auth.md)) hoặc dependency `require_role("admin")` của FastAPI ([Python Bài 19](../Python/19_auth.md)) — nhưng khai báo trực tiếp trong CDS, tự động áp dụng cho mọi operation (`READ`, `CREATE`, `UPDATE`, `DELETE`) mà không cần viết code kiểm tra thủ công.

## 5. Ownership check — logic vẫn cần code (giống Go/Python)

`@restrict` chỉ kiểm tra được role/scope, KHÔNG tự biết "task này có phải của user hiện tại không" — phần đó vẫn cần code trong event handler, giống hệt nguyên tắc ở [Go Bài 19 mục 7](../Go/19_auth.md) và [Python Bài 19](../Python/19_auth.md):

```javascript
// srv/task-service.js
srv.before(['UPDATE', 'DELETE'], Tasks, async (req) => {
  const user = req.user;  // CAP tự "giải mã" JWT, gắn sẵn thông tin user vào req.user
  const isAdmin = user.is('Admin');

  if (!isAdmin) {
    const task = await SELECT.one.from(Tasks).where({ ID: req.params[0].ID });
    if (task.owner !== user.id) {
      req.reject(403, 'Bạn không có quyền sửa/xóa task này');
    }
  }
});
```

`req.user.is('Admin')` kiểm tra role — tương đương `current_user["role"] == "admin"` của FastAPI hoặc `requesterRole != domain.RoleAdmin` của Go.

## 6. Mock user để test local — không cần deploy XSUAA thật

```json
// .cdsrc.json
{
  "requires": {
    "auth": {
      "kind": "mocked",
      "users": {
        "alice": { "password": "pass", "roles": ["User"] },
        "bob": { "password": "pass", "roles": ["Admin"] }
      }
    }
  }
}
```

Chạy `cds watch`, gọi API kèm HTTP Basic Auth (`alice:pass` hoặc `bob:pass`) — CAP tự giả lập JWT tương ứng, cho phép test toàn bộ `@requires`/`@restrict`/ownership check **mà không cần deploy XSUAA service thật lên BTP** — tương đương việc bạn tự tạo JWT giả trong test (`Go Bài 14`/`Python Bài 14`) nhưng CAP hỗ trợ ngay ở tầng framework.

## 7. Bật XSUAA thật khi deploy (liên hệ [Bài 20](./20_fiori_deployment.md))

```json
// package.json
{
  "cds": {
    "requires": {
      "auth": {
        "kind": "xsuaa"
      }
    }
  }
}
```

Khi deploy lên BTP (production profile), đổi `kind` từ `"mocked"` sang `"xsuaa"`, bind app với XSUAA service instance đã tạo từ `xs-security.json`.

## Bài tập

1. **`xs-security.json`**: viết file định nghĩa 2 scope `Admin`/`User` như ví dụ.
2. **`@requires`**: thêm `@requires: 'authenticated-user'` cho `TaskService`, test bằng mock user — xác nhận request không có auth bị từ chối.
3. **`@restrict`**: thêm `@restrict` chỉ cho `Admin` xóa task, `User` chỉ đọc/tạo, test với cả 2 mock user.
4. **Ownership check**: viết handler `before UPDATE/DELETE` kiểm tra `owner`, viết test (dùng `cds test` hoặc gọi thủ công) cho 3 case: admin sửa task người khác (được phép), user sửa task của mình (được phép), user sửa task người khác (bị từ chối) — giống bài tập RBAC ở [Go Bài 19](../Go/19_auth.md) và [Python Bài 19](../Python/19_auth.md).

## Tiếp theo
→ [Bài 20: Fiori Elements & Deployment](./20_fiori_deployment.md)
