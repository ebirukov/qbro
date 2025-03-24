package service

import (
	"context"
	"ebirukov/qbro/internal/model"
	"fmt"
)

type Queue interface {
	Put(ctx context.Context, msg model.Message) error
	Get(ctx context.Context) (model.Message, error)
}

type QueueConnCreator interface {
	Create(ctx context.Context, queueID model.QueueID, size int) (<-chan model.Message, func(context.Context, model.QueueID, model.Message) error, error)
}

type QueueRegistry interface {
	Get(model.QueueID) (*QueueConnetion, bool)
	Create(context.Context, model.QueueID) (*QueueConnetion, error)
}

type BrokerSvc struct {
	config        model.Config
	queueRegistry QueueRegistry
}

func NewBrokerSvc(cfg model.Config, registry QueueRegistry) *BrokerSvc {
	return &BrokerSvc{
		config:        cfg,
		queueRegistry: registry,
	}
}

// Put добавляет сообщение в очередь с заданным идентификатором или создает очередь, если очередь не существует.
// Возвращает ошибку если очередь переполнена или превышено максимально возможное кол-во очередей.
func (b *BrokerSvc) Put(ctx context.Context, queueID model.QueueID, msg model.Message) error {
	qconn, err := b.queueRegistry.Create(ctx, queueID)
	if err != nil {
		return fmt.Errorf("can't put message to queue %s; err: %w", queueID, err)
	}

	return qconn.putMsgFn(ctx, queueID, msg)
}

// Get получает сообщение из очереди с заданным идентификатором.
// Возвращает ошибку если очереди не существует, истекло время ожидания на пустой очереди или переполнен буфер запросов.
func (b *BrokerSvc) Get(ctx context.Context, queueID model.QueueID) (model.Message, error) {
	qconn, ok := b.queueRegistry.Get(queueID)
	if !ok {
		return nil, fmt.Errorf("can't get message from queue %s; err: %w", queueID, ErrQueueNotFound)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-qconn.getMsgChan:
		return msg, nil
	}
}
