package service

import (
	"context"
	"ebirukov/qbro/internal/model"
	"ebirukov/qbro/pkg/queue/memory"
	"errors"
	"log"
	"testing"
	"time"
)

type CreateFn func(ctx context.Context, queueID model.QueueID, size int) (*memory.ChanQueue, error)

func (c CreateFn) Create(ctx context.Context, queueID model.QueueID, size int) (Queue, error) {
	return c(ctx, queueID, size)
}


func Test_Lifecycle(t *testing.T) {
	appCtx, cancelApp := context.WithCancelCause(context.Background())
	config := model.Config{
		QueueLimit: 1,
		MaxQueueSize: 10,
	}

	reg := NewQueueRegistry(appCtx, config, CreateFn(memory.Create))

	defer reg.Shutdown()

	qbro := NewBrokerSvc(config, reg)

	ctx, cancel := context.WithTimeout(appCtx, time.Second)
	defer cancel()
	if err := qbro.Put(ctx, "test", model.Message("message")); err != nil {
		log.Fatal(err)
	}

	msg, err := qbro.Get(ctx, "test")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(msg)

	go func() {
		time.Sleep(300*time.Millisecond)
		cancelApp(errors.New("shutdown"))
	}()

	_, err = qbro.Get(ctx, "test")
	if err != nil {
		t.Logf("%v\n", err)
	}

	_, err = qbro.Get(ctx, "test")
	if err != nil {
		t.Error(err)
	}
}