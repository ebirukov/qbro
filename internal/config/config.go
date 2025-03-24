package config

import (
	"errors"
	"fmt"
	"time"
)

type BrokerCfg struct {
	QueueLimit   int
	MaxQueueSize int

	DefaultTimeout time.Duration

	HttpPort int
}

func (c BrokerCfg) Validate() error {
	switch {
	case c.QueueLimit == 0:
		return errors.New("required number of queues")
	case c.MaxQueueSize == 0:
		return errors.New("required size of queues")
	default:
		return nil
	}
}

func (c BrokerCfg) HttpAddr() string {
	return fmt.Sprintf(":%d", c.HttpPort)
}
