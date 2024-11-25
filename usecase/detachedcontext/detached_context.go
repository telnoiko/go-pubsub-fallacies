package detachedcontext

import (
	"context"
	"go-concurrency-story/usecase"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

type useCase struct {
	l *slog.Logger

	consumerN      int
	itemsToProcess int
	buffer         int

	p usecase.Producer
	t usecase.Transformer
	c usecase.Consumer
}

func New(l *slog.Logger, cfg usecase.Config, p usecase.Producer, t usecase.Transformer, c usecase.Consumer) *useCase {
	return &useCase{
		l:              l,
		consumerN:      cfg.ConsumerNumber,
		itemsToProcess: cfg.ItemsToProcess,
		buffer:         cfg.Buffer,
		p:              p,
		t:              t,
		c:              c}
}

func (u *useCase) Process(ctx context.Context) {
	ctx, cancelConsumers := context.WithCancel(context.Background())

	input := make(chan int, u.buffer)
	defer close(input)
	results := make(chan int, u.buffer)
	defer close(results)

	// start consumers
	go u.c.Consume(ctx, results)
	consumerGroup, ctx := errgroup.WithContext(ctx)
	for i := range u.consumerN {
		consumerGroup.Go(u.t.TransformFunc(ctx, i, input, results))
	}

	// detach context, start producer
	producerCtx, cancelProducer := context.WithCancel(context.WithoutCancel(ctx))
	producerGroup, producerCtx := errgroup.WithContext(producerCtx)
	producerGroup.Go(u.p.ProduceFunc(producerCtx, u.itemsToProcess, input))

	time.Sleep(time.Second)

	// simulate shutdown
	// wait till producer is stopped
	u.l.Info("stopping producer")
	cancelProducer()
	_ = producerGroup.Wait()

	u.l.Info("stopping consumers")
	cancelConsumers()
	err := consumerGroup.Wait()

	u.l.Info("finished", "error", err)
}
