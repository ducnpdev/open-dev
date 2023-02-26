# golang-concurrentcy

## error-handling

## leak

## prevent-leak
- Trong ví dụ trước đã nói về leak cũng như debug resource bằng `pprof`
- Thì bài viết này hỗ trợ làm sao để không bị leak khi sử dụng `goroutine` trong `golang`

### Ví dụ 1;
- Để thành công trong việc giảm thiểu leak trong goroutine thì dùng `channel` giữa các `goroutine cha và con`. Bởi theo quy định, signal luôn luôn chỉ đọc, và `goroutine` cha pass `channel` đến goroutine con. Khi `channel` close, nó sẽ close cả goroutine con.
- Code example:
```go
func Preven() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} { // (1)
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				case <-done: // (2)
					fmt.Println("done in work")
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan interface{})
	terminated := doWork(done, nil)
	go func() { // (3)
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()
	<-terminated // (4)
	fmt.Println("Done.")
}
``` 
- Ghi chú:
  - (1) `doWork()` là một function bình thường, khai báo trong func `Prevent()`. Nhận vào 2 parameter và return 1 parameter
  - (2) Trong line này, dùng `for-select` pattern, trong case `<-done` là kiểm tra `channel` có được báo tín hiệu chưa, nếu có thì sẽ return `goroutine`
  - (3) Tạo một `goroutine` khác, mục đích để cancel `doWork` sau thời gian 1 giây
  - (4) Để merge 2 goroutine lại với nhau, tiếp tục process những phần khác.
- Kêt quả:
```
Canceling doWork goroutine...
done in work
doWork exited.
Done.
```
- Như kết quả, mặc dù trong function `doWork()` truyền `string=nil` nhưng goroutine vẫn có thể exit, và clean-up resource. 
- Để có thể join 2 `goroutine` lại với nhau, tạo thêm 1 `goroutine` thứ 3, mục đích để `cancel` func `doWork()` sau 1 giây.

### Ví dụ 2:
- Trong ví dụ này, thử nghiệm thêm trường hợp đó là nhận `value` từ `channel`
- Code example:
```go
func LeakReceive() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.") // (1)
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
		}()
		return randStream
	}
	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}
```

- Ghi chú:
  - (1) khi dòng này xuất hiện thì `goroutine` đã được remove thành công.
- Kêt quả:
```
3 random ints:
1: 5577006791947779410
2: 8674665223082153551
3: 6129484611666145821
```
  - Trong print out, không có hàm `defer fmt.Println`, điều này đồng nghĩa nó không được thực thi => `leak`.
  - Sau khi 3 lần lặp, `goroutine` đã bị block và cố gắng send `random number` ra bên ngoài, nhưng có có `channel` read. Có nghĩa không có cách nào để `stop goroutine` đang chạy `random-number`

- Giải pháp, code:
```go
func PreventLeakReceive() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}
	done := make(chan interface{})

	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	// Simulate ongoing work
	time.Sleep(1 * time.Second)
}
```
- Ghi Chú:
  - Như ví dụ trước, cũng tạo thêm 1 `channle`, 1 `goroutine` thứ 3 => `terminates goroutine` thành công
- Kết quả:
```
3 random ints:
1: 5577006791947779410
2: 8674665223082153551
3: 6129484611666145821
newRandStream closure exited.
```
- Như đã nhìn thấy, goroutine thực sữ đã được clean-up.