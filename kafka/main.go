package kafka

import (
	"context"
	"open-dev/kafka/segmentio"
)

func MainKafka() {
	s := make(chan bool)
	// go segmentio.Produce(context.Background())
	// time.Sleep(3 * time.Second)
	go segmentio.Consume(context.Background(), segmentio.TopicLogging, 0)

	<-s
}
