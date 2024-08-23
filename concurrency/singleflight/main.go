package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/singleflight"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	concurrencyControl = &singleflight.Group{}
	n                  = 10
)

func locking(_ int) (int64, error) {
	now := time.Now()
	randomNumber := rand.Intn(1000) + 1
	time.Sleep(time.Millisecond * time.Duration(randomNumber))
	return time.Since(now).Milliseconds(), nil
}
func main() {
	now := time.Now()
	keyS := "flight"
	for i := 1; i <= n; i++ {
		go func() {
			value, err, s := concurrencyControl.Do(keyS, func() (interface{}, error) {
				return locking(i)
			})
			if err != nil {
				fmt.Printf("Goroutine %d: error: %v\n", i, err)
				return
			}
			fmt.Printf("Goroutine %d: result: %v (shared: %v)\n", i, value, s)

		}()
	}
	time.Sleep(time.Second * 5)
	fmt.Println("Done", time.Since(now).Milliseconds())
}
