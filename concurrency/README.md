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

## select statement
- Để có thể biết `select` như thế nào, cách sử dụng, hãy xem ví dụ dưới.
```go
func example() {
	var ch1, ch2 <-chan interface{}
	var ch3 chan<- interface{}
	select {
	case <-ch1:
		// Do something
	case <-ch2:
		// Do something
	case ch3 <- struct{}{}:
		// Do something
	}
}
```
  - nhìn code thì nó gần như giống với `switch case`, nhưng thực tế là không, nó là list của các case statement, và là điểm kết thúc của `select` là các `case`.
  - Không như `switch case` các case không thực hiện tuần tự và cũng sẽ không tự động thực thi nếu không có tín hiệu đầu vào, ở đây là 1 `channel`
- Thay vào đó, tất cả `channel` sẽ được đọc và ghi đồng thời, nếu trong trường hợp không có `channel` nào ready, thì tất cả sẽ block. Đợi đến khi có 1 channel sắn sàng, thì sẽ lựa chọn case tương ứng để thực thi. 
- ví dụ:
```go
func Example2() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c) // (1)
	}()
	fmt.Println("Blocking on read...")
	select {
	case <-c: // (2)
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}
```
- ghi chú:
  - (1) close `channel` sau khi đợi 5 giây.
  - (2) đợi để read value từ `channel`.
- output:
```
Blocking on read...
Unblocked 5.001199292s later.
```
  - như kết quả thì process sẽ unblock sau 5 giây.
  - trong trường hợp này, cũng là một cách đơn giản và hiệu quả để đợi một process đang xảy ra cho đến khi kết thúc. Nhưng sẽ có một số vấn đặt ra
    - Điều gì sẽ xảy ra khi nhiều `channel` cùng read.
    - Điều gì sẽ xảy ra khi không có `channel` ready.
    - Sẽ như thế nào nếu bạn muốn process 1 logic khác mà không có `channel` ready.
  - Cùng đi giải quyết từng vấn đề
1. nhiều channel cùng read:
- code example:
```go
func MulChannel() {
	ch1 := make(chan interface{})
	close(ch1)
	ch2 := make(chan interface{})
	close(ch2)
	var ch1Count, ch2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-ch1:
			ch1Count++
		case <-ch2:
			ch2Count++
		}
	}
	fmt.Printf("ch1Count: %d\n ch2Count: %d\n", ch1Count, ch2Count)
}
```
- output:
```
- lan1
ch1Count: 498
ch2Count: 503

- lan2
ch1Count: 479
ch2Count: 522
```
  - Như kết quả, thì sau vòng lặp sấp xỉ một nữa được read bởi `ch1` và một nữa `ch2`. Có vẻ ngẫu nhiên nhỉ, vì trong `Golang` các case sẽ random giữa các statement nếu như có `channel` ready. Vậy trong một chương trình thực tế không thể để các state như thế được.

2. Không có channel là ready.
- Điều gì sẽ xảy ra khi không có channel nào ready, trong thực thế thì không thể để select block `forever`, do đó cần có timeout. Trong golang có package time có thể xử lý timeout được. 
- code:
```go
var ch <-chan int
select {
case <-ch: // (1) sẽ không bao giờ unblock bởi vì read từ 1 channel nill
case <-time.After(2 * time.Second):
  fmt.Println("Timed out.")
}
```
- output:
```
Timed out.
```
3. điều gì xảy ra khi không có channel nào ready, và cần làm một process khác sau 1 thời gian.
- như `switch-case` thì `select-case` cũng có giá trị `default`
- `default` là cho bạn lựa chọn khi tất cả các select case còn lại luôn luôn `block`
- code example 1:
```go
func SelectDefault() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("Value default %v\n\n", time.Since(start))
	}
}
```
- output:
```
Value default 334ns
--
Value default 208ns
```
- Như kết quả thì `default` sẽ giúp ta ngăn chặn block
- Cũng có thể được xem là waiting để đợi một `goroutine` khác process, và nhận kết quả
- code-example-2:
```go
func WaitingOther() {
	done := make(chan interface{})
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("done")
		close(done)
	}()
	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}
		fmt.Println(workCounter)
		workCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("workCounter %v .\n", workCounter)
}
```
- output:
```
0
1
2
done
workCounter 3
```
- Như ví dụ, thì chúng ta có 1 vòng lặp sẽ đợi 1 goroutine khác process, trong ví dụ thì process đó là `sleep 3s`

