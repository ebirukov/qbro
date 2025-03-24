package service

import (
	"context"
	"ebirukov/qbro/internal/config"
	"ebirukov/qbro/internal/model"
	"fmt"
	"sync"
)

type QueueConneсtion struct {
	getMsgChan <-chan model.Message
	putMsgFn   func(context.Context, model.QueueID, model.Message) error
}

type QConnRegistry struct {
	mx sync.RWMutex

	store        map[model.QueueID]*QueueConneсtion
	queueLimit   int
	maxQueueSize int
	connCreator  QueueConnCreator
}

func NewQueueConnRegistry(cfg config.Config, creator QueueConnCreator) *QConnRegistry {
	return &QConnRegistry{
		mx:           sync.RWMutex{},
		store:        make(map[model.QueueID]*QueueConneсtion, cfg.QueueLimit),
		queueLimit:   cfg.QueueLimit,
		maxQueueSize: cfg.MaxQueueSize,
		connCreator:  creator,
	}
}

func (r *QConnRegistry) Create(ctx context.Context, queueID model.QueueID) (*QueueConneсtion, error) {
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

	qconn.getMsgChan, qconn.putMsgFn, err = r.connCreator.Create(ctx, queueID, r.maxQueueSize)
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
