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
	redis.Set(ctx, "demo-key", "", time.Duration(time.Second)*2)
}
