package main

import (
	"context"
	"fmt"
	"open-dev/redis/redisPkg"
	"time"
)

func main() {
	ctx := context.Background()
	fmt.Println("")
	redis := redisPkg.InitRedis()
	fmt.Println(redis)
	err := redis.Set(ctx, "demo-key-123", "", time.Duration(time.Second)*2).Err()
	fmt.Println(err)
}
