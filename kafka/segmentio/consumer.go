package segmentio

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func Consume(ctx context.Context, topic string, pa int) {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	l := log.New(os.Stdout, "kafka reader: ", 0)
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	// mechanism, err := scram.Mechanism(scram.SHA512, "username", "password")
	// if err != nil {
	// 	panic(err)
	// }

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		// SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{BrokerAddress},
		Topic:   topic,
		GroupID: Group,
		// assign the logger to the reader
		Logger: l,
		Dialer: dialer,
		// Partition: pa,
		// GroupTopics: []string{},
		// MinBytes: 10e3, // 10KB
		// MaxBytes: 10e6, // 10MB
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

	// for {
	// 	// the `ReadMessage` method blocks until we receive the next event
	// 	msg, err := r.ReadMessage(ctx)
	// 	// now := time.Now()
	// 	// time.Sleep(time.Second * 3)
	// 	// fmt.Println(time.Since(now).Milliseconds())
	// 	if err != nil {
	// 		panic("could not read message " + err.Error())
	// 	}
	// 	// after receiving the message, log its value
	// 	fmt.Println("received: ", string(msg.Value))
	// }
}
