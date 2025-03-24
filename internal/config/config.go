package config

import (
	"errors"
	"fmt"
	"time"
)

type Config struct {
	QueueLimit int
	MaxQueueSize int

	DefaultTimeout time.Duration

	HttpPort int
}

func (c Config) Validate() error {
	switch {
	case c.QueueLimit == 0:
		return errors.New("required number of queues")
	case c.MaxQueueSize == 0:
		return errors.New("required size of queues")
	default:
		return nil
	}
}

func (c Config) HttpAddr() string {
	if c.HttpPort == 0 {
		return ":8080"
	}

	return fmt.Sprintf(":%d", c.HttpPort)
}