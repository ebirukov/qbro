package service

import (
	"context"
	"ebirukov/qbro/internal/model"
	"fmt"
	"sync"
)

type QRegistry struct {
	mx sync.RWMutex

	store        map[model.QueueID]*FifoGetQueue
	queueLimit   int
	maxQueueSize int
	creator      QueueCreator

	appCtx context.Context
}

func NewQueueRegistry(appCtx context.Context, cfg model.Config, creator QueueCreator) *QRegistry {
	return &QRegistry{
		mx:           sync.RWMutex{},
		store:        make(map[model.QueueID]*FifoGetQueue, cfg.QueueLimit),
		queueLimit:   cfg.QueueLimit,
		maxQueueSize: cfg.MaxQueueSize,
		creator:      creator,
		appCtx:       appCtx,
	}
}

func (r *QRegistry) Get(queueID model.QueueID) (Queue, bool) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	q, ok := r.store[queueID]

	return q, ok
}

func (r *QRegistry) Create(ctx context.Context, queueID model.QueueID) (Queue, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	q, ok := r.store[queueID]
	if ok {
		return q, nil
	}

	if r.creator == nil {
		return nil, fmt.Errorf("can't create queue %s; err: %w", queueID, ErrUnsupportedOperation)
	}

	if len(r.store) > r.queueLimit {
		return nil, fmt.Errorf("can't create queue %s; err: %w", queueID, ErrQueueLimitExceeded)
	}

	queue, err := r.creator.Create(ctx, queueID, r.maxQueueSize)
	if err != nil {
		return nil, fmt.Errorf("can't create queue %s; err: %w", queueID, err)
	}

	q = NewFifoGetQueue(r.appCtx, queueID, queue, min(r.maxQueueSize, MaxPoolSize))

	r.store[queueID] = q

	return q, nil
}

func (r *QRegistry) Shutdown() error {
	r.mx.Lock()
	defer r.mx.Unlock()

	clear(r.store)

	return nil
}
