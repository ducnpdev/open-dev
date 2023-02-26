package concurrency

import "fmt"

// MainLeaks goroutine leaks memory
func MainLeaks() {
	Leak()
	Leak()
}

func Leak() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer func() {
				fmt.Println("doWork exited.")
			}()
			defer func() {
				close(completed)
			}()
			for s := range strings {
				// Do something interesting
				fmt.Println(s)
			}
		}()
		fmt.Println("end function doWork")
		return completed
	}
	fmt.Println("dowork")
	doWork(nil)
	// Perhaps more work is done here
	fmt.Println("Done.")
}
