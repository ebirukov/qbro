package main

import (
	broker "ebirukov/qbro/internal/app"
	"ebirukov/qbro/internal/model"
)

func main() {

	cfg := model.Config{
		QueueLimit: 1,
		MaxQueueSize: 10,
		DefaultTimeout: 30,
		HttpPort: ":8080",
	}
	
	qbro := broker.NewBrokerApp(cfg)

	qbro.Run()
}