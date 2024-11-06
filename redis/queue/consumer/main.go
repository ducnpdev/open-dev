package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {

	ctx := context.TODO()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	for {

		result, err := redisClient.BLPop(ctx, 0*time.Second, "payments").Result()

		if err != nil {
			fmt.Println(err.Error())
		} else {
			for item := range result {
				str := result[item]
				fmt.Println(str)
			}
		}
	}
}
