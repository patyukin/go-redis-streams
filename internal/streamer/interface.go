package streamer

import (
	"context"
	goredis "github.com/redis/go-redis/v9"
)

type StreamerInterface interface {
	Publish(ctx context.Context) error
	LimitConsume(ctx context.Context, stream string, processMessage func(ctx context.Context, m goredis.XMessage) error)
}
