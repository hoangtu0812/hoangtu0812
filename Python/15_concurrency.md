# Bài 15: Concurrency (threading, multiprocessing, asyncio)

## Mục tiêu
- Hiểu GIL (Global Interpreter Lock) — khác biệt lớn nhất so với goroutine của Go.
- Biết khi nào dùng `threading`, `multiprocessing`, hay `asyncio`.
- Viết code bất đồng bộ cơ bản với `async`/`await`.

## 1. GIL là gì, và vì sao Python "concurrency" khác Go

Go's goroutine chạy song song thật sự trên nhiều CPU core (`GOMAXPROCS`). Python (CPython — bản triển khai phổ biến nhất) có **Global Interpreter Lock (GIL)**: tại 1 thời điểm, chỉ 1 thread được thực thi Python bytecode, dù máy có nhiều core.

Hệ quả thực tế — chọn công cụ theo loại tác vụ:

| Loại tác vụ | Công cụ phù hợp | Lý do |
|---|---|---|
| I/O-bound (gọi API, đọc file, query DB) | `asyncio` hoặc `threading` | Thread "nhường" CPU trong lúc chờ I/O — GIL không phải vấn đề |
| CPU-bound (tính toán nặng) | `multiprocessing` | Mỗi process có GIL RIÊNG, chạy thật sự song song trên nhiều core |

## 2. `threading` — dùng cho I/O-bound

```python
import threading
import time

def download(name: str):
	print(f"Bắt đầu tải {name}")
	time.sleep(2)  # giả lập I/O (gọi mạng)
	print(f"Đã tải xong {name}")

threads = []
for i in range(3):
	t = threading.Thread(target=download, args=(f"file{i}",))
	threads.append(t)
	t.start()

for t in threads:
	t.join()  # chờ tất cả thread hoàn thành — tương đương wg.Wait() của Go

print("Tất cả đã tải xong")
```

`threading.Lock` tương đương `sync.Mutex` của Go ([Go Bài 11](../Go/11_goroutines_channels.md)):

```python
import threading

class SafeCounter:
	def __init__(self):
		self.count = 0
		self.lock = threading.Lock()

	def increment(self):
		with self.lock:  # tương đương mu.Lock() ... defer mu.Unlock()
			self.count += 1

counter = SafeCounter()
threads = [threading.Thread(target=counter.increment) for _ in range(1000)]
for t in threads:
	t.start()
for t in threads:
	t.join()

print(counter.count)  # 1000 nhờ Lock
```

## 3. `multiprocessing` — dùng cho CPU-bound

```python
from multiprocessing import Pool
import time

def cpu_heavy_task(n: int) -> int:
	total = 0
	for i in range(n):
		total += i ** 2
	return total

if __name__ == "__main__":
	numbers = [10_000_000] * 4

	start = time.perf_counter()
	with Pool(processes=4) as pool:  # tạo 4 process con, mỗi process có GIL RIÊNG
		results = pool.map(cpu_heavy_task, numbers)
	print(f"Song song: {time.perf_counter() - start:.2f}s")
	print(results)
```

`multiprocessing` tốn overhead khởi tạo process cao hơn thread nhiều — chỉ dùng khi tác vụ thực sự nặng về CPU và đủ lớn để bù lại chi phí đó.

## 4. `asyncio` — cách hiện đại & phổ biến nhất cho I/O-bound trong backend

```python
import asyncio

async def fetch_data(name: str, delay: float) -> str:
	print(f"Bắt đầu gọi {name}")
	await asyncio.sleep(delay)  # "nhường" event loop cho tác vụ khác trong lúc chờ — giống <-time.After của Go
	print(f"Hoàn thành {name}")
	return f"Kết quả từ {name}"

async def main():
	# Chạy TUẦN TỰ (chậm — tổng ~4.5s)
	# r1 = await fetch_data("API-1", 1)
	# r2 = await fetch_data("API-2", 1.5)
	# r3 = await fetch_data("API-3", 2)

	# Chạy SONG SONG bằng asyncio.gather — tương đương chạy nhiều goroutine rồi WaitGroup.Wait()
	results = await asyncio.gather(
		fetch_data("API-1", 1),
		fetch_data("API-2", 1.5),
		fetch_data("API-3", 2),
	)
	print(results)  # tổng thời gian chỉ ~2s (thời gian của tác vụ CHẬM NHẤT), không phải tổng cộng

asyncio.run(main())
```

## 5. `asyncio` với timeout — tương đương `context.WithTimeout` của Go ([Go Bài 12](../Go/12_context.md))

```python
import asyncio

async def call_external_api():
	await asyncio.sleep(3)  # giả lập API chậm
	return "dữ liệu"

async def main():
	try:
		result = await asyncio.wait_for(call_external_api(), timeout=1.0)
		print(result)
	except asyncio.TimeoutError:
		print("Timeout: API mất quá lâu để phản hồi")

asyncio.run(main())
```

## 6. Vì sao `async def` quan trọng khi viết API (liên hệ [Bài 18](./18_rest_api.md))

FastAPI hỗ trợ cả `def` (chạy trong thread pool) và `async def` (chạy trực tiếp trên event loop). Dùng `async def` khi hàm handler có gọi I/O bất đồng bộ (query DB async, gọi API ngoài bằng `httpx.AsyncClient`) để server phục vụ được nhiều request đồng thời mà không cần tạo nhiều thread.

## Bài tập

1. **`threading` cơ bản**: chạy 3 thread giả lập tải file (dùng `time.sleep`), dùng `.join()` để chờ tất cả hoàn thành.
2. **`SafeCounter` với Lock**: viết như ví dụ trên, chạy 1000 thread cùng tăng counter, in kết quả cuối; thử bỏ Lock để quan sát kết quả sai (race condition).
3. **`asyncio.gather`**: viết 3 hàm `async def` giả lập gọi API với độ trễ khác nhau, chạy song song bằng `gather`, đo tổng thời gian và so sánh với chạy tuần tự.
4. **Nâng cao**: viết ví dụ dùng `asyncio.wait_for` để giới hạn timeout cho 1 tác vụ giả lập chậm, xử lý đúng `asyncio.TimeoutError`.

## Tổng kết Giai đoạn 2
Bạn đã nắm iterator/generator, decorator/context manager, file/JSON I/O, type hints, testing, và concurrency. Giai đoạn 3 sẽ áp dụng tất cả để xây dựng backend thực tế bằng FastAPI.

## Tiếp theo
→ [Bài 16: Kiến trúc project & Clean Architecture](./16_clean_architecture.md)
