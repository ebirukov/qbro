package service

import "errors"

var (
	ErrQueueLimitExceeded   = errors.New("queue limit exceeded")
	ErrTooManyConnection    = errors.New("too many connection")
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrQueueNotFound        = errors.New("queue not found")
)
