package waitgroup

import (
	"context"
	"go-concurrency-story/usecase"
	"log/slog"
	"sync"
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
	ctx, cancel := context.WithCancel(ctx)

	input := make(chan int, u.buffer)
	defer close(input)
	results := make(chan int, u.buffer)
	defer close(results)

	// start consumers in reversed order
	go u.c.Consume(ctx, results)

	var wg sync.WaitGroup
	for i := range u.consumerN {
		wg.Add(1)
		// ignoring returned error intentionally
		// TODO: no cancellation on error from consumer
		go u.t.Transform(ctx, &wg, i, input, results)
	}

	// start producer
	// TODO: potential deadlock on blocking write after consumers has been stopped
	go u.p.Produce(ctx, u.itemsToProcess, input)

	time.Sleep(time.Second)

	// simulate shutdown
	u.l.Info("cancelling")
	cancel()

	// wait for consumer to shut down
	wg.Wait()
	u.l.Info("finished")
}
