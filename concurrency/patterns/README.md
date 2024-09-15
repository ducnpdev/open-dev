# Concurrency Pattern

## worker pool
- worker pool là 1 khái niệm khá phổ biến trong khi lập trình, còn trong `golang` mục đích để hạn chế số lượng goroutines được tạo ra để phục vụ một mục đích nhất định.
- pattern này cực kì hữu ích trong quá trình sử dụng, nhắm giới hạn số lượng process cần thực hiện, và quản lý tài nguyên.
  - để handle request từ client
  - để process các job backgrount
  - để ...
### Code:
```golang
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numberJobConst  = 4
	numWorkersConst = 2
)

func doWorker(id int, jobs <-chan int, kqs chan<- int, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Printf("Do Worker %d started job %d\n", id, job)
		time.Sleep(time.Second)
		fmt.Printf("Do Worker %d finished job %d\n", id, job)
		kqs <- job * 2
	}
	wg.Done()
}

func main() {
	jobs := make(chan int, numberJobConst)
	kqs := make(chan int, numberJobConst)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkersConst; i++ {
		wg.Add(1)
		go doWorker(i, jobs, kqs, &wg)
	}

	for j := 1; j <= numberJobConst; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(kqs)

	for kq := range kqs {
		fmt.Println("Worker Result:", kq)
	}
}
```
### kết quả
```sonsole
Do Worker 2 started job 1
Do Worker 1 started job 2
Do Worker 1 finished job 2
Do Worker 1 started job 3
Do Worker 2 finished job 1
Do Worker 2 started job 4
Do Worker 1 finished job 3
Do Worker 2 finished job 4
Worker Result: 4
Worker Result: 2
Worker Result: 6
Worker Result: 8
```