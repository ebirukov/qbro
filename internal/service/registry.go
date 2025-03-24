package service

import (
	"context"
	"ebirukov/qbro/internal/config"
	"ebirukov/qbro/internal/model"
	"errors"
	"fmt"
	"sync"
)

var ErrQueueLimitExceeded = errors.New("queue limit exceeded")

// Интерфейс для создания соединений к реализации очереди
type QueueConnCreator interface {
	Connect(ctx context.Context, queueID model.QueueID) (<-chan model.Message, func(context.Context, model.QueueID, model.Message) error, error)
}

type QueueConneсtion struct {
	// канал для чтения из очереди
	getMsgChan <-chan model.Message
	// функция для отправки сообщения в очередь
	putMsgFn func(context.Context, model.QueueID, model.Message) error
}

type QConnRegistry struct {
	mx sync.RWMutex

	store        map[model.QueueID]*QueueConneсtion
	queueLimit   int
	connCreator  QueueConnCreator
}

// NewQueueConnRegistry создает и хранит соединение к очереди
func NewQueueConnRegistry(cfg config.BrokerCfg, creator QueueConnCreator) *QConnRegistry {
	return &QConnRegistry{
		mx:           sync.RWMutex{},
		store:        make(map[model.QueueID]*QueueConneсtion, cfg.QueueLimit),
		queueLimit:   cfg.QueueLimit,
		connCreator:  creator,
	}
}

// GetOrRegister получает или регистрирует соединение для очереди с идентификатором queueID
func (r *QConnRegistry) GetOrRegister(ctx context.Context, queueID model.QueueID) (*QueueConneсtion, error) {
	r.mx.RLock()

	qconn, ok := r.store[queueID]
	if ok {
		r.mx.RUnlock()

		return qconn, nil
	}

	r.mx.RUnlock()

	r.mx.Lock()
	defer r.mx.Unlock()

	qconn, ok = r.store[queueID]
	if ok {
		return qconn, nil
	}

	if len(r.store) >= r.queueLimit {
		return nil, fmt.Errorf("can't connect to queue %s; err: %w", queueID, ErrQueueLimitExceeded)
	}

	var err error

	qconn = &QueueConneсtion{}

	qconn.getMsgChan, qconn.putMsgFn, err = r.connCreator.Connect(ctx, queueID)
	if err != nil {
		return nil, fmt.Errorf("can't connect to queue %s; err: %w", queueID, err)
	}

	r.store[queueID] = qconn

	return qconn, nil
}

func (r *QConnRegistry) Shutdown() error {
	r.mx.Lock()
	defer r.mx.Unlock()

	clear(r.store)

	return nil
}
