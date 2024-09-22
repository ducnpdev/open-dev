package main

import "fmt"

func pipelineStep1(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			fmt.Println("pipelineStep1", n)
			out <- n
		}
		close(out)
	}()
	return out
}

func pipelineStep2(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Println("pipelineStep2", n)
			out <- n * 2
		}
		close(out)
	}()
	return out
}

func pipelineStep3(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Println("pipelineStep3", n)
			out <- n + 1
		}
		close(out)
	}()
	return out
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	c1 := pipelineStep1(nums)
	c2 := pipelineStep2(c1)
	c3 := pipelineStep3(c2)
	for result := range c3 {
		fmt.Println(result)
	}
}
