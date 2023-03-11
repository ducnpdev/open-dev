package segmentio

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

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
		// Balancer: &kafka.Hash{},
		// Dialer:   dialer,
	})

	for {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		msgs := kafka.Message{
			Topic: TopicLogging,
			Key:   []byte(strconv.Itoa(i)),
			Value: []byte("this is message" + strconv.Itoa(i)),
		}
		err := w.WriteMessages(ctx, msgs)
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", i)
		i++
		// sleep for a second
		time.Sleep(time.Second)
	}
}
