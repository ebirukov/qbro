package memory

import (
	"context"
	"ebirukov/qbro/internal/model"
	"math/rand"
	"time"
)

type ChanQueue struct {
	Buf chan []byte
}

func Create(ctx context.Context, queueID model.QueueID, size int) (*ChanQueue, error) {
	return &ChanQueue{Buf: make(chan []byte, size)}, nil
}

func (q *ChanQueue) Put(ctx context.Context, msg model.Message) error {
	select {
	case <- ctx.Done():
		return ctx.Err()
	case q.Buf <- msg:
		return nil
	}
}

func (q *ChanQueue) Get(ctx context.Context) (model.Message, error) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case msg := <- q.Buf:
			return msg, nil
	}
}