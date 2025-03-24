package http

import (
	"context"
	"ebirukov/qbro/internal/config"
	"ebirukov/qbro/internal/model"
	"errors"
	"net/http"
)

type Broker interface {
	Put(ctx context.Context, queueID model.QueueID, msg model.Message) error
	Get(ctx context.Context, queueID model.QueueID) (model.Message, error)
}

type BrokerHandler struct {
	broker Broker
	cfg    config.BrokerCfg
}

func NewBrokerHandler(cfg config.BrokerCfg, broker Broker) *BrokerHandler {
	return &BrokerHandler{
		broker: broker,
		cfg:    cfg,
	}
}

func translateError(err error) int {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
