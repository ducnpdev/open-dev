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
	TopicLogging  = "topic-123"
	BrokerAddress = "localhost:9092"
	Group         = "topic-123-group-1"
	workerNumber  = 21
)

func main() {
	s := make(chan bool)
	Consume(context.Background(), TopicLogging, 0)
	<-s
}

func Consume(ctx context.Context, topic string, p1a int) {
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
		Logger:  l,
		Dialer:  dialer,
	})
	var waiting chan bool
	for ii := 1; ii <= workerNumber; ii++ {
		go func() {
			for {
				handleLogic(ctx, r)
			}
		}()
	}

	<-waiting

}

func handleLogic(ctx context.Context, r *kafka.Reader) {
	m, err := r.FetchMessage(ctx)
	fmt.Println("start process handle message")
	defer fmt.Println("end process handle message")
	// sleep handle logic
	time.Sleep(time.Millisecond * 1000)
	if err != nil {
		return
	}

	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	if err := r.CommitMessages(ctx, m); err != nil {
		log.Fatal("failed to commit messages:", err)
	}
}
