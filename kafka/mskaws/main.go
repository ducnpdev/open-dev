package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Replace with your MSK broker addresses and topic name
	brokers := []string{"your-msk-broker-1:9092", "your-msk-broker-2:9092"}
	topic := "your-kafka-topic"
	groupID := "my-consumer-group"

	config := kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		MaxWait:     500 * time.Millisecond,
		MinBytes:    1,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset, // Start consuming from the beginning of the topic
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	// Setup signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ctx := context.Background()

ConsumerLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled. Shutting down...")
			break ConsumerLoop
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error fetching message: %v\n", err)
				continue
			}
			fmt.Printf("Received message: key=%s, value=%s\n", string(msg.Key), string(msg.Value))
			reader.CommitMessages(ctx, msg)
		}
	}

	fmt.Println("Consumer closed. Exiting.")
}
