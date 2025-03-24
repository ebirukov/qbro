package service

import (
	"context"
	"ebirukov/qbro/internal/config"
	"ebirukov/qbro/internal/model"
	"fmt"
)

type QueueConnCreator interface {
	Create(ctx context.Context, queueID model.QueueID, size int) (<-chan model.Message, func(context.Context, model.QueueID, model.Message) error, error)
}

type QueueRegistry interface {
	Create(context.Context, model.QueueID) (*QueueConneсtion, error)
}

type BrokerSvc struct {
	config        config.Config
	queueRegistry QueueRegistry
	appCtx        context.Context
}

func NewBrokerSvc(appCtx context.Context, cfg config.Config, registry QueueRegistry) *BrokerSvc {
	return &BrokerSvc{
		config:        cfg,
		queueRegistry: registry,
		appCtx:        appCtx,
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
	qconn, err := b.queueRegistry.Create(ctx, queueID)
	if err != nil {
		return nil, fmt.Errorf("can't get message from queue %s; err: %w", queueID, err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-b.appCtx.Done():
		return nil, b.appCtx.Err()
	case msg := <-qconn.getMsgChan:
		return msg, nil
	}
}
