package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var (
	KafkaHost = []string{"10.9.22.23:9092", "10.9.22.22:9092", "10.9.22.25:9092"} // kafka集群节点
)

func main() {
	kconfig := sarama.NewConfig()

	// WaitForAll waits for all in-sync replicas to commit before responding.
	kconfig.Producer.RequiredAcks = sarama.WaitForAll
	kconfig.Producer.Retry.Max = 5

	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
	kconfig.Producer.Partitioner = sarama.NewRandomPartitioner

	kconfig.Producer.Return.Successes = true
	kconfig.Version = sarama.V0_11_0_2
	client, err := sarama.NewSyncProducer(KafkaHost, kconfig)
	fmt.Println(err)
	msg := &sarama.ProducerMessage{
		Topic: "cmdb-test",
		Value: sarama.StringEncoder("message2"),
	}
	client.SendMessage(msg)
}
