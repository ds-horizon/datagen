package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

var __datagen_multiple_types_postgres_connection *sql.DB

// Init___datagen_multiple_types_postgres_connection initializes a shared Postgres connection for __datagen_multiple_types.
func Init___datagen_multiple_types_postgres_connection(req *__dgi_PostgresConfig) error {
	if _, err := Get___datagen_multiple_types_postgres_connection(); err == nil {
		return nil
	}

	sslMode := req.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	port := req.Port
	if port == 0 {
		port = 5432
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		req.Host, port, req.Username, req.Password, req.Database, sslMode)

	// Optional timeout
	if d, err := time.ParseDuration(req.Timeout); err == nil && d > 0 {
		dsn += fmt.Sprintf(" connect_timeout=%d", int(d.Seconds()))
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("ping db: %w", err)
	}

	__datagen_multiple_types_postgres_connection = db
	return nil
}

// Get___datagen_multiple_types_postgres_connection returns the shared Postgres DB or an error if not initialized.
func Get___datagen_multiple_types_postgres_connection() (*sql.DB, error) {
	if __datagen_multiple_types_postgres_connection == nil {
		return nil, fmt.Errorf("postgres connection for __datagen_multiple_types is not initialized")
	}
	return __datagen_multiple_types_postgres_connection, nil
}

// Close___datagen_multiple_types_postgres_connection closes the shared Postgres DB for __datagen_multiple_types if initialized.
func Close___datagen_multiple_types_postgres_connection() error {
	if __datagen_multiple_types_postgres_connection == nil {
		slog.Warn(fmt.Sprintf("Attempted to close Postgres connection for %s, but connection was never initialized or already closed", "multiple_types"))
		return nil
	}
	err := __datagen_multiple_types_postgres_connection.Close()
	__datagen_multiple_types_postgres_connection = nil
	return err
}
