package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	err := client.Ping(context.TODO()).Err()
	if err != nil {
		panic(err)
	}
	// Đăng ký nhận thông báo từ channel "mychannel"
	pubsub := client.Subscribe(context.TODO(), "opendev-channel")

	// Nhận thông báo đầu tiên
	for {
		msg, err := pubsub.ReceiveMessage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received message from channel: %s", msg.Payload)
	}

}
