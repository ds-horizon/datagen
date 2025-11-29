package main

import (
	"errors"
)

type __dgi_KafkaConfig struct {
	Topic               string   `json:"topic"`
	Key                 string   `json:"key,omitempty"`
	IncludeKeyInMessage string   `json:"include_key_in_message,omitempty"`
	BootstrapServers    []string `json:"bootstrap_servers"`
	KeySerializer       string   `json:"key_serializer"`
	ValueSerializer     string   `json:"value_serializer"`
	BatchSize           int      `json:"batch_size,omitempty"`
	Throttle            int      `json:"throttle,omitempty"`
}

func (c *__dgi_KafkaConfig) Validate() error {
	if c.Topic == "" || len(c.BootstrapServers) == 0 {
		return errors.New("kafka: topic and bootstrap_servers are required")
	}
	return nil
}
