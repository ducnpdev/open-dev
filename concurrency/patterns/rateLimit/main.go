package main

import (
	"fmt"
	"time"
)

func main() {
	rate := time.Millisecond * 500
	ticker := time.NewTicker(rate)
	defer ticker.Stop()

	requests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)

	for req := range requests {
		<-ticker.C // Esperar el siguiente tick
		fmt.Println("Processing event:", req)
	}
}
