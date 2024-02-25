package streamer

import (
	"context"
)

type StreamerInterface interface {
	Publish(ctx context.Context) error
	Consume(ctx context.Context, stream string)
	LimitConsume(ctx context.Context, stream string)
}
