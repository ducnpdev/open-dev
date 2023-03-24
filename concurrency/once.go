package concurrency

import (
	"fmt"
	"sync"
)

func Example1Once() {
	var count int
	increment := func() {
		count++
	}
	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}
	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

func Example2Once() {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }
	var once sync.Once
	once.Do(increment)
	fmt.Printf("Count: %d\n", count)
	once.Do(decrement)
	fmt.Printf("Count: %d\n", count)
}
