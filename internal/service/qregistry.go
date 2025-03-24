package service

import (
	"context"
	"ebirukov/qbro/internal/model"
	"fmt"
	"sync"
)

type QueueConnetion struct {
	getMsgChan <-chan model.Message
	putMsgFn   func(context.Context, model.QueueID, model.Message) error
}

type QConnRegistry struct {
	mx sync.RWMutex

	store        map[model.QueueID]*QueueConnetion
	queueLimit   int
	maxQueueSize int
	connCreator  QueueConnCreator

	appCtx context.Context
}

func NewQueueConnRegistry(appCtx context.Context, cfg model.Config, creator QueueConnCreator) *QConnRegistry {
	return &QConnRegistry{
		mx:           sync.RWMutex{},
		store:        make(map[model.QueueID]*QueueConnetion, cfg.QueueLimit),
		queueLimit:   cfg.QueueLimit,
		maxQueueSize: cfg.MaxQueueSize,
		connCreator:  creator,
		appCtx:       appCtx,
	}
}

func (r *QConnRegistry) Get(queueID model.QueueID) (*QueueConnetion, bool) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	q, ok := r.store[queueID]

	return q, ok
}

func (r *QConnRegistry) Create(ctx context.Context, queueID model.QueueID) (*QueueConnetion, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	qconn, ok := r.store[queueID]
	if ok {
		return qconn, nil
	}

	if r.connCreator == nil {
		return nil, fmt.Errorf("can't create queue %s; err: %w", queueID, ErrUnsupportedOperation)
	}

	if len(r.store) > r.queueLimit {
		return nil, fmt.Errorf("can't create queue %s; err: %w", queueID, ErrQueueLimitExceeded)
	}

	qconn = &QueueConnetion{}

	var err error

	qconn.getMsgChan, qconn.putMsgFn, err = r.connCreator.Create(ctx, queueID, r.maxQueueSize)
	if err != nil {
		return nil, fmt.Errorf("can't create queue %s; err: %w", queueID, err)
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
