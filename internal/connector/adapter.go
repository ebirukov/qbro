package connector

import (
	"context"
	"ebirukov/qbro/internal/model"
)

type QueueAdapter func(ctx context.Context, queueID model.QueueID, size int) (<-chan model.Message, func(context.Context, model.QueueID, model.Message) error, error)

func (c QueueAdapter) Create(ctx context.Context, queueID model.QueueID, size int) (<-chan model.Message, func(context.Context, model.QueueID, model.Message) error, error) {
	return c(ctx, queueID, size)
}