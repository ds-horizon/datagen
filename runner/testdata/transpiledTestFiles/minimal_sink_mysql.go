package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// Sink_mysql___datagen_minimal_data loads __datagen_minimal data into MySQL
func Sink_mysql___datagen_minimal_data(modelName string, records []*__datagen_minimal, config *MySQLConfig) error {
	slog.Debug(fmt.Sprintf("initializing MySQL connection for %s with %d records", modelName, len(records)))
	if err := Init___datagen_minimal_mysql_connection(config); err != nil {
		return fmt.Errorf("✘ [MySQL] %s: FAILED\n   └─ Rows inserted: 0/%d\n   └─ Error: %v\n",
			modelName, len(records), err)
	}
	defer func() {
		err := Close___datagen_minimal_mysql_connection()
		if err != nil {
			slog.Warn(fmt.Sprintf("failed to close DB connection for %s: %s", modelName, err.Error()))
		}
	}()

	db, err := Get___datagen_minimal_mysql_connection()
	if err != nil {
		return fmt.Errorf("✘ [MySQL] %s: FAILED\n   └─ Rows inserted: 0/%d\n   └─ Error: %v\n",
			modelName, len(records), err)
	}

	batchSize := config.BatchSize
	if batchSize <= 0 {
		batchSize = len(records)
	}

	totalInserted := 0

	slog.Debug(fmt.Sprintf("starting MySQL transaction for %s with batch size %d", modelName, batchSize))
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("✘ [MySQL] %s: FAILED\n   └─ Rows inserted: %d/%d\n   └─ Error: %v\n",
			modelName, totalInserted, len(records), err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				slog.Error(fmt.Sprintf("error rolling back transaction for %s: %s", modelName, err.Error()))
			}
		}
	}()

	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		slog.Debug(fmt.Sprintf("loading batch starting at %d of size %d for %s into MySQL", i, len(batch), modelName))
		if err := Load___datagen_minimal_mysql(batch, tx); err != nil {
			return fmt.Errorf("✘ [MySQL] %s: FAILED\n   └─ Rows inserted: %d/%d\n   └─ Error: %v\n",
				modelName, totalInserted, len(records), err)
		}

		totalInserted += len(batch)

		if config.Throttle != "" && end < len(records) {
			if throttleDuration, err := time.ParseDuration(config.Throttle); err == nil {
				slog.Debug(fmt.Sprintf("throttling %s between batches for %s", throttleDuration, modelName))
				time.Sleep(throttleDuration)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("✘ [MySQL] %s: FAILED\n   └─ Rows inserted: %d/%d\n   └─ Error: %v\n",
			modelName, totalInserted, len(records), err)
	}

	slog.Info(fmt.Sprintf("successfully loaded %d/%d rows for %s into MySQL", totalInserted, len(records), modelName))
	return nil
}

// Clear_mysql___datagen_minimal_data clears __datagen_minimal data from MySQL
func Clear_mysql___datagen_minimal_data(modelName string, config *MySQLConfig) error {
	slog.Debug(fmt.Sprintf("initializing MySQL connection for clearing data for %s", modelName))
	if err := Init___datagen_minimal_mysql_connection(config); err != nil {
		return fmt.Errorf("MySQL connection failed: %w", err)
	}

	defer func() {
		err := Close___datagen_minimal_mysql_connection()
		if err != nil {
			slog.Warn(fmt.Sprintf("failed to close DB connection: %s", err.Error()))
		}
	}()

	db, err := Get___datagen_minimal_mysql_connection()
	if err != nil {
		return fmt.Errorf("failed to get MySQL connection: %w", err)
	}

	slog.Debug(fmt.Sprintf("starting MySQL transaction for clearing data for %s", modelName))
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction for clearing model %s: %w", modelName, err)
	}

	if err := Truncate___datagen_minimal_mysql(tx); err != nil {
		return fmt.Errorf("failed to truncate table for model %s: %w", modelName, err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				slog.Error(fmt.Sprintf("error rolling back transaction for %s: %s", modelName, err.Error()))
			}
		}
	}()

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for clearing model %s: %w", modelName, err)
	}

	slog.Info(fmt.Sprintf("successfully cleared data for %s from MySQL", modelName))
	return nil
}
