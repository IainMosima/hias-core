package workers

import "time"

type ConsumerConfig struct {
	MaxConcurrency    int
	ProcessingTimeout time.Duration
	MaxRetries        int
}

func DefaultConsumerConfig() ConsumerConfig {
	return ConsumerConfig{
		MaxConcurrency:    5,
		ProcessingTimeout: 30 * time.Second,
		MaxRetries:        3,
	}
}
