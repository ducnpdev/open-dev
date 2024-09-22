package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Địa chỉ Redis server
	})
	err := client.Ping(context.TODO()).Err()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		msg := uuid.New().String()
		err = client.Publish(context.TODO(), "opendev-channel", msg).Err()
		if err != nil {
			panic(err)
		}
		log.Println("Message published!", msg)
		time.Sleep(time.Second)
	}
}