## livelock
- livelock là những chương trình tích cực thực hiện các hoạt động đồng thời, nhưng nhữn hoạt động này không làm gì để di chuyển các trạng thái đến các phần tiếp.
- livelock are programs that are actively performing concurrent operations, but these operations do nothing to move the state of the program forward.

- Bạn đã bao giờ trong hành lang và đi theo hướng một người khác? Người phía trước di chuyển sang 1 hướng để bạn vượt qua, nhưng bạn cũng vừa di chuyển sang hướng tương tự. Vì vậy bạn cũng di chuyển sang hướng khác và người phía trước cũng làm tương tự. Hãy tưởng tượng hành động đó lặp lại mãi mãi, thì bạn sẽ hiểu được livelooks.
- have you ever been in a hallway walking toward another person? She moves to one side to let you pass, but you are just done the same. So you move to the other side, and she is also done the same. Imagine this going on forever, and you understand livelooks.

- Để có thể hiểu hơn, chúng ta sẽ viết một đoạn code chứng minh. Trong đoạn code sẽ có một số hàm khó hiểu, mình nghĩa các bạn chỉ cẩn hiểu sơ thôi, không cần biết chi tiết nó làm gì.
- let's actually write some code that will help demonstrate this scenario. First, We will set up a few helper function, that will simplify the example. In order to have a working example, the code here utilize several topics we have not yet covered, I don't advise attemping to understand it. Instead, I recommend following the code callouts to understand the high-lights.


## sync-cond
- theo định nghĩa, `event` là bất kì tính hiệu nào giữa 2 hoặc nhiều `goroutine` mà thực tế nó đã xảy ra. Thông thường bạn sẽ muốn một đợi một tính hiệu trước khi thực hiện một goroutine khác. Nếu chúng ta xem xét để thực hiện điều này mà ko dùng `Cond` thì đơn giản là dùng vòng lặp vô tận.
```go
for conditionTrue() == false { 
	time.Sleep(1*time.Millisecond)
}
```
- Cách giải quyết này cũng tốt, nhưng thực sự là sẽ không hiệu quả, bời vì bạn phải tìm ra thời gian để cấu hình hàm `sleep`. Nếu sleep quá lâu thì sẽ ảnh hưởng đến performance, còn nếu quá ngắn sẽ thì sẽ tốn thời gian của `CPU` một cách không cần thiết. Nó sẽ tốt nếu có một loại function hoặc cách gì đó mà `goroutine` có thể `sleep` cho đến khi có 1 tín hiệu thực thi. Và `Cond` sẽ giúp thực thi điều đó.
### các function trong cond:
- NewCond:
```go
func NewCond(l Locker) *Cond  // Create a new Cond conditional variable.
```
- Broadcast
```go
func (c *Cond) Broadcast() 
// Broadcast will wake up all goroutines waiting for c.
// Broadcast can be called with or without locking.
```
- Signal:
```go
func (c *Cond) Signal() 
// Signal wakes up only 1 goroutine waiting for c.
// Signal can be called with or without locking.
```
- Wait:
```go
func (c *Cond) Wait()
// does not return unless it is woken up by Signal or Broadcast.
```

