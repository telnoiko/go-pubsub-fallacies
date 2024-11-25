package usecase

import (
	"context"
	"sync"
)

type (
	Producer interface {
		Produce(ctx context.Context, amount int, input chan<- int)
		ProduceFunc(ctx context.Context, amount int, input chan<- int) func() error
	}
	Transformer interface {
		Transform(
			ctx context.Context,
			wg *sync.WaitGroup,
			transformerID int,
			input, results chan int,
		) error
		TransformFunc(
			ctx context.Context,
			id int,
			input, results chan int,
		) func() error
	}
	Consumer interface {
		Consume(ctx context.Context, results <-chan int)
	}

	Processor interface {
		Process(ctx context.Context)
	}
)

type Config struct {
	ConsumerNumber int // number of consumers
	ItemsToProcess int // amount of work
	Buffer         int // back pressure limit
}
