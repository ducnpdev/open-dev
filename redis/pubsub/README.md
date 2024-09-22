# redis pub-sub 
## install redis use docker:
```sh
docker run redis
```

## code
- code example pub:
```go
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
```

- code example sub:
```go
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
```