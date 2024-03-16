package processor

import (
	"context"
	"github.com/redis/go-redis/v9"
)

const Books = "books"

type BookStorage interface {
	Save() error
}

type BookProcessor struct {
	streamer StreamReader
	storage  BookStorage
}

func NewBookProcessor(streamer StreamReader, storage BookStorage) *BookProcessor {
	return &BookProcessor{
		streamer: streamer,
		storage:  storage,
	}
}

func (p *BookProcessor) Run(ctx context.Context) error {
	p.streamer.LimitConsume(ctx, Books, p.processMessage)

	return nil
}

func (p *BookProcessor) processMessage(ctx context.Context, m redis.XMessage) error {
	/////// DB
	return nil
}
