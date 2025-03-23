package service

import (
	"context"
	"ebirukov/qbro/internal/model"
)

const MaxPoolSize = 100

type InvokeResponse struct {
	value model.Message
	err   error
}

type InvokeContext struct {
	res    chan InvokeResponse
	reqctx context.Context
}

type FifoGetQueue struct {
	pool   chan InvokeContext
	appCtx context.Context
	queue  Queue
}

func NewFifoGetQueue(appCtx context.Context, queueID model.QueueID, queue Queue, poolSize int) *FifoGetQueue {
	pool := make(chan InvokeContext, poolSize)

	context.AfterFunc(appCtx, func() {
		close(pool)
	})

	go func() {
		for invoke := range pool {
			select {
			case <-invoke.reqctx.Done():
				continue
			case <-appCtx.Done():
				return
			default:
			}

			v, err := queue.Get(invoke.reqctx)

			select {
			case <-invoke.reqctx.Done():
				continue
			case <-appCtx.Done():
				return
			default:
			}

			select {
			case invoke.res <- InvokeResponse{value: v, err: err}:
			case <-appCtx.Done():
				return
			}

		}
	}()

	return &FifoGetQueue{
		pool:   pool,
		appCtx: appCtx,
		queue:  queue,
	}
}

func (co *FifoGetQueue) Put(ctx context.Context, msg model.Message) error {
	return co.queue.Put(ctx, msg)
}

func (co *FifoGetQueue) Get(ctx context.Context) (model.Message, error) {
	resp := make(chan InvokeResponse)

	defer close(resp)

	select {
	case <-co.appCtx.Done():
		return nil, co.appCtx.Err()
	default:
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-co.appCtx.Done():
		return nil, co.appCtx.Err()
	case co.pool <- InvokeContext{res: resp, reqctx: ctx}:
	default:
		return nil, ErrTooManyConnection
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-co.appCtx.Done():
		return nil, co.appCtx.Err()
	case res := <-resp:
		return res.value, res.err
	}
}
