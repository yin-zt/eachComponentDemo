package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

var (
	resp any
)

func main() {
	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		resp = err
		panic(resp)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			resp = err
			panic(resp)
		}
	}()

	// Set up a signal handler to gracefully handle SIGINT and SIGTERM signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages from the "my_topic" Kafka topic
	partitionConsumer, err := consumer.ConsumePartition("my_topic", 0, sarama.OffsetNewest)
	if err != nil {
		resp = err
		panic(resp)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			resp = err
			panic(resp)
		}
	}()

	// Continuously consume messages from the partition consumer
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error while consuming message: %v\n", err)
		case <-signals:
			return
		}
	}
}
