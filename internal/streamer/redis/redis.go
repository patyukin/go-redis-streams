package redis

import (
	"context"
	"errors"
	"github.com/patyukin/go-redis-streams/internal/config"
	"github.com/patyukin/go-redis-streams/internal/streamer"
	"github.com/patyukin/go-redis-streams/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"math/rand"
)

type Streamer struct {
	c *redis.Client
}

func NewRedisStreamer(ctx context.Context, cfg *config.Config) streamer.StreamerInterface {
	client := redis.NewClient(&redis.Options{
		Addr:           cfg.Redis.DNS,
		PoolSize:       10,
		MaxActiveConns: 100,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis", err)
	}

	return &Streamer{
		c: client,
	}
}

func (s *Streamer) LimitConsume(ctx context.Context, stream string, processMessage func(ctx context.Context, m redis.XMessage) error) {
	workerPool := make(chan struct{}, 6)

	cursor := "0"

	for {
		logger.GetLogger(ctx).Debug("Consume", zap.String("stream", stream), zap.String("cursor", cursor))
		result, err := s.c.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, cursor},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Fatalf("Failed to read from stream: %v\n", err)
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