- Sử dụng `Cond` để viết một ví dụ đơn giản trước:
```go
c := sync.NewCond(&sync.Mutex{}) // (1) new cond instantiate, func NewCond sẽ đáp ứng sync.Locker
c.L.Lock() // (2) Lock process
for conditionTrue() == false {
	c.Wait() // (3)  chúng ta sẽ đợi khi có một thông báo điều kiện được xả ra. Và sẽ block tất cả các goroutine khác.
}
c.L.Unlock() // (4) UnLock process
```
- Cách tiếp cận này là hiệu quả hơn. Ghi chú, func `Wait` không chỉ block, nó còn treo `goroutine` hiện tại, và cho phép các goroutine khác vẫn chạy trên `OS thread`. 
- Để giải thích thêm, xem thêm ví dụ sau:
```go
func CondExample1() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]int, 0, 10)
	removeFromQueue := func(delay time.Duration, i int) {
		time.Sleep(delay)
		c.L.Lock()
		fmt.Println("before remove:", queue)
		queue = queue[1:]
		fmt.Println("after remove:", queue)
		c.L.Unlock()
		c.Signal()
	}
	for i := 0; i < 10; i++ {
		fmt.Println("start loop;", i)
		c.L.Lock()
		for len(queue) == 2 {
			fmt.Println("len  equal 2, waiting", i)
			c.Wait()
		}
		fmt.Println("Adding to queue", i)
		queue = append(queue, i)
		go removeFromQueue(1*time.Second, i)
		c.L.Unlock()
		fmt.Println()
	}
	fmt.Println("after processing, len queue:", len(queue), queue)
}
```
- output:
```
start loop; 0
Adding to queue 0

start loop; 1
Adding to queue 1

start loop; 2
len  equal 2, waiting 2
before remove: [0 1]
after remove: [1]
Adding to queue 2

start loop; 3
len  equal 2, waiting 3
before remove: [1 2]
after remove: [2]
Adding to queue 3

start loop; 4
len  equal 2, waiting 4
before remove: [2 3]
after remove: [3]
Adding to queue 4

start loop; 5
len  equal 2, waiting 5
before remove: [3 4]
after remove: [4]
Adding to queue 5

start loop; 6
len  equal 2, waiting 6
before remove: [4 5]
after remove: [5]
Adding to queue 6

start loop; 7
before remove: [5 6]
after remove: [6]
Adding to queue 7

start loop; 8
len  equal 2, waiting 8
before remove: [6 7]
after remove: [7]
Adding to queue 8

start loop; 9
len  equal 2, waiting 9
before remove: [7 8]
after remove: [8]
Adding to queue 9

after processing, len queue: 2 [8 9]
```
- Như kết quả, thì chương trình đã add 10 item đến queue, nhưng nó luôn luôn đợi cho cho 1 item được `dequeue` trước khi `enqueue` một item khác
- Trong ví dụ có một function `Signal`, nó là một method mà Cond cung cấp để notifying một goroutine đã được block trên wait trước đó.

## Single Flight
- giúp tránh việc thực thi cùng một công việc nhiều lần bởi nhiều goroutine cùng lúc. Khi nhiều goroutine yêu cầu cùng một công việc, singleflight sẽ đảm bảo rằng chỉ có một goroutine thực hiện công việc đó và các goroutine khác sẽ nhận được kết quả từ lần thực hiện đó.
1. Code example:
```go
package main
import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/singleflight"
)
func init() {
	rand.Seed(time.Now().UnixNano())
}
var (
	concurrencyControl = &singleflight.Group{}
	n                  = 10
)
func locking(_ int) (int64, error) {
	now := time.Now()
	randomNumber := rand.Intn(1000) + 1
	time.Sleep(time.Millisecond * time.Duration(randomNumber))
	return time.Since(now).Milliseconds(), nil
}
func main() {
	now := time.Now()
	keyS := "flight"
	for i := 1; i <= n; i++ {
		go func() {
			value, err, s := concurrencyControl.Do(keyS, func() (interface{}, error) {
				return locking(i)
			})
			if err != nil {
				fmt.Printf("Goroutine %d: error: %v\n", i, err)
				return
			}
			fmt.Printf("Goroutine %d: result: %v (shared: %v)\n", i, value, s)

		}()
	}
	time.Sleep(time.Second * 5)
	fmt.Println("Done", time.Since(now).Milliseconds())
}
```
2. output
```console
ducnp@nguyens-MacBook-Pro-4 singleflight % go run main.go
Goroutine 9: result: 556 (shared: true)
Goroutine 8: result: 556 (shared: true)
Goroutine 10: result: 556 (shared: true)
Goroutine 3: result: 556 (shared: true)
Goroutine 6: result: 556 (shared: true)
Goroutine 5: result: 556 (shared: true)
Goroutine 7: result: 556 (shared: true)
Goroutine 2: result: 556 (shared: true)
Goroutine 4: result: 556 (shared: true)
Goroutine 1: result: 556 (shared: true)
Done 5001
```