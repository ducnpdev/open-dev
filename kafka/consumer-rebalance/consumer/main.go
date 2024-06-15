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
	BrokerAddress = "localhost:9092"
	TopicName     = "consumer-rebalance-3p"
	Group         = "rebalance-1"
)

func main() {
	s := make(chan bool)
	go Consume(context.Background(), TopicName, 0)
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
		// assign the logger to the reader
		Logger: l,
		Dialer: dialer,
	})
	var ch chan bool
	for ii := 1; ii <= 1; ii++ {
		go func() {
			for {
				m, err := r.FetchMessage(ctx)
				if err != nil {
					break
				}
				time.Sleep(time.Second)
				fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
				if err := r.CommitMessages(ctx, m); err != nil {
					log.Fatal("failed to commit messages:", err)
				}
			}
		}()
	}

	<-ch

}
