package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var (
	BrokerAddress = "localhost:9092"
	TopicLogging  = "dba4"
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
		Topic:   TopicLogging,
		// assign the logger to the writer
		Logger: l,
		Async:  false,
	})
	var ch chan bool
	for ii := 1; ii <= 10000; ii++ {
		go func() {
			for {
				writeMsg(ctx, w, 1)
				time.Sleep(time.Second * 3)
			}
		}()
	}

	<-ch
}

func writeMsg(ctx context.Context, w *kafka.Writer, i int) {
	key := strconv.Itoa(i)
	value := "timeNow:" + time.Now().Format(time.RFC3339Nano) + " " + "uuid:" + uuid.New().String()
	msgs := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}
	fmt.Println("write message: key:", key, "value %s", value)
	err := w.WriteMessages(ctx, msgs)
	if err != nil {
		panic("could not write message " + err.Error())
	}

}
