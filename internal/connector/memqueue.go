package connector

import (
	"context"
	"ebirukov/qbro/internal/model"
	"ebirukov/qbro/internal/service"
	"errors"
)

func NewChanQueueConnCreator() service.QueueConnCreator {
	return QueueAdapter(NewQueueConn)
}

func NewQueueConn(ctx context.Context, queueID model.QueueID, size int) (
	<-chan model.Message, func(context.Context, model.QueueID, model.Message) error, error) {
	queue := make(chan model.Message, size)
	
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