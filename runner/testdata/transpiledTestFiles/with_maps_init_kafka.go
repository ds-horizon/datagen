package main

import (
	"fmt"

	"github.com/IBM/sarama"
)

var __datagen_with_maps_kafka_producer sarama.SyncProducer

// Init___datagen_with_maps_kafka_producer initializes a shared Kafka producer for __datagen_with_maps.
func Init___datagen_with_maps_kafka_producer(req *__dgi_KafkaConfig) error {
	if _, err := Get___datagen_with_maps_kafka_producer(); err == nil {
		return nil
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Version = sarama.V2_6_0_0

	producer, err := sarama.NewSyncProducer(req.BootstrapServers, config)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	__datagen_with_maps_kafka_producer = producer
	return nil
}

// Get___datagen_with_maps_kafka_producer returns the shared Kafka producer or an error if not initialized.
func Get___datagen_with_maps_kafka_producer() (sarama.SyncProducer, error) {
	if __datagen_with_maps_kafka_producer == nil {
		return nil, fmt.Errorf("kafka producer for __datagen_with_maps is not initialized")
	}
	return __datagen_with_maps_kafka_producer, nil
}

// Close___datagen_with_maps_kafka_producer closes the shared Kafka producer for __datagen_with_maps if initialized.
func Close___datagen_with_maps_kafka_producer() error {
	if __datagen_with_maps_kafka_producer == nil {
		return nil
	}
	err := __datagen_with_maps_kafka_producer.Close()
	__datagen_with_maps_kafka_producer = nil
	return err
}
