package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/zroubalik/keda-openshift-examples/kafka/kafka-app/internal/kafkaconsumer"
)

func main() {
	config := kafkaconsumer.NewConsumerConfig()
	log.Printf("Go consumer starting, connecting to Kafka Server: bootstrapServer=%s, topic=%s, group=%s, sasl=%v\n", config.BootstrapServers, config.Topic, config.GroupID, config.SaslEnabled)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	consumerConfig := sarama.NewConfig()

	if config.SaslEnabled {
		consumerConfig.Net.SASL.Enable = true
		consumerConfig.Net.SASL.User = config.SaslUser
		consumerConfig.Net.SASL.Password = config.SaslPassword
		consumerConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext

		tlsConfig := &tls.Config{}
		consumerConfig.Net.TLS.Enable = true
		consumerConfig.Net.TLS.Config = tlsConfig
	}

	cgh := &consumerGroupHandler{
		ready: make(chan bool),
	}
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	consumerGroup, err := sarama.NewConsumerGroup([]string{config.BootstrapServers}, config.GroupID, consumerConfig)
	if err != nil {
		log.Printf("Error creating the Sarama consumer: %v", err)
		os.Exit(1)
	}

	defer func() {
		cancel()
		wg.Wait()
		err = consumerGroup.Close()
		if err != nil {
			log.Printf("Error closing the Sarama consumer: %v", err)
			os.Exit(1)
		}
		log.Printf("Consumer closed")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// this method calls the methods handler on each stage: setup, consume and cleanup
			if err := consumerGroup.Consume(ctx, []string{config.Topic}, cgh); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			cgh.ready = make(chan bool)
		}
	}()

	<-cgh.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	// run till the we terminate
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm
}

// struct defining the handler for the consuming Sarama method
type consumerGroupHandler struct {
	ready chan bool
}

func (cgh *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group handler setup\n")
	close(cgh.ready)
	return nil
}

func (cgh *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group handler cleanup\n")
	return nil
}

func (cgh *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message received: value=%s, topic=%s, partition=%d, offset=%d", string(message.Value), message.Topic, message.Partition, message.Offset)
		session.MarkMessage(message, "")
	}
	return nil
}
