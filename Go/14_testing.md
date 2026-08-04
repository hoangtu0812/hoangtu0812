# Bài 14: Testing trong Go

## Mục tiêu
- Viết unit test với package `testing` chuẩn của Go.
- Dùng table-driven test — pattern phổ biến nhất trong cộng đồng Go.
- Dùng thư viện `testify` để assertion dễ đọc hơn.
- Đo test coverage.

## 1. Quy ước file test

File test luôn có hậu tố `_test.go`, nằm cùng thư mục với code được test:

```
shapes/
├── rectangle.go
└── rectangle_test.go
```

## 2. Test cơ bản

```go
// rectangle.go
package shapes

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("không thể chia cho 0")
	}
	return a / b, nil
}
```

```go
// rectangle_test.go
package shapes

import "testing"

func TestDivide(t *testing.T) {
	result, err := Divide(10, 2)
	if err != nil {
		t.Fatalf("không mong đợi lỗi, nhưng nhận được: %v", err)
	}
	if result != 5 {
		t.Errorf("mong đợi 5, nhận được %d", result)
	}
}
```

Chạy: `go test ./...` (chạy toàn bộ project) hoặc `go test ./shapes/...` (1 package).

`t.Errorf` báo lỗi nhưng **tiếp tục chạy** test còn lại trong hàm; `t.Fatalf` báo lỗi và **dừng ngay** hàm test đó (dùng khi lỗi tiếp theo sẽ gây panic, vd nil pointer).

## 3. Table-driven test — pattern chuẩn của Go

Đây là cách viết test phổ biến nhất trong Go, giúp thêm case mới cực nhanh mà không lặp code:

```go
func TestDivideTableDriven(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		{name: "chia hết", a: 10, b: 2, want: 5, wantErr: false},
		{name: "chia có dư", a: 7, b: 2, want: 3, wantErr: false},
		{name: "chia cho 0", a: 10, b: 0, want: 0, wantErr: true},
		{name: "số âm", a: -10, b: 2, want: -5, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // subtest — chạy độc lập, tên hiện rõ khi fail
			got, err := Divide(tt.a, tt.b)

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr = %v, got err = %v", tt.wantErr, err)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Divide(%d, %d) = %d, muốn %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
```

Chạy `go test -v ./...` sẽ in ra tên từng subtest (`TestDivideTableDriven/chia_hết`, `.../chia_cho_0`...) — cực kỳ hữu ích khi debug case nào fail.

## 4. Dùng `testify` cho assertion gọn hơn

```powershell
go get github.com/stretchr/testify
```

```go
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDivideWithTestify(t *testing.T) {
	result, err := Divide(10, 2)
	require.NoError(t, err) // require: dừng ngay nếu fail (giống Fatalf)
	assert.Equal(t, 5, result) // assert: báo lỗi nhưng tiếp tục chạy (giống Errorf)

	_, err = Divide(10, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "không thể chia")
}
```

**Quy tắc:** dùng `require` cho các bước tiền đề (vd "phải không lỗi thì mới kiểm tra tiếp được"), dùng `assert` cho các kiểm tra độc lập nhau.

## 5. Test coverage

```powershell
go test -cover ./...              # in % coverage tổng quát
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # mở trình duyệt xem chi tiết dòng nào chưa test
```

## 6. Mock interface để test tầng service (liên hệ Bài 10 & 16)

```go
type UserRepository interface {
	FindByID(id int) (*User, error)
}

// Mock implementation dùng riêng cho test
type mockUserRepo struct {
	users map[int]*User
}

func (m *mockUserRepo) FindByID(id int) (*User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("không tìm thấy user %d", id)
	}
	return u, nil
}

func TestUserService_GetUser(t *testing.T) {
	repo := &mockUserRepo{users: map[int]*User{1: {ID: 1, Name: "Ben"}}}
	service := NewUserService(repo) // service phụ thuộc vào INTERFACE, không phải Postgres thật

	user, err := service.GetUser(1)
	require.NoError(t, err)
	assert.Equal(t, "Ben", user.Name)
}
```

Đây chính là lý do tầng service nên nhận `UserRepository` (interface) qua constructor thay vì phụ thuộc trực tiếp vào 1 implementation cụ thể — giúp test không cần database thật.

## 7. Benchmark cơ bản

```go
func BenchmarkDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Divide(10, 2)
	}
}
```

Chạy: `go test -bench=. -benchmem`

## Bài tập

1. **Table-driven test cho `divide`**: viết test bao phủ ít nhất 4 case (chia hết, chia dư, chia 0, số âm) theo pattern table-driven như trên.
2. **Coverage >80%**: viết đủ test cho package `shapes` (Bài 9) để đạt >80% coverage khi chạy `go test -cover`.
3. **Mock interface**: viết interface `Notifier` (liên hệ Bài 10), 1 struct triển khai thật + 1 mock cho test, viết test dùng mock để kiểm tra logic gọi `Notify` đúng số lần mong đợi.
4. **Nâng cao**: cài `testify`, viết lại toàn bộ test ở bài 1 bằng `assert`/`require`, so sánh độ ngắn gọn với cách viết `testing` thuần.

## Tiếp theo
→ [Bài 15: Generics](./15_generics.md)
