# Bài 11: Goroutines & Channels

## Mục tiêu
- Chạy hàm song song bằng `go func()`.
- Giao tiếp giữa goroutine bằng channel.
- Đồng bộ hóa bằng `sync.WaitGroup`, bảo vệ dữ liệu dùng chung bằng `sync.Mutex`.
- Phát hiện race condition bằng `go run -race`.

## 1. Goroutine — luồng nhẹ của Go

```go
func sayHello() {
	fmt.Println("Xin chào từ goroutine")
}

func main() {
	go sayHello() // chạy sayHello() song song, KHÔNG chờ nó xong
	fmt.Println("Xin chào từ main")
	time.Sleep(100 * time.Millisecond) // chờ tạm để goroutine kịp chạy (KHÔNG NÊN dùng trong code thật)
}
```

Goroutine cực kỳ nhẹ (vài KB stack ban đầu, tự động scale) — chương trình Go thực tế có thể chạy hàng chục nghìn goroutine cùng lúc. Nhưng `time.Sleep` để "chờ" là **anti-pattern** — cách đúng là dùng `sync.WaitGroup` hoặc channel.

## 2. `sync.WaitGroup` — chờ nhiều goroutine hoàn thành

```go
func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) // báo "có thêm 1 goroutine cần chờ"
		go func(id int) {
			defer wg.Done() // báo "goroutine này đã xong" khi hàm return
			fmt.Println("Goroutine", id, "đang chạy")
		}(i) // truyền i vào làm tham số để tránh bug closure bắt biến vòng lặp
	}

	wg.Wait() // block cho tới khi tất cả goroutine đã gọi Done()
	fmt.Println("Tất cả goroutine đã hoàn thành")
}
```

**Lưu ý quan trọng:** truyền `i` làm tham số của closure (`func(id int) {...}(i)`) thay vì dùng trực tiếp biến vòng lặp `i` bên trong closure — đây là lỗi kinh điển khi mới học goroutine (dù từ Go 1.22 vòng lặp đã tạo biến mới mỗi lần lặp nên bug này giảm bớt, vẫn nên viết tường minh cho dễ đọc và tương thích ngược).

## 3. Channel — giao tiếp giữa các goroutine

Triết lý Go: *"Don't communicate by sharing memory; share memory by communicating."*

```go
func main() {
	ch := make(chan string) // channel không buffer

	go func() {
		ch <- "Xin chào từ goroutine" // gửi giá trị vào channel
	}()

	msg := <-ch // nhận giá trị từ channel (block cho tới khi có dữ liệu)
	fmt.Println(msg)
}
```

Channel không buffer sẽ **block** cả bên gửi lẫn bên nhận cho tới khi cả hai sẵn sàng — đây là cơ chế đồng bộ hóa tự nhiên.

### Buffered channel

```go
ch := make(chan int, 3) // buffer chứa tối đa 3 giá trị
ch <- 1
ch <- 2
ch <- 3
// ch <- 4 // sẽ block vì buffer đầy, không ai nhận

close(ch) // đóng channel khi không còn gửi thêm nữa

for v := range ch { // duyệt hết giá trị còn lại rồi tự dừng khi channel đóng
	fmt.Println(v)
}
```

### `select` — chờ nhiều channel cùng lúc

```go
func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() { time.Sleep(1 * time.Second); ch1 <- "từ ch1" }()
	go func() { time.Sleep(2 * time.Second); ch2 <- "từ ch2" }()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		}
	}
}
```

## 4. Producer-Consumer pattern

```go
func producer(ch chan<- int, count int) { // chan<- : channel chỉ dùng để GỬI
	for i := 1; i <= count; i++ {
		ch <- i
	}
	close(ch)
}

func consumer(ch <-chan int, done chan<- bool) { // <-chan : channel chỉ dùng để NHẬN
	for v := range ch {
		fmt.Println("Nhận được:", v)
	}
	done <- true
}

func main() {
	ch := make(chan int)
	done := make(chan bool)

	go producer(ch, 5)
	go consumer(ch, done)

	<-done // chờ consumer báo xong
}
```

## 5. `sync.Mutex` — bảo vệ dữ liệu dùng chung

Khi nhiều goroutine cùng đọc/ghi 1 biến, phải khóa bằng Mutex để tránh **race condition**:

```go
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // đảm bảo unlock dù có panic xảy ra
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("Kết quả:", counter.Value()) // luôn = 1000 nhờ Mutex
}
```

**Kiểm tra race condition:** chạy `go run -race main.go` — nếu bỏ Mutex ra, Go race detector sẽ báo lỗi ngay khi có truy cập đồng thời không an toàn.

## Bài tập

1. **5 goroutine song song**: chạy 5 goroutine in số 1–5, dùng `sync.WaitGroup` để `main` chờ tất cả hoàn thành trước khi kết thúc.
2. **Producer-consumer**: dùng code mẫu trên làm nền, tự viết lại producer gửi 10 số, consumer in ra và tính tổng.
3. **Counter an toàn**: viết `SafeCounter` với `sync.Mutex` như trên, chạy 1000 goroutine cùng tăng counter, in kết quả cuối. Sau đó **cố tình xóa Mutex** và chạy `go run -race` để quan sát race condition được phát hiện.
4. **Nâng cao — worker pool**: viết chương trình có 3 "worker" goroutine cùng nhận job từ 1 channel chung (`jobs chan int`), xử lý (giả lập bằng in ra + `time.Sleep`), gửi kết quả vào channel `results chan int`, dùng `select` hoặc `WaitGroup` để thu thập đủ kết quả.

## Tiếp theo
→ [Bài 12: Context](./12_context.md)
