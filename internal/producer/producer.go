package producer

import (
	"context"
	"log/slog"
)

type producer struct {
	l *slog.Logger
}

func New(l *slog.Logger) *producer {
	return &producer{l}
}

func (p *producer) Produce(ctx context.Context, amount int, input chan<- int) {
	for i := range amount {
		select {
		case <-ctx.Done():
			p.l.Info("producer stopped", "last_seen_value", i)
			return
		default:
			input <- i
		}
	}
	p.l.Info("producer: sent", "values", amount)
}

func (p *producer) ProduceFunc(ctx context.Context, amount int, input chan<- int) func() error {
	return func() error {
		p.Produce(ctx, amount, input)
		return nil
	}
}
