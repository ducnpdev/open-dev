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
