package main

import (
	"github.com/Shopify/sarama"
	"log"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("failed to create producer: %s", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("failed to close producer: %s", err)
		}
	}()

	msg := &sarama.ProducerMessage{
		Topic: "my_topic",
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("failed to send message to Kafka: %s", err)
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
