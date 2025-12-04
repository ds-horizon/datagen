package main

import (
	"fmt"
	"log/slog"
	"time"
)

// Sink_kafka___datagen_minimal_data loads __datagen_minimal data into Kafka
func Sink_kafka___datagen_minimal_data(modelName string, records []*__datagen_minimal, config *__dgi_KafkaConfig) error {
	slog.Debug(fmt.Sprintf("initializing Kafka producer for %s with %d records", modelName, len(records)))
	if err := Init___datagen_minimal_kafka_producer(config); err != nil {
		return fmt.Errorf("✘ [Kafka] %s: FAILED\n   └─ Messages sent: 0/%d\n   └─ Error: %v\n",
			modelName, len(records), err)
	}
	defer func() {
		err := Close___datagen_minimal_kafka_producer()
		if err != nil {
			slog.Warn(fmt.Sprintf("failed to close Kafka producer for %s: %s", modelName, err.Error()))
		}
	}()

	_, err := Get___datagen_minimal_kafka_producer()
	if err != nil {
		return fmt.Errorf("✘ [Kafka] %s: FAILED\n   └─ Messages sent: 0/%d\n   └─ Error: %v\n",
			modelName, len(records), err)
	}

	batchSize := config.BatchSize
	if batchSize <= 0 {
		batchSize = len(records)
	}

	totalSent := 0

	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		slog.Debug(fmt.Sprintf("sending batch starting at %d of size %d for %s to Kafka", i, len(batch), modelName))
		if err := Load___datagen_minimal_kafka(batch, config); err != nil {
			return fmt.Errorf("✘ [Kafka] %s: FAILED\n   └─ Messages sent: %d/%d\n   └─ Error: %v\n",
				modelName, totalSent, len(records), err)
		}

		totalSent += len(batch)

		if config.Throttle > 0 && end < len(records) {
			throttleDuration := time.Duration(config.Throttle) * time.Millisecond
			slog.Debug(fmt.Sprintf("throttling %s between batches for %s", throttleDuration, modelName))
			time.Sleep(throttleDuration)
		}
	}

	slog.Info(fmt.Sprintf("successfully sent %d/%d messages for %s to Kafka topic '%s'", totalSent, len(records), modelName, config.Topic))
	return nil
}

// Clear_kafka___datagen_minimal_data is a no-op for Kafka as it's an append-only log
func Clear_kafka___datagen_minimal_data(modelName string, config *__dgi_KafkaConfig) error {
	slog.Info(fmt.Sprintf("Kafka clear operation skipped for %s (Kafka is append-only)", modelName))
	// No-op: Kafka doesn't support truncation like SQL databases
	// If you need to clear data, you would need to delete and recreate the topic,
	// which requires admin permissions and is typically not done during normal operations.
	return nil
}
