package service_test

import (
	"context"
	"ebirukov/qbro/internal/connector"
	"ebirukov/qbro/internal/model"
	"ebirukov/qbro/internal/service"
	"errors"
	"log"
	"testing"
	"time"
)

func Test_Lifecycle(t *testing.T) {
	appCtx, cancelApp := context.WithCancelCause(context.Background())
	config := model.Config{
		QueueLimit:   1,
		MaxQueueSize: 10,
	}

	reg := service.NewQueueConnRegistry(appCtx, config, connector.NewChanQueueConnCreator())

	defer reg.Shutdown()

	qbro := service.NewBrokerSvc(config, reg)

	ctx, cancel := context.WithTimeout(appCtx, time.Second)
	defer cancel()
	if err := qbro.Put(ctx, "test", model.Message("testmessage")); err != nil {
		log.Fatal(err)
	}

	msg, err := qbro.Get(ctx, "test")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(msg)

	go func() {
		time.Sleep(300 * time.Millisecond)
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
