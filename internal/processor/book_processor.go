package processor

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/streamer"
	"github.com/redis/go-redis/v9"
)

type BookProcessor struct {
	streamer streamer.StreamerInterface
}

func NewBookProcessor(streamer streamer.StreamerInterface) Processor {
	return &BookProcessor{
		streamer: streamer,
	}
}

func (p *BookProcessor) Run(ctx context.Context) error {
	p.streamer.LimitConsume(ctx, "books", p.processMessage)

	return nil
}

func (p *BookProcessor) processMessage(ctx context.Context, m redis.XMessage) error {

	return nil
}
