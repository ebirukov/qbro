package connector

import (
	"context"
	"ebirukov/qbro/internal/model"
	"errors"
)

type ChanQueueConnector struct {
	size int
}

func NewChanQueueConnector(size int) *ChanQueueConnector {
	return &ChanQueueConnector{size: size}
}

func (c *ChanQueueConnector) Connect(ctx context.Context, queueID model.QueueID) (
	<-chan model.Message, func(context.Context, model.QueueID, model.Message) error, error) {
	queue := make(chan model.Message, c.size)
	
	return queue, func(ctx context.Context, queueID model.QueueID, msg model.Message) error {
		select {
		case <- ctx.Done():
			return ctx.Err()
		case queue <- msg:
			return nil
		default:
			return errors.New("queue max size reached")
		}
	}, nil
}