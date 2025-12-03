package main

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

// Load___datagen_with_metadata_kafka sends a batch of records to the configured Kafka topic.
func Load___datagen_with_metadata_kafka(records []*__datagen_with_metadata, config *__dgi_KafkaConfig) error {
	if len(records) == 0 {
		return nil
	}

	producer, err := Get___datagen_with_metadata_kafka_producer()
	if err != nil {
		return fmt.Errorf("failed to get kafka producer: %w", err)
	}

	for _, record := range records {
		// Serialize the record to JSON (default serialization)
		valueBytes, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to serialize record: %w", err)
		}

		// Prepare the Kafka message
		msg := &sarama.ProducerMessage{
			Topic: config.Topic,
			Value: sarama.ByteEncoder(valueBytes),
		}

		// Set the key if provided in config
		if config.Key != "" {
			msg.Key = sarama.StringEncoder(config.Key)
		}

		// Send the message synchronously
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			return fmt.Errorf("failed to send message to kafka: %w", err)
		}

		// Optional: Log successful send (can be removed or made conditional)
		_ = partition
		_ = offset
	}

	return nil
}

// Truncate___datagen_with_metadata_kafka is a no-op for Kafka as it's an append-only log.
// Kafka topics cannot be truncated in the traditional sense without deleting and recreating them.
func Truncate___datagen_with_metadata_kafka(config *__dgi_KafkaConfig) error {
	// No-op: Kafka doesn't support truncation like SQL databases
	// If you need to clear data, you would need to delete and recreate the topic,
	// which requires admin permissions and is typically not done during normal operations.
	return nil
}
