package app

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/config"
	"github.com/patyukin/go-redis-streams/internal/streamer"
	"github.com/patyukin/go-redis-streams/internal/streamer/redis"
	"log/slog"
)

type serviceProvider struct {
	cfg      *config.Config
	logger   *slog.Logger
	streamer streamer.StreamerInterface
}

func newServiceProvider(cfg *config.Config, logger *slog.Logger) *serviceProvider {
	return &serviceProvider{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *serviceProvider) Streamer(ctx context.Context) streamer.StreamerInterface {
	if s.streamer == nil {
		s.streamer = redis.NewRedisStreamer(ctx, s.cfg)
	}

	return s.streamer
}
