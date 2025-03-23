package model

type Config struct {
	QueueLimit int
	MaxQueueSize int

	DefaultTimeout int

	HttpPort string
}