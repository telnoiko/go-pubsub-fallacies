package consumer

import (
	"context"
	"log/slog"
)

type consumer struct {
	l *slog.Logger
}

func New(l *slog.Logger) *consumer {
	return &consumer{l: l}
}

func (c *consumer) Consume(_ context.Context, results <-chan int) {
	for i := range results {
		c.l.Info("consumer: consumed", "result", i)
	}
}
