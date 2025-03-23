package service

import (
	"context"
	"ebirukov/qbro/internal/model"
	"ebirukov/qbro/pkg/queue/memory"
	"fmt"

	//"log"
	"math/rand"

	//"math/rand"
	"reflect"
	"runtime"
	"sync"

	"sync/atomic"
	"testing"
	"time"
)

func Test_CallOrder(t *testing.T) {
	size := 1
	queue := &memory.ChanQueue{Buf: make(chan []byte, size)}

	appCtx := context.Background()

	testMsg := model.Message("testMessage1")

	err := queue.Put(context.Background(), testMsg)
	if err != nil {
		t.Fatal(err)
	}

	orderCallQueue := NewFifoGetQueue(appCtx, "testqueue", queue, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)

	defer cancel()

	msg, err := orderCallQueue.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if string(msg) != string(testMsg) {
		t.Errorf("msg = %s want %s", msg, testMsg)
	}

	_, err = orderCallQueue.Get(ctx)
	if err != context.DeadlineExceeded {
		t.Fatal(err)
	}
}

func Test_CallOrderParallel(t *testing.T) {
	runtime.GOMAXPROCS(1)

	var wantMsgs []model.Message
	for i := range 10 {
		wantMsgs = append(wantMsgs, model.Message(fmt.Sprint("msg", i)))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testqueue := &TestQueue{Buf: make(chan []byte, len(wantMsgs))}
	for _, msg := range wantMsgs {
		if err := testqueue.Put(ctx, msg); err != nil {
			t.Error(err)
		}
	}

	queue := NewFifoGetQueue(ctx, "testqueue", testqueue, len(wantMsgs))

	if gotMsgs, err := callInParallel(t, wantMsgs, queue); err != nil || !reflect.DeepEqual(gotMsgs, wantMsgs) {
		if err != nil {
			t.Errorf("got err: %s", err)
			return
		}

		t.Errorf("want=%v, got = %v", wantMsgs, gotMsgs)
	}

}

func callInParallel(t *testing.T, orderedMsg []model.Message, queue Queue) ([]model.Message, error) {
	numG := len(orderedMsg)

	var counter int32

	gotOrderedMsg := make([]model.Message, len(orderedMsg))

	var wg sync.WaitGroup

	sema := make(chan struct{})

	wg.Add(numG)
	for range numG {
		go func() {
			defer wg.Done()

			atomic.AddInt32(&counter, 1)

			if int(atomic.LoadInt32(&counter)) == numG {
				atomic.StoreInt32(&counter, 0)
				close(sema)
			}

			<-sema

			ordIdx := atomic.AddInt32(&counter, 1) - 1
			t.Log(ordIdx, "get")
			msg, err := queue.Get(context.TODO())
			if err != nil {
				t.Error(err)
			}
			t.Log(ordIdx, "got", msg)

			gotOrderedMsg[ordIdx] = msg

		}()
	}

	wg.Wait()

	return gotOrderedMsg, nil
}

type TestQueue struct {
	Buf chan []byte
}

func (q *TestQueue) Put(ctx context.Context, msg model.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.Buf <- msg:
		return nil
	}
}

func (q *TestQueue) Get(ctx context.Context) (model.Message, error) {
	time.Sleep(time.Duration(rand.Intn(2)) * time.Millisecond)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-q.Buf:
		return msg, nil
	}
}
