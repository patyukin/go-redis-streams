package redis

import (
	"context"
	"fmt"
	"github.com/patyukin/go-redis-streams/internal/config"
	"github.com/patyukin/go-redis-streams/internal/streamer"
	"github.com/patyukin/go-redis-streams/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
)

type Streamer struct {
	c *redis.Client
}

func NewRedisStreamer(ctx context.Context, cfg *config.Config) (streamer.StreamerInterface, error) {
	wrappedCtx := logger.WithLoggerFields(ctx, zap.String("func", "NewRedisStreamer"))
	client := redis.NewClient(&redis.Options{
		Addr:           cfg.Redis.DNS,
		PoolSize:       10,
		MaxActiveConns: 100,
	})

	_, err := client.Ping(wrappedCtx).Result()
	if err != nil {
		logger.GetLogger(wrappedCtx).Error("Failed to ping redis", zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.GetLogger(wrappedCtx).Info("Connected to redis", zap.String("address", cfg.Redis.DNS))
	return &Streamer{c: client}, nil
}

func (s *Streamer) LimitConsume(ctx context.Context, stream string, processMessage func(ctx context.Context, m redis.XMessage) error) {
	wrappedCtx := logger.WithLoggerFields(ctx, zap.String("LimitConsume with stream:", stream))
	workerPool := make(chan struct{}, 6)

	cursor := "0"

	for {
		logger.GetLogger(wrappedCtx).Debug("Consume", zap.String("stream", stream), zap.String("cursor", cursor))
		result, err := s.c.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, cursor},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil {
			logger.GetLogger(wrappedCtx).Error("Failed to XRead", zap.String("error", err.Error()))
			return
		}

		for _, xStream := range result {
			for _, message := range xStream.Messages {

				workerPool <- struct{}{}
				go func(ctx context.Context, m redis.XMessage) {

					defer func() {
						<-workerPool
					}()

					err = processMessage(ctx, m)
					if err != nil {
						// TODO - обработка ошибки
						return
					}
				}(ctx, message)
			}
		}

		if len(result) > 0 {
			cursor = result[0].Messages[len(result[0].Messages)-1].ID
		}
	}
}

func (s *Streamer) Publish(ctx context.Context, stream string) error {
	err := s.c.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			"whatHappened": stream + " received",
			"ticketID":     rand.Intn(100000000),
			"ticketData":   "some " + stream + " data",
		},
	}).Err()

	return err
}
