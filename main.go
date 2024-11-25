package main

import (
	"context"
	"go-concurrency-story/internal/consumer"
	"go-concurrency-story/internal/producer"
	"go-concurrency-story/internal/transformer"
	"go-concurrency-story/usecase"
	"go-concurrency-story/usecase/waitgroup"
	"log/slog"
)

const (
	consumerN      = 2   // number of consumers
	itemsToProcess = 100 // amount of work
	buffer         = 10  // back pressure limit
)

func main() {
	l := slog.Default()

	cfg := usecase.Config{
		ConsumerNumber: consumerN,
		ItemsToProcess: itemsToProcess,
		Buffer:         buffer,
	}
	p := producer.New(l)
	t := transformer.New(l)
	c := consumer.New(l)

	processor := waitgroup.New(l, cfg, p, t, c)
	//processor := errgroup.New(l, cfg, p, t, c)
	//processor := detachedcontext.New(l, cfg, p, t, c)

	processor.Process(context.Background())
}
