package main

import (
	"errors"
)

type __dgi_PostgresConfig struct {
	Host           string `json:"host"`
	Database       string `json:"database"`
	Port           int    `json:"port,omitempty"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	BatchSize      int    `json:"batch_size,omitempty"`
	SSLMode        string `json:"ssl_mode,omitempty"`
	Timeout        string `json:"timeout,omitempty"`
	Throttle       string `json:"throttle,omitempty"`
}

func (c *__dgi_PostgresConfig) Validate() error {
	if c.Host == "" || c.Database == "" || c.Username == "" || c.Password == "" {
		return errors.New("postgres: host, database, user and password are required")
	}
	return nil
}

