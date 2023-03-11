package concurrency

import (
	"fmt"
	"sync"
	"time"
)

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
			time.Sleep(time.Second * 2)
			fmt.Println("len  equal 2, waiting", i)
			c.Wait()
		}
		fmt.Println("Adding to queue", i)
		queue = append(queue, i)
		go removeFromQueue(1*time.Second, i)
		c.L.Unlock()
		fmt.Println()
		fmt.Println()
	}
	fmt.Println("after processing, len queue:", len(queue), queue)
}
