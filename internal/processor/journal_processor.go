package processor

import (
	"context"
	"github.com/redis/go-redis/v9"
)

const Journals = "journals"

type JournalStorage interface {
	Save() error
}

type JournalProcessor struct {
	streamer StreamReader
	storage  JournalStorage
}

func NewJournalProcessor(streamer StreamReader, storage JournalStorage) *JournalProcessor {
	return &JournalProcessor{
		streamer: streamer,
		storage:  storage,
	}
}

func (p *JournalProcessor) Run(ctx context.Context) error {
	p.streamer.LimitConsume(ctx, Journals, p.processMessage)

	return nil
}

func (p *JournalProcessor) processMessage(ctx context.Context, m redis.XMessage) error {

	return nil
}
