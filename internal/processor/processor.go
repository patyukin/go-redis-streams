package processor

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type StreamReader interface {
	LimitConsume(ctx context.Context, stream string, processMessage func(ctx context.Context, m redis.XMessage) error)
	Close() error
}
