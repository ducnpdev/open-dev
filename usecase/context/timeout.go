package main

import (
	"context"
	"fmt"
	"time"
)

func TimeOut() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second)*2)
	defer cancel()
	var (
		value string
	)
	go func() {
		value = handle(ctx)
	}()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		fmt.Println("cancelling...")
		return
	}
	fmt.Println("value:", value)
}

func handle(ctx context.Context) string {
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}
	return "ok"
}
