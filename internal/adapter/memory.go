package adapter

import (
	"context"
	"ebirukov/qbro/internal/model"
	"ebirukov/qbro/internal/service"
	"ebirukov/qbro/pkg/queue/memory"
)


type CreatorAdapter func(ctx context.Context, queueID model.QueueID, size int) (*memory.ChanQueue, error)

func (c CreatorAdapter) Create(ctx context.Context, queueID model.QueueID, size int) (service.Queue, error) {
	return c(ctx, queueID, size)
}

func NewChanQueueCreator() service.QueueCreator {
	return CreatorAdapter(memory.Create)
}