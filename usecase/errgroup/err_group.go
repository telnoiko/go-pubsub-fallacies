package errgroup

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
	rootCtx, cancel := context.WithCancel(context.Background())

	input := make(chan int, u.buffer)
	defer close(input)
	results := make(chan int, u.buffer)
	defer close(results)

	g, ctx := errgroup.WithContext(rootCtx)

	// start consumers
	go u.c.Consume(ctx, results)

	for i := range u.consumerN {
		g.Go(u.t.TransformFunc(ctx, i, input, results))
	}

	// start producer
	// TODO: potential deadlock on blocking write after consumers has been stopped
	go u.p.Produce(rootCtx, u.itemsToProcess, input)

	time.Sleep(time.Second)

	// simulate shutdown
	u.l.Info("cancelling")
	cancel()

	// wait for consumer to shutdown
	err := g.Wait()
	u.l.Info("finished", "error", err)
}
