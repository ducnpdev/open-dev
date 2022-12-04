package main

import (
	"context"
	"fmt"
	"open-dev/usecase/ecdsa"
	"open-dev/usecase/rsa"
	"time"
)

func main() {
	ecdsa.MainECDSA()
	rsa.MainRsa()

	//
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Millisecond)*10)
	defer cancel()

	ch := make(chan bool)
	defer close(ch)

	var (
		errC  error
		value string
	)
	go func() {
		value = handle(ctx, ch)
	}()

	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		fmt.Println("cancelling...")
		return
	case <-ch:
	}
	fmt.Println("errC:", errC)
	fmt.Println("value:", value)

}

func handle(ctx context.Context, ch chan bool) string {
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}
	ch <- true
	return "ok"
}
