package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var (
	BrokerAddress = "localhost:9092"
	TopicName     = "consumer-rebalance-3p"
)

func main() {
	Produce(context.Background())
}
func Produce(ctx context.Context) {
	// i := 0

	l := log.New(os.Stdout, "kafka writer: ", 0)
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{BrokerAddress},
		Topic:   TopicName,
		// assign the logger to the writer
		Logger: l,
		Async:  false,
	})
	var ch chan bool
	for ii := 1; ii <= 1000; ii++ {
		go func() {
			for {
				for ii := 1; ii <= 100; ii++ {
					writeMsg(ctx, w, ii)
				}
				// time.Sleep(time.Second * 10)
			}
		}()
	}

	<-ch
}

func writeMsg(ctx context.Context, w *kafka.Writer, i int) {
	key := uuid.New().String()
	value := "timeNow:" + time.Now().Format(time.RFC3339Nano) + " " + "uuid:" + uuid.New().String()
	msgs := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}
	fmt.Println("write message: key:", key, "value:", value)
	err := w.WriteMessages(ctx, msgs)
	if err != nil {
		panic("could not write message " + err.Error())
	}

}
