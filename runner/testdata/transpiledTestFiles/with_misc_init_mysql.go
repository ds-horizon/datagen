package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

var __datagen_with_misc_mysql_connection *sql.DB

// Init___datagen_with_misc_mysql_connection initializes a shared MySQL connection for __datagen_with_misc.
func Init___datagen_with_misc_mysql_connection(req *MySQLConfig) error {
	if _, err := Get___datagen_with_misc_mysql_connection(); err == nil {
		return nil
	}

	cfg := mysql.Config{
		User:            req.Username,
		Passwd:          req.Password,
		Net:             "tcp",
		Addr:            fmt.Sprintf("%s:%d", req.Host, req.Port),
		DBName:          req.Database,
		ParseTime:       true,
		MultiStatements: true,
		Params:          map[string]string{"charset": "utf8mb4"},
	}
	// Optional timeouts: accept ms strings; ignore if empty or invalid
	if d, err := time.ParseDuration(req.Timeout); err == nil && d > 0 {
		cfg.Timeout = d
	}
	if d, err := time.ParseDuration(req.WriteTimeout); err == nil && d > 0 {
		cfg.WriteTimeout = d
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("ping db: %w", err)
	}

	__datagen_with_misc_mysql_connection = db
	return nil
}

// Get___datagen_with_misc_mysql_connection returns the shared MySQL DB or an error if not initialized.
func Get___datagen_with_misc_mysql_connection() (*sql.DB, error) {
	if __datagen_with_misc_mysql_connection == nil {
		return nil, fmt.Errorf("mysql connection for __datagen_with_misc is not initialized")
	}
	return __datagen_with_misc_mysql_connection, nil
}

// Close___datagen_with_misc_mysql_connection closes the shared MySQL DB for __datagen_with_misc if initialized.
func Close___datagen_with_misc_mysql_connection() error {
	if __datagen_with_misc_mysql_connection == nil {
		return nil
	}
	err := __datagen_with_misc_mysql_connection.Close()
	__datagen_with_misc_mysql_connection = nil
	return err
}
