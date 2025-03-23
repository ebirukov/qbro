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

type QueueCreator interface {
	Create(ctx context.Context, queueID model.QueueID, size int) (Queue, error)
}

type QueueRegistry interface {
	Get(model.QueueID) (Queue, bool)
	Create(context.Context, model.QueueID) (Queue, error)
}

type QBro struct {
	config        model.Config
	queueRegistry QueueRegistry
}

func NewQBro(cfg model.Config, registry QueueRegistry) *QBro {
	return &QBro{
		config:        cfg,
		queueRegistry: registry,
	}
}

// Put добавляет сообщение в очередь с заданным идентификатором или создает очередь, если очередь не существует.
// Возвращает ошибку если очередь переполнена или превышено максимально возможное кол-во очередей.
func (b *QBro) Put(ctx context.Context, queueID model.QueueID, msg model.Message) error {
	queue, err := b.queueRegistry.Create(ctx, queueID)
	if err != nil {
		return fmt.Errorf("can't put message to queue %s; err: %w", queueID, err)
	}

	return queue.Put(ctx, msg)
}

// Get получает сообщение из очереди с заданным идентификатором.
// Возвращает ошибку если очереди не существует, истекло время ожидания на пустой очереди или переполнен буфер запросов.
func (b *QBro) Get(ctx context.Context, queueID model.QueueID) (model.Message, error) {
	queue, ok := b.queueRegistry.Get(queueID)
	if !ok {
		return nil, fmt.Errorf("can't get message from queue %s; err: %w", queueID, ErrQueueNotFound)
	}

	return queue.Get(ctx)
}
