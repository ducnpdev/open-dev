package concurrency

import (
	"fmt"
	"time"
)

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

// multiple channel read
func MulChannel() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)
	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func SelectDefault() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}

func WaitingOther() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
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
		// Simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("workCounter %v cycles of work before signalled to stop.\n", workCounter)
}
