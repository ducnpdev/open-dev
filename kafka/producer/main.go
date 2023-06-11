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
	TopicLogging  = "logging"
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
