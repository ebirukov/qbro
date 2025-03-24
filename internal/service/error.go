package service

import "errors"

var (
	ErrQueueLimitExceeded  = errors.New("queue limit exceeded")
)
