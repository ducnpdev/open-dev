# kafka golang example
- install kafka-go
```go
go get https://github.com/segmentio/kafka-go
```
## producer
- create file main.go
- copy and paste:
```go
package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var (
	BrokerAddress = "localhost:9092"
	TopicLogging  = "cmak"
)

func main() {
	Produce(context.Background())
}
func Produce(ctx context.Context) {
	i := 0

	l := log.New(os.Stdout, "kafka writer: ", 0)
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{BrokerAddress},
		Topic:   TopicLogging,
		// assign the logger to the writer
		Logger: l,
		Async:  false,
	})

	for {
		writeMsg(ctx, w, i)
		time.Sleep(time.Millisecond)
	}
}

func writeMsg(ctx context.Context, w *kafka.Writer, i int) {
	msgs := kafka.Message{
		Key:   []byte(strconv.Itoa(i)),
		Value: []byte("timeNow:" + time.Now().Format(time.RFC3339Nano) + " " + "uuid:" + uuid.New().String()),
	}
	err := w.WriteMessages(ctx, msgs)
	if err != nil {
		panic("could not write message " + err.Error())
	}

}
```
## consumer
- create file main.go
- copy and paste:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "topic"
	TopicLogging  = "cmak"
	Topic1        = "topic1"
	Topic2        = "topic2"
	BrokerAddress = "localhost:9092"
	Group         = "cmak-consumer-group-1"
)

func MainKafka() {
	s := make(chan bool)
	go Consume(context.Background(), TopicLogging, 0)
	<-s
}

func Consume(ctx context.Context, topic string, pa int) {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	l := log.New(os.Stdout, "kafka reader: ", 0)
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{BrokerAddress},
		Topic:   topic,
		GroupID: Group,
		// assign the logger to the reader
		Logger: l,
		Dialer: dialer,
	})
	for {

		m, err := r.FetchMessage(ctx)

		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}
```