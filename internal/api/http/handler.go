package http

import (
	"context"
	"ebirukov/qbro/internal/model"
	"ebirukov/qbro/internal/service"
	"errors"
	"net/http"
)

type Broker interface {
	Put(ctx context.Context, queueID model.QueueID, msg model.Message) error
	Get(ctx context.Context, queueID model.QueueID) (model.Message, error)
}

type BrokerHandler struct {
	broker Broker
	cfg model.Config
}

func NewBrokerHandler(cfg model.Config, broker Broker) *BrokerHandler {
	return &BrokerHandler{
		broker: broker,
		cfg: cfg,
	}
}

func translateError(err error) int {
	switch {
	case errors.Is(err, context.DeadlineExceeded),errors.Is(err, service.ErrQueueNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrTooManyRequest):
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
