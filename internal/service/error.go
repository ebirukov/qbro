package service

import "errors"

var (
	ErrQueueLimitExceeded   = errors.New("queue limit exceeded")
	ErrTooManyRequest       = errors.New("too many requests")
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrQueueNotFound        = errors.New("queue not found")
)
