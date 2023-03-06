package concurrency

import (
	"fmt"
	"math/rand"
	"time"
)

func MainFanInFanOut() {

}

func faninfanout() {
	rand := func() interface{} { return rand.Intn(50000000) }
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}
