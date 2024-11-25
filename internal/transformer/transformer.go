package transformer

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type transformer struct {
	l *slog.Logger
}

func New(l *slog.Logger) *transformer {
	return &transformer{l: l}
}

func (t *transformer) Transform(
	ctx context.Context,
	wg *sync.WaitGroup,
	transformerID int,
	input, results chan int,
) error {
	if wg != nil {
		defer wg.Done()
	}

	for {
		select {
		case <-ctx.Done():
			t.l.Info("transformer: context cancelled", "transformer_id", transformerID, "error", ctx.Err())
			return ctx.Err()
		case v, ok := <-input:
			if !ok {
				t.l.Info("transformer: channel closed", "transformer_id", transformerID)
				return fmt.Errorf("transformer: channel closed")
			}
			time.Sleep(time.Second)
			t.l.Info("transformer: consumed", "value", v, "transformer_id", transformerID)
			results <- v * v
		}
	}
}

func (t *transformer) TransformFunc(
	ctx context.Context, id int, input, results chan int,
) func() error {
	return func() error {
		return t.Transform(ctx, nil, id, input, results)
	}
}
