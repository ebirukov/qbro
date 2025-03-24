package main

import (
	broker "ebirukov/qbro/internal/app"
	"ebirukov/qbro/internal/config"
	"flag"
	"log"
)

func main() {
	httpPort := flag.Int("http_port", 0, "http server port")
	queueLimit := flag.Int("queues", 0, "maximum number of queues")
	maxQueueSize := flag.Int("queue_size", 0, "maximum queue size")
	defaultTimeout := flag.Duration("read_timeout", 0, "timeout for waiting data from queue")

	flag.Parse()

	cfg := config.Config{
		QueueLimit: *queueLimit,
		MaxQueueSize: *maxQueueSize,
		DefaultTimeout: *defaultTimeout,
		HttpPort: *httpPort,
	}
	
	qbro, err := broker.NewBrokerApp(cfg)
	if err != nil {
		log.Fatalf("can't create app: %s", err)
	}

	if err := qbro.Run(); err != nil {
		log.Fatalf("stop app with err: %s", err)
	}
}