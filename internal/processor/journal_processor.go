package processor

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/streamer"
	"github.com/redis/go-redis/v9"
)

type JournalProcessor struct {
	streamer streamer.StreamerInterface
}

func NewJournalProcessor(streamer streamer.StreamerInterface) Processor {
	return &JournalProcessor{
		streamer: streamer,
	}
}

func (p *JournalProcessor) Run(ctx context.Context) error {
	p.streamer.LimitConsume(ctx, "books", p.processMessage)

	return nil
}

func (p *JournalProcessor) processMessage(ctx context.Context, m redis.XMessage) error {

	return nil
}
