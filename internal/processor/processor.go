package processor

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Processor interface {
	Run(ctx context.Context) error
	processMessage(ctx context.Context, m redis.XMessage) error
}
