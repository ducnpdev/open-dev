package main

import (
	"fmt"
	"time"
)

func main() {
	// kafka.MainKafka() /
	// Example2()
	// MulChannel()
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
