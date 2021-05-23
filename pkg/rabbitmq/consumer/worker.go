package consumer

import (
	"context"
	logger "github.com/ereb-or-od/kenobi/pkg/logging/interfaces"
	"sync"

	"github.com/streadway/amqp"
)

type Worker interface {
	Serve(ctx context.Context, h Handler, msgCh <-chan amqp.Delivery)
}

type DefaultWorker struct {
}

func (dw *DefaultWorker) Serve(ctx context.Context, h Handler, msgCh <-chan amqp.Delivery) {
	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				return
			}

			if res := h.Handle(ctx, msg); res != nil {

			}
		case <-ctx.Done():
			return
		}
	}
}

type ParallelWorker struct {
	Num    int
	Logger logger.Logger
}

func NewParallelWorker(num int) *ParallelWorker {
	if num < 1 {
		panic("num workers must be greater than zero")
	}

	return &ParallelWorker{
		Num: num,
	}
}

func (pw *ParallelWorker) Serve(ctx context.Context, h Handler, msgCh <-chan amqp.Delivery) {
	wg := &sync.WaitGroup{}
	for i := 0; i < pw.Num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case msg, ok := <-msgCh:
					if !ok {
						return
					}

					if res := h.Handle(ctx, msg); res != nil {
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()
}
