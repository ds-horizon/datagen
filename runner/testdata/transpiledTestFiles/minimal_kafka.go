package main

import (
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
)

// __dgi_GetFieldAsString returns the value of a field as a string for use as Kafka key
func __datagen_minimal_GetFieldAsString(record *__datagen_minimal, fieldName string) (string, error) {
	switch fieldName {
	case "id":
		return fmt.Sprintf("%v", record.id), nil
	default:
		return "", fmt.Errorf("field '%s' not found in minimal", fieldName)
	}
}

// Load___datagen_minimal_kafka sends a batch of records to the configured Kafka topic.
func Load___datagen_minimal_kafka(records []*__datagen_minimal, config *__dgi_KafkaConfig) error {
	if len(records) == 0 {
		return nil
	}

	producer, err := Get___datagen_minimal_kafka_producer()
	if err != nil {
		return fmt.Errorf("failed to get kafka producer: %w", err)
	}

	for _, record := range records {
		// Serialize the record to JSON (default serialization)
		valueBytes := record.__dgi_Serialise()

		// Prepare the Kafka message
		msg := &sarama.ProducerMessage{
			Topic: config.Topic,
			Value: sarama.ByteEncoder(valueBytes),
		}

		// Set the key if provided in config
		// config.Key represents the field name, we need to extract the value of that field
		if config.Key != "" {
			keyValue, err := __datagen_minimal_GetFieldAsString(record, config.Key)
			if err != nil {
				return fmt.Errorf("failed to extract key field '%s': %w", config.Key, err)
			}
			msg.Key = sarama.StringEncoder(keyValue)
		}

		// Send the message synchronously
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			return fmt.Errorf("failed to send message to kafka: %w", err)
		}

		slog.Debug(fmt.Sprintf("sent kafka message on partition %d at offset %d", partition, offset))
	}

	return nil
}

// Truncate___datagen_minimal_kafka is a no-op for Kafka as it's an append-only log.
func Truncate___datagen_minimal_kafka(config *__dgi_KafkaConfig) error {
	return nil
}
